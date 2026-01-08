// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets/sprites/*.png
var spritesFS embed.FS

type Sprite struct {
	Pos Vec2
	Img *ebiten.Image
}

// LoadSprite returns the sprite image for the given filename.
func LoadSprite(filename string) (*ebiten.Image, error) {
	fullPath := filepath.ToSlash(filepath.Join(SpriteFolder, filename))
	f, err := spritesFS.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", fullPath, err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("decode %q: %w", fullPath, err)
	}
	return img, nil
}
