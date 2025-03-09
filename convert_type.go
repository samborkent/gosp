package gsp

import (
	"errors"
	"math"
	"unsafe"
)

var (
	errUnknownInputType     = errors.New("gsp: convertType: unknown input type encountered")
	errUnknownInputBitSize  = errors.New("gsp: convertType: unknown input bit size encountered")
	errUnknownOutputBitSize = errors.New("gsp: convertType: unknown output bit size encountered")
)

// ConvertType converts any audio sample type to any other sample type.
func ConvertType[O Type, I Type](in I) (out O) {
	switch unsafe.Sizeof(in) {
	case 1: // 8-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 8-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint8 -> uint8
					return O(in)
				}

				// uint8 -> int8
				return O(int8(int16(in) - int16(zeroUint8)))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint8
					return O(uint8(int16(in) + int16(zeroUint8)))
				}

				// int8 -> int8
				return O(in)
			default:
				panic(errUnknownInputType.Error())
			}
		case 2: // 8-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint8 -> uint16
					return O(uint16(in) << 8)
				}

				// uint8 -> int16
				return O((int16(in) - int16(zeroUint8)) << 8)
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint16
					return O((int16(in) + int16(zeroUint8)) << 8)
				}

				// int8 -> int16
				return O(int16(in) << 8)
			default:
				panic(errUnknownInputType.Error())
			}
		case 4: // 8-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint8 -> uint32
					return O(uint32(in) << 24)
				case isSigned[O]():
					// uint8 -> int32
					return O(int32(int16(in)-int16(zeroUint8)) << 24)
				default:
					// uint8 -> float32
					return O(float32(int16(in)-int16(zeroUint8)) / (-math.MinInt8))
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int8 -> uint32
					return O(uint32(int16(in)+int16(zeroUint8)) << 24)
				case isSigned[O]():
					// int8 -> int32
					return O(int32(in) << 24)
				default:
					// int8 -> float32
					return O(float32(in) / (-math.MinInt8))
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 8-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint8 -> float64
				return O(float64(int16(in)-int16(zeroUint8)) / (-math.MinInt8))
			case isSigned[I]():
				// int8 -> float64
				return O(float64(in) / (-math.MinInt8))
			default:
				panic(errUnknownInputType.Error())
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 2: // 16-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 16-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint16 -> uint8
					return O(uint8(uint16(in) >> 8))
				}

				// uint16 -> int8
				return O(int8(int16(int32(in)-int32(zeroUint16)) >> 8))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint8
					return O(uint8(uint16(int32(in)+int32(zeroUint16)) >> 8))
				}

				// int16 -> int8
				return O(int8(int16(in) >> 8))
			default:
				panic(errUnknownInputType.Error())
			}
		case 2: // 16-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint16 -> uint16
					return O(in)
				}

				// uint16 -> int16
				return O(int16(int32(in) - int32(zeroUint16)))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint16
					return O(uint16(int32(in) + int32(zeroUint16)))
				}

				// int16 -> int16
				return O(in)
			default:
				panic(errUnknownInputType.Error())
			}
		case 4: // 16-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint16 -> uint32
					return O(uint32(in) << 16)
				case isSigned[O]():
					// uint16 -> int32
					return O(int32(int32(in)-int32(zeroUint16)) << 16)
				default:
					// uint16 -> float32
					return O(float32(int32(in)-int32(zeroUint16)) / (-math.MinInt16))
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int16 -> uint32
					return O(uint32(int32(in)+int32(zeroUint16)) << 16)
				case isSigned[O]():
					// int16 -> int32
					return O(int32(in) << 16)
				default:
					// int16 -> float32
					return O(float32(in) / (-math.MinInt16))
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 16-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint16 -> float64
				return O(float64(int32(in)-int32(zeroUint16)) / (-math.MinInt16))
			case isSigned[I]():
				// int16 -> float64
				return O(float64(in) / (-math.MinInt16))
			default:
				panic(errUnknownInputType.Error())
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 4: // 32-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 32-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint8
					return O(uint8(uint32(in) >> 24))
				}

				// uint32 -> int8
				return O(int8(int32(int64(in)-int64(zeroUint32)) >> 24))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint8
					return O(uint8(uint32(int64(in)+int64(zeroUint32)) >> 24))
				}

				// int32 -> int8
				return O(int8(int32(in) >> 24))
			default:
				if isUnsigned[O]() {
					// float32 -> uint8
					return O(uint8(int16(float32(in)*math.MaxInt8) + int16(zeroUint8)))
				}

				// float32 -> int8
				return O(int8(float32(in) * math.MaxInt8))
			}
		case 2: // 32-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint16
					return O(uint16(uint32(in) >> 16))
				}

				// uint32 -> int16
				return O(int16(int32(int64(in)-int64(zeroUint32)) >> 16))
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint16
					return O(uint16(uint32(int64(in)+int64(zeroUint32)) >> 16))
				}

				// int32 -> int16
				return O(int16(int32(in) >> 16))
			default:
				if isUnsigned[O]() {
					// float32 -> uint16
					return O(uint16(int32(float32(in)*math.MaxInt16) + int32(zeroUint16)))
				}

				// float32 -> int16
				return O(int16(float32(in) * math.MaxInt16))
			}
		case 4: // 32-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint32 -> uint32
					return O(in)
				case isSigned[O]():
					// uint32 -> int32
					return O(int32(int64(in) - int64(zeroUint32)))
				default:
					// uint32 -> float32
					return O(float32(float64(int64(in)-int64(zeroUint32)) / (-math.MinInt32)))
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int32 -> uint32
					return O(uint32(int64(in) + int64(zeroUint32)))
				case isSigned[O]():
					// int32 -> int32
					return O(in)
				default:
					// int32 -> float32
					return O(float32(float64(in) / (-math.MinInt32)))
				}
			default:
				switch {
				case isUnsigned[O]():
					// float32 -> uint32
					return O(uint32(int64(float64(in)*math.MaxInt32) + int64(zeroUint32)))
				case isSigned[O]():
					// float32 -> int32
					return O(int32(float64(in) * math.MaxInt32))
				default:
					// float32 -> float32
					return O(in)
				}
			}
		case 8: // 32-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint32 -> float64
				return O(float64(int64(in)-int64(zeroUint32)) / (-math.MinInt32))
			case isSigned[I]():
				// int32 -> float64
				return O(float64(in) / (-math.MinInt32))
			default:
				// float32 -> float64
				return O(float64(in))
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 8: // 64-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 64-bit -> 8-bit
			if isUnsigned[O]() {
				// float64 -> uint8
				return O(uint8(int16(float64(in)*math.MaxInt8) + int16(zeroUint8)))
			}

			// float64 -> int8
			return O(int8(float64(in) * math.MaxInt8))
		case 2: // 64-bit -> 16-bit
			if isUnsigned[O]() {
				// float64 -> uint16
				return O(uint16(int32(float64(in)*math.MaxInt16) + int32(zeroUint16)))
			}

			// float64 -> int16
			return O(int16(float64(in) * math.MaxInt16))
		case 4: // 64-bit -> 32-bit
			switch {
			case isUnsigned[O]():
				// float64 -> uint32
				return O(uint32(int64(float64(in)*math.MaxInt32) + int64(zeroUint32)))
			case isSigned[O]():
				// float64 -> int32
				return O(int32(float64(in) * math.MaxInt32))
			default:
				// float64 -> float32
				return O(float32(in))
			}
		case 8: // 64-bit -> 64-bit
			// float64 -> float64
			return O(in)
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	default:
		panic(errUnknownInputBitSize.Error())
	}
}
