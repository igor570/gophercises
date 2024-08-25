package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type QuestionAnswer struct {
	Question string
	Answer   string
}

type UserAnswers struct {
	correct   int
	incorrect int
}

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Channel for timer
	done := make(chan bool)
	// Channel for userInput
	answerChan := make(chan string)

	go timer(done)

	var questionsAndAnswers []QuestionAnswer
	var userAnswers UserAnswers

	for _, record := range records {
		questionsAndAnswers = append(questionsAndAnswers, QuestionAnswer{
			Question: record[0],
			Answer:   record[1],
		})

		question := questionsAndAnswers[0].Question
		answer := questionsAndAnswers[0].Answer

		fmt.Println("Answer this: ", question)

		// Start a goroutine to read user input
		go func() {
			var userInput string
			fmt.Scan(&userInput)
			answerChan <- userInput
		}()

		select {
		case <-done:
			fmt.Println("\nTime's up!")
			fmt.Printf("Correct Answers: %d, Incorrect Answers: %d.\n", userAnswers.correct, userAnswers.incorrect)
			return
		case userInput := <-answerChan:
			if userInput == answer {
				userAnswers.correct += 1
			} else {
				userAnswers.incorrect += 1
			}
		}
		questionsAndAnswers = []QuestionAnswer{} // reset slice
	}

	fmt.Printf("Correct Answers: %d, Incorrect Answers: %d.\n", userAnswers.correct, userAnswers.incorrect)
}

func timer(done chan bool) {
	time.Sleep(time.Second * 30)
	done <- true
}

func openCSVFile(fileName string) *os.File {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)

	}

	fmt.Println("Correctly opened CSV file")

	return file
}

func readCSVFile(file *os.File) [][]string {
	fileReader := csv.NewReader(file)
	questions, err := fileReader.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	return questions
}
