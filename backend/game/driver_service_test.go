package game

import (
	"context"
	"testing"
	"time"

	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// MockDriverWSManager implements WSManagerInterface for testing driver service
type MockDriverWSManager struct {
	broadcasts map[string][]*websocket.WSMessage
	messages   map[string][]*websocket.WSMessage
}

func NewMockDriverWSManager() *MockDriverWSManager {
	return &MockDriverWSManager{
		broadcasts: make(map[string][]*websocket.WSMessage),
		messages:   make(map[string][]*websocket.WSMessage),
	}
}

func (m *MockDriverWSManager) BroadcastToRoom(roomID string, message *websocket.WSMessage) {
	if m.broadcasts[roomID] == nil {
		m.broadcasts[roomID] = make([]*websocket.WSMessage, 0)
	}
	m.broadcasts[roomID] = append(m.broadcasts[roomID], message)
}

func (m *MockDriverWSManager) SendToPlayer(playerID string, message *websocket.WSMessage) error {
	if m.messages[playerID] == nil {
		m.messages[playerID] = make([]*websocket.WSMessage, 0)
	}
	m.messages[playerID] = append(m.messages[playerID], message)
	return nil
}

func (m *MockDriverWSManager) GetBroadcasts(roomID string) []*websocket.WSMessage {
	return m.broadcasts[roomID]
}

func (m *MockDriverWSManager) GetMessages(playerID string) []*websocket.WSMessage {
	return m.messages[playerID]
}

func TestDriverService_StartGameWithDriver(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create driver service
	service := NewDriverService(wsManager)
	
	// Create test players
	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}
	
	// Test successful game start
	roomID := "test-room-1"
	err := service.StartGameWithDriver(roomID, players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	
	// Verify game was created
	status, err := service.GetGameStatus(roomID)
	if err != nil {
		t.Fatalf("Failed to get game status: %v", err)
	}
	
	if status["room_id"] != roomID {
		t.Errorf("Expected room_id %s, got %v", roomID, status["room_id"])
	}
	
	// Test duplicate game creation before stopping
	err = service.StartGameWithDriver(roomID, players)
	if err == nil {
		t.Error("Expected error when creating duplicate game")
	}
	
	// Stop the game to prevent background goroutine crashes
	service.StopGame(roomID)
	
	// Test with invalid player count
	invalidPlayers := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
	}
	
	err = service.StartGameWithDriver("test-room-2", invalidPlayers)
	if err == nil {
		t.Error("Expected error with invalid player count")
	}
}

func TestDriverService_SubmitPlayDecision(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create driver service
	service := NewDriverService(wsManager)
	
	// Create test players
	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}
	
	// Start game
	roomID := "test-room-play"
	err := service.StartGameWithDriver(roomID, players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	
	// Wait for game to initialize
	time.Sleep(100 * time.Millisecond)
	
	// Test submitting play decision for non-existent room
	err = service.SubmitPlayDecision("non-existent-room", 0, &sdk.PlayDecision{
		Action: sdk.ActionPass,
	})
	if err == nil {
		t.Error("Expected error for non-existent room")
	}
	
	// Stop the game to prevent background goroutine crashes
	service.StopGame(roomID)
	
	// Note: Testing actual play decisions would require waiting for the game
	// to request input, which would need more complex test setup
}

func TestDriverService_StopGame(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create driver service
	service := NewDriverService(wsManager)
	
	// Create test players
	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}
	
	// Start game
	roomID := "test-room-stop"
	err := service.StartGameWithDriver(roomID, players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	
	// Stop game
	err = service.StopGame(roomID)
	if err != nil {
		t.Fatalf("Failed to stop game: %v", err)
	}
	
	// Verify game was stopped
	_, err = service.GetGameStatus(roomID)
	if err == nil {
		t.Error("Expected error when getting status of stopped game")
	}
	
	// Test stopping non-existent game
	err = service.StopGame("non-existent-room")
	if err == nil {
		t.Error("Expected error when stopping non-existent game")
	}
}

func TestRoomInputProvider_RequestPlayDecision(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create input provider
	provider := NewRoomInputProvider("test-room", wsManager)
	
	// Test cards
	cards := []*sdk.Card{
		mustNewCard(t, 3, "Heart", 2),
		mustNewCard(t, 4, "Spade", 2),
		mustNewCard(t, 5, "Diamond", 2),
	}
	
	// Create a goroutine to submit decision
	go func() {
		time.Sleep(50 * time.Millisecond)
		decision := &sdk.PlayDecision{
			Action: sdk.ActionPlay,
			Cards:  []*sdk.Card{cards[0]},
		}
		err := provider.SubmitPlayDecision(0, decision)
		if err != nil {
			t.Errorf("Failed to submit play decision: %v", err)
		}
	}()
	
	// Request play decision
	trickInfo := &sdk.TrickInfo{IsLeader: true}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	decision, err := provider.RequestPlayDecision(ctx, 0, cards, trickInfo)
	if err != nil {
		t.Fatalf("Failed to request play decision: %v", err)
	}
	
	// Verify decision
	if decision.Action != sdk.ActionPlay {
		t.Errorf("Expected play action, got %v", decision.Action)
	}
	if len(decision.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(decision.Cards))
	}
	
	// Verify WebSocket message was sent
	broadcasts := wsManager.GetBroadcasts("test-room")
	if len(broadcasts) == 0 {
		t.Error("Expected WebSocket broadcast for play request")
	}
}

func TestRoomInputProvider_CancelAll(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create input provider
	provider := NewRoomInputProvider("test-room", wsManager)
	
	// Create pending requests
	provider.playDecisions[0] = make(chan *sdk.PlayDecision, 1)
	provider.tributeSelections[1] = make(chan *sdk.Card, 1)
	provider.returnTributes[2] = make(chan *sdk.Card, 1)
	
	// Cancel all
	provider.CancelAll()
	
	// Verify all channels were cleared
	if len(provider.playDecisions) != 0 {
		t.Error("Expected play decisions to be cleared")
	}
	if len(provider.tributeSelections) != 0 {
		t.Error("Expected tribute selections to be cleared")
	}
	if len(provider.returnTributes) != 0 {
		t.Error("Expected return tributes to be cleared")
	}
}

func TestWebSocketObserver_OnGameEvent(t *testing.T) {
	// Create mock WebSocket manager
	wsManager := NewMockDriverWSManager()
	
	// Create observer
	observer := NewWebSocketObserver("test-room", wsManager)
	
	// Create test event
	event := &sdk.GameEvent{
		Type: sdk.EventDealStarted,
		Data: map[string]interface{}{
			"deal_level": 2,
		},
		Timestamp:  time.Now(),
		PlayerSeat: 0,
	}
	
	// Send event
	observer.OnGameEvent(event)
	
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
	if data["event_type"] != string(sdk.EventDealStarted) {
		t.Errorf("Expected event_type %s, got %v", sdk.EventDealStarted, data["event_type"])
	}
}

// Helper function to create cards
func mustNewCard(t *testing.T, number int, color string, level int) *sdk.Card {
	card, err := sdk.NewCard(number, color, level)
	if err != nil {
		t.Fatalf("Failed to create card: %v", err)
	}
	return card
}