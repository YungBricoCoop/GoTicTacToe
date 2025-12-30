// déclaraer une fois k
// cameraX est endehors de la boucle for
// les couleurs dans le fichier constant

// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct{}

func (w *World) Update(_ *Game) {
	//  nothing dynamic to update
}

// drawRaycastView renders a classic raycaster view (vertical wall slices)
// into a square area on the left side of the screen.

func (w *World) Draw(screen *ebiten.Image, g *Game) {
	// Render area = left square (720x720) to match WindowSizeY
	viewX := float32(0)
	viewY := float32(0)
	viewW := WindowSizeY
	viewH := WindowSizeY

	fillRect(screen, viewX, viewY, float32(viewW), float32(viewH/2), ColorCeil)
	fillRect(screen, viewX, viewY+float32(viewH/2), float32(viewW), float32(viewH/2), ColorFloor)

	p := g.currentPlayer
	if p == nil {
		return
	}

	k := GetK(PlayerFOV) // PlayerFOV is radians
	maxIter := MapGridSize * MapGridSize

	for x := 0; x < viewW; x++ {
		cameraX := GetCameraX(x, viewW)
		rayDir := GetRayDirection(p.dir, k, cameraX)

		hit := CastRay(p.pos, rayDir, g.worldMap.Tiles, maxIter)
		if !hit.hit || math.IsInf(hit.distance, 1) || hit.distance <= 0 {
			continue
		}

		// Classic height = screenHeight / distance
		lineH := float64(viewH) / hit.distance

		drawStart := float64(viewH)/2 - lineH/2
		drawEnd := float64(viewH)/2 + lineH/2

		if drawStart < 0 {
			drawStart = 0
		}
		if drawEnd > float64(viewH) {
			drawEnd = float64(viewH)
		}

		// basic shading
		wallCol := color.RGBA{180, 180, 180, 255}
		if hit.side == 1 {
			wallCol = color.RGBA{130, 130, 130, 255}
		}

		fillRect(
			screen,
			viewX+float32(x),
			viewY+float32(drawStart),
			1,
			float32(drawEnd-drawStart),
			wallCol,
		)
	}
}
