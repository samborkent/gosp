package gosp

import (
	"sync"
)

// BytesPool is a pool for retrieving variable sized sample buffers.
// Useful for processing decoded audio streams.
type BufferPool[S SampleType[T], T Type] struct {
	pool sync.Pool
}

func NewBufferPool[S SampleType[T], T Type]() *BufferPool[S, T] {
	return &BufferPool[S, T]{
		pool: sync.Pool{
			New: func() any {
				return new(Buffer[S, T])
			},
		},
	}
}

func (p *BufferPool[S, T]) Get() *Buffer[S, T] {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	buf, ok := ptr.(*Buffer[S, T])
	if !ok {
		return nil
	}

	return buf
}

func (p *BufferPool[S, T]) Put(buffer *Buffer[S, T]) {
	buffer.Reset()
	p.pool.Put(buffer)
}
