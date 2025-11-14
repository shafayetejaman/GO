package main

import (
	"sync"
	"time"
)

type safeCounter struct {
	counts map[string]int
	mu     *sync.Mutex
}

// var mut sync.Mutex
mut := sync.Mutex{}

func (sc safeCounter) inc(key string) {
	// ?
	// mut.Lock()
	sc.slowIncrement(key)
	// mut.Unlock()
}

func (sc safeCounter) val(key string) int {
	// ?
	// defer mut.Unlock()
	// mut.Lock()
	return sc.slowVal(key)
}

// don't touch below this line

func (sc safeCounter) slowIncrement(key string) {
	tempCounter := sc.counts[key]
	time.Sleep(time.Microsecond)
	tempCounter++
	sc.counts[key] = tempCounter
}

func (sc safeCounter) slowVal(key string) int {
	time.Sleep(time.Microsecond)
	return sc.counts[key]
}
