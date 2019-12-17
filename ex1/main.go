package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

var correct = 0
var csvFilename = flag.String("csv", "problem.csv", "a csv file in the format of 'question,answer'")
var timeLimit = flag.Int("limit", 10, "the time limit for the quiz in seconds")

func main() {
	Initialize(*csvFilename, *timeLimit)
}

// Initialize the quiz.
func Initialize(csvFilename string, timeLimit int) error {
	fmt.Printf("Time limit is %d seconds \n", timeLimit)
	flag.Parse()

	file, err := os.Open(csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", csvFilename))
		return err
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
		return err
	}
	problems := parseLines(lines)
	output(problems, timeLimit)
	return err
}

func output(problems []problem, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime UP!!!!!!")
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
}
