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
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Spade", Name: "10"}}    // rank1
	deal.PlayerCards[1] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}}   // rank3
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Heart", Name: "9"}}      // rank2
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

	// Verify selection was recorded in SelectionResults
	if deal.TributePhase.SelectionResults == nil {
		t.Errorf("SelectionResults was not initialized")
	} else if originalGiver, exists := deal.TributePhase.SelectionResults[0]; !exists {
		t.Errorf("Selection was not recorded in SelectionResults")
	} else {
		// Verify the selection points to the correct original giver
		if deal.TributePhase.TributeCards[originalGiver] != selectedCard {
			t.Errorf("Selection does not point to correct original giver's card")
		}
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
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Heart", Name: "10"}}     // rank1
	deal.PlayerCards[1] = []*Card{{Number: 14, Color: "Spade", Name: "Ace"}}    // rank3
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Diamond", Name: "9"}}     // rank2
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
	if deal.TributePhase.SelectionResults == nil || deal.TributePhase.SelectionResults[0] == 0 {
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
	deal.PlayerCards[0] = []*Card{{Number: 10, Color: "Heart", Name: "10"}}  // rank1
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
	deal.PlayerCards[2] = []*Card{{Number: 9, Color: "Spade", Name: "9"}}    // rank3
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

// TestCompleteDoubleDownFlow tests the complete double down tribute flow
func TestCompleteDoubleDownFlow(t *testing.T) {
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

	// Set up player hands with specific cards
	deal.PlayerCards[0] = []*Card{ // rank1 - will select first
		{Number: 2, Color: "Heart", Name: "2"},
		{Number: 3, Color: "Heart", Name: "3"},
	}
	deal.PlayerCards[1] = []*Card{ // rank3 - contributes to pool
		{Number: 14, Color: "Spade", Name: "Ace"},
		{Number: 13, Color: "Heart", Name: "King"},
	}
	deal.PlayerCards[2] = []*Card{ // rank2 - gets remaining card
		{Number: 4, Color: "Heart", Name: "4"},
		{Number: 5, Color: "Heart", Name: "5"},
	}
	deal.PlayerCards[3] = []*Card{ // rank4 - contributes to pool
		{Number: 16, Color: "Joker", Name: "BJ"}, // Big Joker
		{Number: 12, Color: "Spade", Name: "Queen"},
	}

	engine.currentMatch.CurrentDeal = deal
	deal.Status = DealStatusTribute

	if deal.TributePhase != nil {
		deal.TributePhase.Status = TributeStatusWaiting
	}

	// Step 1: Process tribute phase to create pool
	action, err := engine.ProcessTributePhase()
	if err != nil {
		t.Fatalf("ProcessTributePhase failed: %v", err)
	}

	if action == nil || action.Type != TributeActionSelect {
		t.Fatal("Expected tribute selection action")
	}

	if action.PlayerID != 0 {
		t.Errorf("Expected player 0 to select first, got %d", action.PlayerID)
	}

	// Verify pool contains tribute cards from rank3 and rank4
	if len(deal.TributePhase.PoolCards) != 2 {
		t.Fatalf("Expected 2 cards in pool, got %d", len(deal.TributePhase.PoolCards))
	}

	// Step 2: Player 0 selects Big Joker
	var selectedCardID string
	for _, card := range deal.TributePhase.PoolCards {
		if card.Number == 16 && card.Color == "Joker" {
			selectedCardID = card.GetID()
			break
		}
	}

	if selectedCardID == "" {
		t.Fatal("Big Joker not found in pool")
	}

	err = engine.SubmitTributeSelection(0, selectedCardID)
	if err != nil {
		t.Fatalf("SubmitTributeSelection failed: %v", err)
	}

	// Step 3: Process tribute phase again for player 2's selection
	action, err = engine.ProcessTributePhase()
	if err != nil {
		t.Fatalf("ProcessTributePhase failed for second selection: %v", err)
	}

	if action == nil || action.Type != TributeActionSelect {
		t.Fatal("Expected second tribute selection action")
	}

	if action.PlayerID != 2 {
		t.Errorf("Expected player 2 to select second, got %d", action.PlayerID)
	}

	// Step 4: Player 2 gets the remaining card (Ace)
	var remainingCardID string
	for _, card := range deal.TributePhase.PoolCards {
		remainingCardID = card.GetID()
		break // Should only be one card left
	}

	err = engine.SubmitTributeSelection(2, remainingCardID)
	if err != nil {
		t.Fatalf("SubmitTributeSelection failed for second selection: %v", err)
	}

	// Step 5: Verify return tribute relationships were established
	if len(deal.TributePhase.TributeMap) != 2 {
		t.Errorf("Expected 2 return tribute relationships, got %d", len(deal.TributePhase.TributeMap))
	}

	// Player 3 should return to player 0 (who got player 3's Big Joker)
	if receiver, exists := deal.TributePhase.TributeMap[3]; !exists || receiver != 0 {
		t.Errorf("Expected player 3 to return tribute to player 0, got %v", receiver)
	}

	// Player 1 should return to player 2 (who got player 1's Ace)
	if receiver, exists := deal.TributePhase.TributeMap[1]; !exists || receiver != 2 {
		t.Errorf("Expected player 1 to return tribute to player 2, got %v", receiver)
	}

	// Step 6: Process return tribute phase
	action, err = engine.ProcessTributePhase()
	if err != nil {
		t.Fatalf("ProcessTributePhase failed for return phase: %v", err)
	}

	if action == nil || action.Type != TributeActionReturn {
		t.Fatal("Expected return tribute action")
	}

	// Step 7: Handle all return tribute actions
	maxReturnActions := 4 // Safety counter
	returnActionCount := 0

	for action != nil && action.Type == TributeActionReturn && returnActionCount < maxReturnActions {
		t.Logf("Processing return tribute action for player %d", action.PlayerID)

		if len(deal.PlayerCards[action.PlayerID]) == 0 {
			t.Fatalf("Player %d has no cards to return", action.PlayerID)
		}

		returnCardID := deal.PlayerCards[action.PlayerID][0].GetID() // Return first card
		err = engine.SubmitReturnTribute(action.PlayerID, returnCardID)
		if err != nil {
			t.Fatalf("SubmitReturnTribute failed for player %d: %v", action.PlayerID, err)
		}

		returnActionCount++

		// Process next action
		action, err = engine.ProcessTributePhase()
		if err != nil {
			t.Fatalf("ProcessTributePhase failed after return #%d: %v", returnActionCount, err)
		}
	}

	if returnActionCount >= maxReturnActions {
		t.Fatalf("Too many return actions processed")
	}

	// Step 8: Verify tribute phase is now complete
	if action != nil {
		t.Errorf("Expected no more tribute actions, got %v for player %d", action.Type, action.PlayerID)
	}

	// Debug: Print tribute phase status
	t.Logf("Tribute phase status: %v", deal.TributePhase.Status)
	t.Logf("Selection results: %v", deal.TributePhase.SelectionResults)
	t.Logf("Tribute map: %v", deal.TributePhase.TributeMap)
	t.Logf("Return cards: %v", deal.TributePhase.ReturnCards)

	// Step 9: Verify hand changes
	// Player 0 should have gained Big Joker and lost a return card
	foundBigJoker := false
	for _, card := range deal.PlayerCards[0] {
		if card.Number == 16 && card.Color == "Joker" {
			foundBigJoker = true
			break
		}
	}
	if !foundBigJoker {
		t.Error("Player 0 should have Big Joker after tribute")
	}

	// Player 2 should have gained Ace
	foundAce := false
	for _, card := range deal.PlayerCards[2] {
		if card.Number == 14 && card.Color == "Spade" {
			foundAce = true
			break
		}
	}
	if !foundAce {
		t.Error("Player 2 should have Ace after tribute")
	}

	// Player 3 should have lost Big Joker
	stillHasBigJoker := false
	for _, card := range deal.PlayerCards[3] {
		if card.Number == 16 && card.Color == "Joker" {
			stillHasBigJoker = true
			break
		}
	}
	if stillHasBigJoker {
		t.Error("Player 3 should have lost Big Joker after tribute")
	}

	// Player 1 should have lost Ace
	stillHasAce := false
	for _, card := range deal.PlayerCards[1] {
		if card.Number == 14 && card.Color == "Spade" {
			stillHasAce = true
			break
		}
	}
	if stillHasAce {
		t.Error("Player 1 should have lost Ace after tribute")
	}

	t.Logf("✅ Complete double down tribute flow verified successfully")
}
