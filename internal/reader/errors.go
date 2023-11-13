package reader

import (
	"fmt"
	"io"
)

var (
	ErrEndOfFile = fmt.Errorf("no more lines %w", io.EOF)
)
