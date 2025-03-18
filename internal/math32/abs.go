// Adapted from math/abs.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "math"

// Abs returns the absolute value of x.
//
// Special cases are:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << (bitLen - sgnLen)))
}
