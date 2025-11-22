package sdk

import (
	"fmt"
	"testing"
)

func TestSortCardsForConsecutiveBehavior(t *testing.T) {
	level := 5
	
	fmt.Println()
	fmt.Println("========== 测试 sortCardsForConsecutive 的实际行为 ==========")
	fmt.Println()
	
	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Heart", level)
	cardK1, _ := NewCard(13, "Club", level)
	cardK2, _ := NewCard(13, "Diamond", level)
	cardA1, _ := NewCard(14, "Spade", level)
	cardA2, _ := NewCard(14, "Club", level)
	
	cards := []*Card{cardQ1, cardQ2, cardK1, cardK2, cardA1, cardA2}
	
	fmt.Println("输入顺序: QQ, KK, AA")
	for i, c := range cards {
		fmt.Printf("  [%d] %s: RawNumber=%d\n", i, c.Name, c.RawNumber)
	}
	
	sorted := sortCardsForConsecutive(cards)
	
	fmt.Println("\nsortCardsForConsecutive 排序后:")
	for i, c := range sorted {
		fmt.Printf("  [%d] %s: RawNumber=%d\n", i, c.Name, c.RawNumber)
	}
	
	fmt.Println("\n提取 cardNumbers 数组:")
	cardNumbers := make([]int, len(sorted))
	for i, card := range sorted {
		cardNumbers[i] = card.RawNumber
	}
	fmt.Printf("  cardNumbers = %v\n", cardNumbers)
	
	fmt.Println("\n获取 uniqueNumbers:")
	uniqueNumbers := make(map[int]bool)
	for _, num := range cardNumbers {
		uniqueNumbers[num] = true
	}
	fmt.Printf("  uniqueNumbers = %v\n", uniqueNumbers)
	fmt.Printf("  len(uniqueNumbers) = %d\n", len(uniqueNumbers))
	
	fmt.Println("\ncomputeRelativeDiffs:")
	diffs := computeRelativeDiffs(cardNumbers, len(cardNumbers))
	fmt.Printf("  基准值 = cardNumbers[0] = %d\n", cardNumbers[0])
	fmt.Printf("  diffs = %v\n", diffs)
	
	fmt.Println("\nTUBE_PATTERN_TRIPLET:")
	fmt.Printf("  pattern = %v\n", TUBE_PATTERN_TRIPLET)
	
	fmt.Println("\n匹配结果:")
	matches := matchesPattern(diffs, TUBE_PATTERN_TRIPLET)
	fmt.Printf("  matchesPattern(diffs, TUBE_PATTERN_TRIPLET) = %v\n", matches)
	
	if !matches {
		fmt.Println("\n❌ 不匹配！这就是为什么 QQ,KK,AA 无法被识别为钢管")
	}
}
