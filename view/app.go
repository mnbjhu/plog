package view

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	Select   SelectModel
	Logs     TableModel
	Selected bool
	Scanner  *bufio.Scanner
}

var (
	logRegex = `^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{2}:\d{2})\s+(\w+)\s+(\d+)\s+---\s+\[\s*([^\]]+)\s*\]\s+([^ ]+)\s*:\s*(.*)$`
	matcher  = regexp.MustCompile(logRegex)
)

func ReadStdIn(s *bufio.Scanner, c chan table.Row) tea.Cmd {
	return func() tea.Msg {
		// for s.Scan() {
		// 	input := s.Text()
		// 	groups := matcher.FindStringSubmatch(input)
		// 	if len(groups) != 7 {
		// 		c <- table.Row{"---", "ERROR", "---", "---", "---", "Invalid log line"}
		// 		continue
		// 	}
		// 	c <- table.Row{groups[1], groups[2], groups[3], groups[4], groups[5], groups[6]}
		// }
		// return nil
		for {
			c <- table.Row{"---", "ERROR", "---", "---", "---", "Invalid log line"}
			time.Sleep(1 * time.Second)
		}
	}
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
			m.Select.TextArea.SetValue(*msg.Text)
			m.Selected = true
			m.Logs.Table.Blur()
		} else {
			m.Select.TextArea.SetValue("")
			m.Selected = false
			m.Logs.Table.Focus()
		}
	}
	var logsCmd tea.Cmd
	m.Logs, logsCmd = m.Logs.Update(msg)
	var areaCmd tea.Cmd
	var s tea.Model
	s, areaCmd = m.Select.Update(msg)
	m.Select = s.(SelectModel)
	return m, tea.Batch(logsCmd, areaCmd)
}

func (m AppModel) View() string {
	if m.Selected {
		return baseStyle.Render(m.Select.View())
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
