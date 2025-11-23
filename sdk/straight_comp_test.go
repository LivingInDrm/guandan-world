package sdk

import (
	"testing"
)

// TestStraightSatisfyNew_TJQKA_NoWildcard 测试 10-J-Q-K-A 顺子（无变化牌）
func TestStraightSatisfyNew_TJQKA_NoWildcard(t *testing.T) {
	level := 5 // 避免10/J/Q/K/A是变化牌

	card10, _ := NewCard(10, "Spade", level)
	cardJ, _ := NewCard(11, "Club", level)
	cardQ, _ := NewCard(12, "Heart", level)
	cardK, _ := NewCard(13, "Diamond", level)
	cardA, _ := NewCard(14, "Spade", level)

	cards := []*Card{card10, cardJ, cardQ, cardK, cardA}

	t.Log("=== 10-J-Q-K-A 顺子（无变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized, comparisonKey := StraightSatisfy(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
		}
	}

	if !isValid {
		t.Error("10-J-Q-K-A 应该是合法的顺子")
	}

	if comparisonKey != 10 {
		t.Errorf("comparisonKey 应该是 10 (10-J-Q-K-A)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "10S,11C,12H,13D,1S"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
	}
}

// TestStraightSatisfyNew_TJQKA_OneWildcard 测试 10-J-Q-K-A 顺子（1个变化牌）
func TestStraightSatisfyNew_TJQKA_OneWildcard(t *testing.T) {
	level := 5

	t.Run("缺10", func(t *testing.T) {
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wildcard, _ := NewCard(5, "Heart", level) // 5是变化牌

		cards := []*Card{cardJ, cardQ, cardK, cardA, wildcard}

		t.Log("J-Q-K-A + wild (补10)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("J-Q-K-A+wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,11C,12H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺J", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardQ, cardK, cardA, wildcard}

		t.Log("10-Q-K-A + wild (补J)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-Q-K-A+wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,5H,12H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺Q", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardK, cardA, wildcard}

		t.Log("10-J-K-A + wild (补Q)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-K-A+wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,5H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺K", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardA, _ := NewCard(14, "Spade", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardQ, cardA, wildcard}

		t.Log("10-J-Q-A + wild (补K)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-Q-A+wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,12H,5H,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺A", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardQ, cardK, wildcard}

		t.Log("10-J-Q-K + wild (补A)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-Q-K+wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,12H,13D,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_TJQKA_TwoWildcards 测试 10-J-Q-K-A 顺子（2个变化牌）
func TestStraightSatisfyNew_TJQKA_TwoWildcards(t *testing.T) {
	level := 5

	t.Run("缺10+J", func(t *testing.T) {
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{cardQ, cardK, cardA, wild1, wild2}

		t.Log("Q-K-A + 2wild (补10和J)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("Q-K-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,5H,12H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺10+Q", func(t *testing.T) {
		cardJ, _ := NewCard(11, "Club", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{cardJ, cardK, cardA, wild1, wild2}

		t.Log("J-K-A + 2wild (补10和Q)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("J-K-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,11C,5H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺10+K", func(t *testing.T) {
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{cardJ, cardQ, cardA, wild1, wild2}

		t.Log("J-Q-A + 2wild (补10和K)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("J-Q-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,11C,12H,5H,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺10+A", func(t *testing.T) {
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{cardJ, cardQ, cardK, wild1, wild2}

		t.Log("J-Q-K + 2wild (补10和A)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("J-Q-K+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,11C,12H,13D,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺J+Q", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardK, _ := NewCard(13, "Diamond", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardK, cardA, wild1, wild2}

		t.Log("10-K-A + 2wild (补J和Q)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-K-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,5H,5H,13D,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺J+K", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardQ, cardA, wild1, wild2}

		t.Log("10-Q-A + 2wild (补J和K)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-Q-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,5H,12H,5H,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺J+A", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardQ, _ := NewCard(12, "Heart", level)
		cardK, _ := NewCard(13, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardQ, cardK, wild1, wild2}

		t.Log("10-Q-K + 2wild (补J和A)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-Q-K+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,5H,12H,13D,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺Q+K", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardA, _ := NewCard(14, "Spade", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardA, wild1, wild2}

		t.Log("10-J-A + 2wild (补Q和K)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-A+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,5H,5H,1S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺Q+A", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardK, _ := NewCard(13, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardK, wild1, wild2}

		t.Log("10-J-K + 2wild (补Q和A)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-K+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,5H,13D,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺K+A", func(t *testing.T) {
		card10, _ := NewCard(10, "Spade", level)
		cardJ, _ := NewCard(11, "Club", level)
		cardQ, _ := NewCard(12, "Heart", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)

		cards := []*Card{card10, cardJ, cardQ, wild1, wild2}

		t.Log("10-J-Q + 2wild (补K和A)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("10-J-Q+2wild 应该是合法的顺子")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,11C,12H,5H,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_A2345_NoWildcard 测试 A-2-3-4-5 顺子（无变化牌）
func TestStraightSatisfyNew_A2345_NoWildcard(t *testing.T) {
	level := 7 // 避免A/2/3/4/5是变化牌

	cardA, _ := NewCard(14, "Spade", level)
	card2, _ := NewCard(2, "Club", level)
	card3, _ := NewCard(3, "Heart", level)
	card4, _ := NewCard(4, "Diamond", level)
	card5, _ := NewCard(5, "Spade", level)

	cards := []*Card{cardA, card2, card3, card4, card5}

	t.Log("=== A-2-3-4-5 顺子（无变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (Number=%d, RawNumber=%d, IsWildcard=%v)", i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized, comparisonKey := StraightSatisfy(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (Number=%d, RawNumber=%d)", i, c.Name, c.Number, c.RawNumber)
		}
	}

	if !isValid {
		t.Error("A-2-3-4-5 应该是合法的顺子")
	}

	if comparisonKey != 1 {
		t.Errorf("comparisonKey 应该是 1 (A-2-3-4-5)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "1S,2C,3H,4D,5S"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
	}
}

// TestStraightSatisfyNew_A2345_OneWildcard 测试 A-2-3-4-5 顺子（1个变化牌）
func TestStraightSatisfyNew_A2345_OneWildcard(t *testing.T) {
	level := 7

	t.Run("缺A", func(t *testing.T) {
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level) // 7是变化牌

		cards := []*Card{card2, card3, card4, card5, wildcard}

		t.Log("2-3-4-5 + wild (可以补A或补6，应选更大的2-3-4-5-6)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-3-4-5+wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (选择2-3-4-5-6而非A-2-3-4-5)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2C,3H,4D,5S,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺2", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card3, card4, card5, wildcard}

		t.Log("A-3-4-5 + wild (补2)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-3-4-5+wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,7H,3H,4D,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺3", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card4, card5, wildcard}

		t.Log("A-2-4-5 + wild (补3)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-4-5+wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,7H,4D,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺4", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		card5, _ := NewCard(5, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card3, card5, wildcard}

		t.Log("A-2-3-5 + wild (补4)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-3-5+wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,3H,7H,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺5", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card3, card4, wildcard}

		t.Log("A-2-3-4 + wild (补5)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-3-4+wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,3H,4D,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_A2345_TwoWildcards 测试 A-2-3-4-5 顺子（2个变化牌）
func TestStraightSatisfyNew_A2345_TwoWildcards(t *testing.T) {
	level := 7

	t.Run("缺A+2", func(t *testing.T) {
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card3, card4, card5, wild1, wild2}

		t.Log("3-4-5 + 2wild (可补A+2, 2+6, 6+7，应选最大的3-4-5-6-7)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("3-4-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 3 {
			t.Errorf("comparisonKey 应该是 3 (选择3-4-5-6-7)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "3H,4D,5S,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺A+3", func(t *testing.T) {
		card2, _ := NewCard(2, "Club", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card4, card5, wild1, wild2}

		t.Log("2-4-5 + 2wild (可补A+3或3+6，应选更大的2-3-4-5-6)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-4-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (选择2-3-4-5-6)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2C,7H,4D,5S,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺A+4", func(t *testing.T) {
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card3, card5, wild1, wild2}

		t.Log("2-3-5 + 2wild (可补A+4或4+6，应选更大的2-3-4-5-6)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-3-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (选择2-3-4-5-6)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2C,3H,7H,5S,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺A+5", func(t *testing.T) {
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card3, card4, wild1, wild2}

		t.Log("2-3-4 + 2wild (可补A+5或5+6，应选更大的2-3-4-5-6)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-3-4+2wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (选择2-3-4-5-6)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2C,3H,4D,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺2+3", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card4, _ := NewCard(4, "Diamond", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card4, card5, wild1, wild2}

		t.Log("A-4-5 + 2wild (补2和3)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-4-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,7H,7H,4D,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺2+4", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card3, _ := NewCard(3, "Heart", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card3, card5, wild1, wild2}

		t.Log("A-3-5 + 2wild (补2和4)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-3-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,7H,3H,7H,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺2+5", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card3, _ := NewCard(3, "Heart", level)
		card4, _ := NewCard(4, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card3, card4, wild1, wild2}

		t.Log("A-3-4 + 2wild (补2和5)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-3-4+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,7H,3H,4D,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺3+4", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card5, _ := NewCard(5, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card5, wild1, wild2}

		t.Log("A-2-5 + 2wild (补3和4)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-5+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,7H,7H,5S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺3+5", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card4, _ := NewCard(4, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card4, wild1, wild2}

		t.Log("A-2-4 + 2wild (补3和5)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-4+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,7H,4D,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("缺4+5", func(t *testing.T) {
		cardA, _ := NewCard(14, "Spade", level)
		card2, _ := NewCard(2, "Club", level)
		card3, _ := NewCard(3, "Heart", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{cardA, card2, card3, wild1, wild2}

		t.Log("A-2-3 + 2wild (补4和5)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("A-2-3+2wild 应该是合法的顺子")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,2C,3H,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_Regular_NoWildcard 测试普通顺子（无变化牌）
func TestStraightSatisfyNew_Regular_NoWildcard(t *testing.T) {
	level := 5

	t.Run("2-3-4-5-6", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)
		card6, _ := NewCard(6, "Spade", level)

		cards := []*Card{card2, card3, card4, card5, card6}

		t.Log("2-3-4-5-6 (普通顺子)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-3-4-5-6 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2S,3C,4H,5D,6S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("3-4-5-6-7", func(t *testing.T) {
		card3, _ := NewCard(3, "Spade", level)
		card4, _ := NewCard(4, "Club", level)
		card5, _ := NewCard(5, "Heart", level)
		card6, _ := NewCard(6, "Diamond", level)
		card7, _ := NewCard(7, "Spade", level)

		cards := []*Card{card3, card4, card5, card6, card7}

		t.Log("3-4-5-6-7")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("3-4-5-6-7 应该是合法的顺子")
		}

		if comparisonKey != 3 {
			t.Errorf("comparisonKey 应该是 3，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "3S,4C,5H,6D,7S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("5-6-7-8-9", func(t *testing.T) {
		card5, _ := NewCard(5, "Spade", level)
		card6, _ := NewCard(6, "Club", level)
		card7, _ := NewCard(7, "Heart", level)
		card8, _ := NewCard(8, "Diamond", level)
		card9, _ := NewCard(9, "Spade", level)

		cards := []*Card{card5, card6, card7, card8, card9}

		t.Log("5-6-7-8-9")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("5-6-7-8-9 应该是合法的顺子")
		}

		if comparisonKey != 5 {
			t.Errorf("comparisonKey 应该是 5，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5S,6C,7H,8D,9S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("7-8-9-10-J", func(t *testing.T) {
		card7, _ := NewCard(7, "Spade", level)
		card8, _ := NewCard(8, "Club", level)
		card9, _ := NewCard(9, "Heart", level)
		card10, _ := NewCard(10, "Diamond", level)
		cardJ, _ := NewCard(11, "Spade", level)

		cards := []*Card{card7, card8, card9, card10, cardJ}

		t.Log("7-8-9-10-J")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("7-8-9-10-J 应该是合法的顺子")
		}

		if comparisonKey != 7 {
			t.Errorf("comparisonKey 应该是 7，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,8C,9H,10D,11S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("9-10-J-Q-K", func(t *testing.T) {
		card9, _ := NewCard(9, "Spade", level)
		card10, _ := NewCard(10, "Club", level)
		cardJ, _ := NewCard(11, "Heart", level)
		cardQ, _ := NewCard(12, "Diamond", level)
		cardK, _ := NewCard(13, "Spade", level)

		cards := []*Card{card9, card10, cardJ, cardQ, cardK}

		t.Log("9-10-J-Q-K")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("9-10-J-Q-K 应该是合法的顺子")
		}

		if comparisonKey != 9 {
			t.Errorf("comparisonKey 应该是 9，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "9S,10C,11H,12D,13S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_Regular_OneWildcard 测试普通顺子（1个变化牌）
func TestStraightSatisfyNew_Regular_OneWildcard(t *testing.T) {
	level := 7

	t.Run("2-3-4-5-6 缺4", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card5, _ := NewCard(5, "Heart", level)
		card6, _ := NewCard(6, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card3, card5, card6, wildcard}

		t.Log("2-3-5-6 + wild (补4)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-3-5-6+wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2S,3C,7H,5H,6S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("3-4-5-6-7 缺5", func(t *testing.T) {
		card3, _ := NewCard(3, "Spade", level)
		card4, _ := NewCard(4, "Club", level)
		card6, _ := NewCard(6, "Diamond", level)
		card7, _ := NewCard(7, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card3, card4, card6, card7, wildcard}

		t.Log("3-4-6-7 + wild (补5，注意7是变化牌)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("3-4-6-7+wild 应该是合法的顺子")
		}

		if comparisonKey != 3 {
			t.Errorf("comparisonKey 应该是 3，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "3S,4C,7H,6D,7S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("5-6-7-8-9 缺7", func(t *testing.T) {
		card5, _ := NewCard(5, "Spade", level)
		card6, _ := NewCard(6, "Club", level)
		card8, _ := NewCard(8, "Diamond", level)
		card9, _ := NewCard(9, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card5, card6, card8, card9, wildcard}

		t.Log("5-6-8-9 + wild (补7)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("5-6-8-9+wild 应该是合法的顺子")
		}

		if comparisonKey != 5 {
			t.Errorf("comparisonKey 应该是 5，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5S,6C,7H,8D,9S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("7-8-9-10-J 缺9", func(t *testing.T) {
		card7, _ := NewCard(7, "Spade", level)
		card8, _ := NewCard(8, "Club", level)
		card10, _ := NewCard(10, "Diamond", level)
		cardJ, _ := NewCard(11, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card7, card8, card10, cardJ, wildcard}

		t.Log("7-8-10-J + wild (补9)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("7-8-10-J+wild 应该是合法的顺子")
		}

		if comparisonKey != 7 {
			t.Errorf("comparisonKey 应该是 7，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,8C,7H,10D,11S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("9-10-J-Q-K 缺Q", func(t *testing.T) {
		card9, _ := NewCard(9, "Spade", level)
		card10, _ := NewCard(10, "Club", level)
		cardJ, _ := NewCard(11, "Heart", level)
		cardK, _ := NewCard(13, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card9, card10, cardJ, cardK, wildcard}

		t.Log("9-10-J-K + wild (补Q)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("9-10-J-K+wild 应该是合法的顺子")
		}

		if comparisonKey != 9 {
			t.Errorf("comparisonKey 应该是 9，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "9S,10C,11H,7H,13S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_Regular_TwoWildcards 测试普通顺子（2个变化牌）
func TestStraightSatisfyNew_Regular_TwoWildcards(t *testing.T) {
	level := 7

	t.Run("2-3-4-5-6 缺3+5", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card4, _ := NewCard(4, "Heart", level)
		card6, _ := NewCard(6, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card4, card6, wild1, wild2}

		t.Log("2-4-6 + 2wild (补3和5)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("2-4-6+2wild 应该是合法的顺子")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2S,7H,4H,7H,6S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("3-4-5-6-7 缺4+6", func(t *testing.T) {
		card3, _ := NewCard(3, "Spade", level)
		card5, _ := NewCard(5, "Heart", level)
		card7, _ := NewCard(7, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card3, card5, card7, wild1, wild2}

		t.Log("3-5-7 + 2wild (补4和6)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("3-5-7+2wild 应该是合法的顺子")
		}

		if comparisonKey != 3 {
			t.Errorf("comparisonKey 应该是 3，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "3S,7H,5H,7H,7S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("5-6-7-8-9 缺6+8", func(t *testing.T) {
		card5, _ := NewCard(5, "Spade", level)
		card7, _ := NewCard(7, "Spade", level)
		card9, _ := NewCard(9, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card5, card7, card9, wild1, wild2}

		t.Log("5-7-9 + 2wild (补6和8)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("5-7-9+2wild 应该是合法的顺子")
		}

		if comparisonKey != 5 {
			t.Errorf("comparisonKey 应该是 5，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5S,7H,7S,7H,9S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("7-8-9-10-J 缺8+10", func(t *testing.T) {
		card7, _ := NewCard(7, "Spade", level)
		card9, _ := NewCard(9, "Diamond", level)
		cardJ, _ := NewCard(11, "Spade", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card7, card9, cardJ, wild1, wild2}

		t.Log("7-9-J + 2wild (补8和10)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("7-9-J+2wild 应该是合法的顺子")
		}

		if comparisonKey != 7 {
			t.Errorf("comparisonKey 应该是 7，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,7H,9D,7H,11S"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("9-10-J-Q-K 缺10+K", func(t *testing.T) {
		card9, _ := NewCard(9, "Spade", level)
		cardJ, _ := NewCard(11, "Heart", level)
		cardQ, _ := NewCard(12, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card9, cardJ, cardQ, wild1, wild2}

		t.Log("9-J-Q + 2wild (补10和K)")
		isValid, normalized, comparisonKey := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("9-J-Q+2wild 应该是合法的顺子")
		}

		if comparisonKey != 9 {
			t.Errorf("comparisonKey 应该是 9，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "9S,7H,11H,12D,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestStraightSatisfyNew_Invalid 测试非法顺子
func TestStraightSatisfyNew_Invalid(t *testing.T) {
	level := 7

	t.Run("重复数字: 3-3-4-5-6", func(t *testing.T) {
		card31, _ := NewCard(3, "Spade", level)
		card32, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)
		card6, _ := NewCard(6, "Spade", level)

		cards := []*Card{card31, card32, card4, card5, card6}

		t.Log("3-3-4-5-6 (重复数字)")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("3-3-4-5-6 不应该是合法的顺子（重复数字）")
		}
	})

	t.Run("不连续: 2-3-5-6-7 (缺4)", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card5, _ := NewCard(5, "Heart", level)
		card6, _ := NewCard(6, "Diamond", level)
		card7, _ := NewCard(7, "Spade", level)

		cards := []*Card{card2, card3, card5, card6, card7}

		t.Log("2-3-5-6-7 (不连续，缺4)")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2-3-5-6-7 不应该是合法的顺子（缺4）")
		}
	})

	t.Run("间隔太大: 2-4-6-8-10", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card4, _ := NewCard(4, "Club", level)
		card6, _ := NewCard(6, "Heart", level)
		card8, _ := NewCard(8, "Diamond", level)
		card10, _ := NewCard(10, "Spade", level)

		cards := []*Card{card2, card4, card6, card8, card10}

		t.Log("2-4-6-8-10 (间隔太大)")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2-4-6-8-10 不应该是合法的顺子（间隔太大）")
		}
	})

	t.Run("包含大王", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)
		joker, _ := NewCard(16, "Joker", level) // 大王

		cards := []*Card{card2, card3, card4, card5, joker}

		t.Log("2-3-4-5-大王")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("包含大王的牌组不应该是合法的顺子")
		}
	})

	t.Run("包含小王", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)
		joker, _ := NewCard(15, "Joker", level) // 小王

		cards := []*Card{card2, card3, card4, card5, joker}

		t.Log("2-3-4-5-小王")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("包含小王的牌组不应该是合法的顺子")
		}
	})

	t.Run("长度不对: 4张", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)

		cards := []*Card{card2, card3, card4, card5}

		t.Log("长度不对: 只有4张")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("4张牌不应该是合法的顺子（长度不对）")
		}
	})

	t.Run("长度不对: 6张", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card4, _ := NewCard(4, "Heart", level)
		card5, _ := NewCard(5, "Diamond", level)
		card6, _ := NewCard(6, "Spade", level)
		card7, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card3, card4, card5, card6, card7}

		t.Log("长度不对: 6张")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("6张牌不应该是合法的顺子（长度不对）")
		}
	})

	t.Run("全是同一张牌: 5-5-5-5-5", func(t *testing.T) {
		card51, _ := NewCard(5, "Spade", level)
		card52, _ := NewCard(5, "Club", level)
		card53, _ := NewCard(5, "Heart", level)
		card54, _ := NewCard(5, "Diamond", level)
		wild, _ := NewCard(7, "Heart", level) // 变化牌（相当于第五张5）

		cards := []*Card{card51, card52, card53, card54, wild}

		t.Log("5-5-5-5-5 (全是同一张牌)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("全是同一数字不应该是合法的顺子")
		}
	})

	t.Run("跨度超范围: Q-K-A-2-3", func(t *testing.T) {
		cardQ, _ := NewCard(12, "Spade", level)
		cardK, _ := NewCard(13, "Club", level)
		cardA, _ := NewCard(14, "Heart", level)
		card2, _ := NewCard(2, "Diamond", level)
		card3, _ := NewCard(3, "Spade", level)

		cards := []*Card{cardQ, cardK, cardA, card2, card3}

		t.Log("Q-K-A-2-3 (跨度超范围)")
		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("Q-K-A-2-3 不应该是合法的顺子（A不能同时作为高低位）")
		}
	})

	t.Run("无法凑成顺子: 2+3+6+7+wild", func(t *testing.T) {
		card2, _ := NewCard(2, "Spade", level)
		card3, _ := NewCard(3, "Club", level)
		card6, _ := NewCard(6, "Heart", level)
		card7, _ := NewCard(7, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level)

		cards := []*Card{card2, card3, card6, card7, wildcard}

		t.Log("2-3-6-7 + wild (缺4和5，只有1个变化牌)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, _, _ := StraightSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2-3-6-7+wild 不应该是合法的顺子（缺两张但只有1个变化牌）")
		}
	})
}
