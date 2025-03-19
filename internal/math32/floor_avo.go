//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	// . "github.com/mmcloughlin/avo/reg"
)

func main() {
	TEXT("·archFloor", NOSPLIT, "func(x float64) float64")
	Doc("Floor round a float down towards zero")
	x := Load(Param("x"), GP64())
	a := GP64()
	MOVQ(F64(1<<63), a)
	ANDQ(x, a)
	Store(x, ReturnIndex(0))
	RET()
	Generate()
}
