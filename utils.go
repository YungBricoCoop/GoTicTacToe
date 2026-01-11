// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/fonts/*.ttf
var fontsFS embed.FS

type GameState int

const (
	StateNameInput GameState = iota
	StatePlaying
	StateGameOver
)

type Winner int

const (
	WinnerNone Winner = iota
	WinnerX
	WinnerO
	WinnerDraw
)

type Assets struct {
	XImage         *ebiten.Image
	OImage         *ebiten.Image
	NormalTextFace *text.GoTextFace
	BigTextFace    *text.GoTextFace
	Textures       map[uint8][]*ebiten.Image

	// HUD images
	PlayerHUDImage *ebiten.Image
	WASDHUDImage   *ebiten.Image
}

type Game struct {
	// state
	state         GameState
	board         [GridSize][GridSize]PlayerSymbol
	currentPlayer *Player
	winner        Winner

	// input state
	inputBuffer    string
	editingPlayerX bool

	// visuals
	assets *Assets

	// EbitenUI (HUD)
	ui *ebitenui.UI

	// raycasting world (minimap)
	worldMap Map
	minimap  Minimap

	gameObjects []GameObject
}

func NewGame() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}

	spawnX := defaultPlayerXSpawn()
	spawnO := defaultPlayerOSpawn()
	pX := NewPlayer(spawnX.X, spawnX.Y, PlayerSymbolX, "X")
	pO := NewPlayer(spawnO.X, spawnO.Y, PlayerSymbolO, "O")
	world := World{}

	g := &Game{
		state:          StateNameInput,
		currentPlayer:  pX,
		editingPlayerX: true,
		assets:         assets,
		worldMap:       NewMap(),
	}

	g.gameObjects = append(g.gameObjects, &world, &g.minimap, pX, pO)

	// UI init (ここで WASDHUDImage が nil だと widget.NewGraphic が panic する)
	g.initUI()

	return g, nil
}

// loadEmbeddedImage loads an image from the embedded texturesFS (defined in texture.go).
func loadEmbeddedImage(path string) (*ebiten.Image, error) {
	f, err := texturesFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("decode %q: %w", path, err)
	}
	return img, nil
}

func loadAssets() (*Assets, error) {
	// ---- pixel font ----
	fontBytes, err := fontsFS.ReadFile("assets/fonts/PressStart2P-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("load pixel font: %w", err)
	}

	pixelFontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return nil, fmt.Errorf("create pixel font: %w", err)
	}

	// ---- textures (raycasting) ----
	textures, err := LoadTextures()
	if err != nil {
		return nil, fmt.Errorf("load textures: %w", err)
	}

	// ---- HUD images (IMPORTANT) ----
	// NOTE: these files MUST exist in assets/textures and are embedded by texture.go's texturesFS.
	playerHUD, err := loadEmbeddedImage("assets/textures/1.png")
	if err != nil {
		return nil, fmt.Errorf("load PlayerHUDImage: %w", err)
	}

	wasdHUD, err := loadEmbeddedImage("assets/textures/wasd.png")
	if err != nil {
		return nil, fmt.Errorf("load WASDHUDImage: %w", err)
	}

	return &Assets{
		NormalTextFace: &text.GoTextFace{
			Source: pixelFontSource,
			Size:   14, // pixel perfect
		},
		BigTextFace: &text.GoTextFace{
			Source: pixelFontSource,
			Size:   20,
		},
		Textures:       textures,
		PlayerHUDImage: playerHUD,
		WASDHUDImage:   wasdHUD,
	}, nil
}

func (g *Game) resetBoard() {
	g.board = [GridSize][GridSize]PlayerSymbol{}
	g.winner = WinnerNone
	g.state = StatePlaying
	g.currentPlayer = g.getPlayer(PlayerSymbolX)
}

func (g *Game) fullReset() {
	g.resetBoard()
	pX := g.getPlayer(PlayerSymbolX)
	pO := g.getPlayer(PlayerSymbolO)
	pX.score = 0
	pO.score = 0
	g.state = StateNameInput
	g.editingPlayerX = true
	g.inputBuffer = ""
	pX.name = "X"
	pO.name = "O"
}

func (g *Game) getPlayer(s PlayerSymbol) *Player {
	for _, obj := range g.gameObjects {
		if p, ok := obj.(*Player); ok && p.symbol == s {
			return p
		}
	}
	return nil
}

func (g *Game) checkWinner() Winner {
	for i := range GridSize {
		if w := g.checkLine(g.board[i][0], g.board[i][1], g.board[i][2]); w != WinnerNone {
			return w
		}
		if w := g.checkLine(g.board[0][i], g.board[1][i], g.board[2][i]); w != WinnerNone {
			return w
		}
	}

	if w := g.checkLine(g.board[0][0], g.board[1][1], g.board[2][2]); w != WinnerNone {
		return w
	}
	if w := g.checkLine(g.board[0][2], g.board[1][1], g.board[2][0]); w != WinnerNone {
		return w
	}

	if g.isBoardFull() {
		return WinnerDraw
	}
	return WinnerNone
}

func (g *Game) checkLine(a, b, c PlayerSymbol) Winner {
	if a != PlayerSymbolNone && a == b && b == c {
		if a == PlayerSymbolX {
			return WinnerX
		}
		return WinnerO
	}
	return WinnerNone
}

func (g *Game) isBoardFull() bool {
	for y := range GridSize {
		for x := range GridSize {
			if g.board[y][x] == PlayerSymbolNone {
				return false
			}
		}
	}
	return true
}
