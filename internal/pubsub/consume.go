package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// creates a durable or transient queue
func createQueueType(queueName string, connCh *amqp.Channel, simpleQueueType SimpleQueueType, args amqp.Table) (amqp.Queue, error) {
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
		args,
	)
}

// declares and binds queue to exchange and subscribes to queue
func DBSubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	// declaring and binding the queueName to an exchange or ensuring that the queue exists
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	// deliveryCh is a 'consumer' for queueName.
	// leaving consumer argument as "" so that it is randomly assigned a unique name
	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// handling each message from 'queueName'
	go func() {
		defer ch.Close()
		for msg := range deliveryCh {
			var data T
			err = json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Println(err)
			}

			ackType := handler(data)

			switch ackType {
			case Ack:
				log.Println("\nmsg ack")
				msg.Ack(false)
			case NackRequeue:
				log.Println("\nmsg nack requeue")
				msg.Nack(false, true)
			case NackDiscard:
				log.Println("\nmsg nack discard")
				msg.Nack(false, false)
			}
		}
	}()
	return nil
}
