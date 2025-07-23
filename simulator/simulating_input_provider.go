package simulator

import (
	"context"
	"fmt"

	"guandan-world/ai"
	"guandan-world/sdk"
)

// SimulatingInputProvider 模拟输入提供者
// 将现有的AutoPlayAlgorithm适配到新的PlayerInputProvider接口
type SimulatingInputProvider struct {
	algorithms map[int]ai.AutoPlayAlgorithm // 每个玩家的自动算法
}

// NewSimulatingInputProvider 创建新的模拟输入提供者
func NewSimulatingInputProvider() *SimulatingInputProvider {
	return &SimulatingInputProvider{
		algorithms: make(map[int]ai.AutoPlayAlgorithm),
	}
}

// SetPlayerAlgorithm 为指定玩家设置自动算法
func (sip *SimulatingInputProvider) SetPlayerAlgorithm(playerSeat int, algorithm ai.AutoPlayAlgorithm) {
	sip.algorithms[playerSeat] = algorithm
}

// RequestPlayDecision 实现PlayerInputProvider接口 - 请求出牌决策
func (sip *SimulatingInputProvider) RequestPlayDecision(
	ctx context.Context,
	playerSeat int,
	hand []*sdk.Card,
	trickInfo *sdk.TrickInfo,
) (*sdk.PlayDecision, error) {
	// 检查是否有对应的算法
	algorithm, exists := sip.algorithms[playerSeat]
	if !exists {
		return &sdk.PlayDecision{Action: sdk.ActionPass}, nil
	}

	// 使用算法选择卡牌
	selectedCards := algorithm.SelectCardsToPlay(hand, trickInfo)

	// 构造决策
	if selectedCards == nil || len(selectedCards) == 0 {
		return &sdk.PlayDecision{Action: sdk.ActionPass}, nil
	}

	return &sdk.PlayDecision{
		Action: sdk.ActionPlay,
		Cards:  selectedCards,
	}, nil
}

// RequestTributeSelection 实现PlayerInputProvider接口 - 请求贡牌选择
func (sip *SimulatingInputProvider) RequestTributeSelection(
	ctx context.Context,
	playerSeat int,
	options []*sdk.Card,
) (*sdk.Card, error) {
	algorithm, exists := sip.algorithms[playerSeat]
	if !exists || len(options) == 0 {
		return options[0], nil // 默认选择第一张
	}

	// 使用算法从选项中选择
	// 这里简化处理，实际可能需要更复杂的逻辑
	selected := algorithm.SelectTributeCard(options, false)
	if selected == nil {
		return options[0], nil
	}

	return selected, nil
}

// RequestReturnTribute 实现PlayerInputProvider接口 - 请求还贡选择
func (sip *SimulatingInputProvider) RequestReturnTribute(
	ctx context.Context,
	playerSeat int,
	hand []*sdk.Card,
) (*sdk.Card, error) {
	algorithm, exists := sip.algorithms[playerSeat]
	if !exists || len(hand) == 0 {
		return hand[0], nil // 默认选择第一张
	}

	// 这里简化处理，传递nil作为receivedCard
	selected := algorithm.SelectReturnTributeCard(hand, nil)
	if selected == nil {
		return hand[0], nil
	}

	return selected, nil
}

// BatchSetAlgorithms 批量设置算法
func (sip *SimulatingInputProvider) BatchSetAlgorithms(algorithms []ai.AutoPlayAlgorithm) error {
	if len(algorithms) != 4 {
		return fmt.Errorf("exactly 4 algorithms required, got %d", len(algorithms))
	}

	for i, algorithm := range algorithms {
		sip.SetPlayerAlgorithm(i, algorithm)
	}

	return nil
}
