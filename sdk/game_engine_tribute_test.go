package sdk

import (
	"testing"
	"time"
)

// TestProcessTributePhase_NoTribute tests when there's no tribute phase (first deal)
func TestProcessTributePhase_NoTribute(t *testing.T) {
	engine := NewGameEngine()
	
	// Start a match
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
	
	// First deal should not have tribute phase, so it should return error or nil action
	action, err := engine.ProcessTributePhase()
	
	// Either no error with nil action, or error saying not in tribute phase
	if err == nil && action != nil {
		t.Errorf("ProcessTributePhase should return nil action for first deal")
	}
	// It's OK if it returns an error saying not in tribute phase
	// since first deal doesn't have tribute phase
}

// TestProcessTributePhase_Immunity tests tribute immunity conditions
func TestProcessTributePhase_Immunity(t *testing.T) {
	engine := NewGameEngine()
	
	// Set up a match with a previous deal result
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
	
	// Simulate a previous deal result with Double Down victory
	engine.currentMatch.DealHistory = []*Deal{{
		ID:       "deal1",
		Rankings: []int{0, 2, 1, 3}, // Team 0 (seats 0,2) won with double down
	}}
	
	// Create a new deal with tribute phase
	lastResult := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		VictoryType: VictoryTypeDoubleDown,
	}
	
	deal, err := NewDeal(2, lastResult)
	if err != nil {
		t.Fatalf("Failed to create deal: %v", err)
	}
	
	// Give rank 3 and rank 4 two big jokers for immunity
	deal.PlayerCards[1] = []*Card{ // rank 3 (seat 1)
		{Number: 16, Color: "Joker", Name: "Red Joker"}, // Big Joker
	}
	deal.PlayerCards[3] = []*Card{ // rank 4 (seat 3)
		{Number: 16, Color: "Joker", Name: "Red Joker"}, // Big Joker
	}
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Check immunity
	tm := NewTributeManager(2)
	isImmune := tm.CheckTributeImmunity(lastResult, deal.PlayerCards)
	if !isImmune {
		t.Errorf("Should be immune when rank 3 and rank 4 have 2 big jokers total")
	}
}

// TestProcessTributePhase_DoubleDownSelection tests double down selection scenario
func TestProcessTributePhase_DoubleDownSelection(t *testing.T) {
	engine := NewGameEngine()
	
	// Set up match
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
	
	// Create deal with double down tribute
	lastResult := &DealResult{
		Rankings:    []int{0, 2, 1, 3}, // rank1=0, rank2=2, rank3=1, rank4=3
		VictoryType: VictoryTypeDoubleDown,
	}
	
	deal, err := NewDeal(2, lastResult)
	if err != nil {
		t.Fatalf("Failed to create deal: %v", err)
	}
	
	// Set up player hands BEFORE starting deal
	// Need to set all 4 players' hands to avoid issues
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	
	deal.PlayerCards[0] = []*Card{ // rank 1
		{Number: 10, Color: "Spade", Name: "10"},
	}
	deal.PlayerCards[1] = []*Card{ // rank 3
		{Number: 14, Color: "Spade", Name: "Ace"},
		{Number: 13, Color: "Heart", Name: "King"},
		{Number: 2, Color: "Diamond", Name: "2"},
	}
	deal.PlayerCards[2] = []*Card{ // rank 2
		{Number: 9, Color: "Heart", Name: "9"},
	}
	deal.PlayerCards[3] = []*Card{ // rank 4
		{Number: 14, Color: "Diamond", Name: "Ace"},
		{Number: 12, Color: "Spade", Name: "Queen"},
		{Number: 3, Color: "Club", Name: "3"},
	}
	
	engine.currentMatch.CurrentDeal = deal
	
	// Don't call StartDeal as it will override our hand setup
	// Instead, manually set the status
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status so it can be properly initialized
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process tribute phase - should return selection action for rank 1
	action, err := engine.ProcessTributePhase()
	if err != nil {
		t.Fatalf("ProcessTributePhase failed: %v", err)
	}
	
	// Debug output
	if deal.TributePhase != nil {
		t.Logf("Tribute phase status: %v", deal.TributePhase.Status)
		t.Logf("Pool cards: %d", len(deal.TributePhase.PoolCards))
		t.Logf("Selecting player: %d", deal.TributePhase.SelectingPlayer)
	}
	
	if action == nil {
		t.Fatal("Expected tribute action for double down selection")
	}
	
	if action.Type != TributeActionSelect {
		t.Errorf("Expected select action, got %v", action.Type)
	}
	
	if action.PlayerID != 0 { // rank 1
		t.Errorf("Expected player 0 (rank 1) to select, got %v", action.PlayerID)
	}
	
	if len(action.Options) != 2 {
		t.Errorf("Expected 2 cards in pool, got %d", len(action.Options))
	}
}

// TestSubmitTributeSelection tests submitting tribute selection
func TestSubmitTributeSelection(t *testing.T) {
	engine := NewGameEngine()
	
	// Set up match with double down tribute scenario
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	engine.StartMatch(players)
	
	lastResult := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		VictoryType: VictoryTypeDoubleDown,
	}
	
	deal, _ := NewDeal(2, lastResult)
	
	// Set up all player hands
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Spade", Name: "10"}} // rank1
	deal.PlayerCards[1] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}}  // rank3
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Heart", Name: "9"}}    // rank2
	deal.PlayerCards[3] = []*Card{{Number: 14, Color: "Diamond", Name: "Ace"}} // rank4
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process to get selection state
	action, _ := engine.ProcessTributePhase()
	if action == nil || len(deal.TributePhase.PoolCards) == 0 {
		t.Fatal("No pool cards available for selection")
	}
	
	// Submit selection
	selectedCard := deal.TributePhase.PoolCards[0]
	err := engine.SubmitTributeSelection(0, selectedCard.GetID())
	if err != nil {
		t.Errorf("SubmitTributeSelection failed: %v", err)
	}
	
	// Verify selection was recorded
	if deal.TributePhase.TributeCards[0] != selectedCard {
		t.Errorf("Selection was not recorded correctly")
	}
}

// TestSubmitReturnTribute tests submitting return tribute
func TestSubmitReturnTribute(t *testing.T) {
	engine := NewGameEngine()
	
	// Set up match with single last tribute scenario
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	engine.StartMatch(players)
	
	lastResult := &DealResult{
		Rankings:    []int{0, 1, 2, 3}, // rank1=0, rank4=3
		VictoryType: VictoryTypeSingleLast,
	}
	
	deal, _ := NewDeal(2, lastResult)
	
	// Set up all player hands
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	deal.PlayerCards[0] = []*Card{ // rank 1
		{Number: 14, Color: "Spade", Name: "Ace"},
		{Number: 2, Color: "Heart", Name: "2"},
		{Number: 3, Color: "Diamond", Name: "3"},
	}
	deal.PlayerCards[1] = []*Card{ // rank 2
		{Number: 10, Color: "Heart", Name: "10"},
	}
	deal.PlayerCards[2] = []*Card{ // rank 3
		{Number: 9, Color: "Club", Name: "9"},
	}
	deal.PlayerCards[3] = []*Card{ // rank 4
		{Number: 13, Color: "Heart", Name: "King"},
		{Number: 12, Color: "Diamond", Name: "Queen"},
	}
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process tribute to get to return phase
	action, _ := engine.ProcessTributePhase()
	
	// Should be return tribute action
	if action == nil || action.Type != TributeActionReturn {
		t.Fatal("Expected return tribute action")
	}
	
	// Submit return tribute
	returnCard := deal.PlayerCards[0][2] // Return the 3 of Diamond
	err := engine.SubmitReturnTribute(0, returnCard.GetID())
	if err != nil {
		t.Errorf("SubmitReturnTribute failed: %v", err)
	}
}

// TestSkipTributeAction tests timeout handling
func TestSkipTributeAction(t *testing.T) {
	engine := NewGameEngine()
	
	// Set up match with double down selection
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	engine.StartMatch(players)
	
	lastResult := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		VictoryType: VictoryTypeDoubleDown,
	}
	
	deal, _ := NewDeal(2, lastResult)
	
	// Set up all player hands
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Heart", Name: "10"}} // rank1
	deal.PlayerCards[1] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}} // rank3
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Diamond", Name: "9"}}  // rank2
	deal.PlayerCards[3] = []*Card{{Number: 13, Color: "Diamond", Name: "King"}} // rank4
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process to get selection state
	engine.ProcessTributePhase()
	
	// Skip (timeout) the selection
	err := engine.SkipTributeAction()
	if err != nil {
		t.Errorf("SkipTributeAction failed: %v", err)
	}
	
	// Verify auto-selection happened (should select highest card)
	if deal.TributePhase.TributeCards[0] == nil {
		t.Errorf("Auto-selection did not happen")
	}
}

// TestGetTributeStatus tests getting tribute status
func TestGetTributeStatus(t *testing.T) {
	engine := NewGameEngine()
	
	// When no deal is active
	status := engine.GetTributeStatus()
	if status != nil {
		t.Errorf("GetTributeStatus should return nil when no deal is active")
	}
	
	// Set up match with tribute
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	engine.StartMatch(players)
	
	lastResult := &DealResult{
		Rankings:    []int{0, 1, 2, 3},
		VictoryType: VictoryTypeSingleLast,
	}
	
	deal, _ := NewDeal(2, lastResult)
	
	// Set up all player hands
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Heart", Name: "10"}} // rank1
	deal.PlayerCards[3] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}} // rank4
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process once to initialize
	engine.ProcessTributePhase()
	
	// Get status during tribute phase
	status = engine.GetTributeStatus()
	if status == nil {
		t.Fatal("GetTributeStatus should return status during tribute phase")
	}
	
	if status.Phase != deal.TributePhase.Status {
		t.Errorf("Status phase mismatch: expected %v, got %v", deal.TributePhase.Status, status.Phase)
	}
}

// TestTributeEventHandling tests that tribute events are properly sent
func TestTributeEventHandling(t *testing.T) {
	engine := NewGameEngine()
	
	// Track events
	var receivedEvents []*GameEvent
	engine.RegisterEventHandler(EventTributeCompleted, func(event *GameEvent) {
		receivedEvents = append(receivedEvents, event)
	})
	engine.RegisterEventHandler(EventTributeSelected, func(event *GameEvent) {
		receivedEvents = append(receivedEvents, event)
	})
	engine.RegisterEventHandler(EventReturnTribute, func(event *GameEvent) {
		receivedEvents = append(receivedEvents, event)
	})
	
	// Set up match
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	engine.StartMatch(players)
	
	// Create simple tribute scenario
	lastResult := &DealResult{
		Rankings:    []int{0, 1, 2, 3},
		VictoryType: VictoryTypeSingleLast,
	}
	
	deal, _ := NewDeal(2, lastResult)
	
	// Set up all player hands
	for i := 0; i < 4; i++ {
		deal.PlayerCards[i] = []*Card{}
	}
	deal.PlayerCards[0] = []*Card{ // rank1
		{Number: 2, Color: "Heart", Name: "2"},
		{Number: 3, Color: "Diamond", Name: "3"},
	}
	deal.PlayerCards[1] = []*Card{{Number: 10, Color: "Club", Name: "10"}}   // rank2
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Spade", Name: "9"}}     // rank3
	deal.PlayerCards[3] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}} // rank4
	
	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute
	
	// Reset tribute phase to waiting status
	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}
	
	// Process tribute completely
	engine.ProcessTributePhase() // Initial processing
	
	// Submit return tribute
	returnCard := deal.PlayerCards[0][1]
	engine.SubmitReturnTribute(0, returnCard.GetID())
	
	// Process again to complete
	engine.ProcessTributePhase()
	
	// Give time for events to be processed
	time.Sleep(100 * time.Millisecond)
	
	// Check events were received
	if len(receivedEvents) < 2 {
		t.Errorf("Expected at least 2 events, got %d", len(receivedEvents))
	}
}