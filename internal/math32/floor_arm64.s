// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// func archFloor(x float32) float32
TEXT ·archFloor(SB),NOSPLIT,$0
    FMOVS   x+0(FP), F0     // Load 32-bit float into F0
    FRINTMS F0, F0          // Round to integer toward -inf (floor)
    FMOVS   F0, ret+4(FP)   // Store 32-bit result
    RET
