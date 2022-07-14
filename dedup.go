// -*- tab-width: 2 -*-

// Package dedupcount lets you add a bunch of keys in and then
// quickly.  get out the keys that are added >1 times
package dedupcount

import (
	"sync"
)

// Dedup lets you feed in a lot of strings and later et the list of
// dups. The API is non-blocking (channel send) and all the keys are
// kept in memory (not some blooming thing)
type Dedup struct {
	lock    sync.RWMutex
	map1    map[string]int // for keys seen 1 time
	mapN    map[string]int // for keys seen >1 time
	addChan chan string
	done    chan bool
	name    string
}

// New give you a new Dedup context
func New(Name string) (d *Dedup) {
	// Bolt DB for done Objects for copy
	d = &Dedup{}
	d.lock = sync.RWMutex{}
	d.name = Name
	d.map1 = make(map[string]int)
	d.mapN = make(map[string]int)
	d.done = make(chan bool, 2)
	d.addChan = make(chan string, 1000000) // should be config driven

	go func() { // async writer
		for {
			select {
			case <-d.done:
				return
			case s := <-d.addChan:
				d.lock.Lock()
				_, ok2 := d.mapN[s]
				if ok2 {
					d.mapN[s]++
				} else {
					_, ok1 := d.map1[s]
					if ok1 {
						delete(d.map1, s)
						d.mapN[s] = 1
					} else {
						d.map1[s] = 1
					}
				}
				d.lock.Unlock()
			}
		}
	}()
	return d
}

// Close shutdowns the goroutine
func (d *Dedup) Close() {
	d.done <- true
}

// Add puts the string into the set eventually
func (d *Dedup) Add(s string) {
	d.addChan <- s
}

// GetDups returns the map with the duplicated things
func (d *Dedup) GetDups() map[string]int {
	return d.mapN
}

// InSet returns true if the string is in the set DB
// currently it doesn't look at pending Set
func (d *Dedup) InSet(s string) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	_, ok1 := d.map1[s]
	_, ok2 := d.mapN[s]
	return bool(ok1 || ok2)
}
