package sdk

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// GameEventType 定义游戏事件的类型
// 游戏引擎通过这些事件类型来通知外部系统游戏状态的变化
type GameEventType string

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

// GameEvent 表示游戏中发生的事件及其相关数据
// 游戏引擎通过事件系统来通知外部关于游戏状态变化的信息
type GameEvent struct {
	Type       GameEventType `json:"type"`                  // 事件类型，标识这是什么类型的事件
	Data       interface{}   `json:"data"`                  // 事件数据，包含与事件相关的具体信息
	Timestamp  time.Time     `json:"timestamp"`             // 事件发生的时间戳
	PlayerSeat int           `json:"player_seat,omitempty"` // 触发事件的玩家座位号（如果适用）
}

// GameEventHandler 是处理游戏事件的函数类型
// 参数:
//
//	*GameEvent: 要处理的游戏事件
//
// 功能说明:
//   - 事件处理器在独立的协程中运行，不会阻塞游戏主流程
//   - 可以用于日志记录、统计分析、UI更新等
//   - 处理器应该快速执行，避免长时间阻塞
type GameEventHandler func(*GameEvent)

// GameState 表示游戏的完整状态
// 包含游戏的全局信息，适用于管理员或观察者视角
type GameState struct {
	ID           string     `json:"id"`                      // 游戏的唯一标识符
	Status       GameStatus `json:"status"`                  // 当前游戏状态（等待中、进行中、已结束）
	CurrentMatch *Match     `json:"current_match,omitempty"` // 当前活跃的比赛实例（如果有）
	CreatedAt    time.Time  `json:"created_at"`              // 游戏创建时间
	UpdatedAt    time.Time  `json:"updated_at"`              // 最后更新时间
}

// PlayerGameState 表示从特定玩家视角看到的游戏状态
// 包含该玩家的私有信息（如手牌）和公共可见信息
type PlayerGameState struct {
	PlayerSeat   int        `json:"player_seat"`   // 玩家的座位号(0-3)
	GameState    *GameState `json:"game_state"`    // 游戏的公共状态信息
	PlayerCards  []*Card    `json:"player_cards"`  // 该玩家的手牌（只对该玩家可见）
	VisibleCards []*Card    `json:"visible_cards"` // 当前可见的牌（已出的牌）
}

// GameStatus 表示游戏的当前状态
// 用于跟踪游戏的生命周期
type GameStatus string

// 游戏状态常量定义
const (
	GameStatusWaiting  GameStatus = "waiting"  // 等待开始状态，等待玩家加入
	GameStatusStarted  GameStatus = "started"  // 游戏进行中状态
	GameStatusFinished GameStatus = "finished" // 游戏已结束状态
)

// GameEngine 是管理完整游戏生命周期的主要游戏引擎
// 它协调所有游戏组件，处理玩家操作，管理游戏状态，并发送事件通知
type GameEngine struct {
	id            string                               // 游戏引擎的唯一标识符
	status        GameStatus                           // 当前游戏状态
	currentMatch  *Match                               // 当前活跃的比赛实例
	eventHandlers map[GameEventType][]GameEventHandler // 事件处理器映射，按事件类型分组
	mutex         sync.RWMutex                         // 读写锁，保护并发访问游戏状态
	createdAt     time.Time                            // 游戏引擎创建时间
	updatedAt     time.Time                            // 最后更新时间
}

// GameEngineInterface 定义了游戏引擎的公共接口
// 这个接口封装了掼蛋游戏的所有核心功能，包括游戏生命周期管理、
// 游戏操作、状态查询、事件处理和玩家管理等功能
type GameEngineInterface interface {
	// 游戏生命周期管理

	// StartMatch 开始一局新的比赛
	// 参数:
	//   players: 参与比赛的4个玩家列表，必须包含4个玩家
	// 返回值:
	//   error: 如果玩家数量不是4个或游戏状态不允许开始，返回错误
	// 功能说明:
	//   - 初始化新的比赛实例
	//   - 验证玩家数量和游戏状态
	//   - 设置游戏状态为已开始
	//   - 触发比赛开始事件
	StartMatch(players []Player) error

	// StartDeal 开始新的一局牌局
	// 返回值:
	//   error: 如果没有活跃的比赛或无法开始新局，返回错误
	// 功能说明:
	//   - 洗牌并发牌给4个玩家
	//   - 初始化进贡阶段（如果需要）
	//   - 设置首家和当前轮次
	//   - 触发牌局开始事件
	StartDeal() error

	// 游戏操作

	// PlayCards 玩家出牌
	// 参数:
	//   playerSeat: 玩家座位号(0-3)
	//   cards: 要出的牌的列表
	// 返回值:
	//   *GameEvent: 出牌成功时返回的游戏事件
	//   error: 如果出牌无效或不是该玩家回合，返回错误
	// 功能说明:
	//   - 验证玩家是否轮到出牌
	//   - 验证出牌组合的合法性
	//   - 更新游戏状态和玩家手牌
	//   - 检查是否需要进入下一轮或结束牌局
	//   - 触发玩家出牌事件
	PlayCards(playerSeat int, cards []*Card) (*GameEvent, error)

	// PassTurn 玩家选择不出牌（过牌）
	// 参数:
	//   playerSeat: 玩家座位号(0-3)
	// 返回值:
	//   *GameEvent: 过牌成功时返回的游戏事件
	//   error: 如果不是该玩家回合或不允许过牌，返回错误
	// 功能说明:
	//   - 验证玩家是否轮到出牌
	//   - 验证是否允许过牌（非首家出牌时可以过牌）
	//   - 更新当前轮次到下一个玩家
	//   - 触发玩家过牌事件
	PassTurn(playerSeat int) (*GameEvent, error)

	// 贡牌相关接口

	// ProcessTributePhase 处理贡牌阶段
	// 返回值:
	//   *TributeAction: 需要玩家执行的动作（选牌或还贡），如果为nil表示贡牌阶段已完成或免贡
	//   error: 如果处理失败，返回错误
	// 功能说明:
	//   - 委托给 TributeManager 处理所有贡牌逻辑
	//   - 自动完成免贡判定
	//   - 自动完成上贡选择（除红桃主牌外最大的牌）
	//   - 返回需要玩家输入的动作（双下选牌或还贡）
	//   - 当贡牌完成时，自动应用效果并更新游戏状态
	ProcessTributePhase() (*TributeAction, error)

	// SubmitTributeSelection 提交贡牌选择（用于双下选牌）
	// 参数:
	//   playerID: 选择的玩家座位号(0-3)
	//   cardID: 选择的牌ID
	// 返回值:
	//   error: 如果选择无效或不是该玩家的选择回合，返回错误
	// 功能说明:
	//   - 用于双下情况下rank1玩家从贡牌池选择一张牌
	//   - 验证玩家是否有权选择
	//   - 验证选择的牌是否在贡牌池中
	//   - 自动将剩余的牌分配给rank2
	SubmitTributeSelection(playerID int, cardID string) error

	// SubmitReturnTribute 提交还贡
	// 参数:
	//   playerID: 还贡的玩家座位号(0-3)
	//   cardID: 还贡的牌ID
	// 返回值:
	//   error: 如果还贡无效或不是该玩家的还贡回合，返回错误
	// 功能说明:
	//   - 用于收到贡牌的玩家选择一张牌还给贡牌方
	//   - 验证玩家是否需要还贡
	//   - 验证选择的牌是否在该玩家手中
	//   - 完成牌的交换
	SubmitReturnTribute(playerID int, cardID string) error

	// SkipTributeAction 跳过当前贡牌动作（超时处理）
	// 返回值:
	//   error: 如果当前没有待处理的贡牌动作，返回错误
	// 功能说明:
	//   - 用于处理玩家超时的情况
	//   - 双下选牌时自动选择最大的牌
	//   - 还贡时自动选择最小的牌
	SkipTributeAction() error

	// GetTributeStatus 获取当前贡牌状态
	// 返回值:
	//   *TributeStatusInfo: 当前贡牌阶段的详细信息，如果不在贡牌阶段返回nil
	// 功能说明:
	//   - 查询当前贡牌阶段的状态
	//   - 包含已确定的贡牌、还贡信息
	//   - 包含待执行的动作列表
	GetTributeStatus() *TributeStatusInfo

	// 状态查询

	// GetGameState 获取当前完整的游戏状态
	// 返回值:
	//   *GameState: 包含游戏ID、状态、当前比赛等信息的完整状态
	// 功能说明:
	//   - 返回游戏的全局状态信息
	//   - 包括比赛进度、当前牌局状态等
	//   - 适用于管理员或观察者视角
	GetGameState() *GameState

	// GetPlayerView 获取特定玩家视角的游戏状态
	// 参数:
	//   playerSeat: 玩家座位号(0-3)
	// 返回值:
	//   *PlayerGameState: 包含该玩家手牌和可见信息的游戏状态
	// 功能说明:
	//   - 返回从指定玩家角度看到的游戏状态
	//   - 包含该玩家的手牌信息
	//   - 包含当前可见的牌（已出的牌）
	//   - 隐藏其他玩家的手牌信息
	GetPlayerView(playerSeat int) *PlayerGameState

	// IsGameFinished 检查游戏是否已结束
	// 返回值:
	//   bool: 如果游戏已结束返回true，否则返回false
	// 功能说明:
	//   - 快速检查游戏是否处于结束状态
	//   - 用于判断是否还可以进行游戏操作
	IsGameFinished() bool

	// 事件处理

	// RegisterEventHandler 注册游戏事件处理器
	// 参数:
	//   eventType: 要监听的事件类型
	//   handler: 事件处理函数
	// 功能说明:
	//   - 允许外部系统监听游戏事件
	//   - 支持多个处理器监听同一事件类型
	//   - 事件处理器在独立的协程中执行，不会阻塞游戏进程
	//   - 常用事件包括：出牌、过牌、牌局结束、比赛结束等
	RegisterEventHandler(eventType GameEventType, handler GameEventHandler)

	// ProcessTimeouts 处理超时情况
	// 返回值:
	//   []*GameEvent: 因超时产生的事件列表
	// 功能说明:
	//   - 检查是否有玩家操作超时
	//   - 对超时玩家执行自动操作
	//   - 返回因超时处理产生的事件
	//   - 需要定期调用以维护游戏进度
	ProcessTimeouts() []*GameEvent

	// 玩家管理

	// HandlePlayerDisconnect 处理玩家断线
	// 参数:
	//   playerSeat: 断线玩家的座位号(0-3)
	// 返回值:
	//   *GameEvent: 断线处理事件
	//   error: 如果处理断线失败，返回错误
	// 功能说明:
	//   - 标记玩家为断线状态
	//   - 启用该玩家的自动托管模式
	//   - 触发玩家断线事件
	//   - 游戏继续进行，由系统代为操作
	HandlePlayerDisconnect(playerSeat int) (*GameEvent, error)

	// HandlePlayerReconnect 处理玩家重连
	// 参数:
	//   playerSeat: 重连玩家的座位号(0-3)
	// 返回值:
	//   *GameEvent: 重连处理事件
	//   error: 如果处理重连失败，返回错误
	// 功能说明:
	//   - 恢复玩家的在线状态
	//   - 关闭自动托管模式
	//   - 触发玩家重连事件
	//   - 玩家可以重新手动操作
	HandlePlayerReconnect(playerSeat int) (*GameEvent, error)

	// SetPlayerAutoPlay 设置玩家的自动托管状态
	// 参数:
	//   playerSeat: 玩家座位号(0-3)
	//   enabled: 是否启用自动托管
	// 返回值:
	//   error: 如果设置失败，返回错误
	// 功能说明:
	//   - 手动控制玩家的托管状态
	//   - 启用托管时，系统将自动为该玩家做决策
	//   - 可用于处理长时间未操作的玩家
	SetPlayerAutoPlay(playerSeat int, enabled bool) error

	// 新增状态查询接口

	// GetCurrentDealStatus 获取当前牌局状态
	// 返回值:
	//   DealStatus: 当前牌局的状态（waiting/dealing/tribute/playing/finished）
	// 功能说明:
	//   - 提供当前牌局状态的快速查询
	//   - 替代直接访问deal.Status的需求
	//   - 如果没有活跃牌局，返回DealStatusWaiting
	GetCurrentDealStatus() DealStatus

	// GetCurrentTurnInfo 获取当前轮次信息
	// 返回值:
	//   *TurnInfo: 当前轮次的详细信息，如果没有活跃轮次返回nil
	// 功能说明:
	//   - 提供当前轮次的完整信息
	//   - 包括当前玩家、是否为首出、是否为新trick等
	//   - 替代直接访问deal.CurrentTrick的需求
	GetCurrentTurnInfo() *TurnInfo

	// GetMatchDetails 获取比赛详细信息
	// 返回值:
	//   *MatchDetails: 比赛的详细信息，如果没有活跃比赛返回nil
	// 功能说明:
	//   - 提供比赛级别的信息
	//   - 包括队伍等级、玩家信息等
	//   - 替代直接访问match对象的需求
	GetMatchDetails() *MatchDetails
}

// NewGameEngine creates a new game engine instance
func NewGameEngine() *GameEngine {
	now := time.Now()
	return &GameEngine{
		id:            generateID(),
		status:        GameStatusWaiting,
		eventHandlers: make(map[GameEventType][]GameEventHandler),
		createdAt:     now,
		updatedAt:     now,
	}
}

// StartMatch initializes a new match with the given players
func (ge *GameEngine) StartMatch(players []Player) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if len(players) != 4 {
		return errors.New("exactly 4 players are required")
	}

	if ge.status != GameStatusWaiting {
		return errors.New("game is not in waiting status")
	}

	// Create new match
	match, err := NewMatch(players)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	ge.currentMatch = match
	ge.status = GameStatusStarted
	ge.updatedAt = time.Now()

	// Emit match started event
	event := &GameEvent{
		Type:      EventMatchStarted,
		Data:      match,
		Timestamp: time.Now(),
	}
	ge.emitEvent(event)

	return nil
}

// StartDeal starts a new deal in the current match
func (ge *GameEngine) StartDeal() error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil {
		return errors.New("no active match")
	}

	err := ge.currentMatch.StartNewDeal()
	if err != nil {
		return fmt.Errorf("failed to start deal: %w", err)
	}

	ge.updatedAt = time.Now()

	// Emit deal started event
	event := &GameEvent{
		Type:      EventDealStarted,
		Data:      ge.currentMatch.CurrentDeal,
		Timestamp: time.Now(),
	}
	ge.emitEvent(event)

	// If there's a tribute phase, emit tribute rules set event first
	if ge.currentMatch.CurrentDeal.TributePhase != nil {
		tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
		tributeMap, isDoubleDown, _ := tm.DetermineTributeRequirements(ge.currentMatch.CurrentDeal.LastResult)

		// Create rule description based on victory type
		var ruleDescription string
		lastResult := ge.currentMatch.CurrentDeal.LastResult
		switch lastResult.VictoryType {
		case VictoryTypeDoubleDown:
			ruleDescription = fmt.Sprintf("双下：Player%d和Player%d上贡到池，Player%d优先选择",
				lastResult.Rankings[2], lastResult.Rankings[3], lastResult.Rankings[0])
		case VictoryTypeSingleLast:
			ruleDescription = fmt.Sprintf("单下：Player%d上贡给Player%d",
				lastResult.Rankings[3], lastResult.Rankings[0])
		case VictoryTypePartnerLast:
			ruleDescription = fmt.Sprintf("对下：Player%d上贡给Player%d",
				lastResult.Rankings[2], lastResult.Rankings[0])
		}

		// Emit tribute rules set event
		rulesEvent := &GameEvent{
			Type: EventTributeRulesSet,
			Data: map[string]interface{}{
				"last_result":  lastResult,
				"victory_type": lastResult.VictoryType,
				"tribute_rules": map[string]interface{}{
					"tribute_map":    tributeMap,
					"is_double_down": isDoubleDown,
					"description":    ruleDescription,
				},
				"player_rankings": lastResult.Rankings,
			},
			Timestamp: time.Now(),
		}
		ge.emitEvent(rulesEvent)
	}

	// Check if tribute phase was skipped due to immunity
	if ge.currentMatch.CurrentDeal.TributePhase != nil &&
		ge.currentMatch.CurrentDeal.TributePhase.IsImmune {
		// Get detailed immunity information
		tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
		_, immunityDetails := tm.GetTributeImmunityDetails(ge.currentMatch.CurrentDeal.LastResult,
			ge.currentMatch.CurrentDeal.PlayerCards)

		// Emit immunity event with detailed information
		immunityEvent := &GameEvent{
			Type: EventTributeImmunity,
			Data: map[string]interface{}{
				"tribute_phase":   ge.currentMatch.CurrentDeal.TributePhase,
				"immunity_reason": immunityDetails,
			},
			Timestamp: time.Now(),
		}
		ge.emitEvent(immunityEvent)
	}

	return nil
}

// PlayCards handles a player's card play action
func (ge *GameEngine) PlayCards(playerSeat int, cards []*Card) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal

	// Use validator to validate the play
	validator := NewPlayValidator(deal.Level)
	err := validator.ValidatePlay(playerSeat, cards, deal.PlayerCards[playerSeat], deal.CurrentTrick)
	if err != nil {
		return nil, fmt.Errorf("invalid play: %w", err)
	}

	// Check for pre-action state transitions (e.g., trick starting) BEFORE executing the play
	preEvents := ge.checkPreActionStateTransitions()
	for _, evt := range preEvents {
		ge.emitEvent(evt)
	}

	// Execute the play
	err = deal.PlayCards(playerSeat, cards)
	if err != nil {
		return nil, fmt.Errorf("failed to play cards: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create and emit player played event
	event := &GameEvent{
		Type: EventPlayerPlayed,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"cards":       cards,
			"deal_state":  deal,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)

	// Check for post-action state transitions (e.g., trick ending, deal ending)
	postEvents := ge.checkPostActionStateTransitions()
	for _, evt := range postEvents {
		ge.emitEvent(evt)
	}

	return event, nil
}

// PassTurn handles a player's pass action
func (ge *GameEngine) PassTurn(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal

	// Use validator to validate the pass
	validator := NewPlayValidator(deal.Level)
	err := validator.ValidatePass(playerSeat, deal.CurrentTrick)
	if err != nil {
		return nil, fmt.Errorf("invalid pass: %w", err)
	}

	// Check for pre-action state transitions (e.g., trick starting) BEFORE executing the pass
	preEvents := ge.checkPreActionStateTransitions()
	for _, evt := range preEvents {
		ge.emitEvent(evt)
	}

	// Execute the pass
	err = deal.PassTurn(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to pass turn: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create and emit player passed event
	event := &GameEvent{
		Type: EventPlayerPassed,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"deal_state":  deal,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)

	// Check for post-action state transitions (e.g., trick ending, deal ending)
	postEvents := ge.checkPostActionStateTransitions()
	for _, evt := range postEvents {
		ge.emitEvent(evt)
	}

	return event, nil
}

// GetGameState returns the current complete game state
func (ge *GameEngine) GetGameState() *GameState {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	return &GameState{
		ID:           ge.id,
		Status:       ge.status,
		CurrentMatch: ge.currentMatch,
		CreatedAt:    ge.createdAt,
		UpdatedAt:    ge.updatedAt,
	}
}

// GetPlayerView returns the game state from a specific player's perspective
func (ge *GameEngine) GetPlayerView(playerSeat int) *PlayerGameState {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	gameState := ge.GetGameState()
	playerView := &PlayerGameState{
		PlayerSeat: playerSeat,
		GameState:  gameState,
	}

	// Add player-specific information if there's an active deal
	if ge.currentMatch != nil && ge.currentMatch.CurrentDeal != nil {
		if playerSeat >= 0 && playerSeat < 4 {
			playerView.PlayerCards = ge.currentMatch.CurrentDeal.PlayerCards[playerSeat]
		}

		// Add visible cards (cards played in current trick)
		if ge.currentMatch.CurrentDeal.CurrentTrick != nil {
			playerView.VisibleCards = ge.getVisibleCardsForPlayer(playerSeat)
		}
	}

	return playerView
}

// IsGameFinished checks if the game is finished
func (ge *GameEngine) IsGameFinished() bool {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	return ge.status == GameStatusFinished
}

// RegisterEventHandler registers an event handler for a specific event type
func (ge *GameEngine) RegisterEventHandler(eventType GameEventType, handler GameEventHandler) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.eventHandlers[eventType] == nil {
		ge.eventHandlers[eventType] = make([]GameEventHandler, 0)
	}
	ge.eventHandlers[eventType] = append(ge.eventHandlers[eventType], handler)
}

// ProcessTimeouts processes any pending timeouts and returns resulting events
func (ge *GameEngine) ProcessTimeouts() []*GameEvent {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	events := make([]*GameEvent, 0)

	if ge.currentMatch != nil && ge.currentMatch.CurrentDeal != nil {
		timeoutEvents := ge.currentMatch.CurrentDeal.ProcessTimeouts()
		events = append(events, timeoutEvents...)

		// Emit all timeout events
		for _, event := range timeoutEvents {
			ge.emitEvent(event)
		}
	}

	return events
}

// HandlePlayerDisconnect handles a player disconnection
func (ge *GameEngine) HandlePlayerDisconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil {
		return nil, errors.New("no active match")
	}

	err := ge.currentMatch.HandlePlayerDisconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle disconnect: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create disconnect event
	event := &GameEvent{
		Type: EventPlayerDisconnect,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   true,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)

	return event, nil
}

// HandlePlayerReconnect handles a player reconnection
func (ge *GameEngine) HandlePlayerReconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil {
		return nil, errors.New("no active match")
	}

	err := ge.currentMatch.HandlePlayerReconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle reconnect: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create reconnect event
	event := &GameEvent{
		Type: EventPlayerReconnect,
		Data: map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   false,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)

	return event, nil
}

// SetPlayerAutoPlay sets the auto-play status for a player
func (ge *GameEngine) SetPlayerAutoPlay(playerSeat int, enabled bool) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.currentMatch == nil {
		return errors.New("no active match")
	}

	return ge.currentMatch.SetPlayerAutoPlay(playerSeat, enabled)
}

// emitEvent emits an event to all registered handlers
func (ge *GameEngine) emitEvent(event *GameEvent) {
	handlers, exists := ge.eventHandlers[event.Type]
	if !exists {
		return
	}

	// Call all handlers for this event type synchronously to maintain order
	for _, handler := range handlers {
		handler(event) // 同步调用确保事件按顺序处理
	}
}

// getVisibleCardsForPlayer returns the cards visible to a specific player
func (ge *GameEngine) getVisibleCardsForPlayer(playerSeat int) []*Card {
	visibleCards := make([]*Card, 0)

	if ge.currentMatch.CurrentDeal.CurrentTrick != nil {
		for _, play := range ge.currentMatch.CurrentDeal.CurrentTrick.Plays {
			if play.Cards != nil {
				visibleCards = append(visibleCards, play.Cards...)
			}
		}
	}

	return visibleCards
}

// checkPreActionStateTransitions checks for state transitions that should happen before player actions
// Currently only handles trick starting (TrickStatusWaiting -> TrickStatusPlaying)
func (ge *GameEngine) checkPreActionStateTransitions() []*GameEvent {
	events := make([]*GameEvent, 0)

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return events
	}

	deal := ge.currentMatch.CurrentDeal

	// Check if there's a waiting trick that needs to be started
	if deal.CurrentTrick != nil && deal.CurrentTrick.Status == TrickStatusWaiting {
		// Start the new trick
		err := deal.CurrentTrick.StartTrick()
		if err == nil {
			// 收集所有玩家的手牌信息，避免在事件处理器中调用引擎方法
			playerHands := make(map[int][]*Card)
			for i := 0; i < 4; i++ {
				if deal.PlayerCards[i] != nil {
					// 创建手牌的副本，避免并发访问问题
					handCopy := make([]*Card, len(deal.PlayerCards[i]))
					copy(handCopy, deal.PlayerCards[i])
					playerHands[i] = handCopy
				}
			}

			trickStartedEvent := &GameEvent{
				Type: EventTrickStarted,
				Data: map[string]interface{}{
					"trick":        deal.CurrentTrick,
					"leader":       deal.CurrentTrick.Leader,
					"current_turn": deal.CurrentTrick.CurrentTurn,
					"player_hands": playerHands,
				},
				Timestamp: time.Now(),
			}
			events = append(events, trickStartedEvent)
		}
	}

	return events
}

// checkPostActionStateTransitions checks for state transitions that should happen after player actions
// Handles trick ending and deal ending
func (ge *GameEngine) checkPostActionStateTransitions() []*GameEvent {
	events := make([]*GameEvent, 0)

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return events
	}

	deal := ge.currentMatch.CurrentDeal

	// Check if deal is finished first (can happen at any time)
	if deal.Status == DealStatusFinished {
		// Calculate deal result using the new result system
		dealResult, err := deal.CalculateResult(ge.currentMatch)
		if err != nil {
			// Log error but continue - create a basic result
			winningTeam := ge.currentMatch.GetTeamForPlayer(deal.Rankings[0])
			upgrades := [2]int{0, 0}
			upgrades[winningTeam] = 1 // 确保获胜队伍能够升级

			dealResult = &DealResult{
				Rankings:    deal.Rankings,
				WinningTeam: winningTeam,
				VictoryType: VictoryTypePartnerLast,
				Upgrades:    upgrades,
				Duration:    time.Since(deal.StartTime),
				TrickCount:  len(deal.TrickHistory),
			}
		}

		// Emit deal ended event
		dealEndedEvent := &GameEvent{
			Type: EventDealEnded,
			Data: map[string]interface{}{
				"deal":       deal,
				"result":     dealResult,
				"rankings":   deal.Rankings,
				"statistics": dealResult.Statistics,
			},
			Timestamp: time.Now(),
		}
		events = append(events, dealEndedEvent)

		// Update match with deal result
		err = ge.currentMatch.FinishDeal(dealResult)
		if err == nil {
			// Check if match is finished
			if ge.currentMatch.Status == MatchStatusFinished {
				ge.status = GameStatusFinished

				// Create match result
				matchResult := ge.createMatchResult()

				// Emit match ended event
				matchEndedEvent := &GameEvent{
					Type: EventMatchEnded,
					Data: map[string]interface{}{
						"match":        ge.currentMatch,
						"result":       matchResult,
						"winner":       ge.currentMatch.Winner,
						"final_levels": ge.currentMatch.TeamLevels,
					},
					Timestamp: time.Now(),
				}
				events = append(events, matchEndedEvent)
			}
		}
	} else if deal.CurrentTrick != nil && deal.CurrentTrick.Status == TrickStatusFinished {
		// Check if current trick is finished
		// Emit trick ended event
		finishedTrick := deal.CurrentTrick
		trickEndedEvent := &GameEvent{
			Type: EventTrickEnded,
			Data: map[string]interface{}{
				"trick":       finishedTrick,
				"winner":      finishedTrick.Winner,
				"next_leader": finishedTrick.NextLeader,
			},
			Timestamp: time.Now(),
		}
		events = append(events, trickEndedEvent)

		// Add finished trick to history
		deal.TrickHistory = append(deal.TrickHistory, finishedTrick)

		// Create new trick with the next leader
		nextTrick, err := NewTrick(finishedTrick.NextLeader)
		if err == nil {
			// Set the new trick but leave it in TrickStatusWaiting
			deal.CurrentTrick = nextTrick
		}
	}

	return events
}

// checkStateTransitions checks for and handles automatic state transitions (legacy method)
// Now delegates to pre-action and post-action methods for backward compatibility
func (ge *GameEngine) checkStateTransitions() []*GameEvent {
	events := make([]*GameEvent, 0)

	// Check pre-action transitions first
	preEvents := ge.checkPreActionStateTransitions()
	events = append(events, preEvents...)

	// Then check post-action transitions
	postEvents := ge.checkPostActionStateTransitions()
	events = append(events, postEvents...)

	return events
}

// createMatchResult creates a MatchResult from a finished match
func (ge *GameEngine) createMatchResult() *MatchResult {
	if ge.currentMatch == nil || ge.currentMatch.Status != MatchStatusFinished {
		return nil
	}

	// Calculate total duration
	duration := time.Duration(0)
	if ge.currentMatch.EndTime != nil {
		duration = ge.currentMatch.EndTime.Sub(ge.currentMatch.StartTime)
	}

	// Calculate match statistics
	stats := &MatchStatistics{
		TotalDeals:    len(ge.currentMatch.DealHistory),
		TotalDuration: duration,
		FinalLevels:   ge.currentMatch.TeamLevels,
		TeamStats:     [2]*TeamMatchStats{},
	}

	// Initialize team stats
	for team := 0; team < 2; team++ {
		stats.TeamStats[team] = &TeamMatchStats{
			Team:        team,
			DealsWon:    0,
			TotalTricks: 0,
			Upgrades:    0,
		}
	}

	// Calculate team statistics from deal history
	for _, deal := range ge.currentMatch.DealHistory {
		if result, err := deal.CalculateResult(ge.currentMatch); err == nil {
			// Count deals won
			stats.TeamStats[result.WinningTeam].DealsWon++

			// Count upgrades
			for team := 0; team < 2; team++ {
				stats.TeamStats[team].Upgrades += result.Upgrades[team]
			}

			// Count tricks won by each team
			if result.Statistics != nil {
				for _, playerStats := range result.Statistics.PlayerStats {
					if playerStats != nil {
						team := ge.currentMatch.GetTeamForPlayer(playerStats.PlayerSeat)
						stats.TeamStats[team].TotalTricks += playerStats.TricksWon
					}
				}
			}
		}
	}

	return &MatchResult{
		Winner:      ge.currentMatch.Winner,
		FinalLevels: ge.currentMatch.TeamLevels,
		Duration:    duration,
		Statistics:  stats,
	}
}

// AutoPlayForPlayer executes an automatic play for a player (used for disconnected/timeout players)
func (ge *GameEngine) AutoPlayForPlayer(playerSeat int) (*GameEvent, error) {
	// Don't lock here since we'll call PlayCards or PassTurn which will lock
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal

	// Check if it's the player's turn
	if deal.CurrentTrick == nil || deal.CurrentTrick.CurrentTurn != playerSeat {
		return nil, errors.New("not player's turn")
	}

	// Auto-play strategy: if trick leader, play smallest card; otherwise pass
	if deal.CurrentTrick.LeadComp == nil {
		// Player is trick leader - play smallest single card
		playerCards := deal.PlayerCards[playerSeat]
		if len(playerCards) > 0 {
			// Find smallest card
			smallestCard := playerCards[0]
			for _, card := range playerCards {
				if card.LessThan(smallestCard) {
					smallestCard = card
				}
			}

			return ge.PlayCards(playerSeat, []*Card{smallestCard})
		}
	} else {
		// Player is not leader - pass
		return ge.PassTurn(playerSeat)
	}

	return nil, errors.New("unable to auto-play")
}

// sendEvent 发送游戏事件到所有注册的处理器
func (ge *GameEngine) sendEvent(event *GameEvent) {
	// 获取该事件类型的所有处理器
	handlers, exists := ge.eventHandlers[event.Type]
	if !exists || len(handlers) == 0 {
		return
	}

	// 在独立的协程中执行每个处理器
	for _, handler := range handlers {
		go handler(event)
	}
}

// ProcessTributePhase 处理贡牌阶段
func (ge *GameEngine) ProcessTributePhase() (*TributeAction, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// 验证基本状态
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute {
		return nil, errors.New("not in tribute phase")
	}

	if deal.TributePhase == nil {
		return nil, nil
	}

	// 记录处理前的状态和贡牌情况
	previousStatus := deal.TributePhase.Status
	previousTributeCards := make(map[int]*Card)
	for giver, card := range deal.TributePhase.TributeCards {
		previousTributeCards[giver] = card
	}

	// 获取 TributeManager 并处理
	tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
	action, err := tm.ProcessTributePhaseAction(deal.TributePhase, deal.PlayerCards)
	if err != nil {
		return nil, err
	}

	// 检测状态变化并触发相应事件
	if previousStatus == TributeStatusWaiting && deal.TributePhase.Status == TributeStatusSelecting {
		// 双下场景：贡牌池已创建
		var contributors []map[string]interface{}
		selectionOrder := []int{deal.TributePhase.SelectingPlayer}

		// 根据贡牌映射找出贡献者
		for giver := range deal.TributePhase.TributeMap {
			if deal.TributePhase.TributeMap[giver] == -1 {
				// 贡献到池子的玩家
				if tributeCard := deal.TributePhase.TributeCards[giver]; tributeCard != nil {
					contributors = append(contributors, map[string]interface{}{
						"player_seat": giver,
						"card":        tributeCard,
					})
				}
			}
		}

		// 确定选择顺序（第二名是第一名的队友）
		if len(selectionOrder) > 0 {
			secondPlace := (selectionOrder[0] + 2) % 4 // 队友
			selectionOrder = append(selectionOrder, secondPlace)
		}

		// 触发贡牌池创建事件
		poolEvent := &GameEvent{
			Type: EventTributePoolCreated,
			Data: map[string]interface{}{
				"description":      fmt.Sprintf("双下贡牌池已创建，包含%d张贡牌", len(contributors)),
				"contributors":     contributors,
				"selection_order":  selectionOrder,
				"pool_cards":       deal.TributePhase.PoolCards,
				"selecting_player": deal.TributePhase.SelectingPlayer,
			},
			Timestamp: time.Now(),
		}
		ge.emitEvent(poolEvent)
	}

	// 检测上贡卡牌是否刚刚被确定（适用于所有场景）
	for giver, receiver := range deal.TributePhase.TributeMap {
		if receiver != -1 { // 不是贡献到池子的情况
			// 检查是否有新的贡牌被确定（从nil变为非nil）
			currentCard := deal.TributePhase.TributeCards[giver]
			previousCard := previousTributeCards[giver]

			if currentCard != nil && previousCard == nil {
				// 触发上贡完成事件
				givenEvent := &GameEvent{
					Type: EventTributeGiven,
					Data: map[string]interface{}{
						// 保留现有字段格式
						"giver":    giver,
						"receiver": receiver,
						"card":     currentCard,

						// 新增字段
						"tribute_type":     "normal",
						"is_auto_selected": true,
						"selection_reason": "除红桃Trump外最大牌",
					},
					Timestamp: time.Now(),
				}
				ge.emitEvent(givenEvent)
			}
		}
	}

	// 如果贡牌阶段完成，更新状态并发送事件
	if deal.TributePhase.Status == TributeStatusFinished {
		// 应用贡牌效果
		err = tm.ApplyTributeToHands(deal.TributePhase, &deal.PlayerCards)
		if err != nil {
			return nil, fmt.Errorf("apply tribute failed: %w", err)
		}

		// 发送完成事件（同步发送以确保日志顺序正确）
		ge.emitEvent(&GameEvent{
			Type:      EventTributeCompleted,
			Data:      deal.TributePhase,
			Timestamp: time.Now(),
		})

		// 启动游戏阶段（包括创建第一个trick和设置状态）
		err = deal.StartPlayingPhase()
		if err != nil {
			return nil, fmt.Errorf("failed to start playing phase: %w", err)
		}
	}

	return action, nil
}

// SubmitTributeSelection 提交贡牌选择（用于双下选牌）
func (ge *GameEngine) SubmitTributeSelection(playerID int, cardID string) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// 验证基本状态
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute || deal.TributePhase == nil {
		return errors.New("not in tribute phase")
	}

	// 调用 TributeManager 处理选择
	tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
	err := tm.SubmitSelection(deal.TributePhase, playerID, cardID)
	if err != nil {
		return err
	}

	// 收集增强的选择事件数据
	var selectedCard *Card
	for _, card := range deal.TributePhase.PoolCards {
		if card.GetID() == cardID {
			selectedCard = card
			break
		}
	}

	// 获取当前池中剩余的卡牌
	remainingOptions := make([]*Card, len(deal.TributePhase.PoolCards))
	copy(remainingOptions, deal.TributePhase.PoolCards)

	// 确定选择顺序
	selectionOrder := 1 // 默认为第一次选择
	if deal.TributePhase.SelectingPlayer != deal.TributePhase.TributeMap[deal.TributePhase.SelectingPlayer] {
		// 如果选择者不是第一名，说明是第二次选择
		selectionOrder = 2
	}

	// 发送增强的选择事件
	ge.emitEvent(&GameEvent{
		Type: EventTributeSelected,
		Data: map[string]interface{}{
			// 保留现有字段
			"action": "select",
			"player": playerID,
			"cardID": cardID,

			// 新增字段
			"selected_card":     selectedCard,
			"remaining_options": remainingOptions,
			"selection_order":   selectionOrder,
			"is_timeout":        false, // 正常选择，非超时
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerID,
	})

	return nil
}

// SubmitReturnTribute 提交还贡
func (ge *GameEngine) SubmitReturnTribute(playerID int, cardID string) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// 验证基本状态
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute || deal.TributePhase == nil {
		return errors.New("not in tribute phase")
	}

	// 调用 TributeManager 处理还贡
	tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
	err := tm.SubmitReturn(deal.TributePhase, playerID, cardID, deal.PlayerCards[playerID])
	if err != nil {
		return err
	}

	// 收集增强的还贡事件数据
	var returnCard *Card
	for _, card := range deal.PlayerCards[playerID] {
		if card.GetID() == cardID {
			returnCard = card
			break
		}
	}

	// 找到还贡的目标玩家和原来收到的贡牌
	var targetPlayer int = -1
	var originalTribute *Card
	for giver, receiver := range deal.TributePhase.TributeMap {
		if receiver == playerID && receiver != -1 {
			targetPlayer = giver
			originalTribute = deal.TributePhase.TributeCards[giver]
			break
		}
	}

	// 发送增强的还贡事件
	ge.emitEvent(&GameEvent{
		Type: EventReturnTribute,
		Data: map[string]interface{}{
			// 保留现有字段
			"action": "return",
			"player": playerID,
			"cardID": cardID,

			// 新增字段
			"return_card":      returnCard,
			"target_player":    targetPlayer,
			"original_tribute": originalTribute,
			"is_auto_selected": false, // 正常选择，非自动
			"selection_reason": "玩家手动选择",
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerID,
	})

	return nil
}

// SkipTributeAction 跳过当前贡牌动作（超时处理）
func (ge *GameEngine) SkipTributeAction() error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// 验证基本状态
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return errors.New("no active deal")
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute || deal.TributePhase == nil {
		return errors.New("not in tribute phase")
	}

	// 调用 TributeManager 处理超时
	tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
	err := tm.HandleTimeoutAction(deal.TributePhase, deal.PlayerCards)
	if err != nil {
		return err
	}

	// 发送超时事件
	ge.emitEvent(&GameEvent{
		Type:      EventPlayerTimeout,
		Data:      map[string]interface{}{"action": "tribute_timeout", "phase": deal.TributePhase.Status},
		Timestamp: time.Now(),
	})

	return nil
}

// GetTributeStatus 获取当前贡牌状态
func (ge *GameEngine) GetTributeStatus() *TributeStatusInfo {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	// 验证基本状态
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute || deal.TributePhase == nil {
		return nil
	}

	// 调用 TributeManager 获取状态信息
	tm := NewTributeManager(ge.currentMatch.TeamLevels[0])
	return tm.GetTributeStatusInfo(deal.TributePhase, deal.PlayerCards)
}

// GetCurrentDealStatus 获取当前牌局状态
func (ge *GameEngine) GetCurrentDealStatus() DealStatus {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return DealStatusWaiting
	}

	return ge.currentMatch.CurrentDeal.Status
}

// GetCurrentTurnInfo 获取当前轮次信息
func (ge *GameEngine) GetCurrentTurnInfo() *TurnInfo {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.CurrentTrick == nil {
		return &TurnInfo{
			CurrentPlayer:  -1,
			IsLeader:       false,
			IsNewTrick:     false,
			HasActiveTrick: false,
			LeadComp:       nil,
		}
	}

	trick := deal.CurrentTrick
	return &TurnInfo{
		CurrentPlayer:  trick.CurrentTurn,
		IsLeader:       trick.LeadComp == nil,
		IsNewTrick:     len(trick.Plays) == 0,
		HasActiveTrick: true,
		LeadComp:       trick.LeadComp,
	}
}

// GetMatchDetails 获取比赛详细信息
func (ge *GameEngine) GetMatchDetails() *MatchDetails {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	if ge.currentMatch == nil {
		return nil
	}

	match := ge.currentMatch
	players := make([]*PlayerInfo, 4)

	for i := 0; i < 4; i++ {
		if match.Players[i] != nil {
			players[i] = &PlayerInfo{
				Seat:     match.Players[i].Seat,
				Username: match.Players[i].Username,
				TeamNum:  i % 2, // 座位号0,2为team0；座位号1,3为team1
			}
		}
	}

	return &MatchDetails{
		TeamLevels: match.TeamLevels,
		Players:    players,
	}
}

// generateID generates a unique ID for the game engine
func generateID() string {
	return fmt.Sprintf("game_%d", time.Now().UnixNano())
}
