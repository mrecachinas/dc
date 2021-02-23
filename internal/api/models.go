package api

import (
	"net/http"
	"sync"

	"github.com/mrecachinas/dcserver/internal/config"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"
)

// Api is a wrapper around various state including
// the MongoDB connection, the AMQP connection and channel,
// and the passed CLI parameters.
type Api struct {
	MongoClient *mongo.Client
	DB          DB
	AMQPClient  *amqp.Connection
	AMQPChannel *amqp.Channel
	HTTPClient  *http.Client
	Websocket   *WebsocketConnectionPool
	Cfg         *config.Config
}

// WebsocketConnectionPool holds a map of every websocket
// connection, so we can broadcast updates to everyone,
// thus requiring only one pull from the database.
type WebsocketConnectionPool struct {
	sync.RWMutex
	Connections map[*websocket.Conn]struct{}
}

// TODO: Status and Task can probably be combined into just Task or status
// Status is the primary data structure for DC. It includes information
// such as start time, stop time, etc.
type Status struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StartTime primitive.DateTime `json:"start_time" bson:"start_time"`
	StopTime  primitive.DateTime `json:"stop_time,omitempty" bson:"stop_time"`
	StopFlag  bool               `json:"stop_flag" bson:"stop_flag"`
}

type Task struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StartTime primitive.DateTime `json:"start_time" bson:"start_time"`
	StopTime  primitive.DateTime `json:"stop_time,omitempty" bson:"stop_time,omitempty"`
}

// Response is a fairly generic response struct to handle common responses,
// for example, responses that include a message, specific ID, or result
type Response struct {
	Msg    string `json:"msg,omitempty"`
	Result string `json:"result,omitempty"`
	Id     string `json:"id,omitempty"`
}
