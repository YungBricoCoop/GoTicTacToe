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

//TODO: implement sounds (walking, placing symbol, win, lose, draw)
//TODO: maybe: implement monsters that chase players

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
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace
	Textures       map[uint8][]*ebiten.Image
	Sprites        map[PlayerSymbol]Sprite
	XSymbolImg     *ebiten.Image
	OSymbolImg     *ebiten.Image
}

type Game struct {
	// state
	state         GameState
	board         Board
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

// NewGame creates a new Game instance with initialized assets and players.
func NewGame() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}

	spawnX := defaultPlayerXSpawn()
	spawnO := defaultPlayerOSpawn()
	pX := NewPlayer(spawnX.X, spawnX.Y, PlayerSymbolX, "X")
	pO := NewPlayer(spawnO.X, spawnO.Y, PlayerSymbolO, "O")
	world := World{
		fovScale: GetK(PlayerFOV),
	}

	g := &Game{
		state:          StateNameInput,
		currentPlayer:  pX,
		editingPlayerX: true,
		assets:         assets,

		worldMap: NewMap(),
	}

	g.gameObjects = append(g.gameObjects, &world, &g.minimap, pX, pO)

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

	xSprite, err := LoadSprite("x-player.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite x-player.png: %w", err)
	}
	oSprite, err := LoadSprite("o-player.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite o-player.png: %w", err)
	}

	xSymbol, err := LoadSprite("x.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite x.png: %w", err)
	}

	oSymbol, err := LoadSprite("o.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite o.png: %w", err)
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
		Sprites: map[PlayerSymbol]Sprite{
			PlayerSymbolX:    {Img: xSprite, Pos: Vec2{X: InitialSpriteXPosX, Y: InitialSpriteXPosY}},
			PlayerSymbolO:    {Img: oSprite, Pos: Vec2{X: InitialSpriteOPosX, Y: InitialSpriteOPosY}},
			PlayerSymbolNone: {Img: nil, Pos: Vec2{}},
		},
		XSymbolImg: xSymbol,
		OSymbolImg: oSymbol,
	}, nil
}

func (g *Game) resetBoard() {
	g.board.Reset()
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

// FIXME: This might not need to exists.
func (g *Game) getPlayer(s PlayerSymbol) *Player {
	for _, obj := range g.gameObjects {
		if p, ok := obj.(*Player); ok && p.symbol == s {
			return p
		}
	}
	return nil
}
