package sdk

import (
	"fmt"
	"testing"
)

// 测试钢管在有变化牌时对 A-2 和 K-A 循环的处理
func TestTubeWithWildcardACycle(t *testing.T) {
	
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("钢管（Tube）变化牌情况下的 A-2 和 K-A 循环处理")
	fmt.Println("==================================================")
	fmt.Println()
	
	// ========== 测试 1: A-2-3 钢管 + 1个变化牌 ==========
	fmt.Println("========== 测试 1: A-2-3 钢管 + 1个变化牌 ==========")
	level := 5
	
	cardA1, _ := NewCard(14, "Spade", level)   // A♠
	cardA2, _ := NewCard(14, "Club", level)    // A♣
	card2_1, _ := NewCard(2, "Heart", level)   // 2♥
	card2_2, _ := NewCard(2, "Diamond", level) // 2♦
	card3_1, _ := NewCard(3, "Spade", level)   // 3♠
	wildcard1, _ := NewCard(5, "Heart", level) // 5♥ 变化牌
	
	a23Cards := []*Card{cardA1, cardA2, card2_1, card2_2, card3_1, wildcard1}
	
	fmt.Println("输入: AA + 22 + 3 + 变化牌(5♥)")
	for i, c := range a23Cards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d, IsWildcard=%v\n",
			i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}
	
	a23Tube := NewTube(a23Cards)
	
	fmt.Printf("\n排序后:\n")
	if a23Tube.Cards != nil {
		for i, c := range a23Tube.Cards {
			fmt.Printf("  [%d] %s: RawNumber=%d, IsWildcard=%v\n",
				i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}
	
	fmt.Printf("\nIsValid: %v\n", a23Tube.IsValid())
	if !a23Tube.IsValid() {
		fmt.Println("❌ A-2-3 钢管未被识别（缺少 A-2 循环处理）")
	} else {
		fmt.Println("✅ A-2-3 钢管被识别")
	}
	fmt.Println()
	
	// ========== 测试 2: Q-K-A 钢管 + 1个变化牌 ==========
	fmt.Println("========== 测试 2: Q-K-A 钢管 + 1个变化牌 ==========")
	
	cardQ1, _ := NewCard(12, "Spade", level)   // Q♠
	cardQ2, _ := NewCard(12, "Club", level)    // Q♣
	cardK1, _ := NewCard(13, "Heart", level)   // K♥
	cardK2, _ := NewCard(13, "Diamond", level) // K♦
	cardA3, _ := NewCard(14, "Spade", level)   // A♠
	wildcard2, _ := NewCard(5, "Heart", level) // 5♥ 变化牌
	
	qkaCards := []*Card{cardQ1, cardQ2, cardK1, cardK2, cardA3, wildcard2}
	
	fmt.Println("输入: QQ + KK + A + 变化牌(5♥)")
	for i, c := range qkaCards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d, IsWildcard=%v\n",
			i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}
	
	qkaTube := NewTube(qkaCards)
	
	fmt.Printf("\n排序后:\n")
	if qkaTube.Cards != nil {
		for i, c := range qkaTube.Cards {
			fmt.Printf("  [%d] %s: RawNumber=%d, IsWildcard=%v\n",
				i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}
	
	fmt.Printf("\nIsValid: %v\n", qkaTube.IsValid())
	if !qkaTube.IsValid() {
		fmt.Println("❌ Q-K-A 钢管未被识别（缺少 K-A 循环处理）")
	} else {
		fmt.Println("✅ Q-K-A 钢管被识别")
	}
	fmt.Println()
	
	// ========== 测试 3: A-K 钢管 + 2个变化牌 ==========
	fmt.Println("========== 测试 3: A-K 钢管 + 2个变化牌 ==========")
	
	cardA4, _ := NewCard(14, "Spade", level)   // A♠
	cardA5, _ := NewCard(14, "Club", level)    // A♣
	cardK3, _ := NewCard(13, "Club", level)    // K♣ (不是红心，避免混淆)
	cardK4, _ := NewCard(13, "Diamond", level) // K♦
	wildcard3, _ := NewCard(5, "Heart", level) // 5♥ 变化牌 (红心+level)
	wildcard4, _ := NewCard(5, "Heart", level) // 5♥ 变化牌 (需要两张红心5)
	
	akCards := []*Card{cardA4, cardA5, cardK3, cardK4, wildcard3, wildcard4}
	
	fmt.Println("输入: AA + KK + 2个变化牌(5♥5♦)")
	for i, c := range akCards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d, IsWildcard=%v\n",
			i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}
	
	akTube := NewTube(akCards)
	
	fmt.Printf("\n排序后:\n")
	if akTube.Cards != nil {
		for i, c := range akTube.Cards {
			fmt.Printf("  [%d] %s: RawNumber=%d, IsWildcard=%v\n",
				i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}
	
	fmt.Printf("\nIsValid: %v\n", akTube.IsValid())
	if akTube.IsValid() {
		fmt.Println("✅ A-K 钢管被识别（有特殊处理！）")
		fmt.Println("   代码位置: comp.go:1637-1651")
	} else {
		fmt.Println("❌ A-K 钢管未被识别")
	}
	fmt.Println()
	
	// ========== 测试 4: A-2 钢管 + 2个变化牌 ==========
	fmt.Println("========== 测试 4: A-2 钢管 + 2个变化牌 ==========")
	
	cardA6, _ := NewCard(14, "Spade", level)   // A♠
	cardA7, _ := NewCard(14, "Club", level)    // A♣
	card2_3, _ := NewCard(2, "Club", level)    // 2♣ (不是红心)
	card2_4, _ := NewCard(2, "Diamond", level) // 2♦
	wildcard5, _ := NewCard(5, "Heart", level) // 5♥ 变化牌
	wildcard6, _ := NewCard(5, "Heart", level) // 5♥ 变化牌
	
	a2Cards := []*Card{cardA6, cardA7, card2_3, card2_4, wildcard5, wildcard6}
	
	fmt.Println("输入: AA + 22 + 2个变化牌(5♥5♦)")
	for i, c := range a2Cards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d, IsWildcard=%v\n",
			i, c.Name, c.Number, c.RawNumber, c.IsWildcard())
	}
	
	a2Tube := NewTube(a2Cards)
	
	fmt.Printf("\n排序后:\n")
	if a2Tube.Cards != nil {
		for i, c := range a2Tube.Cards {
			fmt.Printf("  [%d] %s: RawNumber=%d, IsWildcard=%v\n",
				i, c.Name, c.RawNumber, c.IsWildcard())
		}
	}
	
	fmt.Printf("\nIsValid: %v\n", a2Tube.IsValid())
	if !a2Tube.IsValid() {
		fmt.Println("❌ A-2 钢管未被识别（缺少 A-2 循环处理）")
	} else {
		fmt.Println("✅ A-2 钢管被识别")
	}
	fmt.Println()
	
	fmt.Println("==================================================")
	fmt.Println("总结：钢管的变化牌循环处理情况")
	fmt.Println("==================================================")
	fmt.Println("1个变化牌:")
	fmt.Println("  - A-2-3 钢管: ?")
	fmt.Println("  - Q-K-A 钢管: ?")
	fmt.Println()
	fmt.Println("2个变化牌:")
	fmt.Println("  - A-2-3 钢管: ?")
	fmt.Println("  - A-K 钢管: ✅ 有特殊处理 (comp.go:1637-1651)")
	fmt.Println()
}
