# 开发指南 - Development Guide

## 🚀 快速开始

### 1. 启动完整应用
```bash
# 从项目根目录运行
./scripts/start.sh
```

这将启动：
- ✅ PostgreSQL 数据库 (端口 5432)
- ✅ Go 后端 API (端口 8080)
- ✅ React 前端 (端口 3000)  
- ✅ Swagger 编辑器 (端口 3020)

### 2. 验证 API
```bash
# 测试 API 是否正常工作
./scripts/test_api.sh
```

### 3. 访问应用
- **前端**: http://localhost:3000
- **API**: http://localhost:8080/api/v1
- **健康检查**: http://localhost:8080/health
- **Swagger**: http://localhost:3020

## 🛠️ 后端开发

### 本地开发环境
```bash
cd backend

# 安装依赖和工具
./scripts/setup.sh

# 启动数据库
docker-compose up postgres -d

# 运行迁移
./scripts/migrate.sh up

# 生成代码
./scripts/generate_models.sh
./scripts/generate_api.sh

# 启动开发服务器
go run main.go
```

### 数据库操作
```bash
# 查看迁移状态
./scripts/migrate.sh status

# 创建新迁移
./scripts/migrate.sh create add_new_feature

# 应用迁移
./scripts/migrate.sh up

# 回滚迁移
./scripts/migrate.sh down
```

### 代码生成
```bash
# 从数据库生成 Go-Jet 模型
./scripts/generate_models.sh

# 从 OpenAPI 规范生成 API 代码
./scripts/generate_api.sh
```

## 🎨 前端开发

### 本地开发
```bash
cd frontend

# 安装依赖
bun install

# 生成 API 客户端
bun run generate

# 启动开发服务器
bun run dev
```

### API 客户端生成
前端使用 Orval 从 OpenAPI 规范自动生成 TypeScript API 客户端：
```bash
# 生成 API 客户端代码
bunx orval --config orval.config.json
```

## 📊 常用命令

### Docker 操作
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
docker-compose logs -f frontend

# 停止服务
docker-compose down

# 重新构建
docker-compose up --build
```

### API 测试
```bash
# 健康检查
curl http://localhost:8080/health

# 获取用户列表
curl http://localhost:8080/api/v1/users

# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com"}'

# 获取测验列表
curl http://localhost:8080/api/v1/quizzes
```

## 🔧 开发工作流

### 1. 添加新的 API 端点
1. 更新 `backend/api/openapi.yaml`
2. 运行 `./scripts/generate_api.sh` 生成 Go 代码
3. 实现业务逻辑 (services, handlers)
4. 在前端运行 `bun run generate` 更新客户端代码

### 2. 数据库变更
1. 创建迁移: `./scripts/migrate.sh create description`
2. 编辑迁移文件: `backend/migrations/xxx_description.sql`
3. 应用迁移: `./scripts/migrate.sh up`
4. 重新生成模型: `./scripts/generate_models.sh`

### 3. 调试
- **Go 后端日志**: `docker-compose logs -f backend`
- **前端日志**: `docker-compose logs -f frontend`
- **数据库连接**: `docker-compose exec database psql -U postgres -d genshinquiz`

## 📁 重要文件

### 配置文件
- `docker-compose.yml` - 服务编排
- `backend/.env` - 后端环境变量
- `backend/api/openapi.yaml` - API 规范
- `frontend/orval.config.json` - API 客户端生成配置

### 脚本文件
- `scripts/start.sh` - 启动完整应用
- `scripts/test_api.sh` - API 测试
- `backend/scripts/setup.sh` - 后端环境设置
- `backend/scripts/migrate.sh` - 数据库迁移
- `backend/scripts/generate_*.sh` - 代码生成

## 🐛 故障排除

### 常见问题
1. **端口冲突**: 确保 3000、8080、5432 端口未被占用
2. **数据库连接失败**: 检查 PostgreSQL 是否启动
3. **代码生成失败**: 确保安装了必要的工具
4. **Docker 构建失败**: 尝试 `docker system prune` 清理缓存

### 重置环境
```bash
# 停止所有服务
docker-compose down -v

# 清理 Docker 资源
docker system prune -f

# 重新启动
./scripts/start.sh
```