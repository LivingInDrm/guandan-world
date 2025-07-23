package sdk

import "time"

// Player represents a game player
type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Seat     int    `json:"seat"`
	Online   bool   `json:"online"`
	AutoPlay bool   `json:"auto_play"`
}

// MatchStatus represents the status of a match
type MatchStatus string

const (
	MatchStatusWaiting  MatchStatus = "waiting"
	MatchStatusPlaying  MatchStatus = "playing"
	MatchStatusFinished MatchStatus = "finished"
)

// DealStatus represents the status of a deal
type DealStatus string

const (
	DealStatusWaiting  DealStatus = "waiting"
	DealStatusDealing  DealStatus = "dealing"
	DealStatusTribute  DealStatus = "tribute"
	DealStatusPlaying  DealStatus = "playing"
	DealStatusFinished DealStatus = "finished"
)

// TrickStatus represents the status of a trick
type TrickStatus string

const (
	TrickStatusWaiting  TrickStatus = "waiting"
	TrickStatusPlaying  TrickStatus = "playing"
	TrickStatusFinished TrickStatus = "finished"
)

// Match represents a complete match (multiple deals until someone reaches A level)
type Match struct {
	ID          string      `json:"id"`
	Status      MatchStatus `json:"status"`
	Players     [4]*Player  `json:"players"`
	CurrentDeal *Deal       `json:"current_deal"`
	DealHistory []*Deal     `json:"deal_history"`
	TeamLevels  [2]int      `json:"team_levels"` // Team 0: seats 0,2; Team 1: seats 1,3
	Winner      int         `json:"winner"`      // -1 if not finished, 0 or 1 for winning team
	StartTime   time.Time   `json:"start_time"`
	EndTime     *time.Time  `json:"end_time,omitempty"`
}

// Deal represents a single deal (one round of the game)
type Deal struct {
	ID           string        `json:"id"`
	Level        int           `json:"level"` // Current level for this deal
	Status       DealStatus    `json:"status"`
	CurrentTrick *Trick        `json:"current_trick"`
	TrickHistory []*Trick      `json:"trick_history"`
	TributePhase *TributePhase `json:"tribute_phase,omitempty"`
	PlayerCards  [4][]*Card    `json:"player_cards"` // Each player's current hand
	Rankings     []int         `json:"rankings"`     // Order of players finishing (seat numbers)
	StartTime    time.Time     `json:"start_time"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	LastResult   *DealResult   `json:"-"`             // Previous deal result (not serialized)
}

// Trick represents a single trick (one round of card plays)
type Trick struct {
	ID          string        `json:"id"`
	Leader      int           `json:"leader"`       // Seat number of trick leader
	CurrentTurn int           `json:"current_turn"` // Seat number of current player
	Plays       []*PlayAction `json:"plays"`        // All plays in this trick
	Winner      int           `json:"winner"`       // Seat number of trick winner (-1 if not finished)
	LeadComp    CardComp      `json:"lead_comp"`    // Leading card combination
	Status      TrickStatus   `json:"status"`
	StartTime   time.Time     `json:"start_time"`
	TurnTimeout time.Time     `json:"turn_timeout"` // When current player's turn times out
	NextLeader  int           `json:"next_leader"`  // Seat number of next trick leader (set when trick finishes)
}

// PlayAction represents a single play action by a player
type PlayAction struct {
	PlayerSeat int       `json:"player_seat"`
	Cards      []*Card   `json:"cards,omitempty"` // nil for pass
	Comp       CardComp  `json:"comp,omitempty"`  // Card combination played
	Timestamp  time.Time `json:"timestamp"`
	IsPass     bool      `json:"is_pass"`
}

// TributePhase represents the tribute phase of a deal
type TributePhase struct {
	Status          TributeStatus `json:"status"`
	TributeMap      map[int]int   `json:"tribute_map"`      // giver -> receiver
	TributeCards    map[int]*Card `json:"tribute_cards"`    // giver -> card
	ReturnCards     map[int]*Card `json:"return_cards"`     // receiver -> card
	PoolCards       []*Card       `json:"pool_cards"`       // Cards in tribute pool (for double-down)
	SelectingPlayer int           `json:"selecting_player"` // Player selecting from pool (-1 if none)
	SelectTimeout   time.Time     `json:"select_timeout"`   // When selection times out
	IsImmune        bool          `json:"is_immune"`        // Whether tribute was skipped due to immunity
}

// TributeStatus represents the status of tribute phase
type TributeStatus string

const (
	TributeStatusWaiting   TributeStatus = "waiting"
	TributeStatusSelecting TributeStatus = "selecting" // Selecting from pool
	TributeStatusReturning TributeStatus = "returning" // Returning tribute
	TributeStatusFinished  TributeStatus = "finished"
)

// TributeAction represents an action that needs player input during tribute phase
type TributeAction struct {
	Type         TributeActionType `json:"type"`
	PlayerID     int               `json:"player_id"`
	Options      []*Card           `json:"options"`       // Available cards to choose from
	TargetPlayer int               `json:"target_player"` // Target player for return tribute
	Timeout      time.Duration     `json:"timeout"`
}

// TributeActionType represents the type of tribute action
type TributeActionType string

const (
	TributeActionNone   TributeActionType = "none"   // No action needed
	TributeActionSelect TributeActionType = "select" // Select tribute card from pool (double-down)
	TributeActionReturn TributeActionType = "return" // Return tribute card
)

// TributeStatusInfo provides detailed information about tribute phase status
type TributeStatusInfo struct {
	Phase          TributeStatus    `json:"phase"`
	TributeCards   map[int]*Card    `json:"tribute_cards"`   // Already determined tribute cards
	ReturnCards    map[int]*Card    `json:"return_cards"`    // Already determined return cards
	PendingActions []*TributeAction `json:"pending_actions"` // Actions waiting for player input
	IsImmune       bool             `json:"is_immune"`       // Whether tribute was skipped due to immunity
}

// Methods for Match and Deal are implemented in their respective files
