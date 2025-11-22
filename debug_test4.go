package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	cards := []*sdk.Card{
		mustNewCard(2, "Spade", 5),
		mustNewCard(2, "Heart", 5),
		mustNewCard(3, "Spade", 5),
		mustNewCard(5, "Heart", 5),
		mustNewCard(5, "Diamond", 5),
	}
	
	fmt.Println("Cards:")
	for i, c := range cards {
		fmt.Printf("[%d] %v - IsWildcard=%v\n", i, c, c.IsWildcard())
	}
}

func mustNewCard(number int, color string, level int) *sdk.Card {
	card, _ := sdk.NewCard(number, color, level)
	return card
}
