// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/hajimehoshi/ebiten/v2"

type Updatable interface {
	Update(g *Game)
}

type Drawable interface {
	Draw(screen *ebiten.Image, g *Game)
}
