package gsp

import (
	"unsafe"
)

// ConvertSlice converts a slice of audio samples from one type to another.
func ConvertSlice[O Type, I Type](out []O, in []I) (length int) {
	length = min(len(in), len(out))
	if length == 0 {
		return
	}

	switch unsafe.Sizeof(in[0]) {
	case 1: // 8-bit input
		switch unsafe.Sizeof(out[0]) {
		case 1: // 8-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint8 -> uint8
					return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
				}

				// uint8 -> int8
				for i := range length {
					out[i] = O(int8(in[i]) - maxInt8 - 1)
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint8
					for i := range length {
						out[i] = O(uint8(in[i]) + zeroUint8)
					}

					return length
				}

				// int8 -> int8
				return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
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

					return length
				}

				// uint8 -> int16
				for i := range length {
					out[i] = O(int16(int8(in[i])-maxInt8-1) << 8)
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int8 -> uint16
					for i := range length {
						out[i] = O(uint16(uint8(in[i])+zeroUint8) << 8)
					}

					return length
				}

				// int8 -> int16
				for i := range length {
					out[i] = O(int16(in[i]) << 8)
				}

				return length
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

					return length
				case isSigned[O]():
					// uint8 -> int32
					for i := range length {
						out[i] = O(int32(int8(in[i])-maxInt8-1) << 24)
					}

					return length
				default:
					// uint8 -> float32
					for i := range length {
						if in[i] > 0 {
							out[i] = O(float32(int8(in[i])-maxInt8-1) * invMaxInt8_32)
						} else {
							out[i] = O(minFloat32)
						}
					}

					return length
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int8 -> uint32
					for i := range length {
						out[i] = O(uint32(uint8(in[i])+zeroUint8) << 24)
					}

					return length
				case isSigned[O]():
					// int8 -> int32
					for i := range length {
						out[i] = O(int32(in[i]) << 24)
					}

					return length
				default:
					// int8 -> float32
					for i := range length {
						if in[i] >= I(minInt8) {
							out[i] = O(float32(in[i]) * invMaxInt8_32)
						} else {
							out[i] = O(minFloat32)
						}
					}

					return length
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 8-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint8 -> float64
				for i := range length {
					if in[i] > 0 {
						out[i] = O(float64(int8(in[i])-maxInt8-1) * invMaxInt8_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return length
			case isSigned[I]():
				// int8 -> float64
				for i := range length {
					if in[i] >= I(minInt8) {
						out[i] = O(float64(in[i]) * invMaxInt8_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return length
			default:
				panic(errUnknownInputType.Error())
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 2: // 16-bit input
		switch unsafe.Sizeof(out[0]) {
		case 1: // 16-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint16 -> uint8
					for i := range length {
						out[i] = O(uint8(uint16(in[i]) >> 8))
					}

					return length
				}

				// uint16 -> int8
				for i := range length {
					out[i] = O(int8((int16(in[i]) - maxInt16 - 1) >> 8))
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint8
					for i := range length {
						out[i] = O(uint8((uint16(in[i]) + zeroUint16) >> 8))
					}

					return length
				}

				// int16 -> int8
				for i := range length {
					out[i] = O(int8(int16(in[i]) >> 8))
				}

				return length
			default:
				panic(errUnknownInputType.Error())
			}
		case 2: // 16-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint16 -> uint16
					return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
				}

				// uint16 -> int16
				for i := range length {
					out[i] = O(int16(in[i]) - maxInt16 - 1)
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int16 -> uint16
					for i := range length {
						out[i] = O(uint16(in[i]) + zeroUint16)
					}

					return length
				}

				// int16 -> int16
				return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
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
						out[i] = O(int32(int16(in[i])-maxInt16-1) << 16)
					}

					return length
				default:
					// uint16 -> float32
					for i := range length {
						out[i] = O(float32(int16(in[i])-maxInt16-1) * invMaxInt16_32)
					}

					return length
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int16 -> uint32
					for i := range length {
						out[i] = O(uint32(uint16(in[i])+zeroUint16) << 16)
					}

					return length
				case isSigned[O]():
					// int16 -> int32
					for i := range length {
						out[i] = O(int32(in[i]) << 16)
					}

					return length
				default:
					// int16 -> float32
					for i := range length {
						if in[i] >= I(minInt16) {
							out[i] = O(float32(in[i]) * invMaxInt16_32)
						} else {
							out[i] = O(minFloat32)
						}
					}

					return length
				}
			default:
				panic(errUnknownInputType.Error())
			}
		case 8: // 16-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint16 -> float64
				for i := range length {
					if in[i] > 0 {
						out[i] = O(float64(int16(in[i])-maxInt16-1) * invMaxInt16_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return length
			case isSigned[I]():
				// int16 -> float64
				for i := range length {
					if in[i] >= I(minInt16) {
						out[i] = O(float64(in[i]) * invMaxInt16_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return
			default:
				panic(errUnknownInputType.Error())
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 4: // 32-bit input
		switch unsafe.Sizeof(out[0]) {
		case 1: // 32-bit -> 8-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint8
					for i := range length {
						out[i] = O(uint8(uint32(in[i]) >> 24))
					}

					return length
				}

				// uint32 -> int8
				for i := range length {
					out[i] = O(int8((int32(in[i]) - maxInt32 - 1) >> 24))
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint8
					for i := range length {
						out[i] = O(uint8((uint32(in[i]) + zeroUint32) >> 24))
					}

					return length
				}

				// int32 -> int8
				for i := range length {
					out[i] = O(int8(int32(in[i]) >> 24))
				}

				return length
			default:
				if isUnsigned[O]() {
					// float32 -> uint8
					for i := range length {
						if in[i] >= I(minFloat32) && in[i] <= 1 {
							out[i] = O(uint8(quantize32[int8](float32(in[i])*maxInt8_32)) + zeroUint8)
						} else if in[i] > 1 {
							out[i] = O(maxUint8)
						} else {
							out[i] = 1
						}
					}

					return length
				}

				// float32 -> int8
				for i := range length {
					if in[i] >= I(minFloat32) && in[i] <= 1 {
						out[i] = O(quantize32[int8](float32(in[i]) * maxInt8_32))
					} else if in[i] > 1 {
						out[i] = O(maxInt8)
					} else {
						out[i] = O(minInt8)
					}
				}

				return length
			}
		case 2: // 32-bit -> 16-bit
			switch {
			case isUnsigned[I]():
				if isUnsigned[O]() {
					// uint32 -> uint16
					for i := range length {
						out[i] = O(uint16(uint32(in[i]) >> 16))
					}

					return length
				}

				// uint32 -> int16
				for i := range length {
					out[i] = O(int16((int32(in[i]) - maxInt32 - 1) >> 16))
				}

				return length
			case isSigned[I]():
				if isUnsigned[O]() {
					// int32 -> uint16
					for i := range length {
						out[i] = O(uint16((uint32(in[i]) + zeroUint32) >> 16))
					}

					return length
				}

				// int32 -> int16
				for i := range length {
					out[i] = O(int16(int32(in[i]) >> 16))
				}

				return length
			default:
				if isUnsigned[O]() {
					// float32 -> uint16
					for i := range length {
						if in[i] >= I(minFloat32) && in[i] <= 1 {
							out[i] = O(uint16(quantize32[int16](float32(in[i])*maxInt16_32)) + zeroUint16)
						} else if in[i] > 1 {
							out[i] = O(maxUint16)
						} else {
							out[i] = 1
						}
					}

					return length
				}

				// float32 -> int16
				for i := range length {
					if in[i] >= I(minFloat32) && in[i] <= 1 {
						out[i] = O(quantize32[int16](float32(in[i]) * maxInt16_32))
					} else if in[i] > 1 {
						out[i] = O(maxInt16)
					} else {
						out[i] = O(minInt16)
					}
				}

				return
			}
		case 4: // 32-bit -> 32-bit
			switch {
			case isUnsigned[I]():
				switch {
				case isUnsigned[O]():
					// uint32 -> uint32
					return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
				case isSigned[O]():
					// uint32 -> int32
					for i := range length {
						out[i] = O(int32(in[i]) - maxInt32 - 1)
					}

					return length
				default:
					// uint32 -> float32
					for i := range length {
						if in[i] > 0 {
							out[i] = O(float32(int32(in[i])-maxInt32-1) * invMaxInt32_32)
						} else {
							out[i] = O(minFloat32)
						}
					}

					return length
				}
			case isSigned[I]():
				switch {
				case isUnsigned[O]():
					// int32 -> uint32
					for i := range length {
						out[i] = O(uint32(in[i]) + zeroUint32)
					}

					return length
				case isSigned[O]():
					// int32 -> int32
					return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
				default:
					// int32 -> float32
					for i := range length {
						if in[i] >= I(minInt32) {
							out[i] = O(float32(in[i]) * invMaxInt32_32)
						} else {
							out[i] = O(minFloat32)
						}
					}

					return length
				}
			default:
				switch {
				case isUnsigned[O]():
					// float32 -> uint32
					for i := range length {
						if in[i] >= I(minFloat32) && in[i] <= 1 {
							out[i] = O(uint32(quantize64[int32](float64(in[i])*maxInt32_64)) + zeroUint32)
						} else if in[i] > 1 {
							out[i] = O(maxUint32)
						} else {
							out[i] = 1
						}
					}

					return length
				case isSigned[O]():
					// float32 -> int32
					for i := range length {
						if in[i] >= I(minFloat32) && in[i] <= 1 {
							out[i] = O(quantize64[int32](float64(in[i]) * maxInt32_64))
						} else if in[i] > 1 {
							out[i] = O(maxInt32)
						} else {
							out[i] = O(minInt32)
						}
					}

					return length
				default:
					// float32 -> float32
					return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
				}
			}
		case 8: // 32-bit -> 64-bit
			switch {
			case isUnsigned[I]():
				// uint32 -> float64
				for i := range length {
					if in[i] > 0 {
						out[i] = O(float64(int32(in[i])-maxInt32-1) * invMaxInt32_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return length
			case isSigned[I]():
				// int32 -> float64
				for i := range length {
					if in[i] >= I(minInt32) {
						out[i] = O(float64(in[i]) * invMaxInt32_64)
					} else {
						out[i] = O(minFloat64)
					}
				}

				return length
			default:
				// float32 -> float64
				for i := range length {
					out[i] = O(float64(in[i]))
				}

				return length
			}
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	case 8: // 64-bit input
		switch unsafe.Sizeof(out[0]) {
		case 1: // 64-bit -> 8-bit
			if isUnsigned[O]() {
				// float64 -> uint8
				for i := range length {
					if in[i] >= I(minFloat64) && in[i] <= 1 {
						out[i] = O(uint8(quantize64[int8](float64(in[i])*maxInt8_64)) + zeroUint8)
					} else if in[i] > 1 {
						out[i] = O(maxUint8)
					} else {
						out[i] = 1
					}
				}

				return length
			}

			// float64 -> int8
			for i := range length {
				if in[i] >= I(minFloat64) && in[i] <= 1 {
					out[i] = O(quantize64[int8](float64(in[i]) * maxInt8_64))
				} else if in[i] > 1 {
					out[i] = O(maxInt8)
				} else {
					out[i] = O(minInt8)
				}
			}

			return length
		case 2: // 64-bit -> 16-bit
			if isUnsigned[O]() {
				// float64 -> uint16
				for i := range length {
					if in[i] >= I(minFloat64) && in[i] <= 1 {
						out[i] = O(uint16(quantize64[int16](float64(in[i])*maxInt16_64)) + zeroUint16)
					} else if in[i] > 1 {
						out[i] = O(maxUint16)
					} else {
						out[i] = 1
					}
				}

				return length
			}

			// float64 -> int16
			for i := range length {
				if in[i] >= I(minFloat64) && in[i] <= 1 {
					out[i] = O(quantize64[int16](float64(in[i]) * maxInt16_64))
				} else if in[i] > 1 {
					out[i] = O(maxInt16)
				} else {
					out[i] = O(minInt16)
				}
			}

			return length
		case 4: // 64-bit -> 32-bit
			switch {
			case isUnsigned[O]():
				// float64 -> uint32
				for i := range length {
					if in[i] >= I(minFloat64) && in[i] <= 1 {
						out[i] = O(uint32(quantize64[int32](float64(in[i])*maxInt32_64)) + zeroUint32)
					} else if in[i] > 1 {
						out[i] = O(maxUint32)
					} else {
						out[i] = 1
					}
				}

				return length
			case isSigned[O]():
				// float64 -> int32
				for i := range length {
					if in[i] >= I(minFloat64) && in[i] <= 1 {
						out[i] = O(quantize64[int32](float64(in[i]) * maxInt32_64))
					} else if in[i] > 1 {
						out[i] = O(maxInt32)
					} else {
						out[i] = O(minInt32)
					}
				}

				return length
			default:
				// float64 -> float32
				for i := range length {
					out[i] = O(float32(in[i]))
				}

				return length
			}
		case 8: // 64-bit -> 64-bit
			// float64 -> float64
			return copy(out, unsafe.Slice((*O)(unsafe.Pointer(&in[0])), length))
		default:
			panic(errUnknownOutputBitSize.Error())
		}
	default:
		panic(errUnknownInputBitSize.Error())
	}
}
