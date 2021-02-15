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
)

// GetAllStatus queries the tasks collection
// for every record and returns it as JSON.
func (a *Api) GetAllStatus(c echo.Context) error {
	collection := a.DB.Database(a.Cfg.MongoDatabaseName).Collection("tasks")
	var statusList []Status
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = cursor.All(context.TODO(), &statusList)
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
		return err
	}

	// TODO: Ask Charlie if I can just write to DB and send ID/start request
	//       with a field {"status": "submitted"} in the DB
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

	collection := a.DB.Database(a.Cfg.MongoDatabaseName).Collection("tasks")
	updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"stop_flag": true})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if updateResult.ModifiedCount != 1 {
		errorMsg := fmt.Sprintf("Task %s was not successfully updated", id)
		c.Logger().Error(errorMsg)
		return c.JSON(http.StatusInternalServerError, Response{Msg: errorMsg})
	}
	return c.JSON(http.StatusOK, Response{Msg: "Successfully submitted stop task request"})
}

var (
	upgrader = websocket.Upgrader{}
)

// HelloWebsocket
func (a *Api) HelloWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Get ALL records every N seconds

		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
