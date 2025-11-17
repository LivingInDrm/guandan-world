package sdk

import (
	eventpb "guandan-world/proto/event"
)

// ConvertCardToProto converts SDK Card to proto Card
func ConvertCardToProto(sdkCard *Card) *eventpb.Card {
	if sdkCard == nil {
		return nil
	}
	// Convert Color string to suit int
	// Proto suit: 0=Spade, 1=Heart, 2=Club, 3=Diamond, -1=Joker
	suit := -1
	switch sdkCard.Color {
	case "Spade":
		suit = 0
	case "Heart":
		suit = 1
	case "Club":
		suit = 2
	case "Diamond":
		suit = 3
	case "":
		// Joker cards have empty color
		suit = -1
	default:
		// Unknown color - log and use -1
		suit = -1
	}
	
	return &eventpb.Card{
		Suit:      int32(suit),
		Rank:      int32(sdkCard.Number),
		DeckIndex: int32(sdkCard.DeckIndex),
	}
}

// ConvertCardFromProto converts proto Card to SDK Card
func ConvertCardFromProto(protoCard *eventpb.Card) *Card {
	if protoCard == nil {
		return nil
	}
	// Convert suit int to Color string
	color := ""
	switch protoCard.Suit {
	case 0:
		color = "Spade"
	case 1:
		color = "Heart"
	case 2:
		color = "Club"
	case 3:
		color = "Diamond"
	}
	
	return &Card{
		Number:    int(protoCard.Rank),
		RawNumber: int(protoCard.Rank),
		Color:     color,
		DeckIndex: int(protoCard.DeckIndex),
	}
}

// ConvertCardsToProto converts slice of SDK Cards to proto Cards
func ConvertCardsToProto(sdkCards []*Card) []*eventpb.Card {
	if sdkCards == nil {
		return nil
	}
	protoCards := make([]*eventpb.Card, len(sdkCards))
	for i, card := range sdkCards {
		protoCards[i] = ConvertCardToProto(card)
	}
	return protoCards
}

// ConvertCardsFromProto converts slice of proto Cards to SDK Cards
func ConvertCardsFromProto(protoCards []*eventpb.Card) []*Card {
	if protoCards == nil {
		return nil
	}
	sdkCards := make([]*Card, len(protoCards))
	for i, card := range protoCards {
		sdkCards[i] = ConvertCardFromProto(card)
	}
	return sdkCards
}

// ConvertVictoryTypeToProto converts SDK VictoryType to proto
func ConvertVictoryTypeToProto(vt VictoryType) eventpb.VictoryType {
	switch vt {
	case VictoryTypeDoubleDown:
		return eventpb.VictoryType_VICTORY_TYPE_DOUBLE_DOWN
	case VictoryTypeSingleLast:
		return eventpb.VictoryType_VICTORY_TYPE_SINGLE_LAST
	case VictoryTypePartnerLast:
		return eventpb.VictoryType_VICTORY_TYPE_PARTNER_LAST
	default:
		return eventpb.VictoryType_VICTORY_TYPE_UNSPECIFIED
	}
}

// ConvertVictoryTypeFromProto converts proto VictoryType to SDK
func ConvertVictoryTypeFromProto(vt eventpb.VictoryType) VictoryType {
	switch vt {
	case eventpb.VictoryType_VICTORY_TYPE_DOUBLE_DOWN:
		return VictoryTypeDoubleDown
	case eventpb.VictoryType_VICTORY_TYPE_SINGLE_LAST:
		return VictoryTypeSingleLast
	case eventpb.VictoryType_VICTORY_TYPE_PARTNER_LAST:
		return VictoryTypePartnerLast
	case eventpb.VictoryType_VICTORY_TYPE_UNSPECIFIED:
		// Unspecified - use partner_last as safe default (lowest upgrade)
		return VictoryTypePartnerLast
	default:
		// Unknown proto value - use partner_last as safe default
		// In production, this should be logged as it indicates a schema mismatch
		return VictoryTypePartnerLast
	}
}

// ConvertTributeTypeToProto converts tribute type string to proto enum
func ConvertTributeTypeToProto(tt string) eventpb.TributeType {
	switch tt {
	case "double_down":
		return eventpb.TributeType_TRIBUTE_TYPE_DOUBLE_DOWN
	case "single_last":
		return eventpb.TributeType_TRIBUTE_TYPE_SINGLE_LAST
	case "partner_last":
		return eventpb.TributeType_TRIBUTE_TYPE_PARTNER_LAST
	case "none":
		return eventpb.TributeType_TRIBUTE_TYPE_NONE
	default:
		return eventpb.TributeType_TRIBUTE_TYPE_UNSPECIFIED
	}
}

// ConvertTributeTypeFromProto converts proto TributeType to string
func ConvertTributeTypeFromProto(tt eventpb.TributeType) string {
	switch tt {
	case eventpb.TributeType_TRIBUTE_TYPE_DOUBLE_DOWN:
		return "double_down"
	case eventpb.TributeType_TRIBUTE_TYPE_SINGLE_LAST:
		return "single_last"
	case eventpb.TributeType_TRIBUTE_TYPE_PARTNER_LAST:
		return "partner_last"
	case eventpb.TributeType_TRIBUTE_TYPE_NONE:
		return "none"
	default:
		return "none"
	}
}

// ConvertPlayerTimeoutActionTypeToProto converts timeout action string to proto enum
func ConvertPlayerTimeoutActionTypeToProto(actionType string) eventpb.PlayerTimeoutActionType {
	switch actionType {
	case "play_decision":
		return eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_PLAY_DECISION
	case "tribute_select":
		return eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_TRIBUTE_SELECT
	case "return_tribute":
		return eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_RETURN_TRIBUTE
	default:
		return eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_UNSPECIFIED
	}
}

// ConvertPlayerTimeoutActionTypeFromProto converts proto PlayerTimeoutActionType to string
func ConvertPlayerTimeoutActionTypeFromProto(actionType eventpb.PlayerTimeoutActionType) string {
	switch actionType {
	case eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_PLAY_DECISION:
		return "play_decision"
	case eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_TRIBUTE_SELECT:
		return "tribute_select"
	case eventpb.PlayerTimeoutActionType_PLAYER_TIMEOUT_ACTION_RETURN_TRIBUTE:
		return "return_tribute"
	default:
		return ""
	}
}
