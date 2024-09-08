package input

import (
	"bufio"
	"io"
	"regexp"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type CustomHandlerConfig struct {
	Regex   string   `json:"regex"`
	Columns []string `json:"columns"`
	Level   string   `json:"levelIndex"`
}

func NewCustomHandler(reader io.Reader, config CustomHandlerConfig) CustomHandler {
	return CustomHandler{
		Reader: reader,
		Config: config,
	}
}

type CustomHandler struct {
	MsgAppender chan string
	RowAppender chan table.Row
	Reader      io.Reader
	Config      CustomHandlerConfig
}

func (h CustomHandler) GetColumns() []string {
	return h.Config.Columns
}

func (h CustomHandler) GetLevelColumnIndex() int {
	for i, col := range h.Config.Columns {
		if col == h.Config.Level {
			return i
		}
	}
	panic("Unable to find column")
}

func (h CustomHandler) GetMsgColumnIndex() int {
	return len(h.Config.Columns) - 1
}

func (h CustomHandler) HandleLog() tea.Cmd {
	return func() tea.Msg {
		scanner := bufio.NewScanner(h.Reader)
		expr := regexp.MustCompile(h.Config.Regex)
		for scanner.Scan() {
			line := scanner.Text()
			// TODO: Re-add log appending
			groups := expr.FindStringSubmatch(line)
			row := groups[1:]
			h.RowAppender <- row
		}
		return nil
	}
}
