package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// GetStatus returns a single status object
// to the client given an id.
func (a *Api) GetStatus(c echo.Context) error {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	collection := a.DB.Database(a.Cfg.MongoDatabaseName).Collection("tasks")
	var status Status
	err = collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&status)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, status)
}

// GetAllStatus queries the tasks collection
// for every record and returns it as JSON.
func (a *Api) GetAllStatus(c echo.Context) error {
	statusList, err := GetAllStatus(a.DB.Database(a.Cfg.MongoDatabaseName))
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, statusList)
}

// GetTasks hits the external API for the tasks list XML
// and returns it serialized as JSON.
func (a *Api) GetTasks(c echo.Context) error {
	tasks, err := QueryExternal(a.Cfg.TaskURL, a.HTTPClient)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, &tasks)
}

// CreateTask inserts the client-requested task
// into the tasks collection and pushes a message
// onto the RabbitMQ exchange for creation.
func (a *Api) CreateTask(c echo.Context) error {
	// Deserialize the task JSON request into a `Task` object
	var task Task
	err := json.NewDecoder(c.Request().Body).Decode(&task)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	collection := a.DB.Database(a.Cfg.MongoDatabaseName).Collection("tasks")
	insertResult, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	oid, ok := insertResult.InsertedID.(primitive.ObjectID)
	if ok {
		task.Id = oid

		// Serialize the task object augmented with the new
		// ObjectId and push it onto RabbitMQ
		taskJson, err := json.Marshal(task)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		err = a.AMQPChannel.Publish(
			a.Cfg.AMQPOutputExchange,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/json",
				Body:        taskJson,
			},
		)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, Response{
			Msg: "Successfully submitted start task request",
			Id:  oid.Hex(),
		})
	} else {
		return c.JSON(http.StatusInternalServerError, Response{
			Msg: "Error parsing inserted MongoDB ObjectId",
		})
	}
}

// StopTask sets the `stop_flag` field in the requested task
// in the tasks collection to `true`, which indicates it
// has been requested to stop.
func (a *Api) StopTask(c echo.Context) error {
	id := c.Param("id")
	err := StopTask(a.DB.Database(a.Cfg.MongoDatabaseName), id)
	if err != nil {
		msg := fmt.Sprintf("An error occurred when trying to delete Task %s", id)
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, Response{Msg: msg})
	}
	return c.JSON(http.StatusOK, Response{Msg: "Successfully submitted stop task request"})
}

var (
	upgrader = websocket.Upgrader{}
)

// UpdaterWebsocket handles broadcasting updates to the database;
// Note: this just means querying the collection every
// `api.Cfg.PollingInterval` seconds and returning the *entire*
// contents of the collection.
func (a *Api) UpdaterWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer ws.Close()

	// Create a new connection entry in our connection map
	a.Websocket.Lock()
	a.Websocket.Connections[ws] = struct{}{}
	defer a.Websocket.CloseWebsocketConnection(ws)
	a.Websocket.Unlock()

	msg := fmt.Sprintf(
		"Client %s joined. %d total connections.",
		c.Request().Host,
		len(a.Websocket.Connections),
	)
	c.Logger().Info(msg)

	delay := time.Duration(a.Cfg.PollingInterval) * time.Second

	// TODO: Make error handling/logging more useful? Break if error occurs?
	for {
		select {
		case <-time.After(delay):
			// Get ALL records every N seconds
			statusList, err := GetAllStatus(a.DB.Database(a.Cfg.MongoDatabaseName))
			if err != nil {
				c.Logger().Error("Error getting all status from database")
			}

			// Broadcast result to all connections
			err = a.Websocket.SendMessageToPool(statusList)
			if err != nil {
				c.Logger().Error("Error sending message to pool")
			}

			// Read inbound websocket messages;
			// if we get an error back, that means the client
			// closed the connection
			_, _, err = ws.ReadMessage()
			if err != nil {
				msg := fmt.Sprintf(
					"Client %s quitting. %d connections remain.",
					c.Request().Host,
					len(a.Websocket.Connections),
				)
				c.Logger().Info(msg)
				break
			}
		}
	}
}
