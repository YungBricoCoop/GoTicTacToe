// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameObject represents any object in the game that can be updated and drawn.
// Update is called every frame to update the object's state.
// Draw is called every frame to render the object on the screen.
type GameObject interface {
	Update(g *Game)
	Draw(screen *ebiten.Image, g *Game)
}
