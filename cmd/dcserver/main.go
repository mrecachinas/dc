package main

import (
	"github.com/mrecachinas/dcserver/internal/app"
	"github.com/spf13/pflag"
)

func main() {
	cfg := &Config{}
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
	pflag.IntVar(&cfg.PollingInterval, "polling-interval", 5, "Number of seconds between database polls")
	pflag.Parse()

	app.Run(cfg)
}
