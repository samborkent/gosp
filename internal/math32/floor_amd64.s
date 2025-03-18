// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// func archFloor(x float32) float32
TEXT ·archFloor(SB),NOSPLIT,$0
    MOVL    x+0(FP), AX          // Load x into AX
    MOVL    $0x7FFFFFFF, DX      // Sign bit mask for 32-bit floats
    ANDL    AX, DX               // DX = |x|
    SUBL    $1, DX               // Prepare for boundary check
    MOVL    $(0x7FFFFF - 1), CX  // (2²³-1 - 1) = 0x7FFFFE
    CMPL    DX, CX               // Compare |x|-1 with boundary
    JAE     isBig_floor          // Jump if |x| >= 2²³-1 or special cases
    
    MOVL    AX, X0               // X0 = x
    CVTTSS2SL X0, AX             // Convert to integer (truncate)
    CVTSL2SS AX, X1              // Convert back to float32
    CMPSS   X1, X0, $1           // Compare x < integer_part (LT)
    MOVSS   $(-1.0), X2          // Load -1.0 (0xBF800000)
    ANDPS   X2, X0               // Mask -1.0 if needed
    ADDSS   X1, X0               // Adjust integer part
    MOVSS   X0, ret+4(FP)        // Store result
    RET

isBig_floor:
    MOVL    AX, ret+4(FP)        // Return original value
    RET
