package main

import (
	"encoding/json"
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
	channel.DeclareQueue(true)
	channel.BindQueue()
	message, _ := json.Marshal(map[string]string{"go": "test"})
	if err := channel.Publish(message, "application/json"); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
