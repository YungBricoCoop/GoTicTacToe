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
	sWidth      = 480
	sHeight     = 600
	fontSize    = 15
	bigFontSize = 100
	dpi         = 72
)

//go:embed images/*
var imageFS embed.FS

var (
	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
	boardImage     *ebiten.Image
	symbolImage    *ebiten.Image
	textImage      = ebiten.NewImage(sWidth, sWidth)
	gameImage      = ebiten.NewImage(sWidth, sWidth)
	fontSource     *text.GoTextFaceSource
)

type Game struct {
	playing   string
	state     int
	gameBoard [3][3]string
	round     int
	pointsO   int
	pointsX   int
	win       string
	alter     int
}

func (g *Game) Update() error {
	switch g.state {
	case 0:
		g.Init()
		break

	case 1:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if mx/160 < 3 && mx >= 0 && my/160 < 3 && my >= 0 && g.gameBoard[mx/160][my/160] == "" {
				if g.round%2 == 0+g.alter {
					g.DrawSymbol(mx/160, my/160, "O")
					g.gameBoard[mx/160][my/160] = "O"
					g.playing = "X"
				} else {
					g.DrawSymbol(mx/160, my/160, "X")
					g.gameBoard[mx/160][my/160] = "X"
					g.playing = "O"
				}
				g.wins(g.CheckWin())
				g.round++
			}
		}
		break
	case 2:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Load()
		}
		break
	}
	if inpututil.KeyPressDuration(ebiten.KeyR) == 60 {
		g.Load()
		g.ResetPoints()
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == 60 {
		os.Exit(0)
	}
	return nil
}

func keyChangeColor(key ebiten.Key, screen *ebiten.Image) {
	if inpututil.KeyPressDuration(key) > 1 {
		var msgText string
		var colorText color.RGBA
		colorChange := 255 - (255 / 60 * uint8(inpututil.KeyPressDuration(key)))
		if key == ebiten.KeyEscape {
			msgText = fmt.Sprintf("CLOSING...")
			colorText = color.RGBA{R: 255, G: colorChange, B: colorChange, A: 255}
		} else if key == ebiten.KeyR {
			msgText = fmt.Sprintf("RESETING...")
			colorText = color.RGBA{R: colorChange, G: 255, B: 255, A: 255}
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(sWidth/2, sHeight-30)
		op.ColorScale.ScaleWithColor(colorText)
		text.Draw(screen, msgText, normalTextFace, op)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(boardImage, nil)
	screen.DrawImage(gameImage, nil)
	mx, my := ebiten.CursorPosition()

	msgFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	opFPS := &text.DrawOptions{}
	opFPS.GeoM.Translate(0, sHeight-30)
	opFPS.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msgFPS, normalTextFace, opFPS)

	keyChangeColor(ebiten.KeyEscape, screen)
	keyChangeColor(ebiten.KeyR, screen)
	msgOX := fmt.Sprintf("O: %v | X: %v", g.pointsO, g.pointsX)
	opOX := &text.DrawOptions{}
	opOX.GeoM.Translate(sWidth/2, sHeight-5)
	opOX.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msgOX, normalTextFace, opOX)
	if g.win != "" {
		msgWin := fmt.Sprintf("%v wins!", g.win)
		opWin := &text.DrawOptions{}
		opWin.GeoM.Translate(70, 200)
		opWin.ColorScale.ScaleWithColor(color.RGBA{G: 50, B: 200, A: 255})
		text.Draw(screen, msgWin, bigTextFace, opWin)
	}
	msg := fmt.Sprintf("%v", g.playing)
	opPlaying := &text.DrawOptions{}
	opPlaying.GeoM.Translate(float64(mx), float64(my))
	opPlaying.ColorScale.ScaleWithColor(color.RGBA{G: 255, A: 255})
	text.Draw(screen, msg, normalTextFace, opPlaying)
}

func (g *Game) DrawSymbol(x, y int, sym string) {
	imageBytes, err := imageFS.ReadFile(fmt.Sprintf("images/%v.png", sym))
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	symbolImage = ebiten.NewImageFromImage(decoded)
	opSymbol := &ebiten.DrawImageOptions{}
	opSymbol.GeoM.Translate(float64((160*(x+1)-160)+7), float64((160*(y+1)-160)+7))

	gameImage.DrawImage(symbolImage, opSymbol)
}

func (g *Game) Init() {
	imageBytes, err := imageFS.ReadFile("images/board.png")
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	boardImage = ebiten.NewImageFromImage(decoded)
	re := newRandom().Intn(2)
	if re == 0 {
		g.playing = "O"
		g.alter = 0
	} else {
		g.playing = "X"
		g.alter = 1
	}
	g.Load()
	g.ResetPoints()
}

func (g *Game) Load() {
	gameImage.Clear()
	g.gameBoard = [3][3]string{{"", "", ""}, {"", "", ""}, {"", "", ""}}
	g.round = 0
	if g.alter == 0 {
		g.playing = "X"
		g.alter = 1
	} else if g.alter == 1 {
		g.playing = "O"
		g.alter = 0
	}
	g.win = ""
	g.state = 1
}

func (g *Game) wins(winner string) {
	if winner == "O" {
		g.win = "O"
		g.pointsO++
		g.state = 2
	} else if winner == "X" {
		g.win = "X"
		g.pointsX++
		g.state = 2
	} else if winner == "tie" {
		g.win = "No one\n"
		g.state = 2
	}
}

func (g *Game) CheckWin() string {
	for i, _ := range g.gameBoard {
		if g.gameBoard[i][0] == g.gameBoard[i][1] && g.gameBoard[i][1] == g.gameBoard[i][2] {
			return g.gameBoard[i][0]
		}
	}
	for i, _ := range g.gameBoard {
		if g.gameBoard[0][i] == g.gameBoard[1][i] && g.gameBoard[1][i] == g.gameBoard[2][i] {
			return g.gameBoard[0][i]
		}
	}
	if (g.gameBoard[0][0] == g.gameBoard[1][1] && g.gameBoard[1][1] == g.gameBoard[2][2]) || (g.gameBoard[0][2] == g.gameBoard[1][1] && g.gameBoard[1][1] == g.gameBoard[2][0]) {
		return g.gameBoard[1][1]
	}
	if g.round == 8 {
		return "tie"
	}
	return ""
}

func (g *Game) ResetPoints() {
	g.pointsO = 0
	g.pointsX = 0
}

func init() {
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	normalTextFace = &text.GoTextFace{
		Source: fontSource,
		Size:   fontSize,
	}
	bigTextFace = &text.GoTextFace{
		Source: fontSource,
		Size:   bigFontSize,
	}
}

func newRandom() *rand.Rand {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1)
}

func (g *Game) Layout(int, int) (int, int) {
	return sWidth, sHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
