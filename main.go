package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	loadImages()

	game := NewGame()
	ebiten.SetWindowSize(WindowSize, WindowSize)
	ebiten.SetWindowTitle("Tic-Tac-Toe (images + names + scores)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
