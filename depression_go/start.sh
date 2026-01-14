#!/bin/bash

# 抑郁倾向检测系统启动脚本

echo "=========================================="
echo "    抑郁倾向检测系统 - Go后端启动脚本"
echo "=========================================="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go 1.21+"
    exit 1
fi

echo "✅ Go环境检查通过"

# 检查配置文件
if [ ! -f "config.env" ]; then
    echo "❌ 错误: 未找到config.env配置文件"
    echo "请复制config.env.example为config.env并配置相关参数"
    exit 1
fi

echo "✅ 配置文件检查通过"

# 检查数据库连接
echo "🔍 检查数据库连接..."
# 这里可以添加数据库连接检查逻辑

# 安装依赖
echo "📦 安装Go依赖..."
go mod tidy

# 创建uploads目录
if [ ! -d "uploads" ]; then
    echo "📁 创建uploads目录..."
    mkdir -p uploads
fi

# 设置环境变量
export GIN_MODE=debug

# 启动服务
echo "🚀 启动服务..."
echo "服务将在 http://localhost:8088 启动"
echo "按 Ctrl+C 停止服务"
echo "=========================================="

go run main.go 