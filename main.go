package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	sWidth     = 480
	sHeight    = 600
	sWidthMid  = sWidth / 2
	sHeightMid = sHeight / 2
	gridCells  = 3
	gridCellSz = sWidth / gridCells
	cellPad    = 7

	fontSize               = 15
	fontSizeLineSpacing    = 5
	bigFontSize            = 100
	bigFontSizeLineSpacing = 20

	pressTicksToReset = 60
	pressTicksToExit  = 60
	keyHoldMinFrames  = 2
)

const (
	cross  = "X"
	circle = "O"
	tie    = "tie"
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

//go:embed images/*
var imageFS embed.FS

var (
	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
	fontSource     *text.GoTextFaceSource
)

type GameState int

const (
	stateBoot GameState = iota
	statePlay
	stateDone
)

type Game struct {
	state   GameState
	board   [gridCells][gridCells]string
	round   int
	pointsO int
	pointsX int
	playing string
	win     string
	alter   int

	boardImg *ebiten.Image
	gameImg  *ebiten.Image
	imgCache map[string]*ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.KeyPressDuration(ebiten.KeyR) == pressTicksToReset {
		g.resetRound()
		g.resetPoints()
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == pressTicksToExit {
		os.Exit(0)
	}

	switch g.state {

	case stateBoot:
		g.initAssets()
		g.seedFirstPlayer()
		g.resetRound()
		g.resetPoints()
		g.state = statePlay
		return nil

	case statePlay:
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return nil
		}
		cx, cy, ok := getCursorIdxFromCell()
		if !ok {
			return nil
		}
		if g.board[cx][cy] != "" {
			return nil
		}

		cur := g.currentSymbol()
		nxt := otherSymbol(cur)

		g.drawSymbol(cx, cy, cur)
		g.board[cx][cy] = cur
		g.playing = nxt

		w := g.checkWin()
		if w != "" {
			g.applyWinner(w)
			return nil
		}

		g.round++
		return nil

	case stateDone:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.resetRound()
			g.state = statePlay
		}
		return nil
	}

	return nil
}

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

func (g *Game) initAssets() {
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	normalTextFace = &text.GoTextFace{Source: fontSource, Size: fontSize}
	bigTextFace = &text.GoTextFace{Source: fontSource, Size: bigFontSize}

	g.boardImg = loadImage("images/board.png")
	g.gameImg = ebiten.NewImage(sWidth, sWidth)
	g.imgCache = map[string]*ebiten.Image{
		circle: loadImage("images/O.png"),
		cross:  loadImage("images/X.png"),
	}
}

func loadImage(path string) *ebiten.Image {
	data, err := imageFS.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func (g *Game) seedFirstPlayer() {
	if newRandom().Intn(2) == 0 {
		g.playing = circle
		g.alter = 0
		return
	}
	g.playing = cross
	g.alter = 1
}

func (g *Game) resetRound() {
	g.gameImg.Clear()
	g.board = [gridCells][gridCells]string{}
	g.round = 0

	// alternate the round starter
	if g.alter == 0 {
		g.playing = cross
		g.alter = 1
	} else {
		g.playing = circle
		g.alter = 0
	}

	g.win = ""
	g.state = statePlay
}

func (g *Game) resetPoints() {
	g.pointsO = 0
	g.pointsX = 0
}

func (g *Game) currentSymbol() string {
	if g.round%2 == g.alter {
		return circle
	}
	return cross
}

func otherSymbol(s string) string {
	if s == circle {
		return cross
	}
	return circle
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

func (g *Game) applyWinner(w string) {
	switch w {
	case circle:
		g.win = circle
		g.pointsO++
		g.state = stateDone
	case cross:
		g.win = cross
		g.pointsX++
		g.state = stateDone
	case tie:
		g.win = "No one"
		g.state = stateDone
	}
}

func (g *Game) checkWin() string {
	for i := 0; i < gridCells; i++ {
		a, b, c := g.board[i][0], g.board[i][1], g.board[i][2]
		if a != "" && a == b && b == c {
			return a
		}
	}
	for i := 0; i < gridCells; i++ {
		a, b, c := g.board[0][i], g.board[1][i], g.board[2][i]
		if a != "" && a == b && b == c {
			return a
		}
	}
	m := g.board[1][1]
	if m != "" {
		if g.board[0][0] == m && g.board[2][2] == m {
			return m
		}
		if g.board[0][2] == m && g.board[2][0] == m {
			return m
		}
	}

	if g.round == gridCells*gridCells-1 {
		return tie
	}
	return ""
}

func getCursorIdxFromCell() (int, int, bool) {
	mx, my := ebiten.CursorPosition()
	if mx < 0 || my < 0 {
		return 0, 0, false
	}
	cx := mx / gridCellSz
	cy := my / gridCellSz
	if cx >= gridCells || cy >= gridCells {
		return 0, 0, false
	}
	return cx, cy, true
}

func newRandom() *rand.Rand {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1)
}

func (g *Game) Layout(int, int) (int, int) { return sWidth, sHeight }

func main() {
	game := &Game{state: stateBoot}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
