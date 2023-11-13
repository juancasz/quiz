package reader

import (
	"context"
	"encoding/csv"
	"io"
	"os"
)

type Reader interface {
	ReadLine(ctx context.Context) ([]string, error)
}

type reader struct {
	*csv.Reader
}

func New(path string) (*reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &reader{
		Reader: csv.NewReader(file),
	}, nil
}

func (r *reader) ReadLine(ctx context.Context) ([]string, error) {
	record, err := r.Reader.Read()
	if err == io.EOF {
		return nil, ErrEndOfFile
	}
	return record, nil
}
