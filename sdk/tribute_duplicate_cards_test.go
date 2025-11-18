package sdk

import (
	"testing"
)

// TestTributeWithDuplicateCards 验证修复：当手牌中有两张相同点数相同花色的牌时，
// 贡牌应该精确移除指定的那一张（基于DeckIndex），而不是错误地删除另一张副本
func TestTributeWithDuplicateCards(t *testing.T) {
	// 创建两张黑桃A，DeckIndex不同
	aceSpade1, _ := NewCard(14, "Spade", 5)
	aceSpade1.DeckIndex = 12 // 第一副牌的黑桃A

	aceSpade2, _ := NewCard(14, "Spade", 5)
	aceSpade2.DeckIndex = 64 // 第二副牌的黑桃A

	otherCard1, _ := NewCard(13, "Heart", 5)
	otherCard1.DeckIndex = 20

	otherCard2, _ := NewCard(12, "Diamond", 5)
	otherCard2.DeckIndex = 30

	// 玩家手牌包含两张黑桃A
	playerHands := [4][]*Card{
		{},                                    // 玩家0: 接收者
		{otherCard1, otherCard2},              // 玩家1
		{},                                    // 玩家2
		{aceSpade1, aceSpade2, otherCard1.Clone()}, // 玩家3: 贡者，有两张黑桃A
	}

	// 设置上贡阶段：玩家3上贡给玩家0
	tributePhase := &TributePhase{
		Status:          TributeStatusFinished,
		TributePairs:    make([]*TributePair, 0),
		PoolCards:       make([]*Card, 0),
		SelectingPlayer: -1,
	}

	// 明确指定贡牌为第一张黑桃A (DeckIndex=12)
	pair := &TributePair{
		Giver:       3,
		Receiver:    0,
		TributeCard: aceSpade1, // 指定第一张
		ReturnCard:  otherCard2,
	}
	tributePhase.TributePairs = append(tributePhase.TributePairs, pair)

	// 记录原始手牌数量
	player3OriginalCount := len(playerHands[3])
	player0OriginalCount := len(playerHands[0])

	// 应用上贡效果
	tm := NewTributeManager(5)
	err := tm.ApplyTributeToHands(tributePhase, &playerHands)
	if err != nil {
		t.Fatalf("应用上贡失败: %v", err)
	}

	// 验证1: 玩家3应该失去指定的那张黑桃A (DeckIndex=12)
	found12 := false
	found64 := false
	for _, card := range playerHands[3] {
		if card.DeckIndex == 12 {
			found12 = true
		}
		if card.DeckIndex == 64 {
			found64 = true
		}
	}

	if found12 {
		t.Errorf("❌ 玩家3手牌中仍有 DeckIndex=12 的黑桃A，应该被移除")
	}
	if !found64 {
		t.Errorf("❌ 玩家3手牌中缺少 DeckIndex=64 的黑桃A，不应该被移除")
	}

	// 验证2: 玩家0应该获得指定的那张黑桃A (DeckIndex=12)
	foundInReceiver := false
	for _, card := range playerHands[0] {
		if card.DeckIndex == 12 {
			foundInReceiver = true
			break
		}
	}

	if !foundInReceiver {
		t.Errorf("❌ 玩家0手牌中没有 DeckIndex=12 的黑桃A，贡牌未正确转移")
	}

	// 验证3: 手牌数量变化正确（玩家3少1张贡牌+多1张还贡，玩家0多1张贡牌-1张还贡）
	if len(playerHands[3]) != player3OriginalCount {
		t.Errorf("玩家3手牌数量错误: 期望 %d, 实际 %d", player3OriginalCount, len(playerHands[3]))
	}
	if len(playerHands[0]) != player0OriginalCount {
		t.Errorf("玩家0手牌数量错误: 期望 %d, 实际 %d", player0OriginalCount, len(playerHands[0]))
	}

	// 验证4: 还贡也应该精确匹配
	foundReturnCard := false
	for _, card := range playerHands[3] {
		if card.DeckIndex == 30 { // otherCard2 的 DeckIndex
			foundReturnCard = true
			break
		}
	}
	if !foundReturnCard {
		t.Errorf("❌ 玩家3未收到还贡卡 (DeckIndex=30)")
	}

	t.Logf("✅ 重复牌场景验证通过:")
	t.Logf("   - 玩家3正确失去 DeckIndex=12 的黑桃A")
	t.Logf("   - 玩家3保留了 DeckIndex=64 的黑桃A")
	t.Logf("   - 玩家0正确获得 DeckIndex=12 的黑桃A")
	t.Logf("   - 手牌数量正确")
}

// TestRemoveCardFromHandPrecision 测试 removeCardFromHand 的精确性
func TestRemoveCardFromHandPrecision(t *testing.T) {
	tm := NewTributeManager(5)

	// 创建三张相同的牌（点数和花色都相同，但DeckIndex不同）
	card1, _ := NewCard(10, "Club", 5)
	card1.DeckIndex = 10

	card2, _ := NewCard(10, "Club", 5)
	card2.DeckIndex = 58

	card3, _ := NewCard(10, "Club", 5)
	card3.DeckIndex = 80

	hand := []*Card{card1, card2, card3}

	// 移除第二张 (DeckIndex=58)
	result := tm.removeCardFromHand(hand, card2)

	// 验证只移除了指定的那张
	if len(result) != 2 {
		t.Fatalf("期望手牌剩余2张，实际 %d 张", len(result))
	}

	// 检查剩余的是哪两张
	remainingIndexes := make(map[int]bool)
	for _, card := range result {
		remainingIndexes[card.DeckIndex] = true
	}

	if !remainingIndexes[10] {
		t.Errorf("❌ DeckIndex=10 的牌应该保留")
	}
	if remainingIndexes[58] {
		t.Errorf("❌ DeckIndex=58 的牌应该被移除")
	}
	if !remainingIndexes[80] {
		t.Errorf("❌ DeckIndex=80 的牌应该保留")
	}

	t.Logf("✅ removeCardFromHand 精确匹配验证通过")
}

// TestCardsEqualByDeckIndex 测试 cardsEqual 使用 DeckIndex 比较
func TestCardsEqualByDeckIndex(t *testing.T) {
	tm := NewTributeManager(5)

	// 创建两张"看起来相同"的牌（点数和花色相同）
	card1, _ := NewCard(14, "Heart", 5)
	card1.DeckIndex = 5

	card2, _ := NewCard(14, "Heart", 5)
	card2.DeckIndex = 5 // 相同的 DeckIndex

	card3, _ := NewCard(14, "Heart", 5)
	card3.DeckIndex = 53 // 不同的 DeckIndex

	// 测试相等情况
	if !tm.cardsEqual(card1, card2) {
		t.Errorf("相同 DeckIndex 的牌应该判定为相等")
	}

	// 测试不等情况
	if tm.cardsEqual(card1, card3) {
		t.Errorf("不同 DeckIndex 的牌不应该判定为相等，即使点数和花色相同")
	}

	t.Logf("✅ cardsEqual 基于 DeckIndex 比较验证通过")
}
