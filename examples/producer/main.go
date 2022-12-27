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
	rmqClient := rmq.New(connectionUrl, "test", logger)
	rmqClient.Connect()
	defer rmqClient.CloseConnection()
	rmqClient.CreateChannel()
	defer rmqClient.CloseChannel()
	rmqClient.DeclareExchange("test")
	rmqClient.BindQueue("test", "text")
	message := "foobar"
	if err := rmqClient.PublishTask([]byte(message), "text/plain"); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
