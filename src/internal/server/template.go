package server

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
