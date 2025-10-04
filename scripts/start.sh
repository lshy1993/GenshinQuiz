#!/bin/bash

# Genshin Quiz 完整应用启动脚本
# 启动前端、Go 后端和 PostgreSQL 数据库

set -e

echo "🚀 启动 Genshin Quiz 应用"

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

# 检查是否在项目根目录
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

echo "📦 构建并启动所有服务..."

# 启动所有服务
docker-compose up --build -d

echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 服务状态："
docker-compose ps

echo ""
echo "🎉 应用启动完成！"
echo ""
echo "📱 服务地址："
echo "  • 前端应用:     http://localhost:3000"
echo "  • Go API:      http://localhost:8080"
echo "  • API 文档:     http://localhost:8080/health"
echo "  • Swagger:     http://localhost:3020"
echo "  • PostgreSQL:  localhost:5432"
echo ""
echo "🔧 有用的命令："
echo "  • 查看日志:     docker-compose logs -f"
echo "  • 停止服务:     docker-compose down"
echo "  • 重启服务:     docker-compose restart"
echo ""
echo "🐛 调试："
echo "  • Go API 日志:  docker-compose logs -f backend"
echo "  • 前端日志:     docker-compose logs -f frontend"
echo "  • 数据库日志:   docker-compose logs -f database"