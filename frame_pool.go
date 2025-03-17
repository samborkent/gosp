package gsp

import (
	"sync"
)

// FramePool is a pool for retrieving fixed size sample frame buffers.
// Useful for audio processing where the buffer size is fixed.
type FramePool[F Frame[T], T Type] struct {
	pool sync.Pool
}

func NewFramePool[F Frame[T], T Type](size int) *FramePool[F, T] {
	return &FramePool[F, T]{
		pool: sync.Pool{
			New: func() any {
				samples := make([]F, size)
				return &samples
			},
		},
	}
}

func (p *FramePool[F, T]) Get() *[]F {
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

func (p *FramePool[F, T]) Put(samples *[]F) {
	p.pool.Put(samples)
}
