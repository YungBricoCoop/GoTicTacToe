// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type World struct {
	zBuffer  []float64
	cameraX  []float64
	fovScale float64
}

// Draw renders the world view.
// it draws ceiling and floor, then raycasts each screen column to draw textured walls,
// then draws sprites using the z-buffer for correct occlusion.
func (w *World) Draw(screen *ebiten.Image, g *Game) {
	if screen == nil || g == nil || g.assets == nil {
		return
	}

	drawCeiling(screen)
	drawFloor(screen)

	p := g.currentPlayer
	if p == nil {
		return
	}

	// initialize buffers once
	w.ensureRenderBuffers()

	// walls write to w.zBuffer
	w.raycastColumnsAndDrawWalls(screen, g, p)

	// sprites read w.zBuffer to hide behind walls
	w.drawSprites(screen, g, p)
}

// ensureRenderBuffers allocates and fills per frame constant buffers.
// zBuffer stores the distance per screen column.
// cameraX stores camera x coordinates per screen column (between -1 and 1).
func (w *World) ensureRenderBuffers() {
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

// raycastColumnsAndDrawWalls casts one ray per screen column and draws the corresponding wall slice.
// this is the main raycasting render loop.
func (w *World) raycastColumnsAndDrawWalls(screen *ebiten.Image, g *Game, p *Player) {
	if screen == nil || g == nil || p == nil {
		return
	}

	for x := range WindowSizeX {
		hit, ok := w.castRayForScreenColumn(g, p, x)
		if !ok {
			w.zBuffer[x] = math.Inf(1)
			continue
		}

		w.zBuffer[x] = hit.distance

		strip, ok := w.resolveTextureStripFromHit(g, hit)
		if !ok {
			continue
		}

		lineH := w.wallSliceHeightOnScreen(hit.distance)
		drawStart := w.wallSliceTopY(lineH)

		w.drawTexturedWallSlice(screen, strip, x, drawStart, lineH, hit.distance)
	}
}

// castRayForScreenColumn builds the ray direction for the given screen column and runs the dda cast.
// it returns ok=false if there is no valid hit.
func (w *World) castRayForScreenColumn(g *Game, p *Player, x int) (RayHit, bool) {
	if g == nil || p == nil {
		return RayHit{}, false
	}
	if x < 0 || x >= WindowSizeX {
		return RayHit{}, false
	}

	rayDir := GetRayDirection(p.dir, w.fovScale, w.cameraX[x])

	hit := CastRay(p.pos, rayDir, g.worldMap, MaxRayIter)
	if !hit.hit || math.IsInf(hit.distance, 1) || hit.distance <= 0 {
		return RayHit{}, false
	}

	return hit, true
}

// resolveTextureStripFromHit converts the map hit into a texture strip image.
// it uses the hit cell to read the tile id, maps it to a texture id, then uses wallX to pick the strip.
func (w *World) resolveTextureStripFromHit(g *Game, hit RayHit) (*ebiten.Image, bool) {
	if g == nil || g.assets == nil {
		return nil, false
	}

	// get the texture id
	tileID, outOfBounds := g.worldMap.GetTileID(hit.cellX, hit.cellY)
	if outOfBounds {
		return nil, false
	}

	textureID, ok := tileID.TextureID()
	if !ok {
		return nil, false
	}

	texture := g.assets.Textures[textureID]
	if len(texture.Strips) == 0 {
		return nil, false
	}

	stripIndex := w.textureStripIndexFromWallX(hit.wallX, len(texture.Strips))
	return texture.Strips[stripIndex], true
}

// textureStripIndexFromWallX converts wallX (0..1) into a strip index for the texture.
// it clamps the index to avoid out of range due to floating point rounding.
func (w *World) textureStripIndexFromWallX(wallX float64, stripCount int) int {
	if stripCount <= 1 {
		return 0
	}

	// stripIndex is the x coordinate inside the texture (converted to a strip index)
	stripIndex := int(wallX * float64(stripCount))

	// clamp stripIndex to valid range
	return clampInt(stripIndex, 0, stripCount-1)
}

// wallSliceHeightOnScreen returns the wall slice height in pixels for a given hit distance.
// it uses the classic projection formula: screenHeight / distance.
func (w *World) wallSliceHeightOnScreen(distance float64) float64 {
	if distance <= 0 || math.IsInf(distance, 1) || math.IsNaN(distance) {
		return 0
	}

	// classic height = screenHeight / distance
	return float64(WindowSizeY) / distance
}

// wallSliceTopY returns the y coordinate where the wall slice should start so it is vertically centered.
func (w *World) wallSliceTopY(lineH float64) float64 {
	return float64(WindowSizeYDiv2) - lineH/Two
}

// drawTexturedWallSlice draws one vertical textured strip on screen.
// it scales the strip to the projected wall height, applies distance shading, and draws it at column x.
func (w *World) drawTexturedWallSlice(
	screen *ebiten.Image,
	textureStrip *ebiten.Image,
	x int,
	drawStart float64,
	lineH float64,
	distance float64,
) {
	if screen == nil || textureStrip == nil {
		return
	}
	if x < 0 || x >= WindowSizeX {
		return
	}
	if lineH <= 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// scale the strip to match the wall height on screen
	scaleY := lineH / float64(TextureSize)
	op.GeoM.Scale(1, scaleY)

	// put shading based on the distance
	colorScale := w.distanceShadeScale(distance)
	op.ColorScale.Scale(colorScale, colorScale, colorScale, 1)

	op.GeoM.Translate(float64(x), drawStart)
	screen.DrawImage(textureStrip, op)
}

// distanceShadeScale returns a grayscale multiplier for distance shading.
// farther walls get darker.
func (w *World) distanceShadeScale(distance float64) float32 {
	if distance <= 0 || math.IsInf(distance, 1) || math.IsNaN(distance) {
		return 1
	}

	return float32(1.0) / (1.0 + float32(distance)/20.0)
}

// drawSprites renders all world sprites.
// it uses the z-buffer to clip sprites behind walls and sorts sprites back-to-front.
func (w *World) drawSprites(screen *ebiten.Image, g *Game, p *Player) {
	if screen == nil || g == nil || p == nil || g.assets == nil {
		return
	}
	if w.zBuffer == nil || len(w.zBuffer) != WindowSizeX {
		return
	}

	// build camera plane used by ray direction: rayDir = dir + plane * cameraX
	plane := Vec2{
		X: -p.dir.Y * w.fovScale,
		Y: p.dir.X * w.fovScale,
	}

	// collect visible sprites with distance for sorting
	allSprites := make([]*Sprite, 0, len(g.sprites)+1)
	allSprites = append(allSprites, g.sprites...)

	// add other player as a sprite
	var other *Player
	if g.currentPlayer == g.playerX {
		other = g.playerO
	} else {
		other = g.playerX
	}

	if other != nil {
		allSprites = append(allSprites, &Sprite{
			Position:  other.pos,
			TextureID: other.characterTextureID,
			Scale:     1.0,
			Z:         0.0,
			Hidden:    false,
		})
	}

	sortedSprites := SortSpritesByDistance(allSprites, p.pos)

	for _, s := range sortedSprites {
		w.drawSingleSprite(screen, g, p, plane, s)
	}
}

// drawSingleSprite projects one sprite into the screen and draws it column by column.
func (w *World) drawSingleSprite(screen *ebiten.Image, g *Game, p *Player, plane Vec2, s *Sprite) {
	texture := g.assets.Textures[s.TextureID]
	if len(texture.Strips) == 0 {
		return
	}

	// sprite position relative to player
	spriteX := s.Position.X - p.pos.X
	spriteY := s.Position.Y - p.pos.Y

	// inverse determinant for camera transform
	det := plane.X*p.dir.Y - p.dir.X*plane.Y
	if det == 0 {
		return
	}
	invDet := 1.0 / det

	// transform sprite into camera space
	transformX := invDet * (p.dir.Y*spriteX - p.dir.X*spriteY)
	transformY := invDet * (-plane.Y*spriteX + plane.X*spriteY)

	// transformY is depth (in front of camera must be > 0)
	if transformY <= 0 || math.IsInf(transformY, 1) || math.IsNaN(transformY) {
		return
	}

	// sprite scale guard
	scale := s.Scale
	if scale <= 0 {
		scale = 1.0
	}

	// screen x of the sprite center
	spriteScreenX := int(float64(WindowSizeXDiv2) * (1.0 + transformX/transformY))

	// projected sprite size (classic)
	spriteHeight := int((float64(WindowSizeY) / transformY) * scale)
	if spriteHeight <= 0 {
		return
	}

	spriteWidth := spriteHeight

	// vertical placement
	// z shifts the sprite up in world units, scaled by depth into pixels
	zOffsetPx := int((s.Z / transformY) * float64(WindowSizeYDiv2))
	drawStartY := -spriteHeight/2 + WindowSizeYDiv2 - zOffsetPx

	// horizontal placement
	drawStartX := -spriteWidth/Two + spriteScreenX
	drawEndX := spriteWidth/Two + spriteScreenX

	// clip to screen bounds
	if drawStartX < 0 {
		drawStartX = 0
	}
	if drawEndX >= WindowSizeX {
		drawEndX = WindowSizeX - 1
	}

	// precompute shading from depth
	shade := w.distanceShadeScale(transformY)

	stripCount := len(texture.Strips)
	if stripCount <= 0 {
		return
	}

	// draw one screen column at a time, selecting the matching texture strip
	for x := drawStartX; x <= drawEndX; x++ {
		// z-buffer test: if wall is closer, skip this sprite column
		if transformY >= w.zBuffer[x] {
			continue
		}

		// map current screen x into [0..1] inside sprite
		u := (float64(x - (spriteScreenX - spriteWidth/2))) / float64(spriteWidth)
		if u < 0 || u > 1 {
			continue
		}

		stripIndex := int(u * float64(stripCount))
		stripIndex = clampInt(stripIndex, 0, stripCount-1)

		strip := texture.Strips[stripIndex]
		if strip == nil {
			continue
		}

		// scale strip vertically to match projected sprite height
		op := &ebiten.DrawImageOptions{}
		scaleY := float64(spriteHeight) / float64(TextureSize)
		op.GeoM.Scale(1, scaleY)

		// apply distance shading
		op.ColorScale.Scale(shade, shade, shade, 1)

		op.GeoM.Translate(float64(x), float64(drawStartY))
		screen.DrawImage(strip, op)
	}
}

// clampInt clamps an int value in the inclusive range [lo, hi].
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// drawCeiling draws the ceiling color rectangle.
func drawCeiling(screen *ebiten.Image) {
	if screen == nil {
		return
	}
	vector.FillRect(screen, float32(0), float32(0), float32(WindowSizeX), float32(WindowSizeYDiv2), ColorCeiling, false)
}

// drawFloor draws the floor color rectangle.
func drawFloor(screen *ebiten.Image) {
	if screen == nil {
		return
	}
	vector.FillRect(
		screen,
		float32(0),
		float32(WindowSizeYDiv2),
		float32(WindowSizeX),
		float32(WindowSizeYDiv2),
		ColorFloor,
		false,
	)
}
