package fmp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

type FmpFile struct {
	VersionDate time.Time
	CreatorName string
	FileSize    uint
	NumSectors  uint
	Stream      io.ReadSeeker
	Sectors     []*FmpSector
}

type FmpSector struct {
	Deleted      bool
	Level        uint8
	PrevSectorID uint32
	NextSectorID uint32
	Chunks       []*FmpChunk
}

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

	ctx := &FmpFile{Stream: stream}
	if err := ctx.readHeader(); err != nil {
		return nil, err
	}

	ctx.FileSize = uint(info.Size())
	ctx.NumSectors = ctx.FileSize / sectorSize
	ctx.Sectors = make([]*FmpSector, ctx.NumSectors)

	for i := uint(0); i < ctx.NumSectors; i++ {
		println("reading sector ", i)
		sector, err := ctx.readSector()
		if err != nil {
			return nil, err
		}
		ctx.Sectors[i] = sector
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
		PrevSectorID: binary.BigEndian.Uint32(buf[2:6]),
		NextSectorID: binary.BigEndian.Uint32(buf[6:10]),
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
		println(hex.EncodeToString(payload[0:2]))
		chunk, err := ctx.readChunk(payload)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		sector.Chunks = append(sector.Chunks, chunk)
		if chunk == nil {
			break
		}
		if chunk.Length == 0 {
			panic("chunk length not set")
		}
		payload = payload[min(chunk.Length, uint32(len(payload))):]
		if len(payload) == 0 {
			break
		}
	}

	return sector, nil
}
