package sdk

import (
	"errors"
	"fmt"
	"time"
)

// TributeManager handles all tribute-related operations independently
type TributeManager struct {
	level int
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
		TributeMap:      make(map[int]int),
		TributeCards:    make(map[int]*Card),
		ReturnCards:     make(map[int]*Card),
		PoolCards:       make([]*Card, 0),
		SelectingPlayer: -1,
	}

	// 按排名获取玩家
	rank1 := rankings[0] // 第1名
	rank3 := rankings[2] // 第3名
	rank4 := rankings[3] // 第4名

	// 根据胜利类型确定上贡规则
	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		// Double Down: rank1, rank2 同队
		// Rank3 和 Rank4 各上交 1 张贡牌，放入贡牌池
		// Rank1 优先从贡牌池中挑选其一；Rank2 获得剩下的一张贡牌
		tributePhase.Status = TributeStatusSelecting
		tributePhase.SelectingPlayer = rank1 // Rank1 先选
		tributePhase.SelectTimeout = time.Now().Add(30 * time.Second)

		tributePhase.TributeMap[rank3] = -1 // -1 表示贡献到池子
		tributePhase.TributeMap[rank4] = -1 // -1 表示贡献到池子

	case VictoryTypeSingleLast:
		// Single Last: rank1, rank3 同队
		// Rank4 上交 1 张贡牌，直接交给 Rank1
		tributePhase.Status = TributeStatusReturning
		tributePhase.TributeMap[rank4] = rank1

	case VictoryTypePartnerLast:
		// Partner Last: rank1, rank4 同队
		// Rank3 上交 1 张贡牌，直接交给 Rank1
		tributePhase.Status = TributeStatusReturning
		tributePhase.TributeMap[rank3] = rank1

	default:
		return nil, fmt.Errorf("unknown victory type: %v", lastResult.VictoryType)
	}

	return tributePhase, nil
}

// CheckTributeImmunity 检查免贡条件
func (tm *TributeManager) CheckTributeImmunity(lastResult *DealResult, playerHands [4][]*Card) bool {
	if lastResult == nil {
		return false
	}

	rankings := lastResult.Rankings
	if len(rankings) < 4 {
		return false
	}

	rank3 := rankings[2]
	rank4 := rankings[3]

	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		// Double Down：若 Rank3 和 Rank4 合计持有 两张 Big Joker，则触发免贡
		bigJokerCount := tm.countBigJokers(playerHands[rank3]) + tm.countBigJokers(playerHands[rank4])
		return bigJokerCount >= 2

	case VictoryTypeSingleLast:
		// Single Last：若 Rank4 单独持有 两张 Big Joker，则触发免贡
		return tm.countBigJokers(playerHands[rank4]) >= 2

	case VictoryTypePartnerLast:
		// Partner Last：若 Rank3 单独持有 两张 Big Joker，则触发免贡
		return tm.countBigJokers(playerHands[rank3]) >= 2

	default:
		return false
	}
}

// countBigJokers 统计手牌中大王的数量
func (tm *TributeManager) countBigJokers(hand []*Card) int {
	count := 0
	for _, card := range hand {
		if card.Number == 16 && card.Color == "Joker" { // Red Joker = Big Joker
			count++
		}
	}
	return count
}

// ProcessTribute processes the complete tribute phase with player hands
func (tm *TributeManager) ProcessTribute(tributePhase *TributePhase, playerHands [4][]*Card) error {
	if tributePhase == nil {
		return nil // No tribute phase
	}

	switch tributePhase.Status {
	case TributeStatusWaiting:
		return tm.startTributePhase(tributePhase, playerHands)
	case TributeStatusSelecting:
		return tm.handleDoubleDownTribute(tributePhase, playerHands)
	case TributeStatusReturning:
		// For normal tribute, we need to select tribute cards first if not already done
		if len(tributePhase.TributeCards) == 0 {
			for giver := range tributePhase.TributeMap {
				if tributePhase.TributeMap[giver] != -1 {
					// 自动选取除红桃Trump外最大的一张牌
					tributeCard := tm.getHighestCardExcludingHeartTrump(playerHands[giver])
					if tributeCard != nil {
						tributePhase.TributeCards[giver] = tributeCard
					}
				}
			}
		}
		return tm.processReturnCards(tributePhase, playerHands)
	default:
		return nil // Already finished
	}
}

// startTributePhase starts the tribute phase by determining tribute cards
func (tm *TributeManager) startTributePhase(tributePhase *TributePhase, playerHands [4][]*Card) error {
	// Check if this is a double down scenario
	isDoubleDown := false
	for giver := range tributePhase.TributeMap {
		if tributePhase.TributeMap[giver] == -1 {
			isDoubleDown = true
			break
		}
	}

	if isDoubleDown {
		// Create tribute pool from losing players' highest cards (excluding heart trump)
		poolCards := make([]*Card, 0)

		for giver := range tributePhase.TributeMap {
			if tributePhase.TributeMap[giver] == -1 {
				// Get highest card excluding heart trump from this player
				tributeCard := tm.getHighestCardExcludingHeartTrump(playerHands[giver])
				if tributeCard != nil {
					poolCards = append(poolCards, tributeCard)
					tributePhase.TributeCards[giver] = tributeCard
				}
			}
		}

		tributePhase.SetPoolCards(poolCards)
		tributePhase.Status = TributeStatusSelecting
	} else {
		// Normal tribute: automatically select highest cards (excluding heart trump)
		for giver := range tributePhase.TributeMap {
			if tributePhase.TributeMap[giver] != -1 {
				tributeCard := tm.getHighestCardExcludingHeartTrump(playerHands[giver])
				if tributeCard != nil {
					tributePhase.TributeCards[giver] = tributeCard
				}
			}
		}
		tributePhase.Status = TributeStatusReturning
	}

	return nil
}

// handleDoubleDownTribute handles the double down tribute selection
func (tm *TributeManager) handleDoubleDownTribute(tributePhase *TributePhase, playerHands [4][]*Card) error {
	// This method would be called when players make selections
	// For now, we'll implement auto-selection logic
	return nil
}

// processReturnCards processes the return cards phase
func (tm *TributeManager) processReturnCards(tributePhase *TributePhase, playerHands [4][]*Card) error {
	// For each tribute card received, receiver chooses a card to return
	// For simulation, we use the lowest card
	for giver, receiver := range tributePhase.TributeMap {
		if receiver != -1 && tributePhase.TributeCards[giver] != nil {
			// Find lowest card from receiver to return
			lowestCard := tm.getLowestCard(playerHands[receiver])
			if lowestCard != nil {
				tributePhase.AddReturnCard(receiver, lowestCard)
			}
		}
	}

	tributePhase.Status = TributeStatusFinished
	return nil
}

// ApplyTributeToHands applies tribute effects to player hands
func (tm *TributeManager) ApplyTributeToHands(tributePhase *TributePhase, playerHands *[4][]*Card) error {
	if tributePhase == nil || tributePhase.Status != TributeStatusFinished {
		return nil
	}

	// Apply tribute card exchanges
	for giver, receiver := range tributePhase.TributeMap {
		if receiver == -1 {
			continue // Skip pool contributors in double down
		}

		tributeCard := tributePhase.TributeCards[giver]
		if tributeCard == nil {
			continue
		}

		// Remove tribute card from giver
		playerHands[giver] = tm.removeCardFromHand(playerHands[giver], tributeCard)

		// Add tribute card to receiver
		playerHands[receiver] = append(playerHands[receiver], tributeCard)

		// Apply return card if exists
		if returnCard := tributePhase.ReturnCards[receiver]; returnCard != nil {
			// Remove return card from receiver
			playerHands[receiver] = tm.removeCardFromHand(playerHands[receiver], returnCard)

			// Add return card to giver
			playerHands[giver] = append(playerHands[giver], returnCard)
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

// getHighestCard returns the highest card from a hand
func (tm *TributeManager) getHighestCard(hand []*Card) *Card {
	if len(hand) == 0 {
		return nil
	}

	highest := hand[0]
	for _, card := range hand[1:] {
		if card.GreaterThan(highest) {
			highest = card
		}
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

// cardsEqual checks if two cards are equal
func (tm *TributeManager) cardsEqual(card1, card2 *Card) bool {
	return card1.Number == card2.Number && card1.Color == card2.Color
}

// DetermineTributeRequirements determines tribute requirements based on deal result
func (tm *TributeManager) DetermineTributeRequirements(lastResult *DealResult) (map[int]int, bool, error) {
	if lastResult == nil {
		return nil, false, nil // No tribute needed for first deal
	}

	rankings := lastResult.Rankings
	if len(rankings) < 4 {
		return nil, false, errors.New("invalid rankings for tribute determination")
	}

	tributeMap := make(map[int]int)
	isDoubleDown := lastResult.VictoryType == VictoryTypeDoubleDown

	// 按排名获取玩家
	rank1 := rankings[0] // 第1名
	rank3 := rankings[2] // 第3名
	rank4 := rankings[3] // 第4名

	// 根据胜利类型确定上贡规则
	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		// Double Down: Rank3 和 Rank4 贡献到池子
		tributeMap[rank3] = -1
		tributeMap[rank4] = -1

	case VictoryTypeSingleLast:
		// Single Last: Rank4 -> Rank1
		tributeMap[rank4] = rank1

	case VictoryTypePartnerLast:
		// Partner Last: Rank3 -> Rank1
		tributeMap[rank3] = rank1
	}

	return tributeMap, isDoubleDown, nil
}

// Start starts the tribute phase
func (tp *TributePhase) Start() error {
	if tp.Status == TributeStatusSelecting {
		// For double down, we need to create the pool from losing players' cards
		// This would be done by the Deal when it has access to player hands
		return nil
	} else if tp.Status == TributeStatusReturning {
		// For normal tribute, players automatically give their highest cards
		// This would be handled by the Deal when it has access to player hands
		return nil
	}

	return nil
}

// SelectTribute handles tribute selection from pool (double down scenario)
func (tp *TributePhase) SelectTribute(playerSeat int, card *Card) error {
	if tp.Status != TributeStatusSelecting {
		return fmt.Errorf("not in selecting status: %s", tp.Status)
	}

	if playerSeat != tp.SelectingPlayer {
		return fmt.Errorf("not player %d's turn to select", playerSeat)
	}

	// Validate card is in pool
	found := false
	for i, poolCard := range tp.PoolCards {
		if tp.cardsEqual(card, poolCard) {
			// Remove card from pool
			tp.PoolCards = append(tp.PoolCards[:i], tp.PoolCards[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("card not found in tribute pool")
	}

	// Record the selection - the card goes to the selecting player
	tp.TributeCards[tp.SelectingPlayer] = card

	// Move to next selector or finish
	if len(tp.PoolCards) > 0 {
		// Find the second place player (teammate of current selector)
		secondPlace := tp.getSecondPlace()
		tp.SelectingPlayer = secondPlace
		tp.SelectTimeout = time.Now().Add(30 * time.Second)
	} else {
		// Selection finished, move to return phase
		tp.Status = TributeStatusReturning
		tp.SelectingPlayer = -1
	}

	return nil
}

// HandleTimeout handles timeout during tribute selection
func (tp *TributePhase) HandleTimeout() error {
	if tp.Status != TributeStatusSelecting {
		return errors.New("not in selecting status")
	}

	if len(tp.PoolCards) == 0 {
		return errors.New("no cards in pool to auto-select")
	}

	// Auto-select the highest card
	highestCard := tp.PoolCards[0]
	for _, card := range tp.PoolCards {
		if card.GreaterThan(highestCard) {
			highestCard = card
		}
	}

	return tp.SelectTribute(tp.SelectingPlayer, highestCard)
}

// FinishTribute finishes the tribute phase
func (tp *TributePhase) FinishTribute() error {
	tp.Status = TributeStatusFinished
	return nil
}

// SetPoolCards sets the pool cards for double down scenario
func (tp *TributePhase) SetPoolCards(cards []*Card) {
	tp.PoolCards = make([]*Card, len(cards))
	copy(tp.PoolCards, cards)
}

// AddReturnCard adds a return card from receiver to giver
func (tp *TributePhase) AddReturnCard(receiver int, card *Card) {
	tp.ReturnCards[receiver] = card
}

// cardsEqual checks if two cards are equal
func (tp *TributePhase) cardsEqual(card1, card2 *Card) bool {
	return card1.Number == card2.Number && card1.Color == card2.Color
}

// getSecondPlace returns the seat number of second place
func (tp *TributePhase) getSecondPlace() int {
	// Find the teammate of current selecting player
	// In 4-player game: 0<->2, 1<->3 are teammates
	return (tp.SelectingPlayer + 2) % 4
}
