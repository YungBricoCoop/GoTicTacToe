// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "image/color"

const (
	WindowSizeX = 1280
	WindowSizeY = 720
	TPS         = 60
	DeltaTime   = 1.0 / TPS

	GridSize        = 3
	MapGridSize     = 22
	Margin          = 10
	LineWidth       = 2
	HeaderY         = 20
	BottomY         = WindowSizeY - 10
	TextLineSpacing = 5

	BigTextLineSpacing = 20
	DefaultFontSize    = 15
	BigFontSize        = 100

	NameInputX          = 10
	NameInputY          = 40
	NameInputLineHeight = 40

	// MinimapGridCellSize is the minimap cell size in pixels.
	MinimapGridCellSize      = 8
	MinimapWidth             = MapGridSize * MinimapGridCellSize
	MinimapHeight            = MapGridSize * MinimapGridCellSize
	MinimapPadding           = 10
	MinimapBorderWidth       = 2
	MinimapPosX              = WindowSizeX - MinimapWidth - MinimapPadding - MinimapBorderWidth
	MinimapPosY              = MinimapPadding
	MinimapPlayerRadius      = 2
	MinimapPlayerDiameter    = MinimapPlayerRadius * 2
	MinimapWallValue         = 1
	MinimapPlayerArrowAAngle = 0.5
	MinimapPlayerArrowLength = 20
	MinimapPlayerArrowWidth  = 2

	PlayerFOV = 1.58

	DefaultPlayerXSpawnX = 11.5
	DefaultPlayerXSpawnY = 11.5
	DefaultPlayerOSpawnX = 15.5
	DefaultPlayerOSpawnY = 15.5

	TextureSize   = 64
	TextureFolder = "assets/textures"
)

func defaultPlayerXSpawn() Vec2 {
	return Vec2{X: DefaultPlayerXSpawnX, Y: DefaultPlayerXSpawnY}
}

func defaultPlayerOSpawn() Vec2 {
	return Vec2{X: DefaultPlayerOSpawnX, Y: DefaultPlayerOSpawnY}
}

//nolint:gochecknoglobals // colors
var (
	ColorBackground = color.RGBA{30, 30, 30, 255}

	ColorMinimapBorder = color.RGBA{0, 0, 0, 255}
	ColorMinimapWall   = color.RGBA{200, 200, 200, 230}

	ColorMinimapPlayerX = color.RGBA{80, 80, 255, 255}
	ColorMinimapPlayerO = color.RGBA{80, 255, 80, 255}
)
