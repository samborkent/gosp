package gmath

import (
	"math"
	"unsafe"

	"golang.org/x/exp/constraints"

	math32 "github.com/samborkent/math32"
)

// Round returns the nearest integer, rounding half away from zero.
//
// Special cases are:
//
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round[float constraints.Float](x float) float {
	switch unsafe.Sizeof(x) {
	case 4:
		return float(math32.Round(float32(x)))
	case 8:
		return float(math.Round(float64(x)))
	default:
		panic("gmath: Round: unknown float bit size detected")
	}
}
