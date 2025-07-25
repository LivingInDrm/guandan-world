package sdk

import (
	"testing"
)

// TestWildcardNormalization tests wildcard replacement functionality
func TestWildcardNormalization(t *testing.T) {
	// Test setup: level 3, so red heart 3 is wildcard
	level := 3

	t.Run("Pair with wildcard", func(t *testing.T) {
		// Create a pair with one wildcard (red heart 3) and one normal 5
		wildcard, _ := NewCard(3, "Heart", level) // This will be a wildcard
		normal5, _ := NewCard(5, "Spade", level)
		cards := []*Card{wildcard, normal5}

		// Debug: Check if first card is actually a wildcard
		if !cards[0].IsWildcard() {
			t.Errorf("First card should be wildcard. Number=%d, Color=%s, Level=%d", 
				cards[0].Number, cards[0].Color, cards[0].Level)
		}

		pair := NewPair(cards)
		if !pair.Valid {
			t.Error("Pair should be valid")
		}

		// Debug: Print original cards in pair
		t.Logf("Original pair cards:")
		for i, card := range pair.Cards {
			t.Logf("  Card %d: Number=%d, RawNumber=%d, Color=%s, IsWildcard=%v", 
				i, card.Number, card.RawNumber, card.Color, card.IsWildcard())
		}

		// Check normalized cards
		if pair.NormalizedCards == nil {
			t.Error("NormalizedCards should not be nil")
			return
		}

		// Debug: Print detailed card info
		for i, card := range pair.NormalizedCards {
			t.Logf("Normalized card %d: Number=%d, RawNumber=%d, Color=%s, Level=%d", 
				i, card.Number, card.RawNumber, card.Color, card.Level)
		}

		// Both cards should be 5s in normalized version
		for i, card := range pair.NormalizedCards {
			if card.Number != 5 {
				t.Errorf("Expected normalized card %d to be 5, got %d", i, card.Number)
			}
		}
	})

	t.Run("Triple with two wildcards", func(t *testing.T) {
		// Create a triple with two wildcards and one normal 7
		wildcard1, _ := NewCard(3, "Heart", level)
		wildcard2, _ := NewCard(3, "Heart", level) 
		normal7, _ := NewCard(7, "Club", level)
		cards := []*Card{wildcard1, wildcard2, normal7}

		triple := NewTriple(cards)
		if !triple.Valid {
			t.Error("Triple should be valid")
		}

		// Check normalized cards
		if triple.NormalizedCards == nil {
			t.Error("NormalizedCards should not be nil")
		}

		// All cards should be 7s in normalized version
		for _, card := range triple.NormalizedCards {
			if card.Number != 7 {
				t.Errorf("Expected normalized card to be 7, got %d", card.Number)
			}
		}
	})

	t.Run("FullHouse with wildcard", func(t *testing.T) {
		// Create a full house: 888+99 with one wildcard replacing a 9
		card8_1, _ := NewCard(8, "Spade", level)
		card8_2, _ := NewCard(8, "Club", level)
		card8_3, _ := NewCard(8, "Diamond", level)
		card9, _ := NewCard(9, "Spade", level)
		wildcard, _ := NewCard(3, "Heart", level)
		cards := []*Card{card8_1, card8_2, card8_3, card9, wildcard}

		fullHouse := NewFullHouse(cards)
		if !fullHouse.Valid {
			t.Error("FullHouse should be valid")
		}

		// Check normalized cards
		if fullHouse.NormalizedCards == nil {
			t.Error("NormalizedCards should not be nil")
		}

		// Should have three 8s and two 9s
		countMap := make(map[int]int)
		for _, card := range fullHouse.NormalizedCards {
			countMap[card.Number]++
		}

		if countMap[8] != 3 {
			t.Errorf("Expected 3 cards with value 8, got %d", countMap[8])
		}
		if countMap[9] != 2 {
			t.Errorf("Expected 2 cards with value 9, got %d", countMap[9])
		}
	})

	t.Run("Straight with wildcard", func(t *testing.T) {
		// Create a straight: 5,6,7,8,wildcard (as 9)
		card5, _ := NewCard(5, "Spade", level)
		card6, _ := NewCard(6, "Club", level)
		card7, _ := NewCard(7, "Diamond", level)
		card8, _ := NewCard(8, "Spade", level)
		wildcard, _ := NewCard(3, "Heart", level)
		cards := []*Card{card5, card6, card7, card8, wildcard}

		straight := NewStraight(cards)
		if !straight.Valid {
			t.Error("Straight should be valid")
		}

		// Check normalized cards
		if straight.NormalizedCards == nil {
			t.Error("NormalizedCards should not be nil")
		}

		// Should have 5,6,7,8,9 in order
		expectedValues := []int{5, 6, 7, 8, 9}
		for i, card := range straight.NormalizedCards {
			if card.Number != expectedValues[i] {
				t.Errorf("Expected card at position %d to be %d, got %d", 
					i, expectedValues[i], card.Number)
			}
		}
	})

	t.Run("NaiveBomb with wildcard", func(t *testing.T) {
		// Create a bomb: KKKK with one wildcard
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Club", level)
		cardK3, _ := NewCard(13, "Diamond", level)
		wildcard, _ := NewCard(3, "Heart", level)
		cards := []*Card{cardK1, cardK2, cardK3, wildcard}

		bomb := NewNaiveBomb(cards)
		if !bomb.Valid {
			t.Error("Bomb should be valid")
		}

		// Check normalized cards
		if bomb.NormalizedCards == nil {
			t.Error("NormalizedCards should not be nil")
		}

		// All cards should be Kings
		for _, card := range bomb.NormalizedCards {
			if card.Number != 13 {
				t.Errorf("Expected normalized card to be K (13), got %d", card.Number)
			}
		}
	})

	t.Run("Comparison with wildcards", func(t *testing.T) {
		// Create two pairs: one with wildcard (3,5) and one without (4,4)
		wildcard, _ := NewCard(3, "Heart", level)
		card5, _ := NewCard(5, "Spade", level)
		pairWithWildcard := NewPair([]*Card{wildcard, card5})

		card4_1, _ := NewCard(4, "Spade", level)
		card4_2, _ := NewCard(4, "Club", level)
		pairWithoutWildcard := NewPair([]*Card{card4_1, card4_2})

		// Pair with wildcard (normalized to 5,5) should be greater than (4,4)
		if !pairWithWildcard.GreaterThan(pairWithoutWildcard) {
			t.Error("Pair (5,5) should be greater than (4,4)")
		}

		// (4,4) should not be greater than (5,5)
		if pairWithoutWildcard.GreaterThan(pairWithWildcard) {
			t.Error("Pair (4,4) should not be greater than (5,5)")
		}
	})
}