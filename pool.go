package gosp

import (
	"sync"
)

type Pool[S SampleType[T], T Type] struct {
	pool *sync.Pool
}

func NewPool[S SampleType[T], T Type](size int) *Pool[S, T] {
	pool := sync.Pool{
		New: func() any {
			samples := make([]S, size)
			return &samples
		},
	}

	return &Pool[S, T]{
		pool: &pool,
	}
}

func (p *Pool[S, T]) Get() []S {
	val, ok := p.pool.Get().(*[]S)
	if !ok {
		return []S{}
	}

	return *val
}

func (p *Pool[S, T]) Put(samples []S) {
	p.pool.Put(&samples)
}
