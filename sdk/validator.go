package sdk

import (
	"errors"
	"fmt"
)

// ValidationError 表示可恢复的验证错误（非法出牌）
type ValidationError struct {
	PlayerSeat int
	Reason     string
	Detail     string
}

func (e *ValidationError) Error() string {
	return e.Detail
}

// 错误原因常量
const (
	ReasonNotYourTurn        = "not_your_turn"
	ReasonInvalidCards       = "invalid_cards"
	ReasonCannotBeat         = "cannot_beat"
	ReasonCardOwnership      = "card_ownership"
	ReasonInvalidCombination = "invalid_combination"
)

// PlayValidator handles all play validation logic
type PlayValidator struct {
	level int // Current deal level for wildcard validation
}

// NewPlayValidator creates a new play validator for the given level
func NewPlayValidator(level int) *PlayValidator {
	return &PlayValidator{
		level: level,
	}
}

// ValidatePlay validates a player's card play action
func (pv *PlayValidator) ValidatePlay(playerSeat int, cards []*Card, playerCards []*Card, currentTrick *Trick) error {
	// Basic validation
	if len(cards) == 0 {
		return &ValidationError{
			PlayerSeat: playerSeat,
			Reason:     ReasonInvalidCards,
			Detail:     "cannot play empty cards",
		}
	}

	// Validate it's the player's turn
	if currentTrick != nil && currentTrick.CurrentTurn != playerSeat {
		return &ValidationError{
			PlayerSeat: playerSeat,
			Reason:     ReasonNotYourTurn,
			Detail:     fmt.Sprintf("not player %d's turn, current turn is %d", playerSeat, currentTrick.CurrentTurn),
		}
	}

	// Validate cards are from player's hand
	err := pv.validatePlayerOwnsCards(cards, playerCards)
	if err != nil {
		return &ValidationError{
			PlayerSeat: playerSeat,
			Reason:     ReasonCardOwnership,
			Detail:     fmt.Sprintf("invalid card ownership: %v", err),
		}
	}

	// Create and validate card combination
	var prevComp CardComp
	if currentTrick != nil {
		prevComp = currentTrick.LeadComp
	}

	comp := FromCardList(cards, prevComp)
	if !comp.IsValid() {
		return &ValidationError{
			PlayerSeat: playerSeat,
			Reason:     ReasonInvalidCombination,
			Detail:     "invalid card combination",
		}
	}

	// If this is not the first play in the trick, validate against lead combination
	if currentTrick != nil && currentTrick.LeadComp != nil {
		err = pv.validateAgainstLeadCombination(comp, currentTrick.LeadComp)
		if err != nil {
			return &ValidationError{
				PlayerSeat: playerSeat,
				Reason:     ReasonCannotBeat,
				Detail:     fmt.Sprintf("cannot beat lead combination: %v", err),
			}
		}
	}

	return nil
}

// ValidatePass validates a player's pass action
func (pv *PlayValidator) ValidatePass(playerSeat int, currentTrick *Trick) error {
	// Validate it's the player's turn
	if currentTrick == nil {
		return errors.New("no active trick")
	}

	if currentTrick.CurrentTurn != playerSeat {
		return fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, currentTrick.CurrentTurn)
	}

	// Cannot pass if no one has played yet (must play as leader)
	if currentTrick.LeadComp == nil {
		return errors.New("cannot pass as trick leader - must play cards")
	}

	return nil
}

// ValidateTurn validates if it's a specific player's turn
func (pv *PlayValidator) ValidateTurn(playerSeat int, currentTrick *Trick) error {
	if currentTrick == nil {
		return errors.New("no active trick")
	}

	if currentTrick.Status != TrickStatusPlaying {
		return fmt.Errorf("trick is not in playing status: %s", currentTrick.Status)
	}

	if currentTrick.CurrentTurn != playerSeat {
		return fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, currentTrick.CurrentTurn)
	}

	return nil
}

// ValidateGameState validates the current game state for play actions
func (pv *PlayValidator) ValidateGameState(deal *Deal) error {
	if deal == nil {
		return errors.New("no active deal")
	}

	if deal.Status != DealStatusPlaying {
		return fmt.Errorf("deal is not in playing status: %s", deal.Status)
	}

	if deal.CurrentTrick == nil {
		return errors.New("no active trick in deal")
	}

	return nil
}

// validatePlayerOwnsCards checks if the player owns all the specified cards
func (pv *PlayValidator) validatePlayerOwnsCards(playedCards []*Card, playerCards []*Card) error {
	// Create a map of player's cards for efficient lookup
	playerCardMap := make(map[string]int) // card string -> count

	for _, card := range playerCards {
		cardKey := pv.getCardKey(card)
		playerCardMap[cardKey]++
	}

	// Check if player has all the played cards
	playedCardMap := make(map[string]int)
	for _, card := range playedCards {
		cardKey := pv.getCardKey(card)
		playedCardMap[cardKey]++
	}

	// Validate each played card
	for cardKey, playedCount := range playedCardMap {
		ownedCount, exists := playerCardMap[cardKey]
		if !exists {
			return fmt.Errorf("player does not own card: %s", cardKey)
		}
		if playedCount > ownedCount {
			return fmt.Errorf("player does not have enough of card %s: need %d, have %d",
				cardKey, playedCount, ownedCount)
		}
	}

	return nil
}

// validateAgainstLeadCombination validates if a combination can beat the current lead
func (pv *PlayValidator) validateAgainstLeadCombination(comp CardComp, leadComp CardComp) error {
	// Bombs can always be played
	if comp.IsBomb() {
		return nil
	}

	// Must be same type unless it's a bomb
	if comp.GetType() != leadComp.GetType() {
		return fmt.Errorf("combination type %s cannot beat lead type %s",
			comp.GetType().String(), leadComp.GetType().String())
	}

	// Must be greater than current lead
	if !comp.GreaterThan(leadComp) {
		return errors.New("combination is not greater than current lead")
	}

	return nil
}

// getCardKey creates a unique key for a card for comparison
func (pv *PlayValidator) getCardKey(card *Card) string {
	return fmt.Sprintf("%d_%s", card.Number, card.Color)
}
