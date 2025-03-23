package gsp

// Implementations must not retain p.
type Reader[F Frame[T], T Type] interface {
	Read(buffer []F) (framesRead int, err error)
}

// Implementations must not retain p.
type Writer[F Frame[T], T Type] interface {
	Write(buffer []F) (framesWritten int, err error)
}

type ReaderFrom[F Frame[T], T Type] interface {
	ReadFrom(r Reader[F, T]) (framesRead int64, err error)
}

type WriterTo[F Frame[T], T Type] interface {
	WriteTo(w Writer[F, T]) (framesWritten int64, err error)
}

type FrameReader[F Frame[T], T Type] interface {
	ReadFrame() (frame F)
}

type FrameWriter[F Frame[T], T Type] interface {
	WriteFrame(frame F)
}
