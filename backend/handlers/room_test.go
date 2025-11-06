package handlers

import (
	"bytes"
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

// createRoom creates a room and returns the room ID
func createRoom(t *testing.T, router *gin.Engine, token string) string {
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var roomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &roomResp)
	assert.NoError(t, err)

	return roomResp.Room.ID
}

// joinRoom makes a player join a room and returns the response
func joinRoom(t *testing.T, router *gin.Engine, roomID, token string) *RoomResponse {
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var roomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &roomResp)
	assert.NoError(t, err)

	return &roomResp
}

// leaveRoom makes a player leave a room and returns the response
func leaveRoom(t *testing.T, router *gin.Engine, roomID, token string) *RoomResponse {
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var roomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &roomResp)
	assert.NoError(t, err)

	return &roomResp
}

// startGame starts a game in a room
func startGame(t *testing.T, router *gin.Engine, roomID, token string) *RoomResponse {
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var roomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &roomResp)
	assert.NoError(t, err)

	return &roomResp
}

// getRoom retrieves a room by ID
func getRoom(t *testing.T, router *gin.Engine, roomID, token string) *RoomResponse {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/rooms/%s", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var roomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &roomResp)
	assert.NoError(t, err)

	return &roomResp
}

// createFullRoom creates a room with 4 players (ready status)
func createFullRoom(t *testing.T, router *gin.Engine, userPrefix string) (string, []string) {
	tokens := make([]string, 4)
	for i := 0; i < 4; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("%s%d", userPrefix, i))
		tokens[i] = token
	}

	// Create room with first user
	roomID := createRoom(t, router, tokens[0])

	// Add 3 more players
	for i := 1; i < 4; i++ {
		joinRoom(t, router, roomID, tokens[i])
	}

	return roomID, tokens
}

func TestRoomHandler_CreateRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create and login user
	token, user := createTestUserAndLogin(t, router, "testuser")

	// Create room
	roomID := createRoom(t, router, token)

	// Get room to verify properties
	response := getRoom(t, router, roomID, token)

	// Verify room properties
	assert.NotEmpty(t, response.Room.ID)
	assert.Equal(t, room.RoomStatusWaiting, response.Room.Status)
	assert.Equal(t, user.ID, response.Room.Owner)
	assert.Equal(t, 1, response.Room.PlayerCount)
	assert.NotNil(t, response.Room.Players[0])
	assert.Equal(t, user.ID, response.Room.Players[0].ID)
	assert.Equal(t, 0, response.Room.Players[0].Seat)
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
		createRoom(t, router, token)
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
		createRoom(t, router, token)
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
	createRoom(t, router, tokens[0])

	// Create a ready room (4 players)
	readyRoomID := createRoom(t, router, tokens[1])

	// Add 3 more players to make it ready
	for i := 2; i < 5; i++ {
		joinRoom(t, router, readyRoomID, tokens[i])
	}

	// Get only waiting rooms
	req := httptest.NewRequest("GET", "/api/rooms?status=waiting", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, 1, len(response.Rooms))
	assert.Equal(t, room.RoomStatusWaiting, response.Rooms[0].Status)
}

func TestRoomHandler_GetRooms_InvalidPagination(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user
	token, _ := createTestUserAndLogin(t, router, "paginuser")

	// Test with invalid page (0)
	req := httptest.NewRequest("GET", "/api/rooms?page=0&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Page 0 should be treated as page 1
	assert.Equal(t, 1, response.Page)

	// Test with negative limit
	req2 := httptest.NewRequest("GET", "/api/rooms?page=1&limit=-5", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response2 RoomListResponse
	err2 := json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err2)

	// Negative limit should use default
	assert.True(t, response2.Limit > 0)
}

func TestRoomHandler_GetRooms_EmptyList(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user but no rooms
	token, _ := createTestUserAndLogin(t, router, "emptyuser")

	// Get room list
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
}

func TestRoomHandler_GetRooms_ReadyStatusFilter(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create waiting room with separate user
	token1, _ := createTestUserAndLogin(t, router, "readyfilter_waiting")
	createRoom(t, router, token1)

	// Create ready room (4 players)
	roomID, tokens := createFullRoom(t, router, "readyfilter")

	// Get only ready rooms
	req := httptest.NewRequest("GET", "/api/rooms?status=ready", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, 1, len(response.Rooms))
	assert.Equal(t, room.RoomStatusReady, response.Rooms[0].Status)
	assert.Equal(t, roomID, response.Rooms[0].ID)
	assert.Equal(t, 4, response.Rooms[0].PlayerCount)
	assert.False(t, response.Rooms[0].CanJoin)
}

func TestRoomHandler_GetRooms_PlayingStatusFilter(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create waiting room with separate user
	token1, _ := createTestUserAndLogin(t, router, "playfilter_waiting")
	createRoom(t, router, token1)

	// Create playing room (4 players + start game)
	roomID, tokens := createFullRoom(t, router, "playfilter")
	startGame(t, router, roomID, tokens[0])

	// Get only playing rooms
	req := httptest.NewRequest("GET", "/api/rooms?status=playing", nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RoomListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.TotalCount)
	assert.Equal(t, 1, len(response.Rooms))
	assert.Equal(t, room.RoomStatusPlaying, response.Rooms[0].Status)
	assert.Equal(t, roomID, response.Rooms[0].ID)
	assert.Equal(t, 4, response.Rooms[0].PlayerCount)
	assert.False(t, response.Rooms[0].CanJoin)
}

func TestRoomHandler_JoinRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, player := createTestUserAndLogin(t, router, "player")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Join room
	joinResp := joinRoom(t, router, roomID, playerToken)

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

func TestRoomHandler_JoinRoom_AlreadyInOtherRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create two room owners and one player
	owner1Token, _ := createTestUserAndLogin(t, router, "owner1")
	owner2Token, _ := createTestUserAndLogin(t, router, "owner2")
	playerToken, _ := createTestUserAndLogin(t, router, "player")

	// Create two rooms
	room1ID := createRoom(t, router, owner1Token)
	room2ID := createRoom(t, router, owner2Token)

	// Player joins room1
	joinRoom(t, router, room1ID, playerToken)

	// Try to join room2 while already in room1
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", room2ID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Contains(t, errorResp.Message, "already in room")
}

func TestRoomHandler_JoinRoom_AlreadyInSameRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, player := createTestUserAndLogin(t, router, "player")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Player joins room
	joinRoom(t, router, roomID, playerToken)

	// Try to join same room again (should be idempotent)
	joinResp := joinRoom(t, router, roomID, playerToken)

	// Should succeed and return same room state
	assert.Equal(t, 2, joinResp.Room.PlayerCount)

	// Verify player still in room
	var foundPlayer *room.Player
	for _, p := range joinResp.Room.Players {
		if p != nil && p.ID == player.ID {
			foundPlayer = p
			break
		}
	}
	assert.NotNil(t, foundPlayer)
}

func TestRoomHandler_JoinRoom_RoomFull(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room with 4 players (becomes Ready status)
	roomID, tokens := createFullRoom(t, router, "fulluser")

	// Verify room is ready (full)
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusReady, roomResp.Room.Status)
	assert.Equal(t, 4, roomResp.Room.PlayerCount)

	// Create a 5th player
	fifthPlayerToken, _ := createTestUserAndLogin(t, router, "fifthplayer")

	// Try to join full room (Ready status doesn't accept new players)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+fifthPlayerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	// Room is full (Ready status), so it's not accepting new players
	assert.Equal(t, "room_not_accepting", errorResp.Error)
}

func TestRoomHandler_JoinRoom_PlayingStatus(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room and start game
	roomID, tokens := createFullRoom(t, router, "playuser")
	startGame(t, router, roomID, tokens[0])

	// Create a new player
	newPlayerToken, _ := createTestUserAndLogin(t, router, "newplayer")

	// Try to join playing room
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/join", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+newPlayerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "room_not_accepting", errorResp.Error)
}

func TestRoomHandler_JoinRoom_StatusChangeToReady(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create users
	tokens := make([]string, 4)
	for i := 0; i < 4; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("readyuser%d", i))
		tokens[i] = token
	}

	// Create room (waiting status)
	roomID := createRoom(t, router, tokens[0])

	// Verify initial waiting status
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusWaiting, roomResp.Room.Status)
	assert.Equal(t, 1, roomResp.Room.PlayerCount)

	// Add second player - still waiting
	joinResp := joinRoom(t, router, roomID, tokens[1])
	assert.Equal(t, room.RoomStatusWaiting, joinResp.Room.Status)
	assert.Equal(t, 2, joinResp.Room.PlayerCount)

	// Add third player - still waiting
	joinResp = joinRoom(t, router, roomID, tokens[2])
	assert.Equal(t, room.RoomStatusWaiting, joinResp.Room.Status)
	assert.Equal(t, 3, joinResp.Room.PlayerCount)

	// Add fourth player - should become ready
	joinResp = joinRoom(t, router, roomID, tokens[3])
	assert.Equal(t, room.RoomStatusReady, joinResp.Room.Status)
	assert.Equal(t, 4, joinResp.Room.PlayerCount)
}

func TestRoomHandler_JoinRoom_VerifySeatAssignment(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create users
	tokens := make([]string, 4)
	users := make([]*auth.User, 4)
	for i := 0; i < 4; i++ {
		token, user := createTestUserAndLogin(t, router, fmt.Sprintf("seatuser%d", i))
		tokens[i] = token
		users[i] = user
	}

	// Create room
	roomID := createRoom(t, router, tokens[0])

	// Track assigned seats
	assignedSeats := make(map[int]bool)

	// First player should be in seat 0
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, 0, roomResp.Room.Players[0].Seat)
	assignedSeats[0] = true

	// Add remaining players and verify seat assignments
	for i := 1; i < 4; i++ {
		joinResp := joinRoom(t, router, roomID, tokens[i])

		// Find the newly joined player
		var joinedPlayer *room.Player
		for _, p := range joinResp.Room.Players {
			if p != nil && p.ID == users[i].ID {
				joinedPlayer = p
				break
			}
		}

		assert.NotNil(t, joinedPlayer)
		assert.True(t, joinedPlayer.Seat >= 0 && joinedPlayer.Seat <= 3, "Seat should be between 0-3")
		assert.False(t, assignedSeats[joinedPlayer.Seat], "Seat %d should not be already assigned", joinedPlayer.Seat)
		assignedSeats[joinedPlayer.Seat] = true
	}

	// Verify all 4 seats are assigned
	assert.Equal(t, 4, len(assignedSeats))
}

// ============ LeaveRoom Test Cases ============

func TestRoomHandler_LeaveRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, _ := createTestUserAndLogin(t, router, "player")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Join room
	joinRoom(t, router, roomID, playerToken)

	// Leave room
	leaveResp := leaveRoom(t, router, roomID, playerToken)

	assert.Equal(t, 1, leaveResp.Room.PlayerCount)
}

func TestRoomHandler_LeaveRoom_NotInRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, _ := createTestUserAndLogin(t, router, "player")

	// Create room (player doesn't join)
	roomID := createRoom(t, router, ownerToken)

	// Try to leave room player hasn't joined
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "not_in_room", errorResp.Error)
}

func TestRoomHandler_LeaveRoom_OwnerTransfer(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and players
	ownerToken, owner := createTestUserAndLogin(t, router, "owner")
	player2Token, player2 := createTestUserAndLogin(t, router, "player2")
	player3Token, _ := createTestUserAndLogin(t, router, "player3")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Add two more players
	joinRoom(t, router, roomID, player2Token)
	joinRoom(t, router, roomID, player3Token)

	// Verify initial owner
	roomResp := getRoom(t, router, roomID, ownerToken)
	assert.Equal(t, owner.ID, roomResp.Room.Owner)
	assert.Equal(t, 3, roomResp.Room.PlayerCount)

	// Owner leaves room
	leaveResp := leaveRoom(t, router, roomID, ownerToken)

	// Verify ownership transferred to first remaining player
	assert.Equal(t, player2.ID, leaveResp.Room.Owner)
	assert.Equal(t, 2, leaveResp.Room.PlayerCount)

	// Verify original owner not in room
	var foundOwner *room.Player
	for _, p := range leaveResp.Room.Players {
		if p != nil && p.ID == owner.ID {
			foundOwner = p
			break
		}
	}
	assert.Nil(t, foundOwner)
}

func TestRoomHandler_LeaveRoom_LastPlayerClosesRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner (only player)
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Last player leaves room
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response should indicate room was closed
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have a message field (not Room field)
	_, hasMessage := response["message"]
	assert.True(t, hasMessage)

	// Verify room no longer exists
	req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/rooms/%s", roomID), nil)
	req2.Header.Set("Authorization", "Bearer "+ownerToken)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestRoomHandler_LeaveRoom_ReadyToWaiting(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room (ready status)
	roomID, tokens := createFullRoom(t, router, "readyuser")

	// Verify ready status
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusReady, roomResp.Room.Status)
	assert.Equal(t, 4, roomResp.Room.PlayerCount)

	// One player leaves
	leaveResp := leaveRoom(t, router, roomID, tokens[3])

	// Should change back to waiting status
	assert.Equal(t, room.RoomStatusWaiting, leaveResp.Room.Status)
	assert.Equal(t, 3, leaveResp.Room.PlayerCount)
}

func TestRoomHandler_LeaveRoom_PlayingRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room and start game
	roomID, tokens := createFullRoom(t, router, "playuser")
	startGame(t, router, roomID, tokens[0])

	// Verify playing status
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusPlaying, roomResp.Room.Status)

	// Player leaves playing room
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/leave", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[1])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response should indicate room was closed (game ended due to insufficient players)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have a message field indicating game ended
	_, hasMessage := response["message"]
	assert.True(t, hasMessage)

	// Verify room no longer exists (game ended, room closed)
	req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/rooms/%s", roomID), nil)
	req2.Header.Set("Authorization", "Bearer "+tokens[0])

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestRoomHandler_LeaveRoom_VerifySeatCleared(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and players
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	player2Token, player2 := createTestUserAndLogin(t, router, "player2")
	player3Token, _ := createTestUserAndLogin(t, router, "player3")

	// Create room and add players
	roomID := createRoom(t, router, ownerToken)
	joinRoom(t, router, roomID, player2Token)
	joinResp := joinRoom(t, router, roomID, player3Token)

	// Get player2's seat
	var player2Seat int
	for _, p := range joinResp.Room.Players {
		if p != nil && p.ID == player2.ID {
			player2Seat = p.Seat
			break
		}
	}

	// Player2 leaves
	leaveResp := leaveRoom(t, router, roomID, player2Token)

	// Verify seat is now nil
	assert.Nil(t, leaveResp.Room.Players[player2Seat])
	assert.Equal(t, 2, leaveResp.Room.PlayerCount)
}

func TestRoomHandler_StartGame(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room with 4 players
	roomID, tokens := createFullRoom(t, router, "gameuser")

	// Start game
	startResp := startGame(t, router, roomID, tokens[0])

	assert.Equal(t, room.RoomStatusPlaying, startResp.Room.Status)
	assert.Equal(t, 4, startResp.Room.PlayerCount)
}

func TestRoomHandler_StartGame_NotOwner(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create owner and player
	ownerToken, _ := createTestUserAndLogin(t, router, "owner")
	playerToken, _ := createTestUserAndLogin(t, router, "player")

	// Create room
	roomID := createRoom(t, router, ownerToken)

	// Try to start game as non-owner
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)

	w := httptest.NewRecorder()
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
	roomID := createRoom(t, router, ownerToken)

	// Try to start game with only 1 player
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "room_not_ready", errorResp.Error)
}

func TestRoomHandler_StartGame_NotMember(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room
	roomID, _ := createFullRoom(t, router, "startmember")

	// Create a different user not in the room
	outsiderToken, _ := createTestUserAndLogin(t, router, "outsider")

	// Try to start game as non-member
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+outsiderToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "not_room_owner", errorResp.Error)
}

func TestRoomHandler_StartGame_AlreadyPlaying(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room and start game
	roomID, tokens := createFullRoom(t, router, "startalready")
	startGame(t, router, roomID, tokens[0])

	// Verify playing status
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusPlaying, roomResp.Room.Status)

	// Try to start game again
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%s/start", roomID), nil)
	req.Header.Set("Authorization", "Bearer "+tokens[0])

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "room_not_ready", errorResp.Error)
}

func TestRoomHandler_StartGame_VerifyStatusChange(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room
	roomID, tokens := createFullRoom(t, router, "startverify")

	// Verify initial Ready status
	roomResp := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusReady, roomResp.Room.Status)

	// Start game
	startResp := startGame(t, router, roomID, tokens[0])

	// Verify status changed to Playing
	assert.Equal(t, room.RoomStatusPlaying, startResp.Room.Status)
	assert.Equal(t, 4, startResp.Room.PlayerCount)

	// Verify by getting room again
	roomResp2 := getRoom(t, router, roomID, tokens[0])
	assert.Equal(t, room.RoomStatusPlaying, roomResp2.Room.Status)
}

func TestRoomHandler_GetMyRoom(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user
	token, _ := createTestUserAndLogin(t, router, "testuser")

	// Create room
	expectedRoomID := createRoom(t, router, token)

	// Get my room
	req := httptest.NewRequest("GET", "/api/rooms/my", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
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

// ============ Additional CreateRoom Test Cases ============

func TestRoomHandler_CreateRoom_JoinThenCreate(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create two users
	tokenA, _ := createTestUserAndLogin(t, router, "userA")
	tokenB, _ := createTestUserAndLogin(t, router, "userB")

	// User A creates room
	roomID := createRoom(t, router, tokenA)

	// User B joins room A
	joinRoom(t, router, roomID, tokenB)

	// User B tries to create a new room while in room A
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should fail with conflict
	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "already_in_room", errorResp.Error)
}

func TestRoomHandler_CreateRoom_LeaveAndCreate(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create two users
	tokenA, _ := createTestUserAndLogin(t, router, "userAA")
	tokenB, _ := createTestUserAndLogin(t, router, "userBB")

	// User A creates room
	roomID := createRoom(t, router, tokenA)

	// User B joins room
	joinRoom(t, router, roomID, tokenB)

	// User B leaves room
	leaveRoom(t, router, roomID, tokenB)

	// User B creates a new room after leaving
	req := httptest.NewRequest("POST", "/api/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed
	assert.Equal(t, http.StatusCreated, w.Code)

	var newRoomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &newRoomResp)
	assert.NoError(t, err)

	// Verify it's a different room
	assert.NotEqual(t, roomID, newRoomResp.Room.ID)
	assert.Equal(t, room.RoomStatusWaiting, newRoomResp.Room.Status)
	assert.Equal(t, 1, newRoomResp.Room.PlayerCount)
}

func TestRoomHandler_CreateRoom_VerifyInitialState(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user
	token, user := createTestUserAndLogin(t, router, "stateuser")

	// Create room
	roomID := createRoom(t, router, token)

	// Get room to verify state
	response := getRoom(t, router, roomID, token)

	// Verify all initial state properties
	r := response.Room

	// Room properties
	assert.NotEmpty(t, r.ID)
	assert.Contains(t, r.ID, "room_")
	assert.Equal(t, room.RoomStatusWaiting, r.Status)
	assert.Equal(t, user.ID, r.Owner)
	assert.Equal(t, 1, r.PlayerCount)

	// Timestamps
	assert.False(t, r.CreatedAt.IsZero())
	assert.False(t, r.UpdatedAt.IsZero())
	timeDiff := r.UpdatedAt.Sub(r.CreatedAt)
	assert.True(t, timeDiff >= 0 && timeDiff < time.Second, "created_at and updated_at should be very close")

	// Players array
	assert.NotNil(t, r.Players[0], "First player should exist")
	assert.Nil(t, r.Players[1], "Second seat should be empty")
	assert.Nil(t, r.Players[2], "Third seat should be empty")
	assert.Nil(t, r.Players[3], "Fourth seat should be empty")

	// First player properties
	p := r.Players[0]
	assert.Equal(t, user.ID, p.ID)
	assert.Equal(t, user.Username, p.Username)
	assert.Equal(t, 0, p.Seat)
	assert.True(t, p.Online)
}

func TestRoomHandler_CreateRoom_UniqueRoomIDs(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create 10 users and rooms
	roomIDs := make(map[string]bool)

	for i := 0; i < 10; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("uniqueuser%d", i))
		roomID := createRoom(t, router, token)

		// Check for duplicate room ID
		assert.False(t, roomIDs[roomID], "Room ID %s is duplicate", roomID)
		roomIDs[roomID] = true
	}

	// Verify we have 10 unique room IDs
	assert.Equal(t, 10, len(roomIDs))
}

func TestRoomHandler_CreateRoom_MultipleUsersSimultaneous(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create 5 users
	userCount := 5
	tokens := make([]string, userCount)
	for i := 0; i < userCount; i++ {
		token, _ := createTestUserAndLogin(t, router, fmt.Sprintf("simuser%d", i))
		tokens[i] = token
	}

	// Create rooms simultaneously
	var wg sync.WaitGroup
	results := make([]int, userCount)
	roomIDs := make([]string, userCount)

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			req := httptest.NewRequest("POST", "/api/rooms", nil)
			req.Header.Set("Authorization", "Bearer "+tokens[idx])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			results[idx] = w.Code

			if w.Code == http.StatusCreated {
				var response RoomResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				roomIDs[idx] = response.Room.ID
			}
		}(i)
	}

	wg.Wait()

	// Verify all requests succeeded
	for i := 0; i < userCount; i++ {
		assert.Equal(t, http.StatusCreated, results[i], "User %d failed to create room", i)
	}

	// Verify all room IDs are unique
	uniqueIDs := make(map[string]bool)
	for _, id := range roomIDs {
		if id != "" {
			uniqueIDs[id] = true
		}
	}
	assert.Equal(t, userCount, len(uniqueIDs), "Not all room IDs are unique")
}

func TestRoomHandler_CreateRoom_VerifyPlayerMapping(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user and room
	token, _ := createTestUserAndLogin(t, router, "mapuser")

	// Create room
	createdRoomID := createRoom(t, router, token)

	// Verify player-room mapping by getting "my room"
	req := httptest.NewRequest("GET", "/api/rooms/my", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var myRoomResp RoomResponse
	err := json.Unmarshal(w.Body.Bytes(), &myRoomResp)
	assert.NoError(t, err)

	// Verify it's the same room
	assert.Equal(t, createdRoomID, myRoomResp.Room.ID)
}

// ============ GetRoom Test Cases ============

func TestRoomHandler_GetRoom_Success(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user and room
	token, user := createTestUserAndLogin(t, router, "getuser")

	// Create room
	roomID := createRoom(t, router, token)

	// Get room by ID
	getResp := getRoom(t, router, roomID, token)

	// Verify room data
	assert.Equal(t, roomID, getResp.Room.ID)
	assert.Equal(t, room.RoomStatusWaiting, getResp.Room.Status)
	assert.Equal(t, user.ID, getResp.Room.Owner)
	assert.Equal(t, 1, getResp.Room.PlayerCount)
	assert.NotNil(t, getResp.Room.Players[0])
	assert.Equal(t, user.ID, getResp.Room.Players[0].ID)
}

func TestRoomHandler_GetRoom_WaitingStatus(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user and room
	token, _ := createTestUserAndLogin(t, router, "waitinguser")

	// Create room (will be in waiting status with 1 player)
	roomID := createRoom(t, router, token)

	// Get room by ID
	getResp := getRoom(t, router, roomID, token)

	// Verify waiting status
	assert.Equal(t, room.RoomStatusWaiting, getResp.Room.Status)
	assert.Equal(t, 1, getResp.Room.PlayerCount)
}

func TestRoomHandler_GetRoom_ReadyStatus(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room with 4 players
	roomID, tokens := createFullRoom(t, router, "readyuser")

	// Get room by ID
	getResp := getRoom(t, router, roomID, tokens[0])

	// Verify ready status
	assert.Equal(t, room.RoomStatusReady, getResp.Room.Status)
	assert.Equal(t, 4, getResp.Room.PlayerCount)

	// Verify all 4 players are present
	playerCount := 0
	for _, p := range getResp.Room.Players {
		if p != nil {
			playerCount++
		}
	}
	assert.Equal(t, 4, playerCount)
}

func TestRoomHandler_GetRoom_PlayingStatus(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create full room with 4 players
	roomID, tokens := createFullRoom(t, router, "playuser")

	// Start game
	startGame(t, router, roomID, tokens[0])

	// Get room by ID
	getResp := getRoom(t, router, roomID, tokens[0])

	// Verify playing status
	assert.Equal(t, room.RoomStatusPlaying, getResp.Room.Status)
	assert.Equal(t, 4, getResp.Room.PlayerCount)
}

func TestRoomHandler_GetRoom_VerifyAllFields(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create 2 users
	token1, user1 := createTestUserAndLogin(t, router, "fielduser1")
	token2, user2 := createTestUserAndLogin(t, router, "fielduser2")

	// Create room
	roomID := createRoom(t, router, token1)

	// Get created time
	createResp := getRoom(t, router, roomID, token1)
	createdAt := createResp.Room.CreatedAt

	// Add second player
	joinRoom(t, router, roomID, token2)

	// Get room by ID
	getResp := getRoom(t, router, roomID, token1)

	r := getResp.Room

	// Verify basic fields
	assert.Equal(t, roomID, r.ID)
	assert.Equal(t, room.RoomStatusWaiting, r.Status)
	assert.Equal(t, user1.ID, r.Owner)
	assert.Equal(t, 2, r.PlayerCount)

	// Verify timestamps
	assert.False(t, r.CreatedAt.IsZero())
	assert.False(t, r.UpdatedAt.IsZero())
	assert.Equal(t, createdAt, r.CreatedAt)
	assert.True(t, r.UpdatedAt.After(r.CreatedAt) || r.UpdatedAt.Equal(r.CreatedAt))

	// Verify players array structure
	assert.NotNil(t, r.Players[0], "First player should exist")
	assert.NotNil(t, r.Players[1], "Second player should exist")
	assert.Nil(t, r.Players[2], "Third seat should be empty")
	assert.Nil(t, r.Players[3], "Fourth seat should be empty")

	// Verify first player
	p1 := r.Players[0]
	assert.Equal(t, user1.ID, p1.ID)
	assert.Equal(t, user1.Username, p1.Username)
	assert.Equal(t, 0, p1.Seat)
	assert.True(t, p1.Online)

	// Verify second player
	p2 := r.Players[1]
	assert.Equal(t, user2.ID, p2.ID)
	assert.Equal(t, user2.Username, p2.Username)
	assert.Equal(t, 1, p2.Seat)
	assert.True(t, p2.Online)
}

func TestRoomHandler_GetRoom_DifferentUserCanView(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create two users
	token1, _ := createTestUserAndLogin(t, router, "roomowner")
	token2, _ := createTestUserAndLogin(t, router, "viewer")

	// User 1 creates room
	roomID := createRoom(t, router, token1)

	// User 2 (not in room) tries to view room
	getResp := getRoom(t, router, roomID, token2)

	// Should succeed - any authenticated user can view room details
	assert.Equal(t, roomID, getResp.Room.ID)
}

// ============ Table-Driven Tests for Common Error Scenarios ============

func TestRoomHandler_UnauthorizedAccess(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Prepare test data (create room for endpoints that need it)
	token, _ := createTestUserAndLogin(t, router, "setup_user")
	roomID := createRoom(t, router, token)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "CreateRoom without auth",
			method:         "POST",
			path:           "/api/rooms",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetRooms without auth",
			method:         "GET",
			path:           "/api/rooms",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "JoinRoom without auth",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/join", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "LeaveRoom without auth",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/leave", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "StartGame without auth",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/start", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetMyRoom without auth",
			method:         "GET",
			path:           "/api/rooms/my",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetRoom without auth",
			method:         "GET",
			path:           fmt.Sprintf("/api/rooms/%s", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			// No Authorization header

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRoomHandler_InvalidToken(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Prepare test data
	token, _ := createTestUserAndLogin(t, router, "setup_user")
	roomID := createRoom(t, router, token)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "CreateRoom with invalid token",
			method:         "POST",
			path:           "/api/rooms",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "JoinRoom with invalid token",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/join", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "StartGame with invalid token",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/start", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetMyRoom with invalid token",
			method:         "GET",
			path:           "/api/rooms/my",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetRoom with invalid token",
			method:         "GET",
			path:           fmt.Sprintf("/api/rooms/%s", roomID),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Authorization", "Bearer invalid-token-xyz")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRoomHandler_RoomNotFound(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user for authentication
	token, _ := createTestUserAndLogin(t, router, "test_user")
	nonExistentRoomID := "nonexistent-room-id"

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "JoinRoom not found",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/join", nonExistentRoomID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "room_not_found",
		},
		{
			name:           "LeaveRoom not found",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/leave", nonExistentRoomID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "room_not_found",
		},
		{
			name:           "StartGame not found",
			method:         "POST",
			path:           fmt.Sprintf("/api/rooms/%s/start", nonExistentRoomID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "room_not_found",
		},
		{
			name:           "GetRoom not found",
			method:         "GET",
			path:           fmt.Sprintf("/api/rooms/%s", nonExistentRoomID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "room_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var errorResp ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &errorResp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedError, errorResp.Error)
		})
	}
}

func TestRoomHandler_EmptyRoomID(t *testing.T) {
	router, _, _, _, _ := setupRoomTestRouter()

	// Create user for authentication
	token, _ := createTestUserAndLogin(t, router, "test_user")

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "JoinRoom with empty ID",
			method: "POST",
			path:   "/api/rooms//join",
		},
		{
			name:   "LeaveRoom with empty ID",
			method: "POST",
			path:   "/api/rooms//leave",
		},
		{
			name:   "StartGame with empty ID",
			method: "POST",
			path:   "/api/rooms//start",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 404 (route not found) or 400
			assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusBadRequest,
				"Expected 404 or 400, got %d", w.Code)
		})
	}
}
