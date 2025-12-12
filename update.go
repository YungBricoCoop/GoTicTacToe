package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	g.handleGlobalInput()

	switch g.state {
	case stateBoot:
		return g.updateBoot()
	case statePlay:
		return g.updatePlay()
	case stateDone:
		return g.updateDone()
	}

	return nil
}

func (g *Game) handleGlobalInput() {
	if inpututil.KeyPressDuration(ebiten.KeyR) == pressTicksToReset {
		g.resetRound()
		g.resetPoints()
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == pressTicksToExit {
		os.Exit(0)
	}
}

func (g *Game) updateBoot() error {
	g.initAssets()
	g.seedFirstPlayer()
	g.resetRound()
	g.resetPoints()
	g.state = statePlay
	return nil
}

func (g *Game) updatePlay() error {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	cx, cy, ok := getCursorIdxFromCell()
	if !ok {
		return nil
	}

	if g.board[cx][cy] != PlayerNone {
		return nil
	}

	cur := g.currentSymbol()
	nxt := otherSymbol(cur)

	g.board[cx][cy] = cur
	g.drawSymbol(g.boardImage, cx, cy, cur)
	g.playing = nxt

	w := g.checkWin()
	if w != PlayerNone {
		g.applyWinner(w)
		return nil
	}

	g.round++
	return nil
}

func (g *Game) updateDone() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.resetRound()
		g.state = statePlay
	}
	return nil
}
