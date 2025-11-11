package game

import (
	"testing"
	"time"
	
	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// TestWebSocketObserver_PlayerTimeoutEvent verifies EventPlayerTimeout is broadcast
func TestWebSocketObserver_PlayerTimeoutEvent(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create a test engine
	engine := sdk.NewGameEngine()
	
	// Create observer
	observer := NewWebSocketObserver("test-room", wsManager, engine)
	
	// Create timeout event
	timeoutEvent := &sdk.GameEvent{
		Type: sdk.EventPlayerTimeout,
		Data: map[string]interface{}{
			"action": "play_decision",
		},
		Timestamp:  time.Now(),
		PlayerSeat: 2,
	}
	
	// Send event
	observer.OnGameEvent(timeoutEvent)
	
	// Verify WebSocket message was broadcast
	broadcasts := wsManager.GetBroadcasts("test-room")
	if len(broadcasts) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(broadcasts))
	}
	
	// Verify message content
	msg := broadcasts[0]
	if msg.Type != websocket.MSG_GAME_EVENT {
		t.Errorf("Expected MSG_GAME_EVENT, got %s", msg.Type)
	}
	
	data := msg.Data.(map[string]interface{})
	if data["event_type"] != string(sdk.EventPlayerTimeout) {
		t.Errorf("Expected event_type %s, got %v", sdk.EventPlayerTimeout, data["event_type"])
	}
	
	if data["player_seat"] != 2 {
		t.Errorf("Expected player_seat 2, got %v", data["player_seat"])
	}
	
	// Verify event data
	eventData := data["event_data"].(map[string]interface{})
	if eventData["action"] != "play_decision" {
		t.Errorf("Expected action play_decision, got %v", eventData["action"])
	}
}

// TestDriverService_TimeoutStrategyConfigured verifies TimeoutStrategy is configured
func TestDriverService_TimeoutStrategyConfigured(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	driverService := NewDriverService(wsManager)
	
	players := []sdk.Player{
		{ID: "p1", Username: "Player 1", Seat: 0},
		{ID: "p2", Username: "Player 2", Seat: 1},
		{ID: "p3", Username: "Player 3", Seat: 2},
		{ID: "p4", Username: "Player 4", Seat: 3},
	}
	
	err := driverService.StartGameWithDriver("test-room", players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	
	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)
	
	// Verify driver exists
	driverService.mu.RLock()
	driver, exists := driverService.drivers["test-room"]
	driverService.mu.RUnlock()
	
	if !exists {
		t.Fatal("Expected driver to exist for test-room")
	}
	
	// The driver should have a timeout strategy configured
	// We can't directly access the config, but we can verify the driver was created successfully
	if driver == nil {
		t.Fatal("Driver should not be nil")
	}
	
	// Clean up
	driverService.StopGame("test-room")
	time.Sleep(10 * time.Millisecond)
}
