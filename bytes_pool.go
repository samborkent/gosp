package gosp

import (
	"bytes"
	"sync"
)

// ByteBufferPool is a pool for retrieving variable sized byte buffers.
// Useful for processing encoded audio streams.
type ByteBufferPool struct {
	pool sync.Pool
}

func NewByteBufferPool() *ByteBufferPool {
	return &ByteBufferPool{
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (p *ByteBufferPool) Get() *bytes.Buffer {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	buf, ok := ptr.(*bytes.Buffer)
	if !ok {
		return nil
	}

	return buf
}

func (p *ByteBufferPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()
	p.pool.Put(buffer)
}
