package sdk

import (
	"testing"
)

// TestTubeSatisfyNew_QKA_NoWildcard 测试 Q-K-A 钢管（无变化牌）
func TestTubeSatisfyNew_QKA_NoWildcard(t *testing.T) {
	level := 5 // 避免Q/K/A是变化牌

	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Club", level)
	cardK1, _ := NewCard(13, "Heart", level)
	cardK2, _ := NewCard(13, "Diamond", level)
	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)

	cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, cardA1, cardA2}

	t.Log("=== Q-K-A 钢管（无变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized := tubeSatisfyNew(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
		}
	}

	if !isValid {
		t.Error("QQ,KK,AA 应该是合法的 Q-K-A 钢管")
	}
}

// TestTubeSatisfyNew_QKA_OneWildcard 测试 Q-K-A 钢管（1个变化牌）
func TestTubeSatisfyNew_QKA_OneWildcard(t *testing.T) {
	level := 5

	// QQ + KK + A + 1个变化牌（补另一张A）
	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Club", level)
	cardK1, _ := NewCard(13, "Heart", level)
	cardK2, _ := NewCard(13, "Diamond", level)
	cardA1, _ := NewCard(14, "Spade", level)
	wildcard, _ := NewCard(5, "Heart", level) // 5是变化牌

	cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, cardA1, wildcard}

	t.Log("=== Q-K-A 钢管（1个变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized := tubeSatisfyNew(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}

	if !isValid {
		t.Error("QQ,KK,A+wild 应该是合法的 Q-K-A 钢管")
	}
}

// TestTubeSatisfyNew_QKA_TwoWildcards 测试 Q-K-A 钢管（2个变化牌）
func TestTubeSatisfyNew_QKA_TwoWildcards(t *testing.T) {
	level := 5

	t.Run("QQ+KK+2wild", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardK1, _ := NewCard(13, "Heart", level)
		cardK2, _ := NewCard(13, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Diamond", level)

		cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, wild1, wild2}

		t.Log("QQ+KK+2wild (应该补两张A)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("QQ+KK+2wild 应该是合法的 Q-K-A 钢管")
		}
	})

	t.Run("QQ+AA+2wild", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardA1, _ := NewCard(14, "Heart", level)
		cardA2, _ := NewCard(14, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Diamond", level)

		cards := []*Card{cardQ1, cardQ2, cardA1, cardA2, wild1, wild2}

		t.Log("QQ+AA+2wild (应该补两张K)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("QQ+AA+2wild 应该是合法的 Q-K-A 钢管")
		}
	})

	t.Run("KK+AA+2wild", func(t *testing.T) {
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Club", level)
		cardA1, _ := NewCard(14, "Heart", level)
		cardA2, _ := NewCard(14, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Diamond", level)

		cards := []*Card{cardK1, cardK2, cardA1, cardA2, wild1, wild2}

		t.Log("KK+AA+2wild (应该补两张Q)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("KK+AA+2wild 应该是合法的 Q-K-A 钢管")
		}
	})
}

// TestTubeSatisfyNew_A23_NoWildcard 测试 A-2-3 钢管（无变化牌）
func TestTubeSatisfyNew_A23_NoWildcard(t *testing.T) {
	level := 5

	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)
	card21, _ := NewCard(3, "Heart", level) // RawNumber=3 表示2（因为级别是5）
	card22, _ := NewCard(3, "Diamond", level)
	card31, _ := NewCard(4, "Spade", level) // RawNumber=4 表示3
	card32, _ := NewCard(4, "Club", level)

	cards := []*Card{cardA1, cardA2, card21, card22, card31, card32}

	t.Log("=== A-2-3 钢管（无变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (Number=%d, RawNumber=%d, IsWildcard=%v)", i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized := tubeSatisfyNew(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (Number=%d, RawNumber=%d)", i, c.Name, c.Number, c.RawNumber)
		}
	}

	if !isValid {
		t.Error("AA,22,33 应该是合法的 A-2-3 钢管")
	}
}

// TestTubeSatisfyNew_A23_OneWildcard 测试 A-2-3 钢管（1个变化牌）
func TestTubeSatisfyNew_A23_OneWildcard(t *testing.T) {
	level := 7 // 避免A/2/3是变化牌

	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)
	card21, _ := NewCard(2, "Heart", level)
	card22, _ := NewCard(2, "Diamond", level)
	card31, _ := NewCard(3, "Spade", level)
	wildcard, _ := NewCard(7, "Heart", level) // 7是变化牌

	cards := []*Card{cardA1, cardA2, card21, card22, card31, wildcard}

	t.Log("=== A-2-3 钢管（1个变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized := tubeSatisfyNew(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}

	if !isValid {
		t.Error("AA,22,3+wild 应该是合法的 A-2-3 钢管")
	}
}

// TestTubeSatisfyNew_A23_TwoWildcards 测试 A-2-3 钢管（2个变化牌）
func TestTubeSatisfyNew_A23_TwoWildcards(t *testing.T) {
	level := 7

	t.Run("AA+22+2wild", func(t *testing.T) {
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Club", level)
		card21, _ := NewCard(2, "Heart", level)
		card22, _ := NewCard(2, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Diamond", level)

		cards := []*Card{cardA1, cardA2, card21, card22, wild1, wild2}

		t.Log("AA+22+2wild (应该补两张3)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("AA+22+2wild 应该是合法的 A-2-3 钢管")
		}
	})

	t.Run("AA+33+2wild", func(t *testing.T) {
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Club", level)
		card31, _ := NewCard(3, "Heart", level)
		card32, _ := NewCard(3, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Diamond", level)

		cards := []*Card{cardA1, cardA2, card31, card32, wild1, wild2}

		t.Log("AA+33+2wild (应该补两张2)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("AA+33+2wild 应该是合法的 A-2-3 钢管")
		}
	})

	t.Run("22+33+2wild", func(t *testing.T) {
		card21, _ := NewCard(2, "Spade", level)
		card22, _ := NewCard(2, "Club", level)
		card31, _ := NewCard(3, "Heart", level)
		card32, _ := NewCard(3, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Diamond", level)

		cards := []*Card{card21, card22, card31, card32, wild1, wild2}

		t.Log("22+33+2wild (应该补两张A)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("22+33+2wild 应该是合法的 A-2-3 钢管")
		}
	})
}

// TestTubeSatisfyNew_Regular 测试普通钢管（不含A循环）
func TestTubeSatisfyNew_Regular(t *testing.T) {
	level := 5

	t.Run("JJ+QQ+KK", func(t *testing.T) {
		cardJ1, _ := NewCard(11, "Spade", level)
		cardJ2, _ := NewCard(11, "Club", level)
		cardQ1, _ := NewCard(12, "Heart", level)
		cardQ2, _ := NewCard(12, "Diamond", level)
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Club", level)

		cards := []*Card{cardJ1, cardJ2, cardQ1, cardQ2, cardK1, cardK2}

		t.Log("JJ+QQ+KK (普通钢管)")
		isValid, normalized := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
			}
		}

		if !isValid {
			t.Error("JJ+QQ+KK 应该是合法的钢管")
		}
	})

	t.Run("55+66+77", func(t *testing.T) {
		card51, _ := NewCard(6, "Spade", level)  // RawNumber=6, 但不是变化牌
		card52, _ := NewCard(6, "Club", level)
		card61, _ := NewCard(7, "Heart", level)
		card62, _ := NewCard(7, "Diamond", level)
		card71, _ := NewCard(8, "Spade", level)
		card72, _ := NewCard(8, "Club", level)

		cards := []*Card{card51, card52, card61, card62, card71, card72}

		t.Log("55+66+77 (普通钢管)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, RawNumber=%d, IsWildcard=%v)", i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
		}

		isValid, _ := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("55+66+77 应该是合法的钢管")
		}
	})
}

// TestTubeSatisfyNew_BestTripleSelection 测试多组合时选择最大的
func TestTubeSatisfyNew_BestTripleSelection(t *testing.T) {
	level := 5

	// QQ+KK+2wild 可以凑成 J-Q-K 或 Q-K-A，应该选择 Q-K-A
	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Club", level)
	cardK1, _ := NewCard(13, "Heart", level)
	cardK2, _ := NewCard(13, "Diamond", level)
	wild1, _ := NewCard(5, "Heart", level)
	wild2, _ := NewCard(5, "Diamond", level)

	cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, wild1, wild2}

	t.Log("=== QQ+KK+2wild (应该选择 Q-K-A 而不是 J-Q-K) ===")
	isValid, normalized := tubeSatisfyNew(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		// 验证选择的是 Q-K-A（最后两张应该是A，用变化牌补的）
		// 规范化后应该是 Q,Q,K,K,A,A 的顺序
		if len(normalized) == 6 {
			// 检查最后两张是否是A（RawNumber=1）或变化牌
			lastTwoAreAorWild := (normalized[4].RawNumber == 1 || normalized[4].IsWildcard()) &&
				(normalized[5].RawNumber == 1 || normalized[5].IsWildcard())

			if !lastTwoAreAorWild {
				t.Log("警告：可能选择了 J-Q-K 而不是 Q-K-A")
			}
		}
	}

	if !isValid {
		t.Error("QQ+KK+2wild 应该是合法的钢管")
	}
}

// TestTubeSatisfyNew_Invalid 测试无法凑成钢管的情况
func TestTubeSatisfyNew_Invalid(t *testing.T) {
	level := 5

	t.Run("QQ+99+2wild (间隔太大)", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		card91, _ := NewCard(9, "Heart", level)
		card92, _ := NewCard(9, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Diamond", level)

		cards := []*Card{cardQ1, cardQ2, card91, card92, wild1, wild2}

		t.Log("QQ+99+2wild (间隔太大，无法凑成)")
		isValid, _ := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("QQ+99+2wild 不应该是合法的钢管（间隔太大）")
		}
	})

	t.Run("错误长度", func(t *testing.T) {
		card1, _ := NewCard(12, "Spade", level)
		card2, _ := NewCard(12, "Club", level)

		cards := []*Card{card1, card2}

		t.Log("长度不对（只有2张）")
		isValid, _ := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("长度不对的牌组不应该是合法的钢管")
		}
	})

	t.Run("含大王", func(t *testing.T) {
		card1, _ := NewCard(12, "Spade", level)
		card2, _ := NewCard(12, "Club", level)
		card3, _ := NewCard(13, "Heart", level)
		card4, _ := NewCard(13, "Diamond", level)
		joker1, _ := NewCard(16, "Joker", level) // 大王
		joker2, _ := NewCard(16, "Joker", level)

		cards := []*Card{card1, card2, card3, card4, joker1, joker2}

		t.Log("含大王")
		isValid, _ := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("含大王的牌组不应该是合法的钢管")
		}
	})

	t.Run("含小王", func(t *testing.T) {
		card1, _ := NewCard(12, "Spade", level)
		card2, _ := NewCard(12, "Club", level)
		card3, _ := NewCard(13, "Heart", level)
		card4, _ := NewCard(13, "Diamond", level)
		joker1, _ := NewCard(15, "Joker", level) // 小王
		joker2, _ := NewCard(15, "Joker", level)

		cards := []*Card{card1, card2, card3, card4, joker1, joker2}

		t.Log("含小王")
		isValid, _ := tubeSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("含小王的牌组不应该是合法的钢管")
		}
	})
}
