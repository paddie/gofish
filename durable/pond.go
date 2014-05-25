package main

import (
	// "container/list"
	// "os"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	"github.com/paddie/statedb"
	"github.com/paddie/statedb/fs"
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
	db         *statedb.StateDB
	restored   bool
}

func (p *Pond) Mutable() interface{} {
	return &p.I
}

func (p *Pond) Restore() error {
	if err := p.db.RestoreSingle(p); err != nil {
		panic(err)
	}

	t := statedb.ReflectTypeM(Fish{})
	iter, err := p.db.RestoreIter(t)
	if err != nil {
		return err
	}

	fish := new(Fish)
	for {
		if _, ok := iter.Next(fish); !ok {
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
	p.db.Insert(p)

	// make sure the fish stop and wait
	// for the pond no initialize
	// insert all fish into the map
	for _, f := range fish {
		// start fish processes here..
		p.db.Insert(f)
		p.fish[f.ID] = f
	}
}

func NewPond(bound *Bound, fish []*Fish, steps int, statPath string) (*Pond, error) {

	fs, err := fs.NewFS_OS("test")
	if err != nil {
		return nil, err
	}

	// signals a checkpoint every time the price goes up
	// - a checkpoint is signalled until the checkpoint is taken
	mdl := statedb.NewRisingEdge()
	// sends a pre-defined list of price-values to the framework
	mon := statedb.NewTestMonitor(time.Second * 5)

	db, restored, err := statedb.NewStateDB(fs, mdl, mon, 2.0, statPath)
	if err != nil {
		panic(err)
	}

	pond := &Pond{
		ID:    1,
		Bound: bound,
		N:     steps,
		I:     0,
		fish:  make(map[int]*Fish),
		db:    db,
	}

	// var pond *Pond
	if restored {
		if err := pond.Restore(); err != nil {
			return nil, err
		}
		return pond, nil
	}

	pond.InitPond(bound, fish)

	return pond, nil
}

func (p *Pond) Simulate(procs, freq int) {

	if len(p.fish) == 0 || procs == 0 {
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

		if p.I%100 == 0 && p.I > 0 {
			fmt.Printf("Round %d: Completed. Updating positions..\n", p.I)
		}

		for _, f := range p.fish {
			f.UpdatePosition()
		}

		// if p.I == 500 && p.restored == false {
		// 	fmt.Print("Exit(1)\n")
		// 	os.Exit(1)
		// }

		// if i == 1 {
		// 	fmt.Println("full")
		// 	fmt.Println(p.db.FullCheckpoint())
		// } else {
		// 	fmt.Println("delta cpt")
		// 	fmt.Println(p.db.DeltaCheckpoint())
		// }
		// p.db.Sync()
		if p.I%freq == 0 {
			if err := p.db.Sync(); err != nil {
				panic(err)
			}
		}
		// signal a new round
		sync <- true
	}
	close(sync)

	p.db.Commit()

	for i := 0; i < procs; i++ {
		_ = <-p.quit
	}
	fmt.Println("GoFish: All workers has shut down")
}

func Worker(p *Pond) {
	for f := range p.jobs {
		f.Step(p)
		// i++
		p.done <- f
	}

	p.quit <- time.Now()
}
