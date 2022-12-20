package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vspaz/rmqclient/pkg/rmq"
)

func main() {
	logger := logrus.New()
	// default test configuration for local testing
	connectionUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		"user",
		"password",
		"host",
		"5672",
	)
	rmqClient := rmq.New(connectionUrl, "test_queue", logger)
	connection := rmqClient.Connect()
	defer rmqClient.CloseConnection(connection)
	channel := rmqClient.CreateChannel(connection)
	defer rmqClient.CloseChannel(channel)
	rmqClient.DeclareExchange(channel)
	rmqClient.DeclareQueue(channel)
	rmqClient.BindQueue(channel)
	for message := range rmqClient.Consume(channel, "consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}