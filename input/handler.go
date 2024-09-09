package input

import (
	"bufio"
	"io"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oriser/regroup"
)

var (
	logRegex = `^(?<Date>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{2}:\d{2})\s+(?<Level>\w+)\s+(?<PID>\d+)\s+---\s+\[\s*(?<Thread>[^\]]+)\s*\]\s+(?<Class>[^ ]+)\s*:\s*(?<Msg>.*)$`
	Matcher  = regroup.MustCompile(logRegex)
)

type Log4jHandler struct {
	MsgAppender chan string
	RowAppender chan table.Row
	Reader      io.Reader
	Columns     []string
	Regex       *regroup.ReGroup
	LeadingSize int
}

func (h Log4jHandler) GetColumns() []string {
	return h.Columns
}

func (h Log4jHandler) GetLevelColumnIndex() int {
	for i, col := range h.GetColumns() {
		if strings.ToLower(col) == "level" {
			return i
		}
	}
	return -1
}

func (h Log4jHandler) GetMsgColumnIndex() int {
	return len(h.GetColumns()) - 1
}

func (h Log4jHandler) HandleLog() tea.Cmd {
	return func() tea.Msg {
		scanner := bufio.NewScanner(h.Reader)
		for scanner.Scan() {
			line := scanner.Text()
			line = stripansi.Strip(line)
			groups, err := h.Regex.Groups(line)
			if err != nil {
				h.MsgAppender <- line
				// h.RowAppender <- table.Row{"-", "ERROR", "-", line}
			} else {
				row := table.Row{}
				for _, col := range h.Columns {
					val, ok := groups[col]
					if ok {
						row = append(row, val)
					} else {
						row = append(row, "")
					}
				}
				h.RowAppender <- row
			}
		}
		return nil
	}
}
