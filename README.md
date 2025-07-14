# 掼蛋世界 (Guandan World)

一个基于 Go + React 的掼蛋游戏平台，支持实时多人在线对战。

## 🏗️ 项目结构

```
guandan-world/
├── backend/               # Go 后端服务
│   ├── main.go           # 主服务入口
│   ├── Dockerfile        # 后端 Docker 配置
│   └── go.mod            # Go 模块配置
├── frontend/             # React 前端应用
│   ├── src/              # 前端源码
│   ├── Dockerfile        # 前端 Docker 配置
│   └── package.json      # 前端依赖配置
├── sdk/                  # 游戏核心逻辑
│   ├── card.go           # 卡牌结构体和逻辑
│   ├── card_test.go      # 卡牌测试用例
│   ├── card_example.go   # 卡牌功能演示
│   └── go.mod            # SDK 模块配置
├── infra/                # DevOps 相关脚本
├── .github/workflows/    # GitHub Actions 工作流
├── docker-compose.yml    # 本地开发环境编排
└── README.md             # 项目说明
```

## 🚀 快速开始

### 前置要求

- Docker 和 Docker Compose
- Go 1.22+
- Node.js 18+

### 一键启动

```bash
# 克隆项目
git clone https://github.com/LivingInDrm/guandan-world.git
cd guandan-world

# 复制环境变量配置
cp .env.example .env

# 启动所有服务
docker-compose up --build
```

### 访问应用

- 前端应用：http://localhost:3000
- 后端 API：http://localhost:8080
- 健康检查：http://localhost:8080/healthz

## 🧰 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 后端 | Go + Gin | 高性能 API 服务 |
| 前端 | React + TypeScript + Vite | 现代前端开发 |
| 数据库 | PostgreSQL | 用户数据存储 |
| 缓存 | Redis | 游戏状态缓存 |
| 容器化 | Docker + Docker Compose | 环境一致性 |
| CI/CD | GitHub Actions | 自动化构建测试 |

## 🎮 游戏核心模块

### 卡牌系统 (Card)

完整的掼蛋卡牌系统实现，支持：

- **卡牌类型**：普通牌(2-10)、人头牌(J/Q/K/A)、大小王
- **特殊规则**：级别牌、变化牌（红桃级别牌）
- **比较逻辑**：牌的大小比较、顺子比较
- **功能特性**：卡牌克隆、JSON 编码、字符串表示

```go
// 创建卡牌
card, err := NewCard(3, "Spade", 2)  // 3 of Spade, 级别为2
ace, err := NewCard(1, "Heart", 2)   // Ace of Heart
joker, err := NewCard(16, "Joker", 2) // Red Joker

// 比较卡牌
if card1.GreaterThan(card2) {
    fmt.Printf("%s 比 %s 大\n", card1.String(), card2.String())
}

// 变化牌判断
if card.IsWildcard() {
    fmt.Println("这是一张变化牌（红桃级别牌）")
}
```

## 🔧 开发环境

### 本地开发

```bash
# 启动后端服务
cd backend
go run main.go

# 启动前端服务
cd frontend
npm install
npm run dev

# 启动数据库（可选）
docker-compose up postgres redis
```

### 运行测试

```bash
# 后端测试
cd backend
go test ./...

# SDK 测试
cd sdk
go test ./...

# 前端测试
cd frontend
npm test
```

### 构建部署

```bash
# 构建 Docker 镜像
docker-compose build

# 部署到生产环境
docker-compose -f docker-compose.yml up -d
```

## 🎮 游戏特性

- ✅ 实时多人对战
- ✅ 完整的掼蛋规则实现
- ✅ 用户认证系统
- ✅ 游戏房间管理
- ✅ 实时聊天功能
- ✅ 游戏回放功能
- ✅ 完整的卡牌系统

## 📚 API 文档

### 健康检查

```bash
GET /healthz
```

响应：
```json
{
  "status": "pong"
}
```

## 🧪 测试覆盖

- ✅ 卡牌创建与验证
- ✅ 卡牌比较逻辑
- ✅ 变化牌判断
- ✅ 顺子比较
- ✅ JSON 编码
- ✅ 克隆与相等判断
- ✅ 完整的单元测试

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📞 联系方式

- 项目链接：https://github.com/LivingInDrm/guandan-world
- 问题反馈：https://github.com/LivingInDrm/guandan-world/issues 