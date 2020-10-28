package main

import (
	"fmt"
	"sync"
)

func main() {
	mux := &sync.Mutex{}
	h2o := h2oBuilder{currH: 0, currO: 0, mux: mux, cond: sync.NewCond(mux)}

	molecules := []string{"O", "O", "H", "H", "H", "H", "O", "H", "O", "H", "H", "H"}
	wg := sync.WaitGroup{}
	wg.Add(len(molecules))
	for _, m := range molecules {
		if m == "H" {
			go h2o.releaseHydrogen(&wg)
			continue
		}
		go h2o.releaseOxygen(&wg)
	}

	wg.Wait()
	fmt.Println()
}

type h2oBuilder struct {
	currH int
	currO int
	mux   *sync.Mutex
	cond  *sync.Cond
}

func (b *h2oBuilder) releaseHydrogen(wg *sync.WaitGroup) {
	defer wg.Done()
	b.mux.Lock()
	defer b.mux.Unlock()
	for b.currH >= 2 {
		b.cond.Wait()
	}

	fmt.Printf("H")
	b.currH++
	if b.currH == 2 && b.currO == 1 {
		b.currH = 0
		b.currO = 0
		b.cond.Broadcast()
	}
}

func (b *h2oBuilder) releaseOxygen(wg *sync.WaitGroup) {
	defer wg.Done()
	b.mux.Lock()
	defer b.mux.Unlock()
	for b.currO >= 1 {
		b.cond.Wait()
	}

	fmt.Printf("O")
	b.currO++
	if b.currH == 2 && b.currO == 1 {
		b.currH = 0
		b.currO = 0
		b.cond.Broadcast()
	}
}
