package singleflight

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

var errGoexit = errors.New("runtime.Goexit was called")

type panicError struct {
	value interface{}
	stack []byte
}

func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func (p *panicError) Unwrap() error {
	err, ok := p.value.(error)
	if !ok {
		return nil
	}

	return err
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

type call[V any] struct {
	wg sync.WaitGroup

	val V
	err error

	dups  int
	chans []chan<- Result[V]
}

type Group[V any] struct {
	mu sync.Mutex
	m  map[string]*call[V]
}

type Result[V any] struct {
	Val    V
	Err    error
	Shared bool
}

func (g *Group[V]) Do(key fmt.Stringer, fn func() (V, error)) (v V, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[V])
	}
	ks := key.String()
	if c, ok := g.m[ks]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call[V])
	c.wg.Add(1)
	g.m[ks] = c
	g.mu.Unlock()

	g.doCall(c, ks, fn)
	return c.val, c.err, c.dups > 0
}

func (g *Group[V]) DoChan(key fmt.Stringer, fn func() (V, error)) <-chan Result[V] {
	ch := make(chan Result[V], 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[V])
	}
	ks := key.String()
	if c, ok := g.m[ks]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call[V]{chans: []chan<- Result[V]{ch}}
	c.wg.Add(1)
	g.m[ks] = c
	g.mu.Unlock()

	go g.doCall(c, ks, fn)

	return ch
}

func (g *Group[V]) doCall(c *call[V], key string, fn func() (V, error)) {
	normalReturn := false
	recovered := false

	defer func() {
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				go panic(e)
				select {}
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
		} else {
			for _, ch := range c.chans {
				ch <- Result[V]{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}

func (g *Group[V]) Forget(key fmt.Stringer) {
	g.mu.Lock()
	delete(g.m, key.String())
	g.mu.Unlock()
}
