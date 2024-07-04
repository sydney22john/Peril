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
	fmt.Println("Starting Peril client...")
	rabbitMQconn := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitMQconn)
	if err != nil {
		log.Fatalf("couldn't connect to rabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection was successful...")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("couldn't open channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("couldn't get username: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

	// subscribing to the pause/resume game queue
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gameState.GetUsername(),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("couldn't subscribe to pause: %v", err)
	}

	// subscribing to the army_moves.* queues. Consumes messages from other players and itself
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gameState.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
		handlerMoves(gameState),
	)
	if err != nil {
		log.Fatalf("couldn't subscribe to moves: %v", err)
	}

repl:
	for true {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			move, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+gameState.GetUsername(),
				move,
			)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break repl
		default:
			log.Printf("didn't recognize command %s", words[0])
		}
	}
}
