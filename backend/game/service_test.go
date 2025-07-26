package game

import (
	"testing"
	"time"

	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// MockWSManager is a mock WebSocket manager for testing
type MockWSManager struct {
	broadcastMessages []MockBroadcastMessage
	playerMessages    []MockPlayerMessage
}

// Ensure MockWSManager implements WSManagerInterface
var _ WSManagerInterface = (*MockWSManager)(nil)

type MockBroadcastMessage struct {
	RoomID  string
	Message *websocket.WSMessage
}

type MockPlayerMessage struct {
	PlayerID string
	Message  *websocket.WSMessage
}

func NewMockWSManager() *MockWSManager {
	return &MockWSManager{
		broadcastMessages: make([]MockBroadcastMessage, 0),
		playerMessages:    make([]MockPlayerMessage, 0),
	}
}

func (m *MockWSManager) BroadcastToRoom(roomID string, message *websocket.WSMessage) {
	m.broadcastMessages = append(m.broadcastMessages, MockBroadcastMessage{
		RoomID:  roomID,
		Message: message,
	})
}

func (m *MockWSManager) SendToPlayer(playerID string, message *websocket.WSMessage) error {
	m.playerMessages = append(m.playerMessages, MockPlayerMessage{
		PlayerID: playerID,
		Message:  message,
	})
	return nil
}

func (m *MockWSManager) GetBroadcastCount() int {
	return len(m.broadcastMessages)
}

func (m *MockWSManager) GetPlayerMessageCount() int {
	return len(m.playerMessages)
}

func (m *MockWSManager) GetLastBroadcast() *MockBroadcastMessage {
	if len(m.broadcastMessages) == 0 {
		return nil
	}
	return &m.broadcastMessages[len(m.broadcastMessages)-1]
}

func (m *MockWSManager) GetLastPlayerMessage() *MockPlayerMessage {
	if len(m.playerMessages) == 0 {
		return nil
	}
	return &m.playerMessages[len(m.playerMessages)-1]
}

func TestNewGameService(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	
	if service == nil {
		t.Fatal("NewGameService returned nil")
	}
	
	if service.engines == nil {
		t.Error("engines map not initialized")
	}
	
	// Note: We can't directly compare interfaces, so we'll skip this check
	
	// Cleanup
	service.Stop()
}

func TestStartGame(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test valid game start
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
	
	// Verify engine was created
	service.mu.RLock()
	engine, exists := service.engines["room1"]
	service.mu.RUnlock()
	
	if !exists {
		t.Error("Game engine not created for room1")
	}
	
	if engine == nil {
		t.Error("Game engine is nil")
	}
	
	// Verify events were broadcast
	if mockWS.GetBroadcastCount() == 0 {
		t.Error("No events were broadcast")
	}
	
	// Test invalid inputs
	err = service.StartGame("", players)
	if err == nil {
		t.Error("Expected error for empty room ID")
	}
	
	err = service.StartGame("room2", []sdk.Player{})
	if err == nil {
		t.Error("Expected error for empty players")
	}
	
	err = service.StartGame("room2", players[:3])
	if err == nil {
		t.Error("Expected error for insufficient players")
	}
	
	// Test duplicate room
	err = service.StartGame("room1", players)
	if err == nil {
		t.Error("Expected error for duplicate room")
	}
}

func TestPlayCards(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game first
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
	
	// Test invalid room
	err = service.PlayCards("nonexistent", 0, []string{"Heart_5"})
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
	
	// Test invalid card IDs
	err = service.PlayCards("room1", 0, []string{"invalid_card"})
	if err == nil {
		t.Error("Expected error for invalid card ID")
	}
	
	// Test empty cards
	err = service.PlayCards("room1", 0, []string{})
	if err == nil {
		t.Error("Expected error for empty cards")
	}
}

func TestPassTurn(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game first
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
	
	// Test invalid room
	err = service.PassTurn("nonexistent", 0)
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestSubmitTributeSelection(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test invalid room
	err := service.SubmitTributeSelection("nonexistent", 0, "Heart_5")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestSubmitReturnTribute(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test invalid room
	err := service.SubmitReturnTribute("nonexistent", 0, "Heart_5")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestHandlePlayerDisconnect(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test invalid room
	err := service.HandlePlayerDisconnect("nonexistent", 0)
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestHandlePlayerReconnect(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test invalid room
	err := service.HandlePlayerReconnect("nonexistent", 0)
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestGetPlayerView(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game first
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
	
	// Test valid player view
	playerView, err := service.GetPlayerView("room1", 0)
	if err != nil {
		t.Errorf("GetPlayerView failed: %v", err)
	}
	
	if playerView == nil {
		t.Error("Player view is nil")
	}
	
	if playerView.PlayerSeat != 0 {
		t.Errorf("Expected player seat 0, got %d", playerView.PlayerSeat)
	}
	
	// Test invalid room
	_, err = service.GetPlayerView("nonexistent", 0)
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestGetGameState(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game first
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
	
	// Test valid game state
	gameState, err := service.GetGameState("room1")
	if err != nil {
		t.Errorf("GetGameState failed: %v", err)
	}
	
	if gameState == nil {
		t.Error("Game state is nil")
	}
	
	if gameState.Status != sdk.GameStatusStarted {
		t.Errorf("Expected game status started, got %s", gameState.Status)
	}
	
	// Test invalid room
	_, err = service.GetGameState("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestEndGame(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game first
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
	
	// End the game
	err = service.EndGame("room1")
	if err != nil {
		t.Errorf("EndGame failed: %v", err)
	}
	
	// Verify engine was removed
	service.mu.RLock()
	_, exists := service.engines["room1"]
	service.mu.RUnlock()
	
	if exists {
		t.Error("Game engine still exists after EndGame")
	}
	
	// Test ending nonexistent game
	err = service.EndGame("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent room")
	}
}

func TestParseCardFromID(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test valid card IDs
	testCases := []struct {
		cardID   string
		expected struct {
			number int
			color  string
		}
	}{
		{"Heart_5", struct{ number int; color string }{5, "Heart"}},
		{"Spade_13", struct{ number int; color string }{13, "Spade"}},
		{"Club_1", struct{ number int; color string }{14, "Club"}}, // 1 gets converted to 14 (Ace)
		{"Diamond_14", struct{ number int; color string }{14, "Diamond"}},
		{"Joker_15", struct{ number int; color string }{15, "Joker"}},
		{"Joker_16", struct{ number int; color string }{16, "Joker"}},
	}
	
	for _, tc := range testCases {
		card, err := service.parseCardFromID(tc.cardID)
		if err != nil {
			t.Errorf("parseCardFromID(%s) failed: %v", tc.cardID, err)
			continue
		}
		
		if card.Number != tc.expected.number {
			t.Errorf("parseCardFromID(%s): expected number %d, got %d", 
				tc.cardID, tc.expected.number, card.Number)
		}
		
		if card.Color != tc.expected.color {
			t.Errorf("parseCardFromID(%s): expected color %s, got %s", 
				tc.cardID, tc.expected.color, card.Color)
		}
	}
	
	// Test invalid card IDs
	invalidCases := []string{
		"",
		"Heart",
		"Heart_",
		"_5",
		"Heart_17",
		"Heart_0",
		"InvalidColor_5",
		"Heart_invalid",
	}
	
	for _, invalidID := range invalidCases {
		_, err := service.parseCardFromID(invalidID)
		if err == nil {
			t.Errorf("parseCardFromID(%s): expected error but got none", invalidID)
		}
	}
}

func TestConvertCardIDs(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test valid conversion
	cardIDs := []string{"Heart_5", "Spade_13", "Joker_15"}
	cards, err := service.convertCardIDs(cardIDs)
	if err != nil {
		t.Fatalf("convertCardIDs failed: %v", err)
	}
	
	if len(cards) != 3 {
		t.Errorf("Expected 3 cards, got %d", len(cards))
	}
	
	// Test empty input
	_, err = service.convertCardIDs([]string{})
	if err == nil {
		t.Error("Expected error for empty card IDs")
	}
	
	// Test invalid card ID
	_, err = service.convertCardIDs([]string{"invalid_card"})
	if err == nil {
		t.Error("Expected error for invalid card ID")
	}
}

func TestEventHandlerRegistration(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Start a game to trigger event handler registration
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
	
	// Wait a bit for events to be processed
	time.Sleep(100 * time.Millisecond)
	
	// Verify that events were broadcast
	if mockWS.GetBroadcastCount() == 0 {
		t.Error("No events were broadcast after starting game")
	}
	
	// Check that match started event was sent
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

func TestShouldSendPlayerViews(t *testing.T) {
	mockWS := NewMockWSManager()
	service := NewGameService(mockWS)
	defer service.Stop()
	
	// Test events that should trigger player views
	shouldSendEvents := []sdk.GameEventType{
		sdk.EventCardsDealt,
		sdk.EventTributeCompleted,
		sdk.EventPlayerPlayed,
		sdk.EventPlayerPassed,
		sdk.EventTrickEnded,
		sdk.EventDealEnded,
	}
	
	for _, eventType := range shouldSendEvents {
		if !service.shouldSendPlayerViews(eventType) {
			t.Errorf("shouldSendPlayerViews(%s) should return true", eventType)
		}
	}
	
	// Test events that should not trigger player views
	shouldNotSendEvents := []sdk.GameEventType{
		sdk.EventMatchStarted,
		sdk.EventTributePhase,
		sdk.EventTrickStarted,
	}
	
	for _, eventType := range shouldNotSendEvents {
		if service.shouldSendPlayerViews(eventType) {
			t.Errorf("shouldSendPlayerViews(%s) should return false", eventType)
		}
	}
}