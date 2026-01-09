// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

type Board [GridSize][GridSize]PlayerSymbol

// Reset clears the board to its initial empty state.
func (b *Board) Reset() {
	*b = Board{}
}

// CheckWinner checks the board for a winner.
func (b *Board) CheckWinner() Winner {
	for i := range GridSize {
		// check rows
		if w := b.checkLine(b[i][0], b[i][1], b[i][2]); w != WinnerNone {
			return w
		}

		// check columns
		if w := b.checkLine(b[0][i], b[1][i], b[2][i]); w != WinnerNone {
			return w
		}
	}

	// check diagonals
	if w := b.checkLine(b[0][0], b[1][1], b[2][2]); w != WinnerNone {
		return w
	}
	if w := b.checkLine(b[0][2], b[1][1], b[2][0]); w != WinnerNone {
		return w
	}

	if b.IsFull() {
		return WinnerDraw
	}
	return WinnerNone
}

// IsFull returns true if the board is full (no empty cells).
func (b *Board) IsFull() bool {
	for y := range GridSize {
		for x := range GridSize {
			if b[y][x] == PlayerSymbolNone {
				return false
			}
		}
	}
	return true
}

// checkLine checks if the three given cells form a winning line.
func (b *Board) checkLine(c1, c2, c3 PlayerSymbol) Winner {
	if c1 != PlayerSymbolNone && c1 == c2 && c2 == c3 {
		if c1 == PlayerSymbolX {
			return WinnerX
		}
		return WinnerO
	}
	return WinnerNone
}
