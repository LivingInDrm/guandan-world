package sdk

import (
	"testing"
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
			name:        "player does not own card",
			playerSeat:  0,
			cards:       []*Card{{Number: 10, Color: "Diamond", Level: level}},
			playerCards: playerCards,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying},
			wantErr:     true,
			errContains: "player does not own card",
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
			trick:      &Trick{CurrentTurn: 1, Status: TrickStatusPlaying, Leader: 0},
			wantErr:    false,
		},
		{
			name:        "no active trick",
			playerSeat:  1,
			trick:       nil,
			wantErr:     true,
			errContains: "no active trick",
		},
		{
			name:        "trick not playing",
			playerSeat:  1,
			trick:       &Trick{CurrentTurn: 1, Status: TrickStatusWaiting},
			wantErr:     true,
			errContains: "trick is not in playing status",
		},
		{
			name:        "not player's turn",
			playerSeat:  1,
			trick:       &Trick{CurrentTurn: 2, Status: TrickStatusPlaying},
			wantErr:     true,
			errContains: "not player 1's turn",
		},
		{
			name:        "cannot pass as trick leader",
			playerSeat:  0,
			trick:       &Trick{CurrentTurn: 0, Status: TrickStatusPlaying, Leader: 0},
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
			playerSeat: 1,
			trick:      &Trick{CurrentTurn: 1, Status: TrickStatusPlaying},
			wantErr:    false,
		},
		{
			name:        "no active trick",
			playerSeat:  1,
			trick:       nil,
			wantErr:     true,
			errContains: "no active trick",
		},
		{
			name:        "trick not playing",
			playerSeat:  1,
			trick:       &Trick{CurrentTurn: 1, Status: TrickStatusWaiting},
			wantErr:     true,
			errContains: "trick is not in playing status",
		},
		{
			name:        "not player's turn",
			playerSeat:  2,
			trick:       &Trick{CurrentTurn: 1, Status: TrickStatusPlaying},
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

func TestPlayValidator_ValidateGameState(t *testing.T) {
	level := 5
	pv := NewPlayValidator(level)

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
				CurrentTrick: &Trick{},
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
			name: "deal is not in playing status",
			deal: &Deal{
				Status: DealStatusWaiting,
			},
			wantErr:     true,
			errContains: "deal is not in playing status",
		},
		{
			name: "no active trick in deal",
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

// Helper function to check if a string contains a substring (for error checking)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (substr == "" ||
		containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
