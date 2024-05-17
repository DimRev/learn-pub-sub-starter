package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType simpleQueueType,
	handler func(T),
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
			}
			handler(msg)
			dat.Ack(false)
		}
	}()

	return nil
}
