// proto_adapter_enums.go - 枚举类型 Proto 适配器
//
// 职责:
// - 所有枚举类型的双向转换（SDK ↔ Proto）
// - 包含 7 种枚举：VictoryType, DealStatus, MatchStatus, TrickStatus,
//   TributeStatus, CompType, GameEventType
//
// 依赖:
// - 无外部依赖（基础层）
//
// 被依赖:
// - proto_adapter_comp.go: 使用 ToProtoCompType, FromProtoCompType
// - proto_adapter_action.go: 使用 ToProtoTrickStatus, FromProtoTrickStatus
// - proto_adapter_tribute.go: 使用 ToProtoTributeStatus, FromProtoTributeStatus
// - proto_adapter_result.go: 使用 ToProtoVictoryType, FromProtoVictoryType
package sdk

import (
	pb "guandan-world/proto/gen/go/common"
)

// ==================== Enum Adapters ====================

// ToProtoVictoryType 转换 SDK VictoryType 到 Proto VictoryType
func ToProtoVictoryType(vt VictoryType) pb.VictoryType {
	switch vt {
	case VictoryTypeDoubleDown:
		return pb.VictoryType_VICTORY_TYPE_DOUBLE_DOWN
	case VictoryTypeSingleLast:
		return pb.VictoryType_VICTORY_TYPE_SINGLE_LAST
	case VictoryTypePartnerLast:
		return pb.VictoryType_VICTORY_TYPE_PARTNER_LAST
	default:
		return pb.VictoryType_VICTORY_TYPE_UNSPECIFIED
	}
}

// FromProtoVictoryType 转换 Proto VictoryType 到 SDK VictoryType
func FromProtoVictoryType(pvt pb.VictoryType) VictoryType {
	switch pvt {
	case pb.VictoryType_VICTORY_TYPE_DOUBLE_DOWN:
		return VictoryTypeDoubleDown
	case pb.VictoryType_VICTORY_TYPE_SINGLE_LAST:
		return VictoryTypeSingleLast
	case pb.VictoryType_VICTORY_TYPE_PARTNER_LAST:
		return VictoryTypePartnerLast
	default:
		return ""
	}
}

// ToProtoDealStatus 转换 SDK DealStatus 到 Proto DealStatus
func ToProtoDealStatus(ds DealStatus) pb.DealStatus {
	switch ds {
	case DealStatusWaiting:
		return pb.DealStatus_DEAL_STATUS_WAITING
	case DealStatusDealing:
		return pb.DealStatus_DEAL_STATUS_DEALING
	case DealStatusTribute:
		return pb.DealStatus_DEAL_STATUS_TRIBUTE
	case DealStatusPlaying:
		return pb.DealStatus_DEAL_STATUS_PLAYING
	case DealStatusFinished:
		return pb.DealStatus_DEAL_STATUS_FINISHED
	default:
		return pb.DealStatus_DEAL_STATUS_UNSPECIFIED
	}
}

// FromProtoDealStatus 转换 Proto DealStatus 到 SDK DealStatus
func FromProtoDealStatus(pds pb.DealStatus) DealStatus {
	switch pds {
	case pb.DealStatus_DEAL_STATUS_WAITING:
		return DealStatusWaiting
	case pb.DealStatus_DEAL_STATUS_DEALING:
		return DealStatusDealing
	case pb.DealStatus_DEAL_STATUS_TRIBUTE:
		return DealStatusTribute
	case pb.DealStatus_DEAL_STATUS_PLAYING:
		return DealStatusPlaying
	case pb.DealStatus_DEAL_STATUS_FINISHED:
		return DealStatusFinished
	default:
		return ""
	}
}

// ToProtoMatchStatus 转换 SDK MatchStatus 到 Proto MatchStatus
func ToProtoMatchStatus(ms MatchStatus) pb.MatchStatus {
	switch ms {
	case MatchStatusWaiting:
		return pb.MatchStatus_MATCH_STATUS_WAITING
	case MatchStatusPlaying:
		return pb.MatchStatus_MATCH_STATUS_PLAYING
	case MatchStatusFinished:
		return pb.MatchStatus_MATCH_STATUS_FINISHED
	default:
		return pb.MatchStatus_MATCH_STATUS_UNSPECIFIED
	}
}

// FromProtoMatchStatus 转换 Proto MatchStatus 到 SDK MatchStatus
func FromProtoMatchStatus(pms pb.MatchStatus) MatchStatus {
	switch pms {
	case pb.MatchStatus_MATCH_STATUS_WAITING:
		return MatchStatusWaiting
	case pb.MatchStatus_MATCH_STATUS_PLAYING:
		return MatchStatusPlaying
	case pb.MatchStatus_MATCH_STATUS_FINISHED:
		return MatchStatusFinished
	default:
		return ""
	}
}

// ToProtoTrickStatus 转换 SDK TrickStatus 到 Proto TrickStatus
func ToProtoTrickStatus(ts TrickStatus) pb.TrickStatus {
	switch ts {
	case TrickStatusWaiting:
		return pb.TrickStatus_TRICK_STATUS_WAITING
	case TrickStatusPlaying:
		return pb.TrickStatus_TRICK_STATUS_PLAYING
	case TrickStatusFinished:
		return pb.TrickStatus_TRICK_STATUS_FINISHED
	default:
		return pb.TrickStatus_TRICK_STATUS_UNSPECIFIED
	}
}

// FromProtoTrickStatus 转换 Proto TrickStatus 到 SDK TrickStatus
func FromProtoTrickStatus(pts pb.TrickStatus) TrickStatus {
	switch pts {
	case pb.TrickStatus_TRICK_STATUS_WAITING:
		return TrickStatusWaiting
	case pb.TrickStatus_TRICK_STATUS_PLAYING:
		return TrickStatusPlaying
	case pb.TrickStatus_TRICK_STATUS_FINISHED:
		return TrickStatusFinished
	default:
		return ""
	}
}

// ToProtoTributeStatus 转换 SDK TributeStatus 到 Proto TributeStatus
func ToProtoTributeStatus(ts TributeStatus) pb.TributeStatus {
	switch ts {
	case TributeStatusWaiting:
		return pb.TributeStatus_TRIBUTE_STATUS_WAITING
	case TributeStatusSelecting:
		return pb.TributeStatus_TRIBUTE_STATUS_SELECTING
	case TributeStatusReturning:
		return pb.TributeStatus_TRIBUTE_STATUS_RETURNING
	case TributeStatusFinished:
		return pb.TributeStatus_TRIBUTE_STATUS_FINISHED
	default:
		return pb.TributeStatus_TRIBUTE_STATUS_UNSPECIFIED
	}
}

// FromProtoTributeStatus 转换 Proto TributeStatus 到 SDK TributeStatus
func FromProtoTributeStatus(pts pb.TributeStatus) TributeStatus {
	switch pts {
	case pb.TributeStatus_TRIBUTE_STATUS_WAITING:
		return TributeStatusWaiting
	case pb.TributeStatus_TRIBUTE_STATUS_SELECTING:
		return TributeStatusSelecting
	case pb.TributeStatus_TRIBUTE_STATUS_RETURNING:
		return TributeStatusReturning
	case pb.TributeStatus_TRIBUTE_STATUS_FINISHED:
		return TributeStatusFinished
	default:
		return ""
	}
}

// ToProtoCompType 转换 SDK CompType 到 Proto CompType
func ToProtoCompType(ct CompType) pb.CompType {
	switch ct {
	case TypeFold:
		return pb.CompType_COMP_TYPE_FOLD
	case TypeIllegal:
		return pb.CompType_COMP_TYPE_ILLEGAL
	case TypeSingle:
		return pb.CompType_COMP_TYPE_SINGLE
	case TypePair:
		return pb.CompType_COMP_TYPE_PAIR
	case TypeTriple:
		return pb.CompType_COMP_TYPE_TRIPLE
	case TypeFullHouse:
		return pb.CompType_COMP_TYPE_FULL_HOUSE
	case TypeStraight:
		return pb.CompType_COMP_TYPE_STRAIGHT
	case TypePlate:
		return pb.CompType_COMP_TYPE_PLATE
	case TypeTube:
		return pb.CompType_COMP_TYPE_TUBE
	case TypeJokerBomb:
		return pb.CompType_COMP_TYPE_JOKER_BOMB
	case TypeNaiveBomb:
		return pb.CompType_COMP_TYPE_NAIVE_BOMB
	case TypeStraightFlush:
		return pb.CompType_COMP_TYPE_STRAIGHT_FLUSH
	default:
		return pb.CompType_COMP_TYPE_UNSPECIFIED
	}
}

// FromProtoCompType 转换 Proto CompType 到 SDK CompType
func FromProtoCompType(pct pb.CompType) CompType {
	switch pct {
	case pb.CompType_COMP_TYPE_FOLD:
		return TypeFold
	case pb.CompType_COMP_TYPE_ILLEGAL:
		return TypeIllegal
	case pb.CompType_COMP_TYPE_SINGLE:
		return TypeSingle
	case pb.CompType_COMP_TYPE_PAIR:
		return TypePair
	case pb.CompType_COMP_TYPE_TRIPLE:
		return TypeTriple
	case pb.CompType_COMP_TYPE_FULL_HOUSE:
		return TypeFullHouse
	case pb.CompType_COMP_TYPE_STRAIGHT:
		return TypeStraight
	case pb.CompType_COMP_TYPE_PLATE:
		return TypePlate
	case pb.CompType_COMP_TYPE_TUBE:
		return TypeTube
	case pb.CompType_COMP_TYPE_JOKER_BOMB:
		return TypeJokerBomb
	case pb.CompType_COMP_TYPE_NAIVE_BOMB:
		return TypeNaiveBomb
	case pb.CompType_COMP_TYPE_STRAIGHT_FLUSH:
		return TypeStraightFlush
	default:
		return -1 // Invalid type
	}
}

// ToProtoGameEventType 转换 SDK GameEventType 到 Proto GameEventType
func ToProtoGameEventType(get GameEventType) pb.GameEventType {
	switch get {
	case EventMatchStarted:
		return pb.GameEventType_GAME_EVENT_TYPE_MATCH_STARTED
	case EventDealStarted:
		return pb.GameEventType_GAME_EVENT_TYPE_DEAL_STARTED
	case EventCardsDealt:
		return pb.GameEventType_GAME_EVENT_TYPE_CARDS_DEALT
	case EventTributePhase:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_PHASE
	case EventTributeRulesSet:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_RULES_SET
	case EventTributeImmunity:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_IMMUNITY
	case EventTributePoolCreated:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_POOL_CREATED
	case EventTributeStarted:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_STARTED
	case EventTributeGiven:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_GIVEN
	case EventTributeSelected:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_SELECTED
	case EventReturnTribute:
		return pb.GameEventType_GAME_EVENT_TYPE_RETURN_TRIBUTE
	case EventTributeCompleted:
		return pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_COMPLETED
	case EventTrickStarted:
		return pb.GameEventType_GAME_EVENT_TYPE_TRICK_STARTED
	case EventPlayerPlayed:
		return pb.GameEventType_GAME_EVENT_TYPE_PLAYER_PLAYED
	case EventPlayerPassed:
		return pb.GameEventType_GAME_EVENT_TYPE_PLAYER_PASSED
	case EventTrickEnded:
		return pb.GameEventType_GAME_EVENT_TYPE_TRICK_ENDED
	case EventDealEnded:
		return pb.GameEventType_GAME_EVENT_TYPE_DEAL_ENDED
	case EventMatchEnded:
		return pb.GameEventType_GAME_EVENT_TYPE_MATCH_ENDED
	case EventPlayerTimeout:
		return pb.GameEventType_GAME_EVENT_TYPE_PLAYER_TIMEOUT
	case EventPlayerDisconnect:
		return pb.GameEventType_GAME_EVENT_TYPE_PLAYER_DISCONNECT
	case EventPlayerReconnect:
		return pb.GameEventType_GAME_EVENT_TYPE_PLAYER_RECONNECT
	default:
		return pb.GameEventType_GAME_EVENT_TYPE_UNSPECIFIED
	}
}

// FromProtoGameEventType 转换 Proto GameEventType 到 SDK GameEventType
func FromProtoGameEventType(pget pb.GameEventType) GameEventType {
	switch pget {
	case pb.GameEventType_GAME_EVENT_TYPE_MATCH_STARTED:
		return EventMatchStarted
	case pb.GameEventType_GAME_EVENT_TYPE_DEAL_STARTED:
		return EventDealStarted
	case pb.GameEventType_GAME_EVENT_TYPE_CARDS_DEALT:
		return EventCardsDealt
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_PHASE:
		return EventTributePhase
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_RULES_SET:
		return EventTributeRulesSet
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_IMMUNITY:
		return EventTributeImmunity
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_POOL_CREATED:
		return EventTributePoolCreated
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_STARTED:
		return EventTributeStarted
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_GIVEN:
		return EventTributeGiven
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_SELECTED:
		return EventTributeSelected
	case pb.GameEventType_GAME_EVENT_TYPE_RETURN_TRIBUTE:
		return EventReturnTribute
	case pb.GameEventType_GAME_EVENT_TYPE_TRIBUTE_COMPLETED:
		return EventTributeCompleted
	case pb.GameEventType_GAME_EVENT_TYPE_TRICK_STARTED:
		return EventTrickStarted
	case pb.GameEventType_GAME_EVENT_TYPE_PLAYER_PLAYED:
		return EventPlayerPlayed
	case pb.GameEventType_GAME_EVENT_TYPE_PLAYER_PASSED:
		return EventPlayerPassed
	case pb.GameEventType_GAME_EVENT_TYPE_TRICK_ENDED:
		return EventTrickEnded
	case pb.GameEventType_GAME_EVENT_TYPE_DEAL_ENDED:
		return EventDealEnded
	case pb.GameEventType_GAME_EVENT_TYPE_MATCH_ENDED:
		return EventMatchEnded
	case pb.GameEventType_GAME_EVENT_TYPE_PLAYER_TIMEOUT:
		return EventPlayerTimeout
	case pb.GameEventType_GAME_EVENT_TYPE_PLAYER_DISCONNECT:
		return EventPlayerDisconnect
	case pb.GameEventType_GAME_EVENT_TYPE_PLAYER_RECONNECT:
		return EventPlayerReconnect
	default:
		return ""
	}
}
