package sdk

import (
	"testing"
	"time"
)

// TestMatchSimulator 测试匹配模拟器
func TestMatchSimulator(t *testing.T) {
	simulator := NewMatchSimulator(true) // 启用详细输出

	result, err := simulator.SimulateMatch()
	if err != nil {
		t.Fatalf("模拟比赛失败: %v", err)
	}

	if result == nil {
		t.Fatal("比赛结果为空")
	}

	t.Logf("比赛成功完成，获胜队伍: %d", result.Winner)
	t.Logf("最终等级: 队伍0=%d级, 队伍1=%d级", result.FinalLevels[0], result.FinalLevels[1])
	t.Logf("比赛时长: %v", result.Duration)

	if result.Statistics != nil {
		t.Logf("总局数: %d", result.Statistics.TotalDeals)
	}
}

// TestMatchSimulatorQuiet 测试匹配模拟器（安静模式）
func TestMatchSimulatorQuiet(t *testing.T) {
	simulator := NewMatchSimulator(false) // 关闭详细输出

	startTime := time.Now()
	result, err := simulator.SimulateMatch()
	duration := time.Since(startTime)

	if err != nil {
		t.Fatalf("模拟比赛失败: %v", err)
	}

	if result == nil {
		t.Fatal("比赛结果为空")
	}

	t.Logf("安静模式比赛完成，耗时: %v", duration)
	t.Logf("获胜队伍: %d", result.Winner)
	t.Logf("最终等级: 队伍0=%d级, 队伍1=%d级", result.FinalLevels[0], result.FinalLevels[1])
}

// TestMultipleMatches 测试多次模拟比赛
func TestMultipleMatches(t *testing.T) {
	const numMatches = 5

	results := make([]*MatchResult, numMatches)
	durations := make([]time.Duration, numMatches)

	for i := 0; i < numMatches; i++ {
		simulator := NewMatchSimulator(false) // 关闭详细输出以加快速度

		startTime := time.Now()
		result, err := simulator.SimulateMatch()
		durations[i] = time.Since(startTime)

		if err != nil {
			t.Fatalf("第%d场比赛模拟失败: %v", i+1, err)
		}

		results[i] = result
		t.Logf("第%d场比赛: 获胜队伍=%d, 等级=%d-%d, 耗时=%v",
			i+1, result.Winner, result.FinalLevels[0], result.FinalLevels[1], durations[i])
	}

	// 统计结果
	team0Wins := 0
	team1Wins := 0
	totalDuration := time.Duration(0)

	for i, result := range results {
		if result.Winner == 0 {
			team0Wins++
		} else {
			team1Wins++
		}
		totalDuration += durations[i]
	}

	t.Logf("=== 多场比赛统计 ===")
	t.Logf("总场数: %d", numMatches)
	t.Logf("队伍0获胜: %d场", team0Wins)
	t.Logf("队伍1获胜: %d场", team1Wins)
	t.Logf("平均耗时: %v", totalDuration/time.Duration(numMatches))
	t.Logf("总耗时: %v", totalDuration)
}

// BenchmarkMatchSimulator 性能测试
func BenchmarkMatchSimulator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulator := NewMatchSimulator(false) // 关闭详细输出

		_, err := simulator.SimulateMatch()
		if err != nil {
			b.Fatalf("模拟比赛失败: %v", err)
		}
	}
}
