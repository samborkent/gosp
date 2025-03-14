package gsp

// Implementations must not retain p.
type Reader[F Frame[T], T Type] interface {
	Read(p []F) (framesRead int, err error)
}

// Implementations must not retain p.
type Writer[F Frame[T], T Type] interface {
	Write(p []F) (framesWrtie int, err error)
}

type SampleReader[F Frame[T], T Type] interface {
	ReadSample() (sample F)
}

type SampleWriter[F Frame[T], T Type] interface {
	WriteSample(sample F)
}
