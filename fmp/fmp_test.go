package fmp

import "testing"

func TestOpenFile(t *testing.T) {
	f, err := OpenFile("../files/Untitled.fmp12")
	if err != nil {
		t.Fatal(err)
	}
	if f.FileSize != 229376 {
		t.Errorf("expected file size to be 393216, got %d", f.FileSize)
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
	tables := f.Tables()

	expectedNames := []string{"Untitled"}
	tableNames := []string{}
	for _, table := range tables {
		tableNames = append(tableNames, table.Name)
	}

	if !slicesHaveSameElements(tableNames, expectedNames) {
		t.Errorf("tables do not match")
	}

	var field FmpColumn
	for _, table := range tables {
		for _, column := range table.Columns {
			if column.Name == "PrimaryKey" {
				field = column
				break
			}
		}
	}

	if field.Type != FmpFieldSimple {
		t.Errorf("expected field type to be simple, but it is not")
	}
	if field.DataType != FmpDataText {
		t.Errorf("expected field data type to be text, but it is not")
	}
	if field.StorageType != FmpFieldStorageRegular {
		t.Errorf("expected field storage type to be regular, but it is not")
	}
	if field.Repetitions != 1 {
		t.Errorf("expected field repetition count to be 1, but it is %d", field.Repetitions)
	}
	if !field.Indexed {
		t.Errorf("expected field to be indexed, but it is not")
	}
	if field.AutoEnter != FmpAutoEnterCalculationReplacingExistingValue {
		t.Errorf("expected field to have auto enter calculation replacing existing value, but it does not")
	}
	if len(tables[0].Records) != 3 {
		t.Errorf("expected table to have 3 records, but it has %d", len(tables[0].Records))
	}
	if tables[0].Records[1].Values[1] != "629FAA83-50D8-401F-A560-C8D45217D17B" {
		t.Errorf("first record has an incorrect ID '%s'", tables[0].Records[0].Values[0])
	}
}
