package gsp

import (
	"sync"
)

// Pool is a pool for retrieving variable sized sample buffers.
// Useful for processing decoded audio streams.
type Pool[F Frame[T], T Type] struct {
	pool sync.Pool
}

func NewPool[F Frame[T], T Type]() *Pool[F, T] {
	return &Pool[F, T]{
		pool: sync.Pool{
			New: func() any {
				return new(Buffer[F, T])
			},
		},
	}
}

func (p *Pool[F, T]) Get() *Buffer[F, T] {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	buf, ok := ptr.(*Buffer[F, T])
	if !ok {
		return nil
	}

	return buf
}

func (p *Pool[F, T]) Put(buffer *Buffer[F, T]) {
	buffer.Reset()
	p.pool.Put(buffer)
}
