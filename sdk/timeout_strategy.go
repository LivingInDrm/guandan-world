package sdk

// TimeoutStrategy 定义超时处理策略接口
// SDK GameDriver 使用该接口在玩家超时时生成默认决策
type TimeoutStrategy interface {
	// GetDefaultPlayDecision 获取出牌超时的默认决策
	// 参数:
	//   hand: 玩家当前手牌
	//   trickInfo: 当前trick信息（是否为首出、领先牌型等）
	// 返回:
	//   默认的出牌决策（出牌或过牌）
	GetDefaultPlayDecision(hand []*Card, trickInfo *TrickInfo) *PlayDecision

	// GetDefaultTributeCard 获取贡牌超时的默认选择
	// 参数:
	//   options: 可选择的牌列表（双下情况）
	// 返回:
	//   选择的默认贡牌，如果options为空或全为nil则返回nil
	GetDefaultTributeCard(options []*Card) *Card

	// GetDefaultReturnCard 获取还贡超时的默认选择
	// 参数:
	//   hand: 玩家当前手牌
	// 返回:
	//   选择的默认还贡牌，如果hand为空或全为nil则返回nil
	GetDefaultReturnCard(hand []*Card) *Card
}

// DefaultTimeoutStrategy 默认超时策略实现
// 策略规则：
// - Leader超时：出最小单牌
// - 非Leader超时：PASS
// - 贡牌超时：选最大牌
// - 还贡超时：选最小牌
type DefaultTimeoutStrategy struct{}

// NewDefaultTimeoutStrategy 创建默认超时策略实例
func NewDefaultTimeoutStrategy() *DefaultTimeoutStrategy {
	return &DefaultTimeoutStrategy{}
}

// GetDefaultPlayDecision 实现出牌超时默认策略
func (s *DefaultTimeoutStrategy) GetDefaultPlayDecision(hand []*Card, trickInfo *TrickInfo) *PlayDecision {
	// 防御性检查
	if trickInfo == nil {
		return &PlayDecision{
			Action: ActionPass,
		}
	}

	// Leader必须出牌，选择最小的单牌
	if trickInfo.IsLeader {
		if len(hand) == 0 {
			// 没有手牌，返回PASS（理论上不应该发生）
			return &PlayDecision{
				Action: ActionPass,
			}
		}

		// 找到最小的牌
		smallestCard := s.findSmallestCard(hand)
		if smallestCard == nil {
			// 找不到最小牌（理论上不应该发生），返回PASS
			return &PlayDecision{
				Action: ActionPass,
			}
		}

		return &PlayDecision{
			Action: ActionPlay,
			Cards:  []*Card{smallestCard},
		}
	}

	// 非Leader默认PASS
	return &PlayDecision{
		Action: ActionPass,
	}
}

// GetDefaultTributeCard 实现贡牌超时默认策略
func (s *DefaultTimeoutStrategy) GetDefaultTributeCard(options []*Card) *Card {
	if len(options) == 0 {
		return nil
	}

	// 选择最大的牌，跳过nil元素
	var largestCard *Card
	for _, card := range options {
		if card == nil {
			continue
		}
		if largestCard == nil || card.GreaterThan(largestCard) {
			largestCard = card
		}
	}

	return largestCard
}

// GetDefaultReturnCard 实现还贡超时默认策略
func (s *DefaultTimeoutStrategy) GetDefaultReturnCard(hand []*Card) *Card {
	if len(hand) == 0 {
		return nil
	}

	// 选择最小的牌
	return s.findSmallestCard(hand)
}

// findSmallestCard 找到手牌中最小的牌
// 注意：会跳过nil元素，如果所有元素都是nil则返回nil
func (s *DefaultTimeoutStrategy) findSmallestCard(hand []*Card) *Card {
	if len(hand) == 0 {
		return nil
	}

	var smallestCard *Card
	for _, card := range hand {
		if card == nil {
			continue
		}
		if smallestCard == nil || card.LessThan(smallestCard) {
			smallestCard = card
		}
	}

	return smallestCard
}
