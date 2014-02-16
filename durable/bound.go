package gofish

import (
	"fmt"
)

const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

const (
	NO = 0
	L  = 16
	R  = 32
	U  = 4
	D  = 8
	LU = 20
	LD = 24
	RU = 36
	RD = 40
)

// Bound
// - Origo: bottom left corner
// -   End: top right corner
type Bound struct {
	O, E Point
}

// Assumes that Point p is contained by the Bound
// returns a uint8 that indicates which zones it overlaps with:
// -------------------------------
// - 0000|0000 = No Overlap = 0
// - 0001|0000 = LEFT = 16
// - 0010|0000 = RIGHT = 32
// - 0000|0100 = UP = 4
// - 0000|1000 = DOWN = 8
// -------------------------------
// - 0001|0100 = LEFT  | UP   = 20
// - 0001|1000 = LEFT  | DOWN = 24
// - 0010|0100 = RIGHT | UP   = 36
// - 0010|1000 = RIGHT | DOWN = 40
// -------------------------------
// - etc..
// TODO:
// ## Observation:
// - A point can _at most_ overlap with two zones at any given time
// - currently neglecting top left|right, bottom left|right

//  -------- -------- --------
// |        |        |        |
// |   %    |   X    |   %    |
// |        |        |        |
//  -------- -------- --------
// |        |        |        |
// |   X    |   X    |   X    |
// |        |        |        |
//  -------- -------- --------
// |        |        |        |
// |   %    |   X    |   %    |
// |        |        |        |
//  -------- -------- --------
func (b *Bound) OverlapZones(p Point, Horizon float64) uint8 {
	var zones uint8 = 0
	// left zone
	if p.X < b.O.X+Horizon {
		zones = zones | 1<<LEFT
		zones = zones << 4
	}
	// right zone
	if p.X >= b.E.X-Horizon {
		zones = zones | 1<<RIGHT
		zones = zones << 4
	}
	// upper zone
	if p.Y >= b.E.Y-Horizon {
		return zones | 1<<UP
	}
	// lower zone
	if p.Y < b.O.Y+Horizon {
		return zones | 1<<DOWN
	}

	return zones
}

func NewBound(o, e Point) (*Bound, error) {
	// TODO: make sure that origo < endpoint

	bound := &Bound{o, e}
	if !bound.IsValid() {
		return nil, fmt.Errorf("Invalid Bound: %v", bound)
	}
	return bound, nil
}

func (b *Bound) IsValid() bool {
	return (b.O.X < b.E.X && b.O.Y < b.E.Y)
}

func (b *Bound) Origo() Point {
	return b.O
}

func (b *Bound) End() Point {
	return b.E
}

// Given a position, return true if the position is not contained
// by the bound.
func (b *Bound) OutOfBound(p Point) bool {
	x, y := p.X, p.Y
	return (x < b.O.X || x >= b.E.X) || (y < b.O.Y || y >= b.E.Y)
}

func (b *Bound) WithinBound(p Point) bool {
	x, y := p.X, p.Y
	return (b.O.X <= x && x < b.E.X) && (b.O.Y <= y && y < b.E.Y)
}
