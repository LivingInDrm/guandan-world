package sdk

import (
	"testing"
)

// TestFullHouseSatisfyNew_Invalid 测试不合法的葫芦场景
func TestFullHouseSatisfyNew_Invalid(t *testing.T) {
	t.Run("牌张数错误（6张）", func(t *testing.T) {
		level := 5

		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardA3, _ := NewCard(14, "Club", level)
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Heart", level)
		cardQ1, _ := NewCard(12, "Spade", level)

		cards := []*Card{cardA1, cardA2, cardA3, cardK1, cardK2, cardQ1}

		t.Log("=== 6张牌（应该是5张）===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("6张牌不应该是合法的葫芦")
		}

		if len(normalized) != 6 {
			t.Errorf("返回的牌数应该是6，实际: %d", len(normalized))
		}
	})

	t.Run("不同joker（大王+小王）", func(t *testing.T) {
		level := 5

		bigJoker, _ := NewCard(16, "Joker", level)
		smallJoker, _ := NewCard(15, "Joker", level)
		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)
		card3, _ := NewCard(3, "Spade", level)

		cards := []*Card{bigJoker, smallJoker, wild1, wild2, card3}

		t.Log("=== 1大王+1小王+2wildcard+1普通牌 ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, _ := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("大王+小王（不同number的joker）不应该是合法的葫芦")
		}
	})

	t.Run("2张大王但无法构成triple", func(t *testing.T) {
		level := 5

		bigJoker1, _ := NewCard(16, "Joker", level)
		bigJoker2, _ := NewCard(16, "Joker", level)
		wild1, _ := NewCard(5, "Heart", level)
		cardA, _ := NewCard(14, "Spade", level)
		cardK, _ := NewCard(13, "Spade", level)

		cards := []*Card{bigJoker1, bigJoker2, wild1, cardA, cardK}

		t.Log("=== 2张大王+1wildcard+2个不同single ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, _ := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2张大王+1wildcard+A+K 不应该是合法的葫芦（剩余3张有2种number）")
		}
	})

	t.Run("2张小王+1对A+1个K", func(t *testing.T) {
		level := 5

		smallJoker1, _ := NewCard(15, "Joker", level)
		smallJoker2, _ := NewCard(15, "Joker", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardK, _ := NewCard(13, "Spade", level)

		cards := []*Card{smallJoker1, smallJoker2, cardA1, cardA2, cardK}

		t.Log("=== 2张小王+1对A+1个K ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, _ := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2张小王+1对A+1个K 不应该是合法的葫芦（剩余3张有2种number）")
		}
	})

	t.Run("只有1种number（会形成炸弹）", func(t *testing.T) {
		level := 5

		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardA3, _ := NewCard(14, "Club", level)

		cards := []*Card{wild1, wild2, cardA1, cardA2, cardA3}

		t.Log("=== 2个wildcard+3个同number的A ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, _ := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)

		if isValid {
			t.Error("2wildcard+3个A 不应该是合法的葫芦（只有1种number，应该是炸弹）")
		}
	})
}

// TestFullHouseSatisfyNew_Valid 测试合法的葫芦场景
func TestFullHouseSatisfyNew_Valid(t *testing.T) {
	t.Run("2wild+1对A+1个K", func(t *testing.T) {
		level := 5

		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardK, _ := NewCard(13, "Spade", level)

		cards := []*Card{wild1, wild2, cardA1, cardA2, cardK}

		t.Log("=== 2wild+1对A+1个K ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("2wild+1对A+1个K 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1H,5H,13S,5H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("2wild+1对K+1个A", func(t *testing.T) {
		level := 5

		wild1, _ := NewCard(5, "Heart", level)
		wild2, _ := NewCard(5, "Heart", level)
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Heart", level)
		cardA, _ := NewCard(14, "Spade", level)

		cards := []*Card{wild1, wild2, cardK1, cardK2, cardA}

		t.Log("=== 2wild+1对K+1个A ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("2wild+1对K+1个A 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,5H,5H,13S,13H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("2wild+1个level card+1对A", func(t *testing.T) {
		level := 7

		wild1, _ := NewCard(7, "Heart", level)
		wild2, _ := NewCard(7, "Heart", level)
		card3, _ := NewCard(3, "Spade", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)

		cards := []*Card{wild1, wild2, card3, cardA1, cardA2}

		t.Log("=== 2wild+1个3+1对A ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("2wild+1个3+1对A 应该是合法的葫芦（triple=A，更大）")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1H,7H,3S,7H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("1wild+1对A+1对K", func(t *testing.T) {
		level := 5

		wild1, _ := NewCard(5, "Heart", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardK1, _ := NewCard(13, "Spade", level)
		cardK2, _ := NewCard(13, "Heart", level)

		cards := []*Card{wild1, cardA1, cardA2, cardK1, cardK2}

		t.Log("=== 1wild+1对A+1对K ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("1wild+1对A+1对K 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1H,5H,13S,13H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("1wild+1对level card+1对A", func(t *testing.T) {
		level := 7

		wild1, _ := NewCard(7, "Heart", level)
		card3_1, _ := NewCard(7, "Spade", level)
		card3_2, _ := NewCard(7, "Club", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)

		cards := []*Card{wild1, card3_1, card3_2, cardA1, cardA2}

		t.Log("=== 1wild+1对level+1对A ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("1wild+1对level+1对A 应该是合法的葫芦（triple=A，更大）")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,7C,7H,1S,1H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("1wild+1对A+1对小王", func(t *testing.T) {
		level := 5

		wild1, _ := NewCard(5, "Heart", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		smallJoker1, _ := NewCard(15, "Joker", level)
		smallJoker2, _ := NewCard(15, "Joker", level)

		cards := []*Card{wild1, cardA1, cardA2, smallJoker1, smallJoker2}

		t.Log("=== 1wild+1对A+1对小王 ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("1wild+1对A+1对小王 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1H,5H,SJ,SJ"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("0wild+3个A+2个3", func(t *testing.T) {
		level := 5

		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)
		cardA3, _ := NewCard(14, "Club", level)
		card3_1, _ := NewCard(3, "Spade", level)
		card3_2, _ := NewCard(3, "Heart", level)

		cards := []*Card{cardA1, cardA2, cardA3, card3_1, card3_2}

		t.Log("=== 0wild+3个A+2个3 ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("0wild+3个A+2个3 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "1S,1C,1H,3S,3H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("0wild+3个level card+2个A", func(t *testing.T) {
		level := 7

		card7_1, _ := NewCard(7, "Spade", level)
		card7_2, _ := NewCard(7, "Club", level)
		card7_3, _ := NewCard(7, "Diamond", level)
		cardA1, _ := NewCard(14, "Spade", level)
		cardA2, _ := NewCard(14, "Heart", level)

		cards := []*Card{card7_1, card7_2, card7_3, cardA1, cardA2}

		t.Log("=== 0wild+3个7（level card）+2个A ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("0wild+3个7+2个A 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,7C,7D,1S,1H"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})

	t.Run("0wild+3个level card+2个大王", func(t *testing.T) {
		level := 7

		card7_1, _ := NewCard(7, "Spade", level)
		card7_2, _ := NewCard(7, "Club", level)
		card7_3, _ := NewCard(7, "Diamond", level)
		bigJoker1, _ := NewCard(16, "Joker", level)
		bigJoker2, _ := NewCard(16, "Joker", level)

		cards := []*Card{card7_1, card7_2, card7_3, bigJoker1, bigJoker2}

		t.Log("=== 0wild+3个7（level card）+2个大王 ===")
		for i, c := range cards {
			t.Logf("  牌[%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
		}

		isValid, normalized := FullHouseSatisfyNew(cards)

		t.Logf("结果: IsValid=%v", isValid)
		if normalized != nil {
			t.Log("规范化后的牌:")
			for i, c := range normalized {
				t.Logf("  [%d]: %s (Number=%d, IsWildcard=%v)", i, c.Name, c.Number, c.IsWildcard())
			}
		}

		if !isValid {
			t.Error("0wild+3个7+2个大王 应该是合法的葫芦")
		}

		actualStr := FormatCardsSimple(normalized)
		expectedStr := "7S,7C,7D,BJ,BJ"
		if actualStr != expectedStr {
			t.Errorf("normalized cards 不符合预期\n期望: %s\n实际: %s", expectedStr, actualStr)
		}
	})
}
