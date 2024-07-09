package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
 * publishes json data to a given exchange and queue key
 */
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	payload, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
}
