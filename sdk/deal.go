package sdk

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// NewDeal creates a new deal with the specified level
func NewDeal(level int, lastResult *DealResult) (*Deal, error) {
	if level < 2 || level > 14 {
		return nil, fmt.Errorf("invalid level: %d", level)
	}

	deal := &Deal{
		ID:           generateDealID(),
		Level:        level,
		Status:       DealStatusWaiting,
		PlayerCards:  [4][]*Card{},
		Rankings:     make([]int, 0),
		StartTime:    time.Now(),
		TrickHistory: make([]*Trick, 0),
		LastResult:   lastResult,
	}

	// Initialize tribute phase if needed
	if lastResult != nil {
		tributePhase, err := NewTributePhase(lastResult)
		if err != nil {
			return nil, fmt.Errorf("failed to create tribute phase: %w", err)
		}
		deal.TributePhase = tributePhase
	}

	return deal, nil
}

// StartDeal starts the deal by dealing cards to all players
func (d *Deal) StartDeal() error {
	if d.Status != DealStatusWaiting {
		return fmt.Errorf("deal is not in waiting status: %s", d.Status)
	}

	// Deal cards to all players
	err := d.dealCards()
	if err != nil {
		return fmt.Errorf("failed to deal cards: %w", err)
	}

	d.Status = DealStatusDealing

	// If there's a tribute phase, check for immunity first
	if d.TributePhase != nil {
		// Check if tribute should be skipped due to immunity
		tributeManager := NewTributeManager(d.Level)
		isImmune, _ := tributeManager.GetTributeImmunityDetails(d.LastResult, d.PlayerCards)
		if isImmune {
			// 抗贡：设置状态为 Finished，但保持 Deal.Status 为 Tribute
			// 让 StepTribute 统一处理事件发送和状态转换
			d.TributePhase.Status = TributeStatusFinished
			d.TributePhase.IsImmune = true
			d.Status = DealStatusTribute
		} else {
			// Normal tribute phase
			err = d.startTributePhase()
			if err != nil {
				return fmt.Errorf("failed to start tribute phase: %w", err)
			}
			d.Status = DealStatusTribute
		}
	} else {
		// No tribute phase, start playing directly
		err = d.startFirstTrick()
		if err != nil {
			return fmt.Errorf("failed to start first trick: %w", err)
		}
		d.Status = DealStatusPlaying
	}

	return nil
}

// PlayCards handles a player playing cards
func (d *Deal) PlayCards(playerSeat int, cards []*Card) (CompType, error) {
	if d.Status != DealStatusPlaying {
		return TypeIllegal, fmt.Errorf("deal is not in playing status: %s", d.Status)
	}

	if d.CurrentTrick == nil {
		return TypeIllegal, errors.New("no active trick")
	}

	if len(cards) == 0 {
		return TypeIllegal, errors.New("cannot play empty cards")
	}

	// Validate it's the player's turn
	if d.CurrentTrick.CurrentTurn != playerSeat {
		return TypeIllegal, fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, d.CurrentTrick.CurrentTurn)
	}

	// Validate cards are from player's hand
	err := d.validatePlayerCards(playerSeat, cards)
	if err != nil {
		return TypeIllegal, fmt.Errorf("invalid cards: %w", err)
	}

	// Create card combination and validate it
	comp := FromCardList(cards, d.CurrentTrick.LeadComp)
	if !comp.IsValid() {
		return TypeIllegal, errors.New("invalid card combination")
	}

	// If this is not the first play in trick, validate against lead combination
	if d.CurrentTrick.LeadComp != nil && !comp.GreaterThan(d.CurrentTrick.LeadComp) {
		return TypeIllegal, errors.New("card combination cannot beat current lead")
	}

	// Remove cards from player's hand
	d.removeCardsFromPlayer(playerSeat, cards)

	// Add play to current trick
	play := &PlayAction{
		PlayerSeat: playerSeat,
		Cards:      cards,
		Comp:       comp,
		Timestamp:  time.Now(),
		IsPass:     false,
	}
	d.CurrentTrick.Plays = append(d.CurrentTrick.Plays, play)

	// Update trick state
	if d.CurrentTrick.LeadComp == nil {
		// This is the first play, set as lead
		d.CurrentTrick.LeadComp = comp
		d.CurrentTrick.Leader = playerSeat
	} else if comp.GreaterThan(d.CurrentTrick.LeadComp) {
		// This play beats the current lead
		d.CurrentTrick.LeadComp = comp
		d.CurrentTrick.Leader = playerSeat
	}

	// Update PlayState: current player → Played, others → Waiting
	for i := 0; i < 4; i++ {
		if i == playerSeat {
			// Check if this player finished all cards
			if len(d.PlayerCards[i]) == 0 {
				d.PlayState[i] = PlayStateFinished
				d.ActivePlayers[i] = false
				d.Rankings = append(d.Rankings, playerSeat)
				
				// Check if deal is finished
				if d.isDealFinished() {
					return comp.GetType(), d.finishDeal()
				}
			} else {
				d.PlayState[i] = PlayStatePlayed
			}
		} else if d.ActivePlayers[i] {
			// Other active players reset to waiting (need to respond to new lead)
			d.PlayState[i] = PlayStateWaiting
		}
		// Already-finished players remain PlayStateFinished
	}

	// Move to next player
	d.CurrentTrick.CurrentTurn = d.getNextPlayer(playerSeat)

	// Check if trick is finished (all players played or passed)
	if d.isTrickFinished() {
		err = d.finishCurrentTrick()
		if err != nil {
			return TypeIllegal, fmt.Errorf("failed to finish trick: %w", err)
		}
	}

	return comp.GetType(), nil
}

// PassTurn handles a player passing their turn
func (d *Deal) PassTurn(playerSeat int) error {
	if d.Status != DealStatusPlaying {
		return fmt.Errorf("deal is not in playing status: %s", d.Status)
	}

	if d.CurrentTrick == nil {
		return errors.New("no active trick")
	}

	// Validate it's the player's turn
	if d.CurrentTrick.CurrentTurn != playerSeat {
		return fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, d.CurrentTrick.CurrentTurn)
	}

	// Cannot pass if no one has played yet (must play as leader)
	if d.CurrentTrick.LeadComp == nil {
		return errors.New("cannot pass as trick leader")
	}

	// Add pass to current trick
	play := &PlayAction{
		PlayerSeat: playerSeat,
		Cards:      nil,
		Comp:       nil,
		Timestamp:  time.Now(),
		IsPass:     true,
	}
	d.CurrentTrick.Plays = append(d.CurrentTrick.Plays, play)

	// Update state: this player has passed
	d.PlayState[playerSeat] = PlayStatePassed

	// Move to next player
	d.CurrentTrick.CurrentTurn = d.getNextPlayer(playerSeat)

	// Check if trick is finished
	if d.isTrickFinished() {
		return d.finishCurrentTrick()
	}

	return nil
}

// dealCards deals 27 cards to each player
func (d *Deal) dealCards() error {
	// Create full deck (108 cards)
	deck := d.createFullDeck()

	// Shuffle deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	// Deal 27 cards to each player
	for player := 0; player < 4; player++ {
		d.PlayerCards[player] = make([]*Card, 27)
		for card := 0; card < 27; card++ {
			d.PlayerCards[player][card] = deck[player*27+card]
		}

		// Sort player's hand
		d.PlayerCards[player] = sortCards(d.PlayerCards[player])
	}

	return nil
}

// createFullDeck creates a full deck of 108 cards
func (d *Deal) createFullDeck() []*Card {
	deck := make([]*Card, 0, 108)
	deckIndex := 0 // 唯一索引计数器

	// Add regular cards (2-A) for each suit, 2 copies each
	for _, color := range Colors {
		for number := 2; number <= 14; number++ {
			for copy := 0; copy < 2; copy++ {
				card, _ := NewCard(number, color, d.Level)
				card.DeckIndex = deckIndex // 分配唯一索引
				deckIndex++
				deck = append(deck, card)
			}
		}
	}

	// Add jokers (2 small jokers + 2 big jokers)
	for copy := 0; copy < 2; copy++ {
		smallJoker, _ := NewCard(15, "Joker", d.Level)
		smallJoker.DeckIndex = deckIndex
		deckIndex++

		bigJoker, _ := NewCard(16, "Joker", d.Level)
		bigJoker.DeckIndex = deckIndex
		deckIndex++

		deck = append(deck, smallJoker, bigJoker)
	}

	return deck
}

// startTributePhase starts the tribute phase
func (d *Deal) startTributePhase() error {
	if d.TributePhase == nil {
		return errors.New("no tribute phase to start")
	}

	// No special start logic needed, tribute phase is ready to use
	return nil
}

// StartPlayingPhase 开始游戏阶段（公开方法，用于贡牌阶段结束后启动游戏）
func (d *Deal) StartPlayingPhase() error {
	if d.Status != DealStatusTribute {
		return fmt.Errorf("can only start playing phase from tribute status, current status: %s", d.Status)
	}

	// 启动第一个trick
	err := d.startFirstTrick()
	if err != nil {
		return fmt.Errorf("failed to start first trick: %w", err)
	}

	// 更新状态为playing
	d.Status = DealStatusPlaying
	return nil
}

// startFirstTrick starts the first trick of the deal
func (d *Deal) startFirstTrick() error {
	// Initialize ActivePlayers (all players have cards at start)
	d.ActivePlayers = [4]bool{true, true, true, true}

	// Initialize PlayState (all active players are waiting)
	for i := 0; i < 4; i++ {
		d.PlayState[i] = PlayStateWaiting
	}

	// Determine first player (usually the player with lowest level card or specific rule)
	firstPlayer := d.determineFirstPlayer()

	trick, err := NewTrick(firstPlayer)
	if err != nil {
		return fmt.Errorf("failed to create first trick: %w", err)
	}

	d.CurrentTrick = trick
	return nil
}

// validatePlayerCards validates that the cards belong to the player
func (d *Deal) validatePlayerCards(playerSeat int, cards []*Card) error {
	if playerSeat < 0 || playerSeat > 3 {
		return fmt.Errorf("invalid player seat: %d", playerSeat)
	}

	playerHand := d.PlayerCards[playerSeat]

	for _, card := range cards {
		found := false
		for _, handCard := range playerHand {
			if d.cardsEqual(card, handCard) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("card %s not in player %d's hand", card.String(), playerSeat)
		}
	}

	return nil
}

// removeCardsFromPlayer removes cards from a player's hand
func (d *Deal) removeCardsFromPlayer(playerSeat int, cards []*Card) {
	playerHand := d.PlayerCards[playerSeat]

	for _, cardToRemove := range cards {
		for i := len(playerHand) - 1; i >= 0; i-- {
			if d.cardsEqual(cardToRemove, playerHand[i]) {
				// Remove card by swapping with last and truncating
				playerHand[i] = playerHand[len(playerHand)-1]
				playerHand = playerHand[:len(playerHand)-1]
				break
			}
		}
	}

	d.PlayerCards[playerSeat] = playerHand
}

// cardsEqual checks if two cards are equal by DeckIndex
func (d *Deal) cardsEqual(card1, card2 *Card) bool {
	return card1.DeckIndex == card2.DeckIndex
}

// getNextPlayer returns the next player in turn order
func (d *Deal) getNextPlayer(currentPlayer int) int {
	// Find next player who is still active
	for i := 1; i <= 4; i++ {
		nextPlayer := (currentPlayer + i) % 4
		if d.ActivePlayers[nextPlayer] {
			return nextPlayer
		}
	}
	// This should not happen in normal gameplay
	return (currentPlayer + 1) % 4
}

// isDealFinished checks if the deal is finished
func (d *Deal) isDealFinished() bool {
	// 双下判断：前两名是否是同一队
	// 如果前两名都完成且属于同一队，Deal立即结束（双下胜利）
	if len(d.Rankings) >= 2 {
		rank1 := d.Rankings[0]
		rank2 := d.Rankings[1]
		// Team 0: seats 0, 2; Team 1: seats 1, 3
		team1 := rank1 % 2
		team2 := rank2 % 2
		if team1 == team2 {
			// 双下：同一队的两名玩家分别获得第1、第2名
			return true
		}
	}
	
	// 常规判断：3个玩家完成（第4个自动完成）
	return len(d.Rankings) >= 3
}

// isTrickFinished checks if the current trick is finished
// Trick ends when there are no players in "waiting" state
func (d *Deal) isTrickFinished() bool {
	if d.CurrentTrick == nil {
		return false
	}

	// Check if any player is still waiting
	for i := 0; i < 4; i++ {
		if d.PlayState[i] == PlayStateWaiting {
			return false // Still has waiting players, trick not finished
		}
	}

	// No waiting players, trick is finished
	return true
}

// finishCurrentTrick finishes the current trick and sets it up for GameEngine to handle
func (d *Deal) finishCurrentTrick() error {
	if d.CurrentTrick == nil {
		return errors.New("no current trick to finish")
	}

	// Set trick winner (Winner >= 0 indicates trick is finished)
	d.CurrentTrick.Winner = d.CurrentTrick.Leader

	// Check if deal is finished
	if d.isDealFinished() {
		// Add to history before finishing deal
		d.TrickHistory = append(d.TrickHistory, d.CurrentTrick)
		return d.finishDeal()
	}

	// Find next leader for the next trick
	nextLeader := d.CurrentTrick.Winner

	// Check if the winner finished their cards and no one followed
	if !d.ActivePlayers[nextLeader] {
		// Check if anyone followed (played non-pass after the winner)
		anyoneFollowed := false
		winnerPlayIndex := -1

		// Find when the winner played
		for i, play := range d.CurrentTrick.Plays {
			if play.PlayerSeat == d.CurrentTrick.Winner && !play.IsPass {
				winnerPlayIndex = i
				break
			}
		}

		// Check if anyone played (not passed) after the winner
		if winnerPlayIndex >= 0 {
			for i := winnerPlayIndex + 1; i < len(d.CurrentTrick.Plays); i++ {
				if !d.CurrentTrick.Plays[i].IsPass {
					anyoneFollowed = true
					break
				}
			}
		}

		// If no one followed, give priority to teammate
		if !anyoneFollowed {
			// Find teammate (0<->2, 1<->3)
			teammate := (d.CurrentTrick.Winner + 2) % 4
			if d.ActivePlayers[teammate] {
				nextLeader = teammate
			} else {
				// Teammate has no cards, find next player with cards
				for i := 1; i < 4; i++ {
					candidate := (d.CurrentTrick.Winner + i) % 4
					if d.ActivePlayers[candidate] {
						nextLeader = candidate
						break
					}
				}
			}
		} else {
			// Someone followed, use default order
			for i := 1; i < 4; i++ {
				candidate := (d.CurrentTrick.Winner + i) % 4
				if d.ActivePlayers[candidate] {
					nextLeader = candidate
					break
				}
			}
		}
	}

	// Store next leader info for GameEngine to use
	d.CurrentTrick.NextLeader = nextLeader

	// Reset PlayState for next trick
	for i := 0; i < 4; i++ {
		if d.ActivePlayers[i] {
			d.PlayState[i] = PlayStateWaiting
		} else {
			d.PlayState[i] = PlayStateFinished
		}
	}

	// Don't create new trick here - let GameEngine handle the transition
	// This ensures TrickEnded event can be properly fired
	return nil
}

// finishDeal finishes the deal and calculates the result
func (d *Deal) finishDeal() error {
	// Add remaining players to rankings
	for seat := 0; seat < 4; seat++ {
		found := false
		for _, rankedSeat := range d.Rankings {
			if rankedSeat == seat {
				found = true
				break
			}
		}
		if !found {
			d.Rankings = append(d.Rankings, seat)
		}
	}

	d.Status = DealStatusFinished
	now := time.Now()
	d.EndTime = &now

	return nil
}

// CalculateResult calculates the deal result using the result calculator
func (d *Deal) CalculateResult(match *Match) (*DealResult, error) {
	if d.Status != DealStatusFinished {
		return nil, fmt.Errorf("deal is not finished")
	}

	calculator := NewDealResultCalculator(d.Level)
	return calculator.CalculateDealResult(d, match)
}

// determineFirstPlayer determines who plays first in the deal
func (d *Deal) determineFirstPlayer() int {
	// First deal in match: truly random selection
	if d.LastResult == nil {
		return rand.Intn(4) // Random player 0-3
	}

	// Subsequent deals based on tribute results
	if d.TributePhase != nil && d.TributePhase.Status == TributeStatusFinished {
		// Check for immunity (anti-tribute situation)
		if d.TributePhase.IsImmune {
			// Anti-tribute: rank1 player starts
			return d.LastResult.Rankings[0]
		}

		// Check victory type for tribute rules
		switch d.LastResult.VictoryType {
		case VictoryTypeDoubleDown:
			// Double Down: determine who gave the bigger tribute card
			rank3 := d.LastResult.Rankings[2]
			rank4 := d.LastResult.Rankings[3]
			
			// Find tribute cards from TributePairs
			var tribute3, tribute4 *Card
			for _, pair := range d.TributePhase.TributePairs {
				if pair.Giver == rank3 {
					tribute3 = pair.TributeCard
				} else if pair.Giver == rank4 {
					tribute4 = pair.TributeCard
				}
			}

			// Compare tribute cards to determine who starts
			if tribute3 != nil && tribute4 != nil {
				if tribute3.GreaterThan(tribute4) {
					return rank3
				} else {
					return rank4
				}
			}
			// Fallback to rank3 if tribute cards not available
			return rank3

		case VictoryTypeSingleLast:
			// Single Last: rank4 gives tribute, so rank4 starts
			return d.LastResult.Rankings[3]

		case VictoryTypePartnerLast:
			// Partner Last: rank3 gives tribute, so rank3 starts
			return d.LastResult.Rankings[2]
		}
	}

	// Default: rank1 player starts (covers any edge cases)
	return d.LastResult.Rankings[0]
}

// NewTrick creates a new trick with the specified leader
func NewTrick(leader int) (*Trick, error) {
	if leader < 0 || leader > 3 {
		return nil, fmt.Errorf("invalid leader seat: %d", leader)
	}

	return &Trick{
		ID:          generateTrickID(),
		Leader:      leader,
		CurrentTurn: leader,
		Plays:       make([]*PlayAction, 0),
		Winner:      -1,
		LeadComp:    nil,
		StartTime:   time.Now(),
		NextLeader:  -1,
	}, nil
}

// generateTrickID generates a unique ID for a trick
func generateTrickID() string {
	return fmt.Sprintf("trick_%d", time.Now().UnixNano())
}
