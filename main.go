package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"encoding/json"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const first = 0
const quizModel sessionState = iota

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	current   = quizModel
	models    []tea.Model
)

func HelpMenu(view ...string) string {
	if len(view) != 0 {
		return helpStyle(fmt.Sprintf("right/l: next • left/h: previous • enter: new %s", view[first]))
	}
	return helpStyle("right/l: next • left/h: previous")
}

func NextModel() (tea.Model, tea.Cmd) {
	if int(current) == len(models)-1 {
		current = first
	} else {
		current++
	}
	return models[current], models[current].Init()
}

func PrevModel() (tea.Model, tea.Cmd) {
	if int(current) == first {
		current = sessionState(len(models) - 1)
	} else {
		current--
	}
	return models[current], models[current].Init()
}

// parseJSONFile will parse a valid JSON file with the following format:
// [
//  {"question":"Meaning of life"?, "answer": "42"},
//  ...
// ]
func parseJSONFile(filename string) ([]QuestionPair, error) {

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var questionPairs []QuestionPair
	err = json.Unmarshal(fileContents, &questionPairs)
	if err != nil {
		return nil, err
	}

	return questionPairs, nil
}

// Accepted sources are:
// JSON file - see parseJSONFile() for file format
func parseQuestionSource(filename string) ([]QuestionPair, error) {
	switch {
	case strings.HasSuffix(filename, ".json"):
		return parseJSONFile(filename)
	default:
		return nil, errors.New("Invalid input type, should be json")
	}
}

func main() {
	questionSource := flag.String("f", "cards.json", "json file to read")
	shuffle := flag.Bool("s", false, "shuffle the deck")

	flag.Parse()

	if QuestionPairs, err := parseQuestionSource(*questionSource); err != nil {
		log.Fatal(err)
	} else {

		if *shuffle { // randomise QuestionPairs
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(QuestionPairs), func(i, j int) {
				QuestionPairs[i], QuestionPairs[j] = QuestionPairs[j], QuestionPairs[i]
			})
		}
		models = []tea.Model{}
		models = append(models, NewQuiz(QuestionPairs))

		p := tea.NewProgram(models[current])
		if err := p.Start(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}
