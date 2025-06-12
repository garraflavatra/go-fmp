package fmp

func (ctx *FmpFile) Tables() []*FmpTable {
	tables := make([]*FmpTable, 0)
	metaDict := ctx.Dictionary.get([]uint16{4, 1, 7})
	for _, meta := range *metaDict.Children {
		name := decodeByteSeq(meta.Children.get([]uint16{16}).Value)
		table := &FmpTable{Name: name}
		tables = append(tables, table)
	}
	return tables
}
