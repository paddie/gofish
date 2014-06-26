package main

import (
	// "container/list"
	// "os"
	// "encoding/json"
	"fmt"
	// "github.com/ajstarks/svgo"
	"github.com/paddie/statedb"
	"github.com/paddie/statedb/fs"
	"github.com/paddie/statedb/monitor"
	"github.com/paddie/statedb/schedular"
	// "sync"
	"math"
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

func NewPond(bound *Bound, fish []*Fish, steps int, statPath string) (*Pond, error) {

	fs, err := fs.NewFS_OS("test")
	if err != nil {
		return nil, err
	}

	// signals a checkpoint every time the price goes up
	// - a checkpoint is signalled until the checkpoint is taken
	mdl := schedular.NewNever()
	// sends a pre-defined list of price-values to the framework
	mon := monitor.NewTestMonitor(time.Second * 5)

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

	if restored {
		fmt.Println("Restoring!")
		if err := pond.Restore(); err != nil {
			return nil, err
		}
		return pond, nil
	}

	pond.InitPond(bound, fish)

	return pond, nil
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
	sum := 0
	for {
		fish := new(Fish)
		if _, ok := iter.Next(fish); !ok {
			break
		}
		sum += fish.val
		p.fish[fish.ID] = fish
	}
	n := float64(len(p.fish))
	exp := math.Floor((n * (n + 1)) / 2)
	sumf := float64(sum)
	if sumf != exp {
		fmt.Printf("Incomplete restore: ID sum == %.2f, Expected: %.2f\n", sumf, exp)
	} else {
		fmt.Printf("Successfully restored fish with ID sum = %d\n", sum)
	}

	p.restored = true

	return nil
}

func (p *Pond) InitPond(bound *Bound, fish []*Fish) {
	fmt.Println("init")

	// insert into db
	p.db.Register(p)

	// make sure the fish stop and wait
	// for the pond no initialize
	// insert all fish into the map
	for _, f := range fish {
		// start fish processes here..
		p.db.Register(f)
		p.fish[f.ID] = f
	}
}

func (p *Pond) Simulate(freq int) {
	if len(p.fish) == 0 {
		return
	}
	for ; p.I < p.N; p.I++ {
		for _, f := range p.fish {
			f.Step(p)
		}
		if p.I%100 == 0 && p.I > 0 {
			fmt.Printf("Round %d: Completed. Updating positions..\n", p.I)
		}
		for _, f := range p.fish {
			f.UpdatePosition()
		}
		if p.I%freq == 0 {
			if err := p.db.PointOfConsistency(); err != nil {
				panic(err)
			}
		}
	}
	p.db.FinalCommit()
	fmt.Println("Simulation completed.")
}
