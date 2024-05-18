package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, pubCh *amqp.Channel) func(armyMove gamelogic.ArmyMove) pubsub.Acktype {
	return func(armyMove gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(armyMove)
		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			log.Print("Ack message")
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			key := fmt.Sprintf("%v.%v", routing.WarRecognitionsPrefix, gs.GetUsername())
			err := pubsub.PublishJSON(pubCh, routing.ExchangePerilTopic, key, gamelogic.RecognitionOfWar{
				Attacker: armyMove.Player,
				Defender: gs.GetPlayerSnap(),
			})
			if err != nil {
				log.Print("Nack message")
				return pubsub.NackRequeue
			}
			log.Print("Ack message")
			return pubsub.Ack
		default:
			log.Print("Nack message")
			return pubsub.NackDiscard
		}
	}
}
