package quiz

import (
	"context"
	"strings"
)

type question struct {
	question string
	answer   string
}

func newQuestion(ctx context.Context, line []string) (*question, error) {
	if err := validateLine(ctx, line); err != nil {
		return nil, err
	}
	return &question{
		question: line[0],
		answer:   line[1],
	}, nil
}

func validateLine(ctx context.Context, line []string) error {
	if len(line) != 2 {
		return ErrInvalidNumberOfColumns
	}
	return nil
}

func (q *question) IsAnswerCorrect(ctx context.Context, userAnswer string) bool {
	return strings.EqualFold(q.answer, userAnswer)
}
