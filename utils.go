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
	ScreenSize  = 640
	WindowSizeX = 1280
	WindowSizeY = 720

	MapGridSize = 22

	GridSize = 3
	CellSize = ScreenSize / GridSize

	Margin              = 10
	LineWidth           = 2
	HeaderY             = 20
	BottomY             = ScreenSize - 10
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

type Player int

const (
	PlayerNone Player = iota
	PlayerX
	PlayerO
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
	board         [GridSize][GridSize]Player
	currentPlayer Player
	winner        Winner

	// player info
	playerXName string
	playerOName string
	scoreX      int
	scoreO      int

	// input state
	inputBuffer    string
	editingPlayerX bool

	// visuals
	assets *VisualAssets

	// raycasting world (minimap)
	worldMap  Map
	playerPos Vec2

	minimap Minimap
}

func NewGame() *Game {
	assets := loadAssets()

	return &Game{
		state:          StateNameInput,
		currentPlayer:  PlayerX,
		playerXName:    "X",
		playerOName:    "O",
		editingPlayerX: true,
		assets:         assets,

		worldMap:  NewMap(),
		playerPos: Vec2{X: 11.5, Y: 11.5},
	}
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
	g.board = [GridSize][GridSize]Player{}
	g.winner = WinnerNone
	g.state = StatePlaying
	g.currentPlayer = PlayerX
}

func (g *Game) fullReset() {
	g.resetBoard()
	g.scoreX = 0
	g.scoreO = 0
	g.state = StateNameInput
	g.editingPlayerX = true
	g.inputBuffer = ""
	g.playerXName = "X"
	g.playerOName = "O"
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

func (g *Game) checkLine(a, b, c Player) Winner {
	if a != PlayerNone && a == b && b == c {
		if a == PlayerX {
			return WinnerX
		}
		return WinnerO
	}
	return WinnerNone
}

func (g *Game) isBoardFull() bool {
	for y := range GridSize {
		for x := range GridSize {
			if g.board[y][x] == PlayerNone {
				return false
			}
		}
	}
	return true
}
