// Copyright (c) 2025 Elwan Mayencourt, Masami Morimura
// SPDX-License-Identifier: Apache-2.0

package main

import "math"

// Vec2 represents a 2D vector with X and Y components.
// Methods never modify the original vector, they always return a new vector.
type Vec2 struct {
	X, Y float64
}

// Add returns a new vector addition of a and b.
func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Sub returns a new vector subtraction of a and b.
func (a Vec2) Sub(b Vec2) Vec2 { return Vec2{a.X - b.X, a.Y - b.Y} }

// Scale returns a new vector which is the original vector multiplied by scalar s.
func (a Vec2) Scale(s float64) Vec2 { return Vec2{a.X * s, a.Y * s} }

// Dot returns the dot product of vectors a and b.
func (a Vec2) Dot(b Vec2) float64 { return a.X*b.X + a.Y*b.Y }

// Len2 returns the squared length of the vector a.
func (a Vec2) Len2() float64 { return a.Dot(a) }

// Len returns the magnitude of the vector a.
func (a Vec2) Len() float64 { return math.Sqrt(a.Len2()) }

// Perp returns a new vector that is perpendicular to vector a.
func (a Vec2) Perp() Vec2 { return Vec2{-a.Y, a.X} }

// Normalize returns a new vector in the same direction as a but with a length of 1.
func (a Vec2) Normalize() Vec2 {
	l := a.Len()
	if l == 0 {
		return Vec2{}
	}
	return a.Scale(1 / l)
}

// Rotate returns a new vector which is the original vector rotated by rad radians.
func (a Vec2) Rotate(rad float64) Vec2 {
	c, s := math.Cos(rad), math.Sin(rad)
	return Vec2{
		X: a.X*c - a.Y*s,
		Y: a.X*s + a.Y*c,
	}
}

// Abs returns a new vector with the absolute values of the components of vector a.
func (a Vec2) Abs() Vec2 { return Vec2{math.Abs(a.X), math.Abs(a.Y)} }

// Sign returns a new vector with the sign of the components of vector a.
func (a Vec2) Sign() Vec2 {
	return Vec2{sign(a.X), sign(a.Y)}
}

// sign returns -1 if v is negative, 1 if v is positive, and 0 if v is zero.
func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}
