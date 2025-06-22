package fmp

import "slices"

type FmpDict map[uint64]*FmpDictEntry

type FmpDictEntry struct {
	Value    []byte
	Children *FmpDict
}

func (dict *FmpDict) GetEntry(path ...uint64) *FmpDictEntry {
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

func (dict *FmpDict) GetValue(path ...uint64) []byte {
	ent := dict.GetEntry(path...)
	if ent != nil {
		return ent.Value
	}
	return nil
}

func (dict *FmpDict) GetChildren(path ...uint64) *FmpDict {
	ent := dict.GetEntry(path...)
	if ent != nil {
		return ent.Children
	}
	return &FmpDict{}
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

func parseVarUint64(payload []byte) uint64 {
	var length uint64
	n := min(len(payload), 8) // clamp to uint64
	for i := range n {
		length <<= 8
		length |= uint64(payload[i])
	}
	return length
}

func decodeFmpString(payload []byte) string {
	result := ""
	for i := range payload {
		result += string(payload[i] ^ 0x5A)
	}
	return result
}

func addIf(cond bool, val uint64) uint64 {
	if cond {
		return val
	}
	return 0
}

func slicesHaveSameElements[Type comparable](a, b []Type) bool {
	if len(a) != len(b) {
		return false
	}
	for _, av := range a {
		found := slices.Contains(b, av)
		if !found {
			return false
		}
	}
	return true
}
