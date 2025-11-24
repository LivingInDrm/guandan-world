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

**开发环境：**
- Docker 20.10+ 和 Docker Compose 2.0+（运行后端）
- Node.js 18+ 和 npm（运行前端）
- 至少 4GB RAM
- 至少 10GB 磁盘空间

**生产环境：**
- Docker 20.10+ 和 Docker Compose 2.0+
- 至少 4GB RAM
- 至少 10GB 磁盘空间

### 开发环境部署

**开发架构：后端Docker + 前端本地运行**

这样做的好处：
- ✅ 前端热重载更快，无Service Worker干扰
- ✅ 调试方便，直接使用浏览器DevTools
- ✅ 后端和数据库隔离在Docker中

```bash
# 1. 克隆项目
git clone https://github.com/your-username/guandan-world.git
cd guandan-world

# 2. 安装根目录依赖并生成 Proto 文件（首次运行必需）
npm install      # 会自动运行 make proto-ts 生成 TypeScript 类型

# 3. 启动后端服务（Docker）
./deploy.sh dev deploy

# 4. 在新终端启动前端（本地）
cd frontend
npm install      # 首次运行需要
npm run dev      # 启动Vite开发服务器（会自动检查并生成 Proto 文件）

# 5. 访问应用
# 前端: http://localhost:5173 (Vite本地服务)
# 后端: http://localhost:8080
```

**开发模式特性：**
- ✅ 前端本地运行，热重载快速
- ✅ 后端Docker容器，支持Air热重载
- ✅ 详细的调试日志
- ✅ Vite自动代理API请求到后端

### 生产环境部署

```bash
# 配置环境变量
# 编辑 .env.production 文件（首次运行会自动创建）

# 部署生产环境
./deploy.sh prod deploy

# 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8080
# 监控: http://localhost:3001 (Grafana)
```

### 访问地址

| 服务 | 开发模式 | 生产模式 |
|------|---------|---------|
| 前端 | http://localhost:5173 (本地npm) | http://localhost:3000 (Docker) |
| 后端 API | http://localhost:8080 (Docker) | http://localhost:8080 (Docker) |
| WebSocket | ws://localhost:8080/ws | ws://localhost:8080/ws |
| Grafana | - | http://localhost:3001 |
| Prometheus | - | http://localhost:9090 |

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
│   ├── deal.go             # 牌局管理（包含发牌逻辑）
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

### 详细文档

- [DEPLOYMENT.md](DEPLOYMENT.md) - 完整部署文档

## 🛠️ 开发

### 开发工作流

```bash
# 1. 启动开发环境
./deploy.sh dev deploy

# 2. 修改代码（自动热重载）
# - 后端代码: backend/*.go → Air 自动重启
# - 前端代码: frontend/src/* → Vite HMR

# 3. 查看日志
./deploy.sh dev logs backend  # 后端日志
./deploy.sh dev logs frontend # 前端日志

# 4. 重启服务
./deploy.sh dev restart

# 5. 停止服务
./deploy.sh dev stop
```

### 常用命令

```bash
# 部署管理
./deploy.sh dev deploy    # 首次部署
./deploy.sh dev start      # 启动服务
./deploy.sh dev stop       # 停止服务
./deploy.sh dev restart    # 重启服务
./deploy.sh dev rebuild    # 重新构建镜像

# 日志和调试
./deploy.sh dev logs       # 查看所有日志
./deploy.sh dev logs -f    # 实时跟踪日志
./deploy.sh dev status     # 查看服务状态
./deploy.sh dev health     # 健康检查
```

### 代码规范

- Go: 使用 `gofmt` 和 `golint`
- TypeScript: 使用 ESLint 和 Prettier
- 提交信息: 使用 Conventional Commits

### 为什么所有服务都在 Docker 中？

✅ **环境一致性**: 开发环境 = 生产环境  
✅ **快速上手**: 新成员只需一条命令  
✅ **依赖隔离**: 不污染本地环境  
✅ **问题重现**: 避免"在我机器上能运行"

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