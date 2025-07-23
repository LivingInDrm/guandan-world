package sdk

import (
	"context"
	"fmt"
)

// SimulatingInputProvider 模拟输入提供者
// 将现有的AutoPlayAlgorithm适配到新的PlayerInputProvider接口
type SimulatingInputProvider struct {
	algorithms map[int]AutoPlayAlgorithm // 每个玩家的自动算法
}

// NewSimulatingInputProvider 创建新的模拟输入提供者
func NewSimulatingInputProvider() *SimulatingInputProvider {
	return &SimulatingInputProvider{
		algorithms: make(map[int]AutoPlayAlgorithm),
	}
}

// SetPlayerAlgorithm 为指定玩家设置自动算法
func (sip *SimulatingInputProvider) SetPlayerAlgorithm(playerSeat int, algorithm AutoPlayAlgorithm) {
	sip.algorithms[playerSeat] = algorithm
}

// RequestPlayDecision 实现PlayerInputProvider接口 - 请求出牌决策
func (sip *SimulatingInputProvider) RequestPlayDecision(
	ctx context.Context,
	playerSeat int,
	hand []*Card,
	trickInfo *TrickInfo,
) (*PlayDecision, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 获取玩家的算法
	algorithm, exists := sip.algorithms[playerSeat]
	if !exists {
		return nil, fmt.Errorf("no algorithm set for player %d", playerSeat)
	}

	// 使用算法选择要出的牌
	selectedCards := algorithm.SelectCardsToPlay(hand, trickInfo)

	// 构造决策结果
	if selectedCards != nil && len(selectedCards) > 0 {
		return &PlayDecision{
			Action: ActionPlay,
			Cards:  selectedCards,
		}, nil
	} else {
		// 如果算法返回nil或空数组，表示过牌
		// 但首出时不能过牌，强制出最小单张
		if trickInfo != nil && trickInfo.IsLeader {
			if len(hand) > 0 {
				// 找最小的牌强制出牌
				smallest := hand[0]
				for _, card := range hand {
					if card.LessThan(smallest) {
						smallest = card
					}
				}
				return &PlayDecision{
					Action: ActionPlay,
					Cards:  []*Card{smallest},
				}, nil
			}
			return nil, fmt.Errorf("player %d has no cards to play as leader", playerSeat)
		}

		return &PlayDecision{
			Action: ActionPass,
		}, nil
	}
}

// RequestTributeSelection 实现PlayerInputProvider接口 - 请求贡牌选择
func (sip *SimulatingInputProvider) RequestTributeSelection(
	ctx context.Context,
	playerSeat int,
	options []*Card,
) (*Card, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 检查是否为该玩家设置了算法
	if _, exists := sip.algorithms[playerSeat]; !exists {
		return nil, fmt.Errorf("no algorithm set for player %d", playerSeat)
	}

	// 双下选牌：选择最大的牌（贪心策略）
	if len(options) == 0 {
		return nil, fmt.Errorf("no tribute options for player %d", playerSeat)
	}

	// 双下选牌使用贪心策略：选择最大的牌
	selectedCard := options[0]
	for _, card := range options {
		if card.GreaterThan(selectedCard) {
			selectedCard = card
		}
	}

	return selectedCard, nil
}

// RequestReturnTribute 实现PlayerInputProvider接口 - 请求还贡选择
func (sip *SimulatingInputProvider) RequestReturnTribute(
	ctx context.Context,
	playerSeat int,
	hand []*Card,
) (*Card, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 获取玩家的算法
	algorithm, exists := sip.algorithms[playerSeat]
	if !exists {
		return nil, fmt.Errorf("no algorithm set for player %d", playerSeat)
	}

	// 使用算法的还贡逻辑
	returnCard := algorithm.SelectReturnCard(hand, true) // 避免破坏炸弹

	if returnCard == nil {
		return nil, fmt.Errorf("algorithm returned nil for return tribute for player %d", playerSeat)
	}

	return returnCard, nil
}

// BatchSetAlgorithms 批量设置所有玩家的算法（便捷方法）
func (sip *SimulatingInputProvider) BatchSetAlgorithms(algorithms []AutoPlayAlgorithm) error {
	if len(algorithms) != 4 {
		return fmt.Errorf("expected 4 algorithms, got %d", len(algorithms))
	}

	for i, algorithm := range algorithms {
		if algorithm == nil {
			return fmt.Errorf("algorithm for player %d is nil", i)
		}
		sip.algorithms[i] = algorithm
	}

	return nil
}

// GetPlayerAlgorithm 获取指定玩家的算法（用于调试）
func (sip *SimulatingInputProvider) GetPlayerAlgorithm(playerSeat int) AutoPlayAlgorithm {
	return sip.algorithms[playerSeat]
}

// HasAlgorithmForPlayer 检查是否为指定玩家设置了算法
func (sip *SimulatingInputProvider) HasAlgorithmForPlayer(playerSeat int) bool {
	_, exists := sip.algorithms[playerSeat]
	return exists
}

// Reset 重置所有算法设置
func (sip *SimulatingInputProvider) Reset() {
	sip.algorithms = make(map[int]AutoPlayAlgorithm)
}
