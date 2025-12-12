package main

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) initAssets() {
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	normalTextFace = &text.GoTextFace{Source: fontSource, Size: fontSize}
	bigTextFace = &text.GoTextFace{Source: fontSource, Size: bigFontSize}

	g.boardImage = ebiten.NewImage(sWidth, sHeight)
}

func (g *Game) seedFirstPlayer() {
	if newRandom().Intn(2) == 0 {
		g.playing = circle
		g.alter = 0
		return
	}
	g.playing = cross
	g.alter = 1
}

func (g *Game) resetRound() {
	g.board = [gridCells][gridCells]string{}
	g.round = 0

	// reset board image
	g.boardImage.Fill(colBg)
	g.drawGrid(g.boardImage)

	// alternate the round starter
	if g.alter == 0 {
		g.playing = cross
		g.alter = 1
	} else {
		g.playing = circle
		g.alter = 0
	}

	g.win = ""
	g.state = statePlay
}

func (g *Game) resetPoints() {
	g.pointsO = 0
	g.pointsX = 0
}

func (g *Game) currentSymbol() string {
	if g.round%2 == g.alter {
		return circle
	}
	return cross
}

func otherSymbol(s string) string {
	if s == circle {
		return cross
	}
	return circle
}

func (g *Game) applyWinner(w string) {
	switch w {
	case circle:
		g.win = circle
		g.pointsO++
		g.state = stateDone
	case cross:
		g.win = cross
		g.pointsX++
		g.state = stateDone
	case tie:
		g.win = "No one"
		g.state = stateDone
	}
}

func (g *Game) checkWin() string {
	for i := 0; i < gridCells; i++ {
		a, b, c := g.board[i][0], g.board[i][1], g.board[i][2]
		if a != "" && a == b && b == c {
			return a
		}
	}
	for i := 0; i < gridCells; i++ {
		a, b, c := g.board[0][i], g.board[1][i], g.board[2][i]
		if a != "" && a == b && b == c {
			return a
		}
	}
	m := g.board[1][1]
	if m != "" {
		if g.board[0][0] == m && g.board[2][2] == m {
			return m
		}
		if g.board[0][2] == m && g.board[2][0] == m {
			return m
		}
	}

	if g.round == gridCells*gridCells-1 {
		return tie
	}
	return ""
}

func getCursorIdxFromCell() (int, int, bool) {
	mx, my := ebiten.CursorPosition()
	if mx < 0 || my < 0 {
		return 0, 0, false
	}
	cx := mx / gridCellSz
	cy := my / gridCellSz
	if cx >= gridCells || cy >= gridCells {
		return 0, 0, false
	}
	return cx, cy, true
}

func newRandom() *rand.Rand {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1)
}
