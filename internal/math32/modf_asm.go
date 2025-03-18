//go:build amd64

package math32

const haveArchModf = true

func archModf(f float32) (integer, frac float32)
