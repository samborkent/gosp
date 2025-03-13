// Implementation from conner.dev/blog/circular-buffers
package gsp

import (
	"context"
)

// TODO: improve this generics magic
type buffered[F Frame[T], T Type] interface {
	Frame[T] | []F
}

type circularBuffer[B buffered[F, T], F Frame[T], T Type] struct {
	stream         <-chan B
	bufferedStream chan B
	overflowCount  uint64
}

func newCircularBuffer[B buffered[F, T], F Frame[T], T Type](stream <-chan B, bufferedStream chan B) *circularBuffer[B, F, T] {
	return &circularBuffer[B, F, T]{
		stream:         stream,
		bufferedStream: bufferedStream,
	}
}

func (b *circularBuffer[B, F, T]) OverflowCount() uint64 {
	return b.overflowCount
}

func (b *circularBuffer[B, F, T]) Run(ctx context.Context) {
	for sample := range b.stream {
		select {
		case <-ctx.Done():
			return
		case b.bufferedStream <- sample:
		default:
			// Buffer is full, drop oldest value
			<-b.bufferedStream
			b.overflowCount++
			b.bufferedStream <- sample
		}
	}
}
