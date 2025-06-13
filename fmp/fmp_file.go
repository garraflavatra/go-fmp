package fmp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	sectorSize        = 4096
	sectorHeaderSize  = 20
	sectorPayloadSize = sectorSize - sectorHeaderSize

	magicSequence = "\x00\x01\x00\x00\x00\x02\x00\x01\x00\x05\x00\x02\x00\x02\xC0"
	hbamSequence  = "HBAM7"

	headerSize = sectorSize
	magicSize  = len(magicSequence)
	hbamSize   = len(hbamSequence)
)

func OpenFile(path string) (*FmpFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	stream, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	ctx := &FmpFile{stream: stream, Dictionary: &FmpDict{}}
	if err := ctx.readHeader(); err != nil {
		return nil, err
	}

	ctx.FileSize = uint(info.Size())
	ctx.numSectors = uint64((ctx.FileSize / sectorSize) - 1)
	ctx.Sectors = make([]*FmpSector, 0)
	currentPath := make([]uint64, 0)

	for {
		sector, err := ctx.readSector()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ctx.Sectors = append(ctx.Sectors, sector)
		ctx.Chunks = append(ctx.Chunks, sector.Chunks...)

		for _, chunk := range sector.Chunks {
			switch chunk.Type {
			case FMP_CHUNK_PATH_PUSH:
				currentPath = append(currentPath, uint64(chunk.Value[0]))

			case FMP_CHUNK_PATH_POP:
				if len(currentPath) > 0 {
					currentPath = currentPath[:len(currentPath)-1]
				}

			case FMP_CHUNK_SIMPLE_DATA:
				ctx.Dictionary.set(currentPath, chunk.Value)

			case FMP_CHUNK_SEGMENTED_DATA:
				// Todo: take index into account
				ctx.Dictionary.set(
					currentPath,
					append(ctx.Dictionary.getValue(currentPath), chunk.Value...),
				)

			case FMP_CHUNK_SIMPLE_KEY_VALUE:
				ctx.Dictionary.set(
					append(currentPath, uint64(chunk.Key)),
					chunk.Value,
				)

			case FMP_CHUNK_LONG_KEY_VALUE:
				ctx.Dictionary.set(
					append(currentPath, uint64(chunk.Key)), // todo: ??
					chunk.Value,
				)

			case FMP_CHUNK_NOOP:
				// noop
			}

			if chunk.Delayed {
				if len(currentPath) == 0 {
					println("warning: delayed pop without path")
				} else {
					currentPath = currentPath[:len(currentPath)-1]
				}
			}
		}

		ctx.currentSectorID = sector.NextID
		if sector.NextID == 0 {
			break
		} else if sector.NextID > ctx.numSectors {
			return nil, ErrBadHeader
		} else {
			ctx.stream.Seek(int64(sector.NextID*sectorSize), 0)
		}
	}

	return ctx, nil
}

func (ctx *FmpFile) readHeader() error {
	buf := make([]byte, headerSize)
	_, err := ctx.stream.Read(buf)
	if err != nil {
		return err
	}
	if !bytes.Equal(buf[:magicSize], []byte(magicSequence)) {
		return ErrBadMagic
	}
	if !bytes.Equal(buf[magicSize:magicSize+hbamSize], []byte(hbamSequence)) {
		return ErrBadMagic
	}

	ctx.VersionDate, err = time.Parse("06JAN02", string(buf[531:538]))
	if err != nil {
		return ErrBadHeader
	}

	appNameLength := int(buf[541])
	ctx.CreatorName = string(buf[542 : 542+appNameLength])

	return nil
}

func (ctx *FmpFile) readSector() (*FmpSector, error) {
	println("------- Reading sector", ctx.currentSectorID)
	buf := make([]byte, sectorHeaderSize)
	n, err := ctx.stream.Read(buf)

	if n == 0 {
		return nil, io.EOF
	}
	if err != nil {
		return nil, ErrRead
	}

	sector := &FmpSector{
		ID:      ctx.currentSectorID,
		Deleted: buf[0] > 0,
		Level:   uint8(buf[1]),
		PrevID:  parseVarUint64(buf[4 : 4+4]),
		NextID:  parseVarUint64(buf[8 : 8+4]),
	}

	if ctx.currentSectorID == 0 && sector.PrevID > 0 {
		return nil, ErrBadSectorHeader
	}

	payload := make([]byte, sectorPayloadSize)
	n, err = ctx.stream.Read(payload)

	if n != sectorPayloadSize {
		return nil, ErrRead
	}
	if err != nil {
		return nil, ErrRead
	}

	sector.Chunks = make([]*FmpChunk, 0)

	for {
		chunk, err := ctx.readChunk(payload)
		fmt.Printf("0x%02X", payload[0])
		if chunk != nil {
			fmt.Printf(" (type %v)\n", int(chunk.Type))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if chunk == nil {
			break
		}
		if chunk.Length == 0 {
			panic("chunk length not set")
		}
		sector.Chunks = append(sector.Chunks, chunk)
		payload = payload[min(chunk.Length, uint64(len(payload))):]
		if len(payload) == 0 || (len(payload) == 1 && payload[0] == 0x00) {
			break
		}
	}

	return sector, nil
}
