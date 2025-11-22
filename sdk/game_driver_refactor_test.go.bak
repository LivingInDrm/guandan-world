package sdk

import (
	"context"
	"testing"
	"time"
)

// TestGameDriver_NoDoubleRegistration 测试重复调用 RunMatch 不会重复注册观察者
func TestGameDriver_NoDoubleRegistration(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// 设置一个简单的输入提供者
	mockProvider := &mockInputProviderQuick{}
	driver.SetInputProvider(mockProvider)

	// 创建一个观察者来计数事件
	eventCount := 0
	eventChan := make(chan *GameEvent, 100)
	observer := &mockObserver{
		onEvent: func(event *GameEvent) {
			if event.Type == EventMatchStarted {
				eventChan <- event
			}
		},
	}
	driver.AddObserver(observer)

	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}

	// 第一次运行（会立即因为 mock provider 而结束）
	go driver.RunMatch(players)
	time.Sleep(50 * time.Millisecond)
	driver.CancelMatch()

	// 收集第一次的事件
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case <-eventChan:
			eventCount++
		case <-timeout:
			goto firstDone
		}
	}
firstDone:

	firstCount := eventCount
	if firstCount != 1 {
		t.Logf("First run event count: %d (expected 1)", firstCount)
	}

	// 第二次运行 - 如果有重复注册，会收到2个事件
	eventCount = 0
	engine2 := NewGameEngine()
	driver2 := NewGameDriver(engine2, nil)
	driver2.SetInputProvider(mockProvider)
	driver2.AddObserver(observer)

	go driver2.RunMatch(players)
	time.Sleep(50 * time.Millisecond)
	driver2.CancelMatch()

	// 收集第二次的事件
	timeout2 := time.After(200 * time.Millisecond)
	for {
		select {
		case <-eventChan:
			eventCount++
		case <-timeout2:
			goto secondDone
		}
	}
secondDone:

	if eventCount != 1 {
		t.Errorf("Second run received %d events, expected 1 (no duplicate registration)", eventCount)
	}
}

// TestGameDriver_AddObserver_NoDuplicates 测试 AddObserver 不会添加重复观察者
func TestGameDriver_AddObserver_NoDuplicates(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	eventCount := 0
	observer := &mockObserver{
		onEvent: func(event *GameEvent) {
			eventCount++
		},
	}

	// 添加同一个观察者3次
	driver.AddObserver(observer)
	driver.AddObserver(observer)
	driver.AddObserver(observer)

	// 触发一个事件
	event := NewMatchStartedEvent(nil)
	driver.notifyObservers(event)

	// 给异步处理一点时间
	time.Sleep(50 * time.Millisecond)

	// 应该只收到1次通知
	if eventCount != 1 {
		t.Errorf("Received %d events, expected 1 (observer should not be duplicated)", eventCount)
	}
}

// TestGameDriver_ObserverPanicRecovery 测试观察者 panic 不会影响其他观察者
func TestGameDriver_ObserverPanicRecovery(t *testing.T) {
	engine := NewGameEngine()
	config := DefaultGameDriverConfig()
	config.AsyncEventHandling = false // 使用同步模式更容易测试
	driver := NewGameDriver(engine, config)

	panicObserverCalled := false
	normalObserverCalled := false

	// 第一个观察者会 panic
	panicObserver := &mockObserver{
		onEvent: func(event *GameEvent) {
			panicObserverCalled = true
			panic("intentional panic for testing")
		},
	}

	// 第二个观察者正常处理
	normalObserver := &mockObserver{
		onEvent: func(event *GameEvent) {
			normalObserverCalled = true
		},
	}

	driver.AddObserver(panicObserver)
	driver.AddObserver(normalObserver)

	// 触发事件
	event := NewMatchStartedEvent(nil)
	driver.notifyObservers(event)

	// 验证两个观察者都被调用了
	if !panicObserverCalled {
		t.Error("Panic observer was not called")
	}
	if !normalObserverCalled {
		t.Error("Normal observer was not called (panic was not recovered)")
	}
}

// TestGameDriver_AsyncObserverPanicRecovery 测试异步模式下的 panic recovery
func TestGameDriver_AsyncObserverPanicRecovery(t *testing.T) {
	engine := NewGameEngine()
	config := DefaultGameDriverConfig()
	config.AsyncEventHandling = true // 使用异步模式
	driver := NewGameDriver(engine, config)

	panicChan := make(chan bool, 1)
	normalChan := make(chan bool, 1)

	// 第一个观察者会 panic
	panicObserver := &mockObserver{
		onEvent: func(event *GameEvent) {
			panicChan <- true
			panic("intentional async panic for testing")
		},
	}

	// 第二个观察者正常处理
	normalObserver := &mockObserver{
		onEvent: func(event *GameEvent) {
			normalChan <- true
		},
	}

	driver.AddObserver(panicObserver)
	driver.AddObserver(normalObserver)

	// 触发事件
	event := NewMatchStartedEvent(nil)
	driver.notifyObservers(event)

	// 等待两个观察者都执行
	timeout := time.After(500 * time.Millisecond)
	panicCalled := false
	normalCalled := false

	for i := 0; i < 2; i++ {
		select {
		case <-panicChan:
			panicCalled = true
		case <-normalChan:
			normalCalled = true
		case <-timeout:
			t.Fatal("Timeout waiting for observers")
		}
	}

	if !panicCalled {
		t.Error("Panic observer was not called")
	}
	if !normalCalled {
		t.Error("Normal observer was not called (async panic was not recovered)")
	}
}

// mockInputProviderQuick 快速返回的 mock provider，用于测试
type mockInputProviderQuick struct{}

func (m *mockInputProviderQuick) RequestPlayDecision(ctx context.Context, playerSeat int, hand []*Card, trickInfo *TrickInfo) (*PlayDecision, error) {
	return &PlayDecision{Action: ActionPass}, nil
}

func (m *mockInputProviderQuick) RequestTributeSelection(ctx context.Context, playerSeat int, options []*Card) (*Card, error) {
	if len(options) > 0 {
		return options[0], nil
	}
	return nil, nil
}

func (m *mockInputProviderQuick) RequestReturnTribute(ctx context.Context, playerSeat int, hand []*Card) (*Card, error) {
	if len(hand) > 0 {
		return hand[0], nil
	}
	return nil, nil
}
