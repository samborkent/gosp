package gsp

import (
	"context"
	"runtime"
	"unsafe"
)

var (
	_ Reader[uint8, uint8] = &Converter[uint8, uint8, uint8, uint8]{}
	_ Writer[uint8, uint8] = &Converter[uint8, uint8, uint8, uint8]{}
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

	runtime.AddCleanup(converter, func(_ int) {
		cancel()
		close(input)
		close(output)
	}, 0)

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
	switch any(*new(In)).(type) {
	case I: // mono input
		switch any(*new(Out)).(type) {
		case O: // mono -> mono
			out := convertMono[O](*(*I)(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case [2]O, Stereo[O]: // mono -> stereo
			out := convertMonoToStereo[O](*(*I)(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case []O, MultiChannel[O]: // mono -> multi-channel
			out := convertMonoToMultiChannel[O](*(*I)(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default:
			panic("gsp: convert: unknown output audio type")
		}
	case [2]I, Stereo[I]: // stereo input
		switch any(*new(Out)).(type) {
		case O: // stereo -> mono
			out := convertStereoToMono[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case [2]O, Stereo[O]: // stereo -> stereo
			out := convertStereo[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case []O, MultiChannel[O]: // stereo -> multi-channel
			out := convertStereoToMultiChannel[O](*(*Stereo[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default:
			panic("gsp: convert: unknown output audio type")
		}
	case []I, MultiChannel[I]: // multi-channel input
		switch any(*new(Out)).(type) {
		case O: // multi-channel -> mono
			out := convertMultiChannelToMono[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case [2]O, Stereo[O]: // multi-channel -> stereo
			out := convertMultiChannelToStereo[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		case []O, MultiChannel[O]: // multi-channel -> multi-channel
			out := convertMultiChannel[O](*(*MultiChannel[I])(unsafe.Pointer(&input)))
			return *(*Out)(unsafe.Pointer(&out))
		default:
			panic("gsp: convert: unknown output audio type")
		}
	default:
		panic("gsp: convert: unknown input audio type")
	}
}

func convertMono[O Type, I Type](input I) (output O) {
	return ConvertType[O](input)
}

func convertMonoToStereo[O Type, I Type](input I) (output Stereo[O]) {
	return MonoToStereo(ConvertType[O](input))
}

func convertMonoToMultiChannel[O Type, I Type](input I) (output MultiChannel[O]) {
	if len(output) == 0 {
		return MultiChannel[O]{}
	}

	out := ZeroMultiChannel[O](len(output))
	out[0] = ConvertType[O](input)

	return out
}

func convertStereoToMono[O Type, I Type](input Stereo[I]) (output O) {
	return ConvertType[O](input.M())
}

func convertStereo[O Type, I Type](input Stereo[I]) (output Stereo[O]) {
	return ToStereo(ConvertType[O](input.L()), ConvertType[O](input.R()))
}

func convertStereoToMultiChannel[O Type, I Type](input Stereo[I]) (output MultiChannel[O]) {
	if len(output) == 0 {
		return MultiChannel[O]{}
	}

	out := ZeroMultiChannel[O](len(output))

	switch len(output) {
	case 1:
		out[0] = ConvertType[O](input.M())
	default:
		out[L] = ConvertType[O](input.L())
		out[R] = ConvertType[O](input.R())
	}

	return out
}

func convertMultiChannelToMono[O Type, I Type](input MultiChannel[I]) (output O) {
	if len(input) == 0 {
		return Zero[O]()
	}

	return ConvertType[O](input.M())
}

func convertMultiChannelToStereo[O Type, I Type](input MultiChannel[I]) (output Stereo[O]) {
	if len(input) == 0 {
		return ZeroStereo[O]()
	}

	switch len(input) {
	case 1:
		return MonoToStereo(ConvertType[O](input[0]))
	case 2:
		return ToStereo(ConvertType[O](input[L]), ConvertType[O](input[R]))
	default:
		out := ToStereo(ConvertType[O](input[L]), ConvertType[O](input[R]))

		sides := O(0)
		for i := 2; i < len(input); i++ {
			sides += ConvertType[O](input[i])
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
		return convertMonoToMultiChannel[O](input[0])
	case 2:
		return convertStereoToMultiChannel[O](ToStereo(input[L], input[R]))
	default:
		switch len(output) {
		case 1:
			return MultiChannel[O]{convertMultiChannelToMono[O](input)}
		case 2:
			out := convertMultiChannelToStereo[O](input)
			return MultiChannel[O]{out[L], out[R]}
		default:
			out := ZeroMultiChannel[O](len(output))

			for i := range min(len(input), len(output)) {
				out[i] = ConvertType[O](input[i])
			}

			return out
		}
	}
}
