package config

import (
	"github.com/spf13/pflag"
)

// Config is a struct that encapsulates CLI-parsed arguments.
// Note that it is also JSON serializable and deserializable.
type Config struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	Debug              bool   `json:"debug"`
	MongoHost          string `json:"mongo_host"`
	MongoPort          int    `json:"mongo_port"`
	MongoDatabaseName  string `json:"mongo_dbname"`
	AMQPHost           string `json:"amqp_host"`
	AMQPPort           int    `json:"amqp_port"`
	AMQPUser           string `json:"amqp_user"`
	AMQPPassword       string `json:"amqp_password"`
	AMQPOutputExchange string `json:"amqp_output_exchange"`
	TaskURL            string `json:"task_url"`
	ClientCertFile     string `json:"client_certfile"`
	ClientKeyFile      string `json:"client_keyfile"`
	CACertFile         string `json:"cacert_file"`
}

// NewConfigFromCLI parses the command-line and returns a Config.
func NewConfigFromCLI() Config {
	var cfg Config
	pflag.StringVarP(&cfg.Host, "host", "i", "localhost", "The URL to listen on")
	pflag.IntVarP(&cfg.Port, "port", "p", 1337, "The port to run on")
	pflag.BoolVarP(&cfg.Debug, "debug", "d", false, "Whether or not to enable debug logging")
	pflag.StringVar(&cfg.MongoHost, "mongo-host", "localhost", "The host MongoDB is running on")
	pflag.IntVar(&cfg.MongoPort, "mongo-port", 27017, "The port MongoDB is running on")
	pflag.StringVar(&cfg.MongoDatabaseName, "mongo-dbname", "dc", "Name of MongoDB database to use")
	pflag.StringVar(&cfg.AMQPHost, "amqp-url", "localhost", "The URL RabbitMQ is running on")
	pflag.IntVar(&cfg.AMQPPort, "amqp-port", 5672, "The port RabbitMQ is running on")
	pflag.StringVar(&cfg.AMQPUser, "amqp-user", "guest", "Username for RabbitMQ")
	pflag.StringVar(&cfg.AMQPPassword, "amqp-password", "guest", "Password for RabbitMQ")
	pflag.StringVar(&cfg.ClientCertFile, "client-cert", "", "Client public key file")
	pflag.StringVar(&cfg.ClientKeyFile, "client-key", "", "Client private key file")
	pflag.StringVar(&cfg.CACertFile, "cacert", "", "CA Certificate file")
	pflag.Parse()
	return cfg
}
