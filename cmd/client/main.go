package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril client...")
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not welcome client: %v", err)
	}

	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)

	pubsub.DeclareAndBind(
		conn,                        // conn
		routing.ExchangePerilDirect, // exchange
		queueName,                   // queueName
		routing.PauseKey,            // key
		pubsub.SimpleQueueTransient, // simpleQueueType
	)

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("\nReceived interrupt signal. Exiting...")
}
