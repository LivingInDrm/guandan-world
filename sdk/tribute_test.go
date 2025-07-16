package sdk

import (
	"testing"
	"time"
)

func TestNewTributeManager(t *testing.T) {
	tm := NewTributeManager(7)
	if tm == nil {
		t.Fatal("Expected tribute manager, got nil")
	}
	if tm.level != 7 {
		t.Errorf("Expected level 7, got %d", tm.level)
	}
}

func TestNewTributePhase(t *testing.T) {
	tests := []struct {
		name           string
		lastResult     *DealResult
		expectNil      bool
		expectError    bool
		expectedStatus TributeStatus
	}{
		{
			name:       "No last result",
			lastResult: nil,
			expectNil:  true,
		},
		{
			name: "Normal tribute scenario",
			lastResult: &DealResult{
				Rankings: []int{0, 1, 2, 3}, // Team 0 wins (seats 0,2), Team 1 loses (seats 1,3)
			},
			expectedStatus: TributeStatusReturning,
		},
		{
			name: "Double down scenario",
			lastResult: &DealResult{
				Rankings: []int{1, 3, 0, 2}, // Team 1: 1(1st), 3(2nd); Team 0: 0(3rd), 2(4th) - both 0,2 from Team 0 last
			},
			expectedStatus: TributeStatusSelecting,
		},
		{
			name: "Invalid rankings",
			lastResult: &DealResult{
				Rankings: []int{0, 1}, // Too few rankings
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp, err := NewTributePhase(tt.lastResult)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if tt.expectNil {
				if tp != nil {
					t.Error("Expected nil tribute phase, got non-nil")
				}
				return
			}
			
			if tp == nil {
				t.Fatal("Expected tribute phase, got nil")
			}
			
			if tp.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, tp.Status)
			}
		})
	}
}

func TestDetermineTributeRequirements(t *testing.T) {
	tm := NewTributeManager(7)
	
	tests := []struct {
		name           string
		lastResult     *DealResult
		expectNil      bool
		expectError    bool
		expectedMap    map[int]int
		isDoubleDown   bool
	}{
		{
			name:       "No last result",
			lastResult: nil,
			expectNil:  true,
		},
		{
			name: "Normal tribute",
			lastResult: &DealResult{
				Rankings: []int{0, 1, 2, 3}, // 0 first, 1 second, 2 third, 3 fourth
			},
			expectedMap: map[int]int{
				3: 0, // Fourth gives to first
				2: 1, // Third gives to second
			},
			isDoubleDown: false,
		},
		{
			name: "Double down - same team last",
			lastResult: &DealResult{
				Rankings: []int{1, 3, 0, 2}, // Team 1: 1,3; Team 0: 0,2 - both 0,2 from Team 0 last
			},
			expectedMap: map[int]int{
				0: -1, // Pool contributor (3rd place)
				2: -1, // Pool contributor (4th place)
			},
			isDoubleDown: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tributeMap, isDoubleDown, err := tm.DetermineTributeRequirements(tt.lastResult)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if tt.expectNil {
				if tributeMap != nil {
					t.Error("Expected nil tribute map, got non-nil")
				}
				return
			}
			
			if isDoubleDown != tt.isDoubleDown {
				t.Errorf("Expected isDoubleDown %v, got %v", tt.isDoubleDown, isDoubleDown)
			}
			
			if len(tributeMap) != len(tt.expectedMap) {
				t.Errorf("Expected tribute map length %d, got %d", len(tt.expectedMap), len(tributeMap))
			}
			
			for giver, receiver := range tt.expectedMap {
				if tributeMap[giver] != receiver {
					t.Errorf("Expected tribute map[%d] = %d, got %d", giver, receiver, tributeMap[giver])
				}
			}
		})
	}
}

func TestTributeManagerGetHighestCard(t *testing.T) {
	tm := NewTributeManager(7)
	
	// Create test cards
	card2, _ := NewCard(2, "Spade", 7)
	card5, _ := NewCard(5, "Heart", 7)
	cardA, _ := NewCard(14, "Club", 7)
	joker, _ := NewCard(15, "Joker", 7)
	
	tests := []struct {
		name     string
		hand     []*Card
		expected *Card
	}{
		{
			name:     "Empty hand",
			hand:     []*Card{},
			expected: nil,
		},
		{
			name:     "Single card",
			hand:     []*Card{card5},
			expected: card5,
		},
		{
			name:     "Multiple cards",
			hand:     []*Card{card2, cardA, card5},
			expected: cardA,
		},
		{
			name:     "With joker",
			hand:     []*Card{card2, cardA, joker},
			expected: joker,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.getHighestCard(tt.hand)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}
			
			if result == nil {
				t.Fatal("Expected card, got nil")
			}
			
			if !tm.cardsEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTributeManagerGetLowestCard(t *testing.T) {
	tm := NewTributeManager(7)
	
	// Create test cards
	card2, _ := NewCard(2, "Spade", 7)
	card5, _ := NewCard(5, "Heart", 7)
	cardA, _ := NewCard(14, "Club", 7)
	joker, _ := NewCard(15, "Joker", 7)
	
	tests := []struct {
		name     string
		hand     []*Card
		expected *Card
	}{
		{
			name:     "Empty hand",
			hand:     []*Card{},
			expected: nil,
		},
		{
			name:     "Single card",
			hand:     []*Card{card5},
			expected: card5,
		},
		{
			name:     "Multiple cards",
			hand:     []*Card{card2, cardA, card5},
			expected: card2,
		},
		{
			name:     "With joker",
			hand:     []*Card{card2, cardA, joker},
			expected: card2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.getLowestCard(tt.hand)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}
			
			if result == nil {
				t.Fatal("Expected card, got nil")
			}
			
			if !tm.cardsEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTributeManagerRemoveCardFromHand(t *testing.T) {
	tm := NewTributeManager(7)
	
	// Create test cards
	card2, _ := NewCard(2, "Spade", 7)
	card5, _ := NewCard(5, "Heart", 7)
	cardA, _ := NewCard(14, "Club", 7)
	
	tests := []struct {
		name           string
		hand           []*Card
		cardToRemove   *Card
		expectedLength int
		shouldContain  []*Card
	}{
		{
			name:           "Remove from middle",
			hand:           []*Card{card2, card5, cardA},
			cardToRemove:   card5,
			expectedLength: 2,
			shouldContain:  []*Card{card2, cardA},
		},
		{
			name:           "Remove first card",
			hand:           []*Card{card2, card5, cardA},
			cardToRemove:   card2,
			expectedLength: 2,
			shouldContain:  []*Card{card5, cardA},
		},
		{
			name:           "Remove last card",
			hand:           []*Card{card2, card5, cardA},
			cardToRemove:   cardA,
			expectedLength: 2,
			shouldContain:  []*Card{card2, card5},
		},
		{
			name:           "Card not in hand",
			hand:           []*Card{card2, card5},
			cardToRemove:   cardA,
			expectedLength: 2,
			shouldContain:  []*Card{card2, card5},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.removeCardFromHand(tt.hand, tt.cardToRemove)
			
			if len(result) != tt.expectedLength {
				t.Errorf("Expected length %d, got %d", tt.expectedLength, len(result))
			}
			
			// Check that expected cards are still present
			for _, expectedCard := range tt.shouldContain {
				found := false
				for _, resultCard := range result {
					if tm.cardsEqual(resultCard, expectedCard) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected card %v not found in result", expectedCard)
				}
			}
		})
	}
}

func TestProcessTribute(t *testing.T) {
	tm := NewTributeManager(7)
	
	// Create test hands
	var playerHands [4][]*Card
	for player := 0; player < 4; player++ {
		playerHands[player] = make([]*Card, 0)
		// Add some test cards
		for i := 2; i <= 5; i++ {
			card, _ := NewCard(i, "Spade", 7)
			playerHands[player] = append(playerHands[player], card)
		}
	}
	
	// Test normal tribute scenario
	lastResult := &DealResult{
		Rankings: []int{0, 1, 2, 3}, // Normal ranking
	}
	
	tributePhase, err := NewTributePhase(lastResult)
	if err != nil {
		t.Fatalf("Failed to create tribute phase: %v", err)
	}
	
	// Process tribute
	err = tm.ProcessTribute(tributePhase, playerHands)
	if err != nil {
		t.Errorf("Failed to process tribute: %v", err)
	}
	
	// Verify tribute phase progressed
	if tributePhase.Status != TributeStatusFinished {
		t.Errorf("Expected status %s, got %s", TributeStatusFinished, tributePhase.Status)
	}
	
	// Verify tribute cards were selected
	if len(tributePhase.TributeCards) == 0 {
		t.Error("Expected tribute cards to be selected")
	}
}

func TestApplyTributeToHands(t *testing.T) {
	tm := NewTributeManager(7)
	
	// Create test hands
	var playerHands [4][]*Card
	
	// Player 0 (receiver): gets tribute
	card2_0, _ := NewCard(2, "Spade", 7)
	card3_0, _ := NewCard(3, "Spade", 7)
	playerHands[0] = []*Card{card2_0, card3_0}
	
	// Player 1 (receiver): gets tribute
	card2_1, _ := NewCard(2, "Heart", 7)
	card3_1, _ := NewCard(3, "Heart", 7)
	playerHands[1] = []*Card{card2_1, card3_1}
	
	// Player 2 (giver): gives tribute
	card4_2, _ := NewCard(4, "Club", 7)
	card5_2, _ := NewCard(5, "Club", 7)
	playerHands[2] = []*Card{card4_2, card5_2}
	
	// Player 3 (giver): gives tribute
	card4_3, _ := NewCard(4, "Diamond", 7)
	card5_3, _ := NewCard(5, "Diamond", 7)
	playerHands[3] = []*Card{card4_3, card5_3}
	
	// Create tribute phase
	tributePhase := &TributePhase{
		Status:       TributeStatusFinished,
		TributeMap:   map[int]int{3: 0, 2: 1}, // 3->0, 2->1
		TributeCards: map[int]*Card{3: card5_3, 2: card5_2},
		ReturnCards:  map[int]*Card{0: card2_0, 1: card2_1},
	}
	
	originalLen0 := len(playerHands[0])
	originalLen3 := len(playerHands[3])
	
	// Apply tribute
	err := tm.ApplyTributeToHands(tributePhase, &playerHands)
	if err != nil {
		t.Errorf("Failed to apply tribute: %v", err)
	}
	
	// Verify hand sizes remain the same
	if len(playerHands[0]) != originalLen0 {
		t.Errorf("Player 0 hand size changed: expected %d, got %d", originalLen0, len(playerHands[0]))
	}
	if len(playerHands[3]) != originalLen3 {
		t.Errorf("Player 3 hand size changed: expected %d, got %d", originalLen3, len(playerHands[3]))
	}
	
	// Verify tribute card was transferred
	found := false
	for _, card := range playerHands[0] {
		if tm.cardsEqual(card, card5_3) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tribute card not found in receiver's hand")
	}
}

func TestTributePhaseSelectTribute(t *testing.T) {
	// Create test cards for pool
	card5, _ := NewCard(5, "Spade", 7)
	cardA, _ := NewCard(14, "Heart", 7)
	
	tributePhase := &TributePhase{
		Status:          TributeStatusSelecting,
		SelectingPlayer: 0,
		PoolCards:       []*Card{card5, cardA},
		TributeCards:    make(map[int]*Card),
		SelectTimeout:   time.Now().Add(3 * time.Second),
	}
	
	// Test valid selection
	err := tributePhase.SelectTribute(0, card5)
	if err != nil {
		t.Errorf("Failed to select tribute: %v", err)
	}
	
	// Verify card was selected
	if tributePhase.TributeCards[0] != card5 {
		t.Error("Tribute card not recorded correctly")
	}
	
	// Verify card was removed from pool
	if len(tributePhase.PoolCards) != 1 {
		t.Errorf("Expected 1 card in pool, got %d", len(tributePhase.PoolCards))
	}
	
	// Test invalid player
	err = tributePhase.SelectTribute(1, cardA)
	if err == nil {
		t.Error("Expected error for wrong player, got none")
	}
}

func TestTributePhaseHandleTimeout(t *testing.T) {
	// Create test cards for pool
	card5, _ := NewCard(5, "Spade", 7)
	cardA, _ := NewCard(14, "Heart", 7)
	
	tributePhase := &TributePhase{
		Status:          TributeStatusSelecting,
		SelectingPlayer: 0,
		PoolCards:       []*Card{card5, cardA},
		TributeCards:    make(map[int]*Card),
		SelectTimeout:   time.Now().Add(-1 * time.Second), // Already expired
	}
	
	// Test timeout handling
	err := tributePhase.HandleTimeout()
	if err != nil {
		t.Errorf("Failed to handle timeout: %v", err)
	}
	
	// Verify highest card was selected
	selectedCard := tributePhase.TributeCards[0]
	if selectedCard == nil {
		t.Fatal("No card was selected on timeout")
	}
	
	// Ace should be higher than 5
	if !selectedCard.GreaterThan(card5) {
		t.Error("Expected highest card to be selected on timeout")
	}
}

func TestCardsEqual(t *testing.T) {
	tm := NewTributeManager(7)
	
	card1, _ := NewCard(5, "Spade", 7)
	card2, _ := NewCard(5, "Spade", 7)
	card3, _ := NewCard(5, "Heart", 7)
	card4, _ := NewCard(6, "Spade", 7)
	
	tests := []struct {
		name     string
		card1    *Card
		card2    *Card
		expected bool
	}{
		{"Same card", card1, card1, true},
		{"Equal cards", card1, card2, true},
		{"Different suit", card1, card3, false},
		{"Different number", card1, card4, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.cardsEqual(tt.card1, tt.card2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}