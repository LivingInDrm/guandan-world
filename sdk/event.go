package sdk

import "time"

// ==================== 核心类型定义 ====================

// GameEventType 定义游戏事件的类型
// 游戏引擎通过这些事件类型来通知外部系统游戏状态的变化
type GameEventType string

// GameEvent 表示游戏中发生的事件及其相关数据
// 游戏引擎通过事件系统来通知外部关于游戏状态变化的信息
type GameEvent struct {
	Type       GameEventType `json:"type"`                  // 事件类型，标识这是什么类型的事件
	Data       interface{}   `json:"data"`                  // 事件数据，包含与事件相关的具体信息
	Timestamp  time.Time     `json:"timestamp"`             // 事件发生的时间戳
	PlayerSeat int           `json:"player_seat,omitempty"` // 触发事件的玩家座位号（如果适用）
}

// EventObserver 定义事件观察者接口
// 用于观察和响应游戏事件，但不影响游戏流程
type EventObserver interface {
	// OnGameEvent 处理游戏事件
	// 参数:
	//   event: 游戏事件
	// 功能说明:
	//   - 该方法应该快速执行，不应阻塞游戏流程
	//   - 主要用于日志记录、统计分析、UI更新等
	OnGameEvent(event *GameEvent)
}

// ==================== 函数式适配器 ====================

// EventHandlerFunc 是将函数适配为 EventObserver 接口的适配器类型
type EventHandlerFunc func(*GameEvent)

// OnGameEvent 实现 EventObserver 接口
func (f EventHandlerFunc) OnGameEvent(event *GameEvent) {
	f(event)
}

// GameEventHandler 是处理游戏事件的函数类型（向后兼容旧 API）
// Deprecated: 使用 EventHandlerFunc 代替
type GameEventHandler = EventHandlerFunc

// ==================== 事件常量 ====================

// 游戏事件类型常量定义
// 这些常量用于标识不同类型的游戏事件，外部系统可以通过监听这些事件来响应游戏状态变化
const (
	EventMatchStarted       GameEventType = "match_started"        // 比赛开始事件
	EventDealStarted        GameEventType = "deal_started"         // 牌局开始事件
	EventCardsDealt         GameEventType = "cards_dealt"          // 发牌完成事件
	EventTributePhase       GameEventType = "tribute_phase"        // 进贡阶段事件
	EventTributeRulesSet    GameEventType = "tribute_rules_set"    // 上贡规则确定事件
	EventTributeImmunity    GameEventType = "tribute_immunity"     // 免贡事件
	EventTributePoolCreated GameEventType = "tribute_pool_created" // 贡牌池创建事件（双下）
	EventTributeStarted     GameEventType = "tribute_started"      // 贡牌开始事件
	EventTributeGiven       GameEventType = "tribute_given"        // 上贡完成事件
	EventTributeSelected    GameEventType = "tribute_selected"     // 选牌完成事件（双下）
	EventReturnTribute      GameEventType = "return_tribute"       // 还贡完成事件
	EventTributeCompleted   GameEventType = "tribute_completed"    // 贡牌阶段结束事件
	EventTrickStarted       GameEventType = "trick_started"        // 新轮次开始事件
	EventPlayerPlayed       GameEventType = "player_played"        // 玩家出牌事件
	EventPlayerPassed       GameEventType = "player_passed"        // 玩家过牌事件
	EventTrickEnded         GameEventType = "trick_ended"          // 轮次结束事件
	EventDealEnded          GameEventType = "deal_ended"           // 牌局结束事件
	EventMatchEnded         GameEventType = "match_ended"          // 比赛结束事件
	EventPlayerTimeout      GameEventType = "player_timeout"       // 玩家超时事件
	EventPlayerDisconnect   GameEventType = "player_disconnect"    // 玩家断线事件
	EventPlayerReconnect    GameEventType = "player_reconnect"     // 玩家重连事件
)

// ==================== 事件工厂方法 ====================

// NewMatchStartedEvent 创建比赛开始事件
// Data: *Match 对象
func NewMatchStartedEvent(match *Match) *GameEvent {
	return &GameEvent{
		Type:      EventMatchStarted,
		Data:      match,
		Timestamp: time.Now(),
	}
}

// NewDealStartedEvent 创建牌局开始事件
// Data: {"deal_level": int, "team_levels": [2]int}
func NewDealStartedEvent(dealLevel int, teamLevels [2]int) *GameEvent {
	return &GameEvent{
		Type: EventDealStarted,
		Data: map[string]interface{}{
			"deal_level":  dealLevel,
			"team_levels": teamLevels,
		},
		Timestamp: time.Now(),
	}
}

// NewCardsDealtEvent 创建发牌完成事件
// 不包含具体的手牌信息，避免泄露敏感数据
// 玩家手牌通过 PlayerView 单独获取
func NewCardsDealtEvent() *GameEvent {
	return &GameEvent{
		Type: EventCardsDealt,
		Data: map[string]interface{}{
			"message": "Cards dealt to all players",
		},
		Timestamp: time.Now(),
	}
}

// NewTributePhaseEvent 创建进贡阶段事件
// Data: *TributePhase 对象
func NewTributePhaseEvent(tributePhase *TributePhase) *GameEvent {
	return &GameEvent{
		Type:      EventTributePhase,
		Data:      tributePhase,
		Timestamp: time.Now(),
	}
}

// NewTributeStartedEvent 创建贡牌开始事件
// Data: *TributePhase 对象
func NewTributeStartedEvent(tributePhase *TributePhase) *GameEvent {
	return &GameEvent{
		Type:      EventTributeStarted,
		Data:      tributePhase,
		Timestamp: time.Now(),
	}
}

// NewTributeRulesSetEvent 创建上贡规则确定事件
// Data: {
//   "last_result": *DealResult,
//   "victory_type": VictoryType,
//   "tribute_rules": {"tribute_map": map[int]int, "is_double_down": bool, "description": string},
//   "player_rankings": []int
// }
func NewTributeRulesSetEvent(lastResult *DealResult, victoryType VictoryType, tributeMap map[int]int, isDoubleDown bool, ruleDescription string, playerRankings []int) *GameEvent {
	return &GameEvent{
		Type: EventTributeRulesSet,
		Data: map[string]interface{}{
			"last_result":  lastResult,
			"victory_type": victoryType,
			"tribute_rules": map[string]interface{}{
				"tribute_map":    tributeMap,
				"is_double_down": isDoubleDown,
				"description":    ruleDescription,
			},
			"player_rankings": playerRankings,
		},
		Timestamp: time.Now(),
	}
}

// NewTributeImmunityEvent 创建免贡事件
// Data: {
//   "tribute_phase": *TributePhase,
//   "immunity_reason": map[string]interface{} (包含免贡原因的详细信息)
// }
func NewTributeImmunityEvent(tributePhase *TributePhase, immunityReason map[string]interface{}) *GameEvent {
	return &GameEvent{
		Type: EventTributeImmunity,
		Data: map[string]interface{}{
			"tribute_phase":   tributePhase,
			"immunity_reason": immunityReason,
		},
		Timestamp: time.Now(),
	}
}

// NewTributePoolCreatedEvent 创建贡牌池创建事件（双下）
// Data: {
//   "description": string,
//   "contributors": []map[string]interface{} (每个元素包含 "player_seat": int, "card": *Card),
//   "selection_order": []int,
//   "pool_cards": []*Card,
//   "selecting_player": int
// }
func NewTributePoolCreatedEvent(description string, contributors []map[string]interface{}, selectionOrder []int, poolCards []*Card, selectingPlayer int) *GameEvent {
	return &GameEvent{
		Type: EventTributePoolCreated,
		Data: map[string]interface{}{
			"description":      description,
			"contributors":     contributors,
			"selection_order":  selectionOrder,
			"pool_cards":       poolCards,
			"selecting_player": selectingPlayer,
		},
		Timestamp: time.Now(),
	}
}

// NewTributeGivenEvent 创建上贡完成事件
// Data: {"giver": int, "receiver": int, "card": *Card}
func NewTributeGivenEvent(giver int, receiver int, card *Card) *GameEvent {
	return &GameEvent{
		Type: EventTributeGiven,
		Data: map[string]interface{}{
			"giver":    giver,
			"receiver": receiver,
			"card":     card,
		},
		Timestamp: time.Now(),
	}
}

// NewTributeSelectedEvent 创建选牌完成事件（双下）
// Data: {
//   "action": string,
//   "player_seat": int,
//   "card_id": string,
//   "selected_card": *Card,
//   "remaining_options": []*Card,
//   "selection_order": int,
//   "is_timeout": bool
// }
func NewTributeSelectedEvent(action string, playerID int, cardID string, selectedCard *Card, remainingOptions []*Card, selectionOrder int, isTimeout bool) *GameEvent {
	return &GameEvent{
		Type: EventTributeSelected,
		Data: map[string]interface{}{
			"action":            action,
			"player_seat":       playerID,
			"card_id":           cardID,
			"selected_card":     selectedCard,
			"remaining_options": remainingOptions,
			"selection_order":   selectionOrder,
			"is_timeout":        isTimeout,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerID,
	}
}

// NewReturnTributeEvent 创建还贡完成事件
// Data: {
//   "player_seat": int,
//   "return_card": *Card,
//   "target_player": int
// }
func NewReturnTributeEvent(playerID int, returnCard *Card, targetPlayer int) *GameEvent {
	return &GameEvent{
		Type: EventReturnTribute,
		Data: map[string]interface{}{
			"player_seat":   playerID,
			"return_card":   returnCard,
			"target_player": targetPlayer,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerID,
	}
}

// NewTributeCompletedEvent 创建贡牌阶段结束事件
// Data: *TributePhase 对象
func NewTributeCompletedEvent(tributePhase *TributePhase) *GameEvent {
	return &GameEvent{
		Type:      EventTributeCompleted,
		Data:      tributePhase,
		Timestamp: time.Now(),
	}
}

// NewTrickStartedEvent 创建新轮次开始事件
// Data: {
//   "trick": *Trick,
//   "leader": int,
//   "current_turn": int,
// }
func NewTrickStartedEvent(trick *Trick, leader int, currentTurn int) *GameEvent {
	return &GameEvent{
		Type: EventTrickStarted,
		Data: map[string]interface{}{
			"trick":        trick,
			"leader":       leader,
			"current_turn": currentTurn,
		},
		Timestamp: time.Now(),
	}
}

// NewPlayerPlayedEvent 创建玩家出牌事件
// Data: {"player_seat": int, "cards": []*Card}
func NewPlayerPlayedEvent(playerSeat int, cards []*Card) *GameEvent {
	return &GameEvent{
		Type: EventPlayerPlayed,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"cards":       cards,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
}

// NewPlayerPassedEvent 创建玩家过牌事件
// Data: {"player_seat": int}
func NewPlayerPassedEvent(playerSeat int) *GameEvent {
	return &GameEvent{
		Type: EventPlayerPassed,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
}

// NewTrickEndedEvent 创建轮次结束事件
// Data: {"trick": *Trick, "winner": int, "next_leader": int}
func NewTrickEndedEvent(trick *Trick, winner int, nextLeader int) *GameEvent {
	return &GameEvent{
		Type: EventTrickEnded,
		Data: map[string]interface{}{
			"trick":       trick,
			"winner":      winner,
			"next_leader": nextLeader,
		},
		Timestamp: time.Now(),
	}
}

// NewDealEndedEvent 创建牌局结束事件
// Data: {
//   "deal": *Deal,
//   "result": *DealResult,
//   "rankings": []int,
//   "statistics": *DealStatistics
// }
func NewDealEndedEvent(deal *Deal, result *DealResult, rankings []int, statistics *DealStatistics) *GameEvent {
	return &GameEvent{
		Type: EventDealEnded,
		Data: map[string]interface{}{
			"deal":       deal,
			"result":     result,
			"rankings":   rankings,
			"statistics": statistics,
		},
		Timestamp: time.Now(),
	}
}

// NewMatchEndedEvent 创建比赛结束事件
// Data: {
//   "match": *Match,
//   "result": *MatchResult,
//   "winner": int,
//   "final_levels": [2]int
// }
func NewMatchEndedEvent(match *Match, result *MatchResult, winner int, finalLevels [2]int) *GameEvent {
	return &GameEvent{
		Type: EventMatchEnded,
		Data: map[string]interface{}{
			"match":        match,
			"result":       result,
			"winner":       winner,
			"final_levels": finalLevels,
		},
		Timestamp: time.Now(),
	}
}

// NewPlayerTimeoutEvent 创建玩家超时事件
// Data: {"action": string} - action类型: "play_decision", "tribute_select", "return_tribute"
func NewPlayerTimeoutEvent(seat int, actionType string) *GameEvent {
	return &GameEvent{
		Type: EventPlayerTimeout,
		Data: map[string]interface{}{
			"action": actionType,
		},
		Timestamp:  time.Now(),
		PlayerSeat: seat,
	}
}

// NewPlayerDisconnectEvent 创建玩家断线事件
// Data: {"player_seat": int, "auto_play": bool}
func NewPlayerDisconnectEvent(playerSeat int) *GameEvent {
	return &GameEvent{
		Type: EventPlayerDisconnect,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   true,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
}

// NewPlayerReconnectEvent 创建玩家重连事件
// Data: {"player_seat": int, "auto_play": bool}
func NewPlayerReconnectEvent(playerSeat int) *GameEvent {
	return &GameEvent{
		Type: EventPlayerReconnect,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   false,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
}
