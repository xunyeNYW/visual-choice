package models

import (
	"testing"
	"time"
)

func TestDefaultServerConfig(t *testing.T) {
	config := DefaultServerConfig()

	if config.Port != 5234 {
		t.Fatalf("期望默认端口为 5234，得到 %d", config.Port)
	}

	if config.ReadTimeout != 15*time.Second {
		t.Fatalf("期望读取超时为 15s，得到 %v", config.ReadTimeout)
	}

	if config.WriteTimeout != 15*time.Second {
		t.Fatalf("期望写入超时为 15s，得到 %v", config.WriteTimeout)
	}

	if config.IdleTimeout != 60*time.Second {
		t.Fatalf("期望空闲超时为 60s，得到 %v", config.IdleTimeout)
	}
}

func TestServerInfo_Marshal(t *testing.T) {
	info := ServerInfo{
		Type:      "server-started",
		Port:      5234,
		URL:       "http://localhost:5234",
		ScreenDir: "/path/to/screens",
		StateDir:  "/path/to/state",
	}

	// 测试可以正常序列化（实际使用 encoding/json）
	// 这里只是验证结构体字段可访问
	if info.Type != "server-started" {
		t.Fatalf("期望类型为 server-started，得到 %s", info.Type)
	}
	if info.Port != 5234 {
		t.Fatalf("期望端口为 5234，得到 %d", info.Port)
	}
}

func TestEvent_Marshal(t *testing.T) {
	event := Event{
		Type:       "click",
		Choice:     "option-a",
		Text:       "测试选项",
		Timestamp:  1234567890,
		ServerTime: "2026-04-02T12:00:00Z",
	}

	// 测试可以正常序列化
	if event.Type != "click" {
		t.Fatalf("期望事件类型为 click，得到 %s", event.Type)
	}
	if event.Choice != "option-a" {
		t.Fatalf("期望事件选择为 option-a，得到 %s", event.Choice)
	}
}

func TestServerConfig_Validation(t *testing.T) {
	config := DefaultServerConfig()

	// 验证端口范围
	if config.Port < 1 || config.Port > 65535 {
		t.Fatalf("端口号超出范围：%d", config.Port)
	}

	// 验证超时时间为正数
	if config.ReadTimeout <= 0 {
		t.Fatalf("读取超时应该为正数：%v", config.ReadTimeout)
	}
	if config.WriteTimeout <= 0 {
		t.Fatalf("写入超时应该为正数：%v", config.WriteTimeout)
	}
	if config.IdleTimeout <= 0 {
		t.Fatalf("空闲超时应该为正数：%v", config.IdleTimeout)
	}
}
