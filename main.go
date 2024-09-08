package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mnbjhu/plog/input"
	"github.com/mnbjhu/plog/view"
)

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		fmt.Println("Usage: plog <command>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]
	p := exec.CommandContext(ctx, cmd, args...)

	out, err := p.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	defer p.Cancel()

	config := input.GetConfig()

	m := view.NewAppModel(out, config)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
