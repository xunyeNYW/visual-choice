package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// ServerInfo 服务器启动信息
type ServerInfo struct {
	Type      string `json:"type"`
	Port      int    `json:"port"`
	URL       string `json:"url"`
	ScreenDir string `json:"screen_dir"`
	StateDir  string `json:"state_dir"`
}

// PIDFile 保存服务器进程 ID
func writePIDFile(dir string, pid int) error {
	pidFile := filepath.Join(dir, "server.pid")
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

// readPIDFile 读取服务器进程 ID
func readPIDFile(dir string) (int, error) {
	pidFile := filepath.Join(dir, "server.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

// removePIDFile 删除 PID 文件
func removePIDFile(dir string) error {
	pidFile := filepath.Join(dir, "server.pid")
	return os.Remove(pidFile)
}

// writeServerInfo 写入服务器信息
func writeServerInfo(dir string, info ServerInfo) error {
	infoFile := filepath.Join(dir, "server-info.json")
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(infoFile, data, 0644)
}

// readServerInfo 读取服务器信息
func readServerInfo(dir string) (*ServerInfo, error) {
	infoFile := filepath.Join(dir, "server-info.json")
	data, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var info ServerInfo
	err = json.Unmarshal(data, &info)
	return &info, err
}

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStart(os.Args[2:])
	case "status":
		handleStatus(os.Args[2:])
	case "stop":
		handleStop(os.Args[2:])
	case "events":
		handleEvents(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "未知命令：%s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`visual-choice - 可视化选择 MVP 应用

用法:
  visual-choice <command> [选项]

命令:
  start     启动服务器
  status    查看服务器状态
  stop      停止服务器
  events    查看事件记录

启动选项:
  --port    端口号 (默认：5234)
  --dir     会话目录 (默认：./session)

示例:
  visual-choice start --port 5234 --dir ./session
  visual-choice status --dir ./session
  visual-choice stop --dir ./session
  visual-choice events --dir ./session`)
}

func handleStart(args []string) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	port := fs.Int("port", 5234, "端口号")
	dir := fs.String("dir", "./session", "会话目录")
	fs.Parse(args)

	// 确保目录存在
	sessionDir := *dir
	if err := ensureDir(sessionDir); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败：%v\n", err)
		os.Exit(1)
	}

	// 检查服务器是否已在运行
	if _, err := readPIDFile(sessionDir); err == nil {
		fmt.Println("服务器已在运行中")
		info, err := readServerInfo(sessionDir)
		if err == nil {
			fmt.Printf("URL: %s\n", info.URL)
		}
		return
	}

	// 创建子目录结构
	screenDir := filepath.Join(sessionDir, "screens")
	stateDir := filepath.Join(sessionDir, "state")
	if err := ensureDir(screenDir); err != nil {
		fmt.Fprintf(os.Stderr, "创建 screens 目录失败：%v\n", err)
		os.Exit(1)
	}
	if err := ensureDir(stateDir); err != nil {
		fmt.Fprintf(os.Stderr, "创建 state 目录失败：%v\n", err)
		os.Exit(1)
	}

	// 启动 HTTP 服务器
	server := NewServer(*port, screenDir, stateDir)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "启动服务器失败：%v\n", err)
		os.Exit(1)
	}

	// 写入服务器信息
	info := ServerInfo{
		Type:      "server-started",
		Port:      *port,
		URL:       fmt.Sprintf("http://localhost:%d", *port),
		ScreenDir: screenDir,
		StateDir:  stateDir,
	}
	if err := writeServerInfo(stateDir, info); err != nil {
		fmt.Fprintf(os.Stderr, "写入服务器信息失败：%v\n", err)
		server.Stop()
		os.Exit(1)
	}

	// 写入 PID 文件
	pid := os.Getpid()
	if err := writePIDFile(sessionDir, pid); err != nil {
		fmt.Fprintf(os.Stderr, "写入 PID 文件失败：%v\n", err)
		server.Stop()
		os.Exit(1)
	}

	fmt.Printf("服务器已启动\n")
	fmt.Printf("URL: %s\n", info.URL)
	fmt.Printf("Screen 目录：%s\n", screenDir)
	fmt.Printf("State 目录：%s\n", stateDir)
	fmt.Printf("\n按 Ctrl+C 停止服务器\n")

	// 等待退出信号
	server.Wait()

	// 清理 PID 文件
	removePIDFile(sessionDir)
	fmt.Println("\n服务器已停止")
}

func handleStatus(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	dir := fs.String("dir", "./session", "会话目录")
	fs.Parse(args)

	sessionDir := *dir
	pidFile := filepath.Join(sessionDir, "server.pid")

	// 检查 PID 文件
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		fmt.Println("服务器未运行")
		return
	}

	// 读取 PID
	pid, err := readPIDFile(sessionDir)
	if err != nil {
		fmt.Println("服务器未运行 (PID 文件损坏)")
		return
	}

	// 检查进程是否存在
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("服务器未运行 (进程不存在)")
		removePIDFile(sessionDir)
		return
	}

	// 发送信号 0 检查进程是否存活
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		fmt.Println("服务器未运行 (进程已终止)")
		removePIDFile(sessionDir)
		return
	}

	fmt.Println("服务器正在运行")
	fmt.Printf("PID: %d\n", pid)

	// 读取服务器信息
	info, err := readServerInfo(filepath.Join(sessionDir, "state"))
	if err == nil {
		fmt.Printf("URL: %s\n", info.URL)
		fmt.Printf("端口：%d\n", info.Port)
	}
}

func handleStop(args []string) {
	fs := flag.NewFlagSet("stop", flag.ExitOnError)
	dir := fs.String("dir", "./session", "会话目录")
	fs.Parse(args)

	sessionDir := *dir
	pidFile := filepath.Join(sessionDir, "server.pid")

	// 检查 PID 文件
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		fmt.Println("服务器未运行")
		return
	}

	// 读取 PID
	pid, err := readPIDFile(sessionDir)
	if err != nil {
		fmt.Println("服务器未运行 (PID 文件损坏)")
		removePIDFile(sessionDir)
		return
	}

	// 发送 SIGTERM 信号
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("无法找到进程 %d: %v\n", pid, err)
		return
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Printf("无法停止进程 %d: %v\n", pid, err)
		return
	}

	// 等待进程退出
	done := make(chan error, 1)
	go func() {
		_, err := proc.Wait()
		done <- err
	}()

	select {
	case <-time.After(5 * time.Second):
		proc.Signal(syscall.SIGKILL)
		fmt.Println("服务器已强制停止")
	case err := <-done:
		if err != nil {
			fmt.Printf("服务器停止时出错：%v\n", err)
		} else {
			fmt.Println("服务器已停止")
		}
	}

	// 清理 PID 文件
	removePIDFile(sessionDir)

	// 写入停止标记
	stopFile := filepath.Join(sessionDir, "state", "server-stopped")
	os.WriteFile(stopFile, []byte(time.Now().Format(time.RFC3339)), 0644)
}

func handleEvents(args []string) {
	fs := flag.NewFlagSet("events", flag.ExitOnError)
	dir := fs.String("dir", "./session", "会话目录")
	fs.Parse(args)

	stateDir := filepath.Join(*dir, "state")
	eventsFile := filepath.Join(stateDir, "events.jsonl")

	// 检查事件文件
	data, err := os.ReadFile(eventsFile)
	if os.IsNotExist(err) {
		fmt.Println("暂无事件记录")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取事件文件失败：%v\n", err)
		os.Exit(1)
	}

	// 解析并显示事件
	lines := string(data)
	if lines == "" {
		fmt.Println("暂无事件记录")
		return
	}

	fmt.Println("事件记录:")
	fmt.Println("---------")

	for i, line := range splitLines(lines) {
		if line == "" {
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			fmt.Printf("[%d] 解析失败：%s\n", i+1, line)
			continue
		}

		timestamp := ""
		if ts, ok := event["timestamp"].(float64); ok {
			timestamp = time.Unix(int64(ts), 0).Format("15:04:05")
		}

		eventType := "click"
		if t, ok := event["type"].(string); ok {
			eventType = t
		}

		choice := ""
		if c, ok := event["choice"].(string); ok {
			choice = c
		}

		text := ""
		if t, ok := event["text"].(string); ok {
			text = t
		}

		fmt.Printf("[%d] %s %s - 选择：%s - %s\n", i+1, timestamp, eventType, choice, text)
	}
}

// splitLines 简单的按行分割
func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
