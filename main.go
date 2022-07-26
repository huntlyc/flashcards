package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
// JSON file - see parseJSONFile() for file format
// @TODO hook into charm kv?
func parseQuestionSource(filename string) ([]questionPair, error) {
	switch {
	case strings.HasSuffix(filename, ".json"):
		return parseJSONFile(filename)
	default:
		return nil, errors.New("Invalid input type, should be .csv or .json")
	}
}

func main() {
	numQuestionsAsked := 0
	numCorrectAnswers := 0

	questionSource := flag.String("f", "cards.json", "json file to read")
	shuffle := flag.Bool("s", false, "shuffle the deck")

	flag.Parse()

	if questionPairs, err := parseQuestionSource(*questionSource); err != nil {
		log.Fatal(err)
	} else {

		if *shuffle { // randomise questionPairs
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(questionPairs), func(i, j int) {
				questionPairs[i], questionPairs[j] = questionPairs[j], questionPairs[i]
			})
		}

		var userInput = ""
		for _, question := range questionPairs {
			fmt.Printf("%s\n", question.Question)

			if _, err := fmt.Scanln(&userInput); err == nil { // ignore err, blank input

				isCorrect := "❌"
				if question.Answer == strings.ToLower(strings.TrimSpace(userInput)) {
					numCorrectAnswers += 1
					isCorrect = "✅"
				}
				fmt.Printf("Answer: %s (%s)\n", question.Answer, isCorrect)
			}
			numQuestionsAsked++
		}

		outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers)
	}
}

func outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers int) {
	scoreAsPercentage := math.Floor(float64(numCorrectAnswers) / float64(numQuestionsAsked) * 100)
	fmt.Printf("\n\nYour score was: %d/%d (%.0f%%)\n\n\n", numCorrectAnswers, numQuestionsAsked, scoreAsPercentage)
	os.Exit(1)
}
