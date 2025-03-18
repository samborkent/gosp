package processors

import (
	"math"
	"unsafe"

	"github.com/samborkent/gsp"
	"github.com/samborkent/gsp/internal/gmath"
)

// TODO: multi-channel processing
// TODO: add buffer processing

type CompanderAlgorithm string

const (
	CompanderAlgorithmALaw  CompanderAlgorithm = "A-law"
	CompanderAlgorithmMuLaw CompanderAlgorithm = "mu-law"
	CompanderAlgorithmSine  CompanderAlgorithm = "sine"
)

const (
	paramA = 87.6
	invA   = 1 / paramA

	paramMu = 255
	invMu   = 1 / paramMu

	halfPi    = 0.5 * math.Pi
	invHalfPi = 1 / halfPi
)

var (
	calcA1 = 1 / (1 + math.Log(paramA))
	calcA2 = 1 + math.Log(paramA)
	calcMu = 1 / (1 + math.Log(paramMu))
)

type Compander[F gsp.Frame[T], T gsp.Float] struct {
	algorithm CompanderAlgorithm
	expand    bool
	mode      Mode
}

func NewCompander[F gsp.Frame[T], T gsp.Float](algorithm CompanderAlgorithm, expand bool) *Compander[F, T] {
	compander := &Compander[F, T]{
		algorithm: algorithm,
		expand:    expand,
	}

	switch any(*new(F)).(type) {
	case T:
		compander.mode = ModeMono
	case [2]T, gsp.Stereo[T]:
		compander.mode = ModeStereo
	case []T, gsp.MultiChannel[T]:
		compander.mode = ModeMultiChannel
	default:
		panic("gsp: NewCompander: unknown audio frame type")
	}

	return compander
}

func (p *Compander[F, T]) Process(sample F) F {
	switch p.mode {
	case ModeMono:
		monoSample := *(*T)(unsafe.Pointer(&sample))

		var processedSample T

		switch p.algorithm {
		case CompanderAlgorithmALaw:
			processedSample = p.processALaw(monoSample)
		case CompanderAlgorithmMuLaw:
			processedSample = p.processMuLaw(monoSample)
		case CompanderAlgorithmSine:
			processedSample = p.processSine(monoSample)
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}

		return *(*F)(unsafe.Pointer(&processedSample))
	case ModeStereo:
		stereoSample := *(*gsp.Stereo[T])(unsafe.Pointer(&sample))

		var processedSamples [2]T

		switch p.algorithm {
		case CompanderAlgorithmALaw:
			processedSamples[gsp.L] = p.processALaw(stereoSample.L())
			processedSamples[gsp.R] = p.processALaw(stereoSample.R())
		case CompanderAlgorithmMuLaw:
			processedSamples[gsp.L] = p.processMuLaw(stereoSample.L())
			processedSamples[gsp.R] = p.processMuLaw(stereoSample.R())
		case CompanderAlgorithmSine:
			processedSamples[gsp.L] = p.processSine(stereoSample.L())
			processedSamples[gsp.R] = p.processSine(stereoSample.R())
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}

		return *(*F)(unsafe.Pointer(&processedSamples))
	case ModeMultiChannel:
		multiChannelSample := *(*gsp.MultiChannel[T])(unsafe.Pointer(&sample))

		processedSamples := make([]T, len(multiChannelSample))

		switch p.algorithm {
		case CompanderAlgorithmALaw:
			for i := range multiChannelSample {
				processedSamples[i] = p.processALaw(multiChannelSample[i])
			}
		case CompanderAlgorithmMuLaw:
			for i := range multiChannelSample {
				processedSamples[i] = p.processMuLaw(multiChannelSample[i])
			}
		case CompanderAlgorithmSine:
			for i := range multiChannelSample {
				processedSamples[i] = p.processSine(multiChannelSample[i])
			}
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}

		return *(*F)(unsafe.Pointer(&processedSamples))
	default:
		return *new(F)
	}
}

func (p *Compander[F, T]) ProcessBuffer(output, input []F) {
	size := min(len(output), len(input))
	if size == 0 {
		return
	}

	switch p.mode {
	case ModeMono:
		samplePtr := (*T)(unsafe.Pointer(&input[0]))
		monoSamples := unsafe.Slice(samplePtr, len(input))

		switch p.algorithm {
		case CompanderAlgorithmALaw:
			for i := range size {
				processedSample := p.processALaw(monoSamples[i])
				output[i] = *(*F)(unsafe.Pointer(&processedSample))
			}
		case CompanderAlgorithmMuLaw:
			for i := range size {
				processedSample := p.processMuLaw(monoSamples[i])
				output[i] = *(*F)(unsafe.Pointer(&processedSample))
			}
		case CompanderAlgorithmSine:
			for i := range size {
				processedSample := p.processSine(monoSamples[i])
				output[i] = *(*F)(unsafe.Pointer(&processedSample))
			}
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}
	case ModeStereo:
		samplePtr := (*gsp.Stereo[T])(unsafe.Pointer(&input[0]))
		stereoSamples := unsafe.Slice(samplePtr, len(input))

		switch p.algorithm {
		case CompanderAlgorithmALaw:
			for i := range size {
				processedSamples := [2]T{p.processALaw(stereoSamples[i].L()), p.processALaw(stereoSamples[i].R())}
				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		case CompanderAlgorithmMuLaw:
			for i := range size {
				processedSamples := [2]T{p.processMuLaw(stereoSamples[i].L()), p.processMuLaw(stereoSamples[i].R())}
				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		case CompanderAlgorithmSine:
			for i := range size {
				processedSamples := [2]T{p.processSine(stereoSamples[i].L()), p.processSine(stereoSamples[i].R())}
				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}
	case ModeMultiChannel:
		samplePtr := (*gsp.MultiChannel[T])(unsafe.Pointer(&input[0]))
		multiChannelSamples := unsafe.Slice(samplePtr, len(input))

		switch p.algorithm {
		case CompanderAlgorithmALaw:

			for i := range size {
				processedSamples := make([]T, len(multiChannelSamples[i])) // TODO: use pool

				for j := range multiChannelSamples[i] {
					processedSamples[j] = p.processALaw(multiChannelSamples[i][j])
				}

				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		case CompanderAlgorithmMuLaw:
			for i := range size {
				processedSamples := make([]T, len(multiChannelSamples[i])) // TODO: use pool

				for j := range multiChannelSamples[i] {
					processedSamples[j] = p.processMuLaw(multiChannelSamples[i][j])
				}

				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		case CompanderAlgorithmSine:
			for i := range size {
				processedSamples := make([]T, len(multiChannelSamples[i])) // TODO: use pool

				for j := range multiChannelSamples[i] {
					processedSamples[j] = p.processSine(multiChannelSamples[i][j])
				}

				output[i] = *(*F)(unsafe.Pointer(&processedSamples))
			}
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}
	}
}

func (p *Compander[F, T]) processALaw(sample T) T {
	abs, sgn := absSgn(sample)

	if p.expand {
		switch {
		case abs < T(calcA1):
			return sample * T(calcA2) * invA
		case (abs >= T(calcA1)) && (abs < 1):
			return sgn * T(math.Exp(-1+float64(abs)*calcA2)*invA)
		default:
			return sgn
		}
	}

	switch {
	case abs < invA:
		return paramA * sample * T(calcA1)
	case (abs >= invA) && (abs < 1):
		return sgn * T((1+math.Log(paramA*float64(abs)))*calcA1)
	default:
		return sgn
	}
}

func (p *Compander[F, T]) processMuLaw(sample T) T {
	abs, sgn := absSgn(sample)

	if p.expand {
		if abs < 1 {
			return sgn * (T(math.Pow(1+paramMu, float64(abs))) - 1) * invMu
		}

		return sgn
	}

	if abs < 1 {
		return sgn * T(math.Log(1+paramMu*float64(abs))*calcMu)
	}

	return sgn
}

func (p *Compander[F, T]) processSine(sample T) T {
	abs, sgn := absSgn(sample)

	if p.expand {
		if abs < 1 {
			return T(math.Asin(float64(sample))) * invHalfPi
		}

		return sgn
	}

	if abs < 1 {
		return T(math.Sin(halfPi * float64(sample)))
	}

	return sgn
}

func absSgn[T gsp.Float](sample T) (T, T) {
	absX := gmath.Abs(sample)

	if gmath.Signbit(sample) {
		return absX, -1
	} else {
		return absX, 1
	}
}
