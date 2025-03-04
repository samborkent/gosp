package gosp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sync/atomic"
	"unsafe"
)

// TODO: support little-endian and big-endian modes.
// TODO: separate float decoder to separate type.

var ErrLPCMDecoderNotinitialized = errors.New("gosp: LPCMDecoder: not initialized")

// LPCMDecoder is a linear pulse-code modulation decoder.
// This decoder can convert binary data into samples.
type LPCMDecoder[S SampleType[T], T Type] struct {
	r                         io.Reader
	channels, byteSize        int
	bytesRead, samplesDecoded atomic.Int64
	initialized               bool
}

func NewLPCMDecoder[S SampleType[T], T Type](r io.Reader) *LPCMDecoder[S, T] {
	return &LPCMDecoder[S, T]{
		r:           r,
		channels:    len(*new(S)),
		byteSize:    int(unsafe.Sizeof(*new(T))),
		initialized: true,
	}
}

// Decode reads from the internal [io.Reader] and decodes into the sample slice.
func (d *LPCMDecoder[S, T]) Decode(s []S) error {
	if !d.initialized {
		return ErrLPCMDecoderNotinitialized
	}

	buf := make([]byte, len(s)*d.channels*d.byteSize)

	bytesRead, err := d.r.Read(buf)
	if err != nil {
		return fmt.Errorf("gosp: LPCMDecoder.Decode: reading bytes into buffer: %w", err)
	}

	d.bytesRead.Add(int64(bytesRead))

	switch d.channels {
	case 1:
		// We only have a single channel, so we can safely perform this unsafe type casting.
		samplesDecoded := d.convertMono(unsafe.Slice((*Mono[T])(unsafe.Pointer(&s[0])), len(s)), buf)
		d.samplesDecoded.Add(int64(samplesDecoded))
	case 2:
		panic("gosp: LPCMDecoder.Decode: implement stereo decoding")
	default:
		panic("gosp: LPCMDecoder.Decode: implement multi-channel decoding")
	}

	return nil
}

// convertMono returns the number of samples converted.
func (d *LPCMDecoder[S, T]) convertMono(dst []Mono[T], src []byte) int {
	switch d.byteSize {
	case 1: // 8 bit
		minLen := min(len(dst), len(src))

		// Abuse overflow rules to deduce specific type.
		if *new(T)-1 > 0 {
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

		// Abuse overflow rules to deduce specific type.
		if *new(T)-1 > 0 {
			// uint16
			buf := unsafe.Slice((*uint16)(unsafe.Pointer(&dst[0])), minLen)
			for i := range minLen {
				buf[i] = binary.LittleEndian.Uint16(src[d.byteSize*i : d.byteSize*i+d.byteSize])
			}

			return minLen
		}

		// int16
		buf := unsafe.Slice((*int16)(unsafe.Pointer(&dst[0])), minLen)
		for i := range minLen {
			buf[i] = int16(binary.LittleEndian.Uint16(src[d.byteSize*i : d.byteSize*i+d.byteSize]))
		}

		return minLen
	case 4: // 32 bit
		minLen := min(len(dst), len(src)/d.byteSize)

		// Abuse overflow rules to deduce specific type.
		if *new(T)-1 > 0 {
			// uint32
			buf := unsafe.Slice((*uint32)(unsafe.Pointer(&dst[0])), minLen)
			for i := range minLen {
				buf[i] = binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*i+d.byteSize])
			}

			return minLen
		}

		// Abuse underflow rules to deduce specific type.
		maxInt32 := math.MaxInt32
		if T(maxInt32)+1 < 0 {
			// int32
			buf := unsafe.Slice((*int32)(unsafe.Pointer(&dst[0])), minLen)
			for i := range minLen {
				buf[i] = int32(binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*i+d.byteSize]))
			}

			return minLen
		}

		// float32
		buf := unsafe.Slice((*float32)(unsafe.Pointer(&dst[0])), minLen)
		for i := range minLen {
			buf[i] = math.Float32frombits(binary.LittleEndian.Uint32(src[d.byteSize*i : d.byteSize*i+d.byteSize]))
		}

		return minLen
	case 8: // 64 bit
		minLen := min(len(dst), len(src)/d.byteSize)

		// Abuse overflow rules to deduce specific type.
		if *new(T)-1 > 0 {
			// uint32
			buf := unsafe.Slice((*uint64)(unsafe.Pointer(&dst[0])), minLen)
			for i := range minLen {
				buf[i] = binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*i+d.byteSize])
			}

			return minLen
		}

		// Abuse underflow rules to deduce specific type.
		maxInt32 := math.MaxInt64
		if T(maxInt32)+1 < 0 {
			// int64
			buf := unsafe.Slice((*int64)(unsafe.Pointer(&dst[0])), minLen)
			for i := range minLen {
				buf[i] = int64(binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*i+d.byteSize]))
			}

			return minLen
		}

		// float64
		buf := unsafe.Slice((*float64)(unsafe.Pointer(&dst[0])), minLen)
		for i := range minLen {
			buf[i] = math.Float64frombits(binary.LittleEndian.Uint64(src[d.byteSize*i : d.byteSize*i+d.byteSize]))
		}

		return minLen
	default:
		panic("gosp: LPCMDecoder.convertMono: unknown bit size encountered")
	}
}
