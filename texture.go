// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"fmt"
	"image"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets/textures/*.png
var texturesFS embed.FS

// TextureStrips represents an array of vertical strips of a texture image.
type TextureStrips [TextureSize]*ebiten.Image

// Texture represents a texture with its source image and vertical strips.
type Texture struct {
	Source *ebiten.Image
	Strips TextureStrips
}

// TextureMap maps texture IDs to their corresponding Texture.
type TextureMap map[TextureID]Texture

// TextureID represents the ID of a texture.
// 1-127 are reserved for wall textures.
// 128-255 are reserved for sprite textures.
type TextureID uint8

const (
	// WallBrick represents a standard brick wall texture.
	WallBrick       TextureID = 1
	WallBrickHole   TextureID = 2
	WallBrickGopher TextureID = 3

	// PlayerXSymbol represents the symbol for Player X.
	PlayerXSymbol    TextureID = 128
	PlayerXCharacter TextureID = 129
	PlayerOSymbol    TextureID = 130
	PlayerOCharacter TextureID = 131
	SkeletonSkull    TextureID = 132
	Chains           TextureID = 133
	Light            TextureID = 134
	WasdKeys         TextureID = 135
)

// LoadTextures loads all textures defined in imageManifest.
// Each file is decoded once into Source, and Strips are derived for raycasting.
func LoadTextures() (TextureMap, error) {
	out := make(TextureMap, len(imageManifest))

	for id, filename := range imageManifest {
		fullPath := filepath.ToSlash(filepath.Join(TextureFolder, filename))

		f, err := texturesFS.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", fullPath, err)
		}

		img, _, err := ebitenutil.NewImageFromReader(f)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("decode %q: %w", fullPath, err)
		}

		strips, err := sliceIntoVerticalStrips(img)
		if err != nil {
			return nil, fmt.Errorf("slice %q: %w", fullPath, err)
		}

		out[id] = Texture{
			Source: img,
			Strips: strips,
		}
	}

	return out, nil
}

func LoadHUDImage(name string) (*ebiten.Image, error) {
	fullPath := filepath.ToSlash(filepath.Join(TextureFolder, name))
	f, errO := texturesFS.Open(fullPath)
	if errO != nil {
		return nil, fmt.Errorf("open HUD image %q: %w", fullPath, errO)
	}

	img, _, errI := ebitenutil.NewImageFromReader(f)
	_ = f.Close()
	if errI != nil {
		return nil, fmt.Errorf("decode HUD image %q: %w", fullPath, errI)
	}
	return img, nil
}

// sliceIntoVerticalStrips returns TextureSize images of width 1.
// The source texture must be exactly TextureSize pixels wide.
func sliceIntoVerticalStrips(src *ebiten.Image) (TextureStrips, error) {
	// get source image bounds
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// enforce fixed texture width for raycasting
	if width != TextureSize {
		return TextureStrips{}, fmt.Errorf("expected texture width %d, got %d", TextureSize, width)
	}

	// height must be valid
	if height <= 0 {
		return TextureStrips{}, fmt.Errorf("invalid texture height %d", height)
	}

	// array of 1 pixel wide vertical texture strips
	var strips TextureStrips

	for x := range TextureSize {
		// take a 1 pixel wide slice from the source texture
		sub, ok := src.SubImage(image.Rect(x, 0, x+1, height)).(*ebiten.Image)
		if !ok || sub == nil {
			return TextureStrips{}, fmt.Errorf("failed to slice texture column x=%d", x)
		}

		// store the strip directly, subimage shares the same underlying texture
		strips[x] = sub
	}

	return strips, nil
}
