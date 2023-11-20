package reader_test

import (
	"context"
	"os"
	"quiz/internal/reader"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	// Create a mock CSV file
	content := `header1,header2`

	dir := t.TempDir()
	tmpfile, err := os.CreateTemp(dir, "example.*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test ReadLine function
	readerTest, err := reader.New(tmpfile.Name())
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}
	ctx := context.Background()
	record, err := readerTest.ReadLine(ctx)
	if err != nil {
		t.Errorf("ReadLine() returned an error: %v", err)
	}
	if !strings.EqualFold(record[0], "header1") || !strings.EqualFold(record[1], "header2") {
		t.Errorf("ReadLine() returned unexpected record: %v", record)
	}

	// Test ReadLine function when the end of the file is reached
	record, err = readerTest.ReadLine(ctx)
	if err != reader.ErrEndOfFile {
		t.Errorf("ReadLine() returned unexpected error: %v", err)
	}
}
