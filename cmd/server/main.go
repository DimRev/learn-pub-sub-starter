package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conStr := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	conn, err := amqp.Dial(conStr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Println("Peril game server connected to RabbitMQ!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

}
