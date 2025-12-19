package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	switch g.state {
	case StateNameInput:
		g.drawNameInput(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	case StateGameOver:
		g.drawPlaying(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenSize, ScreenSize
}

func (g *Game) drawNameInput(screen *ebiten.Image) {
	// title
	g.drawText(screen, "Enter player names", NameInputX, NameInputY, TopLeft, color.White)

	// current input label
	label := "Player X: "
	if !g.editingPlayerX {
		label = "Player O: "
	}
	g.drawText(screen, label+g.inputBuffer, NameInputX, NameInputY+NameInputLineHeight, TopLeft, color.White)

	// instructions
	info := "Type name, Enter = OK, Backspace = delete"
	g.drawText(screen, info, NameInputX, NameInputY+NameInputLineHeight*2, TopLeft, color.White)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	g.drawGrid(screen)
	g.drawPieces(screen)
	g.drawScoreAndShortcuts(screen)

	if g.state == StatePlaying {
		g.drawTurnInfo(screen)
	}
}

func (g *Game) drawScoreAndShortcuts(screen *ebiten.Image) {
	scoreMsg := "Score " + g.playerXName + ": " + strconv.Itoa(g.scoreX) +
		"  " + g.playerOName + ": " + strconv.Itoa(g.scoreO)
	g.drawText(screen, scoreMsg, Margin, HeaderY, TopLeft, color.White)

	g.drawText(screen, "ESC = quit", ScreenSize-Margin, HeaderY, TopRight, color.White)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	lineColor := color.RGBA{200, 200, 200, 255}

	for i := 1; i < GridSize; i++ {
		// horizontal lines
		h := ebiten.NewImage(ScreenSize, LineWidth)
		h.Fill(lineColor)
		opH := &ebiten.DrawImageOptions{}
		opH.GeoM.Translate(0, float64(i*CellSize))
		screen.DrawImage(h, opH)

		// vertical lines
		v := ebiten.NewImage(LineWidth, ScreenSize)
		v.Fill(lineColor)
		opV := &ebiten.DrawImageOptions{}
		opV.GeoM.Translate(float64(i*CellSize), 0)
		screen.DrawImage(v, opV)
	}
}

func (g *Game) drawPieces(screen *ebiten.Image) {
	for y := range GridSize {
		for x := range GridSize {
			cell := g.board[y][x]
			if cell == PlayerNone {
				continue
			}

			var img *ebiten.Image
			if cell == PlayerX {
				img = g.assets.XImage
			} else {
				img = g.assets.OImage
			}

			g.drawPieceAt(screen, img, x, y)
		}
	}
}

func (g *Game) drawPieceAt(screen *ebiten.Image, img *ebiten.Image, x, y int) {
	const (
		doubleMargin = 2 * Margin
		half         = 2
	)
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	targetSize := CellSize - doubleMargin
	scale := float64(targetSize) / float64(w)
	if h > w {
		scale = float64(targetSize) / float64(h)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	// center the image in the cell
	imgW := float64(w) * scale
	imgH := float64(h) * scale

	cellX := float64(x * CellSize)
	cellY := float64(y * CellSize)

	px := cellX + (float64(CellSize)-imgW)/half
	py := cellY + (float64(CellSize)-imgH)/half

	op.GeoM.Translate(px, py)
	screen.DrawImage(img, op)
}

func (g *Game) drawTurnInfo(screen *ebiten.Image) {
	msg := "Turn: "
	if g.currentPlayer == PlayerX {
		msg += g.playerXName
	} else {
		msg += g.playerOName
	}
	msg += "   (R = full reset)"
	g.drawText(screen, msg, Margin, BottomY, BottomLeft, color.White)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	var msg string
	switch g.winner {
	case WinnerX:
		msg = g.playerXName + " wins!"
	case WinnerO:
		msg = g.playerOName + " wins!"
	case WinnerDraw:
		msg = "Draw!"
	case WinnerNone:
		// should not happen in game over state
	}

	msg += " Click to restart"
	g.drawText(screen, msg, Margin, BottomY, BottomLeft, color.White)
}

func (g *Game) drawText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color) {
	const half = 2
	op := &text.DrawOptions{}
	op.LineSpacing = g.assets.NormalTextFace.Size + TextLineSpacing
	op.ColorScale.ScaleWithColor(col)

	w, h := text.Measure(msg, g.assets.NormalTextFace, op.LineSpacing)

	var ox float64
	switch align {
	case TopLeft, CenterLeft, BottomLeft:
		ox = 0
	case TopCenter, Center, BottomCenter:
		ox = -w / half
	case TopRight, CenterRight, BottomRight:
		ox = -w
	}

	var oy float64
	switch align {
	case TopLeft, TopCenter, TopRight:
		oy = 0
	case CenterLeft, Center, CenterRight:
		oy = -h / half
	case BottomLeft, BottomCenter, BottomRight:
		oy = -h
	}

	op.GeoM.Translate(x+ox, y+oy)
	text.Draw(screen, msg, g.assets.NormalTextFace, op)
}
