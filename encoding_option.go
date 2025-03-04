package gosp

// EncodingOption is a functional option which is used for both [LPCMEncoder] and [LPCMDecoder].
type EncodingOption func(cfg *EncodingConfig)

// EncodingConfig contains all configuration options for [LPCMEncoder] and [LPCMDecoder].
type EncodingConfig struct {
	BigEndian bool // Enable big-endian encoding for bit sizes above 8-bit.
}

// EncodingBigEndian enables big-endian encoding instead of little-endian encoding.
func EncodingBigEndian(cfg *EncodingConfig) {
	cfg.BigEndian = true
}
