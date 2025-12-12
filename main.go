package main

import (
	"image/color"
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
	cellPad    = 20

	fontSize               = 15
	fontSizeLineSpacing    = 5
	bigFontSize            = 100
	bigFontSizeLineSpacing = 20

	pressTicksToReset = 60
	pressTicksToExit  = 60
	keyHoldMinFrames  = 2

	lineWidth = 10
)

var (
	colBg    = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
	colGrid  = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
	colX     = color.RGBA{0x77, 0x6e, 0x65, 0xff}
	colO     = color.RGBA{0xf2, 0xb1, 0x79, 0xff}
	colText  = color.RGBA{0x77, 0x6e, 0x65, 0xff}
	colWin   = color.RGBA{0xed, 0xc2, 0x2e, 0xff}
	colReset = color.RGBA{0x8f, 0x7a, 0x66, 0xff}
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

	boardImage *ebiten.Image
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
