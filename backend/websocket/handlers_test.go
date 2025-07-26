package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"guandan-world/backend/room"
)

func TestHandlePing(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		send:     make(chan []byte, 256),
		manager:  manager,
		lastPing: time.Now().Add(-time.Minute), // Old ping time
	}
	
	// Create ping message
	message := &WSMessage{
		Type:      MSG_PING,
		Data:      nil,
		Timestamp: time.Now(),
	}
	
	// Handle ping
	err := manager.handlePing(conn, message)
	assert.NoError(t, err)
	
	// Check that pong was sent
	select {
	case pongData := <-conn.send:
		var pongMsg WSMessage
		err := json.Unmarshal(pongData, &pongMsg)
		assert.NoError(t, err)
		assert.Equal(t, MSG_PONG, pongMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Pong message not received")
	}
	
	// Check that last ping time was updated
	conn.mu.RLock()
	assert.True(t, time.Since(conn.lastPing) < time.Second)
	conn.mu.RUnlock()
}

func TestHandleJoinRoom_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Mock room service
	testRoom := &room.Room{
		ID:          "room1",
		Status:      room.RoomStatusWaiting,
		PlayerCount: 1,
	}
	mockRoom.On("JoinRoom", "room1", "player1").Return(testRoom, nil)
	
	// Create join room message
	message := &WSMessage{
		Type: MSG_JOIN_ROOM,
		Data: map[string]interface{}{
			"room_id": "room1",
		},
		Timestamp: time.Now(),
	}
	
	// Handle join room
	err := manager.handleJoinRoom(conn, message)
	assert.NoError(t, err)
	
	// Check that connection room ID was updated
	assert.Equal(t, "room1", conn.roomID)
	
	// Check that connection was added to room
	manager.mu.RLock()
	roomConns, exists := manager.rooms["room1"]
	assert.True(t, exists)
	assert.Contains(t, roomConns, "player1")
	manager.mu.RUnlock()
	
	mockRoom.AssertExpectations(t)
}

func TestHandleJoinRoom_InvalidData(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create join room message with missing room_id
	message := &WSMessage{
		Type:      MSG_JOIN_ROOM,
		Data:      map[string]interface{}{},
		Timestamp: time.Now(),
	}
	
	// Handle join room
	err := manager.handleJoinRoom(conn, message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "room ID is required")
}

func TestHandleLeaveRoom_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Add connection to room manually
	manager.mu.Lock()
	manager.rooms["room1"] = map[string]*WSConnection{
		"player1": conn,
	}
	manager.mu.Unlock()
	
	// Mock room service
	testRoom := &room.Room{
		ID:          "room1",
		Status:      room.RoomStatusWaiting,
		PlayerCount: 0,
	}
	mockRoom.On("LeaveRoom", "room1", "player1").Return(testRoom, nil)
	
	// Create leave room message
	message := &WSMessage{
		Type: MSG_LEAVE_ROOM,
		Data: map[string]interface{}{
			"room_id": "room1",
		},
		Timestamp: time.Now(),
	}
	
	// Handle leave room
	err := manager.handleLeaveRoom(conn, message)
	assert.NoError(t, err)
	
	// Check that connection room ID was cleared
	assert.Equal(t, "", conn.roomID)
	
	// Check that connection was removed from room
	manager.mu.RLock()
	roomConns, exists := manager.rooms["room1"]
	if exists {
		assert.NotContains(t, roomConns, "player1")
	}
	manager.mu.RUnlock()
	
	mockRoom.AssertExpectations(t)
}

func TestHandleStartGame_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Mock room service
	testRoom := &room.Room{
		ID:          "room1",
		Status:      room.RoomStatusPlaying,
		PlayerCount: 4,
	}
	mockRoom.On("StartGame", "room1", "player1").Return(nil)
	mockRoom.On("GetRoom", "room1").Return(testRoom, nil)
	
	// Create start game message
	message := &WSMessage{
		Type: MSG_START_GAME,
		Data: map[string]interface{}{
			"room_id": "room1",
		},
		Timestamp: time.Now(),
	}
	
	// Handle start game
	err := manager.handleStartGame(conn, message)
	assert.NoError(t, err)
	
	mockRoom.AssertExpectations(t)
}

func TestHandlePlayCards_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create play cards message
	message := &WSMessage{
		Type: MSG_PLAY_CARDS,
		Data: map[string]interface{}{
			"cards": []string{"card1", "card2"},
		},
		Timestamp: time.Now(),
	}
	
	// Handle play cards
	err := manager.handlePlayCards(conn, message)
	assert.NoError(t, err)
}

func TestHandlePlayCards_NoCards(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create play cards message with no cards
	message := &WSMessage{
		Type: MSG_PLAY_CARDS,
		Data: map[string]interface{}{
			"cards": []string{},
		},
		Timestamp: time.Now(),
	}
	
	// Handle play cards
	err := manager.handlePlayCards(conn, message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no cards specified")
}

func TestHandlePass_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create pass message
	message := &WSMessage{
		Type:      MSG_PASS,
		Data:      nil,
		Timestamp: time.Now(),
	}
	
	// Handle pass
	err := manager.handlePass(conn, message)
	assert.NoError(t, err)
}

func TestHandleTributeSelect_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create tribute select message
	message := &WSMessage{
		Type: MSG_TRIBUTE_SELECT,
		Data: map[string]interface{}{
			"card_id": "card1",
		},
		Timestamp: time.Now(),
	}
	
	// Handle tribute select
	err := manager.handleTributeSelect(conn, message)
	assert.NoError(t, err)
}

func TestHandleTributeReturn_Success(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockRoom := &MockRoomService{}
	manager := NewWSManager(mockAuth, mockRoom)
	
	// Create mock connection
	conn := &WSConnection{
		playerID: "player1",
		roomID:   "room1",
		send:     make(chan []byte, 256),
		manager:  manager,
	}
	
	// Create tribute return message
	message := &WSMessage{
		Type: MSG_TRIBUTE_RETURN,
		Data: map[string]interface{}{
			"card_id": "card1",
		},
		Timestamp: time.Now(),
	}
	
	// Handle tribute return
	err := manager.handleTributeReturn(conn, message)
	assert.NoError(t, err)
}

func TestParseMessageData(t *testing.T) {
	// Test successful parsing
	data := map[string]interface{}{
		"room_id": "room1",
	}
	
	var target JoinRoomData
	err := parseMessageData(data, &target)
	assert.NoError(t, err)
	assert.Equal(t, "room1", target.RoomID)
	
	// Test parsing with missing field
	data2 := map[string]interface{}{
		"other_field": "value",
	}
	
	var target2 JoinRoomData
	err = parseMessageData(data2, &target2)
	assert.NoError(t, err)
	assert.Equal(t, "", target2.RoomID) // Should be empty string for missing field
}