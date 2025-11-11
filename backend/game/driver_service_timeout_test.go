package game

import (
	"context"
	"testing"
	"time"
	
	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// TestRequestPlayDecision_TimeoutInMessage verifies timeout is derived from context deadline
func TestRequestPlayDecision_TimeoutInMessage(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	// Create context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Start request in goroutine (will block)
	go func() {
		provider.RequestPlayDecision(ctx, 0, []*sdk.Card{}, &sdk.TrickInfo{IsLeader: true})
	}()
	
	// Wait for message to be sent
	time.Sleep(50 * time.Millisecond)
	
	// Check broadcast
	broadcasts := wsManager.GetBroadcasts("test-room")
	if len(broadcasts) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(broadcasts))
	}
	
	msg := broadcasts[0]
	if msg.Type != websocket.MSG_GAME_ACTION {
		t.Fatalf("Expected MSG_GAME_ACTION, got %s", msg.Type)
	}
	
	data := msg.Data.(map[string]interface{})
	timeout, ok := data["timeout"]
	if !ok {
		t.Fatal("Expected timeout field in message")
	}
	
	timeoutSecs := timeout.(int)
	// Should be around 10 seconds (9-10 due to slight delay)
	if timeoutSecs < 9 || timeoutSecs > 10 {
		t.Errorf("Expected timeout around 10 seconds, got %d", timeoutSecs)
	}
	
	// Cancel to cleanup
	provider.CancelAll()
}

// TestRequestPlayDecision_NoDeadline verifies timeout field when no deadline
func TestRequestPlayDecision_NoDeadline(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	// Create context without deadline
	ctx := context.Background()
	
	// Start request in goroutine (will block)
	go func() {
		provider.RequestPlayDecision(ctx, 0, []*sdk.Card{}, &sdk.TrickInfo{IsLeader: true})
	}()
	
	// Wait for message to be sent
	time.Sleep(50 * time.Millisecond)
	
	// Check broadcast
	broadcasts := wsManager.GetBroadcasts("test-room")
	if len(broadcasts) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(broadcasts))
	}
	
	msg := broadcasts[0]
	data := msg.Data.(map[string]interface{})
	
	// timeout field should not exist when no deadline
	_, ok := data["timeout"]
	if ok {
		t.Error("Expected no timeout field when context has no deadline")
	}
	
	// Cancel to cleanup
	provider.CancelAll()
}
