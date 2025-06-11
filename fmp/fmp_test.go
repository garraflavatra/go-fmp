package fmp

import "testing"

func TestOpenFile(t *testing.T) {
	f, err := OpenFile("../files/Untitled.fmp12")
	if err != nil {
		t.Fatal(err)
	}
	if f.ApplicationName != "Pro 12.0" {
		t.Errorf("expected application name to be 'Pro 12.0', got '%s'", f.ApplicationName)
	}
	if f.VersionDate.Format("2006-01-02") != "2025-01-11" {
		t.Errorf("expected version date to be '2025-01-11', got '%s'", f.VersionDate.Format("2006-01-02"))
	}
}
