package gsp

import (
	"math"
	"testing"

	"github.com/samborkent/gsp/internal/math32"
)

func TestQuantize32(t *testing.T) {
	tests := [][2]float32{
		{-0.49999998, -0},
		{-0.5, -1},
		{-0.5000000000000001, -1},
		{0, 0},
		{0.49999998, 0},
		{0.5, 1},
		{0.5000000000000001, 1},
		{1.390671161567e-309, 0},
		{-0, -0},
	}

	for i, test := range tests {
		res := math32.Round(test[0])

		if res != test[1] {
			t.Errorf("fail %d: got '%g', want '%g'", i, res, test[1])
		}
	}
}

func TestQuantize64(t *testing.T) {
	tests := [][2]float64{
		{-0.49999999999999994, -0},
		{-0.5, -1},
		{-0.5000000000000001, -1},
		{0, 0},
		{0.49999999999999994, 0},
		{0.5, 1},
		{0.5000000000000001, 1},
		{1.390671161567e-309, 0},
		{-0, -0},
	}

	for i, test := range tests {
		res := math.Round(test[0])

		if res != test[1] {
			t.Errorf("fail %d: got '%g', want '%g'", i, res, test[1])
		}
	}
}
