package sdk

import (
	"fmt"
	"testing"
)

// 测试分析 A 的循环特性和排序行为
func TestACycleAnalysis(t *testing.T) {
	level := 5
	
	fmt.Println("========================================")
	fmt.Println("A 的双重特性分析")
	fmt.Println("========================================")
	fmt.Println()
	
	// 创建 A
	cardA1, _ := NewCard(14, "Spade", level)
	
	fmt.Printf("Ace 的数字属性：\n")
	fmt.Printf("  Number:    %d (用于大小比较，A最大)\n", cardA1.Number)
	fmt.Printf("  RawNumber: %d (用于连续性判断，A可以接在K后或在2前)\n\n", cardA1.RawNumber)
	
	// ========== 测试 1: A-2-3 钢板（A在低位） ==========
	fmt.Println("========== 测试 1: A-2-3 钢板（A在低位） ==========")
	cardA3, _ := NewCard(14, "Spade", level)
	cardA4, _ := NewCard(14, "Heart", level)
	cardA5, _ := NewCard(14, "Club", level)
	card2_1, _ := NewCard(2, "Spade", level)
	card2_2, _ := NewCard(2, "Heart", level)
	card2_3, _ := NewCard(2, "Club", level)
	
	a2Cards := []*Card{cardA3, cardA4, cardA5, card2_1, card2_2, card2_3}
	
	fmt.Println("输入牌组: AAA + 222")
	for i, c := range a2Cards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
	}
	
	a2Plate := NewPlate(a2Cards)
	
	fmt.Printf("\n排序后 (sortCardsForConsecutive):\n")
	if a2Plate.Cards != nil {
		for i, c := range a2Plate.Cards {
			fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
		}
	}
	
	fmt.Printf("\n验证结果: IsValid=%v\n", a2Plate.IsValid())
	fmt.Printf("逻辑: A(RawNumber=1) 排在最前，2(RawNumber=2) 排在后面\n")
	fmt.Printf("      检测到 card1Num=1, card2Num=2，匹配 A-2 特殊规则\n")
	fmt.Printf("      代码位置: comp.go:1367-1370\n\n")
	
	// ========== 测试 2: K-A 钢板（A在高位） ==========
	fmt.Println("========== 测试 2: K-A 钢板（A在高位） ==========")
	cardK1, _ := NewCard(13, "Spade", level)
	cardK2, _ := NewCard(13, "Heart", level)
	cardK3, _ := NewCard(13, "Club", level)
	cardA6, _ := NewCard(14, "Spade", level)
	cardA7, _ := NewCard(14, "Heart", level)
	cardA8, _ := NewCard(14, "Club", level)
	
	kaCards := []*Card{cardK1, cardK2, cardK3, cardA6, cardA7, cardA8}
	
	fmt.Println("输入牌组: KKK + AAA")
	for i, c := range kaCards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
	}
	
	kaPlate := NewPlate(kaCards)
	
	fmt.Printf("\n排序后 (sortCardsForConsecutive):\n")
	if kaPlate.Cards != nil {
		for i, c := range kaPlate.Cards {
			fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
		}
	}
	
	fmt.Printf("\n验证结果: IsValid=%v\n", kaPlate.IsValid())
	fmt.Printf("逻辑: A(RawNumber=1) 排在最前，K(RawNumber=13) 排在后面\n")
	fmt.Printf("      检测到 card1Num=1, card2Num=13，匹配 A-K 特殊规则\n")
	fmt.Printf("      代码位置: comp.go:1377-1380\n\n")
	
	// ========== 测试 3: Q-K-A 钢管（应该合法但当前失败） ==========
	fmt.Println("========== 测试 3: Q-K-A 钢管（A在高位，6张对子） ==========")
	cardQ1, _ := NewCard(12, "Spade", level)
	cardQ2, _ := NewCard(12, "Heart", level)
	cardK4, _ := NewCard(13, "Club", level)
	cardK5, _ := NewCard(13, "Diamond", level)
	cardA9, _ := NewCard(14, "Spade", level)
	cardA10, _ := NewCard(14, "Club", level)
	
	qkaCards := []*Card{cardQ1, cardQ2, cardK4, cardK5, cardA9, cardA10}
	
	fmt.Println("输入牌组: QQ + KK + AA")
	for i, c := range qkaCards {
		fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
	}
	
	qkaTube := NewTube(qkaCards)
	
	fmt.Printf("\n排序后 (sortCardsForConsecutive):\n")
	if qkaTube.Cards != nil {
		for i, c := range qkaTube.Cards {
			fmt.Printf("  [%d] %s: Number=%d, RawNumber=%d\n", i, c.Name, c.Number, c.RawNumber)
		}
	}
	
	fmt.Printf("\n验证结果: IsValid=%v ❌ 错误！\n", qkaTube.IsValid())
	fmt.Printf("问题分析:\n")
	fmt.Printf("  1. 排序后: A(1), A(1), Q(12), Q(12), K(13), K(13)\n")
	fmt.Printf("  2. cardNumbers: [1, 1, 12, 12, 13, 13]\n")
	fmt.Printf("  3. uniqueNumbers: {1, 12, 13} - 3个唯一数字 ✓\n")
	fmt.Printf("  4. computeRelativeDiffs([1,1,12,12,13,13], 6):\n")
	fmt.Printf("     基准值 = cardNumbers[0] = 1\n")
	fmt.Printf("     结果 = [1-1, 1-1, 12-1, 12-1, 13-1, 13-1]\n")
	fmt.Printf("          = [0, 0, 11, 11, 12, 12]\n")
	fmt.Printf("  5. TUBE_PATTERN_TRIPLET = [0, 0, 1, 1, 2, 2]\n")
	fmt.Printf("  6. [0,0,11,11,12,12] ≠ [0,0,1,1,2,2] ❌ 不匹配\n")
	fmt.Printf("  7. 缺少特殊处理：应该像钢板一样检测 1,12,13 组合\n\n")
	
	// ========== 对比: J-Q-K 钢管（正常工作） ==========
	fmt.Println("========== 对比: J-Q-K 钢管（正常连续） ==========")
	cardJ1, _ := NewCard(11, "Spade", level)
	cardJ2, _ := NewCard(11, "Heart", level)
	cardQ3, _ := NewCard(12, "Club", level)
	cardQ4, _ := NewCard(12, "Diamond", level)
	cardK6, _ := NewCard(13, "Spade", level)
	cardK7, _ := NewCard(13, "Club", level)
	
	jqkCards := []*Card{cardJ1, cardJ2, cardQ3, cardQ4, cardK6, cardK7}
	jqkTube := NewTube(jqkCards)
	
	fmt.Println("输入牌组: JJ + QQ + KK")
	fmt.Printf("\n排序后: J(11), J(11), Q(12), Q(12), K(13), K(13)\n")
	fmt.Printf("cardNumbers: [11, 11, 12, 12, 13, 13]\n")
	fmt.Printf("computeRelativeDiffs: [0, 0, 1, 1, 2, 2]\n")
	fmt.Printf("匹配 TUBE_PATTERN_TRIPLET ✓\n")
	fmt.Printf("验证结果: IsValid=%v\n\n", jqkTube.IsValid())
	
	// ========== 设计原理总结 ==========
	fmt.Println("========================================")
	fmt.Println("设计原理总结")
	fmt.Println("========================================")
	fmt.Println("1. RawNumber 的双重性:")
	fmt.Println("   - Number=14 (A最大，用于牌力比较)")
	fmt.Println("   - RawNumber=1 (A可在最低位，用于连续性)")
	fmt.Println("")
	fmt.Println("2. 循环的本质:")
	fmt.Println("   - A-2循环: A接在2前面 (A-2-3-4-5)")
	fmt.Println("   - K-A循环: A接在K后面 (J-Q-K-A 或 Q-K-A)")
	fmt.Println("")
	fmt.Println("3. 排序导致的问题:")
	fmt.Println("   - sortCardsForConsecutive按RawNumber排序")
	fmt.Println("   - A(RawNumber=1)总是排在最前面")
	fmt.Println("   - 破坏了K-A的视觉顺序")
	fmt.Println("")
	fmt.Println("4. 特殊处理的必要性:")
	fmt.Println("   - 钢板: 有 A-2 和 K-A 特殊检测 ✓")
	fmt.Println("   - 钢管: 缺少 Q-K-A 特殊检测 ❌")
	fmt.Println("   - 顺子: 有 10-J-Q-K-A 特殊重排 ✓")
	fmt.Println("")
	fmt.Println("5. 修复方案:")
	fmt.Println("   在 tubeSatisfy 的 wildcardCount==0 分支中:")
	fmt.Println("   添加对 uniqueNumbers={1,12,13} 的检测")
	fmt.Println("   识别为 Q-K-A 钢管并返回 true")
}
