// proto_adapter_event.go - GameEvent Proto 适配器
//
// 职责:
// - GameEvent 类型的双向转换（SDK ↔ Proto）
// - 各种具体事件类型的转换（MatchStartedEvent, DealStartedEvent 等）
//
// 依赖:
// - proto_adapter_basic.go: timeFromMillis
// - proto_adapter_match.go: ToProtoMatch, FromProtoMatch
// - proto_adapter_deal.go: ToProtoDeal, FromProtoDeal
// - proto_adapter_result.go: ToProtoDealResult, FromProtoDealResult, ToProtoMatchResult, FromProtoMatchResult, ToProtoTeamUpgrades, FromProtoTeamUpgrades, ToProtoDealStatistics, FromProtoDealStatistics
// - proto_adapter_enums.go: ToProtoGameEventType, FromProtoGameEventType
//
// 被依赖:
// - backend event adapters
package sdk

import (
	pbmsg "guandan-world/proto/gen/go/messages"
)

// ==================== GameEvent Adapters ====================

// ToProtoGameEvent 转换 SDK GameEvent 到 Proto GameEvent
// 特殊处理：
// - Data 字段根据 Type 解析为不同的具体事件类型
// - Timestamp 转换为毫秒时间戳
// - PlayerSeat 保留原值（-1 表示无关联玩家）
func ToProtoGameEvent(e *GameEvent) *pbmsg.GameEvent {
	if e == nil {
		return nil
	}

	var timestampMs int64
	if !e.Timestamp.IsZero() {
		timestampMs = e.Timestamp.UnixMilli()
	}

	result := &pbmsg.GameEvent{
		Type:        ToProtoGameEventType(e.Type),
		TimestampMs: timestampMs,
		PlayerSeat:  int32(e.PlayerSeat),
	}

	// 根据事件类型解析 Data 字段
	switch e.Type {
	case EventMatchStarted:
		if match, ok := e.Data.(*Match); ok {
			result.Payload = &pbmsg.GameEvent_MatchStarted{
				MatchStarted: &pbmsg.MatchStartedEvent{
					Match: ToProtoMatch(match),
				},
			}
		}

	case EventDealStarted:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			// 数据类型不匹配，跳过此事件
			break
		}
		
		deal, hasDeal := data["deal"].(*Deal)
		if !hasDeal || deal == nil {
			// 缺少必要的deal字段
			break
		}
		
		// team_levels 可以从 team0_level 和 team1_level 构建
		var teamLevels [2]int
		if team0Level, ok := data["team0_level"].(int); ok {
			teamLevels[0] = team0Level
		}
		if team1Level, ok := data["team1_level"].(int); ok {
			teamLevels[1] = team1Level
		}
		
		result.Payload = &pbmsg.GameEvent_DealStarted{
			DealStarted: &pbmsg.DealStartedEvent{
				Deal:       ToProtoDeal(deal),
				TeamLevels: ToProtoTeamUpgrades(teamLevels),
			},
		}

	case EventCardsDealt:
		// 当前代码中未实际触发此事件，但预留转换逻辑
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		event := &pbmsg.CardsDealtEvent{}
		if handSizes, ok := data["hand_sizes"].([]int); ok {
			event.HandSizes = make([]int32, len(handSizes))
			for i, size := range handSizes {
				event.HandSizes[i] = int32(size)
			}
		}
		if dealer, ok := data["dealer"].(int); ok {
			event.Dealer = int32(dealer)
		}
		result.Payload = &pbmsg.GameEvent_CardsDealt{
			CardsDealt: event,
		}

	case EventDealEnded:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			// 数据类型不匹配，跳过此事件
			break
		}
		
		deal, hasDeal := data["deal"].(*Deal)
		dealResult, hasResult := data["result"].(*DealResult)
		if !hasDeal || deal == nil || !hasResult || dealResult == nil {
			// 缺少必要字段
			break
		}
		
		result.Payload = &pbmsg.GameEvent_DealEnded{
			DealEnded: &pbmsg.DealEndedEvent{
				Deal:   ToProtoDeal(deal),
				Result: ToProtoDealResult(dealResult),
			},
		}

	case EventMatchEnded:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			// 数据类型不匹配，跳过此事件
			break
		}
		
		match, hasMatch := data["match"].(*Match)
		matchResult, hasResult := data["result"].(*MatchResult)
		if !hasMatch || match == nil || !hasResult || matchResult == nil {
			// 缺少必要字段
			break
		}
		
		result.Payload = &pbmsg.GameEvent_MatchEnded{
			MatchEnded: &pbmsg.MatchEndedEvent{
				Match:  ToProtoMatch(match),
				Result: ToProtoMatchResult(matchResult),
			},
		}
	}

	return result
}

// FromProtoGameEvent 转换 Proto GameEvent 到 SDK GameEvent
// 特殊处理：
// - Payload oneof 字段根据具体类型转换为 Data interface{}
// - TimestampMs 转换为 time.Time
func FromProtoGameEvent(pe *pbmsg.GameEvent) *GameEvent {
	if pe == nil {
		return nil
	}

	result := &GameEvent{
		Type:       FromProtoGameEventType(pe.Type),
		Timestamp:  timeFromMillis(pe.TimestampMs),
		PlayerSeat: int(pe.PlayerSeat),
	}

	// 根据 payload 类型设置 Data
	switch payload := pe.Payload.(type) {
	case *pbmsg.GameEvent_MatchStarted:
		if payload.MatchStarted != nil && payload.MatchStarted.Match != nil {
			result.Data = FromProtoMatch(payload.MatchStarted.Match)
		}

	case *pbmsg.GameEvent_DealStarted:
		if payload.DealStarted != nil && payload.DealStarted.Deal != nil {
			teamLevels := FromProtoTeamUpgrades(payload.DealStarted.TeamLevels)
			result.Data = map[string]interface{}{
				"deal":        FromProtoDeal(payload.DealStarted.Deal),
				"deal_level":  int(payload.DealStarted.Deal.Level), // 从deal中提取
				"team0_level": teamLevels[0],
				"team1_level": teamLevels[1],
			}
		}

	case *pbmsg.GameEvent_CardsDealt:
		if payload.CardsDealt != nil {
			handSizes := make([]int, len(payload.CardsDealt.HandSizes))
			for i, size := range payload.CardsDealt.HandSizes {
				handSizes[i] = int(size)
			}
			result.Data = map[string]interface{}{
				"hand_sizes": handSizes,
				"dealer":     int(payload.CardsDealt.Dealer),
			}
		}

	case *pbmsg.GameEvent_DealEnded:
		if payload.DealEnded != nil && payload.DealEnded.Deal != nil && payload.DealEnded.Result != nil {
			deal := FromProtoDeal(payload.DealEnded.Deal)
			dealResult := FromProtoDealResult(payload.DealEnded.Result)
			
			result.Data = map[string]interface{}{
				"deal":       deal,
				"result":     dealResult,
				"rankings":   deal.Rankings,              // 从deal中获取
				"statistics": dealResult.Statistics,      // 从result中获取
			}
		}

	case *pbmsg.GameEvent_MatchEnded:
		if payload.MatchEnded != nil && payload.MatchEnded.Match != nil && payload.MatchEnded.Result != nil {
			match := FromProtoMatch(payload.MatchEnded.Match)
			result.Data = map[string]interface{}{
				"match":        match,
				"result":       FromProtoMatchResult(payload.MatchEnded.Result),
				"winner":       match.Winner,       // 从match中获取
				"final_levels": match.TeamLevels,   // 从match中获取
			}
		}

	default:
		// 未知的 payload 类型
		result.Data = nil
	}

	return result
}
