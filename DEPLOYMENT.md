# 部署文档

## 🎯 部署模式

项目只有两种部署模式：

| 模式 | 用途 | 命令 |
|------|------|------|
| **开发模式** | 本地开发、调试 | `./deploy.sh dev deploy` |
| **生产模式** | 生产部署 | `./deploy.sh prod deploy` |

**核心原则**：所有服务（前端、后端、数据库）都在 Docker 中运行。

---

## 🚀 快速开始

### 环境要求

```bash
# 必需
Docker 20.10+
Docker Compose 2.0+

# 无需安装
Go、Node.js、Redis、PostgreSQL（都在 Docker 中）
```

### 开发环境

```bash
# 1. 启动（首次会自动构建）
./deploy.sh dev deploy

# 2. 访问
前端: http://localhost:5173
后端: http://localhost:8080

# 3. 查看日志
./deploy.sh dev logs backend
./deploy.sh dev logs frontend

# 4. 停止
./deploy.sh dev stop
```

**开发模式特性**：
- ✅ 代码热重载（修改代码自动生效）
- ✅ 本地代码挂载到容器
- ✅ 详细的调试日志

### 生产环境

```bash
# 1. 配置环境变量（首次部署会自动创建 .env.production）
# 编辑 .env.production，修改所有密码和密钥

# 2. 部署
./deploy.sh prod deploy

# 3. 访问
前端: http://localhost:3000
后端: http://localhost:8080
监控: http://localhost:3001 (Grafana)
```

---

## 📋 常用命令

```bash
# 部署管理
./deploy.sh [dev|prod] deploy    # 部署（构建+启动）
./deploy.sh [dev|prod] start      # 启动
./deploy.sh [dev|prod] stop       # 停止
./deploy.sh [dev|prod] restart    # 重启
./deploy.sh [dev|prod] rebuild    # 重新构建

# 日志和状态
./deploy.sh [dev|prod] logs       # 查看所有日志
./deploy.sh [dev|prod] logs backend  # 查看后端日志
./deploy.sh [dev|prod] status     # 查看状态
./deploy.sh [dev|prod] health     # 健康检查

# 数据管理（仅生产）
./deploy.sh prod backup           # 备份数据
./deploy.sh prod restore [dir]    # 恢复数据

# 清理
./deploy.sh [dev|prod] clean      # 清理容器和卷
```

---

## 🔧 开发工作流

```bash
# 1. 启动开发环境
./deploy.sh dev deploy

# 2. 修改代码
# - 后端 (*.go): Air 自动重启
# - 前端 (*.tsx): Vite HMR 即时更新

# 3. 查看日志调试
./deploy.sh dev logs -f

# 4. 完成后停止
./deploy.sh dev stop
```

### 访问容器

```bash
# 进入后端容器
docker exec -it guandan-backend sh

# 访问数据库
docker exec -it guandan-postgres psql -U guandan

# 访问 Redis
docker exec -it guandan-redis redis-cli
```

---

## ⚙️ 配置说明

### 生产环境变量（.env.production）

```bash
# JWT 配置（必须修改）
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h

# CORS 配置
CORS_ORIGINS=https://yourdomain.com

# 数据库密码（必须修改）
REDIS_PASSWORD=your-redis-password
POSTGRES_PASSWORD=your-postgres-password

# 监控密码（必须修改）
GRAFANA_PASSWORD=your-grafana-password

# 前端配置
REACT_APP_API_URL=https://yourdomain.com
REACT_APP_WS_URL=wss://yourdomain.com
```

### 文件结构

```
guandan-world/
├── backend/
│   ├── Dockerfile          # 生产环境
│   ├── Dockerfile.dev      # 开发环境（热重载）
│   └── .air.toml          # Air 配置
├── frontend/
│   ├── Dockerfile          # 生产环境
│   └── Dockerfile.dev      # 开发环境（Vite）
├── docker-compose.dev.yml          # 开发配置
├── docker-compose.production.yml   # 生产配置
└── deploy.sh                       # 部署脚本
```

---

## 🔍 故障排查

### 服务启动失败

```bash
# 查看日志
./deploy.sh dev logs backend

# 检查配置
docker-compose -f docker-compose.dev.yml config

# 重新构建
./deploy.sh dev rebuild
```

### 端口被占用

```bash
# 检查端口占用
lsof -i :8080
lsof -i :5173

# 停止占用端口的服务
./deploy.sh dev stop
```

### 热重载不工作

```bash
# 重新构建镜像
./deploy.sh dev rebuild

# 检查 Volume 挂载
docker-compose -f docker-compose.dev.yml ps
```

### 数据库连接失败

```bash
# 检查服务状态
./deploy.sh dev status

# 查看数据库日志
./deploy.sh dev logs postgres

# 重启数据库
docker-compose -f docker-compose.dev.yml restart postgres
```

---

## 📊 对比：开发 vs 生产

| 特性 | 开发模式 | 生产模式 |
|------|---------|---------|
| 热重载 | ✅ | ❌ |
| 代码挂载 | ✅ Volume | ❌ 打包到镜像 |
| 前端端口 | 5173 (Vite) | 3000 (Nginx) |
| 日志级别 | DEBUG | INFO |
| 监控栈 | ❌ | ✅ Prometheus + Grafana |
| 构建时间 | 快 | 慢（优化构建） |

---

## 💡 最佳实践

### 开发环境

1. 定期重建镜像以更新依赖
   ```bash
   ./deploy.sh dev rebuild
   ```

2. 清理磁盘空间
   ```bash
   docker system prune -a
   ```

3. 查看资源使用
   ```bash
   docker stats
   ```

### 生产环境

1. **修改所有默认密码**
2. **配置 HTTPS**
3. **设置定时备份**
4. **配置监控告警**

---

## ❓ 常见问题

**Q: 为什么所有服务都在 Docker 中？**  
A: 确保开发环境 = 生产环境，避免"在我机器上能运行"的问题。

**Q: 如何在 IDE 中使用代码补全？**  
A: 代码在本地，IDE 可以正常解析。VSCode 会自动识别 Go 和 TypeScript 项目。

**Q: 数据会丢失吗？**  
A: 不会。使用 Docker Volume 持久化。只有执行 `clean` 命令才会删除数据。

**Q: 如何切换模式？**  
A: 
```bash
./deploy.sh dev stop    # 停止开发环境
./deploy.sh prod deploy # 启动生产环境
```

---

**需要帮助？** 运行 `./deploy.sh --help` 查看所有命令。
