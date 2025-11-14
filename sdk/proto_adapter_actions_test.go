package sdk

import (
	"testing"
	"time"
)

// TestPlayActionRoundTrip 测试 PlayAction 的序列化/反序列化
func TestPlayActionRoundTrip(t *testing.T) {
	now := time.Now()
	card := &Card{
		Number:    3,
		RawNumber: 3,
		Color:     "Spade",
		Level:     2,
		Name:      "3♠",
		DeckIndex: 0,
	}

	tests := []struct {
		name string
		play *PlayAction
	}{
		{
			name: "正常出牌",
			play: &PlayAction{
				PlayerSeat: 0,
				Cards:      []*Card{card},
				Comp:       NewSingle([]*Card{card}),
				Timestamp:  now,
				IsPass:     false,
			},
		},
		{
			name: "弃牌",
			play: &PlayAction{
				PlayerSeat: 1,
				Cards:      nil,
				Comp:       &Fold{BaseComp: BaseComp{Valid: true, Type: TypeFold}},
				Timestamp:  now,
				IsPass:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SDK → Proto
			proto := ToProtoPlayAction(tt.play)
			if proto == nil {
				t.Fatal("ToProtoPlayAction 返回 nil")
			}

			// Proto → SDK
			result := FromProtoPlayAction(proto)
			if result == nil {
				t.Fatal("FromProtoPlayAction 返回 nil")
			}

			// 验证字段
			if result.PlayerSeat != tt.play.PlayerSeat {
				t.Errorf("PlayerSeat 不匹配: got %d, want %d", result.PlayerSeat, tt.play.PlayerSeat)
			}
			if result.IsPass != tt.play.IsPass {
				t.Errorf("IsPass 不匹配: got %v, want %v", result.IsPass, tt.play.IsPass)
			}
			// 时间戳精度到毫秒
			if result.Timestamp.UnixMilli() != tt.play.Timestamp.UnixMilli() {
				t.Errorf("Timestamp 不匹配: got %v, want %v", result.Timestamp, tt.play.Timestamp)
			}

			// 验证 Cards
			if len(result.Cards) != len(tt.play.Cards) {
				t.Errorf("Cards 数量不匹配: got %d, want %d", len(result.Cards), len(tt.play.Cards))
			}
			for i := range tt.play.Cards {
				if i >= len(result.Cards) {
					break
				}
				if result.Cards[i].Number != tt.play.Cards[i].Number ||
					result.Cards[i].Color != tt.play.Cards[i].Color ||
					result.Cards[i].Name != tt.play.Cards[i].Name {
					t.Errorf("Cards[%d] 不匹配: got %+v, want %+v", i, result.Cards[i], tt.play.Cards[i])
				}
			}

			// 验证 Comp
			if result.Comp == nil && tt.play.Comp != nil {
				t.Error("Comp 不应为 nil")
			}
			if result.Comp != nil && tt.play.Comp != nil {
				if result.Comp.GetType() != tt.play.Comp.GetType() {
					t.Errorf("Comp 类型不匹配: got %v, want %v", result.Comp.GetType(), tt.play.Comp.GetType())
				}
				if result.Comp.IsValid() != tt.play.Comp.IsValid() {
					t.Errorf("Comp 有效性不匹配: got %v, want %v", result.Comp.IsValid(), tt.play.Comp.IsValid())
				}
			}
		})
	}
}

// TestTrickRoundTrip 测试 Trick 的序列化/反序列化
func TestTrickRoundTrip(t *testing.T) {
	now := time.Now()
	card := &Card{
		Number:    3,
		RawNumber: 3,
		Color:     "Spade",
		Level:     2,
		Name:      "3♠",
		DeckIndex: 0,
	}

	play := &PlayAction{
		PlayerSeat: 0,
		Cards:      []*Card{card},
		Comp:       NewSingle([]*Card{card}),
		Timestamp:  now,
		IsPass:     false,
	}

	trick := &Trick{
		ID:          "trick-1",
		Leader:      0,
		CurrentTurn: 1,
		Plays:       []*PlayAction{play},
		Winner:      -1,
		LeadComp:    NewSingle([]*Card{card}),
		Status:      TrickStatusPlaying,
		StartTime:   now,
		NextLeader:  0,
	}

	// SDK → Proto
	proto := ToProtoTrick(trick)
	if proto == nil {
		t.Fatal("ToProtoTrick 返回 nil")
	}

	// Proto → SDK
	result := FromProtoTrick(proto)
	if result == nil {
		t.Fatal("FromProtoTrick 返回 nil")
	}

	// 验证字段
	if result.ID != trick.ID {
		t.Errorf("ID 不匹配: got %s, want %s", result.ID, trick.ID)
	}
	if result.Leader != trick.Leader {
		t.Errorf("Leader 不匹配: got %d, want %d", result.Leader, trick.Leader)
	}
	if result.CurrentTurn != trick.CurrentTurn {
		t.Errorf("CurrentTurn 不匹配: got %d, want %d", result.CurrentTurn, trick.CurrentTurn)
	}
	if result.Winner != trick.Winner {
		t.Errorf("Winner 不匹配: got %d, want %d", result.Winner, trick.Winner)
	}
	if result.Status != trick.Status {
		t.Errorf("Status 不匹配: got %s, want %s", result.Status, trick.Status)
	}
	if len(result.Plays) != len(trick.Plays) {
		t.Errorf("Plays 数量不匹配: got %d, want %d", len(result.Plays), len(trick.Plays))
	}
	// 时间戳精度到毫秒
	if result.StartTime.UnixMilli() != trick.StartTime.UnixMilli() {
		t.Errorf("StartTime 不匹配: got %v, want %v", result.StartTime, trick.StartTime)
	}

	// 验证 LeadComp
	if result.LeadComp == nil && trick.LeadComp != nil {
		t.Error("LeadComp 不应为 nil")
	}
	if result.LeadComp != nil && trick.LeadComp != nil {
		if result.LeadComp.GetType() != trick.LeadComp.GetType() {
			t.Errorf("LeadComp 类型不匹配: got %v, want %v", result.LeadComp.GetType(), trick.LeadComp.GetType())
		}
		if result.LeadComp.IsValid() != trick.LeadComp.IsValid() {
			t.Errorf("LeadComp 有效性不匹配: got %v, want %v", result.LeadComp.IsValid(), trick.LeadComp.IsValid())
		}
	}

	// 验证 Plays 内容
	for i := range trick.Plays {
		if i >= len(result.Plays) {
			break
		}
		if result.Plays[i].PlayerSeat != trick.Plays[i].PlayerSeat {
			t.Errorf("Plays[%d].PlayerSeat 不匹配: got %d, want %d", i, result.Plays[i].PlayerSeat, trick.Plays[i].PlayerSeat)
		}
		if result.Plays[i].IsPass != trick.Plays[i].IsPass {
			t.Errorf("Plays[%d].IsPass 不匹配: got %v, want %v", i, result.Plays[i].IsPass, trick.Plays[i].IsPass)
		}
	}
}

// TestTributePhaseRoundTrip 测试 TributePhase 的序列化/反序列化
func TestTributePhaseRoundTrip(t *testing.T) {
	card1 := &Card{
		Number:    3,
		RawNumber: 3,
		Color:     "Spade",
		Level:     2,
		Name:      "3♠",
		DeckIndex: 0,
	}
	card2 := &Card{
		Number:    4,
		RawNumber: 4,
		Color:     "Heart",
		Level:     2,
		Name:      "4♥",
		DeckIndex: 1,
	}

	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributeMap: map[int]int{
			0: 1,
			2: 3,
		},
		TributeCards: map[int]*Card{
			0: card1,
			2: card2,
		},
		ReturnCards: map[int]*Card{
			1: card2,
		},
		PoolCards:       []*Card{card1, card2},
		SelectingPlayer: 1,
		IsImmune:        false,
		SelectionResults: map[int]int{
			1: 0,
		},
	}

	// SDK → Proto
	proto := ToProtoTributePhase(phase)
	if proto == nil {
		t.Fatal("ToProtoTributePhase 返回 nil")
	}

	// Proto → SDK
	result := FromProtoTributePhase(proto)
	if result == nil {
		t.Fatal("FromProtoTributePhase 返回 nil")
	}

	// 验证字段
	if result.Status != phase.Status {
		t.Errorf("Status 不匹配: got %s, want %s", result.Status, phase.Status)
	}
	if result.SelectingPlayer != phase.SelectingPlayer {
		t.Errorf("SelectingPlayer 不匹配: got %d, want %d", result.SelectingPlayer, phase.SelectingPlayer)
	}
	if result.IsImmune != phase.IsImmune {
		t.Errorf("IsImmune 不匹配: got %v, want %v", result.IsImmune, phase.IsImmune)
	}

	// 验证 map 字段
	if len(result.TributeMap) != len(phase.TributeMap) {
		t.Errorf("TributeMap 大小不匹配: got %d, want %d", len(result.TributeMap), len(phase.TributeMap))
	}
	for k, v := range phase.TributeMap {
		if result.TributeMap[k] != v {
			t.Errorf("TributeMap[%d] 不匹配: got %d, want %d", k, result.TributeMap[k], v)
		}
	}

	if len(result.TributeCards) != len(phase.TributeCards) {
		t.Errorf("TributeCards 大小不匹配: got %d, want %d", len(result.TributeCards), len(phase.TributeCards))
	}

	if len(result.ReturnCards) != len(phase.ReturnCards) {
		t.Errorf("ReturnCards 大小不匹配: got %d, want %d", len(result.ReturnCards), len(phase.ReturnCards))
	}

	if len(result.SelectionResults) != len(phase.SelectionResults) {
		t.Errorf("SelectionResults 大小不匹配: got %d, want %d", len(result.SelectionResults), len(phase.SelectionResults))
	}

	if len(result.PoolCards) != len(phase.PoolCards) {
		t.Errorf("PoolCards 数量不匹配: got %d, want %d", len(result.PoolCards), len(phase.PoolCards))
	}
}

// TestNilHandling 测试 nil 值处理
func TestNilHandling(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{
			name: "ToProtoPlayAction(nil)",
			fn: func(t *testing.T) {
				result := ToProtoPlayAction(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
		{
			name: "FromProtoPlayAction(nil)",
			fn: func(t *testing.T) {
				result := FromProtoPlayAction(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
		{
			name: "ToProtoTrick(nil)",
			fn: func(t *testing.T) {
				result := ToProtoTrick(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
		{
			name: "FromProtoTrick(nil)",
			fn: func(t *testing.T) {
				result := FromProtoTrick(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
		{
			name: "ToProtoTributePhase(nil)",
			fn: func(t *testing.T) {
				result := ToProtoTributePhase(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
		{
			name: "FromProtoTributePhase(nil)",
			fn: func(t *testing.T) {
				result := FromProtoTributePhase(nil)
				if result != nil {
					t.Error("期望返回 nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}
