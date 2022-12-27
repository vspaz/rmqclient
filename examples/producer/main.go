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
	conn := rmq.NewConnection(connectionUrl, logger)
	conn.Connect()
	defer conn.CloseConnection()

	broker := rmq.NewChannel("test", "test", "test", conn)
	broker.Create()
	defer broker.CloseChannel()
	broker.DeclareExchange()
	broker.DeclareQueue()
	broker.BindQueue()
	message := "foobar"
	if err := broker.PublishTask([]byte(message), "text/plain"); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
