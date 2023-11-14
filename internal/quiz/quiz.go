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
	questions []question
	*counter
	chanErr           chan error
	chanCorrectAnswer chan struct{}
}

type counter struct {
	correctAnswers int
	totalQuestions int
}

func New(reader reader.Reader, scanner scanner.Scanner, timer timer.Timer) *quiz {
	return &quiz{
		Reader:            reader,
		Scanner:           scanner,
		Timer:             timer,
		counter:           &counter{},
		chanErr:           make(chan error),
		chanCorrectAnswer: make(chan struct{}),
	}
}

func (q *quiz) Read(ctx context.Context) error {
	if err := q.firstRead(ctx); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func(ctx context.Context) {
		if err := q.solveQuiz(ctx); err != nil {
			q.chanErr <- err
		}
	}(ctx)
	if err := q.wait(ctx); err != nil {
		return err
	}
	return nil
}

func (q *quiz) firstRead(ctx context.Context) error {
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
		q.questions = append(q.questions, *question)
		q.counter.totalQuestions += 1
	}
	return nil
}

func (q *quiz) solveQuiz(ctx context.Context) error {
	for _, question := range q.questions {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Println("question: ", question.question)
			fmt.Printf("enter your answer and press enter: ")

			userAnswer, err := q.Scanner.ReadInput(ctx)
			if err != nil {
				return err
			}

			if question.IsAnswerCorrect(ctx, userAnswer) {
				q.chanCorrectAnswer <- struct{}{}
			}

			fmt.Printf("\n\n")
		}
	}
	q.Timer.Finish(ctx)
	return nil
}

func (q *quiz) wait(ctx context.Context) error {
	for {
		select {
		case <-q.chanCorrectAnswer:
			q.counter.correctAnswers += 1
		case err := <-q.chanErr:
			return err
		case <-q.Done(ctx):
			fmt.Printf("\n\nall answers sent\n\n")
			q.results(ctx)
			return nil
		case <-q.Timer.Tick(ctx):
			fmt.Printf("\n\ntime completed\n\n")
			q.results(ctx)
			return nil
		}
	}
}

func (q *quiz) results(ctx context.Context) {
	fmt.Println("total questions --->", q.totalQuestions)
	fmt.Println("correct answers --->", q.correctAnswers)
}
