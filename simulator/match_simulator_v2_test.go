package simulator

import (
	"testing"
	"time"

	"guandan-world/ai"
	"guandan-world/sdk"
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
	engine := sdk.NewGameEngine()
	driver := sdk.NewGameDriver(engine, sdk.DefaultGameDriverConfig())

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

	// 设置算法
	algorithms := make([]ai.AutoPlayAlgorithm, 4)
	for i := 0; i < 4; i++ {
		algorithms[i] = ai.NewSimpleAutoPlayAlgorithm(2)
	}
	err := inputProvider.BatchSetAlgorithms(algorithms)
	if err != nil {
		t.Fatalf("Failed to set algorithms: %v", err)
	}

	t.Log("GameDriver basic test completed")
}

// TestSimulatingInputProvider 测试模拟输入提供者
func TestSimulatingInputProvider(t *testing.T) {
	provider := NewSimulatingInputProvider()

	// 设置算法
	algorithm := ai.NewSimpleAutoPlayAlgorithm(2)
	provider.SetPlayerAlgorithm(0, algorithm)

	// 测试批量设置
	algorithms := make([]ai.AutoPlayAlgorithm, 4)
	for i := 0; i < 4; i++ {
		algorithms[i] = ai.NewSimpleAutoPlayAlgorithm(2)
	}

	err := provider.BatchSetAlgorithms(algorithms)
	if err != nil {
		t.Fatalf("Failed to batch set algorithms: %v", err)
	}

	t.Log("SimulatingInputProvider test completed")
}

// TestMatchSimulatorObserver 测试事件观察者
func TestMatchSimulatorObserver(t *testing.T) {
	engine := sdk.NewGameEngine()

	eventCount := 0
	logger := func(message string) {
		eventCount++
		t.Logf("Event: %s", message)
	}

	observer := NewMatchSimulatorObserver(engine, false, logger)

	// 模拟一些事件
	observer.OnGameEvent(&sdk.GameEvent{
		Type: sdk.EventMatchStarted,
		Data: nil,
	})

	observer.OnGameEvent(&sdk.GameEvent{
		Type: sdk.EventDealStarted,
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
