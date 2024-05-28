package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sjohn/Peril/internal/helpers"
	"sjohn/Peril/internal/pubsub"
	"sjohn/Peril/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	rabbitMQconn := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMQconn)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connection was successful...")
	connChan, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	// wait for ctrl+c
	helpers.BlockUntilSignal(os.Interrupt)
	fmt.Println("The program is shutting down. Connection is being closed.")
}
