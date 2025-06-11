package fmp

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"time"
)

const (
	magicSequence = "\x00\x01\x00\x00\x00\x02\x00\x01\x00\x05\x00\x02\x00\x02\xC0"
	hbamSequence  = "HBAM7"

	magicSize  = len(magicSequence)
	hbamSize   = len(hbamSequence)
	sectorSize = 4096
)

type FmpFile struct {
	Stream io.ReadSeeker

	FileSize   uint
	NumSectors uint

	VersionDate     time.Time
	ApplicationName string
}

type FmpSector struct {
	Deleted      bool
	Level        uint8
	PrevSectorID uint32
	NextSectorID uint32
	Payload      []byte
}

func OpenFile(path string) (*FmpFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	stream, err := os.Open(path)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return nil, err
	}
	ctx := &FmpFile{Stream: stream}
	if err := ctx.readHeader(); err != nil {
		stream.Close()
		return nil, err
	}
	ctx.FileSize = uint(info.Size())
	ctx.NumSectors = ctx.FileSize / sectorSize
	return ctx, nil
}

func (ctx *FmpFile) readHeader() error {
	buf := make([]byte, sectorSize)
	_, err := ctx.Stream.Read(buf)
	if err != nil {
		return ErrRead
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
	ctx.ApplicationName = string(buf[542 : 542+appNameLength])

	return nil
}

func (ctx *FmpFile) readSector() (*FmpSector, error) {
	buf := make([]byte, sectorSize)
	_, err := ctx.Stream.Read(buf)
	if err != nil {
		return nil, ErrRead
	}
	sector := &FmpSector{
		Deleted:      buf[0] != 0,
		Level:        uint8(buf[1]),
		PrevSectorID: binary.BigEndian.Uint32(buf[2:6]),
		NextSectorID: binary.BigEndian.Uint32(buf[6:10]),
		Payload:      buf[6:4076],
	}
	return sector, nil
}
