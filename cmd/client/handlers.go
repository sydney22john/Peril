package main

import (
	"fmt"
	"sjohn/Peril/internal/gamelogic"
	"sjohn/Peril/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(playingState routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
	}
}

func handlerMoves(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(mo gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(mo)
	}
}
