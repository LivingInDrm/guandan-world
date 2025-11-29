package sdk

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ActionType 定义玩家可执行的动作类型
type ActionType string

const (
	ActionPlay ActionType = "play" // 出牌
	ActionPass ActionType = "pass" // 过牌
)

// PlayDecision 表示玩家的出牌决策
type PlayDecision struct {
	Action ActionType `json:"action"`          // 动作类型：出牌或过牌
	Cards  []*Card    `json:"cards,omitempty"` // 出牌时包含的Card对象
}

// PlayerInputProvider 定义玩家输入提供者接口
// 这个接口将游戏逻辑与具体的输入来源（AI、人工、网络等）解耦
type PlayerInputProvider interface {
	// RequestPlayDecision 请求玩家的出牌决策
	// 参数:
	//   ctx: 上下文，可用于超时控制
	//   playerSeat: 玩家座位号(0-3)
	//   hand: 玩家当前手牌
	//   trickInfo: 当前trick的信息
	// 返回值:
	//   *PlayDecision: 玩家的决策（出牌或过牌）
	//   error: 如果获取决策失败，返回错误
	RequestPlayDecision(ctx context.Context, playerSeat int, hand []*Card, trickInfo *TrickInfo) (*PlayDecision, error)

	// RequestTributeSelection 请求贡牌选择（用于双下情况）
	// 参数:
	//   ctx: 上下文，可用于超时控制
	//   playerSeat: 玩家座位号(0-3)
	//   options: 可选择的牌列表
	// 返回值:
	//   *Card: 选择的牌
	//   error: 如果选择失败，返回错误
	RequestTributeSelection(ctx context.Context, playerSeat int, options []*Card) (*Card, error)

	// RequestReturnTribute 请求还贡选择
	// 参数:
	//   ctx: 上下文，可用于超时控制
	//   playerSeat: 玩家座位号(0-3)
	//   hand: 玩家当前手牌
	// 返回值:
	//   *Card: 选择还贡的牌
	//   error: 如果选择失败，返回错误
	RequestReturnTribute(ctx context.Context, playerSeat int, hand []*Card) (*Card, error)
}

// GameDriverConfig 游戏驱动器配置
type GameDriverConfig struct {
	// 超时配置
	PlayDecisionTimeout time.Duration `json:"play_decision_timeout"` // 出牌决策超时时间
	TributeTimeout      time.Duration `json:"tribute_timeout"`       // 贡牌选择超时时间
	TurnTimeoutSeconds  int           `json:"turn_timeout_seconds"`  // 单次出牌超时秒数

	// 超时策略
	TimeoutStrategy TimeoutStrategy `json:"-"` // 超时默认决策策略（不序列化）

	// 并发控制
	MaxConcurrentPlayers int `json:"max_concurrent_players"` // 最大并发处理玩家数

	// 事件处理
	AsyncEventHandling bool `json:"async_event_handling"` // 是否异步处理事件

	// 上贡阶段延时配置（用于前端动画展示）
	TributeStartedDelay       time.Duration `json:"tribute_started_delay"`        // TributeStarted 事件后延时
	TributeCardSubmittedDelay time.Duration `json:"tribute_card_submitted_delay"` // TributeCardSubmitted 事件后延时
	TributeCardSelectedDelay  time.Duration `json:"tribute_card_selected_delay"`  // TributeCardSelected 事件后延时
	TributeCardReturnedDelay  time.Duration `json:"tribute_card_returned_delay"`  // TributeCardReturned 事件后延时
	TributeCompletedDelay     time.Duration `json:"tribute_completed_delay"`      // TributeCompleted 事件后延时
}

// DefaultGameDriverConfig 返回默认的游戏驱动器配置
func DefaultGameDriverConfig() *GameDriverConfig {
	return &GameDriverConfig{
		PlayDecisionTimeout:       30 * time.Second,            // 30秒出牌超时
		TributeTimeout:            20 * time.Second,            // 20秒贡牌超时
		TurnTimeoutSeconds:        20,                          // 20秒单次出牌超时
		TimeoutStrategy:           NewDefaultTimeoutStrategy(), // 使用默认超时策略
		MaxConcurrentPlayers:      4,                           // 最多4个玩家
		AsyncEventHandling:        false,                       // 同步事件处理确保顺序
		TributeStartedDelay:       2000 * time.Millisecond,     // 上贡开始后延时
		TributeCardSubmittedDelay: 2000 * time.Millisecond,     // 贡牌提交后延时
		TributeCardSelectedDelay:  2000 * time.Millisecond,     // 选牌后延时
		TributeCardReturnedDelay:  2000 * time.Millisecond,     // 还贡后延时
		TributeCompletedDelay:     2000 * time.Millisecond,     // 上贡完成后延时
	}
}

// GameDriver 游戏驱动器，负责协调游戏引擎和输入提供者
// 这是新架构的核心组件，将游戏循环逻辑封装在SDK内部
type GameDriver struct {
	engine        GameEngineInterface // 游戏引擎接口
	inputProvider PlayerInputProvider // 玩家输入提供者
	observers     []EventObserver     // 事件观察者列表
	observersMu   sync.RWMutex        // 保护观察者列表的读写锁
	config        *GameDriverConfig   // 驱动器配置

	// 事件系统
	registeredWithEngine bool       // 标记是否已向引擎注册为观察者（防止重复注册）
	registrationMu       sync.Mutex // 保护注册状态的互斥锁

	// 超时管理
	timeoutStats  map[int]*PlayerTimeoutStats // 超时统计 (座位号 -> 统计信息)
	timeoutMu     sync.RWMutex                // 保护超时统计的读写锁
	gameCancelCtx context.Context             // 游戏取消信号上下文（将在RunMatch中初始化）
	cancelFunc    context.CancelFunc          // 取消函数（将在RunMatch中初始化）
	cancelMu      sync.Mutex                  // 保护cancelFunc的互斥锁
}

// NewGameDriver 创建新的游戏驱动器
func NewGameDriver(engine GameEngineInterface, config *GameDriverConfig) *GameDriver {
	if config == nil {
		config = DefaultGameDriverConfig()
	}

	// 如果配置中没有设置超时策略，使用默认策略
	if config.TimeoutStrategy == nil {
		config.TimeoutStrategy = NewDefaultTimeoutStrategy()
	}

	return &GameDriver{
		engine:       engine,
		observers:    make([]EventObserver, 0),
		config:       config,
		timeoutStats: make(map[int]*PlayerTimeoutStats),
	}
}

// SetInputProvider 设置输入提供者
func (gd *GameDriver) SetInputProvider(provider PlayerInputProvider) {
	gd.inputProvider = provider
}

// AddObserver 添加事件观察者
// 如果观察者已存在，不会重复添加
func (gd *GameDriver) AddObserver(observer EventObserver) {
	gd.observersMu.Lock()
	defer gd.observersMu.Unlock()

	// 检查是否已存在，避免重复添加
	for _, obs := range gd.observers {
		if obs == observer {
			return
		}
	}

	gd.observers = append(gd.observers, observer)
}

// RemoveObserver 移除事件观察者
func (gd *GameDriver) RemoveObserver(observer EventObserver) {
	gd.observersMu.Lock()
	defer gd.observersMu.Unlock()
	for i, obs := range gd.observers {
		if obs == observer {
			gd.observers = append(gd.observers[:i], gd.observers[i+1:]...)
			break
		}
	}
}

// OnGameEvent 实现 EventObserver 接口
// GameDriver 作为 EventObserver，接收来自 GameEngine 的事件并转发给自己的观察者
func (gd *GameDriver) OnGameEvent(event *GameEvent) {
	gd.notifyObservers(event)
}

// notifyObservers 通知所有观察者
func (gd *GameDriver) notifyObservers(event *GameEvent) {
	// Safely read and copy the observer list to avoid holding the lock during callback execution
	gd.observersMu.RLock()
	observersCopy := make([]EventObserver, len(gd.observers))
	copy(observersCopy, gd.observers)
	gd.observersMu.RUnlock()

	for _, observer := range observersCopy {
		if gd.config.AsyncEventHandling {
			// 异步处理事件，使用 goroutine
			obs := observer
			evt := event
			go func() {
				defer func() {
					if r := recover(); r != nil {
						// 捕获 panic，防止崩溃
						fmt.Printf("[GameDriver] Observer panic for %s: %v\n", evt.Type, r)
					}
				}()
				obs.OnGameEvent(evt)
			}()
		} else {
			// 同步处理事件，确保顺序
			// 即使是同步模式也要防止 panic 影响后续观察者
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("[GameDriver] Observer panic for %s: %v\n", event.Type, r)
					}
				}()
				observer.OnGameEvent(event)
			}()
		}
	}
}

// GetEngine 获取游戏引擎（只读访问）
func (gd *GameDriver) GetEngine() GameEngineInterface {
	return gd.engine
}

// GetConfig 获取驱动器配置
func (gd *GameDriver) GetConfig() *GameDriverConfig {
	return gd.config
}

// GetTimeoutStats 获取超时统计信息（返回深拷贝，避免外部修改）
func (gd *GameDriver) GetTimeoutStats() map[int]*PlayerTimeoutStats {
	gd.timeoutMu.RLock()
	defer gd.timeoutMu.RUnlock()

	// 返回深拷贝以防止外部修改内部状态
	result := make(map[int]*PlayerTimeoutStats, len(gd.timeoutStats))
	for seat, stats := range gd.timeoutStats {
		if stats != nil {
			// 复制值而不是共享指针
			statsCopy := *stats
			result[seat] = &statsCopy
		}
	}
	return result
}

// CancelMatch 取消当前正在运行的比赛（线程安全）
func (gd *GameDriver) CancelMatch() {
	gd.cancelMu.Lock()
	defer gd.cancelMu.Unlock()

	if gd.cancelFunc != nil {
		gd.cancelFunc()
	}
}

// incrementPlayDecisionTimeout 增加出牌决策超时次数（内部方法，线程安全）
func (gd *GameDriver) incrementPlayDecisionTimeout(seat int) {
	gd.timeoutMu.Lock()
	defer gd.timeoutMu.Unlock()

	if gd.timeoutStats[seat] == nil {
		gd.timeoutStats[seat] = &PlayerTimeoutStats{}
	}
	gd.timeoutStats[seat].PlayDecisionTimeouts++
	gd.timeoutStats[seat].TotalTimeouts++
}

// incrementTributeTimeout 增加贡牌超时次数（内部方法，线程安全）
func (gd *GameDriver) incrementTributeTimeout(seat int) {
	gd.timeoutMu.Lock()
	defer gd.timeoutMu.Unlock()

	if gd.timeoutStats[seat] == nil {
		gd.timeoutStats[seat] = &PlayerTimeoutStats{}
	}
	gd.timeoutStats[seat].TributeTimeouts++
	gd.timeoutStats[seat].TotalTimeouts++
}

// GameDriverResult 游戏驱动器返回的扩展结果
type GameDriverResult struct {
	*MatchResult                      // 嵌入现有的MatchResult
	DealCount    int                  `json:"deal_count"`   // 总局数
	PlayerStats  map[int]*PlayerStats `json:"player_stats"` // 玩家统计
}

// PlayerStats 玩家统计信息
type PlayerStats struct {
	CardsPlayed int           `json:"cards_played"` // 出牌次数
	TricksWon   int           `json:"tricks_won"`   // 赢得的trick数
	AverageTime time.Duration `json:"average_time"` // 平均出牌时间
}

// RunMatch 运行完整比赛
// 这是新架构的核心方法，将整个比赛循环封装在SDK内部
func (gd *GameDriver) RunMatch(players []Player) (*GameDriverResult, error) {
	if gd.inputProvider == nil {
		return nil, fmt.Errorf("input provider not set")
	}

	// 使用互斥锁保护cancelFunc的读写
	gd.cancelMu.Lock()
	// 如果存在之前的取消函数，先调用它防止资源泄漏
	if gd.cancelFunc != nil {
		gd.cancelFunc()
	}

	// 创建游戏取消上下文，用于游戏结束时取消所有待处理的操作
	gd.gameCancelCtx, gd.cancelFunc = context.WithCancel(context.Background())
	gd.cancelMu.Unlock()

	defer func() {
		gd.cancelMu.Lock()
		if gd.cancelFunc != nil {
			gd.cancelFunc() // 确保游戏结束时取消所有操作
		}
		gd.cancelMu.Unlock()
	}()

	// 注册所有引擎事件，让 GameDriver 作为 EventObserver 统一接收
	// GameDriver 会将这些事件转发给自己的观察者
	// 使用一次性注册机制避免重复注册造成的内存泄漏和重复通知
	gd.registrationMu.Lock()
	if !gd.registeredWithEngine {
		allEventTypes := []GameEventType{
			EventMatchStarted,
			EventDealStarted,
			EventCardsDealt,     // 添加发牌完成事件（未来可能使用）
			EventTributeStarted, // 添加贡牌开始事件（未来可能使用）
			EventPlayerPlayed,
			EventPlayerPassed,
			EventTrickStarted,
			EventTrickEnded,
			EventDealEnded,
			EventMatchEnded,
			EventTributeRulesSet,
			EventTributeImmunity,
			EventTributePoolCreated,
			EventTributeGiven,
			EventTributeSelected,
			EventReturnTribute,
			EventTributeCompleted,
			EventPlayerDisconnect,
			EventPlayerReconnect,
		}
		for _, eventType := range allEventTypes {
			gd.engine.RegisterObserver(eventType, gd)
		}
		gd.registeredWithEngine = true
	}
	gd.registrationMu.Unlock()

	// 开始比赛
	if err := gd.engine.StartMatch(players); err != nil {
		return nil, fmt.Errorf("failed to start match: %w", err)
	}

	startTime := time.Now()
	dealCount := 0

	// 比赛主循环
	for !gd.engine.IsGameFinished() {
		dealCount++

		if err := gd.runDeal(); err != nil {
			return nil, fmt.Errorf("failed to run deal %d: %w", dealCount, err)
		}
	}

	// 构建比赛结果
	gameState := gd.engine.GetGameState()
	matchDetails := gd.engine.GetMatchDetails()

	// 创建基础MatchResult
	baseResult := &MatchResult{
		Winner:   -1,
		Duration: time.Since(startTime),
	}

	if gameState.CurrentMatch != nil {
		baseResult.Winner = gameState.CurrentMatch.Winner
	}

	if matchDetails != nil {
		baseResult.FinalLevels = matchDetails.TeamLevels

		// 创建统计信息
		baseResult.Statistics = &MatchStatistics{
			TotalDeals:    dealCount,
			TotalDuration: time.Since(startTime),
			FinalLevels:   matchDetails.TeamLevels,
			TeamStats: [2]*TeamMatchStats{
				{Team: 0, DealsWon: 0, TotalTricks: 0, Upgrades: 0},
				{Team: 1, DealsWon: 0, TotalTricks: 0, Upgrades: 0},
			},
		}
	}

	// 创建扩展结果
	result := &GameDriverResult{
		MatchResult: baseResult,
		DealCount:   dealCount,
		PlayerStats: make(map[int]*PlayerStats),
	}

	// 初始化玩家统计
	for i := 0; i < 4; i++ {
		result.PlayerStats[i] = &PlayerStats{
			CardsPlayed: 0,
			TricksWon:   0,
			AverageTime: 0,
		}
	}

	return result, nil
}

// runDeal 运行一局牌
func (gd *GameDriver) runDeal() error {
	// 开始新局
	if err := gd.engine.StartDeal(); err != nil {
		return fmt.Errorf("failed to start deal: %w", err)
	}

	// 处理贡牌阶段
	if gd.engine.GetCurrentDealStatus() == DealStatusTribute {
		// TributeStarted 事件已发送，延时让前端展示动画
		if !gd.sleepWithContext(gd.config.TributeStartedDelay) {
			return fmt.Errorf("game cancelled during tribute started delay")
		}

		if err := gd.runTributePhase(); err != nil {
			return fmt.Errorf("failed to run tribute phase: %w", err)
		}
	}

	// 处理游戏阶段
	if gd.engine.GetCurrentDealStatus() == DealStatusPlaying {
		if err := gd.runPlayingPhase(); err != nil {
			return fmt.Errorf("failed to run playing phase: %w", err)
		}
	}

	return nil
}

// runTributePhase 运行贡牌阶段
func (gd *GameDriver) runTributePhase() error {
	maxActions := 20 // 安全计数器，防止无限循环
	actionCount := 0

	for gd.engine.GetCurrentDealStatus() == DealStatusTribute && actionCount < maxActions {
		actionCount++

		// 记录处理前的状态
		prevStatus := gd.getTributeStatus()

		// 处理贡牌动作
		// ProcessTributePhase 会在一次调用中完成所有自动状态转换
		// 返回 action 表示需要用户输入，返回 nil 表示贡牌阶段完成
		action, err := gd.engine.ProcessTributePhase()
		if err != nil {
			return fmt.Errorf("failed to process tribute phase: %w", err)
		}

		currStatus := gd.getTributeStatus()

		// 添加延时让前端展示动画
		// 注意：ProcessTributePhase 可能在一次调用中跨越多个状态
		// 例如单下/末游：Waiting → Selecting → Returning

		// 如果从 Waiting 开始有状态变化，说明发送了 TributeCardSubmitted 事件
		if prevStatus == TributeStatusWaiting && currStatus != TributeStatusWaiting {
			if !gd.sleepWithContext(gd.config.TributeCardSubmittedDelay) {
				return fmt.Errorf("game cancelled during tribute card submitted delay")
			}
		}

		// 如果当前是 Returning 且之前不是，说明发送了 TributeCardSelected 事件（单下/末游自动选贡）
		if prevStatus != TributeStatusReturning && currStatus == TributeStatusReturning {
			if !gd.sleepWithContext(gd.config.TributeCardSelectedDelay) {
				return fmt.Errorf("game cancelled during tribute card selected delay")
			}
		}

		// 贡牌阶段完成（ProcessTributePhase 返回 nil）
		if action == nil {
			// TributeCompleted 事件已发送，添加延时
			if !gd.sleepWithContext(gd.config.TributeCompletedDelay) {
				return fmt.Errorf("game cancelled during tribute completed delay")
			}
			break
		}

		// 根据动作类型请求玩家输入
		// 使用gameCancelCtx作为基础，这样游戏结束时会自动取消所有请求
		ctx, cancel := context.WithTimeout(gd.gameCancelCtx, gd.config.TributeTimeout)

		switch action.Type {
		case TributeActionSelect:
			// 双下选牌
			selectedCard, err := gd.inputProvider.RequestTributeSelection(ctx, action.PlayerID, action.Options)
			ctxErr := ctx.Err() // 在cancel()之前捕获上下文错误
			cancel()

			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded) || ctxErr == context.DeadlineExceeded:
					// 超时，使用默认策略
					gd.handleTimeout(action.PlayerID, "tribute_select")
					selectedCard = gd.config.TimeoutStrategy.GetDefaultTributeCard(action.Options)
					if selectedCard == nil {
						// Fallback：选择第一个非nil选项
						for _, c := range action.Options {
							if c != nil {
								selectedCard = c
								break
							}
						}
						if selectedCard == nil {
							return fmt.Errorf("timeout strategy returned nil and no valid options, player %d", action.PlayerID)
						}
					}
				case errors.Is(err, context.Canceled) || ctxErr == context.Canceled:
					// 游戏被取消
					return fmt.Errorf("game cancelled during tribute selection for player %d", action.PlayerID)
				default:
					return fmt.Errorf("failed to get tribute selection from player %d: %w", action.PlayerID, err)
				}
			}

			if err := gd.engine.SubmitTributeSelection(action.PlayerID, selectedCard); err != nil {
				return fmt.Errorf("failed to submit tribute selection: %w", err)
			}

			// TributeCardSelected 事件已发送，延时让前端展示动画
			if !gd.sleepWithContext(gd.config.TributeCardSelectedDelay) {
				return fmt.Errorf("game cancelled during tribute card selected delay")
			}

		case TributeActionReturn:
			// 还贡
			returnCard, err := gd.inputProvider.RequestReturnTribute(ctx, action.PlayerID, action.Options)
			ctxErr := ctx.Err() // 在cancel()之前捕获上下文错误
			cancel()

			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded) || ctxErr == context.DeadlineExceeded:
					// 超时，使用默认策略
					gd.handleTimeout(action.PlayerID, "return_tribute")
					returnCard = gd.config.TimeoutStrategy.GetDefaultReturnCard(action.Options)
					if returnCard == nil {
						// Fallback：选择第一个非nil选项
						for _, c := range action.Options {
							if c != nil {
								returnCard = c
								break
							}
						}
						if returnCard == nil {
							return fmt.Errorf("timeout strategy returned nil and no valid return card, player %d", action.PlayerID)
						}
					}
				case errors.Is(err, context.Canceled) || ctxErr == context.Canceled:
					// 游戏被取消
					return fmt.Errorf("game cancelled during return tribute for player %d", action.PlayerID)
				default:
					return fmt.Errorf("failed to get return tribute from player %d: %w", action.PlayerID, err)
				}
			}

			if err := gd.engine.SubmitReturnTribute(action.PlayerID, returnCard); err != nil {
				return fmt.Errorf("failed to submit return tribute: %w", err)
			}

			// TributeCardReturned 事件已发送，延时让前端展示动画
			if !gd.sleepWithContext(gd.config.TributeCardReturnedDelay) {
				return fmt.Errorf("game cancelled during tribute card returned delay")
			}
		default:
			cancel()
			return fmt.Errorf("unknown tribute action type: %v", action.Type)
		}
	}

	if actionCount >= maxActions {
		return fmt.Errorf("tribute phase exceeded maximum actions limit")
	}

	return nil
}

// runPlayingPhase 运行游戏阶段
func (gd *GameDriver) runPlayingPhase() error {
	maxTricks := 200 // 安全计数器
	trickCount := 0

	for gd.engine.GetCurrentDealStatus() == DealStatusPlaying && trickCount < maxTricks {
		trickCount++

		if err := gd.runTrick(); err != nil {
			return fmt.Errorf("failed to run trick %d: %w", trickCount, err)
		}
	}

	if trickCount >= maxTricks {
		return fmt.Errorf("playing phase exceeded maximum tricks limit")
	}

	return nil
}

// runTrick 运行单个trick
func (gd *GameDriver) runTrick() error {
	maxTurns := 50 // 安全计数器，考虑到复杂情况下可能需要更多轮
	turnCount := 0

	// 在开始trick前检查状态
	initialStatus := gd.engine.GetCurrentDealStatus()
	if initialStatus != DealStatusPlaying {
		return nil // Deal已经不在playing状态
	}

	for turnCount < maxTurns {
		turnCount++

		// 每轮开始前重新检查deal状态
		dealStatus := gd.engine.GetCurrentDealStatus()
		if dealStatus != DealStatusPlaying {
			// Deal已结束
			break
		}

		// 获取当前轮次信息
		turnInfo := gd.engine.GetCurrentTurnInfo()
		if turnInfo == nil || !turnInfo.HasActiveTrick {
			// 当前trick结束
			break
		}

		// 检查是否开始了新的trick
		if turnInfo.IsNewTrick && turnCount > 1 {
			// 新trick已经开始，当前trick结束
			break
		}

		currentPlayer := turnInfo.CurrentPlayer
		if currentPlayer < 0 || currentPlayer > 3 {
			// 无效的玩家ID，trick可能已结束
			break
		}

		// 获取玩家手牌（SDK 原生类型，包含完整 Level）
		hand := gd.engine.GetPlayerHand(currentPlayer)
		if hand == nil {
			return fmt.Errorf("failed to get player hand for player %d", currentPlayer)
		}

		// 如果玩家没有手牌了，检查游戏是否结束
		if len(hand) == 0 {
			// 玩家已出完牌，等待引擎处理状态转换
			break
		}

		// 构造TrickInfo
		trickInfo := &TrickInfo{
			IsLeader: turnInfo.IsLeader,
			LeadComp: turnInfo.LeadComp,
		}

		// 请求玩家决策（带超时检测）
		// 使用gameCancelCtx作为基础，这样游戏结束时会自动取消所有请求
		ctx, cancel := context.WithTimeout(gd.gameCancelCtx, gd.config.PlayDecisionTimeout)
		decision, err := gd.inputProvider.RequestPlayDecision(ctx, currentPlayer, hand, trickInfo)
		ctxErr := ctx.Err() // 在cancel()之前捕获上下文错误
		cancel()

		// 处理超时情况
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded) || ctxErr == context.DeadlineExceeded:
				// 超时，使用默认策略生成决策
				gd.handleTimeout(currentPlayer, "play_decision")
				decision = gd.config.TimeoutStrategy.GetDefaultPlayDecision(hand, trickInfo)
				if decision == nil {
					// 如果策略返回nil，使用PASS作为后备
					decision = &PlayDecision{Action: ActionPass}
				}
			case errors.Is(err, context.Canceled) || ctxErr == context.Canceled:
				// 游戏被取消（例如：游戏结束或用户中止）
				return fmt.Errorf("game cancelled during play decision for player %d", currentPlayer)
			default:
				// 其他错误
				return fmt.Errorf("failed to get play decision from player %d: %w", currentPlayer, err)
			}
		}

		// 如果decision仍然为nil（可能是provider返回了nil），使用默认策略
		if decision == nil {
			decision = gd.config.TimeoutStrategy.GetDefaultPlayDecision(hand, trickInfo)
			if decision == nil {
				// 如果策略也返回nil，使用PASS作为后备
				decision = &PlayDecision{Action: ActionPass}
			}
		}

		// 执行决策前再次检查状态
		if gd.engine.GetCurrentDealStatus() != DealStatusPlaying {
			break
		}

		// 执行决策
		if decision.Action == ActionPlay {
			if decision.Cards == nil || len(decision.Cards) == 0 {
				return fmt.Errorf("player %d chose to play but provided no cards", currentPlayer)
			}

			_, err = gd.engine.PlayCards(currentPlayer, decision.Cards)
			if err != nil {
				// 检测到非法出牌
				// 记录为超时统计（统一处理非法行为）
				gd.handleTimeout(currentPlayer, "invalid_play")

				// 使用超时策略自动生成合法决策
				autoDecision := gd.config.TimeoutStrategy.GetDefaultPlayDecision(hand, trickInfo)
				if autoDecision == nil {
					// 策略返回nil，使用PASS作为后备
					autoDecision = &PlayDecision{Action: ActionPass}
				}

				// 执行自动决策
				if autoDecision.Action == ActionPlay {
					_, autoErr := gd.engine.PlayCards(currentPlayer, autoDecision.Cards)
					if autoErr != nil {
						// 自动策略也失败了，这是严重系统错误
						return fmt.Errorf("auto-play failed for player %d after invalid play: %w", currentPlayer, autoErr)
					}
				} else {
					_, autoErr := gd.engine.PassTurn(currentPlayer)
					if autoErr != nil {
						return fmt.Errorf("auto-pass failed for player %d after invalid play: %w", currentPlayer, autoErr)
					}
				}

				// 继续游戏循环，不退出
				continue
			}
		} else if decision.Action == ActionPass {
			_, err = gd.engine.PassTurn(currentPlayer)
			if err != nil {
				return fmt.Errorf("failed to pass turn for player %d: %w", currentPlayer, err)
			}
		} else {
			return fmt.Errorf("invalid action type from player %d: %v", currentPlayer, decision.Action)
		}

		// 执行后检查是否有状态变化
		// 注意：不要过于严格地检查状态变化，因为某些情况下（如trick结束）状态可能不会立即变化
		newTurnInfo := gd.engine.GetCurrentTurnInfo()
		if newTurnInfo != nil && newTurnInfo.HasActiveTrick && newTurnInfo.CurrentPlayer == currentPlayer {
			// 轮次没有变化，可能trick已结束或有其他状态转换
			// 我们不应该在这里报错，让外层循环来处理
		}
	}

	if turnCount >= maxTurns {
		return fmt.Errorf("trick exceeded maximum turns limit (%d turns)", maxTurns)
	}

	return nil
}

// handleTimeout 处理玩家超时
// actionType: "play_decision", "tribute_select", "return_tribute"
func (gd *GameDriver) handleTimeout(playerSeat int, actionType string) {
	// 记录超时统计
	switch actionType {
	case "play_decision":
		gd.incrementPlayDecisionTimeout(playerSeat)
	case "tribute_select", "return_tribute":
		gd.incrementTributeTimeout(playerSeat)
	}

	// 获取当前游戏状态
	var match *Match
	var deal *Deal
	var trick *Trick

	gameState := gd.engine.GetGameState()
	if gameState != nil && gameState.CurrentMatch != nil {
		match = gameState.CurrentMatch
		if match.CurrentDeal != nil {
			deal = match.CurrentDeal
			if deal.CurrentTrick != nil {
				trick = deal.CurrentTrick
			}
		}
	}

	// 使用engine的EventMetadataProvider发出超时事件
	eventMeta := gd.engine.GetEventMetadataProvider()
	event := NewPlayerTimeoutEvent(eventMeta, match, deal, trick, playerSeat, actionType)
	gd.notifyObservers(event)
}

// sleepWithContext 在指定时间内等待，支持 context 取消
// 返回 true 表示正常完成延时，false 表示被 context 取消
func (gd *GameDriver) sleepWithContext(d time.Duration) bool {
	if d <= 0 {
		return true
	}

	select {
	case <-time.After(d):
		return true
	case <-gd.gameCancelCtx.Done():
		return false
	}
}

// getTributeStatus 获取当前上贡阶段的状态
func (gd *GameDriver) getTributeStatus() TributeStatus {
	gameState := gd.engine.GetGameState()
	if gameState == nil || gameState.CurrentMatch == nil {
		return ""
	}
	if gameState.CurrentMatch.CurrentDeal == nil {
		return ""
	}
	if gameState.CurrentMatch.CurrentDeal.TributePhase == nil {
		return ""
	}
	return gameState.CurrentMatch.CurrentDeal.TributePhase.Status
}
