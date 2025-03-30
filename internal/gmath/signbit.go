package gmath

import (
	"math"
	"unsafe"

	"golang.org/x/exp/constraints"

	math32 "github.com/samborkent/math32"
)

// Signbit reports whether x is negative or negative zero.
func Signbit[float constraints.Float](x float) bool {
	switch unsafe.Sizeof(x) {
	case 4:
		return math32.Signbit(float32(x))
	case 8:
		return math.Signbit(float64(x))
	default:
		panic("gmath: Signbit: unknown float bit size detected")
	}
}
