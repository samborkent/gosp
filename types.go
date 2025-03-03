package gosp

import "golang.org/x/exp/constraints"

type SignedType interface {
	constraints.Float | constraints.Signed
}

type Type interface {
	SignedType | constraints.Unsigned
}

type (
	Mono                 Type
	Stereo[T Type]       [2]T
	MultiChannel[T Type] []T
)

type SampleType[T Type] interface {
	Mono | Stereo[T] | MultiChannel[T]
}
