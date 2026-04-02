package server

import (
	"context"
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

const (
	// 请求体最大大小限制 (10MB)
	MaxRequestBodySize = 10 << 20

	// 默认超时配置
	DefaultReadTimeout  = 15 * time.Second
	DefaultWriteTimeout = 15 * time.Second
	DefaultIdleTimeout  = 60 * time.Second
)

type Server struct {
	port       int
	screenDir  string
	stateDir   string
	httpServer *http.Server
	watcher    *fsnotify.Watcher
	wg         sync.WaitGroup
	stopChan   chan struct{}
	mu         sync.RWMutex
	latestFile string
}

func NewServer(port int, screenDir, stateDir string) *Server {
	return &Server{
		port:      port,
		screenDir: screenDir,
		stateDir:  stateDir,
		stopChan:  make(chan struct{}),
	}
}

func (s *Server) Start() error {
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败：%w", err)
	}

	if err := s.watcher.Add(s.screenDir); err != nil {
		return fmt.Errorf("添加监听目录失败：%w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/event", s.handleEvent)
	mux.HandleFunc("/latest", s.handleLatest)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	s.wg.Add(1)
	go s.watchFiles()

	go func() {
		defer s.wg.Done()
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP 服务器错误：%v\n", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.watcher != nil {
		s.watcher.Close()
	}

	close(s.stopChan)

	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) watchFiles() {
	defer s.wg.Done()

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			if !strings.HasSuffix(event.Name, ".html") {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				s.mu.Lock()
				s.latestFile = event.Name
				s.mu.Unlock()

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

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	latestFile, err := s.getLatestHTML()
	if err != nil {
		s.serveWelcome(w)
		return
	}

	cleanPath := filepath.Clean(latestFile)
	cleanScreenDir := filepath.Clean(s.screenDir)
	if !strings.HasPrefix(cleanPath, cleanScreenDir) {
		http.Error(w, "非法的文件路径", http.StatusForbidden)
		return
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		http.Error(w, "读取文件失败", http.StatusInternalServerError)
		return
	}

	contentStr := string(content)
	isFullDoc := strings.HasPrefix(strings.TrimSpace(contentStr), "<!DOCTYPE") ||
		strings.HasPrefix(strings.TrimSpace(contentStr), "<html")

	var html string
	if isFullDoc {
		html = s.injectScript(contentStr)
	} else {
		html = strings.Replace(FrameTemplate, "{{CONTENT}}", contentStr, 1)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var event map[string]interface{}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxRequestBodySize))
	if err := decoder.Decode(&event); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "请求体过大", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "无效的 JSON", http.StatusBadRequest)
		}
		return
	}

	if eventType, ok := event["type"].(string); !ok || eventType == "" {
		event["type"] = "unknown"
	}

	if _, ok := event["choice"].(string); !ok {
		event["choice"] = ""
	}

	if text, ok := event["text"].(string); ok && len(text) > 1000 {
		event["text"] = text[:1000] + "..."
	}

	event["server_time"] = time.Now().Format(time.RFC3339)

	if err := s.appendEvent(event); err != nil {
		http.Error(w, "记录事件失败", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

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

func (s *Server) getLatestHTML() (string, error) {
	s.mu.RLock()
	latest := s.latestFile
	s.mu.RUnlock()

	if latest != "" {
		return latest, nil
	}

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

func (s *Server) injectScript(html string) string {
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

func (s *Server) appendEvent(event map[string]interface{}) error {
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

func (s *Server) clearEvents() {
	eventsFile := filepath.Join(s.stateDir, "events.jsonl")
	os.WriteFile(eventsFile, []byte{}, 0600)
}

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

	sort.Slice(files, func(i, j int) bool {
		return files[i].time.After(files[j].time)
	})

	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.path
	}

	return result, nil
}
