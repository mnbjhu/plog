package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mnbjhu/plog/input"
	"github.com/mnbjhu/plog/view"
)

func main() {
	stdInChan := make(chan table.Row)
	go input.ReadInput(stdInChan)

	m := view.NewTableModel(stdInChan)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
