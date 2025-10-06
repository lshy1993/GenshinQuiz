# Genshin Quiz Backend - Microservices Architecture

这是一个基于 Go 的现代微服务架构后端，用于 Genshin Impact 问答应用。

## 🔧 技术栈

### 核心技术
- **编程语言**: Go 1.25.0
- **架构模式**: 微服务架构 (5个独立服务)

### Web 框架与 API
- **HTTP 路由**: Chi Router (github.com/go-chi/chi/v5)
- **API 规范**: OpenAPI/Swagger (使用 oapi-codegen 生成代码)
- **认证**: JWT (github.com/go-chi/jwtauth/v5, github.com/golang-jwt/jwt/v5)

### 数据存储
- **主数据库**: PostgreSQL 16.4
- **数据库访问**: Jet ORM (github.com/go-jet/jet/v2)
- **数据库驱动**: github.com/lib/pq

### 服务架构组件
项目包含5个主要服务组件：
1. **Web Server (server)** - 主要的 HTTP API 服务
2. **Worker (worker)** - 后台任务处理器
3. **Cron Job (cronjob)** - 定时任务服务
4. **Console (console)** - 命令行工具
5. **DB Migration** - 数据库迁移服务

### 任务队列与调度
- **异步任务**: Asynq (github.com/hibiken/asynq)
- **任务监控**: Asynqmon (github.com/hibiken/asynqmon)

### 云服务集成
- **Azure Storage**: Azure Blob Storage (github.com/Azure/azure-sdk-for-go)
- **容器化**: Docker + Docker Compose
- **部署**: Kubernetes (AKS) with Helm Charts

### 监控与日志
- **错误追踪**: Sentry (github.com/getsentry/sentry-go)
- **日志**: Zap (go.uber.org/zap)

### 开发工具
- **构建工具**: Task (Taskfile.yaml)
- **测试**:
  - github.com/stretchr/testify (单元测试)
  - github.com/DATA-DOG/go-txdb (事务测试)
  - github.com/go-faker/faker/v4 (测试数据生成)
- **数据验证**: github.com/gookit/validate

## 🚀 快速开始

### 前置条件
- Go 1.25.0+
- Docker & Docker Compose
- Task (https://taskfile.dev/)

### 环境配置

1. 复制环境变量配置：
```bash
cp .env.example .env
```

2. 修改 `.env` 文件中的配置（数据库密码、JWT 密钥等）

### 本地开发

使用 Docker Compose 启动完整的开发环境：

```bash
# 启动所有服务 (PostgreSQL, Redis, 所有微服务)
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f server
docker-compose logs -f worker
```

### 使用 Task 进行开发

```bash
# 安装依赖
task deps

# 本地运行单个服务
task dev:server   # 启动 API 服务器
task dev:worker   # 启动后台任务处理器

# 数据库迁移
task migrate:up   # 应用迁移
task migrate:down # 回滚迁移

# 构建所有服务
task build

# 运行测试
task test
task test:unit
task test:integration

# 代码格式化和检查
task fmt
task lint
```

## 📁 项目结构

```
backend/
├── cmd/                    # 服务入口点
│   ├── server/            # API 服务器
│   ├── worker/            # 后台任务处理器
│   ├── cronjob/           # 定时任务服务
│   ├── console/           # 命令行工具
│   └── migration/         # 数据库迁移工具
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── handlers/         # HTTP 处理器
│   ├── services/         # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── models/           # 数据模型
│   ├── middleware/       # 中间件
│   ├── tasks/            # 异步任务定义
│   ├── infrastructure/   # 基础设施层
│   ├── validation/       # 数据验证
│   ├── azure/            # Azure 集成
│   ├── cron/             # 定时任务调度
│   ├── console/          # 控制台命令
│   └── migration/        # 迁移管理
├── deployments/          # 部署配置
│   ├── docker/           # Docker 配置
│   └── helm/             # Kubernetes Helm Charts
├── migrations/           # 数据库迁移文件
├── scripts/              # 脚本文件
├── Dockerfile           # 多阶段 Docker 构建
├── docker-compose.yml   # 开发环境编排
└── Taskfile.yaml        # 任务定义
```

## 🔧 服务详情

### 1. Server (API 服务)
- 端口: 8080
- 功能: REST API, JWT 认证, 用户管理, 问答功能
- 健康检查: `GET /health`

### 2. Worker (后台任务处理器)
- 功能: 邮件发送, 文件处理, 数据分析
- 队列: critical, default, low
- 监控: Asynqmon (端口 8081)

### 3. Cronjob (定时任务)
- 功能: 数据清理, 统计报告, 系统维护
- 调度: 每日、每周、每小时任务

### 4. Console (命令行工具)
```bash
# 创建用户
go run ./cmd/console -command=user create user@example.com

# 数据种子
go run ./cmd/console -command=seed

# 系统统计
go run ./cmd/console -command=stats
```

### 5. Migration (数据库迁移)
```bash
# 应用迁移
go run ./cmd/migration -action=up

# 创建新迁移
go run ./cmd/migration -action=create -name=add_new_table
```

## 🐳 Docker 部署

### 构建特定服务
```bash
# 构建 API 服务器
docker build --build-arg SERVICE=server -t genshin-quiz-server .

# 构建工作器
docker build --build-arg SERVICE=worker -t genshin-quiz-worker .
```

### 生产环境部署
```bash
# 使用 Task 构建 Docker 镜像
task docker:build

# 部署到 Kubernetes
task deploy:prod
```

## ☸️ Kubernetes 部署

```bash
# 安装 Helm Chart
helm install genshin-quiz ./deployments/helm/genshin-quiz

# 升级部署
helm upgrade genshin-quiz ./deployments/helm/genshin-quiz

# 查看状态
kubectl get pods -l app=genshin-quiz
```

## 📊 监控和调试

### Asynq 监控
访问 http://localhost:8081 查看任务队列状态

### 日志查看
```bash
# Docker Compose 环境
docker-compose logs -f server
docker-compose logs -f worker

# Kubernetes 环境
kubectl logs -f deployment/genshin-quiz-server
kubectl logs -f deployment/genshin-quiz-worker
```

### 健康检查
```bash
curl http://localhost:8080/health
```

## 🧪 测试

```bash
# 运行所有测试
task test

# 运行单元测试
task test:unit

# 运行集成测试
task test:integration

# 测试覆盖率
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔒 安全性

- 使用非 root 用户运行容器
- JWT 令牌认证
- 输入验证和清理
- SQL 注入防护 (通过 ORM)
- HTTPS 支持 (生产环境)

## 🚀 性能优化

- 数据库连接池配置
- Redis 缓存
- 异步任务处理
- 水平扩展支持
- 健康检查和自动重启

## 📝 API 文档

API 文档通过 OpenAPI/Swagger 生成，可在以下位置访问：
- 开发环境: http://localhost:3020 (Swagger Editor)
- 生产环境: 集成到应用中

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

此项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。