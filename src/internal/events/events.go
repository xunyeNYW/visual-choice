package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"visual-choice/internal/models"
)

type Store struct {
	stateDir string
}

func NewStore(stateDir string) *Store {
	return &Store{
		stateDir: stateDir,
	}
}

func (s *Store) ReadEvents() ([]models.Event, error) {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")

	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		return []models.Event{}, nil
	}

	file, err := os.Open(eventsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []models.Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var event models.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			fmt.Fprintf(os.Stderr, "警告：解析事件失败：%v\n", err)
			continue
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Store) Append(event map[string]interface{}) error {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(eventsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}

func (s *Store) Clear() error {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")
	return os.WriteFile(eventsFile, []byte{}, 0600)
}

func FormatEvents(events []models.Event) string {
	if len(events) == 0 {
		return "暂无事件"
	}

	result := ""
	for i, event := range events {
		timestamp := time.Unix(event.Timestamp, 0).Format("15:04:05")
		result += fmt.Sprintf("[%d] %s %s - 选择：%s - %s\n",
			i+1, timestamp, event.Type, event.Choice, event.Text)
	}

	return result
}

func GetLastChoice(events []models.Event) *models.Event {
	if len(events) == 0 {
		return nil
	}

	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Type == "click" {
			return &events[i]
		}
	}

	return nil
}

func GetAllChoices(events []models.Event) []string {
	seen := make(map[string]bool)
	var choices []string

	for _, event := range events {
		if event.Type == "click" && !seen[event.Choice] {
			seen[event.Choice] = true
			choices = append(choices, event.Choice)
		}
	}

	return choices
}
