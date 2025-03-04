package gsp

// Floating-point type.
type Float interface {
	float32 | float64
}

// Signed integer type.
type Int interface {
	int8 | int16 | int32 | int64
}

// Unsigned integer type.
type Uint interface {
	uint8 | uint16 | uint32 | uint64
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

// Frame is mono, stereo, or multi-channel sample.
type Frame[T Type] interface {
	Mono[T] | Stereo[T] | MultiChannel[T]
}

// Implementations must not retain p.
type Reader[F Frame[T], T Type] interface {
	Read(p []F) (n int, err error)
}

// Implementations must not retain p.
type Writer[F Frame[T], T Type] interface {
	Write(p []F) (n int, err error)
}
