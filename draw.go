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

// fillRect fills a rectangle on the given destination image with the specified color
func fillRect(dst *ebiten.Image, x, y, w, h float32, col color.Color) {
	c := color.RGBAModel.Convert(col)
	rgba, ok := c.(color.RGBA)
	if !ok {
		rgba = color.RGBA{}
	}
	vector.FillRect(dst, x, y, w, h, rgba, false)
}

func drawFrame(dst *ebiten.Image, x, y, w, h int) {
	bw := HUDBorderWidth
	col := ColorHUDBorder

	// top
	fillRect(dst, float32(x), float32(y), float32(w), bw, col)
	// bottom
	fillRect(dst, float32(x), float32(y+h)-bw, float32(w), bw, col)
	// left
	fillRect(dst, float32(x), float32(y), bw, float32(h), col)
	// right
	fillRect(dst, float32(x+w)-bw, float32(y), bw, float32(h), col)
}

// Draw an image as big as possible inside a box, keeping aspect ratio ("contain")
func drawImageContain(dst *ebiten.Image, img *ebiten.Image, x, y, w, h float64) {
	if img == nil {
		return
	}

	iw := float64(img.Bounds().Dx())
	ih := float64(img.Bounds().Dy())
	if iw <= 0 || ih <= 0 {
		return
	}

	sx := w / iw
	sy := h / ih
	scale := sx
	if sy < scale {
		scale = sy
	}

	drawW := iw * scale
	drawH := ih * scale

	ox := x + (w-drawW)/2
	oy := y + (h-drawH)/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(ox, oy)

	dst.DrawImage(img, op)
}

// draw functions
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

	// HUD bar
	if g.state == StatePlaying || g.state == StateGameOver {
		g.drawHUD(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return WindowSizeX, WindowSizeY
}

// UI drawing functions
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
	for _, obj := range g.gameObjects {
		obj.Draw(screen, g)
	}
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

func (g *Game) drawHUD(screen *ebiten.Image) {
	// HUD background
	fillRect(
		screen,
		0,
		float32(HUDY),
		float32(WindowSizeX),
		float32(HUDHeight),
		color.RGBA{10, 10, 16, 220},
	)

	// Layout widths
	leftW := HUDLeftW
	infoW := HUDInfoW
	rightW := HUDRightW
	scoreW := HUDCenterW - infoW

	// Frames
	drawFrame(screen, 0, HUDY, leftW, HUDHeight)
	drawFrame(screen, leftW, HUDY, scoreW, HUDHeight)
	drawFrame(screen, leftW+scoreW, HUDY, infoW, HUDHeight)
	drawFrame(screen, leftW+scoreW+infoW, HUDY, rightW, HUDHeight)

	// LEFT : player icon
	margin := 10.0
	boxX := float64(HUDBorderWidth) + margin
	boxY := float64(HUDY) + float64(HUDBorderWidth) + margin
	boxW := float64(leftW) - float64(2*HUDBorderWidth) - 2*margin
	boxH := float64(HUDHeight) - float64(2*HUDBorderWidth) - 2*margin
	drawImageContain(screen, g.assets.PlayerHUDImage, boxX, boxY, boxW, boxH)

	// WASD image: draw it proportionally inside the right block (contain)
	marginR := 10.0
	rx := float64(HUDLeftW+HUDCenterW) + float64(HUDBorderWidth) + marginR
	ry := float64(HUDY) + float64(HUDBorderWidth) + marginR
	rw := float64(HUDRightW) - float64(2*HUDBorderWidth) - 2*marginR
	rh := float64(HUDHeight) - float64(2*HUDBorderWidth) - 2*marginR

	drawImageContain(screen, g.assets.WASDHUDImage, rx, ry, rw, rh)

	// Score (center block)
	pX := g.getPlayer(PlayerSymbolX)
	pO := g.getPlayer(PlayerSymbolO)
	if pX != nil && pO != nil {
		scoreCenterX := float64(leftW) + float64(scoreW)/2
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
	infoCenterX := float64(leftW+scoreW) + float64(infoW)/2

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
		float64(HUDY+HUDTextPadY+28),
		TopCenter,
		color.RGBA{180, 180, 180, 255},
	)
}

// drawText draws text on the given screen at (x,y) with the specified alignment and color
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
