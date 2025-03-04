package gosp

type floats interface {
	float32 | float64
}

type ints interface {
	int8 | int16 | int32 | int64
}

type SignedType interface {
	ints | floats
}

type UnsignedType interface {
	uint8 | uint16 | uint32 | uint64
}

type Type interface {
	SignedType | UnsignedType
}

type (
	Mono[T Type]         [1]T
	Stereo[T Type]       [2]T
	MultiChannel[T Type] []T
)

type SampleType[T Type] interface {
	Mono[T] | Stereo[T] | MultiChannel[T]
}

// Implementations must not retain p.
type Reader[S SampleType[T], T Type] interface {
	Read(p []S) (n int, err error)
}

// Implementations must not retain p.
type Writer[S SampleType[T], T Type] interface {
	Write(p []S) (n int, err error)
}
