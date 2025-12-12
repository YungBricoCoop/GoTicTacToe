package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawText(
	screen *ebiten.Image,
	msg string,
	x, y float64,
	align TextAlign,
	col color.Color,
	face *text.GoTextFace,
	lineSpacingPaddingPx float64,
) {
	op := &text.DrawOptions{}
	op.LineSpacing = face.Size + lineSpacingPaddingPx

	w, h := text.Measure(msg, face, op.LineSpacing)

	var ox float64
	switch align {
	case TopLeft, CenterLeft, BottomLeft:
		ox = 0
	case TopCenter, Center, BottomCenter:
		ox = -w / 2
	case TopRight, CenterRight, BottomRight:
		ox = -w
	}

	var oy float64
	switch align {
	case TopLeft, TopCenter, TopRight:
		oy = 0
	case CenterLeft, Center, CenterRight:
		oy = -h / 2
	case BottomLeft, BottomCenter, BottomRight:
		oy = -h
	}

	op.GeoM.Translate(x+ox, y+oy)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(screen, msg, face, op)
}

func drawKeyHoldFeedback(screen *ebiten.Image, key ebiten.Key, label string, base color.RGBA) {
	frames := inpututil.KeyPressDuration(key)
	if frames < keyHoldMinFrames {
		return
	}
	if frames > pressTicksToReset {
		frames = pressTicksToReset
	}
	fade := uint8(255 - (255*frames)/pressTicksToReset)

	c := base
	if key == ebiten.KeyEscape {
		c.G = fade
		c.B = fade
	}
	if key == ebiten.KeyR {
		c.R = fade
	}

	drawText(
		screen,
		label,
		float64(sWidthMid),
		float64(sHeight-30),
		BottomCenter,
		c,
		normalTextFace,
		fontSizeLineSpacing,
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the cached board (background + grid + symbols)
	screen.DrawImage(g.boardImage, nil)

	mx, my := ebiten.CursorPosition()

	// metrics
	msgFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	drawText(screen, msgFPS, 0, float64(sHeight), BottomLeft, colText, normalTextFace, fontSizeLineSpacing)

	// pressed key info
	drawKeyHoldFeedback(screen, ebiten.KeyEscape, "CLOSING...", colReset)
	drawKeyHoldFeedback(screen, ebiten.KeyR, "RESETING...", colReset)

	// score
	msgOX := fmt.Sprintf("O: %d | X: %d", g.pointsO, g.pointsX)
	drawText(
		screen,
		msgOX,
		float64(sWidthMid),
		float64(sHeight),
		BottomCenter,
		colText,
		normalTextFace,
		fontSizeLineSpacing,
	)

	// winner banner
	if g.win != "" {
		msgWin := fmt.Sprintf("%s wins!", g.win)
		drawText(
			screen,
			msgWin,
			float64(sWidthMid),
			float64(sHeightMid),
			Center,
			colWin,
			bigTextFace,
			bigFontSizeLineSpacing,
		)
	}

	// cursor for the current player
	drawText(
		screen,
		g.playing,
		float64(mx),
		float64(my),
		TopLeft,
		colText,
		normalTextFace,
		fontSizeLineSpacing,
	)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for i := 1; i < gridCells; i++ {
		// vertical lines
		x := float32(i * gridCellSz)
		vector.StrokeLine(screen, x, 0, x, float32(sWidth), lineWidth, colGrid, true)

		// horizontal lines
		y := float32(i * gridCellSz)
		vector.StrokeLine(screen, 0, y, float32(sWidth), y, lineWidth, colGrid, true)
	}
}

func (g *Game) drawSymbol(screen *ebiten.Image, cx, cy int, symbol string) {
	size := float32(gridCellSz/2 - cellPad)
	x := float32(cx*gridCellSz + gridCellSz/2)
	y := float32(cy*gridCellSz + gridCellSz/2)

	if symbol == cross {
		g.drawX(screen, x, y, size)
	} else if symbol == circle {
		g.drawO(screen, x, y, size)
	}
}

func (g *Game) drawX(screen *ebiten.Image, cx, cy, size float32) {
	vector.StrokeLine(screen, cx-size, cy-size, cx+size, cy+size, lineWidth, colX, true)
	vector.StrokeLine(screen, cx+size, cy-size, cx-size, cy+size, lineWidth, colX, true)
}

func (g *Game) drawO(screen *ebiten.Image, cx, cy, size float32) {
	vector.StrokeCircle(screen, cx, cy, size, lineWidth, colO, true)
}
