package simulator

import (
	"fmt"
	"strings"

	"guandan-world/sdk"
)

// MatchSimulatorObserver 比赛模拟器事件观察者
// 负责观察游戏事件并输出日志信息
type MatchSimulatorObserver struct {
	engine  sdk.GameEngineInterface // 引擎引用，用于查询状态
	verbose bool                    // 是否详细输出
	logger  func(string)            // 日志输出函数
}

// NewMatchSimulatorObserver 创建新的观察者
func NewMatchSimulatorObserver(engine sdk.GameEngineInterface, verbose bool, logger func(string)) *MatchSimulatorObserver {
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
func (mso *MatchSimulatorObserver) OnGameEvent(event *sdk.GameEvent) {
	switch event.Type {
	case sdk.EventMatchStarted:
		mso.handleMatchStarted(event)
	case sdk.EventDealStarted:
		mso.handleDealStarted(event)
	case sdk.EventCardsDealt:
		mso.handleCardsDealt(event)
	case sdk.EventTributeStarted:
		mso.handleTributeStarted(event)
	case sdk.EventTributeExempted:
		mso.handleTributeExempted(event)
	case sdk.EventTributeCardSubmitted:
		mso.handleTributeCardSubmitted(event)
	case sdk.EventTributeCardSelected:
		mso.handleTributeCardSelected(event)
	case sdk.EventReturnTribute:
		mso.handleReturnTribute(event)
	case sdk.EventTributeCompleted:
		mso.handleTributeCompleted(event)
	case sdk.EventTrickStarted:
		mso.handleTrickStarted(event)
	case sdk.EventPlayerPlayed:
		mso.handlePlayerPlayed(event)
	case sdk.EventPlayerPassed:
		mso.handlePlayerPassed(event)
	case sdk.EventTrickEnded:
		mso.handleTrickEnded(event)
	case sdk.EventDealEnded:
		mso.handleDealEnded(event)
	case sdk.EventMatchEnded:
		mso.handleMatchEnded(event)
	default:
		// 忽略未知事件类型
	}
}

// 事件处理方法（从原MatchSimulator移植过来）

func (mso *MatchSimulatorObserver) handleMatchStarted(event *sdk.GameEvent) {
	mso.log("Event: Match Started")
}

func (mso *MatchSimulatorObserver) handleDealStarted(event *sdk.GameEvent) {
	mso.log("Event: Deal Started")

	payload := event.GetDealStarted()
	if payload == nil {
		return
	}

	dealLevel := int(payload.DealLevel)
	team0Level := int(payload.TeamLevels[0])
	team1Level := int(payload.TeamLevels[1])

	mso.log(fmt.Sprintf("=== Deal Started (Index: %d) ===", event.GetDealIndex()))
	mso.log(fmt.Sprintf("当前Deal Level: %d", dealLevel))
	mso.log(fmt.Sprintf("队伍0 Level: %d (玩家 0,2)", team0Level))
	mso.log(fmt.Sprintf("队伍1 Level: %d (玩家 1,3)", team1Level))
	mso.log("=======================")

	if mso.engine != nil && mso.verbose {
		gameState := mso.engine.GetGameState()
		if gameState != nil && gameState.CurrentMatch != nil && gameState.CurrentMatch.CurrentDeal != nil {
			deal := gameState.CurrentMatch.CurrentDeal
			mso.log("=== 发牌完成，玩家手牌 ===")
			for playerSeat := 0; playerSeat < 4; playerSeat++ {
				cards := deal.PlayerCards[playerSeat]
				var cardStrs []string
				for _, card := range cards {
					cardStrs = append(cardStrs, card.ToShortString())
				}
				mso.log(fmt.Sprintf("Player %d (%d cards): [%s]",
					playerSeat, len(cards), strings.Join(cardStrs, ",")))
			}
			mso.log("===========================")
		}
	}
}

func (mso *MatchSimulatorObserver) handleCardsDealt(event *sdk.GameEvent) {
	mso.log("Event: Cards Dealt")
}

func (mso *MatchSimulatorObserver) handleTributeStarted(event *sdk.GameEvent) {
	mso.log("=== 上贡阶段开始 ===")

	payload := event.GetTributeStarted()
	if payload == nil {
		return
	}

	tributeType := sdk.ConvertTributeTypeFromProto(payload.TributeType)
	mso.log(fmt.Sprintf("上贡类型: %s", tributeType))

	if len(payload.Givers) > 0 {
		var giverStrs []string
		for _, giver := range payload.Givers {
			giverStrs = append(giverStrs, fmt.Sprintf("Player %d", giver))
		}
		mso.log(fmt.Sprintf("进贡方: %s", strings.Join(giverStrs, ", ")))
	}

	if len(payload.Receivers) > 0 {
		var receiverStrs []string
		for _, receiver := range payload.Receivers {
			receiverStrs = append(receiverStrs, fmt.Sprintf("Player %d", receiver))
		}
		mso.log(fmt.Sprintf("收贡方: %s", strings.Join(receiverStrs, ", ")))
	}

	mso.log("===================")
}

func (mso *MatchSimulatorObserver) handleTributeExempted(event *sdk.GameEvent) {
	mso.log("=== 免贡（抗贡） ===")

	payload := event.GetTributeExempted()
	if payload == nil {
		return
	}

	if len(payload.BigJokerHolders) > 0 {
		mso.log("大王持有者详情：")
		totalJokers := int32(0)
		for seat, count := range payload.BigJokerHolders {
			mso.log(fmt.Sprintf("  Player %d: %d张大王", seat, count))
			totalJokers += count
		}
		if totalJokers >= 2 {
			mso.log("结果：触发抗贡，本局跳过上贡阶段")
		}
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributeCardSubmitted(event *sdk.GameEvent) {
	mso.log("=== 上贡 ===")

	payload := event.GetTributeCardSubmitted()
	if payload == nil {
		return
	}

	actorSeat := int(event.GetActorSeat())
	card := sdk.ConvertCardFromProto(payload.SubmittedCard)

	if card != nil {
		mso.log(fmt.Sprintf("Player %d 提交贡牌：%s", actorSeat, card.ToShortString()))
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributeCardSelected(event *sdk.GameEvent) {
	mso.log("=== 选牌（双下） ===")

	payload := event.GetTributeCardSelected()
	if payload == nil {
		return
	}

	actorSeat := int(event.GetActorSeat())
	selectedCard := sdk.ConvertCardFromProto(payload.SelectedCard)

	if selectedCard != nil {
		mso.log(fmt.Sprintf("Player %d 选择：%s", actorSeat, selectedCard.ToShortString()))
	}

	if payload.IsAuto {
		mso.log("注意：此次选择为自动选择")
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleReturnTribute(event *sdk.GameEvent) {
	mso.log("=== 还贡 ===")

	payload := event.GetTributeCardReturned()
	if payload == nil {
		return
	}

	actorSeat := int(event.GetActorSeat())
	returnCard := sdk.ConvertCardFromProto(payload.ReturnedCard)
	targetPlayer := int(payload.TargetPlayer)

	if returnCard != nil {
		mso.log(fmt.Sprintf("Player %d 还贡给 Player %d：%s",
			actorSeat, targetPlayer, returnCard.ToShortString()))
	}

	if payload.IsAuto {
		mso.log("选择方式：自动选择")
	} else {
		mso.log("选择方式：玩家手动选择")
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributeCompleted(event *sdk.GameEvent) {
	mso.log("=== 进贡阶段完成 ===")
	mso.log("所有上贡和还贡流程已完成，游戏阶段即将开始")
	mso.log("===================")
}

func (mso *MatchSimulatorObserver) handleTrickStarted(event *sdk.GameEvent) {
	payload := event.GetTrickStarted()
	if payload == nil {
		return
	}

	leader := int(payload.Leader)
	mso.log(fmt.Sprintf("Event: New Trick Started, Leader: Player %d", leader))

	if mso.verbose && mso.engine != nil {
		mso.log("=== Player Hands at Trick Start ===")
		for playerSeat := 0; playerSeat < 4; playerSeat++ {
			playerView := mso.engine.GetPlayerView(playerSeat)
			if playerView != nil && playerView.PlayerCards != nil {
				cards := sdk.ConvertCardsFromProto(playerView.PlayerCards)
				var cardStrs []string
				for _, card := range cards {
					cardStrs = append(cardStrs, card.ToShortString())
				}
				mso.log(fmt.Sprintf("Player %d (%d cards): [%s]",
					playerSeat, len(cards), strings.Join(cardStrs, ",")))
			} else {
				mso.log(fmt.Sprintf("Player %d: No cards", playerSeat))
			}
		}
		mso.log("====================================")
	}
}

func (mso *MatchSimulatorObserver) handlePlayerPlayed(event *sdk.GameEvent) {
	payload := event.GetPlayerPlayed()
	if payload == nil {
		return
	}

	playerSeat := int(event.GetActorSeat())
	cards := sdk.ConvertCardsFromProto(payload.Cards)

	var cardStrs []string
	for _, card := range cards {
		cardStrs = append(cardStrs, card.ToShortString())
	}

	mso.log(fmt.Sprintf("Event: Player %d played %d cards: [%s]",
		playerSeat, len(cards), strings.Join(cardStrs, ",")))
}

func (mso *MatchSimulatorObserver) handlePlayerPassed(event *sdk.GameEvent) {
	playerSeat := int(event.GetActorSeat())
	mso.log(fmt.Sprintf("Event: Player %d passed", playerSeat))
}

func (mso *MatchSimulatorObserver) handleTrickEnded(event *sdk.GameEvent) {
	payload := event.GetTrickEnded()
	if payload == nil {
		return
	}

	winner := int(payload.TrickWinner)
	mso.log(fmt.Sprintf("Event: Trick Ended, Winner: Player %d", winner))
}

func (mso *MatchSimulatorObserver) handleDealEnded(event *sdk.GameEvent) {
	payload := event.GetDealEnded()
	if payload == nil {
		return
	}

	var rankings []int
	for _, r := range payload.Rankings {
		rankings = append(rankings, int(r))
	}

	victoryType := sdk.ConvertVictoryTypeFromProto(payload.VictoryType)

	mso.log(fmt.Sprintf("Event: Deal Ended, Rankings: %v, Victory Type: %v",
		rankings, victoryType))
}

func (mso *MatchSimulatorObserver) handleMatchEnded(event *sdk.GameEvent) {
	payload := event.GetMatchEnded()
	if payload == nil {
		return
	}

	winner := int(payload.Winner)
	mso.log(fmt.Sprintf("Event: Match Ended, Winner: Team %d", winner))
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
