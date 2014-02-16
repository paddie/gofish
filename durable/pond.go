package gofish

import (
	// "container/list"
	"os"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	. "github.com/paddie/statedb"
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
	ID         int
	Bound      *Bound // 2 x Point
	N, I       int
	fish       map[int]*Fish // doubly linked list
	quit       chan time.Time
	jobs, done chan *Fish
	db         *StateDB
	restored   bool
}

func (p *Pond) Restore() error {
	fmt.Println("restoring..")
	if err := p.db.RestoreSingle(p, &p.I); err != nil {
		panic(err)
	}

	iter, err := p.db.RestoreIter(ReflectType(Fish{}))
	if err != nil {
		return err
	}

	fish := new(Fish)
	for {
		if _, ok := iter.Next(fish, &fish.C); !ok {
			break
		}
		p.fish[fish.ID] = fish
		fish = new(Fish)
	}
	p.restored = true

	return nil
}

func (p *Pond) InitPond(bound *Bound, fish []*Fish) {
	fmt.Println("init")

	// insert into db
	p.db.Insert(p, &p.I)

	// make sure the fish stop and wait
	// for the pond no initialize
	// insert all fish into the map
	for _, f := range fish {
		// start fish processes here..
		p.db.Insert(f, &f.C)
		p.fish[f.ID] = f
	}
}

func NewPond(bound *Bound, fish []*Fish) (*Pond, error) {

	db, err := NewStateDB("", "tmp", "friday")
	if err != nil {
		panic(err)
	}

	pond := &Pond{
		ID:    1,
		Bound: bound,
		N:     1000,
		I:     0,
		fish:  make(map[int]*Fish),
		db:    db,
	}

	// var pond *Pond
	if db.Restored {
		if err := pond.Restore(); err != nil {
			return nil, err
		}
		return pond, nil
	}

	pond.InitPond(bound, fish)

	return pond, nil
}

func (p *Pond) Simulate(n, procs int) {

	if len(p.fish) == 0 || n == 0 || procs == 0 {
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
		for {
			for _, f := range p.fish {
				p.jobs <- f
			}
			if _, ok := <-sync; !ok {
				close(p.jobs)
				break
			}
		}
		fmt.Println("Generator shutting down")
	}()

	for ; p.I < p.N; p.I++ {
		// receive results
		for m := 0; m < len(p.fish); m++ {
			_ = <-p.done
			// fmt.Printf("fish %d was processed..\n", f.ID)
		}

		if p.I%10 == 0 && p.I > 0 {
			fmt.Printf("Round %d: Completed. Updating positions..\n", p.I)
		}

		for _, f := range p.fish {
			f.UpdatePosition()
		}

		if p.I == 500 && p.restored == false {
			fmt.Print("Exit(1)\n")
			os.Exit(1)
		}

		// if i == 1 {
		// 	fmt.Println("full")
		// 	fmt.Println(p.db.FullCheckpoint())
		// } else {
		// 	fmt.Println("delta cpt")
		// 	fmt.Println(p.db.DeltaCheckpoint())
		// }
		// p.db.Sync()
		if err := p.db.Checkpoint(); err != nil {
			panic(err)
		}

		// signal a new round
		sync <- true
	}
	close(sync)

	p.db.Commit()

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
