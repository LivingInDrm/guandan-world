package game

import (
	"testing"
	"time"

	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// TestEventDrivenStateSynchronization tests the complete event-driven state sync flow
func TestEventDrivenStateSynchronization(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Wait for initial events to be processed
	time.Sleep(100 * time.Millisecond)
	
	initialBroadcastCount := mockWS.GetBroadcastCount()
	if initialBroadcastCount == 0 {
		t.Error("No initial events were broadcast")
	}
	
	// Test that match started and deal started events were sent
	broadcasts := mockWS.broadcastMessages
	foundMatchStarted := false
	foundDealStarted := false
	
	for _, broadcast := range broadcasts {
		if broadcast.Message.Type == websocket.MSG_GAME_EVENT {
			if data, ok := broadcast.Message.Data.(map[string]interface{}); ok {
				eventType, exists := data["event_type"]
				if exists {
					switch eventType {
					case "match_started":
						foundMatchStarted = true
					case "deal_started":
						foundDealStarted = true
					}
				}
			}
		}
	}
	
	if !foundMatchStarted {
		t.Error("Match started event was not broadcast")
	}
	
	if !foundDealStarted {
		t.Error("Deal started event was not broadcast")
	}
}

// TestPlayerViewFiltering tests that player views are properly filtered
func TestPlayerViewFiltering(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Test individual player views
	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		playerView, err := service.GetPlayerView("room1", playerSeat)
		if err != nil {
			t.Errorf("GetPlayerView failed for player %d: %v", playerSeat, err)
			continue
		}
		
		if playerView == nil {
			t.Errorf("Player view is nil for player %d", playerSeat)
			continue
		}
		
		// Verify player seat matches
		if playerView.PlayerSeat != playerSeat {
			t.Errorf("Expected player seat %d, got %d", playerSeat, playerView.PlayerSeat)
		}
		
		// Verify player has their own cards (should not be empty after deal)
		if len(playerView.PlayerCards) == 0 {
			t.Errorf("Player %d has no cards", playerSeat)
		}
		
		// Verify game state is present
		if playerView.GameState == nil {
			t.Errorf("Game state is nil for player %d", playerSeat)
		}
	}
}

// TestTimeoutProcessing tests the timeout processing functionality
func TestTimeoutProcessing(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Clear initial messages
	mockWS.broadcastMessages = []MockBroadcastMessage{}
	
	// Manually trigger timeout processing
	service.processTimeouts()
	
	// Note: Since we don't have actual timeouts in the test scenario,
	// we mainly verify that the timeout processing doesn't crash
	// and that the method can be called safely
	
	// The timeout processing should not generate events if there are no actual timeouts
	// This is expected behavior
}

// TestSyncGameState tests the complete game state synchronization
func TestSyncGameState(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Clear initial messages
	mockWS.broadcastMessages = []MockBroadcastMessage{}
	mockWS.playerMessages = []MockPlayerMessage{}
	
	// Test sync game state
	err = service.SyncGameState("room1")
	if err != nil {
		t.Errorf("SyncGameState failed: %v", err)
	}
	
	// Verify that sync event was broadcast
	if mockWS.GetBroadcastCount() == 0 {
		t.Error("No sync events were broadcast")
	}
	
	// Verify that player views were sent
	if mockWS.GetPlayerMessageCount() == 0 {
		t.Error("No player view messages were sent during sync")
	}
	
	// Test sync for nonexistent room
	err = service.SyncGameState("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

// TestBroadcastGameEvent tests custom event broadcasting
func TestBroadcastGameEvent(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Clear any initial messages
	mockWS.broadcastMessages = []MockBroadcastMessage{}
	
	// Test broadcasting custom event
	eventData := map[string]interface{}{
		"message": "test event",
		"value":   42,
	}
	
	err := service.BroadcastGameEvent("room1", "custom_event", eventData)
	if err != nil {
		t.Errorf("BroadcastGameEvent failed: %v", err)
	}
	
	// Verify event was broadcast
	if mockWS.GetBroadcastCount() != 1 {
		t.Errorf("Expected 1 broadcast message, got %d", mockWS.GetBroadcastCount())
	}
	
	lastBroadcast := mockWS.GetLastBroadcast()
	if lastBroadcast == nil {
		t.Error("No broadcast message found")
	} else {
		if lastBroadcast.RoomID != "room1" {
			t.Errorf("Expected room ID 'room1', got '%s'", lastBroadcast.RoomID)
		}
		
		if lastBroadcast.Message.Type != websocket.MSG_GAME_EVENT {
			t.Errorf("Expected message type '%s', got '%s'", 
				websocket.MSG_GAME_EVENT, lastBroadcast.Message.Type)
		}
	}
}

// TestGetRealTimeGameStatus tests real-time status retrieval
func TestGetRealTimeGameStatus(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Test getting real-time status
	status, err := service.GetRealTimeGameStatus("room1")
	if err != nil {
		t.Errorf("GetRealTimeGameStatus failed: %v", err)
	}
	
	if status == nil {
		t.Error("Status is nil")
	} else {
		// Verify required fields are present
		requiredFields := []string{"game_status", "room_id", "timestamp"}
		for _, field := range requiredFields {
			if _, exists := status[field]; !exists {
				t.Errorf("Required field '%s' missing from status", field)
			}
		}
		
		// Verify room ID matches
		if roomID, exists := status["room_id"]; exists {
			if roomID != "room1" {
				t.Errorf("Expected room ID 'room1', got '%v'", roomID)
			}
		}
		
		// Verify game status is valid
		if gameStatus, exists := status["game_status"]; exists {
			if gameStatus != sdk.GameStatusStarted {
				t.Errorf("Expected game status '%s', got '%v'", 
					sdk.GameStatusStarted, gameStatus)
			}
		}
	}
	
	// Test getting status for nonexistent room
	_, err = service.GetRealTimeGameStatus("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

// TestCreateFilteredState tests the player state filtering functionality
func TestCreateFilteredState(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Get player view for testing
	playerView, err := service.GetPlayerView("room1", 0)
	if err != nil {
		t.Fatalf("GetPlayerView failed: %v", err)
	}
	
	// Test filtered state creation
	filteredState := service.createFilteredState(playerView, 0)
	
	if filteredState == nil {
		t.Error("Filtered state is nil")
	} else {
		// Verify required fields are present
		requiredFields := []string{"player_seat", "player_cards", "visible_cards"}
		for _, field := range requiredFields {
			if _, exists := filteredState[field]; !exists {
				t.Errorf("Required field '%s' missing from filtered state", field)
			}
		}
		
		// Verify player seat matches
		if playerSeat, exists := filteredState["player_seat"]; exists {
			if playerSeat != 0 {
				t.Errorf("Expected player seat 0, got %v", playerSeat)
			}
		}
		
		// Verify players information is present and properly filtered
		if players, exists := filteredState["players"]; exists {
			if playersSlice, ok := players.([]map[string]interface{}); ok {
				if len(playersSlice) != 4 {
					t.Errorf("Expected 4 players in filtered state, got %d", len(playersSlice))
				}
				
				// Verify each player has required fields but no private data
				for i, player := range playersSlice {
					if player != nil {
						requiredPlayerFields := []string{"id", "username", "seat", "online", "auto_play"}
						for _, field := range requiredPlayerFields {
							if _, exists := player[field]; !exists {
								t.Errorf("Player %d missing required field '%s'", i, field)
							}
						}
						
						// Verify no private data like cards is included
						if _, exists := player["cards"]; exists {
							t.Errorf("Player %d should not have cards in filtered state", i)
						}
					}
				}
			} else {
				t.Error("Players field is not in expected format")
			}
		}
	}
}

// TestEventHandlerIntegration tests the complete event handler integration
func TestEventHandlerIntegration(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game
	players := []sdk.Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	err := service.StartGame("room1", players)
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}
	
	// Wait for events to be processed
	time.Sleep(100 * time.Millisecond)
	
	// Verify that events were properly converted and broadcast
	broadcasts := mockWS.broadcastMessages
	if len(broadcasts) == 0 {
		t.Error("No events were broadcast")
	}
	
	// Verify event structure
	for _, broadcast := range broadcasts {
		if broadcast.Message.Type != websocket.MSG_GAME_EVENT {
			continue
		}
		
		// Verify event data structure
		if data, ok := broadcast.Message.Data.(map[string]interface{}); ok {
			requiredFields := []string{"event_type", "event_data", "timestamp"}
			for _, field := range requiredFields {
				if _, exists := data[field]; !exists {
					t.Errorf("Event missing required field '%s'", field)
				}
			}
		} else {
			t.Error("Event data is not in expected format")
		}
	}
	
	// Verify that player views were sent for appropriate events
	playerMessages := mockWS.playerMessages
	if len(playerMessages) == 0 {
		t.Error("No player view messages were sent")
	}
	
	// Verify player view message structure
	for _, message := range playerMessages {
		if message.Message.Type != websocket.MSG_PLAYER_VIEW {
			continue
		}
		
		if data, ok := message.Message.Data.(map[string]interface{}); ok {
			requiredFields := []string{"player_view", "event_type", "player_seat", "filtered_state"}
			for _, field := range requiredFields {
				if _, exists := data[field]; !exists {
					t.Errorf("Player view message missing required field '%s'", field)
				}
			}
		} else {
			t.Error("Player view message data is not in expected format")
		}
	}
}