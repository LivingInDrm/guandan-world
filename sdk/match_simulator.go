package sdk

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// MatchSimulator 模拟器，用于完整模拟掼蛋牌局
type MatchSimulator struct {
	gameEngine *GameEngine
	logger     *log.Logger
	verbose    bool
}

// NewMatchSimulator 创建新的匹配模拟器
func NewMatchSimulator(verbose bool) *MatchSimulator {
	return &MatchSimulator{
		gameEngine: NewGameEngine(),
		verbose:    verbose,
	}
}

// formatCard 将卡牌格式化为简化表示 (如9H, QS, SJ, BJ)
func (ms *MatchSimulator) formatCard(card *Card) string {
	if card == nil {
		return "??"
	}

	// 处理王牌
	if card.Color == "Joker" {
		if card.Number == 15 {
			return "SJ" // 小王
		} else if card.Number == 16 {
			return "BJ" // 大王
		}
	}

	// 处理普通牌
	var rank string
	switch card.Number {
	case 11:
		rank = "J"
	case 12:
		rank = "Q"
	case 13:
		rank = "K"
	case 14:
		rank = "A"
	default:
		rank = fmt.Sprintf("%d", card.Number)
	}

	// 花色首字母
	var suit string
	switch card.Color {
	case "Heart":
		suit = "H"
	case "Diamond":
		suit = "D"
	case "Club":
		suit = "C"
	case "Spade":
		suit = "S"
	default:
		suit = "?"
	}

	return rank + suit
}

// formatCards 将卡牌列表格式化并排序
func (ms *MatchSimulator) formatCards(cards []*Card) string {
	if len(cards) == 0 {
		return "[]"
	}

	cardStrs := make([]string, len(cards))
	for i, card := range cards {
		cardStrs[i] = ms.formatCard(card)
	}

	// 按字母序排序
	sort.Strings(cardStrs)

	return "[" + strings.Join(cardStrs, " ") + "]"
}

// logSDKCall 记录SDK接口调用
func (ms *MatchSimulator) logSDKCall(methodName string, params ...interface{}) {
	if ms.verbose {
		if len(params) > 0 {
			fmt.Printf("🔧 SDK调用: %s(%v)\n", methodName, params)
		} else {
			fmt.Printf("🔧 SDK调用: %s()\n", methodName)
		}
	}
}

// logPlayerHands 记录所有玩家的手牌
func (ms *MatchSimulator) logPlayerHands(title string, playerCards [4][]*Card) {
	if ms.verbose {
		fmt.Printf("📋 %s:\n", title)
		for i := 0; i < 4; i++ {
			fmt.Printf("  玩家%d: %s (共%d张)\n", i, ms.formatCards(playerCards[i]), len(playerCards[i]))
		}
		fmt.Println()
	}
}

// logTeamLevels 记录队伍等级信息
func (ms *MatchSimulator) logTeamLevels(match *Match, dealLevel int) {
	if ms.verbose {
		fmt.Printf("📊 队伍等级信息:\n")
		fmt.Printf("  队伍0等级: %d\n", match.TeamLevels[0])
		fmt.Printf("  队伍1等级: %d\n", match.TeamLevels[1])
		fmt.Printf("  当前局等级: %d\n", dealLevel)
		fmt.Println()
	}
}

// logPlayAction 记录出牌行为
func (ms *MatchSimulator) logPlayAction(playerSeat int, cards []*Card, isPass bool) {
	if ms.verbose {
		if isPass {
			fmt.Printf("🎯 玩家%d: 过牌\n", playerSeat)
		} else {
			fmt.Printf("🎯 玩家%d: 出牌 %s\n", playerSeat, ms.formatCards(cards))
		}
	}
}

// logTributeAction 记录上贡行为
func (ms *MatchSimulator) logTributeAction(actionType string, fromPlayer, toPlayer int, card *Card) {
	if ms.verbose {
		fmt.Printf("🎁 %s: 玩家%d -> 玩家%d, 牌: %s\n",
			actionType, fromPlayer, toPlayer, ms.formatCard(card))
	}
}

// logDealStart 记录Deal开始信息
func (ms *MatchSimulator) logDealStart(dealNumber int, level int) {
	if ms.verbose {
		fmt.Printf("\n🎮 === 第%d局开始 ===\n", dealNumber)
		fmt.Printf("🎯 牌局等级: %d\n", level)
		fmt.Println()
	}
}

// logDealResult 记录Deal结果
func (ms *MatchSimulator) logDealResult(dealNumber int, rankings []int, winningTeam int, upgrades [2]int, newLevels [2]int) {
	if ms.verbose {
		fmt.Printf("🏆 第%d局结果:\n", dealNumber)
		fmt.Printf("  排名: %v\n", rankings)
		fmt.Printf("  获胜队伍: %d\n", winningTeam)
		fmt.Printf("  等级升级: 队伍0升%d级, 队伍1升%d级\n", upgrades[0], upgrades[1])
		fmt.Printf("  新等级: 队伍0=%d级, 队伍1=%d级\n", newLevels[0], newLevels[1])
		fmt.Println()
	}
}

// SimulateMatch 模拟完整的match过程
func (ms *MatchSimulator) SimulateMatch() (*MatchResult, error) {
	// 创建4个模拟玩家
	players := ms.createSimulatedPlayers()

	if ms.verbose {
		fmt.Println("🀄 开始模拟掼蛋牌局...")
		fmt.Println("玩家信息:")
		for _, player := range players {
			fmt.Printf("  座位%d: %s (队伍%d)\n", player.Seat, player.Username, player.Seat%2)
		}
		fmt.Println()
	}

	// 启动匹配
	ms.logSDKCall("gameEngine.StartMatch", players)
	err := ms.gameEngine.StartMatch(players)
	if err != nil {
		return nil, fmt.Errorf("启动匹配失败: %w", err)
	}

	// 模拟多个deal直到比赛结束
	dealNumber := 1
	for !ms.gameEngine.IsGameFinished() {
		gameState := ms.gameEngine.GetGameState()
		if gameState.CurrentMatch != nil {
			ms.logDealStart(dealNumber, gameState.CurrentMatch.GetCurrentLevel())
			ms.logTeamLevels(gameState.CurrentMatch, gameState.CurrentMatch.GetCurrentLevel())
		}

		err = ms.simulateDeal(dealNumber)
		if err != nil {
			return nil, fmt.Errorf("模拟第%d局失败: %w", dealNumber, err)
		}

		// 检查是否有队伍达到A级
		gameState = ms.gameEngine.GetGameState()
		if gameState.CurrentMatch != nil {
			if gameState.CurrentMatch.IsAnyTeamAtALevel() {
				if ms.verbose {
					fmt.Printf("✨ 有队伍达到A级，比赛结束！\n")
				}
				break
			}
		}

		dealNumber++

		// 防止无限循环，最多模拟50局
		if dealNumber > 50 {
			if ms.verbose {
				fmt.Println("⚠️ 达到最大局数限制，强制结束比赛")
			}
			break
		}

		// 短暂延迟，便于观察
		if ms.verbose {
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 获取最终结果
	gameState := ms.gameEngine.GetGameState()
	if gameState.CurrentMatch != nil && gameState.CurrentMatch.Status == MatchStatusFinished {
		result := gameState.CurrentMatch.GetMatchResult()
		if ms.verbose {
			ms.printMatchResult(result)
		}
		return result, nil
	}

	return nil, fmt.Errorf("比赛未正常结束")
}

// simulateDeal 模拟一个完整的deal
func (ms *MatchSimulator) simulateDeal(dealNumber int) error {
	// 开始新deal
	ms.logSDKCall("gameEngine.StartDeal")
	err := ms.gameEngine.StartDeal()
	if err != nil {
		return fmt.Errorf("开始deal失败: %w", err)
	}

	gameState := ms.gameEngine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal

	// 记录开局手牌
	ms.logPlayerHands("开局手牌", deal.PlayerCards)

	// 处理上贡阶段
	if deal.TributePhase != nil {
		err = ms.simulateTributePhase(deal, dealNumber)
		if err != nil {
			return fmt.Errorf("模拟上贡阶段失败: %w", err)
		}
	}

	// 模拟出牌阶段
	trickNumber := 1
	for deal.Status == DealStatusPlaying {
		if ms.verbose {
			fmt.Printf("🎲 --- 第%d轮Trick ---\n", trickNumber)
		}

		err = ms.simulateTrick(deal, trickNumber)
		if err != nil {
			return fmt.Errorf("模拟第%d轮出牌失败: %w", trickNumber, err)
		}

		trickNumber++

		// 防止无限循环
		if trickNumber > 100 {
			return fmt.Errorf("单局超过最大轮数限制")
		}
	}

	// 确保deal结束后正确结算
	if deal.Status == DealStatusFinished {
		result, err := deal.CalculateResult(gameState.CurrentMatch)
		if err != nil {
			if ms.verbose {
				fmt.Printf("❌ 结算失败: %v\n", err)
			}
			return fmt.Errorf("deal结算失败: %w", err)
		}

		// 手动调用FinishDeal来更新match状态
		ms.logSDKCall("match.FinishDeal", result)
		err = gameState.CurrentMatch.FinishDeal(result)
		if err != nil {
			if ms.verbose {
				fmt.Printf("❌ FinishDeal失败: %v\n", err)
			}
			return fmt.Errorf("FinishDeal失败: %w", err)
		}

		// 记录Deal结果
		ms.logDealResult(dealNumber, deal.Rankings, result.WinningTeam, result.Upgrades, gameState.CurrentMatch.TeamLevels)
	} else {
		if ms.verbose {
			fmt.Printf("⚠️ 第%d局未正常结束，状态: %s\n", dealNumber, deal.Status)
			fmt.Println()
		}
	}

	return nil
}

// simulateTributePhase 模拟上贡阶段
func (ms *MatchSimulator) simulateTributePhase(deal *Deal, dealNumber int) error {
	if deal.TributePhase == nil {
		return nil
	}

	if ms.verbose {
		fmt.Println("🎁 === 上贡阶段 ===")
	}

	// 记录上贡前手牌
	ms.logPlayerHands("上贡前手牌", deal.PlayerCards)

	// 检查免贡条件
	gameState := ms.gameEngine.GetGameState()
	match := gameState.CurrentMatch

	// 获取上一局结果来检查免贡
	lastResult := match.getLastDealResult()
	if lastResult != nil {
		// 创建TributeManager来检查免贡
		tm := NewTributeManager(deal.Level)
		isImmune := tm.CheckTributeImmunity(lastResult, deal.PlayerCards)

		if isImmune {
			if ms.verbose {
				fmt.Println("🔰 触发免贡条件，跳过上贡阶段")
				ms.logImmunityReason(lastResult, deal.PlayerCards)
			}
			// 跳过上贡，直接开始出牌
			deal.Status = DealStatusPlaying
			err := deal.startFirstTrick()
			if err != nil {
				return fmt.Errorf("开始第一个trick失败: %w", err)
			}
			return nil
		}
	}

	// 继续正常上贡流程
	maxIterations := 10 // 防止无限循环
	iterations := 0

	for deal.Status == DealStatusTribute && iterations < maxIterations {
		tributePhase := deal.TributePhase
		hasAction := false

		if tributePhase.Status == TributeStatusSelecting {
			// 处理选贡（从池中选择） - 双下场景
			if tributePhase.SelectingPlayer >= 0 && len(tributePhase.PoolCards) > 0 {
				// 选择最大的牌
				selectedCard := ms.selectBestCard(tributePhase.PoolCards)
				ms.logTributeAction("选贡", tributePhase.SelectingPlayer, -1, selectedCard)
				ms.logSDKCall("gameEngine.SelectTribute", tributePhase.SelectingPlayer, selectedCard)

				_, err := ms.gameEngine.SelectTribute(tributePhase.SelectingPlayer, selectedCard)
				if err != nil {
					if ms.verbose {
						fmt.Printf("⚠️ 选贡失败: %v\n", err)
					}
					return err
				}
				hasAction = true
			} else {
				// 如果没有pool cards或selecting player，需要先创建pool
				err := ms.createTributePool(deal)
				if err != nil {
					if ms.verbose {
						fmt.Printf("⚠️ 创建上贡池失败: %v\n", err)
					}
					return err
				}
				hasAction = true
			}
		} else if tributePhase.Status == TributeStatusReturning {
			// 在returning状态下，处理还贡牌选择
			// 上贡牌已由系统自动选择（使用新的排除红桃Trump逻辑）
			if ms.verbose {
				fmt.Printf("🎁 处理还贡阶段...\n")
				fmt.Printf("   TributeMap: %v\n", tributePhase.TributeMap)

				// 显示具体的上贡牌
				if len(tributePhase.TributeCards) > 0 {
					tributeDetails := make([]string, 0)
					for giver, card := range tributePhase.TributeCards {
						tributeDetails = append(tributeDetails, fmt.Sprintf("%d→%s", giver, ms.formatCard(card)))
					}
					fmt.Printf("   TributeCards: %v\n", tributeDetails)
				} else {
					fmt.Printf("   TributeCards: []\n")
				}

				// 显示具体的还贡牌
				if len(tributePhase.ReturnCards) > 0 {
					returnDetails := make([]string, 0)
					for receiver, card := range tributePhase.ReturnCards {
						returnDetails = append(returnDetails, fmt.Sprintf("%d→%s", receiver, ms.formatCard(card)))
					}
					fmt.Printf("   ReturnCards: %v\n", returnDetails)
				} else {
					fmt.Printf("   ReturnCards: []\n")
				}

				// 显示已选择的上贡牌的详细动作
				for giver, card := range tributePhase.TributeCards {
					if receiver := tributePhase.TributeMap[giver]; receiver != -1 {
						ms.logTributeAction("上贡（系统自动）", giver, receiver, card)
					}
				}

				// 显示还贡牌的详细动作
				for receiver, card := range tributePhase.ReturnCards {
					if card != nil {
						// 找到对应的给方
						for giver, rec := range tributePhase.TributeMap {
							if rec == receiver {
								ms.logTributeAction("还贡（系统自动）", receiver, giver, card)
								break
							}
						}
					}
				}
			}

			// 直接完成还贡阶段
			err := ms.completeReturningPhase(deal)
			if err != nil {
				if ms.verbose {
					fmt.Printf("⚠️ 完成还贡阶段失败: %v\n", err)
				}
				return err
			}
			hasAction = true
		} else if tributePhase.Status == TributeStatusFinished {
			// 上贡阶段已经完成
			break
		}

		// 处理超时，这会推进上贡阶段的状态
		timeoutEvents := ms.gameEngine.ProcessTimeouts()
		if len(timeoutEvents) > 0 {
			hasAction = true
		}

		// 如果是returning状态，强制推进一次来完成自动还贡
		if tributePhase.Status == TributeStatusReturning && !hasAction {
			// 尝试手动完成还贡阶段
			err := ms.completeReturningPhase(deal)
			if err != nil {
				if ms.verbose {
					fmt.Printf("⚠️ 完成还贡阶段失败: %v\n", err)
				}
			} else {
				hasAction = true
			}
		}

		if !hasAction {
			// 如果没有任何动作，避免无限循环
			if ms.verbose {
				fmt.Printf("⚠️ 上贡阶段无进展，状态: %s\n", tributePhase.Status)
			}
			break
		}

		iterations++
	}

	// 记录上贡后手牌
	ms.logPlayerHands("上贡后手牌", deal.PlayerCards)

	return nil
}

// completeReturningPhase 完成还贡阶段
func (ms *MatchSimulator) completeReturningPhase(deal *Deal) error {
	if deal.TributePhase == nil || deal.TributePhase.Status != TributeStatusReturning {
		return fmt.Errorf("不在还贡状态")
	}

	tributePhase := deal.TributePhase

	// 使用TributeManager来处理还贡逻辑
	tm := NewTributeManager(deal.Level)

	// 确保上贡牌已经选择完成
	err := tm.ProcessTribute(tributePhase, deal.PlayerCards)
	if err != nil {
		return fmt.Errorf("处理上贡逻辑失败: %w", err)
	}

	// 在处理完成后显示具体的上贡牌
	if ms.verbose && len(tributePhase.TributeCards) > 0 {
		fmt.Printf("🎁 上贡牌选择完成:\n")
		for giver, card := range tributePhase.TributeCards {
			if receiver := tributePhase.TributeMap[giver]; receiver != -1 {
				fmt.Printf("   玩家%d → 玩家%d: %s\n", giver, receiver, ms.formatCard(card))
			}
		}
	}

	// 自动为每个接受者选择还贡卡牌（如果还没选择的话）
	for giver, receiver := range tributePhase.TributeMap {
		if receiver != -1 && tributePhase.TributeCards[giver] != nil {
			if tributePhase.ReturnCards[receiver] == nil {
				// 选择还贡卡牌
				returnCard := ms.selectReturnCard(deal.PlayerCards[receiver])
				if returnCard != nil {
					// 直接设置还贡卡牌
					tributePhase.ReturnCards[receiver] = returnCard
					if ms.verbose {
						ms.logTributeAction("还贡（自动）", receiver, giver, returnCard)
					}
				}
			}
		}
	}

	// 显示完整的还贡情况
	if ms.verbose && len(tributePhase.ReturnCards) > 0 {
		fmt.Printf("🎁 还贡牌选择完成:\n")
		for receiver, card := range tributePhase.ReturnCards {
			if card != nil {
				// 找到对应的给方
				for giver, rec := range tributePhase.TributeMap {
					if rec == receiver {
						fmt.Printf("   玩家%d → 玩家%d: %s\n", receiver, giver, ms.formatCard(card))
						break
					}
				}
			}
		}
	}

	// 完成还贡阶段
	tributePhase.Status = TributeStatusFinished

	// 应用上贡效果到手牌
	err = tm.ApplyTributeToHands(tributePhase, &deal.PlayerCards)
	if err != nil {
		return fmt.Errorf("应用上贡效果失败: %w", err)
	}

	// 开始第一个trick
	err = deal.startFirstTrick()
	if err != nil {
		return fmt.Errorf("开始第一个trick失败: %w", err)
	}

	deal.Status = DealStatusPlaying

	return nil
}

// simulateTrick 模拟一轮出牌
func (ms *MatchSimulator) simulateTrick(deal *Deal, trickNumber int) error {
	if deal.CurrentTrick == nil {
		return fmt.Errorf("当前没有活跃的trick")
	}

	trick := deal.CurrentTrick
	maxPlays := 16 // 每轮最多16次出牌（4个玩家每人最多4次）
	playCount := 0

	for trick.Status == TrickStatusPlaying && playCount < maxPlays {
		// 处理超时事件和自动状态转换
		timeoutEvents := ms.gameEngine.ProcessTimeouts()
		if len(timeoutEvents) > 0 && ms.verbose {
			fmt.Printf("🕐 处理了%d个超时事件\n", len(timeoutEvents))
		}

		currentPlayer := trick.CurrentTurn

		// 检查玩家是否还有牌
		if len(deal.PlayerCards[currentPlayer]) == 0 {
			// 玩家已经出完牌，需要跳过这个玩家
			// 但不能直接过牌，因为如果是trick leader就会出错
			// 直接跳到下一个玩家

			// 更新到下一个有牌的玩家
			playersChecked := 0
			for playersChecked < 4 {
				currentPlayer = (currentPlayer + 1) % 4
				playersChecked++
				if len(deal.PlayerCards[currentPlayer]) > 0 {
					break
				}
			}

			// 如果所有玩家都没牌了，trick结束
			if playersChecked == 4 && len(deal.PlayerCards[currentPlayer]) == 0 {
				break
			}

			// 更新trick的当前玩家
			trick.CurrentTurn = currentPlayer
			continue
		}

		// 自动出牌逻辑
		if trick.LeadComp == nil {
			// 首出玩家必须出牌
			cards := ms.selectFirstPlayCards(deal.PlayerCards[currentPlayer])
			if len(cards) == 0 && len(deal.PlayerCards[currentPlayer]) > 0 {
				// 保底：如果选择逻辑失败，出最小的单张
				cards = []*Card{ms.findSmallestCard(deal.PlayerCards[currentPlayer])}
			}
			if len(cards) > 0 {
				ms.logPlayAction(currentPlayer, cards, false)
				ms.logSDKCall("gameEngine.PlayCards", currentPlayer, cards)

				_, err := ms.gameEngine.PlayCards(currentPlayer, cards)
				if err != nil {
					// 如果还是失败，强制出第一张牌
					if len(deal.PlayerCards[currentPlayer]) > 0 {
						ms.logPlayAction(currentPlayer, []*Card{deal.PlayerCards[currentPlayer][0]}, false)
						ms.logSDKCall("gameEngine.PlayCards", currentPlayer, []*Card{deal.PlayerCards[currentPlayer][0]})

						_, err = ms.gameEngine.PlayCards(currentPlayer, []*Card{deal.PlayerCards[currentPlayer][0]})
						if err != nil {
							return err
						}
					}
				}
			}
		} else {
			// 跟牌玩家：尝试压过，否则过牌
			cards := ms.selectFollowPlayCards(deal.PlayerCards[currentPlayer], trick.LeadComp)
			if len(cards) > 0 {
				// 尝试出牌
				ms.logPlayAction(currentPlayer, cards, false)
				ms.logSDKCall("gameEngine.PlayCards", currentPlayer, cards)

				_, err := ms.gameEngine.PlayCards(currentPlayer, cards)
				if err != nil {
					// 出牌失败，过牌
					ms.logPlayAction(currentPlayer, nil, true)
					ms.logSDKCall("gameEngine.PassTurn", currentPlayer)

					_, err = ms.gameEngine.PassTurn(currentPlayer)
					if err != nil {
						if ms.verbose {
							fmt.Printf("❌ PassTurn失败 - Player %d, LeadComp: %v, Error: %v\n",
								currentPlayer, trick.LeadComp, err)
						}
						return err
					}
				}
			} else {
				// 不能压过，过牌
				ms.logPlayAction(currentPlayer, nil, true)
				ms.logSDKCall("gameEngine.PassTurn", currentPlayer)

				_, err := ms.gameEngine.PassTurn(currentPlayer)
				if err != nil {
					if ms.verbose {
						fmt.Printf("❌ PassTurn失败 - Player %d, LeadComp: %v, Error: %v\n",
							currentPlayer, trick.LeadComp, err)
					}
					return err
				}
			}
		}

		playCount++

		// 检查deal是否结束
		if deal.Status != DealStatusPlaying {
			break
		}

		// 检查trick是否结束
		if trick.Status != TrickStatusPlaying {
			// Trick结束，记录剩余手牌
			ms.logPlayerHands(fmt.Sprintf("第%d轮Trick结束后手牌", trickNumber), deal.PlayerCards)
			break
		}

		// 短暂延迟
		if ms.verbose {
			time.Sleep(5 * time.Millisecond)
		}
	}

	return nil
}

// createSimulatedPlayers 创建4个模拟玩家
func (ms *MatchSimulator) createSimulatedPlayers() []Player {
	players := make([]Player, 4)

	playerNames := []string{"Alice", "Bob", "Charlie", "Diana"}

	for i := 0; i < 4; i++ {
		players[i] = Player{
			ID:       fmt.Sprintf("player_%d", i),
			Username: playerNames[i],
			Seat:     i,
			Online:   true,
			AutoPlay: true, // 启用自动出牌
		}
	}

	return players
}

// selectAutoPlayCards 自动选择出牌
func (ms *MatchSimulator) selectAutoPlayCards(playerCards []*Card, leadComp CardComp) []*Card {
	if len(playerCards) == 0 {
		return nil
	}

	if leadComp == nil {
		// 首出：选择张数尽可能多的非炸弹牌
		return ms.selectFirstPlayCards(playerCards)
	} else {
		// 跟牌：如果能压过，则出牌
		return ms.selectFollowPlayCards(playerCards, leadComp)
	}
}

// selectFirstPlayCards 首出时选择牌
func (ms *MatchSimulator) selectFirstPlayCards(playerCards []*Card) []*Card {
	// 尝试各种牌型，优先选择张数多的非炸弹牌

	// 尝试钢板(6张)
	if plateCards := ms.findPlateCards(playerCards); len(plateCards) > 0 {
		return plateCards
	}

	// 尝试顺子(5张)
	if straightCards := ms.findStraightCards(playerCards); len(straightCards) > 0 {
		return straightCards
	}

	// 尝试葫芦(5张)
	if fullHouseCards := ms.findFullHouseCards(playerCards); len(fullHouseCards) > 0 {
		return fullHouseCards
	}

	// 尝试三张
	if tripleCards := ms.findTripleCards(playerCards); len(tripleCards) > 0 {
		return tripleCards
	}

	// 尝试对子
	if pairCards := ms.findPairCards(playerCards); len(pairCards) > 0 {
		return pairCards
	}

	// 最后出单张（最小的）
	return []*Card{ms.findSmallestCard(playerCards)}
}

// selectFollowPlayCards 跟牌时选择牌
func (ms *MatchSimulator) selectFollowPlayCards(playerCards []*Card, leadComp CardComp) []*Card {
	leadType := leadComp.GetType()

	// 根据leadComp的类型，找能压过的最小组合
	switch leadType {
	case TypeSingle:
		return ms.findBeatingCard(playerCards, leadComp, 1)
	case TypePair:
		return ms.findBeatingCard(playerCards, leadComp, 2)
	case TypeTriple:
		return ms.findBeatingCard(playerCards, leadComp, 3)
	case TypeStraight:
		return ms.findBeatingStraight(playerCards, leadComp)
	case TypeFullHouse:
		return ms.findBeatingFullHouse(playerCards, leadComp)
	case TypePlate:
		return ms.findBeatingPlate(playerCards, leadComp)
	default:
		// 对于炸弹等特殊牌型，只有更大的炸弹能压过
		return ms.findBeatingBomb(playerCards, leadComp)
	}
}

// 辅助方法实现（这些是简化的实现）

func (ms *MatchSimulator) findSmallestCard(cards []*Card) *Card {
	if len(cards) == 0 {
		return nil
	}

	smallest := cards[0]
	for _, card := range cards[1:] {
		if card.LessThan(smallest) {
			smallest = card
		}
	}
	return smallest
}

func (ms *MatchSimulator) findPairCards(cards []*Card) []*Card {
	// 简化实现：找到第一个对子
	cardCounts := make(map[int]int)
	cardMap := make(map[int][]*Card)

	for _, card := range cards {
		cardCounts[card.Number]++
		cardMap[card.Number] = append(cardMap[card.Number], card)
	}

	for number, count := range cardCounts {
		if count >= 2 {
			return cardMap[number][:2]
		}
	}

	return nil
}

func (ms *MatchSimulator) findTripleCards(cards []*Card) []*Card {
	// 简化实现：找到第一个三张
	cardCounts := make(map[int]int)
	cardMap := make(map[int][]*Card)

	for _, card := range cards {
		cardCounts[card.Number]++
		cardMap[card.Number] = append(cardMap[card.Number], card)
	}

	for number, count := range cardCounts {
		if count >= 3 {
			return cardMap[number][:3]
		}
	}

	return nil
}

func (ms *MatchSimulator) findStraightCards(cards []*Card) []*Card {
	// 简化实现：尝试找5张连续的牌
	// 这是一个基础实现，实际情况会更复杂
	if len(cards) < 5 {
		return nil
	}

	sorted := make([]*Card, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Number < sorted[j].Number
	})

	for i := 0; i <= len(sorted)-5; i++ {
		sequence := make([]*Card, 5)
		isValid := true

		for j := 0; j < 5; j++ {
			if j == 0 {
				sequence[j] = sorted[i+j]
			} else {
				if sorted[i+j].Number != sorted[i+j-1].Number+1 {
					isValid = false
					break
				}
				sequence[j] = sorted[i+j]
			}
		}

		if isValid {
			return sequence
		}
	}

	return nil
}

func (ms *MatchSimulator) findFullHouseCards(cards []*Card) []*Card {
	// 简化实现：找三张+对子
	tripleCards := ms.findTripleCards(cards)
	if len(tripleCards) == 0 {
		return nil
	}

	// 从剩余牌中找对子
	remainingCards := []*Card{}
	tripleNumber := tripleCards[0].Number

	for _, card := range cards {
		if card.Number != tripleNumber {
			remainingCards = append(remainingCards, card)
		}
	}

	pairCards := ms.findPairCards(remainingCards)
	if len(pairCards) == 0 {
		return nil
	}

	result := make([]*Card, 0, 5)
	result = append(result, tripleCards...)
	result = append(result, pairCards...)
	return result
}

func (ms *MatchSimulator) findPlateCards(cards []*Card) []*Card {
	// 简化实现：找6张牌的钢板（三个对子）
	// 这是一个非常简化的实现
	cardCounts := make(map[int]int)
	cardMap := make(map[int][]*Card)

	for _, card := range cards {
		cardCounts[card.Number]++
		cardMap[card.Number] = append(cardMap[card.Number], card)
	}

	pairs := [][]*Card{}
	for number, count := range cardCounts {
		if count >= 2 {
			pairs = append(pairs, cardMap[number][:2])
		}
	}

	if len(pairs) >= 3 {
		result := make([]*Card, 0, 6)
		for i := 0; i < 3; i++ {
			result = append(result, pairs[i]...)
		}
		return result
	}

	return nil
}

func (ms *MatchSimulator) findBeatingCard(cards []*Card, leadComp CardComp, count int) []*Card {
	// 简化实现：找能压过leadComp的牌
	for _, card := range cards {
		testCards := []*Card{card}
		if count == 2 {
			// 找对子
			for _, card2 := range cards {
				if card2 != card && card2.Number == card.Number {
					testCards = append(testCards, card2)
					break
				}
			}
		} else if count == 3 {
			// 找三张
			foundCount := 0
			for _, card2 := range cards {
				if card2 != card && card2.Number == card.Number && foundCount < 2 {
					testCards = append(testCards, card2)
					foundCount++
				}
			}
		}

		if len(testCards) == count {
			comp := FromCardList(testCards, leadComp)
			if comp.IsValid() && comp.GreaterThan(leadComp) {
				return testCards
			}
		}
	}

	return nil
}

func (ms *MatchSimulator) findBeatingStraight(cards []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时返回nil
	return nil
}

func (ms *MatchSimulator) findBeatingFullHouse(cards []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时返回nil
	return nil
}

func (ms *MatchSimulator) findBeatingPlate(cards []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时返回nil
	return nil
}

func (ms *MatchSimulator) findBeatingBomb(cards []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时返回nil
	return nil
}

// selectBestCard 选择最佳的牌（上贡时选最大的）
func (ms *MatchSimulator) selectBestCard(cards []*Card) *Card {
	if len(cards) == 0 {
		return nil
	}

	best := cards[0]
	for _, card := range cards[1:] {
		if card.GreaterThan(best) {
			best = card
		}
	}
	return best
}

// selectReturnCard 选择还贡的牌（选不破坏炸弹的最小牌）
func (ms *MatchSimulator) selectReturnCard(cards []*Card) *Card {
	if len(cards) == 0 {
		return nil
	}

	// 简化实现：选择最小的牌
	// 实际应该考虑不破坏炸弹等组合
	return ms.findSmallestCard(cards)
}

// printMatchResult 打印比赛结果
func (ms *MatchSimulator) printMatchResult(result *MatchResult) {
	if result == nil {
		return
	}

	fmt.Println("🎉 === 比赛结果 ===")
	fmt.Printf("🏆 获胜队伍: %d\n", result.Winner)
	fmt.Printf("📊 最终等级: 队伍0=%d级, 队伍1=%d级\n", result.FinalLevels[0], result.FinalLevels[1])
	fmt.Printf("⏰ 比赛时长: %v\n", result.Duration)

	if result.Statistics != nil {
		fmt.Printf("🎯 总局数: %d\n", result.Statistics.TotalDeals)
		for i, teamStats := range result.Statistics.TeamStats {
			if teamStats != nil {
				fmt.Printf("   队伍%d: 获胜%d局, 总升级%d级, 总墩数%d\n",
					i, teamStats.DealsWon, teamStats.Upgrades, teamStats.TotalTricks)
			}
		}
	}
	fmt.Println("========================")
}

// createTributePool 创建上贡池（双下场景）
func (ms *MatchSimulator) createTributePool(deal *Deal) error {
	if deal.TributePhase == nil {
		return fmt.Errorf("没有上贡阶段")
	}

	tributePhase := deal.TributePhase

	// 使用TributeManager的逻辑来处理上贡池
	tm := NewTributeManager(deal.Level)
	err := tm.ProcessTribute(tributePhase, deal.PlayerCards)
	if err != nil {
		return fmt.Errorf("处理上贡失败: %w", err)
	}

	// 记录贡献到池的牌
	if ms.verbose {
		for giver, card := range tributePhase.TributeCards {
			if tributePhase.TributeMap[giver] == -1 { // -1 表示贡献到池
				ms.logTributeAction("贡献到池", giver, -1, card)
			}
		}

		fmt.Printf("🎁 上贡池包含 %d 张牌: %s\n",
			len(tributePhase.PoolCards), ms.formatCards(tributePhase.PoolCards))

		if tributePhase.SelectingPlayer >= 0 {
			fmt.Printf("🎁 当前选贡玩家: %d\n", tributePhase.SelectingPlayer)
		}
	}

	return nil
}

// logImmunityReason 记录免贡原因
func (ms *MatchSimulator) logImmunityReason(lastResult *DealResult, playerHands [4][]*Card) {
	if !ms.verbose || lastResult == nil {
		return
	}

	rankings := lastResult.Rankings
	if len(rankings) < 4 {
		return
	}

	switch lastResult.VictoryType {
	case VictoryTypeDoubleDown:
		rank3 := rankings[2]
		rank4 := rankings[3]
		bigJokersRank3 := ms.countBigJokersInHand(playerHands[rank3])
		bigJokersRank4 := ms.countBigJokersInHand(playerHands[rank4])
		fmt.Printf("   Double Down免贡: Rank3(玩家%d)有%d张大王, Rank4(玩家%d)有%d张大王, 合计%d张\n",
			rank3, bigJokersRank3, rank4, bigJokersRank4, bigJokersRank3+bigJokersRank4)

	case VictoryTypeSingleLast:
		rank4 := rankings[3]
		bigJokers := ms.countBigJokersInHand(playerHands[rank4])
		fmt.Printf("   Single Last免贡: Rank4(玩家%d)有%d张大王\n", rank4, bigJokers)

	case VictoryTypePartnerLast:
		rank3 := rankings[2]
		bigJokers := ms.countBigJokersInHand(playerHands[rank3])
		fmt.Printf("   Partner Last免贡: Rank3(玩家%d)有%d张大王\n", rank3, bigJokers)
	}
}

// countBigJokersInHand 统计手牌中大王的数量 (模拟器专用)
func (ms *MatchSimulator) countBigJokersInHand(hand []*Card) int {
	count := 0
	for _, card := range hand {
		if card.Number == 16 && card.Color == "Joker" { // Red Joker = Big Joker
			count++
		}
	}
	return count
}
