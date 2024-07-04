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
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connection was successful...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalln(err)
	}
	_, _, err = pubsub.DeclareAndBind(conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalln(err)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatalln(err)
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

			fmt.Printf("moved to: %s by player: %s of %v unit(s) successful\n",
				move.ToLocation,
				move.Player.Username,
				len(move.Units),
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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(playingState routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
	}
}
