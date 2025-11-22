package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	level := 5
	fmt.Println("=== 测试 Straight 集成 ===\n")

	// 测试1: 10-J-Q-K-A (最大顺子)
	fmt.Println("测试1: 10-J-Q-K-A 顺子")
	card10, _ := sdk.NewCard(10, "Spade", level)
	cardJ, _ := sdk.NewCard(11, "Club", level)
	cardQ, _ := sdk.NewCard(12, "Heart", level)
	cardK, _ := sdk.NewCard(13, "Diamond", level)
	cardA, _ := sdk.NewCard(14, "Spade", level)
	
	straight1 := sdk.NewStraight([]*sdk.Card{card10, cardJ, cardQ, cardK, cardA})
	fmt.Printf("  Valid: %v\n", straight1.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 10)\n", straight1.ComparisonKey)
	
	if !straight1.IsValid() || straight1.ComparisonKey != 10 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试2: A-2-3-4-5 (最小顺子)
	fmt.Println("测试2: A-2-3-4-5 顺子")
	level2 := 7
	cardA2, _ := sdk.NewCard(14, "Spade", level2)
	card2, _ := sdk.NewCard(2, "Club", level2)
	card3, _ := sdk.NewCard(3, "Heart", level2)
	card4, _ := sdk.NewCard(4, "Diamond", level2)
	card5, _ := sdk.NewCard(5, "Spade", level2)
	
	straight2 := sdk.NewStraight([]*sdk.Card{cardA2, card2, card3, card4, card5})
	fmt.Printf("  Valid: %v\n", straight2.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 1)\n", straight2.ComparisonKey)
	
	if !straight2.IsValid() || straight2.ComparisonKey != 1 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试3: 比较逻辑
	fmt.Println("测试3: 顺子比较")
	fmt.Printf("  10-J-Q-K-A > A-2-3-4-5: %v (期望: true)\n", straight1.GreaterThan(straight2))
	fmt.Printf("  A-2-3-4-5 > 10-J-Q-K-A: %v (期望: false)\n", straight2.GreaterThan(straight1))
	
	if !straight1.GreaterThan(straight2) || straight2.GreaterThan(straight1) {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试4: 含万能牌的顺子
	fmt.Println("测试4: 含万能牌的顺子 (10-J-K-A+wild)")
	card10b, _ := sdk.NewCard(10, "Spade", level)
	cardJb, _ := sdk.NewCard(11, "Club", level)
	cardKb, _ := sdk.NewCard(13, "Diamond", level)
	cardAb, _ := sdk.NewCard(14, "Spade", level)
	wildcard, _ := sdk.NewCard(5, "Heart", level) // 5是变化牌
	
	straight3 := sdk.NewStraight([]*sdk.Card{card10b, cardJb, cardKb, cardAb, wildcard})
	fmt.Printf("  Valid: %v\n", straight3.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 10)\n", straight3.ComparisonKey)
	
	if !straight3.IsValid() || straight3.ComparisonKey != 10 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	// 测试5: 中间顺子 5-6-7-8-9
	fmt.Println("测试5: 5-6-7-8-9 顺子")
	card5b, _ := sdk.NewCard(5, "Spade", level)
	card6, _ := sdk.NewCard(6, "Club", level)
	card7, _ := sdk.NewCard(7, "Heart", level)
	card8, _ := sdk.NewCard(8, "Diamond", level)
	card9, _ := sdk.NewCard(9, "Spade", level)
	
	straight4 := sdk.NewStraight([]*sdk.Card{card5b, card6, card7, card8, card9})
	fmt.Printf("  Valid: %v\n", straight4.IsValid())
	fmt.Printf("  ComparisonKey: %d (期望: 5)\n", straight4.ComparisonKey)
	
	if !straight4.IsValid() || straight4.ComparisonKey != 5 {
		fmt.Println("  ❌ 失败")
		return
	}
	fmt.Println("  ✓ 通过\n")

	fmt.Println("=== 所有测试通过! ✓ ===")
}
