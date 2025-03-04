package gsp

import (
	"sync"
)

// BytePool is a pool for retrieving fixed size byte buffers.
// Useful for data stream processing where the buffer size is fixed.
type BytePool struct {
	pool sync.Pool
}

func NewBytePool(size int) *BytePool {
	return &BytePool{
		pool: sync.Pool{
			New: func() any {
				samples := make([]byte, size)
				return &samples
			},
		},
	}
}

func (p *BytePool) Get() *[]byte {
	ptr := p.pool.Get()
	if ptr == nil {
		return nil
	}

	slice, ok := ptr.(*[]byte)
	if !ok {
		return nil
	}

	return slice
}

func (p *BytePool) Put(b *[]byte) {
	p.pool.Put(b)
}
