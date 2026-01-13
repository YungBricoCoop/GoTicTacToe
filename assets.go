// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Assets holds all game resources like fonts and images.
// NormalTextFace is the standard font face used for most text rendering.
// BigTextFace is a larger font face used for titles and important messages.
// Textures holds all loaded textures mapped by their TextureId.
type Assets struct {
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace
	Textures       TextureMap
}

//go:embed assets/fonts/*.ttf
var fontsFS embed.FS

// loadAssets loads all game resources. After this returns, Assets should be treated as read only.
func loadAssets() (*Assets, error) {
	fontBytes, err := fontsFS.ReadFile("assets/fonts/PressStart2P-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
	}

	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return nil, fmt.Errorf("create font source: %w", err)
	}

	textures, err := LoadTextures()
	if err != nil {
		return nil, fmt.Errorf("load textures: %w", err)
	}

	return &Assets{
		NormalTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   DefaultFontSize,
		},
		BigTextFace: &text.GoTextFace{
			Source: fontSource,
			Size:   BigFontSize,
		},

		Textures: textures,
	}, nil
}
