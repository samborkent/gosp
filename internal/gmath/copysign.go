package gmath

import (
	"math"
	"unsafe"

	"golang.org/x/exp/constraints"

	math32 "github.com/samborkent/math32"
)

// Copysign returns a value with the magnitude of f
// and the sign of sign.
func Copysign[float constraints.Float](f, sign float) float {
	switch unsafe.Sizeof(f) {
	case 4:
		return float(math32.Copysign(float32(f), float32(sign)))
	case 8:
		return float(math.Copysign(float64(f), float64(sign)))
	default:
		panic("gmath: Copysign: unknown float bit size detected")
	}
}
