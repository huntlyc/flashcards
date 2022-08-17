package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QuestionPair struct {
	Question string
	Answer   string
}

type RunQuizModel struct {
	questionPairs      []QuestionPair
	userAnswers        []string
	questionsAsked     int
	correctAnswers     int
	currentQuestionIdx int
	isGameFinished     bool

	textInput textinput.Model
}

func NewQuiz() *RunQuizModel {
	ti := textinput.New()
	ti.Placeholder = "..your answer"
	ti.Focus()
	ti.Width = 100

	questionPairs, err := parseQuestionSource("cards.json")
	if err != nil {
		log.Fatal(err)
	} else {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questionPairs), func(i, j int) {
			questionPairs[i], questionPairs[j] = questionPairs[j], questionPairs[i]
		})
	}
	return &RunQuizModel{
		questionPairs:      questionPairs,
		userAnswers:        []string{},
		questionsAsked:     0,
		correctAnswers:     0,
		currentQuestionIdx: 0,
		isGameFinished:     false,
		textInput:          ti,
	}
}

func (m RunQuizModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RunQuizModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return SelectModel(0)
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

func (m RunQuizModel) View() string {
	if !m.isGameFinished {

		boldStyle := lipgloss.NewStyle().Bold(true)
		question := boldStyle.Render(m.questionPairs[m.currentQuestionIdx].Question)

		s := fmt.Sprintf("%s\n\n%s", question, m.textInput.View())

		return lipgloss.JoinVertical(lipgloss.Left, s, HelpMenu("spinner"))
	} else {
		return outputFinalQuizResults(m)
	}
}

func outputFinalQuizResults(m RunQuizModel) string {

	// Lipgloss styles
	headingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Margin(1, 0)

	borderedStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#626262")).
		BorderTop(true).
		BorderBottom(true)

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	boldStyle := lipgloss.NewStyle().Bold(true)

	// show correct/incorrect answers
	answersText := headingStyle.Render("Score Card")
	for i, qp := range m.questionPairs {

		question := boldStyle.Render(qp.Question)
		answer := boldStyle.Render(qp.Answer)

		if qp.Answer == m.userAnswers[i] {
			isCorrect := "✅"
			answersText += fmt.Sprintf("\n%s\n%s %s \n", question, answer, isCorrect)
		} else {
			isCorrect := "❌"
			answersText += fmt.Sprintf("\n%s\n%s %s (%s)\n", question, m.userAnswers[i], isCorrect, answer)
		}
	}

	scorePercentage := math.Floor(float64(m.correctAnswers) / float64(m.questionsAsked) * 100)
	scoreText := fmt.Sprintf("Your score was: %d/%d (%.0f%%)", m.correctAnswers, m.questionsAsked, scorePercentage)

	helpText := helpStyle.Render("Press Ctrl+C to exit - Enter to restart - Escape for main menu")

	return fmt.Sprintf("%s\n%s\n\n%s", answersText, borderedStyle.Render(scoreText), helpText)
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
