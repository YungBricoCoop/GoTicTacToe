// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState int

const (
	StateNameInput GameState = iota
	StatePlaying
	StateGameOver
)

type Winner int

const (
	WinnerNone Winner = iota
	WinnerX
	WinnerO
	WinnerDraw
)

type Assets struct {
	XImage         *ebiten.Image
	OImage         *ebiten.Image
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace
	Textures       map[uint8][]*ebiten.Image
}

type Game struct {
	// state
	state         GameState
	board         [GridSize][GridSize]PlayerSymbol
	currentPlayer *Player
	winner        Winner

	// input state
	inputBuffer    string
	editingPlayerX bool

	// visuals
	assets *Assets

	// raycasting world (minimap)
	worldMap Map

	minimap Minimap

	gameObjects []GameObject
}

func NewGame() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}

	spawnX := defaultPlayerXSpawn()
	spawnO := defaultPlayerOSpawn()
	pX := NewPlayer(spawnX.X, spawnX.Y, PlayerSymbolX, "X")
	pO := NewPlayer(spawnO.X, spawnO.Y, PlayerSymbolO, "O")

	g := &Game{
		state:          StateNameInput,
		currentPlayer:  pX,
		editingPlayerX: true,
		assets:         assets,

		worldMap: NewMap(),
	}

	g.gameObjects = append(g.gameObjects, &g.minimap, pX, pO)

	return g, nil
}

func loadAssets() (*Assets, error) {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
	}

	textures, err := LoadTextures()
	if err != nil {
		return nil, fmt.Errorf("load textures: %w", err)
	}

	return &Assets{
		NormalTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   DefaultFontSize,
		},
		BigTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   BigFontSize,
		},
		Textures: textures,
	}, nil
}

func (g *Game) resetBoard() {
	g.board = [GridSize][GridSize]PlayerSymbol{}
	g.winner = WinnerNone
	g.state = StatePlaying
	g.currentPlayer = g.getPlayer(PlayerSymbolX)
}

func (g *Game) fullReset() {
	g.resetBoard()
	pX := g.getPlayer(PlayerSymbolX)
	pO := g.getPlayer(PlayerSymbolO)
	pX.score = 0
	pO.score = 0
	g.state = StateNameInput
	g.editingPlayerX = true
	g.inputBuffer = ""
	pX.name = "X"
	pO.name = "O"
}

func (g *Game) getPlayer(s PlayerSymbol) *Player {
	for _, obj := range g.gameObjects {
		if p, ok := obj.(*Player); ok && p.symbol == s {
			return p
		}
	}
	return nil
}

func (g *Game) checkWinner() Winner {
	for i := range GridSize {
		if w := g.checkLine(g.board[i][0], g.board[i][1], g.board[i][2]); w != WinnerNone {
			return w
		}
		if w := g.checkLine(g.board[0][i], g.board[1][i], g.board[2][i]); w != WinnerNone {
			return w
		}
	}

	if w := g.checkLine(g.board[0][0], g.board[1][1], g.board[2][2]); w != WinnerNone {
		return w
	}
	if w := g.checkLine(g.board[0][2], g.board[1][1], g.board[2][0]); w != WinnerNone {
		return w
	}

	if g.isBoardFull() {
		return WinnerDraw
	}
	return WinnerNone
}

func (g *Game) checkLine(a, b, c PlayerSymbol) Winner {
	if a != PlayerSymbolNone && a == b && b == c {
		if a == PlayerSymbolX {
			return WinnerX
		}
		return WinnerO
	}
	return WinnerNone
}

func (g *Game) isBoardFull() bool {
	for y := range GridSize {
		for x := range GridSize {
			if g.board[y][x] == PlayerSymbolNone {
				return false
			}
		}
	}
	return true
}
