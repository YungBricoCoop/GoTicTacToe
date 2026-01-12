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
	TextureID TextureID
	Scale     float64
	Z         float64
	Hidden    bool
}

const SkeletonSkullScale = 0.5
const SkeletonSkullZ = -1.0
const ChainsScale = 0.5
const ChainsZ = 0.5

func createSprites() []*Sprite {
	return []*Sprite{
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 2.2, Y: 4.7},
			TextureID: Light,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 11.5, Y: 1.4},
			TextureID: Light,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 8.4, Y: 14.6},
			TextureID: Light,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 18.2, Y: 18.6},
			TextureID: Light,
			Hidden:    false,
		},

		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 2.9, Y: 5.3},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 10.6, Y: 5.1},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 15.9, Y: 4.4},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 5.2, Y: 11.7},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 12.7, Y: 12.4},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 17.1, Y: 17.3},
			TextureID: SkeletonSkull,
			Scale:     SkeletonSkullScale,
			Z:         SkeletonSkullZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 6.4, Y: 3.2},
			TextureID: Chains,
			Scale:     ChainsScale,
			Z:         ChainsZ,
			Hidden:    false,
		},
		{
			//nolint:mnd // position on the map
			Position:  Vec2{X: 14.8, Y: 10.6},
			TextureID: Chains,
			Scale:     ChainsScale,
			Z:         ChainsZ,
			Hidden:    false,
		},
	}
}
