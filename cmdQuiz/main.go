package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type vocabulary struct {
	v string
}

func main() {
	// demoFromCourse()
	demoFromMyself()
}

func demoFromMyself() {
	// flag.String returns a string pointer
	csvFileName := flag.String("csv", "vocabulary.csv", "a csv file in the format of word")
	timeLimit := flag.Int("t", 20, "time limit in seconds")
	flag.Parse()

	// read file
	// dereference
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse file\n"))
	}

	// parse file
	vocabs := parseLine(lines)

	// get ready
	fmt.Printf("Ready?\n")
	for _, i := range []int{3, 2, 1} {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

	// game start
	// start time count down

	// self-defined timer
	ch := make(chan bool, 1)
	correct := make(chan int, 1)

	fmt.Printf("Go!!!\n\n")

	go timeCounter(*timeLimit, ch)
	go question(vocabs, correct)

	// loop until game over or finish
	for {
		select {
		case _, ok := <-ch:
			if ok {
				fmt.Printf("\n\nGame Over!")
			}
			return
		case e, ok := <-correct:
			if ok {
				fmt.Printf("Congratulations!!!\n You finished!!!\n\n")
				fmt.Printf("You've got %d out of %d.\n", e, len(vocabs))
			}
			return
		}
	}

}

func timeCounter(limit int, ch chan<- bool) {
	for i := 0; i < limit; i++ {
		time.Sleep(1 * time.Second)
	}
	ch <- true
}

func question(vocabs []vocabulary, correct chan<- int) {
	ret := 0
	for _, p := range vocabs {
		fmt.Printf("%s : ", p.v)
		input := bufio.NewReader(os.Stdin)
		text, _ := input.ReadString('\n')

		if strings.TrimSpace(text) == p.v {
			ret += 1
		}
	}
	correct <- ret
}

func parseLine(lines [][]string) []vocabulary {

	// if already know the length
	// declare it
	problems := make([]vocabulary, len(lines))

	for i, l := range lines {
		problems[i] = vocabulary{
			v: strings.TrimSpace(l[0]),
		}
	}
	return problems
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func demoFromCourse() {
	// flag.String returns a string pointer
	csvFileName := flag.String("csv", "vocabulary.csv", "a csv file in the format of word")
	timeLimit := flag.Int("t", 20, "time limit in seconds")
	flag.Parse()

	// read file
	// dereference
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse file\n"))
	}

	// parse file
	vocabs := parseLine(lines)

	// get ready
	fmt.Printf("Ready?\n")
	for _, i := range []int{3, 2, 1} {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Go!!!\n\n")

	// game start
	// start time count down

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for _, p := range vocabs {
		fmt.Printf("%s : ", p.v)
		ansCh := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s\n", &ans)
			ansCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(vocabs))
			return
		case ans, ok := <-ansCh:
			if ok {
				if ans == p.v {
					correct++
				}
			}
		}
	}
	fmt.Printf("\nCongratulations!!!\nYou finished!!!\n")
	fmt.Printf("You've got %d out of %d.\n", correct, len(vocabs))
}
