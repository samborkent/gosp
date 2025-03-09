package gsp

import (
	"math"
	"unsafe"
)

// ConvertSlice converts a slice of audio samples from one type to another.
func ConvertSlice[O Type, I Type](out []O, in []I) (length int) {
	length = min(len(in), len(out))
	if length == 0 {
		return
	}

	switch unsafe.Sizeof(in) {
	case 1: // 8-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 8-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint8 -> uint8
					for i := range length {
						out[i] = O(in[i])
					}
					return
				}

				// uint8 -> int8
				for i := range length {
					out[i] = O(int8(int16(in[i]) - int16(zeroUint8)))
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint8
					for i := range length {
						out[i] = O(uint8(int16(in[i]) + int16(zeroUint8)))
					}
					return
				}

				// int8 -> int8
				for i := range length {
					out[i] = O(in[i])
				}
				return
			default:
				panic(errUnknownInputType.Error())
			}
		case 2: // 8-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint8 -> uint16
					for i := range length {
						out[i] = O(uint16(in[i]) << 8)
					}
					return
				}

				// uint8 -> int16
				for i := range length {
					out[i] = O((int16(in[i]) - int16(zeroUint8)) << 8)
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint16
					for i := range length {
						out[i] = O((int16(in[i]) + int16(zeroUint8)) << 8)
					}
					return
				}

				// int8 -> int16
				for i := range length {
					out[i] = O(int16(in[i]) << 8)
				}
				return
			default:
				panic(errUnknownInputType.Error())
			}
		case 4: // 8-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint8 -> uint32
					for i := range length {
						out[i] = O(uint32(in[i]) << 24)
					}
					return
				case isSigned[O]():
					// uint8 -> int32
					for i := range length {
						out[i] = O(int32(int16(in[i])-int16(zeroUint8)) << 24)
					}
					return
				default:
					// uint8 -> float32
					for i := range length {
						out[i] = O(float32(int16(in[i])-int16(zeroUint8)) / (-math.MinInt8))
					}
					return
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int8 -> uint32
					for i := range length {
						out[i] = O(uint32(int16(in[i])+int16(zeroUint8)) << 24)
					}
					return
				case isSigned[O]():
					// int8 -> int32
					for i := range length {
						out[i] = O(int32(in[i]) << 24)
					}
					return
				default:
					// int8 -> float32
					for i := range length {
						out[i] = O(float32(in[i]) / (-math.MinInt8))
					}
					return
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 8-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint8 -> float64
				for i := range length {
					out[i] = O(float64(int16(in[i])-int16(zeroUint8)) / (-math.MinInt8))
				}
				return
			case isSigned[I]():
				// int8 -> float64
				for i := range length {
					out[i] = O(float64(in[i]) / (-math.MinInt8))
				}
				return
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
					for i := range length {
						out[i] = O(uint8(uint16(in[i]) >> 8))
					}
					return
				}

				// uint16 -> int8
				for i := range length {
					out[i] = O(int8(int16(int32(in[i])-int32(zeroUint16)) >> 8))
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint8
					for i := range length {
						out[i] = O(uint8(uint16(int32(in[i])+int32(zeroUint16)) >> 8))
					}
					return
				}

				// int16 -> int8
				for i := range length {
					out[i] = O(int8(int16(in[i]) >> 8))
				}
				return
			default:
				panic(errUnknownInputType.Error())
			}
		case 2: // 16-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint16 -> uint16
					for i := range length {
						out[i] = O(in[i])
					}
					return
				}

				// uint16 -> int16
				for i := range length {
					out[i] = O(int16(int32(in[i]) - int32(zeroUint16)))
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint16
					for i := range length {
						out[i] = O(uint16(int32(in[i]) + int32(zeroUint16)))
					}
					return
				}

				// int16 -> int16
				for i := range length {
					out[i] = O(in[i])
				}
				return
			default:
				panic(errUnknownInputType.Error())
			}
		case 4: // 16-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint16 -> uint32
					for i := range length {
						out[i] = O(uint32(in[i]) << 16)
					}
					return
				case isSigned[O]():
					// uint16 -> int32
					for i := range length {
						out[i] = O(int32(int32(in[i])-int32(zeroUint16)) << 16)
					}
					return
				default:
					// uint16 -> float32
					for i := range length {
						out[i] = O(float32(int32(in[i])-int32(zeroUint16)) / (-math.MinInt16))
					}
					return
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int16 -> uint32
					for i := range length {
						out[i] = O(uint32(int32(in[i])+int32(zeroUint16)) << 16)
					}
					return
				case isSigned[O]():
					// int16 -> int32
					for i := range length {
						out[i] = O(int32(in[i]) << 16)
					}
					return
				default:
					// int16 -> float32
					for i := range length {
						out[i] = O(float32(in[i]) / (-math.MinInt16))
					}
					return
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 16-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint16 -> float64
				for i := range length {
					out[i] = O(float64(int32(in[i])-int32(zeroUint16)) / (-math.MinInt16))
				}
				return
			case isSigned[I]():
				// int16 -> float64
				for i := range length {
					out[i] = O(float64(in[i]) / (-math.MinInt16))
				}
				return
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
					for i := range length {
						out[i] = O(uint8(uint32(in[i]) >> 24))
					}
					return
				}

				// uint32 -> int8
				for i := range length {
					out[i] = O(int8(int32(int64(in[i])-int64(zeroUint32)) >> 24))
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint8
					for i := range length {
						out[i] = O(uint8(uint32(int64(in[i])+int64(zeroUint32)) >> 24))
					}
					return
				}

				// int32 -> int8
				for i := range length {
					out[i] = O(int8(int32(in[i]) >> 24))
				}
				return
			default:
				if isUnsigned[O]() {
					// float32 -> uint8
					for i := range length {
						out[i] = O(uint8(int16(float32(in[i])*math.MaxInt8) + int16(zeroUint8)))
					}
					return
				}

				// float32 -> int8
				for i := range length {
					out[i] = O(int8(float32(in[i]) * math.MaxInt8))
				}
				return
			}
		case 2: // 32-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint16
					for i := range length {
						out[i] = O(uint16(uint32(in[i]) >> 16))
					}
					return
				}

				// uint32 -> int16
				for i := range length {
					out[i] = O(int16(int32(int64(in[i])-int64(zeroUint32)) >> 16))
				}
				return
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint16
					for i := range length {
						out[i] = O(uint16(uint32(int64(in[i])+int64(zeroUint32)) >> 16))
					}
					return
				}

				// int32 -> int16
				for i := range length {
					out[i] = O(int16(int32(in[i]) >> 16))
				}
				return
			default:
				if isUnsigned[O]() {
					// float32 -> uint16
					for i := range length {
						out[i] = O(uint16(int32(float32(in[i])*math.MaxInt16) + int32(zeroUint16)))
					}
					return
				}

				// float32 -> int16
				for i := range length {
					out[i] = O(int16(float32(in[i]) * math.MaxInt16))
				}
				return
			}
		case 4: // 32-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint32 -> uint32
					for i := range length {
						out[i] = O(in[i])
					}
					return
				case isSigned[O]():
					// uint32 -> int32
					for i := range length {
						out[i] = O(int32(int64(in[i]) - int64(zeroUint32)))
					}
					return
				default:
					// uint32 -> float32
					for i := range length {
						out[i] = O(float32(float64(int64(in[i])-int64(zeroUint32)) / (-math.MinInt32)))
					}
					return
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int32 -> uint32
					for i := range length {
						out[i] = O(uint32(int64(in[i]) + int64(zeroUint32)))
					}
					return
				case isSigned[O]():
					// int32 -> int32
					for i := range length {
						out[i] = O(in[i])
					}
					return
				default:
					// int32 -> float32
					for i := range length {
						out[i] = O(float32(float64(in[i]) / (-math.MinInt32)))
					}
					return
				}
			default:
				switch {
				case isUnsigned[O]():
					// float32 -> uint32
					for i := range length {
						out[i] = O(uint32(int64(float64(in[i])*math.MaxInt32) + int64(zeroUint32)))
					}
					return
				case isSigned[O]():
					// float32 -> int32
					for i := range length {
						out[i] = O(int32(float64(in[i]) * math.MaxInt32))
					}
					return
				default:
					// float32 -> float32
					for i := range length {
						out[i] = O(in[i])
					}
					return
				}
			}
		case 8: // 32-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint32 -> float64
				for i := range length {
					out[i] = O(float64(int64(in[i])-int64(zeroUint32)) / (-math.MinInt32))
				}
				return
			case isSigned[I]():
				// int32 -> float64
				for i := range length {
					out[i] = O(float64(in[i]) / (-math.MinInt32))
				}
				return
			default:
				// float32 -> float64
				for i := range length {
					out[i] = O(float64(in[i]))
				}
				return
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 8: // 64-bit input
		switch unsafe.Sizeof(out) {
		case 1: // 64-bit -> 8-bit
			if isUnsigned[O]() {
				// float64 -> uint8
				for i := range length {
					out[i] = O(uint8(int16(float64(in[i])*math.MaxInt8) + int16(zeroUint8)))
				}
				return
			}

			// float64 -> int8
			for i := range length {
				out[i] = O(int8(float64(in[i]) * math.MaxInt8))
			}
			return
		case 2: // 64-bit -> 16-bit
			if isUnsigned[O]() {
				// float64 -> uint16
				for i := range length {
					out[i] = O(uint16(int32(float64(in[i])*math.MaxInt16) + int32(zeroUint16)))
				}
				return
			}

			// float64 -> int16
			for i := range length {
				out[i] = O(int16(float64(in[i]) * math.MaxInt16))
			}
			return
		case 4: // 64-bit -> 32-bit
			switch {
			case isUnsigned[O]():
				// float64 -> uint32
				for i := range length {
					out[i] = O(uint32(int64(float64(in[i])*math.MaxInt32) + int64(zeroUint32)))
				}
				return
			case isSigned[O]():
				// float64 -> int32
				for i := range length {
					out[i] = O(int32(float64(in[i]) * math.MaxInt32))
				}
				return
			default:
				// float64 -> float32
				for i := range length {
					out[i] = O(float32(in[i]))
				}
				return
			}
		case 8: // 64-bit -> 64-bit
			// float64 -> float64
			for i := range length {
				out[i] = O(in[i])
			}
			return
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	default:
		panic(errUnknownInputBitSize.Error())
	}
}
