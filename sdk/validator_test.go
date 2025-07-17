package sdk

import (
	"testing"
	"time"
)

func TestPlayValidator_ValidatePlay(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)
	
	// Create test cards
	card1, _ := NewCard(3, "Spade", level)
	card2, _ := NewCard(3, "Heart", level)
	card3, _ := NewCard(4, "Club", level)
	card4, _ := NewCard(4, "Diamond", level)
	
	playerCards := []*Card{card1, card2, card3, card4}
	
	tests := []struct {
		name        string
		playerSeat  int
		cards       []*Card
		playerCards []*Card
		trick       *Trick
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid single card play",
			playerSeat:  0,
			cards:       []*Card{card1},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying},
			wantErr:     false,
		},
		{
			name:        "empty cards",
			playerSeat:  0,
			cards:       []*Card{},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying},
			wantErr:     true,
			errContains: "cannot play empty cards",
		},
		{
			name:        "not player's turn",
			playerSeat:  0,
			cards:       []*Card{card1},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 1, Status: TrickStatusPlaying},
			wantErr:     true,
			errContains: "not player 0's turn",
		},
		{
			name:        "player doesn't own card",
			playerSeat:  0,
			cards:       []*Card{&Card{Number: 10, Color: "Spade", Level: level}},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying},
			wantErr:     true,
			errContains: "player does not own card",
		},
		{
			name:        "valid pair play",
			playerSeat:  0,
			cards:       []*Card{card1, card2},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying},
			wantErr:     false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.ValidatePlay(tt.playerSeat, tt.cards, tt.playerCards, tt.trick)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePlay() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidatePlay() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePlay() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPlayValidator_ValidatePass(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)
	
	// Create a lead combination for testing
	leadCard, _ := NewCard(7, "Spade", level)
	leadComp := NewSingle([]*Card{leadCard})
	
	tests := []struct {
		name        string
		playerSeat  int
		trick       *Trick
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid pass",
			playerSeat: 1,
			trick: &Trick{
				CurrentTurn: 1,
				Status:      TrickStatusPlaying,
				LeadComp:    leadComp,
			},
			wantErr: false,
		},
		{
			name:        "no active trick",
			playerSeat:  1,
			trick:       nil,
			wantErr:     true,
			errContains: "no active trick",
		},
		{
			name:       "not player's turn",
			playerSeat: 1,
			trick: &Trick{
				CurrentTurn: 2,
				Status:      TrickStatusPlaying,
				LeadComp:    leadComp,
			},
			wantErr:     true,
			errContains: "not player 1's turn",
		},
		{
			name:       "cannot pass as leader",
			playerSeat: 0,
			trick: &Trick{
				CurrentTurn: 0,
				Status:      TrickStatusPlaying,
				LeadComp:    nil, // No lead combination yet
			},
			wantErr:     true,
			errContains: "cannot pass as trick leader",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.ValidatePass(tt.playerSeat, tt.trick)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePass() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidatePass() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePass() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPlayValidator_ValidateTurn(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)
	
	tests := []struct {
		name        string
		playerSeat  int
		trick       *Trick
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid turn",
			playerSeat: 2,
			trick: &Trick{
				CurrentTurn: 2,
				Status:      TrickStatusPlaying,
			},
			wantErr: false,
		},
		{
			name:        "no active trick",
			playerSeat:  2,
			trick:       nil,
			wantErr:     true,
			errContains: "no active trick",
		},
		{
			name:       "trick not playing",
			playerSeat: 2,
			trick: &Trick{
				CurrentTurn: 2,
				Status:      TrickStatusWaiting,
			},
			wantErr:     true,
			errContains: "trick is not in playing status",
		},
		{
			name:       "not player's turn",
			playerSeat: 2,
			trick: &Trick{
				CurrentTurn: 3,
				Status:      TrickStatusPlaying,
			},
			wantErr:     true,
			errContains: "not player 2's turn",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.ValidateTurn(tt.playerSeat, tt.trick)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTurn() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateTurn() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTurn() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPlayValidator_ValidateCardCombination(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)
	
	// Create test cards
	card1, _ := NewCard(3, "Spade", level)
	card2, _ := NewCard(3, "Heart", level)
	card3, _ := NewCard(4, "Club", level)
	
	tests := []struct {
		name        string
		cards       []*Card
		wantErr     bool
		errContains string
		wantType    CompType
	}{
		{
			name:     "valid single",
			cards:    []*Card{card1},
			wantErr:  false,
			wantType: TypeSingle,
		},
		{
			name:     "valid pair",
			cards:    []*Card{card1, card2},
			wantErr:  false,
			wantType: TypePair,
		},
		{
			name:        "empty cards",
			cards:       []*Card{},
			wantErr:     true,
			errContains: "cannot validate empty cards",
		},
		{
			name:        "invalid combination",
			cards:       []*Card{card1, card3}, // 3 and 4, not a valid pair
			wantErr:     true,
			errContains: "invalid card combination",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := pv.ValidateCardCombination(tt.cards)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCardCombination() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateCardCombination() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCardCombination() unexpected error = %v", err)
				} else if comp.GetType() != tt.wantType {
					t.Errorf("ValidateCardCombination() type = %v, want %v", comp.GetType(), tt.wantType)
				}
			}
		})
	}
}

func TestPlayValidator_ValidateGameState(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)
	
	// Create a valid trick
	trick, _ := NewTrick(0)
	trick.Status = TrickStatusPlaying
	
	tests := []struct {
		name        string
		deal        *Deal
		wantErr     bool
		errContains string
	}{
		{
			name: "valid game state",
			deal: &Deal{
				Status:       DealStatusPlaying,
				CurrentTrick: trick,
			},
			wantErr: false,
		},
		{
			name:        "no active deal",
			deal:        nil,
			wantErr:     true,
			errContains: "no active deal",
		},
		{
			name: "deal not playing",
			deal: &Deal{
				Status:       DealStatusWaiting,
				CurrentTrick: trick,
			},
			wantErr:     true,
			errContains: "deal is not in playing status",
		},
		{
			name: "no active trick",
			deal: &Deal{
				Status:       DealStatusPlaying,
				CurrentTrick: nil,
			},
			wantErr:     true,
			errContains: "no active trick in deal",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pv.ValidateGameState(tt.deal)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateGameState() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateGameState() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateGameState() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestTributeValidator_ValidateTributeSelection(t *testing.T) {
	level := 5
	tv := NewTributeValidator(level)
	
	// Create test cards
	card1, _ := NewCard(3, "Spade", level)
	card2, _ := NewCard(4, "Heart", level)
	poolCards := []*Card{card1, card2}
	playerCards := []*Card{card1, card2}
	
	tests := []struct {
		name         string
		playerSeat   int
		card         *Card
		tributePhase *TributePhase
		playerCards  []*Card
		wantErr      bool
		errContains  string
	}{
		{
			name:       "valid pool selection",
			playerSeat: 0,
			card:       card1,
			tributePhase: &TributePhase{
				Status:          TributeStatusSelecting,
				SelectingPlayer: 0,
				PoolCards:       poolCards,
			},
			playerCards: playerCards,
			wantErr:     false,
		},
		{
			name:         "no tribute phase",
			playerSeat:   0,
			card:         card1,
			tributePhase: nil,
			playerCards:  playerCards,
			wantErr:      true,
			errContains:  "no active tribute phase",
		},
		{
			name:       "wrong status",
			playerSeat: 0,
			card:       card1,
			tributePhase: &TributePhase{
				Status:          TributeStatusWaiting,
				SelectingPlayer: 0,
				PoolCards:       poolCards,
			},
			playerCards: playerCards,
			wantErr:     true,
			errContains: "tribute phase is not in selecting status",
		},
		{
			name:       "wrong selecting player",
			playerSeat: 1,
			card:       card1,
			tributePhase: &TributePhase{
				Status:          TributeStatusSelecting,
				SelectingPlayer: 0,
				PoolCards:       poolCards,
			},
			playerCards: playerCards,
			wantErr:     true,
			errContains: "player 1 is not the selecting player",
		},
		{
			name:       "card not in pool",
			playerSeat: 0,
			card:       &Card{Number: 10, Color: "Spade", Level: level},
			tributePhase: &TributePhase{
				Status:          TributeStatusSelecting,
				SelectingPlayer: 0,
				PoolCards:       poolCards,
			},
			playerCards: playerCards,
			wantErr:     true,
			errContains: "selected card is not in tribute pool",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tv.ValidateTributeSelection(tt.playerSeat, tt.card, tt.tributePhase, tt.playerCards)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTributeSelection() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateTributeSelection() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTributeSelection() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestTributeValidator_ValidateTributeTimeout(t *testing.T) {
	tv := NewTributeValidator(5)
	
	pastTime := time.Now().Add(-1 * time.Second)
	futureTime := time.Now().Add(1 * time.Second)
	
	tests := []struct {
		name         string
		tributePhase *TributePhase
		wantErr      bool
		errContains  string
	}{
		{
			name: "valid timeout",
			tributePhase: &TributePhase{
				Status:        TributeStatusSelecting,
				SelectTimeout: pastTime,
			},
			wantErr: false,
		},
		{
			name:         "no tribute phase",
			tributePhase: nil,
			wantErr:      true,
			errContains:  "no active tribute phase",
		},
		{
			name: "wrong status",
			tributePhase: &TributePhase{
				Status:        TributeStatusWaiting,
				SelectTimeout: pastTime,
			},
			wantErr:     true,
			errContains: "tribute phase is not in selecting status",
		},
		{
			name: "timeout not reached",
			tributePhase: &TributePhase{
				Status:        TributeStatusSelecting,
				SelectTimeout: futureTime,
			},
			wantErr:     true,
			errContains: "tribute selection timeout not reached yet",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tv.ValidateTributeTimeout(tt.tributePhase)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTributeTimeout() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateTributeTimeout() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTributeTimeout() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestGameFlowValidator_ValidatePlayerAction(t *testing.T) {
	level := 5
	gfv := NewGameFlowValidator(level)
	
	// Create test data
	card1, _ := NewCard(3, "Spade", level)
	cards := []*Card{card1}
	
	// Create valid game state
	trick, _ := NewTrick(0)
	trick.Status = TrickStatusPlaying
	trick.CurrentTurn = 0
	
	deal := &Deal{
		Level:        level,
		Status:       DealStatusPlaying,
		CurrentTrick: trick,
		PlayerCards:  [4][]*Card{{card1}, {}, {}, {}},
	}
	
	match := &Match{
		Status:      MatchStatusPlaying,
		CurrentDeal: deal,
	}
	
	gameState := &GameState{
		Status:       GameStatusStarted,
		CurrentMatch: match,
	}
	
	tests := []struct {
		name        string
		action      string
		playerSeat  int
		data        interface{}
		gameState   *GameState
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid play_cards action",
			action:     "play_cards",
			playerSeat: 0,
			data:       cards,
			gameState:  gameState,
			wantErr:    false,
		},
		{
			name:        "no active game",
			action:      "play_cards",
			playerSeat:  0,
			data:        cards,
			gameState:   nil,
			wantErr:     true,
			errContains: "no active game",
		},
		{
			name:       "invalid action type",
			action:     "invalid_action",
			playerSeat: 0,
			data:       cards,
			gameState:  gameState,
			wantErr:    true,
			errContains: "unknown action type",
		},
		{
			name:       "invalid data type for play_cards",
			action:     "play_cards",
			playerSeat: 0,
			data:       "invalid_data",
			gameState:  gameState,
			wantErr:    true,
			errContains: "invalid data type for play_cards action",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gfv.ValidatePlayerAction(tt.action, tt.playerSeat, tt.data, tt.gameState)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePlayerAction() expected error but got none")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidatePlayerAction() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePlayerAction() unexpected error = %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring (for error checking)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}