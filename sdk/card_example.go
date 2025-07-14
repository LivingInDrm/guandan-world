package sdk

import (
	"fmt"
	"log"
)

// CardExample 演示 Card 功能的示例函数
func CardExample() {
	fmt.Println("=== 掼蛋 Card 功能演示 ===")

	// 创建一些牌
	fmt.Println("\n1. 创建各种牌：")

	// 普通牌
	card1, err := NewCard(3, "Spade", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("普通牌: %s\n", card1.String())

	// Ace 牌
	ace, err := NewCard(1, "Heart", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ace 牌: %s (Number=%d, RawNumber=%d)\n", ace.String(), ace.Number, ace.RawNumber)

	// 人头牌
	jack, err := NewCard(11, "Diamond", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("人头牌: %s\n", jack.String())

	// 大王
	redJoker, err := NewCard(16, "Joker", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("大王: %s\n", redJoker.String())

	// 小王
	blackJoker, err := NewCard(15, "Joker", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("小王: %s\n", blackJoker.String())

	// 级别牌
	levelCard, err := NewCard(2, "Heart", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("级别牌: %s\n", levelCard.String())

	// 2. 演示变化牌（红桃级别牌）
	fmt.Println("\n2. 变化牌判断：")
	fmt.Printf("%s 是否为变化牌: %v\n", levelCard.String(), levelCard.IsWildcard())

	normalCard, _ := NewCard(2, "Spade", 2)
	fmt.Printf("%s 是否为变化牌: %v\n", normalCard.String(), normalCard.IsWildcard())

	// 3. 演示牌的比较
	fmt.Println("\n3. 牌的大小比较：")
	fmt.Printf("%s vs %s: %s 更大 = %v\n", card1.String(), jack.String(), card1.String(), card1.GreaterThan(jack))
	fmt.Printf("%s vs %s: %s 更大 = %v\n", jack.String(), card1.String(), jack.String(), jack.GreaterThan(card1))
	fmt.Printf("%s vs %s: %s 更大 = %v\n", levelCard.String(), jack.String(), levelCard.String(), levelCard.GreaterThan(jack))
	fmt.Printf("%s vs %s: %s 更大 = %v\n", redJoker.String(), blackJoker.String(), redJoker.String(), redJoker.GreaterThan(blackJoker))

	// 4. 演示顺子比较
	fmt.Println("\n4. 顺子比较（用于判断连续性）：")
	two, _ := NewCard(2, "Spade", 3)
	three, _ := NewCard(3, "Heart", 3)
	fmt.Printf("%s vs %s: %s 顺子更大 = %v\n", three.String(), two.String(), three.String(), three.ConsecutiveGreaterThan(two))
	fmt.Printf("%s vs %s: %s 顺子更大 = %v\n", ace.String(), two.String(), ace.String(), ace.ConsecutiveGreaterThan(two))

	// 5. 演示相等判断
	fmt.Println("\n5. 相等判断：")
	sameCard, _ := NewCard(3, "Heart", 2)
	fmt.Printf("%s vs %s: 相等 = %v\n", card1.String(), sameCard.String(), card1.Equals(sameCard))

	// 6. 演示 JSON 编码
	fmt.Println("\n6. JSON 编码：")
	json := card1.JSONEncode()
	fmt.Printf("%s 的 JSON: %v\n", card1.String(), json)

	// 7. 演示克隆
	fmt.Println("\n7. 克隆牌：")
	clonedCard := card1.Clone()
	fmt.Printf("原牌: %s\n", card1.String())
	fmt.Printf("克隆牌: %s\n", clonedCard.String())
	fmt.Printf("是否为同一对象: %v\n", card1 == clonedCard)

	fmt.Println("\n=== 演示完成 ===")
}

// 如果需要运行示例，可以在 main 函数中调用
func ExampleMain() {
	CardExample()
}
