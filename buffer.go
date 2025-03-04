// Adapted from bytes.Buffer

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gosp

// Simple sample buffer for marshaling data.

import (
	"errors"
	"fmt"
	"io"
)

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// A Buffer is a variable-sized buffer of samples with [Buffer.Read] and [Buffer.Write] methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer[S SampleType[T], T Type] struct {
	buf      []S    // contents are the samples buf[off : len(buf)]
	off      int    // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadSample can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var (
	ErrTooLarge     = errors.New("gosp: Buffer: too large")
	errNegativeRead = errors.New("gosp: Buffer: reader returned negative count from Read")
)

const maxInt = int(^uint(0) >> 1)

// Samples returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like [Buffer.Read], [Buffer.Write], [Buffer.Reset], or [Buffer.Truncate]).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer[S, T]) Samples() []S { return b.buf[b.off:] }

// AvailableBuffer returns an empty buffer with b.Available() capacity.
// This buffer is intended to be appended to and
// passed to an immediately succeeding [Buffer.Write] call.
// The buffer is only valid until the next write operation on b.
func (b *Buffer[S, T]) AvailableBuffer() []S { return b.buf[len(b.buf):] }

// empty reports whether the unread portion of the buffer is empty.
func (b *Buffer[S, T]) empty() bool { return len(b.buf) <= b.off }

// Len returns the number of samples of the unread portion of the buffer;
// b.Len() == len(b.Samples()).
func (b *Buffer[S, T]) Len() int { return len(b.buf) - b.off }

// Cap returns the capacity of the buffer's underlying sample slice, that is, the
// total space allocated for the buffer's data.
func (b *Buffer[S, T]) Cap() int { return cap(b.buf) }

// Available returns how many samples are unused in the buffer.
func (b *Buffer[S, T]) Available() int { return cap(b.buf) - len(b.buf) }

// Truncate discards all but the first n unread samples from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (b *Buffer[S, T]) Truncate(n int) {
	if n == 0 {
		b.Reset()
		return
	}
	b.lastRead = opInvalid
	if n < 0 || n > b.Len() {
		panic("gosp: Buffer.Truncate: truncation out of range")
	}
	b.buf = b.buf[:b.off+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as [Buffer.Truncate](0).
func (b *Buffer[S, T]) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
	b.lastRead = opInvalid
}

// tryGrowByReslice is an inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where samples should be written and whether it succeeded.
func (b *Buffer[S, T]) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more samples.
// It returns the index where samples should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer[S, T]) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]S, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Add b.off to account for b.buf[:b.off] being sliced off the front.
		b.buf = growSlice[S, T](b.buf[b.off:], b.off+n)
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n samples. After Grow(n), at least n samples can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with [ErrTooLarge].
func (b *Buffer[S, T]) Grow(n int) {
	if n < 0 {
		panic("gosp: Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with [ErrTooLarge].
func (b *Buffer[S, T]) Write(p []S) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

// MinRead is the minimum slice size passed to a [Buffer.Read] call by
// [Buffer.ReadFrom]. As long as the [Buffer] has at least MinRead samples beyond
// what is required to hold the contents of r, [Buffer.ReadFrom] will not grow the
// underlying buffer.
const MinRead = 512

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of samples read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
func (b *Buffer[S, T]) ReadFrom(r Reader[S, T]) (n int64, err error) {
	b.lastRead = opInvalid
	for {
		i := b.grow(MinRead)
		b.buf = b.buf[:i]
		m, e := r.Read(b.buf[i:cap(b.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		b.buf = b.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice[S SampleType[T], T Type](b []S, n int) []S {
	defer func() {
		if recover() != nil {
			fmt.Printf("%s\n", ErrTooLarge.Error())
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]S, n)...)
	// This avoids unnecessary zero-ing of the first len(b) samples of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	// The growth rate has historically always been 2x. In the future,
	// we could rely purely on append to determine the growth rate.
	c = max(c, 2*cap(b))
	b2 := append([]S(nil), make([]S, c)...)
	i := copy(b2, b)
	return b2[:i]
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of samples written; it always fits into an
// int, but it is int64 to match the [io.WriterTo] interface. Any error
// encountered during the write is also returned.
func (b *Buffer[S, T]) WriteTo(w Writer[S, T]) (n int64, err error) {
	b.lastRead = opInvalid
	if nSamples := b.Len(); nSamples > 0 {
		m, e := w.Write(b.buf[b.off:])
		if m > nSamples {
			panic("gosp: Buffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all samples should have been written, by definition of
		// Write method in io.Writer
		if m != nSamples {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	b.Reset()
	return n, nil
}

// WriteSample appends the sample s to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match [bufio.Writer]'s
// WriteSample. If the buffer becomes too large, WriteSample will panic with
// [ErrTooLarge].
func (b *Buffer[S, T]) WriteSample(s S) error {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(1)
	if !ok {
		m = b.grow(1)
	}
	b.buf[m] = s
	return nil
}

// Read reads the next len(p) samples from the buffer or until the buffer
// is drained. The return value n is the number of samples read. If the
// buffer has no data to return, err is [io.EOF] (unless len(p) is zero);
// otherwise it is nil.
func (b *Buffer[S, T]) Read(p []S) (n int, err error) {
	b.lastRead = opInvalid
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return n, nil
}

// Next returns a slice containing the next n samples from the buffer,
// advancing the buffer as if the samples had been returned by [Buffer.Read].
// If there are fewer than n samples in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer[S, T]) Next(n int) []S {
	b.lastRead = opInvalid
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}

// ReadSample reads and returns the next sample from the buffer.
// If no sample is available, it returns error [io.EOF].
func (b *Buffer[S, T]) ReadSample() (S, error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		return *new(S), io.EOF
	}
	c := b.buf[b.off]
	b.off++
	b.lastRead = opRead
	return c, nil
}

var errUnreadSample = errors.New("gosp: Buffer.UnreadSample: previous operation was not a successful read")

// UnreadSample unreads the last sample returned by the most recent successful
// read operation that read at least one sample. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// samples, UnreadSample returns an error.
func (b *Buffer[S, T]) UnreadSample() error {
	if b.lastRead == opInvalid {
		return errUnreadSample
	}
	b.lastRead = opInvalid
	if b.off > 0 {
		b.off--
	}
	return nil
}

// NewBuffer creates and initializes a new [Buffer] using buf as its
// initial contents. The new [Buffer] takes ownership of buf, and the
// caller should not use buf after this call. NewBuffer is intended to
// prepare a [Buffer] to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new([Buffer]) (or just declaring a [Buffer] variable) is
// sufficient to initialize a [Buffer].
func NewBuffer[S SampleType[T], T Type](buf []S) *Buffer[S, T] { return &Buffer[S, T]{buf: buf} }
