package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	level := 5
	fmt.Println("=== 测试 StraightFlush 升级 ===\n")

	// 测试1: 同花顺 10-J-Q-K-A (黑桃)
	fmt.Println("测试1: 黑桃 10-J-Q-K-A 同花顺")
	card10, _ := sdk.NewCard(10, "Spade", level)
	cardJ, _ := sdk.NewCard(11, "Spade", level)
	cardQ, _ := sdk.NewCard(12, "Spade", level)
	cardK, _ := sdk.NewCard(13, "Spade", level)
	cardA, _ := sdk.NewCard(14, "Spade", level)
	
	sf1 := sdk.NewStraightFlush([]*sdk.Card{card10, cardJ, cardQ, cardK, cardA})
	fmt.Printf("  Valid: %v\n", sf1.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 10)\n", sf1.ComparisonKey)
	
	if !sf1.IsValid() || sf1.ComparisonKey != 10 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试2: 同花顺 A-2-3-4-5 (红心)
	fmt.Println("测试2: 红心 A-2-3-4-5 同花顺")
	level2 := 7
	cardA2, _ := sdk.NewCard(14, "Heart", level2)
	card2, _ := sdk.NewCard(2, "Heart", level2)
	card3, _ := sdk.NewCard(3, "Heart", level2)
	card4, _ := sdk.NewCard(4, "Heart", level2)
	card5, _ := sdk.NewCard(5, "Heart", level2)
	
	sf2 := sdk.NewStraightFlush([]*sdk.Card{cardA2, card2, card3, card4, card5})
	fmt.Printf("  Valid: %v\n", sf2.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 1)\n", sf2.ComparisonKey)
	
	if !sf2.IsValid() || sf2.ComparisonKey != 1 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试3: 比较逻辑
	fmt.Println("测试3: 同花顺比较")
	fmt.Printf("  黑桃10-J-Q-K-A > 红心A-2-3-4-5: %v (期望: true)\n", sf1.GreaterThan(sf2))
	fmt.Printf("  红心A-2-3-4-5 > 黑桃10-J-Q-K-A: %v (期望: false)\n", sf2.GreaterThan(sf1))
	
	if !sf1.GreaterThan(sf2) || sf2.GreaterThan(sf1) {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试4: 非同花顺（花色不一致）
	fmt.Println("测试4: 非同花顺 (10黑-J红-Q黑-K红-A黑)")
	card10b, _ := sdk.NewCard(10, "Spade", level)
	cardJb, _ := sdk.NewCard(11, "Heart", level)
	cardQb, _ := sdk.NewCard(12, "Spade", level)
	cardKb, _ := sdk.NewCard(13, "Heart", level)
	cardAb, _ := sdk.NewCard(14, "Spade", level)
	
	notSF := sdk.NewStraightFlush([]*sdk.Card{card10b, cardJb, cardQb, cardKb, cardAb})
	fmt.Printf("  Valid: %v (期望: false)\n", notSF.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 0)\n", notSF.ComparisonKey)
	
	if notSF.IsValid() {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试5: 含万能牌的同花顺
	fmt.Println("测试5: 含万能牌的同花顺 (黑桃10-J-K-A+wild)")
	card10c, _ := sdk.NewCard(10, "Spade", level)
	cardJc, _ := sdk.NewCard(11, "Spade", level)
	cardKc, _ := sdk.NewCard(13, "Spade", level)
	cardAc, _ := sdk.NewCard(14, "Spade", level)
	wildcard, _ := sdk.NewCard(5, "Heart", level) // 5是变化牌
	
	sf3 := sdk.NewStraightFlush([]*sdk.Card{card10c, cardJc, cardKc, cardAc, wildcard})
	fmt.Printf("  Valid: %v\n", sf3.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 10)\n", sf3.ComparisonKey)
	
	if !sf3.IsValid() || sf3.ComparisonKey != 10 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试6: 中间同花顺 5-6-7-8-9 (方片)
	fmt.Println("测试6: 方片 5-6-7-8-9 同花顺")
	card5d, _ := sdk.NewCard(5, "Diamond", level)
	card6d, _ := sdk.NewCard(6, "Diamond", level)
	card7d, _ := sdk.NewCard(7, "Diamond", level)
	card8d, _ := sdk.NewCard(8, "Diamond", level)
	card9d, _ := sdk.NewCard(9, "Diamond", level)
	
	sf4 := sdk.NewStraightFlush([]*sdk.Card{card5d, card6d, card7d, card8d, card9d})
	fmt.Printf("  Valid: %v\n", sf4.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 5)\n", sf4.ComparisonKey)
	
	if !sf4.IsValid() || sf4.ComparisonKey != 5 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	fmt.Println("=== 所有测试通过! ✓ ===")
}
