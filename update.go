package main

import (
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/v2"
)

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
