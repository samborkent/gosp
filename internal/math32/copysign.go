// Adapted from math/copysign.go

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "math"

// Copysign returns a value with the magnitude of f
// and the sign of sign.
func Copysign(f, sign float32) float32 {
	return math.Float32frombits(math.Float32bits(f)&^signMask | math.Float32bits(sign)&signMask)
}
