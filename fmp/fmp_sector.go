package fmp

type FmpSector struct {
	ID      uint64
	Level   uint8
	Deleted bool
	PrevID  uint64
	NextID  uint64
	Prev    *FmpSector
	Next    *FmpSector
	Payload []byte
	Chunks  []*FmpChunk
}

type FmpChunk struct {
	Type   FmpChunkType
	Length uint64
	Key    uint64 // If Type == FMP_CHUNK_SHORT_KEY_VALUE or FMP_CHUNK_LONG_KEY_VALUE
	Index  uint64 // Segment index, if Type == FMP_CHUNK_SEGMENTED_DATA
	Value  []byte
}

func (sect *FmpSector) readChunks() error {
	if len(sect.Chunks) > 0 {
		panic("chunks already read")
	}
	for {
		pos := (sect.ID+1)*sectorSize - uint64(len(sect.Payload))

		if sect.Payload[0] == 0x00 && sect.Payload[1] == 0x00 {
			break
		}

		chunk, err := sect.readChunk(sect.Payload)
		if chunk == nil {
			debug("0x%02x (pos %v, unknown)\n", sect.Payload[0], pos)
		} else {
			debug("0x%02x (pos %v, type %v)\n", sect.Payload[0], pos, int(chunk.Type))
		}

		if err != nil {
			debug("chunk error at sector %d", sect.ID)
			dump(sect.Payload)
			return err
		}
		if chunk == nil {
			break
		}
		if chunk.Length == 0 {
			panic("chunk length not set")
		}

		sect.Chunks = append(sect.Chunks, chunk)
		sect.Payload = sect.Payload[min(chunk.Length, uint64(len(sect.Payload))):]

		if len(sect.Payload) == 0 || (len(sect.Payload) == 1 && sect.Payload[0] == 0x00) {
			break
		}
	}
	return nil
}

func (sect *FmpSector) processChunks(dict *FmpDict) error {
	err := sect.readChunks()
	if err != nil {
		return err
	}

	currentPath := make([]uint64, 0)
	for _, chunk := range sect.Chunks {
		switch chunk.Type {
		case FmpChunkPathPush, FmpChunkPathPushLong:
			currentPath = append(currentPath, decodeVarUint64(chunk.Value))
			dumpPath(currentPath)

		case FmpChunkPathPop:
			if len(currentPath) > 0 {
				currentPath = (currentPath)[:len(currentPath)-1]
			}

		case FmpChunkSimpleData:
			dict.set(currentPath, chunk.Value)

		case FmpChunkSegmentedData:
			// Todo: take index into account
			dict.set(
				currentPath,
				append(dict.GetValue(currentPath...), chunk.Value...),
			)

		case FmpChunkSimpleKeyValue:
			dict.set(
				append(currentPath, uint64(chunk.Key)),
				chunk.Value,
			)

		case FmpChunkLongKeyValue:
			dict.set(
				append(currentPath, uint64(chunk.Key)), // todo: ??
				chunk.Value,
			)

		case FmpChunkNoop:
			// noop
		}
	}
	return nil
}

func (sect *FmpSector) readChunk(payload []byte) (*FmpChunk, error) {

	// https://github.com/evanmiller/fmptools/blob/02eb770e59e0866dab213d80e5f7d88e17648031/HACKING
	// https://github.com/Rasmus20B/fmplib/blob/66245e5269275724bacfe1437fb1f73bc587a2f3/src/fmp_format/chunk.rs#L57-L60

	chunk := &FmpChunk{}
	chunkCode := payload[0]

	switch chunkCode {
	case 0x00:
		chunk.Length = 2
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[1:chunk.Length]

	case 0x01, 0x02, 0x03, 0x04, 0x05:
		chunk.Length = 2 + 2*uint64(chunkCode-0x01) + addIf(chunkCode == 0x01, 1)
		chunk.Type = FmpChunkSimpleKeyValue
		chunk.Key = uint64(payload[1])
		chunk.Value = payload[2:chunk.Length]

	case 0x06:
		chunk.Length = 3 + uint64(payload[2])
		chunk.Type = FmpChunkSimpleKeyValue
		chunk.Key = uint64(payload[1])
		chunk.Value = payload[3:chunk.Length]

	case 0x07:
		valueLength := decodeVarUint64(payload[2 : 2+2])
		chunk.Length = min(4+valueLength, uint64(len(payload)))
		chunk.Type = FmpChunkSegmentedData
		chunk.Index = uint64(payload[1])
		chunk.Value = payload[4:chunk.Length]

	case 0x08:
		chunk.Length = 3
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[1:chunk.Length]

	case 0x09:
		chunk.Length = 4
		chunk.Type = FmpChunkSimpleKeyValue
		chunk.Key = decodeVarUint64(payload[1 : 1+2])
		chunk.Value = payload[3:chunk.Length]

	case 0x0A, 0x0B, 0x0C, 0x0D:
		chunk.Length = 3 + 2*uint64(chunkCode-0x09)
		chunk.Type = FmpChunkSimpleKeyValue
		chunk.Key = decodeVarUint64(payload[1 : 1+2])
		chunk.Value = payload[3:chunk.Length]

	case 0x0E:
		if payload[1] == 0xFF {
			chunk.Length = 7
			chunk.Type = FmpChunkSimpleData
			chunk.Value = payload[2:chunk.Length]
			break
		}

		chunk.Length = 4 + uint64(payload[3])
		chunk.Type = FmpChunkSimpleKeyValue
		chunk.Key = decodeVarUint64(payload[1 : 1+2])
		chunk.Value = payload[4:chunk.Length]

	case 0x0F:
		valueLength := decodeVarUint64(payload[3 : 3+2])
		chunk.Length = uint64(len(payload))
		if chunk.Length > 5+valueLength {
			return nil, ErrBadChunk
		}
		chunk.Type = FmpChunkSegmentedData
		chunk.Index = decodeVarUint64(payload[1 : 1+2])
		chunk.Value = payload[5:chunk.Length]

	case 0x10, 0x11:
		chunk.Length = 4 + addIf(chunkCode == 0x11, 1)
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[1:chunk.Length]

	case 0x12, 0x13, 0x14, 0x15:
		chunk.Length = 4 + 2*(uint64(chunkCode)-0x11)
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[1:chunk.Length]

	case 0x16:
		chunk.Length = 5 + uint64(payload[4])
		chunk.Type = FmpChunkLongKeyValue
		chunk.Key = decodeVarUint64(payload[1 : 1+3])
		chunk.Value = payload[5:chunk.Length]

	case 0x17:
		chunk.Length = 6 + decodeVarUint64(payload[4:4+2])
		chunk.Type = FmpChunkLongKeyValue
		chunk.Key = decodeVarUint64(payload[1 : 1+3])
		chunk.Value = payload[6:chunk.Length]

	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D:
		valueLength := uint64(payload[1])
		chunk.Length = 2 + valueLength + 2*uint64(chunkCode-0x19) + addIf(chunkCode == 0x19, 1)
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[2 : 2+valueLength]

	case 0x1E:
		keyLength := uint64(payload[1])
		valueLength := uint64(payload[2+keyLength])
		chunk.Length = 2 + keyLength + 1 + valueLength
		chunk.Type = FmpChunkLongKeyValue
		chunk.Key = decodeVarUint64(payload[2 : 2+keyLength])
		chunk.Value = payload[2+keyLength+1 : chunk.Length]

	case 0x1F:
		keyLength := uint64(payload[1])
		valueLength := decodeVarUint64(payload[2+keyLength : 2+keyLength+2+1])
		chunk.Length = 2 + keyLength + 2 + valueLength
		chunk.Type = FmpChunkLongKeyValue
		chunk.Key = decodeVarUint64(payload[2 : 2+keyLength])
		chunk.Value = payload[2+keyLength+2 : chunk.Length]

	case 0x20, 0xE0:
		if payload[1] == 0xFE {
			chunk.Length = 10
			chunk.Type = FmpChunkPathPush
			chunk.Value = payload[2:chunk.Length]
			break
		}

		chunk.Length = 2
		chunk.Type = FmpChunkPathPush
		chunk.Value = payload[1:chunk.Length]

	case 0x23:
		chunk.Length = 2 + uint64(payload[1])
		chunk.Type = FmpChunkSimpleData
		chunk.Value = payload[1:chunk.Length]

	case 0x28, 0x30:
		chunk.Length = 3 + addIf(chunkCode == 0x30, 1)
		chunk.Type = FmpChunkPathPush
		chunk.Value = payload[1:chunk.Length]

	case 0x38:
		valueLength := uint64(payload[1])
		chunk.Length = 2 + valueLength
		chunk.Type = FmpChunkPathPushLong
		chunk.Value = payload[2:chunk.Length]

	case 0x3D, 0x40:
		chunk.Type = FmpChunkPathPop
		chunk.Length = 1

	case 0x80:
		chunk.Type = FmpChunkNoop
		chunk.Length = 1

	default:
		return nil, ErrBadChunk
	}

	return chunk, nil
}
