package game

import (
	"context"
	"testing"
	"time"
	
	"guandan-world/sdk"
)

// TestRequestPlayDecision_NilContext verifies nil context check
func TestRequestPlayDecision_NilContext(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	_, err := provider.RequestPlayDecision(nil, 0, nil, nil)
	if err == nil {
		t.Fatal("Expected error for nil context, got nil")
	}
	if err.Error() != "context must not be nil" {
		t.Errorf("Expected 'context must not be nil', got %v", err)
	}
}

// TestRequestTributeSelection_NilContext verifies nil context check
func TestRequestTributeSelection_NilContext(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	_, err := provider.RequestTributeSelection(nil, 0, nil)
	if err == nil {
		t.Fatal("Expected error for nil context, got nil")
	}
	if err.Error() != "context must not be nil" {
		t.Errorf("Expected 'context must not be nil', got %v", err)
	}
}

// TestRequestReturnTribute_NilContext verifies nil context check
func TestRequestReturnTribute_NilContext(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	_, err := provider.RequestReturnTribute(nil, 0, nil)
	if err == nil {
		t.Fatal("Expected error for nil context, got nil")
	}
	if err.Error() != "context must not be nil" {
		t.Errorf("Expected 'context must not be nil', got %v", err)
	}
}

// TestSubmitTributeSelection_NilCard verifies nil card check
func TestSubmitTributeSelection_NilCard(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	err := provider.SubmitTributeSelection(0, nil)
	if err == nil {
		t.Fatal("Expected error for nil card, got nil")
	}
	if err.Error() != "card cannot be nil for player 0" {
		t.Errorf("Expected 'card cannot be nil for player 0', got %v", err)
	}
}

// TestSubmitReturnTribute_NilCard verifies nil card check
func TestSubmitReturnTribute_NilCard(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	err := provider.SubmitReturnTribute(0, nil)
	if err == nil {
		t.Fatal("Expected error for nil card, got nil")
	}
	if err.Error() != "card cannot be nil for player 0" {
		t.Errorf("Expected 'card cannot be nil for player 0', got %v", err)
	}
}

// TestRequestPlayDecision_CanceledChannel verifies closed channel handling
func TestRequestPlayDecision_CanceledChannel(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	provider := NewRoomInputProvider("test-room", wsManager)
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Start a goroutine that will cancel all requests
	go func() {
		time.Sleep(20 * time.Millisecond)
		provider.CancelAll()
	}()
	
	_, err := provider.RequestPlayDecision(ctx, 0, []*sdk.Card{}, &sdk.TrickInfo{IsLeader: true})
	if err == nil {
		t.Fatal("Expected error for canceled request, got nil")
	}
	if err.Error() != "play decision request canceled for player 0" {
		t.Errorf("Expected 'play decision request canceled for player 0', got %v", err)
	}
}
