package gosp

import (
	"io"
)

// Encoder is a linear PCM and floating-point encoder.
// This encoder can convert samples into their binary representation.
type Encoder[F Frame[T], T Type] struct {
	w io.Writer
}

func NewEncoder[F Frame[T], T Type](w io.Writer) *Encoder[F, T] {
	return &Encoder[F, T]{
		w: w,
	}
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *Encoder[F, T]) Encode(s []F) error {
	panic("gosp: Encoder.Encode: implement")
}
