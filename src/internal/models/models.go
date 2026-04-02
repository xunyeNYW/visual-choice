package models

import "time"

type ServerInfo struct {
	Type      string `json:"type"`
	Port      int    `json:"port"`
	URL       string `json:"url"`
	ScreenDir string `json:"screen_dir"`
	StateDir  string `json:"state_dir"`
}

type Event struct {
	Type       string `json:"type"`
	Choice     string `json:"choice"`
	Text       string `json:"text"`
	Timestamp  int64  `json:"timestamp"`
	ServerTime string `json:"server_time,omitempty"`
}

type ServerConfig struct {
	Port         int
	ScreenDir    string
	StateDir     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         5234,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
