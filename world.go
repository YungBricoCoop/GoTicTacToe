// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type World struct {
	zBuffer  []float64
	cameraX  []float64
	fovScale float64
}

func (w *World) Draw(screen *ebiten.Image, g *Game) {
	viewX := float64(0)
	viewY := float64(0)

	drawCeiling(screen, viewX, viewY, WindowSizeX, WindowSizeYDiv2)
	drawFloor(screen, viewX, viewY+float64(WindowSizeY)/HalfDivisor, WindowSizeX, WindowSizeYDiv2)

	p := g.currentPlayer
	if p == nil {
		return
	}

	// initialize buffers once
	w.ensureBuffers()

	w.rayCastAndDrawWalls(screen, g, viewX, viewY)

	spritesToDraw := w.gatherSprites(g)

	drawSprites(screen, w.zBuffer, spritesToDraw, p, w.fovScale, viewX, viewY, WindowSizeX, WindowSizeY)
}

func (w *World) ensureBuffers() {
	if w.zBuffer == nil || len(w.zBuffer) != WindowSizeX {
		w.zBuffer = make([]float64, WindowSizeX)
	}
	if w.cameraX == nil || len(w.cameraX) != WindowSizeX {
		w.cameraX = make([]float64, WindowSizeX)
		for i := range WindowSizeX {
			w.cameraX[i] = GetCameraX(i, WindowSizeX)
		}
	}
}

func (w *World) rayCastAndDrawWalls(screen *ebiten.Image, g *Game, viewX, viewY float64) {
	p := g.currentPlayer
	for x := range WindowSizeX {
		rayDir := GetRayDirection(p.dir, w.fovScale, w.cameraX[x])

		hit := CastRay(p.pos, rayDir, g.worldMap.Tiles, MaxRayIter)
		if !hit.hit || math.IsInf(hit.distance, 1) || hit.distance <= 0 {
			w.zBuffer[x] = math.Inf(1)
			continue
		}

		w.zBuffer[x] = hit.distance

		// get the texture id
		textureID, outOfBounds := g.worldMap.GetTile(hit.cellX, hit.cellY)
		if outOfBounds {
			continue
		}
		textureSlices := g.assets.Textures[textureID]
		textureSliceIndex := int(hit.wallX * float64(TextureSize))

		// clamp textureSliceIndex to valid range
		if textureSliceIndex < 0 {
			textureSliceIndex = 0
		} else if textureSliceIndex >= len(textureSlices) {
			textureSliceIndex = len(textureSlices) - 1
		}

		texture := textureSlices[textureSliceIndex]

		// classic height = screenHeight / distance
		lineH := float64(WindowSizeY) / hit.distance

		drawStart := float64(WindowSizeYDiv2) - lineH/HalfDivisor

		// display the texture slice
		op := &ebiten.DrawImageOptions{}
		scaleY := lineH / float64(texture.Bounds().Dy())

		// put shading based on the distance
		//TODO: use constants
		colorScale := float32(1.0) / (1.0 + float32(hit.distance)/10.0)
		op.ColorScale.Scale(colorScale, colorScale, colorScale, 1)

		op.GeoM.Scale(1, scaleY)
		op.GeoM.Translate(viewX+float64(x), viewY+drawStart)
		screen.DrawImage(texture, op)
	}
}

func (w *World) gatherSprites(g *Game) []Sprite {
	sprites := w.gatherPlayerSprites(g)
	sprites = append(sprites, w.gatherBoardSprites(g)...)
	return sprites
}

func (w *World) gatherPlayerSprites(g *Game) []Sprite {
	var sprites []Sprite
	p := g.currentPlayer
	// 1. Other players
	otherPlayer := g.playerX
	if p.symbol != PlayerSymbolX {
		otherPlayer = g.playerO
	}
	// Find the sprite for this player
	if s, found := g.assets.PlayerImg[otherPlayer.symbol]; found {
		if s != nil {
			sprites = append(sprites, Sprite{
				Pos: otherPlayer.pos,
				Img: s,
			})
		}
	}
	return sprites
}

func (w *World) gatherBoardSprites(g *Game) []Sprite {
	var sprites []Sprite
	// 2. Board symbols (X and O placed on the grid)
	for y := range GridSize {
		for x := range GridSize {
			sym := g.board[y][x]

			img := g.assets.SymbolImg[sym]

			if img != nil {
				// Calculate world position for the symbol
				// Center of the room
				posX := float64(x)*MapRoomStride + MapRoomOffset
				posY := float64(y)*MapRoomStride + MapRoomOffset
				sprites = append(sprites, Sprite{
					Pos: Vec2{X: posX, Y: posY},
					Img: img,
				})
			}
		}
	}
	return sprites
}

func drawCeiling(screen *ebiten.Image, viewX, viewY float64, width, viewH int) {
	vector.FillRect(screen, float32(viewX), float32(viewY), float32(width), float32(viewH), ColorCeil, false)
}

func drawFloor(screen *ebiten.Image, viewX, viewY float64, width, viewH int) {
	vector.FillRect(screen, float32(viewX), float32(viewY), float32(width), float32(viewH), ColorFloor, false)
}

func drawSprites(
	screen *ebiten.Image,
	zBuffer []float64,
	sprites []Sprite,
	player *Player,
	fovScale float64,
	viewportX, viewportY float64,
	viewportWidth, viewportHeight int,
) {
	if player == nil || len(sprites) == 0 || len(zBuffer) != viewportWidth {
		return
	}

	// calculate the camera plane vector
	plane := player.dir.Perp().Scale(fovScale)

	// sort sprites from far to near to handle transparency correctly
	sortedSprites := make([]Sprite, len(sprites))
	copy(sortedSprites, sprites)

	slices.SortFunc(sortedSprites, func(a, b Sprite) int {
		distA := a.Pos.Sub(player.pos).Len2()
		distB := b.Pos.Sub(player.pos).Len2()
		if distA > distB {
			return -1
		}
		if distA < distB {
			return 1
		}
		return 0
	})

	for _, sprite := range sortedSprites {
		if sprite.Img == nil {
			continue
		}
		drawSingleSprite(screen, zBuffer, sprite, player, plane, viewportX, viewportY, viewportWidth, viewportHeight)
	}
}

func drawSingleSprite(
	screen *ebiten.Image,
	zBuffer []float64,
	sprite Sprite,
	player *Player,
	plane Vec2,
	viewportX, viewportY float64,
	viewportWidth, viewportHeight int,
) {
	// translate sprite position relative to camera
	relativePos := sprite.Pos.Sub(player.pos)

	// transform sprite with the inverse camera matrix
	inverseDeterminant := 1.0 / (plane.X*player.dir.Y - player.dir.X*plane.Y)

	transformX := inverseDeterminant * (player.dir.Y*relativePos.X - player.dir.X*relativePos.Y)
	transformY := inverseDeterminant * (-plane.Y*relativePos.X + plane.X*relativePos.Y) // this is actually the depth inside the screen

	if transformY <= 0 {
		return // sprite is behind the camera
	}

	spriteScreenX := int((float64(viewportWidth) / HalfDivisor) * (1 + transformX/transformY))

	// calculate sprite dimensions on screen, keeping aspect ratio square
	spriteHeight := math.Abs(float64(viewportHeight) / transformY)
	spriteWidth := spriteHeight
	if spriteHeight < 1 || spriteWidth < 1 {
		return
	}

	// calculate drawing bounds on screen
	spriteTopY := -spriteHeight/HalfDivisor + float64(viewportHeight)/HalfDivisor

	drawStartX := int(-spriteWidth/HalfDivisor + float64(spriteScreenX))
	drawEndX := int(spriteWidth/HalfDivisor + float64(spriteScreenX))

	imageBounds := sprite.Img.Bounds()
	imageWidth := imageBounds.Dx()
	imageHeight := imageBounds.Dy()
	if imageWidth <= 0 || imageHeight <= 0 {
		return
	}

	// loop through every vertical stripe of the sprite on screen
	for screenColumn := drawStartX; screenColumn < drawEndX; screenColumn++ {
		if screenColumn < 0 || screenColumn >= viewportWidth {
			continue
		}
		// check z-buffer to see if sprite is visible (not hidden by wall)
		if transformY >= zBuffer[screenColumn] {
			continue
		}

		textureX := int(float64(screenColumn-drawStartX) * float64(imageWidth) / (spriteWidth))
		if textureX < 0 || textureX >= imageWidth {
			continue
		}

		// draw the vertical slice of the sprite
		subImage, ok := sprite.Img.SubImage(image.Rect(textureX, 0, textureX+1, imageHeight)).(*ebiten.Image)
		if !ok {
			continue
		}

		drawOptions := &ebiten.DrawImageOptions{}
		scaleY := spriteHeight / float64(imageHeight)
		drawOptions.GeoM.Scale(1, scaleY)
		drawOptions.GeoM.Translate(viewportX+float64(screenColumn), viewportY+spriteTopY)
		screen.DrawImage(subImage, drawOptions)
	}
}
