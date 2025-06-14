package fmp

func (ctx *FmpFile) Tables() []*FmpTable {
	tables := make([]*FmpTable, 0)
	ent := ctx.Dictionary.GetEntry([]uint64{3, 16, 5})

	for path, tableEnt := range *ent.Children {
		if path < 128 {
			continue
		}
		name := decodeByteSeq(tableEnt.Children.GetValue([]uint64{16}))
		table := &FmpTable{ID: path, Name: name}
		tables = append(tables, table)
	}

	return tables
}
