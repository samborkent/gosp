//go:build amd64 || arm64 || wasm

package math32

const haveArchFloor = true

func archFloor(x float32) float32
