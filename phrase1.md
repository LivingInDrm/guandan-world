以下是针对 **第 1 阶段「框架 & DevOps 搭建」** 的详细设计，拆分为结构设计、技术选型、关键任务和交付验收标准，可直接用于开发实施。

---

# 🧱 第 1 阶段：框架 & DevOps 搭建（预计 1 周）

## 🎯 目标

建立一个可运行的基础工程框架，支持：

* 本地一键运行全栈环境（含数据库、缓存、前后端）
* GitHub Actions 自动化流程：测试 → 构建 → 镜像推送
* 预部署到 staging 环境（即将来的测试服）

---

## 📦 1. 工程目录结构

```shell
guandan-project/
├── backend/               # Go 服务代码
├── frontend/              # React + Vite 前端
├── sdk/                   # 游戏核心逻辑，无 I/O 依赖
├── infra/                 # DevOps / 部署脚本 / Terraform 等
├── docker-compose.yml     # 本地多服务编排
└── README.md              # 工程说明
```

---

## 🧰 2. 技术选型

| 模块       | 技术栈                               | 说明               |
| -------- | --------------------------------- | ---------------- |
| 后端 API   | Go 1.22 + Gin + Gorilla WebSocket | 小巧、高性能，适合游戏实时性要求 |
| 前端 UI    | React + TypeScript + Vite         | 快速构建开发体验好        |
| 游戏逻辑 SDK | 纯 Go 模块                           | 无副作用、可单元测试、后期可移植 |
| 数据存储     | PostgreSQL + Redis                | 用户数据、房间状态 / 缓存等  |
| 身份认证     | JWT + bcrypt                      | 简单易用、安全可靠        |
| DevOps   | GitHub Actions + Docker + Compose | 一键部署与 CI/CD      |

---

## 🏗️ 3. 关键任务拆解

### ✅ 3.1 初始化项目结构

```bash
mkdir guandan-project && cd guandan-project
mkdir backend frontend sdk infra
touch docker-compose.yml README.md
```

* `backend/` 内初始化 Go 模块 `go mod init`
* `frontend/` 使用 Vite 快速启动：`npm create vite@latest frontend -- --template react-ts`

---

### ✅ 3.2 实现 `/healthz` 接口

**路径**：`backend/main.go`

```go
r.GET("/healthz", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "pong"})
})
```

前端在首页调用 `/healthz`，验证后端联通性。

---

### ✅ 3.3 配置 Docker 支持

**前端 Dockerfile 示例**：`frontend/Dockerfile`

```dockerfile
FROM node:18 AS builder
WORKDIR /app
COPY . .
RUN npm install && npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
```

**后端 Dockerfile 示例**：`backend/Dockerfile`

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server main.go
EXPOSE 8080
CMD ["./server"]
```

---

### ✅ 3.4 编写 `docker-compose.yml`

```yaml
version: "3.9"

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:80"
  
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: guandan
      POSTGRES_PASSWORD: guandan
      POSTGRES_DB: guandan
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"
```

---

### ✅ 3.5 设置 GitHub Actions 自动化流程

**路径**：`.github/workflows/ci.yml`

```yaml
name: CI Build & Test

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: guandan
          POSTGRES_PASSWORD: guandan
        ports:
          - 5432:5432

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Build backend
        run: |
          cd backend
          go build -v ./...

      - name: Run tests
        run: |
          cd sdk
          go test ./... -v
```

---

### ✅ 3.6 添加环境变量管理

使用 `.env` 和 `.env.example` 配置数据库连接、JWT 密钥等，避免写死配置。

---

### ✅ 3.7 部署 Staging 环境（可选）

* 若已购买域名和服务器，可通过 `docker-compose -f docker-compose.yml up` 部署 Staging
* 后续阶段可引入 Terraform / Pulumi 自动化部署脚本

---

## 📌 4. 验收标准（Definition of Done）

| 验收项               | 检查点                                                          |
| ----------------- | ------------------------------------------------------------ |
| 本地联调              | 执行 `docker-compose up`，前端访问 `/healthz` 返回 `{ status: pong }` |
| GitHub Actions CI | 提交代码后，自动测试 / 构建流程跑通                                          |
| 基础服务启动            | PostgreSQL / Redis 正常连接、服务容器正常运行                             |
| 前端显示              | 首页可见“服务正常运行”的提示                                              |

---

## 🔚 阶段结束后应具备能力：

* 每位开发者一键启动全栈环境
* 每次提交代码都可验证构建、测试、部署是否通过
* 可以在浏览器中验证后端接口联通性
* 可以开始编写 SDK 模块及业务逻辑开发

---

