package gosp

import (
	"sync"
)

// SamplePool is a pool for retrieving fixed size sample buffers.
// Useful for audio processing where the buffer size if fixed.
type SamplePool[S SampleType[T], T Type] struct {
	pool sync.Pool
}

func NewSamplePool[S SampleType[T], T Type](size int) *SamplePool[S, T] {
	return &SamplePool[S, T]{
		pool: sync.Pool{
			New: func() any {
				samples := make([]S, size)
				return &samples
			},
		},
	}
}

func (p *SamplePool[S, T]) Get() *[]S {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	slice, ok := ptr.(*[]S)
	if !ok {
		return nil
	}

	return slice
}

func (p *SamplePool[S, T]) Put(samples *[]S) {
	p.pool.Put(&samples)
}
