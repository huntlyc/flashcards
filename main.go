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
	questionPairs      []questionPair
	userAnswers        []string
	questionsAsked     int
	correctAnswers     int
	currentQuestionIdx int
	isGameFinished     bool

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
	ti.Placeholder = "..your answer"
	ti.Focus()
	ti.Width = 100

	return model{
		questionPairs:      questionPairs,
		userAnswers:        []string{},
		questionsAsked:     0,
		correctAnswers:     0,
		currentQuestionIdx: 0,
		isGameFinished:     false,
		textInput:          ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:

			if m.isGameFinished { // restart
				m.questionsAsked = 0
				m.correctAnswers = 0
				m.currentQuestionIdx = 0
				m.userAnswers = []string{}
				m.isGameFinished = false
				m.textInput.SetValue("")

				return m, nil
			}

			correctAnswer := m.questionPairs[m.currentQuestionIdx].Answer
			userAnswer := strings.ToLower(strings.TrimSpace(m.textInput.Value()))

			m.userAnswers = append(m.userAnswers, userAnswer)

			if userAnswer == correctAnswer {
				m.correctAnswers += 1
			}

			m.questionsAsked += 1

			// check end
			if m.currentQuestionIdx+1 >= len(m.questionPairs) {
				m.isGameFinished = true
				return m, nil
			} else {
				m.currentQuestionIdx += 1
				m.textInput.SetValue("")
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg) // update the textinput with the keypress
	return m, cmd
}

func (m model) View() string {
	if !m.isGameFinished {
		return fmt.Sprintf("%s\n\n%s", m.questionPairs[m.currentQuestionIdx].Question, m.textInput.View())
	} else {
		return outputFinalQuizResults(m)
	}
}

func main() {
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
	}
}

func outputFinalQuizResults(m model) string {
	scorePercentage := math.Floor(float64(m.correctAnswers) / float64(m.questionsAsked) * 100)
	scoreText := fmt.Sprintf("All done!\nYour score was: %d/%d (%.0f%%)\n", m.correctAnswers, m.questionsAsked, scorePercentage)

	// show correct/incorrect answers
	answersText := ""
	for i, qp := range m.questionPairs {
		if qp.Answer == m.userAnswers[i] {
			isCorrect := "✅"
			answersText += fmt.Sprintf("%s\n%s %s \n", qp.Question, qp.Answer, isCorrect)
		} else {
			isCorrect := "❌"
			answersText += fmt.Sprintf("%s\n%s %s (%s)\n", qp.Question, m.userAnswers[i], isCorrect, qp.Answer)
		}
	}

	return fmt.Sprintf("%s\n\nScore Card:\n%s\n\nPress Ctrl+C to exit - Enter to restart", scoreText, answersText)
}
