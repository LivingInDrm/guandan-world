package sdk

import (
	"testing"
)

// TestStraightIntegration_Basic 测试 Straight 集成基础功能
func TestStraightIntegration_Basic(t *testing.T) {
	level := 5 // 避免10/J/Q/K/A是变化牌

	card10, _ := NewCard(10, "Spade", level)
	cardJ, _ := NewCard(11, "Club", level)
	cardQ, _ := NewCard(12, "Heart", level)
	cardK, _ := NewCard(13, "Diamond", level)
	cardA, _ := NewCard(14, "Spade", level)

	cards := []*Card{card10, cardJ, cardQ, cardK, cardA}

	// 使用 NewStraight 创建顺子
	straight := NewStraight(cards)

	if !straight.IsValid() {
		t.Error("10-J-Q-K-A 应该是合法的顺子")
	}

	if straight.ComparisonKey != 10 {
		t.Errorf("ComparisonKey 应该是 10，实际: %d", straight.ComparisonKey)
	}

	t.Logf("✓ 顺子有效: %v", straight.IsValid())
	t.Logf("✓ ComparisonKey: %d", straight.ComparisonKey)
}

// TestStraightIntegration_Comparison 测试顺子比较功能
func TestStraightIntegration_Comparison(t *testing.T) {
	level := 7

	// 创建 A-2-3-4-5 (最小顺子, key=1)
	cardA1, _ := NewCard(14, "Spade", level)
	card2, _ := NewCard(2, "Club", level)
	card3, _ := NewCard(3, "Heart", level)
	card4, _ := NewCard(4, "Diamond", level)
	card5, _ := NewCard(5, "Spade", level)
	straight1 := NewStraight([]*Card{cardA1, card2, card3, card4, card5})

	// 创建 10-J-Q-K-A (最大顺子, key=10)
	card10, _ := NewCard(10, "Spade", level)
	cardJ, _ := NewCard(11, "Club", level)
	cardQ, _ := NewCard(12, "Heart", level)
	cardK, _ := NewCard(13, "Diamond", level)
	cardA2, _ := NewCard(14, "Heart", level)
	straight2 := NewStraight([]*Card{card10, cardJ, cardQ, cardK, cardA2})

	// 验证 ComparisonKey
	if straight1.ComparisonKey != 1 {
		t.Errorf("A-2-3-4-5 的 ComparisonKey 应该是 1，实际: %d", straight1.ComparisonKey)
	}
	if straight2.ComparisonKey != 10 {
		t.Errorf("10-J-Q-K-A 的 ComparisonKey 应该是 10，实际: %d", straight2.ComparisonKey)
	}

	// 测试比较
	if !straight2.GreaterThan(straight1) {
		t.Error("10-J-Q-K-A 应该大于 A-2-3-4-5")
	}
	if straight1.GreaterThan(straight2) {
		t.Error("A-2-3-4-5 不应该大于 10-J-Q-K-A")
	}

	t.Logf("✓ A-2-3-4-5 ComparisonKey: %d", straight1.ComparisonKey)
	t.Logf("✓ 10-J-Q-K-A ComparisonKey: %d", straight2.ComparisonKey)
	t.Logf("✓ 比较逻辑正确")
}

// TestStraightIntegration_Wildcard 测试含万能牌的顺子
func TestStraightIntegration_Wildcard(t *testing.T) {
	level := 5

	// 创建 10-J-Q-K-A 顺子，用万能牌替代 Q
	card10, _ := NewCard(10, "Spade", level)
	cardJ, _ := NewCard(11, "Club", level)
	cardK, _ := NewCard(13, "Diamond", level)
	cardA, _ := NewCard(14, "Spade", level)
	wildcard, _ := NewCard(5, "Heart", level) // 5是变化牌

	cards := []*Card{card10, cardJ, cardK, cardA, wildcard}
	straight := NewStraight(cards)

	if !straight.IsValid() {
		t.Error("10-J-K-A+wild 应该是合法的顺子")
	}

	if straight.ComparisonKey != 10 {
		t.Errorf("ComparisonKey 应该是 10，实际: %d", straight.ComparisonKey)
	}

	t.Logf("✓ 含万能牌的顺子有效")
	t.Logf("✓ ComparisonKey: %d", straight.ComparisonKey)
}
