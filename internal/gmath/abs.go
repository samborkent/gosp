package gmath

import (
	"math"
	"unsafe"

	"golang.org/x/exp/constraints"

	math32 "github.com/samborkent/gsp/internal/math32"
)

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs[float constraints.Float](x float) float {
	switch unsafe.Sizeof(x) {
	case 4:
		return float(math32.Abs(float32(x)))
	case 8:
		return float(math.Abs(float64(x)))
	default:
		panic("gmath: Abs: unknown float bit size detected")
	}
}
