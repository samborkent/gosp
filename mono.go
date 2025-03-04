package gosp

// Mono channel type.
//
//	len(*new(Mono[T])) == 1
type Mono[T Type] [1]T

func ToMono[T Type](sample T) Mono[T] {
	return Mono[T]{sample}
}

// S return the mono sample value, equivalent to m[0].
func (m Mono[T]) S() T {
	return m[0]
}
