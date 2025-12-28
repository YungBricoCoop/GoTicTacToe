// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Minimap struct{}

func (m *Minimap) Update(g *Game) {
	// nothing dynamic to update
}

func (m *Minimap) Draw(screen *ebiten.Image, g *Game) {
	const (
		tilePx = 8.0
		pad    = 10.0
		border = 2.0
		r      = 2.0
		lineW  = float32(2.0)

		// arrow tuning
		dirLen    = 3.0  // in tiles
		headLen   = 1.2  // in tiles
		headAngle = 0.55 // radians (~31°)
	)

	mapW := float64(g.worldMap.Width()) * tilePx
	mapH := float64(g.worldMap.Height()) * tilePx

	originX := float64(WindowSizeX) - pad - mapW
	originY := pad

	// minimap background / border
	fillRect(
		screen,
		float32(originX-border), float32(originY-border),
		float32(mapW+border*2), float32(mapH+border*2),
		color.RGBA{0, 0, 0, 255},
	)

	// draw walls
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

	// draw each player + direction arrow
	for _, obj := range g.gameObjects {
		p, ok := obj.(*Player)
		if !ok {
			continue
		}

		// player color
		var col color.RGBA
		if p.symbol == PlayerSymbolX {
			col = color.RGBA{80, 80, 255, 255}
		} else if p.symbol == PlayerSymbolO {
			col = color.RGBA{80, 255, 80, 255}
		} else {
			col = color.RGBA{200, 200, 200, 255}
		}

		// player dot position (in minimap pixels)
		px := originX + p.pos.X*tilePx
		py := originY + p.pos.Y*tilePx

		// draw player dot
		fillRect(
			screen,
			float32(px-r), float32(py-r),
			float32(r*2), float32(r*2),
			col,
		)

		// draw direction arrow
		dir := p.dir.Normalize()
		if dir.Len2() == 0 {
			continue
		}

		endX := px + dir.X*(dirLen*tilePx)
		endY := py + dir.Y*(dirLen*tilePx)

		// main line
		vector.StrokeLine(screen, float32(px), float32(py), float32(endX), float32(endY), lineW, col, true)

		// arrow head: two small lines rotated around the direction
		angle := math.Atan2(dir.Y, dir.X)

		leftA := angle + math.Pi - headAngle
		rightA := angle + math.Pi + headAngle

		hx1 := endX + math.Cos(leftA)*(headLen*tilePx)
		hy1 := endY + math.Sin(leftA)*(headLen*tilePx)
		hx2 := endX + math.Cos(rightA)*(headLen*tilePx)
		hy2 := endY + math.Sin(rightA)*(headLen*tilePx)

		vector.StrokeLine(screen, float32(endX), float32(endY), float32(hx1), float32(hy1), lineW, col, true)
		vector.StrokeLine(screen, float32(endX), float32(endY), float32(hx2), float32(hy2), lineW, col, true)
	}
}
