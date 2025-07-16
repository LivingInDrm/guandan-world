package sdk

import (
	"errors"
	"fmt"
	"time"
)

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
		Status:      TrickStatusWaiting,
		StartTime:   time.Now(),
		TurnTimeout: time.Now().Add(20 * time.Second),
	}, nil
}

// StartTrick starts the trick and sets it to playing status
func (t *Trick) StartTrick() error {
	if t.Status != TrickStatusWaiting {
		return fmt.Errorf("trick is not in waiting status: %s", t.Status)
	}
	
	t.Status = TrickStatusPlaying
	t.TurnTimeout = time.Now().Add(20 * time.Second)
	
	return nil
}

// PlayCards handles a player playing cards in this trick
func (t *Trick) PlayCards(playerSeat int, cards []*Card, comp CardComp) error {
	if t.Status != TrickStatusPlaying {
		return fmt.Errorf("trick is not in playing status: %s", t.Status)
	}
	
	if playerSeat != t.CurrentTurn {
		return fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, t.CurrentTurn)
	}
	
	if comp == nil || !comp.IsValid() {
		return errors.New("invalid card combination")
	}
	
	// If this is not the first play, validate against lead combination
	if t.LeadComp != nil && !t.canPlayCombination(comp) {
		return errors.New("card combination cannot beat current lead")
	}
	
	// Create play action
	play := &PlayAction{
		PlayerSeat: playerSeat,
		Cards:      cards,
		Comp:       comp,
		Timestamp:  time.Now(),
		IsPass:     false,
	}
	
	// Add play to trick
	t.Plays = append(t.Plays, play)
	
	// Update trick state
	if t.LeadComp == nil {
		// This is the first play, set as lead
		t.LeadComp = comp
		t.Leader = playerSeat
	} else if comp.GreaterThan(t.LeadComp) {
		// This play beats the current lead
		t.LeadComp = comp
		t.Leader = playerSeat
	}
	
	// Move to next player
	t.CurrentTurn = t.getNextPlayer(playerSeat)
	t.TurnTimeout = time.Now().Add(20 * time.Second)
	
	// Check if trick is finished
	if t.isTrickFinished() {
		return t.finishTrick()
	}
	
	return nil
}

// PassTurn handles a player passing their turn in this trick
func (t *Trick) PassTurn(playerSeat int) error {
	if t.Status != TrickStatusPlaying {
		return fmt.Errorf("trick is not in playing status: %s", t.Status)
	}
	
	if playerSeat != t.CurrentTurn {
		return fmt.Errorf("not player %d's turn, current turn is %d", playerSeat, t.CurrentTurn)
	}
	
	// Cannot pass if no one has played yet (must play as leader)
	if t.LeadComp == nil {
		return errors.New("cannot pass as trick leader")
	}
	
	// Create pass action
	play := &PlayAction{
		PlayerSeat: playerSeat,
		Cards:      nil,
		Comp:       nil,
		Timestamp:  time.Now(),
		IsPass:     true,
	}
	
	// Add pass to trick
	t.Plays = append(t.Plays, play)
	
	// Move to next player
	t.CurrentTurn = t.getNextPlayer(playerSeat)
	t.TurnTimeout = time.Now().Add(20 * time.Second)
	
	// Check if trick is finished
	if t.isTrickFinished() {
		return t.finishTrick()
	}
	
	return nil
}

// GetWinner returns the winner of the trick (-1 if not finished)
func (t *Trick) GetWinner() int {
	if t.Status != TrickStatusFinished {
		return -1
	}
	return t.Winner
}

// GetLeadingPlayer returns the current leading player
func (t *Trick) GetLeadingPlayer() int {
	return t.Leader
}

// GetCurrentPlayer returns the player whose turn it is
func (t *Trick) GetCurrentPlayer() int {
	return t.CurrentTurn
}

// GetLeadingCombination returns the current leading card combination
func (t *Trick) GetLeadingCombination() CardComp {
	return t.LeadComp
}

// GetPlays returns all plays in this trick
func (t *Trick) GetPlays() []*PlayAction {
	return t.Plays
}

// IsFinished returns true if the trick is finished
func (t *Trick) IsFinished() bool {
	return t.Status == TrickStatusFinished
}

// HasPlayerPlayed returns true if the specified player has played in this trick
func (t *Trick) HasPlayerPlayed(playerSeat int) bool {
	for _, play := range t.Plays {
		if play.PlayerSeat == playerSeat {
			return true
		}
	}
	return false
}

// GetPlayerPlay returns the play action for a specific player, or nil if not played
func (t *Trick) GetPlayerPlay(playerSeat int) *PlayAction {
	for _, play := range t.Plays {
		if play.PlayerSeat == playerSeat {
			return play
		}
	}
	return nil
}

// ProcessTimeout handles timeout for the current player
func (t *Trick) ProcessTimeout() error {
	if t.Status != TrickStatusPlaying {
		return errors.New("trick is not in playing status")
	}
	
	if time.Now().Before(t.TurnTimeout) {
		return errors.New("timeout not reached yet")
	}
	
	// Auto-pass on timeout
	return t.PassTurn(t.CurrentTurn)
}

// canPlayCombination checks if a combination can be played against the current lead
func (t *Trick) canPlayCombination(comp CardComp) bool {
	if t.LeadComp == nil {
		return true // First play, anything is valid
	}
	
	// Bombs can always be played
	if comp.IsBomb() {
		return true
	}
	
	// Must be same type unless it's a bomb
	if comp.GetType() != t.LeadComp.GetType() {
		return false
	}
	
	// Must be greater than current lead
	return comp.GreaterThan(t.LeadComp)
}

// getNextPlayer returns the next player in turn order
func (t *Trick) getNextPlayer(currentPlayer int) int {
	return (currentPlayer + 1) % 4
}

// isTrickFinished checks if the trick is finished
func (t *Trick) isTrickFinished() bool {
	// Trick is finished when all players have played or passed,
	// and we're back to the leader or all others have passed
	playCount := len(t.Plays)
	if playCount < 4 {
		return false
	}
	
	// Check if last 3 plays were all passes (everyone passed after leader)
	passCount := 0
	for i := playCount - 3; i < playCount; i++ {
		if t.Plays[i].IsPass {
			passCount++
		}
	}
	
	return passCount == 3
}

// finishTrick finishes the trick and determines the winner
func (t *Trick) finishTrick() error {
	if t.Status == TrickStatusFinished {
		return errors.New("trick is already finished")
	}
	
	// Set winner to the current leader (who played the winning combination)
	t.Winner = t.Leader
	t.Status = TrickStatusFinished
	
	return nil
}

// GetTrickSummary returns a summary of the trick for logging/debugging
func (t *Trick) GetTrickSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":           t.ID,
		"status":       t.Status,
		"leader":       t.Leader,
		"winner":       t.Winner,
		"current_turn": t.CurrentTurn,
		"play_count":   len(t.Plays),
		"start_time":   t.StartTime,
	}
	
	if t.LeadComp != nil {
		summary["lead_combination"] = t.LeadComp.String()
		summary["lead_combination_type"] = t.LeadComp.GetType().String()
	}
	
	// Add play summary
	plays := make([]map[string]interface{}, len(t.Plays))
	for i, play := range t.Plays {
		playInfo := map[string]interface{}{
			"player_seat": play.PlayerSeat,
			"is_pass":     play.IsPass,
			"timestamp":   play.Timestamp,
		}
		
		if !play.IsPass && play.Comp != nil {
			playInfo["combination"] = play.Comp.String()
			playInfo["combination_type"] = play.Comp.GetType().String()
			playInfo["card_count"] = len(play.Cards)
		}
		
		plays[i] = playInfo
	}
	summary["plays"] = plays
	
	return summary
}

// generateTrickID generates a unique ID for a trick
func generateTrickID() string {
	return fmt.Sprintf("trick_%d", time.Now().UnixNano())
}