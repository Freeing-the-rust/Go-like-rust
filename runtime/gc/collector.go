package gc

import (
	"fmt"
	"sync"
)

type Object struct {
	Value interface{}
	mark  bool
}

type GC struct {
	heap []*Object
	lock sync.Mutex
}

func NewGC() *GC {
	return &GC{heap: make([]*Object, 0, 1024)}
}

func (g *GC) Allocate(v interface{}) *Object {
	g.lock.Lock()
	defer g.lock.Unlock()
	obj := &Object{Value: v}
	g.heap = append(g.heap, obj)
	return obj
}

func (g *GC) Collect() {
	g.lock.Lock()
	defer g.lock.Unlock()
	var alive []*Object
	for _, o := range g.heap {
		if o.mark {
			alive = append(alive, o)
		}
	}
	g.heap = alive
	fmt.Printf("[GC] Collected, remaining %d objects\n", len(g.heap))
}

func (g *GC) MarkAll() {
	for _, o := range g.heap {
		o.mark = true
	}
}
