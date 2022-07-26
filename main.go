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

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type questionPair struct {
	Question string
	Answer   string
}

type model struct {
	questionPais    []questionPair
	questionsAsked  int
	correctAnswers  int
	currentQuestion int

	textInput textinput.Model
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
		return nil, errors.New("Invalid input type, should be json")
	}
}

func (m model) Init() tea.Cmd {

	return textinput.Blink
}

func initialModel(questionPairs []questionPair) model {
	ti := textinput.New()
	ti.Placeholder = "..answer"
	ti.Focus()
	ti.Width = 100

	return model{
		questionPais:    questionPairs,
		questionsAsked:  0,
		correctAnswers:  0,
		currentQuestion: 0,
		textInput:       ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.currentQuestion > len(m.questionPais) {
				return m, nil
			}

			if m.questionPais[m.currentQuestion].Answer == strings.ToLower(strings.TrimSpace(m.textInput.Value())) {
				m.questionsAsked += 1
			}

			m.currentQuestion += 1
			/*
				                isCorrect := "❌"
								isCorrect = "✅"
							//fmt.Printf("Answer: %s (%s)\n", question.Answer, isCorrect)
			*/
		}
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.questionPais[m.currentQuestion].Question, m.textInput.View())

}

func main() {
	/*
		numQuestionsAsked := 0
		numCorrectAnswers := 0
	*/

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

		p := tea.NewProgram(initialModel(questionPairs))
		if err := p.Start(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		/* ask questions
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
		*/

	}
}

func outputFinalQuizResults(numQuestionsAsked, numCorrectAnswers int) string {
	scoreAsPercentage := math.Floor(float64(numCorrectAnswers) / float64(numQuestionsAsked) * 100)
	return fmt.Sprintf("\n\nYour score was: %d/%d (%.0f%%)\n\n\n", numCorrectAnswers, numQuestionsAsked, scoreAsPercentage)
}
