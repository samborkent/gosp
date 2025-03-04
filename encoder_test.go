package gosp_test

import (
	"bytes"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/samborkent/gosp"
)

func TestEncoderEncode(t *testing.T) {
	t.Parallel()

	N := 10

	t.Run("uint8 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]gosp.Mono[uint8], N)
		for i := range N {
			input[i] = gosp.Mono[uint8]{uint8(rand.UintN(math.MaxUint8))}
		}

		want := make([]byte, N)
		for i := range N {
			want[i] = input[i][0]
		}

		testEncodeMono(t, input, want)
	})

	// t.Run("int8 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]int8, N)
	// 	for i := range N {
	// 		data[i] = int8((2*rand.Float32() - 1) * math.MinInt8)
	// 	}

	// 	input := make([]byte, N)
	// 	expected := make([]gosp.Mono[uint8], N)
	// 	for i := range N {
	// 		input[i] = byte(data[i])
	// 		expected[i] = gosp.Mono[uint8]{input[i]}
	// 	}

	// 	testDecodeMono(t, input, expected)
	// })

	// t.Run("uint16 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]uint16, N)
	// 	for i := range N {
	// 		data[i] = uint16(rand.UintN(math.MaxUint16))
	// 	}

	// 	input := make([]byte, 2*N)
	// 	expected := make([]gosp.Mono[uint16], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint16(input[2*i:2*i+2], data[i])
	// 		expected[i] = gosp.Mono[uint16]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("int16 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]int16, N)
	// 	for i := range N {
	// 		data[i] = int16((2*rand.Float32() - 1) * math.MinInt16)
	// 	}

	// 	input := make([]byte, 2*N)
	// 	expected := make([]gosp.Mono[int16], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint16(input[2*i:2*i+2], uint16(data[i]))
	// 		expected[i] = gosp.Mono[int16]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("uint32 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]uint32, N)
	// 	for i := range N {
	// 		data[i] = rand.Uint32()
	// 	}

	// 	input := make([]byte, 4*N)
	// 	expected := make([]gosp.Mono[uint32], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint32(input[4*i:4*i+4], data[i])
	// 		expected[i] = gosp.Mono[uint32]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("int32 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]int32, N)
	// 	for i := range N {
	// 		data[i] = int32((2*rand.Float32() - 1) * math.MinInt32)
	// 	}

	// 	input := make([]byte, 4*N)
	// 	expected := make([]gosp.Mono[int32], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint32(input[4*i:4*i+4], uint32(data[i]))
	// 		expected[i] = gosp.Mono[int32]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("float32 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]float32, N)
	// 	for i := range N {
	// 		data[i] = 2*rand.Float32() - 1
	// 	}

	// 	input := make([]byte, 4*N)
	// 	expected := make([]gosp.Mono[float32], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint32(input[4*i:4*i+4], math.Float32bits(data[i]))
	// 		expected[i] = gosp.Mono[float32]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("uint64 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]uint64, N)
	// 	for i := range N {
	// 		data[i] = rand.Uint64()
	// 	}

	// 	input := make([]byte, 8*N)
	// 	expected := make([]gosp.Mono[uint64], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint64(input[8*i:8*i+8], data[i])
	// 		expected[i] = gosp.Mono[uint64]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("int64 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]int64, N)
	// 	for i := range N {
	// 		data[i] = int64((2*rand.Float64() - 1) * math.MinInt64)
	// 	}

	// 	input := make([]byte, 8*N)
	// 	expected := make([]gosp.Mono[int64], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint64(input[8*i:8*i+8], uint64(data[i]))
	// 		expected[i] = gosp.Mono[int64]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })

	// t.Run("float64 mono", func(t *testing.T) {
	// 	t.Parallel()

	// 	data := make([]float64, N)
	// 	for i := range N {
	// 		data[i] = 2*rand.Float64() - 1
	// 	}

	// 	input := make([]byte, 8*N)
	// 	expected := make([]gosp.Mono[float64], N)
	// 	for i := range N {
	// 		binary.LittleEndian.PutUint64(input[8*i:8*i+8], math.Float64bits(data[i]))
	// 		expected[i] = gosp.Mono[float64]{data[i]}
	// 	}

	// 	testEncodeMono(t, input, expected)
	// })
}

func testEncodeMono[T gosp.Type](t *testing.T, input []gosp.Mono[T], want []byte) {
	t.Helper()

	buf := new(bytes.Buffer)
	encoder := gosp.NewEncoder[gosp.Mono[T], T](buf)

	err := encoder.Encode(input)
	if err != nil {
		t.Errorf("encoding samples: error: %s", err.Error())
	}

	output := make([]byte, len(want))

	_, err = buf.Read(output)
	if err != nil {
		t.Errorf("reading buffer: error: %s", err.Error())
	}

	for i := range output {
		if output[i] != want[i] {
			t.Errorf("binary mismatch at index '%d': got '%v', want '%v'", i, output[i], want[i])
		}
	}
}
