#!/bin/bash

# API 测试脚本
# 用于验证 Go 后端 API 是否正常工作

set -e

API_BASE_URL="http://localhost:8080"

echo "🧪 测试 Genshin Quiz Go API"
echo "API Base URL: $API_BASE_URL"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local description=$4
    local data=$5

    echo -n "Testing $description... "
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE_URL$endpoint")
    else
        response=$(curl -s -w "%{http_code}" -X "$method" "$API_BASE_URL$endpoint")
    fi
    
    status_code="${response: -3}"
    body="${response%???}"
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ PASS${NC} (Status: $status_code)"
        if [ -n "$body" ] && [ "$body" != "null" ]; then
            echo "   Response: $(echo "$body" | head -c 100)..."
        fi
    else
        echo -e "${RED}❌ FAIL${NC} (Expected: $expected_status, Got: $status_code)"
        echo "   Response: $body"
    fi
    echo ""
}

# 等待 API 启动
echo "⏳ 等待 API 启动..."
for i in {1..30}; do
    if curl -s "$API_BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ API 已启动${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ API 启动超时${NC}"
        exit 1
    fi
    sleep 2
done

echo ""

# 开始测试
echo "🚀 开始 API 测试"
echo "================================"

# 1. 健康检查
test_endpoint "GET" "/health" 200 "Health Check"

# 2. 获取用户列表
test_endpoint "GET" "/api/v1/users" 200 "Get Users List"

# 3. 获取用户列表（带分页）
test_endpoint "GET" "/api/v1/users?limit=5&offset=0" 200 "Get Users with Pagination"

# 4. 创建用户
user_data='{"username":"test_user","email":"test@example.com","display_name":"Test User"}'
test_endpoint "POST" "/api/v1/users" 201 "Create User" "$user_data"

# 5. 获取测验列表
test_endpoint "GET" "/api/v1/quizzes" 200 "Get Quizzes List"

# 6. 获取测验列表（带筛选）
test_endpoint "GET" "/api/v1/quizzes?category=characters&difficulty=easy" 200 "Get Quizzes with Filters"

# 7. 获取特定测验
test_endpoint "GET" "/api/v1/quizzes/1" 200 "Get Quiz by ID"

# 8. 测试不存在的资源
test_endpoint "GET" "/api/v1/users/999999" 404 "Get Non-existent User"

# 9. 测试无效的请求
test_endpoint "POST" "/api/v1/users" 400 "Create User with Invalid Data" '{"invalid":"data"}'

echo "================================"
echo -e "${YELLOW}🎉 API 测试完成！${NC}"
echo ""
echo "💡 提示："
echo "  • 如果有测试失败，请检查 API 服务是否正常运行"
echo "  • 查看详细日志: docker-compose logs backend"
echo "  • 重启服务: docker-compose restart backend"