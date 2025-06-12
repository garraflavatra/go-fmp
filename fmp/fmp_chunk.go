package fmp

import (
	"encoding/binary"
)

type FmpChunk struct {
	Type   FmpChunkType
	Length uint32
	Key    uint32 // If Type == FMP_CHUNK_SHORT_KEY_VALUE or FMP_CHUNK_LONG_KEY_VALUE
	Index  uint32 // Segment index, if Type == FMP_CHUNK_SEGMENTED_DATA
	Value  []byte
}

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
		length := 3 + (payload[0] - 0x10)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+length],
			Length: 1 + uint32(length),
		}, nil
	}
	if 0x12 <= payload[0] && payload[0] <= 0x15 {
		length := 1 + 2*(payload[0]-0x10)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+length],
			Length: 1 + uint32(length),
		}, nil
	}
	if 0x1A <= payload[0] && payload[0] <= 0x1D {
		length := 2 * (payload[0] - 0x19)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_DATA,
			Value:  payload[1 : 1+length],
			Length: 1 + uint32(length),
		}, nil
	}

	// Simple key-value

	if payload[0] == 0x01 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(payload[1]),
			Value:  payload[2 : 2+1],
			Length: 3,
		}, nil
	}
	if 0x02 <= payload[0] && payload[0] <= 0x05 {
		length := 2 * (payload[0] - 1)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(payload[1]),
			Value:  payload[2 : 2+length],
			Length: 2 + uint32(length),
		}, nil
	}
	if payload[0] == 0x06 {
		length := payload[2]
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(payload[1]),
			Value:  payload[3 : 3+length], // docs say offset 2?
			Length: 3 + uint32(length),
		}, nil
	}
	if payload[0] == 0x09 {
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(binary.BigEndian.Uint16(payload[1:3])),
			Value:  payload[3 : 3+1], // docs say offset 2?
			Length: 4,
		}, nil
	}
	if 0x0A <= payload[0] && payload[0] <= 0x0D {
		length := 2 * (payload[0] - 9)
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(binary.BigEndian.Uint16(payload[1:3])),
			Value:  payload[3 : 3+length], // docs say offset 2?
			Length: 2 + uint32(length),
		}, nil
	}
	if payload[0] == 0x0E {
		length := payload[2]
		return &FmpChunk{
			Type:   FMP_CHUNK_SIMPLE_KEY_VALUE,
			Key:    uint32(binary.BigEndian.Uint16(payload[1:3])),
			Value:  payload[4 : 4+length], // docs say offset 2?
			Length: 4 + uint32(length),
		}, nil
	}

	// Long key-value

	if payload[0] == 0x16 {
		length := payload[4]
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint32(payload[1 : 1+3]),
			Value:  payload[5 : 5+length],
			Length: 5 + uint32(length),
		}, nil
	}
	if payload[0] == 0x17 {
		length := uint32(binary.BigEndian.Uint16(payload[4 : 4+2]))
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    uint32(binary.BigEndian.Uint16(payload[1 : 1+2])),
			Value:  payload[6 : 6+length],
			Length: 6 + uint32(length),
		}, nil
	}
	if payload[0] == 0x1E {
		keyLength := payload[1]
		valueLength := payload[2+keyLength]
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint32(payload[2 : 2+keyLength]),
			Value:  payload[2+keyLength+1 : 2+keyLength+1+valueLength],
			Length: 3 + uint32(keyLength) + uint32(valueLength),
		}, nil
	}
	if payload[0] == 0x1F {
		keyLength := uint32(payload[1])
		valueLength := parseVarUint32(payload[2+keyLength : 2+keyLength+2+1])
		return &FmpChunk{
			Type:   FMP_CHUNK_LONG_KEY_VALUE,
			Key:    parseVarUint32(payload[2 : 2+keyLength]),
			Value:  payload[2+keyLength+2 : 2+keyLength+2+valueLength],
			Length: 4 + uint32(keyLength) + uint32(valueLength),
		}, nil
	}

	// Segmented data

	if payload[0] == 0x07 {
		length := binary.BigEndian.Uint16(payload[2 : 2+2])
		payloadLimit := min(4+length, uint16(len(payload)))
		return &FmpChunk{
			Type:   FMP_CHUNK_SEGMENTED_DATA,
			Index:  uint32(payload[1]),
			Value:  payload[4:payloadLimit],
			Length: 4 + uint32(length),
		}, nil
	}
	if payload[0] == 0x0F {
		length := uint32(binary.BigEndian.Uint16(payload[3 : 3+2]))
		payloadLimit := min(5+length, uint32(len(payload)))
		return &FmpChunk{
			Type:   FMP_CHUNK_SEGMENTED_DATA,
			Index:  uint32(binary.BigEndian.Uint16(payload[1 : 1+2])),
			Value:  payload[5:payloadLimit],
			Length: 5 + length,
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
		length := payload[1]
		return &FmpChunk{
			Type:   FMP_CHUNK_PATH_PUSH,
			Value:  payload[2 : 2+length],
			Length: 2 + uint32(length),
		}, nil
	}

	// Path pop

	if payload[0] == 0x3D && payload[1] == 0x40 {
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

func parseVarUint32(payload []byte) uint32 {
	var length uint32
	n := len(payload)
	if n > 4 {
		n = 4 // clamp to max uint32
	}
	for i := range n {
		length <<= 8
		length |= uint32(payload[i])
	}
	return length
}
