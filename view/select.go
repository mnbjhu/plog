package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectModel struct {
	Viewport viewport.Model
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func NewSelectModel() SelectModel {
	return SelectModel{
		Viewport: viewport.New(0, 0),
	}
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, SelectMsg(nil)
		}
	}
	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m SelectModel) View() string {
	return baseStyle.Render(m.Viewport.View())
}

func (m SelectModel) Resize(width, height int) SelectModel {
	m.Viewport.Width = width - 2
	m.Viewport.Height = height
	return m
}
