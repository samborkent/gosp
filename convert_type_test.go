package gsp

import (
	"math"
	"math/rand/v2"
	"testing"
)

func checkConvertType[I, O Type](t *testing.T, input I, want O, name string) {
	t.Helper()

	if res := ConvertType[O](input); res != want {
		t.Errorf("wrong %s: got '%v', want '%v'", name, res, want)
	}
}

func TestConvertTypeUint8(t *testing.T) {
	t.Parallel()

	minimum := uint8(0)
	zero := zeroUint8
	maximum := uint8(math.MaxUint8)
	random := uint8(rand.UintN(math.MaxUint8))

	t.Run("uint8->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, minimum, "minimum")
		checkConvertType(t, zero, zero, "zero")
		checkConvertType(t, maximum, maximum, "maximum")
		checkConvertType(t, random, random, "random")
	})
	t.Run("uint8->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(math.MinInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8(int16(random)-int16(zeroUint8)), "random")
	})
	t.Run("uint8->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(0), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint8)<<8, "maximum")
		checkConvertType(t, random, uint16(random)<<8, "random")
	})
	t.Run("uint8->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(math.MinInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt8)<<8, "maximum")
		checkConvertType(t, random, int16(int16(random)-int16(zeroUint8))<<8, "random")
	})
	t.Run("uint8->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint8)<<24, "maximum")
		checkConvertType(t, random, uint32(random)<<24, "random")
	})
	t.Run("uint8->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt8)<<24, "maximum")
		checkConvertType(t, random, int32(int16(random)-int16(zeroUint8))<<24, "random")
	})
	t.Run("uint8->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(int16(maximum)-int16(zeroUint8))/(-math.MinInt8), "maximum")
		checkConvertType(t, random, float32(int16(random)-int16(zeroUint8))/(-math.MinInt8), "random")
	})
	t.Run("uint8->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(int16(maximum)-int16(zeroUint8))/(-math.MinInt8), "maximum")
		checkConvertType(t, random, float64(int16(random)-int16(zeroUint8))/(-math.MinInt8), "random")
	})
}

func TestConvertTypeInt8(t *testing.T) {
	t.Parallel()

	minimum := int8(math.MinInt8)
	zero := int8(0)
	maximum := int8(math.MaxInt8)
	random := int8(int16(rand.UintN(math.MaxUint8)) - int16(zeroUint8))

	t.Run("int8->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(0), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8(int16(random)+int16(zeroUint8)), "random")
	})
	t.Run("int8->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, minimum, "minimum")
		checkConvertType(t, zero, zero, "zero")
		checkConvertType(t, maximum, maximum, "maximum")
		checkConvertType(t, random, random, "random")
	})
	t.Run("int8->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(0), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint8)<<8, "maximum")
		checkConvertType(t, random, uint16(int16(random)+int16(zeroUint8))<<8, "random")
	})
	t.Run("int8->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(math.MinInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt8)<<8, "maximum")
		checkConvertType(t, random, int16(random)<<8, "random")
	})
	t.Run("int8->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint8)<<24, "maximum")
		checkConvertType(t, random, uint32(int16(random)+int16(zeroUint8))<<24, "random")
	})
	t.Run("int8->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt8)<<24, "maximum")
		checkConvertType(t, random, int32(random)<<24, "random")
	})
	t.Run("int8->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(maximum)/(-math.MinInt8), "maximum")
		checkConvertType(t, random, float32(random)/(-math.MinInt8), "random")
	})
	t.Run("int8->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(maximum)/(-math.MinInt8), "maximum")
		checkConvertType(t, random, float64(random)/(-math.MinInt8), "random")
	})
}

func TestConvertTypeUint16(t *testing.T) {
	t.Parallel()

	minimum := uint16(0)
	zero := zeroUint16
	maximum := uint16(math.MaxUint16)
	random := uint16(rand.UintN(math.MaxUint16))

	t.Run("uint16->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(0), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint16(random)>>8, "random")
	})
	t.Run("uint16->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(math.MinInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8((int32(random)-int32(zeroUint16))>>8), "random")
	})
	t.Run("uint16->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, minimum, "minimum")
		checkConvertType(t, zero, zero, "zero")
		checkConvertType(t, maximum, maximum, "maximum")
		checkConvertType(t, random, random, "random")
	})
	t.Run("uint16->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(math.MinInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt16), "maximum")
		checkConvertType(t, random, int16(int32(random)-int32(zeroUint16)), "random")
	})
	t.Run("uint16->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint16)<<16, "maximum")
		checkConvertType(t, random, uint32(random)<<16, "random")
	})
	t.Run("uint16->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt16)<<16, "maximum")
		checkConvertType(t, random, (int32(random)-int32(zeroUint16))<<16, "random")
	})
	t.Run("uint16->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(int32(maximum)-int32(zeroUint16))/(-math.MinInt16), "maximum")
		checkConvertType(t, random, float32(int32(random)-int32(zeroUint16))/(-math.MinInt16), "random")
	})
	t.Run("uint16->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(int32(maximum)-int32(zeroUint16))/(-math.MinInt16), "maximum")
		checkConvertType(t, random, float64(int32(random)-int32(zeroUint16))/(-math.MinInt16), "random")
	})
}

func TestConvertTypeInt16(t *testing.T) {
	t.Parallel()

	minimum := int16(math.MinInt16)
	zero := int16(0)
	maximum := int16(math.MaxInt16)
	random := int16(int32(rand.UintN(math.MaxUint16)) - int32(zeroUint16))

	t.Run("int16->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(0), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8(int16(random>>8)+int16(zeroUint8)), "random")
	})
	t.Run("int16->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(math.MinInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8(random>>8), "random")
	})
	t.Run("int16->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(0), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint16), "maximum")
		checkConvertType(t, random, uint16(int32(random)+int32(zeroUint16)), "random")
	})
	t.Run("int16->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, minimum, "minimum")
		checkConvertType(t, zero, zero, "zero")
		checkConvertType(t, maximum, maximum, "maximum")
		checkConvertType(t, random, random, "random")
	})
	t.Run("int16->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint16)<<16, "maximum")
		checkConvertType(t, random, uint32(int32(random)+int32(zeroUint16))<<16, "random")
	})
	t.Run("int16->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt16)<<16, "maximum")
		checkConvertType(t, random, int32(random)<<16, "random")
	})
	t.Run("int16->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(maximum)/(-math.MinInt16), "maximum")
		checkConvertType(t, random, float32(random)/(-math.MinInt16), "random")
	})
	t.Run("int16->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(maximum)/(-math.MinInt16), "maximum")
		checkConvertType(t, random, float64(random)/(-math.MinInt16), "random")
	})
}

func TestConvertTypeUint32(t *testing.T) {
	t.Parallel()

	minimum := uint32(0)
	zero := zeroUint32
	maximum := uint32(math.MaxUint32)
	random := uint32(rand.Uint32())

	t.Run("uint32->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(0), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8(random>>24), "random")
	})
	t.Run("uint32->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(math.MinInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8((int64(random)-int64(zeroUint32))>>24), "random")
	})
	t.Run("uint32->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(0), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint16), "maximum")
		checkConvertType(t, random, uint16(random>>16), "random")
	})
	t.Run("uint32->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(math.MinInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt16), "maximum")
		checkConvertType(t, random, int16((int64(random)-int64(zeroUint32))>>16), "random")
	})
	t.Run("uint32->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint32), "maximum")
		checkConvertType(t, random, uint32(random), "random")
	})
	t.Run("uint32->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt32), "maximum")
		checkConvertType(t, random, int32(int64(random)-int64(zeroUint32)), "random")
	})
	t.Run("uint32->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(int64(maximum)-int64(zeroUint32))/(-math.MinInt32), "maximum")
		checkConvertType(t, random, float32(int64(random)-int64(zeroUint32))/(-math.MinInt32), "random")
	})
	t.Run("uint32->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(int64(maximum)-int64(zeroUint32))/(-math.MinInt32), "maximum")
		checkConvertType(t, random, float64(int64(random)-int64(zeroUint32))/(-math.MinInt32), "random")
	})
}

func TestConvertTypeInt32(t *testing.T) {
	t.Parallel()

	minimum := int32(math.MinInt32)
	zero := int32(0)
	maximum := int32(math.MaxInt32)
	random := int32(int64(rand.Uint32()) - int64(zeroUint32))

	t.Run("int32->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(0), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8((int64(random)+int64(zeroUint32))>>24), "random")
	})
	t.Run("int32->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(math.MinInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8(random>>24), "random")
	})
	t.Run("int32->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(0), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint16), "maximum")
		checkConvertType(t, random, uint16((int64(random)+int64(zeroUint32))>>16), "random")
	})
	t.Run("int32->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(math.MinInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt16), "maximum")
		checkConvertType(t, random, int16(random>>16), "random")
	})
	t.Run("int32->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(0), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint32), "maximum")
		checkConvertType(t, random, uint32(int64(random)+int64(zeroUint32)), "random")
	})
	t.Run("int32->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(math.MinInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt32), "maximum")
		checkConvertType(t, random, int32(random), "random")
	})
	t.Run("int32->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(maximum)/(-math.MinInt32), "maximum")
		checkConvertType(t, random, float32(random)/(-math.MinInt32), "random")
	})
	t.Run("int32->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(maximum)/(-math.MinInt32), "maximum")
		checkConvertType(t, random, float64(random)/(-math.MinInt32), "random")
	})
}

func TestConvertTypeFloat32(t *testing.T) {
	t.Parallel()

	minimum := float32(-1)
	zero := float32(0)
	maximum := float32(1)
	random := 2*rand.Float32() - 1

	t.Run("float32->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(1), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8(int16(random*math.MaxInt8)+int16(zeroUint8)), "random")
	})
	t.Run("float32->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(-math.MaxInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8(random*math.MaxInt8), "random")
	})
	t.Run("float32->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(1), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint16), "maximum")
		checkConvertType(t, random, uint16(int32(random*math.MaxInt16)+int32(zeroUint16)), "random")
	})
	t.Run("float32->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(-math.MaxInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt16), "maximum")
		checkConvertType(t, random, int16(random*math.MaxInt16), "random")
	})
	t.Run("float32->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(1), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint32), "maximum")
		checkConvertType(t, random, uint32(int64(float64(random)*math.MaxInt32)+int64(zeroUint32)), "random")
	})
	t.Run("float32->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(-math.MaxInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt32), "maximum")
		checkConvertType(t, random, int32(float64(random)*math.MaxInt32), "random")
	})
	t.Run("float32->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(1), "maximum")
		checkConvertType(t, random, float32(random), "random")
	})
	t.Run("float32->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(1), "maximum")
		checkConvertType(t, random, float64(random), "random")
	})
}

func TestConvertTypeFloat64(t *testing.T) {
	t.Parallel()

	minimum := float64(-1)
	zero := float64(0)
	maximum := float64(1)
	random := 2*rand.Float64() - 1

	t.Run("float64->uint8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint8(1), "minimum")
		checkConvertType(t, zero, zeroUint8, "zero")
		checkConvertType(t, maximum, uint8(math.MaxUint8), "maximum")
		checkConvertType(t, random, uint8(int16(random*math.MaxInt8)+int16(zeroUint8)), "random")
	})
	t.Run("float64->int8", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int8(-math.MaxInt8), "minimum")
		checkConvertType(t, zero, int8(0), "zero")
		checkConvertType(t, maximum, int8(math.MaxInt8), "maximum")
		checkConvertType(t, random, int8(random*math.MaxInt8), "random")
	})
	t.Run("float64->uint16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint16(1), "minimum")
		checkConvertType(t, zero, zeroUint16, "zero")
		checkConvertType(t, maximum, uint16(math.MaxUint16), "maximum")
		checkConvertType(t, random, uint16(int32(random*math.MaxInt16)+int32(zeroUint16)), "random")
	})
	t.Run("float64->int16", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int16(-math.MaxInt16), "minimum")
		checkConvertType(t, zero, int16(0), "zero")
		checkConvertType(t, maximum, int16(math.MaxInt16), "maximum")
		checkConvertType(t, random, int16(random*math.MaxInt16), "random")
	})
	t.Run("float64->uint32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, uint32(1), "minimum")
		checkConvertType(t, zero, zeroUint32, "zero")
		checkConvertType(t, maximum, uint32(math.MaxUint32), "maximum")
		checkConvertType(t, random, uint32(int64(float64(random)*math.MaxInt32)+int64(zeroUint32)), "random")
	})
	t.Run("float64->int32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, int32(-math.MaxInt32), "minimum")
		checkConvertType(t, zero, int32(0), "zero")
		checkConvertType(t, maximum, int32(math.MaxInt32), "maximum")
		checkConvertType(t, random, int32(float64(random)*math.MaxInt32), "random")
	})
	t.Run("float64->float32", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float32(-1), "minimum")
		checkConvertType(t, zero, float32(0), "zero")
		checkConvertType(t, maximum, float32(1), "maximum")
		checkConvertType(t, random, float32(random), "random")
	})
	t.Run("float64->float64", func(t *testing.T) {
		t.Parallel()

		checkConvertType(t, minimum, float64(-1), "minimum")
		checkConvertType(t, zero, float64(0), "zero")
		checkConvertType(t, maximum, float64(1), "maximum")
		checkConvertType(t, random, float64(random), "random")
	})
}
