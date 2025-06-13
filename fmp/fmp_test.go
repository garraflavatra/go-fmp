package fmp

import "testing"

func TestOpenFile(t *testing.T) {
	f, err := OpenFile("../files/Untitled.fmp12")
	if err != nil {
		t.Fatal(err)
	}
	if f.FileSize != 393216 {
		t.Errorf("expected file size to be 393216, got %d", f.FileSize)
	}
	if f.numSectors != 95 {
		t.Errorf("expected number of sectors to be 95, got %d", f.numSectors)
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

	if len(tables) != 1 || tables[0].Name != "Untitled" {
		tablesString := ""
		for i, table := range tables {
			tablesString += table.Name
			if i < len(tables)-1 {
				tablesString += ", "
			}
		}
		t.Errorf("expected tables to be 'Untitled', got '%s'", tablesString)
	}
}
