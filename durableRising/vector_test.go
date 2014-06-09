package main

import (
	// "fmt"
	// "math"
	"testing"
)

// Compare two vectors
func isEqualVector2D(v1, v2 Vector2D) bool {
	return isEqualFloat64(v1.X, v2.X) && isEqualFloat64(v1.Y, v2.Y)
}

func notEqualVector2D(v1, v2 Vector2D) bool {
	return !(isEqualFloat64(v1.X, v2.X) && isEqualFloat64(v1.Y, v2.Y))
}

func TestRotate(t *testing.T) {

	v := Vector2D{1.0, 1.0}
	r := v.Rotate(20)

	if notEqualVector2D(v, r.Rotate(-20)) {
		t.Errorf("Rotating %v by 20º, followed by -20º did not produce the original vector.\n", v)
	}

	if isEqualVector2D(v, r.Rotate(-21)) {
		t.Errorf("Rotating %v by 20º, followed by -21º produced the original vector.\n", v)
	}
}

func TestNormalize(t *testing.T) {

	vs := []Vector2D{
		{0.2, 9.0},
		{1.0, 1.0},
		{5.0, 4.4},
		{-3.6, -9.0},
		// {0, 0},
	}

	for _, v := range vs {
		n := v.Normalize()
		act := n.Norm()
		// fmt.Printf("Normalized: %v => |%v| = %f\n", v, n, act)
		if notEqualFloat64(act, 1.0) {
			t.Error("Normalized vector |v| != 1.0")
		}
	}
}

func TestAdd(t *testing.T) {
	v1, v2 := Vector2D{1.0, 1.0}, Vector2D{1.0, 1.0}
	expFail := Vector2D{2.000005, 2.000000000001}
	expSuccess := Vector2D{2.0000005, 2.0000000007}

	if notEqualVector2D(v1.Add(v2), expSuccess) {
		t.Errorf("%v + %v != %v", v1, v2, expSuccess)
	}

	if isEqualVector2D(v1.Add(v2), expFail) {
		t.Errorf("%v + %v == %v", v1, v2, expFail)
	}
}

func TestMult(t *testing.T) {
	v1 := Vector2D{2.0, 2.0}
	factor := 2.0
	exp := Vector2D{4.0, 4.0}

	if notEqualVector2D(v1.Mult(factor), exp) {
		t.Errorf("%v * %v != %v", v1, factor, exp)
	}
}
