// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Assets holds all game resources like fonts and images.
// NormalTextFace is the standard font face used for most text rendering.
// BigTextFace is a larger font face used for titles and important messages.
// Textures maps tile types to their corresponding images.
// PlayerImg holds the images for each player character.
// SymbolImg holds the images for each symbol used in the game board.
type Assets struct {
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace

	Textures map[uint8][]*ebiten.Image

	PlayerImg map[PlayerSymbol]*ebiten.Image
	SymbolImg map[PlayerSymbol]*ebiten.Image
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

	xPlayer, err := LoadSprite("x-player.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite x-player.png: %w", err)
	}
	oPlayer, err := LoadSprite("o-player.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite o-player.png: %w", err)
	}

	xSymbol, err := LoadSprite("x.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite x.png: %w", err)
	}
	oSymbol, err := LoadSprite("o.png")
	if err != nil {
		return nil, fmt.Errorf("load sprite o.png: %w", err)
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

		PlayerImg: map[PlayerSymbol]*ebiten.Image{
			PlayerSymbolX: xPlayer,
			PlayerSymbolO: oPlayer,
		},
		SymbolImg: map[PlayerSymbol]*ebiten.Image{
			PlayerSymbolX: xSymbol,
			PlayerSymbolO: oSymbol,
		},
	}, nil
}
