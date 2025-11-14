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

	// ==================== 进贡事件 (6B组) ====================

	case EventTributePhase:
		if tributePhase, ok := e.Data.(*TributePhase); ok {
			result.Payload = &pbmsg.GameEvent_TributePhase{
				TributePhase: &pbmsg.TributePhaseEvent{
					TributePhase: ToProtoTributePhase(tributePhase),
				},
			}
		}

	case EventTributeRulesSet:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		lastResult, _ := data["last_result"].(*DealResult)
		victoryType, _ := data["victory_type"].(VictoryType)
		playerRankings, _ := data["player_rankings"].([]int)
		
		event := &pbmsg.TributeRulesSetEvent{
			LastResult:  ToProtoDealResult(lastResult),
			VictoryType: ToProtoVictoryType(victoryType),
		}
		
		// 转换 player_rankings
		if len(playerRankings) > 0 {
			event.PlayerRankings = make([]int32, len(playerRankings))
			for i, rank := range playerRankings {
				event.PlayerRankings[i] = int32(rank)
			}
		}
		
		// 转换 tribute_rules
		if tributeRulesData, ok := data["tribute_rules"].(map[string]interface{}); ok {
			event.TributeRules = &pbmsg.TributeRules{}
			
			if tributeMap, ok := tributeRulesData["tribute_map"].(map[int]int); ok {
				event.TributeRules.TributeMap = make(map[int32]int32)
				for k, v := range tributeMap {
					event.TributeRules.TributeMap[int32(k)] = int32(v)
				}
			}
			
			if isDoubleDown, ok := tributeRulesData["is_double_down"].(bool); ok {
				event.TributeRules.IsDoubleDown = isDoubleDown
			}
			
			if description, ok := tributeRulesData["description"].(string); ok {
				event.TributeRules.Description = description
			}
		}
		
		result.Payload = &pbmsg.GameEvent_TributeRulesSet{
			TributeRulesSet: event,
		}

	case EventTributeImmunity:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		tributePhase, _ := data["tribute_phase"].(*TributePhase)
		immunityReason, _ := data["immunity_reason"].(string)
		
		result.Payload = &pbmsg.GameEvent_TributeImmunity{
			TributeImmunity: &pbmsg.TributeImmunityEvent{
				TributePhase:   ToProtoTributePhase(tributePhase),
				ImmunityReason: immunityReason,
			},
		}

	case EventTributePoolCreated:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		event := &pbmsg.TributePoolCreatedEvent{}
		
		if description, ok := data["description"].(string); ok {
			event.Description = description
		}
		
		if selectingPlayer, ok := data["selecting_player"].(int); ok {
			event.SelectingPlayer = int32(selectingPlayer)
		}
		
		if poolCards, ok := data["pool_cards"].([]*Card); ok {
			event.PoolCards = ToProtoCards(poolCards)
		}
		
		if selectionOrder, ok := data["selection_order"].([]int); ok {
			event.SelectionOrder = make([]int32, len(selectionOrder))
			for i, order := range selectionOrder {
				event.SelectionOrder[i] = int32(order)
			}
		}
		
		// 转换 contributors
		if contributors, ok := data["contributors"].([]map[string]interface{}); ok {
			event.Contributors = make([]*pbmsg.TributeContributor, len(contributors))
			for i, contrib := range contributors {
				playerSeat, _ := contrib["player_seat"].(int)
				card, _ := contrib["card"].(*Card)
				event.Contributors[i] = &pbmsg.TributeContributor{
					PlayerSeat: int32(playerSeat),
					Card:       ToProtoCard(card),
				}
			}
		}
		
		result.Payload = &pbmsg.GameEvent_TributePoolCreated{
			TributePoolCreated: event,
		}

	case EventTributeGiven:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		event := &pbmsg.TributeGivenEvent{}
		
		if giver, ok := data["giver"].(int); ok {
			event.Giver = int32(giver)
		}
		if receiver, ok := data["receiver"].(int); ok {
			event.Receiver = int32(receiver)
		}
		if card, ok := data["card"].(*Card); ok {
			event.Card = ToProtoCard(card)
		}
		if tributeType, ok := data["tribute_type"].(string); ok {
			event.TributeType = tributeType
		}
		if isAutoSelected, ok := data["is_auto_selected"].(bool); ok {
			event.IsAutoSelected = isAutoSelected
		}
		if selectionReason, ok := data["selection_reason"].(string); ok {
			event.SelectionReason = selectionReason
		}
		
		result.Payload = &pbmsg.GameEvent_TributeGiven{
			TributeGiven: event,
		}

	case EventTributeSelected:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		event := &pbmsg.TributeSelectedEvent{}
		
		if action, ok := data["action"].(string); ok {
			event.Action = action
		}
		if player, ok := data["player"].(int); ok {
			event.Player = int32(player)
		}
		if cardID, ok := data["cardID"].(string); ok {
			event.CardId = cardID
		}
		if selectedCard, ok := data["selected_card"].(*Card); ok {
			event.SelectedCard = ToProtoCard(selectedCard)
		}
		if remainingOptions, ok := data["remaining_options"].([]*Card); ok {
			event.RemainingOptions = ToProtoCards(remainingOptions)
		}
		if selectionOrder, ok := data["selection_order"].(int); ok {
			event.SelectionOrder = int32(selectionOrder)
		}
		if isTimeout, ok := data["is_timeout"].(bool); ok {
			event.IsTimeout = isTimeout
		}
		
		result.Payload = &pbmsg.GameEvent_TributeSelected{
			TributeSelected: event,
		}

	case EventReturnTribute:
		data, ok := e.Data.(map[string]interface{})
		if !ok {
			break
		}
		
		event := &pbmsg.ReturnTributeEvent{}
		
		if action, ok := data["action"].(string); ok {
			event.Action = action
		}
		if player, ok := data["player"].(int); ok {
			event.Player = int32(player)
		}
		if cardID, ok := data["cardID"].(string); ok {
			event.CardId = cardID
		}
		if returnCard, ok := data["return_card"].(*Card); ok {
			event.ReturnCard = ToProtoCard(returnCard)
		}
		if targetPlayer, ok := data["target_player"].(int); ok {
			event.TargetPlayer = int32(targetPlayer)
		}
		if originalTribute, ok := data["original_tribute"].(*Card); ok {
			event.OriginalTribute = ToProtoCard(originalTribute)
		}
		if isAutoSelected, ok := data["is_auto_selected"].(bool); ok {
			event.IsAutoSelected = isAutoSelected
		}
		if selectionReason, ok := data["selection_reason"].(string); ok {
			event.SelectionReason = selectionReason
		}
		
		result.Payload = &pbmsg.GameEvent_ReturnTribute{
			ReturnTribute: event,
		}

	case EventTributeStarted:
		if tributePhase, ok := e.Data.(*TributePhase); ok {
			result.Payload = &pbmsg.GameEvent_TributeStarted{
				TributeStarted: &pbmsg.TributeStartedEvent{
					TributePhase: ToProtoTributePhase(tributePhase),
				},
			}
		}

	case EventTributeCompleted:
		if tributePhase, ok := e.Data.(*TributePhase); ok {
			result.Payload = &pbmsg.GameEvent_TributeCompleted{
				TributeCompleted: &pbmsg.TributeCompletedEvent{
					TributePhase: ToProtoTributePhase(tributePhase),
				},
			}
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

	// ==================== 进贡事件 (6B组) ====================

	case *pbmsg.GameEvent_TributePhase:
		if payload.TributePhase != nil && payload.TributePhase.TributePhase != nil {
			result.Data = FromProtoTributePhase(payload.TributePhase.TributePhase)
		}

	case *pbmsg.GameEvent_TributeRulesSet:
		if payload.TributeRulesSet != nil {
			data := make(map[string]interface{})
			
			if payload.TributeRulesSet.LastResult != nil {
				data["last_result"] = FromProtoDealResult(payload.TributeRulesSet.LastResult)
			}
			
			data["victory_type"] = FromProtoVictoryType(payload.TributeRulesSet.VictoryType)
			
			// 转换 player_rankings
			if len(payload.TributeRulesSet.PlayerRankings) > 0 {
				rankings := make([]int, len(payload.TributeRulesSet.PlayerRankings))
				for i, rank := range payload.TributeRulesSet.PlayerRankings {
					rankings[i] = int(rank)
				}
				data["player_rankings"] = rankings
			}
			
			// 转换 tribute_rules
			if payload.TributeRulesSet.TributeRules != nil {
				tributeRules := make(map[string]interface{})
				
				if len(payload.TributeRulesSet.TributeRules.TributeMap) > 0 {
					tributeMap := make(map[int]int)
					for k, v := range payload.TributeRulesSet.TributeRules.TributeMap {
						tributeMap[int(k)] = int(v)
					}
					tributeRules["tribute_map"] = tributeMap
				}
				
				tributeRules["is_double_down"] = payload.TributeRulesSet.TributeRules.IsDoubleDown
				tributeRules["description"] = payload.TributeRulesSet.TributeRules.Description
				
				data["tribute_rules"] = tributeRules
			}
			
			result.Data = data
		}

	case *pbmsg.GameEvent_TributeImmunity:
		if payload.TributeImmunity != nil {
			result.Data = map[string]interface{}{
				"tribute_phase":   FromProtoTributePhase(payload.TributeImmunity.TributePhase),
				"immunity_reason": payload.TributeImmunity.ImmunityReason,
			}
		}

	case *pbmsg.GameEvent_TributePoolCreated:
		if payload.TributePoolCreated != nil {
			data := make(map[string]interface{})
			data["description"] = payload.TributePoolCreated.Description
			data["selecting_player"] = int(payload.TributePoolCreated.SelectingPlayer)
			
			if len(payload.TributePoolCreated.PoolCards) > 0 {
				data["pool_cards"] = FromProtoCards(payload.TributePoolCreated.PoolCards)
			}
			
			if len(payload.TributePoolCreated.SelectionOrder) > 0 {
				selectionOrder := make([]int, len(payload.TributePoolCreated.SelectionOrder))
				for i, order := range payload.TributePoolCreated.SelectionOrder {
					selectionOrder[i] = int(order)
				}
				data["selection_order"] = selectionOrder
			}
			
			// 转换 contributors
			if len(payload.TributePoolCreated.Contributors) > 0 {
				contributors := make([]map[string]interface{}, len(payload.TributePoolCreated.Contributors))
				for i, contrib := range payload.TributePoolCreated.Contributors {
					contributors[i] = map[string]interface{}{
						"player_seat": int(contrib.PlayerSeat),
						"card":        FromProtoCard(contrib.Card),
					}
				}
				data["contributors"] = contributors
			}
			
			result.Data = data
		}

	case *pbmsg.GameEvent_TributeGiven:
		if payload.TributeGiven != nil {
			result.Data = map[string]interface{}{
				"giver":            int(payload.TributeGiven.Giver),
				"receiver":         int(payload.TributeGiven.Receiver),
				"card":             FromProtoCard(payload.TributeGiven.Card),
				"tribute_type":     payload.TributeGiven.TributeType,
				"is_auto_selected": payload.TributeGiven.IsAutoSelected,
				"selection_reason": payload.TributeGiven.SelectionReason,
			}
		}

	case *pbmsg.GameEvent_TributeSelected:
		if payload.TributeSelected != nil {
			data := make(map[string]interface{})
			data["action"] = payload.TributeSelected.Action
			data["player"] = int(payload.TributeSelected.Player)
			data["cardID"] = payload.TributeSelected.CardId
			data["selected_card"] = FromProtoCard(payload.TributeSelected.SelectedCard)
			data["selection_order"] = int(payload.TributeSelected.SelectionOrder)
			data["is_timeout"] = payload.TributeSelected.IsTimeout
			
			if len(payload.TributeSelected.RemainingOptions) > 0 {
				data["remaining_options"] = FromProtoCards(payload.TributeSelected.RemainingOptions)
			}
			
			result.Data = data
		}

	case *pbmsg.GameEvent_ReturnTribute:
		if payload.ReturnTribute != nil {
			result.Data = map[string]interface{}{
				"action":           payload.ReturnTribute.Action,
				"player":           int(payload.ReturnTribute.Player),
				"cardID":           payload.ReturnTribute.CardId,
				"return_card":      FromProtoCard(payload.ReturnTribute.ReturnCard),
				"target_player":    int(payload.ReturnTribute.TargetPlayer),
				"original_tribute": FromProtoCard(payload.ReturnTribute.OriginalTribute),
				"is_auto_selected": payload.ReturnTribute.IsAutoSelected,
				"selection_reason": payload.ReturnTribute.SelectionReason,
			}
		}

	case *pbmsg.GameEvent_TributeStarted:
		if payload.TributeStarted != nil && payload.TributeStarted.TributePhase != nil {
			result.Data = FromProtoTributePhase(payload.TributeStarted.TributePhase)
		}

	case *pbmsg.GameEvent_TributeCompleted:
		if payload.TributeCompleted != nil && payload.TributeCompleted.TributePhase != nil {
			result.Data = FromProtoTributePhase(payload.TributeCompleted.TributePhase)
		}

	default:
		// 未知的 payload 类型
		result.Data = nil
	}

	return result
}
