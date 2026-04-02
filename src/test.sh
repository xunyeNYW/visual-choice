#!/bin/bash
# Visual Choice 测试脚本

set -e

echo "=== Visual Choice 测试脚本 ==="
echo ""

# 清理之前的测试
if [ -d "./test-session" ]; then
    echo "清理之前的测试..."
    ./visual-choice stop --dir ./test-session 2>/dev/null || true
    rm -rf ./test-session
    sleep 1
fi

# 1. 启动服务器
echo "1. 启动服务器..."
./visual-choice start --port 5234 --dir ./test-session &
SERVER_PID=$!
sleep 2

# 检查服务器是否启动
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "❌ 服务器启动失败"
    exit 1
fi
echo "✓ 服务器已启动 (PID: $SERVER_PID)"
echo ""

# 2. 测试根路径
echo "2. 测试根路径（欢迎页面）..."
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:5234/)
if [ "$RESPONSE" != "200" ]; then
    echo "❌ 根路径访问失败 (HTTP $RESPONSE)"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo "✓ 根路径访问正常"
echo ""

# 3. 写入测试 HTML
echo "3. 写入测试 HTML 文件..."
cp example.html test-session/screens/
sleep 1
echo "✓ 文件已写入"
echo ""

# 4. 测试 HTML 服务
echo "4. 测试 HTML 内容服务..."
RESPONSE=$(curl -s http://localhost:5234/)
if echo "$RESPONSE" | grep -q "哪个布局更适合"; then
    echo "✓ HTML 内容服务正常"
else
    echo "❌ HTML 内容服务失败"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# 5. 测试事件 API
echo "5. 测试事件记录 API..."
curl -s -X POST http://localhost:5234/event \
    -H "Content-Type: application/json" \
    -d '{"type":"click","choice":"test-a","text":"测试选项 A","timestamp":1706000001}' > /dev/null
curl -s -X POST http://localhost:5234/event \
    -H "Content-Type: application/json" \
    -d '{"type":"click","choice":"test-b","text":"测试选项 B","timestamp":1706000002}' > /dev/null

if [ -f "test-session/state/events.jsonl" ]; then
    echo "✓ 事件记录文件已创建"
    EVENT_COUNT=$(wc -l < test-session/state/events.jsonl)
    echo "✓ 已记录 $EVENT_COUNT 个事件"
else
    echo "❌ 事件记录文件未创建"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# 6. 测试 /latest 端点
echo "6. 测试 /latest 端点..."
LATEST=$(curl -s http://localhost:5234/latest)
if echo "$LATEST" | grep -q "example.html"; then
    echo "✓ /latest 端点正常"
    echo "  最新文件：example.html"
else
    echo "❌ /latest 端点失败"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# 7. 测试 status 命令
echo "7. 测试 status 命令..."
STATUS=$(./visual-choice status --dir ./test-session)
if echo "$STATUS" | grep -q "服务器正在运行"; then
    echo "✓ status 命令正常"
else
    echo "❌ status 命令失败"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# 8. 测试 events 命令
echo "8. 测试 events 命令..."
EVENTS=$(./visual-choice events --dir ./test-session)
if echo "$EVENTS" | grep -q "测试选项"; then
    echo "✓ events 命令正常"
    echo "$EVENTS" | head -5
else
    echo "❌ events 命令失败"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# 9. 停止服务器
echo "9. 停止服务器..."
./visual-choice stop --dir ./test-session
sleep 1

STATUS=$(./visual-choice status --dir ./test-session)
if echo "$STATUS" | grep -q "服务器未运行"; then
    echo "✓ 服务器已停止"
else
    echo "❌ 服务器停止失败"
    exit 1
fi
echo ""

echo "=== 所有测试通过！ ==="
echo ""
echo "测试会话目录：./test-session"
echo "  - screens/: HTML 文件目录"
echo "  - state/: 状态文件目录"
echo ""
echo "清理测试：rm -rf ./test-session"
