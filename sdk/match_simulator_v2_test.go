package sdk

import (
	"testing"
	"time"
)

// TestMatchSimulatorV2 测试新架构的比赛模拟器
func TestMatchSimulatorV2(t *testing.T) {
	// 设置超时时间
	timeout := time.After(5 * time.Minute)
	done := make(chan bool)

	go func() {
		simulator := NewMatchSimulatorV2(false)
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
		t.Log("Match simulation V2 completed successfully")
	}
}

// TestMatchSimulatorV2Verbose 测试详细模式
func TestMatchSimulatorV2Verbose(t *testing.T) {
	simulator := NewMatchSimulatorV2(true)
	err := simulator.SimulateMatch()
	if err != nil {
		t.Errorf("Failed to simulate match in verbose mode: %v", err)
	}

	// 检查事件日志
	eventLog := simulator.GetEventLog()
	t.Logf("Generated %d events", len(eventLog))

	if len(eventLog) == 0 {
		t.Error("Expected some events to be logged")
	}
}

// TestGameDriverBasic 测试GameDriver基础功能
func TestGameDriverBasic(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, DefaultGameDriverConfig())

	// 测试基础设置
	if driver.GetEngine() != engine {
		t.Error("Engine mismatch")
	}

	config := driver.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	// 测试输入提供者设置
	inputProvider := NewSimulatingInputProvider()
	driver.SetInputProvider(inputProvider)

	// 创建玩家
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}

	// 设置算法
	algorithms := make([]AutoPlayAlgorithm, 4)
	for i := 0; i < 4; i++ {
		algorithms[i] = NewSimpleAutoPlayAlgorithm(2)
	}
	err := inputProvider.BatchSetAlgorithms(algorithms)
	if err != nil {
		t.Fatalf("Failed to set algorithms: %v", err)
	}

	// 运行比赛（简化版，只运行一局）
	err = engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// 启动一局
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	t.Log("GameDriver basic test completed")
}

// TestSimulatingInputProvider 测试模拟输入提供者
func TestSimulatingInputProvider(t *testing.T) {
	provider := NewSimulatingInputProvider()

	// 设置算法
	algorithm := NewSimpleAutoPlayAlgorithm(2)
	provider.SetPlayerAlgorithm(0, algorithm)

	// 检查算法设置
	if !provider.HasAlgorithmForPlayer(0) {
		t.Error("Algorithm should be set for player 0")
	}

	if provider.HasAlgorithmForPlayer(1) {
		t.Error("Algorithm should not be set for player 1")
	}

	// 测试批量设置
	algorithms := make([]AutoPlayAlgorithm, 4)
	for i := 0; i < 4; i++ {
		algorithms[i] = NewSimpleAutoPlayAlgorithm(2)
	}

	err := provider.BatchSetAlgorithms(algorithms)
	if err != nil {
		t.Fatalf("Failed to batch set algorithms: %v", err)
	}

	// 检查所有玩家都有算法
	for i := 0; i < 4; i++ {
		if !provider.HasAlgorithmForPlayer(i) {
			t.Errorf("Algorithm should be set for player %d", i)
		}
	}

	t.Log("SimulatingInputProvider test completed")
}

// TestMatchSimulatorObserver 测试事件观察者
func TestMatchSimulatorObserver(t *testing.T) {
	engine := NewGameEngine()

	eventCount := 0
	logger := func(message string) {
		eventCount++
		t.Logf("Event: %s", message)
	}

	observer := NewMatchSimulatorObserver(engine, false, logger)

	// 模拟一些事件
	observer.OnGameEvent(&GameEvent{
		Type: EventMatchStarted,
		Data: nil,
	})

	observer.OnGameEvent(&GameEvent{
		Type: EventDealStarted,
		Data: nil,
	})

	if eventCount != 2 {
		t.Errorf("Expected 2 events, got %d", eventCount)
	}

	t.Log("MatchSimulatorObserver test completed")
}

// BenchmarkMatchSimulatorV2 性能测试
func BenchmarkMatchSimulatorV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simulator := NewMatchSimulatorV2(false) // 非详细模式
		err := simulator.SimulateMatch()
		if err != nil {
			b.Fatalf("Simulation failed: %v", err)
		}
	}
}

// TestArchitectureComparison 架构对比测试
func TestArchitectureComparison(t *testing.T) {
	t.Log("Testing old architecture...")
	startTime := time.Now()
	err := RunMatchSimulation(false)
	oldDuration := time.Since(startTime)
	if err != nil {
		t.Logf("Old architecture failed: %v", err)
	} else {
		t.Logf("Old architecture completed in: %v", oldDuration)
	}

	t.Log("Testing new architecture...")
	startTime = time.Now()
	err = RunMatchSimulationV2(false)
	newDuration := time.Since(startTime)
	if err != nil {
		t.Errorf("New architecture failed: %v", err)
	} else {
		t.Logf("New architecture completed in: %v", newDuration)
	}

	if oldDuration > 0 && newDuration > 0 {
		improvement := float64(oldDuration-newDuration) / float64(oldDuration) * 100
		if improvement > 0 {
			t.Logf("New architecture is %.1f%% faster", improvement)
		} else {
			t.Logf("Old architecture is %.1f%% faster", -improvement)
		}
	}
}
