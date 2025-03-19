// Adapted from math/floor.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN
func Floor(x float32) float32 {
	// if haveArchFloor {
	// 	return archFloor(x)
	// }

	return floor(x)
}

func floor(x float32) float32 {
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}

	if x < 0 {
		d, fract := Modf(-x)
		if fract != 0.0 {
			d = d + 1
		}

		return -d
	}

	d, _ := Modf(x)
	return d
}
