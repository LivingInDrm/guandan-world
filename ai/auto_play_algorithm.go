package ai

import "guandan-world/sdk"

// AutoPlayAlgorithm 自动出牌算法接口
type AutoPlayAlgorithm interface {
	// SelectCardsToPlay 选择要出的牌
	// 参数:
	//   hand: 手牌
	//   trickInfo: 当前轮次信息 (包含是否为首出和领先牌组合)
	// 返回值:
	//   []*Card: 要出的牌，如果返回nil表示过牌
	SelectCardsToPlay(hand []*sdk.Card, trickInfo *sdk.TrickInfo) []*sdk.Card

	// SelectTributeCard 选择要进贡的牌
	// 参数:
	//   hand: 手牌
	//   excludeHeartTrump: 是否排除红桃主牌
	// 返回值:
	//   *Card: 选择的贡牌
	SelectTributeCard(hand []*sdk.Card, excludeHeartTrump bool) *sdk.Card

	// SelectReturnTributeCard 选择要还贡的牌
	// 参数:
	//   hand: 手牌
	//   receivedCard: 收到的贡牌
	// 返回值:
	//   *Card: 选择的还贡牌
	SelectReturnTributeCard(hand []*sdk.Card, receivedCard *sdk.Card) *sdk.Card
}

// SimpleAutoPlayAlgorithm 简单的自动出牌算法实现
type SimpleAutoPlayAlgorithm struct {
	level int // 当前级别
}

// NewSimpleAutoPlayAlgorithm 创建简单算法实例
func NewSimpleAutoPlayAlgorithm(level int) AutoPlayAlgorithm {
	return &SimpleAutoPlayAlgorithm{
		level: level,
	}
}

// SelectCardsToPlay 实现自动出牌逻辑
func (algo *SimpleAutoPlayAlgorithm) SelectCardsToPlay(hand []*sdk.Card, trickInfo *sdk.TrickInfo) []*sdk.Card {
	if hand == nil || len(hand) == 0 {
		return nil
	}

	// 如果是首出，选择最佳牌型
	if trickInfo.IsLeader {
		return algo.selectBestFirstPlay(hand)
	}

	// 如果不是首出，尝试跟牌
	return algo.tryToFollow(hand, trickInfo)
}

// selectBestFirstPlay 选择最佳首出牌型
func (algo *SimpleAutoPlayAlgorithm) selectBestFirstPlay(hand []*sdk.Card) []*sdk.Card {
	// 按张数优先的策略：钢板 > 顺子 > 葫芦 > 三张 > 对子 > 单张

	// 1. 检查钢板（三个对子）
	if steelPlate := algo.findSteelPlate(hand); steelPlate != nil {
		return steelPlate
	}

	// 2. 检查顺子（5张）
	if straight := algo.findStraight(hand, 5); straight != nil {
		return straight
	}

	// 3. 检查葫芦（三张+对子）
	if fullHouse := algo.findFullHouse(hand); fullHouse != nil {
		return fullHouse
	}

	// 4. 检查三张
	if triplet := algo.findTriplet(hand); triplet != nil {
		return triplet
	}

	// 5. 检查对子
	if pair := algo.findPair(hand); pair != nil {
		return pair
	}

	// 6. 单张（选择最小的）
	return []*sdk.Card{algo.findSmallestCard(hand)}
}

// tryToFollow 尝试跟牌
func (algo *SimpleAutoPlayAlgorithm) tryToFollow(hand []*sdk.Card, trickInfo *sdk.TrickInfo) []*sdk.Card {
	// 简化版：尝试找到能beat当前领先牌组合的牌
	leadComp := trickInfo.LeadComp
	if leadComp == nil {
		return nil // 过牌
	}

	// 简单的单张跟牌逻辑
	leadCards := leadComp.GetCards()
	if len(leadCards) == 1 {
		// 找比领先牌更大的单张
		for _, card := range hand {
			if card.GreaterThan(leadCards[0]) {
				return []*sdk.Card{card}
			}
		}
	}

	return nil // 过牌
}

// SelectTributeCard 选择贡牌
func (algo *SimpleAutoPlayAlgorithm) SelectTributeCard(hand []*sdk.Card, excludeHeartTrump bool) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}

	// 选择最大的非关键牌
	var bestCard *sdk.Card
	for _, card := range hand {
		// 如果需要排除红桃主牌，跳过红桃主牌
		if excludeHeartTrump && card.Color == "H" && card.Number == algo.level {
			continue
		}

		// 选择最大的牌
		if bestCard == nil || card.GreaterThan(bestCard) {
			bestCard = card
		}
	}

	return bestCard
}

// SelectReturnTributeCard 选择还贡牌
func (algo *SimpleAutoPlayAlgorithm) SelectReturnTributeCard(hand []*sdk.Card, receivedCard *sdk.Card) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}

	// 选择最小的牌进行还贡
	smallestCard := hand[0]
	for _, card := range hand[1:] {
		if !card.GreaterThan(smallestCard) {
			smallestCard = card
		}
	}

	return smallestCard
}

// 辅助方法 - 简化实现

func (algo *SimpleAutoPlayAlgorithm) findSteelPlate(hand []*sdk.Card) []*sdk.Card {
	// 简化：这里应该实现真正的钢板检测逻辑
	return nil
}

func (algo *SimpleAutoPlayAlgorithm) findStraight(hand []*sdk.Card, length int) []*sdk.Card {
	// 简化：这里应该实现真正的顺子检测逻辑
	return nil
}

func (algo *SimpleAutoPlayAlgorithm) findFullHouse(hand []*sdk.Card) []*sdk.Card {
	// 简化：这里应该实现真正的葫芦检测逻辑
	return nil
}

func (algo *SimpleAutoPlayAlgorithm) findTriplet(hand []*sdk.Card) []*sdk.Card {
	// 简化：这里应该实现真正的三张检测逻辑
	return nil
}

func (algo *SimpleAutoPlayAlgorithm) findPair(hand []*sdk.Card) []*sdk.Card {
	// 简化：这里应该实现真正的对子检测逻辑
	return nil
}

func (algo *SimpleAutoPlayAlgorithm) findSmallestCard(hand []*sdk.Card) *sdk.Card {
	if len(hand) == 0 {
		return nil
	}

	smallest := hand[0]
	for _, card := range hand[1:] {
		if !card.GreaterThan(smallest) {
			smallest = card
		}
	}

	return smallest
}
