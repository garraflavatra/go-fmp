package fmp

import (
	"encoding/hex"
	"fmt"
	"io"
)

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
			fmt.Printf("0x%02x (pos %v, unknown)\n", sect.Payload[0], pos)
		} else {
			fmt.Printf("0x%02x (pos %v, type %v)\n", sect.Payload[0], pos, int(chunk.Type))
		}

		if err == io.EOF {
			println("break1")
			break
		}
		if err != nil {
			println(hex.EncodeToString(sect.Payload))
			println("break2")
			return err
		}
		if chunk == nil {
			println("break3")
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

func (sect *FmpSector) processChunks(dict *FmpDict, currentPath *[]uint64) error {
	err := sect.readChunks()
	if err != nil {
		return err
	}

	for _, chunk := range sect.Chunks {
		switch chunk.Type {
		case FMP_CHUNK_PATH_PUSH:
			*currentPath = append(*currentPath, uint64(chunk.Value[0]))

		case FMP_CHUNK_PATH_POP:
			if len(*currentPath) > 0 {
				*currentPath = (*currentPath)[:len(*currentPath)-1]
			}

		case FMP_CHUNK_SIMPLE_DATA:
			dict.set(*currentPath, chunk.Value)

		case FMP_CHUNK_SEGMENTED_DATA:
			// Todo: take index into account
			dict.set(
				*currentPath,
				append(dict.getValue(*currentPath), chunk.Value...),
			)

		case FMP_CHUNK_SIMPLE_KEY_VALUE:
			dict.set(
				append(*currentPath, uint64(chunk.Key)),
				chunk.Value,
			)

		case FMP_CHUNK_LONG_KEY_VALUE:
			dict.set(
				append(*currentPath, uint64(chunk.Key)), // todo: ??
				chunk.Value,
			)

		case FMP_CHUNK_NOOP:
			// noop
		}

		if chunk.Delayed {
			if len(*currentPath) == 0 {
				println("warning: delayed pop without path")
			} else {
				*currentPath = (*currentPath)[:len(*currentPath)-1]
			}
		}
	}
	return nil
}

func (sect *FmpSector) readChunk(payload []byte) (*FmpChunk, error) {

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
