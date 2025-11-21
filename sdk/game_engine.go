package sdk

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	viewpb "guandan-world/proto/view"
)

// GameState 表示游戏的完整状态
// 包含游戏的全局信息，适用于管理员或观察者视角
type GameState struct {
	ID           string     `json:"id"`                      // 游戏的唯一标识符
	Status       GameStatus `json:"status"`                  // 当前游戏状态（等待中、进行中、已结束）
	CurrentMatch *Match     `json:"current_match,omitempty"` // 当前活跃的比赛实例（如果有）
	CreatedAt    time.Time  `json:"created_at"`              // 游戏创建时间
	UpdatedAt    time.Time  `json:"updated_at"`              // 最后更新时间
}

// PlayerView 表示从特定玩家视角看到的游戏状态
// 包含该玩家的私有信息（如手牌）和公共可见信息
// 该结构采用扁平化设计，移除冗余字段，减少嵌套层级，提升传输和解析效率
type PlayerView struct {
	PlayerSeat   int           `json:"player_seat"`             // 玩家的座位号(0-3)
	PlayerCards  []*Card       `json:"player_cards"`            // 该玩家的手牌（只对该玩家可见）
	TeamLevels   [2]int        `json:"team_levels"`             // 两队当前等级 [team0, team1]
	DealLevel    int           `json:"deal_level"`              // 当前局的等级
	DealStatus   DealStatus    `json:"deal_status"`             // 当前局的状态
	TrickID      string        `json:"trick_id,omitempty"`      // 当前Trick的ID（playing阶段）
	CurrentTurn  *int          `json:"current_turn,omitempty"`  // 当前轮到的玩家座位号（playing阶段，使用指针避免0值被omitempty省略）
	Leader       *int          `json:"leader,omitempty"`        // 当前Trick的首家玩家座位号（playing阶段，使用指针避免0值被omitempty省略）
	Plays        []*PlayAction `json:"plays,omitempty"`         // 当前Trick的所有出牌记录（playing阶段）
	TributePhase *TributePhase `json:"tribute_phase,omitempty"` // 上贡阶段信息（tribute阶段）
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
	id              string                            // 游戏引擎的唯一标识符
	status          GameStatus                        // 当前游戏状态
	currentMatch    *Match                            // 当前活跃的比赛实例
	observers       map[GameEventType][]EventObserver // 事件观察者映射，按事件类型分组
	eventMeta       *EventMetadataProvider            // 事件元数据提供者
	currentStateSeq int64                             // 当前状态版本号（由事件seq驱动，用于视图版本控制）
	mutex           sync.RWMutex                      // 读写锁，保护并发访问游戏状态
	createdAt       time.Time                         // 游戏引擎创建时间
	updatedAt       time.Time                         // 最后更新时间
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
	//   selectedCard: 选择的牌对象
	// 返回值:
	//   error: 如果选择无效或不是该玩家的选择回合，返回错误
	// 功能说明:
	//   - 用于双下情况下rank1玩家从贡牌池选择一张牌
	//   - 验证玩家是否有权选择
	//   - 验证选择的牌是否在贡牌池中
	//   - 自动将剩余的牌分配给rank2
	SubmitTributeSelection(playerID int, selectedCard *Card) error

	// SubmitReturnTribute 提交还贡
	// 参数:
	//   playerID: 还贡的玩家座位号(0-3)
	//   returnCard: 还贡的牌对象
	// 返回值:
	//   error: 如果还贡无效或不是该玩家的还贡回合，返回错误
	// 功能说明:
	//   - 用于收到贡牌的玩家选择一张牌还给贡牌方
	//   - 验证玩家是否需要还贡
	//   - 验证选择的牌是否在该玩家手中
	//   - 完成牌的交换
	SubmitReturnTribute(playerID int, returnCard *Card) error

	// 状态查询

	// GetGameState 获取当前完整的游戏状态
	// 返回值:
	//   *GameState: 包含游戏ID、状态、当前比赛等信息的完整状态
	// 功能说明:
	//   - 返回游戏的全局状态信息
	//   - 包括比赛进度、当前牌局状态等
	//   - 适用于管理员或观察者视角
	GetGameState() *GameState

	// GetEventMetadataProvider 获取事件元数据提供者
	// 返回值:
	//   *EventMetadataProvider: 用于生成事件序列号和元数据的提供者
	// 功能说明:
	//   - 返回引擎的事件元数据提供者
	//   - 确保所有事件使用统一的序列号生成器
	GetEventMetadataProvider() *EventMetadataProvider

	// GetPlayerView 获取玩家视角的游戏状态
	// 参数:
	//   playerSeat: 玩家座位号(0-3)
	// 返回值:
	//   *viewpb.PlayerView: 玩家视角游戏状态（proto类型）
	// 功能说明:
	//   - 返回从指定玩家角度看到的游戏状态
	//   - 包含玩家手牌、比赛等级、当前局信息、Trick信息等
	//   - 隐藏其他玩家的手牌信息
	//   - 返回proto定义的PlayerView消息
	GetPlayerView(playerSeat int) *viewpb.PlayerView

	// GetTributeView 获取进贡阶段视角的游戏状态
	// 参数:
	//   playerSeat: 玩家座位号(0-3)（预留，当前未使用）
	// 返回值:
	//   *viewpb.TributeView: 进贡阶段视角游戏状态（proto类型）
	// 功能说明:
	//   - 仅在tribute阶段返回有效数据
	//   - 包含进贡关系、贡牌池、选牌玩家等信息
	//   - 返回proto定义的TributeView消息
	GetTributeView(playerSeat int) *viewpb.TributeView

	// IsGameFinished 检查游戏是否已结束
	// 返回值:
	//   bool: 如果游戏已结束返回true，否则返回false
	// 功能说明:
	//   - 快速检查游戏是否处于结束状态
	//   - 用于判断是否还可以进行游戏操作
	IsGameFinished() bool

	// 事件处理

	// RegisterObserver 注册事件观察者（推荐使用）
	// 参数:
	//   eventType: 要监听的事件类型
	//   observer: 事件观察者接口
	// 功能说明:
	//   - 允许外部系统监听游戏事件
	//   - 支持多个观察者监听同一事件类型
	//   - 事件观察者在独立的协程中执行，不会阻塞游戏进程
	RegisterObserver(eventType GameEventType, observer EventObserver)

	// On 语法糖方法，方便函数式注册事件处理器（推荐使用）
	// 参数:
	//   eventType: 要监听的事件类型
	//   handler: 事件处理函数
	// 功能说明:
	//   - 提供更简洁的事件注册方式
	//   - 内部自动包装为EventObserver
	//   - 功能等同于RegisterObserver
	On(eventType GameEventType, handler func(*GameEvent))

	// RegisterEventHandler 注册游戏事件处理器（已废弃，使用On代替）
	// Deprecated: 使用 RegisterObserver 或 On 代替
	// 参数:
	//   eventType: 要监听的事件类型
	//   handler: 事件处理函数
	// 功能说明:
	//   - 允许外部系统监听游戏事件
	//   - 支持多个处理器监听同一事件类型
	//   - 事件处理器在独立的协程中执行，不会阻塞游戏进程
	//   - 常用事件包括：出牌、过牌、牌局结束、比赛结束等
	RegisterEventHandler(eventType GameEventType, handler GameEventHandler)

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
		id:        generateID(),
		status:    GameStatusWaiting,
		observers: make(map[GameEventType][]EventObserver),
		eventMeta: NewEventMetadataProvider(),
		createdAt: now,
		updatedAt: now,
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

	// Emit match started event (use the players parameter directly)
	event := NewMatchStartedEvent(ge.eventMeta, match, players, match.TeamLevels)
	ge.emitEventLocked(event)

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
	event := NewDealStartedEvent(ge.eventMeta, ge.currentMatch, ge.currentMatch.CurrentDeal, 
		ge.currentMatch.CurrentDeal.Level, ge.currentMatch.TeamLevels)
	ge.emitEventLocked(event)

	// 如果有上贡阶段，发送 TributeStarted 事件
	if ge.currentMatch.CurrentDeal.TributePhase != nil {
		lastResult := ge.currentMatch.CurrentDeal.LastResult
		
		// 确定上贡类型
		var tributeType string
		switch lastResult.VictoryType {
		case VictoryTypeDoubleDown:
			tributeType = "double_down"
		case VictoryTypeSingleLast:
			tributeType = "single_last"
		case VictoryTypePartnerLast:
			tributeType = "partner_last"
		default:
			tributeType = "none"
		}

		// 提取 givers 和 receivers
		var givers []int
		var receivers []int
		tributePhase := ge.currentMatch.CurrentDeal.TributePhase
		
		// 从 TributePairs 提取 givers（需要排序以保证事件一致性）
		for _, pair := range tributePhase.TributePairs {
			givers = append(givers, pair.Giver)
		}
		// 对 givers 排序，确保事件中的数组顺序一致
		sort.Ints(givers)
		
		// 确定 receivers（根据上贡类型）
		rank1 := lastResult.Rankings[0]
		switch lastResult.VictoryType {
		case VictoryTypeDoubleDown:
			rank2 := lastResult.Rankings[1]
			receivers = []int{rank1, rank2}
		case VictoryTypeSingleLast, VictoryTypePartnerLast:
			receivers = []int{rank1}
		}

		// 发送 TributeStarted 事件
		tributeStartedEvent := NewTributeStartedEvent(ge.eventMeta, ge.currentMatch, 
			ge.currentMatch.CurrentDeal, tributeType, givers, receivers)
		ge.emitEventLocked(tributeStartedEvent)

		// 检查是否触发抗贡
		if tributePhase.IsImmune {
			// 获取详细的抗贡信息
			tm := NewTributeManager(ge.currentMatch.CurrentDeal.Level)
			_, immunityDetails := tm.GetTributeImmunityDetails(lastResult, 
				ge.currentMatch.CurrentDeal.PlayerCards)

			// 提取 joker_holders
			jokerHolders := make(map[int]int)
			if immunityDetails != nil {
				if holders, ok := immunityDetails["big_joker_holders"].([]map[string]interface{}); ok {
					for _, holder := range holders {
						if seat, ok := holder["player_seat"].(int); ok {
							if count, ok := holder["big_joker_count"].(int); ok {
								jokerHolders[seat] = count
							}
						}
					}
				}
			}

			// 发送 TributeExempted 事件
			tributeExemptedEvent := NewTributeExemptedEvent(ge.eventMeta, ge.currentMatch, 
				ge.currentMatch.CurrentDeal, jokerHolders)
			ge.emitEventLocked(tributeExemptedEvent)

			// 注意：不在这里发送 TributeCompleted，统一在 ProcessTributePhase 中发送
		}
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
		return nil, err
	}

	// Check for pre-action state transitions (e.g., trick starting) BEFORE executing the play
	preEvents := ge.checkPreActionStateTransitions()
	for _, evt := range preEvents {
		ge.emitEventLocked(evt)
	}

	// Execute the play
	err = deal.PlayCards(playerSeat, cards)
	if err != nil {
		return nil, fmt.Errorf("failed to play cards: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create and emit player played event
	event := NewPlayerPlayedEvent(ge.eventMeta, ge.currentMatch, deal, deal.CurrentTrick, playerSeat, cards)
	ge.emitEventLocked(event)

	// Check for post-action state transitions (e.g., trick ending, deal ending)
	postEvents := ge.checkPostActionStateTransitions()
	for _, evt := range postEvents {
		ge.emitEventLocked(evt)
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
		ge.emitEventLocked(evt)
	}

	// Execute the pass
	err = deal.PassTurn(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to pass turn: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create and emit player passed event
	event := NewPlayerPassedEvent(ge.eventMeta, ge.currentMatch, deal, deal.CurrentTrick, playerSeat)
	ge.emitEventLocked(event)

	// Check for post-action state transitions (e.g., trick ending, deal ending)
	postEvents := ge.checkPostActionStateTransitions()
	for _, evt := range postEvents {
		ge.emitEventLocked(evt)
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

// GetEventMetadataProvider returns the event metadata provider
func (ge *GameEngine) GetEventMetadataProvider() *EventMetadataProvider {
	return ge.eventMeta
}

// GetPlayerView 返回玩家视角的游戏状态（proto类型）
func (ge *GameEngine) GetPlayerView(playerSeat int) *viewpb.PlayerView {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	// 1. 构建SDK内部视图
	sdkView := ge.buildInternalPlayerView(playerSeat)
	if sdkView == nil {
		return nil
	}

	// 2. 读取当前状态版本号
	seq := atomic.LoadInt64(&ge.currentStateSeq)

	// 3. 转换为proto
	return ConvertPlayerViewToProto(
		sdkView,
		ge.currentMatch.ID,
		len(ge.currentMatch.DealHistory),
		seq,
	)
}

// buildInternalPlayerView 构建SDK内部的PlayerView（保留原有逻辑）
func (ge *GameEngine) buildInternalPlayerView(playerSeat int) *PlayerView {
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	deal := ge.currentMatch.CurrentDeal
	
	view := &PlayerView{
		PlayerSeat: playerSeat,
		TeamLevels: ge.currentMatch.TeamLevels,
		DealLevel:  deal.Level,
		DealStatus: deal.Status,
	}

	if playerSeat >= 0 && playerSeat < 4 {
		cards := deal.PlayerCards[playerSeat]
		if cards == nil {
			cards = make([]*Card, 0)
		}
		view.PlayerCards = cards
	}

	// 根据 deal.Status 填充不同的字段
	if deal.Status == DealStatusPlaying && deal.CurrentTrick != nil {
		view.TrickID = deal.CurrentTrick.ID
		turn := deal.CurrentTrick.CurrentTurn
		view.CurrentTurn = &turn
		leader := deal.CurrentTrick.Leader
		view.Leader = &leader
		
		plays := deal.CurrentTrick.Plays
		if plays == nil {
			plays = make([]*PlayAction, 0)
		}
		view.Plays = plays
	} else if deal.Status == DealStatusTribute && deal.TributePhase != nil {
		view.TributePhase = deal.TributePhase
	}

	return view
}

// GetTributeView 返回进贡阶段视角的游戏状态（proto类型）
func (ge *GameEngine) GetTributeView(playerSeat int) *viewpb.TributeView {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	deal := ge.currentMatch.CurrentDeal
	if deal.Status != DealStatusTribute || deal.TributePhase == nil {
		return nil
	}

	// 读取当前状态版本号
	seq := atomic.LoadInt64(&ge.currentStateSeq)

	// 转换为proto
	return ConvertTributeViewToProto(
		deal.TributePhase,
		ge.currentMatch.ID,
		len(ge.currentMatch.DealHistory),
		seq,
	)
}

// IsGameFinished checks if the game is finished
func (ge *GameEngine) IsGameFinished() bool {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	return ge.status == GameStatusFinished
}

// RegisterObserver registers an event observer for a specific event type
//
// Concurrency contract for observers:
//   - Observers are called asynchronously in separate goroutines
//   - Observers should be non-blocking and fast-executing
//   - Observers must not make synchronous calls back into GameEngine methods
//     while expecting immediate state changes during event processing
//   - Event delivery order is not guaranteed across different event types
//   - Multiple observers for the same event type will all be called
//   - Observers should handle panics internally or they will be recovered and logged
func (ge *GameEngine) RegisterObserver(eventType GameEventType, observer EventObserver) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if ge.observers[eventType] == nil {
		ge.observers[eventType] = make([]EventObserver, 0)
	}
	ge.observers[eventType] = append(ge.observers[eventType], observer)
}

// On is a convenience method for registering event handlers using functions
// See RegisterObserver for concurrency contract details
func (ge *GameEngine) On(eventType GameEventType, handler func(*GameEvent)) {
	ge.RegisterObserver(eventType, EventHandlerFunc(handler))
}

// RegisterEventHandler registers an event handler (legacy interface for backward compatibility)
// Deprecated: Use RegisterObserver or On instead
// See RegisterObserver for concurrency contract details
func (ge *GameEngine) RegisterEventHandler(eventType GameEventType, handler GameEventHandler) {
	ge.RegisterObserver(eventType, EventHandlerFunc(handler))
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
	var deal *Deal
	var trick *Trick
	if ge.currentMatch.CurrentDeal != nil {
		deal = ge.currentMatch.CurrentDeal
		trick = deal.CurrentTrick
	}
	event := NewPlayerDisconnectEvent(ge.eventMeta, ge.currentMatch, deal, trick, playerSeat, true)
	ge.emitEventLocked(event)

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
	var deal *Deal
	var trick *Trick
	if ge.currentMatch.CurrentDeal != nil {
		deal = ge.currentMatch.CurrentDeal
		trick = deal.CurrentTrick
	}
	event := NewPlayerReconnectEvent(ge.eventMeta, ge.currentMatch, deal, trick, playerSeat, false)
	ge.emitEventLocked(event)

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

// emitEvent emits an event to all registered observers
func (ge *GameEngine) emitEvent(event *GameEvent) {
	// Safely read and copy the observer list to avoid holding the lock during callback execution
	ge.mutex.RLock()
	observers, exists := ge.observers[event.Type]
	if !exists || len(observers) == 0 {
		ge.mutex.RUnlock()
		return
	}
	// Create a copy of the observer slice to prevent race conditions
	observersCopy := make([]EventObserver, len(observers))
	copy(observersCopy, observers)
	ge.mutex.RUnlock()

	// Call all observers for this event type asynchronously to avoid deadlock
	// Each observer runs in its own goroutine to prevent blocking the engine
	// This is crucial to avoid deadlock when observers call back into the engine
	for _, observer := range observersCopy {
		// Create a copy of observer for the goroutine closure
		obs := observer
		go func() {
			// Use defer to recover from any panic in the observer
			defer func() {
				if r := recover(); r != nil {
					// Log the panic but don't crash the engine
					// TODO: Replace with pluggable logger interface
					fmt.Printf("[GameEngine] Event observer panic for %s: %v\n", event.Type, r)
				}
			}()
			obs.OnGameEvent(event)
		}()
	}
}

// Note: Observer contract
// - Observers run asynchronously and should be non-blocking
// - Observers should not re-enter engine methods that depend on immediate state changes while locks are held
// - Event delivery is not guaranteed to be ordered across different event types

// emitEventLocked is called when the caller already holds ge.mutex lock
// It reads the observers without acquiring additional locks to avoid deadlock
// Note: All calls to this method must be inside a ge.mutex.Lock() block
func (ge *GameEngine) emitEventLocked(event *GameEvent) {
	// 更新当前状态版本号
	if event != nil && event.Seq > 0 {
		atomic.StoreInt64(&ge.currentStateSeq, event.Seq)
	}

	// Caller already holds the lock, so we can safely read observers
	observers, exists := ge.observers[event.Type]
	if !exists || len(observers) == 0 {
		return
	}
	// Create a copy of the observer slice
	observersCopy := make([]EventObserver, len(observers))
	copy(observersCopy, observers)

	// Call all observers asynchronously outside the lock
	for _, observer := range observersCopy {
		obs := observer
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// TODO: Replace with pluggable logger interface
					fmt.Printf("[GameEngine] Event observer panic for %s: %v\n", event.Type, r)
				}
			}()
			obs.OnGameEvent(event)
		}()
	}
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
			// Calculate remaining players (all 4 players at start of trick)
			remainingPlayers := []int{}
			for i := 0; i < 4; i++ {
				if len(deal.PlayerCards[i]) > 0 {
					remainingPlayers = append(remainingPlayers, i)
				}
			}
			isFirstTrick := len(deal.TrickHistory) == 0
			// Note: Timeout is now managed by GameDriver, not set here
			trickStartedEvent := NewTrickStartedEvent(ge.eventMeta, ge.currentMatch, deal, 
				deal.CurrentTrick, deal.CurrentTrick.Leader, isFirstTrick, remainingPlayers)
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
		durationMs := dealResult.Duration.Milliseconds()
		dealEndedEvent := NewDealEndedEvent(ge.eventMeta, ge.currentMatch, deal,
			deal.Level, deal.Rankings, dealResult.VictoryType, dealResult.WinningTeam,
			dealResult.Upgrades, durationMs, len(deal.TrickHistory))
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
				durationMs := matchResult.Duration.Milliseconds()
				matchEndedEvent := NewMatchEndedEvent(ge.eventMeta, ge.currentMatch,
					ge.currentMatch.Winner, ge.currentMatch.TeamLevels, durationMs, len(ge.currentMatch.DealHistory))
				events = append(events, matchEndedEvent)
			}
		}
	} else if deal.CurrentTrick != nil && deal.CurrentTrick.Status == TrickStatusFinished {
		// Check if current trick is finished
		// Emit trick ended event
		finishedTrick := deal.CurrentTrick
		trickEndedEvent := NewTrickEndedEvent(ge.eventMeta, ge.currentMatch, deal, finishedTrick, finishedTrick.Winner)
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

	// 记录处理前的状态
	previousStatus := deal.TributePhase.Status

	// 获取 TributeManager 并处理
	tm := NewTributeManager(deal.Level)
	action, err := tm.ProcessTributePhaseAction(deal.TributePhase, deal.PlayerCards)
	if err != nil {
		return nil, err
	}

	// 检测状态变化并触发相应事件
	currentStatus := deal.TributePhase.Status

	// 1. Waiting → Selecting: 发送 TributeCardSubmitted 事件（每张贡牌）
	if previousStatus == TributeStatusWaiting && currentStatus == TributeStatusSelecting {
		// 遍历所有贡牌，发送提交事件
		for _, pair := range deal.TributePhase.TributePairs {
			if pair.TributeCard != nil {
				event := NewTributeCardSubmittedEvent(ge.eventMeta, ge.currentMatch, deal, pair.Giver, pair.TributeCard)
				ge.emitEventLocked(event)
			}
		}
	}

	// 2. Selecting → Returning: 单下/末游自动选贡，发送 TributeCardSelected 事件
	if previousStatus == TributeStatusSelecting && currentStatus == TributeStatusReturning {
		// 检查是否是自动选择（单下/末游）
		// TributePairs 中应该只有1个条目
		if len(deal.TributePhase.TributePairs) == 1 {
			pair := deal.TributePhase.TributePairs[0]
			if pair.Receiver != -1 && pair.TributeCard != nil {
				event := NewTributeCardSelectedEvent(ge.eventMeta, ge.currentMatch, deal, pair.Receiver, pair.TributeCard, true)
				ge.emitEventLocked(event)
			}
		}
	}

	// 3. 如果贡牌阶段完成，应用贡牌效果并发送完成事件
	if currentStatus == TributeStatusFinished {
		// 抗贡场景：需要先启动游戏阶段（创建第一个trick）
		if deal.TributePhase.IsImmune {
			// 启动游戏阶段
			err = deal.startFirstTrick()
			if err != nil {
				return nil, fmt.Errorf("failed to start first trick: %w", err)
			}
			deal.Status = DealStatusPlaying
		} else {
			// 正常上贡场景：应用贡牌效果
			err = tm.ApplyTributeToHands(deal.TributePhase, &deal.PlayerCards)
			if err != nil {
				return nil, fmt.Errorf("apply tribute failed: %w", err)
			}

			// 启动游戏阶段（包括创建第一个trick和设置状态）
			err = deal.StartPlayingPhase()
			if err != nil {
				return nil, fmt.Errorf("failed to start playing phase: %w", err)
			}
		}

		// 统一发送完成事件（抗贡和正常上贡都在这里发送）
		ge.emitEventLocked(NewTributeCompletedEvent(ge.eventMeta, ge.currentMatch, deal))
	}

	return action, nil
}

// SubmitTributeSelection 提交贡牌选择（用于双下选牌）
func (ge *GameEngine) SubmitTributeSelection(playerID int, selectedCard *Card) error {
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

	if selectedCard == nil {
		return errors.New("selected card is nil")
	}

	// 记录选择前的已分配贡牌数量，用于检测 rank2 的自动分配
	var previousReceivedCount int
	for _, pair := range deal.TributePhase.TributePairs {
		if pair.Receiver != -1 {
			previousReceivedCount++
		}
	}

	// 调用 TributeManager 处理选择
	tm := NewTributeManager(deal.Level)
	err := tm.SubmitSelection(deal.TributePhase, playerID, selectedCard)
	if err != nil {
		return err
	}

	// 发送 rank1 的选择事件
	ge.emitEventLocked(NewTributeCardSelectedEvent(ge.eventMeta, ge.currentMatch, deal, 
		playerID, selectedCard, false))

	// 检测是否是双下场景（rank2 自动获得剩余牌）
	var currentReceivedCount int
	for _, pair := range deal.TributePhase.TributePairs {
		if pair.Receiver != -1 {
			currentReceivedCount++
		}
	}
	if currentReceivedCount > previousReceivedCount && currentReceivedCount == 2 {
		// rank2 自动获得了剩余牌，发送 rank2 的选择事件
		// 基于 TributePairs 最终状态找到 rank2 和对应的卡片
		for _, pair := range deal.TributePhase.TributePairs {
			if pair.Receiver != -1 && pair.Receiver != playerID && pair.TributeCard != nil {
				// 找到了 rank2 的 pair
				ge.emitEventLocked(NewTributeCardSelectedEvent(ge.eventMeta, ge.currentMatch, deal, 
					pair.Receiver, pair.TributeCard, true))
				break
			}
		}
	}

	return nil
}

// SubmitReturnTribute 提交还贡
func (ge *GameEngine) SubmitReturnTribute(playerID int, returnCard *Card) error {
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

	if returnCard == nil {
		return errors.New("return card is nil")
	}

	// 从 TributePairs 获取 target_player（原贡者）
	var targetPlayer int = -1
	for _, pair := range deal.TributePhase.TributePairs {
		if pair.Receiver == playerID {
			targetPlayer = pair.Giver
			break
		}
	}
	if targetPlayer == -1 {
		return fmt.Errorf("player %d is not a tribute receiver", playerID)
	}

	// 调用 TributeManager 处理还贡
	tm := NewTributeManager(deal.Level)
	err := tm.SubmitReturn(deal.TributePhase, playerID, returnCard, deal.PlayerCards[playerID])
	if err != nil {
		return err
	}

	// 发送还贡事件
	ge.emitEventLocked(NewTributeCardReturnedEvent(ge.eventMeta, ge.currentMatch, deal, 
		playerID, returnCard, targetPlayer, false))

	return nil
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
