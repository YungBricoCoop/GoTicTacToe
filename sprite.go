// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

// Sprite represents a 2D image in the game world.
// Position is the world coordinates of the sprite's center.
// TextureID is the ID of the texture to use for rendering the sprite.
// Scale is a multiplier for the sprite's size.
// Z is the vertical offset of the sprite in the world.
// Hidden indicates whether the sprite should be rendered or not.
type Sprite struct {
	Position  Vec2
	TextureID TextureId
	Scale     float64
	Z         float64
	Hidden    bool
}
