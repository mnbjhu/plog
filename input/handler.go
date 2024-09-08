package input

import (
	"bufio"
	"os"
	"regexp"

	"github.com/charmbracelet/bubbles/table"
)

// 2024-09-07T15:24:45.382+01:00  INFO 36146 --- [           main] org.example.App                          : Starting App using Java 17.0.12 with PID 36146 (/home/james/projects/java_thing/app/build/libs/app.jar started by james in /home/james/projects/java_thing)

var (
	logRegex = `^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{2}:\d{2})\s+(\w+)\s+(\d+)\s+---\s+\[\s*([^\]]+)\s*\]\s+([^ ]+)\s*:\s*(.*)$`
	matcher  = regexp.MustCompile(logRegex)
)

func ReadInput(stdInChan chan table.Row) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		groups := matcher.FindStringSubmatch(input)
		if len(groups) != 7 {

			stdInChan <- table.Row{"---", "ERROR", "---", "---", "---", "Invalid log line"}
			continue
		}
		stdInChan <- table.Row{groups[1], groups[2], groups[3], groups[4], groups[5], groups[6]}
	}
}
