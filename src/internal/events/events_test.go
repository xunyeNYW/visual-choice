package events

import (
	"os"
	"path/filepath"
	"testing"

	"visual-choice/internal/models"
)

func TestStore_Append(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	// 测试追加事件
	event := map[string]interface{}{
		"type":      "click",
		"choice":    "option-a",
		"text":      "测试选项 A",
		"timestamp": 1234567890,
	}

	err := store.Append(event)
	if err != nil {
		t.Fatalf("追加事件失败：%v", err)
	}

	// 验证文件是否创建
	eventsFile := filepath.Join(tmpDir, "events.jsonl")
	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		t.Fatalf("事件文件未创建")
	}
}

func TestStore_ReadEvents(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	// 测试读取空事件
	events, err := store.ReadEvents()
	if err != nil {
		t.Fatalf("读取事件失败：%v", err)
	}
	if len(events) != 0 {
		t.Fatalf("期望 0 个事件，得到 %d 个", len(events))
	}

	// 追加一个事件
	event := map[string]interface{}{
		"type":      "click",
		"choice":    "option-a",
		"text":      "测试选项 A",
		"timestamp": int64(1234567890),
	}

	err = store.Append(event)
	if err != nil {
		t.Fatalf("追加事件失败：%v", err)
	}

	// 读取事件
	events, err = store.ReadEvents()
	if err != nil {
		t.Fatalf("读取事件失败：%v", err)
	}
	if len(events) != 1 {
		t.Fatalf("期望 1 个事件，得到 %d 个", len(events))
	}

	// 验证事件内容
	if events[0].Type != "click" {
		t.Fatalf("期望事件类型为 click，得到 %s", events[0].Type)
	}
	if events[0].Choice != "option-a" {
		t.Fatalf("期望事件选择为 option-a，得到 %s", events[0].Choice)
	}
}

func TestStore_Clear(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	// 追加一个事件
	event := map[string]interface{}{
		"type":      "click",
		"choice":    "option-a",
		"text":      "测试选项 A",
		"timestamp": int64(1234567890),
	}

	err := store.Append(event)
	if err != nil {
		t.Fatalf("追加事件失败：%v", err)
	}

	// 清空事件
	err = store.Clear()
	if err != nil {
		t.Fatalf("清空事件失败：%v", err)
	}

	// 验证事件已清空
	events, err := store.ReadEvents()
	if err != nil {
		t.Fatalf("读取事件失败：%v", err)
	}
	if len(events) != 0 {
		t.Fatalf("期望 0 个事件，得到 %d 个", len(events))
	}
}

func TestFormatEvents(t *testing.T) {
	events := []models.Event{
		{
			Type:      "click",
			Choice:    "option-a",
			Text:      "测试选项 A",
			Timestamp: 1234567890,
		},
		{
			Type:      "click",
			Choice:    "option-b",
			Text:      "测试选项 B",
			Timestamp: 1234567891,
		},
	}

	result := FormatEvents(events)
	if result == "" {
		t.Fatal("期望格式化结果不为空")
	}
}

func TestGetLastChoice(t *testing.T) {
	events := []models.Event{
		{
			Type:      "click",
			Choice:    "option-a",
			Text:      "测试选项 A",
			Timestamp: 1234567890,
		},
		{
			Type:      "click",
			Choice:    "option-b",
			Text:      "测试选项 B",
			Timestamp: 1234567891,
		},
	}

	lastChoice := GetLastChoice(events)
	if lastChoice == nil {
		t.Fatal("期望最后一个选择不为空")
	}
	if lastChoice.Choice != "option-b" {
		t.Fatalf("期望最后一个选择为 option-b，得到 %s", lastChoice.Choice)
	}
}

func TestGetAllChoices(t *testing.T) {
	events := []models.Event{
		{
			Type:      "click",
			Choice:    "option-a",
			Text:      "测试选项 A",
			Timestamp: 1234567890,
		},
		{
			Type:      "click",
			Choice:    "option-b",
			Text:      "测试选项 B",
			Timestamp: 1234567891,
		},
		{
			Type:      "click",
			Choice:    "option-a",
			Text:      "测试选项 A（重复）",
			Timestamp: 1234567892,
		},
	}

	choices := GetAllChoices(events)
	if len(choices) != 2 {
		t.Fatalf("期望 2 个唯一选择，得到 %d 个", len(choices))
	}
}

func TestStore_EmptyFile(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	// 创建空文件
	eventsFile := filepath.Join(tmpDir, "events.jsonl")
	err := os.WriteFile(eventsFile, []byte{}, 0600)
	if err != nil {
		t.Fatalf("创建空文件失败：%v", err)
	}

	// 读取空文件
	events, err := store.ReadEvents()
	if err != nil {
		t.Fatalf("读取空文件失败：%v", err)
	}
	if len(events) != 0 {
		t.Fatalf("期望 0 个事件，得到 %d 个", len(events))
	}
}

func TestStore_InvalidJSON(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	store := NewStore(tmpDir)

	// 写入无效 JSON
	eventsFile := filepath.Join(tmpDir, "events.jsonl")
	err := os.WriteFile(eventsFile, []byte("invalid json\n"), 0600)
	if err != nil {
		t.Fatalf("写入文件失败：%v", err)
	}

	// 读取文件（应该跳过无效行）
	events, err := store.ReadEvents()
	if err != nil {
		t.Fatalf("读取文件失败：%v", err)
	}
	if len(events) != 0 {
		t.Fatalf("期望 0 个事件，得到 %d 个", len(events))
	}
}
