package gsp

import "math"

const reciprocal20 = float64(1.0 / 20.0)

func DBToLinear[T Float](gain T) T {
	return T(math.Pow(10, float64(gain)*reciprocal20))
}

func LinearToDB[T Float](gain T) T {
	return T(20 * math.Log10(float64(gain)))
}
