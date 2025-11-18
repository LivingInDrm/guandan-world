package main

import (
	"fmt"
	"guandan-world/sdk"
)

func main() {
	fmt.Println("=== 验证现有贡牌功能未被破坏 ===\n")

	// 测试1: 正常的贡牌流程（没有重复牌）
	fmt.Println("测试1: 正常贡牌流程")
	testNormalTribute()

	// 测试2: 抗贡检测
	fmt.Println("\n测试2: 抗贡检测")
	testTributeImmunity()

	// 测试3: 选择最大牌（排除红桃Trump）
	fmt.Println("\n测试3: 选择最大牌")
	testGetHighestCard()

	fmt.Println("\n=== 所有测试通过 ✅ ===")
}

func testNormalTribute() {
	// 创建不重复的牌
	bigJoker, _ := sdk.NewCard(16, "Joker", 5)
	bigJoker.DeckIndex = 104

	aceSpade, _ := sdk.NewCard(14, "Spade", 5)
	aceSpade.DeckIndex = 12

	kingHeart, _ := sdk.NewCard(13, "Heart", 5)
	kingHeart.DeckIndex = 20

	playerHands := [4][]*sdk.Card{
		{},                       // 玩家0: 接收者
		{},                       // 玩家1
		{},                       // 玩家2
		{bigJoker, aceSpade, kingHeart}, // 玩家3: 贡者
	}

	// 创建上贡阶段
	tributePhase := &sdk.TributePhase{
		Status:          sdk.TributeStatusFinished,
		TributePairs:    make([]*sdk.TributePair, 0),
		PoolCards:       make([]*sdk.Card, 0),
		SelectingPlayer: -1,
	}

	pair := &sdk.TributePair{
		Giver:       3,
		Receiver:    0,
		TributeCard: bigJoker,
		ReturnCard:  kingHeart,
	}
	tributePhase.TributePairs = append(tributePhase.TributePairs, pair)

	// 应用上贡
	tm := sdk.NewTributeManager(5)
	err := tm.ApplyTributeToHands(tributePhase, &playerHands)
	if err != nil {
		panic(fmt.Sprintf("❌ 错误: %v", err))
	}

	// 验证
	hasBigJoker := false
	for _, card := range playerHands[0] {
		if card.DeckIndex == 104 {
			hasBigJoker = true
			break
		}
	}

	if !hasBigJoker {
		panic("❌ 玩家0应该收到大王")
	}

	hasKingHeart := false
	for _, card := range playerHands[3] {
		if card.DeckIndex == 20 {
			hasKingHeart = true
			break
		}
	}

	if !hasKingHeart {
		panic("❌ 玩家3应该收到还贡")
	}

	fmt.Println("  ✅ 贡牌和还贡正确转移")
}

func testTributeImmunity() {
	tm := sdk.NewTributeManager(5)

	// 创建手牌（败方Team 1有2张大王）
	// Team 0: 玩家0, 玩家2
	// Team 1: 玩家1, 玩家3
	bigJoker1, _ := sdk.NewCard(16, "Joker", 5)
	bigJoker2, _ := sdk.NewCard(16, "Joker", 5)

	handsWithBigJokers := [4][]*sdk.Card{
		{},          // 玩家0 (Team 0)
		{bigJoker1}, // 玩家1 (Team 1) - 败方
		{},          // 玩家2 (Team 0)
		{bigJoker2}, // 玩家3 (Team 1) - 败方
	}

	// 双下场景：Team 0 获胜，Team 1 败方
	lastResult := &sdk.DealResult{
		Rankings:    []int{0, 2, 1, 3}, // rank1=0(Team0), rank2=2(Team0), rank3=1(Team1), rank4=3(Team1)
		WinningTeam: 0,                 // Team 0 获胜
		VictoryType: sdk.VictoryTypeDoubleDown,
	}

	isImmune, details := tm.GetTributeImmunityDetails(lastResult, handsWithBigJokers)
	if !isImmune {
		fmt.Printf("  详情: %+v\n", details)
		panic("❌ 应该触发抗贡（败方有2张大王）")
	}

	fmt.Println("  ✅ 抗贡检测正确")
}

func testGetHighestCard() {
	// 测试通过上贡自动选择来验证getHighestCardExcludingHeartTrump逻辑
	// 创建一个实际的上贡场景
	lastResult := &sdk.DealResult{
		Rankings:    []int{0, 1, 2, 3},
		WinningTeam: 0,
		VictoryType: sdk.VictoryTypeSingleLast,
	}

	tributePhase, _ := sdk.NewTributePhase(lastResult)
	tm := sdk.NewTributeManager(5)

	// 玩家3需要上贡，手牌包含红桃Trump和其他牌
	heartTrump, _ := sdk.NewCard(5, "Heart", 5)
	heartTrump.DeckIndex = 1
	
	aceSpade, _ := sdk.NewCard(14, "Spade", 5)
	aceSpade.DeckIndex = 12
	
	kingHeart, _ := sdk.NewCard(13, "Heart", 5)
	kingHeart.DeckIndex = 20

	playerHands := [4][]*sdk.Card{
		{}, {}, {},
		{heartTrump, aceSpade, kingHeart}, // 玩家3
	}

	// 自动选择贡牌
	tm.ProcessTributePhaseAction(tributePhase, playerHands)

	// 检查选择的贡牌是否是黑桃A（而不是红桃5 Trump）
	if len(tributePhase.TributePairs) == 0 || tributePhase.TributePairs[0].TributeCard == nil {
		panic("❌ 未选择贡牌")
	}

	tributeCard := tributePhase.TributePairs[0].TributeCard
	if tributeCard.Number != 14 || tributeCard.Color != "Spade" {
		panic(fmt.Sprintf("❌ 应该选择黑桃A（排除红桃Trump），实际选择了 %s", tributeCard.String()))
	}

	fmt.Println("  ✅ 最大牌选择正确（排除红桃Trump）")
}
