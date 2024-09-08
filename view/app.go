package view

import (
	"bufio"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppModel struct {
	Select   SelectModel
	Logs     TableModel
	Selected bool
	Scanner  *bufio.Scanner
	Width    int
	Height   int
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(ReadStdIn(m.Scanner, m.Logs.Channel), Wait(m.Logs.Channel))
}

func NewAppModel(out io.Reader) AppModel {
	return AppModel{
		Select:   NewSelectModel(),
		Logs:     NewTableModel(),
		Selected: false,
		Scanner:  bufio.NewScanner(out),
	}
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case selectMsg:
		if msg.Text != nil {
			m.Select.Viewport.SetContent(*msg.Text)
			m.Selected = true
			m.Logs.Table.Blur()
		} else {
			m.Select.Viewport.SetContent("")
			m.Selected = false
			m.Logs.Table.Focus()
		}
		m.Resize(m.Width, m.Height)
	case tea.WindowSizeMsg:
		m.Resize(msg.Width-2, msg.Height-2)
	}
	var logsCmd tea.Cmd
	m.Logs, logsCmd = m.Logs.Update(msg)
	var areaCmd tea.Cmd
	var s tea.Model
	s, areaCmd = m.Select.Update(msg)
	m.Select = s.(SelectModel)
	return m, tea.Batch(logsCmd, areaCmd)
}

func (m *AppModel) Resize(width, height int) {
	m.Width = width
	m.Height = height
	w := width
	if m.Selected {
		w = w / 2
	}
	m.Logs = m.Logs.Resize(w, height)
	m.Select = m.Select.Resize(w, height)
}

func (m AppModel) View() string {
	if m.Selected {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.Logs.View(), m.Select.View())
	}
	return m.Logs.View()
}

type selectMsg struct {
	Text *string
}

func SelectMsg(text *string) tea.Cmd {
	return func() tea.Msg {
		return selectMsg{Text: text}
	}
}

func ReadStdIn(s *bufio.Scanner, c chan string) tea.Cmd {
	return func() tea.Msg {
		for s.Scan() {
			c <- s.Text()
		}
		return nil
	}
}
