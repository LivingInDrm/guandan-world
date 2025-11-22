package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	fmt.Println("=== FullHouse Manual Test ===\n")

	// 测试用例1: 情况1 - 0 wildcard, triple + pair
	fmt.Println("Test 1: 333 + 22 (情况1)")
	cards1 := []*sdk.Card{
		mustNewCard(3, "Spade", 5),
		mustNewCard(3, "Heart", 5),
		mustNewCard(3, "Diamond", 5),
		mustNewCard(2, "Spade", 5),
		mustNewCard(2, "Heart", 5),
	}
	testFullHouse("333+22", cards1, 3)

	// 测试用例2: 情况2 - 1 wildcard, triple + single
	fmt.Println("Test 2: KKK + 2 + 1wildcard (情况2)")
	cards2 := []*sdk.Card{
		mustNewCard(13, "Spade", 5),
		mustNewCard(13, "Heart", 5),
		mustNewCard(13, "Diamond", 5),
		mustNewCard(2, "Spade", 5),
		mustNewCard(5, "Heart", 5), // 5是变化牌
	}
	testFullHouse("KKK+2+wildcard", cards2, 13)

	// 测试用例3: 情况3 - 1 wildcard, 2 pairs
	fmt.Println("Test 3: 77 + 55 + 1wildcard (情况3)")
	cards3 := []*sdk.Card{
		mustNewCard(7, "Spade", 8),
		mustNewCard(7, "Club", 8),
		mustNewCard(8, "Heart", 8), // 8♥是变化牌（level=8）
		mustNewCard(5, "Diamond", 8),
		mustNewCard(5, "Club", 8),
	}
	testFullHouse("77+55+wildcard", cards3, 7)

	// 测试用例4: 情况4 - 2 wildcard, pair + single
	fmt.Println("Test 4: 22 + 3 + 2wildcard (情况4)")
	cards4 := []*sdk.Card{
		mustNewCard(2, "Spade", 5),
		mustNewCard(2, "Heart", 5),
		mustNewCard(3, "Spade", 5),
		mustNewCard(5, "Heart", 5),   // 5是变化牌
		mustNewCard(5, "Diamond", 5), // 5是变化牌
	}
	testFullHouse("22+3+2wildcard", cards4, 3)

	// 测试用例5: 王对作为pair
	fmt.Println("Test 5: 222 + 大王大王 (王对作为pair)")
	bigJoker1, _ := sdk.NewCard(16, "Joker", 5)
	bigJoker2, _ := sdk.NewCard(16, "Joker", 5)
	cards5 := []*sdk.Card{
		mustNewCard(2, "Spade", 5),
		mustNewCard(2, "Heart", 5),
		mustNewCard(2, "Diamond", 5),
		bigJoker1,
		bigJoker2,
	}
	testFullHouse("222+大王大王", cards5, 2)

	// 测试用例6: 非法 - 1大王+1小王（不能配对）
	fmt.Println("Test 6: 333 + 大王小王 (非法-王不能混合配对)")
	bigJoker, _ := sdk.NewCard(16, "Joker", 5)
	smallJoker, _ := sdk.NewCard(15, "Joker", 5)
	cards6 := []*sdk.Card{
		mustNewCard(3, "Spade", 5),
		mustNewCard(3, "Heart", 5),
		mustNewCard(3, "Diamond", 5),
		bigJoker,
		smallJoker,
	}
	fmt.Printf("Expected: Invalid\n")
	valid6, _ := sdk.FullHouseSatisfyNew(cards6)
	fmt.Printf("Result: Valid=%v\n\n", valid6)

	// 测试用例7: 情况3 - 王对不能升级为triple
	fmt.Println("Test 7: 王王 + 33 + 1wildcard (王对只能作pair)")
	cards7 := []*sdk.Card{
		bigJoker1,
		bigJoker2,
		mustNewCard(3, "Spade", 5),
		mustNewCard(3, "Club", 5),
		mustNewCard(5, "Heart", 5), // 5♥是变化牌（level=5）
	}
	testFullHouse("王王+33+wildcard", cards7, 3)

	fmt.Println("=== Test Complete ===")
}

func testFullHouse(name string, cards []*sdk.Card, expectedKey int) {
	fmt.Printf("Expected ComparisonKey: %d\n", expectedKey)
	valid, norm := sdk.FullHouseSatisfyNew(cards)
	fmt.Printf("Result: Valid=%v\n", valid)
	if valid {
		fmt.Printf("Normalized order: ")
		for _, c := range norm {
			fmt.Printf("%v ", c)
		}
		fmt.Println()
	}
	fmt.Println()
}

func mustNewCard(number int, color string, level int) *sdk.Card {
	card, _ := sdk.NewCard(number, color, level)
	return card
}
