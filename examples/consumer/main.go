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
		"guest",
		"guest",
		"localhost",
		"5672",
	)
	conn := rmq.NewConnection(connectionUrl, logger)
	conn.Connect()
	defer conn.CloseConnection()

	channel := rmq.NewChannel(conn, "test", "test", "test")
	channel.Create()
	defer channel.CloseChannel()
	channel.DeclareExchange("direct", true)
	channel.DeclareQueue()
	channel.BindQueue()
	for message := range channel.Consume("consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
