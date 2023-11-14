package scanner

import (
	"context"
	"fmt"
	"strings"
)

type Scanner interface {
	ReadInput(ctx context.Context) (string, error)
}

func New() *scanner {
	return &scanner{}
}

type scanner struct{}

func (s *scanner) ReadInput(ctx context.Context) (string, error) {
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return "", err
	}
	return strings.Trim(input, " "), nil
}
