package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type questionPair struct {
	Question string
	Answer   string
}

// parseCSVFile will parse a valid csv file with the following format
// "question","answer"
func parseCSVFile(filename string) ([]questionPair, error) {
	var questionPairs []questionPair

	csvFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(strings.NewReader(string(csvFile)))

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			questionPairs = append(questionPairs, questionPair{Question: record[0], Answer: record[1]})
		}
	}

	return questionPairs, nil
}

// parseJSONFile will parse a valid JSON file with the following format:
// [
//  {"question":"Meaning of life"?, "answer": "42"},
//  ...
// ]
func parseJSONFile(filename string) ([]questionPair, error) {

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var questionPairs []questionPair
	err = json.Unmarshal(fileContents, &questionPairs)
	if err != nil {
		return nil, err
	}

	return questionPairs, nil
}

// Accepted sources are:
// CSV file  - see parseCSVFile() for file format
// JSON file - see parseJSONFile() for file format
func parseQuestionSource(filename string) ([]questionPair, error) {
	switch {
	case strings.HasSuffix(filename, ".csv"):
		return parseCSVFile(filename)
	case strings.HasSuffix(filename, ".json"):
		return parseJSONFile(filename)
	default:
		return nil, errors.New("Invalid input type, should be .csv or .json")
	}
}

func main() {
	numQuestionsAsked := 0
	numCorrectAnswers := 0

	questionSource := flag.String("f", "problems.csv", "csv/json file to read")
	duration := flag.Int("d", 30, "time in seconds to run quiz for")
	shuffle := flag.Bool("s", false, "shuffle the deck")

	flag.Parse()

	durationStr := fmt.Sprintf("%ds", *duration)
	timerDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatal(err)
	}

	if questionPairs, err := parseQuestionSource(*questionSource); err != nil {
		log.Fatal(err)
	} else {

		if *shuffle { // randomise questionPairs
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(questionPairs), func(i, j int) {
				questionPairs[i], questionPairs[j] = questionPairs[j], questionPairs[i]
			})
		}

		fmt.Printf("\n\nYou have %s to answer all questions - press enter to begin\n\n", durationStr)
		fmt.Scanln()

		timer := time.NewTimer(timerDuration)
		go func() { // seperate threaded "goroutine" function that sits and waits for timer channel to fire
			<-timer.C

			fmt.Println("\n\nTime's up!!!")
			outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers)
		}()

		var userInput = ""
		for _, question := range questionPairs {
			fmt.Printf("%s=", question.Question)

			if _, err := fmt.Scanln(&userInput); err == nil { // ignore err, blank input
				if question.Answer == strings.ToLower(strings.TrimSpace(userInput)) {
					numCorrectAnswers += 1
				}
			}
			numQuestionsAsked++
		}

		fmt.Println("\n\nWell done - you answered all questions in the allowed time!!!")
		outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers)
	}
}

func outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers int) {
	scoreAsPercentage := math.Floor(float64(numCorrectAnswers) / float64(numQuestionsAsked) * 100)
	fmt.Printf("\n\nYour score was: %d/%d (%.0f%%)\n\n\n", numCorrectAnswers, numQuestionsAsked, scoreAsPercentage)
	os.Exit(1)
}
