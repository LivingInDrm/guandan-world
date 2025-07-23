package sdk

import (
	"fmt"
	"strings"
)

// MatchSimulatorObserver 比赛模拟器事件观察者
// 负责观察游戏事件并输出日志信息
type MatchSimulatorObserver struct {
	engine  GameEngineInterface // 引擎引用，用于查询状态
	verbose bool                // 是否详细输出
	logger  func(string)        // 日志输出函数
}

// NewMatchSimulatorObserver 创建新的观察者
func NewMatchSimulatorObserver(engine GameEngineInterface, verbose bool, logger func(string)) *MatchSimulatorObserver {
	if logger == nil {
		// 默认日志函数
		logger = func(message string) {
			if verbose {
				fmt.Println(message)
			}
		}
	}

	return &MatchSimulatorObserver{
		engine:  engine,
		verbose: verbose,
		logger:  logger,
	}
}

// OnGameEvent 实现EventObserver接口
func (mso *MatchSimulatorObserver) OnGameEvent(event *GameEvent) {
	switch event.Type {
	case EventMatchStarted:
		mso.handleMatchStarted(event)
	case EventDealStarted:
		mso.handleDealStarted(event)
	case EventCardsDealt:
		mso.handleCardsDealt(event)
	case EventTributePhase:
		mso.handleTributePhase(event)
	case EventTributeImmunity:
		mso.handleTributeImmunity(event)
	case EventTributeStarted:
		mso.handleTributeStarted(event)
	case EventTributeGiven:
		mso.handleTributeGiven(event)
	case EventTributeSelected:
		mso.handleTributeSelected(event)
	case EventReturnTribute:
		mso.handleReturnTribute(event)
	case EventTributeCompleted:
		mso.handleTributeCompleted(event)
	case EventTrickStarted:
		mso.handleTrickStarted(event)
	case EventPlayerPlayed:
		mso.handlePlayerPlayed(event)
	case EventPlayerPassed:
		mso.handlePlayerPassed(event)
	case EventTrickEnded:
		mso.handleTrickEnded(event)
	case EventDealEnded:
		mso.handleDealEnded(event)
	case EventMatchEnded:
		mso.handleMatchEnded(event)
	default:
		// 忽略未知事件类型
	}
}

// 事件处理方法（从原MatchSimulator移植过来）

func (mso *MatchSimulatorObserver) handleMatchStarted(event *GameEvent) {
	mso.log("Event: Match Started")
}

func (mso *MatchSimulatorObserver) handleDealStarted(event *GameEvent) {
	mso.log("Event: Deal Started")
}

func (mso *MatchSimulatorObserver) handleCardsDealt(event *GameEvent) {
	mso.log("Event: Cards Dealt")
}

func (mso *MatchSimulatorObserver) handleTributePhase(event *GameEvent) {
	mso.log("Event: Tribute Phase")
}

func (mso *MatchSimulatorObserver) handleTributeImmunity(event *GameEvent) {
	mso.log("Event: Tribute Immunity triggered - No tribute required this deal")
}

func (mso *MatchSimulatorObserver) handleTributeStarted(event *GameEvent) {
	mso.log("Event: Tribute Started - Tribute phase begins")
}

func (mso *MatchSimulatorObserver) handleTributeGiven(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if giver, ok := data["giver"].(int); ok {
			if receiver, ok := data["receiver"].(int); ok {
				if card, ok := data["card"].(*Card); ok {
					mso.log(fmt.Sprintf("Event: Tribute Given - Player %d gives %s to Player %d",
						giver, card.ToShortString(), receiver))
				}
			}
		}
	}
}

func (mso *MatchSimulatorObserver) handleTributeSelected(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if player, ok := data["player"].(int); ok {
			if cardID, ok := data["cardID"].(string); ok {
				mso.log(fmt.Sprintf("Event: Tribute Selected - Player %d selected card %s (Double-down selection)",
					player, cardID))
			}
		}
	}
}

func (mso *MatchSimulatorObserver) handleReturnTribute(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if returner, ok := data["player"].(int); ok {
			if cardID, ok := data["cardID"].(string); ok {
				mso.log(fmt.Sprintf("Event: Return Tribute - Player %d returns card %s",
					returner, cardID))
			}
		}
	}
}

func (mso *MatchSimulatorObserver) handleTributeCompleted(event *GameEvent) {
	mso.log("Event: Tribute Completed")
}

func (mso *MatchSimulatorObserver) handleTrickStarted(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if leader, ok := data["leader"].(int); ok {
			mso.log(fmt.Sprintf("Event: New Trick Started, Leader: Player %d", leader))

			// 注意：不能在事件处理器中调用引擎方法，会导致死锁
			// 手牌信息应该通过其他方式获取或在事件数据中提供
			if mso.verbose {
				mso.log(fmt.Sprintf("New Trick Started (Leader: Player %d) - Hand details omitted to avoid deadlock", leader))
			}
		}
	}
}

func (mso *MatchSimulatorObserver) handlePlayerPlayed(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		cards := data["cards"].([]*Card)

		// 将出牌转换为简化格式
		var cardStrs []string
		for _, card := range cards {
			cardStrs = append(cardStrs, card.ToShortString())
		}

		mso.log(fmt.Sprintf("Event: Player %d played %d cards: [%s]",
			playerSeat, len(cards), strings.Join(cardStrs, ",")))
	}
}

func (mso *MatchSimulatorObserver) handlePlayerPassed(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		mso.log(fmt.Sprintf("Event: Player %d passed", playerSeat))
	}
}

func (mso *MatchSimulatorObserver) handleTrickEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if winner, ok := data["winner"].(int); ok {
			mso.log(fmt.Sprintf("Event: Trick Ended, Winner: Player %d", winner))
		}
	}
}

func (mso *MatchSimulatorObserver) handleDealEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if result, ok := data["result"].(*DealResult); ok {
			mso.log(fmt.Sprintf("Event: Deal Ended, Rankings: %v, Victory Type: %v",
				result.Rankings, result.VictoryType))
		}
	}
}

func (mso *MatchSimulatorObserver) handleMatchEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if winner, ok := data["winner"].(int); ok {
			mso.log(fmt.Sprintf("Event: Match Ended, Winner: Team %d", winner))
		}
	}
}

// 辅助方法

func (mso *MatchSimulatorObserver) log(message string) {
	mso.logger(message)
}

// logPlayerHands 输出所有玩家的手牌
// 注意：此方法可能导致死锁，仅在确定没有锁竞争时使用
func (mso *MatchSimulatorObserver) logPlayerHands(context string) {
	if !mso.verbose {
		return
	}

	// 为了避免死锁，这里只记录上下文信息
	// 实际的手牌信息需要在没有锁竞争的时候获取
	mso.log(fmt.Sprintf("%s - Player hands details omitted to avoid deadlock", context))
}

// logTeamStatus 输出队伍状态
func (mso *MatchSimulatorObserver) logTeamStatus() {
	matchDetails := mso.engine.GetMatchDetails()
	if matchDetails == nil {
		return
	}

	mso.log("=== Team Status Before Next Deal ===")
	mso.log(fmt.Sprintf("Team 0 (Players 0,2): Level %d", matchDetails.TeamLevels[0]))
	mso.log(fmt.Sprintf("Team 1 (Players 1,3): Level %d", matchDetails.TeamLevels[1]))
	mso.log("Players:")

	for i, player := range matchDetails.Players {
		teamNum := (i % 2)
		mso.log(fmt.Sprintf("  Player %d (%s) - Team %d", i, player.Username, teamNum))
	}
}

// SetVerbose 设置详细输出模式
func (mso *MatchSimulatorObserver) SetVerbose(verbose bool) {
	mso.verbose = verbose
}

// SetLogger 设置日志输出函数
func (mso *MatchSimulatorObserver) SetLogger(logger func(string)) {
	if logger != nil {
		mso.logger = logger
	}
}
