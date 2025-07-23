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
		// Keep status as waiting - tribute phase will be started after dealing cards
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
			// Skip tribute phase due to immunity
			d.TributePhase.Status = TributeStatusFinished
			d.TributePhase.IsImmune = true
			// Start playing directly
			err = d.startFirstTrick()
			if err != nil {
				return fmt.Errorf("failed to start first trick: %w", err)
			}
			d.Status = DealStatusPlaying
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
func (d *Deal) PlayCards(playerSeat int, cards []*Card) error {
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

	// Validate cards are from player's hand
	err := d.validatePlayerCards(playerSeat, cards)
	if err != nil {
		return fmt.Errorf("invalid cards: %w", err)
	}

	// Create card combination and validate it
	comp := FromCardList(cards, d.CurrentTrick.LeadComp)
	if !comp.IsValid() {
		return errors.New("invalid card combination")
	}

	// If this is not the first play in trick, validate against lead combination
	if d.CurrentTrick.LeadComp != nil && !d.canPlayCombination(comp, d.CurrentTrick.LeadComp) {
		return errors.New("card combination cannot beat current lead")
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

	// Check if player finished (no more cards)
	if len(d.PlayerCards[playerSeat]) == 0 {
		d.Rankings = append(d.Rankings, playerSeat)

		// Check if deal is finished
		if d.isDealFinished() {
			return d.finishDeal()
		}
	}

	// Move to next player
	d.CurrentTrick.CurrentTurn = d.getNextPlayer(playerSeat)
	d.CurrentTrick.TurnTimeout = time.Now().Add(20 * time.Second)

	// Check if trick is finished (all players played or passed)
	if d.isTrickFinished() {
		err = d.finishCurrentTrick()
		if err != nil {
			return fmt.Errorf("failed to finish trick: %w", err)
		}
	}

	return nil
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

	// Move to next player
	d.CurrentTrick.CurrentTurn = d.getNextPlayer(playerSeat)
	d.CurrentTrick.TurnTimeout = time.Now().Add(20 * time.Second)

	// Check if trick is finished
	if d.isTrickFinished() {
		return d.finishCurrentTrick()
	}

	return nil
}

// ProcessTimeouts processes any pending timeouts and returns resulting events
func (d *Deal) ProcessTimeouts() []*GameEvent {
	events := make([]*GameEvent, 0)
	now := time.Now()

	// Check tribute phase timeout
	if d.Status == DealStatusTribute && d.TributePhase != nil {
		if d.TributePhase.Status == TributeStatusSelecting && now.After(d.TributePhase.SelectTimeout) {
			// Auto-select tribute on timeout
			err := d.TributePhase.handleTimeout()
			if err == nil {
				event := &GameEvent{
					Type: EventPlayerTimeout,
					Data: map[string]interface{}{
						"player_seat": d.TributePhase.SelectingPlayer,
						"action":      "tribute_select",
					},
					Timestamp:  now,
					PlayerSeat: d.TributePhase.SelectingPlayer,
				}
				events = append(events, event)

				// Check if tribute phase finished
				if d.TributePhase.Status == TributeStatusFinished {
					// Note: Tribute effects are applied by GameEngine, not here
					d.startFirstTrick()
					d.Status = DealStatusPlaying
				}
			}
		}
	}

	// Check trick timeout
	if d.Status == DealStatusPlaying && d.CurrentTrick != nil && d.CurrentTrick.Status == TrickStatusPlaying {
		if now.After(d.CurrentTrick.TurnTimeout) {
			currentPlayer := d.CurrentTrick.CurrentTurn

			// Auto-pass on timeout
			err := d.PassTurn(currentPlayer)
			if err == nil {
				event := &GameEvent{
					Type: EventPlayerTimeout,
					Data: map[string]interface{}{
						"player_seat": currentPlayer,
						"action":      "pass",
					},
					Timestamp:  now,
					PlayerSeat: currentPlayer,
				}
				events = append(events, event)
			} else {
				// If pass fails, try to play a card instead (for trick leader)
				if d.CurrentTrick.LeadComp == nil && len(d.PlayerCards[currentPlayer]) > 0 {
					// Find smallest card to play
					smallestCard := d.PlayerCards[currentPlayer][0]
					for _, card := range d.PlayerCards[currentPlayer] {
						if card.LessThan(smallestCard) {
							smallestCard = card
						}
					}

					playErr := d.PlayCards(currentPlayer, []*Card{smallestCard})
					if playErr == nil {
						event := &GameEvent{
							Type: EventPlayerTimeout,
							Data: map[string]interface{}{
								"player_seat": currentPlayer,
								"action":      "auto_play",
							},
							Timestamp:  now,
							PlayerSeat: currentPlayer,
						}
						events = append(events, event)
					}
				}
			}
		}
	}

	return events
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

	// Add regular cards (2-A) for each suit, 2 copies each
	for _, color := range Colors {
		for number := 2; number <= 14; number++ {
			for copy := 0; copy < 2; copy++ {
				card, _ := NewCard(number, color, d.Level)
				deck = append(deck, card)
			}
		}
	}

	// Add jokers (2 small jokers + 2 big jokers)
	for copy := 0; copy < 2; copy++ {
		smallJoker, _ := NewCard(15, "Joker", d.Level)
		bigJoker, _ := NewCard(16, "Joker", d.Level)
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
	// Determine first player (usually the player with lowest level card or specific rule)
	firstPlayer := d.determineFirstPlayer()

	trick, err := NewTrick(firstPlayer)
	if err != nil {
		return fmt.Errorf("failed to create first trick: %w", err)
	}

	// Do NOT start the trick immediately - leave it in Waiting status
	// The trick will be started by checkPreActionStateTransitions when first player acts
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

// canPlayCombination checks if a combination can be played against the lead
func (d *Deal) canPlayCombination(comp, leadComp CardComp) bool {
	// Must be same type unless it's a bomb
	if comp.IsBomb() {
		return true
	}

	if comp.GetType() != leadComp.GetType() {
		return false
	}

	return comp.GreaterThan(leadComp)
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

// cardsEqual checks if two cards are equal
func (d *Deal) cardsEqual(card1, card2 *Card) bool {
	return card1.Number == card2.Number && card1.Color == card2.Color
}

// getNextPlayer returns the next player in turn order
func (d *Deal) getNextPlayer(currentPlayer int) int {
	// Find next player who still has cards
	for i := 1; i <= 4; i++ {
		nextPlayer := (currentPlayer + i) % 4
		if len(d.PlayerCards[nextPlayer]) > 0 {
			return nextPlayer
		}
	}
	// This should not happen in normal gameplay
	return (currentPlayer + 1) % 4
}

// isDealFinished checks if the deal is finished
func (d *Deal) isDealFinished() bool {
	// Deal is finished when 3 players have finished (4th is automatic)
	return len(d.Rankings) >= 3
}

// isTrickFinished checks if the current trick is finished
func (d *Deal) isTrickFinished() bool {
	if d.CurrentTrick == nil {
		return false
	}

	playCount := len(d.CurrentTrick.Plays)

	// Need at least 4 plays for a trick to be finished
	if playCount < 4 {
		return false
	}

	// Case 1: Last 3 plays were all passes (everyone passed after leader)
	passCount := 0
	for i := playCount - 3; i < playCount; i++ {
		if d.CurrentTrick.Plays[i].IsPass {
			passCount++
		}
	}
	if passCount == 3 {
		return true
	}

	// Case 2: Current turn is back to the current leading player and everyone has played
	// This means a complete round has happened and it's back to the leader
	if d.CurrentTrick.CurrentTurn == d.CurrentTrick.Leader && playCount >= 4 {
		// Check if all 4 players have played at least once
		playersPlayed := make(map[int]bool)
		for _, play := range d.CurrentTrick.Plays {
			playersPlayed[play.PlayerSeat] = true
		}

		// If all 4 players have played and we're back to the leader, trick is finished
		if len(playersPlayed) == 4 {
			return true
		}
	}

	return false
}

// finishCurrentTrick finishes the current trick and sets it up for GameEngine to handle
func (d *Deal) finishCurrentTrick() error {
	if d.CurrentTrick == nil {
		return errors.New("no current trick to finish")
	}

	// Set trick winner and status
	d.CurrentTrick.Winner = d.CurrentTrick.Leader
	d.CurrentTrick.Status = TrickStatusFinished

	// Check if deal is finished
	if d.isDealFinished() {
		// Add to history before finishing deal
		d.TrickHistory = append(d.TrickHistory, d.CurrentTrick)
		return d.finishDeal()
	}

	// Find next leader for the next trick
	nextLeader := d.CurrentTrick.Winner
	if len(d.PlayerCards[nextLeader]) == 0 {
		// Winner has no cards, find next player with cards
		for i := 1; i < 4; i++ {
			candidate := (d.CurrentTrick.Winner + i) % 4
			if len(d.PlayerCards[candidate]) > 0 {
				nextLeader = candidate
				break
			}
		}
	}

	// Store next leader info for GameEngine to use
	d.CurrentTrick.NextLeader = nextLeader

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
	// Simplified: return player 0 for now
	// In full implementation, this would consider tribute phase results
	// or find player with specific card (like 2 of hearts)
	return 0
}
