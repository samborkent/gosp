package processors

import (
	"math"
	"unsafe"

	"github.com/samborkent/gsp"
)

// TODO: stereo and multi-channel processing

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

	switch len(*new(F)) {
	case 1:
		compander.mode = ModeMono
	case 2:
		compander.mode = ModeStereo
	default:
		compander.mode = ModeMultiChannel
	}

	return compander
}

func (p *Compander[F, T]) Process(sample F) F {
	switch p.mode {
	case ModeMono:
		monoSample := *(*gsp.Mono[T])(unsafe.Pointer(&sample))

		absX := math.Abs(float64(monoSample.M()))

		var sgnX float64

		if math.Signbit(float64(monoSample.M())) {
			sgnX = -1
		} else {
			sgnX = 1
		}

		processedSample := sgnX

	algoSwitch:
		switch p.algorithm {
		case CompanderAlgorithmALaw:
			if p.expand {
				switch {
				case absX < calcA1:
					processedSample = float64(monoSample.M()) * calcA2 * invA
				case (absX >= calcA1) && (absX < 1):
					processedSample = sgnX * math.Exp(-1+absX*calcA2) * invA
				}

				break algoSwitch
			}

			switch {
			case absX < invA:
				processedSample = paramA * float64(monoSample.M()) * calcA1
			case (absX >= invA) && (absX < 1):
				processedSample = sgnX * (1 + math.Log(paramA*absX)) * calcA1
			}
		case CompanderAlgorithmMuLaw:
			if p.expand {
				if absX < 1 {
					processedSample = sgnX * ((math.Pow(1+paramMu, absX)) - 1) * invMu
				}

				break algoSwitch
			}

			if absX < 1 {
				processedSample = sgnX * math.Log(1+paramMu*absX) * calcMu
			}
		case CompanderAlgorithmSine:
			if p.expand {
				if absX < 1 {
					processedSample = math.Asin(float64(monoSample.M())) * invHalfPi
				}

				break algoSwitch
			}

			if absX < 1 {
				processedSample = math.Sin(halfPi * float64(monoSample.M()))
			}
		default:
			panic("gsp: processors: Compander: algorithm not implemented")
		}

		monoSample = gsp.ToMono(T(processedSample))
		return *(*F)(unsafe.Pointer(&monoSample))
	case ModeStereo:
		panic("gsp: processors: Compander: stereo processing not implemented")
	case ModeMultiChannel:
		panic("gsp: processors: Compander: multi-channel processing not implemented")
	default:
		return *new(F)
	}
}

// func (p *Compander[float]) ProcessBuffer(samples []float) []float {
// 	for i := range samples {
// 		samples[i] = p.Process(samples[i])
// 	}

// 	return samples
// }
