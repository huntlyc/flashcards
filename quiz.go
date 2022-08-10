package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QuestionPair struct {
	Question string
	Answer   string
}

type runQuizModel struct {
	questionPairs      []QuestionPair
	userAnswers        []string
	questionsAsked     int
	correctAnswers     int
	currentQuestionIdx int
	isGameFinished     bool

	textInput textinput.Model
}

func (m runQuizModel) Init() tea.Cmd {
	return textinput.Blink
}

func NewQuiz(questionPairs []QuestionPair) tea.Model {
	ti := textinput.New()
	ti.Placeholder = "..your answer"
	ti.Focus()
	ti.Width = 100

	return runQuizModel{
		questionPairs:      questionPairs,
		userAnswers:        []string{},
		questionsAsked:     0,
		correctAnswers:     0,
		currentQuestionIdx: 0,
		isGameFinished:     false,
		textInput:          ti,
	}
}

func (m runQuizModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m runQuizModel) View() string {
	if !m.isGameFinished {

		boldStyle := lipgloss.NewStyle().Bold(true)
		question := boldStyle.Render(m.questionPairs[m.currentQuestionIdx].Question)

		s := fmt.Sprintf("%s\n\n%s", question, m.textInput.View())

		return lipgloss.JoinVertical(lipgloss.Left, s, HelpMenu("spinner"))
	} else {
		return outputFinalQuizResults(m)
	}
}

func outputFinalQuizResults(m runQuizModel) string {

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

	helpText := helpStyle.Render("Press Ctrl+C to exit - Enter to restart")

	return fmt.Sprintf("%s\n%s\n\n%s", answersText, borderedStyle.Render(scoreText), helpText)
}
