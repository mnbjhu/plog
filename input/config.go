package input

import (
	"encoding/json"
	"fmt"
	"os"
)

var configFile = ".plog.json"

type ColumnDef struct {
	Title string `json:"name"`
	Width int    `json:"width"`
}

type Config struct {
	Columns []ColumnDef `json:"columns"`
	Regex   string      `json:"regex"`
	Input   string      `json:"input"`
}

func (c Config) Save() {
	file, err := os.Create(configFile)
	if err != nil {
		panic(fmt.Errorf("failed to open config file: %w", err))
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		panic(fmt.Errorf("failed to encode config file: %w", err))
	}
}

func DefaultConfig() Config {
	return Config{
		Columns: []ColumnDef{
			{Title: "Date", Width: 10},
			{Title: "Level", Width: 5},
			{Title: "PID", Width: 7},
			{Title: "Thread", Width: 6},
			{Title: "Class", Width: 10},
			{Title: "Msg", Width: 18},
		},
		Regex: logRegex,
		Input: "stdout",
	}
}

func LeadingRowSize(config Config) int {
	size := 0
	for _, col := range config.Columns {
		size += col.Width
	}
	return size + config.Columns[len(config.Columns)-1].Width
}

func GetConfig() Config {
	file, err := os.Open(configFile)
	if err != nil {
		return DefaultConfig()
	}
	config := Config{}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic(fmt.Errorf("failed to decode config file: %w", err))
	}
	return config
}
