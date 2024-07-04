package pubsub

import (
	"context"
	"encoding/json"
	"log"

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

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
	if err != nil {
		return err
	}

	return nil
}

/**
 * creates a queue bound to the passed in exchange.
 * The function supports durable and transient queues.
 */
func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	connCh, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	q, err := createQueueType(queueName, connCh, SimpleQueueType(simpleQueueType))
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = connCh.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return connCh, q, nil
}

// creates a durable or transient queue
func createQueueType(queueName string, connCh *amqp.Channel, simpleQueueType SimpleQueueType) (amqp.Queue, error) {
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

// subscribes a queue to listen to an exchange
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T),
) error {
	// declaring and binding the queueName to an exchange or ensuring that the queue exists
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	// deliveryCh is a 'consumer' for queueName.
	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// handling each message from 'queueName'
	go func() {
		for msg := range deliveryCh {
			var data T
			err = json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Println(err)
			}

			handler(data)

			msg.Ack(false)
		}
	}()
	return nil
}
