//go:build !amd64 && !arm64 && !wasm

package math32

const haveArchFloor = false

func archFloor(x float32) float32 {
	panic("math32: archFloor: not implemented")
}
