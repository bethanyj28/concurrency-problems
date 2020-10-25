package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	mux := &sync.Mutex{}
	p := printer{mux: mux, cond: sync.NewCond(mux), curr: 0}
	go p.printFirst(&wg)
	go p.printSecond(&wg)
	go p.printThird(&wg)
	wg.Wait()
}

type printer struct {
	mux  *sync.Mutex
	cond *sync.Cond
	curr int
}

func (p *printer) printFirst(wg *sync.WaitGroup) {
	defer wg.Done()
	p.mux.Lock()
	defer p.mux.Unlock()
	fmt.Println("first!")
	p.curr = 2
	p.cond.Broadcast()
}

func (p *printer) printSecond(wg *sync.WaitGroup) {
	defer wg.Done()
	p.mux.Lock()
	defer p.mux.Unlock()
	for p.curr != 2 {
		p.cond.Wait()
	}
	fmt.Println("second")
	p.curr = 3
	p.cond.Broadcast()
}

func (p *printer) printThird(wg *sync.WaitGroup) {
	defer wg.Done()
	p.mux.Lock()
	defer p.mux.Unlock()
	for p.curr != 3 {
		p.cond.Wait()
	}
	fmt.Println("third")
}
