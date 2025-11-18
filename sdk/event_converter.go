package sdk

import (
	commonpb "guandan-world/proto/common"
	eventpb "guandan-world/proto/event"
	viewpb "guandan-world/proto/view"
	"time"
)

// ConvertCardToProto converts SDK Card to proto Card
func ConvertCardToProto(sdkCard *Card) *commonpb.Card {
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
	
	return &commonpb.Card{
		Suit:      int32(suit),
		Rank:      int32(sdkCard.Number),
		DeckIndex: int32(sdkCard.DeckIndex),
	}
}

// ConvertCardFromProto converts proto Card to SDK Card
func ConvertCardFromProto(protoCard *commonpb.Card) *Card {
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
func ConvertCardsToProto(sdkCards []*Card) []*commonpb.Card {
	if sdkCards == nil {
		return nil
	}
	protoCards := make([]*commonpb.Card, len(sdkCards))
	for i, card := range sdkCards {
		protoCards[i] = ConvertCardToProto(card)
	}
	return protoCards
}

// ConvertCardsFromProto converts slice of proto Cards to SDK Cards
func ConvertCardsFromProto(protoCards []*commonpb.Card) []*Card {
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

// ConvertDealStatusToProto converts SDK DealStatus to proto DealStatus
func ConvertDealStatusToProto(status DealStatus) viewpb.DealStatus {
	switch status {
	case DealStatusWaiting:
		return viewpb.DealStatus_DEAL_STATUS_WAITING
	case DealStatusDealing:
		return viewpb.DealStatus_DEAL_STATUS_DEALING
	case DealStatusTribute:
		return viewpb.DealStatus_DEAL_STATUS_TRIBUTE
	case DealStatusPlaying:
		return viewpb.DealStatus_DEAL_STATUS_PLAYING
	case DealStatusFinished:
		return viewpb.DealStatus_DEAL_STATUS_FINISHED
	default:
		return viewpb.DealStatus_DEAL_STATUS_UNSPECIFIED
	}
}

// ConvertDealStatusFromProto converts proto DealStatus to SDK DealStatus
func ConvertDealStatusFromProto(status viewpb.DealStatus) DealStatus {
	switch status {
	case viewpb.DealStatus_DEAL_STATUS_WAITING:
		return DealStatusWaiting
	case viewpb.DealStatus_DEAL_STATUS_DEALING:
		return DealStatusDealing
	case viewpb.DealStatus_DEAL_STATUS_TRIBUTE:
		return DealStatusTribute
	case viewpb.DealStatus_DEAL_STATUS_PLAYING:
		return DealStatusPlaying
	case viewpb.DealStatus_DEAL_STATUS_FINISHED:
		return DealStatusFinished
	default:
		return DealStatusWaiting
	}
}

// ConvertTributeStatusToProto converts SDK TributeStatus to proto TributeStatus
func ConvertTributeStatusToProto(status TributeStatus) viewpb.TributeStatus {
	switch status {
	case TributeStatusWaiting:
		return viewpb.TributeStatus_TRIBUTE_STATUS_WAITING
	case TributeStatusSelecting:
		return viewpb.TributeStatus_TRIBUTE_STATUS_SELECTING
	case TributeStatusReturning:
		return viewpb.TributeStatus_TRIBUTE_STATUS_RETURNING
	case TributeStatusFinished:
		return viewpb.TributeStatus_TRIBUTE_STATUS_FINISHED
	default:
		return viewpb.TributeStatus_TRIBUTE_STATUS_UNSPECIFIED
	}
}

// ConvertTributeStatusFromProto converts proto TributeStatus to SDK TributeStatus
func ConvertTributeStatusFromProto(status viewpb.TributeStatus) TributeStatus {
	switch status {
	case viewpb.TributeStatus_TRIBUTE_STATUS_WAITING:
		return TributeStatusWaiting
	case viewpb.TributeStatus_TRIBUTE_STATUS_SELECTING:
		return TributeStatusSelecting
	case viewpb.TributeStatus_TRIBUTE_STATUS_RETURNING:
		return TributeStatusReturning
	case viewpb.TributeStatus_TRIBUTE_STATUS_FINISHED:
		return TributeStatusFinished
	default:
		return TributeStatusWaiting
	}
}

// ConvertCompTypeToProto converts SDK CompType to proto CompType
func ConvertCompTypeToProto(compType CompType) commonpb.CompType {
	switch compType {
	case TypeFold:
		return commonpb.CompType_COMP_TYPE_FOLD
	case TypeIllegal:
		return commonpb.CompType_COMP_TYPE_ILLEGAL
	case TypeSingle:
		return commonpb.CompType_COMP_TYPE_SINGLE
	case TypePair:
		return commonpb.CompType_COMP_TYPE_PAIR
	case TypeTriple:
		return commonpb.CompType_COMP_TYPE_TRIPLE
	case TypeFullHouse:
		return commonpb.CompType_COMP_TYPE_FULL_HOUSE
	case TypeStraight:
		return commonpb.CompType_COMP_TYPE_STRAIGHT
	case TypePlate:
		return commonpb.CompType_COMP_TYPE_PLATE
	case TypeTube:
		return commonpb.CompType_COMP_TYPE_TUBE
	case TypeJokerBomb:
		return commonpb.CompType_COMP_TYPE_JOKER_BOMB
	case TypeNaiveBomb:
		return commonpb.CompType_COMP_TYPE_NAIVE_BOMB
	case TypeStraightFlush:
		return commonpb.CompType_COMP_TYPE_STRAIGHT_FLUSH
	default:
		return commonpb.CompType_COMP_TYPE_UNSPECIFIED
	}
}

// ConvertCompTypeFromProto converts proto CompType to SDK CompType
func ConvertCompTypeFromProto(compType commonpb.CompType) CompType {
	switch compType {
	case commonpb.CompType_COMP_TYPE_FOLD:
		return TypeFold
	case commonpb.CompType_COMP_TYPE_ILLEGAL:
		return TypeIllegal
	case commonpb.CompType_COMP_TYPE_SINGLE:
		return TypeSingle
	case commonpb.CompType_COMP_TYPE_PAIR:
		return TypePair
	case commonpb.CompType_COMP_TYPE_TRIPLE:
		return TypeTriple
	case commonpb.CompType_COMP_TYPE_FULL_HOUSE:
		return TypeFullHouse
	case commonpb.CompType_COMP_TYPE_STRAIGHT:
		return TypeStraight
	case commonpb.CompType_COMP_TYPE_PLATE:
		return TypePlate
	case commonpb.CompType_COMP_TYPE_TUBE:
		return TypeTube
	case commonpb.CompType_COMP_TYPE_JOKER_BOMB:
		return TypeJokerBomb
	case commonpb.CompType_COMP_TYPE_NAIVE_BOMB:
		return TypeNaiveBomb
	case commonpb.CompType_COMP_TYPE_STRAIGHT_FLUSH:
		return TypeStraightFlush
	default:
		return TypeIllegal
	}
}

// ConvertPlayActionToProto converts SDK PlayAction to proto PlayAction
func ConvertPlayActionToProto(action *PlayAction) *commonpb.PlayAction {
	if action == nil {
		return nil
	}

	protoAction := &commonpb.PlayAction{
		PlayerSeat:  int32(action.PlayerSeat),
		Cards:       ConvertCardsToProto(action.Cards),
		TimestampMs: action.Timestamp.UnixMilli(),
		IsPass:      action.IsPass,
	}

	// Convert CompType if available
	if action.Comp != nil {
		protoAction.CompType = ConvertCompTypeToProto(action.Comp.GetType())
	} else if action.IsPass {
		protoAction.CompType = commonpb.CompType_COMP_TYPE_FOLD
	} else {
		protoAction.CompType = commonpb.CompType_COMP_TYPE_UNSPECIFIED
	}

	return protoAction
}

// ConvertPlayActionFromProto converts proto PlayAction to SDK PlayAction
func ConvertPlayActionFromProto(protoAction *commonpb.PlayAction) *PlayAction {
	if protoAction == nil {
		return nil
	}

	return &PlayAction{
		PlayerSeat: int(protoAction.PlayerSeat),
		Cards:      ConvertCardsFromProto(protoAction.Cards),
		Timestamp:  time.UnixMilli(protoAction.TimestampMs),
		IsPass:     protoAction.IsPass,
		// Note: Comp field is not reconstructed as it requires full CardComp implementation
	}
}

// ConvertTributePairToProto converts SDK TributePair to proto TributePair
func ConvertTributePairToProto(pair *TributePair) *viewpb.TributePair {
	if pair == nil {
		return nil
	}

	protoPair := &viewpb.TributePair{
		Giver:    int32(pair.Giver),
		Receiver: int32(pair.Receiver),
	}

	// Handle optional tribute card
	if pair.TributeCard != nil {
		protoPair.TributeCard = ConvertCardToProto(pair.TributeCard)
	}

	// Handle optional return card
	if pair.ReturnCard != nil {
		protoPair.ReturnCard = ConvertCardToProto(pair.ReturnCard)
	}

	return protoPair
}

// ConvertTributePairFromProto converts proto TributePair to SDK TributePair
func ConvertTributePairFromProto(protoPair *viewpb.TributePair) *TributePair {
	if protoPair == nil {
		return nil
	}

	return &TributePair{
		Giver:       int(protoPair.Giver),
		Receiver:    int(protoPair.Receiver),
		TributeCard: ConvertCardFromProto(protoPair.TributeCard),
		ReturnCard:  ConvertCardFromProto(protoPair.ReturnCard),
	}
}

// ConvertPlayerViewToProto converts SDK PlayerView to proto PlayerView
func ConvertPlayerViewToProto(
	sdkView *PlayerView,
	matchID string,
	dealIndex int,
	seq int64,
) *viewpb.PlayerView {
	if sdkView == nil {
		return nil
	}

	protoView := &viewpb.PlayerView{
		// Metadata fields
		MatchId:     matchID,
		DealIndex:   int32(dealIndex),
		Seq:         seq,
		UpdatedAtMs: time.Now().UnixMilli(),

		// Basic info
		PlayerSeat: int32(sdkView.PlayerSeat),
		PlayerCards: ConvertCardsToProto(sdkView.PlayerCards),
		TeamLevels:  []int32{int32(sdkView.TeamLevels[0]), int32(sdkView.TeamLevels[1])},
		DealLevel:   int32(sdkView.DealLevel),
		DealStatus:  ConvertDealStatusToProto(sdkView.DealStatus),
	}

	// Conditional fields based on DealStatus
	if sdkView.DealStatus == DealStatusPlaying {
		// Fill playing phase fields
		if sdkView.CurrentTurn != nil {
			currentTurn := int32(*sdkView.CurrentTurn)
			protoView.CurrentTurn = &currentTurn
		}

		// Set leader from sdkView.Leader (authoritative source)
		if sdkView.Leader != nil {
			leader := int32(*sdkView.Leader)
			protoView.Leader = &leader
		}

		// Convert plays
		protoPlays := make([]*commonpb.PlayAction, len(sdkView.Plays))
		for i, play := range sdkView.Plays {
			protoPlays[i] = ConvertPlayActionToProto(play)
		}
		protoView.Plays = protoPlays
	}
	// Note: TributePhase is not included in PlayerView proto
	// Use GetTributeView() separately for tribute phase info

	return protoView
}

// ConvertTributeViewToProto converts SDK TributePhase to proto TributeView
func ConvertTributeViewToProto(
	phase *TributePhase,
	matchID string,
	dealIndex int,
	seq int64,
) *viewpb.TributeView {
	if phase == nil {
		return nil
	}

	protoView := &viewpb.TributeView{
		// Metadata fields
		MatchId:     matchID,
		DealIndex:   int32(dealIndex),
		Seq:         seq,
		UpdatedAtMs: time.Now().UnixMilli(),

		// Tribute phase info
		Status:          ConvertTributeStatusToProto(phase.Status),
		SelectingPlayer: int32(phase.SelectingPlayer),
		IsImmune:        phase.IsImmune,
		PoolCards:       ConvertCardsToProto(phase.PoolCards),
	}

	// Convert tribute pairs
	protoPairs := make([]*viewpb.TributePair, len(phase.TributePairs))
	for i, pair := range phase.TributePairs {
		protoPairs[i] = ConvertTributePairToProto(pair)
	}
	protoView.TributePairs = protoPairs

	return protoView
}
