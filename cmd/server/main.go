package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("The program is shutting down. Connection is being closed.")
}
