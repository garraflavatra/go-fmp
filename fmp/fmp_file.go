package fmp

import (
	"bytes"
	"io"
	"os"
	"time"
)

const (
	magicSequence = "\x00\x01\x00\x00\x00\x02\x00\x01\x00\x05\x00\x02\x00\x02\xC0"
	hbamSequence  = "HBAM7"

	magicSize        = len(magicSequence)
	hbamSize         = len(hbamSequence)
	sectorSize       = 4096
	sectorHeaderSize = 20
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

	ctx := &FmpFile{Stream: stream, Dictionary: &FmpDict{}}
	if err := ctx.readHeader(); err != nil {
		return nil, err
	}

	ctx.FileSize = uint(info.Size())
	ctx.NumSectors = ctx.FileSize / sectorSize
	ctx.Sectors = make([]*FmpSector, ctx.NumSectors)

	currentPath := make([]uint64, 0)

	for i := uint(0); i < ctx.NumSectors; i++ {
		sector, err := ctx.readSector()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ctx.Sectors[i] = sector
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
		}
	}

	return ctx, nil
}

func (ctx *FmpFile) readHeader() error {
	buf := make([]byte, sectorSize)
	_, err := ctx.Stream.Read(buf)
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
	buf := make([]byte, sectorHeaderSize)
	n, err := ctx.Stream.Read(buf)

	if n == 0 {
		return nil, io.EOF
	}
	if err != nil {
		return nil, ErrRead
	}

	sector := &FmpSector{
		Deleted:      buf[0] != 0,
		Level:        uint8(buf[1]),
		PrevSectorID: parseVarUint64(buf[2:6]),
		NextSectorID: parseVarUint64(buf[6:10]),
	}

	payload := make([]byte, sectorSize-sectorHeaderSize)
	n, err = ctx.Stream.Read(payload)
	if n != sectorSize-sectorHeaderSize {
		return nil, ErrRead
	}
	if err != nil {
		return nil, ErrRead
	}
	sector.Chunks = make([]*FmpChunk, 0)

	for {
		chunk, err := ctx.readChunk(payload)
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
