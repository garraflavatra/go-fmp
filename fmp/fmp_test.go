package fmp

import (
	"slices"
	"testing"
)

func TestOpenFile(t *testing.T) {
	f, err := OpenFile("../files/Untitled.fmp12")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if f.FileSize != 229376 {
		t.Errorf("expected file size to be 229376, got %d", f.FileSize)
	}
	if f.numSectors != 55 {
		t.Errorf("expected number of sectors to be 55, got %d", f.numSectors)
	}
	if f.CreatorName != "Pro 12.0" {
		t.Errorf("expected application name to be 'Pro 12.0', got '%s'", f.CreatorName)
	}
	if f.VersionDate.Format("2006-01-02") != "2025-01-11" {
		t.Errorf("expected version date to be '2025-01-11', got '%s'", f.VersionDate.Format("2006-01-02"))
	}

	f.ToDebugFile("../private/output")
}

func TestTables(t *testing.T) {
	f, err := OpenFile("../files/Untitled.fmp12")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	expectedNames := []string{"Untitled"}
	tableNames := []string{}
	for _, table := range f.tables {
		tableNames = append(tableNames, table.Name)
	}
	if !slicesHaveSameElements(tableNames, expectedNames) {
		t.Errorf("tables do not match")
	}

	table := f.Table("Untitled")
	if table == nil {
		t.Errorf("expected table to exist, but it does not")
		return
	}
	if table.Name != "Untitled" {
		t.Errorf("expected table name to be 'Untitled', but it is '%s'", table.Name)
	}
	if len(table.Records) != 3 {
		t.Errorf("expected table to have 3 records, but it has %d", len(table.Records))
	}
	if table.Records[1].Values[1] != "629FAA83-50D8-401F-A560-C8D45217D17B" {
		t.Errorf("first record has an incorrect ID '%s'", table.Records[0].Values[0])
	}

	col := table.Column("PrimaryKey")
	if col == nil {
		t.Errorf("expected column to exist, but it does not")
		return
	}
	if col.Name != "PrimaryKey" {
		t.Errorf("expected column name to be 'PrimaryKey', but it is '%s'", col.Name)
	}
	if col.Type != FmpFieldSimple {
		t.Errorf("expected field type to be simple, but it is not")
	}
	if col.DataType != FmpDataText {
		t.Errorf("expected field data type to be text, but it is not")
	}
	if col.StorageType != FmpFieldStorageRegular {
		t.Errorf("expected field storage type to be regular, but it is not")
	}
	if col.Repetitions != 1 {
		t.Errorf("expected field repetition count to be 1, but it is %d", col.Repetitions)
	}
	if !col.Indexed {
		t.Errorf("expected field to be indexed, but it is not")
	}
	if col.AutoEnter != FmpAutoEnterCalculationReplacingExistingValue {
		t.Errorf("expected field to have auto enter calculation replacing existing value, but it does not")
	}

	newRecord, err := table.NewRecord(map[string]string{"PrimaryKey": "629FAA83-50D8-401F-A560-C8D45217D17B"})
	if newRecord == nil || err != nil {
		t.Errorf("expected new record to be created, but it is nil")
		return
	}
	if newRecord.Index != 4 {
		t.Errorf("expected new record index to be 4, but it is %d", newRecord.Index)
	}
	if newRecord.Value("PrimaryKey") != "629FAA83-50D8-401F-A560-C8D45217D17B" {
		t.Errorf("expected new record primary key to be '629FAA83-50D8-401F-A560-C8D45217D17B', but it is '%s'", newRecord.Value("PrimaryKey"))
	}
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
