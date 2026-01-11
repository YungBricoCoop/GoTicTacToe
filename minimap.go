// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Minimap struct{}

func drawPlayer(player *Player, screen *ebiten.Image) {

	px := MinimapPosX + player.pos.X*MinimapGridCellSize
	py := MinimapPosY + player.pos.Y*MinimapGridCellSize

	col := ColorMinimapWall
	switch player.symbol {
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

	endX := px + player.dir.X*(MinimapPlayerArrowLength)
	endY := py + player.dir.Y*(MinimapPlayerArrowLength)

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

	// draw each player
	drawPlayer(g.playerX, screen)
	drawPlayer(g.playerO, screen)
}
