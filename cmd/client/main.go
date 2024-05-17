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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could establish channel to RabbitMQ: %v", err)
	}

	fmt.Println("Peril game server connected to RabbitMQ!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not welcome client: %v", err)
	}

	pauseUserRoute := fmt.Sprintf("%v.%v", routing.PauseKey, username)
	armyMoveUserRoute := fmt.Sprintf("%v.%v", routing.ArmyMovesPrefix, username)
	armyMoveAnyRoute := fmt.Sprintf("%v.*", routing.ArmyMovesPrefix)

	pubsub.DeclareAndBind(
		conn,                        // conn
		routing.ExchangePerilDirect, // exchange
		pauseUserRoute,              // queueName
		routing.PauseKey,            // key
		pubsub.SimpleQueueTransient, // simpleQueueType
	)

	gameState := gamelogic.NewGameState(username)

	// army_move.*
	err = pubsub.SubscribeJSON(
		conn,                        // conn
		routing.ExchangePerilTopic,  // exchange
		armyMoveUserRoute,           // queueName
		armyMoveAnyRoute,            // key
		pubsub.SimpleQueueTransient, // simpleQueueType
		handlerMove(gameState),      // handler
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to army_move queue: %v", err)
	}

	// pause.*
	err = pubsub.SubscribeJSON(
		conn,                        // conn
		routing.ExchangePerilDirect, // exchange
		pauseUserRoute,              // queueName
		routing.PauseKey,            // key
		pubsub.SimpleQueueTransient, // simpleQueueType
		handlerPause(gameState),     // handler
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to pause queue: %v", err)
	}

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
			armyMove, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Printf("Failed to move unit: %v\n", err)
				continue
			}
			err = pubsub.PublishJSON(
				ch,                         // ch
				routing.ExchangePerilTopic, // exchange
				armyMoveUserRoute,          // key
				armyMove,                   // val
			)
			if err != nil {
				fmt.Printf("Failed to publish unit move: %v\n", err)
				continue
			}
			fmt.Println("Move published successfully")
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
