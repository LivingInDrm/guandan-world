// proto_adapter_action.go - 游戏动作 Proto 适配器
//
// 职责:
// - PlayAction 类型的双向转换（SDK ↔ Proto）
// - Trick 类型的双向转换（SDK ↔ Proto）
//
// 依赖:
// - proto_adapter_basic.go: ToProtoCards, FromProtoCards, timeFromMillis
// - proto_adapter_comp.go: ToProtoCardComp, FromProtoCardComp
// - proto_adapter_enums.go: ToProtoTrickStatus, FromProtoTrickStatus
//
// 被依赖:
// - 无（叶子节点）
package sdk

import (
	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== PlayAction Adapters ====================

// ToProtoPlayAction 转换 SDK PlayAction 到 Proto PlayAction
// 特殊处理：
// - Timestamp: time.Time → int64 毫秒（零值转换为 0）
// - Cards: 可能为 nil（表示弃牌）
// - Comp: 可能为 nil（表示弃牌）
func ToProtoPlayAction(pa *PlayAction) *pbgame.PlayAction {
	if pa == nil {
		return nil
	}
	
	// 处理时间戳：零值转换为 0
	var timestampMs int64
	if !pa.Timestamp.IsZero() {
		timestampMs = pa.Timestamp.UnixMilli()
	}
	
	return &pbgame.PlayAction{
		PlayerSeat:  int32(pa.PlayerSeat),
		Cards:       ToProtoCards(pa.Cards),
		Comp:        ToProtoCardComp(pa.Comp),
		TimestampMs: timestampMs,
		IsPass:      pa.IsPass,
	}
}

// ToProtoPlayActions 批量转换 SDK PlayActions 到 Proto PlayActions
// 注意: 如果输入包含 nil 元素，输出对应位置也为 nil（保持索引一致性）
func ToProtoPlayActions(plays []*PlayAction) []*pbgame.PlayAction {
	if plays == nil {
		return nil
	}
	result := make([]*pbgame.PlayAction, len(plays))
	for i, play := range plays {
		if play != nil {
			result[i] = ToProtoPlayAction(play)
		}
		// nil 元素保持为 nil，保留索引顺序
	}
	return result
}

// FromProtoPlayAction 转换 Proto PlayAction 到 SDK PlayAction
// 特殊处理：
// - TimestampMs: int64 毫秒 → time.Time
func FromProtoPlayAction(ppa *pbgame.PlayAction) *PlayAction {
	if ppa == nil {
		return nil
	}
	return &PlayAction{
		PlayerSeat: int(ppa.PlayerSeat),
		Cards:      FromProtoCards(ppa.Cards),
		Comp:       FromProtoCardComp(ppa.Comp),
		Timestamp:  timeFromMillis(ppa.TimestampMs),
		IsPass:     ppa.IsPass,
	}
}

// FromProtoPlayActions 批量转换 Proto PlayActions 到 SDK PlayActions
// 注意: 如果输入包含 nil 元素，输出对应位置也为 nil（保持索引一致性）
func FromProtoPlayActions(ppas []*pbgame.PlayAction) []*PlayAction {
	if ppas == nil {
		return nil
	}
	result := make([]*PlayAction, len(ppas))
	for i, ppa := range ppas {
		if ppa != nil {
			result[i] = FromProtoPlayAction(ppa)
		}
		// nil 元素保持为 nil，保留索引顺序
	}
	return result
}

// ==================== Trick Adapters ====================

// ToProtoTrick 转换 SDK Trick 到 Proto Trick
// 特殊处理：
// - StartTime: time.Time → int64 毫秒（零值转换为 0）
// - Plays: 使用 ToProtoPlayActions
// - LeadComp: 可能为 nil
func ToProtoTrick(t *Trick) *pbgame.Trick {
	if t == nil {
		return nil
	}
	
	// 处理时间戳：零值转换为 0
	var startTimeMs int64
	if !t.StartTime.IsZero() {
		startTimeMs = t.StartTime.UnixMilli()
	}
	
	return &pbgame.Trick{
		Id:          t.ID,
		Leader:      int32(t.Leader),
		CurrentTurn: int32(t.CurrentTurn),
		Plays:       ToProtoPlayActions(t.Plays),
		Winner:      int32(t.Winner),
		LeadComp:    ToProtoCardComp(t.LeadComp),
		Status:      ToProtoTrickStatus(t.Status),
		StartTimeMs: startTimeMs,
		NextLeader:  int32(t.NextLeader),
	}
}

// FromProtoTrick 转换 Proto Trick 到 SDK Trick
// 特殊处理：
// - StartTimeMs: int64 毫秒 → time.Time
func FromProtoTrick(pt *pbgame.Trick) *Trick {
	if pt == nil {
		return nil
	}
	return &Trick{
		ID:          pt.Id,
		Leader:      int(pt.Leader),
		CurrentTurn: int(pt.CurrentTurn),
		Plays:       FromProtoPlayActions(pt.Plays),
		Winner:      int(pt.Winner),
		LeadComp:    FromProtoCardComp(pt.LeadComp),
		Status:      FromProtoTrickStatus(pt.Status),
		StartTime:   timeFromMillis(pt.StartTimeMs),
		NextLeader:  int(pt.NextLeader),
	}
}
