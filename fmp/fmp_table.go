package fmp

type FmpTable struct {
	ID      uint64
	Name    string
	Columns map[uint64]*FmpColumn
	Records map[uint64]*FmpRecord

	lastRecordID uint64
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
	Table  *FmpTable
	Index  uint64
	Values map[uint64]string
}

func (ctx *FmpFile) Table(name string) *FmpTable {
	for _, table := range ctx.tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func (t *FmpTable) Column(name string) *FmpColumn {
	for _, column := range t.Columns {
		if column.Name == name {
			return column
		}
	}
	return nil
}

func (t *FmpTable) NewRecord(values map[string]string) (*FmpRecord, error) {
	vals := make(map[uint64]string)
	for k, v := range values {
		col := t.Column(k)
		vals[col.Index] = v
	}

	id := t.lastRecordID + 1
	t.lastRecordID = id
	t.Records[id] = &FmpRecord{Table: t, Index: id, Values: vals}

	return t.Records[id], nil
}

func (r *FmpRecord) Value(name string) string {
	return r.Values[r.Table.Column(name).Index]
}
