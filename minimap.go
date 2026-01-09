// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Minimap struct{}

func (m *Minimap) Update(_ *Game) {
	//  nothing dynamic to update
}

func (m *Minimap) Draw(screen *ebiten.Image, g *Game) {
	vector.FillRect(
		screen,
		float32(MinimapPosX),
		float32(MinimapPosY-MinimapBorderWidth),
		float32(MinimapWidth+2*MinimapBorderWidth),
		float32(MinimapHeight+2*MinimapBorderWidth),
		ColorMinimapBorder,
		false,
	)

	mapHCells := g.worldMap.Height()
	mapWCells := g.worldMap.Width()

	for y := range mapHCells {
		for x := range mapWCells {
			if g.worldMap.Tiles[y][x] >= MinimapWallValue {
				vector.FillRect(
					screen,
					float32(MinimapPosX+float64(x)*MinimapGridCellSize),
					float32(MinimapPosY+float64(y)*MinimapGridCellSize),
					float32(MinimapGridCellSize), float32(MinimapGridCellSize),
					ColorMinimapWall,
					false,
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

		px := MinimapPosX + p.pos.X*MinimapGridCellSize
		py := MinimapPosY + p.pos.Y*MinimapGridCellSize

		col := ColorMinimapWall
		switch p.symbol {
		case PlayerSymbolNone:
			continue
		case PlayerSymbolX:
			col = ColorMinimapPlayerX
		case PlayerSymbolO:
			col = ColorMinimapPlayerO
		}

		vector.FillRect(
			screen,
			float32(px-MinimapPlayerRadius),
			float32(py-MinimapPlayerRadius),
			float32(MinimapPlayerDiameter),
			float32(MinimapPlayerDiameter),
			col,
			false,
		)

		endX := px + p.dir.X*(MinimapPlayerArrowLength)
		endY := py + p.dir.Y*(MinimapPlayerArrowLength)

		// main line
		vector.StrokeLine(
			screen,
			float32(px),
			float32(py),
			float32(endX),
			float32(endY),
			MinimapPlayerArrowWidth,
			col,
			true,
		)
	}
}
