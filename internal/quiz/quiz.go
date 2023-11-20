package quiz

import (
	"context"
	"fmt"
	"math/rand"
	"quiz/internal/reader"
	"quiz/internal/scanner"
	"quiz/internal/timer"
)

type quiz struct {
	reader.Reader
	scanner.Scanner
	timer.Timer
	questions   []question
	shuffleFlag bool
	*counter
	chanErr           chan error
	chanCorrectAnswer chan struct{}
}

type counter struct {
	correctAnswers int
	totalQuestions int
}

func New(reader reader.Reader, scanner scanner.Scanner, timer timer.Timer, shuffle bool) *quiz {
	return &quiz{
		Reader:            reader,
		Scanner:           scanner,
		Timer:             timer,
		shuffleFlag:       shuffle,
		counter:           &counter{},
		chanErr:           make(chan error),
		chanCorrectAnswer: make(chan struct{}, 1),
	}
}

func (q *quiz) Read(ctx context.Context) error {
	if err := q.firstRead(ctx); err != nil {
		return err
	}
	if q.shuffleFlag {
		q.shuffle(ctx)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go q.solveQuiz(ctx)
	if err := q.wait(ctx); err != nil {
		return err
	}
	q.results(ctx)
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

func (q *quiz) shuffle(ctx context.Context) {
	for i := len(q.questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		q.questions[i], q.questions[j] = q.questions[j], q.questions[i]
	}
}

func (q *quiz) solveQuiz(ctx context.Context) {
	for _, question := range q.questions {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("question: ", question.question)
			fmt.Printf("enter your answer and press enter: ")

			userAnswer, err := q.Scanner.ReadInput(ctx)
			if err != nil {
				q.chanErr <- err
			}

			if question.IsAnswerCorrect(ctx, userAnswer) {
				q.chanCorrectAnswer <- struct{}{}
			}

			fmt.Printf("\n\n")
		}
	}
	close(q.chanErr)
	close(q.chanCorrectAnswer)
}

func (q *quiz) wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-q.chanCorrectAnswer:
			if ok {
				q.counter.correctAnswers += 1
			} else {
				q.chanCorrectAnswer = nil
			}
		case err, ok := <-q.chanErr:
			if ok {
				return err
			} else {
				q.chanErr = nil
			}
		case <-q.Timer.Tick(ctx):
			fmt.Printf("\n\ntime completed\n\n")
			return nil
		}

		if q.chanCorrectAnswer == nil && q.chanErr == nil {
			break
		}
	}
	return nil
}

func (q *quiz) results(ctx context.Context) {
	fmt.Println("total questions --->", q.totalQuestions)
	fmt.Println("correct answers --->", q.correctAnswers)
}
