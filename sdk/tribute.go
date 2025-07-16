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
	
	// Determine tribute requirements based on rankings
	tributeMap := make(map[int]int)
	
	// Basic tribute rules:
	// - Last place (rank 4) gives tribute to first place (rank 1)
	// - Second to last (rank 3) gives tribute to second place (rank 2)
	// - Exception: if both last places are from same team (double down), 
	//   winners choose from a pool
	
	rankings := lastResult.Rankings
	if len(rankings) < 4 {
		return nil, errors.New("invalid rankings for tribute phase")
	}
	
	rank1 := rankings[0] // First place
	rank2 := rankings[1] // Second place
	rank3 := rankings[2] // Third place
	rank4 := rankings[3] // Fourth place
	
	// Check for double down (both last places from same team)
	team3 := rank3 % 2
	team4 := rank4 % 2
	isDoubleDown := team3 == team4
	
	tributePhase := &TributePhase{
		Status:          TributeStatusWaiting,
		TributeMap:      tributeMap,
		TributeCards:    make(map[int]*Card),
		ReturnCards:     make(map[int]*Card),
		PoolCards:       make([]*Card, 0),
		SelectingPlayer: -1,
	}
	
	if isDoubleDown {
		// Double down: both losers contribute to pool, winners select
		tributePhase.Status = TributeStatusSelecting
		tributePhase.SelectingPlayer = rank1 // First place selects first
		tributePhase.SelectTimeout = time.Now().Add(3 * time.Second)
		
		// Store the rankings for later use
		tributePhase.TributeMap[rank3] = -1 // Mark as pool contributor
		tributePhase.TributeMap[rank4] = -1 // Mark as pool contributor
	} else {
		// Normal tribute: direct exchange
		tributeMap[rank4] = rank1 // Last gives to first
		tributeMap[rank3] = rank2 // Third gives to second
		tributePhase.Status = TributeStatusReturning
	}
	
	return tributePhase, nil
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
				highestCard := tm.getHighestCard(playerHands[giver])
				if highestCard != nil {
					tributePhase.TributeCards[giver] = highestCard
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
		// Create tribute pool from losing players' highest cards
		poolCards := make([]*Card, 0)
		
		for giver := range tributePhase.TributeMap {
			if tributePhase.TributeMap[giver] == -1 {
				// Get highest card from this player
				highestCard := tm.getHighestCard(playerHands[giver])
				if highestCard != nil {
					poolCards = append(poolCards, highestCard)
				}
			}
		}
		
		tributePhase.SetPoolCards(poolCards)
		tributePhase.Status = TributeStatusSelecting
	} else {
		// Normal tribute: automatically select highest cards
		for giver := range tributePhase.TributeMap {
			highestCard := tm.getHighestCard(playerHands[giver])
			if highestCard != nil {
				tributePhase.TributeCards[giver] = highestCard
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
	// For each tribute card received, return the lowest card
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
	
	rank1 := rankings[0] // First place
	rank2 := rankings[1] // Second place
	rank3 := rankings[2] // Third place
	rank4 := rankings[3] // Fourth place
	
	// Check for double down (both last places from same team)
	team3 := rank3 % 2
	team4 := rank4 % 2
	isDoubleDown := team3 == team4
	
	tributeMap := make(map[int]int)
	
	if isDoubleDown {
		// Double down: both losers contribute to pool
		tributeMap[rank3] = -1 // Mark as pool contributor
		tributeMap[rank4] = -1 // Mark as pool contributor
	} else {
		// Normal tribute: direct exchange
		tributeMap[rank4] = rank1 // Last gives to first
		tributeMap[rank3] = rank2 // Third gives to second
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
	
	// Record the selection
	tp.TributeCards[tp.SelectingPlayer] = card
	
	// Move to next selector or finish
	if len(tp.PoolCards) > 0 && tp.SelectingPlayer%2 == 0 {
		// First place selected, now second place selects
		tp.SelectingPlayer = tp.getSecondPlace()
		tp.SelectTimeout = time.Now().Add(3 * time.Second)
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
// This is a simplified implementation - in reality this would be tracked
func (tp *TributePhase) getSecondPlace() int {
	// Find the teammate of current selecting player
	if tp.SelectingPlayer%2 == 0 {
		return tp.SelectingPlayer + 2 // Teammate
	}
	return tp.SelectingPlayer - 2 // Teammate
}