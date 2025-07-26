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
	"github.com/stretchr/testify/suite"

	"guandan-world/backend/auth"
	"guandan-world/backend/game"
	"guandan-world/backend/handlers"
	"guandan-world/backend/room"
	"guandan-world/backend/websocket"
)

// BasicIntegrationTestSuite 基础集成测试套件
type BasicIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	authService auth.AuthService
	roomService room.RoomService
	gameService *game.GameService
	wsManager   *websocket.WSManager
	server      *httptest.Server
	userTokens  map[string]string // username -> token
}

func (suite *BasicIntegrationTestSuite) SetupSuite() {
	// 初始化所有服务
	suite.authService = auth.NewAuthService("test-secret-key", 24*time.Hour)
	suite.roomService = room.NewRoomService(suite.authService)
	suite.wsManager = websocket.NewWSManager(suite.authService, suite.roomService)
	suite.gameService = game.NewGameService(suite.wsManager)

	// 创建路由
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.setupRoutes()

	// 启动测试服务器
	suite.server = httptest.NewServer(suite.router)
	suite.userTokens = make(map[string]string)

	fmt.Printf("🚀 测试服务器启动: %s\n", suite.server.URL)
}

func (suite *BasicIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
	fmt.Println("🧹 测试服务器已关闭")
}

func (suite *BasicIntegrationTestSuite) setupRoutes() {
	// 设置基本路由（简化版）
	authHandler := handlers.NewAuthHandler(suite.authService)
	roomHandler := handlers.NewRoomHandler(suite.roomService, suite.authService)

	// 注册认证路由
	auth := suite.router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.JWTMiddleware(), authHandler.Me)
	}

	// 注册房间路由
	rooms := suite.router.Group("/api/rooms")
	rooms.Use(authHandler.JWTMiddleware())
	{
		rooms.GET("/", roomHandler.GetRooms)
		rooms.POST("/create", roomHandler.CreateRoom)
		rooms.POST("/join", roomHandler.JoinRoom)
		rooms.POST("/leave", roomHandler.LeaveRoom)
		rooms.POST("/:id/start", roomHandler.StartGame)
		rooms.GET("/my", roomHandler.GetMyRoom)
	}

	// 健康检查
	suite.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Level 1: 基础API集成测试
func (suite *BasicIntegrationTestSuite) TestBasicAPIFlow() {
	fmt.Println("📊 Level 1: 基础API集成测试")

	// 1. 用户注册
	users := []string{"alice", "bob", "charlie", "david"}
	for _, username := range users {
		resp := suite.registerUser(username, "password123")
		suite.Equal(http.StatusCreated, resp.Code, "用户 %s 注册应该成功", username)
	}

	// 2. 用户登录
	for _, username := range users {
		token := suite.loginUser(username, "password123")
		suite.NotEmpty(token, "用户 %s 应该获得有效token", username)
		suite.userTokens[username] = token
	}

	// 3. 创建房间
	roomResp := suite.createRoom("alice")
	suite.Equal(http.StatusCreated, roomResp.Code, "创建房间应该成功")

	var roomData map[string]interface{}
	json.Unmarshal(roomResp.Body.Bytes(), &roomData)
	roomID := roomData["room"].(map[string]interface{})["id"].(string)

	// 4. 验证房间创建成功
	suite.NotEmpty(roomID, "房间ID应该不为空")

	fmt.Println("✅ 基础API流程测试通过")
}

// Level 2: 使用现有API游戏测试器进行完整流程测试
func (suite *BasicIntegrationTestSuite) TestCompleteGameFlow() {
	fmt.Println("🎮 Level 2: 完整游戏流程测试")

	// 注意：现有的APIGameTester需要一个真实运行的服务器，而不是测试服务器
	// 所以这个测试暂时跳过，直到我们有了完整的服务器集成
	suite.T().Skip("需要真实的后端服务器运行，当前使用的是测试服务器")

	fmt.Println("✅ 完整游戏流程测试通过")
}

// Level 3: 简单性能测试
func (suite *BasicIntegrationTestSuite) TestBasicPerformance() {
	fmt.Println("⚡ Level 3: 基础性能测试")

	token := suite.getValidToken()

	// 测试API响应时间
	endpoints := []struct {
		name   string
		method string
		path   string
	}{
		{"GetRooms", "GET", "/api/rooms"},
		{"GetMe", "GET", "/api/auth/me"},
	}

	for _, endpoint := range endpoints {
		const numRequests = 10
		var totalDuration time.Duration

		for i := 0; i < numRequests; i++ {
			start := time.Now()
			resp := suite.makeRequest(endpoint.method, endpoint.path, nil, token)
			duration := time.Since(start)

			suite.Less(resp.Code, 500, "请求不应该返回服务器错误")
			totalDuration += duration
		}

		avgDuration := totalDuration / numRequests
		fmt.Printf("📈 %s 平均响应时间: %v\n", endpoint.name, avgDuration)
		suite.Less(avgDuration.Milliseconds(), int64(200),
			"%s 平均响应时间应该 < 200ms", endpoint.name)
	}

	fmt.Println("✅ 基础性能测试通过")
}

// Helper方法
func (suite *BasicIntegrationTestSuite) registerUser(username, password string) *httptest.ResponseRecorder {
	data := map[string]string{
		"username": username,
		"password": password,
	}
	return suite.makeRequest("POST", "/api/auth/register", data, "")
}

func (suite *BasicIntegrationTestSuite) loginUser(username, password string) string {
	data := map[string]string{
		"username": username,
		"password": password,
	}
	resp := suite.makeRequest("POST", "/api/auth/login", data, "")
	if resp.Code != http.StatusOK {
		return ""
	}

	var loginResp map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &loginResp)
	tokenData := loginResp["token"].(map[string]interface{})
	return tokenData["token"].(string)
}

func (suite *BasicIntegrationTestSuite) createRoom(username string) *httptest.ResponseRecorder {
	data := map[string]string{
		"name": fmt.Sprintf("%s的房间", username),
	}
	return suite.makeRequest("POST", "/api/rooms/create", data, suite.userTokens[username])
}

func (suite *BasicIntegrationTestSuite) joinRoom(username, roomID string) *httptest.ResponseRecorder {
	data := map[string]string{
		"room_id": roomID,
	}
	return suite.makeRequest("POST", "/api/rooms/join", data, suite.userTokens[username])
}

func (suite *BasicIntegrationTestSuite) getRoom(roomID string) map[string]interface{} {
	resp := suite.makeRequest("GET", "/api/rooms", nil, suite.userTokens["alice"])

	if resp.Code != http.StatusOK {
		return nil
	}

	var roomList map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &roomList); err != nil {
		return nil
	}

	rooms, ok := roomList["rooms"]
	if !ok || rooms == nil {
		return nil
	}

	roomsArray, ok := rooms.([]interface{})
	if !ok {
		return nil
	}

	for _, room := range roomsArray {
		r, ok := room.(map[string]interface{})
		if !ok {
			continue
		}
		if id, exists := r["id"]; exists && id.(string) == roomID {
			return r
		}
	}
	return nil
}

func (suite *BasicIntegrationTestSuite) getValidToken() string {
	if len(suite.userTokens) == 0 {
		// 创建一个测试用户
		suite.registerUser("testuser", "testpass")
		token := suite.loginUser("testuser", "testpass")
		suite.userTokens["testuser"] = token
		return token
	}

	for _, token := range suite.userTokens {
		return token
	}
	return ""
}

func (suite *BasicIntegrationTestSuite) makeRequest(method, path string, data interface{}, token string) *httptest.ResponseRecorder {
	var body []byte
	if data != nil {
		body, _ = json.Marshal(data)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)
	return resp
}

// 运行测试套件
func TestBasicIntegrationSuite(t *testing.T) {
	suite.Run(t, new(BasicIntegrationTestSuite))
}
