package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type question struct {
	problem string
	result  string
}

type quiz struct {
	questions []question
	score     int
}

func (q *quiz) ask() {
	for _, question := range q.questions {
		fmt.Println(question.problem)
		var answer string
		fmt.Scanln(&answer)
		if answer == question.result {
			q.score++
		}
	}
}

func parseCSV(p string) []question {
	file, err := os.Open(p)
	if err != nil {
		log.Fatalf("open file failed: %v\n", err)
	}
	defer file.Close()

	csvRead := csv.NewReader(file)
	records, err := csvRead.ReadAll()
	if err != nil {
		log.Fatalf("csv parse failed: %v\n", err)
	}

	var questions []question
	for _, record := range records {
		q := question{problem: record[0], result: record[1]}
		questions = append(questions, q)
	}
	return questions
}

var (
	CSV     string
	timeout int
)

func init() {
	flag.StringVar(&CSV, "q", "", "A .csv file with questions and answers.")
	flag.IntVar(&timeout, "t", 30, " quiz time limit.")
	flag.Parse()
}

func main() {
	if CSV == "" {
		log.Fatalln("a .csv file must be provided through the -q flag")
	}

	questions := parseCSV(CSV)
	qz := quiz{questions: questions}
	timeoutCh := time.After(time.Duration(timeout) * time.Second)
	resultCh := make(chan quiz)
	go func() {
		qz.ask()
		resultCh <- qz
	}()

	select {
	case <-resultCh:
		fmt.Println("Test finished!")
	case <-timeoutCh:
		fmt.Println("Timeout!")
	}
	fmt.Printf("Your score: %d/%d", qz.score, len(qz.questions))
}
