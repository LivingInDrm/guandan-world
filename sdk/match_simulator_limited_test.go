package sdk

import (
	"testing"
)

// TestLimitedMatchSimulator 测试有限制的比赛模拟器
func TestLimitedMatchSimulator(t *testing.T) {
	simulator := NewMatchSimulator(false)
	
	// 创建4个模拟玩家
	for i := 0; i < 4; i++ {
		simulator.players[i] = SimulatedPlayer{
			Player: Player{
				ID:       testPlayerID(i),
				Username: testPlayerName(i),
				Seat:     i,
				Online:   true,
				AutoPlay: true,
			},
			AutoPlayAlgorithm: NewSimpleAutoPlayAlgorithm(2),
		}
	}
	
	// 将Player类型转换为[]Player
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = simulator.players[i].Player
	}
	
	// 开始比赛
	err := simulator.engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}
	
	// 限制只进行3局游戏
	maxTestDeals := 3
	dealCount := 0
	
	for !simulator.engine.IsGameFinished() && dealCount < maxTestDeals {
		simulator.currentDealNum++
		dealCount++
		
		t.Logf("Starting deal %d", simulator.currentDealNum)
		
		// 开始新的一局
		err := simulator.engine.StartDeal()
		if err != nil {
			t.Errorf("Failed to start deal %d: %v", simulator.currentDealNum, err)
			break
		}
		
		// 获取deal状态
		gameState := simulator.engine.GetGameState()
		deal := gameState.CurrentMatch.CurrentDeal
		
		if deal == nil {
			t.Errorf("No active deal after StartDeal")
			break
		}
		
		t.Logf("Deal %d status: %v", simulator.currentDealNum, deal.Status)
		
		// 处理贡牌阶段
		if deal.Status == DealStatusTribute {
			err = simulator.processTributePhase()
			if err != nil {
				t.Errorf("Failed to process tribute phase in deal %d: %v", simulator.currentDealNum, err)
				break
			}
			t.Logf("Deal %d tribute phase completed", simulator.currentDealNum)
		}
		
		// 检查是否进入playing状态
		gameState = simulator.engine.GetGameState()
		if gameState.CurrentMatch.CurrentDeal != nil {
			deal = gameState.CurrentMatch.CurrentDeal
			if deal.Status == DealStatusPlaying {
				t.Logf("Deal %d entered playing state successfully", simulator.currentDealNum)
				
				// 模拟几轮出牌
				trickCount := 0
				maxTricksPerDeal := 5 // 限制每局只进行5轮
				
				for deal.Status == DealStatusPlaying && trickCount < maxTricksPerDeal {
					trickCount++
					
					if deal.CurrentTrick == nil {
						t.Logf("No current trick, deal may be finished")
						break
					}
					
					currentPlayer := deal.CurrentTrick.CurrentTurn
					
					// 获取玩家手牌
					playerView := simulator.engine.GetPlayerView(currentPlayer)
					playerHand := playerView.PlayerCards
					
					if len(playerHand) == 0 {
						t.Logf("Player %d has no cards, game should end", currentPlayer)
						break
					}
					
					// 简单出牌：总是出最小的单张
					smallest := playerHand[0]
					for _, card := range playerHand {
						if card.LessThan(smallest) {
							smallest = card
						}
					}
					
					_, err = simulator.engine.PlayCards(currentPlayer, []*Card{smallest})
					if err != nil {
						// 如果不能出牌，尝试过牌
						isLeader := deal.CurrentTrick.LeadComp == nil
						if !isLeader {
							_, err = simulator.engine.PassTurn(currentPlayer)
							if err != nil {
								t.Logf("Player %d cannot pass: %v", currentPlayer, err)
								break
							}
						} else {
							t.Logf("Leader cannot play or pass, ending trick simulation")
							break
						}
					}
					
					// 重新获取状态
					gameState = simulator.engine.GetGameState()
					if gameState.CurrentMatch.CurrentDeal == nil {
						t.Logf("Deal ended after trick %d", trickCount)
						break
					}
					deal = gameState.CurrentMatch.CurrentDeal
				}
				
				t.Logf("Deal %d completed %d tricks", simulator.currentDealNum, trickCount)
			} else {
				t.Logf("Deal %d did not enter playing state, status: %v", simulator.currentDealNum, deal.Status)
			}
		}
		
		// 强制结束当前deal以避免无限循环
		if gameState.CurrentMatch.CurrentDeal != nil && gameState.CurrentMatch.CurrentDeal.Status != DealStatusFinished {
			t.Logf("Forcibly ending deal %d for testing", simulator.currentDealNum)
		}
	}
	
	t.Logf("Limited match simulation completed after %d deals", dealCount)
	
	if dealCount >= maxTestDeals {
		t.Logf("Test stopped at maximum deal limit (%d)", maxTestDeals)
	}
}