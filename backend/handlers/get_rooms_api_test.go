package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/room"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============ 1. 正常流程测试 ============

func TestGetRoomsAPI_DefaultParameters(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建3个用户和房间
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("getrooms_user%d", i))
		tokens[i] = token

		// 创建房间
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 获取房间列表（使用默认参数）
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, len(response.Rooms))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 12, response.Limit)

	// 验证每个房间信息的完整性
	for _, roomInfo := range response.Rooms {
		assert.NotEmpty(t, roomInfo.ID)
		assert.Equal(t, room.RoomStatusWaiting, roomInfo.Status)
		assert.Equal(t, 1, roomInfo.PlayerCount)
		assert.NotEmpty(t, roomInfo.Owner)
		assert.True(t, roomInfo.CanJoin)
		assert.False(t, roomInfo.CreatedAt.IsZero())
		assert.NotNil(t, roomInfo.Players)
		assert.Equal(t, 1, len(roomInfo.Players))
	}
}

func TestGetRoomsAPI_EmptyList(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建用户但不创建房间
	token, _ := createTestUserAndLogin(t, router, "getrooms_empty_user")

	// 获取房间列表
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 0, response.TotalCount)
	assert.Equal(t, 0, len(response.Rooms))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 12, response.Limit)
}

func TestGetRoomsAPI_RoomSorting(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 12)
	for i := 0; i < 12; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("sort_user%d", i))
		tokens[i] = token
	}

	// 1. 创建1个waiting房间（1人）
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	waitingRoom1ID := ""
	var resp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	waitingRoom1ID = resp.Room.ID

	time.Sleep(10 * time.Millisecond) // 确保时间戳不同

	// 2. 创建1个waiting房间（3人）
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	waitingRoom3ID := resp.Room.ID

	// 加入2个玩家
	for i := 2; i < 4; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", waitingRoom3ID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	time.Sleep(10 * time.Millisecond)

	// 3. 创建1个ready房间（4人）
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	readyRoomID := resp.Room.ID

	for i := 5; i < 8; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", readyRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	time.Sleep(10 * time.Millisecond)

	// 4. 创建1个playing房间
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[8])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	playingRoomID := resp.Room.ID

	for i := 9; i < 12; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", playingRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 开始游戏
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", playingRoomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[8])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 获取房间列表
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 4, response.TotalCount)
	assert.Equal(t, 4, len(response.Rooms))

	// 验证排序: waiting房间应该在前面
	// waiting房间中，玩家数多的排在前面
	// waiting(3人) > waiting(1人) > ready(4人) > playing(4人)

	// 第1个应该是waiting且玩家最多的
	assert.Equal(t, room.RoomStatusWaiting, response.Rooms[0].Status)
	assert.Equal(t, 3, response.Rooms[0].PlayerCount)
	assert.Equal(t, waitingRoom3ID, response.Rooms[0].ID)

	// 第2个应该是waiting且玩家较少的
	assert.Equal(t, room.RoomStatusWaiting, response.Rooms[1].Status)
	assert.Equal(t, 1, response.Rooms[1].PlayerCount)
	assert.Equal(t, waitingRoom1ID, response.Rooms[1].ID)

	// 后面是ready和playing（顺序可能因为玩家数相同而由创建时间决定）
	// 但都应该不是waiting状态
	assert.NotEqual(t, room.RoomStatusWaiting, response.Rooms[2].Status)
	assert.NotEqual(t, room.RoomStatusWaiting, response.Rooms[3].Status)
}

// ============ 2. 认证相关测试 ============

func TestGetRoomsAPI_NoAuthToken(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 不携带Authorization header
	req := httptest.NewRequest("GET", "/api/rooms", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "missing_token", errorResp.Error)
}

func TestGetRoomsAPI_InvalidToken(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 使用无效token
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-12345")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetRoomsAPI_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个很短过期时间的auth service
	authService := auth.NewTestAuthService("test-secret", 1*time.Millisecond)
	roomService := room.NewRoomService(authService)
	mockDriverService := &MockDriverService{}
	mockWSManager := &MockWSManager{}

	authHandler := NewAuthHandler(authService)
	roomHandler := NewRoomHandler(roomService, authService, mockDriverService, mockWSManager)

	router := gin.New()
	authHandler.RegisterRoutes(router)
	
	// Register room routes manually
	rooms := router.Group("/api/rooms")
	rooms.Use(authHandler.JWTMiddleware())
	{
		rooms.GET("", roomHandler.GetRooms)
	}

	// 注册并登录用户获取token
	_, _ = authService.Register("expireduser", "password123")
	authToken, _ := authService.Login("expireduser", "password123")
	token := authToken.AccessToken

	// 等待token过期
	time.Sleep(10 * time.Millisecond)

	// 使用过期token
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ============ 3. 分页功能测试 ============

func TestGetRoomsAPI_Pagination(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建10个房间
	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("page_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 测试第1页（limit=3）
	req := httptest.NewRequest("GET", "/api/rooms?page=1&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp1 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp1)
	assert.Equal(t, 10, resp1.TotalCount)
	assert.Equal(t, 3, len(resp1.Rooms))
	assert.Equal(t, 1, resp1.Page)
	assert.Equal(t, 3, resp1.Limit)

	// 测试第2页
	req = httptest.NewRequest("GET", "/api/rooms?page=2&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp2 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp2)
	assert.Equal(t, 10, resp2.TotalCount)
	assert.Equal(t, 3, len(resp2.Rooms))
	assert.Equal(t, 2, resp2.Page)

	// 测试第3页
	req = httptest.NewRequest("GET", "/api/rooms?page=3&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp3 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp3)
	assert.Equal(t, 10, resp3.TotalCount)
	assert.Equal(t, 3, len(resp3.Rooms))
	assert.Equal(t, 3, resp3.Page)

	// 测试第4页（只有1个）
	req = httptest.NewRequest("GET", "/api/rooms?page=4&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp4 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp4)
	assert.Equal(t, 10, resp4.TotalCount)
	assert.Equal(t, 1, len(resp4.Rooms))
	assert.Equal(t, 4, resp4.Page)
}

func TestGetRoomsAPI_PageOutOfRange(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建5个房间
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("outofrange_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 请求超出范围的页码
	req := httptest.NewRequest("GET", "/api/rooms?page=10&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 5, response.TotalCount)
	assert.Equal(t, 0, len(response.Rooms)) // 空数组
	assert.Equal(t, 10, response.Page)
	assert.Equal(t, 10, response.Limit)
}

func TestGetRoomsAPI_InvalidPaginationParams(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建若干房间
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("invalid_param_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	testCases := []struct {
		name          string
		url           string
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "Invalid page and limit (page=0, limit=-5)",
			url:           "/api/rooms?page=0&limit=-5",
			expectedPage:  1,
			expectedLimit: 12,
		},
		{
			name:          "Non-numeric parameters",
			url:           "/api/rooms?page=abc&limit=xyz",
			expectedPage:  1,
			expectedLimit: 12,
		},
		{
			name:          "Limit exceeds maximum (100 > 50)",
			url:           "/api/rooms?limit=100",
			expectedPage:  1,
			expectedLimit: 12, // 超过50使用默认值12
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			req.Header.Set("Authorization", "Bearer "+tokens[0])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response RoomListResponse
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedPage, response.Page)
			assert.Equal(t, tc.expectedLimit, response.Limit)
		})
	}
}

// ============ 4. 状态过滤测试 ============

func TestGetRoomsAPI_FilterByStatusWaiting(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 8)
	for i := 0; i < 8; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("filter_wait_user%d", i))
		tokens[i] = token
	}

	// 创建2个waiting房间
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 创建1个ready房间（4人）
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[2])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	readyRoomID := resp.Room.ID

	for i := 3; i < 6; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", readyRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 创建1个playing房间
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[6])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var playingResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &playingResp)

	// (由于只有1人，暂时保持waiting状态)

	// 过滤waiting状态
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// 应该有3个waiting房间（2个初始 + 1个playing创建但只有1人）
	assert.Equal(t, 3, response.TotalCount)
	for _, r := range response.Rooms {
		assert.Equal(t, room.RoomStatusWaiting, r.Status)
	}
}

func TestGetRoomsAPI_FilterByStatusReady(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("filter_ready_user%d", i))
		tokens[i] = token
	}

	// 创建2个ready房间（每个4人）
	for roomIdx := 0; roomIdx < 2; roomIdx++ {
		startIdx := roomIdx * 4
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[startIdx])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomID := resp.Room.ID

		// 加入3个玩家
		for i := 1; i < 4; i++ {
			req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[startIdx+i])
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}

	// 创建1个waiting房间
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[8])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 过滤ready状态
	req = httptest.NewRequest("GET", "/api/rooms?status=ready", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2, response.TotalCount)
	for _, r := range response.Rooms {
		assert.Equal(t, room.RoomStatusReady, r.Status)
		assert.Equal(t, 4, r.PlayerCount)
	}
}

func TestGetRoomsAPI_FilterByStatusPlaying(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 8)
	for i := 0; i < 8; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("filter_play_user%d", i))
		tokens[i] = token
	}

	// 创建1个playing房间
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	playingRoomID := resp.Room.ID

	// 加入3个玩家
	for i := 1; i < 4; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", playingRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 开始游戏
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", playingRoomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 创建1个waiting房间
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 过滤playing状态
	req = httptest.NewRequest("GET", "/api/rooms?status=playing", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, room.RoomStatusPlaying, response.Rooms[0].Status)
	assert.False(t, response.Rooms[0].CanJoin)
}

func TestGetRoomsAPI_InvalidStatusFilter(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建若干房间
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("invalid_status_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	testCases := []string{
		"invalid",
		"closed",
		"unknown",
	}

	for _, status := range testCases {
		t.Run(fmt.Sprintf("status=%s", status), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/rooms?status=%s", status), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[0])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response RoomListResponse
			json.Unmarshal(w.Body.Bytes(), &response)

			// 无效状态应该被忽略，返回所有房间
			assert.Equal(t, 3, response.TotalCount)
		})
	}
}

func TestGetRoomsAPI_CombinedPaginationAndFilter(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 21)
	for i := 0; i < 21; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("combined_user%d", i))
		tokens[i] = token
	}

	// 创建8个waiting房间
	for i := 0; i < 8; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 创建3个ready房间（每个4人，需要12个用户，已用8个，还需12个）
	for roomIdx := 0; roomIdx < 3; roomIdx++ {
		startIdx := 8 + roomIdx*4
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[startIdx])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomID := resp.Room.ID

		for i := 1; i < 4; i++ {
			req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[startIdx+i])
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}

	// 测试: waiting状态，第1页，limit=3
	req := httptest.NewRequest("GET", "/api/rooms?status=waiting&page=1&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp1 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp1)

	assert.Equal(t, 8, resp1.TotalCount)
	assert.Equal(t, 3, len(resp1.Rooms))
	assert.Equal(t, 1, resp1.Page)
	for _, r := range resp1.Rooms {
		assert.Equal(t, room.RoomStatusWaiting, r.Status)
	}

	// 测试: waiting状态，第2页，limit=3
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting&page=2&limit=3", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp2 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp2)

	assert.Equal(t, 8, resp2.TotalCount)
	assert.Equal(t, 3, len(resp2.Rooms))
	assert.Equal(t, 2, resp2.Page)

	// 测试: ready状态，第1页，limit=2
	req = httptest.NewRequest("GET", "/api/rooms?status=ready&page=1&limit=2", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp3 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp3)

	assert.Equal(t, 3, resp3.TotalCount)
	assert.Equal(t, 2, len(resp3.Rooms))
	for _, r := range resp3.Rooms {
		assert.Equal(t, room.RoomStatusReady, r.Status)
	}
}

// ============ 5. 数据完整性测试 ============

func TestGetRoomsAPI_RoomInfoIntegrity(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("integrity_user%d", i))
		tokens[i] = token
	}

	// 创建1个waiting房间（1人）
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 创建1个ready房间（4人）
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	readyRoomID := resp.Room.ID

	for i := 2; i < 5; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", readyRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 获取房间列表
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2, len(response.Rooms))

	// 验证每个房间信息的完整性
	for _, roomInfo := range response.Rooms {
		// 必需字段检查
		assert.NotEmpty(t, roomInfo.ID)
		assert.Contains(t, roomInfo.ID, "room_")
		assert.NotEmpty(t, roomInfo.Owner)
		assert.False(t, roomInfo.CreatedAt.IsZero())

		// 状态值有效性
		assert.True(t, roomInfo.Status == room.RoomStatusWaiting ||
			roomInfo.Status == room.RoomStatusReady ||
			roomInfo.Status == room.RoomStatusPlaying)

		// 玩家数与数组一致性
		assert.NotNil(t, roomInfo.Players)
		assert.Equal(t, roomInfo.PlayerCount, len(roomInfo.Players))

		// can_join是布尔值
		_ = roomInfo.CanJoin // 确保字段存在
	}
}

func TestGetRoomsAPI_CanJoinAccuracy(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("canjoin_user%d", i))
		tokens[i] = token
	}

	// 创建1个waiting房间（1人）
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 创建1个ready房间（4人）
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var readyResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &readyResp)
	readyRoomID := readyResp.Room.ID

	for i := 2; i < 5; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", readyRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 创建1个playing房间
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[5])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var playingResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &playingResp)
	playingRoomID := playingResp.Room.ID

	for i := 6; i < 9; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", playingRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 开始游戏
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", playingRoomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[5])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 获取房间列表
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[9])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// 验证can_join字段
	for _, roomInfo := range response.Rooms {
		if roomInfo.Status == room.RoomStatusWaiting && roomInfo.PlayerCount < 4 {
			assert.True(t, roomInfo.CanJoin, "Waiting room with < 4 players should be joinable")
		} else {
			assert.False(t, roomInfo.CanJoin, "Ready or Playing room should not be joinable")
		}
	}
}

func TestGetRoomsAPI_PlayerInfoAccuracy(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建3个用户
	users := make([]*auth.User, 3)
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, user := createTestUserAndLogin(t, router, fmt.Sprintf("playerinfo_user%d", i))
		tokens[i] = token
		users[i] = user
	}

	// 用户A创建房间
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	roomID := createResp.Room.ID

	// 用户B和C加入
	for i := 1; i < 3; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 获取房间列表
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, len(response.Rooms))
	roomInfo := response.Rooms[0]

	// 验证玩家信息
	assert.Equal(t, 3, roomInfo.PlayerCount)
	assert.Equal(t, 3, len(roomInfo.Players))

	// 检查每个玩家的信息
	for _, player := range roomInfo.Players {
		assert.NotEmpty(t, player.ID)
		assert.NotEmpty(t, player.Username)
		assert.True(t, player.Seat >= 0 && player.Seat <= 3)
		// online字段存在（无需验证值，因为测试环境中都是online）
	}

	// 验证房主在玩家列表中
	ownerFound := false
	for _, player := range roomInfo.Players {
		if player.ID == roomInfo.Owner {
			ownerFound = true
			break
		}
	}
	assert.True(t, ownerFound, "Owner should be in players list")
}

// ============ 6. 边界条件测试 ============

func TestGetRoomsAPI_LargeNumberOfRooms(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建100个房间
	roomCount := 100
	tokens := make([]string, roomCount)
	for i := 0; i < roomCount; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("large_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 请求limit=50
	start := time.Now()
	req := httptest.NewRequest("GET", "/api/rooms?limit=50", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	elapsed := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, roomCount, response.TotalCount)
	assert.Equal(t, 50, len(response.Rooms))
	assert.Equal(t, 50, response.Limit)

	// 验证响应时间（应该很快，< 1秒）
	assert.Less(t, elapsed, 1*time.Second, "Response time should be reasonable")
}

func TestGetRoomsAPI_ConcurrentAccess(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建10个房间
	for i := 0; i < 10; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("concurrent_room_user%d", i))

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 创建20个用户用于并发读取
	tokens := make([]string, 20)
	for i := 0; i < 20; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("concurrent_read_user%d", i))
		tokens[i] = token
	}

	// 并发获取房间列表
	var wg sync.WaitGroup
	results := make([]RoomListResponse, 20)
	statusCodes := make([]int, 20)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			req := httptest.NewRequest("GET", "/api/rooms", nil)
			req.Header.Set("Authorization", "Bearer "+tokens[idx])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			statusCodes[idx] = w.Code
			if w.Code == http.StatusOK {
				json.Unmarshal(w.Body.Bytes(), &results[idx])
			}
		}(i)
	}

	wg.Wait()

	// 验证所有请求都成功
	for i := 0; i < 20; i++ {
		assert.Equal(t, http.StatusOK, statusCodes[i], "Request %d failed", i)
	}

	// 验证所有响应的total_count一致
	expectedCount := results[0].TotalCount
	for i := 1; i < 20; i++ {
		assert.Equal(t, expectedCount, results[i].TotalCount, "Total count mismatch")
	}
}

func TestGetRoomsAPI_DynamicRoomUpdates(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	tokens := make([]string, 8)
	for i := 0; i < 8; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("dynamic_user%d", i))
		tokens[i] = token
	}

	// 创建3个waiting房间
	roomIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomIDs[i] = resp.Room.ID
	}

	// 第1次获取：应该有3个waiting
	req := httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp1 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp1)
	assert.Equal(t, 3, resp1.TotalCount)

	// 玩家加入第一个房间使其变为ready
	for i := 3; i < 6; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomIDs[0]), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 第2次获取：应该有2个waiting, 1个ready
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp2 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp2)
	assert.Equal(t, 2, resp2.TotalCount)

	req = httptest.NewRequest("GET", "/api/rooms?status=ready", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp3 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp3)
	assert.Equal(t, 1, resp3.TotalCount)

	// 第一个房间开始游戏
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomIDs[0]), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 第3次获取：应该有2个waiting, 0个ready, 1个playing
	req = httptest.NewRequest("GET", "/api/rooms?status=playing", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[6])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp4 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp4)
	assert.Equal(t, 1, resp4.TotalCount)
}

func TestGetRoomsAPI_RoomClosureHandling(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建3个房间
	tokens := make([]string, 3)
	roomIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("closure_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomIDs[i] = resp.Room.ID
	}

	// 第1次获取：应该有3个房间
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp1 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp1)
	assert.Equal(t, 3, resp1.TotalCount)

	// 用户离开第一个房间（导致房间关闭，因为只有房主）
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomIDs[0]), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 第2次获取：应该有2个房间
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp2 RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &resp2)
	assert.Equal(t, 2, resp2.TotalCount)

	// 验证关闭的房间不在列表中
	for _, roomInfo := range resp2.Rooms {
		assert.NotEqual(t, roomIDs[0], roomInfo.ID, "Closed room should not appear in list")
	}
}

// ============ 7. 边缘情况测试 ============

func TestGetRoomsAPI_MinimumLimit(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建5个房间
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("minlimit_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 请求limit=1
	req := httptest.NewRequest("GET", "/api/rooms?limit=1", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 5, response.TotalCount)
	assert.Equal(t, 1, len(response.Rooms))
	assert.Equal(t, 1, response.Limit)
}

func TestGetRoomsAPI_MaximumLimit(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建60个房间
	roomCount := 60
	tokens := make([]string, roomCount)
	for i := 0; i < roomCount; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("maxlimit_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 请求limit=50
	req := httptest.NewRequest("GET", "/api/rooms?limit=50", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, roomCount, response.TotalCount)
	assert.Equal(t, 50, len(response.Rooms))
	assert.Equal(t, 50, response.Limit)
}

func TestGetRoomsAPI_NegativeOrZeroPage(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建若干房间
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("negpage_user%d", i))
		tokens[i] = token

		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	testCases := []struct {
		name string
		url  string
	}{
		{"Negative page", "/api/rooms?page=-1"},
		{"Zero page", "/api/rooms?page=0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			req.Header.Set("Authorization", "Bearer "+tokens[0])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response RoomListResponse
			json.Unmarshal(w.Body.Bytes(), &response)

			// 应该被当作page=1处理
			assert.Equal(t, 1, response.Page)
			assert.Equal(t, 3, len(response.Rooms))
		})
	}
}

// ============ 8. 集成测试 ============

func TestGetRoomsAPI_CompleteRoomLifecycle(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建4个用户
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("lifecycle_user%d", i))
		tokens[i] = token
	}

	// 用户A创建房间
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	roomID := createResp.Room.ID

	// 验证：waiting状态可见
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var waitingResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &waitingResp)
	assert.Equal(t, 1, waitingResp.TotalCount)
	assert.Equal(t, roomID, waitingResp.Rooms[0].ID)

	// 用户B、C、D加入
	for i := 1; i < 4; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 验证：ready状态可见
	req = httptest.NewRequest("GET", "/api/rooms?status=ready", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var readyResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &readyResp)
	assert.Equal(t, 1, readyResp.TotalCount)
	assert.Equal(t, roomID, readyResp.Rooms[0].ID)

	// 用户A开始游戏
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证：playing状态可见
	req = httptest.NewRequest("GET", "/api/rooms?status=playing", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[4])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var playingResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &playingResp)
	assert.Equal(t, 1, playingResp.TotalCount)
	assert.Equal(t, roomID, playingResp.Rooms[0].ID)
}

func TestGetRoomsAPI_MultiRoomMultiStatus(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建足够的用户
	totalUsers := 5 + 3*4 + 2*4 // 5个waiting + 3个ready(各4人) + 2个playing(各4人)
	tokens := make([]string, totalUsers)
	for i := 0; i < totalUsers; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("multi_user%d", i))
		tokens[i] = token
	}

	currentTokenIdx := 0

	// 创建5个waiting房间（不同玩家数）
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[currentTokenIdx])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		currentTokenIdx++
	}

	// 创建3个ready房间
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[currentTokenIdx])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomID := resp.Room.ID
		currentTokenIdx++

		for j := 0; j < 3; j++ {
			req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[currentTokenIdx])
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			currentTokenIdx++
		}
	}

	// 创建2个playing房间
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[currentTokenIdx])
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp RoomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		roomID := resp.Room.ID
		ownerToken := tokens[currentTokenIdx]
		currentTokenIdx++

		for j := 0; j < 3; j++ {
			req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[currentTokenIdx])
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			currentTokenIdx++
		}

		// 开始游戏
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// 测试无过滤
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var allResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &allResp)
	assert.Equal(t, 10, allResp.TotalCount)

	// 测试waiting过滤
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var waitingResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &waitingResp)
	assert.Equal(t, 5, waitingResp.TotalCount)

	// 测试ready过滤
	req = httptest.NewRequest("GET", "/api/rooms?status=ready", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var readyResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &readyResp)
	assert.Equal(t, 3, readyResp.TotalCount)

	// 测试playing过滤
	req = httptest.NewRequest("GET", "/api/rooms?status=playing", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var playingResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &playingResp)
	assert.Equal(t, 2, playingResp.TotalCount)
}

func TestGetRoomsAPI_IntegrationWithCreateRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// 创建用户
	token, _ := createTestUserAndLogin(t, router, "integration_user")

	// 获取初始列表
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var initialResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &initialResp)
	initialCount := initialResp.TotalCount

	// 创建新房间
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	newRoomID := createResp.Room.ID

	// 立即获取列表
	req = httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var afterResp RoomListResponse
	json.Unmarshal(w.Body.Bytes(), &afterResp)

	// 验证
	assert.Equal(t, initialCount+1, afterResp.TotalCount)

	// 验证新房间在列表中
	foundNewRoom := false
	for _, roomInfo := range afterResp.Rooms {
		if roomInfo.ID == newRoomID {
			foundNewRoom = true
			assert.Equal(t, createResp.Room.Status, roomInfo.Status)
			assert.Equal(t, createResp.Room.PlayerCount, roomInfo.PlayerCount)
			break
		}
	}
	assert.True(t, foundNewRoom, "Newly created room should appear in list")
}
