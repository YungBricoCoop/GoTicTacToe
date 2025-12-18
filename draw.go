package main

import (
	"embed"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

//go:embed images/*.png
var imageFS embed.FS

type TextAlign int

const (
	TopLeft TextAlign = iota
	TopCenter
	TopRight
	CenterLeft
	Center
	CenterRight
	BottomLeft
	BottomCenter
	BottomRight
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
		ox = -w / 2 //nolint:mnd // centering
	case TopRight, CenterRight, BottomRight:
		ox = -w
	}

	var oy float64
	switch align {
	case TopLeft, TopCenter, TopRight:
		oy = 0
	case CenterLeft, Center, CenterRight:
		oy = -h / 2 //nolint:mnd // centering
	case BottomLeft, BottomCenter, BottomRight:
		oy = -h
	}

	op.GeoM.Translate(x+ox, y+oy)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(screen, msg, face, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	face := basicfont.Face7x13

	// === NAME INPUT SCREEN ===
	if g.entering {
		title := "Enter player names"
		drawText(screen, title, 10, 40, TopLeft, color.White, g.normalTextFace, 0)

		label := "Player X: "
		if !g.editingX {
			label = "Player O: "
		}
		drawText(screen, label+g.tempInput, 10, 80, TopLeft, color.White, g.normalTextFace, 0)

		info := "Type name, Enter = OK, Backspace = delete, R = reset"
		// text.Draw(screen, info, face, 10, 120, color.White)
		drawText(screen, info, 10, 120, TopLeft, color.White, g.normalTextFace, 0)
		return
	}

	score := "Score " + g.playerXName + ": " + strconv.Itoa(
		g.pointsX,
	) + "  " + g.playerOName + ": " + strconv.Itoa(
		g.pointsO,
	)
	// text.Draw(screen, score, face, 10, 20, color.White)
	drawText(screen, score, 10, 20, TopLeft, color.White, g.normalTextFace, 0)

	// text.Draw(screen, "ESC = quit", face, ScreenSize-110, 20, color.White)
	drawText(screen, "ESC = quit", ScreenSize, 20, TopRight, color.White, g.normalTextFace, 0)

	// draw the grid
	lineColor := color.RGBA{200, 200, 200, 255}

	for i := 1; i < GridSize; i++ {
		// horizontales
		h := ebiten.NewImage(ScreenSize, 2)
		h.Fill(lineColor)
		opH := &ebiten.DrawImageOptions{}
		opH.GeoM.Translate(0, float64(i*CellSize))
		screen.DrawImage(h, opH)

		// verticales
		v := ebiten.NewImage(2, ScreenSize)
		v.Fill(lineColor)
		opV := &ebiten.DrawImageOptions{}
		opV.GeoM.Translate(float64(i*CellSize), 0)
		screen.DrawImage(v, opV)
	}

	// drawing X/O
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.board[y][x] == 0 {
				continue
			}

			var img *ebiten.Image
			if g.board[y][x] == 1 {
				img = xImage
			} else {
				img = oImage
			}

			w := img.Bounds().Dx()
			h := img.Bounds().Dy()
			_ = w
			_ = h

			targetSize := CellSize - 2*Margin

			scaleW := float64(targetSize) / float64(w)
			scaleH := float64(targetSize) / float64(h)
			scale := scaleW
			if scaleH < scaleW {
				scale = scaleH
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)

			cellX := float64(x * CellSize)
			CellY := float64(y * CellSize)

			imgW := float64(w) * scale
			imgH := float64(h) * scale

			px := cellX + (float64(CellSize)-imgW)/2
			py := CellY + (float64(CellSize)-imgH)/2

			op.GeoM.Translate(px, py)
			screen.DrawImage(img, op)
		}
	}

	// message at the bottom
	if !g.over {
		msg := "Turn: "
		if g.player == 1 {
			msg += g.playerXName
		} else {
			msg += g.playerOName
		}
		msg += "   (R = full reset)"
		text.Draw(screen, msg, face, 10, ScreenSize-10, color.White)
	} else {
		msg := ""
		switch g.winner {
		case 1:
			msg = g.playerXName + " wins! Click to restart (keep score) or R"
		case 2:
			msg = g.playerOName + " wins! Click to restart (keep score) or R"
		case 3:
			msg = "Draw! Click to restart (keep score) or R"
		}
		text.Draw(screen, msg, face, 10, ScreenSize-10, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenSize, ScreenSize
}
