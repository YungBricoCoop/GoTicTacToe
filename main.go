package main

import (
	"image/color"
	"log/slog"
	"os"

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

	feedbackMarginBottom = 30

	pressTicksToReset = 60
	pressTicksToExit  = 60
	keyHoldMinFrames  = 2

	lineWidth = 10
)

type Player int8

const (
	PlayerNone Player = iota
	PlayerX
	PlayerO
	PlayerTie
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

type GameState int

const (
	stateBoot GameState = iota
	statePlay
	stateDone
)

type Game struct {
	state   GameState
	board   [gridCells][gridCells]Player
	round   int
	pointsO int
	pointsX int
	playing Player
	win     Player
	alter   int

	boardImage *ebiten.Image

	colBg    color.RGBA
	colGrid  color.RGBA
	colX     color.RGBA
	colO     color.RGBA
	colText  color.RGBA
	colWin   color.RGBA
	colReset color.RGBA

	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
	fontSource     *text.GoTextFaceSource

	logger *slog.Logger
}

func (g *Game) Layout(int, int) (int, int) { return sWidth, sHeight }

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	game := &Game{
		state:    stateBoot,
		colBg:    color.RGBA{0xfa, 0xf8, 0xef, 0xff},
		colGrid:  color.RGBA{0xbb, 0xad, 0xa0, 0xff},
		colX:     color.RGBA{0x77, 0x6e, 0x65, 0xff},
		colO:     color.RGBA{0xf2, 0xb1, 0x79, 0xff},
		colText:  color.RGBA{0x77, 0x6e, 0x65, 0xff},
		colWin:   color.RGBA{0xed, 0xc2, 0x2e, 0xff},
		colReset: color.RGBA{0x8f, 0x7a, 0x66, 0xff},
		logger:   logger,
	}

	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe Raycast")
	if err := ebiten.RunGame(game); err != nil {
		logger.Error("Game crashed", "error", err)
		os.Exit(1)
	}
}
