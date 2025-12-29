// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerSymbol int

const (
	PlayerSymbolNone PlayerSymbol = iota
	PlayerSymbolX
	PlayerSymbolO
)

const (
	MovementSpeed = 5.0  // units per second
	RotationSpeed = 3.0  // radians per second
	FOV           = 60.0 // degrees
)

type Player struct {
	pos    Vec2
	dir    Vec2
	plane  Vec2
	symbol PlayerSymbol
	name   string
	score  int
}

func NewPlayer(x, y float64, symbol PlayerSymbol, name string) *Player {
	return &Player{
		pos:    Vec2{x, y},
		dir:    Vec2{-1, 0},
		plane:  Vec2{0, 0.66},
		symbol: symbol,
		name:   name,
		score:  0,
	}
}

func (p *Player) Update(g *Game) {
	if g.currentPlayer != p {
		return
	}

	// w/s for forward/backward
	dt := DeltaTime

	moveSpeed := MovementSpeed * dt
	rotSpeed := RotationSpeed * dt

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.move(g, p.dir.Scale(moveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.move(g, p.dir.Scale(-moveSpeed))
	}

	// a/d for rotation
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.rotate(-rotSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.rotate(rotSpeed)
	}
}

func (p *Player) Draw(_ *ebiten.Image, _ *Game) {
	// nothing to draw for now, maybe we will draw an asset for each player
}

func (p *Player) move(g *Game, velocity Vec2) {
	// try to move in X
	nextPos := p.pos.Add(Vec2{X: velocity.X, Y: 0})
	if isValidPosition(g.worldMap, nextPos) {
		p.pos = nextPos
	}

	// try to move in Y
	nextPos = p.pos.Add(Vec2{X: 0, Y: velocity.Y})
	if isValidPosition(g.worldMap, nextPos) {
		p.pos = nextPos
	}
}

func (p *Player) rotate(angle float64) {
	p.dir = p.dir.Rotate(angle)
	p.plane = p.plane.Rotate(angle)
}

func isValidPosition(m Map, pos Vec2) bool {
	x, y := int(pos.X), int(pos.Y)
	if x < 0 || x >= m.Width() || y < 0 || y >= m.Height() {
		return false
	}
	return m.Tiles[y][x] == 0
}
