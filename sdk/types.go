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

// PlayState represents a player's state in the current trick
type PlayState int

const (
	PlayStateWaiting  PlayState = 0 // Waiting to play or pass
	PlayStatePlayed   PlayState = 1 // Has played cards (current leader)
	PlayStatePassed   PlayState = 2 // Has passed
	PlayStateFinished PlayState = 3 // Finished all cards (no longer participating)
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
	ID            string        `json:"id"`
	Level         int           `json:"level"` // Current level for this deal
	Status        DealStatus    `json:"status"`
	CurrentTrick  *Trick        `json:"current_trick"`
	TrickHistory  []*Trick      `json:"trick_history"`
	TributePhase  *TributePhase `json:"tribute_phase,omitempty"`
	PlayerCards   [4][]*Card    `json:"player_cards"` // Each player's current hand
	Rankings      []int         `json:"rankings"`     // Order of players finishing (seat numbers)
	StartTime     time.Time     `json:"start_time"`
	EndTime       *time.Time    `json:"end_time,omitempty"`
	LastResult    *DealResult   `json:"-"`              // Previous deal result (not serialized)
	ActivePlayers [4]bool       `json:"active_players"` // Whether each seat still has cards
	PlayState     [4]PlayState  `json:"play_state"`     // Each player's state in current trick
}

// Trick represents a single trick (one round of card plays)
type Trick struct {
	ID          string        `json:"id"`
	Leader      int           `json:"leader"`       // Seat number of trick leader
	CurrentTurn int           `json:"current_turn"` // Seat number of current player
	Plays       []*PlayAction `json:"plays"`        // All plays in this trick
	Winner      int           `json:"winner"`       // Seat number of trick winner (-1 if not finished)
	LeadComp    CardComp      `json:"lead_comp"`    // Leading card combination
	Started     bool          `json:"started"`      // Whether trick_started event has been sent
	StartTime   time.Time     `json:"start_time"`
	NextLeader  int           `json:"next_leader"` // Seat number of next trick leader (set when trick finishes)
}

// PlayAction represents a single play action by a player
type PlayAction struct {
	PlayerSeat int       `json:"player_seat"`
	Cards      []*Card   `json:"cards,omitempty"` // nil for pass
	Comp       CardComp  `json:"comp,omitempty"`  // Card combination played
	Timestamp  time.Time `json:"timestamp"`
	IsPass     bool      `json:"is_pass"`
}

// TributePair records the complete lifecycle of a single tribute card
type TributePair struct {
	Giver       int   `json:"giver"`        // Seat number of the player giving tribute
	Receiver    int   `json:"receiver"`     // Seat number of the player receiving tribute (-1 if pending/pool)
	TributeCard *Card `json:"tribute_card"` // The tribute card being given
	ReturnCard  *Card `json:"return_card"`  // The return tribute card (nil if not yet returned)
}

// TributePhase represents the tribute phase of a deal
type TributePhase struct {
	Status       TributeStatus  `json:"status"`
	TributeType  string         `json:"tribute_type"`  // Tribute type: "double_down", "single_last", "partner_last"
	Givers       []int          `json:"givers"`        // Tribute givers (losing players)
	Receivers    []int          `json:"receivers"`     // Tribute receivers (winning players)
	TributePairs []*TributePair `json:"tribute_pairs"` // All tribute relationships (single source of truth)
	PoolCards    []*Card        `json:"pool_cards"`    // Cards in tribute pool (for double-down selection UI)
	IsImmune     bool           `json:"is_immune"`     // Whether tribute was skipped due to immunity
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
	Timeout      time.Duration     `json:"timeout"`       // DRIVER-MANAGED: Set by GameDriver, not by Tribute layer
}

// TributeActionType represents the type of tribute action
type TributeActionType string

const (
	TributeActionNone   TributeActionType = "none"   // No action needed
	TributeActionSelect TributeActionType = "select" // Select tribute card from pool (double-down)
	TributeActionReturn TributeActionType = "return" // Return tribute card
)

// TributeInput represents user input for tribute phase
type TributeInput struct {
	PlayerID int   `json:"player_id"` // Player seat number (0-3)
	Card     *Card `json:"card"`      // Selected card
}

// TributeStepResult 纯函数返回值，描述"应该执行的操作"
type TributeStepResult struct {
	NextStatus       TributeStatus        // 状态转换目标（空串表示不变）
	PhaseCompleted   bool                 // 贡牌阶段是否完成
	Events           []TributeEventIntent // 待发送事件
	PendingAction    *TributeAction       // 需要用户输入时返回
	HandChanges      []HandChange         // 手牌变更
	PairUpdates      []PairUpdate         // TributePair 变更
	PoolCardsToSet   []*Card              // 非nil时设置 phase.PoolCards（Waiting 阶段）
	PoolCardToRemove *Card                // 非nil时从 phase.PoolCards 移除
	Error            error
}

// TributeEventIntent 事件意图
type TributeEventIntent struct {
	Type       GameEventType // 复用已有: EventTributeCardSubmitted 等
	PlayerSeat int
	Card       *Card
	TargetSeat int
	IsAuto     bool
}

// HandChange 手牌变更意图
type HandChange struct {
	PlayerSeat int
	Card       *Card
	IsAdd      bool // true=添加, false=移除
}

// PairUpdate TributePair 变更意图
type PairUpdate struct {
	GiverSeat   int   // 用于匹配 TributePair
	Receiver    *int  // 非nil时更新 pair.Receiver
	TributeCard *Card // 非nil时更新 pair.TributeCard
	ReturnCard  *Card // 非nil时更新 pair.ReturnCard
}

// StepResult 单步贡牌处理结果（由 StepTribute 返回）
type StepResult struct {
	Action        *TributeAction `json:"action,omitempty"`  // 非nil表示需要用户输入
	Completed     bool           `json:"completed"`         // 贡牌阶段完成
	PrevStatus    TributeStatus  `json:"prev_status"`       // 本步开始时的状态
	StatusChanged bool           `json:"status_changed"`    // 是否发生状态转换
}

// Methods for Match and Deal are implemented in their respective files

// TurnInfo provides information about the current turn state
// 提供当前轮次状态信息，替代直接访问deal.CurrentTrick的需求
type TurnInfo struct {
	CurrentPlayer  int      `json:"current_player"`   // 当前轮到哪个玩家 (座位号0-3)
	IsLeader       bool     `json:"is_leader"`        // 当前玩家是否为首出
	IsNewTrick     bool     `json:"is_new_trick"`     // 是否为新trick的第一次出牌
	HasActiveTrick bool     `json:"has_active_trick"` // 是否有活跃的trick
	LeadComp       CardComp `json:"lead_comp"`        // 当前领先的牌组合 (如果有的话)
}

// MatchDetails provides comprehensive match information
// 提供完整的比赛信息，替代直接访问match对象的需求
type MatchDetails struct {
	TeamLevels [2]int        `json:"team_levels"` // 两队当前等级 [team0, team1]
	Players    []*PlayerInfo `json:"players"`     // 所有玩家信息
}

// PlayerInfo provides player information with team assignment
// 提供玩家信息包括队伍分配
type PlayerInfo struct {
	Seat     int    `json:"seat"`     // 座位号 (0-3)
	Username string `json:"username"` // 玩家用户名
	TeamNum  int    `json:"team_num"` // 队伍编号 (0 或 1)
}

// TrickInfo provides trick context for AutoPlayAlgorithm
// 为自动出牌算法提供trick上下文信息，避免直接传递Trick对象
type TrickInfo struct {
	IsLeader bool     `json:"is_leader"`           // 是否为首出
	LeadComp CardComp `json:"lead_comp,omitempty"` // 当前领先的牌组合 (如果不是首出)
}

// PlayerTimeoutStats 玩家超时统计信息
type PlayerTimeoutStats struct {
	PlayDecisionTimeouts int `json:"play_decision_timeouts"` // 出牌决策超时次数
	TributeTimeouts      int `json:"tribute_timeouts"`       // 贡牌超时次数
	TotalTimeouts        int `json:"total_timeouts"`         // 总超时次数
}
