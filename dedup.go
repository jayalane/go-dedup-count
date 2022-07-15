// -*- tab-width: 2 -*-

// Package dedupcount lets you add a bunch of keys in and then
// quickly.  get out the keys that are added >1 times
package dedupcount

import (
	"sync"
)

type entry struct {
	k string
	v interface{}
}

// Dedup lets you feed in a lot of strings and later et the list of
// dups. The API is non-blocking (channel send) and all the keys are
// kept in memory (not some blooming thing)
type Dedup struct {
	lock    sync.RWMutex
	map1    map[string]interface{}   // for keys seen 1 time
	mapN    map[string][]interface{} // for keys seen >1 time
	addChan chan entry
	done    chan bool
	name    string
}

// New give you a new Dedup context
func New(Name string) (d *Dedup) {
	d = &Dedup{}
	d.lock = sync.RWMutex{}
	d.name = Name
	d.map1 = make(map[string]interface{})
	d.mapN = make(map[string][]interface{})
	d.done = make(chan bool, 2)
	d.addChan = make(chan entry, 1000000) // should be config driven

	go func() { // async writer
		for {
			select {
			case <-d.done:
				return
			case s := <-d.addChan:
				d.lock.Lock()
				_, ok2 := d.mapN[s.k]
				if ok2 {
					d.mapN[s.k] = append(d.mapN[s.k], s.v)
				} else {
					_, ok1 := d.map1[s.k]
					if ok1 {
						delete(d.map1, s.k)
						l := make([]interface{}, 0)
						d.mapN[s.k] = append(l, s.v)

					} else {
						d.map1[s.k] = s.v
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

// Set puts the string into the set eventually
func (d *Dedup) Set(k string, v interface{}) {
	d.addChan <- entry{k, v}
}

// GetDups returns a copy of the map with the duplicated things
func (d *Dedup) GetDups() map[string][]interface{} {
	rt := make(map[string][]interface{})
	d.lock.RLock()
	defer d.lock.RUnlock()
	for k, v := range d.mapN {
		rt[k] = make([]interface{}, len(v))
		for i := range v {
			rt[k] = append(rt[k], i)
		}
	}
	return rt
}

// KeySetP returns true if the string is in the set DB
// currently it doesn't look at pending Set
func (d *Dedup) KeySetP(k string) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	_, ok1 := d.map1[k]
	_, ok2 := d.mapN[k]
	return bool(ok1 || ok2)
}

// Get returns the value and true if the string is in the maps
func (d *Dedup) Get(s string) (interface{}, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	v1, ok1 := d.map1[s]
	if ok1 {
		return v1, true
	}
	v2, ok2 := d.mapN[s]
	if ok2 {
		return v2, true
	}
	return nil, false
}
