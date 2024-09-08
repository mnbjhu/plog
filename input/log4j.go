package input

import (
	"bufio"
	"io"
	"regexp"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	logRegex = `^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{2}:\d{2})\s+(\w+)\s+(\d+)\s+---\s+\[\s*([^\]]+)\s*\]\s+([^ ]+)\s*:\s*(.*)$`
	Matcher  = regexp.MustCompile(logRegex)
)

type Log4jHandler struct {
	MsgAppender chan string
	RowAppender chan table.Row
	Reader      io.Reader
}

func (h Log4jHandler) GetColumns() []string {
	return []string{"Date", "Level", "PID", "Thread", "Class", "Message"}
}

func (h Log4jHandler) GetLevelColumnIndex() int {
	return 1
}

func (h Log4jHandler) GetMsgColumnIndex() int {
	return 5
}

func (h Log4jHandler) HandleLog() tea.Cmd {
	return func() tea.Msg {
		scanner := bufio.NewScanner(h.Reader)
		for scanner.Scan() {
			line := scanner.Text()
			groups := Matcher.FindStringSubmatch(line)
			if len(groups) != 7 {
				h.MsgAppender <- line
			} else {
				row := table.Row{
					groups[1],
					groups[2],
					groups[3],
					groups[4],
					groups[5],
					groups[6],
				}
				h.RowAppender <- row
			}
		}
		return nil
	}
}
