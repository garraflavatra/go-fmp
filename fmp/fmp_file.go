package fmp

import (
	"bytes"
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

type FmpFile struct {
	VersionDate time.Time
	CreatorName string
	FileSize    uint
	Sectors     []*FmpSector
	Chunks      []*FmpChunk
	Dictionary  *FmpDict

	tables          []*FmpTable
	numSectors      uint64 // Excludes the header sector
	currentSectorID uint64

	stream io.ReadSeeker
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

	ctx := &FmpFile{stream: stream, Dictionary: &FmpDict{}}
	if err := ctx.readHeader(); err != nil {
		return nil, err
	}

	ctx.FileSize = uint(info.Size())
	ctx.numSectors = uint64((ctx.FileSize / sectorSize) - 1)
	ctx.Sectors = make([]*FmpSector, 0)
	ctx.stream.Seek(2*sectorSize, io.SeekStart)

	for {
		sector, err := ctx.readSector()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		ctx.Sectors = append(ctx.Sectors, sector)

		if sector.ID != 0 {
			err = sector.processChunks(ctx.Dictionary)
			if err != nil {
				return nil, err
			}
			ctx.Chunks = append(ctx.Chunks, sector.Chunks...)
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

	ctx.readTables()
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
	debug("---------- Reading sector %d", ctx.currentSectorID)
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
		PrevID:  decodeVarUint64(buf[4 : 4+4]),
		NextID:  decodeVarUint64(buf[8 : 8+4]),
		Chunks:  make([]*FmpChunk, 0),
	}

	if ctx.currentSectorID == 0 && sector.PrevID > 0 {
		return nil, ErrBadSectorHeader
	}

	sector.Payload = make([]byte, sectorPayloadSize)
	n, err = ctx.stream.Read(sector.Payload)

	if n != sectorPayloadSize {
		return nil, ErrRead
	}
	if err != nil {
		return nil, ErrRead
	}
	return sector, nil
}

func (ctx *FmpFile) readTables() {
	tables := make([]*FmpTable, 0)
	ent := ctx.Dictionary.GetEntry(3, 16, 5)

	for path, tableEnt := range *ent.Children {
		if path < 128 {
			continue
		}

		table := &FmpTable{
			ID:      path,
			Name:    decodeString(tableEnt.Children.GetValue(16)),
			Columns: map[uint64]*FmpColumn{},
			Records: map[uint64]*FmpRecord{},
		}

		tables = append(tables, table)

		for colPath, colEnt := range *ctx.Dictionary.GetChildren(table.ID, 3, 5) {
			name := decodeString(colEnt.Children.GetValue(16))
			flags := colEnt.Children.GetValue(2)

			column := &FmpColumn{
				Index:       colPath,
				Name:        name,
				Type:        FmpFieldType(flags[0]),
				DataType:    FmpDataType(flags[1]),
				StorageType: FmpFieldStorageType(flags[9]),
				Repetitions: flags[25],
				Indexed:     flags[8] == 128,
			}

			if flags[11] == 1 {
				column.AutoEnter = autoEnterPresetMap[flags[4]]
			} else {
				column.AutoEnter = autoEnterOptionMap[flags[11]]
			}

			table.Columns[column.Index] = column
		}

		for recPath, recEnt := range *ctx.Dictionary.GetChildren(table.ID, 5) {
			record := &FmpRecord{Index: recPath, Values: make(map[uint64]string)}
			table.Records[record.Index] = record

			if recPath > table.lastRecordID {
				table.lastRecordID = recPath
			}

			for colIndex, value := range *recEnt.Children {
				record.Values[colIndex] = decodeString(value.Value)
			}
		}
	}

	ctx.tables = tables
}
