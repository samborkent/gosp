package processors

import (
	"unsafe"

	"github.com/samborkent/gsp"
)

type Gain[F gsp.Frame[T], T gsp.Float] struct {
	Gain T

	mode Mode
}

func NewGain[F gsp.Frame[T], T gsp.Float](gain T) *Gain[F, T] {
	gainProcessor := &Gain[F, T]{Gain: gsp.DBToLinear(gain)}

	switch len(*new(F)) {
	case 1:
		gainProcessor.mode = ModeMono
	case 2:
		gainProcessor.mode = ModeStereo
	default:
		gainProcessor.mode = ModeMultiChannel
	}

	return gainProcessor
}

func (p *Gain[F, T]) Process(sample F) F {
	switch p.mode {
	case ModeMono:
		monoSample := *(*gsp.Mono[T])(unsafe.Pointer(&sample))
		processedSample := gsp.ToMono(monoSample.M() * p.Gain)
		return *(*F)(unsafe.Pointer(&processedSample))
	case ModeStereo:
		stereoSample := *(*gsp.Stereo[T])(unsafe.Pointer(&sample))
		processedSample := stereoSample.Multiply(p.Gain)
		return *(*F)(unsafe.Pointer(&processedSample))
	case ModeMultiChannel:
		multiChannelSample := *(*gsp.MultiChannel[T])(unsafe.Pointer(&sample))
		processedSample := multiChannelSample.Multiply(p.Gain)
		return *(*F)(unsafe.Pointer(&processedSample))
	default:
		return *new(F)
	}
}

func (p *Gain[F, T]) ProcessBuffer(output, input []F) {
	size := min(len(output), len(input))

	if size == 0 {
		return
	}

	switch p.mode {
	case ModeMono:
		samplePtr := (*gsp.Mono[T])(unsafe.Pointer(&input[0]))
		monoSamples := unsafe.Slice(samplePtr, len(input))

		for i := range size {
			processedSample := gsp.ToMono(monoSamples[i].M() * p.Gain)
			output[i] = *(*F)(unsafe.Pointer(&processedSample))
		}
	case ModeStereo:
		samplePtr := (*gsp.Stereo[T])(unsafe.Pointer(&input[0]))
		stereoSamples := unsafe.Slice(samplePtr, len(input))

		for i := range size {
			processedSample := stereoSamples[i].Multiply(p.Gain)
			output[i] = *(*F)(unsafe.Pointer(&processedSample))
		}
	case ModeMultiChannel:
		samplePtr := (*gsp.MultiChannel[T])(unsafe.Pointer(&input[0]))
		multiChannelSamples := unsafe.Slice(samplePtr, len(input))

		for i := range size {
			processedSample := multiChannelSamples[i].Multiply(p.Gain)
			output[i] = *(*F)(unsafe.Pointer(&processedSample))
		}
	}
}
