package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func drawText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color, face *text.GoTextFace, lineSpacingPaddingPx float64) {
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

	drawText(screen, label, float64(sWidthMid), float64(sHeight-30), BottomCenter, c, normalTextFace, fontSizeLineSpacing)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.boardImg, nil)
	screen.DrawImage(g.gameImg, nil)

	mx, my := ebiten.CursorPosition()

	// metrics
	msgFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	drawText(screen, msgFPS, 0, float64(sHeight), BottomLeft, color.White, normalTextFace, fontSizeLineSpacing)

	// pressed key info
	drawKeyHoldFeedback(screen, ebiten.KeyEscape, "CLOSING...", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	drawKeyHoldFeedback(screen, ebiten.KeyR, "RESETING...", color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// score
	msgOX := fmt.Sprintf("O: %d | X: %d", g.pointsO, g.pointsX)
	drawText(screen, msgOX, float64(sWidthMid), float64(sHeight), BottomCenter, color.White, normalTextFace, fontSizeLineSpacing)

	// winner banner
	if g.win != "" {
		msgWin := fmt.Sprintf("%s wins!", g.win)
		drawText(screen, msgWin, float64(sWidthMid), float64(sHeightMid), Center, color.RGBA{G: 50, B: 200, A: 255}, bigTextFace, bigFontSizeLineSpacing)
	}

	// cursor for the current player
	drawText(screen, g.playing, float64(mx), float64(my), TopLeft, color.RGBA{G: 255, A: 255}, normalTextFace, fontSizeLineSpacing)
}

func (g *Game) drawSymbol(cx, cy int, sym string) {
	img := g.imgCache[sym]
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cx*gridCellSz+cellPad), float64(cy*gridCellSz+cellPad))
	g.gameImg.DrawImage(img, op)
}
