package gofish

import (
	// "container/list"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	"sync"
	// "time"
)

// The pond holds all the information about a part of
// the ocean, bounded by Bound.
// - Overlap: fish in the overlap are to report to the processes
//   that border the pond
// - Transfer: fish in transfer have moved outside the bounds
//   of the pond during this epoch
type Pond struct {
	Bound        *Bound         // 2 x Point
	Fish         map[int]*Fish  // doubly linked list
	fish_barrier sync.WaitGroup // barrier that pond waits on

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

// func (p *Pond) FishCompletedStep() {
// 	p.fish_barrier.Done()
// }

func (p *Pond) Simulate(n, procs int) {

	if len(p.Fish) == 0 || n == 0 || procs == 0 {
		return
	}
	// channel for synchonising fish completion
	done := make(chan *Fish, 1000)
	jobs := make(chan *Fish, 1000)
	quit := make(chan bool, procs)
	for i := 0; i < procs; i++ {
		go Worker(p, jobs, done, quit)
	}

	sync := make(chan bool, 1)

	// Generator routine
	// - sends out a load of jobs, and waits until all results have been processed before sending out a new batch
	go func() {
		for i := 0; i < n; i++ {
			for _, f := range p.Fish {
				jobs <- f
			}
			_ = <-sync
		}
		fmt.Println("Generator shutting down")
	}()

	for i := 0; i < n; i++ {
		// receive results
		for m := 0; m < len(p.Fish); m++ {
			_ = <-done
			// fmt.Printf("fish %d was processed..\n", f.ID)
		}

		if i%100 == 0 && i > 0 {
			fmt.Printf("Round %d: Completed. Updating positions..\n", i)
		}

		for _, f := range p.Fish {
			f.UpdatePosition()
		}
		// signal a new round
		sync <- true
	}

	close(jobs)

	for i := 0; i < procs; i++ {
		_ = <-quit
	}
	fmt.Println("All workers has shut down")

}

func Worker(p *Pond, job <-chan *Fish, done chan<- *Fish, quit chan<- bool) {
	for f := range job {
		f.Step(p)
		// i++
		done <- f
	}

	quit <- true
}
