package main

import (
	"fmt"
	"log"
	"sjohn/Peril/internal/gamelogic"
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
	defer connChan.Close()
	if err != nil {
		log.Fatalln(err)
	}

	gamelogic.PrintServerHelp()
	handleErr := func(err error) {
		if err != nil {
			log.Println(err)
		}
	}
repl:
	for true {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			log.Printf("sending pause message\n")
			handleErr(pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			}))
		case "resume":
			log.Printf("sending resume message\n")
			handleErr(pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			}))
		case "quit":
			log.Printf("exiting program...")
			break repl
		default:
			log.Printf("I don't understand the command '%s'\n", words[0])
		}

	}
	fmt.Println("The program is shutting down. Connection is being closed.")
}
