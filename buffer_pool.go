package gsp

import (
	"sync"
)

// BytesPool is a pool for retrieving variable sized sample buffers.
// Useful for processing decoded audio streams.
type BufferPool[F Frame[T], T Type] struct {
	pool sync.Pool
}

func NewBufferPool[F Frame[T], T Type]() *BufferPool[F, T] {
	return &BufferPool[F, T]{
		pool: sync.Pool{
			New: func() any {
				return new(Buffer[F, T])
			},
		},
	}
}

func (p *BufferPool[F, T]) Get() *Buffer[F, T] {
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

func (p *BufferPool[F, T]) Put(buffer *Buffer[F, T]) {
	buffer.Reset()
	p.pool.Put(buffer)
}
