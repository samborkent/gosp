package gosp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"unsafe"
)

var ErrEncoderNotinitialized = errors.New("gosp: Encoder: not initialized")

// Encoder is a linear PCM and floating-point encoder.
// This encoder can convert samples into their binary representation.
type Encoder[F Frame[T], T Type] struct {
	w                            io.Writer
	pool                         *ByteBufferPool
	samplesEncoded, bytesWritten atomic.Int64
	channels, byteSize           int
	bigEndian                    bool
	initialized                  bool
}

func NewEncoder[F Frame[T], T Type](w io.Writer, opts ...EncodingOption) *Encoder[F, T] {
	// Apply encoding options.
	var cfg EncodingConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Encoder[F, T]{
		w:           w,
		pool:        NewByteBufferPool(),
		channels:    len(*new(F)),
		byteSize:    int(unsafe.Sizeof(*new(T))),
		bigEndian:   cfg.BigEndian,
		initialized: true,
	}
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *Encoder[F, T]) Encode(s []F) error {
	if !e.initialized {
		return ErrDecoderNotinitialized
	}

	switch e.channels {
	case 1:
		// There is only a single channel, so we can safely perform this unsafe type-casting.
		samplesEncoded, err := e.convertMono(unsafe.Slice((*Mono[T])(unsafe.Pointer(&s[0])), len(s)))
		if err != nil {
			return fmt.Errorf("gosp: Encoder.Encode: encoding sample frames: %w", err)
		}

		e.samplesEncoded.Add(int64(samplesEncoded))
	case 2:
		panic("gosp: Encoder.Encode: stereo encoding not implemented")
	default:
		panic("gosp: Encoder.Encode: multi-channel encoding not implemented")
	}

	return nil
}

// convertMono writes the encoded sampled to the internal [io.Writer].
func (e *Encoder[F, T]) convertMono(src []Mono[T]) (int64, error) {
	buf := e.pool.Get()
	if buf == nil {
		// Allocate in-case the pool does not have a pre-allocated entry ready.
		buf = bytes.NewBuffer(make([]byte, len(src)*e.channels*e.byteSize))
	}

	bufio.NewWriterSize(e.w, len(src)*e.channels*e.byteSize)

	switch e.byteSize {
	case 1: // 8 bit
		minLen := len(src)

		// Abuse overflow rules to deduce specific type.
		if T(0)-1 > 0 {
			// uint8
			n, err := e.w.Write(unsafe.Slice((*uint8)(unsafe.Pointer(&src[0])), minLen))
			e.pool.Put(buf)
			return int64(n), err
		}

		// int8
		for i := range minLen {
			_ = buf.WriteByte(byte(src[i][0]))
		}

		n, err := buf.WriteTo(e.w)
		e.pool.Put(buf)
		return n, err
	// case 2: // 16 bit
	// 	minLen := min(len(dst), len(src)/d.byteSize)

	// 	// uint16 & int16
	// 	for i := range minLen {
	// 		if d.bigEndian {
	// 			binary.BigEndian.PutUint16(dst[d.byteSize*i:d.byteSize*(i+1)], uint16(src[i][0]))
	// 		} else {
	// 			binary.LittleEndian.PutUint16(dst[d.byteSize*i:d.byteSize*(i+1)], uint16(src[i][0]))
	// 		}
	// 	}

	// 	return minLen
	// case 4: // 32 bit
	// 	minLen := min(len(dst), len(src)/d.byteSize)

	// 	// Abuse overflow rules to deduce specific type.
	// 	if T(0)-1 > 0 || T(maxInt32)+1 < 0 {
	// 		// uint32 & int32
	// 		for i := range minLen {
	// 			if d.bigEndian {
	// 				binary.BigEndian.PutUint32(dst[d.byteSize*i:d.byteSize*(i+1)], uint32(src[i][0]))
	// 			} else {
	// 				binary.LittleEndian.PutUint32(dst[d.byteSize*i:d.byteSize*(i+1)], uint32(src[i][0]))
	// 			}
	// 		}

	// 		return minLen
	// 	}

	// 	// float32
	// 	for i := range minLen {
	// 		if d.bigEndian {
	// 			binary.BigEndian.PutUint32(dst[d.byteSize*i:d.byteSize*(i+1)], math.Float32bits(float32(src[i][0])))
	// 		} else {
	// 			binary.LittleEndian.PutUint32(dst[d.byteSize*i:d.byteSize*(i+1)], math.Float32bits(float32(src[i][0])))
	// 		}
	// 	}

	// 	return minLen
	// case 8: // 64 bit
	// 	minLen := min(len(dst), len(src)/d.byteSize)

	// 	// Abuse overflow rules to deduce specific type.
	// 	if T(0)-1 > 0 || T(maxInt64)+1 < 0 {
	// 		// uint64 & int64
	// 		for i := range minLen {
	// 			if d.bigEndian {
	// 				binary.BigEndian.PutUint64(dst[d.byteSize*i:d.byteSize*(i+1)], uint64(src[i][0]))
	// 			} else {
	// 				binary.LittleEndian.PutUint64(dst[d.byteSize*i:d.byteSize*(i+1)], uint64(src[i][0]))
	// 			}
	// 		}

	// 		return minLen
	// 	}

	// 	// float64
	// 	for i := range minLen {
	// 		if d.bigEndian {
	// 			binary.BigEndian.PutUint64(dst[d.byteSize*i:d.byteSize*(i+1)], math.Float64bits(float64(src[i][0])))
	// 		} else {
	// 			binary.LittleEndian.PutUint64(dst[d.byteSize*i:d.byteSize*(i+1)], math.Float64bits(float64(src[i][0])))
	// 		}
	// 	}

	// 	return minLen
	default:
		panic("gosp: Encoder.convertMono: unknown bit size encountered")
	}
}
