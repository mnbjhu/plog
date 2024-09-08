package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mnbjhu/plog/view"
)

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]
	p := exec.Command(cmd, args...)
	out, err := p.StdoutPipe()
	p.Start()
	defer p.Cancel()
	if err != nil {
		panic(err)
	}

	m := view.NewAppModel(out)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
