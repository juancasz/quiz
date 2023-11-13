package quiz

import (
	"context"
	"fmt"
	"quiz/internal/reader"
	"quiz/internal/scanner"
	"quiz/internal/timer"
)

type quiz struct {
	reader.Reader
	scanner.Scanner
	timer.Timer
	*counter
}

type counter struct {
	correctAnswers int
	totalQuestions int
}

func New(reader reader.Reader, scanner scanner.Scanner, timer timer.Timer) *quiz {
	return &quiz{
		Reader:  reader,
		Scanner: scanner,
		Timer:   timer,
		counter: &counter{},
	}
}

func (q *quiz) ReadQuiz(ctx context.Context) error {
	go q.Wait(ctx)
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
	q.Timer.Finish(ctx)
	fmt.Println("total questions --->", q.totalQuestions)
	fmt.Println("correct answers --->", q.correctAnswers)
	return nil
}
