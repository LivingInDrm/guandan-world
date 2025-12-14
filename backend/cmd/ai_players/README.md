# AI玩家测试客户端

这是一个自动化测试工具，用于模拟多个AI玩家加入游戏房间并自动进行游戏。

## 功能特性

- 自动注册/登录3个AI玩家账号
- 加入指定的游戏房间
- 使用AI算法自动出牌、过牌
- 自动处理上贡和还贡
- 完整的游戏流程支持

## 使用方法

### 基本用法

```bash
cd backend/cmd/ai_players
go run main.go -room-code <4位房间码>
```

### 命令行参数

| 参数 | 说明 | 默认值 | 必需 |
|------|------|--------|------|
| `-server` | 后端服务器地址(host:port) | localhost:8080 | 否 |
| `-room-code` | 要加入的房间码（4位数字） | 无 | **是** |
| `-verbose` | 启用详细日志输出 | false | 否 |
| `-num-players` | AI玩家数量(1-3) | 3 | 否 |
| `-username-prefix` | AI玩家用户名前缀 | ai_player | 否 |
| `-password` | AI玩家密码 | ai123456 | 否 |

### 示例

1. **基本测试（3个AI玩家）**
```bash
go run main.go -room-code 1234
```

2. **使用详细日志**
```bash
go run main.go -room-code 1234 -verbose
```

3. **自定义服务器地址**
```bash
go run main.go -server 192.168.1.100:8080 -room-code 1234
```

4. **只启动2个AI玩家**
```bash
go run main.go -room-code 1234 -num-players 2
```

## 完整测试流程

### 1. 启动后端服务器

```bash
cd backend
go run main.go
```

确保服务器运行在 `http://localhost:8080`

### 2. 在前端创建房间

1. 打开浏览器访问前端应用
2. 登录或注册账号
3. 创建一个新房间
4. 记下房间码（右上角显示的4位数字，例如：`1234`）

### 3. 启动AI玩家客户端

```bash
cd backend/cmd/ai_players
go run main.go -room-code 1234 -verbose
```

你将看到类似以下的输出：

```
=== AI Players Test Client ===
Server: localhost:8080
Room Code: 1234
Number of players: 3
Verbose: true
==============================

[ai_player_1] Starting AI player client...
[ai_player_1] Authentication successful
[ai_player_1] Joined room 1234 (ID: room_1734567890123), seat: 1
[ai_player_1] WebSocket connected
[ai_player_1] Sent join_room message
[ai_player_1] AI player client started, waiting for game...
[Client 1] Started successfully

[ai_player_2] Starting AI player client...
[ai_player_2] Authentication successful
[ai_player_2] Joined room 1234 (ID: room_1734567890123), seat: 2
[ai_player_2] WebSocket connected
[ai_player_2] Sent join_room message
[ai_player_2] AI player client started, waiting for game...
[Client 2] Started successfully

[ai_player_3] Starting AI player client...
[ai_player_3] Authentication successful
[ai_player_3] Joined room 1234 (ID: room_1734567890123), seat: 3
[ai_player_3] WebSocket connected
[ai_player_3] Sent join_room message
[ai_player_3] AI player client started, waiting for game...
[Client 3] Started successfully

All 3 AI players connected and ready!
Waiting for game to start and complete...
Press Ctrl+C to exit
```

### 4. 开始游戏

在前端界面中，作为房主点击"开始游戏"按钮。

### 5. 观察游戏进行

- AI玩家将自动出牌和过牌
- 如果启用了`-verbose`，你将看到详细的决策日志
- 游戏完成后，AI客户端会自动退出

## 日志说明

### 基本日志（默认）

- 账号注册/登录状态
- 加入房间成功
- WebSocket连接状态
- 游戏操作请求（出牌、贡牌等）
- 错误信息

### 详细日志（-verbose）

额外包括：
- AI决策详情（打出的牌、选择过牌）
- 手牌更新
- 游戏事件详情

## 故障排查

### 问题：无法连接到服务器

**原因**：后端服务器未启动或地址不正确

**解决**：
- 确认后端服务器正在运行
- 检查服务器地址和端口是否正确

### 问题：房间码不存在

**原因**：房间码输入错误或房间已关闭

**解决**：
- 在前端重新创建房间
- 确认房间码拼写正确

### 问题：房间已满

**原因**：房间已有4个玩家

**解决**：
- 减少AI玩家数量：`-num-players 2`
- 或创建新房间

### 问题：AI账号已在其他房间

**原因**：之前的测试未正常退出

**解决**：
- 使用不同的用户名前缀：`-username-prefix test_ai`
- 或重启后端服务器（开发环境）

## 技术细节

### AI算法

使用 `ai.SmartAutoPlayAlgorithm`，从2级开始：
- 智能出牌策略
- 考虑牌型组合
- 优先清空手牌

### 通信协议

- **HTTP API**：用于注册、登录、加入房间、提交决策
- **WebSocket**：用于接收游戏事件和操作请求

### 并发处理

- 每个AI玩家独立运行在单独的goroutine中
- 使用互斥锁保护共享状态
- 异步提交决策，避免阻塞

## 相关文件

- `backend/test/ai_player_client.go` - AI玩家客户端实现
- `backend/test/api_game_tester.go` - Game Driver测试工具
- `ai/smart_algorithm.go` - AI算法实现

## 许可

与主项目相同

