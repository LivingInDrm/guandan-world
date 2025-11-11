package sdk

import (
	"testing"
)

func TestDefaultTimeoutStrategy_GetDefaultPlayDecision(t *testing.T) {
	strategy := NewDefaultTimeoutStrategy()

	// 创建测试用牌
	card1, _ := NewCard(2, "Spade", 2)   // 最小的牌
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(10, "Club", 2)

	tests := []struct {
		name      string
		hand      []*Card
		trickInfo *TrickInfo
		wantAction ActionType
		wantCardsLen int
	}{
		{
			name: "Leader with cards - should play smallest card",
			hand: []*Card{card3, card2, card1},
			trickInfo: &TrickInfo{
				IsLeader: true,
			},
			wantAction: ActionPlay,
			wantCardsLen: 1,
		},
		{
			name: "Non-leader - should pass",
			hand: []*Card{card1, card2, card3},
			trickInfo: &TrickInfo{
				IsLeader: false,
			},
			wantAction: ActionPass,
			wantCardsLen: 0,
		},
		{
			name: "Leader with empty hand - should pass",
			hand: []*Card{},
			trickInfo: &TrickInfo{
				IsLeader: true,
			},
			wantAction: ActionPass,
			wantCardsLen: 0,
		},
		{
			name: "Nil trickInfo - should pass",
			hand: []*Card{card1},
			trickInfo: nil,
			wantAction: ActionPass,
			wantCardsLen: 0,
		},
		{
			name: "Leader with nil cards in hand - should play smallest non-nil card",
			hand: []*Card{nil, card3, nil, card1, card2},
			trickInfo: &TrickInfo{
				IsLeader: true,
			},
			wantAction: ActionPlay,
			wantCardsLen: 1,
		},
		{
			name: "Leader with all nil cards - should pass",
			hand: []*Card{nil, nil, nil},
			trickInfo: &TrickInfo{
				IsLeader: true,
			},
			wantAction: ActionPass,
			wantCardsLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := strategy.GetDefaultPlayDecision(tt.hand, tt.trickInfo)
			
			if decision == nil {
				t.Fatal("GetDefaultPlayDecision returned nil")
			}

			if decision.Action != tt.wantAction {
				t.Errorf("Action = %v, want %v", decision.Action, tt.wantAction)
			}

			if len(decision.Cards) != tt.wantCardsLen {
				t.Errorf("Cards length = %d, want %d", len(decision.Cards), tt.wantCardsLen)
			}

			// 如果应该出牌，验证出的是最小的牌
			if tt.wantAction == ActionPlay && len(decision.Cards) > 0 {
				playedCard := decision.Cards[0]
				if playedCard == nil {
					t.Error("Played card is nil")
				}
				// 验证这是最小的牌
				for _, card := range tt.hand {
					if card != nil && card.LessThan(playedCard) {
						t.Errorf("Played card is not the smallest: played %v but %v is smaller", 
							playedCard.GetID(), card.GetID())
					}
				}
			}
		})
	}
}

func TestDefaultTimeoutStrategy_GetDefaultTributeCard(t *testing.T) {
	strategy := NewDefaultTimeoutStrategy()

	card1, _ := NewCard(2, "Spade", 2)   // 小牌
	card2, _ := NewCard(10, "Heart", 2)  // 中等牌
	card3, _ := NewCard(14, "Club", 2)   // 大牌 (Ace)

	tests := []struct {
		name    string
		options []*Card
		wantNil bool
		wantLargest bool
	}{
		{
			name:    "Normal options - should return largest",
			options: []*Card{card1, card2, card3},
			wantNil: false,
			wantLargest: true,
		},
		{
			name:    "Empty options - should return nil",
			options: []*Card{},
			wantNil: true,
		},
		{
			name:    "Options with nil - should return largest non-nil",
			options: []*Card{nil, card2, nil, card3, card1},
			wantNil: false,
			wantLargest: true,
		},
		{
			name:    "All nil options - should return nil",
			options: []*Card{nil, nil, nil},
			wantNil: true,
		},
		{
			name:    "Single card - should return that card",
			options: []*Card{card2},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.GetDefaultTributeCard(tt.options)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result.GetID())
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// 如果需要验证是最大的牌
			if tt.wantLargest {
				for _, card := range tt.options {
					if card != nil && card.GreaterThan(result) {
						t.Errorf("Result is not the largest: got %v but %v is larger", 
							result.GetID(), card.GetID())
					}
				}
			}
		})
	}
}

func TestDefaultTimeoutStrategy_GetDefaultReturnCard(t *testing.T) {
	strategy := NewDefaultTimeoutStrategy()

	card1, _ := NewCard(2, "Spade", 2)   // 最小
	card2, _ := NewCard(10, "Heart", 2)
	card3, _ := NewCard(14, "Club", 2)   // 最大

	tests := []struct {
		name    string
		hand    []*Card
		wantNil bool
		wantSmallest bool
	}{
		{
			name:    "Normal hand - should return smallest",
			hand:    []*Card{card3, card2, card1},
			wantNil: false,
			wantSmallest: true,
		},
		{
			name:    "Empty hand - should return nil",
			hand:    []*Card{},
			wantNil: true,
		},
		{
			name:    "Hand with nil - should return smallest non-nil",
			hand:    []*Card{nil, card3, nil, card1, card2},
			wantNil: false,
			wantSmallest: true,
		},
		{
			name:    "All nil hand - should return nil",
			hand:    []*Card{nil, nil, nil},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.GetDefaultReturnCard(tt.hand)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result.GetID())
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// 如果需要验证是最小的牌
			if tt.wantSmallest {
				for _, card := range tt.hand {
					if card != nil && card.LessThan(result) {
						t.Errorf("Result is not the smallest: got %v but %v is smaller", 
							result.GetID(), card.GetID())
					}
				}
			}
		})
	}
}

func TestDefaultTimeoutStrategy_findSmallestCard(t *testing.T) {
	strategy := NewDefaultTimeoutStrategy()

	card1, _ := NewCard(2, "Spade", 2)
	card2, _ := NewCard(5, "Heart", 2)
	card3, _ := NewCard(10, "Club", 2)

	tests := []struct {
		name    string
		hand    []*Card
		wantNil bool
	}{
		{
			name:    "Normal hand",
			hand:    []*Card{card3, card2, card1},
			wantNil: false,
		},
		{
			name:    "Empty hand",
			hand:    []*Card{},
			wantNil: true,
		},
		{
			name:    "Hand with nils",
			hand:    []*Card{nil, card3, nil, card1},
			wantNil: false,
		},
		{
			name:    "All nils",
			hand:    []*Card{nil, nil},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.findSmallestCard(tt.hand)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result.GetID())
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// 验证确实是最小的
			for _, card := range tt.hand {
				if card != nil && card.LessThan(result) {
					t.Errorf("Result is not the smallest: got %v but %v is smaller", 
						result.GetID(), card.GetID())
				}
			}
		})
	}
}
