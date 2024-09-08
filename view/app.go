package view

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mnbjhu/plog/input"
	"github.com/oriser/regroup"
)

type AppModel struct {
	Select   SelectModel
	Logs     TableModel
	Selected bool
	Width    int
	Height   int
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(m.Logs.LogHandler.HandleLog(), Wait(m.Logs.LogChannel, m.Logs.MsgChannel))
}

// 2024-09-08T17:35:56.385+01:00 [Runtime]  INFO sqlx::postgres::notice: relation "_sqlx_migrations" already exists, skipping
// 2024-09-08T17:35:56.387+01:00 [Runtime] Starting on 127.0.0.1:8000

func NewAppModel(out io.Reader, config input.Config) AppModel {
	logs := NewTableModel(config)
	columns := make([]string, len(config.Columns))
	for i, col := range config.Columns {
		columns[i] = col.Title
	}
	handler := input.Log4jHandler{
		MsgAppender: logs.MsgChannel,
		RowAppender: logs.LogChannel,
		Reader:      out,
		Regex:       regroup.MustCompile(config.Regex),
		Columns:     columns,
	}
	logs.LogHandler = handler
	return AppModel{
		Select:   NewSelectModel(),
		Logs:     logs,
		Selected: false,
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
