package fmp

type FmpTable struct {
	ID      uint64
	Name    string
	Columns map[uint64]FmpColumn
	Records map[uint64]FmpRecord
}

type FmpColumn struct {
	Index       uint64
	Name        string
	Type        FmpFieldType
	DataType    FmpDataType
	StorageType FmpFieldStorageType
	AutoEnter   FmpAutoEnterOption
	Repetitions uint8
	Indexed     bool
}

type FmpRecord struct {
	Index  uint64
	Values map[uint64]string
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
			Name:    decodeFmpString(tableEnt.Children.GetValue(16)),
			Columns: map[uint64]FmpColumn{},
			Records: map[uint64]FmpRecord{},
		}

		tables = append(tables, table)

		for colPath, colEnt := range *ctx.Dictionary.GetChildren(table.ID, 3, 5) {
			name := decodeFmpString(colEnt.Children.GetValue(16))
			flags := colEnt.Children.GetValue(2)

			column := FmpColumn{
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
			record := FmpRecord{Index: recPath, Values: make(map[uint64]string)}
			table.Records[record.Index] = record

			for colIndex, value := range *recEnt.Children {
				record.Values[colIndex] = decodeFmpString(value.Value)
			}
		}
	}

	return tables
}
