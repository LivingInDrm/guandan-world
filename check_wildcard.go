package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	// IsWildcard要求 Number==Level && Color=="Heart"
	card1, _ := sdk.NewCard(8, "Spade", 8)
	card2, _ := sdk.NewCard(8, "Heart", 8)
	
	fmt.Printf("8♠ (level=8): IsWildcard=%v (Color=%s)\n", card1.IsWildcard(), card1.Color)
	fmt.Printf("8♥ (level=8): IsWildcard=%v (Color=%s)\n", card2.IsWildcard(), card2.Color)
}
