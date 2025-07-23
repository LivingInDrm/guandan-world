package sdk

import (
	"testing"
	"time"
)

func TestNewTrick(t *testing.T) {
	trick, err := NewTrick(0)
	if err != nil {
		t.Fatalf("NewTrick failed: %v", err)
	}

	if trick == nil {
		t.Fatal("NewTrick should return a non-nil trick")
	}

	if trick.ID == "" {
		t.Error("Trick should have a non-empty ID")
	}

	if trick.Leader != 0 {
		t.Errorf("Trick leader should be 0, got %d", trick.Leader)
	}

	if trick.CurrentTurn != 0 {
		t.Errorf("Trick current turn should be 0, got %d", trick.CurrentTurn)
	}

	if trick.Status != TrickStatusWaiting {
		t.Errorf("New trick should have status %v, got %v", TrickStatusWaiting, trick.Status)
	}

	if trick.Winner != -1 {
		t.Error("New trick should have no winner")
	}

	if trick.LeadComp != nil {
		t.Error("New trick should have no lead combination")
	}

	if len(trick.Plays) != 0 {
		t.Error("New trick should have no plays")
	}
}

func TestNewTrickValidation(t *testing.T) {
	// Test with invalid leader seats
	_, err := NewTrick(-1)
	if err == nil {
		t.Error("NewTrick should fail with invalid leader seat -1")
	}

	_, err = NewTrick(4)
	if err == nil {
		t.Error("NewTrick should fail with invalid leader seat 4")
	}
}

func TestTrickStartTrick(t *testing.T) {
	trick, _ := NewTrick(1)

	err := trick.StartTrick()
	if err != nil {
		t.Errorf("StartTrick failed: %v", err)
	}

	if trick.Status != TrickStatusPlaying {
		t.Errorf("Trick should have status %v after start, got %v", TrickStatusPlaying, trick.Status)
	}

	// Test starting already started trick
	err = trick.StartTrick()
	if err == nil {
		t.Error("StartTrick should fail when trick is already playing")
	}
}

func TestTrickPlayCards(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Create test cards and combination
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)

	// Test first play (leader)
	err := trick.PlayCards(0, cards, comp)
	if err != nil {
		t.Errorf("PlayCards failed for leader: %v", err)
	}

	if len(trick.Plays) != 1 {
		t.Error("Trick should have 1 play after first play")
	}

	if trick.LeadComp == nil {
		t.Error("Trick should have lead combination after first play")
	}

	if trick.Leader != 0 {
		t.Error("Leader should still be 0 after first play")
	}

	if trick.CurrentTurn != 1 {
		t.Error("Current turn should be 1 after player 0 plays")
	}

	// Test second play (higher card)
	higherCards := []*Card{}
	higherCard, _ := NewCard(13, "Heart", 5) // King
	higherCards = append(higherCards, higherCard)
	higherComp := NewSingle(higherCards)

	err = trick.PlayCards(1, higherCards, higherComp)
	if err != nil {
		t.Errorf("PlayCards failed for second player: %v", err)
	}

	if trick.Leader != 1 {
		t.Error("Leader should be 1 after higher play")
	}

	if trick.CurrentTurn != 2 {
		t.Error("Current turn should be 2 after player 1 plays")
	}
}

func TestTrickPlayCardsValidation(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)

	// Test playing out of turn
	err := trick.PlayCards(1, cards, comp)
	if err == nil {
		t.Error("PlayCards should fail when not player's turn")
	}

	// Test with invalid combination
	err = trick.PlayCards(0, cards, nil)
	if err == nil {
		t.Error("PlayCards should fail with nil combination")
	}

	// Test with invalid combination that can't be played
	invalidComp := &IllegalComp{BaseComp: BaseComp{Cards: cards, Valid: false, Type: TypeIllegal}}
	err = trick.PlayCards(0, cards, invalidComp)
	if err == nil {
		t.Error("PlayCards should fail with invalid combination")
	}
}

func TestTrickPassTurn(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Leader cannot pass without playing first
	err := trick.PassTurn(0)
	if err == nil {
		t.Error("Leader should not be able to pass without playing first")
	}

	// Leader plays first
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)
	trick.PlayCards(0, cards, comp)

	// Next player passes
	err = trick.PassTurn(1)
	if err != nil {
		t.Errorf("PassTurn failed: %v", err)
	}

	if len(trick.Plays) != 2 {
		t.Error("Trick should have 2 plays after pass")
	}

	lastPlay := trick.Plays[len(trick.Plays)-1]
	if !lastPlay.IsPass {
		t.Error("Last play should be a pass")
	}

	if lastPlay.PlayerSeat != 1 {
		t.Error("Pass should be from player 1")
	}

	if trick.CurrentTurn != 2 {
		t.Error("Current turn should be 2 after player 1 passes")
	}
}

func TestTrickPassTurnValidation(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Test passing out of turn
	err := trick.PassTurn(1)
	if err == nil {
		t.Error("PassTurn should fail when not player's turn")
	}
}

func TestTrickGetMethods(t *testing.T) {
	trick, _ := NewTrick(2)

	if trick.GetWinner() != -1 {
		t.Error("Unfinished trick should have no winner")
	}

	if trick.GetLeadingPlayer() != 2 {
		t.Error("Leading player should be 2")
	}

	if trick.GetCurrentPlayer() != 2 {
		t.Error("Current player should be 2")
	}

	if trick.GetLeadingCombination() != nil {
		t.Error("New trick should have no leading combination")
	}

	if len(trick.GetPlays()) != 0 {
		t.Error("New trick should have no plays")
	}

	if trick.IsFinished() {
		t.Error("New trick should not be finished")
	}
}

func TestTrickHasPlayerPlayed(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Initially no one has played
	for i := 0; i < 4; i++ {
		if trick.HasPlayerPlayed(i) {
			t.Errorf("Player %d should not have played yet", i)
		}
	}

	// Player 0 plays
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)
	trick.PlayCards(0, cards, comp)

	if !trick.HasPlayerPlayed(0) {
		t.Error("Player 0 should have played")
	}

	for i := 1; i < 4; i++ {
		if trick.HasPlayerPlayed(i) {
			t.Errorf("Player %d should not have played yet", i)
		}
	}
}

func TestTrickGetPlayerPlay(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Initially no plays
	for i := 0; i < 4; i++ {
		if trick.GetPlayerPlay(i) != nil {
			t.Errorf("Player %d should have no play yet", i)
		}
	}

	// Player 0 plays
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)
	trick.PlayCards(0, cards, comp)

	play := trick.GetPlayerPlay(0)
	if play == nil {
		t.Error("Player 0 should have a play")
	}

	if play.PlayerSeat != 0 {
		t.Error("Play should be from player 0")
	}

	if play.IsPass {
		t.Error("Play should not be a pass")
	}

	// Other players still have no plays
	for i := 1; i < 4; i++ {
		if trick.GetPlayerPlay(i) != nil {
			t.Errorf("Player %d should have no play yet", i)
		}
	}
}

func TestTrickGetNextPlayer(t *testing.T) {
	trick, _ := NewTrick(0)

	if trick.getNextPlayer(0) != 1 {
		t.Error("Next player after 0 should be 1")
	}

	if trick.getNextPlayer(1) != 2 {
		t.Error("Next player after 1 should be 2")
	}

	if trick.getNextPlayer(2) != 3 {
		t.Error("Next player after 2 should be 3")
	}

	if trick.getNextPlayer(3) != 0 {
		t.Error("Next player after 3 should be 0")
	}
}

func TestTrickCanPlayCombination(t *testing.T) {
	trick, _ := NewTrick(0)

	// First play - anything is valid
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)

	if !trick.canPlayCombination(comp) {
		t.Error("First play should always be valid")
	}

	// Set lead combination
	trick.LeadComp = comp

	// Higher single should be valid
	higherCards := []*Card{}
	higherCard, _ := NewCard(13, "Heart", 5)
	higherCards = append(higherCards, higherCard)
	higherComp := NewSingle(higherCards)

	if !trick.canPlayCombination(higherComp) {
		t.Error("Higher single should be valid")
	}

	// Lower single should not be valid
	lowerCards := []*Card{}
	lowerCard, _ := NewCard(3, "Heart", 5) // 3 is definitely lower than 10
	lowerCards = append(lowerCards, lowerCard)
	lowerComp := NewSingle(lowerCards)

	if trick.canPlayCombination(lowerComp) {
		t.Error("Lower single should not be valid")
	}

	// Different type should not be valid (unless bomb)
	pairCards := []*Card{}
	pairCard1, _ := NewCard(13, "Heart", 5)
	pairCard2, _ := NewCard(13, "Spade", 5)
	pairCards = append(pairCards, pairCard1, pairCard2)
	pairComp := NewPair(pairCards)

	if trick.canPlayCombination(pairComp) {
		t.Error("Different type should not be valid")
	}
}

func TestTrickIsTrickFinished(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Initially not finished
	if trick.isTrickFinished() {
		t.Error("New trick should not be finished")
	}

	// Add some plays but not enough
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)

	trick.PlayCards(0, cards, comp) // Player 0 plays
	if trick.isTrickFinished() {
		t.Error("Trick with 1 play should not be finished")
	}

	trick.PassTurn(1) // Player 1 passes
	if trick.isTrickFinished() {
		t.Error("Trick with 2 plays should not be finished")
	}

	trick.PassTurn(2) // Player 2 passes
	if trick.isTrickFinished() {
		t.Error("Trick with 3 plays should not be finished")
	}

	trick.PassTurn(3) // Player 3 passes
	if !trick.isTrickFinished() {
		t.Error("Trick with 4 plays (3 passes after leader) should be finished")
	}
}

func TestTrickFinishTrick(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Set up a winning scenario
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)
	trick.PlayCards(0, cards, comp)

	// Manually finish trick
	err := trick.finishTrick()
	if err != nil {
		t.Errorf("finishTrick failed: %v", err)
	}

	if trick.Status != TrickStatusFinished {
		t.Error("Trick should be finished")
	}

	if trick.Winner != 0 {
		t.Error("Winner should be player 0")
	}

	if !trick.IsFinished() {
		t.Error("IsFinished should return true")
	}

	if trick.GetWinner() != 0 {
		t.Error("GetWinner should return 0")
	}

	// Test finishing already finished trick
	err = trick.finishTrick()
	if err == nil {
		t.Error("finishTrick should fail on already finished trick")
	}
}

func TestTrickProcessTimeout(t *testing.T) {
	trick, _ := NewTrick(0)
	trick.StartTrick()

	// Set up lead play
	cards := []*Card{}
	card, _ := NewCard(10, "Heart", 5)
	cards = append(cards, card)
	comp := NewSingle(cards)
	trick.PlayCards(0, cards, comp)

	// Test timeout not reached yet
	err := trick.ProcessTimeout()
	if err == nil {
		t.Error("ProcessTimeout should fail when timeout not reached")
	}

	// Set timeout to past
	trick.TurnTimeout = time.Now().Add(-1 * time.Second)

	// Test timeout processing
	err = trick.ProcessTimeout()
	if err != nil {
		t.Errorf("ProcessTimeout failed: %v", err)
	}

	// Should have added a pass play
	if len(trick.Plays) != 2 {
		t.Error("Should have 2 plays after timeout")
	}

	lastPlay := trick.Plays[len(trick.Plays)-1]
	if !lastPlay.IsPass {
		t.Error("Timeout should result in pass")
	}
}
