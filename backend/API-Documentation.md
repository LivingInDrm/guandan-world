# 掼蛋在线对战平台 - 后端接口文档

## 概述

本文档描述了掼蛋在线对战平台后端的完整API接口，包括REST API和WebSocket API。

**服务器地址：** `http://localhost:8080`

**API版本：** v1

## 认证系统

所有需要认证的接口都需要在请求头中包含JWT token：

```
Authorization: Bearer <token>
```

### 数据结构

#### User
```json
{
  "id": "string",
  "username": "string", 
  "online": "boolean"
}
```

#### AuthToken
```json
{
  "token": "string",
  "expires_at": "string (ISO 8601)",
  "user_id": "string"
}
```

#### AuthResponse
```json
{
  "user": "User",
  "token": "AuthToken"
}
```

#### ErrorResponse
```json
{
  "error": "string",
  "message": "string"
}
```

---

## REST API 接口

### 1. 认证接口

#### 1.1 用户注册
**POST** `/api/auth/register`

**请求体：**
```json
{
  "username": "string (required)",
  "password": "string (required, 最少6个字符)"
}
```

**响应：**
- **201 Created** - 注册成功并自动登录
  ```json
  {
    "user": "User",
    "token": "AuthToken"
  }
  ```
- **400 Bad Request** - 请求格式错误
- **409 Conflict** - 用户名已存在

#### 1.2 用户登录
**POST** `/api/auth/login`

**请求体：**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**响应：**
- **200 OK** - 登录成功
  ```json
  {
    "user": "User",
    "token": "AuthToken"
  }
  ```
- **401 Unauthorized** - 认证失败

#### 1.3 用户登出
**POST** `/api/auth/logout`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 登出成功
  ```json
  {
    "message": "Successfully logged out"
  }
  ```
- **400 Bad Request** - Token无效

#### 1.4 获取当前用户信息
**GET** `/api/auth/me`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 成功
  ```json
  {
    "user": "User"
  }
  ```
- **401 Unauthorized** - 未认证

---

### 2. 房间管理接口

> **注意：** 以下房间接口已在代码中实现，但需要在 `main.go` 中注册路由才能使用。

#### 2.1 房间数据结构

#### Player
```json
{
  "id": "string",
  "username": "string",
  "seat": "number (0-3)",
  "online": "boolean"
}
```

#### Room
```json
{
  "id": "string",
  "status": "string (waiting|ready|playing|closed)",
  "players": "Array[4] of Player (可为null)",
  "owner": "string (房主用户ID)",
  "player_count": "number",
  "created_at": "string (ISO 8601)",
  "updated_at": "string (ISO 8601)"
}
```

#### RoomInfo (用于房间列表)
```json
{
  "id": "string",
  "status": "string",
  "player_count": "number",
  "players": "Array of Player",
  "owner": "string",
  "can_join": "boolean",
  "created_at": "string (ISO 8601)"
}
```

#### 2.2 获取房间列表
**GET** `/api/rooms`

**请求头：** `Authorization: Bearer <token>`

**查询参数：**
- `page` (可选): 页码，默认1
- `limit` (可选): 每页数量，默认12，最大50
- `status` (可选): 房间状态过滤 (`waiting|ready|playing`)

**响应：**
- **200 OK** - 成功
  ```json
  {
    "rooms": "Array of RoomInfo",
    "total_count": "number",
    "page": "number", 
    "limit": "number"
  }
  ```

#### 2.3 创建房间
**POST** `/api/rooms`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **201 Created** - 房间创建成功
  ```json
  {
    "room": "Room"
  }
  ```
- **409 Conflict** - 用户已在其他房间中

#### 2.4 获取当前用户的房间
**GET** `/api/rooms/my`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 成功
  ```json
  {
    "room": "Room"
  }
  ```
- **404 Not Found** - 用户不在任何房间中

#### 2.5 获取指定房间信息
**GET** `/api/rooms/:id`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 成功
  ```json
  {
    "room": "Room"
  }
  ```
- **404 Not Found** - 房间不存在

#### 2.6 加入房间
**POST** `/api/rooms/:id/join`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 加入成功
  ```json
  {
    "room": "Room"
  }
  ```
- **404 Not Found** - 房间不存在
- **409 Conflict** - 房间已满或不接受新玩家

#### 2.7 离开房间
**POST** `/api/rooms/:id/leave`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 离开成功
  ```json
  {
    "room": "Room"
  }
  ```
  或房间已关闭时：
  ```json
  {
    "message": "Successfully left room (room was closed)"
  }
  ```
- **404 Not Found** - 房间不存在
- **409 Conflict** - 用户不在此房间中

#### 2.8 开始游戏
**POST** `/api/rooms/:id/start`

**请求头：** `Authorization: Bearer <token>`

**响应：**
- **200 OK** - 游戏开始成功
  ```json
  {
    "room": "Room"
  }
  ```
- **403 Forbidden** - 非房主无权限
- **409 Conflict** - 房间状态不允许开始游戏或人数不足

---

### 3. 健康检查接口

#### 3.1 健康检查
**GET** `/healthz`

**响应：**
- **200 OK**
  ```json
  {
    "status": "pong"
  }
  ```

---

## WebSocket API

### 连接建立

**WebSocket 端点：** `ws://localhost:8080/ws`

**认证方式：** 在连接建立时通过查询参数传递token
```
ws://localhost:8080/ws?token=<jwt_token>
```

### 消息格式

所有WebSocket消息都使用JSON格式：

```json
{
  "type": "string (消息类型)",
  "data": "object (消息数据)",
  "timestamp": "string (ISO 8601)",
  "player_id": "string (可选)"
}
```

### 1. 连接管理消息

#### 1.1 心跳检测
**客户端发送：**
```json
{
  "type": "ping",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**服务器响应：**
```json
{
  "type": "pong", 
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 2. 房间管理消息

#### 2.1 加入房间
**客户端发送：**
```json
{
  "type": "join_room",
  "data": {
    "room_id": "string"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 2.2 离开房间
**客户端发送：**
```json
{
  "type": "leave_room",
  "data": {
    "room_id": "string"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 2.3 开始游戏
**客户端发送：**
```json
{
  "type": "start_game",
  "data": {
    "room_id": "string"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 3. 游戏操作消息

#### 3.1 出牌
**客户端发送：**
```json
{
  "type": "play_cards",
  "data": {
    "cards": ["card_id1", "card_id2", "..."]
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 3.2 不出/过牌
**客户端发送：**
```json
{
  "type": "pass",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 3.3 上贡选牌
**客户端发送：**
```json
{
  "type": "tribute_select",
  "data": {
    "card_id": "string"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 3.4 还贡
**客户端发送：**
```json
{
  "type": "tribute_return",
  "data": {
    "card_id": "string"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 4. 服务器推送消息

#### 4.1 房间状态更新
**服务器推送：**
```json
{
  "type": "room_update",
  "data": {
    "room": "Room对象"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 4.2 游戏事件
**服务器推送：**
```json
{
  "type": "game_event",
  "data": {
    "event_type": "string",
    "event_data": "object (具体事件数据)"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 4.3 玩家视图更新
**服务器推送：**
```json
{
  "type": "player_view",
  "data": {
    "player_seat": "number",
    "player_cards": "Array (玩家手牌)",
    "visible_cards": "object (所有人可见的牌)",
    "game_status": "string",
    "team_levels": "object",
    "players": "Array (玩家信息)",
    "current_turn": "number (可选)",
    "trick_leader": "number (可选)",
    // ... 其他游戏状态信息
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 4.4 错误消息
**服务器推送：**
```json
{
  "type": "error",
  "data": {
    "error": "string (错误代码)",
    "message": "string (错误描述)",
    "context": "object (可选，错误上下文)"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## 错误处理

### HTTP 错误码

- **400 Bad Request** - 请求格式错误或参数无效
- **401 Unauthorized** - 未认证或token无效
- **403 Forbidden** - 权限不足
- **404 Not Found** - 资源不存在
- **409 Conflict** - 操作冲突（如用户名已存在、房间已满等）
- **500 Internal Server Error** - 服务器内部错误

### WebSocket 错误

WebSocket错误通过 `error` 类型消息发送，包含错误代码和描述。

### 常见错误代码

#### 认证相关
- `invalid_request` - 请求格式错误
- `authentication_failed` - 认证失败
- `username_exists` - 用户名已存在
- `invalid_token` - token无效
- `token_expired` - token已过期
- `missing_token` - 缺少token

#### 房间相关
- `room_not_found` - 房间不存在
- `room_full` - 房间已满
- `room_not_accepting` - 房间不接受新玩家
- `already_in_room` - 用户已在房间中
- `not_in_room` - 用户不在房间中
- `not_room_owner` - 非房主权限
- `room_not_ready` - 房间状态不允许操作
- `insufficient_players` - 人数不足

#### 游戏相关
- `game_not_found` - 游戏不存在
- `invalid_move` - 无效操作
- `not_your_turn` - 不是该玩家回合
- `invalid_cards` - 无效牌型

---

## 使用示例

### 1. 用户注册和登录流程

```javascript
// 1. 注册
const registerResponse = await fetch('/api/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'player1',
    password: 'password123'
  })
});

const { user, token } = await registerResponse.json();

// 2. 存储token用于后续请求
localStorage.setItem('authToken', token.token);
```

### 2. 房间操作流程

```javascript
// 获取房间列表
const roomsResponse = await fetch('/api/rooms?page=1&limit=12', {
  headers: { 'Authorization': `Bearer ${token}` }
});

// 创建房间
const createResponse = await fetch('/api/rooms', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` }
});

// 加入房间
const joinResponse = await fetch(`/api/rooms/${roomId}/join`, {
  method: 'POST', 
  headers: { 'Authorization': `Bearer ${token}` }
});
```

### 3. WebSocket 连接和游戏操作

```javascript
// 建立WebSocket连接
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch (message.type) {
    case 'room_update':
      // 处理房间状态更新
      updateRoomUI(message.data.room);
      break;
      
    case 'player_view':
      // 处理游戏状态更新
      updateGameUI(message.data);
      break;
      
    case 'error':
      // 处理错误
      showError(message.data.message);
      break;
  }
};

// 出牌
function playCards(cardIds) {
  ws.send(JSON.stringify({
    type: 'play_cards',
    data: { cards: cardIds },
    timestamp: new Date().toISOString()
  }));
}

// 不出
function pass() {
  ws.send(JSON.stringify({
    type: 'pass', 
    data: {},
    timestamp: new Date().toISOString()
  }));
}
```

---

## 注意事项

1. **房间路由注册**：当前房间相关的HTTP接口已实现但未在 `main.go` 中注册，需要添加以下代码才能使用：

```go
// 在 main.go 中添加
roomService := room.NewRoomService(authService)
roomHandler := handlers.NewRoomHandler(roomService, authService)
roomHandler.RegisterRoutes(r, authHandler)
```

2. **CORS配置**：当前CORS设置为允许所有来源，生产环境需要配置具体的允许来源。

3. **WebSocket认证**：WebSocket连接需要在查询参数中传递有效的JWT token。

4. **超时处理**：游戏中有20秒操作超时，连续2次超时将自动进入托管状态。

5. **数据持久化**：当前实现使用内存存储，重启后数据会丢失，生产环境需要数据库持久化。

6. **游戏状态同步**：游戏状态通过WebSocket实时同步，断线重连后需要重新获取完整状态。 