// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "math"

// GetK returns the camera plane coefficient based on the player's field of view.
func GetK(playerFOV float64) float64 {
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

// CastRay runs the DDA algorithm to find the first wall hit by a ray.
// It returns:
// - hit: true if a wall was hit, false otherwise.
// - hitPos: the exact position (x, y) where the ray hit the wall.
// - hitCellX, hitCellY: the integer grid coordinates of the wall that was hit.
// - distance: the perpendicular distance from the camera plane to the wall (corrected for fisheye effect).
// - side: 0 for vertical wall (x-side), 1 for horizontal wall (y-side).
func CastRay(
	playerPosition Vec2,
	rayDirection Vec2,
	grid [][]int,
	maxIterations int,
) (hit bool, hitPosition Vec2, hitCellX, hitCellY int, distance float64, side int) {

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

	hit = false
	side = 0 // 0=vertical wall hit, 1=horizontal wall hit

	// dda loop
	// each iteration moves the ray to the next closest grid boundary
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
		if mapX < 0 || mapY < 0 || mapY >= len(grid) || mapX >= len(grid[0]) {
			hit = false
			break
		}

		// if the new cell contains a wall, we have a hit
		if grid[mapY][mapX] > 0 {
			hit = true
			break
		}
	}

	if hit {
		// get the distance from the player to the wall
		// we subtract delta step because the last step already crossed the wall
		if side == 0 {
			distance = sideDistX - deltaDistX
		} else {
			distance = sideDistY - deltaDistY
		}

		// compute the world position where the ray hit the wall
		hitPosition = playerPosition.Add(rayDirection.Scale(distance))

		// store which grid cell was hit
		hitCellX = mapX
		hitCellY = mapY
	}

	return hit, hitPosition, hitCellX, hitCellY, distance, side
}
