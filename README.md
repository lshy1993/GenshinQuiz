# GenshinQuiz - Genshin Impact Quiz Application
# GenshinQuiz 原神知识问答

A knowledge quiz application based on Genshin Impact, featuring a React frontend, Go API backend, and PostgreSQL database with OpenAPI specifications.

一个基于原神的知识问答应用，包含前端React应用、后端Go API服务和PostgreSQL数据库，使用OpenAPI规范。

## 🎯 Project Overview / 项目简介

GenshinQuiz is a full-stack quiz application that allows players to test their knowledge of the Genshin Impact game. The project uses modern technology stack with Go backend and supports Docker containerized deployment.

GenshinQuiz 是一个全栈知识问答应用，让玩家可以测试自己对原神游戏的了解程度。项目采用现代化的技术栈和Go后端，支持Docker容器化部署。

## 🛠️ Tech Stack / 技术栈

### Frontend / 前端
- **React 18** - UI Framework / UI框架
- **TypeScript** - Type Safety / 类型安全
- **Vite** - Build Tool / 构建工具
- **Material-UI (MUI)** - UI Component Library / UI组件库
- **Axios** - HTTP Client / HTTP客户端
- **Bun** - Package Manager and Runtime / 包管理器和运行时
- **Orval** - OpenAPI Code Generator / OpenAPI代码生成器

### Backend / 后端
- **Go 1.21** - Programming Language / 编程语言
- **Go-Chi** - Web Framework / Web框架
- **PostgreSQL** - Database / 数据库
- **Go-Jet** - Type-safe SQL Query Builder / 类型安全的SQL查询构建器
- **Goose** - Database Migration Tool / 数据库迁移工具
- **OpenAPI 3.0** - API Specification / API规范
- **oapi-codegen** - OpenAPI Code Generator / OpenAPI代码生成器

### Development Tools / 开发工具
- **Docker & Docker Compose** - Containerization / 容器化
- **Biome** - Code Formatting and Linting / 代码格式化和Lint
- **Orval** - API Client Generation / API客户端生成
- **ESLint** - Code Quality Check / 代码质量检查

## 📁 Project Structure / 项目结构

```
genshin-quiz/
├── backend/                # Go Backend API / Go后端API
│   ├── api/               # OpenAPI Specifications / OpenAPI规范
│   ├── internal/          # Private Application Code / 私有应用代码
│   │   ├── config/       # Configuration / 配置
│   │   ├── database/     # Database Layer / 数据库层
│   │   ├── handlers/     # HTTP Handlers / HTTP处理器
│   │   ├── models/       # Data Models / 数据模型
│   │   ├── repository/   # Data Access Layer / 数据访问层
│   │   └── services/     # Business Logic / 业务逻辑
│   ├── migrations/       # Database Migrations / 数据库迁移
│   ├── scripts/          # Utility Scripts / 实用脚本
│   ├── main.go           # Application Entry / 应用入口
│   ├── go.mod            # Go Module Definition / Go模块定义
│   └── Dockerfile        # Backend Docker Config / 后端Docker配置
├── frontend/             # Frontend Application / 前端应用
│   ├── src/              # Source Code / 源代码
│   │   ├── api/         # Generated API Client / 生成的API客户端
│   │   ├── assets/      # Static Assets / 静态资源
│   │   ├── App.tsx      # Main App Component / 主应用组件
│   │   └── main.tsx     # Application Entry / 应用入口
│   ├── openapi/         # OpenAPI Files / OpenAPI文件
│   ├── package.json     # Frontend Dependencies / 前端依赖
│   └── Dockerfile       # Frontend Docker Config / 前端Docker配置
├── scripts/             # Project Scripts / 项目脚本
├── docker-compose.yml   # Docker Compose Config / Docker编排配置
└── biome.json          # Code Formatting Config / 代码格式化配置
```

## 🚀 Getting Started / 快速开始

### Prerequisites / 前提条件

- [Docker](https://www.docker.com/) and Docker Compose / 和 Docker Compose
- [Go](https://golang.org/) 1.21+ (for local development / 本地开发用)
- [Node.js](https://nodejs.org/) 18+ (for frontend / 前端用)
- [Bun](https://bun.sh/) (recommended package manager / 推荐的包管理器)

### Start with Docker / 使用Docker启动 (Recommended / 推荐)

1. **Clone the project / 克隆项目**

## 🚀 Getting Started / 快速开始

### Prerequisites / 前提条件

- [Docker](https://www.docker.com/) and Docker Compose / 和 Docker Compose
- [Node.js](https://nodejs.org/) 18+ (for local development / 本地开发用)
- [Bun](https://bun.sh/) (recommended package manager / 推荐的包管理器)

### Start with Docker / 使用Docker启动 (Recommended / 推荐)

1. **Clone the project / 克隆项目**
   ```bash
   git clone https://github.com/lshy1993/GenshinQuiz.git
   cd GenshinQuiz
   ```

2. **Start all services / 启动所有服务**
   ```bash
   # Option 1: Use the start script / 选项1：使用启动脚本
   ./scripts/start.sh
   
   # Option 2: Use docker-compose directly / 选项2：直接使用docker-compose
   docker-compose up --build -d
   ```

3. **Access the application / 访问应用**
   - Frontend / 前端应用: http://localhost:3000
   - Go Backend API / Go后端API: http://localhost:8080
   - API Health Check / API健康检查: http://localhost:8080/health
   - Swagger Editor / Swagger编辑器: http://localhost:3020

### Local Development / 本地开发

#### Backend Development / 后端开发

```bash
cd backend

# Install Go dependencies / 安装Go依赖
go mod tidy

# Set up development environment / 设置开发环境
./scripts/setup.sh

# Start PostgreSQL / 启动PostgreSQL
docker-compose up postgres -d

# Run database migrations / 运行数据库迁移
./scripts/migrate.sh up

# Generate models and API code / 生成模型和API代码
./scripts/generate_models.sh
./scripts/generate_api.sh

# Start the Go server / 启动Go服务器
go run main.go
```

Backend server will start at http://localhost:3001

后端服务将在 http://localhost:3001 启动

#### Frontend Development / 前端开发

```bash
cd frontend
bun install
bun run dev
```

Frontend application will start at http://localhost:3000

前端应用将在 http://localhost:3000 启动

## 🔧 Development Commands / 开发命令

### Frontend Commands / 前端命令
```bash
bun run dev          # Start development server / 启动开发服务器
bun run build        # Build for production / 构建生产版本
bun run format       # Format code / 格式化代码
bun run lint         # Lint code / 代码检查
bun run orval        # Generate API client / 生成API客户端
```

### Backend Commands / 后端命令
```bash
npm start            # Start server / 启动服务器
npm run dev          # Start development server with nodemon / 使用nodemon启动开发服务器
npm run db:migrate   # Run database migrations / 运行数据库迁移
npm run db:rollback  # Rollback last migration / 回滚上一次迁移
npm run db:seed      # Run database seeds / 运行数据库种子
npm run db:reset     # Reset database (rollback + migrate + seed) / 重置数据库
```

## 📊 Database / 数据库

The project uses PostgreSQL database with Knex.js for schema management and migrations:

项目使用PostgreSQL数据库，通过Knex.js进行架构管理和迁移：

### Database Schema / 数据库架构
- `users` table - Store user information / 存储用户信息
  - id, name, email, created_at, updated_at
- `quizzes` table - Store quiz questions and answers / 存储问答题目
  - id, question, answer, category, difficulty, explanation, created_at, updated_at

### Database Management / 数据库管理
```bash
# Run migrations / 运行迁移
npm run db:migrate

# Seed database with sample data / 用示例数据填充数据库
npm run db:seed

# Reset database completely / 完全重置数据库
npm run db:reset
```

The database migrations and seeds are located in `backend/db/` directory.

数据库迁移和种子文件位于 `backend/db/` 目录中。

## 🐳 Docker Deployment / Docker部署

The project fully supports Docker containerized deployment:

项目完全支持Docker容器化部署：

```bash
# Start all services / 启动所有服务
docker-compose up -d

# Check running status / 查看运行状态
docker-compose ps

# View logs / 查看日志
docker-compose logs

# Stop services / 停止服务
docker-compose down
```

## 📝 API Documentation / API文档

API documentation is automatically generated through Swagger and can be accessed at:

API文档通过Swagger自动生成，可以在以下地址访问：

- Development Environment / 开发环境: http://localhost:3020 (Swagger Editor)
- API Documentation / API文档: View `frontend/openapi/api-docs.yaml` / 查看 `frontend/openapi/api-docs.yaml`

## 🤝 Contributing / 贡献指南

1. Fork this project / Fork 本项目
2. Create a feature branch / 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. Commit your changes / 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch / 推送到分支 (`git push origin feature/AmazingFeature`)
5. Create a Pull Request / 创建Pull Request

## 🔄 Code Standards / 代码规范

The project uses Biome for code formatting and linting:

项目使用Biome进行代码格式化和检查：

```bash
# Format code / 格式化代码
bun run format

# Lint code / 代码检查
bun run lint
```

## 🛡️ License / 许可证

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

本项目采用MIT许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 🎮 About Genshin Impact / 关于原神

Genshin Impact is an open-world adventure game developed by miHoYo. This project is for educational and entertainment purposes only and is not affiliated with the official game.

原神是miHoYo开发的开放世界冒险游戏。本项目纯属学习和娱乐目的，与官方无关。

## 📞 Contact / 联系方式

- Project Maintainer / 项目维护者: lshy1993
- GitHub: https://github.com/lshy1993/GenshinQuiz

---

**May your journey through the knowledge of Teyvat be filled with discoveries!** 🌟

**愿你在提瓦特大陆的知识之旅中收获满满！** 🌟
