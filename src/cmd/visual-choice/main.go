package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"visual-choice/internal/models"
	"visual-choice/internal/server"
)

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
	fmt.Println(`visual-choice - 可视化选择工具

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

	// 验证端口号范围
	if *port < 1 || *port > 65535 {
		fmt.Fprintf(os.Stderr, "端口号必须在 1-65535 范围内\n")
		os.Exit(1)
	}

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
	srv := server.NewServer(*port, screenDir, stateDir)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "启动服务器失败：%v\n", err)
		os.Exit(1)
	}

	// 写入服务器信息
	info := models.ServerInfo{
		Type:      "server-started",
		Port:      *port,
		URL:       fmt.Sprintf("http://localhost:%d", *port),
		ScreenDir: screenDir,
		StateDir:  stateDir,
	}
	if err := writeServerInfo(stateDir, info); err != nil {
		fmt.Fprintf(os.Stderr, "写入服务器信息失败：%v\n", err)
		srv.Stop()
		os.Exit(1)
	}

	// 写入 PID 文件
	pid := os.Getpid()
	if err := writePIDFile(sessionDir, pid); err != nil {
		fmt.Fprintf(os.Stderr, "写入 PID 文件失败：%v\n", err)
		srv.Stop()
		os.Exit(1)
	}

	fmt.Printf("服务器已启动\n")
	fmt.Printf("URL: %s\n", info.URL)
	fmt.Printf("Screen 目录：%s\n", screenDir)
	fmt.Printf("State 目录：%s\n", stateDir)
	fmt.Printf("\n按 Ctrl+C 停止服务器\n")

	// 等待退出信号
	srv.Wait()

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
	if err := os.WriteFile(stopFile, []byte(time.Now().Format(time.RFC3339)), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "写入停止标记失败：%v\n", err)
	}
}

func handleEvents(args []string) {
	fs := flag.NewFlagSet("events", flag.ExitOnError)
	dir := fs.String("dir", "./session", "会话目录")
	fs.Parse(args)

	stateDir := filepath.Join(*dir, "state")
	eventsFile := filepath.Join(stateDir, "events.jsonl")

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

// splitLines 按行分割字符串
func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

// 辅助函数

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func writePIDFile(dir string, pid int) error {
	pidFile := filepath.Join(dir, "server.pid")
	// 使用更严格的文件权限 (0600)
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0600)
}

func readPIDFile(dir string) (int, error) {
	pidFile := filepath.Join(dir, "server.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}

	// 验证 PID 文件格式
	pidStr := strings.TrimSpace(string(data))
	if pidStr == "" {
		return 0, fmt.Errorf("PID 文件为空")
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("PID 格式错误：%w", err)
	}

	if pid <= 0 {
		return 0, fmt.Errorf("PID 值无效：%d", pid)
	}

	return pid, nil
}

func removePIDFile(dir string) error {
	pidFile := filepath.Join(dir, "server.pid")
	return os.Remove(pidFile)
}

func writeServerInfo(dir string, info models.ServerInfo) error {
	infoFile := filepath.Join(dir, "server-info.json")
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(infoFile, data, 0600)
}

func readServerInfo(dir string) (*models.ServerInfo, error) {
	infoFile := filepath.Join(dir, "server-info.json")
	data, err := os.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var info models.ServerInfo
	err = json.Unmarshal(data, &info)
	return &info, err
}
