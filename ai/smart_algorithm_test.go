package ai

import (
	"guandan-world/sdk"
	"testing"
)

// 创建测试用的牌
func createCard(number int, color string) *sdk.Card {
	card, _ := sdk.NewCard(number, color, 2) // 默认级别为2
	return card
}

func TestSmartAutoPlayAlgorithm_SelectCardsToPlay(t *testing.T) {
	algo := NewSmartAutoPlayAlgorithm(2) // 设置当前级别为2

	tests := []struct {
		name     string
		hand     []*sdk.Card
		trickInfo *sdk.TrickInfo
		expected int // 期望出牌数量
	}{
		{
			name: "首出选择最优牌组",
			hand: []*sdk.Card{
				createCard(3, "Spade"), 
				createCard(3, "Heart"),
				createCard(3, "Diamond"),
				createCard(4, "Spade"),
				createCard(4, "Heart"),
				createCard(5, "Spade"),
				createCard(5, "Heart"),
				createCard(6, "Spade"),
				createCard(6, "Heart"),
				createCard(7, "Spade"),
			},
			trickInfo: &sdk.TrickInfo{IsLeader: true},
			expected: 6, // 应该出钢板（连续三对）
		},
		{
			name: "手中有顺子时优先出顺子",
			hand: []*sdk.Card{
				createCard(3, "Heart"), 
				createCard(4, "Diamond"),
				createCard(5, "Club"),
				createCard(6, "Spade"),
				createCard(7, "Heart"),
				createCard(9, "Heart"),
				createCard(10, "Heart"),
			},
			trickInfo: &sdk.TrickInfo{IsLeader: true},
			expected: 5, // 应该出顺子
		},
		{
			name: "跟牌时选择最小能管上的牌",
			hand: []*sdk.Card{
				createCard(5, "Spade"),
				createCard(5, "Heart"),
				createCard(8, "Spade"),
				createCard(8, "Heart"),
				createCard(10, "Spade"),
				createCard(10, "Heart"),
			},
			trickInfo: &sdk.TrickInfo{
				IsLeader: false,
				LeadComp: sdk.FromCardList([]*sdk.Card{
					createCard(4, "Spade"),
					createCard(4, "Heart"),
				}, nil),
			},
			expected: 2, // 应该出5对
		},
		{
			name: "无法跟牌时过牌",
			hand: []*sdk.Card{
				createCard(3, "Spade"),
				createCard(4, "Heart"),
				createCard(5, "Diamond"),
			},
			trickInfo: &sdk.TrickInfo{
				IsLeader: false,
				LeadComp: sdk.FromCardList([]*sdk.Card{
					createCard(10, "Spade"),
					createCard(10, "Heart"),
				}, nil),
			},
			expected: 0, // 应该过牌
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := algo.SelectCardsToPlay(tt.hand, tt.trickInfo)
			if result == nil && tt.expected == 0 {
				return // 正确过牌
			}
			if result == nil && tt.expected > 0 {
				t.Errorf("Expected %d cards, but got nil", tt.expected)
				return
			}
			if len(result) != tt.expected {
				t.Errorf("Expected %d cards, but got %d cards: %v", tt.expected, len(result), result)
			}
		})
	}
}

func TestSmartAutoPlayAlgorithm_CardGroupIdentification(t *testing.T) {
	algo := &SmartAutoPlayAlgorithm{level: 2}

	tests := []struct {
		name           string
		hand           []*sdk.Card
		expectedGroups []string // 期望识别出的牌型
	}{
		{
			name: "识别多种牌型",
			hand: []*sdk.Card{
				// 三张3
				createCard(3, "Spade"),
				createCard(3, "Heart"),
				createCard(3, "Diamond"),
				// 对4
				createCard(4, "Spade"),
				createCard(4, "Heart"),
				// 顺子 5-9
				createCard(5, "Spade"),
				createCard(6, "Heart"),
				createCard(7, "Diamond"),
				createCard(8, "Club"),
				createCard(9, "Spade"),
			},
			expectedGroups: []string{"Triple", "Pair", "Straight", "FullHouse"},
		},
		{
			name: "识别炸弹",
			hand: []*sdk.Card{
				createCard(5, "Spade"),
				createCard(5, "Heart"),
				createCard(5, "Diamond"),
				createCard(5, "Club"),
				createCard(6, "Spade"),
			},
			expectedGroups: []string{"NaiveBomb", "Pair", "Triple"},
		},
		{
			name: "识别钢板（连续三对）",
			hand: []*sdk.Card{
				createCard(5, "Spade"),
				createCard(5, "Heart"),
				createCard(6, "Diamond"),
				createCard(6, "Club"),
				createCard(7, "Spade"),
				createCard(7, "Heart"),
			},
			expectedGroups: []string{"Tube", "Pair"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := algo.identifyAllPossibleGroups(tt.hand)
			
			// 检查是否识别出期望的牌型
			groupTypes := make(map[string]bool)
			for _, group := range groups {
				groupTypes[group.CompType.String()] = true
			}
			
			for _, expected := range tt.expectedGroups {
				if !groupTypes[expected] {
					t.Errorf("Expected to identify %s, but didn't", expected)
				}
			}
			
			t.Logf("Identified %d groups:", len(groups))
			for _, group := range groups {
				t.Logf("  - %s with %d cards", group.CompType.String(), len(group.Cards))
			}
		})
	}
}

func TestSmartAutoPlayAlgorithm_DamageCalculation(t *testing.T) {
	algo := &SmartAutoPlayAlgorithm{level: 2}

	// 创建一个包含多个重叠牌组的手牌
	hand := []*sdk.Card{
		// 四张5（炸弹）
		createCard(5, "Spade"),
		createCard(5, "Heart"),
		createCard(5, "Diamond"),
		createCard(5, "Club"),
		// 三张6
		createCard(6, "Spade"),
		createCard(6, "Heart"),
		createCard(6, "Diamond"),
		// 对7
		createCard(7, "Spade"),
		createCard(7, "Heart"),
	}

	groups := algo.identifyAllPossibleGroups(hand)
	algo.calculateDamageForGroups(groups, hand)

	// 找出炸弹牌组
	var bombGroup *CardGroup
	for _, group := range groups {
		if group.CompType == sdk.TypeNaiveBomb && len(group.Cards) == 4 {
			bombGroup = group
			break
		}
	}

	// 找出使用了5的对子
	var pairGroup *CardGroup
	for _, group := range groups {
		if group.CompType == sdk.TypePair && group.Cards[0].Number == 5 {
			pairGroup = group
			break
		}
	}

	if bombGroup == nil || pairGroup == nil {
		t.Fatal("Failed to find expected groups")
	}

	// 使用5的对子应该有较高的破坏度，因为会破坏炸弹
	if pairGroup.Damage == 0 {
		t.Error("Expected pair group to have damage > 0 due to breaking bomb")
	}

	t.Logf("Bomb group damage: %.2f", bombGroup.Damage)
	t.Logf("Pair group damage: %.2f", pairGroup.Damage)
}

func TestSmartAutoPlayAlgorithm_ScoreCalculation(t *testing.T) {
	algo := &SmartAutoPlayAlgorithm{level: 2}

	// 创建测试牌组
	groups := []*CardGroup{
		{
			Cards: []*sdk.Card{
				createCard(3, "Spade"),
				createCard(3, "Heart"),
			},
			CompType:  sdk.TypePair,
			CardCount: 2,
			Damage:    0,
			Strength:  3,
		},
		{
			Cards: []*sdk.Card{
				createCard(5, "Spade"),
				createCard(6, "Heart"),
				createCard(7, "Diamond"),
				createCard(8, "Club"),
				createCard(9, "Spade"),
			},
			CompType:  sdk.TypeStraight,
			CardCount: 5,
			Damage:    0.5,
			Strength:  7, // 平均值
		},
		{
			Cards: []*sdk.Card{
				createCard(10, "Spade"),
				createCard(10, "Heart"),
				createCard(11, "Diamond"),
				createCard(11, "Club"),
				createCard(12, "Spade"),
				createCard(12, "Heart"),
			},
			CompType:  sdk.TypePlate,
			CardCount: 6,
			Damage:    0,
			Strength:  11,
		},
	}

	algo.calculateScoresForGroups(groups)

	// 钢板应该有最高的分数（出牌多，有额外加分）
	bestGroup := algo.selectBestGroup(groups)
	if bestGroup == nil || bestGroup.CompType != sdk.TypePlate {
		t.Errorf("Expected Plate to be selected as best group, but got %v", bestGroup)
	}

	t.Logf("Group scores:")
	for _, group := range groups {
		t.Logf("  - %s: %.2f (cards=%d, damage=%.2f, strength=%.2f)",
			group.CompType.String(), group.Score, group.CardCount, group.Damage, group.Strength)
	}
}

func TestSmartAutoPlayAlgorithm_TributeSelection(t *testing.T) {
	algo := NewSmartAutoPlayAlgorithm(5).(*SmartAutoPlayAlgorithm) // 当前级别为5

	tests := []struct {
		name              string
		hand              []*sdk.Card
		excludeHeartTrump bool
		expectedNumber    int
	}{
		{
			name: "选择非王的最大牌",
			hand: []*sdk.Card{
				createCard(14, "Spade"), // A
				createCard(13, "Heart"), // K
				createCard(15, "Joker"), // 小王
			},
			excludeHeartTrump: false,
			expectedNumber:    14, // 应该选A，不选王
		},
		{
			name: "排除红桃主牌",
			hand: []*sdk.Card{
				createCard(5, "Heart"), // 红桃5（主牌）
				createCard(13, "Spade"), // 黑桃K
				createCard(12, "Diamond"), // 方块Q
			},
			excludeHeartTrump: true,
			expectedNumber:    13, // 应该选黑桃K，避开红桃5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := algo.SelectTributeCard(tt.hand, tt.excludeHeartTrump)
			if result == nil {
				t.Fatal("Expected a tribute card, but got nil")
			}
			if result.Number != tt.expectedNumber {
				t.Errorf("Expected card number %d, but got %d", tt.expectedNumber, result.Number)
			}
		})
	}
}

func TestSmartAutoPlayAlgorithm_ReturnTributeSelection(t *testing.T) {
	algo := NewSmartAutoPlayAlgorithm(2).(*SmartAutoPlayAlgorithm)

	hand := []*sdk.Card{
		createCard(3, "Spade"),
		createCard(7, "Heart"),
		createCard(10, "Diamond"),
		createCard(14, "Club"), // A
	}

	receivedCard := createCard(13, "Spade") // 收到K

	result := algo.SelectReturnTributeCard(hand, receivedCard)
	if result == nil {
		t.Fatal("Expected a return tribute card, but got nil")
	}
	if result.Number != 3 {
		t.Errorf("Expected to return the smallest card (3), but got %d", result.Number)
	}
}