package main

import (
	"fmt"
	"rmqclient/pkg/rmq"
)

func main() {
	rmqClient := rmq.New()
	connection := rmqClient.Connect()
	defer rmqClient.CloseConnection(connection)
	channel := rmqClient.CreateChannel(connection)
	defer rmqClient.CloseChannel(channel)
	rmqClient.DeclareExchange(channel)
	message := "foobar"
	if err := rmqClient.PublishTask(channel, []byte(message)); err != nil {
		fmt.Errorf("error occured %s", message)
	}
}
