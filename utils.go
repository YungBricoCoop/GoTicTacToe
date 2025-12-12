package main

import (
	"bytes"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) initAssets() {
	var err error
	g.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		g.logger.Error("Failed to load font", "error", err)
		os.Exit(1)
	}
	g.normalTextFace = &text.GoTextFace{Source: g.fontSource, Size: fontSize}
	g.bigTextFace = &text.GoTextFace{Source: g.fontSource, Size: bigFontSize}

	g.boardImage = ebiten.NewImage(sWidth, sHeight)
}

func (g *Game) seedFirstPlayer() {
	if rand.IntN(2) == 0 { //nolint:gosec,mnd // weak rng is fine
		g.playing = PlayerO
		g.alter = 0
		return
	}
	g.playing = PlayerX
	g.alter = 1
}

func (g *Game) resetRound() {
	g.board = [gridCells][gridCells]Player{}
	g.round = 0

	// reset board image
	g.boardImage.Fill(g.colBg)
	g.drawGrid(g.boardImage)

	// alternate the round starter
	if g.alter == 0 {
		g.playing = PlayerX
		g.alter = 1
	} else {
		g.playing = PlayerO
		g.alter = 0
	}

	g.win = PlayerNone
	g.state = statePlay
}

func (g *Game) resetPoints() {
	g.pointsO = 0
	g.pointsX = 0
}

func (g *Game) currentSymbol() Player {
	if g.round%2 == g.alter {
		return PlayerO
	}
	return PlayerX
}

func otherSymbol(s Player) Player {
	if s == PlayerO {
		return PlayerX
	}
	return PlayerO
}

func (g *Game) applyWinner(w Player) {
	switch w {
	case PlayerO:
		g.win = PlayerO
		g.pointsO++
		g.state = stateDone
	case PlayerX:
		g.win = PlayerX
		g.pointsX++
		g.state = stateDone
	case PlayerTie:
		g.win = PlayerTie
		g.state = stateDone
	case PlayerNone:
		// do nothing
	}
}

func (g *Game) checkWin() Player {
	for i := range gridCells {
		a, b, c := g.board[i][0], g.board[i][1], g.board[i][2]
		if a != PlayerNone && a == b && b == c {
			return a
		}
	}
	for i := range gridCells {
		a, b, c := g.board[0][i], g.board[1][i], g.board[2][i]
		if a != PlayerNone && a == b && b == c {
			return a
		}
	}
	m := g.board[1][1]
	if m != PlayerNone {
		if g.board[0][0] == m && g.board[2][2] == m {
			return m
		}
		if g.board[0][2] == m && g.board[2][0] == m {
			return m
		}
	}

	if g.round == gridCells*gridCells-1 {
		return PlayerTie
	}
	return PlayerNone
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
