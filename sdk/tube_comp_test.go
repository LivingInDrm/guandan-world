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

	isValid, normalized, comparisonKey := TubeSatisfy(cards)

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

	if comparisonKey != 12 {
		t.Errorf("comparisonKey 应该是 12 (Q-K-A)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "12S,12C,13H,13D,1S,1C"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
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

	isValid, normalized, comparisonKey := TubeSatisfy(cards)

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

	if comparisonKey != 12 {
		t.Errorf("comparisonKey 应该是 12 (Q-K-A)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "12S,12C,13H,13D,1S,5H"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
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
		wild2, _ := NewCard(5, "Heart", level) // 两张红桃5都是变化牌

		cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, wild1, wild2}

		t.Log("QQ+KK+2wild (应该补两张A)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("QQ+KK+2wild 应该是合法的 Q-K-A 钢管")
		}

		if comparisonKey != 12 {
			t.Errorf("comparisonKey 应该是 12 (Q-K-A)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "12S,12C,13H,13D,5H,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("QQ+AA+2wild", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardA1, _ := NewCard(14, "Heart", level)
		cardA2, _ := NewCard(14, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level) // 两张红桃5都是变化牌

		cards := []*Card{cardQ1, cardQ2, cardA1, cardA2, wild1, wild2}

		t.Log("QQ+AA+2wild (应该补两张K)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("QQ+AA+2wild 应该是合法的 Q-K-A 钢管")
		}

		if comparisonKey != 12 {
			t.Errorf("comparisonKey 应该是 12 (Q-K-A)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "12S,12C,5H,5H,1H,1D"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("KK+AA+2wild", func(t *testing.T) {
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Club", level)
		cardA1, _ := NewCard(14, "Heart", level)
		cardA2, _ := NewCard(14, "Diamond", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level) // 两张红桃5都是变化牌

		cards := []*Card{cardK1, cardK2, cardA1, cardA2, wild1, wild2}

		t.Log("KK+AA+2wild (应该补两张Q)")
		isValid, normalized, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("KK+AA+2wild 应该是合法的 Q-K-A 钢管")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,5H,13S,13C,1H,1D"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestTubeSatisfyNew_A23_NoWildcard 测试 A-2-3 钢管（无变化牌）
func TestTubeSatisfyNew_A23_NoWildcard(t *testing.T) {
	level := 5

	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)
	card21, _ := NewCard(2, "Heart", level)
	card22, _ := NewCard(2, "Diamond", level)
	card31, _ := NewCard(3, "Spade", level)
	card32, _ := NewCard(3, "Club", level)

	cards := []*Card{cardA1, cardA2, card21, card22, card31, card32}

	t.Log("=== A-2-3 钢管（无变化牌）===")
	for i, c := range cards {
		t.Logf("  牌[%d]: %s (Number=%d, RawNumber=%d, IsWildcard=%v)", i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}

	isValid, normalized, comparisonKey := TubeSatisfy(cards)

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

	if comparisonKey != 1 {
		t.Errorf("comparisonKey 应该是 1 (A-2-3)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "1S,1C,2H,2D,3S,3C"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
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

	isValid, normalized, comparisonKey := TubeSatisfy(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}

	if !isValid {
		t.Error("AA,22,3+wild 应该是合法的 A-2-3 钢管")
	}

	if comparisonKey != 1 {
		t.Errorf("comparisonKey 应该是 1 (A-2-3)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "1S,1C,2H,2D,3S,7H"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
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
		wild2, _ := NewCard(7, "Heart", level) // 两张红桃7都是变化牌

		cards := []*Card{cardA1, cardA2, card21, card22, wild1, wild2}

		t.Log("AA+22+2wild (应该补两张3)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("AA+22+2wild 应该是合法的 A-2-3 钢管")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1 (A-2-3)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1C,2H,2D,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("AA+33+2wild", func(t *testing.T) {
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Club", level)
		card31, _ := NewCard(3, "Heart", level)
		card32, _ := NewCard(3, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level) // 两张红桃7都是变化牌

		cards := []*Card{cardA1, cardA2, card31, card32, wild1, wild2}

		t.Log("AA+33+2wild (应该补两张2)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("AA+33+2wild 应该是合法的 A-2-3 钢管")
		}

		if comparisonKey != 1 {
			t.Errorf("comparisonKey 应该是 1 (A-2-3)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1C,7H,7H,3H,3D"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("22+33+2wild", func(t *testing.T) {
		card21, _ := NewCard(2, "Spade", level)
		card22, _ := NewCard(2, "Club", level)
		card31, _ := NewCard(3, "Heart", level)
		card32, _ := NewCard(3, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level) // 两张红桃7都是变化牌

		cards := []*Card{card21, card22, card31, card32, wild1, wild2}

		t.Log("22+33+2wild (可以补两张A凑成A-2-3，或补两张4凑成2-3-4，应选更大的2-3-4)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("22+33+2wild 应该是合法的 A-2-3 钢管")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (选择 2-3-4 而非 A-2-3)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2S,2C,3H,3D,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
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
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
			}
		}

		if !isValid {
			t.Error("JJ+QQ+KK 应该是合法的钢管")
		}

		if comparisonKey != 11 {
			t.Errorf("comparisonKey 应该是 11 (J-Q-K)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "11S,11C,12H,12D,13S,13C"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("66+77+88", func(t *testing.T) {
		card61, _ := NewCard(6, "Spade", level)
		card62, _ := NewCard(6, "Club", level)
		card71, _ := NewCard(7, "Heart", level)
		card72, _ := NewCard(7, "Diamond", level)
		card81, _ := NewCard(8, "Spade", level)
		card82, _ := NewCard(8, "Club", level)

		cards := []*Card{card61, card62, card71, card72, card81, card82}

		t.Log("66+77+88 (普通钢管，level=5时避开变化牌5)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, RawNumber=%d, IsWildcard=%v)", i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if !isValid {
			t.Error("55+66+77 应该是合法的钢管")
		}

		if comparisonKey != 6 {
			t.Errorf("comparisonKey 应该是 6 (6-7-8)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "6S,6C,7H,7D,8S,8C"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestTubeSatisfyNew_RegularWithWildcards 测试普通钢管 + 变化牌
func TestTubeSatisfyNew_RegularWithWildcards(t *testing.T) {
	level := 7 // 使用level=7，避免3/4/5是变化牌

	t.Run("33+44+5+wild (1个变化牌)", func(t *testing.T) {
		card31, _ := NewCard(3, "Spade", level)
		card32, _ := NewCard(3, "Club", level)
		card41, _ := NewCard(4, "Heart", level)
		card42, _ := NewCard(4, "Diamond", level)
		card51, _ := NewCard(5, "Spade", level)
		wildcard, _ := NewCard(7, "Heart", level) // 7是变化牌，补另一张5

		cards := []*Card{card31, card32, card41, card42, card51, wildcard}

		t.Log("33+44+5+wild (普通钢管 + 1个变化牌，补另一张5)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("33+44+5+wild 应该是合法的钢管（凑成 3-4-5）")
		}

		if comparisonKey != 3 {
			t.Errorf("comparisonKey 应该是 3 (3-4-5)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "3S,3C,4H,4D,5S,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("88+99+2wild (2个变化牌)", func(t *testing.T) {
		card81, _ := NewCard(8, "Spade", level)
		card82, _ := NewCard(8, "Club", level)
		card91, _ := NewCard(9, "Heart", level)
		card92, _ := NewCard(9, "Diamond", level)
		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)

		cards := []*Card{card81, card82, card91, card92, wild1, wild2}

		t.Log("88+99+2wild (普通钢管 + 2个变化牌，补两张10)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("88+99+2wild 应该是合法的钢管（凑成 8-9-10）")
		}

		if comparisonKey != 8 {
			t.Errorf("comparisonKey 应该是 8 (8-9-10)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "8S,8C,9H,9D,7H,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestTubeSatisfyNew_EdgeCases 测试边界钢管
func TestTubeSatisfyNew_EdgeCases(t *testing.T) {
	level := 5

	t.Run("22+33+44 (最小的普通钢管)", func(t *testing.T) {
		card21, _ := NewCard(2, "Spade", level)
		card22, _ := NewCard(2, "Club", level)
		card31, _ := NewCard(3, "Heart", level)
		card32, _ := NewCard(3, "Diamond", level)
		card41, _ := NewCard(4, "Spade", level)
		card42, _ := NewCard(4, "Club", level)

		cards := []*Card{card21, card22, card31, card32, card41, card42}

		t.Log("22+33+44 (最小的普通钢管)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
			}
		}

		if !isValid {
			t.Error("22+33+44 应该是合法的钢管")
		}

		if comparisonKey != 2 {
			t.Errorf("comparisonKey 应该是 2 (2-3-4)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "2S,2C,3H,3D,4S,4C"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("10+JJ+QQ (包含10的钢管)", func(t *testing.T) {
		card101, _ := NewCard(10, "Spade", level)
		card102, _ := NewCard(10, "Club", level)
		cardJ1, _ := NewCard(11, "Heart", level)
		cardJ2, _ := NewCard(11, "Diamond", level)
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)

		cards := []*Card{card101, card102, cardJ1, cardJ2, cardQ1, cardQ2}

		t.Log("10+JJ+QQ (包含10的钢管)")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d)", i, c.Name, c.RawNumber)
			}
		}

		if !isValid {
			t.Error("10+JJ+QQ 应该是合法的钢管")
		}

		if comparisonKey != 10 {
			t.Errorf("comparisonKey 应该是 10 (10-J-Q)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "10S,10C,11H,11D,12S,12C"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestTubeSatisfyNew_MissingPair 测试缺失一对的场景
func TestTubeSatisfyNew_MissingPair(t *testing.T) {
	level := 5

	t.Run("Q+KK+AA+2wild (缺Q)", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardK1, _ := NewCard(13, "Heart", level)
		cardK2, _ := NewCard(13, "Diamond", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Club", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{cardQ1, cardK1, cardK2, cardA1, cardA2, wildcard}

		t.Log("Q+KK+AA+1wild (缺一张Q，用变化牌补)")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			for i, c := range normalized {
				t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("Q+KK+AA+wild 应该是合法的 Q-K-A 钢管")
		}

		if comparisonKey != 12 {
			t.Errorf("comparisonKey 应该是 12 (Q-K-A)，实际: %d", comparisonKey)
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "12S,5H,13H,13D,1S,1C"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("44+6+77+wild (缺5和6)", func(t *testing.T) {
		card41, _ := NewCard(4, "Spade", level)
		card42, _ := NewCard(4, "Club", level)
		card61, _ := NewCard(6, "Heart", level)
		card71, _ := NewCard(7, "Spade", level)
		card72, _ := NewCard(7, "Club", level)
		wildcard, _ := NewCard(5, "Heart", level)

		cards := []*Card{card41, card42, card61, card71, card72, wildcard}

		t.Log("44+6+77+wild (缺失6，不能凑成钢管)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("44+6+77+wild 不应该是合法的钢管（只能补一张6，需要补两张）")
		}
	})
}

// TestTubeSatisfyNew_MoreInvalidCases 测试更多无效场景
func TestTubeSatisfyNew_MoreInvalidCases(t *testing.T) {
	level := 5

	t.Run("三对不连续", func(t *testing.T) {
		card21, _ := NewCard(2, "Spade", level)
		card22, _ := NewCard(2, "Club", level)
		card41, _ := NewCard(4, "Heart", level)
		card42, _ := NewCard(4, "Diamond", level)
		card61, _ := NewCard(6, "Spade", level)
		card62, _ := NewCard(6, "Club", level)

		cards := []*Card{card21, card22, card41, card42, card61, card62}

		t.Log("22+44+66 (三对不连续)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("22+44+66 不应该是合法的钢管（不连续）")
		}
	})

	t.Run("同数字超过2张", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardQ3, _ := NewCard(12, "Heart", level)
		cardK1, _ := NewCard(13, "Diamond", level)
		cardK2, _ := NewCard(13, "Spade", level)
		cardA1, _ := NewCard(14, "Club", level)

		cards := []*Card{cardQ1, cardQ2, cardQ3, cardK1, cardK2, cardA1}

		t.Log("QQQ+KK+A (同数字超过2张)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("QQQ+KK+A 不应该是合法的钢管（Q有3张）")
		}
	})

	t.Run("QQQ+KK+AA (3张Q)", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardQ3, _ := NewCard(12, "Heart", level)
		cardK1, _ := NewCard(13, "Diamond", level)
		cardK2, _ := NewCard(13, "Spade", level)
		cardA1, _ := NewCard(14, "Club", level)
		cardA2, _ := NewCard(14, "Heart", level)

		cards := []*Card{cardQ1, cardQ2, cardQ3, cardK1, cardK2, cardA1, cardA2}

		t.Log("QQQ+KK+AA (同数字超过2张，7张牌)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("QQQ+KK+AA 不应该是合法的钢管（长度不对且Q有3张）")
		}
	})

	t.Run("QQQQ+KK (4张Q)", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardQ3, _ := NewCard(12, "Heart", level)
		cardQ4, _ := NewCard(12, "Diamond", level)
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Club", level)

		cards := []*Card{cardQ1, cardQ2, cardQ3, cardQ4, cardK1, cardK2}

		t.Log("QQQQ+KK (4张Q)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("QQQQ+KK 不应该是合法的钢管（Q有4张）")
		}
	})

	t.Run("QQQ+KKK (两个数字都超过2张)", func(t *testing.T) {
		cardQ1, _ := NewCard(12, "Spade", level)
		cardQ2, _ := NewCard(12, "Club", level)
		cardQ3, _ := NewCard(12, "Heart", level)
		cardK1, _ := NewCard(13, "Diamond", level)
		cardK2, _ := NewCard(13, "Spade", level)
		cardK3, _ := NewCard(13, "Club", level)

		cards := []*Card{cardQ1, cardQ2, cardQ3, cardK1, cardK2, cardK3}

		t.Log("QQQ+KKK (两个数字都超过2张)")
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("QQQ+KKK 不应该是合法的钢管（Q和K都有3张）")
		}
	})

	t.Run("全是变化牌", func(t *testing.T) {
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)
		wild3, _ := NewCard(5, "Heart", level)
		wild4, _ := NewCard(5, "Heart", level)
		wild5, _ := NewCard(5, "Heart", level)
		wild6, _ := NewCard(5, "Heart", level)

		cards := []*Card{wild1, wild2, wild3, wild4, wild5, wild6}

		t.Log("6个变化牌（应该可以凑成任意钢管）")
		isValid, normalized, comparisonKey := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		// 6个变化牌可以凑成任意钢管，算法会选择最大的 Q-K-A
		if !isValid {
			t.Error("6个变化牌应该是合法的钢管（可以凑成任意组合）")
		}

		if comparisonKey != 12 {
			t.Errorf("comparisonKey 应该是 12 (选择最大的 Q-K-A)，实际: %d", comparisonKey)
		}

		// 验证 normalized 长度为6
		if len(normalized) != 6 {
			t.Errorf("normalized cards 长度应该为6，实际: %d", len(normalized))
		}

		// 6个变化牌都是5H
		actualStr := FormatCardsSimple(normalized)
		expectedStr := "5H,5H,5H,5H,5H,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}

// TestTubeSatisfyNew_MultipleChoicesMinimum 测试多组合选择最小的情况
func TestTubeSatisfyNew_MultipleChoicesMinimum(t *testing.T) {
	level := 5

	// 22+33+2wild 可以凑成 A-2-3 或 2-3-4，应该选择 2-3-4（更大）
	card21, _ := NewCard(2, "Spade", level)
	card22, _ := NewCard(2, "Club", level)
	card31, _ := NewCard(3, "Heart", level)
	card32, _ := NewCard(3, "Diamond", level)
	wild1, _ := NewCard(5, "Heart", level)
	wild2, _ := NewCard(5, "Heart", level)

	cards := []*Card{card21, card22, card31, card32, wild1, wild2}

	t.Log("=== 22+33+2wild (可以凑成 A-2-3 或 2-3-4) ===")
	isValid, normalized, comparisonKey := TubeSatisfy(cards)

	t.Logf("结果: IsValid=%v", isValid)
	if normalized != nil {
		t.Log("规范化后的牌:")
		for i, c := range normalized {
			t.Logf("  [%d]: %s (RawNumber=%d, IsWildcard=%v)", i, c.Name, c.RawNumber, c.IsWildcard())
		}

		// 验证选择的是 2-3-4（前两张应该是2）而不是 A-2-3
		if len(normalized) == 6 {
			// 如果选择了 2-3-4，前两张应该是 2
			// 如果选择了 A-2-3，前两张应该是 A (RawNumber=1)
			firstTwoAre2 := (normalized[0].RawNumber == 2) && (normalized[1].RawNumber == 2)
			firstTwoAreA := (normalized[0].RawNumber == 1) && (normalized[1].RawNumber == 1)

			if firstTwoAreA {
				t.Error("应该选择 2-3-4（更大）而不是 A-2-3")
			} else if firstTwoAre2 {
				t.Log("✓ 正确选择了 2-3-4")
			}
		}
	}

	if !isValid {
		t.Error("22+33+2wild 应该是合法的钢管")
	}

	if comparisonKey != 2 {
		t.Errorf("comparisonKey 应该是 2 (选择 2-3-4 而非 A-2-3)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "2S,2C,3H,3D,5H,5H"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期 (应该选择 2-3-4 而不是 A-2-3)\n期望: %s\n实际: %s", expectedStr, actualStr)
	}
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
	wild2, _ := NewCard(5, "Heart", level) // 两张红桃5都是变化牌

	cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, wild1, wild2}

	t.Log("=== QQ+KK+2wild (应该选择 Q-K-A 而不是 J-Q-K) ===")
	isValid, normalized, comparisonKey := TubeSatisfy(cards)

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

	if comparisonKey != 12 {
		t.Errorf("comparisonKey 应该是 12 (选择 Q-K-A 而非 J-Q-K)，实际: %d", comparisonKey)
	}

	actualStr := FormatCardsSimple(normalized)
	expectedStr := "12S,12C,13H,13D,5H,5H"
	if actualStr != expectedStr {
		t.Errorf("normalized cards 不符合预期 (应该选择 Q-K-A 而不是 J-Q-K)\n期望: %s\n实际: %s", expectedStr, actualStr)
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
		wild2, _ := NewCard(5, "Heart", level) // 两张红桃5都是变化牌

		cards := []*Card{cardQ1, cardQ2, card91, card92, wild1, wild2}

		t.Log("QQ+99+2wild (间隔太大，无法凑成)")
		isValid, _, _ := TubeSatisfy(cards)

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
		isValid, _, _ := TubeSatisfy(cards)

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
		isValid, _, _ := TubeSatisfy(cards)

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
		isValid, _, _ := TubeSatisfy(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("含小王的牌组不应该是合法的钢管")
		}
	})
}
