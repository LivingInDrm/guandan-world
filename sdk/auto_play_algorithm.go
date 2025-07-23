package sdk

// AutoPlayAlgorithm 自动出牌算法接口
type AutoPlayAlgorithm interface {
	// SelectCardsToPlay 选择要出的牌
	// 参数:
	//   hand: 手牌
	//   trickInfo: 当前轮次信息 (包含是否为首出和领先牌组合)
	// 返回值:
	//   []*Card: 要出的牌，如果返回nil表示过牌
	SelectCardsToPlay(hand []*Card, trickInfo *TrickInfo) []*Card

	// SelectTributeCard 选择要进贡的牌
	// 参数:
	//   hand: 手牌
	//   excludeHeartTrump: 是否排除红桃主牌
	// 返回值:
	//   *Card: 选择的贡牌
	SelectTributeCard(hand []*Card, excludeHeartTrump bool) *Card

	// SelectReturnCard 选择要还贡的牌
	// 参数:
	//   hand: 手牌
	//   avoidBreakingBomb: 是否避免破坏炸弹
	// 返回值:
	//   *Card: 选择的还贡牌
	SelectReturnCard(hand []*Card, avoidBreakingBomb bool) *Card
}

// SimpleAutoPlayAlgorithm 简单的自动出牌算法实现
type SimpleAutoPlayAlgorithm struct {
	level int // 当前级别，用于判断主牌
}

// NewSimpleAutoPlayAlgorithm 创建简单自动出牌算法
func NewSimpleAutoPlayAlgorithm(level int) *SimpleAutoPlayAlgorithm {
	return &SimpleAutoPlayAlgorithm{
		level: level,
	}
}

// SelectCardsToPlay 实现自动出牌逻辑
func (s *SimpleAutoPlayAlgorithm) SelectCardsToPlay(hand []*Card, trickInfo *TrickInfo) []*Card {
	if trickInfo == nil || trickInfo.IsLeader {
		// 首出：出张数尽可能多的合法非炸弹牌
		return s.selectLeaderPlay(hand)
	} else {
		// 跟牌：如果能压过则出牌，否则过牌
		return s.selectFollowPlay(hand, trickInfo)
	}
}

// selectLeaderPlay 首出选牌逻辑
func (s *SimpleAutoPlayAlgorithm) selectLeaderPlay(hand []*Card) []*Card {
	if len(hand) == 0 {
		return nil
	}

	// 简单策略：找最小的单牌
	var minCard *Card
	for _, card := range hand {
		if minCard == nil || card.LessThan(minCard) {
			minCard = card
		}
	}

	if minCard != nil {
		return []*Card{minCard}
	}

	return nil
}

// selectFollowPlay 跟牌选牌逻辑
func (s *SimpleAutoPlayAlgorithm) selectFollowPlay(hand []*Card, trickInfo *TrickInfo) []*Card {
	// 简单策略：总是过牌
	return nil
}

// SelectTributeCard 选择要进贡的牌
func (s *SimpleAutoPlayAlgorithm) SelectTributeCard(hand []*Card, excludeHeartTrump bool) *Card {
	if len(hand) == 0 {
		return nil
	}

	// 选择最大的牌（通常进贡要给最大的）
	var maxCard *Card
	for _, card := range hand {
		// 如果需要排除红桃主牌
		if excludeHeartTrump && card.Color == "Heart" && card.Number == s.level {
			continue
		}

		if maxCard == nil || card.GreaterThan(maxCard) {
			maxCard = card
		}
	}

	// 如果没有找到合适的牌，则选最大的
	if maxCard == nil && len(hand) > 0 {
		maxCard = hand[0]
		for _, card := range hand {
			if card.GreaterThan(maxCard) {
				maxCard = card
			}
		}
	}

	return maxCard
}

// SelectReturnCard 选择要还贡的牌
func (s *SimpleAutoPlayAlgorithm) SelectReturnCard(hand []*Card, avoidBreakingBomb bool) *Card {
	if len(hand) == 0 {
		return nil
	}

	// 如果需要避免破坏炸弹，先统计牌数
	cardCount := make(map[int]int)
	if avoidBreakingBomb {
		for _, card := range hand {
			if card.Color != "Joker" {
				cardCount[card.Number]++
			}
		}
	}

	// 选择最小的牌（避免破坏炸弹）
	var minCard *Card
	for _, card := range hand {
		// 如果需要避免破坏炸弹，检查这张牌是否会破坏炸弹
		if avoidBreakingBomb && card.Color != "Joker" {
			if cardCount[card.Number] >= 4 {
				continue // 跳过会破坏炸弹的牌
			}
		}

		if minCard == nil || card.LessThan(minCard) {
			minCard = card
		}
	}

	// 如果没有找到合适的牌，则选最小的
	if minCard == nil && len(hand) > 0 {
		minCard = hand[0]
		for _, card := range hand {
			if card.LessThan(minCard) {
				minCard = card
			}
		}
	}

	return minCard
}
