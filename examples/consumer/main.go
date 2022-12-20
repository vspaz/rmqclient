package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vspaz/rmqclient/pkg/rmq"
)

func main() {
	logger := logrus.New()
	connectionUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		"user",
		"password",
		"host",
		"port",
	)
	rmqClient := rmq.New(connectionUrl, "test_queue", logger)
	connection := rmqClient.Connect()
	defer rmqClient.CloseConnection(connection)
	channel := rmqClient.CreateChannel(connection)
	defer rmqClient.CloseChannel(channel)
	rmqClient.DeclareExchange(channel)
	rmqClient.DeclareQueue(channel)
	rmqClient.BindQueue(channel)
	for message := range rmqClient.Consume(channel, "") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
