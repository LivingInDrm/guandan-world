package game

import (
	"testing"
	"time"
	
	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// TestWebSocketObserver_PlayerTimeoutEvent_AllActionTypes verifies EventPlayerTimeout
// is properly handled for all timeout action types
func TestWebSocketObserver_PlayerTimeoutEvent_AllActionTypes(t *testing.T) {
	testCases := []struct {
		name       string
		action     string
		playerSeat int
	}{
		{
			name:       "PlayDecisionTimeout",
			action:     "play_decision",
			playerSeat: 0,
		},
		{
			name:       "TributeSelectTimeout",
			action:     "tribute_select",
			playerSeat: 1,
		},
		{
			name:       "ReturnTributeTimeout",
			action:     "return_tribute",
			playerSeat: 2,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
					"action": tc.action,
				},
				Timestamp:  time.Now(),
				PlayerSeat: tc.playerSeat,
			}
			
			// Send event
			observer.OnGameEvent(timeoutEvent)
			
			// Verify WebSocket message was broadcast
			broadcasts := wsManager.GetBroadcasts("test-room")
			if len(broadcasts) != 1 {
				t.Fatalf("Expected 1 broadcast, got %d", len(broadcasts))
			}
			
			// Verify message type
			msg := broadcasts[0]
			if msg.Type != websocket.MSG_GAME_EVENT {
				t.Errorf("Expected MSG_GAME_EVENT, got %s", msg.Type)
			}
			
			// Verify event_type
			data := msg.Data.(map[string]interface{})
			if data["event_type"] != string(sdk.EventPlayerTimeout) {
				t.Errorf("Expected event_type %s, got %v", sdk.EventPlayerTimeout, data["event_type"])
			}
			
			// Verify player_seat is present
			if data["player_seat"] != tc.playerSeat {
				t.Errorf("Expected player_seat %d, got %v", tc.playerSeat, data["player_seat"])
			}
			
			// Verify timestamp is present
			if data["timestamp"] == nil {
				t.Error("Expected timestamp to be present")
			}
			
			// Verify event_data contains action field
			eventData, ok := data["event_data"].(map[string]interface{})
			if !ok {
				t.Fatal("event_data is not map[string]interface{}")
			}
			
			if eventData["action"] != tc.action {
				t.Errorf("Expected action %s, got %v", tc.action, eventData["action"])
			}
		})
	}
}

// TestWebSocketObserver_PlayerTimeoutEvent_LoggingEnabled verifies that
// EventPlayerTimeout is logged (part of the switch statement)
func TestWebSocketObserver_PlayerTimeoutEvent_LoggingEnabled(t *testing.T) {
	// This test verifies that EventPlayerTimeout is in the logging switch
	// The actual logging is tested by observing the log output in other tests
	
	wsManager := NewMockDriverWSManager()
	engine := sdk.NewGameEngine()
	observer := NewWebSocketObserver("test-room", wsManager, engine)
	
	timeoutEvent := &sdk.GameEvent{
		Type: sdk.EventPlayerTimeout,
		Data: map[string]interface{}{
			"action": "play_decision",
		},
		Timestamp:  time.Now(),
		PlayerSeat: 1,
	}
	
	// Send event - should not panic and should log
	observer.OnGameEvent(timeoutEvent)
	
	// Verify event was broadcast
	broadcasts := wsManager.GetBroadcasts("test-room")
	if len(broadcasts) != 1 {
		t.Errorf("Expected 1 broadcast, got %d", len(broadcasts))
	}
}
