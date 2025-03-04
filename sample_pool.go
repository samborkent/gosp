package gosp

import (
	"sync"
)

// SamplePool is a pool for retrieving fixed size sample buffers.
// Useful for audio processing where the buffer size if fixed.
type SamplePool[F Frame[T], T Type] struct {
	pool sync.Pool
}

func NewSamplePool[F Frame[T], T Type](size int) *SamplePool[F, T] {
	return &SamplePool[F, T]{
		pool: sync.Pool{
			New: func() any {
				samples := make([]F, size)
				return &samples
			},
		},
	}
}

func (p *SamplePool[F, T]) Get() *[]F {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	slice, ok := ptr.(*[]F)
	if !ok {
		return nil
	}

	return slice
}

func (p *SamplePool[F, T]) Put(samples *[]F) {
	p.pool.Put(&samples)
}
