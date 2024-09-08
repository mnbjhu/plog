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

type TableModel struct {
	Table      table.Model
	MsgChannel chan string
	LogChannel chan table.Row
	LogHandler input.LogHandler
}

func NewTableModel(config input.Config) TableModel {
	columns := []table.Column{}
	for _, col := range config.Columns {
		columns = append(columns, table.Column{Width: col.Width, Title: col.Title})
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
	msgChan := make(chan string)
	logChan := make(chan table.Row)
	m := TableModel{Table: t, MsgChannel: msgChan, LogChannel: logChan}
	return m
}

type newLogMsg struct {
	Row table.Row
}

type appendLogMsg struct {
	Text string
}

func Wait(logs chan table.Row, msg chan string) tea.Cmd {
	return func() tea.Msg {
		select {
		case row := <-logs:
			return newLogMsg{Row: row}
		case text := <-msg:
			return appendLogMsg{Text: text}
		}
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
	case appendLogMsg:
		rows := m.Table.Rows()
		if len(rows) > 0 {
			index := m.LogHandler.GetMsgColumnIndex()
			current := rows[len(rows)-1][index]
			rows[len(rows)-1][index] = current + "\n" + msg.Text
			m.Table.SetRows(rows)
			m.Table.GotoBottom()
			m.Table, cmd = m.Table.Update(nil)
		}
		return m, Wait(m.LogChannel, m.MsgChannel)
	case newLogMsg:
		rows := append(m.Table.Rows(), msg.Row)
		m.Table.SetRows(rows)
		m.Table.GotoBottom()
		m.Table, cmd = m.Table.Update(nil)
		return m, Wait(m.LogChannel, m.MsgChannel)
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
	columns[m.LogHandler.GetLevelColumnIndex()].Width = colWidth
	m.Table.SetColumns(columns)
	return m
}
