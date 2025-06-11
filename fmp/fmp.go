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

	magicSize  = len(magicSequence)
	hbamSize   = len(hbamSequence)
	sectorSize = 4096
)

type FmpFile struct {
	Stream          io.ReadSeeker
	VersionDate     time.Time
	ApplicationName string
}

func OpenFile(path string) (*FmpFile, error) {
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
