package sdk

import (
	"errors"
	"fmt"
)

// TributeManager handles all tribute-related operations independently
type TributeManager struct {
	level int // Current level for the game (used for wildcard detection)
}

// NewTributeManager creates a new tribute manager
func NewTributeManager(level int) *TributeManager {
	return &TributeManager{
		level: level,
	}
}

// NewTributePhase creates a new tribute phase based on the last deal result
func NewTributePhase(lastResult *DealResult) (*TributePhase, error) {
	if lastResult == nil {
		return nil, nil // No tribute needed for first deal
	}

	rankings := lastResult.Rankings
	if len(rankings) < 4 {
		return nil, errors.New("invalid rankings for tribute phase")
	}

	tributePhase := &TributePhase{
		Status:       TributeStatusWaiting,
		TributePairs: make([]*TributePair, 0),
		PoolCards:    make([]*Card, 0),
		Winners:      []int{},
	}

	// 按排名获取玩家
	rank1 := rankings[0] // 第1名
	rank2 := rankings[1] // 第2名
	rank3 := rankings[2] // 第3名
	rank4 := rankings[3] // 第4名

	// 统一使用贡牌池机制：所有上贡类型都先贡到池，然后从池中选择
	// 根据胜利类型确定上贡规则
	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		// Double Down: rank1, rank2 同队
		// Rank3 和 Rank4 各上交 1 张贡牌到贡牌池
		// Rank1 优先从贡牌池中挑选其一；Rank2 获得剩下的一张贡牌
		tributePhase.Winners = []int{rank1, rank2}
		tributePhase.TributePairs = append(tributePhase.TributePairs,
			&TributePair{Giver: rank3, Receiver: -1, TributeCard: nil, ReturnCard: nil},
			&TributePair{Giver: rank4, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		)

	case VictoryTypeSingleLast:
		// Single Last: rank1, rank3 同队
		// Rank4 上交 1 张贡牌到贡牌池，Rank1 自动获得
		tributePhase.Winners = []int{rank1}
		tributePhase.TributePairs = append(tributePhase.TributePairs,
			&TributePair{Giver: rank4, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		)

	case VictoryTypePartnerLast:
		// Partner Last: rank1, rank4 同队
		// Rank3 上交 1 张贡牌到贡牌池，Rank1 自动获得
		tributePhase.Winners = []int{rank1}
		tributePhase.TributePairs = append(tributePhase.TributePairs,
			&TributePair{Giver: rank3, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		)

	default:
		return nil, fmt.Errorf("unknown victory type: %v", lastResult.VictoryType)
	}

	return tributePhase, nil
}

// GetTributeImmunityDetails 获取详细的抗贡信息
// 返回是否免贡以及详细的原因说明
func (tm *TributeManager) GetTributeImmunityDetails(lastResult *DealResult, playerHands [4][]*Card) (bool, map[string]interface{}) {
	if lastResult == nil {
		return false, nil
	}

	// 获取输掉的队伍编号
	losingTeam := 1 - lastResult.WinningTeam

	// 统计每个败方玩家的大王详情
	var bigJokerHolders []map[string]interface{}
	totalBigJokers := 0

	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		// 检查该玩家是否属于输掉的队伍
		if playerSeat%2 == losingTeam {
			playerBigJokers := tm.countBigJokers(playerHands[playerSeat])
			if playerBigJokers > 0 {
				bigJokerHolders = append(bigJokerHolders, map[string]interface{}{
					"player_seat":     playerSeat,
					"big_joker_count": playerBigJokers,
				})
			}
			totalBigJokers += playerBigJokers
		}
	}

	// 判断是否触发抗贡
	isImmune := totalBigJokers >= 2

	// 构建详细信息
	details := map[string]interface{}{
		"big_joker_count":   totalBigJokers,
		"big_joker_holders": bigJokerHolders,
		"losing_team":       losingTeam,
		"description": fmt.Sprintf("败方队伍(Team %d)持有%d张大王%s",
			losingTeam, totalBigJokers,
			func() string {
				if isImmune {
					return "，触发抗贡"
				}
				return "，未达到抗贡条件(需要2张)"
			}()),
	}

	return isImmune, details
}

// countBigJokers 统计手牌中大王的数量
func (tm *TributeManager) countBigJokers(hand []*Card) int {
	count := 0
	for _, card := range hand {
		if card.IsBigJoker() {
			count++
		}
	}
	return count
}

// GetNextReceiver 找下一个待分配的 receiver（Winners 中还没成为任何 TributePair.Receiver 的第一个）
func (tm *TributeManager) GetNextReceiver(phase *TributePhase) int {
	assignedReceivers := make(map[int]bool)
	// 收集所有已分配的接收者
	for _, pair := range phase.TributePairs {
		if pair.Receiver != -1 {
			assignedReceivers[pair.Receiver] = true
		}
	}

	// 选择第一个尚未被分配的获胜者
	for _, winner := range phase.Winners {
		if !assignedReceivers[winner] {
			return winner
		}
	}
	return -1
}

// GetPendingReturnReceiver 找下一个需要还贡的 receiver
func (tm *TributeManager) GetPendingReturnReceiver(phase *TributePhase) (receiver int, giver int) {
	for _, pair := range phase.TributePairs {
		if pair.Receiver != -1 && pair.ReturnCard == nil {
			return pair.Receiver, pair.Giver
		}
	}
	return -1, -1
}

// getHighestCardExcludingHeartTrump 获取除红桃Trump外最大的一张牌
func (tm *TributeManager) getHighestCardExcludingHeartTrump(hand []*Card) *Card {
	if len(hand) == 0 {
		return nil
	}

	var highest *Card
	for _, card := range hand {
		// 排除红桃Trump牌（红桃且等于当前级别）
		if card.IsWildcard() { // IsWildcard() 判断是否为红桃Trump
			continue
		}

		if highest == nil || card.GreaterThan(highest) {
			highest = card
		}
	}

	// 如果没有找到合适的牌（全部都是红桃Trump），返回其中一张
	if highest == nil && len(hand) > 0 {
		highest = hand[0]
	}

	return highest
}

// buildTributeMapFromPairs builds a tribute map from TributePairs
// Note: Filters out pending assignments (receiver == -1) to avoid client confusion
func buildTributeMapFromPairs(pairs []*TributePair) map[int]int {
	result := make(map[int]int)
	for _, pair := range pairs {
		// 只包含已确定接收者的映射，过滤掉待定状态（receiver == -1）
		if pair.Receiver != -1 {
			result[pair.Giver] = pair.Receiver
		}
	}
	return result
}

// buildTributeCardsFromPairs builds a tribute cards map from TributePairs
func buildTributeCardsFromPairs(pairs []*TributePair) map[int]*Card {
	result := make(map[int]*Card)
	for _, pair := range pairs {
		if pair.TributeCard != nil {
			result[pair.Giver] = pair.TributeCard
		}
	}
	return result
}

// buildReturnCardsFromPairs builds a return cards map from TributePairs
func buildReturnCardsFromPairs(pairs []*TributePair) map[int]*Card {
	result := make(map[int]*Card)
	for _, pair := range pairs {
		if pair.Receiver != -1 && pair.ReturnCard != nil {
			result[pair.Receiver] = pair.ReturnCard
		}
	}
	return result
}

// ProcessTributeStep 纯函数：根据当前状态计算下一步操作
//
// 特点：
//   - 不修改任何输入参数
//   - 不发送事件
//   - 不依赖 GameEngine/Match/Deal
//
// 返回值包含所有需要执行的变更，由调用者负责应用
func ProcessTributeStep(
	phase *TributePhase,
	playerHands [4][]*Card,
	level int,
	input *TributeInput,
) *TributeStepResult {
	result := &TributeStepResult{}
	tm := NewTributeManager(level)

	switch phase.Status {

	case TributeStatusWaiting:
		result = processTributeWaiting(tm, phase, playerHands)

	case TributeStatusSelecting:
		result = processTributeSelecting(tm, phase, playerHands, input)

	case TributeStatusReturning:
		result = processTributeReturning(tm, phase, playerHands, input)

	case TributeStatusFinished:
		result = processTributeFinished(phase)

	default:
		// Unknown status, do nothing
	}

	return result
}

// processTributeWaiting 处理 Waiting 状态：选择贡牌并放入池
func processTributeWaiting(tm *TributeManager, phase *TributePhase, playerHands [4][]*Card) *TributeStepResult {
	result := &TributeStepResult{
		NextStatus: TributeStatusSelecting,
	}

	poolCards := make([]*Card, 0)

	for _, pair := range phase.TributePairs {
		tributeCard := tm.getHighestCardExcludingHeartTrump(playerHands[pair.Giver])
		if tributeCard == nil {
			result.Error = fmt.Errorf("player %d has no valid tribute card", pair.Giver)
			return result
		}

		poolCards = append(poolCards, tributeCard)

		// 事件意图
		result.Events = append(result.Events, TributeEventIntent{
			Type:       EventTributeCardSubmitted,
			PlayerSeat: pair.Giver,
			Card:       tributeCard,
		})

		// 手牌变更：从 giver 移除
		result.HandChanges = append(result.HandChanges, HandChange{
			PlayerSeat: pair.Giver,
			Card:       tributeCard,
			IsAdd:      false,
		})

		// Pair 变更：设置 TributeCard
		result.PairUpdates = append(result.PairUpdates, PairUpdate{
			GiverSeat:   pair.Giver,
			TributeCard: tributeCard,
		})
	}

	// 设置 PoolCards
	result.PoolCardsToSet = poolCards

	return result
}

// processTributeSelecting 处理 Selecting 状态：分配贡牌给接收者
func processTributeSelecting(tm *TributeManager, phase *TributePhase, playerHands [4][]*Card, input *TributeInput) *TributeStepResult {
	result := &TributeStepResult{}

	// 双下场景：PoolCards == 2，需要第一名选择
	if len(phase.PoolCards) == 2 && len(phase.Winners) > 0 {
		if input == nil {
			// 需要用户输入
			result.PendingAction = &TributeAction{
				Type:     TributeActionSelect,
				PlayerID: phase.Winners[0],
				Options:  phase.PoolCards,
			}
			return result
		}

		// 处理用户选择
		if input.Card == nil {
			result.Error = errors.New("input card is nil")
			return result
		}

		// 在 PoolCards 中查找原始牌
		var selectedCard *Card
		for _, card := range phase.PoolCards {
			if card.DeckIndex == input.Card.DeckIndex {
				selectedCard = card
				break
			}
		}
		if selectedCard == nil {
			result.Error = errors.New("card not found in tribute pool")
			return result
		}

		// 验证是否轮到该玩家选择
		if len(phase.Winners) == 0 || phase.Winners[0] != input.PlayerID {
			result.Error = fmt.Errorf("not player %d's turn to select", input.PlayerID)
			return result
		}

		receiver := input.PlayerID

		// 事件意图
		result.Events = append(result.Events, TributeEventIntent{
			Type:       EventTributeCardSelected,
			PlayerSeat: receiver,
			Card:       selectedCard,
			IsAuto:     false,
		})

		// 手牌变更：添加到 receiver
		result.HandChanges = append(result.HandChanges, HandChange{
			PlayerSeat: receiver,
			Card:       selectedCard,
			IsAdd:      true,
		})

		// Pair 变更：设置 Receiver
		result.PairUpdates = append(result.PairUpdates, PairUpdate{
			GiverSeat: findGiverByCard(phase, selectedCard),
			Receiver:  &receiver,
		})

		// 从 PoolCards 移除
		result.PoolCardToRemove = selectedCard

		return result
	}

	// 单贡场景：PoolCards == 1，自动分配
	if len(phase.PoolCards) == 1 {
		card := phase.PoolCards[0]
		receiver := tm.GetNextReceiver(phase)
		if receiver == -1 {
			result.Error = fmt.Errorf("no receiver found for pool card")
			return result
		}

		// 事件意图
		result.Events = append(result.Events, TributeEventIntent{
			Type:       EventTributeCardSelected,
			PlayerSeat: receiver,
			Card:       card,
			IsAuto:     true,
		})

		// 手牌变更：添加到 receiver
		result.HandChanges = append(result.HandChanges, HandChange{
			PlayerSeat: receiver,
			Card:       card,
			IsAdd:      true,
		})

		// Pair 变更：设置 Receiver
		result.PairUpdates = append(result.PairUpdates, PairUpdate{
			GiverSeat: findGiverByCard(phase, card),
			Receiver:  &receiver,
		})

		// 从 PoolCards 移除
		result.PoolCardToRemove = card

		return result
	}

	// PoolCards 为空，进入还贡阶段
	result.NextStatus = TributeStatusReturning
	return result
}

// processTributeReturning 处理 Returning 状态：还贡
func processTributeReturning(tm *TributeManager, phase *TributePhase, playerHands [4][]*Card, input *TributeInput) *TributeStepResult {
	result := &TributeStepResult{}

	if input == nil {
		// 检查是否有待还贡
		receiver, giver := tm.GetPendingReturnReceiver(phase)
		if receiver == -1 {
			// 全部完成，进入 Finished
			result.NextStatus = TributeStatusFinished
			return result
		}

		// 需要用户输入
		result.PendingAction = &TributeAction{
			Type:         TributeActionReturn,
			PlayerID:     receiver,
			Options:      playerHands[receiver],
			TargetPlayer: giver,
		}
		return result
	}

	// 处理用户还贡
	if input.Card == nil {
		result.Error = errors.New("input card is nil")
		return result
	}

	// 在玩家手牌中查找原始牌
	var returnCard *Card
	for _, card := range playerHands[input.PlayerID] {
		if card.DeckIndex == input.Card.DeckIndex {
			returnCard = card
			break
		}
	}
	if returnCard == nil {
		result.Error = errors.New("card not found in player hand")
		return result
	}

	// 找到对应的 TributePair
	var targetPair *TributePair
	for _, pair := range phase.TributePairs {
		if pair.Receiver == input.PlayerID {
			targetPair = pair
			break
		}
	}
	if targetPair == nil {
		result.Error = fmt.Errorf("player %d is not a tribute receiver", input.PlayerID)
		return result
	}

	if targetPair.ReturnCard != nil {
		result.Error = fmt.Errorf("player %d has already returned tribute", input.PlayerID)
		return result
	}

	giver := targetPair.Giver

	// 事件意图
	result.Events = append(result.Events, TributeEventIntent{
		Type:       EventReturnTribute,
		PlayerSeat: input.PlayerID,
		Card:       returnCard,
		TargetSeat: giver,
		IsAuto:     false,
	})

	// 手牌变更：从 receiver 移除，添加到 giver
	result.HandChanges = append(result.HandChanges, HandChange{
		PlayerSeat: input.PlayerID,
		Card:       returnCard,
		IsAdd:      false,
	})
	result.HandChanges = append(result.HandChanges, HandChange{
		PlayerSeat: giver,
		Card:       returnCard,
		IsAdd:      true,
	})

	// Pair 变更：设置 ReturnCard
	result.PairUpdates = append(result.PairUpdates, PairUpdate{
		GiverSeat:  giver,
		ReturnCard: returnCard,
	})

	return result
}

// processTributeFinished 处理 Finished 状态：验证并完成
func processTributeFinished(phase *TributePhase) *TributeStepResult {
	result := &TributeStepResult{
		PhaseCompleted: true,
	}

	// 验证贡牌完整性（非抗贡情况）
	if !phase.IsImmune {
		for _, pair := range phase.TributePairs {
			if pair.Giver == -1 || pair.Receiver == -1 || pair.TributeCard == nil || pair.ReturnCard == nil {
				result.Error = fmt.Errorf("invalid tribute pair: giver=%d, receiver=%d, tributeCard=%v, returnCard=%v",
					pair.Giver, pair.Receiver, pair.TributeCard, pair.ReturnCard)
				result.PhaseCompleted = false
				return result
			}
		}
	}

	// 发送完成事件
	result.Events = append(result.Events, TributeEventIntent{
		Type: EventTributeCompleted,
	})

	return result
}

// findGiverByCard 根据贡牌找到对应的 Giver
func findGiverByCard(phase *TributePhase, card *Card) int {
	for _, pair := range phase.TributePairs {
		if pair.TributeCard != nil && pair.TributeCard.DeckIndex == card.DeckIndex {
			return pair.Giver
		}
	}
	return -1
}
