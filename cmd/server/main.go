package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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
	blockUntilSignal(os.Interrupt)
	fmt.Println("The program is shutting down. Connection is being closed.")
}

func blockUntilSignal(sig os.Signal) {
	// wait for signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, sig)
	<-signalChan
}
