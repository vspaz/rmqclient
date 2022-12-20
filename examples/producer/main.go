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
