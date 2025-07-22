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

	// If there's a tribute phase, start it
	if d.TributePhase != nil {
		err = d.startTributePhase()
		if err != nil {
			return fmt.Errorf("failed to start tribute phase: %w", err)
		}
		d.Status = DealStatusTribute
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

// SelectTribute handles tribute selection during tribute phase
func (d *Deal) SelectTribute(playerSeat int, card *Card) error {
	if d.Status != DealStatusTribute {
		return fmt.Errorf("deal is not in tribute status: %s", d.Status)
	}

	if d.TributePhase == nil {
		return errors.New("no active tribute phase")
	}

	err := d.TributePhase.SelectTribute(playerSeat, card)
	if err != nil {
		return fmt.Errorf("failed to select tribute: %w", err)
	}

	// Check if tribute phase is finished
	if d.TributePhase.Status == TributeStatusFinished {
		// Apply tribute effects to player hands
		d.applyTributeEffects()

		// Start first trick
		err = d.startFirstTrick()
		if err != nil {
			return fmt.Errorf("failed to start first trick: %w", err)
		}
		d.Status = DealStatusPlaying
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
			err := d.TributePhase.HandleTimeout()
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
					d.applyTributeEffects()
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

	return d.TributePhase.Start()
}

// startFirstTrick starts the first trick of the deal
func (d *Deal) startFirstTrick() error {
	// Determine first player (usually the player with lowest level card or specific rule)
	firstPlayer := d.determineFirstPlayer()

	trick, err := NewTrick(firstPlayer)
	if err != nil {
		return fmt.Errorf("failed to create first trick: %w", err)
	}

	// Start the trick immediately
	err = trick.StartTrick()
	if err != nil {
		return fmt.Errorf("failed to start first trick: %w", err)
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

	// Trick is finished when all players have played or passed,
	// and we're back to the leader or all others have passed
	playCount := len(d.CurrentTrick.Plays)
	if playCount < 4 {
		return false
	}

	// Check if last 3 plays were all passes (everyone passed after leader)
	passCount := 0
	for i := playCount - 3; i < playCount; i++ {
		if d.CurrentTrick.Plays[i].IsPass {
			passCount++
		}
	}

	return passCount == 3
}

// finishCurrentTrick finishes the current trick and starts a new one
func (d *Deal) finishCurrentTrick() error {
	if d.CurrentTrick == nil {
		return errors.New("no current trick to finish")
	}

	// Set trick winner
	d.CurrentTrick.Winner = d.CurrentTrick.Leader
	d.CurrentTrick.Status = TrickStatusFinished

	// Add to history
	d.TrickHistory = append(d.TrickHistory, d.CurrentTrick)

	// Start new trick with winner as leader
	nextTrick, err := NewTrick(d.CurrentTrick.Winner)
	if err != nil {
		return fmt.Errorf("failed to create next trick: %w", err)
	}

	// Start the new trick immediately
	err = nextTrick.StartTrick()
	if err != nil {
		return fmt.Errorf("failed to start next trick: %w", err)
	}

	d.CurrentTrick = nextTrick
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

// applyTributeEffects applies the effects of tribute phase to player hands
func (d *Deal) applyTributeEffects() {
	if d.TributePhase == nil {
		return
	}

	// Apply tribute card exchanges
	for giver, card := range d.TributePhase.TributeCards {
		if receiver, exists := d.TributePhase.TributeMap[giver]; exists {
			// Remove tribute card from giver
			d.removeCardsFromPlayer(giver, []*Card{card})

			// Add tribute card to receiver
			d.PlayerCards[receiver] = append(d.PlayerCards[receiver], card)

			// Apply return card if exists
			if returnCard, hasReturn := d.TributePhase.ReturnCards[receiver]; hasReturn {
				// Remove return card from receiver
				d.removeCardsFromPlayer(receiver, []*Card{returnCard})

				// Add return card to giver
				d.PlayerCards[giver] = append(d.PlayerCards[giver], returnCard)
			}
		}
	}

	// Re-sort all hands
	for player := 0; player < 4; player++ {
		d.PlayerCards[player] = sortCards(d.PlayerCards[player])
	}
}
