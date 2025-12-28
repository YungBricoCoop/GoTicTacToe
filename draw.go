// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

/* ---------- helpers ---------- */

func fillRect(dst *ebiten.Image, x, y, w, h float32, col color.Color) {
	c := color.RGBAModel.Convert(col)
	rgba, ok := c.(color.RGBA)
	if !ok {
		rgba = color.RGBA{}
	}
	vector.FillRect(
		dst,
		x, y, w, h,
		rgba,
		false,
	)
}

/* ---------- draw loop ---------- */

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBackground)

	switch g.state {
	case StateNameInput:
		g.drawNameInput(screen)
	case StatePlaying:
		g.drawPlaying(screen)
		g.minimap.Draw(screen, g)
	case StateGameOver:
		g.drawPlaying(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return WindowSizeX, WindowSizeY
}

/* ---------- UI ---------- */

func (g *Game) drawNameInput(screen *ebiten.Image) {
	g.drawText(screen, "Enter player names", NameInputX, NameInputY, TopLeft, color.White)

	label := "Player X: "
	if !g.editingPlayerX {
		label = "Player O: "
	}
	g.drawText(screen, label+g.inputBuffer, NameInputX, NameInputY+NameInputLineHeight, TopLeft, color.White)

	info := "Type name, Enter = OK, Backspace = delete"
	g.drawText(screen, info, NameInputX, NameInputY+NameInputLineHeight*2, TopLeft, color.White)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	g.drawScoreAndShortcuts(screen)

	if g.state == StatePlaying {
		g.drawTurnInfo(screen)
	}
}

func (g *Game) drawScoreAndShortcuts(screen *ebiten.Image) {
	pX := g.getPlayer(PlayerSymbolX)
	pO := g.getPlayer(PlayerSymbolO)
	score := "Score " + pX.name + ": " + strconv.Itoa(pX.score) +
		"  " + pO.name + ": " + strconv.Itoa(pO.score)

	g.drawText(screen, score, Margin, HeaderY, TopLeft, color.White)
	g.drawText(screen, "ESC = quit", WindowSizeX-Margin, HeaderY, TopRight, color.White)
}

func (g *Game) drawTurnInfo(screen *ebiten.Image) {
	msg := "Turn: " + g.currentPlayer.name
	msg += "   (Ctrl + R = full reset)"
	g.drawText(screen, msg, Margin, BottomY, BottomLeft, color.White)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	msg := "Draw!"
	switch g.winner {
	case WinnerNone, WinnerDraw:
		msg = "Draw!"
	case WinnerX:
		msg = g.getPlayer(PlayerSymbolX).name + " wins!"
	case WinnerO:
		msg = g.getPlayer(PlayerSymbolO).name + " wins!"
	}
	msg += " Click to restart"
	g.drawText(screen, msg, Margin, BottomY, BottomLeft, color.White)
}

/* ---------- text ---------- */

func (g *Game) drawText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color) {
	const half = 2

	op := &text.DrawOptions{}
	op.LineSpacing = g.assets.NormalTextFace.Size + TextLineSpacing
	op.ColorScale.ScaleWithColor(col)

	w, h := text.Measure(msg, g.assets.NormalTextFace, op.LineSpacing)

	switch align {
	case TopLeft, CenterLeft, BottomLeft:
	case TopCenter, Center, BottomCenter:
		x -= w / half
	case TopRight, CenterRight, BottomRight:
		x -= w
	}

	switch align {
	case TopLeft, TopCenter, TopRight:
	case CenterLeft, Center, CenterRight:
		y -= h / half
	case BottomLeft, BottomCenter, BottomRight:
		y -= h
	}

	op.GeoM.Translate(x, y)
	text.Draw(screen, msg, g.assets.NormalTextFace, op)
}
