package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ScreenSize = 480
	WindowSize = 480

	GridSize = 3
	CellSize = ScreenSize / GridSize

	Margin    = 10
	LineWidth = 2
	HeaderY   = 20
	BottomY   = ScreenSize - 10

	// font sizes
	FontSize               = 15
	FontSizeLineSpacing    = 5
	BigFontSize            = 100
	BigFontSizeLineSpacing = 20
)

var (
	xImage *ebiten.Image
	oImage *ebiten.Image
)

type Game struct {
	board  [GridSize][GridSize]int // 0 = empty, 1 = X, 2 = O
	player int                     // current player: 1 (X) or 2 (O)
	winner int                     // 0 = nobody, 1 = X, 2 = O, 3 = draw
	over   bool                    // game over?

	// scores
	pointsX int
	pointsO int

	// Player names
	playerXName string
	playerOName string

	// Name input
	entering  bool   // true while entering names
	editingX  bool   // true = entering X's name, false = O's name
	tempInput string // buffer for current text

	startPlayer int // 1 = x, 2 = o

	// Text faces
	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
	fontSource     *text.GoTextFaceSource
}

func NewGame() *Game {
	fontSource, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

	return &Game{
		player:         1, // X starts
		startPlayer:    1,
		playerXName:    "X",  // default values just in case
		playerOName:    "O",  // default values
		entering:       true, // start by entering names
		editingX:       true, // start with player X
		pointsX:        0,
		pointsO:        0,
		fontSource:     fontSource,
		normalTextFace: &text.GoTextFace{Source: fontSource, Size: FontSize},
		bigTextFace:    &text.GoTextFace{Source: fontSource, Size: BigFontSize},
	}
}

func mustLoadImage(path string) *ebiten.Image {
	data, err := imageFS.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read %s: %v", path, err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("cannot decode %s: %v", path, err)
	}
	return ebiten.NewImageFromImage(img)
}

func loadImages() {
	xImage = mustLoadImage("images/x.png")
	oImage = mustLoadImage("images/o.png")
}
