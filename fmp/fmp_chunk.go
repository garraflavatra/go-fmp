package fmp

func (ctx *FmpFile) readChunk(payload []byte) (*FmpChunk, error) {

	// https://github.com/evanmiller/fmptools/blob/02eb770e59e0866dab213d80e5f7d88e17648031/HACKING
	// https://github.com/Rasmus20B/fmplib/blob/66245e5269275724bacfe1437fb1f73bc587a2f3/src/fmp_format/chunk.rs#L57-L60

	chunk := &FmpChunk{}
	chunkCode := payload[0]

	if (chunkCode & 0xC0) == 0xC0 {
		chunkCode &= 0x3F
		chunk.Delayed = true
	}

	switch chunkCode {
	case 0x00:
		chunk.Length = 2
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x01:
		chunk.Length = 3
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = uint64(payload[1])
		chunk.Value = payload[2:chunk.Length]

	case 0x02, 0x03, 0x04, 0x05:
		valueLength := uint64(2 * (chunkCode - 1))
		chunk.Length = 2 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = uint64(payload[1])
		chunk.Value = payload[2:chunk.Length]

	case 0x06:
		valueLength := uint64(payload[2])
		chunk.Length = 3 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = uint64(payload[1])
		chunk.Value = payload[3:chunk.Length]

	case 0x07:
		valueLength := parseVarUint64(payload[2 : 2+2])
		chunk.Length = min(4+valueLength, uint64(len(payload)))
		chunk.Type = FMP_CHUNK_SEGMENTED_DATA
		chunk.Index = uint64(payload[1])
		chunk.Value = payload[4:chunk.Length]

	case 0x08:
		chunk.Length = 3
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x09:
		chunk.Length = 4
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = parseVarUint64(payload[1 : 1+2])
		chunk.Value = payload[3:chunk.Length]

	case 0x0A, 0x0B, 0x0C, 0x0D:
		valueLength := uint64(2 * (chunkCode - 0x09))
		chunk.Length = 3 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = parseVarUint64(payload[1 : 1+2])
		chunk.Value = payload[3:chunk.Length]

	case 0x0E:
		if payload[1] == 0xFE {
			chunk.Length = 10
			chunk.Type = FMP_CHUNK_PATH_PUSH
			chunk.Value = payload[2:chunk.Length]
			break
		}

		if payload[1] == 0xFF {
			chunk.Length = 7
			chunk.Type = FMP_CHUNK_SIMPLE_DATA
			chunk.Value = payload[2:chunk.Length]
			break
		}

		valueLength := uint64(payload[2])
		chunk.Length = 4 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_KEY_VALUE
		chunk.Key = parseVarUint64(payload[1 : 1+2])
		chunk.Value = payload[4:chunk.Length]

	case 0x0F:
		valueLength := parseVarUint64(payload[3 : 3+2])
		chunk.Length = min(5+valueLength, uint64(len(payload)))
		chunk.Type = FMP_CHUNK_SEGMENTED_DATA
		chunk.Index = parseVarUint64(payload[1 : 1+2])
		chunk.Value = payload[5:chunk.Length]

	case 0x10, 0x11:
		valueLength := 3 + (uint64(chunkCode) - 0x10)
		chunk.Length = 1 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x12, 0x13, 0x14, 0x15:
		valueLength := 1 + 2*(uint64(chunkCode)-0x10)
		chunk.Length = 1 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x16:
		valueLength := uint64(payload[4])
		chunk.Type = FMP_CHUNK_LONG_KEY_VALUE
		chunk.Length = 5 + valueLength
		chunk.Key = parseVarUint64(payload[1 : 1+3])
		chunk.Value = payload[5:chunk.Length]

	case 0x17:
		valueLength := parseVarUint64(payload[4 : 4+2])
		chunk.Length = 6 + valueLength
		chunk.Type = FMP_CHUNK_LONG_KEY_VALUE
		chunk.Key = parseVarUint64(payload[1 : 1+3])
		chunk.Value = payload[6:chunk.Length]

	case 0x19:
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Length = 2
		chunk.Value = payload[1:chunk.Length]

	case 0x1A, 0x1B, 0x1C, 0x1D:
		valueLength := 2 * uint64(chunkCode-0x19)
		chunk.Length = 1 + valueLength
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x1E:
		keyLength := uint64(payload[1])
		valueLength := uint64(payload[2+keyLength])
		chunk.Length = 2 + keyLength + 1 + valueLength
		chunk.Type = FMP_CHUNK_LONG_KEY_VALUE
		chunk.Key = parseVarUint64(payload[2 : 2+keyLength])
		chunk.Value = payload[2+keyLength+1 : chunk.Length]

	case 0x1F:
		keyLength := uint64(uint64(payload[1]))
		valueLength := parseVarUint64(payload[2+keyLength : 2+keyLength+2+1])
		chunk.Length = 2 + keyLength + 2 + valueLength
		chunk.Type = FMP_CHUNK_LONG_KEY_VALUE
		chunk.Key = parseVarUint64(payload[2 : 2+keyLength])
		chunk.Value = payload[2+keyLength+2 : chunk.Length]

	case 0x20:
		if payload[1] == 0xFE {
			chunk.Length = 10
			chunk.Type = FMP_CHUNK_PATH_PUSH
			chunk.Value = payload[1:chunk.Length]
			break
		}

		chunk.Length = 2
		chunk.Type = FMP_CHUNK_PATH_PUSH
		chunk.Value = payload[1:chunk.Length]

	case 0x23:
		chunk.Length = 2
		chunk.Type = FMP_CHUNK_SIMPLE_DATA
		chunk.Value = payload[1:chunk.Length]

	case 0x28:
		chunk.Length = 3
		chunk.Type = FMP_CHUNK_PATH_PUSH
		chunk.Value = payload[1:chunk.Length]

	case 0x30:
		chunk.Length = 4
		chunk.Type = FMP_CHUNK_PATH_PUSH
		chunk.Value = payload[1:chunk.Length]

	case 0x38:
		valueLength := uint64(payload[1])
		chunk.Length = 2 + valueLength
		chunk.Type = FMP_CHUNK_PATH_PUSH
		chunk.Value = payload[2:chunk.Length]

	case 0x3D, 0x40:
		chunk.Type = FMP_CHUNK_PATH_POP
		chunk.Length = 1

	case 0x80:
		chunk.Type = FMP_CHUNK_NOOP
		chunk.Length = 1

	default:
		return nil, ErrBadChunk
	}

	return chunk, nil
}
