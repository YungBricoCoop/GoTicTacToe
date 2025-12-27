// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	WindowSizeX = 1280
	WindowSizeY = 720
	TPS         = 60
	DeltaTime   = 1.0 / TPS

	MapGridSize = 22

	GridSize = 3
	CellSize = WindowSizeY / GridSize

	Margin              = 10
	LineWidth           = 2
	HeaderY             = 20
	BottomY             = WindowSizeY - 10
	TextLineSpacing     = 5
	BigTextLineSpacing  = 20
	DefaultFontSize     = 15
	BigFontSize         = 100
	NameInputX          = 10
	NameInputY          = 40
	NameInputLineHeight = 40
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

type VisualAssets struct {
	XImage         *ebiten.Image
	OImage         *ebiten.Image
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace
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
	assets *VisualAssets

	// raycasting world (minimap)
	worldMap Map

	minimap Minimap

	gameObjects []GameObject
}

func NewGame() *Game {
	assets := loadAssets()

	pX := NewPlayer(11.5, 11.5, PlayerSymbolX, "X")
	pO := NewPlayer(15.5, 15.5, PlayerSymbolO, "O")

	g := &Game{
		state:          StateNameInput,
		currentPlayer:  pX,
		editingPlayerX: true,
		assets:         assets,

		worldMap: NewMap(),
	}

	g.gameObjects = append(g.gameObjects, &g.minimap, pX, pO)

	return g
}

func loadAssets() *VisualAssets {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}

	xImg := createXImage()
	oImg := createOImage()

	return &VisualAssets{
		XImage: xImg,
		OImage: oImg,
		NormalTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   DefaultFontSize,
		},
		BigTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   BigFontSize,
		},
	}
}

func createXImage() *ebiten.Image {
	const (
		strokeWidth  = float32(10)
		padding      = float32(20)
		doubleMargin = 2 * Margin
	)
	size := CellSize - doubleMargin
	img := ebiten.NewImage(size, size)

	s := float32(size)
	colX := color.RGBA{255, 100, 100, 255}

	vector.StrokeLine(img, padding, padding, s-padding, s-padding, strokeWidth, colX, true)
	vector.StrokeLine(img, s-padding, padding, padding, s-padding, strokeWidth, colX, true)

	return img
}

func createOImage() *ebiten.Image {
	const (
		strokeWidth  = float32(10)
		padding      = 40
		doubleMargin = 2 * Margin
		half         = 2
	)
	size := CellSize - doubleMargin
	img := ebiten.NewImage(size, size)

	s := float32(size)
	center := s / half
	radius := (s - padding) / half

	colO := color.RGBA{100, 100, 255, 255}
	vector.StrokeCircle(img, center, center, radius, strokeWidth, colO, true)

	return img
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
