package fmp

func (ctx *FmpFile) readChunk(payload []byte) (*FmpChunk, error) {

	// Simple data

	if payload[0] == 0x00 || payload[0] == 0x19 || payload[0] == 0x23 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+1],
			Length: 2,
		}, nil
	}
	if payload[0] == 0x08 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+2],
			Length: 3,
		}, nil
	}
	if payload[0] == 0x0E && payload[1] == 0xFF {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[2 : 2+5],
			Length: 7,
		}, nil
	}
	if 0x10 <= payload[0] && payload[0] <= 0x11 {
		valueLength := 3 + (payload[0] - 0x10)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+valueLength],
			Length: 1 + uint64(valueLength),
		}, nil
	}
	if 0x12 <= payload[0] && payload[0] <= 0x15 {
		valueLength := 1 + 2*(payload[0]-0x10)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+valueLength],
			Length: 1 + uint64(valueLength),
		}, nil
	}
	if 0x1A <= payload[0] && payload[0] <= 0x1D {
		valueLength := 2 * (payload[0] - 0x19)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+valueLength],
			Length: 1 + uint64(valueLength),
		}, nil
	}

	// Simple key-value

	if payload[0] == 0x01 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint64(payload[1]),
			Value:  payload[2 : 2+1],
			Length: 3,
		}, nil
	}
	if 0x02 <= payload[0] && payload[0] <= 0x05 {
		valueLength := 2 * (payload[0] - 1)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint64(payload[1]),
			Value:  payload[2 : 2+valueLength],
			Length: 2 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x06 {
		valueLength := payload[2]
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint64(payload[1]),
			Value:  payload[2 : 2+valueLength], // docs say offset 2?
			Length: 3 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x09 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    parseVarUint64(payload[1 : 1+2]),
			Value:  payload[3 : 3+1],
			Length: 4,
		}, nil
	}
	if 0x0A <= payload[0] && payload[0] <= 0x0D {
		valueLength := 2 * (payload[0] - 0x09)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    parseVarUint64(payload[1 : 1+2]),
			Value:  payload[3 : 3+valueLength],
			Length: 2 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x0E {
		valueLength := payload[2]
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    parseVarUint64(payload[1 : 1+2]),
			Value:  payload[4 : 4+valueLength],
			Length: 4 + uint64(valueLength),
		}, nil
	}

	// Long key-value

	if payload[0] == 0x16 {
		valueLength := payload[4]
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint64(payload[1 : 1+3]),
			Value:  payload[5 : 5+valueLength],
			Length: 5 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x17 {
		valueLength := parseVarUint64(payload[4 : 4+2])
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint64(payload[1 : 1+3]),
			Value:  payload[6 : 6+valueLength],
			Length: 6 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x1E {
		keyLength := payload[1]
		valueLength := payload[2+keyLength]
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint64(payload[2 : 2+keyLength]),
			Value:  payload[2+keyLength+1 : 2+keyLength+1+valueLength],
			Length: 2 + uint64(keyLength) + 1 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x1F {
		keyLength := uint64(payload[1])
		valueLength := parseVarUint64(payload[2+keyLength : 2+keyLength+2+1])
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint64(payload[2 : 2+keyLength]),
			Value:  payload[2+keyLength+2 : 2+keyLength+2+valueLength],
			Length: 4 + uint64(keyLength) + uint64(valueLength),
		}, nil
	}

	// Segmented data

	if payload[0] == 0x07 {
		valueLength := parseVarUint64(payload[2 : 2+2])
		payloadLimit := min(4+valueLength, uint64(len(payload)))
		return &FmpChunk{
			Type:   FMP_CHUNK_SEGMENTED_DATA,
			Index:  uint64(payload[1]),
			Value:  payload[4:payloadLimit],
			Length: 4 + uint64(valueLength),
		}, nil
	}
	if payload[0] == 0x0F {
		valueLength := parseVarUint64(payload[3 : 3+2])
		payloadLimit := min(5+valueLength, uint64(len(payload)))
		return &FmpChunk{
			Type:   FMP_CHUNK_SEGMENTED_DATA,
			Index:  parseVarUint64(payload[1 : 1+2]),
			Value:  payload[5:payloadLimit],
			Length: 5 + valueLength,
		}, nil
	}

	// Path push

	if payload[0] == 0x20 || payload[0] == 0x0E {
		if payload[1] == 0xFE {
			return &FmpChunk{
				Type:   FMP_CHUNK_PATH_PUSH,
				Value:  payload[1 : 1+8],
				Length: 10,
			}, nil
		}
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_PUSH,
			Value:  payload[1 : 1+1],
			Length: 2,
		}, nil
	}
	if payload[0] == 0x28 {
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_PUSH,
			Value:  payload[1 : 1+2],
			Length: 3,
		}, nil
	}
	if payload[0] == 0x30 {
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_PUSH,
			Value:  payload[1 : 1+3],
			Length: 4,
		}, nil
	}
	if payload[0] == 0x38 {
		valueLength := payload[1]
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_PUSH,
			Value:  payload[2 : 2+valueLength],
			Length: 2 + uint64(valueLength),
		}, nil
	}

	// Path pop

	if payload[0] == 0x3D || payload[0] == 0x40 {
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_POP,
			Length: 1,
		}, nil
	}

	// No-op

	if payload[0] == 0x80 {
		return &FmpChunk{
			Type:   FMP_CHUNK_NOOP,
			Length: 1,
		}, nil
	}

	return nil, nil
}
