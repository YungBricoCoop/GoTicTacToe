// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "math"

// RayHit represents the result of casting a ray in the raycasting engine.
// hit indicates if a wall was hit.
// cellX and cellY are the grid coordinates of the hit cell.
// distance is the distance from the ray origin to the hit point.
// wallX is the exact position along the wall where the ray hit (between 0 and 1).
// side indicates whether a vertical (0) or horizontal (1) wall was hit.
type RayHit struct {
	hit      bool
	cellX    int
	cellY    int
	distance float64
	wallX    float64
	side     int
}

// GetK returns the camera plane coefficient based on the player's field of view.
func GetK(playerFOV float64) float64 {
	//nolint:mnd // dividing the FOV by 2 is part of the standard projection plane computation
	return math.Tan(playerFOV / 2)
}

// GetCameraX returns the camera X coordinate (between -1 and 1) for a given screen X coordinate.
func GetCameraX(screenX int, screenWidth int) float64 {
	return 2*float64(screenX)/float64(screenWidth-1) - 1
}

// GetRayDirection returns the direction vector of the ray based on the player's direction,
// the camera plane coefficient k, and the camera X coordinate.
func GetRayDirection(playerDirection Vec2, k float64, cameraX float64) Vec2 {
	playerDirectionPerp := playerDirection.Perp()
	plane := playerDirectionPerp.Scale(k)
	return playerDirection.Add(plane.Scale(cameraX))
}

// CastRay runs the DDA algorithm to cast a ray from the player position in the given direction.
// It returns a RayHit containing information about the hit.
func CastRay(
	playerPosition Vec2,
	rayDirection Vec2,
	grid [][]uint8,
	maxIterations int,
) RayHit {
	// return no hit if the grid is empty
	if isGridEmpty(grid) {
		return noHit()
	}

	gridWidth := len(grid[0])
	gridHeight := len(grid)

	// grid cell the player is currently standing on
	mapX := int(playerPosition.X)
	mapY := int(playerPosition.Y)

	// deltaDistX tells us how far along the ray we must travel to cross one vertical grid line
	// deltaDistY does the same for horizontal grid lines
	deltaDistX := math.Inf(1)
	if rayDirection.X != 0 {
		deltaDistX = math.Abs(1 / rayDirection.X)
	}

	deltaDistY := math.Inf(1)
	if rayDirection.Y != 0 {
		deltaDistY = math.Abs(1 / rayDirection.Y)
	}

	var stepX, stepY int
	var sideDistX, sideDistY float64

	// stepX define if the ray moves left or right through the grid (1=right, -1=left)
	// sideDistX is the distance from the player to the first vertical grid line
	if rayDirection.X < 0 {
		stepX = -1
		sideDistX = (playerPosition.X - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - playerPosition.X) * deltaDistX
	}

	// stepY define if the ray moves up or down through the grid (1=down, -1=up)
	// sideDistY is the distance from the player to the first horizontal grid line
	if rayDirection.Y < 0 {
		stepY = -1
		sideDistY = (playerPosition.Y - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - playerPosition.Y) * deltaDistY
	}

	hit := false
	side := 0 // 0=vertical wall hit, 1=horizontal wall hit

	// dda loop
	// each iteration moves the ray to the next closest grid cell
	for range maxIterations {
		// define if the ray hits the next vertical or horizontal grid line first
		if sideDistX < sideDistY {
			// cross a vertical grid line
			sideDistX += deltaDistX
			mapX += stepX
			side = 0
		} else {
			// cross a horizontal grid line
			sideDistY += deltaDistY
			mapY += stepY
			side = 1
		}

		// stop if the ray leaves the map
		if isRayOutOfBounds(mapX, mapY, gridWidth, gridHeight) {
			hit = false
			break
		}

		// if the new cell contains a wall, we have a hit
		if isGridCellNotEmpty(grid, mapX, mapY) {
			hit = true
			break
		}
	}

	if !hit {
		// return an infinite distance so callers can safely use it in a zBuffer
		return noHit()
	}

	// get the distance from the player to the wall
	// we subtract delta step because the last step already crossed the wall
	var distance float64
	if side == 0 {
		distance = sideDistX - deltaDistX
	} else {
		distance = sideDistY - deltaDistY
	}

	// store which grid cell was hit
	hitCellX := mapX
	hitCellY := mapY

	// compute the hit position along the wall (between 0 and 1)
	// for a vertical wall we use the y coordinate at the hit point
	// for a horizontal wall we use the x coordinate at the hit point
	var wallX float64
	if side == 0 {
		wallX = playerPosition.Y + distance*rayDirection.Y
	} else {
		wallX = playerPosition.X + distance*rayDirection.X
	}
	wallX -= math.Floor(wallX)

	return RayHit{
		hit:      true,
		cellX:    hitCellX,
		cellY:    hitCellY,
		distance: distance,
		wallX:    wallX,
		side:     side,
	}
}

// noHit returns a RayHit indicating that no wall was hit.
func noHit() RayHit {
	return RayHit{hit: false, distance: math.Inf(1)}
}

// isRayOutOfBounds returns true if the ray's grid coordinates are outside the map boundaries.
func isRayOutOfBounds(mapX, mapY int, gridWidth, gridHeight int) bool {
	return mapX < 0 || mapY < 0 || mapY >= gridHeight || mapX >= gridWidth
}

// isGridEmpty returns true if the grid is empty, if it has zero width or height.
func isGridEmpty(grid [][]uint8) bool {
	return len(grid) == 0 || len(grid[0]) == 0
}

// isGridCellNotEmpty returns true if the grid cell at (x, y) is not empty.
// The cell is considered not empty if its value is greater than zero.
func isGridCellNotEmpty(grid [][]uint8, x, y int) bool {
	return grid[y][x] > 0
}
