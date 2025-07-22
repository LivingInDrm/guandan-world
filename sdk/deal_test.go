package sdk

import (
	"testing"
)

func TestNewDeal(t *testing.T) {
	// Test creating deal without tribute
	deal, err := NewDeal(5, nil)
	if err != nil {
		t.Fatalf("NewDeal failed: %v", err)
	}

	if deal == nil {
		t.Fatal("NewDeal should return a non-nil deal")
	}

	if deal.ID == "" {
		t.Error("Deal should have a non-empty ID")
	}

	if deal.Level != 5 {
		t.Errorf("Deal level should be 5, got %d", deal.Level)
	}

	if deal.Status != DealStatusWaiting {
		t.Errorf("New deal should have status %v, got %v", DealStatusWaiting, deal.Status)
	}

	if deal.TributePhase != nil {
		t.Error("Deal without last result should not have tribute phase")
	}

	if len(deal.Rankings) != 0 {
		t.Error("New deal should have empty rankings")
	}

	if len(deal.TrickHistory) != 0 {
		t.Error("New deal should have empty trick history")
	}
}

func TestNewDealWithTribute(t *testing.T) {
	// Create a mock last result that would trigger tribute
	lastResult := &DealResult{
		Rankings:    []int{0, 1, 2, 3},
		WinningTeam: 0,
		VictoryType: VictoryTypeSingleLast,
		Upgrades:    [2]int{2, 0},
	}

	deal, err := NewDeal(6, lastResult)
	if err != nil {
		t.Fatalf("NewDeal with tribute failed: %v", err)
	}

	if deal.TributePhase == nil {
		t.Error("Deal with last result should have tribute phase")
	}

	if deal.Status != DealStatusTribute {
		t.Errorf("Deal with tribute should have status %v, got %v", DealStatusTribute, deal.Status)
	}
}

func TestNewDealValidation(t *testing.T) {
	// Test with invalid level
	_, err := NewDeal(1, nil)
	if err == nil {
		t.Error("NewDeal should fail with level < 2")
	}

	_, err = NewDeal(15, nil)
	if err == nil {
		t.Error("NewDeal should fail with level > 14")
	}
}

func TestDealCardDealing(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	// Test dealing cards
	err := deal.dealCards()
	if err != nil {
		t.Errorf("dealCards failed: %v", err)
	}

	// Check that each player has 27 cards
	for player := 0; player < 4; player++ {
		if len(deal.PlayerCards[player]) != 27 {
			t.Errorf("Player %d should have 27 cards, got %d", player, len(deal.PlayerCards[player]))
		}
	}

	// Check that all cards are valid
	totalCards := 0
	for player := 0; player < 4; player++ {
		for _, card := range deal.PlayerCards[player] {
			if card == nil {
				t.Errorf("Player %d has nil card", player)
			}
			if card.Level != 5 {
				t.Errorf("Card should have level 5, got %d", card.Level)
			}
		}
		totalCards += len(deal.PlayerCards[player])
	}

	if totalCards != 108 {
		t.Errorf("Total cards should be 108, got %d", totalCards)
	}
}

func TestDealCreateFullDeck(t *testing.T) {
	deal, _ := NewDeal(7, nil)

	deck := deal.createFullDeck()

	if len(deck) != 108 {
		t.Errorf("Full deck should have 108 cards, got %d", len(deck))
	}

	// Count cards by type
	regularCards := 0
	jokers := 0

	for _, card := range deck {
		if card.Color == "Joker" {
			jokers++
		} else {
			regularCards++
		}

		if card.Level != 7 {
			t.Errorf("All cards should have level 7, got %d", card.Level)
		}
	}

	if regularCards != 104 {
		t.Errorf("Should have 104 regular cards, got %d", regularCards)
	}

	if jokers != 4 {
		t.Errorf("Should have 4 jokers, got %d", jokers)
	}
}

func TestDealStartDeal(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	err := deal.StartDeal()
	if err != nil {
		t.Errorf("StartDeal failed: %v", err)
	}

	// Should have dealt cards
	for player := 0; player < 4; player++ {
		if len(deal.PlayerCards[player]) != 27 {
			t.Errorf("Player %d should have 27 cards after StartDeal", player)
		}
	}

	// Should have started first trick
	if deal.CurrentTrick == nil {
		t.Error("Deal should have current trick after StartDeal")
	}

	if deal.Status != DealStatusPlaying {
		t.Errorf("Deal should have status %v after StartDeal, got %v", DealStatusPlaying, deal.Status)
	}
}

func TestDealValidatePlayerCards(t *testing.T) {
	deal, _ := NewDeal(5, nil)
	deal.dealCards()

	// Test with valid cards from player's hand
	playerCards := deal.PlayerCards[0]
	testCards := playerCards[:2] // Take first 2 cards

	err := deal.validatePlayerCards(0, testCards)
	if err != nil {
		t.Errorf("validatePlayerCards should succeed with valid cards: %v", err)
	}

	// Test with cards not in player's hand - create a card that definitely doesn't exist
	fakeCard, _ := NewCard(2, "Heart", 99) // Level 99 doesn't exist in any deck
	err = deal.validatePlayerCards(0, []*Card{fakeCard})
	if err == nil {
		t.Error("validatePlayerCards should fail with cards not in player's hand")
	}

	// Test with invalid player seat
	err = deal.validatePlayerCards(-1, testCards)
	if err == nil {
		t.Error("validatePlayerCards should fail with invalid player seat")
	}

	err = deal.validatePlayerCards(4, testCards)
	if err == nil {
		t.Error("validatePlayerCards should fail with invalid player seat")
	}
}

func TestDealRemoveCardsFromPlayer(t *testing.T) {
	deal, _ := NewDeal(5, nil)
	deal.dealCards()

	player := 0
	originalCount := len(deal.PlayerCards[player])
	cardsToRemove := deal.PlayerCards[player][:3] // Remove first 3 cards

	deal.removeCardsFromPlayer(player, cardsToRemove)

	newCount := len(deal.PlayerCards[player])
	if newCount != originalCount-3 {
		t.Errorf("Player should have %d cards after removal, got %d", originalCount-3, newCount)
	}

	// Check that the correct number of cards was removed
	// (We can't check for exact card absence due to potential duplicates in deck)
}

func TestDealGetNextPlayer(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	if deal.getNextPlayer(0) != 1 {
		t.Error("Next player after 0 should be 1")
	}

	if deal.getNextPlayer(1) != 2 {
		t.Error("Next player after 1 should be 2")
	}

	if deal.getNextPlayer(2) != 3 {
		t.Error("Next player after 2 should be 3")
	}

	if deal.getNextPlayer(3) != 0 {
		t.Error("Next player after 3 should be 0")
	}
}

func TestDealIsDealFinished(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	// Initially not finished
	if deal.isDealFinished() {
		t.Error("New deal should not be finished")
	}

	// Add 2 players to rankings
	deal.Rankings = []int{0, 1}
	if deal.isDealFinished() {
		t.Error("Deal with 2 finished players should not be finished")
	}

	// Add 3rd player
	deal.Rankings = []int{0, 1, 2}
	if !deal.isDealFinished() {
		t.Error("Deal with 3 finished players should be finished")
	}
}

func TestDealCardsEqual(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	card1, _ := NewCard(10, "Heart", 5)
	card2, _ := NewCard(10, "Heart", 5)
	card3, _ := NewCard(10, "Spade", 5)
	card4, _ := NewCard(9, "Heart", 5)

	if !deal.cardsEqual(card1, card2) {
		t.Error("Same cards should be equal")
	}

	if deal.cardsEqual(card1, card3) {
		t.Error("Cards with different colors should not be equal")
	}

	if deal.cardsEqual(card1, card4) {
		t.Error("Cards with different numbers should not be equal")
	}
}

func TestDealProcessTimeouts(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	// Test with no active components
	events := deal.ProcessTimeouts()
	if events == nil {
		t.Error("ProcessTimeouts should return empty slice, not nil")
	}

	if len(events) != 0 {
		t.Error("ProcessTimeouts should return empty slice when no active components")
	}
}

func TestDealFinishDeal(t *testing.T) {
	deal, _ := NewDeal(5, nil)

	// Add some players to rankings
	deal.Rankings = []int{0, 1, 2}

	err := deal.finishDeal()
	if err != nil {
		t.Errorf("finishDeal failed: %v", err)
	}

	if deal.Status != DealStatusFinished {
		t.Errorf("Deal should have status %v after finishing, got %v", DealStatusFinished, deal.Status)
	}

	if len(deal.Rankings) != 4 {
		t.Error("All 4 players should be in rankings after finishing")
	}

	if deal.EndTime == nil {
		t.Error("Deal should have end time after finishing")
	}

	// Check that missing player was added
	found := false
	for _, seat := range deal.Rankings {
		if seat == 3 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Missing player should be added to rankings")
	}
}
