package fmp

import (
	"io"
	"time"
)

type FmpFile struct {
	VersionDate time.Time
	CreatorName string
	FileSize    uint
	NumSectors  uint
	Stream      io.ReadSeeker
	Sectors     []*FmpSector
	Chunks      []*FmpChunk
	Dictionary  *FmpDict
}

type FmpSector struct {
	Deleted      bool
	Level        uint8
	PrevSectorID uint32
	NextSectorID uint32
	Chunks       []*FmpChunk
}

type FmpChunk struct {
	Type   FmpChunkType
	Length uint32
	Key    uint32 // If Type == FMP_CHUNK_SHORT_KEY_VALUE or FMP_CHUNK_LONG_KEY_VALUE
	Index  uint32 // Segment index, if Type == FMP_CHUNK_SEGMENTED_DATA
	Value  []byte
}

type FmpDict map[uint16]*FmpDictEntry

type FmpDictEntry struct {
	Value    []byte
	Children *FmpDict
}

type FmpTable struct {
	Name string
}

func (dict *FmpDict) get(path []uint16) *FmpDictEntry {
	for i, key := range path {
		_, ok := (*dict)[key]
		if !ok {
			(*dict)[key] = &FmpDictEntry{Children: &FmpDict{}}
		}

		if i == len(path)-1 {
			return (*dict)[key]
		} else {
			dict = (*dict)[key].Children
		}
	}
	return nil
}

func (dict *FmpDict) getValue(path []uint16) []byte {
	ent := dict.get(path)
	if ent != nil {
		return ent.Value
	}
	return nil
}

func (dict *FmpDict) set(path []uint16, value []byte) {
	for i, key := range path {
		_, ok := (*dict)[key]
		if !ok {
			(*dict)[key] = &FmpDictEntry{Children: &FmpDict{}}
		}

		if i == len(path)-1 {
			(*dict)[key].Value = value
		} else {
			dict = (*dict)[key].Children
		}
	}
}
