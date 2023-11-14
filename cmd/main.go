package main

import (
	"context"
	"flag"
	"log"
	"quiz/internal/quiz"
	"quiz/internal/reader"
	"quiz/internal/scanner"
	"quiz/internal/timer"
)

func main() {
	path := flag.String("path", "", "path csv quiz file")
	seconds := flag.Int("seconds", 30, "seconds to finish quiz")
	shuffle := flag.Bool("shuffle", false, "shuffle the questions randomly")
	flag.Parse()

	read, err := reader.New(*path)
	if err != nil {
		log.Fatal(err)
	}
	scan := scanner.New()
	time := timer.New(*seconds)
	quizManager := quiz.New(read, scan, time, *shuffle)

	if err = quizManager.Read(context.Background()); err != nil {
		log.Fatal(err)
	}
}
