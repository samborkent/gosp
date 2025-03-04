package gosp

import (
	"io"
)

// LPCMEncoder is a linear pulse-code modulation encoder.
// This encoder can convert samples into their binary representation.
type LPCMEncoder[S SampleType[T], T Type] struct {
	w io.Writer
}

func NewLPCMEncoder[S SampleType[T], T Type](w io.Writer) *LPCMEncoder[S, T] {
	return &LPCMEncoder[S, T]{
		w: w,
	}
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *LPCMEncoder[S, T]) Encode(s []S) error {
	panic("gosp: LPCMEncoder.Encode: implement")
}
