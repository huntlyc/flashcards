package main

import (
	"fmt"
	"log"
	"os"

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
	return helpStyle(fmt.Sprintf("Ctrl+C: exit • Esc: main menu • Enter: submit answer"))
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
func SelectModel(num int) (tea.Model, tea.Cmd) {
	current = sessionState(num)
	return models[current], models[current].Init()
}

func main() {
	models = []tea.Model{}
	models = append(models, NewMenu())
	models = append(models, NewQuiz())
	models = append(models, NewEdit())

	p := tea.NewProgram(models[current])
	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
