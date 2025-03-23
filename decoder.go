package gsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sync/atomic"
	"unsafe"
)

const decoderMinRead = 128

var ErrDecoderNotinitialized = errors.New("gsp: Decoder: not initialized")

// Decoder is a linear PCM and floating-point decoder.
// This decoder can convert binary data into samples.
type Decoder[F Frame[T], T Type] struct {
	r                         io.Reader
	bytePool                  *ByteBufferPool
	samplePool                *Pool[F, T]
	bytesRead, samplesDecoded atomic.Int64
	channels, byteSize        int
	bigEndian                 bool
	initialized               bool
}

func NewDecoder[F Frame[T], T Type](r io.Reader, opts ...EncodingOption) *Decoder[F, T] {
	// Apply encoding options.
	var cfg EncodingConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var channels int

	switch frameType := any(*new(F)).(type) {
	case T:
		channels = 1
	case [2]T, Stereo[T]:
		channels = 2
	case []T:
		channels = len(frameType)
	case MultiChannel[T]:
		channels = len(frameType)
	default:
		panic("gsp: NewDecoder: unknown audio frame type")
	}

	return &Decoder[F, T]{
		r:           r,
		bytePool:    NewByteBufferPool(),
		samplePool:  NewPool[F, T](),
		channels:    channels,
		byteSize:    int(unsafe.Sizeof(T(0))),
		bigEndian:   cfg.BigEndian,
		initialized: true,
	}
}

func (d *Decoder[F, T]) Channels() int {
	return d.channels
}

func (d *Decoder[F, T]) ByteSize() int {
	return d.byteSize
}

// Decode reads from the internal [io.Reader] and decodes into the sample slice.
// It will read until [io.EOF] is reached, so users should ensure to pass an adequately sized sample frame slice to prevent data loss.
func (d *Decoder[F, T]) Decode(s []F) error {
	if !d.initialized {
		return ErrDecoderNotinitialized
	}

	// Retrieve byte buffer from pool.
	buf := d.bytePool.Get()
	if buf == nil {
		// Allocate in-case the pool does not have a pre-allocated entry ready.
		buf = bytes.NewBuffer(make([]byte, len(s)*d.channels*d.byteSize))
	}

	// Put buffer back in pool.
	defer d.bytePool.Put(buf)

	// TODO: check if there is a way to only read up to len(s)
	bytesRead, err := buf.ReadFrom(d.r)
	if err != nil {
		return fmt.Errorf("gsp: Decoder.Decode: reading bytes into buffer: %w", err)
	}

	d.bytesRead.Add(int64(bytesRead))

	switch d.channels {
	case 1:
		// There is only a single channel, so we can safely perform this unsafe type-casting.
		samplesDecoded := d.decodeMono(unsafe.Slice((*T)(unsafe.Pointer(&s[0])), len(s)), buf.Bytes())
		d.samplesDecoded.Add(int64(samplesDecoded))
	case 2:
		panic("gsp: Decoder.Decode: implement stereo decoding")
	default:
		panic("gsp: Decoder.Decode: implement multi-channel decoding")
	}

	return nil
}

// Decode reads from the internal [io.Reader] and decodes into the [Writer].
//
// TODO: can probably be improved, ideally should rely on same logic as Decode (see EncodeFrom).
func (d *Decoder[F, T]) DecodeTo(w Writer[F, T]) error {
	if !d.initialized {
		return ErrDecoderNotinitialized
	}

	// Retrieve byte buffer from pool.
	buf := d.bytePool.Get()
	if buf == nil {
		// Allocate in-case the pool does not have a pre-allocated entry ready.
		buf = bytes.NewBuffer(make([]byte, decoderMinRead*d.channels*d.byteSize))
	}

	// Put buffer back in pool.
	defer d.bytePool.Put(buf)

	bytesRead, err := buf.ReadFrom(d.r)
	if err != nil {
		return fmt.Errorf("gsp: Decoder.DecodeTo: reading bytes into buffer: %w", err)
	}

	d.bytesRead.Add(int64(bytesRead))

	sampleBuf := d.samplePool.Get()
	if sampleBuf == nil {
		// Allocate in-case the pool does not have a pre-allocated entry ready.
		sampleBuf = NewBuffer[F, T](make([]F, decoderMinRead))
	}

	defer d.samplePool.Put(sampleBuf)

	switch d.channels {
	case 1:
		// There is only a single channel, so we can safely perform this unsafe type-casting.
		_ = d.decodeMono(unsafe.Slice((*T)(unsafe.Pointer(&sampleBuf.Frames()[0])), sampleBuf.Len()), buf.Bytes())
	case 2:
		panic("gsp: Decoder.DecodeTo: implement stereo decoding")
	default:
		panic("gsp: Decoder.DecodeTo: implement multi-channel decoding")
	}

	samplesDecoded, err := w.Write(sampleBuf.Frames())
	if err != nil {
		return fmt.Errorf("gsp: Decoder.DecodeTo: writing frames: %w", err)
	}

	d.samplesDecoded.Add(int64(samplesDecoded))

	return nil
}

func (d *Decoder[F, T]) decodeMono(dst []T, src []byte) int {
	switch d.byteSize {
	case 1: // 8 bit
		minLen := min(len(dst), len(src))

		if isUnsigned[T]() {
			// uint8
			return copy(unsafe.Slice((*uint8)(unsafe.Pointer(&dst[0])), minLen), src)
		}

		// int8
		buf := unsafe.Slice((*int8)(unsafe.Pointer(&dst[0])), minLen)
		for i := range minLen {
			buf[i] = int8(src[i])
		}

		return minLen
	case 2: // 16 bit
		minLen := min(len(dst), len(src)/d.byteSize)

		if isUnsigned[T]() {
			// uint16
			buf := unsafe.Slice((*uint16)(unsafe.Pointer(&dst[0])), minLen)
			if d.bigEndian {
				for i := range minLen {
					buf[i] = binary.BigEndian.Uint16(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			} else {
				for i := range minLen {
					buf[i] = binary.LittleEndian.Uint16(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			}

			return minLen
		}

		// int16
		buf := unsafe.Slice((*int16)(unsafe.Pointer(&dst[0])), minLen)
		if d.bigEndian {
			for i := range minLen {
				buf[i] = int16(binary.BigEndian.Uint16(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		} else {
			for i := range minLen {
				buf[i] = int16(binary.LittleEndian.Uint16(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		}

		return minLen
	case 4: // 32 bit
		minLen := min(len(dst), len(src)/d.byteSize)

		if isUnsigned[T]() {
			// uint32
			buf := unsafe.Slice((*uint32)(unsafe.Pointer(&dst[0])), minLen)
			if d.bigEndian {
				for i := range minLen {
					buf[i] = binary.BigEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			} else {
				for i := range minLen {
					buf[i] = binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			}

			return minLen
		}

		if isSigned[T]() {
			// int32
			buf := unsafe.Slice((*int32)(unsafe.Pointer(&dst[0])), minLen)
			if d.bigEndian {
				for i := range minLen {
					buf[i] = int32(binary.BigEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)]))
				}
			} else {
				for i := range minLen {
					buf[i] = int32(binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)]))
				}
			}

			return minLen
		}

		// float32
		buf := unsafe.Slice((*float32)(unsafe.Pointer(&dst[0])), minLen)
		if d.bigEndian {
			for i := range minLen {
				buf[i] = math.Float32frombits(binary.BigEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		} else {
			for i := range minLen {
				buf[i] = math.Float32frombits(binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		}

		return minLen
	case 8: // 64 bit
		minLen := min(len(dst), len(src)/d.byteSize)

		if isUnsigned[T]() {
			// uint32
			buf := unsafe.Slice((*uint64)(unsafe.Pointer(&dst[0])), minLen)
			if d.bigEndian {
				for i := range minLen {
					buf[i] = binary.BigEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			} else {
				for i := range minLen {
					buf[i] = binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)])
				}
			}

			return minLen
		}

		if isSigned[T]() {
			// int64
			buf := unsafe.Slice((*int64)(unsafe.Pointer(&dst[0])), minLen)
			if d.bigEndian {
				for i := range minLen {
					buf[i] = int64(binary.BigEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)]))
				}
			} else {
				for i := range minLen {
					buf[i] = int64(binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)]))
				}
			}

			return minLen
		}

		// float64
		buf := unsafe.Slice((*float64)(unsafe.Pointer(&dst[0])), minLen)
		if d.bigEndian {
			for i := range minLen {
				buf[i] = math.Float64frombits(binary.BigEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		} else {
			for i := range minLen {
				buf[i] = math.Float64frombits(binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*(i+1)]))
			}
		}

		return minLen
	default:
		panic("gsp: Decoder.convertMono: unknown bit size encountered")
	}
}
