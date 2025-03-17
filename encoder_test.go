package gsp_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand/v2"
	"testing"
	"unsafe"

	"github.com/samborkent/gsp"
)

func TestEncoderEncode(t *testing.T) {
	t.Parallel()

	N := 10

	t.Run("uint8 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]uint8, N)
		for i := range N {
			input[i] = uint8(rand.UintN(math.MaxUint8))
		}

		want := make([]byte, N)
		for i := range N {
			want[i] = input[i]
		}

		testEncodeMono(t, input, want)
	})

	t.Run("int8 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]int8, N)
		for i := range N {
			input[i] = int8((2*rand.Float32() - 1) * math.MinInt8)
		}

		want := make([]byte, N)
		for i := range N {
			want[i] = uint8(input[i])
		}

		testEncodeMono(t, input, want)
	})

	t.Run("uint16 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]uint16, N)
		for i := range N {
			input[i] = uint16(rand.UintN(math.MaxUint16))
		}

		want := make([]byte, 2*N)
		for i := range N {
			binary.LittleEndian.PutUint16(want[2*i:2*(i+1)], input[i])
		}

		testEncodeMono(t, input, want)
	})

	t.Run("int16 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]int16, N)
		for i := range N {
			input[i] = int16((2*rand.Float32() - 1) * math.MinInt16)
		}

		want := make([]byte, 2*N)
		for i := range N {
			binary.LittleEndian.PutUint16(want[2*i:2*(i+1)], uint16(input[i]))
		}

		testEncodeMono(t, input, want)
	})

	t.Run("uint32 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]uint32, N)
		for i := range N {
			input[i] = rand.Uint32()
		}

		want := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(want[4*i:4*(i+1)], input[i])
		}

		testEncodeMono(t, input, want)
	})

	t.Run("int32 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]int32, N)
		for i := range N {
			input[i] = int32((2*rand.Float32() - 1) * math.MinInt32)
		}

		want := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(want[4*i:4*(i+1)], uint32(input[i]))
		}

		testEncodeMono(t, input, want)
	})

	t.Run("float32 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]float32, N)
		for i := range N {
			input[i] = 2*rand.Float32() - 1
		}

		want := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(want[4*i:4*(i+1)], math.Float32bits(input[i]))
		}

		testEncodeMono(t, input, want)
	})

	t.Run("float64 mono", func(t *testing.T) {
		t.Parallel()

		input := make([]float64, N)
		for i := range N {
			input[i] = 2*rand.Float64() - 1
		}

		want := make([]byte, 8*N)
		for i := range N {
			binary.LittleEndian.PutUint64(want[8*i:8*(i+1)], math.Float64bits(input[i]))
		}

		testEncodeMono(t, input, want)
	})
}

func testEncodeMono[T gsp.Type](t *testing.T, input []T, want []byte) {
	t.Helper()

	buf := new(bytes.Buffer)
	encoder := gsp.NewEncoder[T, T](buf)

	if encoder.Channels() != 1 {
		t.Errorf("wrong number of channels: got '%d', want '%d'", encoder.Channels(), 1)
	}

	byteSize := int(unsafe.Sizeof(T(0)))

	if encoder.ByteSize() != byteSize {
		t.Errorf("wrong byte size: got '%d', want '%d'", encoder.ByteSize(), byteSize)
	}

	err := encoder.Encode(input)
	if err != nil {
		t.Fatalf("encoding samples: error: %s", err.Error())
	}

	output := make([]byte, len(input)*byteSize)

	_, err = buf.Read(output)
	if err != nil {
		t.Fatalf("reading buffer: error: %s", err.Error())
	}

	if len(output) != len(want) {
		t.Fatalf("missing samples: got '%d', want '%d'", len(output), len(want))
	}

	for i := range output {
		if output[i] != want[i] {
			t.Errorf("binary mismatch at index '%d': got '%v', want '%v'", i, output[i], want[i])
		}
	}
}
