// Adapted from math/all_test.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32_test

import . "github.com/samborkent/gsp/internal/math32"

var vf = []float32{
	4.97901195e+00,
	7.73887245e+00,
	-2.76880059e-01,
	-5.01060369e+00,
	9.63629373e+00,
	2.92637726e+00,
	5.22908346e+00,
	2.72793992e+00,
	1.82530850e+00,
	-8.68592473e+00,
}

var ceilBaseSC = []float32{
	Inf(-1),
	Copysign(0, -1),
	0,
	Inf(1),
	NaN(),
}

func alike(a, b float32) bool {
	switch {
	case IsNaN(a) && IsNaN(b):
		return true
	case a == b:
		return Signbit(a) == Signbit(b)
	}
	return false
}
