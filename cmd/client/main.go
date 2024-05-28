package main

import (
	"fmt"
	"log"
	"os"
	"sjohn/Peril/internal/gamelogic"
	"sjohn/Peril/internal/helpers"
	"sjohn/Peril/internal/pubsub"
	"sjohn/Peril/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	rabbitMQconn := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMQconn)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connection was successful...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalln(err)
	}
	pubsub.DeclareAndBind(conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.Transient,
	)

	helpers.BlockUntilSignal(os.Interrupt)
}
