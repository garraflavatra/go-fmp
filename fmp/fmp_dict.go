package fmp

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
