package main

import (
	"github.com/vspaz/rmqclient/pkg/rmq"
	"github.com/vspaz/simplelogger/pkg/logging"
)

func main() {
	logger := logging.GetTextLogger("info").Logger
	// default test configuration for local testing
	connection := rmq.NewConnection("amqp://guest:guest@localhost:5672", logger)
	connection.Create()
	defer connection.Close()

	channel := rmq.NewChannel(connection, "test", "test", "test")
	channel.Create()
	defer channel.Close()
	channel.DeclareExchange("direct", true)
	channel.DeclareQueue()
	channel.BindQueue()
	for message := range channel.Consume("consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
