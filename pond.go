package gofish

import (
	// "container/list"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	// "sync"
)

// The pond holds all the information about a part of
// the ocean, bounded by Bound.
// - Overlap: fish in the overlap are to report to the processes
//   that border the pond
// - Transfer: fish in transfer have moved outside the bounds
//   of the pond during this epoch
type Pond struct {
	Bound *Bound        // 2 x Point
	Fish  map[int]*Fish // doubly linked list
	// rw_lock      sync.RWMutex   // global lock on the pond
	// pond_cond    sync.Cond      // conditional wait
	// sync_lock    sync.Mutex     // the lock used by the sync.Cond
}

func NewPond(bound *Bound, fish []*Fish) (*Pond, error) {
	pond := &Pond{
		Bound: bound,
		Fish:  make(map[int]*Fish),
	}

	// make sure the fish stop and wait
	// for the pond no initialize
	// insert all fish into the map
	for _, f := range fish {
		// start fish processes here..
		pond.Fish[f.ID] = f
	}
	return pond, nil
}

func (p *Pond) Simulate(n int) {
	if len(p.Fish) == 0 || n == 0 {
		return
	}
	for i := 0; i < n; i++ {
		for _, f := range p.Fish {
			f.Step(p)
		}
		if i%100 == 0 && i > 0 {
			fmt.Printf("Round %d: Completed. Updating positions..\n", i)
		}
		for _, f := range p.Fish {
			f.UpdatePosition()
		}
	}
	fmt.Println("Simulation completed.")
}
