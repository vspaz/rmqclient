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
	conn := rmq.NewConnection(connectionUrl, logger)
	conn.Connect()
	defer conn.CloseConnection()

	channel := rmq.NewChannel("test", "test", "test", conn)
	channel.Create()
	defer channel.CloseChannel()
	channel.DeclareExchange()
	channel.DeclareQueue()
	channel.BindQueue()
	for message := range channel.Consume("consumer1") {
		logger.Infof("message recieved: %s", string(message.Body))
	}
}
