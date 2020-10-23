package main

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)

func main() {
	forks := []*semaphore.Weighted{}
	for i := 0; i < 5; i++ {
		forks = append(forks, semaphore.NewWeighted(int64(10)))
	}
	dt := DiningTable{
		forks: forks,
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		go dt.Eat(context.Background(), int64(i%5), &wg)
	}

	wg.Wait()
}

// DiningTable represents the philosophers and their forks
type DiningTable struct {
	forks []*semaphore.Weighted
}

// Eat simulates a philosopher picking up two forks and eat
func (d *DiningTable) Eat(ctx context.Context, p int64, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	f1 := p - 1
	if p == 0 {
		f1 = int64(len(d.forks)) - 1
	}

	f2 := p

	if err := d.forks[f1].Acquire(ctx, 1); err != nil {
		fmt.Printf("%d failed to aquire fork %d\n", p, f1)
		return
	}
	defer d.forks[f1].Release(1)

	if err := d.forks[f2].Acquire(ctx, 1); err != nil {
		fmt.Printf("%d failed to aquire fork %d\n", p, f2)
		return
	}
	defer d.forks[f2].Release(1)

	fmt.Printf("%d is eating now with forks %d and %d\n", p, f1, f2)
}
