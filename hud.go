// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Hud struct{}

func drawFrame(dst *ebiten.Image, x, w int) {
	bw := float32(HUDBorderWidth)
	col := ColorHUDBorder

	// top
	vector.FillRect(dst, float32(x), float32(HUDY), float32(w), bw, col, false)
	// bottom
	vector.FillRect(dst, float32(x), float32(HUDY+HUDHeight)-bw, float32(w), bw, col, false)
	// left
	vector.FillRect(dst, float32(x), float32(HUDY), bw, float32(HUDHeight), col, false)
	// right
	vector.FillRect(dst, float32(x+w)-bw, float32(HUDY), bw, float32(HUDHeight), col, false)
}

func (w *Hud) Draw(screen *ebiten.Image, g *Game) {
	if screen == nil || g == nil || g.assets == nil {
		return
	}

	// HUD background
	vector.FillRect(
		screen,
		0,
		float32(HUDY),
		float32(WindowSizeX),
		float32(HUDHeight),
		ColorHUDFill,
		false,
	)

	// Layout widths
	leftW := HUDLeftW
	infoW := HUDInfoW
	rightW := HUDRightW
	scoreW := HUDCenterW - infoW

	// Frames
	drawFrame(screen, 0, leftW)
	drawFrame(screen, leftW, scoreW)
	drawFrame(screen, leftW+scoreW, infoW)
	drawFrame(screen, leftW+scoreW+infoW, rightW)

	// LEFT : player icon
	playerSymbolTexture := g.assets.Textures[g.currentPlayer.symbolTextureID]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(HUDTextPadX),
		float64(HUDY+HUDTextPadY),
	)
	screen.DrawImage(
		playerSymbolTexture.Source,
		op,
	)

	// Score (center block)
	pX := g.playerX
	pO := g.playerO
	if pX != nil && pO != nil {
		scoreCenterX := float64(leftW) + float64(scoreW)/HalfDivisor
		score := "Score  X: " + strconv.Itoa(pX.score) + "   O: " + strconv.Itoa(pO.score)
		g.drawText(
			screen,
			score,
			scoreCenterX,
			float64(HUDY+HUDTextPadY),
			TopCenter,
			color.White,
		)
	}

	// CENTER-RIGHT : INFO
	infoCenterX := float64(leftW+scoreW) + float64(infoW)/HalfDivisor

	g.drawText(
		screen,
		"ESC : Quit",
		infoCenterX,
		float64(HUDY+HUDTextPadY),
		TopCenter,
		color.RGBA{220, 220, 220, 255},
	)

	g.drawText(
		screen,
		"Ctrl + R : Restart",
		infoCenterX,
		float64(HUDY+HUDTextPadY+HUDResetYOffset),
		TopCenter,
		color.RGBA{180, 180, 180, 255},
	)

	g.drawText(
		screen,
		"E : Place Marker",
		infoCenterX,
		float64(HUDY+HUDTextPadY+HUDPlaceYOffset),
		TopCenter,
		color.RGBA{180, 180, 180, 255},
	)
}
