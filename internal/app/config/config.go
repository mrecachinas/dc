package config

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
	PollingInterval    int    `json:"polling_interval"`
}
