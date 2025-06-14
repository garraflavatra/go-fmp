package fmp

import (
	"slices"
	"testing"
)

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

	expectedNames := []string{"Untitled", "WayDomains", "WayProcesses"}
	tableNames := []string{}
	for _, table := range tables {
		tableNames = append(tableNames, table.Name)
	}

	if !slicesHaveSameElements(tableNames, expectedNames) {
		t.Errorf("tables do not match")
	}
}
