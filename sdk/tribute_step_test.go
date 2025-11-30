package sdk

import (
	"testing"
)

func TestProcessTributeStep_Waiting_ToSelecting(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusWaiting,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{
		{{Number: 14, Color: "Spade", DeckIndex: 1}},
		{{Number: 13, Color: "Heart", DeckIndex: 2}},
		{{Number: 12, Color: "Diamond", DeckIndex: 3}},
		{{Number: 16, Color: "Joker", DeckIndex: 4}, {Number: 10, Color: "Club", DeckIndex: 5}},
	}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.NextStatus != TributeStatusSelecting {
		t.Errorf("Expected NextStatus=Selecting, got %s", result.NextStatus)
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventTributeCardSubmitted {
		t.Errorf("Expected EventTributeCardSubmitted, got %v", result.Events[0].Type)
	}
	if result.Events[0].PlayerSeat != 3 {
		t.Errorf("Expected PlayerSeat=3, got %d", result.Events[0].PlayerSeat)
	}
	if len(result.HandChanges) != 1 {
		t.Errorf("Expected 1 HandChange, got %d", len(result.HandChanges))
	}
	if result.HandChanges[0].IsAdd {
		t.Error("Expected IsAdd=false for removing tribute card")
	}
	if len(result.PairUpdates) != 1 {
		t.Errorf("Expected 1 PairUpdate, got %d", len(result.PairUpdates))
	}
	if result.PoolCardsToSet == nil || len(result.PoolCardsToSet) != 1 {
		t.Errorf("Expected PoolCardsToSet with 1 card, got %v", result.PoolCardsToSet)
	}
}

func TestProcessTributeStep_Selecting_DoubleDown_NeedsInput(t *testing.T) {
	card1 := &Card{Number: 16, Color: "Joker", DeckIndex: 1}
	card2 := &Card{Number: 14, Color: "Spade", DeckIndex: 2}

	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributePairs: []*TributePair{
			{Giver: 2, Receiver: -1, TributeCard: card1, ReturnCard: nil},
			{Giver: 3, Receiver: -1, TributeCard: card2, ReturnCard: nil},
		},
		PoolCards: []*Card{card1, card2},
		Winners:   []int{0, 1},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.PendingAction == nil {
		t.Fatal("Expected PendingAction, got nil")
	}
	if result.PendingAction.Type != TributeActionSelect {
		t.Errorf("Expected TributeActionSelect, got %v", result.PendingAction.Type)
	}
	if result.PendingAction.PlayerID != 0 {
		t.Errorf("Expected PlayerID=0 (first winner), got %d", result.PendingAction.PlayerID)
	}
	if len(result.PendingAction.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(result.PendingAction.Options))
	}
}

func TestProcessTributeStep_Selecting_DoubleDown_WithInput(t *testing.T) {
	card1 := &Card{Number: 16, Color: "Joker", DeckIndex: 1}
	card2 := &Card{Number: 14, Color: "Spade", DeckIndex: 2}

	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributePairs: []*TributePair{
			{Giver: 2, Receiver: -1, TributeCard: card1, ReturnCard: nil},
			{Giver: 3, Receiver: -1, TributeCard: card2, ReturnCard: nil},
		},
		PoolCards: []*Card{card1, card2},
		Winners:   []int{0, 1},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	input := &TributeInput{
		PlayerID: 0,
		Card:     &Card{DeckIndex: 1},
	}

	result := ProcessTributeStep(phase, playerHands, 7, input)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventTributeCardSelected {
		t.Errorf("Expected EventTributeCardSelected, got %v", result.Events[0].Type)
	}
	if result.Events[0].IsAuto {
		t.Error("Expected IsAuto=false for user selection")
	}
	if len(result.HandChanges) != 1 {
		t.Errorf("Expected 1 HandChange, got %d", len(result.HandChanges))
	}
	if !result.HandChanges[0].IsAdd {
		t.Error("Expected IsAdd=true for adding card to receiver")
	}
	if result.PoolCardToRemove == nil {
		t.Error("Expected PoolCardToRemove to be set")
	}
}

func TestProcessTributeStep_Selecting_SingleTribute_AutoAssign(t *testing.T) {
	card := &Card{Number: 16, Color: "Joker", DeckIndex: 1}

	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: -1, TributeCard: card, ReturnCard: nil},
		},
		PoolCards: []*Card{card},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventTributeCardSelected {
		t.Errorf("Expected EventTributeCardSelected, got %v", result.Events[0].Type)
	}
	if !result.Events[0].IsAuto {
		t.Error("Expected IsAuto=true for auto assignment")
	}
	if result.Events[0].PlayerSeat != 0 {
		t.Errorf("Expected receiver=0, got %d", result.Events[0].PlayerSeat)
	}
}

func TestProcessTributeStep_Selecting_ToReturning(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.NextStatus != TributeStatusReturning {
		t.Errorf("Expected NextStatus=Returning, got %s", result.NextStatus)
	}
}

func TestProcessTributeStep_Returning_NeedsInput(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{
		{{Number: 5, Color: "Heart", DeckIndex: 10}},
		{},
		{},
		{},
	}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.PendingAction == nil {
		t.Fatal("Expected PendingAction, got nil")
	}
	if result.PendingAction.Type != TributeActionReturn {
		t.Errorf("Expected TributeActionReturn, got %v", result.PendingAction.Type)
	}
	if result.PendingAction.PlayerID != 0 {
		t.Errorf("Expected PlayerID=0, got %d", result.PendingAction.PlayerID)
	}
	if result.PendingAction.TargetPlayer != 3 {
		t.Errorf("Expected TargetPlayer=3, got %d", result.PendingAction.TargetPlayer)
	}
}

func TestProcessTributeStep_Returning_WithInput(t *testing.T) {
	returnCard := &Card{Number: 5, Color: "Heart", DeckIndex: 10}

	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{
		{returnCard},
		{},
		{},
		{},
	}

	input := &TributeInput{
		PlayerID: 0,
		Card:     &Card{DeckIndex: 10},
	}

	result := ProcessTributeStep(phase, playerHands, 7, input)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventReturnTribute {
		t.Errorf("Expected EventReturnTribute, got %v", result.Events[0].Type)
	}
	if len(result.HandChanges) != 2 {
		t.Errorf("Expected 2 HandChanges (remove from receiver, add to giver), got %d", len(result.HandChanges))
	}
	if len(result.PairUpdates) != 1 {
		t.Errorf("Expected 1 PairUpdate, got %d", len(result.PairUpdates))
	}
	if result.PairUpdates[0].ReturnCard == nil {
		t.Error("Expected ReturnCard to be set in PairUpdate")
	}
}

func TestProcessTributeStep_Returning_AllDone(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: &Card{DeckIndex: 10}},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.NextStatus != TributeStatusFinished {
		t.Errorf("Expected NextStatus=Finished, got %s", result.NextStatus)
	}
}

func TestProcessTributeStep_Finished_Valid(t *testing.T) {
	phase := &TributePhase{
		Status:   TributeStatusFinished,
		IsImmune: false,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: &Card{DeckIndex: 10}},
		},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if !result.PhaseCompleted {
		t.Error("Expected PhaseCompleted=true")
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventTributeCompleted {
		t.Errorf("Expected EventTributeCompleted, got %v", result.Events[0].Type)
	}
}

func TestProcessTributeStep_Finished_InvalidPair(t *testing.T) {
	phase := &TributePhase{
		Status:   TributeStatusFinished,
		IsImmune: false,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
		},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error == nil {
		t.Error("Expected error for invalid pair, got nil")
	}
	if result.PhaseCompleted {
		t.Error("Expected PhaseCompleted=false when error")
	}
}

func TestProcessTributeStep_Finished_Immune(t *testing.T) {
	phase := &TributePhase{
		Status:   TributeStatusFinished,
		IsImmune: true,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: -1, TributeCard: nil, ReturnCard: nil},
		},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error for immune phase: %v", result.Error)
	}
	if !result.PhaseCompleted {
		t.Error("Expected PhaseCompleted=true for immune phase")
	}
}

func TestProcessTributeStep_Selecting_InvalidCard(t *testing.T) {
	card1 := &Card{Number: 16, Color: "Joker", DeckIndex: 1}
	card2 := &Card{Number: 14, Color: "Spade", DeckIndex: 2}

	phase := &TributePhase{
		Status: TributeStatusSelecting,
		TributePairs: []*TributePair{
			{Giver: 2, Receiver: -1, TributeCard: card1, ReturnCard: nil},
			{Giver: 3, Receiver: -1, TributeCard: card2, ReturnCard: nil},
		},
		PoolCards: []*Card{card1, card2},
		Winners:   []int{0, 1},
	}

	playerHands := [4][]*Card{{}, {}, {}, {}}

	input := &TributeInput{
		PlayerID: 0,
		Card:     &Card{DeckIndex: 999},
	}

	result := ProcessTributeStep(phase, playerHands, 7, input)

	if result.Error == nil {
		t.Error("Expected error for invalid card selection, got nil")
	}
}

func TestProcessTributeStep_Returning_DoubleDown_NeedsInput(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 2, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
			{Giver: 3, Receiver: 1, TributeCard: &Card{DeckIndex: 2}, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0, 1},
	}

	playerHands := [4][]*Card{
		{{Number: 5, Color: "Heart", DeckIndex: 10}},
		{{Number: 6, Color: "Club", DeckIndex: 20}},
		{},
		{},
	}

	result := ProcessTributeStep(phase, playerHands, 7, nil)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.PendingAction == nil {
		t.Fatal("Expected PendingAction, got nil")
	}
	if result.PendingAction.Type != TributeActionReturn {
		t.Errorf("Expected TributeActionReturn, got %v", result.PendingAction.Type)
	}
	if result.PendingAction.PlayerID != 0 {
		t.Errorf("Expected PlayerID=0 (first pending receiver), got %d", result.PendingAction.PlayerID)
	}
	if result.PendingAction.TargetPlayer != 2 {
		t.Errorf("Expected TargetPlayer=2, got %d", result.PendingAction.TargetPlayer)
	}
}

func TestProcessTributeStep_Returning_DoubleDown_WithInput(t *testing.T) {
	returnCard := &Card{Number: 5, Color: "Heart", DeckIndex: 10}

	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 2, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: nil},
			{Giver: 3, Receiver: 1, TributeCard: &Card{DeckIndex: 2}, ReturnCard: nil},
		},
		PoolCards: []*Card{},
		Winners:   []int{0, 1},
	}

	playerHands := [4][]*Card{
		{returnCard},
		{{Number: 6, Color: "Club", DeckIndex: 20}},
		{},
		{},
	}

	input := &TributeInput{
		PlayerID: 0,
		Card:     &Card{DeckIndex: 10},
	}

	result := ProcessTributeStep(phase, playerHands, 7, input)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Type != EventReturnTribute {
		t.Errorf("Expected EventReturnTribute, got %v", result.Events[0].Type)
	}
	if len(result.HandChanges) != 2 {
		t.Errorf("Expected 2 HandChanges, got %d", len(result.HandChanges))
	}
	if len(result.PairUpdates) != 1 {
		t.Errorf("Expected 1 PairUpdate, got %d", len(result.PairUpdates))
	}
	if result.PairUpdates[0].GiverSeat != 2 {
		t.Errorf("Expected GiverSeat=2, got %d", result.PairUpdates[0].GiverSeat)
	}
	if result.PairUpdates[0].ReturnCard == nil {
		t.Error("Expected ReturnCard to be set in PairUpdate")
	}
	if result.NextStatus != "" {
		t.Errorf("Expected NextStatus to be empty (still have pending returns), got %s", result.NextStatus)
	}
	if result.PendingAction != nil {
		t.Error("Expected PendingAction=nil (next iteration will handle)")
	}
}

func TestProcessTributeStep_Returning_AllDone_WithInput(t *testing.T) {
	phase := &TributePhase{
		Status: TributeStatusReturning,
		TributePairs: []*TributePair{
			{Giver: 3, Receiver: 0, TributeCard: &Card{DeckIndex: 1}, ReturnCard: &Card{DeckIndex: 10}},
		},
		PoolCards: []*Card{},
		Winners:   []int{0},
	}

	playerHands := [4][]*Card{
		{{Number: 7, Color: "Diamond", DeckIndex: 20}},
		{},
		{},
		{},
	}

	input := &TributeInput{
		PlayerID: 0,
		Card:     &Card{DeckIndex: 20},
	}

	result := ProcessTributeStep(phase, playerHands, 7, input)

	if result.Error == nil {
		t.Error("Expected error when no pending return but input provided, got nil")
	}
}
