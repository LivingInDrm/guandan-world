package sdk

import (
	"sync/atomic"
	"time"

	commonpb "guandan-world/proto/common"
	eventpb "guandan-world/proto/event"
)

// ==================== 核心类型定义 ====================

// GameEvent is now an alias to proto GameEvent
type GameEvent = eventpb.GameEvent

// GameEventType is an alias to proto EventType for backward compatibility
type GameEventType = eventpb.EventType

// Event type constants for backward compatibility
const (
	EventMatchStarted          = eventpb.EventType_EVENT_TYPE_MATCH_STARTED
	EventMatchEnded            = eventpb.EventType_EVENT_TYPE_MATCH_ENDED
	EventDealStarted           = eventpb.EventType_EVENT_TYPE_DEAL_STARTED
	EventCardsDealt            = eventpb.EventType_EVENT_TYPE_CARDS_DEALT
	EventDealEnded             = eventpb.EventType_EVENT_TYPE_DEAL_ENDED
	EventTributePhaseStarted   = eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED
	EventTributeStarted        = eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED // Alias
	EventTributeExempted       = eventpb.EventType_EVENT_TYPE_TRIBUTE_EXEMPTED
	EventTributeImmunity       = eventpb.EventType_EVENT_TYPE_TRIBUTE_EXEMPTED // Alias
	EventTributeCardSubmitted  = eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SUBMITTED
	EventTributeGiven          = eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SUBMITTED // Alias
	EventTributeCardSelected   = eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SELECTED
	EventTributeSelected       = eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SELECTED // Alias
	EventReturnTribute         = eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_RETURNED
	EventTributeCompleted      = eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED
	EventTrickStarted          = eventpb.EventType_EVENT_TYPE_TRICK_STARTED
	EventTrickEnded            = eventpb.EventType_EVENT_TYPE_TRICK_ENDED
	EventPlayerPlayed          = eventpb.EventType_EVENT_TYPE_PLAYER_PLAYED
	EventPlayerPassed          = eventpb.EventType_EVENT_TYPE_PLAYER_PASSED
	EventPlayerTimeout         = eventpb.EventType_EVENT_TYPE_PLAYER_TIMEOUT
	EventPlayerDisconnect      = eventpb.EventType_EVENT_TYPE_PLAYER_DISCONNECT
	EventPlayerReconnect       = eventpb.EventType_EVENT_TYPE_PLAYER_RECONNECT
	EventPlayerFinished        = eventpb.EventType_EVENT_TYPE_PLAYER_FINISHED
	
	// Legacy constants that don't have proto equivalents (kept for compatibility)
	EventTributeRulesSet    GameEventType = 0  // Deprecated
	EventTributePoolCreated GameEventType = 0  // Deprecated
	EventTributePhase       GameEventType = eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED
)

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

// ==================== 事件元数据提供者 ====================

// EventMetadataProvider manages event metadata generation
type EventMetadataProvider struct {
	seqCounter int64
}

// NewEventMetadataProvider creates a new metadata provider
func NewEventMetadataProvider() *EventMetadataProvider {
	return &EventMetadataProvider{
		seqCounter: 0,
	}
}

// NextSeq returns the next sequence number
func (emp *EventMetadataProvider) NextSeq() int64 {
	return atomic.AddInt64(&emp.seqCounter, 1)
}

// FillMetadata fills event metadata fields
// Parameters:
//   - event: the proto event to fill
//   - match: current match (required for match_id)
//   - deal: current deal (optional, for deal_index)
//   - trick: current trick (optional, for trick_index)
//   - actorSeat: actor seat (-1 if not applicable)
func (emp *EventMetadataProvider) FillMetadata(
	event *eventpb.GameEvent,
	match *Match,
	deal *Deal,
	trick *Trick,
	actorSeat int,
) {
	if event == nil {
		return
	}

	// Fill match_id
	if match != nil {
		event.MatchId = match.ID
	}

	// Fill deal_index
	if deal != nil && match != nil {
		dealIdx := int32(len(match.DealHistory))
		event.DealIndex = &dealIdx
	}

	// Fill trick_index
	if trick != nil && deal != nil {
		trickIdx := int32(len(deal.TrickHistory))
		event.TrickIndex = &trickIdx
	}

	// Fill actor_seat
	if actorSeat >= 0 && actorSeat < 4 {
		actorSeatInt := int32(actorSeat)
		event.ActorSeat = &actorSeatInt
	}

	// Fill seq and timestamp
	event.Seq = emp.NextSeq()
	event.CreatedAtMs = time.Now().UnixMilli()
}

// CreateBaseEvent creates a base event with type and metadata
func (emp *EventMetadataProvider) CreateBaseEvent(
	eventType eventpb.EventType,
	match *Match,
	deal *Deal,
	trick *Trick,
	actorSeat int,
) *eventpb.GameEvent {
	event := &eventpb.GameEvent{
		Type: eventType,
	}
	emp.FillMetadata(event, match, deal, trick, actorSeat)
	return event
}

// ==================== 事件工厂方法 ====================

// NewMatchStartedEvent 创建比赛开始事件
func NewMatchStartedEvent(
	emp *EventMetadataProvider,
	match *Match,
	players []Player,
	initialLevels [2]int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_MATCH_STARTED,
		match, nil, nil, -1,
	)

	// Convert []Player to []*commonpb.PlayerBasicInfo
	protoPlayers := make([]*commonpb.PlayerBasicInfo, len(players))
	for i, p := range players {
		protoPlayers[i] = &commonpb.PlayerBasicInfo{
			Id:       p.ID,
			Username: p.Username,
			Seat:     int32(p.Seat),
			TeamNum:  int32(match.GetTeamForPlayer(p.Seat)),
		}
	}

	event.Payload = &eventpb.GameEvent_MatchStarted{
		MatchStarted: &eventpb.MatchStartedPayload{
			Players:       protoPlayers,
			InitialLevels: []int32{int32(initialLevels[0]), int32(initialLevels[1])},
		},
	}

	return event
}

// NewDealStartedEvent 创建牌局开始事件
func NewDealStartedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	dealLevel int,
	teamLevels [2]int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_DEAL_STARTED,
		match, deal, nil, -1,
	)

	event.Payload = &eventpb.GameEvent_DealStarted{
		DealStarted: &eventpb.DealStartedPayload{
			DealLevel:  int32(dealLevel),
			TeamLevels: []int32{int32(teamLevels[0]), int32(teamLevels[1])},
		},
	}

	return event
}

// NewCardsDealtEvent 创建发牌完成事件
func NewCardsDealtEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	cardCount int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_CARDS_DEALT,
		match, deal, nil, -1,
	)

	event.Payload = &eventpb.GameEvent_CardsDealt{
		CardsDealt: &eventpb.CardsDealtPayload{
			CardCount: int32(cardCount),
		},
	}

	return event
}

// NewTributeStartedEvent 创建进贡阶段开始事件
func NewTributeStartedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	tributeType string,
	givers []int,
	receivers []int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED,
		match, deal, nil, -1,
	)

	givers32 := make([]int32, len(givers))
	for i, g := range givers {
		givers32[i] = int32(g)
	}

	receivers32 := make([]int32, len(receivers))
	for i, r := range receivers {
		receivers32[i] = int32(r)
	}

	event.Payload = &eventpb.GameEvent_TributeStarted{
		TributeStarted: &eventpb.TributeStartedPayload{
			TributeType: ConvertTributeTypeToProto(tributeType),
			Givers:      givers32,
			Receivers:   receivers32,
		},
	}

	return event
}

// NewTributeExemptedEvent 创建免贡事件
func NewTributeExemptedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	jokerHolders map[int]int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_EXEMPTED,
		match, deal, nil, -1,
	)

	jokerHolders32 := make(map[int32]int32)
	for seat, count := range jokerHolders {
		jokerHolders32[int32(seat)] = int32(count)
	}

	event.Payload = &eventpb.GameEvent_TributeExempted{
		TributeExempted: &eventpb.TributeExemptedPayload{
			BigJokerHolders: jokerHolders32,
		},
	}

	return event
}

// NewTributeCardSubmittedEvent 创建贡牌提交事件
func NewTributeCardSubmittedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	giverSeat int,
	card *Card,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SUBMITTED,
		match, deal, nil, giverSeat,
	)

	event.Payload = &eventpb.GameEvent_TributeCardSubmitted{
		TributeCardSubmitted: &eventpb.TributeCardSubmittedPayload{
			SubmittedCard: ConvertCardToProto(card),
		},
	}

	return event
}

// NewTributeCardSelectedEvent 创建选牌事件（双下）
func NewTributeCardSelectedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	actorSeat int,
	selectedCard *Card,
	isAuto bool,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SELECTED,
		match, deal, nil, actorSeat,
	)

	event.Payload = &eventpb.GameEvent_TributeCardSelected{
		TributeCardSelected: &eventpb.TributeCardSelectedPayload{
			SelectedCard: ConvertCardToProto(selectedCard),
			IsAuto:       isAuto,
		},
	}

	return event
}

// NewTributeCardReturnedEvent 创建还贡事件
func NewTributeCardReturnedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	actorSeat int,
	returnedCard *Card,
	targetPlayer int,
	isAuto bool,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_RETURNED,
		match, deal, nil, actorSeat,
	)

	event.Payload = &eventpb.GameEvent_TributeCardReturned{
		TributeCardReturned: &eventpb.TributeCardReturnedPayload{
			ReturnedCard: ConvertCardToProto(returnedCard),
			TargetPlayer: int32(targetPlayer),
			IsAuto:       isAuto,
		},
	}

	return event
}

// NewTributeCompletedEvent 创建进贡阶段完成事件
func NewTributeCompletedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED,
		match, deal, nil, -1,
	)

	event.Payload = &eventpb.GameEvent_TributeCompleted{
		TributeCompleted: &eventpb.TributeCompletedPayload{},
	}

	return event
}

// NewTrickStartedEvent 创建轮次开始事件
func NewTrickStartedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	leader int,
	isFirstTrick bool,
	remainingPlayers []int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRICK_STARTED,
		match, deal, trick, -1,
	)

	remainingPlayers32 := make([]int32, len(remainingPlayers))
	for i, p := range remainingPlayers {
		remainingPlayers32[i] = int32(p)
	}

	event.Payload = &eventpb.GameEvent_TrickStarted{
		TrickStarted: &eventpb.TrickStartedPayload{
			Leader:           int32(leader),
			IsFirstTrick:     isFirstTrick,
			RemainingPlayers: remainingPlayers32,
		},
	}

	return event
}

// NewTrickEndedEvent 创建轮次结束事件
func NewTrickEndedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	trickWinner int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_TRICK_ENDED,
		match, deal, trick, -1,
	)

	event.Payload = &eventpb.GameEvent_TrickEnded{
		TrickEnded: &eventpb.TrickEndedPayload{
			TrickWinner: int32(trickWinner),
		},
	}

	return event
}

// NewPlayerPlayedEvent 创建玩家出牌事件
func NewPlayerPlayedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	playerSeat int,
	cards []*Card,
	compType CompType,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_PLAYED,
		match, deal, trick, playerSeat,
	)

	event.Payload = &eventpb.GameEvent_PlayerPlayed{
		PlayerPlayed: &eventpb.PlayerPlayedPayload{
			Cards:    ConvertCardsToProto(cards),
			CompType: ConvertCompTypeToProto(compType),
		},
	}

	return event
}

// NewPlayerPassedEvent 创建玩家过牌事件
func NewPlayerPassedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	playerSeat int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_PASSED,
		match, deal, trick, playerSeat,
	)

	event.Payload = &eventpb.GameEvent_PlayerPassed{
		PlayerPassed: &eventpb.PlayerPassedPayload{},
	}

	return event
}

// NewPlayerFinishedEvent 创建玩家出完牌事件
func NewPlayerFinishedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	playerSeat int,
	finishRank int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_FINISHED,
		match, deal, trick, playerSeat,
	)

	event.Payload = &eventpb.GameEvent_PlayerFinished{
		PlayerFinished: &eventpb.PlayerFinishedPayload{
			FinishRank: int32(finishRank),
		},
	}

	return event
}

// NewDealEndedEvent 创建牌局结束事件
func NewDealEndedEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	dealLevel int,
	rankings []int,
	victoryType VictoryType,
	winningTeam int,
	levelChange [2]int,
	durationMs int64,
	trickCount int,
	playerStats [4]*PlayerDealStats,
	nextDealDeadlineMs int64,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_DEAL_ENDED,
		match, deal, nil, -1,
	)

	rankings32 := make([]int32, len(rankings))
	for i, r := range rankings {
		rankings32[i] = int32(r)
	}

	// 转换 PlayerDealStats
	protoPlayerStats := make([]*eventpb.PlayerDealStats, 0, 4)
	for _, ps := range playerStats {
		if ps != nil {
			protoPlayerStats = append(protoPlayerStats, &eventpb.PlayerDealStats{
				PlayerSeat:  int32(ps.PlayerSeat),
				CardsPlayed: int32(ps.CardsPlayed),
				TricksWon:   int32(ps.TricksWon),
				PassCount:   int32(ps.PassCount),
				FinishRank:  int32(ps.FinishRank),
			})
		}
	}

	event.Payload = &eventpb.GameEvent_DealEnded{
		DealEnded: &eventpb.DealEndedPayload{
			DealLevel:          int32(dealLevel),
			Rankings:           rankings32,
			VictoryType:        ConvertVictoryTypeToProto(victoryType),
			WinningTeam:        int32(winningTeam),
			LevelChange:        []int32{int32(levelChange[0]), int32(levelChange[1])},
			DurationMs:         durationMs,
			TrickCount:         int32(trickCount),
			PlayerStats:        protoPlayerStats,
			NextDealDeadlineMs: nextDealDeadlineMs,
		},
	}

	return event
}

// NewMatchEndedEvent 创建比赛结束事件
func NewMatchEndedEvent(
	emp *EventMetadataProvider,
	match *Match,
	winner int,
	finalLevels [2]int,
	durationMs int64,
	totalDeals int,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_MATCH_ENDED,
		match, nil, nil, -1,
	)

	event.Payload = &eventpb.GameEvent_MatchEnded{
		MatchEnded: &eventpb.MatchEndedPayload{
			Winner:      int32(winner),
			FinalLevels: []int32{int32(finalLevels[0]), int32(finalLevels[1])},
			DurationMs:  durationMs,
			TotalDeals:  int32(totalDeals),
		},
	}

	return event
}

// NewPlayerTimeoutEvent 创建玩家超时事件
func NewPlayerTimeoutEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	seat int,
	actionType string,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_TIMEOUT,
		match, deal, trick, seat,
	)

	event.Payload = &eventpb.GameEvent_PlayerTimeout{
		PlayerTimeout: &eventpb.PlayerTimeoutPayload{
			ActionType: ConvertPlayerTimeoutActionTypeToProto(actionType),
		},
	}

	return event
}

// NewPlayerDisconnectEvent 创建玩家断线事件
func NewPlayerDisconnectEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	playerSeat int,
	autoPlay bool,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_DISCONNECT,
		match, deal, trick, playerSeat,
	)

	event.Payload = &eventpb.GameEvent_PlayerDisconnect{
		PlayerDisconnect: &eventpb.PlayerDisconnectPayload{
			AutoPlay: autoPlay,
		},
	}

	return event
}

// NewPlayerReconnectEvent 创建玩家重连事件
func NewPlayerReconnectEvent(
	emp *EventMetadataProvider,
	match *Match,
	deal *Deal,
	trick *Trick,
	playerSeat int,
	autoPlay bool,
) *GameEvent {
	event := emp.CreateBaseEvent(
		eventpb.EventType_EVENT_TYPE_PLAYER_RECONNECT,
		match, deal, trick, playerSeat,
	)

	event.Payload = &eventpb.GameEvent_PlayerReconnect{
		PlayerReconnect: &eventpb.PlayerReconnectPayload{
			AutoPlay: autoPlay,
		},
	}

	return event
}
