package main

import (
	"math"
)

const (
	EPSILON float64 = 1E-6
)

// utility functions
// Maybe, instead use (fabs(x-y) < K * FLT_EPSILON * fabs(x+y))
// link: http://stackoverflow.com/a/10335601/205682
func isEqualFloat64(v1, v2 float64) bool {
	if math.Abs(v1-v2) > EPSILON {
		return false
	}
	return true
}

func notEqualFloat64(v1, v2 float64) bool {
	return !isEqualFloat64(v1, v2)
}

// directional vector
type Vector2D struct {
	X, Y float64
}

// normal (magnitude) of the vector
func (v Vector2D) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// the squared normal of the vector
func (v Vector2D) Norm2() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2D) Rotate(a float64) Vector2D {
	cosAngle, sinAngle := math.Cos(a), math.Sin(a)
	x := v.X*cosAngle - v.Y*sinAngle
	y := v.X*sinAngle + v.Y*cosAngle
	return Vector2D{x, y}
}

// angle in radians between the two vectors
func (v1 Vector2D) Angle(v2 Vector2D) float64 {
	a := math.Atan2(v2.Y, v2.X) - math.Atan2(v1.Y, v1.X)

	if a > math.Pi {
		return a - 2.0*math.Pi
	}
	if a < -math.Pi {
		return a + 2.0*math.Pi
	}
	return a
}

// The vector v normalized
func (v Vector2D) Normalize() Vector2D {
	n := v.Norm()

	if n == 0.0 {
		return v
	}
	v.X /= n
	v.Y /= n
	return v
}

func (v Vector2D) IsZero() bool {
	return isEqualFloat64(v.X, 0.0) && isEqualFloat64(v.Y, 0.0)
}

func (v Vector2D) NotZero() bool {
	return v.X != 0.0 || v.Y != 0.0
}

// Add normalized vector
func (v Vector2D) AddNorm(v2 Vector2D, norm float64) Vector2D {
	v.X += v2.X / norm
	v.Y += v2.Y / norm
	return v
}

// Add vector, no normalisation
func (v Vector2D) Add(v2 Vector2D) Vector2D {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (v Vector2D) Mult(f float64) Vector2D {
	v.X = v.X * f
	v.Y = v.Y * f
	return v
}

// subtract vector and normalize
func (v Vector2D) SubNorm(v2 Vector2D, norm float64) Vector2D {
	v.X -= v2.X / norm
	v.Y -= v2.Y / norm
	return v
}

func (v Vector2D) Sub(v2 Vector2D) Vector2D {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}
