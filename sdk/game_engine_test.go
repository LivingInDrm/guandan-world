package sdk

import (
	"testing"
	"time"
)

func TestNewGameEngine(t *testing.T) {
	engine := NewGameEngine()

	if engine == nil {
		t.Fatal("NewGameEngine should return a non-nil engine")
	}

	if engine.id == "" {
		t.Error("GameEngine should have a non-empty ID")
	}

	if engine.status != GameStatusWaiting {
		t.Errorf("New GameEngine should have status %v, got %v", GameStatusWaiting, engine.status)
	}

	if engine.observers == nil {
		t.Error("GameEngine should have initialized observers map")
	}

	if engine.createdAt.IsZero() {
		t.Error("GameEngine should have a creation timestamp")
	}
}

func TestGameEngineEventSystem(t *testing.T) {
	engine := NewGameEngine()

	// Test event handler registration with channel for synchronization
	eventChan := make(chan *GameEvent, 1)

	handler := func(event *GameEvent) {
		eventChan <- event
	}

	engine.On(EventMatchStarted, handler)

	// Check that handler was registered
	if len(engine.observers[EventMatchStarted]) != 1 {
		t.Error("Event handler should be registered")
	}

	// Test event emission
	testEvent := &GameEvent{
		Type:      EventMatchStarted,
		Data:      "test data",
		Timestamp: time.Now(),
	}

	engine.emitEvent(testEvent)

	// Wait for event with timeout
	select {
	case receivedEvent := <-eventChan:
		if receivedEvent == nil {
			t.Error("Event handler should have received the event")
		}
		if receivedEvent.Type != EventMatchStarted {
			t.Errorf("Expected event type %v, got %v", EventMatchStarted, receivedEvent.Type)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Event handler should have been called within timeout")
	}
}

func TestGameEngineMultipleEventHandlers(t *testing.T) {
	engine := NewGameEngine()

	// Register multiple handlers for the same event with channel for synchronization
	callChan := make(chan struct{}, 2)

	handler1 := func(event *GameEvent) {
		callChan <- struct{}{}
	}

	handler2 := func(event *GameEvent) {
		callChan <- struct{}{}
	}

	engine.On(EventPlayerPlayed, handler1)
	engine.On(EventPlayerPlayed, handler2)

	// Emit event
	testEvent := &GameEvent{
		Type:      EventPlayerPlayed,
		Data:      nil,
		Timestamp: time.Now(),
	}

	engine.emitEvent(testEvent)

	// Wait for both handlers to be called
	callCount := 0
	timeout := time.After(500 * time.Millisecond)
	for callCount < 2 {
		select {
		case <-callChan:
			callCount++
		case <-timeout:
			t.Fatalf("Expected 2 handler calls, got %d within timeout", callCount)
		}
	}
}

func TestGameEngineGetGameState(t *testing.T) {
	engine := NewGameEngine()

	state := engine.GetGameState()

	if state == nil {
		t.Fatal("GetGameState should return a non-nil state")
	}

	if state.ID != engine.id {
		t.Error("Game state ID should match engine ID")
	}

	if state.Status != engine.status {
		t.Error("Game state status should match engine status")
	}

	if state.CurrentMatch != engine.currentMatch {
		t.Error("Game state current match should match engine current match")
	}
}

func TestGameEngineGetPlayerView(t *testing.T) {
	engine := NewGameEngine()

	playerView := engine.GetPlayerView(0)

	if playerView != nil {
		t.Error("GetPlayerView should return nil when no match is active")
	}

	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	playerView = engine.GetPlayerView(0)
	if playerView == nil {
		t.Fatal("GetPlayerView should return a non-nil view when match and deal are active")
	}

	if playerView.PlayerSeat != 0 {
		t.Error("Player view should have correct player seat")
	}

	if len(playerView.TeamLevels) != 2 {
		t.Error("Player view should have team levels")
	}

	if playerView.PlayerCards == nil {
		t.Error("Player view should have initialized player cards")
	}
	
	if len(playerView.PlayerCards) != 27 {
		t.Errorf("Player should have 27 cards, got %d", len(playerView.PlayerCards))
	}
}

func TestGameEngineIsGameFinished(t *testing.T) {
	engine := NewGameEngine()

	// Initially should not be finished
	if engine.IsGameFinished() {
		t.Error("New game should not be finished")
	}

	// Set status to finished
	engine.status = GameStatusFinished

	if !engine.IsGameFinished() {
		t.Error("Game with finished status should be finished")
	}
}

func TestGameEngineStartMatchValidation(t *testing.T) {
	engine := NewGameEngine()

	// Test with wrong number of players
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
	}

	err := engine.StartMatch(players)
	if err == nil {
		t.Error("StartMatch should fail with less than 4 players")
	}

	// Test with correct number of players
	players = append(players, Player{ID: "4", Username: "Player4", Seat: 3})

	err = engine.StartMatch(players)
	// This should now succeed since NewMatch is implemented
	if err != nil {
		t.Errorf("StartMatch should succeed with 4 valid players: %v", err)
	}

	// Verify match was created
	if engine.currentMatch == nil {
		t.Error("StartMatch should create a current match")
	}

	if engine.status != GameStatusStarted {
		t.Error("Engine status should be started after successful StartMatch")
	}
}

func TestGameEngineErrorHandling(t *testing.T) {
	engine := NewGameEngine()

	// Test operations without active match
	_, err := engine.PlayCards(0, []*Card{})
	if err == nil {
		t.Error("PlayCards should fail without active match")
	}

	_, err = engine.PassTurn(0)
	if err == nil {
		t.Error("PassTurn should fail without active match")
	}

	_, err = engine.HandlePlayerDisconnect(0)
	if err == nil {
		t.Error("HandlePlayerDisconnect should fail without active match")
	}

	_, err = engine.HandlePlayerReconnect(0)
	if err == nil {
		t.Error("HandlePlayerReconnect should fail without active match")
	}

	err = engine.SetPlayerAutoPlay(0, true)
	if err == nil {
		t.Error("SetPlayerAutoPlay should fail without active match")
	}
}

// TestGameEngineProcessTimeouts removed - timeout handling moved to GameDriver (Task 5)

func TestGameEngineThreadSafety(t *testing.T) {
	engine := NewGameEngine()

	// Test concurrent access to GetGameState
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			state := engine.GetGameState()
			if state == nil {
				t.Error("GetGameState should not return nil")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test concurrent event handler registration
	for i := 0; i < 10; i++ {
		go func(index int) {
			handler := func(event *GameEvent) {
				// Do nothing
			}
			engine.On(EventPlayerPlayed, handler)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Check that all handlers were registered
	if len(engine.observers[EventPlayerPlayed]) != 10 {
		t.Errorf("Expected 10 event handlers, got %d", len(engine.observers[EventPlayerPlayed]))
	}
}

func TestGameEventTypes(t *testing.T) {
	// Test that all event types are defined
	eventTypes := []GameEventType{
		EventMatchStarted,
		EventDealStarted,
		EventCardsDealt,
		EventTributePhase,
		EventTributeRulesSet,
		EventTributeImmunity,
		EventTributePoolCreated,
		EventTributeStarted,
		EventTributeGiven,
		EventTributeSelected,
		EventReturnTribute,
		EventTributeCompleted,
		EventTrickStarted,
		EventPlayerPlayed,
		EventPlayerPassed,
		EventTrickEnded,
		EventDealEnded,
		EventMatchEnded,
		EventPlayerTimeout,
		EventPlayerDisconnect,
		EventPlayerReconnect,
	}

	for _, eventType := range eventTypes {
		if string(eventType) == "" {
			t.Errorf("Event type %v should not be empty", eventType)
		}
	}
}

func TestGameStatusTypes(t *testing.T) {
	// Test that all game status types are defined
	statusTypes := []GameStatus{
		GameStatusWaiting,
		GameStatusStarted,
		GameStatusFinished,
	}

	for _, status := range statusTypes {
		if string(status) == "" {
			t.Errorf("Game status %v should not be empty", status)
		}
	}
}
