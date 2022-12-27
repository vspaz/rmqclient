package main

import (
	"fmt"
	"github.com/vspaz/rmqclient/pkg/rmq"
	"github.com/vspaz/simplelogger/pkg/logging"
)

func main() {
	logger := logging.GetTextLogger("info").Logger
	connectionUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		"user",
		"password",
		"host",
		"5672",
	)
	rmqClient := rmq.New(connectionUrl, logger)
	connection := rmqClient.Connect()
	defer rmqClient.CloseConnection()
	broker := rmqClient.CreateBroker("test", "test", "test")
	broker.CreateChannel(connection)
	defer broker.CloseChannel()
	broker.DeclareExchange()
	broker.BindQueue()
	message := "foobar"
	if err := broker.PublishTask([]byte(message), "text/plain"); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
