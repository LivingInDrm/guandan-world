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
	case sdk.EventTributePhase:
		mso.handleTributePhase(event)
	case sdk.EventTributeRulesSet:
		mso.handleTributeRulesSet(event)
	case sdk.EventTributeImmunity:
		mso.handleTributeImmunity(event)
	case sdk.EventTributePoolCreated:
		mso.handleTributePoolCreated(event)
	case sdk.EventTributeStarted:
		mso.handleTributeStarted(event)
	case sdk.EventTributeGiven:
		mso.handleTributeGiven(event)
	case sdk.EventTributeSelected:
		mso.handleTributeSelected(event)
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

	// 从事件数据中提取level信息
	if eventData, ok := event.Data.(map[string]interface{}); ok {
		dealLevel := eventData["deal_level"].(int)
		team0Level := eventData["team0_level"].(int)
		team1Level := eventData["team1_level"].(int)
		deal := eventData["deal"].(*sdk.Deal)

		// 记录level信息
		mso.log(fmt.Sprintf("=== Deal %s Started ===", deal.ID))
		mso.log(fmt.Sprintf("当前Deal Level: %d", dealLevel))
		mso.log(fmt.Sprintf("队伍0 Level: %d (玩家 0,2)", team0Level))
		mso.log(fmt.Sprintf("队伍1 Level: %d (玩家 1,3)", team1Level))
		mso.log("=======================")

		// 显示所有玩家的手牌信息
		if deal != nil {
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

func (mso *MatchSimulatorObserver) handleTributePhase(event *sdk.GameEvent) {
	mso.log("=== 进贡阶段开始 ===")
	mso.log("准备进行上贡、抗贡检查和还贡流程")
	mso.log("==================")
}

func (mso *MatchSimulatorObserver) handleTributeRulesSet(event *sdk.GameEvent) {
	mso.log("=== 上贡规则确定 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if lastResult, ok := data["last_result"].(*sdk.DealResult); ok {
			mso.log(fmt.Sprintf("上局结果：%v, 胜利类型：%v", lastResult.Rankings, lastResult.VictoryType))
		}

		if tributeRules, ok := data["tribute_rules"].(map[string]interface{}); ok {
			if description, ok := tributeRules["description"].(string); ok {
				mso.log(fmt.Sprintf("上贡规则：%s", description))
			}

			if isDoubleDown, ok := tributeRules["is_double_down"].(bool); ok {
				if isDoubleDown {
					mso.log("类型：双下上贡（贡牌池模式）")
				} else {
					mso.log("类型：直接上贡模式")
				}
			}
		}
	}

	mso.log("====================")
}

func (mso *MatchSimulatorObserver) handleTributeImmunity(event *sdk.GameEvent) {
	mso.log("=== 抗贡检查 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if immunityReason, ok := data["immunity_reason"].(map[string]interface{}); ok {
			// 输出详细抗贡信息
			if description, ok := immunityReason["description"].(string); ok {
				mso.log(fmt.Sprintf("抗贡结果：%s", description))
			}

			// 输出大王持有者详情
			if holders, ok := immunityReason["big_joker_holders"].([]map[string]interface{}); ok {
				if len(holders) > 0 {
					mso.log("大王持有者详情：")
					for _, holder := range holders {
						if playerSeat, ok := holder["player_seat"].(int); ok {
							if count, ok := holder["big_joker_count"].(int); ok {
								mso.log(fmt.Sprintf("  Player %d: %d张大王", playerSeat, count))
							}
						}
					}
				}
			}

			if totalCount, ok := immunityReason["big_joker_count"].(int); ok {
				if totalCount >= 2 {
					mso.log("结果：触发抗贡，本局跳过上贡阶段")
				} else {
					mso.log("结果：大王数量不足，正常进行上贡")
				}
			}
		}
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributePoolCreated(event *sdk.GameEvent) {
	mso.log("=== 双下贡牌池创建 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if description, ok := data["description"].(string); ok {
			mso.log(description)
		}

		if contributors, ok := data["contributors"].([]map[string]interface{}); ok {
			mso.log("贡献详情：")
			for _, contributor := range contributors {
				if playerSeat, ok := contributor["player_seat"].(int); ok {
					if card, ok := contributor["card"].(*sdk.Card); ok {
						mso.log(fmt.Sprintf("  Player %d 贡献：%s", playerSeat, card.ToShortString()))
					}
				}
			}
		}

		if selectionOrder, ok := data["selection_order"].([]int); ok {
			if len(selectionOrder) > 0 {
				orderStr := fmt.Sprintf("Player %d", selectionOrder[0])
				for i := 1; i < len(selectionOrder); i++ {
					orderStr += fmt.Sprintf(" -> Player %d", selectionOrder[i])
				}
				mso.log(fmt.Sprintf("选择顺序：%s", orderStr))
			}
		}

		if poolCards, ok := data["pool_cards"].([]*sdk.Card); ok {
			var cardStrs []string
			for _, card := range poolCards {
				cardStrs = append(cardStrs, card.ToShortString())
			}
			mso.log(fmt.Sprintf("池中牌张：[%s]", strings.Join(cardStrs, ", ")))
		}
	}

	mso.log("=====================")
}

func (mso *MatchSimulatorObserver) handleTributeStarted(event *sdk.GameEvent) {
	mso.log("=== 上贡执行开始 ===")
	mso.log("开始执行具体的上贡流程")
	mso.log("===================")
}

func (mso *MatchSimulatorObserver) handleTributeGiven(event *sdk.GameEvent) {
	mso.log("=== 上贡完成 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if giver, ok := data["giver"].(int); ok {
			if receiver, ok := data["receiver"].(int); ok {
				if card, ok := data["card"].(*sdk.Card); ok {
					mso.log(fmt.Sprintf("Player %d 上贡给 Player %d：%s",
						giver, receiver, card.ToShortString()))
				}
			}
		}

		if tributeType, ok := data["tribute_type"].(string); ok {
			typeText := "普通上贡"
			if tributeType == "double_down_pool" {
				typeText = "双下池贡献"
			}
			mso.log(fmt.Sprintf("上贡类型：%s", typeText))
		}

		if isAutoSelected, ok := data["is_auto_selected"].(bool); ok {
			if isAutoSelected {
				if reason, ok := data["selection_reason"].(string); ok {
					mso.log(fmt.Sprintf("选择方式：自动选择（%s）", reason))
				}
			} else {
				mso.log("选择方式：玩家手动选择")
			}
		}
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributeSelected(event *sdk.GameEvent) {
	mso.log("=== 双下选牌 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if player, ok := data["player"].(int); ok {
			if selectedCard, ok := data["selected_card"].(*sdk.Card); ok && selectedCard != nil {
				if selectionOrder, ok := data["selection_order"].(int); ok {
					orderText := "第一次选择"
					if selectionOrder == 2 {
						orderText = "第二次选择"
					}
					mso.log(fmt.Sprintf("Player %d (%s) 选择：%s",
						player, orderText, selectedCard.ToShortString()))
				}
			}
		}

		if remainingOptions, ok := data["remaining_options"].([]*sdk.Card); ok {
			if len(remainingOptions) > 0 {
				var cardStrs []string
				for _, card := range remainingOptions {
					if card != nil {
						cardStrs = append(cardStrs, card.ToShortString())
					}
				}
				mso.log(fmt.Sprintf("剩余选项：[%s]", strings.Join(cardStrs, ", ")))
			} else {
				mso.log("所有贡牌已选择完毕")
			}
		}

		if isTimeout, ok := data["is_timeout"].(bool); ok && isTimeout {
			mso.log("注意：此次选择为超时自动选择")
		}
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleReturnTribute(event *sdk.GameEvent) {
	mso.log("=== 还贡阶段 ===")

	if data, ok := event.Data.(map[string]interface{}); ok {
		if returner, ok := data["player"].(int); ok {
			if returnCard, ok := data["return_card"].(*sdk.Card); ok {
				if targetPlayer, ok := data["target_player"].(int); ok {
					mso.log(fmt.Sprintf("Player %d 还贡给 Player %d：%s",
						returner, targetPlayer, returnCard.ToShortString()))
				}
			}
		}

		if originalTribute, ok := data["original_tribute"].(*sdk.Card); ok {
			mso.log(fmt.Sprintf("原收到贡牌：%s", originalTribute.ToShortString()))
		}

		if isAutoSelected, ok := data["is_auto_selected"].(bool); ok {
			if isAutoSelected {
				if reason, ok := data["selection_reason"].(string); ok {
					mso.log(fmt.Sprintf("选择方式：自动选择（%s）", reason))
				}
			} else {
				mso.log("选择方式：玩家手动选择")
			}
		}
	}

	mso.log("================")
}

func (mso *MatchSimulatorObserver) handleTributeCompleted(event *sdk.GameEvent) {
	mso.log("=== 进贡阶段完成 ===")
	mso.log("所有上贡和还贡流程已完成，游戏阶段即将开始")
	mso.log("===================")
}

func (mso *MatchSimulatorObserver) handleTrickStarted(event *sdk.GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if leader, ok := data["leader"].(int); ok {
			mso.log(fmt.Sprintf("Event: New Trick Started, Leader: Player %d", leader))

			// 输出每个玩家的手牌信息（从事件数据中获取，避免死锁）
			if mso.verbose {
				if playerHands, ok := data["player_hands"].(map[int][]*sdk.Card); ok {
					mso.log("=== Player Hands at Trick Start ===")
					for playerSeat := 0; playerSeat < 4; playerSeat++ {
						if cards, exists := playerHands[playerSeat]; exists {
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
				} else {
					mso.log(fmt.Sprintf("New Trick Started (Leader: Player %d) - Hand details not available", leader))
				}
			}
		}
	}
}

func (mso *MatchSimulatorObserver) handlePlayerPlayed(event *sdk.GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		cards := data["cards"].([]*sdk.Card)

		// 将出牌转换为简化格式
		var cardStrs []string
		for _, card := range cards {
			cardStrs = append(cardStrs, card.ToShortString())
		}

		mso.log(fmt.Sprintf("Event: Player %d played %d cards: [%s]",
			playerSeat, len(cards), strings.Join(cardStrs, ",")))
	}
}

func (mso *MatchSimulatorObserver) handlePlayerPassed(event *sdk.GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		mso.log(fmt.Sprintf("Event: Player %d passed", playerSeat))
	}
}

func (mso *MatchSimulatorObserver) handleTrickEnded(event *sdk.GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if winner, ok := data["winner"].(int); ok {
			mso.log(fmt.Sprintf("Event: Trick Ended, Winner: Player %d", winner))
		}
	}
}

func (mso *MatchSimulatorObserver) handleDealEnded(event *sdk.GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if result, ok := data["result"].(*sdk.DealResult); ok {
			mso.log(fmt.Sprintf("Event: Deal Ended, Rankings: %v, Victory Type: %v",
				result.Rankings, result.VictoryType))
		}
	}
}

func (mso *MatchSimulatorObserver) handleMatchEnded(event *sdk.GameEvent) {
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
