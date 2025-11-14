// proto_adapter_match.go - Match 和 PlayerView 的 Proto 适配器
//
// 职责:
// - Match 类型的双向转换（SDK ↔ Proto）
// - PlayerView 类型的双向转换（SDK ↔ Proto）
//
// 依赖:
// - proto_adapter_basic.go: ToProtoPlayer, FromProtoPlayer, ToProtoPlayersArray, FromProtoPlayersArray, ToProtoCards, FromProtoCards, timeFromMillis
// - proto_adapter_deal.go: ToProtoDeal, FromProtoDeal, ToProtoDeals, FromProtoDeals
// - proto_adapter_action.go: ToProtoPlayActions, FromProtoPlayActions, ToProtoTributePhase, FromProtoTributePhase
// - proto_adapter_result.go: ToProtoTeamUpgrades, FromProtoTeamUpgrades
//
// 被依赖:
// - 无（顶层适配器）
package sdk

import (
	"time"

	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== Match Adapters ====================

// ToProtoMatch 转换 SDK Match 到 Proto Match
// 特殊处理：
// - StartTime: 零值时转换为 0
// - EndTime: nil 时转换为 0
// - Winner: -1 表示未结束
func ToProtoMatch(m *Match) *pbgame.Match {
	if m == nil {
		return nil
	}

	var startTimeMs int64
	if !m.StartTime.IsZero() {
		startTimeMs = m.StartTime.UnixMilli()
	}

	var endTimeMs int64
	if m.EndTime != nil {
		endTimeMs = m.EndTime.UnixMilli()
	}

	return &pbgame.Match{
		Id:          m.ID,
		Status:      ToProtoMatchStatus(m.Status),
		Players:     ToProtoPlayersArray(m.Players),
		CurrentDeal: ToProtoDeal(m.CurrentDeal),
		DealHistory: ToProtoDeals(m.DealHistory),
		TeamLevels:  ToProtoTeamUpgrades(m.TeamLevels),
		Winner:      int32(m.Winner),
		StartTimeMs: startTimeMs,
		EndTimeMs:   endTimeMs,
	}
}

// FromProtoMatch 转换 Proto Match 到 SDK Match
// 特殊处理：
// - StartTimeMs: <= 0 时转换为 time.Time{} 零值（通过 timeFromMillis）
// - EndTimeMs: 0 时转换为 nil
// 注意: StartTimeMs 为 0 表示时间未设置，会被解析为零值时间
func FromProtoMatch(pm *pbgame.Match) *Match {
	if pm == nil {
		return nil
	}

	var endTime *time.Time
	if pm.EndTimeMs > 0 {
		t := timeFromMillis(pm.EndTimeMs)
		endTime = &t
	}

	return &Match{
		ID:          pm.Id,
		Status:      FromProtoMatchStatus(pm.Status),
		Players:     FromProtoPlayersArray(pm.Players),
		CurrentDeal: FromProtoDeal(pm.CurrentDeal),
		DealHistory: FromProtoDeals(pm.DealHistory),
		TeamLevels:  FromProtoTeamUpgrades(pm.TeamLevels),
		Winner:      int(pm.Winner),
		StartTime:   timeFromMillis(pm.StartTimeMs),
		EndTime:     endTime,
	}
}

// ==================== PlayerView Adapters ====================

// ToProtoPlayerView 转换 SDK PlayerView 到 Proto PlayerView
// 特殊处理：
// - CurrentTurn: nil 时转换为 -1
func ToProtoPlayerView(pv *PlayerView) *pbgame.PlayerView {
	if pv == nil {
		return nil
	}

	currentTurn := int32(-1)
	if pv.CurrentTurn != nil {
		currentTurn = int32(*pv.CurrentTurn)
	}

	return &pbgame.PlayerView{
		PlayerSeat:   int32(pv.PlayerSeat),
		PlayerCards:  ToProtoCards(pv.PlayerCards),
		TeamLevels:   ToProtoTeamUpgrades(pv.TeamLevels),
		DealLevel:    int32(pv.DealLevel),
		DealStatus:   ToProtoDealStatus(pv.DealStatus),
		TrickId:      pv.TrickID,
		CurrentTurn:  currentTurn,
		Plays:        ToProtoPlayActions(pv.Plays),
		TributePhase: ToProtoTributePhase(pv.TributePhase),
	}
}

// FromProtoPlayerView 转换 Proto PlayerView 到 SDK PlayerView
// 特殊处理：
// - CurrentTurn: -1 时转换为 nil
func FromProtoPlayerView(ppv *pbgame.PlayerView) *PlayerView {
	if ppv == nil {
		return nil
	}

	var currentTurn *int
	if ppv.CurrentTurn >= 0 {
		turn := int(ppv.CurrentTurn)
		currentTurn = &turn
	}

	return &PlayerView{
		PlayerSeat:   int(ppv.PlayerSeat),
		PlayerCards:  FromProtoCards(ppv.PlayerCards),
		TeamLevels:   FromProtoTeamUpgrades(ppv.TeamLevels),
		DealLevel:    int(ppv.DealLevel),
		DealStatus:   FromProtoDealStatus(ppv.DealStatus),
		TrickID:      ppv.TrickId,
		CurrentTurn:  currentTurn,
		Plays:        FromProtoPlayActions(ppv.Plays),
		TributePhase: FromProtoTributePhase(ppv.TributePhase),
	}
}
