package gsp

// Mono channel type.
//
//	len(*new(Mono[T])) == 1
type Mono[T Type] [1]T

// M return the mono sample value, equivalent to m[0].
func (m Mono[T]) M() T {
	return m[0]
}

func ToMono[T Type](sample T) Mono[T] {
	return Mono[T]{sample}
}

func ZeroMono[T Type]() Mono[T] {
	return ToMono(T(0))
}
