package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"guandan-world/backend/auth"
	"guandan-world/backend/room"
)

// MockAuthService is a mock implementation of auth.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, password string) (*auth.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockAuthService) Login(username, password string) (*auth.AuthToken, error) {
	args := m.Called(username, password)
	return args.Get(0).(*auth.AuthToken), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*auth.User, error) {
	args := m.Called(token)
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockAuthService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthService) GetUserByID(userID string) (*auth.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

// MockRoomService is a mock implementation of room.RoomService
type MockRoomService struct {
	mock.Mock
}

func (m *MockRoomService) CreateRoom(ownerID string) (*room.Room, error) {
	args := m.Called(ownerID)
	return args.Get(0).(*room.Room), args.Error(1)
}

func (m *MockRoomService) JoinRoom(roomID, playerID string) (*room.Room, error) {
	args := m.Called(roomID, playerID)
	return args.Get(0).(*room.Room), args.Error(1)
}

func (m *MockRoomService) LeaveRoom(roomID, playerID string) (*room.Room, error) {
	args := m.Called(roomID, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*room.Room), args.Error(1)
}

func (m *MockRoomService) GetRoom(roomID string) (*room.Room, error) {
	args := m.Called(roomID)
	return args.Get(0).(*room.Room), args.Error(1)
}

func (m *MockRoomService) GetRoomList(page, limit int, statusFilter *room.RoomStatus) (*room.RoomListResponse, error) {
	args := m.Called(page, limit, statusFilter)
	return args.Get(0).(*room.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) StartGame(roomID, playerID string) error {
	args := m.Called(roomID, playerID)
	return args.Error(0)
}

func (m *MockRoomService) CloseRoom(roomID string) error {
	args := m.Called(roomID)
	return args.Error(0)
}

func (m *MockRoomService) GetPlayerRoom(playerID string) (*room.Room, error) {
	args := m.Called(playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*room.Room), args.Error(1)
}

func TestNewWSManager(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	
	manager := NewWSManager(mockAuth, mockRoom)
	
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.connections)
	assert.NotNil(t, manager.rooms)
	assert.NotNil(t, manager.register)
	assert.NotNil(t, manager.unregister)
	assert.NotNil(t, manager.broadcast)
	assert.Equal(t, 30*time.Second, manager.pingInterval)
	assert.Equal(t, 60*time.Second, manager.pongTimeout)
	assert.NotNil(t, manager.messageHandlers)
	
	// Check default handlers are registered
	assert.Contains(t, manager.messageHandlers, "ping")
	assert.Contains(t, manager.messageHandlers, "join_room")
	assert.Contains(t, manager.messageHandlers, "leave_room")
	assert.Contains(t, manager.messageHandlers, "start_game")
}

func TestWSManager_RegisterHandler(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Register custom handler
	customHandler := func(conn *WSConnection, message *WSMessage) error {
		return nil
	}
	
	manager.RegisterHandler("custom_message", customHandler)
	
	assert.Contains(t, manager.messageHandlers, "custom_message")
}

func TestWSManager_HandleWebSocket_InvalidPlayer(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Mock auth service to return error
	mockAuth.On("GetUserByID", "invalid_player").Return((*auth.User)(nil), assert.AnError)
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := manager.HandleWebSocket(w, r, "invalid_player")
		assert.Error(t, err)
	}))
	defer server.Close()
	
	// Make request
	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	
	mockAuth.AssertExpectations(t)
}

func TestWSManager_HandleWebSocket_ValidPlayer(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Start manager in background
	done := make(chan bool)
	go func() {
		defer func() { done <- true }()
		manager.Run()
	}()
	defer func() {
		close(manager.register)
		close(manager.unregister)
		close(manager.broadcast)
		<-done // Wait for Run to finish
	}()
	
	// Mock auth service to return valid user
	user := &auth.User{
		ID:       "player1",
		Username: "testuser",
		Online:   true,
	}
	mockAuth.On("GetUserByID", "player1").Return(user, nil)
	
	// Mock room service to return no room (player not in any room)
	mockRoom.On("GetPlayerRoom", "player1").Return((*room.Room)(nil), assert.AnError)
	
	// Create test server with WebSocket upgrade
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := manager.HandleWebSocket(w, r, "player1")
		assert.NoError(t, err)
	}))
	defer server.Close()
	
	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	
	// Connect via WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()
	
	// Give some time for connection to be registered
	time.Sleep(200 * time.Millisecond)
	
	// Check that connection is registered
	assert.True(t, manager.IsPlayerConnected("player1"))
	
	mockAuth.AssertExpectations(t)
	mockRoom.AssertExpectations(t)
}

func TestWSManager_BroadcastToRoom(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Start manager in background
	done := make(chan bool)
	go func() {
		defer func() { done <- true }()
		manager.Run()
	}()
	defer func() {
		close(manager.register)
		close(manager.unregister)
		close(manager.broadcast)
		<-done // Wait for Run to finish
	}()
	
	// Create mock connections
	conn1 := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	conn2 := &WSConnection{
		playerID: "player2",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Register connections manually
	manager.mu.Lock()
	manager.connections["player1"] = conn1
	manager.connections["player2"] = conn2
	manager.rooms["room1"] = map[string]*WSConnection{
		"player1": conn1,
		"player2": conn2,
	}
	manager.mu.Unlock()
	
	// Create test message
	message := &WSMessage{
		Type: "test_message",
		Data: map[string]string{"content": "hello"},
		Timestamp: time.Now(),
	}
	
	// Broadcast message
	manager.BroadcastToRoom("room1", message)
	
	// Give some time for broadcast to process
	time.Sleep(100 * time.Millisecond)
	
	// Check that both connections received the message
	select {
	case msg1 := <-conn1.send:
		var receivedMsg WSMessage
		err := json.Unmarshal(msg1, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, "test_message", receivedMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Connection 1 did not receive message")
	}
	
	select {
	case msg2 := <-conn2.send:
		var receivedMsg WSMessage
		err := json.Unmarshal(msg2, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, "test_message", receivedMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Connection 2 did not receive message")
	}
}

func TestWSManager_BroadcastToRoomExcept(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Start manager in background
	done := make(chan bool)
	go func() {
		defer func() { done <- true }()
		manager.Run()
	}()
	defer func() {
		close(manager.register)
		close(manager.unregister)
		close(manager.broadcast)
		<-done // Wait for Run to finish
	}()
	
	// Create mock connections
	conn1 := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	conn2 := &WSConnection{
		playerID: "player2",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Register connections manually
	manager.mu.Lock()
	manager.connections["player1"] = conn1
	manager.connections["player2"] = conn2
	manager.rooms["room1"] = map[string]*WSConnection{
		"player1": conn1,
		"player2": conn2,
	}
	manager.mu.Unlock()
	
	// Create test message
	message := &WSMessage{
		Type: "test_message",
		Data: map[string]string{"content": "hello"},
		Timestamp: time.Now(),
	}
	
	// Broadcast message excluding player1
	manager.BroadcastToRoomExcept("room1", message, "player1")
	
	// Give some time for broadcast to process
	time.Sleep(100 * time.Millisecond)
	
	// Check that only player2 received the message
	select {
	case <-conn1.send:
		t.Fatal("Connection 1 should not have received message")
	case <-time.After(100 * time.Millisecond):
		// Expected - player1 should not receive message
	}
	
	select {
	case msg2 := <-conn2.send:
		var receivedMsg WSMessage
		err := json.Unmarshal(msg2, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, "test_message", receivedMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Connection 2 did not receive message")
	}
}

func TestWSManager_SendToPlayer(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Register connection manually
	manager.mu.Lock()
	manager.connections["player1"] = conn
	manager.mu.Unlock()
	
	// Create test message
	message := &WSMessage{
		Type: "test_message",
		Data: map[string]string{"content": "hello"},
		Timestamp: time.Now(),
	}
	
	// Send message to player
	err := manager.SendToPlayer("player1", message)
	assert.NoError(t, err)
	
	// Check that connection received the message
	select {
	case msg := <-conn.send:
		var receivedMsg WSMessage
		err := json.Unmarshal(msg, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, "test_message", receivedMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Connection did not receive message")
	}
}

func TestWSManager_SendToPlayer_NotConnected(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create test message
	message := &WSMessage{
		Type: "test_message",
		Data: map[string]string{"content": "hello"},
		Timestamp: time.Now(),
	}
	
	// Try to send message to non-connected player
	err := manager.SendToPlayer("nonexistent_player", message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestWSManager_GetRoomConnections(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Initially no connections
	count := manager.GetRoomConnections("room1")
	assert.Equal(t, 0, count)
	
	// Add connections manually
	manager.mu.Lock()
	manager.rooms["room1"] = map[string]*WSConnection{
		"player1": &WSConnection{playerID: "player1"},
		"player2": &WSConnection{playerID: "player2"},
	}
	manager.mu.Unlock()
	
	// Check connection count
	count = manager.GetRoomConnections("room1")
	assert.Equal(t, 2, count)
}

func TestWSManager_IsPlayerConnected(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Initially not connected
	assert.False(t, manager.IsPlayerConnected("player1"))
	
	// Add connection manually
	manager.mu.Lock()
	manager.connections["player1"] = &WSConnection{playerID: "player1"}
	manager.mu.Unlock()
	
	// Check connection status
	assert.True(t, manager.IsPlayerConnected("player1"))
}