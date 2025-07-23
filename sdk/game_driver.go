package sdk

import (
	"context"
	"fmt"
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
	Cards  []*Card    `json:"cards,omitempty"` // 如果是出牌，包含要出的牌
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

// GameDriverConfig 游戏驱动器配置
type GameDriverConfig struct {
	// 超时配置
	PlayDecisionTimeout time.Duration `json:"play_decision_timeout"` // 出牌决策超时时间
	TributeTimeout      time.Duration `json:"tribute_timeout"`       // 贡牌选择超时时间

	// 并发控制
	MaxConcurrentPlayers int `json:"max_concurrent_players"` // 最大并发处理玩家数

	// 事件处理
	AsyncEventHandling bool `json:"async_event_handling"` // 是否异步处理事件
}

// DefaultGameDriverConfig 返回默认的游戏驱动器配置
func DefaultGameDriverConfig() *GameDriverConfig {
	return &GameDriverConfig{
		PlayDecisionTimeout:  30 * time.Second, // 30秒出牌超时
		TributeTimeout:       20 * time.Second, // 20秒贡牌超时
		MaxConcurrentPlayers: 4,                // 最多4个玩家
		AsyncEventHandling:   false,            // 同步事件处理确保顺序
	}
}

// GameDriver 游戏驱动器，负责协调游戏引擎和输入提供者
// 这是新架构的核心组件，将游戏循环逻辑封装在SDK内部
type GameDriver struct {
	engine        GameEngineInterface // 游戏引擎接口
	inputProvider PlayerInputProvider // 玩家输入提供者
	observers     []EventObserver     // 事件观察者列表
	config        *GameDriverConfig   // 驱动器配置
}

// NewGameDriver 创建新的游戏驱动器
func NewGameDriver(engine GameEngineInterface, config *GameDriverConfig) *GameDriver {
	if config == nil {
		config = DefaultGameDriverConfig()
	}

	return &GameDriver{
		engine:    engine,
		observers: make([]EventObserver, 0),
		config:    config,
	}
}

// SetInputProvider 设置输入提供者
func (gd *GameDriver) SetInputProvider(provider PlayerInputProvider) {
	gd.inputProvider = provider
}

// AddObserver 添加事件观察者
func (gd *GameDriver) AddObserver(observer EventObserver) {
	gd.observers = append(gd.observers, observer)
}

// RemoveObserver 移除事件观察者
func (gd *GameDriver) RemoveObserver(observer EventObserver) {
	for i, obs := range gd.observers {
		if obs == observer {
			gd.observers = append(gd.observers[:i], gd.observers[i+1:]...)
			break
		}
	}
}

// notifyObservers 通知所有观察者
func (gd *GameDriver) notifyObservers(event *GameEvent) {
	for _, observer := range gd.observers {
		if gd.config.AsyncEventHandling {
			// 异步处理事件
			go observer.OnGameEvent(event)
		} else {
			// 同步处理事件，确保顺序
			observer.OnGameEvent(event)
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

	// 注册内部事件处理器，用于通知观察者
	gd.engine.RegisterEventHandler(EventMatchStarted, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventDealStarted, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventPlayerPlayed, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventPlayerPassed, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTrickStarted, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTrickEnded, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventDealEnded, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventMatchEnded, gd.handleEngineEvent)

	// 注册上贡相关事件处理器
	gd.engine.RegisterEventHandler(EventTributeRulesSet, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTributeImmunity, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTributePoolCreated, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTributeGiven, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTributeSelected, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventReturnTribute, gd.handleEngineEvent)
	gd.engine.RegisterEventHandler(EventTributeCompleted, gd.handleEngineEvent)

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

		// 处理贡牌动作
		action, err := gd.engine.ProcessTributePhase()
		if err != nil {
			return fmt.Errorf("failed to process tribute phase: %w", err)
		}

		// 如果没有待处理的动作，贡牌阶段完成
		if action == nil {
			break
		}

		// 根据动作类型请求玩家输入
		ctx, cancel := context.WithTimeout(context.Background(), gd.config.TributeTimeout)

		switch action.Type {
		case TributeActionSelect:
			// 双下选牌
			selectedCard, err := gd.inputProvider.RequestTributeSelection(ctx, action.PlayerID, action.Options)
			cancel()
			if err != nil {
				return fmt.Errorf("failed to get tribute selection from player %d: %w", action.PlayerID, err)
			}

			if err := gd.engine.SubmitTributeSelection(action.PlayerID, selectedCard.GetID()); err != nil {
				return fmt.Errorf("failed to submit tribute selection: %w", err)
			}

		case TributeActionReturn:
			// 还贡
			returnCard, err := gd.inputProvider.RequestReturnTribute(ctx, action.PlayerID, action.Options)
			cancel()
			if err != nil {
				return fmt.Errorf("failed to get return tribute from player %d: %w", action.PlayerID, err)
			}

			if err := gd.engine.SubmitReturnTribute(action.PlayerID, returnCard.GetID()); err != nil {
				return fmt.Errorf("failed to submit return tribute: %w", err)
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

		// 获取玩家手牌
		playerView := gd.engine.GetPlayerView(currentPlayer)
		if playerView == nil {
			return fmt.Errorf("failed to get player view for player %d", currentPlayer)
		}

		// 如果玩家没有手牌了，检查游戏是否结束
		if len(playerView.PlayerCards) == 0 {
			// 玩家已出完牌，等待引擎处理状态转换
			break
		}

		// 构造TrickInfo
		trickInfo := &TrickInfo{
			IsLeader: turnInfo.IsLeader,
			LeadComp: turnInfo.LeadComp,
		}

		// 请求玩家决策
		ctx, cancel := context.WithTimeout(context.Background(), gd.config.PlayDecisionTimeout)
		decision, err := gd.inputProvider.RequestPlayDecision(ctx, currentPlayer, playerView.PlayerCards, trickInfo)
		cancel()

		if err != nil {
			return fmt.Errorf("failed to get play decision from player %d: %w", currentPlayer, err)
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
				return fmt.Errorf("failed to play cards for player %d: %w", currentPlayer, err)
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

// handleEngineEvent 处理引擎事件并转发给观察者
func (gd *GameDriver) handleEngineEvent(event *GameEvent) {
	gd.notifyObservers(event)
}
