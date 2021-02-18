package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/websocket"
	"github.com/mrecachinas/dcserver/internal/config"
)

// NewDCAPI creates a new api.Api with MongoDB and AMQP connections
// TODO: Maybe move this into api.go? That would require moving
//       some of the other setup functions too though...
func NewDCAPI(cfg *config.Config) (*Api, error) {
	client, err := SetupMongoConnection(cfg.MongoHost, cfg.MongoPort)
	if err != nil {
		return nil, err
	}

	amqpConnection, amqpChannel, err := SetupAMQP(cfg.AMQPHost, cfg.AMQPPort, cfg.AMQPUser, cfg.AMQPPassword)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if cfg.ClientCertFile != "" && cfg.ClientKeyFile != "" && cfg.CACertFile != "" {
		httpClient, _ = SetupHTTPClient()
	} else {
		httpClient, err = SetupHTTPSClient(cfg.ClientCertFile, cfg.ClientKeyFile, cfg.CACertFile)
		if err != nil {
			return nil, err
		}
	}

	dcapi := &Api{
		DB:          client,
		AMQPClient:  amqpConnection,
		AMQPChannel: amqpChannel,
		HTTPClient:  httpClient,
		Websocket:   SetupWebsocketConnectionPool(),
		Cfg:         cfg,
	}
	return dcapi, nil
}

// SetupMongoConnection simply (attempts to) connects to MongoDB and returns
// a pointer to a mongo.Client.
func SetupMongoConnection(mongohost string, mongoport int) (*mongo.Client, error) {
	mongouri := fmt.Sprintf("mongodb://%s:%d", mongohost, mongoport)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongouri))
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

// SetupAMQP simply (attempts to) connects to RabbitMQ (or some other AMQP broker)
// and returns pointers to amqp.Connection and amqp.Channel.
func SetupAMQP(amqphost string, amqpport int, amqpuser string, amqppassword string) (*amqp.Connection, *amqp.Channel, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", amqpuser, amqppassword, amqphost, amqpport)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func SetupHTTPClient() (*http.Client, error) {
	return &http.Client{}, nil
}

// SetupHTTPSClient sets up an HTTPS client with PKI and TLS support
// for submitting HTTPS requests (e.g., to external APIs).
func SetupHTTPSClient(certfile string, keyfile string, cacertfile string) (*http.Client, error) {
	// Read CA into memory
	cacert, err := ioutil.ReadFile(cacertfile)
	if err != nil {
		return nil, err
	}

	// Load CA
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cacert)

	// Load client PKIs
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}

	// Setup HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	return client, nil
}

// SetupWebsocketConnectionPool establishes a WebsocketConnectionPool object
// and makes an empty *websocket.Conn map.
func SetupWebsocketConnectionPool() *WebsocketConnectionPool {
	return &WebsocketConnectionPool{Connections: make(map[*websocket.Conn]struct{})}
}
