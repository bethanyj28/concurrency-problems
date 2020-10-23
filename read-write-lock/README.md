# ReadWrite Lock

## How to run
`go build main.go` then `./main`

## FAQ

### Why?
Because.

### Should I use this in practice?
No.

### Can you explain what's going on here?
I decided to implement the read-write lock algorithms from [Wikipedia](https://en.wikipedia.org/wiki/Readers%E2%80%93writer_lock) into Go to get a better understanding of them. 

The first algorithm is known as a read-preferring readers-writer lock. This method allows for maximum concurrency because it will allow access to the lock with no preference. This usually starves the writer because multiple readers are allowed to access the lock, but only one writer can access the lock with no other threads accessing it.

The second algorithm is a write-preferring readers-writers lock. This prevents new readers from acquiring a lock if there is a writer queued and waiting. When the current readers are done, the writer will then acquire the lock.

### If I want to use this in practice, what should I use?
[This.](https://golang.org/pkg/sync/#RWMutex)
