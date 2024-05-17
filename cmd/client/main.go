package main

import (
	"fmt"
	"log"

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

	gameState := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Failed to spawn a new unit: %v\n", err)
			}
		case "move":
			_, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Printf("Failed to move unit: %v\n", err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming is not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Command isn't available, try again.")
		}
	}
}
