package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType simpleQueueType,
	handler func(T) Acktype,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for dat := range deliveryCh {
			var msg T
			err = json.Unmarshal(dat.Body, &msg)
			if err != nil {
				fmt.Printf("Failure Unmarshalling data: %v\n", err)
				dat.Nack(false, false) // Discard the message if unmarshalling fails
				continue
			}

			ackType := handler(msg)
			switch ackType {
			case Ack:
				dat.Ack(false)
			case NackRequeue:
				dat.Nack(false, true)
			case NackDiscard:
				dat.Nack(false, false)
			}
		}
	}()

	return nil
}
