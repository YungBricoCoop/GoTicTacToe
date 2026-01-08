// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	zBuffer  []float64
	cameraX  []float64
	fovScale float64
}

func (w *World) Update(_ *Game) {
	//  nothing dynamic to update
}

func (w *World) Draw(screen *ebiten.Image, g *Game) {
	viewX := float64(0)
	viewY := float64(0)

	drawCeiling(screen, viewX, viewY, WindowSizeX, WindowSizeYDiv2)
	drawFloor(screen, viewX, viewY+float64(WindowSizeYDiv2), WindowSizeX, WindowSizeYDiv2)

	p := g.currentPlayer
	if p == nil {
		return
	}

	// initialize buffers once
	if w.zBuffer == nil || len(w.zBuffer) != WindowSizeX {
		w.zBuffer = make([]float64, WindowSizeX)
	}
	if w.cameraX == nil || len(w.cameraX) != WindowSizeX {
		w.cameraX = make([]float64, WindowSizeX)
		for i := 0; i < WindowSizeX; i++ {
			w.cameraX[i] = GetCameraX(i, WindowSizeX)
		}
	}

	for x := 0; x < WindowSizeX; x++ {
		rayDir := GetRayDirection(p.dir, w.fovScale, w.cameraX[x])

		hit := CastRay(p.pos, rayDir, g.worldMap.Tiles, MaxRayIter)
		if !hit.hit || math.IsInf(hit.distance, 1) || hit.distance <= 0 {
			w.zBuffer[x] = math.Inf(1)
			continue
		}

		w.zBuffer[x] = hit.distance

		// get the texture id
		textureID := g.worldMap.GetTile(hit.cellX, hit.cellY)
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

		drawStart := WindowSizeYDiv2 - lineH/2
		drawEnd := WindowSizeYDiv2 + lineH/2

		if drawStart < 0 {
			drawStart = 0
		}
		if drawEnd > float64(WindowSizeY) {
			drawEnd = float64(WindowSizeY)
		}

		//TODO: add shading

		// display the texture slice
		op := &ebiten.DrawImageOptions{}
		scaleY := (drawEnd - drawStart) / float64(texture.Bounds().Dy())
		op.GeoM.Scale(1, scaleY)
		op.GeoM.Translate(viewX+float64(x), viewY+drawStart)
		screen.DrawImage(texture, op)
	}

	drawSprites(screen, w.zBuffer, g.assets, p, w.fovScale, viewX, viewY, WindowSizeX, WindowSizeY)
}

func drawCeiling(screen *ebiten.Image, viewX, viewY float64, WindowSizeX, viewH int) {
	fillRect(screen, float32(viewX), float32(viewY), float32(WindowSizeX), float32(viewH), ColorCeil)
}

func drawFloor(screen *ebiten.Image, viewX, viewY float64, WindowSizeX, viewH int) {
	fillRect(screen, float32(viewX), float32(viewY), float32(WindowSizeX), float32(viewH), ColorFloor)
}

func drawSprites(
	screen *ebiten.Image,
	zBuffer []float64,
	assets *Assets,
	player *Player,
	fovScale float64,
	viewportX, viewportY float64,
	viewportWidth, viewportHeight int,
) {
	if player == nil || assets == nil || len(assets.Sprites) == 0 || len(zBuffer) != viewportWidth {
		return
	}

	// calculate the camera plane vector
	plane := player.dir.Perp().Scale(fovScale)

	// sort sprites from far to near to handle transparency correctly
	sortedSprites := make([]Sprite, 0, len(assets.Sprites))
	for _, sprite := range assets.Sprites {
		sortedSprites = append(sortedSprites, sprite)
	}
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

		// translate sprite position relative to camera
		relativePos := sprite.Pos.Sub(player.pos)

		// transform sprite with the inverse camera matrix
		// [ planeX   dirX ] -1                                       [ dirY      -dirX ]
		// [               ]       =  1/(planeX*dirY-dirX*planeY) *   [                 ]
		// [ planeY   dirY ]                                          [ -planeY  planeX ]
		inverseDeterminant := 1.0 / (plane.X*player.dir.Y - player.dir.X*plane.Y)

		transformX := inverseDeterminant * (player.dir.Y*relativePos.X - player.dir.X*relativePos.Y)
		transformY := inverseDeterminant * (-plane.Y*relativePos.X + plane.X*relativePos.Y) // this is actually the depth inside the screen

		if transformY <= 0 {
			continue // sprite is behind the camera
		}

		spriteScreenX := int((float64(viewportWidth) / 2) * (1 + transformX/transformY))

		// calculate sprite dimensions on screen, keeping aspect ratio square
		spriteHeight := math.Abs(float64(viewportHeight) / transformY)
		spriteWidth := spriteHeight
		if spriteHeight < 1 || spriteWidth < 1 {
			continue
		}

		// calculate drawing bounds on screen
		drawStartY := int(-spriteHeight/2 + float64(viewportHeight)/2)
		drawEndY := int(spriteHeight/2 + float64(viewportHeight)/2)
		if drawStartY < 0 {
			drawStartY = 0
		}
		if drawEndY >= viewportHeight {
			drawEndY = viewportHeight - 1
		}

		drawStartX := int(-spriteWidth/2 + float64(spriteScreenX))
		drawEndX := int(spriteWidth/2 + float64(spriteScreenX))

		imageBounds := sprite.Img.Bounds()
		imageWidth := imageBounds.Dx()
		imageHeight := imageBounds.Dy()
		if imageWidth <= 0 || imageHeight <= 0 {
			continue
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
			scaleY := float64(drawEndY-drawStartY) / float64(imageHeight)
			drawOptions.GeoM.Scale(1, scaleY)
			drawOptions.GeoM.Translate(viewportX+float64(screenColumn), viewportY+float64(drawStartY))
			screen.DrawImage(subImage, drawOptions)
		}
	}
}
