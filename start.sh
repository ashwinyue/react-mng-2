#!/bin/bash

echo "🚀 启动前后端开发环境..."

# 检查是否安装了 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 未检测到 Node.js，请先安装 Node.js"
    exit 1
fi

# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo "❌ 未检测到 Go，请先安装 Go"
    exit 1
fi

# 启动后端
echo "📦 启动后端服务..."
cd backend
go mod download
go run main.go &
BACKEND_PID=$!
cd ..

# 等待后端启动
sleep 3

# 启动前端
echo "🎨 启动前端服务..."
npm install
npm run dev &
FRONTEND_PID=$!

echo ""
echo "✅ 服务启动成功！"
echo "📱 前端地址: http://localhost:5173"
echo "🔧 后端地址: http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待用户中断
wait

# 清理进程
kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
