package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	itemStyle := lipgloss.NewStyle().PaddingLeft(4)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			selectedItemStyle := lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type MenuModel struct {
	list     list.Model
	items    []item
	choice   string
	quitting bool
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i := m.list.Cursor()
			switch i {
			case 0:
				return SelectModel(1)
			case 1:
				return SelectModel(2)
			case 2:
				return m, tea.Quit
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	return "\n" + m.list.View()
}

func NewMenu() *MenuModel {
	items := []list.Item{
		item("Run Quiz"),
		item("Edit Questions"),
		item("Quit"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What's the plan?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	titleStyle := lipgloss.NewStyle().MarginLeft(2)
	l.Styles.Title = titleStyle
	paginationStyle := list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.PaginationStyle = paginationStyle
	helpStyle := list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	l.Styles.HelpStyle = helpStyle

	return &MenuModel{list: l}
}
