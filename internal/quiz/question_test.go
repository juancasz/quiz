package quiz

import (
	"context"
	"reflect"
	"testing"
)

func TestNewQuestion(t *testing.T) {
	tests := map[string]struct {
		line     []string
		question *question
		err      error
	}{
		"nil line":      {line: nil, question: nil, err: ErrInvalidNumberOfColumns},
		"one column":    {line: []string{"a"}, question: nil, err: ErrInvalidNumberOfColumns},
		"two columns":   {line: []string{"a", "b"}, question: &question{question: "a", answer: "b"}, err: nil},
		"three columns": {line: []string{"a", "b", "c"}, question: nil, err: ErrInvalidNumberOfColumns},
	}
	ctx := context.Background()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			question, err := newQuestion(ctx, tc.line)
			if !reflect.DeepEqual(question, tc.question) {
				t.Fatalf("expected: %v, got: %v", tc.question, question)
			}
			if err != tc.err {
				t.Fatalf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}

func TestIsAnswerCorrect(t *testing.T) {
	tests := map[string]struct {
		question   *question
		userAnswer string
		isCorrect  bool
	}{
		"correct answer":                    {question: &question{question: "1+1", answer: "2"}, userAnswer: "2", isCorrect: true},
		"correct answer case insensitivity": {question: &question{question: "are dogs animals?", answer: "yes"}, userAnswer: "YeS", isCorrect: true},
		"wrong answer":                      {question: &question{question: "1+1", answer: "2"}, userAnswer: "3", isCorrect: false},
	}
	ctx := context.Background()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			isCorrect := tc.question.IsAnswerCorrect(ctx, tc.userAnswer)
			if isCorrect != tc.isCorrect {
				t.Fatalf("expected: %v, got: %v", tc.isCorrect, isCorrect)
			}
		})
	}
}
