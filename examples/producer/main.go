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
		"5672",
	)
	rmqClient := rmq.New(connectionUrl, "test_queue", logger)
	connection := rmqClient.Connect()
	defer rmqClient.CloseConnection(connection)
	channel := rmqClient.CreateChannel(connection)
	defer rmqClient.CloseChannel(channel)
	rmqClient.DeclareExchange(channel)
	message := "foobar"
	if err := rmqClient.PublishTask(channel, []byte(message)); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
