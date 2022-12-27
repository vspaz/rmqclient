package main

import (
	"fmt"
	"github.com/vspaz/rmqclient/pkg/rmq"
	"github.com/vspaz/simplelogger/pkg/logging"
)

func main() {
	logger := logging.GetTextLogger("info").Logger
	// default test configuration for local testing
	connectionUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		"user",
		"password",
		"host",
		"5672",
	)
	rmqClient := rmq.NewClient(connectionUrl, logger)
	rmqClient.Connect()
	defer rmqClient.CloseConnection()
	broker := rmq.NewBroker("test", "test", "test", rmqClient)
	broker.CreateChannel()
	defer broker.CloseChannel()
	broker.DeclareExchange()
	broker.DeclareQueue()
	broker.BindQueue()
	for message := range broker.Consume("consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
