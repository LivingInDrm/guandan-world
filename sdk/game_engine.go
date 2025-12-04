package sdk

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"guandan-world/pkg/log"
	eventpb "guandan-world/proto/event"
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
	PlayState    [4]PlayState  `json:"play_state"`              // 每个玩家在当前Trick中的状态（playing阶段）
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
	id                 string                            // 游戏引擎的唯一标识符
	status             GameStatus                        // 当前游戏状态
	currentMatch       *Match                            // 当前活跃的比赛实例
	observers          map[GameEventType][]EventObserver // 事件观察者映射，按事件类型分组
	eventMeta          *EventMetadataProvider            // 事件元数据提供者
	currentStateSeq    int64                             // 当前状态版本号（由事件seq驱动，用于视图版本控制）
	mutex              sync.RWMutex                      // 读写锁，保护并发访问游戏状态
	eventWg            sync.WaitGroup                    // 追踪异步事件完成状态
	createdAt          time.Time                         // 游戏引擎创建时间
	updatedAt          time.Time                         // 最后更新时间
	dealEndedDelay     time.Duration                     // 牌局结束后的延迟时间（用于计算 deadline）
	dealEndedPayload   *eventpb.DealEndedPayload         // 牌局结束结果（用于 PlayerView 展示）
	matchEndedPayload  *eventpb.MatchEndedPayload        // 比赛结束结果（用于 PlayerView 展示）
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

	// StepTribute 单步处理贡牌阶段
	// 参数:
	//   input: 用户输入（选牌或还贡），nil 表示推进状态机
	// 返回值:
	//   *StepResult: 单步处理结果，包含是否需要用户输入、是否完成、状态转换信息
	//   error: 如果处理失败，返回错误
	// 功能说明:
	//   - 每次调用只处理一步，由 GameDriver 负责循环控制
	//   - input == nil: 推进状态机一步
	//   - input != nil: 处理用户输入，然后推进一步
	//   - 返回 Action != nil 表示需要用户输入
	//   - 返回 Completed == true 表示贡牌阶段完成
	//   - 返回 StatusChanged == true 表示发生状态转换（用于延迟控制）
	StepTribute(input *TributeInput) (*StepResult, error)

	// StartPlayingPhase 启动出牌阶段
	// 返回值:
	//   error: 如果启动失败，返回错误
	// 功能说明:
	//   - 在贡牌阶段完成后调用
	//   - 初始化第一个 Trick
	//   - 将 Deal 状态从 tribute 切换到 playing
	StartPlayingPhase() error

	// StartTrickIfNeeded 启动当前 Trick（如果尚未启动）
	// 返回值:
	//   error: 如果启动失败，返回错误
	// 功能说明:
	//   - 检查当前 Trick 是否处于 Waiting 状态
	//   - 如果是，将状态改为 Playing 并发送 trick_started 事件
	//   - 应在请求玩家决策之前调用
	StartTrickIfNeeded() error

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

	// FlushEvents 等待所有异步事件发送完成
	// 功能说明:
	//   - 在 match 结束时调用，确保所有事件在返回前被处理
	//   - 用于解决 match 结束时事件未发送到前端的问题
	FlushEvents()

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

	// GetPlayerHand 返回玩家当前手牌的原始 SDK Card 对象
	//
	// 用途:
	//   - SDK 内部（GameDriver）使用，避免 proto 转换损失 Level 字段
	//   - 保留 Card 对象的完整性（Level、RawNumber 等）
	//
	// 与 GetPlayerView 的区别:
	//   - GetPlayerView: 返回 proto 消息，适合跨进程/序列化传输
	//   - GetPlayerHand:  返回 SDK 原始类型，适合内部逻辑使用
	//
	// 参数:
	//   playerSeat: 玩家座位号 (0-3)
	//
	// 返回值:
	//   []*Card: 玩家手牌列表（如果无牌局或座位号无效返回 nil 或空切片）
	GetPlayerHand(playerSeat int) []*Card

	// SetDealEndedDelay 设置牌局结束后的延迟时间
	// 参数:
	//   delay: 延迟时间
	// 功能说明:
	//   - 用于计算 DealEnded 事件中的 NextDealDeadlineMs 字段
	//   - 应在开始比赛前调用
	SetDealEndedDelay(delay time.Duration)
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

// SetDealEndedDelay 设置牌局结束后的延迟时间
func (ge *GameEngine) SetDealEndedDelay(delay time.Duration) {
	ge.dealEndedDelay = delay
}

// requireActiveMatch 检测是否存在活跃比赛（Public API 用）
func (ge *GameEngine) requireActiveMatch() (*Match, error) {
	if ge.currentMatch == nil {
		log.Warn("require failed: no active match")
		return nil, errors.New("no active match")
	}
	return ge.currentMatch, nil
}

// requireActiveDeal 检测是否存在活跃牌局（Public API 用）
func (ge *GameEngine) requireActiveDeal() (*Deal, error) {
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		log.Warn("require failed: no active deal")
		return nil, errors.New("no active deal")
	}
	return ge.currentMatch.CurrentDeal, nil
}

// mustActiveMatch 断言存在活跃比赛（Internal 用）
func (ge *GameEngine) mustActiveMatch() *Match {
	if ge.currentMatch == nil {
		log.Error("must failed: no active match")
		panic("must: no active match")
	}
	return ge.currentMatch
}

// mustActiveDeal 断言存在活跃牌局（Internal 用）
func (ge *GameEngine) mustActiveDeal() *Deal {
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		log.Error("must failed: no active deal")
		panic("must: no active deal")
	}
	return ge.currentMatch.CurrentDeal
}

// StartMatch initializes a new match with the given players
func (ge *GameEngine) StartMatch(players []Player) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	if len(players) != 4 {
		log.Warn("invalid player count", "player_count", len(players))
		return errors.New("exactly 4 players are required")
	}

	if ge.status != GameStatusWaiting {
		log.Warn("invalid game status", "game_status", ge.status)
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

	match, err := ge.requireActiveMatch()
	if err != nil {
		return err
	}

	err = match.StartNewDeal()
	if err != nil {
		return fmt.Errorf("failed to start deal: %w", err)
	}

	ge.dealEndedPayload = nil
	ge.updatedAt = time.Now()

	deal := match.CurrentDeal

	// Emit deal started event
	event := NewDealStartedEvent(ge.eventMeta, match, deal, deal.Level, match.TeamLevels)
	ge.emitEventLocked(event)

	// 如果有上贡阶段，发送 TributeStarted 事件
	if deal.TributePhase != nil {
		lastResult := deal.LastResult

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
		tributePhase := deal.TributePhase

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
		tributeStartedEvent := NewTributeStartedEvent(ge.eventMeta, match, deal, tributeType, givers, receivers)
		ge.emitEventLocked(tributeStartedEvent)

		// 检查是否触发抗贡
		if tributePhase.IsImmune {
			// 获取详细的抗贡信息
			tm := NewTributeManager(deal.Level)
			_, immunityDetails := tm.GetTributeImmunityDetails(lastResult, deal.PlayerCards)

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
			tributeExemptedEvent := NewTributeExemptedEvent(ge.eventMeta, match, deal, jokerHolders)
			ge.emitEventLocked(tributeExemptedEvent)

			// 注意：不在这里发送 TributeCompleted，统一在 StepTribute 中发送
		}
	}

	return nil
}

// PlayCards handles a player's card play action
func (ge *GameEngine) PlayCards(playerSeat int, cards []*Card) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	deal, err := ge.requireActiveDeal()
	if err != nil {
		return nil, err
	}

	// 根据传入牌的 DeckIndex，从玩家手牌中查找原始牌（确保 Level 等属性正确）
	deckIndexes, err := extractDeckIndexes(cards)
	if err != nil {
		return nil, fmt.Errorf("failed to extract deck indexes: %w", err)
	}
	originalCards, err := findCardsByDeckIndexes(deal.PlayerCards[playerSeat], deckIndexes)
	if err != nil {
		return nil, fmt.Errorf("failed to find cards in player hand: %w", err)
	}

	// Execute the play with original cards
	err = deal.PlayCards(playerSeat, originalCards)
	if err != nil {
		return nil, fmt.Errorf("failed to play cards: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create and emit player played event with original cards
	event := NewPlayerPlayedEvent(ge.eventMeta, ge.currentMatch, deal, deal.CurrentTrick, playerSeat, originalCards)
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

	deal, err := ge.requireActiveDeal()
	if err != nil {
		return nil, err
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
	protoView := ConvertPlayerViewToProto(
		sdkView,
		ge.currentMatch.ID,
		len(ge.currentMatch.DealHistory),
		seq,
	)

	// 4. 填充结果字段
	protoView.DealResult = ge.dealEndedPayload
	protoView.MatchResult = ge.matchEndedPayload

	return protoView
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
		view.PlayState = deal.PlayState
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

// HandlePlayerDisconnect handles a player disconnection
func (ge *GameEngine) HandlePlayerDisconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	match, err := ge.requireActiveMatch()
	if err != nil {
		return nil, err
	}

	err = match.HandlePlayerDisconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle disconnect: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create disconnect event
	var deal *Deal
	var trick *Trick
	if match.CurrentDeal != nil {
		deal = match.CurrentDeal
		trick = deal.CurrentTrick
	}
	event := NewPlayerDisconnectEvent(ge.eventMeta, match, deal, trick, playerSeat, true)
	ge.emitEventLocked(event)

	return event, nil
}

// HandlePlayerReconnect handles a player reconnection
func (ge *GameEngine) HandlePlayerReconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	match, err := ge.requireActiveMatch()
	if err != nil {
		return nil, err
	}

	err = match.HandlePlayerReconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle reconnect: %w", err)
	}

	ge.updatedAt = time.Now()

	// Create reconnect event
	var deal *Deal
	var trick *Trick
	if match.CurrentDeal != nil {
		deal = match.CurrentDeal
		trick = deal.CurrentTrick
	}
	event := NewPlayerReconnectEvent(ge.eventMeta, match, deal, trick, playerSeat, false)
	ge.emitEventLocked(event)

	return event, nil
}

// SetPlayerAutoPlay sets the auto-play status for a player
func (ge *GameEngine) SetPlayerAutoPlay(playerSeat int, enabled bool) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	match, err := ge.requireActiveMatch()
	if err != nil {
		return err
	}

	return match.SetPlayerAutoPlay(playerSeat, enabled)
}

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
		ge.eventWg.Add(1)
		go func() {
			defer ge.eventWg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Error("event observer panic", "event_type", event.Type, "panic", r)
				}
			}()
			obs.OnGameEvent(event)
		}()
	}
}

// FlushEvents 等待所有异步事件发送完成
// 在 match 结束时调用，确保所有事件（如 MatchEnded）在返回前被处理
func (ge *GameEngine) FlushEvents() {
	ge.eventWg.Wait()
}

// StartTrickIfNeeded checks if there's a waiting trick and starts it by emitting trick_started event
// This should be called before requesting player decisions in GameDriver
func (ge *GameEngine) StartTrickIfNeeded() error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	deal, err := ge.requireActiveDeal()
	if err != nil {
		return nil
	}

	if deal.CurrentTrick != nil && !deal.CurrentTrick.Started {
		deal.CurrentTrick.Started = true

		remainingPlayers := []int{}
		for i := 0; i < 4; i++ {
			if len(deal.PlayerCards[i]) > 0 {
				remainingPlayers = append(remainingPlayers, i)
			}
		}
		isFirstTrick := len(deal.TrickHistory) == 0
		trickStartedEvent := NewTrickStartedEvent(ge.eventMeta, ge.currentMatch, deal,
			deal.CurrentTrick, deal.CurrentTrick.Leader, isFirstTrick, remainingPlayers)
		ge.emitEventLocked(trickStartedEvent)
	}

	return nil
}

// checkPostActionStateTransitions checks for state transitions that should happen after player actions
// Handles trick ending and deal ending
func (ge *GameEngine) checkPostActionStateTransitions() []*GameEvent {
	events := make([]*GameEvent, 0)

	match := ge.mustActiveMatch()
	deal := ge.mustActiveDeal()

	// Check if deal is finished first (can happen at any time)
	if deal.Status == DealStatusFinished {
		// Calculate deal result using the new result system
		dealResult, err := deal.CalculateResult(match)
		if err != nil {
			log.Error("failed to calculate deal result", "error", err)
			winningTeam := match.GetTeamForPlayer(deal.Rankings[0])
			upgrades := [2]int{0, 0}
			upgrades[winningTeam] = 1

			dealResult = &DealResult{
				Rankings:    deal.Rankings,
				WinningTeam: winningTeam,
				VictoryType: VictoryTypePartnerLast,
				Upgrades:    upgrades,
				Duration:    time.Since(deal.StartTime),
				TrickCount:  len(deal.TrickHistory),
			}
		}

		// Get player stats (use empty stats if Statistics is nil)
		var playerStats [4]*PlayerDealStats
		if dealResult.Statistics != nil {
			playerStats = dealResult.Statistics.PlayerStats
		} else {
			// Create empty stats for fallback case
			for i := 0; i < 4; i++ {
				playerStats[i] = &PlayerDealStats{
					PlayerSeat:  i,
					CardsPlayed: 0,
					TricksWon:   0,
					PassCount:   0,
					FinishRank:  0,
				}
			}
		}

		// Update match with deal result first to know if match will finish
		err = match.FinishDeal(dealResult)
		matchFinished := err == nil && match.Status == MatchStatusFinished

		// Calculate next deal deadline
		var nextDealDeadlineMs int64 = 0
		if !matchFinished && ge.dealEndedDelay > 0 {
			nextDealDeadlineMs = time.Now().Add(ge.dealEndedDelay).UnixMilli()
		}

		// Emit deal ended event
		durationMs := dealResult.Duration.Milliseconds()
		dealEndedEvent := NewDealEndedEvent(ge.eventMeta, match, deal,
			deal.Level, deal.Rankings, dealResult.VictoryType, dealResult.WinningTeam,
			dealResult.Upgrades, durationMs, len(deal.TrickHistory), playerStats,
			nextDealDeadlineMs)
		events = append(events, dealEndedEvent)
		ge.dealEndedPayload = dealEndedEvent.GetDealEnded()

		// Check if match is finished
		if matchFinished {
			ge.status = GameStatusFinished

			// Create match result
			matchResult := ge.createMatchResult()

			// Emit match ended event
			durationMs := matchResult.Duration.Milliseconds()
			matchEndedEvent := NewMatchEndedEvent(ge.eventMeta, match,
				match.Winner, match.TeamLevels, durationMs, len(match.DealHistory))
			events = append(events, matchEndedEvent)
			ge.matchEndedPayload = matchEndedEvent.GetMatchEnded()
		}
	} else if deal.CurrentTrick != nil && deal.CurrentTrick.Winner >= 0 {
		// Check if current trick is finished (Winner >= 0 means trick has ended)
		// Emit trick ended event
		finishedTrick := deal.CurrentTrick
		trickEndedEvent := NewTrickEndedEvent(ge.eventMeta, match, deal, finishedTrick, finishedTrick.Winner)
		events = append(events, trickEndedEvent)

		// Add finished trick to history
		deal.TrickHistory = append(deal.TrickHistory, finishedTrick)

		// Create new trick with the next leader
		nextTrick, err := NewTrick(finishedTrick.NextLeader)
		if err == nil {
			deal.CurrentTrick = nextTrick
		}
	}

	return events
}

// createMatchResult creates a MatchResult from a finished match
func (ge *GameEngine) createMatchResult() *MatchResult {
	match := ge.mustActiveMatch()
	if match.Status != MatchStatusFinished {
		log.Error("must failed: match not finished", "match_status", match.Status)
		panic("must: match not finished")
	}

	// Calculate total duration
	duration := time.Duration(0)
	if match.EndTime != nil {
		duration = match.EndTime.Sub(match.StartTime)
	}

	// Calculate match statistics
	stats := &MatchStatistics{
		TotalDeals:    len(match.DealHistory),
		TotalDuration: duration,
		FinalLevels:   match.TeamLevels,
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
	for _, deal := range match.DealHistory {
		if result, err := deal.CalculateResult(match); err == nil {
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
						team := match.GetTeamForPlayer(playerStats.PlayerSeat)
						stats.TeamStats[team].TotalTricks += playerStats.TricksWon
					}
				}
			}
		}
	}

	return &MatchResult{
		Winner:      match.Winner,
		FinalLevels: match.TeamLevels,
		Duration:    duration,
		Statistics:  stats,
	}
}

// applyPairUpdates 应用 TributePair 变更
func applyPairUpdates(phase *TributePhase, updates []PairUpdate) {
	for _, u := range updates {
		for _, pair := range phase.TributePairs {
			if pair.Giver == u.GiverSeat {
				if u.Receiver != nil {
					pair.Receiver = *u.Receiver
				}
				if u.TributeCard != nil {
					pair.TributeCard = u.TributeCard
				}
				if u.ReturnCard != nil {
					pair.ReturnCard = u.ReturnCard
				}
				break
			}
		}
	}
}

// applyPoolChanges 应用 PoolCards 变更
func applyPoolChanges(phase *TributePhase, toSet []*Card, toRemove *Card) {
	if toSet != nil {
		phase.PoolCards = make([]*Card, len(toSet))
		copy(phase.PoolCards, toSet)
	}
	if toRemove != nil {
		for i, card := range phase.PoolCards {
			if card.DeckIndex == toRemove.DeckIndex {
				phase.PoolCards = append(phase.PoolCards[:i], phase.PoolCards[i+1:]...)
				break
			}
		}
	}
}

// applyHandChanges 应用手牌变更
func applyHandChanges(hands *[4][]*Card, changes []HandChange) {
	for _, c := range changes {
		if c.IsAdd {
			hands[c.PlayerSeat] = append(hands[c.PlayerSeat], c.Card)
		} else {
			hands[c.PlayerSeat] = removeCardFromSlice(hands[c.PlayerSeat], c.Card)
		}
	}
}

// removeCardFromSlice 从切片中移除指定牌
func removeCardFromSlice(cards []*Card, cardToRemove *Card) []*Card {
	for i, card := range cards {
		if card.DeckIndex == cardToRemove.DeckIndex {
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards
}

// buildTributeEvent 将事件意图转换为 GameEvent
func (ge *GameEngine) buildTributeEvent(intent TributeEventIntent, deal *Deal) *GameEvent {
	switch intent.Type {
	case EventTributeCardSubmitted:
		return NewTributeCardSubmittedEvent(ge.eventMeta, ge.currentMatch, deal, intent.PlayerSeat, intent.Card)
	case EventTributeCardSelected:
		return NewTributeCardSelectedEvent(ge.eventMeta, ge.currentMatch, deal, intent.PlayerSeat, intent.Card, intent.IsAuto)
	case EventReturnTribute:
		return NewTributeCardReturnedEvent(ge.eventMeta, ge.currentMatch, deal, intent.PlayerSeat, intent.Card, intent.TargetSeat, false)
	case EventTributeCompleted:
		return NewTributeCompletedEvent(ge.eventMeta, ge.currentMatch, deal)
	default:
		return nil
	}
}

// StepTribute 单步处理贡牌阶段
// 每次调用只处理一步，由 GameDriver 负责循环控制
// 参数:
//   - input: 用户输入（选牌或还贡），nil 表示推进状态机
//
// 返回值:
//   - StepResult.Action != nil: 需要用户输入
//   - StepResult.Completed: 贡牌阶段完成
//   - StepResult.StatusChanged: 状态发生转换（用于延迟控制）
func (ge *GameEngine) StepTribute(input *TributeInput) (*StepResult, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	deal, err := ge.requireActiveDeal()
	if err != nil {
		return nil, err
	}

	if deal.Status != DealStatusTribute {
		log.Warn("require failed: not in tribute phase", "deal_status", deal.Status)
		return nil, errors.New("not in tribute phase")
	}

	if deal.TributePhase == nil {
		return &StepResult{
			Completed:     true,
			PrevStatus:    TributeStatusFinished,
			StatusChanged: false,
		}, nil
	}

	phase := deal.TributePhase
	prevStatus := phase.Status

	result := ProcessTributeStep(phase, deal.PlayerCards, deal.Level, input)

	if result.Error != nil {
		return nil, result.Error
	}

	applyPairUpdates(phase, result.PairUpdates)
	applyPoolChanges(phase, result.PoolCardsToSet, result.PoolCardToRemove)
	applyHandChanges(&deal.PlayerCards, result.HandChanges)

	if result.NextStatus != "" {
		phase.Status = result.NextStatus
	}

	for _, intent := range result.Events {
		event := ge.buildTributeEvent(intent, deal)
		if event != nil {
			ge.emitEventLocked(event)
		}
	}

	stepResult := &StepResult{
		Action:        result.PendingAction,
		Completed:     result.PhaseCompleted,
		PrevStatus:    prevStatus,
		StatusChanged: result.NextStatus != "" && result.NextStatus != prevStatus,
	}

	if result.PhaseCompleted {
		ge.updatedAt = time.Now()
	}

	return stepResult, nil
}

// StartPlayingPhase 启动出牌阶段
func (ge *GameEngine) StartPlayingPhase() error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	deal, err := ge.requireActiveDeal()
	if err != nil {
		return err
	}

	return deal.StartPlayingPhase()
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

// GetPlayerHand 获取玩家当前手牌（SDK 原生类型）
func (ge *GameEngine) GetPlayerHand(playerSeat int) []*Card {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()

	// 边界检查
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil
	}

	if playerSeat < 0 || playerSeat >= 4 {
		return nil
	}

	// 返回原始引用（零拷贝）
	return ge.currentMatch.CurrentDeal.PlayerCards[playerSeat]
}

// generateID generates a unique ID for the game engine
func generateID() string {
	return fmt.Sprintf("game_%d", time.Now().UnixNano())
}
