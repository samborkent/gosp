package gsp

import (
	"context"
	"math"
	"runtime"
	"unsafe"
)

var (
	_ Reader[Mono[uint8], uint8] = &Converter[Mono[uint8], Mono[uint8], uint8, uint8]{}
	_ Writer[Mono[uint8], uint8] = &Converter[Mono[uint8], Mono[uint8], uint8, uint8]{}
)

// Converter can converts from one data type to another.
type Converter[In Frame[I], Out Frame[O], I Type, O Type] struct {
	input  chan In
	output chan Out
}

func NewConverter[In Frame[I], Out Frame[O], I Type, O Type](bufferSize int) *Converter[In, Out, I, O] {
	input := make(chan In, bufferSize)
	output := make(chan Out, bufferSize)

	converter := &Converter[In, Out, I, O]{
		input:  input,
		output: output,
	}

	ctx, cancel := context.WithCancel(context.Background())

	runtime.AddCleanup(converter, func(_ any) {
		cancel()
		close(input)
		close(output)
	}, nil)

	go converter.run(ctx)

	return converter
}

// Read converted output from converter, blocking until output is filled.
// Error is always nil, so can be safely ignored (present to match to [Reader]).
func (c *Converter[In, Out, I, O]) Read(output []Out) (int, error) {
	for i := range output {
		output[i] = <-c.output
	}

	return len(output), nil
}

// Get tries to get n samples from the output. If no more frames remain, return intermediate result.
func (c *Converter[In, Out, I, O]) Get(n int) []Out {
	output := make([]Out, 0, n)

	select {
	case frame := <-c.output:
		output = append(output, frame)
	default:
		return output
	}

	return output
}

// ReadFrame reads a single converted frame from the converter, blocking if no new frame is available.
func (c *Converter[In, Out, I, O]) ReadFrame() Out {
	return <-c.output
}

// GetFrame is a non-blocking version of [Converter.ReadFrame].
func (c *Converter[In, Out, I, O]) GetFrame() (Out, bool) {
	select {
	case frame := <-c.output:
		return frame, true
	default:
		return *new(Out), false
	}
}

// Write input to converter, blocking if the input channel is full.
// Error is always nil, so can be safely ignored (present to match to [Writer]).
func (c *Converter[In, Out, I, O]) Write(input []In) (int, error) {
	for _, sample := range input {
		c.input <- sample
	}

	return len(input), nil
}

// Put will try to write the input to the converter. If the converter input is full, return early and return number of entries written.
func (c *Converter[In, Out, I, O]) Put(input []In) int {
	n := 0

frameLoop:
	for _, frame := range input {
		select {
		case c.input <- frame:
			n++
		default:
			break frameLoop
		}
	}

	return n
}

// WriteFrame writes a single frame to the converter, blocking if the input channel is full.
func (c *Converter[In, Out, I, O]) WriteFrame(frame In) {
	c.input <- frame
}

// PutFrame puts a single frame to the converter. It returns true is the frame was written to the converter, false if the input is full.
func (c *Converter[In, Out, I, O]) PutFrame(frame In) bool {
	select {
	case c.input <- frame:
		return true
	default:
		return false
	}
}

func (c *Converter[In, Out, I, O]) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-c.input:
			c.output <- convert[In, Out, I, O](input)
		}
	}
}

func convert[In Frame[I], Out Frame[O], I Type, O Type](input In) (output Out) {
	switch len(input) {
	case 1: // mono input
		switch len(output) {
		case 1: // mono output
			out := convertMono[O](*(*Mono[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case 2: // stereo output
			out := convertMonoToStereo[O](*(*Mono[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default: // multi-channel output
			out := convertMonoToMultiChannel[O](*(*Mono[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		}
	case 2: // stereo input
		switch len(output) {
		case 1: // mono output
			out := convertStereoToMono[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case 2: // stereo output
			out := convertStereo[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default: // multi-channel output
			out := convertStereoToMultiChannel[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		}
	default: // multi-channel input
		switch len(output) {
		case 1: // mono output
			out := convertMultiChannelToMono[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case 2: // stereo output
			out := convertMultiChannelToStereo[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default: // multi-channel output
			out := convertMultiChannel[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		}
	}
}

func convertMono[O Type, I Type](input Mono[I]) (output Mono[O]) {
	return ToMono(convertType[O](input.M()))
}

func convertMonoToStereo[O Type, I Type](input Mono[I]) (output Stereo[O]) {
	return MonoToStereo(ToMono(convertType[O](input.M())))
}

func convertMonoToMultiChannel[O Type, I Type](input Mono[I]) (output MultiChannel[O]) {
	if len(output) == 0 {
		return MultiChannel[O]{}
	}

	out := ZeroMultiChannel[O](len(output))
	out[0] = convertType[O](input.M())

	return out
}

func convertStereoToMono[O Type, I Type](input Stereo[I]) (output Mono[O]) {
	return ToMono(convertType[O](input.M()))
}

func convertStereo[O Type, I Type](input Stereo[I]) (output Stereo[O]) {
	return ToStereo(convertType[O](input.L()), convertType[O](input.R()))
}

func convertStereoToMultiChannel[O Type, I Type](input Stereo[I]) (output MultiChannel[O]) {
	if len(output) == 0 {
		return MultiChannel[O]{}
	}

	out := ZeroMultiChannel[O](len(output))

	switch len(output) {
	case 1:
		out[0] = convertType[O](input.M())
	default:
		out[L] = convertType[O](input.L())
		out[R] = convertType[O](input.R())
	}

	return out
}

func convertMultiChannelToMono[O Type, I Type](input MultiChannel[I]) (output Mono[O]) {
	if len(input) == 0 {
		return ZeroMono[O]()
	}

	return ToMono(convertType[O](input.M()))
}

func convertMultiChannelToStereo[O Type, I Type](input MultiChannel[I]) (output Stereo[O]) {
	if len(input) == 0 {
		return ZeroStereo[O]()
	}

	switch len(input) {
	case 1:
		return MonoToStereo(ToMono(convertType[O](input[0])))
	case 2:
		return ToStereo(convertType[O](input[L]), convertType[O](input[R]))
	default:
		out := ToStereo(convertType[O](input[L]), convertType[O](input[R]))

		sides := O(0)
		for i := 2; i < len(input); i++ {
			sides += convertType[O](input[i])
		}
		sides /= O(len(input) - 2)
		sides /= 2

		out[L] += sides
		out[R] += sides

		return out
	}
}

func convertMultiChannel[O Type, I Type](input MultiChannel[I]) (output MultiChannel[O]) {
	if len(output) == 0 {
		return MultiChannel[O]{}
	}

	switch len(input) {
	case 0:
		return ZeroMultiChannel[O](len(output))
	case 1:
		return convertMonoToMultiChannel[O](ToMono(input[0]))
	case 2:
		return convertStereoToMultiChannel[O](ToStereo(input[L], input[R]))
	default:
		switch len(output) {
		case 1:
			return MultiChannel[O]{convertMultiChannelToMono[O](input).M()}
		case 2:
			out := convertMultiChannelToStereo[O](input)
			return MultiChannel[O]{out[L], out[R]}
		default:
			out := ZeroMultiChannel[O](len(output))

			for i := range min(len(input), len(output)) {
				out[i] = convertType[O](input[i])
			}

			return out
		}
	}
}

// TODO: add cases for float input
func convertType[O Type, I Type](in I) (out O) {
	switch unsafe.Sizeof(in) {
	case 1:
		switch unsafe.Sizeof(out) {
		case 1:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint8 -> uint8
					return O(in)
				}

				// uint8 -> int8
				return O(int8(int16(in) - int16(zeroUint8)))
			}

			if isUnsigned[O]() {
				// int8 -> uint8
				return O(uint8(int16(in) + int16(zeroUint8)))
			}

			// int8 -> int8
			return O(in)
		case 2:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint8 -> uint16
					return O(uint16(in) << 8)
				}

				// uint8 -> int16
				return O((int16(in) - int16(zeroUint8)) << 8)
			}

			if isUnsigned[O]() {
				// int8 -> uint16
				return O((int16(in) + int16(zeroUint8)) << 8)
			}

			// int8 -> int16
			return O(int16(in) << 8)
		case 4:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint8 -> uint32
					return O(uint32(in) << 24)
				}

				if isSigned[O]() {
					// uint8 -> int32
					return O(int32(int16(in)-int16(zeroUint8)) << 24)
				}

				// uint8 -> float32
				return O(float32(int16(in)-int16(zeroUint8)) / (-math.MinInt8))
			}

			if isUnsigned[O]() {
				// int8 -> uint32
				return O(uint32(int16(in)+int16(zeroUint8)) << 24)
			}

			if isSigned[O]() {
				// int8 -> int32
				return O(int32(in) << 24)
			}

			// int8 -> float32
			return O(float32(in) / (-math.MinInt8))
		case 8:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint8 -> uint64
					return O(uint64(in) << 56)
				}

				if isSigned[O]() {
					// uint8 -> int64
					return O(int64(int16(in)-int16(zeroUint8)) << 56)
				}

				// uint8 -> float64
				return O(float64(int16(in)-int16(zeroUint8)) / (-math.MinInt8))
			}

			if isUnsigned[O]() {
				// int8 -> uint64
				return O(uint64(int16(in)+int16(zeroUint8)) << 56)
			}

			if isSigned[O]() {
				// int8 -> int64
				return O(int64(in) << 56)
			}

			// int8 -> float64
			return O(float64(in) / (-math.MinInt8))
		default:
			panic("gsp: convertType: unknown output bit size encountered")
		}
	case 2:
		switch unsafe.Sizeof(out) {
		case 1:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint16 -> uint8
					return O(uint8(uint16(in) >> 8))
				}

				// uint16 -> int8
				return O(int8(int16(int32(in)-int32(zeroUint16)) >> 8))
			}

			if isUnsigned[O]() {
				// int16 -> uint8
				return O(uint8(uint16(int32(in)+int32(zeroUint16)) >> 8))
			}

			// int16 -> int8
			return O(int8(int16(in) >> 8))
		case 2:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint16 -> uint16
					return O(in)
				}

				// uint16 -> int16
				return O(int16(int32(in) - int32(zeroUint16)))
			}

			if isUnsigned[O]() {
				// int16 -> uint16
				return O(uint16(int32(in) + int32(zeroUint16)))
			}

			// int16 -> int16
			return O(in)
		case 4:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint16 -> uint32
					return O(uint32(in) << 16)
				}

				if isSigned[O]() {
					// uint16 -> int32
					return O(int32(int32(in)-int32(zeroUint16)) << 16)
				}

				// uint16 -> float32
				return O(float32(int32(in)-int32(zeroUint16)) / (-math.MinInt16))
			}

			if isUnsigned[O]() {
				// int16 -> uint32
				return O(uint32(int32(in)+int32(zeroUint16)) << 16)
			}

			if isSigned[O]() {
				// int16 -> int32
				return O(int32(in) << 16)
			}

			// int16 -> float32
			return O(float32(in) / (-math.MinInt16))
		case 8:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint16 -> uint64
					return O(uint64(in) << 48)
				}

				if isSigned[O]() {
					// uint16 -> int64
					return O(int64(int32(in)-int32(zeroUint16)) << 48)
				}

				// uint16 -> float64
				return O(float64(int32(in)-int32(zeroUint16)) / (-math.MinInt16))
			}

			if isUnsigned[O]() {
				// int16 -> uint64
				return O(uint64(int32(in)+int32(zeroUint16)) << 48)
			}

			if isSigned[O]() {
				// int16 -> int64
				return O(int64(in) << 48)
			}

			// int16 -> float64
			return O(float64(in) / (-math.MinInt16))
		default:
			panic("gsp: convertType: unknown output bit size encountered")
		}
	case 4:
		switch unsafe.Sizeof(out) {
		case 1:
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint8
					return O(uint8(uint32(in) >> 24))
				}

				// uint32 -> int8
				return O(int8(int32(int64(in)-int64(zeroUint32)) >> 24))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint8
					return O(uint8(uint32(int64(in)+int64(zeroUint32)) >> 24))
				}

				// int32 -> int8
				return O(int8(int32(in) >> 24))
			default:
				if isUnsigned[O]() {
					// float32 -> uint8
					return O(uint8(int16(float32(in)*(-math.MinInt8)) + int16(zeroUint8)))
				}

				// float32 -> int8
				return O(int8(float32(in) * (-math.MinInt8)))
			}
		case 2:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint32 -> uint16
					return O(uint16(uint32(in) >> 16))
				}

				// uint32 -> int16
				return O(int16(int32(int64(in)-int64(zeroUint32)) >> 16))
			}

			if isSigned[I]() {
				if isUnsigned[O]() {
					// int32 -> uint16
					return O(uint16(uint32(int64(in)+int64(zeroUint32)) >> 16))
				}

				// int32 -> int16
				return O(int16(int32(in) >> 16))
			}

			if isUnsigned[O]() {
				// float32 -> uint16
				return O(uint8(int16(float32(in)*(-math.MinInt8)) + int16(zeroUint8)))
			}

			// float32 -> int16
			return O(int16(float32(in) * (-math.MinInt16)))
		case 4:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint32 -> uint32
					return O(in)
				}

				if isSigned[O]() {
					// uint32 -> int32
					return O(int32(int64(in) - int64(zeroUint32)))
				}

				// uint32 -> float32
				return O(float32(float64(int64(in)-int64(zeroUint32)) / (-math.MinInt32)))
			}

			if isUnsigned[O]() {
				// int32 -> uint32
				return O(uint32(int64(in) + int64(zeroUint32)))
			}

			if isSigned[O]() {
				// int32 -> int32
				return O(in)
			}

			// int32 -> float32
			return O(float32(float64(in) / (-math.MinInt32)))
		case 8:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint32 -> uint64
					return O(uint64(in) << 32)
				}

				if isSigned[O]() {
					// uint32 -> int64
					return O((int64(in) - int64(zeroUint32)) << 32)
				}

				// uint32 -> float64
				return O(float64(int64(in)-int64(zeroUint32)) / (-math.MinInt32))
			}

			if isUnsigned[O]() {
				// int32 -> uint64
				return O(uint64(int64(in)+int64(zeroUint32)) << 32)
			}

			if isSigned[O]() {
				// int32 -> int64
				return O(int64(in) << 32)
			}

			// int32 -> float64
			return O(float64(in) / (-math.MinInt32))
		default:
			panic("gsp: convertType: unknown output bit size encountered")
		}
	case 8:
		switch unsafe.Sizeof(out) {
		case 1:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint64 -> uint8
					return O(uint8(uint64(in) >> 56))
				}

				// uint64 -> int8
				return O(int8(int64(float64(in)-float64(zeroUint64)) >> 56))
			}

			if isUnsigned[O]() {
				// int64 -> uint8
				return O(uint8(uint64(float64(in)+float64(zeroUint64)) >> 56))
			}

			if isSigned[I]() {
				// int64 -> int8
				return O(int8(int64(in) >> 56))
			}

			// float64 -> int8
			return O(int8(float64(in) * (-math.MinInt8)))
		case 2:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint64 -> uint16
					return O(uint16(uint64(in) >> 48))
				}

				// uint64 -> int16
				return O(int16(int64(float64(in)-float64(zeroUint64)) >> 48))
			}

			if isUnsigned[O]() {
				// int64 -> uint16
				return O(uint16(uint64(float64(in)+float64(zeroUint64)) >> 48))
			}

			if isSigned[I]() {
				// int64 -> int16
				return O(int16(int64(in) >> 48))
			}

			// float64 -> int16
			return O(int16(float64(in) * (-math.MinInt16)))
		case 4:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint64 -> uint32
					return O(uint32(uint64(in) >> 32))
				}

				if isSigned[O]() {
					// uint64 -> int32
					return O(int32(int64(float64(in)-float64(zeroUint64)) >> 32))
				}

				// uint64 -> float32
				return O(float32((float64(in) - float64(zeroUint64)) / (-math.MinInt64)))
			}

			if isUnsigned[O]() {
				// int64 -> uint32
				return O(uint32(uint64(float64(in)+float64(zeroUint64)) >> 32))
			}

			if isSigned[O]() {
				// int64 -> int32
				return O(int32(int64(in) >> 32))
			}

			// int64 -> float32
			return O(float32(float64(in) / (-math.MinInt64)))
		case 8:
			if isUnsigned[I]() {
				if isUnsigned[O]() {
					// uint64 -> uint64
					return O(in)
				}

				if isSigned[O]() {
					// uint64 -> int64
					return O(int64(float64(in) - float64(zeroUint64)))
				}

				// uint64 -> float64
				return O((float64(in) - float64(zeroUint64)) / (-math.MinInt64))
			}

			if isUnsigned[O]() {
				// int64 -> uint64
				return O(uint64(float64(in) + float64(zeroUint64)))
			}

			if isSigned[O]() {
				// int64 -> int64
				return O(in)
			}

			// int64 -> float64
			return O(float64(in) / (-math.MinInt64))
		default:
			panic("gsp: convertType: unknown output bit size encountered")
		}
	default:
		panic("gsp: convertType: unknown input bit size encountered")
	}
}
