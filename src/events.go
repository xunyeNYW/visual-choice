package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Event 用户交互事件
type Event struct {
	Type       string    `json:"type"`
	Choice     string    `json:"choice"`
	Text       string    `json:"text"`
	Timestamp  int64     `json:"timestamp"`
	ServerTime string    `json:"server_time,omitempty"`
}

// ReadEvents 读取事件文件
func ReadEvents(stateDir string) ([]Event, error) {
	eventsFile := filepath.Join(stateDir, "events.jsonl")

	// 检查文件是否存在
	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		return []Event{}, nil
	}

	file, err := os.Open(eventsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			// 跳过解析失败的行
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

// FormatEvents 格式化显示事件
func FormatEvents(events []Event) string {
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

// GetLastChoice 获取最后一次选择
func GetLastChoice(events []Event) *Event {
	if len(events) == 0 {
		return nil
	}

	// 从后往前找最后一个 click 事件
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Type == "click" {
			return &events[i]
		}
	}

	return nil
}

// GetAllChoices 获取所有选择（去重）
func GetAllChoices(events []Event) []string {
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
