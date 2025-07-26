package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/room"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test setup helpers
func setupRoomTestRouter() (*gin.Engine, *AuthHandler, *RoomHandler, auth.AuthService, room.RoomService) {
	gin.SetMode(gin.TestMode)
	
	// Create services
	authService := auth.NewAuthService("test-secret", 24*time.Hour)
	roomService := room.NewRoomService(authService)
	
	// Create handlers
	authHandler := NewAuthHandler(authService)
	roomHandler := NewRoomHandler(roomService, authService)
	
	// Setup router
	router := gin.New()
	authHandler.RegisterRoutes(router)
	roomHandler.RegisterRoutes(router, authHandler)
	
	return router, authHandler, roomHandler, authService, roomService
}

func createTestUserAndLogin(t *testing.T, router *gin.Engine, username string) (string, *auth.User) {
	// Register user
	registerReq := RegisterRequest{
		Username: username,
		Password: "password123",
	}
	
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var authResp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	
	return authResp.Token.Token, authResp.User
}

func TestRoomHandler_CreateRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create and login user
	token, user := createTestUserAndLogin(t, router, "testuser")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify room properties
	assert.NotEmpty(t, response.Room.ID)
	assert.Equal(t, room.RoomStatusWaiting, response.Room.Status)
	assert.Equal(t, user.ID, response.Room.Owner)
	assert.Equal(t, 1, response.Room.PlayerCount)
	assert.NotNil(t, response.Room.Players[0])
	assert.Equal(t, user.ID, response.Room.Players[0].ID)
	assert.Equal(t, 0, response.Room.Players[0].Seat)
}

func TestRoomHandler_CreateRoom_Unauthorized(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Try to create room without authentication
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoomHandler_CreateRoom_AlreadyInRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create and login user
	token, _ := createTestUserAndLogin(t, router, "testuser")
	
	// Create first room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Try to create second room
	req2 := httptest.NewRequest("POST", "/api/rooms", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusConflict, w2.Code)
	
	var errorResp ErrorResponse
	err := json.Unmarshal(w2.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "already_in_room", errorResp.Error)
}

func TestRoomHandler_GetRooms(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create multiple users and rooms
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("user%d", i))
		tokens[i] = token
		
		// Create room
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}
	
	// Get room list
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, 3, response.TotalCount)
	assert.Equal(t, 3, len(response.Rooms))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 12, response.Limit)
	
	// Verify all rooms are waiting status
	for _, roomInfo := range response.Rooms {
		assert.Equal(t, room.RoomStatusWaiting, roomInfo.Status)
		assert.Equal(t, 1, roomInfo.PlayerCount)
		assert.True(t, roomInfo.CanJoin)
	}
}

func TestRoomHandler_GetRooms_WithPagination(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create multiple users and rooms
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("paguser%d", i))
		tokens[i] = token
		
		// Create room
		req := httptest.NewRequest("POST", "/api/rooms", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}
	
	// Get first page with limit 2
	req := httptest.NewRequest("GET", "/api/rooms?page=1&limit=2", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, 5, response.TotalCount)
	assert.Equal(t, 2, len(response.Rooms))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 2, response.Limit)
}

func TestRoomHandler_GetRooms_WithStatusFilter(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create users
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("filteruser%d", i))
		tokens[i] = token
	}
	
	// Create a waiting room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Create a ready room (4 players)
	req = httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	readyRoomID := roomResp.Room.ID
	
	// Add 3 more players to make it ready
	for i := 2; i < 5; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", readyRoomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
	
	// Get only waiting rooms
	req = httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, 1, len(response.Rooms))
	assert.Equal(t, room.RoomStatusWaiting, response.Rooms[0].Status)
}

func TestRoomHandler_JoinRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, player := createTestUserAndLogin(t, router, "player")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp.Room.ID
	
	// Join room
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var joinResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &joinResp)
	assert.NoError(t, err)
	
	assert.Equal(t, 2, joinResp.Room.PlayerCount)
	assert.Equal(t, room.RoomStatusWaiting, joinResp.Room.Status)
	
	// Find the joined player
	var joinedPlayer *room.Player
	for _, p := range joinResp.Room.Players {
		if p != nil && p.ID == player.ID {
			joinedPlayer = p
			break
		}
	}
	
	assert.NotNil(t, joinedPlayer)
	assert.Equal(t, player.ID, joinedPlayer.ID)
	assert.True(t, joinedPlayer.Seat >= 0 && joinedPlayer.Seat <= 3)
}

func TestRoomHandler_JoinRoom_NotFound(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create player
	playerToken, _ := createTestUserAndLogin(t, router, "player")
	
	// Try to join non-existent room
	req := httptest.NewRequest("POST", "/api/rooms/nonexistent/join", nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "room_not_found", errorResp.Error)
}

func TestRoomHandler_LeaveRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, _ := createTestUserAndLogin(t, router, "player")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp.Room.ID
	
	// Join room
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Leave room
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var leaveResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &leaveResp)
	assert.NoError(t, err)
	
	assert.Equal(t, 1, leaveResp.Room.PlayerCount)
}

func TestRoomHandler_StartGame(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create users
	tokens := make([]string, 4)
	for i := 0; i < 4; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("gameuser%d", i))
		tokens[i] = token
	}
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp.Room.ID
	
	// Add 3 more players
	for i := 1; i < 4; i++ {
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens[i])
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
	
	// Start game
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var startResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &startResp)
	assert.NoError(t, err)
	
	assert.Equal(t, room.RoomStatusPlaying, startResp.Room.Status)
	assert.Equal(t, 4, startResp.Room.PlayerCount)
}

func TestRoomHandler_StartGame_NotOwner(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, _ := createTestUserAndLogin(t, router, "player")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp.Room.ID
	
	// Try to start game as non-owner
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "not_room_owner", errorResp.Error)
}

func TestRoomHandler_StartGame_InsufficientPlayers(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create owner
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp.Room.ID
	
	// Try to start game with only 1 player
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusConflict, w.Code)
	
	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "room_not_ready", errorResp.Error)
}

func TestRoomHandler_GetMyRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create user
	token, _ := createTestUserAndLogin(t, router, "testuser")
	
	// Create room
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var roomResp RoomResponse
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	expectedRoomID := roomResp.Room.ID
	
	// Get my room
	req = httptest.NewRequest("GET", "/api/rooms/my", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var myRoomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &myRoomResp)
	assert.NoError(t, err)
	
	assert.Equal(t, expectedRoomID, myRoomResp.Room.ID)
}

func TestRoomHandler_GetMyRoom_NotInRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()
	
	// Create user but don't create room
	token, _ := createTestUserAndLogin(t, router, "testuser")
	
	// Try to get my room
	req := httptest.NewRequest("GET", "/api/rooms/my", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "not_in_room", errorResp.Error)
}