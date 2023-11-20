package quiz

import (
	"context"
	"errors"
	"quiz/internal/reader"
	"quiz/internal/timer"
	"reflect"
	"sync"
	"testing"
	"time"
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
	q := New(nil, nil, nil, true)
	q.questions = []question{
		{"Question 1", "Answer 1"},
		{"Question 2", "Answer 2"},
		{"Question 3", "Answer 3"},
		{"Question 4", "Answer 4"},
		{"Question 5", "Answer 5"},
		{"Question 6", "Answer 6"},
		{"Question 7", "Answer 7"},
		{"Question 8", "Answer 8"},
		{"Question 9", "Answer 9"},
		{"Question 10", "Answer 10"},
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

	// Start a goroutine to receive from the chanCorrectAnswer channel and error channel
	correctAnswers := 0
	var wg sync.WaitGroup
	var errQuiz error
	wg.Add(1)
	go func() {
		for {
			select {
			case _, ok := <-q.chanCorrectAnswer:
				if ok {
					correctAnswers += 1
				} else {
					q.chanCorrectAnswer = nil
				}
			case err, ok := <-q.chanErr:
				if ok {
					errQuiz = err
				} else {
					q.chanErr = nil
				}
			}

			if (q.chanCorrectAnswer == nil && q.chanErr == nil) || errQuiz != nil {
				break
			}
		}
		wg.Done()
	}()

	// Act
	q.solveQuiz(ctx)
	wg.Wait()

	// Assert
	if errQuiz != nil {
		t.Fatalf("Unexpected error: %v", errQuiz)
	}
	if len(mockScanner.Inputs) != correctAnswers {
		t.Fatalf("Expected %d correct answers, but got %d", len(mockScanner.Inputs), len(q.chanCorrectAnswer))
	}
}

type mockTimer struct {
	TickChan chan time.Time
}

func (m *mockTimer) Tick(ctx context.Context) timer.Tick {
	return m.TickChan
}

func TestWait_CorrectAnswer(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTimerTest := &mockTimer{
		TickChan: make(chan time.Time),
	}
	q := New(nil, nil, mockTimerTest, false)

	// Act
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = q.wait(ctx)
		wg.Done()
	}()

	// Send a value to the chanCorrectAnswer channel
	q.chanCorrectAnswer <- struct{}{}
	close(q.chanCorrectAnswer)
	close(q.chanErr)
	wg.Wait()

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if q.counter.correctAnswers != 1 {
		t.Fatalf("Expected 1 correct answer, but got %d", q.counter.correctAnswers)
	}
}

func TestWait_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTimerTest := &mockTimer{
		TickChan: make(chan time.Time),
	}
	q := New(nil, nil, mockTimerTest, false)

	// Act
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = q.wait(ctx)
		wg.Done()
	}()

	// Send an error to the chanErr channel
	q.chanErr <- errors.New("test error")
	close(q.chanCorrectAnswer)
	close(q.chanErr)
	wg.Wait()

	// Assert
	if err == nil || err.Error() != "test error" {
		t.Fatalf("Expected error 'test error', but got: %v", err)
	}
}

func TestWait_Tick(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTimerTest := &mockTimer{
		TickChan: make(chan time.Time),
	}
	q := New(nil, nil, mockTimerTest, false)

	// Act
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = q.wait(ctx)
		wg.Done()
	}()

	// Send a value to the Tick channel
	mockTimerTest.TickChan <- time.Now()
	close(q.chanCorrectAnswer)
	close(q.chanErr)
	wg.Wait()

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRead(t *testing.T) {
	mockReaderTest := &mockReader{
		Lines: [][]string{
			{"Question 1", "Answer 1"},
			{"Question 2", "Answer 2"},
			{"Question 3", "Answer 3"},
		},
	}
	mockScannerTest := &mockScanner{
		Inputs: []string{"Answer 1", "Answer 2", "Answer 3"},
	}
	mockTimerTest := &mockTimer{
		TickChan: make(chan time.Time),
	}
	q := New(mockReaderTest, mockScannerTest, mockTimerTest, true)
	err := q.Read(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
