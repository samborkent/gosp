package gosp

import (
	"io"
)

// LPCMEncoder is a linear pulse-code modulation encoder.
// This encoder can convert samples into their binary representation.
type LPCMEncoder[F Frame[T], T Type] struct {
	w io.Writer
}

func NewLPCMEncoder[F Frame[T], T Type](w io.Writer) *LPCMEncoder[F, T] {
	return &LPCMEncoder[F, T]{
		w: w,
	}
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *LPCMEncoder[F, T]) Encode(s []F) error {
	panic("gosp: LPCMEncoder.Encode: implement")
}
