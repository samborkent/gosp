package gsp

// Floating-point type.
type Float interface {
	float32 | float64
}

// Signed integer type.
type Int interface {
	int8 | int16 | int32
}

// Unsigned integer type.
type Uint interface {
	uint8 | uint16 | uint32
}

// Fixed-point type.
type Fixed interface {
	Int | Uint
}

// Signed type.
type Signed interface {
	Int | Float
}

// Unsigned type (alias of Uint).
type Unsigned Uint

// Combination of all types.
type Type interface {
	Signed | Unsigned
}

// Mono channel type.
type Mono Type

// Frame is mono, stereo, or multi-channel sample.
type Frame[T Type] interface {
	Type | [2]T | Stereo[T] | []T | MultiChannel[T]
}
