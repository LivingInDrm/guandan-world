# 掼蛋在线游戏 (Guandan Online Game)

一个基于Go后端和React前端的在线掼蛋游戏平台，支持实时多人对战、WebSocket通信和完整的游戏管理系统。

## 🎮 项目特性

- **完整的掼蛋游戏逻辑**: 基于标准掼蛋规则实现
- **实时多人对战**: WebSocket支持4人实时游戏
- **用户认证系统**: JWT认证，安全可靠
- **房间管理**: 创建、加入、管理游戏房间
- **断线重连**: 自动托管和重连机制
- **性能优化**: 消息批处理、增量更新、压缩传输
- **监控完整**: Prometheus + Grafana + Loki监控栈
- **容器化部署**: Docker + Docker Compose一键部署

## 🏗️ 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     用户        │    │     Nginx       │    │     后端        │
│   (浏览器)      │◄──►│   (反向代理)    │◄──►│   (Go服务)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     前端        │    │     Redis       │    │   监控系统      │
│  (React应用)    │    │    (缓存)       │    │ (Prometheus)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 快速开始

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 4GB RAM
- 至少 10GB 磁盘空间

### 开发环境部署

```bash
# 克隆项目
git clone https://github.com/your-username/guandan-world.git
cd guandan-world

# 启动开发环境
./deploy.sh development deploy
```

### 生产环境部署

```bash
# 配置环境变量
cp .env.example .env.production
# 编辑 .env.production 文件

# 部署生产环境
./deploy.sh production deploy
```

### 访问应用

- **游戏前端**: http://localhost:3000
- **后端API**: http://localhost:8080
- **监控面板**: http://localhost:3001 (Grafana)
- **指标监控**: http://localhost:9090 (Prometheus)

## 📁 项目结构

```
guandan-world/
├── backend/                 # Go后端服务
│   ├── auth/               # 用户认证
│   ├── room/               # 房间管理
│   ├── game/               # 游戏服务
│   ├── websocket/          # WebSocket管理
│   ├── handlers/           # HTTP处理器
│   └── integration_tests/  # 集成测试
├── frontend/               # React前端应用
│   ├── src/
│   │   ├── components/     # React组件
│   │   ├── services/       # API服务
│   │   ├── store/          # 状态管理
│   │   └── test/           # 测试文件
│   └── public/
├── sdk/                    # 掼蛋游戏引擎
│   ├── game_engine.go      # 游戏引擎
│   ├── dealer.go           # 发牌系统
│   ├── trick.go            # 出牌逻辑
│   └── result.go           # 结算系统
├── monitoring/             # 监控配置
│   ├── prometheus.yml      # Prometheus配置
│   ├── loki.yml           # 日志聚合配置
│   └── grafana/           # Grafana面板
├── nginx/                  # Nginx配置
├── docker-compose.yml      # 开发环境配置
├── docker-compose.production.yml  # 生产环境配置
└── deploy.sh              # 部署脚本
```

## 🎯 功能特性

### 游戏功能

- ✅ **用户认证**: 注册、登录、JWT认证
- ✅ **房间大厅**: 房间列表、创建房间、加入房间
- ✅ **房间等待**: 玩家管理、座位分配、游戏开始
- ✅ **游戏流程**: 发牌、上贡、出牌、结算
- ✅ **实时通信**: WebSocket双向通信
- ✅ **断线托管**: 自动托管、重连恢复
- ✅ **操作控制**: 超时检测、自动操作

### 技术特性

- ✅ **高性能**: 消息优化、批处理、压缩
- ✅ **高可用**: 健康检查、自动重启、负载均衡
- ✅ **可观测**: 完整监控、日志聚合、告警
- ✅ **安全性**: HTTPS、限流、输入验证
- ✅ **可扩展**: 微服务架构、容器化部署

## 🧪 测试

### 运行所有测试

```bash
# 运行综合测试套件
./test_e2e_comprehensive.sh

# 运行单元测试
./test_e2e_comprehensive.sh unit

# 运行集成测试
./test_e2e_comprehensive.sh integration

# 运行E2E测试
./test_e2e_comprehensive.sh e2e
```

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
npm test
```

## 📊 监控

### 监控指标

- **系统指标**: CPU、内存、磁盘、网络
- **应用指标**: 响应时间、错误率、吞吐量
- **业务指标**: 在线用户、活跃游戏、完成率

### 监控面板

访问 http://localhost:3001 查看Grafana监控面板：

- 系统概览
- 在线用户数
- 活跃游戏数
- API响应时间
- WebSocket连接数
- 错误率统计

## 🔧 配置

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `JWT_SECRET` | - | JWT密钥（生产环境必须设置） |
| `JWT_EXPIRY` | `24h` | JWT过期时间 |
| `CORS_ORIGINS` | `http://localhost:3000` | 允许的跨域源 |
| `LOG_LEVEL` | `info` | 日志级别 |
| `REDIS_PASSWORD` | - | Redis密码 |

### 部署配置

详细的部署配置请参考 [DEPLOYMENT.md](DEPLOYMENT.md)

## 🛠️ 开发

### 本地开发

```bash
# 启动后端
cd backend
go run main.go

# 启动前端
cd frontend
npm start

# 启动Redis
docker run -d -p 6379:6379 redis:alpine
```

### 代码规范

- Go: 使用 `gofmt` 和 `golint`
- TypeScript: 使用 ESLint 和 Prettier
- 提交信息: 使用 Conventional Commits

## 📈 性能优化

### WebSocket优化

- **消息批处理**: 50ms间隔批量发送
- **增量更新**: 只发送状态变化
- **消息压缩**: 大消息自动压缩
- **连接复用**: 减少连接开销

### 缓存策略

- **静态资源**: 1年缓存
- **API响应**: 适当缓存
- **游戏状态**: 内存缓存
- **用户会话**: Redis缓存

## 🔒 安全

- **HTTPS**: 生产环境强制HTTPS
- **JWT认证**: 安全的用户认证
- **输入验证**: 防止注入攻击
- **限流保护**: API和WebSocket限流
- **CORS配置**: 跨域请求控制

## 📝 API文档

### 认证接口

```
POST /api/auth/register  # 用户注册
POST /api/auth/login     # 用户登录
POST /api/auth/logout    # 用户登出
GET  /api/auth/me        # 获取用户信息
```

### 房间接口

```
GET  /api/rooms          # 获取房间列表
POST /api/rooms/create   # 创建房间
POST /api/rooms/join     # 加入房间
POST /api/rooms/leave    # 离开房间
POST /api/rooms/:id/start # 开始游戏
```

### WebSocket事件

```
game_prepare    # 游戏准备
game_begin      # 游戏开始
player_turn     # 玩家回合
card_played     # 出牌事件
game_end        # 游戏结束
```

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- 感谢所有贡献者
- 感谢开源社区的支持
- 特别感谢掼蛋游戏的发明者

## 📞 联系

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 加入讨论群

---

**享受掼蛋游戏的乐趣！** 🎉