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
	// good practice to create a queue in case it does not exist.
	channel.DeclareQueue()
	channel.BindQueue()
	message := "foobar"
	if err := channel.Publish([]byte(message), "text/plain"); err != nil {
		logger.Errorf("error occured %s", message)
	}
}
