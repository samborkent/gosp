//go:build !amd64

package math32

const haveArchModf = false

func archModf(f float32) (integer, frac float32) {
	panic("math32: archModf: not implemented")
}
