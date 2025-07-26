# 后端集成测试完整指南

## 📋 测试架构概览

基于当前的测试基础设施，我们设计了一个完整的集成测试框架，涵盖端到端测试、性能测试和可靠性测试。

## 🏗️ 集成测试分层架构

```
集成测试层级
├── Level 1: API接口集成测试
├── Level 2: 服务间协调测试  
├── Level 3: 端到端游戏流程测试
├── Level 4: WebSocket实时通信测试
├── Level 5: 性能和负载测试
└── Level 6: 错误恢复和稳定性测试
```

---

## 🎯 Level 1: API接口集成测试

### 1.1 HTTP REST API集成测试

创建 `backend/integration_tests/api_integration_test.go`:

```go
package integration_tests

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    
    "guandan-world/backend/auth"
    "guandan-world/backend/game"
    "guandan-world/backend/handlers"
    "guandan-world/backend/room"
    "guandan-world/backend/websocket"
)

type APIIntegrationTestSuite struct {
    suite.Suite
    router      *gin.Engine
    authService auth.AuthService
    roomService room.RoomService
    gameService *game.GameService
    wsManager   *websocket.WSManager
    server      *httptest.Server
    userTokens  map[string]string  // username -> token
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
    // 初始化所有服务
    suite.authService = auth.NewAuthService()
    suite.roomService = room.NewRoomService(suite.authService)
    suite.wsManager = websocket.NewWSManager(suite.authService, suite.roomService)
    suite.gameService = game.NewGameService(suite.wsManager)
    
    // 创建路由
    gin.SetMode(gin.TestMode)
    suite.router = gin.New()
    
    // 注册所有路由
    suite.setupRoutes()
    
    // 启动测试服务器
    suite.server = httptest.NewServer(suite.router)
    
    // 创建测试用户
    suite.createTestUsers()
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
    suite.server.Close()
}

func (suite *APIIntegrationTestSuite) TestCompleteUserFlow() {
    // 1. 用户注册
    users := []string{"alice", "bob", "charlie", "david"}
    for _, username := range users {
        suite.registerUser(username, "password123")
    }
    
    // 2. 用户登录
    for _, username := range users {
        token := suite.loginUser(username, "password123")
        suite.userTokens[username] = token
        suite.NotEmpty(token)
    }
    
    // 3. 创建房间
    roomID := suite.createRoom("alice")
    suite.NotEmpty(roomID)
    
    // 4. 其他玩家加入房间
    for _, username := range []string{"bob", "charlie", "david"} {
        suite.joinRoom(username, roomID)
    }
    
    // 5. 验证房间状态
    room := suite.getRoom(roomID)
    suite.Equal(4, len(room.Players))
    suite.Equal("waiting", room.Status)
    
    // 6. 开始游戏
    suite.startGame("alice", roomID)
    
    // 7. 验证游戏已启动
    room = suite.getRoom(roomID)
    suite.Equal("playing", room.Status)
}

func (suite *APIIntegrationTestSuite) TestRoomManagement() {
    // 测试房间的完整生命周期
    
    // 1. 创建多个房间
    rooms := make([]string, 3)
    for i := 0; i < 3; i++ {
        rooms[i] = suite.createRoom("alice")
    }
    
    // 2. 获取房间列表
    roomList := suite.getRoomList(1, 10, nil)
    suite.GreaterOrEqual(len(roomList.Rooms), 3)
    
    // 3. 按状态过滤房间
    waitingStatus := room.RoomStatusWaiting
    waitingRooms := suite.getRoomList(1, 10, &waitingStatus)
    for _, r := range waitingRooms.Rooms {
        suite.Equal(room.RoomStatusWaiting, r.Status)
    }
    
    // 4. 测试房间容量限制
    for i := 0; i < 4; i++ {
        username := fmt.Sprintf("player%d", i+1)
        if i < 3 {
            suite.joinRoom(username, rooms[0])
        } else {
            // 第4个玩家加入应该成功，第5个应该失败
            suite.joinRoom(username, rooms[0])
        }
    }
    
    // 5. 测试房间已满的情况
    resp := suite.makeRequest("POST", "/api/rooms/join", map[string]interface{}{
        "room_id": rooms[0],
    }, "player5")
    suite.Equal(http.StatusBadRequest, resp.Code)
}

// Helper方法
func (suite *APIIntegrationTestSuite) registerUser(username, password string) {
    data := map[string]string{
        "username": username,
        "password": password,
    }
    resp := suite.makeRequest("POST", "/api/auth/register", data, "")
    suite.Equal(http.StatusOK, resp.Code)
}

func (suite *APIIntegrationTestSuite) makeRequest(method, path string, data interface{}, token string) *httptest.ResponseRecorder {
    jsonData, _ := json.Marshal(data)
    req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    
    resp := httptest.NewRecorder()
    suite.router.ServeHTTP(resp, req)
    return resp
}

func TestAPIIntegrationSuite(t *testing.T) {
    suite.Run(t, new(APIIntegrationTestSuite))
}
```

### 1.2 认证流程集成测试

```go
func (suite *APIIntegrationTestSuite) TestAuthenticationFlow() {
    // 1. 测试用户注册
    resp := suite.makeRequest("POST", "/api/auth/register", map[string]string{
        "username": "newuser",
        "password": "securepass123",
    }, "")
    suite.Equal(http.StatusOK, resp.Code)
    
    // 2. 测试重复注册
    resp = suite.makeRequest("POST", "/api/auth/register", map[string]string{
        "username": "newuser",
        "password": "anotherpass",
    }, "")
    suite.Equal(http.StatusConflict, resp.Code)
    
    // 3. 测试登录
    resp = suite.makeRequest("POST", "/api/auth/login", map[string]string{
        "username": "newuser",
        "password": "securepass123",
    }, "")
    suite.Equal(http.StatusOK, resp.Code)
    
    // 4. 解析token
    var loginResp map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &loginResp)
    token := loginResp["token"].(map[string]interface{})["token"].(string)
    
    // 5. 测试受保护的接口
    resp = suite.makeRequest("GET", "/api/auth/me", nil, token)
    suite.Equal(http.StatusOK, resp.Code)
    
    // 6. 测试登出
    resp = suite.makeRequest("POST", "/api/auth/logout", nil, token)
    suite.Equal(http.StatusOK, resp.Code)
    
    // 7. 验证token已失效
    resp = suite.makeRequest("GET", "/api/auth/me", nil, token)
    suite.Equal(http.StatusUnauthorized, resp.Code)
}
```

---

## 🎯 Level 2: WebSocket实时通信集成测试

### 2.1 WebSocket连接管理测试

创建 `backend/integration_tests/websocket_integration_test.go`:

```go
package integration_tests

import (
    "encoding/json"
    "net/url"
    "testing"
    "time"
    
    "github.com/gorilla/websocket"
    "github.com/stretchr/testify/suite"
    
    "guandan-world/backend/websocket"
)

type WebSocketIntegrationTestSuite struct {
    suite.Suite
    server    *httptest.Server
    wsManager *websocket.WSManager
    users     map[string]*TestUser
}

type TestUser struct {
    Username string
    Token    string
    Conn     *websocket.Conn
    Messages chan websocket.WSMessage
}

func (suite *WebSocketIntegrationTestSuite) SetupSuite() {
    // 初始化服务和服务器
    suite.setupServer()
    suite.users = make(map[string]*TestUser)
    
    // 创建测试用户
    usernames := []string{"alice", "bob", "charlie", "david"}
    for _, username := range usernames {
        user := suite.createTestUser(username)
        suite.users[username] = user
    }
}

func (suite *WebSocketIntegrationTestSuite) TestWebSocketConnection() {
    // 1. 测试WebSocket连接建立
    for username, user := range suite.users {
        suite.connectWebSocket(user)
        suite.NotNil(user.Conn, "User %s should have WebSocket connection", username)
    }
    
    // 2. 测试心跳机制
    suite.testHeartbeat()
    
    // 3. 测试消息广播
    suite.testMessageBroadcast()
    
    // 4. 测试连接断开和清理
    suite.testConnectionCleanup()
}

func (suite *WebSocketIntegrationTestSuite) testHeartbeat() {
    user := suite.users["alice"]
    
    // 发送ping消息
    pingMsg := websocket.WSMessage{
        Type: "ping",
        Data: map[string]interface{}{
            "timestamp": time.Now().Format(time.RFC3339),
        },
    }
    
    suite.sendMessage(user, pingMsg)
    
    // 等待pong响应
    select {
    case msg := <-user.Messages:
        suite.Equal("pong", msg.Type)
    case <-time.After(5 * time.Second):
        suite.Fail("Did not receive pong response")
    }
}

func (suite *WebSocketIntegrationTestSuite) testMessageBroadcast() {
    // 1. 所有用户加入同一个房间
    roomID := "test-room-broadcast"
    for _, user := range suite.users {
        joinMsg := websocket.WSMessage{
            Type: "join_room",
            Data: map[string]interface{}{
                "room_id": roomID,
            },
        }
        suite.sendMessage(user, joinMsg)
    }
    
    // 2. 一个用户发送消息
    alice := suite.users["alice"]
    testMsg := websocket.WSMessage{
        Type: "test_broadcast",
        Data: map[string]interface{}{
            "content": "Hello everyone!",
        },
    }
    
    // 广播消息
    suite.wsManager.BroadcastToRoom(roomID, &testMsg)
    
    // 3. 验证所有用户都收到消息
    for username, user := range suite.users {
        select {
        case msg := <-user.Messages:
            suite.Equal("test_broadcast", msg.Type)
            suite.Equal("Hello everyone!", msg.Data.(map[string]interface{})["content"])
        case <-time.After(2 * time.Second):
            suite.Fail("User %s did not receive broadcast message", username)
        }
    }
}
```

---

## 🎯 Level 3: 端到端游戏流程集成测试

### 3.1 完整游戏循环测试

创建 `backend/integration_tests/game_flow_integration_test.go`:

```go
package integration_tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/suite"
    
    "guandan-world/backend/test"
    "guandan-world/ai"
    "guandan-world/sdk"
)

type GameFlowIntegrationTestSuite struct {
    suite.Suite
    tester    *test.APIGameTester
    serverURL string
    authToken string
}

func (suite *GameFlowIntegrationTestSuite) SetupSuite() {
    // 启动后端服务器
    suite.startBackendServer()
    
    // 获取认证token
    suite.authToken = suite.getAuthToken()
    
    // 创建API游戏测试器
    suite.tester = test.NewAPIGameTester(suite.serverURL, suite.authToken, true)
}

func (suite *GameFlowIntegrationTestSuite) TestCompleteGameFlow() {
    // 1. 开始游戏
    roomID := "integration-test-room"
    err := suite.tester.StartGame(roomID)
    suite.NoError(err)
    
    // 2. 运行完整游戏直到结束
    err = suite.tester.RunUntilComplete()
    suite.NoError(err)
    
    // 3. 验证游戏结果
    suite.verifyGameResults()
}

func (suite *GameFlowIntegrationTestSuite) TestMultipleGames() {
    // 测试连续多局游戏
    const numGames = 3
    
    for i := 0; i < numGames; i++ {
        roomID := fmt.Sprintf("multi-game-test-%d", i)
        
        err := suite.tester.StartGame(roomID)
        suite.NoError(err, "Game %d should start successfully", i+1)
        
        err = suite.tester.RunUntilComplete()
        suite.NoError(err, "Game %d should complete successfully", i+1)
        
        // 清理资源
        suite.tester.Close()
        suite.tester = test.NewAPIGameTester(suite.serverURL, suite.authToken, true)
    }
}

func (suite *GameFlowIntegrationTestSuite) TestGameWithTimeouts() {
    // 测试超时处理
    
    // 创建自定义配置的测试器，使用较短的超时时间
    config := &test.TesterConfig{
        PlayTimeout:    2 * time.Second,
        TributeTimeout: 1 * time.Second,
        Verbose:        true,
    }
    
    tester := test.NewAPIGameTesterWithConfig(suite.serverURL, suite.authToken, config)
    defer tester.Close()
    
    err := tester.StartGame("timeout-test-room")
    suite.NoError(err)
    
    // 让一些操作超时，验证系统的超时处理机制
    err = tester.RunUntilComplete()
    suite.NoError(err) // 即使有超时，游戏也应该能继续
}
```

### 3.2 AI对战集成测试

```go
func (suite *GameFlowIntegrationTestSuite) TestAIBattle() {
    // 使用不同的AI算法进行对战测试
    
    algorithms := map[int]ai.AutoPlayAlgorithm{
        0: ai.NewSmartAutoPlayAlgorithm(),      // 智能算法
        1: ai.NewRandomAutoPlayAlgorithm(),     // 随机算法
        2: ai.NewConservativeAutoPlayAlgorithm(), // 保守算法
        3: ai.NewAggressiveAutoPlayAlgorithm(),   // 激进算法
    }
    
    tester := test.NewAPIGameTester(suite.serverURL, suite.authToken, true)
    defer tester.Close()
    
    // 设置不同的AI算法
    for seat, algorithm := range algorithms {
        tester.SetAIAlgorithm(seat, algorithm)
    }
    
    err := tester.StartGame("ai-battle-room")
    suite.NoError(err)
    
    err = tester.RunUntilComplete()
    suite.NoError(err)
    
    // 验证所有AI都有合理的游戏表现
    stats := tester.GetGameStats()
    for seat, stat := range stats {
        suite.Greater(stat.CardsPlayed, 0, "AI %d should have played cards", seat)
        suite.GreaterOrEqual(stat.AverageDecisionTime.Milliseconds(), int64(0))
    }
}
```

---

## 🎯 Level 4: 性能和负载测试

### 4.1 并发连接测试

创建 `backend/integration_tests/performance_test.go`:

```go
package integration_tests

import (
    "sync"
    "testing"
    "time"
    
    "github.com/stretchr/testify/suite"
)

type PerformanceTestSuite struct {
    suite.Suite
    serverURL string
}

func (suite *PerformanceTestSuite) TestConcurrentConnections() {
    // 测试并发WebSocket连接
    const numConnections = 100
    
    var wg sync.WaitGroup
    errors := make(chan error, numConnections)
    
    for i := 0; i < numConnections; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            // 创建用户和连接
            username := fmt.Sprintf("user%d", userID)
            token := suite.createUserAndGetToken(username)
            
            conn, err := suite.createWebSocketConnection(token)
            if err != nil {
                errors <- err
                return
            }
            defer conn.Close()
            
            // 保持连接一段时间
            time.Sleep(5 * time.Second)
            
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    // 检查错误
    var errorCount int
    for err := range errors {
        if err != nil {
            suite.T().Logf("Connection error: %v", err)
            errorCount++
        }
    }
    
    // 允许少量连接失败（<5%）
    suite.Less(errorCount, numConnections/20, "Too many connection failures")
}

func (suite *PerformanceTestSuite) TestConcurrentGames() {
    // 测试并发游戏
    const numGames = 10
    
    var wg sync.WaitGroup
    gameResults := make(chan bool, numGames)
    
    for i := 0; i < numGames; i++ {
        wg.Add(1)
        go func(gameID int) {
            defer wg.Done()
            
            roomID := fmt.Sprintf("perf-test-room-%d", gameID)
            token := suite.getValidToken()
            
            tester := test.NewAPIGameTester(suite.serverURL, token, false)
            defer tester.Close()
            
            err := tester.StartGame(roomID)
            if err != nil {
                gameResults <- false
                return
            }
            
            err = tester.RunUntilComplete()
            gameResults <- err == nil
            
        }(i)
    }
    
    wg.Wait()
    close(gameResults)
    
    // 统计成功率
    var successCount int
    for success := range gameResults {
        if success {
            successCount++
        }
    }
    
    successRate := float64(successCount) / float64(numGames)
    suite.GreaterOrEqual(successRate, 0.8, "Game success rate should be at least 80%")
}
```

### 4.2 API响应时间测试

```go
func (suite *PerformanceTestSuite) TestAPIResponseTimes() {
    // 测试API响应时间
    endpoints := []struct {
        name   string
        method string
        path   string
        data   interface{}
    }{
        {"Login", "POST", "/api/auth/login", map[string]string{"username": "testuser", "password": "testpass"}},
        {"GetRooms", "GET", "/api/rooms", nil},
        {"CreateRoom", "POST", "/api/rooms/create", map[string]string{"name": "Test Room"}},
        {"GetMe", "GET", "/api/auth/me", nil},
    }
    
    token := suite.getValidToken()
    
    for _, endpoint := range endpoints {
        suite.Run(endpoint.name, func() {
            const numRequests = 100
            var totalDuration time.Duration
            
            for i := 0; i < numRequests; i++ {
                start := time.Now()
                resp := suite.makeRequest(endpoint.method, endpoint.path, endpoint.data, token)
                duration := time.Since(start)
                
                suite.Less(resp.Code, 500, "Request should not return server error")
                totalDuration += duration
            }
            
            avgDuration := totalDuration / numRequests
            suite.Less(avgDuration.Milliseconds(), int64(100), 
                "Average response time for %s should be < 100ms, got %v", 
                endpoint.name, avgDuration)
        })
    }
}
```

---

## 🎯 Level 5: 错误恢复和稳定性测试

### 5.1 断线重连测试

创建 `backend/integration_tests/reliability_test.go`:

```go
package integration_tests

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/suite"
)

type ReliabilityTestSuite struct {
    suite.Suite
    serverURL string
}

func (suite *ReliabilityTestSuite) TestDisconnectReconnect() {
    // 测试断线重连
    
    // 1. 建立连接并加入游戏
    token := suite.getValidToken()
    conn1, err := suite.createWebSocketConnection(token)
    suite.NoError(err)
    
    // 2. 加入房间并开始游戏
    roomID := "disconnect-test-room"
    suite.joinRoomAndStartGame(conn1, roomID)
    
    // 3. 模拟断线
    conn1.Close()
    
    // 4. 等待一段时间再重连
    time.Sleep(2 * time.Second)
    
    // 5. 重新连接
    conn2, err := suite.createWebSocketConnection(token)
    suite.NoError(err)
    defer conn2.Close()
    
    // 6. 验证能够恢复游戏状态
    gameState := suite.getGameState(conn2, roomID)
    suite.NotNil(gameState)
    suite.Equal("playing", gameState.Status)
}

func (suite *ReliabilityTestSuite) TestServerRestart() {
    // 测试服务器重启恢复
    
    // 1. 启动游戏
    tester := test.NewAPIGameTester(suite.serverURL, suite.getValidToken(), true)
    defer tester.Close()
    
    err := tester.StartGame("restart-test-room")
    suite.NoError(err)
    
    // 2. 模拟服务器重启（这里我们只是测试状态持久化）
    // 在实际环境中，这需要外部脚本来重启服务器
    
    // 3. 验证游戏状态能够恢复
    // （这需要数据库或其他持久化机制的支持）
}

func (suite *ReliabilityTestSuite) TestResourceCleanup() {
    // 测试资源清理
    
    const numGames = 20
    
    for i := 0; i < numGames; i++ {
        roomID := fmt.Sprintf("cleanup-test-%d", i)
        token := suite.getValidToken()
        
        tester := test.NewAPIGameTester(suite.serverURL, token, false)
        
        err := tester.StartGame(roomID)
        suite.NoError(err)
        
        // 快速结束游戏
        tester.Close()
        
        // 验证资源被正确清理
        // 这里可以检查内存使用、连接数等指标
    }
    
    // 验证系统仍然稳定
    finalTester := test.NewAPIGameTester(suite.serverURL, suite.getValidToken(), true)
    defer finalTester.Close()
    
    err := finalTester.StartGame("final-cleanup-test")
    suite.NoError(err)
}
```

---

## 🚀 实施步骤

### Step 1: 创建集成测试目录结构

```bash
mkdir -p backend/integration_tests
cd backend/integration_tests

# 创建测试文件
touch api_integration_test.go
touch websocket_integration_test.go  
touch game_flow_integration_test.go
touch performance_test.go
touch reliability_test.go
```

### Step 2: 安装依赖

```bash
cd backend
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/assert
go get github.com/gorilla/websocket
```

### Step 3: 创建测试运行脚本

创建 `backend/run-integration-tests.sh`:

```bash
#!/bin/bash

echo "🚀 启动后端集成测试"

# 设置环境变量
export GIN_MODE=test
export TEST_MODE=integration

# 启动后端服务器
echo "📡 启动后端服务器..."
go build -o backend-test main.go
./backend-test &
SERVER_PID=$!

# 等待服务器启动
sleep 3

# 检查服务器是否启动成功
if curl -s http://localhost:8080/healthz > /dev/null; then
    echo "✅ 后端服务器启动成功"
else
    echo "❌ 后端服务器启动失败"
    kill $SERVER_PID
    exit 1
fi

# 运行集成测试
echo "🧪 运行集成测试..."

# Level 1: API集成测试
echo "📊 Level 1: API接口集成测试"
go test -v ./integration_tests -run "APIIntegration" -timeout 10m

# Level 2: WebSocket测试  
echo "🔌 Level 2: WebSocket实时通信测试"
go test -v ./integration_tests -run "WebSocketIntegration" -timeout 5m

# Level 3: 游戏流程测试
echo "🎮 Level 3: 端到端游戏流程测试"
go test -v ./integration_tests -run "GameFlowIntegration" -timeout 15m

# Level 4: 性能测试
echo "⚡ Level 4: 性能和负载测试"
go test -v ./integration_tests -run "Performance" -timeout 10m

# Level 5: 稳定性测试
echo "🛡️ Level 5: 错误恢复和稳定性测试"
go test -v ./integration_tests -run "Reliability" -timeout 10m

# 清理
echo "🧹 清理测试环境..."
kill $SERVER_PID
rm -f backend-test

echo "✅ 集成测试完成"
```

### Step 4: 创建持续集成配置

创建 `.github/workflows/integration-tests.yml`:

```yaml
name: Backend Integration Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Install dependencies
      run: |
        cd backend
        go mod download
        
    - name: Run unit tests first
      run: |
        cd backend
        go test ./auth ./game ./handlers ./room -v
        
    - name: Run integration tests
      run: |
        cd backend
        chmod +x run-integration-tests.sh
        ./run-integration-tests.sh
        
    - name: Upload test results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: integration-test-results
        path: backend/test-results/
```

---

## 📊 监控和报告

### 测试覆盖率报告

创建 `backend/generate-coverage-report.sh`:

```bash
#!/bin/bash

echo "📊 生成集成测试覆盖率报告"

# 运行所有测试并生成覆盖率
go test ./... -coverprofile=coverage.out -covermode=atomic

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 生成文本报告
go tool cover -func=coverage.out

echo "📈 覆盖率报告已生成："
echo "- coverage.out (原始数据)"
echo "- coverage.html (HTML报告)"
```

### 性能基准测试

创建 `backend/benchmark_test.go`:

```go
package main

import (
    "testing"
    "time"
    
    "guandan-world/backend/test"
)

func BenchmarkAPIGameFlow(b *testing.B) {
    serverURL := "http://localhost:8080"
    token := getValidToken() // 实现获取token的方法
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        tester := test.NewAPIGameTester(serverURL, token, false)
        
        roomID := fmt.Sprintf("benchmark-room-%d", i)
        tester.StartGame(roomID)
        tester.RunUntilComplete()
        tester.Close()
    }
}

func BenchmarkConcurrentConnections(b *testing.B) {
    // 基准测试并发连接性能
    // ...
}
```

---

## 🎯 执行方式

### 开发环境执行

```bash
# 1. 启动后端服务
cd backend
go run main.go &

# 2. 运行特定级别的测试
go test -v ./integration_tests -run "APIIntegration" -timeout 10m

# 3. 运行完整集成测试套件
./run-integration-tests.sh

# 4. 运行性能基准测试
go test -bench=. -benchmem
```

### 生产环境验证

```bash
# 1. 指向生产环境
export TEST_SERVER_URL=https://api.guandan.com
export TEST_AUTH_TOKEN=<production_token>

# 2. 运行只读测试
go test -v ./integration_tests -run "APIIntegration.*Read" -timeout 5m

# 3. 运行性能监控测试
go test -v ./integration_tests -run "Performance.*Monitor" -timeout 15m
```

---

## 📈 测试指标和目标

### 性能指标

| 指标 | 目标值 | 当前值 | 状态 |
|------|--------|--------|------|
| API响应时间 | < 100ms | TBD | 🔍 |
| WebSocket连接数 | > 1000 | TBD | 🔍 |
| 游戏并发数 | > 100 | TBD | 🔍 |
| 成功率 | > 99% | TBD | 🔍 |
| 内存使用 | < 1GB | TBD | 🔍 |

### 稳定性指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| 断线重连成功率 | > 95% | 网络中断后能恢复 |
| 资源清理率 | 100% | 无内存泄漏 |
| 错误恢复时间 | < 5s | 从错误状态恢复 |
| 并发游戏稳定性 | > 99% | 多游戏并行运行 |

这个完整的集成测试框架将确保后端系统的可靠性、性能和稳定性，为前端开发和生产部署提供坚实的保障！ 