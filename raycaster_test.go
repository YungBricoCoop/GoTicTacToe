// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math"
	"testing"
)

// mock grid for testing.
type MockGrid struct {
	width, height int
	tiles         [][]TileID
}

func (m MockGrid) Width() int  { return m.width }
func (m MockGrid) Height() int { return m.height }
func (m MockGrid) GetTileID(x, y int) (TileID, bool) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return TileEmpty, true
	}
	return m.tiles[y][x], false
}

func TestRayDirection(t *testing.T) {
	// player direction: east
	dir := Vec2{1, 0}
	k := 1.0

	// 1. center of screen
	ray := GetRayDirection(dir, k, 0)
	if ray.X != 1 || ray.Y != 0 {
		t.Errorf("center ray wrong")
	}

	// 2. right side
	ray = GetRayDirection(dir, k, 1)
	if ray.X != 1 || ray.Y != 1 {
		t.Errorf("right ray wrong")
	}
}

func TestRayCast(t *testing.T) {
	// 5x5 grid with one wall at (2,2)
	w, h := 5, 5
	tiles := make([][]TileID, h)
	for y := range h {
		tiles[y] = make([]TileID, w)
	}
	tiles[2][2] = 1

	grid := MockGrid{width: w, height: h, tiles: tiles}

	// test hitting a horizontal wall
	// player is above the wall looking down
	start := Vec2{2.5, 0.5}
	dir := Vec2{0, 1}

	hit := CastRay(start, dir, grid, 100)

	if !hit.hit {
		t.Fatal("should have hit wall")
	}
	if hit.cellX != 2 || hit.cellY != 2 {
		t.Errorf("wrong cell")
	}
	if hit.side != 1 {
		t.Errorf("wrong side")
	}

	// check dist
	// wall at y=2, player at y=0.5 -> dist 1.5
	if math.Abs(hit.distance-1.5) > 0.001 {
		t.Errorf("wrong distance")
	}

	// test hitting vertical wall
	// player is left of wall looking right
	start = Vec2{0.5, 2.5}
	dir = Vec2{1, 0}

	hit = CastRay(start, dir, grid, 100)

	if !hit.hit {
		t.Fatal("should have hit wall")
	}
	if hit.distance != 1.5 {
		t.Errorf("wrong distance calculation")
	}
}
