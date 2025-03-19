// Adapted from math/all_test.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32_test

import (
	"testing"

	. "github.com/samborkent/gsp/internal/math32"
)

var floor = []float32{
	4.00000000e+00,
	7.00000000e+00,
	-1.00000000e+00,
	-6.00000000e+00,
	9.00000000e+00,
	2.00000000e+00,
	5.00000000e+00,
	2.00000000e+00,
	1.00000000e+00,
	-9.00000000e+00,
}

var vfceilSC = []float32{
	Inf(-1),
	Copysign(0, -1),
	0,
	Inf(1),
	NaN(),
	1<<23 - 1,
	1<<23 - 0.5, // largest fractional float64
	1 << 23,
	-1 << 23,
	-1<<23 + 0.5, // smallest fractional float64
	-1<<23 + 1,
	1 << 24,
	-1 << 24,
}

var floorSC = []float32{
	Inf(-1),
	Copysign(0, -1),
	0,
	Inf(1),
	NaN(),
	1<<23 - 1,
	1<<23 - 1,
	1 << 23,
	-1 << 23,
	-1 << 23,
	-1<<23 + 1,
	1 << 24,
	-1 << 24,
}

func TestFloor(t *testing.T) {
	t.Parallel()

	t.Run("number set", func(t *testing.T) {
		for i := range len(vf) {
			if f := Floor(vf[i]); !alike(floor[i], f) {
				t.Errorf("Floor(%g) = %g, want %g", vf[i], f, floor[i])
			}
		}
	})

	t.Run("special cases", func(t *testing.T) {
		for i := range len(vfceilSC) {
			if f := Floor(vfceilSC[i]); !alike(floorSC[i], f) {
				t.Errorf("%d: Floor(%g) = %g, want %g", i, vfceilSC[i], f, floorSC[i])
			}
		}
	})
}
