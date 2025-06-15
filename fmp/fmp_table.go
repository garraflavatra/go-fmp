package fmp

type FmpTable struct {
	ID      uint64
	Name    string
	Columns []FmpColumn
}

type FmpColumn struct {
	ID          uint64
	Name        string
	Type        FmpFieldType
	DataType    FmpDataType
	StorageType FmpFieldStorageType
	AutoEnter   FmpAutoEnterOption
	Indexed     bool
}

func (ctx *FmpFile) Tables() []*FmpTable {
	tables := make([]*FmpTable, 0)
	ent := ctx.Dictionary.GetEntry(3, 16, 5)

	for path, tableEnt := range *ent.Children {
		if path < 128 {
			continue
		}

		table := &FmpTable{
			ID:      path,
			Name:    decodeByteSeq(tableEnt.Children.GetValue(16)),
			Columns: make([]FmpColumn, 0),
		}

		tables = append(tables, table)
		colEnt := ctx.Dictionary.GetEntry(table.ID, 3, 5)

		for colPath, colEnt := range *colEnt.Children {
			name := decodeByteSeq(colEnt.Children.GetValue(16))
			flags := colEnt.Children.GetValue(2)

			column := FmpColumn{
				ID:          colPath,
				Name:        name,
				Type:        FmpFieldType(flags[0]),
				DataType:    FmpDataType(flags[1]),
				StorageType: FmpFieldStorageType(flags[9]),
				Indexed:     flags[8] == 128,
			}

			if flags[11] == 1 {
				column.AutoEnter = autoEnterPresetMap[flags[4]]
			} else {
				column.AutoEnter = autoEnterOptionMap[flags[11]]
			}

			table.Columns = append(table.Columns, column)
		}
	}

	return tables
}
