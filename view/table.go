package view

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	channel := make(chan table.Row)
	m := TableModel{Table: t, Channel: channel}
	return m
}

type TableModel struct {
	Table   table.Model
	Channel chan table.Row
}

type NewLogLineMsg struct {
	Row table.Row
}

func LogLineMsg(row table.Row) tea.Cmd {
	return func() tea.Msg {
		return NewLogLineMsg{Row: row}
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
	case tea.WindowSizeMsg:
		width := msg.Width - 51
		if width < 5 {
			width = 5
		}
		m.Table.SetWidth(msg.Width - 2)
		m.Table.SetHeight(msg.Height - 4)
		columns := m.Table.Columns()
		columns[5].Width = width
		m.Table.SetColumns(columns)
	case NewLogLineMsg:
		rows := append(m.Table.Rows(), msg.Row)
		m.Table.SetRows(rows)
		m.Table.GotoBottom()
		return m, nil
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return baseStyle.Render(m.Table.View())
}
