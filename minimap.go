// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Minimap struct{}

// drawMinimap renders the minimap at the top-right corner of the screen

func (m *Minimap) Update() {
	//  nothing dynamic to update
}

func (m *Minimap) Draw(screen *ebiten.Image, g *Game) {

	const tilePx = 8.0
	const pad = 10.0
	const border = 2.0
	const r = 2.0

	mapW := float64(g.worldMap.Width()) * tilePx
	mapH := float64(g.worldMap.Height()) * tilePx

	originX := float64(ScreenSize) - pad - mapW
	originY := pad

	fillRect(
		screen,
		float32(originX-border), float32(originY-border),
		float32(mapW+border*2), float32(mapH+border*2),
		color.RGBA{0, 0, 0, 255},
	)

	for y := 0; y < g.worldMap.Height(); y++ {
		for x := 0; x < g.worldMap.Width(); x++ {
			if g.worldMap.Tiles[y][x] == 1 {
				fillRect(
					screen,
					float32(originX+float64(x)*tilePx),
					float32(originY+float64(y)*tilePx),
					float32(tilePx), float32(tilePx),
					color.RGBA{200, 200, 200, 230},
				)
			}
		}
	}

	px := originX + g.playerPos.X*tilePx
	py := originY + g.playerPos.Y*tilePx
	fillRect(
		screen,
		float32(px-r), float32(py-r),
		float32(r*2), float32(r*2),
		color.RGBA{255, 80, 80, 255},
	)
}
