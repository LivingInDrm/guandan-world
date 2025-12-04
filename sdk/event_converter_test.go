package sdk

import (
	"testing"
	"time"

	commonpb "guandan-world/proto/common"
	eventpb "guandan-world/proto/event"
	viewpb "guandan-world/proto/view"
)

// TestConvertDealStatusToProto tests DealStatus enum conversion
func TestConvertDealStatusToProto(t *testing.T) {
	tests := []struct {
		sdk   DealStatus
		proto viewpb.DealStatus
	}{
		{DealStatusWaiting, viewpb.DealStatus_DEAL_STATUS_WAITING},
		{DealStatusDealing, viewpb.DealStatus_DEAL_STATUS_DEALING},
		{DealStatusTribute, viewpb.DealStatus_DEAL_STATUS_TRIBUTE},
		{DealStatusPlaying, viewpb.DealStatus_DEAL_STATUS_PLAYING},
		{DealStatusFinished, viewpb.DealStatus_DEAL_STATUS_FINISHED},
	}

	for _, tt := range tests {
		result := ConvertDealStatusToProto(tt.sdk)
		if result != tt.proto {
			t.Errorf("ConvertDealStatusToProto(%s) = %v, want %v", tt.sdk, result, tt.proto)
		}

		// Test reverse conversion
		reverse := ConvertDealStatusFromProto(tt.proto)
		if reverse != tt.sdk {
			t.Errorf("ConvertDealStatusFromProto(%v) = %s, want %s", tt.proto, reverse, tt.sdk)
		}
	}
}

// TestConvertTributeStatusToProto tests TributeStatus enum conversion
func TestConvertTributeStatusToProto(t *testing.T) {
	tests := []struct {
		sdk   TributeStatus
		proto viewpb.TributeStatus
	}{
		{TributeStatusWaiting, viewpb.TributeStatus_TRIBUTE_STATUS_WAITING},
		{TributeStatusSelecting, viewpb.TributeStatus_TRIBUTE_STATUS_SELECTING},
		{TributeStatusReturning, viewpb.TributeStatus_TRIBUTE_STATUS_RETURNING},
		{TributeStatusFinished, viewpb.TributeStatus_TRIBUTE_STATUS_FINISHED},
	}

	for _, tt := range tests {
		result := ConvertTributeStatusToProto(tt.sdk)
		if result != tt.proto {
			t.Errorf("ConvertTributeStatusToProto(%s) = %v, want %v", tt.sdk, result, tt.proto)
		}

		// Test reverse conversion
		reverse := ConvertTributeStatusFromProto(tt.proto)
		if reverse != tt.sdk {
			t.Errorf("ConvertTributeStatusFromProto(%v) = %s, want %s", tt.proto, reverse, tt.sdk)
		}
	}
}

// TestConvertCompTypeToProto tests CompType enum conversion
func TestConvertCompTypeToProto(t *testing.T) {
	tests := []struct {
		sdk   CompType
		proto commonpb.CompType
	}{
		{TypeFold, commonpb.CompType_COMP_TYPE_FOLD},
		{TypeIllegal, commonpb.CompType_COMP_TYPE_ILLEGAL},
		{TypeSingle, commonpb.CompType_COMP_TYPE_SINGLE},
		{TypePair, commonpb.CompType_COMP_TYPE_PAIR},
		{TypeTriple, commonpb.CompType_COMP_TYPE_TRIPLE},
		{TypeFullHouse, commonpb.CompType_COMP_TYPE_FULL_HOUSE},
		{TypeStraight, commonpb.CompType_COMP_TYPE_STRAIGHT},
		{TypePlate, commonpb.CompType_COMP_TYPE_PLATE},
		{TypeTube, commonpb.CompType_COMP_TYPE_TUBE},
		{TypeJokerBomb, commonpb.CompType_COMP_TYPE_JOKER_BOMB},
		{TypeNaiveBomb, commonpb.CompType_COMP_TYPE_NAIVE_BOMB},
		{TypeStraightFlush, commonpb.CompType_COMP_TYPE_STRAIGHT_FLUSH},
	}

	for _, tt := range tests {
		result := ConvertCompTypeToProto(tt.sdk)
		if result != tt.proto {
			t.Errorf("ConvertCompTypeToProto(%v) = %v, want %v", tt.sdk, result, tt.proto)
		}

		// Test reverse conversion
		reverse := ConvertCompTypeFromProto(tt.proto)
		if reverse != tt.sdk {
			t.Errorf("ConvertCompTypeFromProto(%v) = %v, want %v", tt.proto, reverse, tt.sdk)
		}
	}
}

// TestConvertPlayActionToProto tests PlayAction conversion
func TestConvertPlayActionToProto(t *testing.T) {
	now := time.Now()
	
	// Test with cards
	sdkAction := &PlayAction{
		PlayerSeat: 2,
		Cards: []*Card{
			{Number: 3, Color: "Heart", DeckIndex: 10},
			{Number: 3, Color: "Spade", DeckIndex: 11},
		},
		Timestamp: now,
		IsPass:    false,
	}

	protoAction := ConvertPlayActionToProto(sdkAction)
	if protoAction == nil {
		t.Fatal("ConvertPlayActionToProto returned nil")
	}

	if protoAction.PlayerSeat != 2 {
		t.Errorf("Expected PlayerSeat=2, got %d", protoAction.PlayerSeat)
	}

	if len(protoAction.Cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(protoAction.Cards))
	}

	if protoAction.IsPass != false {
		t.Errorf("Expected IsPass=false, got %v", protoAction.IsPass)
	}

	// Test pass action
	passAction := &PlayAction{
		PlayerSeat: 1,
		Cards:      nil,
		Timestamp:  now,
		IsPass:     true,
	}

	protoPass := ConvertPlayActionToProto(passAction)
	if protoPass.IsPass != true {
		t.Errorf("Expected IsPass=true for pass action")
	}

	if protoPass.CompType != commonpb.CompType_COMP_TYPE_FOLD {
		t.Errorf("Expected CompType=FOLD for pass action, got %v", protoPass.CompType)
	}

	// Test nil
	if ConvertPlayActionToProto(nil) != nil {
		t.Error("ConvertPlayActionToProto(nil) should return nil")
	}
}

// TestConvertTributePairToProto tests TributePair conversion
func TestConvertTributePairToProto(t *testing.T) {
	// Test with cards
	sdkPair := &TributePair{
		Giver:       3,
		Receiver:    0,
		TributeCard: &Card{Number: 14, Color: "Spade", DeckIndex: 50},
		ReturnCard:  &Card{Number: 2, Color: "Heart", DeckIndex: 1},
	}

	protoPair := ConvertTributePairToProto(sdkPair)
	if protoPair == nil {
		t.Fatal("ConvertTributePairToProto returned nil")
	}

	if protoPair.Giver != 3 {
		t.Errorf("Expected Giver=3, got %d", protoPair.Giver)
	}

	if protoPair.Receiver != 0 {
		t.Errorf("Expected Receiver=0, got %d", protoPair.Receiver)
	}

	if protoPair.TributeCard == nil {
		t.Error("Expected TributeCard to be non-nil")
	}

	if protoPair.ReturnCard == nil {
		t.Error("Expected ReturnCard to be non-nil")
	}

	// Test without return card
	sdkPairNoReturn := &TributePair{
		Giver:       3,
		Receiver:    -1,
		TributeCard: &Card{Number: 14, Color: "Spade", DeckIndex: 50},
		ReturnCard:  nil,
	}

	protoPairNoReturn := ConvertTributePairToProto(sdkPairNoReturn)
	if protoPairNoReturn.ReturnCard != nil {
		t.Error("Expected ReturnCard to be nil")
	}

	// Test nil
	if ConvertTributePairToProto(nil) != nil {
		t.Error("ConvertTributePairToProto(nil) should return nil")
	}
}

// TestConvertPlayerViewToProto tests PlayerView conversion
func TestConvertPlayerViewToProto(t *testing.T) {
	sdkView := &PlayerView{
		PlayerSeat: 1,
		PlayerCards: []*Card{
			{Number: 3, Color: "Heart", DeckIndex: 10},
		},
		TeamLevels: [2]int{5, 6},
		DealLevel:  5,
		DealStatus: DealStatusPlaying,
	}

	protoView := ConvertPlayerViewToProto(sdkView, "match123", 2, 100)
	if protoView == nil {
		t.Fatal("ConvertPlayerViewToProto returned nil")
	}

	// Check metadata
	if protoView.MatchId != "match123" {
		t.Errorf("Expected MatchId='match123', got '%s'", protoView.MatchId)
	}

	if protoView.DealIndex != 2 {
		t.Errorf("Expected DealIndex=2, got %d", protoView.DealIndex)
	}

	if protoView.Seq != 100 {
		t.Errorf("Expected Seq=100, got %d", protoView.Seq)
	}

	// Check basic fields
	if protoView.PlayerSeat != 1 {
		t.Errorf("Expected PlayerSeat=1, got %d", protoView.PlayerSeat)
	}

	if len(protoView.PlayerCards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(protoView.PlayerCards))
	}

	if len(protoView.TeamLevels) != 2 {
		t.Errorf("Expected 2 team levels, got %d", len(protoView.TeamLevels))
	}

	if protoView.DealLevel != 5 {
		t.Errorf("Expected DealLevel=5, got %d", protoView.DealLevel)
	}

	if protoView.DealStatus != viewpb.DealStatus_DEAL_STATUS_PLAYING {
		t.Errorf("Expected DealStatus=PLAYING, got %v", protoView.DealStatus)
	}

	// Test nil
	if ConvertPlayerViewToProto(nil, "match", 0, 1) != nil {
		t.Error("ConvertPlayerViewToProto(nil) should return nil")
	}
}

// TestConvertTributeViewToProto tests TributeView conversion
func TestConvertTributeViewToProto(t *testing.T) {
	sdkPhase := &TributePhase{
		Status:      TributeStatusSelecting,
		TributeType: "single_last",
		Givers:      []int{3},
		Receivers:   []int{0},
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: -1, TributeCard: &Card{Number: 14, Color: "Spade", DeckIndex: 50}},
		},
		PoolCards: []*Card{{Number: 14, Color: "Spade", DeckIndex: 50}},
		IsImmune:  false,
	}

	protoView := ConvertTributeViewToProto(sdkPhase, "match456", 1, 200)
	if protoView == nil {
		t.Fatal("ConvertTributeViewToProto returned nil")
	}

	// Check metadata
	if protoView.MatchId != "match456" {
		t.Errorf("Expected MatchId='match456', got '%s'", protoView.MatchId)
	}

	if protoView.DealIndex != 1 {
		t.Errorf("Expected DealIndex=1, got %d", protoView.DealIndex)
	}

	if protoView.Seq != 200 {
		t.Errorf("Expected Seq=200, got %d", protoView.Seq)
	}

	// Check tribute fields
	if protoView.Status != viewpb.TributeStatus_TRIBUTE_STATUS_SELECTING {
		t.Errorf("Expected Status=SELECTING, got %v", protoView.Status)
	}

	if len(protoView.TributePairs) != 1 {
		t.Errorf("Expected 1 tribute pair, got %d", len(protoView.TributePairs))
	}

	if len(protoView.PoolCards) != 1 {
		t.Errorf("Expected 1 pool card, got %d", len(protoView.PoolCards))
	}

	if protoView.IsImmune != false {
		t.Errorf("Expected IsImmune=false, got %v", protoView.IsImmune)
	}

	// Check new fields: TributeType, Givers, Receivers
	if protoView.TributeType != eventpb.TributeType_TRIBUTE_TYPE_SINGLE_LAST {
		t.Errorf("Expected TributeType=SINGLE_LAST, got %v", protoView.TributeType)
	}

	if len(protoView.Givers) != 1 || protoView.Givers[0] != 3 {
		t.Errorf("Expected Givers=[3], got %v", protoView.Givers)
	}

	if len(protoView.Receivers) != 1 || protoView.Receivers[0] != 0 {
		t.Errorf("Expected Receivers=[0], got %v", protoView.Receivers)
	}

	// Test nil
	if ConvertTributeViewToProto(nil, "match", 0, 1) != nil {
		t.Error("ConvertTributeViewToProto(nil) should return nil")
	}
}
