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

func TestDecodeEncode(t *testing.T) {
	t.Parallel()

	N := 1

	t.Run("uint8 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]uint8, N)
		for i := range N {
			data[i] = uint8(rand.UintN(math.MaxUint8))
		}

		input := make([]byte, N)
		for i := range N {
			input[i] = byte(data[i])
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("int8 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int8, N)
		for i := range N {
			data[i] = int8((2*rand.Float32() - 1) * math.MinInt8)
		}

		input := make([]byte, N)
		for i := range N {
			input[i] = byte(data[i])
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("uint16 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]uint16, N)
		for i := range N {
			data[i] = uint16(rand.UintN(math.MaxUint16))
		}

		input := make([]byte, 2*N)
		for i := range N {
			binary.LittleEndian.PutUint16(input[2*i:2*(i+1)], data[i])
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("int16 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int16, N)
		for i := range N {
			data[i] = int16((2*rand.Float32() - 1) * math.MinInt16)
		}

		input := make([]byte, 2*N)
		for i := range N {
			binary.LittleEndian.PutUint16(input[2*i:2*(i+1)], uint16(data[i]))
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("uint32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]uint32, N)
		for i := range N {
			data[i] = rand.Uint32()
		}

		input := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*(i+1)], data[i])
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("int32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]int32, N)
		for i := range N {
			data[i] = int32((2*rand.Float32() - 1) * math.MinInt32)
		}

		input := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*i+4], uint32(data[i]))
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("float32 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]float32, N)
		for i := range N {
			data[i] = 2*rand.Float32() - 1
		}

		input := make([]byte, 4*N)
		for i := range N {
			binary.LittleEndian.PutUint32(input[4*i:4*(i+1)], math.Float32bits(data[i]))
		}

		testDecodeEncodeMono(t, data, input)
	})

	t.Run("float64 mono", func(t *testing.T) {
		t.Parallel()

		data := make([]float64, N)
		for i := range N {
			data[i] = 2*rand.Float64() - 1
		}

		input := make([]byte, 8*N)
		for i := range N {
			binary.LittleEndian.PutUint64(input[8*i:8*(i+1)], math.Float64bits(data[i]))
		}

		testDecodeEncodeMono(t, data, input)
	})
}

func testDecodeEncodeMono[T gsp.Type](t *testing.T, data []T, input []byte) {
	t.Helper()

	decoder := gsp.NewDecoder[T, T](bytes.NewReader(input))

	if decoder.Channels() != 1 {
		t.Errorf("wrong number of channels: got '%d', want '%d'", decoder.Channels(), 1)
	}

	byteSize := int(unsafe.Sizeof(T(0)))

	if decoder.ByteSize() != byteSize {
		t.Errorf("wrong byte size: got '%d', want '%d'", decoder.ByteSize(), byteSize)
	}

	samples := make([]T, len(input)/byteSize)

	err := decoder.Decode(samples)
	if err != nil {
		t.Fatalf("decoding samples: error: %s", err.Error())
	}

	if len(samples) != len(data) {
		t.Fatalf("wrong number of samples: got '%d', want '%d'", len(samples), len(data))
	}

	for i := range samples {
		if samples[i] != data[i] {
			t.Errorf("sample mismatch at index '%d': got '%v', want '%v'", i, samples[i], data[i])
		}
	}

	buf := new(bytes.Buffer)
	encoder := gsp.NewEncoder[T, T](buf)

	err = encoder.Encode(samples)
	if err != nil {
		t.Fatalf("encoding samples: error: %s", err.Error())
	}

	output := make([]byte, buf.Len())

	_, err = buf.Read(output)
	if err != nil {
		t.Fatalf("reading buffer: error: %s", err.Error())
	}

	if len(output) != len(input) {
		t.Fatalf("missing bytes: got '%d', want '%d'", len(output), len(input))
	}

	for i := range output {
		if output[i] != input[i] {
			t.Errorf("binary mismatch at index '%d': got '%v', want '%v'", i, output[i], input[i])
		}
	}
}
