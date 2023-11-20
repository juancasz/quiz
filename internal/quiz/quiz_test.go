package quiz

import (
	"context"
	"errors"
	"quiz/internal/reader"
	"reflect"
	"testing"
)

type mockReader struct {
	Lines [][]string
	Index int
}

func (m *mockReader) ReadLine(ctx context.Context) ([]string, error) {
	if m.Index >= len(m.Lines) {
		return nil, reader.ErrEndOfFile
	}
	line := m.Lines[m.Index]
	m.Index++
	return line, nil
}

func TestFirstRead(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockReader := &mockReader{
		Lines: [][]string{
			{"Question 1", "Answer 1"},
			{"Question 2", "Answer 2"},
			{"Question 3", "Answer 3"},
		},
	}
	q := New(mockReader, nil, nil, false)

	// Act
	err := q.firstRead(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(q.questions) != len(mockReader.Lines) {
		t.Fatalf("Expected %d questions, but got %d", len(mockReader.Lines), len(q.questions))
	}
	for i, question := range q.questions {
		if question.question != mockReader.Lines[i][0] {
			t.Errorf("Expected question %d to be '%s', but was '%s'", i, mockReader.Lines[i][0], question.question)
		}
		if question.answer != mockReader.Lines[i][1] {
			t.Errorf("Expected answer %d to be '%s', but was '%s'", i, mockReader.Lines[i][1], question.answer)
		}
	}
}

func TestShuffle(t *testing.T) {
	// Arrange
	ctx := context.Background()
	q := &quiz{
		questions: []question{
			{"Question 1", "Answer 1"},
			{"Question 2", "Answer 2"},
			{"Question 3", "Answer 3"},
		},
	}
	originalQuestions := make([]question, len(q.questions))
	copy(originalQuestions, q.questions)

	// Act
	q.shuffle(ctx)

	// Assert
	if len(q.questions) != len(originalQuestions) {
		t.Fatalf("Expected %d questions, but got %d", len(originalQuestions), len(q.questions))
	}
	if reflect.DeepEqual(q.questions, originalQuestions) {
		t.Fatalf("Expected questions to be shuffled, but they were not")
	}
}

type mockScanner struct {
	Inputs []string
	Index  int
}

func (m *mockScanner) ReadInput(ctx context.Context) (string, error) {
	if m.Index >= len(m.Inputs) {
		return "", errors.New("no more inputs")
	}
	input := m.Inputs[m.Index]
	m.Index++
	return input, nil
}

func TestSolveQuiz(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockScanner := &mockScanner{
		Inputs: []string{"Answer 1", "Answer 2", "Answer 3"},
	}
	q := New(nil, mockScanner, nil, false)
	q.questions = []question{
		{"Question 1", "Answer 1"},
		{"Question 2", "Answer 2"},
		{"Question 3", "Answer 3"},
	}

	// Start a goroutine to receive from the chanCorrectAnswer channel
	correctAnswers := 0
	go func() {
		for {
			<-q.chanCorrectAnswer
			correctAnswers++
		}
	}()

	// Act
	err := q.solveQuiz(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(mockScanner.Inputs) != correctAnswers {
		t.Fatalf("Expected %d correct answers, but got %d", len(mockScanner.Inputs), len(q.chanCorrectAnswer))
	}
}
