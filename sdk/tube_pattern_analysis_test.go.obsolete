package sdk

import (
	"fmt"
	"testing"
)

// 详细分析1个变化牌时的模式匹配
func TestTubeWildcardPatternAnalysis(t *testing.T) {
	level := 5
	
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("钢管 1个变化牌 模式匹配详细分析")
	fmt.Println("==================================================")
	fmt.Println()
	
	// ========== 测试 1: A-2-3 钢管 ==========
	fmt.Println("========== 测试 1: A-2-3 钢管 + 1个变化牌 ==========")
	
	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)
	card2_1, _ := NewCard(2, "Club", level)
	card2_2, _ := NewCard(2, "Diamond", level)
	card3, _ := NewCard(3, "Spade", level)
	wildcard1, _ := NewCard(5, "Heart", level)
	
	a23Cards := []*Card{cardA1, cardA2, card2_1, card2_2, card3, wildcard1}
	sortedA23 := sortCardsForConsecutive(a23Cards)
	
	fmt.Println("排序后: A(1), A(1), 2(2), 2(2), 3(3), 变化牌(-1)")
	
	cardNumbers := make([]int, len(sortedA23))
	for i, card := range sortedA23 {
		if card.IsWildcard() {
			cardNumbers[i] = -1
		} else {
			cardNumbers[i] = card.RawNumber
		}
	}
	fmt.Printf("cardNumbers: %v\n", cardNumbers)
	
	firstFive := computeRelativeDiffs(cardNumbers, 5)
	fmt.Printf("\ncomputeRelativeDiffs (前5个):\n")
	fmt.Printf("  基准值 = cardNumbers[0] = %d\n", cardNumbers[0])
	fmt.Printf("  firstFive = %v\n", firstFive)
	
	fmt.Printf("\n匹配 TUBE_PATTERN_0112 = %v:\n", TUBE_PATTERN_0112)
	match0112 := matchesPattern(firstFive, TUBE_PATTERN_0112)
	fmt.Printf("  结果: %v\n", match0112)
	if match0112 {
		fmt.Println("  ✅ 匹配！这就是为什么 A-2-3 被识别")
		fmt.Println("  意义: i(1), i(1), i+1(2), i+1(2), i+2(3)")
	}
	
	fmt.Printf("\n匹配 TUBE_PATTERN_0122 = %v:\n", TUBE_PATTERN_0122)
	match0122 := matchesPattern(firstFive, TUBE_PATTERN_0122)
	fmt.Printf("  结果: %v\n", match0122)
	
	fmt.Printf("\n匹配 TUBE_PATTERN_1122 = %v:\n", TUBE_PATTERN_1122)
	match1122 := matchesPattern(firstFive, TUBE_PATTERN_1122)
	fmt.Printf("  结果: %v\n", match1122)
	
	// ========== 测试 2: Q-K-A 钢管 ==========
	fmt.Println()
	fmt.Println("========== 测试 2: Q-K-A 钢管 + 1个变化牌 ==========")
	
	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Club", level)
	cardK1, _ := NewCard(13, "Club", level)
	cardK2, _ := NewCard(13, "Diamond", level)
	cardA_qka, _ := NewCard(14, "Spade", level)
	wildcard2, _ := NewCard(5, "Heart", level)
	
	// 注意：钢管是6张牌，所以应该是 QQ + KK + A + 变化牌 = 6张（变化牌充当另一个A）
	qkaCards := []*Card{cardQ1, cardQ2, cardK1, cardK2, cardA_qka, wildcard2}
	sortedQKA := sortCardsForConsecutive(qkaCards)
	
	fmt.Println("排序后: (只有1个A + 1个变化牌充当A)")
	
	cardNumbers2 := make([]int, len(sortedQKA))
	for i, card := range sortedQKA {
		if card.IsWildcard() {
			cardNumbers2[i] = -1
		} else {
			cardNumbers2[i] = card.RawNumber
		}
	}
	fmt.Printf("cardNumbers: %v\n", cardNumbers2)
	
	firstFive2 := computeRelativeDiffs(cardNumbers2, 5)
	fmt.Printf("\ncomputeRelativeDiffs (前5个):\n")
	fmt.Printf("  基准值 = cardNumbers[0] = %d\n", cardNumbers2[0])
	fmt.Printf("  firstFive = %v\n", firstFive2)
	
	fmt.Printf("\n匹配 TUBE_PATTERN_0112 = %v:\n", TUBE_PATTERN_0112)
	match0112_2 := matchesPattern(firstFive2, TUBE_PATTERN_0112)
	fmt.Printf("  结果: %v\n", match0112_2)
	
	fmt.Printf("\n匹配 TUBE_PATTERN_0122 = %v:\n", TUBE_PATTERN_0122)
	match0122_2 := matchesPattern(firstFive2, TUBE_PATTERN_0122)
	fmt.Printf("  结果: %v\n", match0122_2)
	
	fmt.Printf("\n匹配 TUBE_PATTERN_1122 = %v:\n", TUBE_PATTERN_1122)
	match1122_2 := matchesPattern(firstFive2, TUBE_PATTERN_1122)
	fmt.Printf("  结果: %v\n", match1122_2)
	
	if !match0112_2 && !match0122_2 && !match1122_2 {
		fmt.Println("\n  ❌ 所有模式都不匹配！这就是为什么 Q-K-A 不被识别")
		fmt.Println("  原因: A 的 RawNumber=1 导致排序后破坏了连续性")
		fmt.Println("  需要: 像顺子一样的特殊 K-A 循环检测")
	}
	
	// ========== 分析为什么需要特殊处理 ==========
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("总结：为什么需要特殊处理")
	fmt.Println("==================================================")
	fmt.Println()
	fmt.Println("A-2-3 钢管:")
	fmt.Println("  排序后: A(1), A(1), 2(2), 2(2), 3(3), wild")
	fmt.Println("  相对差: [0, 0, 1, 1, 2] ← 匹配 TUBE_PATTERN_0112")
	fmt.Println("  结论: ✅ 通过通用模式识别（不需要特殊处理）")
	fmt.Println()
	fmt.Println("Q-K-A 钢管:")
	fmt.Println("  排序后: Q(12), Q(12), K(13), K(13), A(1), wild")
	fmt.Println("  相对差: [0, 0, 1, 1, -11] ← 不匹配任何模式")
	fmt.Println("  结论: ❌ 需要特殊的 K-A 循环检测（类似顺子）")
	fmt.Println()
	fmt.Println("设计问题:")
	fmt.Println("  A 的 RawNumber=1 是为了支持 A-2-3 这样的低位循环")
	fmt.Println("  但在 Q-K-A 高位循环中，A 应该被当作 14")
	fmt.Println("  解决方案: 在模式匹配前，检测 {12,13,1} 组合并特殊处理")
}
