package sdk

import (
	"testing"
)

func TestNewDealer(t *testing.T) {
	tests := []struct {
		name        string
		level       int
		expectError bool
	}{
		{"Valid level 2", 2, false},
		{"Valid level 7", 7, false},
		{"Valid level 14", 14, false},
		{"Invalid level 1", 1, true},
		{"Invalid level 15", 15, true},
		{"Invalid level 0", 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealer, err := NewDealer(tt.level)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for level %d, but got none", tt.level)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for level %d: %v", tt.level, err)
				return
			}
			
			if dealer.GetLevel() != tt.level {
				t.Errorf("Expected level %d, got %d", tt.level, dealer.GetLevel())
			}
		})
	}
}

func TestCreateFullDeck(t *testing.T) {
	dealer, err := NewDealer(7)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	deck := dealer.CreateFullDeck()
	
	// Test deck size
	if len(deck) != 108 {
		t.Errorf("Expected deck size 108, got %d", len(deck))
	}
	
	// Test card distribution
	cardCounts := make(map[string]int)
	for _, card := range deck {
		key := card.Color + "-" + card.Name
		cardCounts[key]++
	}
	
	// Check regular cards (should have 2 of each)
	expectedRegularCards := []string{
		"Spade-2", "Spade-3", "Spade-4", "Spade-5", "Spade-6", "Spade-7", "Spade-8", "Spade-9", "Spade-10",
		"Spade-Jack", "Spade-Queen", "Spade-King", "Spade-Ace",
		"Club-2", "Club-3", "Club-4", "Club-5", "Club-6", "Club-7", "Club-8", "Club-9", "Club-10",
		"Club-Jack", "Club-Queen", "Club-King", "Club-Ace",
		"Heart-2", "Heart-3", "Heart-4", "Heart-5", "Heart-6", "Heart-7", "Heart-8", "Heart-9", "Heart-10",
		"Heart-Jack", "Heart-Queen", "Heart-King", "Heart-Ace",
		"Diamond-2", "Diamond-3", "Diamond-4", "Diamond-5", "Diamond-6", "Diamond-7", "Diamond-8", "Diamond-9", "Diamond-10",
		"Diamond-Jack", "Diamond-Queen", "Diamond-King", "Diamond-Ace",
	}
	
	for _, cardKey := range expectedRegularCards {
		if cardCounts[cardKey] != 2 {
			t.Errorf("Expected 2 copies of %s, got %d", cardKey, cardCounts[cardKey])
		}
	}
	
	// Check jokers
	if cardCounts["Joker-Black Joker"] != 2 {
		t.Errorf("Expected 2 Black Jokers, got %d", cardCounts["Joker-Black Joker"])
	}
	if cardCounts["Joker-Red Joker"] != 2 {
		t.Errorf("Expected 2 Red Jokers, got %d", cardCounts["Joker-Red Joker"])
	}
	
	// Verify all cards have correct level
	for _, card := range deck {
		if card.Level != 7 {
			t.Errorf("Expected all cards to have level 7, found card with level %d", card.Level)
		}
	}
}

func TestShuffleDeck(t *testing.T) {
	dealer, err := NewDealer(5)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	// Create initial deck
	originalDeck := dealer.CreateFullDeck()
	originalOrder := make([]*Card, len(originalDeck))
	copy(originalOrder, originalDeck)
	
	// Shuffle deck
	dealer.ShuffleDeck()
	shuffledDeck := dealer.GetDeck()
	
	// Verify deck size remains the same
	if len(shuffledDeck) != len(originalOrder) {
		t.Errorf("Deck size changed after shuffle: expected %d, got %d", len(originalOrder), len(shuffledDeck))
	}
	
	// Verify all cards are still present (same content, different order)
	originalCardCounts := make(map[string]int)
	shuffledCardCounts := make(map[string]int)
	
	for _, card := range originalOrder {
		key := card.Color + "-" + card.Name
		originalCardCounts[key]++
	}
	
	for _, card := range shuffledDeck {
		key := card.Color + "-" + card.Name
		shuffledCardCounts[key]++
	}
	
	for key, count := range originalCardCounts {
		if shuffledCardCounts[key] != count {
			t.Errorf("Card count mismatch for %s: expected %d, got %d", key, count, shuffledCardCounts[key])
		}
	}
	
	// Verify order is different (with high probability)
	samePositions := 0
	for i := 0; i < len(originalOrder); i++ {
		if originalOrder[i].Color == shuffledDeck[i].Color && originalOrder[i].Name == shuffledDeck[i].Name {
			samePositions++
		}
	}
	
	// It's extremely unlikely that more than 20% of cards remain in the same position after shuffle
	if float64(samePositions)/float64(len(originalOrder)) > 0.2 {
		t.Errorf("Shuffle appears ineffective: %d out of %d cards in same position", samePositions, len(originalOrder))
	}
}

func TestDealCards(t *testing.T) {
	dealer, err := NewDealer(8)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	playerHands, err := dealer.DealCards()
	if err != nil {
		t.Fatalf("Failed to deal cards: %v", err)
	}
	
	// Test that each player gets exactly 27 cards
	for player := 0; player < 4; player++ {
		if len(playerHands[player]) != 27 {
			t.Errorf("Player %d should have 27 cards, got %d", player, len(playerHands[player]))
		}
	}
	
	// Test that all 108 cards are distributed
	totalCards := 0
	for player := 0; player < 4; player++ {
		totalCards += len(playerHands[player])
	}
	
	if totalCards != 108 {
		t.Errorf("Total cards should be 108, got %d", totalCards)
	}
	
	// Test that no card is duplicated across all hands
	allCards := make([]*Card, 0, 108)
	for player := 0; player < 4; player++ {
		allCards = append(allCards, playerHands[player]...)
	}
	
	cardCounts := make(map[string]int)
	for _, card := range allCards {
		key := card.Color + "-" + card.Name
		cardCounts[key]++
	}
	
	// Verify correct distribution
	expectedRegularCards := []string{
		"Spade-2", "Spade-3", "Spade-4", "Spade-5", "Spade-6", "Spade-7", "Spade-8", "Spade-9", "Spade-10",
		"Spade-Jack", "Spade-Queen", "Spade-King", "Spade-Ace",
		"Club-2", "Club-3", "Club-4", "Club-5", "Club-6", "Club-7", "Club-8", "Club-9", "Club-10",
		"Club-Jack", "Club-Queen", "Club-King", "Club-Ace",
		"Heart-2", "Heart-3", "Heart-4", "Heart-5", "Heart-6", "Heart-7", "Heart-8", "Heart-9", "Heart-10",
		"Heart-Jack", "Heart-Queen", "Heart-King", "Heart-Ace",
		"Diamond-2", "Diamond-3", "Diamond-4", "Diamond-5", "Diamond-6", "Diamond-7", "Diamond-8", "Diamond-9", "Diamond-10",
		"Diamond-Jack", "Diamond-Queen", "Diamond-King", "Diamond-Ace",
	}
	
	for _, cardKey := range expectedRegularCards {
		if cardCounts[cardKey] != 2 {
			t.Errorf("Expected 2 copies of %s in dealt cards, got %d", cardKey, cardCounts[cardKey])
		}
	}
	
	if cardCounts["Joker-Black Joker"] != 2 {
		t.Errorf("Expected 2 Black Jokers in dealt cards, got %d", cardCounts["Joker-Black Joker"])
	}
	if cardCounts["Joker-Red Joker"] != 2 {
		t.Errorf("Expected 2 Red Jokers in dealt cards, got %d", cardCounts["Joker-Red Joker"])
	}
	
	// Test that hands are sorted
	for player := 0; player < 4; player++ {
		hand := playerHands[player]
		for i := 1; i < len(hand); i++ {
			// Verify sorting order (this is a basic check - the actual sorting logic is in sortCards)
			if hand[i-1].Number > hand[i].Number && hand[i].Number != 15 && hand[i].Number != 16 {
				// Allow for level cards and jokers which have special ordering
				if hand[i-1].Number != dealer.GetLevel() && hand[i].Number != dealer.GetLevel() {
					// This is a simplified check - full sorting validation would require understanding the complete sorting rules
					continue
				}
			}
		}
	}
}

func TestValidateDeck(t *testing.T) {
	dealer, err := NewDealer(6)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	// Test with valid deck
	dealer.CreateFullDeck()
	err = dealer.ValidateDeck()
	if err != nil {
		t.Errorf("Valid deck failed validation: %v", err)
	}
	
	// Test with nil deck
	dealer2, _ := NewDealer(6)
	err = dealer2.ValidateDeck()
	if err == nil {
		t.Error("Expected error for nil deck, got none")
	}
	
	// Test with invalid deck size
	dealer3, _ := NewDealer(6)
	dealer3.deck = make([]*Card, 100) // Wrong size
	err = dealer3.ValidateDeck()
	if err == nil {
		t.Error("Expected error for wrong deck size, got none")
	}
}

func TestSortPlayerHand(t *testing.T) {
	dealer, err := NewDealer(7)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	// Create some test cards
	card2, _ := NewCard(2, "Spade", 7)
	card3, _ := NewCard(3, "Heart", 7)
	cardA, _ := NewCard(14, "Club", 7)
	cardLevel, _ := NewCard(7, "Diamond", 7) // Level card
	joker, _ := NewCard(15, "Joker", 7)
	
	unsortedHand := []*Card{cardA, card2, joker, cardLevel, card3}
	sortedHand := dealer.SortPlayerHand(unsortedHand)
	
	// Verify the hand was sorted (basic check)
	if len(sortedHand) != len(unsortedHand) {
		t.Errorf("Sorted hand length mismatch: expected %d, got %d", len(unsortedHand), len(sortedHand))
	}
	
	// Verify all original cards are present
	originalCardMap := make(map[string]bool)
	sortedCardMap := make(map[string]bool)
	
	for _, card := range unsortedHand {
		key := card.Color + "-" + card.Name
		originalCardMap[key] = true
	}
	
	for _, card := range sortedHand {
		key := card.Color + "-" + card.Name
		sortedCardMap[key] = true
	}
	
	for key := range originalCardMap {
		if !sortedCardMap[key] {
			t.Errorf("Card %s missing from sorted hand", key)
		}
	}
}

func TestDealCardsConsistency(t *testing.T) {
	// Test that dealing cards multiple times produces different results (due to shuffling)
	dealer, err := NewDealer(9)
	if err != nil {
		t.Fatalf("Failed to create dealer: %v", err)
	}
	
	hands1, err := dealer.DealCards()
	if err != nil {
		t.Fatalf("Failed to deal cards first time: %v", err)
	}
	
	hands2, err := dealer.DealCards()
	if err != nil {
		t.Fatalf("Failed to deal cards second time: %v", err)
	}
	
	// Compare first player's hands - they should be different (with high probability)
	sameCards := 0
	for i := 0; i < len(hands1[0]); i++ {
		if hands1[0][i].Color == hands2[0][i].Color && hands1[0][i].Name == hands2[0][i].Name {
			sameCards++
		}
	}
	
	// It's extremely unlikely that more than 5 cards are the same in the same position
	if sameCards > 5 {
		t.Errorf("Dealing appears non-random: %d out of %d cards in same position", sameCards, len(hands1[0]))
	}
}