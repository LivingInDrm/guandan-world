# 掼蛋在线游戏部署文档

## 概述

本文档详细介绍了掼蛋在线游戏的部署流程、配置选项和运维指南。

## 系统架构

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

## 快速开始

### 1. 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 4GB RAM
- 至少 10GB 磁盘空间

### 2. 克隆项目

```bash
git clone <repository-url>
cd guandan-world
```

### 3. 开发环境部署

```bash
# 使用部署脚本
./deploy.sh development deploy

# 或者直接使用docker-compose
docker-compose up -d
```

### 4. 生产环境部署

```bash
# 配置环境变量
cp .env.example .env.production
# 编辑 .env.production 文件

# 部署
./deploy.sh production deploy
```

## 详细配置

### 环境变量配置

#### 后端配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `JWT_SECRET` | - | JWT密钥（生产环境必须设置） |
| `JWT_EXPIRY` | `24h` | JWT过期时间 |
| `CORS_ORIGINS` | `http://localhost:3000` | 允许的跨域源 |
| `LOG_LEVEL` | `info` | 日志级别 |
| `WEBSOCKET_BUFFER_SIZE` | `1024` | WebSocket缓冲区大小 |
| `MAX_CONCURRENT_GAMES` | `100` | 最大并发游戏数 |
| `GAME_TIMEOUT` | `30m` | 游戏超时时间 |

#### 前端配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `REACT_APP_API_URL` | `http://localhost:8080` | 后端API地址 |
| `REACT_APP_WS_URL` | `ws://localhost:8080` | WebSocket地址 |
| `REACT_APP_ENV` | `development` | 运行环境 |

#### Redis配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `REDIS_PASSWORD` | - | Redis密码 |

### Docker Compose配置

#### 开发环境 (docker-compose.yml)

```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=debug
    volumes:
      - ./backend:/app
      - /app/vendor
  
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
```

#### 生产环境 (docker-compose.production.yml)

包含完整的监控栈：
- Nginx反向代理
- Redis缓存
- Prometheus监控
- Grafana仪表板
- Loki日志聚合

## 性能优化

### 1. WebSocket优化

- **消息批处理**: 50ms间隔批量发送消息
- **增量更新**: 只发送状态变化部分
- **消息压缩**: 大于1KB的消息自动压缩
- **连接池**: 复用WebSocket连接

### 2. Nginx优化

- **Gzip压缩**: 启用静态资源压缩
- **缓存策略**: 静态资源1年缓存
- **限流保护**: API和WebSocket限流
- **Keep-Alive**: 连接复用

### 3. 后端优化

- **连接池**: 数据库和Redis连接池
- **内存缓存**: 游戏状态内存缓存
- **异步处理**: 非阻塞IO操作
- **资源限制**: CPU和内存限制

## 监控和日志

### 1. 监控指标

#### 系统指标
- CPU使用率
- 内存使用率
- 磁盘IO
- 网络流量

#### 应用指标
- 在线用户数
- 活跃游戏数
- API响应时间
- WebSocket连接数
- 错误率

#### 业务指标
- 用户注册数
- 游戏完成率
- 平均游戏时长
- 断线重连率

### 2. 日志管理

#### 日志级别
- `ERROR`: 错误信息
- `WARN`: 警告信息
- `INFO`: 一般信息
- `DEBUG`: 调试信息

#### 日志格式
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "service": "backend",
  "message": "User logged in",
  "user_id": "12345",
  "request_id": "req-67890"
}
```

### 3. 告警规则

#### 系统告警
- CPU使用率 > 80%
- 内存使用率 > 85%
- 磁盘使用率 > 90%

#### 应用告警
- API错误率 > 5%
- 响应时间 > 2s
- WebSocket连接失败率 > 10%

## 部署脚本使用

### 基本命令

```bash
# 部署应用
./deploy.sh [environment] deploy

# 启动服务
./deploy.sh [environment] start

# 停止服务
./deploy.sh [environment] stop

# 重启服务
./deploy.sh [environment] restart

# 查看日志
./deploy.sh [environment] logs [service]

# 查看状态
./deploy.sh [environment] status

# 清理资源
./deploy.sh [environment] clean

# 备份数据
./deploy.sh [environment] backup

# 恢复数据
./deploy.sh [environment] restore [backup_dir]
```

### 示例

```bash
# 开发环境部署
./deploy.sh development deploy

# 生产环境启动
./deploy.sh production start

# 查看后端日志
./deploy.sh production logs backend

# 备份生产数据
./deploy.sh production backup
```

## 故障排除

### 常见问题

#### 1. 服务启动失败

**症状**: 容器无法启动或立即退出

**解决方案**:
```bash
# 查看容器日志
docker-compose logs [service]

# 检查配置文件
docker-compose config

# 重新构建镜像
docker-compose build --no-cache [service]
```

#### 2. WebSocket连接失败

**症状**: 前端无法建立WebSocket连接

**解决方案**:
```bash
# 检查后端服务状态
curl http://localhost:8080/healthz

# 检查Nginx配置
docker exec guandan-nginx nginx -t

# 查看WebSocket日志
tail -f logs/nginx/websocket.log
```

#### 3. 数据库连接问题

**症状**: 后端无法连接Redis

**解决方案**:
```bash
# 检查Redis状态
docker exec guandan-redis redis-cli ping

# 检查网络连接
docker network ls
docker network inspect guandan-world_guandan-network
```

#### 4. 性能问题

**症状**: 响应缓慢或超时

**解决方案**:
```bash
# 查看资源使用情况
docker stats

# 检查监控指标
# 访问 http://localhost:3001 (Grafana)

# 分析日志
grep "slow" logs/backend/*.log
```

### 日志分析

#### 后端日志位置
- 应用日志: `logs/backend/app.log`
- 错误日志: `logs/backend/error.log`
- 访问日志: `logs/backend/access.log`

#### 前端日志位置
- Nginx访问日志: `logs/nginx/access.log`
- Nginx错误日志: `logs/nginx/error.log`
- WebSocket日志: `logs/nginx/websocket.log`

## 安全配置

### 1. 网络安全

- 使用HTTPS (生产环境)
- 配置防火墙规则
- 限制不必要的端口暴露

### 2. 应用安全

- JWT密钥定期轮换
- 输入验证和过滤
- SQL注入防护
- XSS防护

### 3. 容器安全

- 使用非root用户运行
- 最小化镜像体积
- 定期更新基础镜像
- 扫描安全漏洞

## 备份和恢复

### 1. 数据备份

```bash
# 自动备份
./deploy.sh production backup

# 手动备份Redis
docker exec guandan-redis redis-cli --rdb /data/backup.rdb
docker cp guandan-redis:/data/backup.rdb ./backups/
```

### 2. 数据恢复

```bash
# 从备份恢复
./deploy.sh production restore backups/20240101_120000

# 手动恢复Redis
docker cp ./backups/redis.rdb guandan-redis:/data/dump.rdb
docker-compose restart redis
```

### 3. 备份策略

- **每日备份**: 自动备份关键数据
- **增量备份**: 只备份变化的数据
- **异地备份**: 备份到云存储
- **定期测试**: 验证备份可用性

## 扩展和升级

### 1. 水平扩展

```yaml
# docker-compose.production.yml
services:
  backend:
    deploy:
      replicas: 3
    
  nginx:
    depends_on:
      - backend
```

### 2. 垂直扩展

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
```

### 3. 滚动更新

```bash
# 无停机更新
docker-compose up -d --no-deps backend
```

## 联系和支持

如有问题，请联系开发团队或查看项目文档。

---

*最后更新: 2024年1月*