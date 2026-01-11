// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// PlayerSymbol represents the symbol used by a player in the game, either X or O.
type PlayerSymbol int

const (
	PlayerSymbolNone PlayerSymbol = iota
	PlayerSymbolX
	PlayerSymbolO
)

// Player represents a player in the game.
// pos is the player's position in the world.
// dir is the player's direction vector.
// symbol is the player's symbol (X or O).
// name is the player's name.
// score is the player's score.
type Player struct {
	pos    Vec2
	dir    Vec2
	symbol PlayerSymbol
	name   string
	score  int
}

// NewPlayer creates a new player with the given position, symbol, and name.
func NewPlayer(x, y float64, symbol PlayerSymbol, name string) *Player {
	return &Player{
		pos:    Vec2{x, y},
		dir:    Vec2{-1, 0},
		symbol: symbol,
		name:   name,
		score:  0,
	}
}

func (p *Player) Update(g *Game) {
	// only update if this is the current player
	if g.currentPlayer != p {
		return
	}

	moveSpeed := PlayerMovementSpeed * DeltaTime
	rotSpeed := PlayerRotationSpeed * DeltaTime

	// w/s for forward/backward
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.move(g, p.dir.Scale(moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.move(g, p.dir.Scale(-moveSpeed))
	}

	// a/d for rotation
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.rotate(rotSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.rotate(-rotSpeed)
	}
}

// Move the player by the given velocity vector, checking for collisions.
func (p *Player) move(g *Game, velocity Vec2) {
	// try to move in X
	nextPos := p.pos.Add(Vec2{X: velocity.X, Y: 0})
	if g.worldMap.IsWalkable(nextPos) {
		p.pos = nextPos
	}

	// try to move in Y
	nextPos = p.pos.Add(Vec2{X: 0, Y: velocity.Y})
	if g.worldMap.IsWalkable(nextPos) {
		p.pos = nextPos
	}
}

// Rotate the player by the given angle in radians.
func (p *Player) rotate(angle float64) {
	p.dir = p.dir.Rotate(angle)
}
