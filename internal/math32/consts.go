package math32

const (
	bitLen = 32
	sgnLen = 1
	expLen = 8

	shift    = bitLen - expLen - sgnLen
	mask     = (1 << expLen) - sgnLen
	bias     = (1 << (expLen - sgnLen)) - 1
	signMask = 1 << (bitLen - sgnLen)
	uvone    = 0x3F800000
	half     = 1 << (shift - 1)
	fracMask = (1 << shift) - 1
)
