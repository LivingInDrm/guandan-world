package sdk

import (
	"fmt"
)

// CompExample 演示牌组功能的示例函数
func CompExample() {
	fmt.Println("=== 掼蛋牌组功能演示 ===")

	// 1. 演示各种牌型的创建
	fmt.Println("\n1. 创建各种牌型：")

	// 单张
	card1, _ := NewCard(5, "Spade", 2)
	single := NewSingle([]*Card{card1})
	fmt.Printf("单张: %s, 有效: %v\n", single.String(), single.IsValid())

	// 对子
	card2, _ := NewCard(5, "Heart", 2)
	pair := NewPair([]*Card{card1, card2})
	fmt.Printf("对子: %s, 有效: %v\n", pair.String(), pair.IsValid())

	// 三张
	card3, _ := NewCard(5, "Diamond", 2)
	triple := NewTriple([]*Card{card1, card2, card3})
	fmt.Printf("三张: %s, 有效: %v\n", triple.String(), triple.IsValid())

	// 葫芦（三带二）
	card4, _ := NewCard(7, "Spade", 2)
	card5, _ := NewCard(7, "Heart", 2)
	fullhouse := NewFullHouse([]*Card{card1, card2, card3, card4, card5})
	fmt.Printf("葫芦: %s, 有效: %v\n", fullhouse.String(), fullhouse.IsValid())

	// 顺子
	sc1, _ := NewCard(3, "Spade", 2)
	sc2, _ := NewCard(4, "Heart", 2)
	sc3, _ := NewCard(5, "Diamond", 2)
	sc4, _ := NewCard(6, "Spade", 2)
	sc5, _ := NewCard(7, "Heart", 2)
	straight := NewStraight([]*Card{sc1, sc2, sc3, sc4, sc5})
	fmt.Printf("顺子: %s, 有效: %v\n", straight.String(), straight.IsValid())

	// 同花顺
	sfc1, _ := NewCard(3, "Spade", 2)
	sfc2, _ := NewCard(4, "Spade", 2)
	sfc3, _ := NewCard(5, "Spade", 2)
	sfc4, _ := NewCard(6, "Spade", 2)
	sfc5, _ := NewCard(7, "Spade", 2)
	straightFlush := NewStraightFlush([]*Card{sfc1, sfc2, sfc3, sfc4, sfc5})
	fmt.Printf("同花顺: %s, 有效: %v, 炸弹: %v\n", straightFlush.String(), straightFlush.IsValid(), straightFlush.IsBomb())

	// 钢板（连续三张）
	pc1, _ := NewCard(5, "Spade", 2)
	pc2, _ := NewCard(5, "Heart", 2)
	pc3, _ := NewCard(5, "Diamond", 2)
	pc4, _ := NewCard(6, "Spade", 2)
	pc5, _ := NewCard(6, "Heart", 2)
	pc6, _ := NewCard(6, "Diamond", 2)
	plate := NewPlate([]*Card{pc1, pc2, pc3, pc4, pc5, pc6})
	fmt.Printf("钢板: %s, 有效: %v\n", plate.String(), plate.IsValid())

	// 钢管（连续对子）
	tc1, _ := NewCard(5, "Spade", 2)
	tc2, _ := NewCard(5, "Heart", 2)
	tc3, _ := NewCard(6, "Spade", 2)
	tc4, _ := NewCard(6, "Heart", 2)
	tc5, _ := NewCard(7, "Spade", 2)
	tc6, _ := NewCard(7, "Heart", 2)
	tube := NewTube([]*Card{tc1, tc2, tc3, tc4, tc5, tc6})
	fmt.Printf("钢管: %s, 有效: %v\n", tube.String(), tube.IsValid())

	// 4张炸弹
	bc1, _ := NewCard(5, "Spade", 2)
	bc2, _ := NewCard(5, "Heart", 2)
	bc3, _ := NewCard(5, "Diamond", 2)
	bc4, _ := NewCard(5, "Club", 2)
	bomb := NewNaiveBomb([]*Card{bc1, bc2, bc3, bc4})
	fmt.Printf("4张炸弹: %s, 有效: %v, 炸弹: %v\n", bomb.String(), bomb.IsValid(), bomb.IsBomb())

	// 王炸
	joker1, _ := NewCard(15, "Joker", 2)
	joker2, _ := NewCard(15, "Joker", 2)
	joker3, _ := NewCard(16, "Joker", 2)
	joker4, _ := NewCard(16, "Joker", 2)
	jokerBomb := NewJokerBomb([]*Card{joker1, joker2, joker3, joker4})
	fmt.Printf("王炸: %s, 有效: %v, 炸弹: %v\n", jokerBomb.String(), jokerBomb.IsValid(), jokerBomb.IsBomb())

	// 2. 演示牌组比较
	fmt.Println("\n2. 牌组大小比较：")

	// 单张比较
	card6, _ := NewCard(7, "Heart", 2)
	single2 := NewSingle([]*Card{card6})
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		single.String(), single2.String(), single.String(), single.GreaterThan(single2))
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		single2.String(), single.String(), single2.String(), single2.GreaterThan(single))

	// 炸弹 vs 普通牌
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		bomb.String(), single.String(), bomb.String(), bomb.GreaterThan(single))
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		bomb.String(), straight.String(), bomb.String(), bomb.GreaterThan(straight))

	// 同花顺 vs 炸弹
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		straightFlush.String(), bomb.String(), straightFlush.String(), straightFlush.GreaterThan(bomb))

	// 王炸 vs 一切
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		jokerBomb.String(), bomb.String(), jokerBomb.String(), jokerBomb.GreaterThan(bomb))
	fmt.Printf("%s vs %s: %s 更大 = %v\n",
		jokerBomb.String(), straightFlush.String(), jokerBomb.String(), jokerBomb.GreaterThan(straightFlush))

	// 3. 演示 FromCardList 自动识别
	fmt.Println("\n3. 自动识别牌型：")

	// 测试各种牌组合
	testCards := [][]*Card{
		{card1},                             // 单张
		{card1, card2},                      // 对子
		{card1, card2, card3},               // 三张
		{joker1, joker2, joker3, joker4},    // 王炸
		{sc1, sc2, sc3, sc4, sc5},           // 顺子
		{sfc1, sfc2, sfc3, sfc4, sfc5},      // 同花顺
		{bc1, bc2, bc3, bc4},                // 炸弹
		{tc1, tc2, tc3, tc4, tc5, tc6},      // 钢管
		{pc1, pc2, pc3, pc4, pc5, pc6},      // 钢板
		{card1, card2, card3, card4, card5}, // 葫芦
	}

	for i, cards := range testCards {
		comp := FromCardList(cards, nil)
		fmt.Printf("牌组 %d: %s, 类型: %v, 有效: %v\n",
			i+1, comp.String(), comp.GetType(), comp.IsValid())
	}

	// 4. 演示变化牌
	fmt.Println("\n4. 变化牌演示：")

	// 级别为2时，红桃2是变化牌
	wildcard, _ := NewCard(2, "Heart", 2)
	normalCard, _ := NewCard(5, "Spade", 2)

	fmt.Printf("变化牌: %s, 是否为变化牌: %v\n", wildcard.String(), wildcard.IsWildcard())
	fmt.Printf("普通牌: %s, 是否为变化牌: %v\n", normalCard.String(), normalCard.IsWildcard())

	// 变化牌和普通牌组成对子
	wildcardPair := NewPair([]*Card{wildcard, normalCard})
	fmt.Printf("变化牌对子: %s, 有效: %v\n", wildcardPair.String(), wildcardPair.IsValid())

	// 5. 演示非法牌组
	fmt.Println("\n5. 非法牌组演示：")

	// 不符合规则的牌组
	invalidCards := []*Card{card1, card2, card6} // 三张不同数字的牌
	invalidComp := FromCardList(invalidCards, nil)
	fmt.Printf("非法牌组: %s, 类型: %v, 有效: %v\n",
		invalidComp.String(), invalidComp.GetType(), invalidComp.IsValid())

	// 三张包含王
	jokerTriple := NewTriple([]*Card{card1, card2, joker1})
	fmt.Printf("包含王的三张: %s, 有效: %v\n", jokerTriple.String(), jokerTriple.IsValid())

	// 6. 演示弃牌
	fmt.Println("\n6. 弃牌演示：")

	fold := FromCardList([]*Card{}, nil)
	fmt.Printf("弃牌: %s, 类型: %v, 有效: %v\n",
		fold.String(), fold.GetType(), fold.IsValid())

	fmt.Println("\n=== 演示完成 ===")
}

// 如果需要运行示例，可以在 main 函数中调用
func CompExampleMain() {
	CompExample()
}

// DemoCompTypes 演示各种牌型的优先级
func DemoCompTypes() {
	fmt.Println("\n=== 掼蛋牌型优先级演示 ===")

	// 创建各种牌型
	comps := make([]CardComp, 0)

	// 单张
	card1, _ := NewCard(5, "Spade", 2)
	single := NewSingle([]*Card{card1})
	comps = append(comps, single)

	// 对子
	card2, _ := NewCard(5, "Heart", 2)
	pair := NewPair([]*Card{card1, card2})
	comps = append(comps, pair)

	// 三张
	card3, _ := NewCard(5, "Diamond", 2)
	triple := NewTriple([]*Card{card1, card2, card3})
	comps = append(comps, triple)

	// 顺子
	sc1, _ := NewCard(3, "Spade", 2)
	sc2, _ := NewCard(4, "Heart", 2)
	sc3, _ := NewCard(5, "Diamond", 2)
	sc4, _ := NewCard(6, "Spade", 2)
	sc5, _ := NewCard(7, "Heart", 2)
	straight := NewStraight([]*Card{sc1, sc2, sc3, sc4, sc5})
	comps = append(comps, straight)

	// 4张炸弹
	bc1, _ := NewCard(5, "Spade", 2)
	bc2, _ := NewCard(5, "Heart", 2)
	bc3, _ := NewCard(5, "Diamond", 2)
	bc4, _ := NewCard(5, "Club", 2)
	bomb := NewNaiveBomb([]*Card{bc1, bc2, bc3, bc4})
	comps = append(comps, bomb)

	// 同花顺
	sfc1, _ := NewCard(3, "Spade", 2)
	sfc2, _ := NewCard(4, "Spade", 2)
	sfc3, _ := NewCard(5, "Spade", 2)
	sfc4, _ := NewCard(6, "Spade", 2)
	sfc5, _ := NewCard(7, "Spade", 2)
	straightFlush := NewStraightFlush([]*Card{sfc1, sfc2, sfc3, sfc4, sfc5})
	comps = append(comps, straightFlush)

	// 王炸
	joker1, _ := NewCard(15, "Joker", 2)
	joker2, _ := NewCard(15, "Joker", 2)
	joker3, _ := NewCard(16, "Joker", 2)
	joker4, _ := NewCard(16, "Joker", 2)
	jokerBomb := NewJokerBomb([]*Card{joker1, joker2, joker3, joker4})
	comps = append(comps, jokerBomb)

	// 比较所有牌型
	fmt.Println("\n各牌型能否压过其他牌型：")
	compNames := []string{"单张", "对子", "三张", "顺子", "4张炸弹", "同花顺", "王炸"}

	for i, comp1 := range comps {
		for j, comp2 := range comps {
			if i != j {
				canWin := comp1.GreaterThan(comp2)
				fmt.Printf("%s vs %s: %v\n", compNames[i], compNames[j], canWin)
			}
		}
	}

	fmt.Println("\n=== 优先级演示完成 ===")
}
