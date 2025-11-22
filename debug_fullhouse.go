package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	fmt.Println("=== FullHouse Debug Test ===\n")

	// 测试用例3: 情况3 - 1 wildcard, 2 pairs
	fmt.Println("Test 3: 77 + 55 + 1wildcard (情况3)")
	cards3 := []*sdk.Card{
		mustNewCard(7, "Spade", 8),
		mustNewCard(7, "Heart", 8),
		mustNewCard(8, "Spade", 8), // 8是变化牌
		mustNewCard(5, "Diamond", 8),
		mustNewCard(5, "Club", 8),
	}
	
	fmt.Println("Cards:")
	for i, c := range cards3 {
		fmt.Printf("  [%d] %v, RawNumber=%d, Number=%d, Level=%d, IsWildcard=%v\n", 
			i, c, c.RawNumber, c.Number, c.Level, c.IsWildcard())
	}
	
	valid, norm := sdk.FullHouseSatisfyNew(cards3)
	fmt.Printf("\nResult: Valid=%v\n", valid)
	if valid {
		fmt.Printf("Normalized: %v\n", norm)
	}
	
	fmt.Println("\n" + "=".repeat(50) + "\n")
	
	// 测试用例4: 情况4 - 2 wildcard, pair + single
	fmt.Println("Test 4: 22 + 3 + 2wildcard (情况4)")
	cards4 := []*sdk.Card{
		mustNewCard(2, "Spade", 5),
		mustNewCard(2, "Heart", 5),
		mustNewCard(3, "Spade", 5),
		mustNewCard(5, "Heart", 5),   // 5是变化牌
		mustNewCard(5, "Diamond", 5), // 5是变化牌
	}
	
	fmt.Println("Cards:")
	for i, c := range cards4 {
		fmt.Printf("  [%d] %v, RawNumber=%d, Number=%d, Level=%d, IsWildcard=%v\n", 
			i, c, c.RawNumber, c.Number, c.Level, c.IsWildcard())
	}
	
	valid4, norm4 := sdk.FullHouseSatisfyNew(cards4)
	fmt.Printf("\nResult: Valid=%v\n", valid4)
	if valid4 {
		fmt.Printf("Normalized: %v\n", norm4)
	}
}

func mustNewCard(number int, color string, level int) *sdk.Card {
	card, _ := sdk.NewCard(number, color, level)
	return card
}
