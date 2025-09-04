# 原神问答投票系统架构设计

## 系统概述
这是一个基于 Node.js + PostgreSQL 的原神主题问答和投票系统，包含用户管理、问答系统、投票系统和管理后台。

## 数据库设计

### 用户系统 (users)
- 用户基本信息：姓名、邮箱、密码
- 权限系统：user/admin 角色
- 状态管理：激活状态、最后登录时间

### 问答系统
- **quizzes**: 题目表，支持多种题型
- **quiz_categories**: 题目分类管理
- **quiz_attempts**: 用户答题记录

### 投票系统
- **votes**: 投票活动表
- **vote_options**: 投票选项表
- **user_votes**: 用户投票记录表

## 功能特性

### 用户功能
- ✅ 用户注册/登录/个人资料管理
- ✅ 参与问答：单选、多选、判断、填空题
- ✅ 参与投票：单选/多选投票
- ✅ 查看个人答题和投票历史
- ✅ 实时投票结果统计

### 管理员功能
- ✅ 题库管理：创建、编辑、删除题目
- ✅ 分类管理：管理题目分类
- ✅ 投票管理：创建投票活动，设置时间限制
- ✅ 用户管理：用户权限控制
- ✅ 数据统计：答题准确率、投票参与度等

### 系统特性
- ✅ 软删除：数据安全，支持恢复
- ✅ 时间控制：投票开始/结束时间管理
- ✅ 防重复：用户每个投票只能参与一次
- ✅ 匿名投票：支持匿名和实名投票模式
- ✅ 多环境：本地MySQL + Docker PostgreSQL

## API 设计建议

### 用户相关
```
POST /api/auth/login
POST /api/auth/register
GET /api/users/profile
PUT /api/users/profile
GET /api/users/history (答题和投票历史)
```

### 问答相关
```
GET /api/quizzes (支持筛选：分类、难度、类型)
GET /api/quizzes/random
GET /api/quizzes/:id
POST /api/quizzes/:id/attempt (提交答案)
GET /api/categories
```

### 投票相关
```
GET /api/votes (活跃投票列表)
GET /api/votes/:id (投票详情+选项)
POST /api/votes/:id/submit (提交投票)
GET /api/votes/:id/results (投票结果)
```

### 管理员相关
```
POST /api/admin/quizzes (创建题目)
PUT /api/admin/quizzes/:id
DELETE /api/admin/quizzes/:id
POST /api/admin/votes (创建投票)
GET /api/admin/stats (统计数据)
```

## 下一步开发建议

1. **认证系统**：实现 JWT 认证和权限中间件
2. **API 路由**：根据 OpenAPI 规范实现完整 API
3. **前端界面**：设计用户友好的问答和投票界面
4. **数据验证**：添加输入验证和错误处理
5. **性能优化**：数据库索引、查询优化
6. **安全加固**：密码加密、CSRF 保护等

## 运行方式
```bash
# 启动 Docker 环境
docker-compose up -d

# 运行数据库迁移
npm run db:migrate

# 运行种子数据
npm run db:seed
```

前端将在 http://localhost:3000 启动
后端 API 在 http://localhost:3082 启动
Swagger Editor 在 http://localhost:3020 启动
