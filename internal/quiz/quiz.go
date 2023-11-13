package quiz

import (
	"context"
	"fmt"
	"quiz/internal/reader"
	"quiz/internal/scanner"
)

type quiz struct {
	reader.Reader
	scanner.Scanner
	*counter
}

type counter struct {
	correctAnswers int
	totalQuestions int
}

func New(reader reader.Reader, scanner scanner.Scanner) *quiz {
	return &quiz{
		Reader:  reader,
		Scanner: scanner,
		counter: &counter{},
	}
}

func (q *quiz) ReadQuiz(ctx context.Context) error {
	for {
		line, err := q.Reader.ReadLine(ctx)
		if err == reader.ErrEndOfFile {
			break
		}
		if err != nil {
			return err
		}

		question, err := newQuestion(ctx, line)
		if err != nil {
			return err
		}
		q.counter.totalQuestions += 1

		fmt.Println("question: ", question.question)
		fmt.Printf("enter your answer and press enter: ")

		userAnswer, err := q.Scanner.ReadInput(ctx)
		if err != nil {
			return err
		}

		if question.IsAnswerCorrect(ctx, userAnswer) {
			q.counter.correctAnswers += 1
		}

		fmt.Printf("\n\n")
	}
	fmt.Println("total questions --->", q.totalQuestions)
	fmt.Println("correct answers --->", q.correctAnswers)
	return nil
}
