// Adapted from math/modf.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "math"

// Modf returns integer and fractional floating-point numbers
// that sum to f. Both values have the same sign as f.
//
// Special cases are:
//
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN) = NaN, NaN
func Modf(f float32) (integer, frac float32) {
	if haveArchModf {
		return archModf(f)
	}
	return modf(f)
}

func modf(f float32) (integer, frac float32) {
	if f < 1 {
		switch {
		case f < 0:
			integer, frac = Modf(-f)
			return -integer, -frac
		case f == 0:
			return f, f // Return -0, -0 when f == -0
		}
		return 0, f
	}

	x := math.Float32bits(f)
	e := uint(x>>shift)&mask - bias

	// Keep the top 12+e bits, the integer part; clear the rest.
	if e < bitLen-expLen-sgnLen {
		x &^= 1<<(bitLen-expLen-sgnLen-e) - 1
	}
	integer = math.Float32frombits(x)
	frac = f - integer
	return
}
