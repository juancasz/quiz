package main

import (
	"context"
	"flag"
	"log"
	"quiz/internal/quiz"
	"quiz/internal/reader"
	"quiz/internal/scanner"
	"quiz/internal/timer"
	"strconv"
)

func main() {
	path := flag.String("path", "", "path csv quiz file")
	seconds := flag.String("seconds", "30", "seconds to finish quiz")
	flag.Parse()

	read, err := reader.New(*path)
	if err != nil {
		log.Fatal(err)
	}
	scan := scanner.New()
	secondsQuiz, err := strconv.Atoi(*seconds)
	if err != nil {
		log.Fatal(err)
	}
	time := timer.New(secondsQuiz)
	quizManager := quiz.New(read, scan, time)

	if err = quizManager.ReadQuiz(context.Background()); err != nil {
		log.Fatal(err)
	}
}
