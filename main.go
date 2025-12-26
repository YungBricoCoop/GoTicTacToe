// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame()

	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetWindowTitle("Tic-Tac-Toe")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
