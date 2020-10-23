package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	rp := ReadPref{rmux: &sync.Mutex{}, gmux: &sync.Mutex{}, r: 0}
	fmt.Println("Read-Preferring RWLock:")
	start := time.Now()
	testRWLocker(&rp)
	fmt.Println("Finished in ", time.Since(start).Microseconds(), " ms")

	mux := &sync.Mutex{}
	cond := sync.NewCond(mux)
	wp := WritePref{mux: mux, cond: cond, activeReaders: 0, activeWriter: false, waitingWriters: 0}
	fmt.Println("Write-Preferring RWLock:")
	start = time.Now()
	testRWLocker(&wp)
	fmt.Println("Finished in ", time.Since(start).Microseconds(), " ms")
}

func testRWLocker(rwl RWLocker) {
	aMap := map[string]int{}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			rwl.AcquireWriteLock()
			aMap[fmt.Sprint(i)] = 100 - i
			rwl.ReleaseWriteLock()
			fmt.Printf("%d written: %d\n", i, 100-i)
		}(i)

		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			rwl.AcquireReadLock()
			num, ok := aMap[fmt.Sprint(i)]
			if !ok {
				num = -1
			}
			rwl.ReleaseReadLock()
			fmt.Printf("%d read: %d\n", i, num)
		}(i)
	}
	wg.Wait()
}

// RWLocker is a read/write lock interface
type RWLocker interface {
	AcquireReadLock()
	ReleaseReadLock()
	AcquireWriteLock()
	ReleaseWriteLock()
}

// ReadPref is a read/write lock implemention that is read-preferring
type ReadPref struct {
	rmux *sync.Mutex
	gmux *sync.Mutex
	r    int64
}

// AcquireReadLock ...
func (l *ReadPref) AcquireReadLock() {
	l.rmux.Lock()
	defer l.rmux.Unlock()
	l.r++
	if l.r == 1 {
		l.gmux.Lock()
	}
}

// ReleaseReadLock ...
func (l *ReadPref) ReleaseReadLock() {
	l.rmux.Lock()
	defer l.rmux.Unlock()
	l.r--
	if l.r == 0 {
		l.gmux.Unlock()
	}
}

// AcquireWriteLock ...
func (l *ReadPref) AcquireWriteLock() {
	l.gmux.Lock()
}

// ReleaseWriteLock ...
func (l *ReadPref) ReleaseWriteLock() {
	l.gmux.Unlock()
}

// WritePref is a read/write lock implementation that is write-preferring
type WritePref struct {
	mux            *sync.Mutex
	cond           *sync.Cond
	activeReaders  int64
	activeWriter   bool
	waitingWriters int64
}

// AcquireReadLock ...
func (l *WritePref) AcquireReadLock() {
	l.mux.Lock()
	for l.waitingWriters > 0 || l.activeWriter {
		l.cond.Wait()
	}

	l.activeReaders++
	l.mux.Unlock()
}

// ReleaseReadLock ...
func (l *WritePref) ReleaseReadLock() {
	l.mux.Lock()
	l.activeReaders--
	if l.activeReaders == 0 {
		l.cond.Broadcast()
	}

	l.mux.Unlock()
}

// AcquireWriteLock ...
func (l *WritePref) AcquireWriteLock() {
	l.mux.Lock()
	l.waitingWriters++
	for l.activeReaders > 0 || l.activeWriter {
		l.cond.Wait()
	}

	l.waitingWriters--
	l.activeWriter = true
	l.mux.Unlock()
}

// ReleaseWriteLock ...
func (l *WritePref) ReleaseWriteLock() {
	l.mux.Lock()
	l.activeWriter = false
	l.cond.Broadcast()
	l.mux.Unlock()
}
