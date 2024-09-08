package input

import (
	tea "github.com/charmbracelet/bubbletea"
)

// 2024-09-07T15:24:45.382+01:00  INFO 36146 --- [           main] org.example.App                          : Starting App using Java 17.0.12 with PID 36146 (/home/james/projects/java_thing/app/build/libs/app.jar started by james in /home/james/projects/java_thing)

type LogHandler interface {
	HandleLog() tea.Cmd
	GetColumns() []string
	GetMsgColumnIndex() int
	GetLevelColumnIndex() int
}
