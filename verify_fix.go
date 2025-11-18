package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	fmt.Println("=== 验证重复牌修复 ===\n")

	// 创建两张相同的黑桃A（点数和花色相同，DeckIndex不同）
	aceSpade1, _ := sdk.NewCard(14, "Spade", 5)
	aceSpade1.DeckIndex = 12

	aceSpade2, _ := sdk.NewCard(14, "Spade", 5)
	aceSpade2.DeckIndex = 64

	otherCard, _ := sdk.NewCard(13, "Heart", 5)
	otherCard.DeckIndex = 20

	fmt.Printf("创建了两张黑桃A:\n")
	fmt.Printf("  - 黑桃A #1: DeckIndex=%d, ID=%s\n", aceSpade1.DeckIndex, aceSpade1.GetID())
	fmt.Printf("  - 黑桃A #2: DeckIndex=%d, ID=%s\n", aceSpade2.DeckIndex, aceSpade2.GetID())

	// 设置手牌：玩家3有两张黑桃A
	playerHands := [4][]*sdk.Card{
		{},                             // 玩家0: 接收者
		{},                             // 玩家1
		{},                             // 玩家2
		{aceSpade1, aceSpade2, otherCard}, // 玩家3: 贡者
	}

	fmt.Printf("\n玩家3初始手牌:\n")
	for i, card := range playerHands[3] {
		fmt.Printf("  %d. %s (DeckIndex=%d)\n", i+1, card.String(), card.DeckIndex)
	}

	// 创建上贡阶段，指定贡牌为第一张黑桃A (DeckIndex=12)
	tributePhase := &sdk.TributePhase{
		Status:          sdk.TributeStatusFinished,
		TributePairs:    make([]*sdk.TributePair, 0),
		PoolCards:       make([]*sdk.Card, 0),
		SelectingPlayer: -1,
	}

	pair := &sdk.TributePair{
		Giver:       3,
		Receiver:    0,
		TributeCard: aceSpade1, // 明确指定第一张 (DeckIndex=12)
		ReturnCard:  otherCard,
	}
	tributePhase.TributePairs = append(tributePhase.TributePairs, pair)

	fmt.Printf("\n上贡设置:\n")
	fmt.Printf("  - 贡者: 玩家3\n")
	fmt.Printf("  - 接收者: 玩家0\n")
	fmt.Printf("  - 贡牌: %s (DeckIndex=%d)\n", aceSpade1.String(), aceSpade1.DeckIndex)
	fmt.Printf("  - 还贡: %s (DeckIndex=%d)\n", otherCard.String(), otherCard.DeckIndex)

	// 应用上贡
	tm := sdk.NewTributeManager(5)
	err := tm.ApplyTributeToHands(tributePhase, &playerHands)
	if err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("\n应用上贡后:\n")

	// 检查玩家3的手牌
	fmt.Printf("\n玩家3手牌 (%d张):\n", len(playerHands[3]))
	has12 := false
	has64 := false
	for i, card := range playerHands[3] {
		fmt.Printf("  %d. %s (DeckIndex=%d)\n", i+1, card.String(), card.DeckIndex)
		if card.DeckIndex == 12 {
			has12 = true
		}
		if card.DeckIndex == 64 {
			has64 = true
		}
	}

	// 检查玩家0的手牌
	fmt.Printf("\n玩家0手牌 (%d张):\n", len(playerHands[0]))
	receivedCorrectCard := false
	for i, card := range playerHands[0] {
		fmt.Printf("  %d. %s (DeckIndex=%d)\n", i+1, card.String(), card.DeckIndex)
		if card.DeckIndex == 12 {
			receivedCorrectCard = true
		}
	}

	// 验证结果
	fmt.Printf("\n=== 验证结果 ===\n")
	success := true

	if has12 {
		fmt.Printf("❌ 失败: 玩家3仍持有DeckIndex=12的黑桃A（应该被移除）\n")
		success = false
	} else {
		fmt.Printf("✅ 通过: 玩家3正确失去DeckIndex=12的黑桃A\n")
	}

	if !has64 {
		fmt.Printf("❌ 失败: 玩家3失去了DeckIndex=64的黑桃A（应该保留）\n")
		success = false
	} else {
		fmt.Printf("✅ 通过: 玩家3保留了DeckIndex=64的黑桃A\n")
	}

	if !receivedCorrectCard {
		fmt.Printf("❌ 失败: 玩家0没有收到DeckIndex=12的黑桃A\n")
		success = false
	} else {
		fmt.Printf("✅ 通过: 玩家0正确收到DeckIndex=12的黑桃A\n")
	}

	if success {
		fmt.Printf("\n🎉 所有验证通过！修复成功！\n")
	} else {
		fmt.Printf("\n⚠️ 验证失败，仍存在问题\n")
	}
}
