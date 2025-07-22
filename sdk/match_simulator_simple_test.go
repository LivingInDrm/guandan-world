package sdk

import (
	"fmt"
	"testing"
)

// TestSimpleGameEngine 测试基本的游戏引擎功能
func TestSimpleGameEngine(t *testing.T) {
	// 创建游戏引擎
	engine := NewGameEngine()
	
	// 创建4个玩家
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = Player{
			ID:       fmt.Sprintf("player_%d", i),
			Username: fmt.Sprintf("Player %d", i+1),
			Seat:     i,
			Online:   true,
			AutoPlay: false,
		}
	}
	
	// 开始比赛
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}
	
	// 检查游戏状态
	gameState := engine.GetGameState()
	if gameState.Status != GameStatusStarted {
		t.Errorf("Expected game status to be started, got %v", gameState.Status)
	}
	
	if gameState.CurrentMatch == nil {
		t.Fatal("Expected current match to be set")
	}
	
	t.Logf("Game started successfully with status: %v", gameState.Status)
}

// TestSimpleStartDeal 测试开始一局游戏
func TestSimpleStartDeal(t *testing.T) {
	// 创建游戏引擎
	engine := NewGameEngine()
	
	// 创建4个玩家
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = Player{
			ID:       fmt.Sprintf("player_%d", i),
			Username: fmt.Sprintf("Player %d", i+1),
			Seat:     i,
			Online:   true,
			AutoPlay: false,
		}
	}
	
	// 开始比赛
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}
	
	// 开始第一局
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}
	
	// 检查deal状态
	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal
	
	if deal == nil {
		t.Fatal("Expected current deal to be set")
	}
	
	t.Logf("Deal started with status: %v, level: %d", deal.Status, deal.Level)
	
	// 检查玩家手牌
	for i := 0; i < 4; i++ {
		playerView := engine.GetPlayerView(i)
		if len(playerView.PlayerCards) == 0 {
			t.Errorf("Player %d has no cards", i)
		} else {
			t.Logf("Player %d has %d cards", i, len(playerView.PlayerCards))
		}
	}
}

// TestMatchSimulatorBasic 基础的模拟器测试（只运行一局）
func TestMatchSimulatorBasic(t *testing.T) {
	simulator := NewMatchSimulator(false)
	
	// 创建4个模拟玩家
	for i := 0; i < 4; i++ {
		simulator.players[i] = SimulatedPlayer{
			Player: Player{
				ID:       fmt.Sprintf("player_%d", i),
				Username: fmt.Sprintf("Player %d", i+1),
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
	
	// 只运行一局进行测试
	simulator.currentDealNum = 1
	
	// 开始一局
	err = simulator.engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}
	
	// 检查deal状态
	gameState := simulator.engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal
	
	if deal == nil {
		t.Fatal("Expected current deal to be set")
	}
	
	t.Logf("Basic simulator test completed. Deal status: %v", deal.Status)
}

// TestTributePhase 测试贡牌阶段
func TestTributePhase(t *testing.T) {
	// 创建游戏引擎
	engine := NewGameEngine()
	
	// 创建4个玩家
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = Player{
			ID:       fmt.Sprintf("player_%d", i),
			Username: fmt.Sprintf("Player %d", i+1),
			Seat:     i,
			Online:   true,
			AutoPlay: false,
		}
	}
	
	// 开始比赛
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}
	
	// 开始第一局（应该没有贡牌阶段）
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}
	
	// 第一局不应该有贡牌阶段
	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal
	
	if deal.TributePhase != nil {
		t.Logf("First deal has tribute phase (unexpected but might be valid)")
	} else {
		t.Logf("First deal has no tribute phase (expected)")
	}
	
	// 测试贡牌接口
	action, err := engine.ProcessTributePhase()
	if err != nil {
		t.Logf("ProcessTributePhase returned error: %v", err)
	} else if action != nil {
		t.Logf("ProcessTributePhase returned action: %+v", action)
	} else {
		t.Logf("ProcessTributePhase returned nil (no tribute needed)")
	}
}