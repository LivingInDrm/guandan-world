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
		Status:          TributeStatusWaiting,
		TributePairs:    make([]*TributePair, 0),
		PoolCards:       make([]*Card, 0),
		SelectingPlayer: -1,
	}

	// 按排名获取玩家
	rank1 := rankings[0] // 第1名
	rank3 := rankings[2] // 第3名
	rank4 := rankings[3] // 第4名

	// 统一使用贡牌池机制：所有上贡类型都先贡到池，然后从池中选择
	// 根据胜利类型确定上贡规则
	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		// Double Down: rank1, rank2 同队
		// Rank3 和 Rank4 各上交 1 张贡牌到贡牌池
		// Rank1 优先从贡牌池中挑选其一；Rank2 获得剩下的一张贡牌
		tributePhase.SelectingPlayer = rank1
		tributePhase.TributePairs = append(tributePhase.TributePairs,
			&TributePair{Giver: rank3, Receiver: -1, TributeCard: nil, ReturnCard: nil},
			&TributePair{Giver: rank4, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		)

	case VictoryTypeSingleLast:
		// Single Last: rank1, rank3 同队
		// Rank4 上交 1 张贡牌到贡牌池，Rank1 自动获得
		tributePhase.SelectingPlayer = rank1
		tributePhase.TributePairs = append(tributePhase.TributePairs,
			&TributePair{Giver: rank4, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		)

	case VictoryTypePartnerLast:
		// Partner Last: rank1, rank4 同队
		// Rank3 上交 1 张贡牌到贡牌池，Rank1 自动获得
		tributePhase.SelectingPlayer = rank1
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

// startTributePhase starts the tribute phase by determining tribute cards
func (tm *TributeManager) startTributePhase(tributePhase *TributePhase, playerHands [4][]*Card) error {
	// 统一处理所有上贡类型：自动选择贡牌并放入贡牌池
	poolCards := make([]*Card, 0)
	seenCardIDs := make(map[string]bool) // 防御性检查：确保不重复

	for _, pair := range tributePhase.TributePairs {
		// 自动选择贡牌（除红桃级牌外的最大牌）
		tributeCard := tm.getHighestCardExcludingHeartTrump(playerHands[pair.Giver])
		if tributeCard != nil {
			// 防御性检查：确保卡 ID 不重复
			cardID := tributeCard.GetID()
			if seenCardIDs[cardID] {
				return fmt.Errorf("duplicate card ID in tribute pool: %s", cardID)
			}
			seenCardIDs[cardID] = true

			// 填充 TributePair
			pair.TributeCard = tributeCard
			// 加入贡牌池
			poolCards = append(poolCards, tributeCard)
		}
	}

	tributePhase.setPoolCards(poolCards)
	tributePhase.Status = TributeStatusSelecting

	return nil
}

// processSelectingPhase processes automatic selection for non-double-down scenarios
// 处理单下/末游场景的自动选贡，双下场景不做处理（等待用户选择）
// 状态转换: Selecting → Returning（仅单下/末游）
// 事件发送: 由 GameEngine 通过状态变化检测发送 TributeCardSelected 事件
func (tm *TributeManager) processSelectingPhase(tributePhase *TributePhase, lastResult *DealResult) error {
	// 判断是否是双下场景（贡牌池有2张牌）
	if len(tributePhase.PoolCards) > 1 {
		// 双下场景，需要等待用户选择，不做处理
		return nil
	}

	// 单下/末游场景：rank1自动获得唯一的贡牌
	if len(tributePhase.PoolCards) == 1 {
		selectedCard := tributePhase.PoolCards[0]
		rank1 := tributePhase.SelectingPlayer

		// 找到对应的 TributePair 并更新 Receiver
		found := false
		for _, pair := range tributePhase.TributePairs {
			if pair.TributeCard != nil && pair.TributeCard.GetID() == selectedCard.GetID() {
				pair.Receiver = rank1
				found = true
				break
			}
		}

		if !found {
			// 提供详细调试信息
			poolCardIDs := make([]string, len(tributePhase.PoolCards))
			for i, card := range tributePhase.PoolCards {
				poolCardIDs[i] = card.GetID()
			}
			pairCardIDs := make([]string, len(tributePhase.TributePairs))
			for i, pair := range tributePhase.TributePairs {
				if pair.TributeCard != nil {
					pairCardIDs[i] = pair.TributeCard.GetID()
				} else {
					pairCardIDs[i] = "nil"
				}
			}
			return fmt.Errorf("could not find tribute pair for pool card: pool=%v, pairs=%v", poolCardIDs, pairCardIDs)
		}

		// 清空贡牌池
		tributePhase.PoolCards = make([]*Card, 0)

		// 状态转换到还贡阶段
		tributePhase.Status = TributeStatusReturning
	}

	return nil
}

// processReturnCards processes the return cards phase
func (tm *TributeManager) processReturnCards(tributePhase *TributePhase, playerHands [4][]*Card) error {
	// Check if all return cards have been submitted
	// 基于 TributePairs 判断是否所有还贡都已完成
	allReturned := true
	hasActivePairs := false
	for _, pair := range tributePhase.TributePairs {
		if pair.Receiver != -1 {
			hasActivePairs = true
			// Check if return card is missing
			if pair.ReturnCard == nil {
				allReturned = false
				break
			}
		}
	}

	// Only finish the phase if all returns are complete
	if allReturned && hasActivePairs {
		tributePhase.Status = TributeStatusFinished
	}

	return nil
}

// ApplyTributeToHands applies tribute effects to player hands
func (tm *TributeManager) ApplyTributeToHands(tributePhase *TributePhase, playerHands *[4][]*Card) error {
	if tributePhase == nil || tributePhase.Status != TributeStatusFinished {
		return nil
	}

	// 统一处理所有上贡场景：基于 TributePairs
	for _, pair := range tributePhase.TributePairs {
		// Step 1: 应用贡牌转移 (from giver to receiver)
		if pair.Receiver != -1 && pair.TributeCard != nil {
			// Remove tribute card from giver
			playerHands[pair.Giver] = tm.removeCardFromHand(playerHands[pair.Giver], pair.TributeCard)

			// Add tribute card to receiver
			playerHands[pair.Receiver] = append(playerHands[pair.Receiver], pair.TributeCard)
		}

		// Step 2: 应用还贡 (from receiver back to giver)
		if pair.Receiver != -1 && pair.ReturnCard != nil {
			// Remove return card from receiver
			playerHands[pair.Receiver] = tm.removeCardFromHand(playerHands[pair.Receiver], pair.ReturnCard)

			// Add return card to giver
			playerHands[pair.Giver] = append(playerHands[pair.Giver], pair.ReturnCard)
		}
	}

	// Re-sort all hands
	for player := 0; player < 4; player++ {
		playerHands[player] = sortCards(playerHands[player])
	}

	return nil
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

// getLowestCard returns the lowest card from a hand
func (tm *TributeManager) getLowestCard(hand []*Card) *Card {
	if len(hand) == 0 {
		return nil
	}

	lowest := hand[0]
	for _, card := range hand[1:] {
		if lowest.GreaterThan(card) {
			lowest = card
		}
	}

	return lowest
}

// removeCardFromHand removes a specific card from a hand
func (tm *TributeManager) removeCardFromHand(hand []*Card, cardToRemove *Card) []*Card {
	for i, card := range hand {
		if tm.cardsEqual(card, cardToRemove) {
			// Remove card by swapping with last and truncating
			hand[i] = hand[len(hand)-1]
			return hand[:len(hand)-1]
		}
	}
	return hand
}

// cardsEqual checks if two cards are equal by DeckIndex
func (tm *TributeManager) cardsEqual(card1, card2 *Card) bool {
	return card1.DeckIndex == card2.DeckIndex
}

// selectTribute handles tribute selection from pool (double down scenario)
func (tp *TributePhase) selectTribute(playerSeat int, card *Card) error {
	if tp.Status != TributeStatusSelecting {
		return fmt.Errorf("not in selecting status: %s", tp.Status)
	}

	if playerSeat != tp.SelectingPlayer {
		return fmt.Errorf("not player %d's turn to select", playerSeat)
	}

	// Validate card is in pool and find the corresponding TributePair
	found := false
	var cardIndex int = -1
	var selectedPair *TributePair = nil
	
	for i, poolCard := range tp.PoolCards {
		if tp.cardsEqual(card, poolCard) {
			cardIndex = i
			// Find the corresponding TributePair
			for _, pair := range tp.TributePairs {
				if pair.TributeCard != nil && pair.TributeCard.GetID() == poolCard.GetID() {
					selectedPair = pair
					found = true
					break
				}
			}
			break
		}
	}

	if !found || selectedPair == nil {
		return errors.New("card not found in tribute pool or no matching tribute pair")
	}

	// 记录选择结果：playerSeat 选择了来自 selectedPair.Giver 的牌
	selectedPair.Receiver = playerSeat

	// 从贡牌池移除已选牌
	tp.PoolCards = append(tp.PoolCards[:cardIndex], tp.PoolCards[cardIndex+1:]...)

	// 处理剩余牌（双下场景）
	if len(tp.PoolCards) > 0 {
		// rank2 自动获得剩余牌
		rank2 := tp.getSecondPlace()
		remainingCard := tp.PoolCards[0]
		
		// Find the remaining TributePair
		for _, pair := range tp.TributePairs {
			if pair.TributeCard != nil && pair.TributeCard.GetID() == remainingCard.GetID() {
				pair.Receiver = rank2
				break
			}
		}
		
		// 清空贡牌池
		tp.PoolCards = make([]*Card, 0)
	}

	// 选贡完成，进入还贡阶段
	tp.Status = TributeStatusReturning
	tp.SelectingPlayer = -1

	return nil
}

// getCardKey returns a unique key for a card
func (tp *TributePhase) getCardKey(card *Card) string {
	return fmt.Sprintf("%d_%s", card.Number, card.Color)
}

// setPoolCards sets the pool cards for double down scenario
func (tp *TributePhase) setPoolCards(cards []*Card) {
	tp.PoolCards = make([]*Card, len(cards))
	copy(tp.PoolCards, cards)
}

// addReturnCard adds a return card from receiver to giver
func (tp *TributePhase) addReturnCard(receiver int, card *Card) {
	// Find the TributePair where receiver matches
	for _, pair := range tp.TributePairs {
		if pair.Receiver == receiver {
			pair.ReturnCard = card
			return
		}
	}
}

// cardsEqual checks if two cards are equal by DeckIndex
func (tp *TributePhase) cardsEqual(card1, card2 *Card) bool {
	return card1.DeckIndex == card2.DeckIndex
}

// getSecondPlace returns the seat number of second place
func (tp *TributePhase) getSecondPlace() int {
	// Find the teammate of current selecting player
	// In 4-player game: 0<->2, 1<->3 are teammates
	return (tp.SelectingPlayer + 2) % 4
}

// ProcessTributePhaseAction processes the tribute phase and returns any required action
func (tm *TributeManager) ProcessTributePhaseAction(phase *TributePhase, playerCards [4][]*Card) (*TributeAction, error) {
	if phase == nil {
		return nil, nil
	}

	// Process the tribute phase based on current status
	switch phase.Status {
	case TributeStatusWaiting:
		// 自动提交贡牌到贡牌池
		err := tm.startTributePhase(phase, playerCards)
		if err != nil {
			return nil, fmt.Errorf("start tribute phase failed: %w", err)
		}
		// 状态已变为 Selecting，继续处理

	case TributeStatusSelecting:
		// 处理选贡阶段
		// 如果是双下（贡牌池有2张牌），等待 rank1 用户选择
		// 如果是单下/末游（贡牌池有1张牌），自动分配给 rank1
		if len(phase.PoolCards) == 1 {
			// 单下/末游：自动选择
			err := tm.processSelectingPhase(phase, nil)
			if err != nil {
				return nil, fmt.Errorf("process selecting phase failed: %w", err)
			}
			// 状态已变为 Returning，GameEngine 会检测状态变化并发送事件
		}
		// 双下场景：等待用户选择，不做自动处理

	case TributeStatusReturning:
		// 还贡阶段：检查是否所有还贡都已完成
		err := tm.processReturnCards(phase, playerCards)
		if err != nil {
			return nil, fmt.Errorf("process return cards failed: %w", err)
		}

	default:
		// Already finished, no processing needed
	}

	// Generate action based on current status
	switch phase.Status {
	case TributeStatusSelecting:
		// 双下选贡场景：rank1 从贡牌池中选择
		if phase.SelectingPlayer >= 0 && len(phase.PoolCards) > 1 {
			return &TributeAction{
				Type:     TributeActionSelect,
				PlayerID: phase.SelectingPlayer,
				Options:  phase.PoolCards,
			}, nil
		}

	case TributeStatusReturning:
		// 还贡阶段：找到需要还贡的 receiver
		for _, pair := range phase.TributePairs {
			if pair.Receiver != -1 && pair.ReturnCard == nil {
				// 需要还贡
				return &TributeAction{
					Type:         TributeActionReturn,
					PlayerID:     pair.Receiver,
					Options:      playerCards[pair.Receiver],
					TargetPlayer: pair.Giver,
				}, nil
			}
		}

	case TributeStatusFinished:
		// No action needed, tribute phase is complete
		return nil, nil
	}

	return nil, nil
}

// SubmitSelection handles tribute selection from pool (double down scenario)
func (tm *TributeManager) SubmitSelection(phase *TributePhase, playerID int, cardID string) error {
	if phase == nil {
		return errors.New("no tribute phase")
	}

	if phase.Status != TributeStatusSelecting {
		return errors.New("not in selecting status")
	}

	if phase.SelectingPlayer != playerID {
		return fmt.Errorf("not player %d's turn to select", playerID)
	}

	// Find the card in pool
	var selectedCard *Card
	for _, card := range phase.PoolCards {
		if card.GetID() == cardID {
			selectedCard = card
			break
		}
	}

	if selectedCard == nil {
		return errors.New("card not found in tribute pool")
	}

	// Execute selection
	return phase.selectTribute(playerID, selectedCard)
}

// SubmitReturn handles return tribute submission
func (tm *TributeManager) SubmitReturn(phase *TributePhase, playerID int, cardID string, playerCards []*Card) error {
	if phase == nil {
		return errors.New("no tribute phase")
	}

	if phase.Status != TributeStatusReturning {
		return errors.New("not in returning status")
	}

	// 验证玩家是否需要还贡（基于 TributePairs）
	var targetPair *TributePair
	for _, pair := range phase.TributePairs {
		if pair.Receiver == playerID {
			targetPair = pair
			break
		}
	}

	if targetPair == nil {
		return fmt.Errorf("player %d does not need to return tribute", playerID)
	}

	// 检查是否已还贡
	if targetPair.ReturnCard != nil {
		return fmt.Errorf("player %d has already returned tribute", playerID)
	}

	// Find the card in player's hand
	var selectedCard *Card
	for _, card := range playerCards {
		if card.GetID() == cardID {
			selectedCard = card
			break
		}
	}

	if selectedCard == nil {
		return errors.New("card not found in player's hand")
	}

	// Record the return
	phase.addReturnCard(playerID, selectedCard)

	// Check if all returns are complete
	allReturned := true
	for _, pair := range phase.TributePairs {
		if pair.Receiver != -1 && pair.ReturnCard == nil {
			allReturned = false
			break
		}
	}

	if allReturned {
		phase.Status = TributeStatusFinished
	}

	return nil
}

// GetTributeStatusInfo returns detailed tribute status information
func (tm *TributeManager) GetTributeStatusInfo(phase *TributePhase, playerCards [4][]*Card) *TributeStatusInfo {
	if phase == nil {
		return nil
	}

	// Build pending actions list
	pendingActions := make([]*TributeAction, 0)

	switch phase.Status {
	case TributeStatusSelecting:
		if phase.SelectingPlayer >= 0 && len(phase.PoolCards) > 0 {
			pendingActions = append(pendingActions, &TributeAction{
				Type:     TributeActionSelect,
				PlayerID: phase.SelectingPlayer,
				Options:  phase.PoolCards,
			})
		}

	case TributeStatusReturning:
		for _, pair := range phase.TributePairs {
			if pair.Receiver != -1 && pair.ReturnCard == nil {
				pendingActions = append(pendingActions, &TributeAction{
					Type:         TributeActionReturn,
					PlayerID:     pair.Receiver,
					Options:      playerCards[pair.Receiver],
					TargetPlayer: pair.Giver,
				})
			}
		}
	}

	return &TributeStatusInfo{
		Phase:          phase.Status,
		TributeCards:   buildTributeCardsFromPairs(phase.TributePairs),
		ReturnCards:    buildReturnCardsFromPairs(phase.TributePairs),
		TributeMap:     buildTributeMapFromPairs(phase.TributePairs),
		PendingActions: pendingActions,
		IsImmune:       phase.IsImmune,
	}
}

// buildTributeMapFromPairs builds a tribute map from TributePairs for backward compatibility
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

// buildTributeCardsFromPairs builds a tribute cards map from TributePairs for backward compatibility
func buildTributeCardsFromPairs(pairs []*TributePair) map[int]*Card {
	result := make(map[int]*Card)
	for _, pair := range pairs {
		if pair.TributeCard != nil {
			result[pair.Giver] = pair.TributeCard
		}
	}
	return result
}

// buildReturnCardsFromPairs builds a return cards map from TributePairs for backward compatibility
func buildReturnCardsFromPairs(pairs []*TributePair) map[int]*Card {
	result := make(map[int]*Card)
	for _, pair := range pairs {
		if pair.Receiver != -1 && pair.ReturnCard != nil {
			result[pair.Receiver] = pair.ReturnCard
		}
	}
	return result
}
