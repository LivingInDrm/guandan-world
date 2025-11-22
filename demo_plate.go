package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("钢板验证演示（使用现有NewPlate）")
	fmt.Println("========================================\n")
	
	level := 5
	
	fmt.Println("测试基本钢板识别:")
	
	// A-2钢板（最小）
	testCase("A-2钢板（最小）", createPlateCards([]int{14, 14, 14, 2, 2, 2}, level), true)
	
	// 5-6钢板
	testCase("5-6钢板", createPlateCards([]int{5, 5, 5, 6, 6, 6}, level), true)
	
	// K-A钢板（最大）
	testCase("K-A钢板（最大）", createPlateCards([]int{13, 13, 13, 14, 14, 14}, level), true)
	
	// Q-K钢板
	testCase("Q-K钢板", createPlateCards([]int{12, 12, 12, 13, 13, 13}, level), true)
	
	fmt.Println("\n测试循环情况:")
	
	// A-2钢板（循环的低端）
	testCaseDetail("A-2钢板", createPlateCards([]int{14, 14, 14, 2, 2, 2}, level), true)
	
	// K-A钢板（循环的高端）
	testCaseDetail("K-A钢板", createPlateCards([]int{13, 13, 13, 14, 14, 14}, level), true)
	
	fmt.Println("\n非法情况验证:")
	
	// 不连续
	testCase("3-5钢板（不连续）", createPlateCards([]int{3, 3, 3, 5, 5, 5}, level), false)
	
	// 数字不足3张
	testCase("5-6钢板（牌数不够）", createPlateCards([]int{5, 5, 6, 6, 6, 7}, level), false)
	
	fmt.Println("\n========================================")
	fmt.Println("演示完成")
	fmt.Println("\n注意：新的plateSatisfyNew函数已实现在plate_comp.go中")
	fmt.Println("      目前NewPlate仍使用旧的plateSatisfy")
	fmt.Println("      后续可以修改comp.go将plateSatisfy替换为plateSatisfyNew")
	fmt.Println("========================================")
}

func formatPair(pair []int) string {
	names := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	name1 := names[pair[0]-1]
	name2 := names[pair[1]-1]
	return fmt.Sprintf("%s-%s", name1, name2)
}

func createPlateCards(numbers []int, level int) []*sdk.Card {
	colors := []string{"Spade", "Heart", "Club", "Diamond"}
	cards := make([]*sdk.Card, 0, len(numbers))
	
	for i, num := range numbers {
		color := colors[i%4]
		card, err := sdk.NewCard(num, color, level)
		if err != nil {
			panic(fmt.Sprintf("Failed to create card: %v", err))
		}
		cards = append(cards, card)
	}
	
	return cards
}

func testCase(name string, cards []*sdk.Card, expectedValid bool) {
	plate := sdk.NewPlate(cards)
	isValid := plate.IsValid()
	
	status := "✗"
	if isValid == expectedValid {
		status = "✓"
	}
	
	if isValid {
		fmt.Printf("   %s %s: 有效\n", status, name)
	} else {
		fmt.Printf("   %s %s: 无效\n", status, name)
	}
}

func testCaseDetail(name string, cards []*sdk.Card, expectedValid bool) {
	plate := sdk.NewPlate(cards)
	
	fmt.Printf("\n   %s:\n", name)
	fmt.Printf("   输入: ")
	for _, card := range cards {
		fmt.Printf("%s ", card.Name)
	}
	fmt.Printf("\n   RawNumbers: ")
	for _, card := range cards {
		fmt.Printf("%d ", card.RawNumber)
	}
	fmt.Printf("\n   验证结果: %v\n", plate.IsValid())
	
	if plate.IsValid() {
		fmt.Printf("   排序后: ")
		for _, card := range plate.Cards {
			fmt.Printf("%s ", card.Name)
		}
		fmt.Println()
	}
}
