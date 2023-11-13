package main

import (
	"context"
	"flag"
	"log"
	"quiz/internal/quiz"
	"quiz/internal/reader"
	"quiz/internal/scanner"
)

func main() {
	path := flag.String("path", "", "path csv quiz file")
	flag.Parse()

	read, err := reader.New(*path)
	if err != nil {
		log.Fatal(err)
	}

	scan := scanner.New()

	quizManager := quiz.New(read, scan)

	if err = quizManager.ReadQuiz(context.Background()); err != nil {
		log.Fatal(err)
	}
}
