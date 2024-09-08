package view

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectModel struct {
	TextArea textarea.Model
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func NewSelectModel() SelectModel {
	return SelectModel{
		TextArea: textarea.New(),
	}
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, SelectMsg(nil)
		}
	case tea.WindowSizeMsg:
		width := msg.Width - 51
		if width < 5 {
			width = 5
		}
		m.TextArea.SetWidth(msg.Width - 2)
		m.TextArea.SetHeight(msg.Height - 4)
	}
	return m, nil
}

func (m SelectModel) View() string {
	return m.TextArea.View()
}
