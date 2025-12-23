package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

/* ---------- helpers ---------- */

func fillRect(dst *ebiten.Image, x, y, w, h float32, col color.Color) {
	r, g, b, a := col.RGBA()
	vector.FillRect(
		dst,
		x, y, w, h,
		color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		},
		false,
	)
}

/* ---------- draw loop ---------- */

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

	g.drawMinimap(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenSize, ScreenSize
}

/* ---------- UI ---------- */

func (g *Game) drawNameInput(screen *ebiten.Image) {
	g.drawText(screen, "Enter player names", NameInputX, NameInputY, TopLeft, color.White)

	label := "Player X: "
	if !g.editingPlayerX {
		label = "Player O: "
	}
	g.drawText(screen, label+g.inputBuffer, NameInputX, NameInputY+NameInputLineHeight, TopLeft, color.White)

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
	score := "Score " + g.playerXName + ": " + strconv.Itoa(g.scoreX) +
		"  " + g.playerOName + ": " + strconv.Itoa(g.scoreO)

	g.drawText(screen, score, Margin, HeaderY, TopLeft, color.White)
	g.drawText(screen, "ESC = quit", ScreenSize-Margin, HeaderY, TopRight, color.White)
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
	msg := "Draw!"
	if g.winner == WinnerX {
		msg = g.playerXName + " wins!"
	} else if g.winner == WinnerO {
		msg = g.playerOName + " wins!"
	}
	msg += " Click to restart"
	g.drawText(screen, msg, Margin, BottomY, BottomLeft, color.White)
}

/* ---------- grid & pieces ---------- */

func (g *Game) drawGrid(screen *ebiten.Image) {
	col := color.RGBA{200, 200, 200, 255}

	for i := 1; i < GridSize; i++ {
		fillRect(
			screen,
			0, float32(i*CellSize),
			float32(ScreenSize), float32(LineWidth),
			col,
		)
		fillRect(
			screen,
			float32(i*CellSize), 0,
			float32(LineWidth), float32(ScreenSize),
			col,
		)
	}
}

func (g *Game) drawPieces(screen *ebiten.Image) {
	for y := range GridSize {
		for x := range GridSize {
			cell := g.board[y][x]
			if cell == PlayerNone {
				continue
			}
			img := g.assets.XImage
			if cell == PlayerO {
				img = g.assets.OImage
			}
			g.drawPieceAt(screen, img, x, y)
		}
	}
}

func (g *Game) drawPieceAt(screen *ebiten.Image, img *ebiten.Image, x, y int) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	scale := float64(CellSize-2*Margin) / float64(w)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	px := float64(x*CellSize) + (float64(CellSize)-float64(w)*scale)/2
	py := float64(y*CellSize) + (float64(CellSize)-float64(h)*scale)/2

	op.GeoM.Translate(px, py)
	screen.DrawImage(img, op)
}

/* ---------- text ---------- */

func (g *Game) drawText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color) {
	op := &text.DrawOptions{}
	op.LineSpacing = g.assets.NormalTextFace.Size + TextLineSpacing
	op.ColorScale.ScaleWithColor(col)

	w, h := text.Measure(msg, g.assets.NormalTextFace, op.LineSpacing)

	if align == TopCenter || align == Center || align == BottomCenter {
		x -= w / 2
	} else if align == TopRight || align == CenterRight || align == BottomRight {
		x -= w
	}

	if align == CenterLeft || align == Center || align == CenterRight {
		y -= h / 2
	} else if align == BottomLeft || align == BottomCenter || align == BottomRight {
		y -= h
	}

	op.GeoM.Translate(x, y)
	text.Draw(screen, msg, g.assets.NormalTextFace, op)
}

/* ---------- minimap ---------- */

func (g *Game) drawMinimap(screen *ebiten.Image) {
	const tilePx = 8.0
	const pad = 10.0
	const border = 2.0
	const r = 2.0

	mapW := float64(g.worldMap.Width()) * tilePx
	mapH := float64(g.worldMap.Height()) * tilePx

	originX := float64(ScreenSize) - pad - mapW
	originY := pad

	fillRect(
		screen,
		float32(originX-border), float32(originY-border),
		float32(mapW+border*2), float32(mapH+border*2),
		color.RGBA{0, 0, 0, 170},
	)

	for y := 0; y < g.worldMap.Height(); y++ {
		for x := 0; x < g.worldMap.Width(); x++ {
			if g.worldMap.Tiles[y][x] == 1 {
				fillRect(
					screen,
					float32(originX+float64(x)*tilePx),
					float32(originY+float64(y)*tilePx),
					float32(tilePx), float32(tilePx),
					color.RGBA{200, 200, 200, 230},
				)
			}
		}
	}

	px := originX + g.playerPos.X*tilePx
	py := originY + g.playerPos.Y*tilePx
	fillRect(
		screen,
		float32(px-r), float32(py-r),
		float32(r*2), float32(r*2),
		color.RGBA{255, 80, 80, 255},
	)
}
