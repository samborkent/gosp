package gosp

import (
	"bytes"
	"sync"
)

// BytesPool is a pool for retrieving variable sized byte buffers.
// Useful for processing encoded audio streams.
type BytesPool struct {
	pool sync.Pool
}

func NewBytesPool() *BytesPool {
	return &BytesPool{
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (p *BytesPool) Get() *bytes.Buffer {
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

func (p *BytesPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()
	p.pool.Put(buffer)
}
