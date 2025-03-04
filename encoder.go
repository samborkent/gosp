package gosp

import (
	"errors"
	"io"
	"unsafe"
)

var ErrEncoderNotinitialized = errors.New("gosp: Encoder: not initialized")

// Encoder is a linear PCM and floating-point encoder.
// This encoder can convert samples into their binary representation.
type Encoder[F Frame[T], T Type] struct {
	w                  io.Writer
	channels, byteSize int
	bigEndian          bool
	initialized        bool
}

func NewEncoder[F Frame[T], T Type](w io.Writer, opts ...EncodingOption) *Encoder[F, T] {
	// Apply encoding options.
	var cfg EncodingConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Encoder[F, T]{
		w:           w,
		channels:    len(*new(F)),
		byteSize:    int(unsafe.Sizeof(*new(T))),
		bigEndian:   cfg.BigEndian,
		initialized: true,
	}
}

// Encode reads from the sample slice and decodes into the internal [io.Writer].
func (e *Encoder[F, T]) Encode(s []F) error {
	if !e.initialized {
		return ErrDecoderNotinitialized
	}

	panic("gosp: Encoder.Encode: implement")
}
