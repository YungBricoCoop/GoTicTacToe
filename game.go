// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type GameState int

const (
	StateNameInput GameState = iota
	StatePlaying
	StateGameOver
)

type Game struct {
	// state
	state  GameState
	board  Board
	winner *Player

	// timer to handle game over transition
	stateTimer float64

	// players
	playerX       *Player
	playerO       *Player
	currentPlayer *Player

	// sprites
	sprites []*Sprite

	// input state
	inputBuffer    string
	editingPlayerX bool

	// visuals
	assets *Assets

	worldMap Map

	// loop lists
	updatables []Updatable
	drawables  []Drawable
}

// NewGame creates a new Game instance with initialized assets and players.
func NewGame() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}

	spawnX := defaultPlayerXSpawn()
	spawnO := defaultPlayerOSpawn()

	pX := NewPlayer(spawnX.X, spawnX.Y, PlayerSymbolX, "X")
	pO := NewPlayer(spawnO.X, spawnO.Y, PlayerSymbolO, "O")

	minimap := &Minimap{}
	hud := &Hud{}
	world := &World{fovScale: GetK(PlayerFOV)}

	sprites := createSprites()

	g := &Game{
		state:          StateNameInput,
		winner:         nil,
		assets:         assets,
		worldMap:       NewMap(),
		playerX:        pX,
		playerO:        pO,
		currentPlayer:  pX,
		editingPlayerX: true,
		inputBuffer:    "",
		updatables:     nil,
		drawables:      nil,
		sprites:        sprites,
	}

	g.updatables = append(g.updatables,
		pX, pO,
	)

	g.drawables = append(g.drawables,
		world,
		minimap,
		hud,
	)

	return g, nil
}

func (g *Game) Update() error {
	for _, obj := range g.updatables {
		obj.Update(g)
	}

	// shortcuts
	// Escape: exit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// R: reset game state
	if inpututil.IsKeyJustPressed(ebiten.KeyR) &&
		(ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) {
		g.fullReset()
		return nil
	}

	switch g.state {
	case StateNameInput:
		return g.updateNameInput()
	case StatePlaying:
		return g.updatePlaying()
	case StateGameOver:
		return g.updateGameOver()
	}

	return nil
}

func (g *Game) updateNameInput() error {
	chars := ebiten.AppendInputChars(nil)
	for _, c := range chars {
		if c == '\n' || c == '\r' {
			continue
		}
		g.inputBuffer += string(c)
	}

	// Backspace: delete last character
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.inputBuffer) > 0 {
			g.inputBuffer = g.inputBuffer[:len(g.inputBuffer)-1]
		}
	}

	// Enter: confirm name
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.confirmName()
	}

	return nil
}

func (g *Game) confirmName() {
	name := g.inputBuffer
	if name == "" {
		// default names if empty
		if g.editingPlayerX {
			name = "X"
		} else {
			name = "O"
		}
	}

	if g.editingPlayerX {
		g.playerX.name = name
		g.state = StateNameInput
		g.editingPlayerX = false
		g.inputBuffer = ""
	} else {
		g.playerO.name = name
		g.state = StatePlaying
		g.inputBuffer = ""
	}
}

func (g *Game) updatePlaying() error {
	if !inpututil.IsKeyJustPressed(ebiten.KeyE) {
		return nil
	}

	// compute the board cell from the current player world position
	pos := g.currentPlayer.pos
	cx := int(pos.X / MapRoomStride)
	cy := int(pos.Y / MapRoomStride)

	// check if within board bounds section
	if cx < 0 || cx >= GridSize || cy < 0 || cy >= GridSize {
		return nil
	}

	// cell must be empty
	if g.board[cy][cx] != PlayerSymbolNone {
		return nil
	}

	// update the board (this is the authoritative game state)
	g.board[cy][cx] = g.currentPlayer.symbol

	// spawn a visual mark sprite at the center of the cell
	// this avoids jitter when the player is not perfectly centered in the room
	cellCenter := Vec2{
		X: (float64(cx) + HalfTile) * MapRoomStride,
		Y: (float64(cy) + HalfTile) * MapRoomStride,
	}

	g.sprites = append(g.sprites, &Sprite{
		Position:  cellCenter,
		TextureID: g.currentPlayer.symbolTextureID,
		Scale:     1.0,
		Z:         0.0,
		Hidden:    false,
	})

	winnerSym := g.board.CheckWinner()
	gameOver := winnerSym != PlayerSymbolNone || g.board.IsFull()
	if gameOver {
		g.handleGameEnd(winnerSym)
		return nil
	}

	g.switchPlayer()
	return nil
}

func (g *Game) updateGameOver() error {
	g.stateTimer -= DeltaTime
	if g.stateTimer <= 0 {
		g.switchPlayer()
		g.resetBoard()
	}
	return nil
}

func (g *Game) switchPlayer() {
	if g.currentPlayer.symbol == PlayerSymbolX {
		g.currentPlayer = g.playerO
	} else {
		g.currentPlayer = g.playerX
	}
}

func (g *Game) handleGameEnd(w PlayerSymbol) {
	g.state = StateGameOver
	g.stateTimer = GameOverDuration

	switch w {
	case PlayerSymbolX:
		g.winner = g.playerX
		g.playerX.score++
	case PlayerSymbolO:
		g.winner = g.playerO
		g.playerO.score++
	case PlayerSymbolNone:
		g.winner = nil
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBackground)

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
	return WindowSizeX, WindowSizeY
}

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
	for _, obj := range g.drawables {
		obj.Draw(screen, g)
	}
}

func (g *Game) drawText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color) {
	g.drawTextWithFace(screen, msg, x, y, align, col, g.assets.NormalTextFace, TextLineSpacing)
}

func (g *Game) drawBigText(screen *ebiten.Image, msg string, x, y float64, align TextAlign, col color.Color) {
	g.drawTextWithFace(screen, msg, x, y, align, col, g.assets.BigTextFace, BigTextLineSpacing)
}

// drawTextWithFace draws text with the given face and line spacing.
func (g *Game) drawTextWithFace(
	screen *ebiten.Image,
	msg string,
	x, y float64,
	align TextAlign,
	col color.Color,
	face *text.GoTextFace,
	lineSpacing float64,
) {
	op := &text.DrawOptions{}
	op.LineSpacing = face.Size + lineSpacing
	op.ColorScale.ScaleWithColor(col)

	w, h := text.Measure(msg, face, op.LineSpacing)

	switch align {
	case TopLeft, CenterLeft, BottomLeft:
	case TopCenter, Center, BottomCenter:
		x -= w / HalfDivisor
	case TopRight, CenterRight, BottomRight:
		x -= w
	}

	switch align {
	case TopLeft, TopCenter, TopRight:
	case CenterLeft, Center, CenterRight:
		y -= h / HalfDivisor
	case BottomLeft, BottomCenter, BottomRight:
		y -= h
	}

	op.GeoM.Translate(x, y)
	text.Draw(screen, msg, face, op)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	var msg string
	if g.winner != nil {
		symbol := "O"
		if g.winner.symbol == PlayerSymbolX {
			symbol = "X"
		}
		msg = fmt.Sprintf("%s WON", symbol)
	} else {
		msg = "DRAW"
	}

	vector.FillRect(screen, 0, 0, float32(WindowSizeX), float32(WindowSizeY), ColorGameOverBackground, false)
	g.drawBigText(
		screen,
		msg,
		float64(WindowSizeXDiv2),
		float64(WindowSizeYDiv2),
		Center,
		ColorGameOverText,
	)
}

func (g *Game) resetBoard() {
	g.board.Reset()
	g.winner = nil
	g.state = StatePlaying

	// remove mark sprites (keeping decorations like lights)
	filtered := g.sprites[:0]
	for _, s := range g.sprites {
		if s.TextureID != PlayerXSymbol && s.TextureID != PlayerOSymbol {
			filtered = append(filtered, s)
		}
	}
	g.sprites = filtered
}

func (g *Game) fullReset() {
	g.resetBoard()
	g.currentPlayer = g.playerX

	g.playerX.score = 0
	g.playerO.score = 0

	g.state = StateNameInput
	g.editingPlayerX = true
	g.inputBuffer = ""

	g.playerX.name = "X"
	g.playerO.name = "O"

	//FIXME: only clear placed marks, not lights/decorations
}
