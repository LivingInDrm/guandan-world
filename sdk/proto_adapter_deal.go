// proto_adapter_deal.go - Deal 和 PlayerHand 的 Proto 适配器
//
// 职责:
// - Deal 类型的双向转换（SDK ↔ Proto）
// - PlayerHand 辅助类型的转换
// - PlayerCards [4][]*Card ↔ repeated PlayerHand 的转换
//
// 依赖:
// - proto_adapter_basic.go: ToProtoCard, FromProtoCard, ToProtoCards, FromProtoCards, timeFromMillis
// - proto_adapter_action.go: ToProtoTrick, FromProtoTrick, ToProtoTributePhase, FromProtoTributePhase
// - proto_adapter_result.go: ToProtoDealResult, FromProtoDealResult
//
// 被依赖:
// - proto_adapter_match.go: 使用 ToProtoDeal, FromProtoDeal
package sdk

import (
	"time"

	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== PlayerHand Adapters ====================

// ToProtoPlayerHands 转换 [4][]*Card 到 repeated PlayerHand
// 每个座位的手牌转换为一个 PlayerHand message
func ToProtoPlayerHands(playerCards [4][]*Card) []*pbgame.PlayerHand {
	result := make([]*pbgame.PlayerHand, 4)
	for seat := 0; seat < 4; seat++ {
		result[seat] = &pbgame.PlayerHand{
			PlayerSeat: int32(seat),
			Cards:      ToProtoCards(playerCards[seat]),
		}
	}
	return result
}

// FromProtoPlayerHands 转换 repeated PlayerHand 到 [4][]*Card
// 注意:
// - 预期输入切片长度为4，对应4个玩家座位
// - 通过 PlayerSeat 字段映射到正确的座位位置
// - 缺失或无效的座位将保持为 nil
// - 重复的座位号会被覆盖（以最后出现的为准）
func FromProtoPlayerHands(phs []*pbgame.PlayerHand) [4][]*Card {
	var result [4][]*Card
	// 验证：预期4个 PlayerHand 元素
	if len(phs) != 4 {
		// 数据异常：长度不符合预期，但仍继续处理
		// 调用方应确保传入正确长度的数据
	}
	for _, ph := range phs {
		if ph != nil && ph.PlayerSeat >= 0 && ph.PlayerSeat < 4 {
			result[ph.PlayerSeat] = FromProtoCards(ph.Cards)
		}
	}
	return result
}

// ==================== Deal Adapters ====================

// ToProtoDeal 转换 SDK Deal 到 Proto Deal
// 特殊处理：
// - StartTime: 零值时转换为 0
// - EndTime: nil 时转换为 0
// - LastResult: 可选字段，nil 时不设置
func ToProtoDeal(d *Deal) *pbgame.Deal {
	if d == nil {
		return nil
	}

	var startTimeMs int64
	if !d.StartTime.IsZero() {
		startTimeMs = d.StartTime.UnixMilli()
	}

	var endTimeMs int64
	if d.EndTime != nil {
		endTimeMs = d.EndTime.UnixMilli()
	}

	return &pbgame.Deal{
		Id:           d.ID,
		Level:        int32(d.Level),
		Status:       ToProtoDealStatus(d.Status),
		CurrentTrick: ToProtoTrick(d.CurrentTrick),
		TrickHistory: ToProtoTricks(d.TrickHistory),
		TributePhase: ToProtoTributePhase(d.TributePhase),
		PlayerHands:  ToProtoPlayerHands(d.PlayerCards),
		Rankings:     toInt32Slice(d.Rankings),
		StartTimeMs:  startTimeMs,
		EndTimeMs:    endTimeMs,
		LastResult:   ToProtoDealResult(d.LastResult),
	}
}

// ToProtoDeals 批量转换 SDK Deals 到 Proto Deals
func ToProtoDeals(deals []*Deal) []*pbgame.Deal {
	if deals == nil {
		return nil
	}
	result := make([]*pbgame.Deal, len(deals))
	for i, deal := range deals {
		result[i] = ToProtoDeal(deal)
	}
	return result
}

// FromProtoDeal 转换 Proto Deal 到 SDK Deal
// 特殊处理：
// - StartTimeMs: <= 0 时转换为 time.Time{} 零值（通过 timeFromMillis）
// - EndTimeMs: 0 时转换为 nil
// 注意: StartTimeMs 为 0 表示时间未设置，会被解析为零值时间
func FromProtoDeal(pd *pbgame.Deal) *Deal {
	if pd == nil {
		return nil
	}

	var endTime *time.Time
	if pd.EndTimeMs > 0 {
		t := timeFromMillis(pd.EndTimeMs)
		endTime = &t
	}

	return &Deal{
		ID:           pd.Id,
		Level:        int(pd.Level),
		Status:       FromProtoDealStatus(pd.Status),
		CurrentTrick: FromProtoTrick(pd.CurrentTrick),
		TrickHistory: FromProtoTricks(pd.TrickHistory),
		TributePhase: FromProtoTributePhase(pd.TributePhase),
		PlayerCards:  FromProtoPlayerHands(pd.PlayerHands),
		Rankings:     fromInt32Slice(pd.Rankings),
		StartTime:    timeFromMillis(pd.StartTimeMs),
		EndTime:      endTime,
		LastResult:   FromProtoDealResult(pd.LastResult),
	}
}

// FromProtoDeals 批量转换 Proto Deals 到 SDK Deals
func FromProtoDeals(pds []*pbgame.Deal) []*Deal {
	if pds == nil {
		return nil
	}
	result := make([]*Deal, len(pds))
	for i, pd := range pds {
		result[i] = FromProtoDeal(pd)
	}
	return result
}

// ==================== Trick Batch Adapters ====================

// ToProtoTricks 批量转换 SDK Tricks 到 Proto Tricks
func ToProtoTricks(tricks []*Trick) []*pbgame.Trick {
	if tricks == nil {
		return nil
	}
	result := make([]*pbgame.Trick, len(tricks))
	for i, trick := range tricks {
		result[i] = ToProtoTrick(trick)
	}
	return result
}

// FromProtoTricks 批量转换 Proto Tricks 到 SDK Tricks
func FromProtoTricks(pts []*pbgame.Trick) []*Trick {
	if pts == nil {
		return nil
	}
	result := make([]*Trick, len(pts))
	for i, pt := range pts {
		result[i] = FromProtoTrick(pt)
	}
	return result
}

// ==================== Helper Functions ====================

// toInt32Slice 转换 []int 到 []int32
func toInt32Slice(ints []int) []int32 {
	if ints == nil {
		return nil
	}
	result := make([]int32, len(ints))
	for i, v := range ints {
		result[i] = int32(v)
	}
	return result
}

// fromInt32Slice 转换 []int32 到 []int
func fromInt32Slice(ints []int32) []int {
	if ints == nil {
		return nil
	}
	result := make([]int, len(ints))
	for i, v := range ints {
		result[i] = int(v)
	}
	return result
}
