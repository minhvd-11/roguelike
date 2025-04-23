package main

import (
	"log"

	"roguelike/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(600, 1080)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
