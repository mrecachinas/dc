package api

import (
	"github.com/mrecachinas/dcserver/internal/app/config"
	"github.com/mrecachinas/dcserver/internal/util"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Api is a wrapper around various state including
// the MongoDB connection, the AMQP connection and channel,
// and the passed CLI parameters.
type Api struct {
	DB          *mongo.Client
	AMQPClient  *amqp.Connection
	AMQPChannel *amqp.Channel
	Cfg         config.Config
}

// TODO: Status and Task can probably be combined into just Task or status
// Status is the primary data structure for DC. It includes information
// such as start time, stop time, etc.
type Status struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StartTime util.JSONTime      `json:"start_time" bson:"start_time"`
	StopTime  util.JSONTime      `json:"stop_time,omitempty" bson:"stop_time"`
	StopFlag  bool               `json:"stop_flag" bson:"stop_flag"`
}

type Task struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StartTime util.JSONTime      `json:"start_time" bson:"start_time"`
	StopTime  util.JSONTime      `json:"stop_time,omitempty" bson:"stop_time"`
}

// Response is a fairly generic response struct to handle common responses,
// for example, responses that include a message, specific ID, or result
type Response struct {
	Msg    string `json:"msg,omitempty"`
	Result string `json:"result,omitempty"`
	Id     string `json:"id,omitempty"`
}
