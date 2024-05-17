package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(armyMove gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(armyMove)
		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			log.Print("Ack message")
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			log.Print("Ack message")
			return pubsub.Ack
		default:
			log.Print("Nack message")
			return pubsub.NackDiscard
		}
	}
}
