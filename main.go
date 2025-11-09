package main

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	sWidth     = 480
	sHeight    = 600
	sWidthMid  = sWidth / 2
	sHeightMid = sHeight / 2
	gridCells  = 3
	gridCellSz = sWidth / gridCells
	cellPad    = 7

	fontSize               = 15
	fontSizeLineSpacing    = 5
	bigFontSize            = 100
	bigFontSizeLineSpacing = 20

	pressTicksToReset = 60
	pressTicksToExit  = 60
	keyHoldMinFrames  = 2
)

const (
	cross  = "X"
	circle = "O"
	tie    = "tie"
)

type TextAlign int

const (
	TopLeft TextAlign = iota
	TopCenter
	TopRight
	CenterLeft
	Center
	CenterRight
	BottomLeft
	BottomCenter
	BottomRight
)

//go:embed images/*
var imageFS embed.FS

var (
	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
	fontSource     *text.GoTextFaceSource
)

type GameState int

const (
	stateBoot GameState = iota
	statePlay
	stateDone
)

type Game struct {
	state   GameState
	board   [gridCells][gridCells]string
	round   int
	pointsO int
	pointsX int
	playing string
	win     string
	alter   int

	boardImg *ebiten.Image
	gameImg  *ebiten.Image
	imgCache map[string]*ebiten.Image
}

func (g *Game) Layout(int, int) (int, int) { return sWidth, sHeight }

func main() {
	game := &Game{state: stateBoot}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
