// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "image/color"

const (
	WindowSizeX     = 1280
	WindowSizeXDiv2 = WindowSizeX / 2
	WindowSizeY     = 720
	WindowSizeYDiv2 = WindowSizeY / 2
	WindowTitle     = "Raycastoe"
	TPS             = 60
	DeltaTime       = 1.0 / TPS

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

	PlayerFOV           = 1.58
	PlayerMovementSpeed = 5.0 // units per second
	PlayerRotationSpeed = 3.0 // radians per second

	DefaultPlayerXSpawnX = 11.5
	DefaultPlayerXSpawnY = 11.5
	DefaultPlayerOSpawnX = 15.5
	DefaultPlayerOSpawnY = 15.5

	TextureSize   = 64
	TextureFolder = "assets/textures"
	SpriteFolder  = "assets/sprites"

	MaxRayIter = MapGridSize * MapGridSize

	MapRoomStride = 7

	InitialSpriteXPosX = 11.5
	InitialSpriteXPosY = 7.0
	InitialSpriteOPosX = 3.0
	InitialSpriteOPosY = 4.0

	HalfDivisor = 2.0
)

func defaultPlayerXSpawn() Vec2 {
	return Vec2{X: DefaultPlayerXSpawnX, Y: DefaultPlayerXSpawnY}
}

func defaultPlayerOSpawn() Vec2 {
	return Vec2{X: DefaultPlayerOSpawnX, Y: DefaultPlayerOSpawnY}
}

//nolint:gochecknoglobals // colors
var (
	ColorBackground = color.RGBA{30, 30, 30, 100}

	ColorMinimapBorder = color.RGBA{0, 0, 0, 100}
	ColorMinimapWall   = color.RGBA{200, 200, 200, 100}

	ColorMinimapPlayerX = color.RGBA{249, 77, 0, 100}
	ColorMinimapPlayerO = color.RGBA{86, 229, 252, 100}

	ColorCeiling = color.RGBA{22, 24, 28, 255}
	ColorFloor   = color.RGBA{17, 15, 15, 255}
)

//nolint:gochecknoglobals // texture manifest
var imageManifest = map[TextureId]string{
	// walls
	WallBrick:       "wall-brick.png",
	WallBrickHole:   "wall-brick-hole.png",
	WallBrickGopher: "wall-brick-gopher.png",

	// sprites
	PlayerXSymbol:    "x.png",
	PlayerXCharacter: "x-player.png",
	PlayerOSymbol:    "o.png",
	PlayerOCharacter: "o-player.png",
	SkeletonSkull:    "skeleton-skull.png",
	Chains:           "chains.png",
	Light:            "lantern.png",
}
