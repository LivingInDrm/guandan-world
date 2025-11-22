package sdk

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockInputProviderWithTimeout 模拟超时的输入提供者
type mockInputProviderWithTimeout struct {
	playDecisionDelay     time.Duration
	tributeSelectionDelay time.Duration
	returnTributeDelay    time.Duration
	simulateTimeout       bool
	returnError           error // 用于测试非上下文错误
}

func (m *mockInputProviderWithTimeout) RequestPlayDecision(ctx context.Context, playerSeat int, hand []*Card, trickInfo *TrickInfo) (*PlayDecision, error) {
	if m.returnError != nil {
		// 返回非上下文错误用于测试
		return nil, m.returnError
	}
	
	if m.simulateTimeout {
		// 模拟超时：等待直到context超时
		<-ctx.Done()
		return nil, ctx.Err()
	}
	
	if m.playDecisionDelay > 0 {
		select {
		case <-time.After(m.playDecisionDelay):
			// 延迟后返回决策
			return &PlayDecision{Action: ActionPass}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	// 立即返回决策
	return &PlayDecision{Action: ActionPass}, nil
}

func (m *mockInputProviderWithTimeout) RequestTributeSelection(ctx context.Context, playerSeat int, options []*Card) (*Card, error) {
	if m.simulateTimeout {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	
	if m.tributeSelectionDelay > 0 {
		select {
		case <-time.After(m.tributeSelectionDelay):
			if len(options) > 0 {
				return options[0], nil
			}
			return nil, errors.New("no options")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	if len(options) > 0 {
		return options[0], nil
	}
	return nil, errors.New("no options")
}

func (m *mockInputProviderWithTimeout) RequestReturnTribute(ctx context.Context, playerSeat int, hand []*Card) (*Card, error) {
	if m.simulateTimeout {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	
	if m.returnTributeDelay > 0 {
		select {
		case <-time.After(m.returnTributeDelay):
			if len(hand) > 0 {
				return hand[0], nil
			}
			return nil, errors.New("no cards")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	if len(hand) > 0 {
		return hand[0], nil
	}
	return nil, errors.New("no cards")
}

// TestHandleTimeout tests the handleTimeout method
func TestHandleTimeout(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// 添加观察者来捕获事件
	var receivedEvent *GameEvent
	observer := &mockObserver{
		onEvent: func(event *GameEvent) {
			if event.Type == EventPlayerTimeout {
				receivedEvent = event
			}
		},
	}
	driver.AddObserver(observer)

	// Test play decision timeout
	driver.handleTimeout(0, "play_decision")
	
	if receivedEvent == nil {
		t.Fatal("Expected timeout event to be sent")
	}
	
	if receivedEvent.Type != EventPlayerTimeout {
		t.Errorf("Event type = %v, want %v", receivedEvent.Type, EventPlayerTimeout)
	}
	
	if receivedEvent.PlayerSeat != 0 {
		t.Errorf("PlayerSeat = %d, want 0", receivedEvent.PlayerSeat)
	}

	// Verify event data uses "action" key (not "action_type")
	data, ok := receivedEvent.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Event data is not map[string]interface{}")
	}
	
	if action, exists := data["action"]; !exists {
		t.Error("Event data should have 'action' key")
	} else if action != "play_decision" {
		t.Errorf("Event data action = %v, want 'play_decision'", action)
	}
	
	// Verify no duplicate player_seat in data
	if _, exists := data["player_seat"]; exists {
		t.Error("Event data should not duplicate player_seat (already in PlayerSeat field)")
	}

	stats := driver.GetTimeoutStats()
	if stats[0].PlayDecisionTimeouts != 1 {
		t.Errorf("PlayDecisionTimeouts = %d, want 1", stats[0].PlayDecisionTimeouts)
	}
	if stats[0].TotalTimeouts != 1 {
		t.Errorf("TotalTimeouts = %d, want 1", stats[0].TotalTimeouts)
	}

	// Test tribute timeout
	receivedEvent = nil
	driver.handleTimeout(1, "tribute_select")
	
	if receivedEvent == nil {
		t.Fatal("Expected timeout event to be sent")
	}
	
	stats = driver.GetTimeoutStats()
	if stats[1].TributeTimeouts != 1 {
		t.Errorf("TributeTimeouts = %d, want 1", stats[1].TributeTimeouts)
	}
}

// mockObserver 模拟观察者
type mockObserver struct {
	onEvent func(*GameEvent)
}

func (m *mockObserver) OnGameEvent(event *GameEvent) {
	if m.onEvent != nil {
		m.onEvent(event)
	}
}

// TestRunMatch_InitializesGameCancelContext tests that RunMatch initializes gameCancelCtx
func TestRunMatch_InitializesGameCancelContext(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)
	
	// 使用快速失败的输入提供者
	mockProvider := &mockInputProviderWithTimeout{
		simulateTimeout: false, // 不模拟超时，快速返回
	}
	driver.SetInputProvider(mockProvider)

	// 创建4个测试玩家
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}

	// 在goroutine中运行比赛，因为可能会持续较长时间
	done := make(chan bool)
	go func() {
		_, _ = driver.RunMatch(players)
		done <- true
	}()

	// 等待一小段时间让RunMatch启动
	time.Sleep(10 * time.Millisecond)

	// 验证gameCancelCtx已初始化
	if driver.gameCancelCtx == nil {
		t.Error("gameCancelCtx should be initialized during RunMatch")
	}

	// 调用CancelMatch来结束比赛
	driver.CancelMatch()

	// 等待比赛结束
	select {
	case <-done:
		// Match ended, which is expected after cancel
	case <-time.After(2 * time.Second):
		t.Fatal("RunMatch did not finish after cancellation")
	}
}

// TestRunMatch_CleansUpContext tests that RunMatch calls cancelFunc on exit
func TestRunMatch_CleansUpContext(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)
	
	mockProvider := &mockInputProviderWithTimeout{}
	driver.SetInputProvider(mockProvider)

	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}

	// Run match in goroutine
	done := make(chan bool)
	go func() {
		driver.RunMatch(players)
		done <- true
	}()

	// Wait a bit for initialization
	time.Sleep(10 * time.Millisecond)

	// Get the context
	ctx := driver.gameCancelCtx
	if ctx == nil {
		t.Fatal("gameCancelCtx not initialized")
	}

	// Cancel the match using thread-safe method
	driver.CancelMatch()

	// Wait for match to end
	select {
	case <-done:
		// Verify context was cancelled
		select {
		case <-ctx.Done():
			// Context is cancelled, as expected
		default:
			t.Error("Context should be cancelled after cancelFunc is called")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunMatch did not finish after cancellation")
	}
}

// TestTimeoutStrategy_Integration tests timeout strategy integration
func TestTimeoutStrategy_Integration(t *testing.T) {
	engine := NewGameEngine()
	
	config := DefaultGameDriverConfig()
	config.PlayDecisionTimeout = 50 * time.Millisecond // Very short timeout for testing
	
	driver := NewGameDriver(engine, config)
	
	// Use provider that will timeout
	mockProvider := &mockInputProviderWithTimeout{
		simulateTimeout: true, // This will cause timeout
	}
	driver.SetInputProvider(mockProvider)

	// Track timeout events
	var timeoutCount int
	observer := &mockObserver{
		onEvent: func(event *GameEvent) {
			if event.Type == EventPlayerTimeout {
				timeoutCount++
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

	// Run match for a short time
	done := make(chan bool)
	go func() {
		driver.RunMatch(players)
		done <- true
	}()

	// Let it run for a bit to trigger some timeouts
	time.Sleep(200 * time.Millisecond)

	// Cancel the match using thread-safe method
	driver.CancelMatch()

	// Wait for completion
	select {
	case <-done:
		// Verify that timeouts were handled
		if timeoutCount == 0 {
			t.Error("Expected some timeout events, but got none")
		}
		
		// Verify stats were recorded
		stats := driver.GetTimeoutStats()
		totalTimeouts := 0
		for _, s := range stats {
			totalTimeouts += s.TotalTimeouts
		}
		
		if totalTimeouts == 0 {
			t.Error("Expected timeout stats to be recorded")
		}
		
		if totalTimeouts != timeoutCount {
			t.Errorf("Timeout stats (%d) doesn't match timeout events (%d)", totalTimeouts, timeoutCount)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout")
	}
}
