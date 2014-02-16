package gofish

import (
	// "fmt"
	"math"
	"math/rand"
	// "sync"
	// "time"
)

const (
	// Personal space
	// - distance is less than alpha
	//   => repelled
	ALPHA  float64 = 5.0
	ALPHA2 float64 = 5.0 * 5.0
	// Attraction Space
	// - distance is Alpha and Rho
	//   => attracts
	RHO         float64 = 160.0         // 160.0
	RHO2        float64 = 160.0 * 160.0 // 160.0 * 160.0
	GAUSS_MU    float64 = 0.0           // 0.0
	GAUSS_SIGMA float64 = 0.01          // 0.01
	THETA       float64 = 0.08          // 0.08
	P           float64 = 0.15          // 0.15
	OMEGA       float64 = 0.5           // 0.5
	STREAM      int     = 3             // 3
	ATTRACTING  int     = 1             // fish is actively attracting
	AVOIDING    int     = -1            // fish is actively avoiding
	NOTHING     int     = 0             // fish is moving along V
)

// eucledian position in
type Point struct {
	X, Y float64
}

func (p1 Point) Diff(p2 Point) Vector2D {
	return Vector2D{p1.X - p2.X, p1.Y - p2.Y}
}

func (p Point) Move(v Vector2D) Point {
	p.X += v.X
	p.Y += v.Y

	return p
}

type Fish struct {
	ID       int      // unique id of fish
	C        Point    // position of fish
	d        Vector2D // new direction
	V        Vector2D // current direction
	G        Vector2D // preferred direction
	S        float64  // speed of fish
	Informed bool     // informed
	N_step   int      // step count
	dirty    bool     // true if C has not been updated after a step
	action   int
	awake    bool
}

func (f *Fish) Pos() *Point {
	return &f.C
}

func (f1 *Fish) Diff(f2 *Fish) Vector2D {
	return f1.C.Diff(f2.C)
}

func NewInformedFish(id int, pos Point, dir Vector2D, info Vector2D, speed float64) *Fish {

	if dir.IsZero() || info.IsZero() {
		panic("Direction vector must not be {0.0, 0.0}")
	}

	v_norm := dir.Normalize()

	return &Fish{
		ID:       id,
		C:        pos,
		V:        v_norm,
		d:        v_norm,
		G:        info.Normalize(), // unit vector
		S:        speed,
		Informed: true,
	}
}

func NewFish(id int, pos Point, dir Vector2D, speed float64) *Fish {

	if dir.IsZero() {
		panic("Direction vector must not be {0.0, 0.0}")
	}

	v_norm := dir.Normalize()

	return &Fish{
		ID:       id,
		C:        pos,
		V:        v_norm,
		G:        v_norm, // unit vector
		d:        v_norm,
		S:        speed,
		Informed: false,
	}
}

func (f1 *Fish) Avoid(pond *Pond) (Vector2D, bool) {
	avoided := 0
	qd := Vector2D{}

	for _, f2 := range pond.fish {
		if f1.ID != f2.ID && WithinRange(f1.C, f2.C, ALPHA) {
			// get difference vector
			d_diff := f2.Diff(f1)
			// normal squared
			d2 := d_diff.Norm2()
			// more precise calculation
			if d2 > 0 && d2 < ALPHA2 {
				// subtract the normalized vector
				// fmt.Printf("f%d: Avoided f%d\n", f1.ID, f2.ID)

				qd = qd.SubNorm(d_diff, math.Sqrt(d2))
				avoided++
			}
		}
	}

	if avoided > 0 {
		f1.action = AVOIDING
		return qd, true
	}

	return qd, false
}

func (f1 *Fish) Attract(pond *Pond) Vector2D {
	qd := Vector2D{}

	attracted := 0

	for _, f2 := range pond.fish {
		if f1.ID != f2.ID && WithinRange(f1.C, f2.C, RHO) {

			// d_diff := f1.Diff(f2)
			d_diff := f2.Diff(f1)
			d2 := d_diff.Norm2()
			if d2 > 0 && d2 < RHO2 {
				qd = qd.AddNorm(d_diff, math.Sqrt(d2))
				attracted++
			}

			qd = qd.Add(f1.V)
		}
	}

	if attracted > 0 {
		f1.action = ATTRACTING
		return qd
	}

	return qd
}

func (f *Fish) MakeInformed(g Vector2D) {

	if g.IsZero() {
		panic("MakeInformed: zero vector")
	}

	f.Informed = true
	f.G = g.Normalize()
}

func (f *Fish) Inform(qd Vector2D) Vector2D {
	return qd.Add(f.G.Mult(OMEGA))
}

// Use UpdatePosition to update the position of the fish after a
// completed Step
func (f *Fish) UpdatePosition() {
	// update direction vector
	// vx, vy = dx, dy
	f.V = f.d
	// update position
	// private state float x : (x + (vx * speed));
	// private state float y : (y + (vy * speed));
	// - Move point using the vector:
	f.C = f.C.Move(f.V.Mult(f.S))

	f.dirty = false
}

// could be iter, but I'll go with step for now
func (f *Fish) Step(pond *Pond) {

	f.action = NOTHING

	qd, avoided := f.Avoid(pond)

	// attract if we're not avoiding
	if !avoided {
		// fmt.Printf("f%d: didn't avoid anyone, attracting!\n", f.ID)
		qd = f.Attract(pond)
	}

	if qd.NotZero() {
		qd = qd.Normalize()
		if f.Informed {
			qd = f.Inform(qd)
			qd = qd.Normalize()
		}
	} else {
		qd = f.V
	}

	// rotate fish with signal noise
	rotationAngle := GAUSS_MU + GAUSS_SIGMA*rand.NormFloat64()

	qd = qd.Rotate(rotationAngle)

	// ceil rotation within +/- THEATA
	// get angle between the current direction, and the new direction
	angle := f.V.Angle(qd)

	// if angle is above the maximum turning speed pr. tick THETA, limit angle to THETA
	if (angle > 0 && angle > THETA) || (angle <= 0 && angle < -THETA) {

		// use maximum rotation
		var maxRot float64
		if angle >= 0 {
			maxRot = THETA
		} else {
			maxRot = -THETA
		}
		qd = f.V.Rotate(maxRot)
	}

	f.d = qd

	// all the meta data updates
	// update step count/iteration id
	f.N_step++
	// mark fish as dirty
	f.dirty = true
}

// square around the fish
func WithinRange(c1, c2 Point, r float64) bool {
	return (c2.X < c1.X+r && c2.X > c1.X-r) && (c2.Y < c1.Y+r && c2.Y > c1.Y-r)
}
