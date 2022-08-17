package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditModel struct {
}

func (m EditModel) Init() tea.Cmd {
	return nil
}

func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return SelectModel(0)
		case tea.KeyEnter:
			return SelectModel(0)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m EditModel) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	helpText := helpStyle.Render("Press Ctrl+C to exit - Escape/Return for main menu")

	return fmt.Sprintf("TODO\n\n%s", helpText)
}

func NewEdit() *EditModel {
	return &EditModel{}
}
