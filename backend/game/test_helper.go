package game

import (
	"log"
	"guandan-world/backend/websocket"
)

// TestWebSocketManager wraps WSManager for testing scenarios
type TestWebSocketManager struct {
	*websocket.WSManager
	testMode bool
	testConnections map[string]bool // roomID -> has test connection
}

// NewTestWebSocketManager creates a test-friendly WebSocket manager
func NewTestWebSocketManager(wsManager *websocket.WSManager) *TestWebSocketManager {
	return &TestWebSocketManager{
		WSManager: wsManager,
		testMode: true,
		testConnections: make(map[string]bool),
	}
}

// EnableTestMode marks a room as having test connections
func (t *TestWebSocketManager) EnableTestMode(roomID string) {
	t.testConnections[roomID] = true
	log.Printf("[TestWebSocketManager] Test mode enabled for room %s", roomID)
}

// BroadcastToRoom overrides to handle test scenarios
func (t *TestWebSocketManager) BroadcastToRoom(roomID string, message *websocket.WSMessage) {
	// Check if this is a test room
	if t.testConnections[roomID] {
		// In test mode, we still broadcast even if no connections
		log.Printf("[TestWebSocketManager] Broadcasting to test room %s: %v", roomID, message.Type)
	}
	
	// Call original method
	t.WSManager.BroadcastToRoom(roomID, message)
}

// GetRoomConnections overrides to return non-zero for test rooms
func (t *TestWebSocketManager) GetRoomConnections(roomID string) int {
	if t.testConnections[roomID] {
		return 1 // Pretend there's at least one connection
	}
	return t.WSManager.GetRoomConnections(roomID)
}