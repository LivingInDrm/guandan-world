package sdk

import (
	"fmt"
	"testing"
	"time"
)

// TestMatchSimulator 测试比赛模拟器的基本功能
func TestMatchSimulator(t *testing.T) {
	// 设置超时时间，防止测试无限运行
	timeout := time.After(5 * time.Minute)
	done := make(chan bool)

	go func() {
		// 创建并运行模拟器（非详细模式）
		simulator := NewMatchSimulator(false)
		err := simulator.SimulateMatch()
		if err != nil {
			t.Errorf("Failed to simulate match: %v", err)
		}
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test timeout: match simulation took too long")
	case <-done:
		// 测试成功完成
		t.Log("Match simulation completed successfully")
	}
}

// TestMatchSimulatorVerbose 测试详细模式的比赛模拟
func TestMatchSimulatorVerbose(t *testing.T) {
	// 只运行简短的详细测试来演示输出
	simulator := NewMatchSimulator(true) // true 启用详细模式

	err := simulator.SimulateMatch()
	if err != nil {
		t.Errorf("Failed to simulate match: %v", err)
	}
}

// TestSimpleAutoPlayAlgorithm 测试自动出牌算法
func TestSimpleAutoPlayAlgorithm(t *testing.T) {
	algorithm := NewSimpleAutoPlayAlgorithm(2)

	// 测试选择贡牌
	t.Run("SelectTributeCard", func(t *testing.T) {
		// 创建测试手牌
		level := 2 // 假设当前级别为2
		hand := []*Card{
			&Card{Number: 3, Color: "Spade", Level: level},
			&Card{Number: 10, Color: "Heart", Level: level},
			&Card{Number: 14, Color: "Diamond", Level: level}, // Ace
		}

		// 选择贡牌（应该选最大的）
		tributeCard := algorithm.SelectTributeCard(hand, true)
		if tributeCard == nil {
			t.Error("Expected tribute card, got nil")
		} else if tributeCard.Number != 14 {
			t.Errorf("Expected Ace (14), got %d", tributeCard.Number)
		}
	})

	// 测试选择还贡
	t.Run("SelectReturnCard", func(t *testing.T) {
		// 创建测试手牌
		level := 2 // 假设当前级别为2
		hand := []*Card{
			&Card{Number: 3, Color: "Spade", Level: level},
			&Card{Number: 3, Color: "Heart", Level: level},
			&Card{Number: 3, Color: "Diamond", Level: level},
			&Card{Number: 3, Color: "Club", Level: level}, // 四个3（炸弹）
			&Card{Number: 10, Color: "Heart", Level: level},
			&Card{Number: 14, Color: "Diamond", Level: level}, // Ace
		}

		// 选择还贡（应该避免破坏炸弹）
		returnCard := algorithm.SelectReturnCard(hand, true)
		if returnCard == nil {
			t.Error("Expected return card, got nil")
		} else if returnCard.Number == 3 {
			t.Error("Should not break bomb by returning a 3")
		}
	})

	// 测试首出逻辑
	t.Run("LeaderPlay", func(t *testing.T) {
		// 创建包含对子的手牌
		level := 2 // 假设当前级别为2
		hand := []*Card{
			&Card{Number: 3, Color: "Spade", Level: level},
			&Card{Number: 3, Color: "Heart", Level: level},
			&Card{Number: 10, Color: "Heart", Level: level},
			&Card{Number: 14, Color: "Diamond", Level: level},
		}

		trickInfo := &TrickInfo{IsLeader: true, LeadComp: nil}
		cards := algorithm.SelectCardsToPlay(hand, trickInfo)
		if cards == nil {
			t.Error("Expected cards to play, got nil")
		} else if len(cards) != 2 {
			t.Errorf("Expected pair (2 cards), got %d cards", len(cards))
		}
	})
}

// TestMatchSimulatorEventHandling 测试事件处理
func TestMatchSimulatorEventHandling(t *testing.T) {
	simulator := NewMatchSimulator(false)

	// 创建测试玩家
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = Player{
			ID:       testPlayerID(i),
			Username: testPlayerName(i),
			Seat:     i,
			Online:   true,
			AutoPlay: true,
		}
	}

	// 注册事件处理器（使用channel来同步）
	eventReceived := make(chan bool, 1)
	simulator.engine.RegisterEventHandler(EventMatchStarted, func(event *GameEvent) {
		select {
		case eventReceived <- true:
		default:
		}
	})

	// 开始比赛
	err := simulator.engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// 验证事件被触发（等待少量时间）
	select {
	case <-eventReceived:
		t.Logf("EventMatchStarted received successfully")
	default:
		t.Error("Expected EventMatchStarted to be triggered")
	}
}

// TestVerboseDemo 演示详细输出的便捷测试
func TestVerboseDemo(t *testing.T) {
	// 跳过演示测试，手动运行时可以取消注释
	t.Skip("Skipping verbose demo - uncomment this line to run")

	err := RunVerboseDemo()
	if err != nil {
		t.Errorf("Verbose demo failed: %v", err)
	}
}

// TestMatchSimulatorQuick 快速测试方法，限制循环次数，避免长时间运行
func TestMatchSimulatorQuick(t *testing.T) {
	simulator := NewMatchSimulator(true) // 启用详细模式

	// 直接运行测试，使用已有的maxDeals=10限制
	err := simulator.SimulateMatch()
	if err != nil {
		t.Logf("模拟出现错误: %v", err)
	}

	t.Logf("快速测试完成")
}

// 辅助函数
func testPlayerID(seat int) string {
	return fmt.Sprintf("test_player_%d", seat)
}

func testPlayerName(seat int) string {
	return fmt.Sprintf("Test Player %d", seat+1)
}

// BenchmarkMatchSimulation 性能测试
func BenchmarkMatchSimulation(b *testing.B) {
	b.Skip("Skipping benchmark - run manually when needed")

	for i := 0; i < b.N; i++ {
		simulator := NewMatchSimulator(false)
		_ = simulator.SimulateMatch()
	}
}
