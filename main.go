package main

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

//go:embed images/*.png
var imageFS embed.FS

const (
	ScreenSize = 480
	WindowSize = 480

	GridSize = 3
	CellSize = ScreenSize / GridSize

	Margin    = 10
	LineWidth = 2
	HeaderY   = 20
	BottomY   = ScreenSize - 10
)

var (
	xImage *ebiten.Image
	oImage *ebiten.Image
)

type Game struct {
	board  [GridSize][GridSize]int // 0 = empty, 1 = X, 2 = O
	player int                     // current player: 1 (X) or 2 (O)
	winner int                     // 0 = nobody, 1 = X, 2 = O, 3 = draw
	over   bool                    // game over?

	// scores
	pointsX int
	pointsO int

	// Player names
	playerXName string
	playerOName string

	// Name input
	entering  bool   // true while entering names
	editingX  bool   // true = entering X's name, false = O's name
	tempInput string // buffer for current text

	startPlayer int // 1 = x, 2 = o
}

func NewGame() *Game {
	return &Game{
		player:      1, // X starts
		startPlayer: 1,
		playerXName: "X",  // default values just in case
		playerOName: "O",  // default values
		entering:    true, // start by entering names
		editingX:    true, // start with player X
		pointsX:     0,
		pointsO:     0,
	}
}

func (g *Game) resetBoard() {
	g.board = [GridSize][GridSize]int{}
	g.winner = 0
	g.over = false

	if g.startPlayer == 1 {
		g.startPlayer = 2
	} else {
		g.startPlayer = 1
	}
	g.player = g.startPlayer
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenSize, ScreenSize
}

func (g *Game) Update() error {
	// ESC = quit the game
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// R = full reset
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *NewGame()
		return nil
	}

	// name input phase I just added
	if g.entering {
		// Get typed characters
		chars := ebiten.AppendInputChars(nil)
		for _, c := range chars {
			if c == '\n' || c == '\r' {
				continue
			}
			g.tempInput += string(c)
		}

		// Backspace
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.tempInput) > 0 {
			g.tempInput = g.tempInput[:len(g.tempInput)-1]
		}

		// Enter -> validate the name
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if g.editingX {
				if g.tempInput != "" {
					g.playerXName = g.tempInput
				}
				g.tempInput = ""
				g.editingX = false // switch to O
			} else {
				if g.tempInput != "" {
					g.playerOName = g.tempInput
				}
				g.tempInput = ""
				g.entering = false // done, we can play
			}
		}

		// While entering names, we don't play
		return nil
	}

	// === GAME OVER ===
	if g.over {
		// Left click = start a new game
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.resetBoard()
		}
		return nil
	}

	// === NORMAL GAME: click on a cell ===
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		if mx >= 0 && mx < ScreenSize && my >= 0 && my < ScreenSize {
			cx := mx / CellSize
			cy := my / CellSize

			if g.board[cy][cx] == 0 { // empty cell
				g.board[cy][cx] = g.player

				// check if someone wins
				if w := g.checkWinner(); w != 0 {
					g.winner = w
					g.over = true

					if w == 1 {
						g.pointsX++
					} else if w == 2 {
						g.pointsO++
					}

				} else if g.isBoardFull() {
					g.winner = 3 // draw
					g.over = true
				} else {
					// switch player
					if g.player == 1 {
						g.player = 2
					} else {
						g.player = 1
					}
				}
			}
		}
	}

	return nil
}

// check if X or O has won
func (g *Game) checkWinner() int {
	b := g.board

	// rows
	for y := 0; y < GridSize; y++ {
		if b[y][0] != 0 && b[y][0] == b[y][1] && b[y][1] == b[y][2] {
			return b[y][0]
		}
	}

	// columns
	for x := 0; x < GridSize; x++ {
		if b[0][x] != 0 && b[0][x] == b[1][x] && b[1][x] == b[2][x] {
			return b[0][x]
		}
	}

	// main diagonal
	if b[0][0] != 0 && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}

	// other diagonal
	if b[0][2] != 0 && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}

	return 0
}

func (g *Game) isBoardFull() bool {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.board[y][x] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	face := basicfont.Face7x13

	// === NAME INPUT SCREEN ===
	if g.entering {
		title := "Enter player names"
		text.Draw(screen, title, face, 10, 40, color.White)

		label := "Player X: "
		if !g.editingX {
			label = "Player O: "
		}
		text.Draw(screen, label+g.tempInput, face, 10, 80, color.RGBA{200, 200, 0, 255})

		info := "Type name, Enter = OK, Backspace = delete, R = reset"
		text.Draw(screen, info, face, 10, 120, color.White)
		return
	}

	score := "Score " + g.playerXName + ": " + strconv.Itoa(
		g.pointsX,
	) + "  " + g.playerOName + ": " + strconv.Itoa(
		g.pointsO,
	)
	text.Draw(screen, score, face, 10, 20, color.White)

	text.Draw(screen, "ESC = quit", face, ScreenSize-110, 20, color.White)

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

func mustLoadImage(path string) *ebiten.Image {
	log.Printf("Loading images from embedded filesystem:")
	files, err := imageFS.ReadDir("images")
	if err != nil {
		log.Fatalf("cannot read images directory: %v", err)
	}
	for _, f := range files {
		log.Printf("Loading image: %s", f.Name())
	}

	data, err := imageFS.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read %s: %v", path, err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("cannot decode %s: %v", path, err)
	}
	return ebiten.NewImageFromImage(img)
}

func loadImages() {
	xImage = mustLoadImage("images/x.png")
	oImage = mustLoadImage("images/o.png")
}

func main() {
	loadImages()

	game := NewGame()
	ebiten.SetWindowSize(WindowSize, WindowSize)
	ebiten.SetWindowTitle("Tic-Tac-Toe (images + names + scores)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
