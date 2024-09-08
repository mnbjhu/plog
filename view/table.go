package view

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mnbjhu/plog/input"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewTableModel() TableModel {
	columns := []table.Column{
		{Title: "Date", Width: 10},
		{Title: "Level", Width: 5},
		{Title: "Pid", Width: 6},
		{Title: "Thread", Width: 6},
		{Title: "Class", Width: 10},
		{Title: "Message", Width: 21},
	}

	rows := []table.Row{}
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithStyleFunc(func(row, col int, value string) lipgloss.Style {
			if col == 1 {
				switch value {
				case "ERROR":
					return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
				case "WARN":
					return lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
				case "INFO":
					return lipgloss.NewStyle().Foreground(lipgloss.Color("74"))
				case "DEBUG":
					return lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				case "TRACE":
					return lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				}
			}
			return s.Cell
		}),
	)
	t.SetStyles(s)
	t.SetColumns(columns)
	channel := make(chan string)
	m := TableModel{Table: t, Channel: channel}
	return m
}

type TableModel struct {
	Table   table.Model
	Channel chan string
}

type NewLogLineMsg struct {
	Text string
}

func LogLineMsg(text string) tea.Cmd {
	return func() tea.Msg {
		return NewLogLineMsg{Text: text}
	}
}

func Wait(c chan string) tea.Cmd {
	return func() tea.Msg {
		return NewLogLineMsg{Text: <-c}
	}
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			text := m.Table.SelectedRow()[5]
			return m, SelectMsg(&text)
		}
	case NewLogLineMsg:
		groups := input.Matcher.FindStringSubmatch(msg.Text)
		if len(groups) == 7 {
			row := table.Row{groups[1], groups[2], groups[3], groups[4], groups[5], groups[6]}
			rows := append(m.Table.Rows(), row)
			m.Table.SetRows(rows)
		} else {
			rows := m.Table.Rows()
			if len(rows) > 0 {
				current := rows[len(rows)-1][5]
				rows[len(rows)-1][5] = current + "\n" + msg.Text
				m.Table.SetRows(rows)
			}
		}
		m.Table.GotoBottom()
		m.Table, cmd = m.Table.Update(nil)
		return m, Wait(m.Channel)
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return baseStyle.Render(m.Table.View())
}

func (m TableModel) Resize(width, height int) TableModel {
	m.Table.SetWidth(width - 2)
	m.Table.SetHeight(height)
	columns := m.Table.Columns()
	colWidth := width - 51
	if colWidth < 4 {
		colWidth = 4
	}
	columns[5].Width = colWidth
	m.Table.SetColumns(columns)
	return m
}
