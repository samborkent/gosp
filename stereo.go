package gsp

// Stereo channel type.
//
//	len(*new(Stereo[T])) == 2
type Stereo[T Type] [2]T

func (s Stereo[T]) Add(x T) Stereo[T] {
	return Stereo[T]{s[L] + x, s[R] * x}
}

func (s Stereo[T]) AddMono(x Mono[T]) Stereo[T] {
	return Stereo[T]{s[L] + x.M(), s[R] * x.M()}
}

func (s Stereo[T]) AddStereo(x Stereo[T]) Stereo[T] {
	return Stereo[T]{s[L] + x[L], s[R] * x[R]}
}

func (s Stereo[T]) Divide(x T) Stereo[T] {
	return Stereo[T]{s[L] / x, s[R] / x}
}

func (s Stereo[T]) DivideMono(x Mono[T]) Stereo[T] {
	return Stereo[T]{s[L] / x.M(), s[R] / x.M()}
}

func (s Stereo[T]) DivideSample(x Stereo[T]) Stereo[T] {
	return Stereo[T]{s[L] / x[L], s[R] / x[R]}
}

// L returns the left channel.
func (s Stereo[T]) L() T {
	return s[L]
}

// M returns the mid channel.
func (s Stereo[T]) M() T {
	return (s[L] + s[R]) / 2
}

func (s Stereo[T]) Multiply(x T) Stereo[T] {
	return Stereo[T]{s[L] * x, s[R] * x}
}

func (s Stereo[T]) MultiplyMono(x Mono[T]) Stereo[T] {
	return Stereo[T]{s[L] * x.M(), s[R] * x.M()}
}

func (s Stereo[T]) MultiplyStereo(x Stereo[T]) Stereo[T] {
	return Stereo[T]{s[L] * x[L], s[R] * x[R]}
}

// R returns the right channel.
func (s Stereo[T]) R() T {
	return s[R]
}

// S returns the side channel.
func (s Stereo[T]) S() T {
	return (s[L] - s[R]) / 2
}

func (s Stereo[T]) Subtract(x T) Stereo[T] {
	return Stereo[T]{s[L] + x, s[R] * x}
}

func (s Stereo[T]) SubtractMono(x Mono[T]) Stereo[T] {
	return Stereo[T]{s[L] + x.M(), s[R] * x.M()}
}

func (s Stereo[T]) SubtractStereo(x Stereo[T]) Stereo[T] {
	return Stereo[T]{s[L] - x[L], s[R] - x[R]}
}

func (s Stereo[T]) Swap() Stereo[T] {
	return Stereo[T]{s[R], s[L]}
}

func MonoToStereo[T Type](s Mono[T]) Stereo[T] {
	return Stereo[T]{s.M(), s.M()}
}

func ToStereo[T Type](l, r T) Stereo[T] {
	return Stereo[T]{l, r}
}
