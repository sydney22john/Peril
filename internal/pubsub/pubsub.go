package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	payload, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
	if err != nil {
		return err
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	connCh, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	q, err := createQueueType(queueName, connCh, simpleQueueType)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = connCh.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return connCh, q, nil
}

func createQueueType(queueName string, connCh *amqp.Channel, simpleQueueType int) (amqp.Queue, error) {
	var durable, autoDelete, exclusive bool
	switch simpleQueueType {
	// durable
	case Durable:
		durable = true
		autoDelete = false
		exclusive = false
		// transient
	case Transient:
		durable = false
		autoDelete = true
		exclusive = true
	}

	return connCh.QueueDeclare(queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
}
