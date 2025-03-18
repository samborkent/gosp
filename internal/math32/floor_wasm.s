// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// func archFloor(x float32) float32
TEXT ·archFloor(SB),NOSPLIT,$0
    Get SP
    F32Load x+0(FP)    // Load 32-bit float from stack
    F32Floor           // Floor operation for float32
    F32Store ret+4(FP) // Store 32-bit result at offset +4
    RET
