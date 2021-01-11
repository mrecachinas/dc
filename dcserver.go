package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"net/http"
)

type JSONTime time.Time

func (j *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05.000000", s)
	if err != nil {
		return err
	}
	*j = JSONTime(t)
	return nil
}

func (j JSONTime) MarshalJSON() ([]byte, error) {
	stamp := time.Time(j).Format("2006-01-02T15:04:05.000000")
	return []byte(stamp), nil
}

type Status struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	StartTime JSONTime           `json:"start_time" bson:"start_time"`
	StopTime  JSONTime           `json:"stop_time,omitempty" bson:"stop_time"`
	StopFlag  bool               `json:"stop_flag" bson:"stop_flag"`
}

type Task struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	StartTime JSONTime           `json:"start_time" bson:"start_time"`
	StopTime  JSONTime           `json:"stop_time,omitempty" bson:"stop_time"`
}

type Config struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Debug             bool   `json:"debug"`
	MongoURL          string `json:"mongo_url"`
	MongoPort         int    `json:"mongo_port"`
	MongoDatabaseName string `json:"mongo_dbname"`
	AMQPURL           string `json:"amqp_url"`
	AMQPPort          int    `json:"amqp_port"`
}

type Handler struct {
	DB          *mongo.Client
	AMQPClient  *amqp.Connection
	AMQPChannel *amqp.Channel
	Cfg         Config
}

func main() {
	cfg := setupConfig()

	handler, err := NewDCHandler(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer handler.DB.Disconnect()
	defer handler.AMQPClient.Close()
	defer handler.AMQPChannel.Close()

	e := SetupEchoServer(cfg, handler)

	// Run server
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	go func() {
		if err := e.Start(address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func NewDCHandler(cfg Config) (*Handler, error) {
	client, mongoConnectionError := SetupMongoConnection(cfg)
	if mongoConnectionError != nil {
		return nil, mongoConnectionError
	}

	amqpChannel, amqpError := SetupAMQP(cfg)
	if amqpError != nil {
		return nil, amqpError
	}
	handler := &Handler{
		DB:          client,
		AMQPChannel: amqpChannel,
		Cfg:         cfg,
	}
	return handler, nil
}

func (handler *Handler) GetStatus(c echo.Context) error {
	collection := handler.DB.Database(handler.Cfg.MongoDatabaseName).Collection("tasks")
	var statusList []Status
	err := collection.Find(ctx.TODO(), bson.D).Decode(&statusList)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Hello, Status!")
}

func (handler *Handler) RetrieveTasks() ([]Task, error) {
	resp, err := http.Get(cfg.TaskURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(body, tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (handler *Handler) GetTasks(c echo.Context) error {
	tasks, err := handler.RetrieveTasks()
	if err != nil {
		c.Logger().Error(err)
	}
	return c.JSON(http.StatusOK, &tasks)
}

const (
	TASK_START = "start"
	TASK_STOP  = "stop"
)

type Response struct {
	msg string
}

func (handler *Handler) StartTask(c echo.Context) error {
	err := handler.AMQPChannel.Publish(
		handler.Cfg.AMQPExchange,
		TASK_START,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        task,
		},
	)
	if err != nil {
		c.Logger().Error(err)
	}
	return c.JSON(http.StatusOK, Response{msg: "Successfully submitted start task request"})
}

func (handler *Handler) StopTask(c echo.Context) error {
	err := handler.AMQPChannel.Publish(
		handler.Cfg.AMQPExchange,
		TASK_STOP,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        task,
		},
	)
	if err != nil {
		c.Logger().Error(err)
	}
	return c.JSON(http.StatusOK, Response{msg: "Successfully submitted stop task request"})
}

var (
	upgrader = websocket.Upgrader{}
)

func (handler *Handler) HelloWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Query MongoDB for new records, keeping
		// track of the last ID

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

func SetupMongoConnection(cfg Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SetupEchoServer(cfg Config, handler *Handler) *echo.Echo {
	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// If we're in debug mode, allow CORS
	if cfg.Debug {
		e.Use(middleware.CORS())
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	// Setup routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/api/status", handler.GetStatus)
	e.GET("/api/tasks", handler.GetTasks)
	e.POST("/api/tasks/start", handler.StartTask)
	e.POST("/api/tasks/stop", handler.StopTask)
	e.POST("/api/tasks/save", handler.SaveTask)
	e.GET("/ws", handler.HelloWebsocket)
	e.File("/", "public/index.html")
	return e
}

func setupConfig() Config {
	var cfg Config
	flag.StringVarP(&cfg.Host, "host", "i", "localhost", "The URL to listen on")
	flag.IntVarP(&cfg.Port, "port", "p", 1337, "The port to run on")
	flag.BoolVarP(&cfg.Debug, "debug", "d", false, "Whether or not to enable debug logging")
	flag.StringVar(&cfg.MongoURL, "mongo-url", "localhost", "The URL MongoDB is running on")
	flag.IntVar(&cfg.MongoPort, "mongo-port", 27017, "The port MongoDB is running on")
	flag.StringVar(&cfg.MongoDatabaseName, "mongo-dbname", "dc", "Name of MongoDB database to use")
	flag.StringVar(&cfg.AMQPURL, "amqp-url", "localhost", "The URL RabbitMQ is running on")
	flag.IntVar(&cfg.AMQPPort, "amqp-port", 5672, "The port RabbitMQ is running on")
	flag.Parse()
	return cfg
}
