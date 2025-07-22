package sdk

import (
	"fmt"
	"testing"
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
			name: "Single Last tribute scenario",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 2, 3}, // Rank1=0, Rank2=1, Rank3=2, Rank4=3
				WinningTeam: 0,
				VictoryType: VictoryTypeSingleLast, // rank1(0), rank3(2) 同队
			},
			expectedStatus: TributeStatusReturning,
		},
		{
			name: "Double down scenario",
			lastResult: &DealResult{
				Rankings:    []int{1, 3, 0, 2}, // Rank1=1, Rank2=3, Rank3=0, Rank4=2
				WinningTeam: 1,
				VictoryType: VictoryTypeDoubleDown, // rank1(1), rank2(3) 同队
			},
			expectedStatus: TributeStatusSelecting,
		},
		{
			name: "Partner Last tribute scenario",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 3, 2}, // Rank1=0, Rank2=1, Rank3=3, Rank4=2
				WinningTeam: 0,
				VictoryType: VictoryTypePartnerLast, // rank1(0), rank4(2) 同队
			},
			expectedStatus: TributeStatusReturning,
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
	tm := NewTributeManager(5)

	tests := []struct {
		name         string
		lastResult   *DealResult
		expectedMap  map[int]int
		isDoubleDown bool
		expectError  bool
	}{
		{
			name:       "No last result",
			lastResult: nil,
		},
		{
			name: "Partner last tribute - rank1,rank4 same team",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 3, 2}, // Rank1=0, Rank2=1, Rank3=3, Rank4=2
				WinningTeam: 0,
				VictoryType: VictoryTypePartnerLast,
			},
			expectedMap: map[int]int{
				3: 0, // Rank3(3) -> Rank1(0)
			},
			isDoubleDown: false,
		},
		{
			name: "Single last tribute - rank1,rank3 same team",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 2, 3}, // Rank1=0, Rank2=1, Rank3=2, Rank4=3
				WinningTeam: 0,
				VictoryType: VictoryTypeSingleLast,
			},
			expectedMap: map[int]int{
				3: 0, // Rank4(3) -> Rank1(0)
			},
			isDoubleDown: false,
		},
		{
			name: "Double down - rank1,rank2 same team",
			lastResult: &DealResult{
				Rankings:    []int{1, 3, 0, 2}, // Rank1=1, Rank2=3, Rank3=0, Rank4=2
				WinningTeam: 1,
				VictoryType: VictoryTypeDoubleDown,
			},
			expectedMap: map[int]int{
				0: -1, // Rank3(0)贡献到池子
				2: -1, // Rank4(2)贡献到池子
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

			if tt.lastResult == nil {
				if tributeMap != nil || isDoubleDown {
					t.Error("Expected nil tribute map and false isDoubleDown for nil result")
				}
				return
			}

			if isDoubleDown != tt.isDoubleDown {
				t.Errorf("Expected isDoubleDown %v, got %v", tt.isDoubleDown, isDoubleDown)
			}

			if len(tributeMap) != len(tt.expectedMap) {
				t.Errorf("Expected tribute map length %d, got %d", len(tt.expectedMap), len(tributeMap))
				return
			}

			for giver, expectedReceiver := range tt.expectedMap {
				if receiver, exists := tributeMap[giver]; !exists {
					t.Errorf("Expected giver %d not found in tribute map", giver)
				} else if receiver != expectedReceiver {
					t.Errorf("For giver %d, expected receiver %d, got %d", giver, expectedReceiver, receiver)
				}
			}
		})
	}
}

func TestCheckTributeImmunity(t *testing.T) {
	tm := NewTributeManager(5)

	// Create test hands with and without Big Jokers
	handsWithBigJokers := [4][]*Card{
		{}, // Player 0: no cards
		{}, // Player 1: no cards
		{}, // Player 2: no cards
		{}, // Player 3: no cards
	}

	// Add 2 Big Jokers to player 2 and player 3
	bigJoker1, _ := NewCard(16, "Joker", 5) // Red Joker = Big Joker
	bigJoker2, _ := NewCard(16, "Joker", 5)
	bigJoker3, _ := NewCard(16, "Joker", 5)

	handsWithBigJokers[2] = []*Card{bigJoker1}            // Player 2: 1 Big Joker
	handsWithBigJokers[3] = []*Card{bigJoker2, bigJoker3} // Player 3: 2 Big Jokers

	handsWithoutBigJokers := [4][]*Card{
		{}, {}, {}, {},
	}

	tests := []struct {
		name           string
		lastResult     *DealResult
		playerHands    [4][]*Card
		expectedResult bool
	}{
		{
			name:           "No last result",
			lastResult:     nil,
			playerHands:    handsWithBigJokers,
			expectedResult: false,
		},
		{
			name: "Double Down immunity - Rank3,Rank4 have 2+ Big Jokers combined",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 2, 3}, // Rank3=2, Rank4=3
				VictoryType: VictoryTypeDoubleDown,
			},
			playerHands:    handsWithBigJokers, // Player 2 has 1, Player 3 has 2 = 3 total
			expectedResult: true,
		},
		{
			name: "Double Down no immunity - Rank3,Rank4 don't have enough Big Jokers",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 2, 3}, // Rank3=2, Rank4=3
				VictoryType: VictoryTypeDoubleDown,
			},
			playerHands:    handsWithoutBigJokers,
			expectedResult: false,
		},
		{
			name: "Single Last immunity - Rank4 has 2+ Big Jokers",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 2, 3}, // Rank4=3
				VictoryType: VictoryTypeSingleLast,
			},
			playerHands:    handsWithBigJokers, // Player 3 has 2 Big Jokers
			expectedResult: true,
		},
		{
			name: "Partner Last immunity - Rank3 has 2+ Big Jokers",
			lastResult: &DealResult{
				Rankings:    []int{0, 1, 3, 2}, // Rank3=3 (player 3)
				VictoryType: VictoryTypePartnerLast,
			},
			playerHands:    handsWithBigJokers, // Player 3 has 2 Big Jokers
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.CheckTributeImmunity(tt.lastResult, tt.playerHands)
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestGetHighestCardExcludingHeartTrump(t *testing.T) {
	tm := NewTributeManager(5) // Level 5

	// Create test cards
	heartTrump, _ := NewCard(5, "Heart", 5)      // Red Trump (level 5 Heart)
	aceSpade, _ := NewCard(14, "Spade", 5)       // Ace of Spades
	kingHeart, _ := NewCard(13, "Heart", 5)      // King of Hearts
	queenDiamond, _ := NewCard(12, "Diamond", 5) // Queen of Diamonds
	bigJoker, _ := NewCard(16, "Joker", 5)       // Big Joker

	tests := []struct {
		name        string
		hand        []*Card
		expectedNum int // Expected card number (-1 if no suitable card)
	}{
		{
			name:        "Empty hand",
			hand:        []*Card{},
			expectedNum: -1,
		},
		{
			name:        "Hand with only Heart Trump",
			hand:        []*Card{heartTrump},
			expectedNum: 5, // Should return the Heart Trump as fallback
		},
		{
			name:        "Hand with mixed cards excluding Heart Trump",
			hand:        []*Card{aceSpade, kingHeart, queenDiamond},
			expectedNum: 14, // Ace is highest
		},
		{
			name:        "Hand with Heart Trump and other cards",
			hand:        []*Card{heartTrump, aceSpade, kingHeart},
			expectedNum: 14, // Ace, excluding Heart Trump
		},
		{
			name:        "Hand with Big Joker",
			hand:        []*Card{heartTrump, aceSpade, bigJoker},
			expectedNum: 16, // Big Joker is highest
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.getHighestCardExcludingHeartTrump(tt.hand)

			if tt.expectedNum == -1 {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected card with number %d, got nil", tt.expectedNum)
				} else if result.Number != tt.expectedNum {
					t.Errorf("Expected card number %d, got %d", tt.expectedNum, result.Number)
				}
			}
		})
	}
}

func TestTributeProcessComplete(t *testing.T) {
	// 创建测试场景：Single Last胜利类型（rank1,rank3同队）
	lastResult := &DealResult{
		Rankings:    []int{0, 1, 2, 3}, // rank1=0, rank2=1, rank3=2, rank4=3
		WinningTeam: 0,
		VictoryType: VictoryTypeSingleLast, // rank1(0), rank3(2) 同队
	}

	// 创建测试手牌
	playerHands := [4][]*Card{
		// 玩家0 (rank1, 胜方): 有红桃Trump和其他牌
		{
			{Number: 7, Color: "Heart"},  // 红桃Trump (level=7)
			{Number: 14, Color: "Spade"}, // A♠
			{Number: 13, Color: "Club"},  // K♣
		},
		// 玩家1 (rank2, 败方): 普通牌
		{
			{Number: 12, Color: "Diamond"}, // Q♦
			{Number: 11, Color: "Spade"},   // J♠
			{Number: 10, Color: "Heart"},   // 10♥
		},
		// 玩家2 (rank3, 胜方): 普通牌
		{
			{Number: 9, Color: "Club"},    // 9♣
			{Number: 8, Color: "Diamond"}, // 8♦
			{Number: 7, Color: "Spade"},   // 7♠
		},
		// 玩家3 (rank4, 败方): 需要上贡，有大王和其他牌
		{
			{Number: 16, Color: "Joker"},   // 大王
			{Number: 14, Color: "Heart"},   // A♥
			{Number: 13, Color: "Diamond"}, // K♦
		},
	}

	// 创建上贡阶段
	tributePhase, err := NewTributePhase(lastResult)
	if err != nil {
		t.Fatalf("创建上贡阶段失败: %v", err)
	}

	// 创建TributeManager
	tm := NewTributeManager(7)

	// 检查免贡（应该不免贡，因为rank4只有1张大王）
	isImmune := tm.CheckTributeImmunity(lastResult, playerHands)
	if isImmune {
		t.Errorf("期望不免贡，但实际免贡了")
	}

	// 处理上贡 - 可能需要多次调用来完成整个流程
	for i := 0; i < 10; i++ { // 最多尝试10次避免无限循环
		err = tm.ProcessTribute(tributePhase, playerHands)
		if err != nil {
			t.Fatalf("处理上贡失败: %v", err)
		}
		if tributePhase.Status == TributeStatusFinished {
			break
		}
	}

	if tributePhase.Status != TributeStatusFinished {
		t.Fatalf("上贡阶段未完成，当前状态: %s", tributePhase.Status)
	}

	// 验证上贡映射
	expectedTributeMap := map[int]int{3: 0} // rank4(3) → rank1(0)
	if len(tributePhase.TributeMap) != len(expectedTributeMap) {
		t.Errorf("上贡映射数量不对，期望 %d，实际 %d", len(expectedTributeMap), len(tributePhase.TributeMap))
	}

	for giver, receiver := range expectedTributeMap {
		if actualReceiver, exists := tributePhase.TributeMap[giver]; !exists || actualReceiver != receiver {
			t.Errorf("上贡映射错误：期望 %d→%d，实际 %d→%d", giver, receiver, giver, actualReceiver)
		}
	}

	// 验证上贡牌选择
	if len(tributePhase.TributeCards) == 0 {
		t.Errorf("未选择上贡牌")
	}

	// 验证上贡牌是除红桃Trump外最大的牌
	tributeCard := tributePhase.TributeCards[3]
	if tributeCard == nil {
		t.Fatalf("玩家3的上贡牌为空")
	}

	// 上贡牌应该是大王 (除红桃Trump外最大的牌)
	expectedTributeCard := &Card{Number: 16, Color: "Joker"} // 大王
	if tributeCard.Number != expectedTributeCard.Number || tributeCard.Color != expectedTributeCard.Color {
		t.Errorf("上贡牌错误：期望 %s，实际 %s",
			formatCardForTest(expectedTributeCard), formatCardForTest(tributeCard))
	}

	// 应用上贡效果到手牌
	err = tm.ApplyTributeToHands(tributePhase, &playerHands)
	if err != nil {
		t.Fatalf("应用上贡效果失败: %v", err)
	}

	// 验证手牌变化
	// 玩家3应该失去大王
	for _, card := range playerHands[3] {
		if card.Number == 16 && card.Color == "Joker" {
			t.Errorf("玩家3手牌中仍有大王，上贡未生效")
		}
	}

	// 玩家0应该获得大王
	foundTributeCard := false
	for _, card := range playerHands[0] {
		if card.Number == 16 && card.Color == "Joker" {
			foundTributeCard = true
			break
		}
	}
	if !foundTributeCard {
		t.Errorf("玩家0手牌中没有大王，上贡未生效")
	}

	t.Logf("✅ 上贡过程验证成功")
	t.Logf("   上贡映射: %v", tributePhase.TributeMap)
	t.Logf("   上贡牌: 玩家3 → 玩家0, 牌: %s", formatCardForTest(tributeCard))
	t.Logf("   玩家3剩余手牌: %d张", len(playerHands[3]))
	t.Logf("   玩家0手牌: %d张", len(playerHands[0]))
}

func formatCardForTest(card *Card) string {
	if card == nil {
		return "nil"
	}
	if card.Color == "Joker" {
		if card.Number == 15 {
			return "小王"
		} else if card.Number == 16 {
			return "大王"
		}
	}
	suits := map[string]string{"Heart": "♥", "Diamond": "♦", "Club": "♣", "Spade": "♠"}
	numbers := map[int]string{11: "J", 12: "Q", 13: "K", 14: "A"}

	suit := suits[card.Color]
	if suit == "" {
		suit = card.Color
	}

	number := numbers[card.Number]
	if number == "" {
		number = fmt.Sprintf("%d", card.Number)
	}

	return number + suit
}
