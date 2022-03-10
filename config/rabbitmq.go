package config

type Configuration struct {
	AMQPConnectionURL string
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://user:user@localhost:5672/",
}


