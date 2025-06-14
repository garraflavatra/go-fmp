package fmp

import (
	"io"
	"time"
)

type FmpFile struct {
	VersionDate time.Time
	CreatorName string
	FileSize    uint
	Sectors     []*FmpSector
	Chunks      []*FmpChunk
	Dictionary  *FmpDict

	numSectors      uint64 // Excludes the header sector
	currentSectorID uint64

	stream io.ReadSeeker
}

type FmpSector struct {
	ID      uint64
	Level   uint8
	Deleted bool
	PrevID  uint64
	NextID  uint64
	Prev    *FmpSector
	Next    *FmpSector
	Payload []byte
	Chunks  []*FmpChunk
}

type FmpChunk struct {
	Type   FmpChunkType
	Length uint64
	Key    uint64 // If Type == FMP_CHUNK_SHORT_KEY_VALUE or FMP_CHUNK_LONG_KEY_VALUE
	Index  uint64 // Segment index, if Type == FMP_CHUNK_SEGMENTED_DATA
	Value  []byte
}

type FmpDict map[uint64]*FmpDictEntry

type FmpDictEntry struct {
	Value    []byte
	Children *FmpDict
}

type FmpTable struct {
	ID   uint64
	Name string
}

func (dict *FmpDict) GetEntry(path []uint64) *FmpDictEntry {
	for i, key := range path {
		_, ok := (*dict)[key]
		if !ok {
			return nil
		}

		if i == len(path)-1 {
			return (*dict)[key]
		} else {
			dict = (*dict)[key].Children
			if dict == nil {
				return nil
			}
		}
	}
	return nil
}

func (dict *FmpDict) GetValue(path []uint64) []byte {
	ent := dict.GetEntry(path)
	if ent != nil {
		return ent.Value
	}
	return nil
}

func (dict *FmpDict) SetValue(path []uint64, value []byte) {
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
