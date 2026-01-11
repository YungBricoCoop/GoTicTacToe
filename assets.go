// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
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

// loadAssets loads all game resources. After this returns, Assets should be treated as read only.
func loadAssets() (*Assets, error) {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
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
