package sdk

import (
	"testing"
	"time"
)

// TestPlayActionEdgeCases 测试 PlayAction 的边界情况
func TestPlayActionEdgeCases(t *testing.T) {
	t.Run("空切片 vs nil 切片", func(t *testing.T) {
		now := time.Now()

		// 空切片
		play1 := &PlayAction{
			PlayerSeat: 0,
			Cards:      []*Card{},
			Timestamp:  now,
			IsPass:     true,
		}
		proto1 := ToProtoPlayAction(play1)
		result1 := FromProtoPlayAction(proto1)
		if result1.Cards == nil {
			t.Error("空切片转换后不应为 nil")
		}
		if len(result1.Cards) != 0 {
			t.Errorf("空切片长度应为 0，got %d", len(result1.Cards))
		}

		// nil 切片
		play2 := &PlayAction{
			PlayerSeat: 0,
			Cards:      nil,
			Timestamp:  now,
			IsPass:     true,
		}
		proto2 := ToProtoPlayAction(play2)
		result2 := FromProtoPlayAction(proto2)
		// proto 中 nil 切片和空切片无法区分，都转换为 nil
		if result2.Cards != nil && len(result2.Cards) != 0 {
			t.Errorf("nil 切片转换后应为 nil 或空切片，got %+v", result2.Cards)
		}
	})

	t.Run("Cards 和 IsPass 不一致", func(t *testing.T) {
		now := time.Now()
		card := &Card{Number: 3, Color: "Spade", Name: "3♠"}

		// IsPass=true 但有 Cards（不一致）
		play := &PlayAction{
			PlayerSeat: 0,
			Cards:      []*Card{card},
			IsPass:     true,
			Timestamp:  now,
		}
		proto := ToProtoPlayAction(play)
		result := FromProtoPlayAction(proto)

		// 验证数据保留（不强制一致性，由业务层检查）
		if result.IsPass != true {
			t.Error("IsPass 应保留原值")
		}
		if len(result.Cards) != 1 {
			t.Error("Cards 应保留原值")
		}
	})

	t.Run("时间戳零值", func(t *testing.T) {
		// time.Time{} 零值
		play := &PlayAction{
			PlayerSeat: 0,
			Timestamp:  time.Time{},
			IsPass:     true,
		}
		proto := ToProtoPlayAction(play)
		if proto.TimestampMs != 0 {
			t.Errorf("零值时间戳应转换为 0，got %d", proto.TimestampMs)
		}

		result := FromProtoPlayAction(proto)
		if !result.Timestamp.IsZero() {
			t.Errorf("0 时间戳应转换为零值，got %v", result.Timestamp)
		}
	})

	t.Run("时间戳负值", func(t *testing.T) {
		// 负值时间戳（1970 年之前）
		play := &PlayAction{
			PlayerSeat: 0,
			Timestamp:  time.UnixMilli(-1000),
			IsPass:     true,
		}
		proto := ToProtoPlayAction(play)
		if proto.TimestampMs != -1000 {
			t.Errorf("负值时间戳应保留，got %d", proto.TimestampMs)
		}

		result := FromProtoPlayAction(proto)
		// 根据修复，负值应转换为零值
		if !result.Timestamp.IsZero() {
			t.Errorf("负值时间戳应转换为零值，got %v", result.Timestamp)
		}
	})
}

// TestPlayActionsBatchEdgeCases 测试批量转换的边界情况
func TestPlayActionsBatchEdgeCases(t *testing.T) {
	t.Run("包含 nil 元素", func(t *testing.T) {
		now := time.Now()
		validPlay := &PlayAction{
			PlayerSeat: 0,
			Timestamp:  now,
			IsPass:     true,
		}

		plays := []*PlayAction{validPlay, nil, validPlay}
		proto := ToProtoPlayActions(plays)

		if len(proto) != 3 {
			t.Errorf("应保留长度，got %d", len(proto))
		}
		if proto[0] == nil {
			t.Error("第一个元素不应为 nil")
		}
		if proto[1] != nil {
			t.Error("第二个元素应为 nil（保留索引）")
		}
		if proto[2] == nil {
			t.Error("第三个元素不应为 nil")
		}

		result := FromProtoPlayActions(proto)
		if len(result) != 3 {
			t.Errorf("应保留长度，got %d", len(result))
		}
		if result[0] == nil {
			t.Error("第一个元素不应为 nil")
		}
		if result[1] != nil {
			t.Error("第二个元素应为 nil（保留索引）")
		}
		if result[2] == nil {
			t.Error("第三个元素不应为 nil")
		}
	})

	t.Run("空切片", func(t *testing.T) {
		plays := []*PlayAction{}
		proto := ToProtoPlayActions(plays)
		if proto == nil {
			t.Error("空切片不应转换为 nil")
		}
		if len(proto) != 0 {
			t.Errorf("空切片长度应为 0，got %d", len(proto))
		}
	})
}

// TestTrickEdgeCases 测试 Trick 的边界情况
func TestTrickEdgeCases(t *testing.T) {
	t.Run("Winner=-1（未结束）", func(t *testing.T) {
		trick := &Trick{
			ID:        "trick-1",
			Leader:    0,
			Winner:    -1,
			Status:    TrickStatusPlaying,
			StartTime: time.Now(),
		}
		proto := ToProtoTrick(trick)
		if proto.Winner != -1 {
			t.Errorf("Winner=-1 应保留，got %d", proto.Winner)
		}

		result := FromProtoTrick(proto)
		if result.Winner != -1 {
			t.Errorf("Winner=-1 应保留，got %d", result.Winner)
		}
	})

	t.Run("Plays 为空", func(t *testing.T) {
		trick := &Trick{
			ID:        "trick-1",
			Leader:    0,
			Plays:     []*PlayAction{},
			StartTime: time.Now(),
		}
		proto := ToProtoTrick(trick)
		result := FromProtoTrick(proto)

		if result.Plays == nil {
			t.Error("空 Plays 不应转换为 nil")
		}
		if len(result.Plays) != 0 {
			t.Errorf("空 Plays 长度应为 0，got %d", len(result.Plays))
		}
	})

	t.Run("LeadComp 为 nil", func(t *testing.T) {
		trick := &Trick{
			ID:        "trick-1",
			Leader:    0,
			LeadComp:  nil,
			StartTime: time.Now(),
		}
		proto := ToProtoTrick(trick)
		result := FromProtoTrick(proto)

		if result.LeadComp != nil {
			t.Error("LeadComp=nil 应保留")
		}
	})
}

// TestTributePhaseEdgeCases 测试 TributePhase 的边界情况
func TestTributePhaseEdgeCases(t *testing.T) {
	t.Run("Map 包含 nil Card", func(t *testing.T) {
		card := &Card{Number: 3, Color: "Spade", Name: "3♠"}

		phase := &TributePhase{
			Status: TributeStatusSelecting,
			TributeCards: map[int]*Card{
				0: card,
				1: nil, // nil 值
			},
		}

		proto := ToProtoTributePhase(phase)
		// nil 值应被过滤
		if len(proto.TributeCards) != 1 {
			t.Errorf("nil Card 应被过滤，got %d 张卡", len(proto.TributeCards))
		}
		if proto.TributeCards[1] != nil {
			t.Error("座位 1 的 nil Card 应被过滤")
		}

		result := FromProtoTributePhase(proto)
		if len(result.TributeCards) != 1 {
			t.Errorf("应只有 1 张卡，got %d", len(result.TributeCards))
		}
	})

	t.Run("空 Map vs nil Map", func(t *testing.T) {
		// 空 Map
		phase1 := &TributePhase{
			Status:       TributeStatusWaiting,
			TributeCards: map[int]*Card{},
		}
		proto1 := ToProtoTributePhase(phase1)
		result1 := FromProtoTributePhase(proto1)
		if result1.TributeCards == nil {
			t.Error("空 Map 不应转换为 nil")
		}

		// nil Map
		phase2 := &TributePhase{
			Status:       TributeStatusWaiting,
			TributeCards: nil,
		}
		proto2 := ToProtoTributePhase(phase2)
		result2 := FromProtoTributePhase(proto2)
		// proto 中 nil map 和空 map 无法区分
		if result2.TributeCards != nil && len(result2.TributeCards) != 0 {
			t.Error("nil Map 应转换为 nil 或空 Map")
		}
	})

	t.Run("SelectingPlayer=-1（无人选牌）", func(t *testing.T) {
		phase := &TributePhase{
			Status:          TributeStatusWaiting,
			SelectingPlayer: -1,
		}
		proto := ToProtoTributePhase(phase)
		if proto.SelectingPlayer != -1 {
			t.Errorf("SelectingPlayer=-1 应保留，got %d", proto.SelectingPlayer)
		}

		result := FromProtoTributePhase(proto)
		if result.SelectingPlayer != -1 {
			t.Errorf("SelectingPlayer=-1 应保留，got %d", result.SelectingPlayer)
		}
	})

	t.Run("多个 Map 同时测试", func(t *testing.T) {
		card1 := &Card{Number: 3, Color: "Spade", Name: "3♠"}
		card2 := &Card{Number: 4, Color: "Heart", Name: "4♥"}

		phase := &TributePhase{
			Status: TributeStatusReturning,
			TributeMap: map[int]int{
				0: 1,
				2: 3,
			},
			TributeCards: map[int]*Card{
				0: card1,
				2: nil, // 应被过滤
			},
			ReturnCards: map[int]*Card{
				1: card2,
				3: nil, // 应被过滤
			},
			SelectionResults: map[int]int{
				1: 0,
			},
		}

		proto := ToProtoTributePhase(phase)
		result := FromProtoTributePhase(proto)

		// 验证 TributeMap（int → int，无 nil）
		if len(result.TributeMap) != 2 {
			t.Errorf("TributeMap 应有 2 项，got %d", len(result.TributeMap))
		}

		// 验证 TributeCards（过滤 nil）
		if len(result.TributeCards) != 1 {
			t.Errorf("TributeCards 应有 1 项（过滤了 nil），got %d", len(result.TributeCards))
		}

		// 验证 ReturnCards（过滤 nil）
		if len(result.ReturnCards) != 1 {
			t.Errorf("ReturnCards 应有 1 项（过滤了 nil），got %d", len(result.ReturnCards))
		}

		// 验证 SelectionResults（int → int，无 nil）
		if len(result.SelectionResults) != 1 {
			t.Errorf("SelectionResults 应有 1 项，got %d", len(result.SelectionResults))
		}
	})
}

// TestTimeFromMillisEdgeCases 测试时间转换的边界情况
func TestTimeFromMillisEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		wantZero bool
	}{
		{"正数时间戳", 1609459200000, false},        // 2021-01-01
		{"零值", 0, true},                         // 应返回 zero time
		{"负值", -1000, true},                     // 应返回 zero time
		{"小负值", -1, true},                      // 应返回 zero time
		{"大正值", 9999999999999, false},          // 未来时间
		{"边界正值", 1, false},                    // 1970-01-01 00:00:00.001
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timeFromMillis(tt.input)
			if tt.wantZero && !result.IsZero() {
				t.Errorf("期望零值时间，got %v", result)
			}
			if !tt.wantZero && result.IsZero() {
				t.Errorf("不期望零值时间，got %v", result)
			}
			if !tt.wantZero && result.UnixMilli() != tt.input {
				t.Errorf("时间转换错误: want %d, got %d", tt.input, result.UnixMilli())
			}
		})
	}
}
