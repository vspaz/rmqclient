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
	rmqClient := rmq.New(connectionUrl, logger)
	rmqClient.Connect()
	defer rmqClient.CloseConnection()
	rmqClient.CreateChannel("test", "test", "test")
	defer rmqClient.CloseChannel()
	rmqClient.DeclareExchange()
	rmqClient.DeclareQueue("test")
	rmqClient.BindQueue()
	for message := range rmqClient.Consume("test", "consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
