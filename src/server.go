package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Server HTTP 服务器结构
type Server struct {
	port      int
	screenDir string
	stateDir  string
	httpServer *http.Server
	watcher    *fsnotify.Watcher
	wg        sync.WaitGroup
	stopChan  chan struct{}
	mu        sync.RWMutex
	latestFile string
}

// FrameTemplate HTML 框架模板
const FrameTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Visual Choice</title>
    <style>
        :root {
            --primary: #2563eb;
            --primary-hover: #1d4ed8;
            --bg: #f8fafc;
            --card-bg: #ffffff;
            --text: #1e293b;
            --text-muted: #64748b;
            --border: #e2e8f0;
            --shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.6;
            padding: 2rem;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        h2 {
            font-size: 1.875rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
            color: var(--text);
        }

        .subtitle {
            color: var(--text-muted);
            font-size: 1.125rem;
            margin-bottom: 2rem;
        }

        /* 选项卡片样式 */
        .options {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        .option {
            background: var(--card-bg);
            border: 2px solid var(--border);
            border-radius: 0.75rem;
            padding: 1.5rem;
            cursor: pointer;
            transition: all 0.2s ease;
            box-shadow: var(--shadow);
            display: flex;
            gap: 1rem;
        }

        .option:hover {
            border-color: var(--primary);
            transform: translateY(-2px);
            box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1);
        }

        .option.selected {
            border-color: var(--primary);
            background: #eff6ff;
        }

        .option .letter {
            width: 3rem;
            height: 3rem;
            background: var(--primary);
            color: white;
            border-radius: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            font-size: 1.25rem;
            flex-shrink: 0;
        }

        .option.selected .letter {
            background: var(--primary-hover);
        }

        .option .content h3 {
            font-size: 1.25rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
        }

        .option .content p {
            color: var(--text-muted);
            font-size: 0.875rem;
        }

        /* 卡片样式 */
        .cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        .card {
            background: var(--card-bg);
            border: 2px solid var(--border);
            border-radius: 0.75rem;
            overflow: hidden;
            cursor: pointer;
            transition: all 0.2s ease;
            box-shadow: var(--shadow);
        }

        .card:hover {
            border-color: var(--primary);
            transform: translateY(-2px);
        }

        .card.selected {
            border-color: var(--primary);
            background: #eff6ff;
        }

        .card-image {
            width: 100%;
            height: 200px;
            background: var(--border);
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--text-muted);
        }

        .card-body {
            padding: 1.5rem;
        }

        .card-body h3 {
            font-size: 1.125rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
        }

        .card-body p {
            color: var(--text-muted);
            font-size: 0.875rem;
        }

        /* Mockup 容器 */
        .mockup {
            background: var(--card-bg);
            border: 2px solid var(--border);
            border-radius: 0.75rem;
            overflow: hidden;
            margin-bottom: 2rem;
            box-shadow: var(--shadow);
        }

        .mockup-header {
            background: var(--border);
            padding: 0.75rem 1rem;
            font-weight: 600;
            font-size: 0.875rem;
            color: var(--text-muted);
            border-bottom: 2px solid var(--border);
        }

        .mockup-body {
            padding: 1.5rem;
        }

        /* 分割视图 */
        .split {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        @media (max-width: 768px) {
            .split {
                grid-template-columns: 1fr;
            }
        }

        /* 优缺点对比 */
        .pros-cons {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1.5rem;
            margin-bottom: 2rem;
        }

        @media (max-width: 768px) {
            .pros-cons {
                grid-template-columns: 1fr;
            }
        }

        .pros, .cons {
            background: var(--card-bg);
            border: 2px solid var(--border);
            border-radius: 0.75rem;
            padding: 1.5rem;
        }

        .pros {
            border-color: #22c55e;
        }

        .cons {
            border-color: #ef4444;
        }

        .pros h4, .cons h4 {
            font-size: 1rem;
            font-weight: 600;
            margin-bottom: 1rem;
        }

        .pros ul, .cons ul {
            list-style: none;
            padding-left: 0;
        }

        .pros li, .cons li {
            padding: 0.5rem 0;
            padding-left: 1.5rem;
            position: relative;
            font-size: 0.875rem;
        }

        .pros li::before {
            content: "✓";
            position: absolute;
            left: 0;
            color: #22c55e;
            font-weight: 700;
        }

        .cons li::before {
            content: "×";
            position: absolute;
            left: 0;
            color: #ef4444;
            font-weight: 700;
        }

        /* Mock 元素 */
        .mock-nav {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 1rem 1.5rem;
            background: var(--border);
            border-radius: 0.5rem;
            margin-bottom: 1rem;
            font-size: 0.875rem;
        }

        .mock-sidebar {
            width: 200px;
            height: 300px;
            background: var(--border);
            border-radius: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--text-muted);
            font-size: 0.875rem;
        }

        .mock-content {
            flex: 1;
            height: 300px;
            background: var(--border);
            border-radius: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--text-muted);
            font-size: 0.875rem;
        }

        .mock-button {
            display: inline-block;
            padding: 0.75rem 1.5rem;
            background: var(--primary);
            color: white;
            border: none;
            border-radius: 0.5rem;
            font-size: 0.875rem;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }

        .mock-button:hover {
            background: var(--primary-hover);
        }

        .mock-input {
            width: 100%;
            padding: 0.75rem;
            border: 2px solid var(--border);
            border-radius: 0.5rem;
            font-size: 0.875rem;
            background: var(--card-bg);
        }

        .placeholder {
            width: 100%;
            height: 200px;
            background: var(--border);
            border-radius: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--text-muted);
            font-size: 0.875rem;
        }

        /* 通用样式 */
        .section {
            margin-bottom: 2rem;
        }

        .label {
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--text-muted);
            font-weight: 600;
        }

        /* 选择指示器 */
        #selection-indicator {
            position: fixed;
            top: 1rem;
            right: 1rem;
            background: var(--card-bg);
            border: 2px solid var(--primary);
            border-radius: 0.5rem;
            padding: 0.75rem 1rem;
            box-shadow: var(--shadow);
            font-size: 0.875rem;
            font-weight: 500;
            z-index: 1000;
            display: none;
        }

        #selection-indicator.visible {
            display: block;
        }

        /* 点击动画 */
        @keyframes click-flash {
            0% { opacity: 1; }
            50% { opacity: 0.7; }
            100% { opacity: 1; }
        }

        .option:active, .card:active {
            animation: click-flash 0.2s ease;
        }
    </style>
</head>
<body>
    <div class="container">
        <div id="content">
            {{CONTENT}}
        </div>
    </div>

    <div id="selection-indicator">已选择：<span id="selection-count">0</span> 个选项</div>

    <script>
        // 切换选择状态
        function toggleSelect(element) {
            const container = element.closest('.options') || element.closest('.cards');
            const isMulti = container && container.dataset.multiselect !== undefined;
            
            const wasSelected = element.classList.contains('selected');
            
            if (!isMulti) {
                // 单选：清除其他选项
                const siblings = container.querySelectorAll('.option, .card');
                siblings.forEach(sib => sib.classList.remove('selected'));
            }
            
            // 切换当前选项
            element.classList.toggle('selected', !wasSelected);
            
            // 更新指示器
            updateIndicator();
            
            // 记录事件
            recordEvent(element);
        }

        // 更新选择指示器
        function updateIndicator() {
            const selected = document.querySelectorAll('.option.selected, .card.selected');
            const indicator = document.getElementById('selection-indicator');
            const count = document.getElementById('selection-count');
            
            if (selected.length > 0) {
                count.textContent = selected.length;
                indicator.classList.add('visible');
            } else {
                indicator.classList.remove('visible');
            }
        }

        // 记录点击事件
        function recordEvent(element) {
            const choice = element.dataset.choice || '';
            const text = element.querySelector('h3')?.textContent || 
                        element.querySelector('p')?.textContent || 
                        element.textContent.trim().split('\n')[0] || '';
            
            const event = {
                type: 'click',
                choice: choice,
                text: text,
                timestamp: Math.floor(Date.now() / 1000)
            };
            
            fetch('/event', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(event)
            }).catch(err => console.error('记录事件失败:', err));
        }

        // 页面加载时检查是否有已选择的选项
        document.addEventListener('DOMContentLoaded', () => {
            updateIndicator();
        });
    </script>
</body>
</html>`

// NewServer 创建新服务器
func NewServer(port int, screenDir, stateDir string) *Server {
	return &Server{
		port:      port,
		screenDir: screenDir,
		stateDir:  stateDir,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 创建文件监听器
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败：%w", err)
	}

	// 添加监听目录
	if err := s.watcher.Add(s.screenDir); err != nil {
		return fmt.Errorf("添加监听目录失败：%w", err)
	}

	// 设置 HTTP 路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/event", s.handleEvent)
	mux.HandleFunc("/latest", s.handleLatest)

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// 启动文件监听协程
	s.wg.Add(1)
	go s.watchFiles()

	// 启动 HTTP 服务器（非阻塞）
	go func() {
		defer s.wg.Done()
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP 服务器错误：%v\n", err)
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	// 关闭文件监听器
	if s.watcher != nil {
		s.watcher.Close()
	}

	// 关闭 HTTP 服务器
	if s.httpServer != nil {
		return s.httpServer.Close()
	}

	return nil
}

// Wait 等待服务器退出
func (s *Server) Wait() {
	s.wg.Wait()
}

// watchFiles 监听文件变化
func (s *Server) watchFiles() {
	defer s.wg.Done()

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			// 只处理 HTML 文件
			if !strings.HasSuffix(event.Name, ".html") {
				continue
			}

			// 文件创建或修改时更新最新文件
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				s.mu.Lock()
				s.latestFile = event.Name
				s.mu.Unlock()

				// 清空事件文件（新屏幕时）
				if event.Op&fsnotify.Create == fsnotify.Create {
					s.clearEvents()
				}

				fmt.Printf("[%s] 检测到文件变化：%s\n", time.Now().Format("15:04:05"), filepath.Base(event.Name))
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "文件监听错误：%v\n", err)

		case <-s.stopChan:
			return
		}
	}
}

// handleIndex 处理根路径请求
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 获取最新 HTML 文件
	latestFile, err := s.getLatestHTML()
	if err != nil {
		// 如果没有文件，显示欢迎页面
		s.serveWelcome(w)
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(latestFile)
	if err != nil {
		http.Error(w, "读取文件失败", http.StatusInternalServerError)
		return
	}

	// 判断是否是完整 HTML 文档
	contentStr := string(content)
	isFullDoc := strings.HasPrefix(strings.TrimSpace(contentStr), "<!DOCTYPE") ||
		strings.HasPrefix(strings.TrimSpace(contentStr), "<html")

	var html string
	if isFullDoc {
		// 完整文档：直接注入脚本
		html = s.injectScript(contentStr)
	} else {
		// 内容片段：使用框架模板包装
		html = strings.Replace(FrameTemplate, "{{CONTENT}}", contentStr, 1)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// handleEvent 处理事件记录
func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var event map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "无效的 JSON", http.StatusBadRequest)
		return
	}

	// 添加服务器时间戳
	event["server_time"] = time.Now().Format(time.RFC3339)

	// 写入事件文件
	if err := s.appendEvent(event); err != nil {
		http.Error(w, "记录事件失败", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleLatest 获取最新屏幕信息
func (s *Server) handleLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	latestFile, err := s.getLatestHTML()
	if err != nil {
		http.Error(w, "暂无屏幕内容", http.StatusNotFound)
		return
	}

	info := map[string]string{
		"file":      latestFile,
		"filename":  filepath.Base(latestFile),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// getLatestHTML 获取最新的 HTML 文件
func (s *Server) getLatestHTML() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 如果已有最新文件记录，直接返回
	if s.latestFile != "" {
		return s.latestFile, nil
	}

	// 否则扫描目录查找最新文件
	entries, err := os.ReadDir(s.screenDir)
	if err != nil {
		return "", err
	}

	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		filePath := filepath.Join(s.screenDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if latestFile == "" || info.ModTime().After(latestTime) {
			latestFile = filePath
			latestTime = info.ModTime()
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("未找到 HTML 文件")
	}

	return latestFile, nil
}

// serveWelcome 显示欢迎页面
func (s *Server) serveWelcome(w http.ResponseWriter) {
	welcomeHTML := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Visual Choice - 欢迎</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f8fafc;
            color: #1e293b;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            margin: 0;
            padding: 2rem;
        }
        .welcome {
            text-align: center;
            max-width: 600px;
        }
        h1 {
            font-size: 2.5rem;
            margin-bottom: 1rem;
            color: #2563eb;
        }
        p {
            font-size: 1.125rem;
            color: #64748b;
            line-height: 1.8;
        }
        .info {
            background: white;
            border: 2px solid #e2e8f0;
            border-radius: 0.75rem;
            padding: 1.5rem;
            margin-top: 2rem;
            text-align: left;
        }
        .info h3 {
            margin-top: 0;
            color: #1e293b;
        }
        .info code {
            background: #f1f5f9;
            padding: 0.25rem 0.5rem;
            border-radius: 0.25rem;
            font-size: 0.875rem;
        }
    </style>
</head>
<body>
    <div class="welcome">
        <h1>🎨 Visual Choice</h1>
        <p>可视化选择工具已启动</p>
        <div class="info">
            <h3>使用说明</h3>
            <p>1. 将 HTML 文件写入目录：<code>` + s.screenDir + `</code></p>
            <p>2. 服务器会自动提供最新的 HTML 文件</p>
            <p>3. 用户点击会被记录到：<code>` + s.stateDir + `/events.jsonl</code></p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(welcomeHTML))
}

// injectScript 向完整 HTML 文档注入脚本
func (s *Server) injectScript(html string) string {
	// 在 </body> 前注入脚本
	script := `
<script>
    function toggleSelect(element) {
        const container = element.closest('.options') || element.closest('.cards');
        const isMulti = container && container.dataset.multiselect !== undefined;
        
        const wasSelected = element.classList.contains('selected');
        
        if (!isMulti) {
            const siblings = container.querySelectorAll('.option, .card');
            siblings.forEach(sib => sib.classList.remove('selected'));
        }
        
        element.classList.toggle('selected', !wasSelected);
        updateIndicator();
        recordEvent(element);
    }

    function updateIndicator() {
        const selected = document.querySelectorAll('.option.selected, .card.selected');
        let indicator = document.getElementById('selection-indicator');
        
        if (!indicator) {
            indicator = document.createElement('div');
            indicator.id = 'selection-indicator';
            indicator.innerHTML = '已选择：<span id="selection-count">0</span> 个选项';
            indicator.style.cssText = 'position:fixed;top:1rem;right:1rem;background:white;border:2px solid #2563eb;border-radius:0.5rem;padding:0.75rem 1rem;box-shadow:0 4px 6px -1px rgb(0 0 0 / 0.1);font-size:0.875rem;font-weight:500;z-index:1000;';
            document.body.appendChild(indicator);
        }
        
        const count = indicator.querySelector('#selection-count');
        count.textContent = selected.length;
        
        if (selected.length === 0) {
            indicator.style.display = 'none';
        } else {
            indicator.style.display = 'block';
        }
    }

    function recordEvent(element) {
        const choice = element.dataset.choice || '';
        const text = element.querySelector('h3')?.textContent || 
                    element.querySelector('p')?.textContent || 
                    element.textContent.trim().split('\\n')[0] || '';
        
        const event = {
            type: 'click',
            choice: choice,
            text: text,
            timestamp: Math.floor(Date.now() / 1000)
        };
        
        fetch('/event', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(event)
        }).catch(err => console.error('记录事件失败:', err));
    }

    document.addEventListener('DOMContentLoaded', () => {
        updateIndicator();
    });
</script>
`

	html = strings.Replace(html, "</body>", script+"</body>", 1)
	return html
}

// appendEvent 追加事件到文件
func (s *Server) appendEvent(event map[string]interface{}) error {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")

	// 转换为 JSON
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 追加写入
	f, err := os.OpenFile(eventsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}

// clearEvents 清空事件文件
func (s *Server) clearEvents() {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")
	os.WriteFile(eventsFile, []byte{}, 0644)
}

// getAllFiles 获取所有 HTML 文件（按时间排序）
func (s *Server) getAllFiles() ([]string, error) {
	entries, err := os.ReadDir(s.screenDir)
	if err != nil {
		return nil, err
	}

	type fileTime struct {
		path string
		time time.Time
	}

	var files []fileTime
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		filePath := filepath.Join(s.screenDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, fileTime{filePath, info.ModTime()})
	}

	// 按时间排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		return files[i].time.After(files[j].time)
	})

	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.path
	}

	return result, nil
}
