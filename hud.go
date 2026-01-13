// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Hud struct{}

// drawPanelFrame draws a rectangular border for one hud panel starting at x with the given width.
func drawPanelFrame(destination *ebiten.Image, x, width int) {
	y := HudTopLeftYPixels
	height := HudHeightPixels
	borderWidth := float32(HudBorderWidthPixels)
	borderColor := ColorHUDBorder

	vector.FillRect(destination, float32(x), float32(y), float32(width), borderWidth, borderColor, false)
	vector.FillRect(
		destination,
		float32(x),
		float32(y+height)-borderWidth,
		float32(width),
		borderWidth,
		borderColor,
		false,
	)
	vector.FillRect(destination, float32(x), float32(y), borderWidth, float32(height), borderColor, false)
	vector.FillRect(
		destination,
		float32(x+width)-borderWidth,
		float32(y),
		borderWidth,
		float32(height),
		borderColor,
		false,
	)
}

// drawImageContained scales and centers image inside the given rectangle, keeping aspect ratio and applying padding.
func drawImageContained(destination *ebiten.Image, image *ebiten.Image, x, y, width, height, padding int) {
	if destination == nil || image == nil {
		return
	}

	innerX := x + padding
	innerY := y + padding
	innerWidth := width - padding*Two
	innerHeight := height - padding*Two
	if innerWidth <= 0 || innerHeight <= 0 {
		return
	}

	bounds := image.Bounds()
	sourceWidth := float64(bounds.Dx())
	sourceHeight := float64(bounds.Dy())
	if sourceWidth == 0 || sourceHeight == 0 {
		return
	}

	scale := math.Min(
		float64(innerWidth)/sourceWidth,
		float64(innerHeight)/sourceHeight,
	)

	drawWidth := sourceWidth * scale
	drawHeight := sourceHeight * scale

	drawX := float64(innerX) + (float64(innerWidth)-drawWidth)/Two
	drawY := float64(innerY) + (float64(innerHeight)-drawHeight)/Two

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(drawX, drawY)
	destination.DrawImage(image, options)
}

// drawTextLines draws multiple lines using a fixed vertical step starting at startY.
func drawTextLines(g *Game, screen *ebiten.Image, x, startY float64, lines []string) {
	if g == nil || screen == nil {
		return
	}

	y := startY
	for _, line := range lines {
		g.drawText(screen, line, x, y, ColorHUDText)
		y += float64(HudTextLineStepPixels)
	}
}

// Draw renders the full hud: player icon, name and scores, action keys, and the wasd image panel.
func (h *Hud) Draw(screen *ebiten.Image, g *Game) {
	if screen == nil || g == nil || g.assets == nil {
		return
	}

	vector.FillRect(
		screen,
		0,
		float32(HudTopLeftYPixels),
		float32(WindowSizeX),
		float32(HudHeightPixels),
		ColorHUDFill,
		false,
	)

	playerPanelX := 0
	playerPanelWidth := HudSquarePanelSizePixels

	namePanelX := playerPanelX + playerPanelWidth
	namePanelWidth := HudNamePanelWidthPixels

	keysPanelX := namePanelX + namePanelWidth
	keysPanelWidth := HudKeysPanelWidthPixels

	wasdPanelX := keysPanelX + keysPanelWidth
	wasdPanelWidth := HudSquarePanelSizePixels

	drawPanelFrame(screen, playerPanelX, playerPanelWidth)
	drawPanelFrame(screen, namePanelX, namePanelWidth)
	drawPanelFrame(screen, keysPanelX, keysPanelWidth)
	drawPanelFrame(screen, wasdPanelX, wasdPanelWidth)

	if g.currentPlayer != nil {
		playerTexture := g.assets.Textures[g.currentPlayer.symbolTextureID]
		drawImageContained(
			screen,
			playerTexture.Source,
			playerPanelX,
			HudTopLeftYPixels,
			playerPanelWidth,
			HudHeightPixels,
			HudImageInnerPaddingPixels,
		)
	}

	nameTextX := float64(namePanelX + HudPanelOuterPaddingXPixels)
	nameTextY := float64(HudTopLeftYPixels + HudPanelOuterPaddingYPixels)

	playerNameLine := "Player: Player"
	if g.currentPlayer != nil && g.currentPlayer.name != "" {
		playerNameLine = "Player: " + g.currentPlayer.name
	}

	currentPlayerScoreLine := "Score: 0"
	if g.currentPlayer != nil {
		currentPlayerScoreLine = "Score: " + strconv.Itoa(g.currentPlayer.score)
	}

	totalScoreLine := "Score  X: 0    O: 0"
	if g.playerX != nil && g.playerO != nil {
		totalScoreLine = "Score  X: " + strconv.Itoa(g.playerX.score) + "    O: " + strconv.Itoa(g.playerO.score)
	}

	drawTextLines(g, screen, nameTextX, nameTextY, []string{
		playerNameLine,
		currentPlayerScoreLine,
		totalScoreLine,
	})

	keysTextX := float64(keysPanelX + HudPanelOuterPaddingXPixels)
	keysTextY := float64(HudTopLeftYPixels + HudPanelOuterPaddingYPixels)

	drawTextLines(g, screen, keysTextX, keysTextY, []string{
		"Esc: Quit",
		"Ctrl + R: Restart",
		"E: Place marker",
	})

	wasdTexture := g.assets.Textures[WasdKeys]
	drawImageContained(
		screen,
		wasdTexture.Source,
		wasdPanelX,
		HudTopLeftYPixels,
		wasdPanelWidth,
		HudHeightPixels,
		HudImageInnerPaddingPixels,
	)
}
