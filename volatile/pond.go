package gofish

import (
	// "container/list"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	// "github.com/paddie/statedb"
	// "sync"
	"time"
)

// The pond holds all the information about a part of
// the ocean, bounded by Bound.
// - Overlap: fish in the overlap are to report to the processes
//   that border the pond
// - Transfer: fish in transfer have moved outside the bounds
//   of the pond during this epoch
type Pond struct {
	Bound      *Bound        // 2 x Point
	Fish       map[int]*Fish // doubly linked list
	quit       chan time.Time
	jobs, done chan *Fish
}

func NewPond(bound *Bound, fish []*Fish) (*Pond, error) {

	// db := statedb.NewStateDB("", "tmp", "friday")

	// if db.rest

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

func (p *Pond) Simulate(n, procs int) {

	if len(p.Fish) == 0 || n == 0 || procs == 0 {
		return
	}
	// channel for synchonising fish completion
	p.done = make(chan *Fish, 1000)
	p.jobs = make(chan *Fish, 1000)
	p.quit = make(chan time.Time, procs)

	for i := 0; i < procs; i++ {
		go Worker(p)
	}

	sync := make(chan bool, 1)

	// Generator routine
	// 1. sends out a load of jobs and..
	// 2. waits on the sync thread until the next tick is signalled.
	go func() {
		for i := 0; i < n; i++ {
			for _, f := range p.Fish {
				p.jobs <- f
			}
			_ = <-sync
		}
		fmt.Println("Generator shutting down")
	}()

	for i := 0; i < n; i++ {
		// receive results
		for m := 0; m < len(p.Fish); m++ {
			_ = <-p.done
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

	close(p.jobs)

	for i := 0; i < procs; i++ {
		_ = <-p.quit
	}
	fmt.Println("All workers has shut down")
}

func Worker(p *Pond) {
	for f := range p.jobs {
		f.Step(p)
		// i++
		p.done <- f
	}

	p.quit <- time.Now()
}
