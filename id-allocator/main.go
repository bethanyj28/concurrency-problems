package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ids := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	a := allocator{free: ids, alloc: map[int64]bool{}, mux: &sync.Mutex{}}
	/*
		for i := 0; i < 20; i++ {
			id := a.allocate()
			time.Sleep(15 * time.Millisecond)
			a.release(id)
		}
	*/

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		a.testConcurrentAllocation(&wg)
	}
	wg.Wait()
}

type allocator struct {
	free  []int64
	alloc map[int64]bool
	mux   *sync.Mutex
}

func (a *allocator) testConcurrentAllocation(wg *sync.WaitGroup) {
	idChan := make(chan int64)
	go a.concurrentAllocate(idChan, wg)
	id := <-idChan

	if id != -1 {
		time.Sleep(15 * time.Millisecond)
		go a.concurrentRelease(id, wg)
	} else {
		fmt.Println("No IDs available")
	}
}

func (a *allocator) concurrentRelease(id int64, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	a.release(id)
}

func (a *allocator) concurrentAllocate(idChan chan int64, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	idChan <- a.allocate()
	close(idChan)
}

func (a *allocator) allocate() int64 { // O(1)
	a.mux.Lock()
	defer a.mux.Unlock()
	if len(a.free) == 0 {
		return -1
	}

	id := a.free[0]
	a.free = a.free[1:]
	a.alloc[id] = true
	fmt.Println("ID allocated: ", id)
	return id
}

func (a *allocator) release(id int64) { // O(1)
	a.mux.Lock()
	defer a.mux.Unlock()
	delete(a.alloc, id)
	a.free = append(a.free, id)
	fmt.Println("ID released: ", id)
}
