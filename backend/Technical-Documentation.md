# 📚 掼蛋在线游戏后端技术文档

基于对backend模块代码的深入分析，这是完整的技术文档：

---

## 🏗️ **系统架构概述**

### **整体架构**
掼蛋在线游戏后端采用分层微服务架构，实现了从单机游戏到在线多人游戏平台的升级：

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend/Client                      │
└─────────────────────────────────────────────────────────┘
                              │
                         HTTP/WebSocket
                              │
┌─────────────────────────────────────────────────────────┐
│                    Backend Services                     │
├─────────────────────────────────────────────────────────┤
│  🔐 Auth Service     🏠 Room Service     🎮 Game Service │
│  📡 WebSocket Manager     🌐 HTTP Handlers             │
└─────────────────────────────────────────────────────────┘
                              │
                         调用SDK接口
                              │
┌─────────────────────────────────────────────────────────┐
│                    Game Engine SDK                      │
├─────────────────────────────────────────────────────────┤
│     🃏 Card Logic     🎯 Game Rules     🤖 AI Players    │
└─────────────────────────────────────────────────────────┘
```

### **核心设计原则**
- **🎯 职责分离**：Backend专注于在线协调，Game SDK专注于游戏逻辑
- **🔄 状态同步**：实时同步游戏状态到所有客户端
- **🛡️ 安全第一**：JWT认证 + CORS配置
- **📡 实时通信**：WebSocket双向通信
- **🏗️ 可扩展性**：模块化设计，易于扩展

---

## 🔧 **核心组件详解**

### **1. 🔐 认证服务 (Auth Service)**

#### **文件位置**: `backend/auth/service.go`

#### **核心功能**:
- **JWT令牌管理**：生成、验证、刷新JWT令牌
- **用户身份验证**：基于用户名/密码的简单认证
- **会话管理**：维护用户登录状态

#### **主要方法**:
```go
type AuthService struct {
    jwtSecret []byte
    tokenExpiry time.Duration
}

// 用户登录，返回JWT令牌
func (a *AuthService) Login(username, password string) (string, error)

// 验证JWT令牌有效性
func (a *AuthService) ValidateToken(tokenString string) (*Claims, error)

// 刷新JWT令牌
func (a *AuthService) RefreshToken(tokenString string) (string, error)
```

#### **JWT Claims结构**:
```go
type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

---

### **2. 🏠 房间服务 (Room Service)**

#### **文件位置**: `backend/room/service.go`

#### **核心功能**:
- **房间生命周期管理**：创建、加入、离开、销毁房间
- **玩家状态追踪**：维护房间内玩家列表和状态
- **游戏准备检查**：确保4人齐全后启动游戏

#### **数据结构**:
```go
type Room struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Players     map[string]Player `json:"players"`
    GameStarted bool              `json:"game_started"`
    CreatedAt   time.Time         `json:"created_at"`
    mu          sync.RWMutex
}

type Player struct {
    Username string `json:"username"`
    Ready    bool   `json:"ready"`
    JoinedAt time.Time `json:"joined_at"`
}
```

#### **主要方法**:
```go
// 创建新房间
func (rs *RoomService) CreateRoom(roomName, username string) (*Room, error)

// 加入房间
func (rs *RoomService) JoinRoom(roomID, username string) error

// 离开房间
func (rs *RoomService) LeaveRoom(roomID, username string) error

// 设置玩家准备状态
func (rs *RoomService) SetPlayerReady(roomID, username string, ready bool) error

// 检查是否可以开始游戏
func (rs *RoomService) CanStartGame(roomID string) bool
```

---

### **3. 🎮 游戏服务 (Game Service)**

#### **文件位置**: `backend/game/service.go`

#### **核心功能**:
- **游戏引擎协调**：管理多个游戏实例
- **玩家输入处理**：接收并转发玩家操作到游戏引擎
- **游戏状态广播**：将游戏状态变化推送给所有玩家

#### **设计亮点**:
```go
type GameService struct {
    engines   map[string]sdk.GameEngineInterface // 游戏引擎实例
    wsManager WSManagerInterface                 // WebSocket管理器
    mu        sync.RWMutex
}

// WebSocket管理器接口
type WSManagerInterface interface {
    BroadcastToRoom(roomID string, message []byte) error
    SendToUser(username string, message []byte) error
}
```

#### **主要方法**:
```go
// 为房间创建游戏实例
func (gs *GameService) CreateGame(roomID string, playerUsernames []string) error

// 处理玩家游戏输入
func (gs *GameService) HandlePlayerInput(roomID, username string, input interface{}) error

// 获取游戏状态
func (gs *GameService) GetGameState(roomID string) (interface{}, error)

// 结束游戏
func (gs *GameService) EndGame(roomID string) error
```

---

### **4. 📡 WebSocket管理器 (WebSocket Manager)**

#### **文件位置**: `backend/websocket/manager.go`

#### **核心功能**:
- **连接生命周期管理**：建立、维护、关闭WebSocket连接
- **消息路由**：支持点对点和广播消息
- **房间分组**：按房间组织连接，支持房间级广播

#### **数据结构**:
```go
type Manager struct {
    connections map[string]*Connection // username -> connection
    rooms       map[string][]string    // roomID -> usernames
    mu          sync.RWMutex
}

type Connection struct {
    Username string
    Conn     *websocket.Conn
    RoomID   string
}
```

#### **消息类型**:
```go
type Message struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

// 消息类型常量
const (
    MessageTypeJoinRoom    = "join_room"
    MessageTypeLeaveRoom   = "leave_room"
    MessageTypeGameState   = "game_state"
    MessageTypePlayerInput = "player_input"
    MessageTypeError       = "error"
)
```

---

## 🌐 **HTTP API接口**

### **认证接口 (handlers/auth.go)**

#### **POST /auth/login**
用户登录接口

**请求体**:
```json
{
    "username": "player1",
    "password": "password123"
}
```

**响应**:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T10:30:00Z"
}
```

#### **POST /auth/refresh**
刷新JWT令牌

**请求头**:
```
Authorization: Bearer <current_token>
```

**响应**:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T11:30:00Z"
}
```

---

### **房间接口 (handlers/room.go)**

#### **POST /rooms**
创建新房间

**请求体**:
```json
{
    "name": "我的掼蛋房间"
}
```

**响应**:
```json
{
    "id": "room_123456",
    "name": "我的掼蛋房间",
    "players": {
        "player1": {
            "username": "player1",
            "ready": false,
            "joined_at": "2024-01-15T10:00:00Z"
        }
    },
    "game_started": false,
    "created_at": "2024-01-15T10:00:00Z"
}
```

#### **GET /rooms**
获取房间列表

**响应**:
```json
{
    "rooms": [
        {
            "id": "room_123456",
            "name": "我的掼蛋房间",
            "player_count": 1,
            "game_started": false
        }
    ]
}
```

#### **POST /rooms/{roomID}/join**
加入房间

**响应**:
```json
{
    "message": "成功加入房间",
    "room": { /* 房间详情 */ }
}
```

#### **DELETE /rooms/{roomID}/leave**
离开房间

**响应**:
```json
{
    "message": "成功离开房间"
}
```

#### **POST /rooms/{roomID}/ready**
设置准备状态

**请求体**:
```json
{
    "ready": true
}
```

---

## 📡 **WebSocket通信协议**

### **连接建立**
```
ws://localhost:8080/ws?token=<jwt_token>
```

### **消息格式**
所有WebSocket消息使用统一格式：

```json
{
    "type": "message_type",
    "data": { /* 具体数据 */ }
}
```

### **客户端发送的消息类型**

#### **1. 加入房间**
```json
{
    "type": "join_room",
    "data": {
        "room_id": "room_123456"
    }
}
```

#### **2. 离开房间**
```json
{
    "type": "leave_room",
    "data": {
        "room_id": "room_123456"
    }
}
```

#### **3. 游戏操作**
```json
{
    "type": "player_input",
    "data": {
        "action": "play_cards",
        "cards": ["AS", "KH", "QD"]
    }
}
```

### **服务器推送的消息类型**

#### **1. 房间状态更新**
```json
{
    "type": "room_state",
    "data": {
        "room_id": "room_123456",
        "players": { /* 玩家列表 */ },
        "game_started": false
    }
}
```

#### **2. 游戏状态更新**
```json
{
    "type": "game_state",
    "data": {
        "current_player": "player1",
        "game_phase": "playing",
        "trick_info": { /* 当前技巧信息 */ }
    }
}
```

#### **3. 错误消息**
```json
{
    "type": "error",
    "data": {
        "message": "房间已满",
        "code": "ROOM_FULL"
    }
}
```

---

## 🚀 **部署配置**

### **环境变量**
```bash
# 服务器配置
PORT=8080
HOST=0.0.0.0

# JWT配置
JWT_SECRET=your-super-secret-key
JWT_EXPIRY_HOURS=24

# 跨域配置
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# WebSocket配置
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
```

### **Docker部署**
项目包含完整的Docker配置：

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**构建和运行**:
```bash
docker build -t guandan-backend .
docker run -p 8080:8080 guandan-backend
```

---

## 🔍 **开发和测试**

### **项目依赖**
```go
module guandan-world/backend

go 1.21

require (
    github.com/gin-contrib/cors v1.5.0
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/gorilla/websocket v1.5.1
    github.com/google/uuid v1.5.0
)
```

### **启动服务**
```bash
cd backend
go mod download
go run main.go
```

### **测试覆盖**
每个核心组件都包含对应的测试文件：
- `auth/service_test.go` - 认证服务测试
- `room/service_test.go` - 房间服务测试  
- `game/service_test.go` - 游戏服务测试
- `websocket/manager_test.go` - WebSocket测试
- `handlers/auth_test.go` - 认证接口测试
- `handlers/room_test.go` - 房间接口测试

**运行测试**:
```bash
go test ./...
```

---

## 🎯 **技术特色**

### **1. 🔄 实时同步架构**
- WebSocket双向通信确保游戏状态实时同步
- 房间级消息广播，支持4人同时在线
- 断线重连机制（待实现）

### **2. 🛡️ 安全机制**
- JWT令牌认证，防止未授权访问
- CORS配置，支持跨域前端访问
- 输入验证和错误处理

### **3. 📈 可扩展设计**
- 接口驱动的模块化架构
- 游戏引擎与在线服务解耦
- 支持水平扩展（添加负载均衡器）

### **4. 🎮 游戏引擎集成**
- 无缝集成现有的Game Engine SDK
- 保持游戏逻辑的纯净性
- 支持AI玩家和人类玩家混合对战

---

## 📝 **下一步发展方向**

### **短期目标**
- [ ] 前端UI开发和集成
- [ ] 断线重连机制
- [ ] 游戏回放功能
- [ ] 观战模式

### **中期目标**
- [ ] 用户注册和持久化存储
- [ ] 排行榜和统计系统
- [ ] 房间密码保护
- [ ] 私人房间功能

### **长期目标**
- [ ] 锦标赛模式
- [ ] 移动端适配
- [ ] 游戏直播功能
- [ ] 社交功能（好友、聊天）

---

## 🎉 **总结**

掼蛋在线游戏后端是一个功能完整、架构优雅的多人在线游戏服务。它成功地将单机游戏引擎升级为在线多人平台，具备以下优势：

- **🎯 专业架构**：清晰的分层设计和职责分离
- **⚡ 高性能**：WebSocket实时通信和高效的状态管理
- **🛡️ 安全可靠**：完整的认证机制和错误处理
- **🔧 易于维护**：模块化设计和完整的测试覆盖
- **📈 可扩展**：支持未来功能扩展和性能优化

这个后端服务为掼蛋游戏提供了坚实的技术基础，能够支撑从小规模测试到大规模商业运营的各种需求。 