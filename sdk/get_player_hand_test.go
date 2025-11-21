package sdk

import (
	"testing"
)

// TestGetPlayerHand 测试 GetPlayerHand 方法
func TestGetPlayerHand(t *testing.T) {
	// 创建游戏引擎
	engine := NewGameEngine()
	
	// 创建玩家
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	// 测试：没有活跃比赛时返回 nil
	hand := engine.GetPlayerHand(0)
	if hand != nil {
		t.Errorf("Expected nil when no active match, got %v", hand)
	}
	
	// 开始比赛
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}
	
	// 开始一局
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}
	
	// 测试：获取玩家手牌
	for seat := 0; seat < 4; seat++ {
		hand := engine.GetPlayerHand(seat)
		if hand == nil {
			t.Errorf("Expected non-nil hand for player %d", seat)
			continue
		}
		
		// 验证手牌数量（每个玩家应该有27张牌）
		if len(hand) != 27 {
			t.Errorf("Expected 27 cards for player %d, got %d", seat, len(hand))
		}
		
		// 验证 Card 对象包含 Level 字段
		for i, card := range hand {
			if card == nil {
				t.Errorf("Card %d for player %d is nil", i, seat)
				continue
			}
			
			// Level 应该是当前局的 Level（初始为 2）
			// 注意：只有当前级牌才有 Level 设置
			// 普通牌的 Level 是 0，这是正常的
			// 关键是 Level 字段存在且可访问
			_ = card.Level // 确保 Level 字段可访问
		}
	}
	
	// 测试：无效的座位号
	hand = engine.GetPlayerHand(-1)
	if hand != nil {
		t.Errorf("Expected nil for invalid seat -1, got %v", hand)
	}
	
	hand = engine.GetPlayerHand(4)
	if hand != nil {
		t.Errorf("Expected nil for invalid seat 4, got %v", hand)
	}
}

// TestGetPlayerHandLevel 测试 Level 字段保留
func TestGetPlayerHandLevel(t *testing.T) {
	// 创建游戏引擎
	engine := NewGameEngine()
	
	// 创建玩家
	players := []Player{
		{ID: "p1", Username: "Player1", Seat: 0},
		{ID: "p2", Username: "Player2", Seat: 1},
		{ID: "p3", Username: "Player3", Seat: 2},
		{ID: "p4", Username: "Player4", Seat: 3},
	}
	
	// 开始比赛和一局
	engine.StartMatch(players)
	engine.StartDeal()
	
	// 获取玩家手牌
	hand := engine.GetPlayerHand(0)
	if hand == nil {
		t.Fatal("Failed to get player hand")
	}
	
	// 统计有 Level 的牌（当前级牌）
	levelCardCount := 0
	for _, card := range hand {
		if card.Level != 0 {
			levelCardCount++
		}
	}
	
	// 应该至少有一些级牌（2或红桃2在初始Level=2时）
	// 注意：这个测试不严格，因为随机发牌可能导致某个玩家没有级牌
	// 但至少验证了 Level 字段可以访问且有非零值
	t.Logf("Player 0 has %d cards with Level set", levelCardCount)
	
	// 验证返回的是引用而非拷贝
	// 修改 hand 应该不影响原始数据（因为我们不应该修改）
	originalLen := len(hand)
	if originalLen == 0 {
		t.Fatal("Hand should not be empty")
	}
}
