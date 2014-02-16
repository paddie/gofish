package gofish

import (
	// "fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	O := Point{0.0, 0.0}
	E := Point{2.0, 2.0}

	b, err := NewBound(O, E)
	if err != nil {
		t.Error(err)
	}

	b, err = NewBound(E, O)
	if err == nil {
		t.Errorf("Invalid bound accepted: %v", b)
	}
}

func TestOverlapZones(t *testing.T) {
	// TODO: test for negative bound

	bound, err := NewBound(Point{0.0, 0.0}, Point{10.0, 10.0})
	if err != nil {
		t.Error(err)
	}

	type BoundTest struct {
		Desc string
		P    Point
		Exp  uint8
	}

	points := []BoundTest{
		{"LEFT", Point{1.0, 5.0}, L},        // LEFT
		{"RIGHT", Point{9.0, 5.0}, R},       // RIGHT
		{"UP", Point{5.0, 9.0}, U},          // UP
		{"DOWN", Point{5.0, 1.0}, D},        // DOWN
		{"LEFT|UP", Point{1.0, 9.0}, LU},    // LEFT | UP
		{"LEFT|DOWN", Point{1.0, 1.0}, LD},  // LEFT | DOWN
		{"RIGHT|UP", Point{9.0, 9.0}, RU},   // RIGHT | UP
		{"RIGHT|DOWN", Point{9.0, 1.0}, RD}, // RIGHT | DOWN
		{"No Overlap", Point{5.0, 5.0}, NO}, // No Overlap
	}

	for _, pt := range points {
		Act := bound.OverlapZones(pt.P, 2.0)

		if Act != pt.Exp {
			t.Errorf("%s: Exp %v != %v Act", pt.Desc, pt.Exp, Act)
		}
	}
}

func TestOutOfBounds(t *testing.T) {

	bound, err := NewBound(Point{0.0, 0.0}, Point{10.0, 10.0})
	if err != nil {
		t.Error(err)
	}

	type OOB struct {
		Desc string
		P    Point
		Exp  bool
	}

	tests := []OOB{
		{"Over", Point{5.0, 11.0}, true},
		{"Under", Point{5.0, -1.0}, true},
		{"In Bounds", Point{5.0, 5.0}, false},
		{"Left", Point{-1.0, 5.0}, true},
		{"Right", Point{11.0, 5.0}, true},
	}

	for _, pt := range tests {
		Act := bound.OutOfBound(pt.P)

		if Act != pt.Exp {
			t.Errorf("%s: Exp %v != %v Act", pt.Desc, pt.Exp, Act)
		}
	}
}

func TestWithinBounds(t *testing.T) {

	bound, err := NewBound(Point{0.0, 0.0}, Point{10.0, 10.0})
	if err != nil {
		t.Error(err)
	}

	type OOB struct {
		Desc string
		P    Point
		Exp  bool
	}

	tests := []OOB{
		{"Over", Point{5.0, 11.0}, false},
		{"Under", Point{5.0, -1.0}, false},
		{"In Bounds", Point{5.0, 5.0}, true},
		{"Left", Point{-1.0, 5.0}, false},
		{"Right", Point{11.0, 5.0}, false},
	}

	for _, pt := range tests {
		Act := bound.WithinBound(pt.P)

		if Act != pt.Exp {
			t.Errorf("%s: Exp %v != %v Act", pt.Desc, pt.Exp, Act)
		}
	}
}
