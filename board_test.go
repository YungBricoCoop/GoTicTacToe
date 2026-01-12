package main

import (
	"testing"
)

func TestBoard_CheckWinner_Rows(t *testing.T) {
	tests := []struct {
		name     string
		board    Board
		expected PlayerSymbol
	}{
		{
			name: "Row 0 Winner X",
			board: Board{
				{PlayerSymbolX, PlayerSymbolX, PlayerSymbolX},
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolNone},
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolNone},
			},
			expected: PlayerSymbolX,
		},
		{
			name: "Row 1 Winner O",
			board: Board{
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolNone},
				{PlayerSymbolO, PlayerSymbolO, PlayerSymbolO},
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolNone},
			},
			expected: PlayerSymbolO,
		},
		{
			name: "No Winner",
			board: Board{
				{PlayerSymbolX, PlayerSymbolO, PlayerSymbolX},
				{PlayerSymbolO, PlayerSymbolX, PlayerSymbolO},
				{PlayerSymbolO, PlayerSymbolX, PlayerSymbolO},
			},
			expected: PlayerSymbolNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.CheckWinner(); got != tt.expected {
				t.Errorf("CheckWinner() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoard_CheckWinner_Columns(t *testing.T) {
	tests := []struct {
		name     string
		board    Board
		expected PlayerSymbol
	}{
		{
			name: "Col 0 Winner X",
			board: Board{
				{PlayerSymbolX, PlayerSymbolNone, PlayerSymbolNone},
				{PlayerSymbolX, PlayerSymbolNone, PlayerSymbolNone},
				{PlayerSymbolX, PlayerSymbolNone, PlayerSymbolNone},
			},
			expected: PlayerSymbolX,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.CheckWinner(); got != tt.expected {
				t.Errorf("CheckWinner() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoard_CheckWinner_Diagonals(t *testing.T) {
	tests := []struct {
		name     string
		board    Board
		expected PlayerSymbol
	}{
		{
			name: "Diagonal TopLeft-BottomRight Winner X",
			board: Board{
				{PlayerSymbolX, PlayerSymbolNone, PlayerSymbolNone},
				{PlayerSymbolNone, PlayerSymbolX, PlayerSymbolNone},
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolX},
			},
			expected: PlayerSymbolX,
		},
		{
			name: "Diagonal TopRight-BottomLeft Winner O",
			board: Board{
				{PlayerSymbolNone, PlayerSymbolNone, PlayerSymbolO},
				{PlayerSymbolNone, PlayerSymbolO, PlayerSymbolNone},
				{PlayerSymbolO, PlayerSymbolNone, PlayerSymbolNone},
			},
			expected: PlayerSymbolO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.CheckWinner(); got != tt.expected {
				t.Errorf("CheckWinner() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBoard_IsFull(t *testing.T) {
	b := Board{
		{PlayerSymbolX, PlayerSymbolO, PlayerSymbolX},
		{PlayerSymbolX, PlayerSymbolO, PlayerSymbolX},
		{PlayerSymbolO, PlayerSymbolX, PlayerSymbolO},
	}
	if !b.IsFull() {
		t.Error("Expected board to be full")
	}

	b[1][1] = PlayerSymbolNone
	if b.IsFull() {
		t.Error("Expected board to NOT be full")
	}
}

func TestBoard_Reset(t *testing.T) {
	b := Board{
		{PlayerSymbolX, PlayerSymbolX, PlayerSymbolX},
		{PlayerSymbolX, PlayerSymbolX, PlayerSymbolX},
		{PlayerSymbolX, PlayerSymbolX, PlayerSymbolX},
	}
	b.Reset()

	for y := range GridSize {
		for x := range GridSize {
			if b[y][x] != PlayerSymbolNone {
				t.Errorf("Cell (%d,%d) is not empty after Reset", x, y)
			}
		}
	}
}
