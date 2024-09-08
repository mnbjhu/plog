package input

import (
	"bufio"
	"context"
	"os/exec"
	"regexp"
	"testing"
)

func TestParseLogLine(t *testing.T) {
	expr := regexp.MustCompile(logRegex)
	res := expr.FindStringSubmatch("2024-09-07T15:24:45.382+01:00  INFO 36146 --- [           main] org.example.App                          : Hello this is an example message")
	if len(res) != 7 {
		t.Errorf("Expected 7 matches, got %d", len(res))
	}
	date := res[1]
	if date != "2024-09-07T15:24:45.382+01:00" {
		t.Errorf("Expected 2024-09-07T15:24:45.382+01:00, got %s", date)
	}

	level := res[2]
	if level != "INFO" {
		t.Errorf("Expected INFO, got %s", level)
	}

	pid := res[3]
	if pid != "36146" {
		t.Errorf("Expected 36146, got %s", pid)
	}

	thread := res[4]
	if thread != "main" {
		t.Errorf("Expected main, got %s", thread)
	}

	class := res[5]
	if class != "org.example.App" {
		t.Errorf("Expected org.example.App, got %s", class)
	}

	message := res[6]
	if message != "Hello this is an example message" {
		t.Errorf("Expected Hello this is an example message, got %s", message)
	}
}

func TestRunCmd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "echo", "Hello, World!")
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(out)
	if scanner.Scan() {
		if scanner.Text() != "Hello, World!" {
			t.Errorf("Expected Hello, World!, got %s", scanner.Text())
		}
	}
	cancel()
}
