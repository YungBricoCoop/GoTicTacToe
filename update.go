// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	// update all game objects
	for _, obj := range g.gameObjects {
		obj.Update(g)
	}

	// shortcuts
	// Escape: exit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// R: reset game state
	if inpututil.IsKeyJustPressed(ebiten.KeyR) &&
		(ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) {
		g.fullReset()
		return nil
	}

	switch g.state {
	case StateNameInput:
		return g.updateNameInput()
	case StatePlaying:
		return g.updatePlaying()
	case StateGameOver:
		return g.updateGameOver()
	}

	return nil
}

func (g *Game) updateNameInput() error {
	chars := ebiten.AppendInputChars(nil)
	for _, c := range chars {
		if c == '\n' || c == '\r' {
			continue
		}
		g.inputBuffer += string(c)
	}

	// Backspace: delete last character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.inputBuffer) > 0 {
			g.inputBuffer = g.inputBuffer[:len(g.inputBuffer)-1]
		}
	}

	// Enter: confirm name
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.confirmName()
	}

	return nil
}

func (g *Game) confirmName() {
	name := g.inputBuffer
	if name == "" {
		// default names if empty
		if g.editingPlayerX {
			name = "X"
		} else {
			name = "O"
		}
	}

	if g.editingPlayerX {
		g.getPlayer(PlayerSymbolX).name = name
		g.editingPlayerX = false
		g.inputBuffer = ""
	} else {
		g.getPlayer(PlayerSymbolO).name = name
		g.state = StatePlaying
		g.inputBuffer = ""
	}
}

func (g *Game) updatePlaying() error {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	mx, my := ebiten.CursorPosition()
	if !isInBounds(mx, my) {
		return nil
	}

	cx, cy := mx/CellSize, my/CellSize

	// cell must be empty
	if g.board[cy][cx] != PlayerSymbolNone {
		return nil
	}

	g.board[cy][cx] = g.currentPlayer.symbol

	winner := g.checkWinner()
	if winner != WinnerNone {
		g.handleGameEnd(winner)
		return nil
	}

	g.switchPlayer()
	return nil
}

func (g *Game) updateGameOver() error {
	// click to restart
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.resetBoard()
	}
	return nil
}

func (g *Game) switchPlayer() {
	if g.currentPlayer.symbol == PlayerSymbolX {
		g.currentPlayer = g.getPlayer(PlayerSymbolO)
	} else {
		g.currentPlayer = g.getPlayer(PlayerSymbolX)
	}
}

func (g *Game) handleGameEnd(w Winner) {
	g.winner = w
	g.state = StateGameOver

	switch w {
	case WinnerX:
		g.getPlayer(PlayerSymbolX).score++
	case WinnerO:
		g.getPlayer(PlayerSymbolO).score++
	case WinnerNone, WinnerDraw:
		// no score update
	}
}

func isInBounds(x, y int) bool {
	return x >= 0 && x < WindowSizeY && y >= 0 && y < WindowSizeY
}
