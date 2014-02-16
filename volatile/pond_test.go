package gofish

import (
	// "github.com/paddie/fish"
	// "github.com/paddie/statedb"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"
)

var pond *Pond

func init() {

	runtime.GOMAXPROCS(2)

	count := 1000
	pct := 10.0
	factor := int(float64(count) / pct)
	speed := 0.3

	fmt.Println("informed factor: ", factor)

	c := Point{500, 160}

	ps := RandomPoints(c, 200, count)

	fish := make([]*Fish, 0, count)

	for i, p := range ps {
		dir := Vector2D{0.1, 0.0}
		info := Vector2D{0.0, 5.0}
		tmp := NewFish(i+1, p, dir, speed)

		if (i+1)%factor == 0 {
			tmp.Inform(info)
		}
		fish = append(fish, tmp)

	}

	bound, _ := NewBound(Point{0, 0}, Point{1000, 1000})
	pond, _ = NewPond(bound, fish)
}

func RandomPoints(c Point, maxRadius float64, count int) []Point {
	if count <= 0 {
		return nil
	}
	points := make([]Point, 0, count)
	var u, v float64

	v = rand.Float64()

	for i := 0; i < count; i++ {
		u = v
		v := rand.Float64()

		w := maxRadius * math.Sqrt(u)
		t := 2 * math.Pi * v

		x := w * math.Cos(t)
		y := w * math.Sin(t)

		p := Point{c.X + x, c.Y + y}
		points = append(points, p)
	}
	return points
}

func TestConcurrentSimulation(t *testing.T) {
	steps := 200

	pond.Simulate(steps, 100)
}
