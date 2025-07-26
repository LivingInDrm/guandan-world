package sdk

import (
	"testing"
)

func TestDetermineFirstPlayer(t *testing.T) {
	t.Run("First deal - random player", func(t *testing.T) {
		// Create deal with no LastResult (first deal)
		deal, err := NewDeal(2, nil)
		if err != nil {
			t.Fatalf("Failed to create deal: %v", err)
		}
		
		// Test multiple times to ensure randomness
		results := make(map[int]int)
		for i := 0; i < 100; i++ {
			firstPlayer := deal.determineFirstPlayer()
			if firstPlayer < 0 || firstPlayer > 3 {
				t.Errorf("Invalid first player: %d", firstPlayer)
			}
			results[firstPlayer]++
		}
		
		// Verify all players can be selected
		if len(results) < 2 {
			t.Errorf("Random selection not working properly, only got players: %v", results)
		}
	})
	
	t.Run("Anti-tribute - rank1 starts", func(t *testing.T) {
		// Create last result with rankings
		lastResult := &DealResult{
			Rankings: []int{2, 0, 1, 3}, // rank1 is player 2
			VictoryType: VictoryTypeSingleLast,
		}
		
		deal, err := NewDeal(3, lastResult)
		if err != nil {
			t.Fatalf("Failed to create deal: %v", err)
		}
		
		// Simulate immunity
		if deal.TributePhase != nil {
			deal.TributePhase.Status = TributeStatusFinished
			deal.TributePhase.IsImmune = true
		}
		
		firstPlayer := deal.determineFirstPlayer()
		if firstPlayer != 2 {
			t.Errorf("Expected player 2 (rank1) to start, got %d", firstPlayer)
		}
	})
	
	t.Run("Double Down - bigger tribute card player starts", func(t *testing.T) {
		// Create last result with Double Down victory
		lastResult := &DealResult{
			Rankings: []int{0, 2, 1, 3}, // rank3=1, rank4=3
			VictoryType: VictoryTypeDoubleDown,
		}
		
		deal, err := NewDeal(4, lastResult)
		if err != nil {
			t.Fatalf("Failed to create deal: %v", err)
		}
		
		// Create tribute phase with cards
		if deal.TributePhase != nil {
			deal.TributePhase.Status = TributeStatusFinished
			deal.TributePhase.TributeCards = make(map[int]*Card)
			// Rank3 (player 1) gives a King
			deal.TributePhase.TributeCards[1], _ = NewCard(13, "Hearts", 4)
			// Rank4 (player 3) gives a Jack
			deal.TributePhase.TributeCards[3], _ = NewCard(11, "Spades", 4)
		}
		
		firstPlayer := deal.determineFirstPlayer()
		if firstPlayer != 1 {
			t.Errorf("Expected player 1 (bigger tribute) to start, got %d", firstPlayer)
		}
	})
	
	t.Run("Single Last - rank4 starts", func(t *testing.T) {
		lastResult := &DealResult{
			Rankings: []int{0, 2, 1, 3}, // rank4 is player 3
			VictoryType: VictoryTypeSingleLast,
		}
		
		deal, err := NewDeal(5, lastResult)
		if err != nil {
			t.Fatalf("Failed to create deal: %v", err)
		}
		
		if deal.TributePhase != nil {
			deal.TributePhase.Status = TributeStatusFinished
		}
		
		firstPlayer := deal.determineFirstPlayer()
		if firstPlayer != 3 {
			t.Errorf("Expected player 3 (rank4) to start, got %d", firstPlayer)
		}
	})
	
	t.Run("Partner Last - rank3 starts", func(t *testing.T) {
		lastResult := &DealResult{
			Rankings: []int{0, 3, 2, 1}, // rank3 is player 2
			VictoryType: VictoryTypePartnerLast,
		}
		
		deal, err := NewDeal(6, lastResult)
		if err != nil {
			t.Fatalf("Failed to create deal: %v", err)
		}
		
		if deal.TributePhase != nil {
			deal.TributePhase.Status = TributeStatusFinished
		}
		
		firstPlayer := deal.determineFirstPlayer()
		if firstPlayer != 2 {
			t.Errorf("Expected player 2 (rank3) to start, got %d", firstPlayer)
		}
	})
}

func TestTrickNextLeaderWithTeammate(t *testing.T) {
	deal, err := NewDeal(2, nil)
	if err != nil {
		t.Fatalf("Failed to create deal: %v", err)
	}
	
	// Deal cards
	err = deal.dealCards()
	if err != nil {
		t.Fatalf("Failed to deal cards: %v", err)
	}
	
	// Start first trick
	deal.CurrentTrick, _ = NewTrick(0)
	deal.CurrentTrick.Status = TrickStatusPlaying
	
	// Simulate player 0 playing and winning
	singleCard := []*Card{deal.PlayerCards[0][0]}
	play0 := &PlayAction{
		PlayerSeat: 0,
		Cards: singleCard,
		Comp: NewSingle(singleCard),
		IsPass: false,
	}
	deal.CurrentTrick.Plays = append(deal.CurrentTrick.Plays, play0)
	deal.CurrentTrick.LeadComp = play0.Comp
	deal.CurrentTrick.Leader = 0
	deal.CurrentTrick.Winner = 0
	
	// Remove all cards from player 0 to simulate finishing
	deal.PlayerCards[0] = []*Card{}
	
	// Other players pass
	for i := 1; i <= 3; i++ {
		pass := &PlayAction{
			PlayerSeat: i,
			IsPass: true,
		}
		deal.CurrentTrick.Plays = append(deal.CurrentTrick.Plays, pass)
	}
	
	// Mark trick as finished
	deal.CurrentTrick.Status = TrickStatusFinished
	
	// Call finishCurrentTrick
	err = deal.finishCurrentTrick()
	if err != nil {
		t.Fatalf("Failed to finish trick: %v", err)
	}
	
	// Verify that player 2 (teammate of player 0) is the next leader
	if deal.CurrentTrick.NextLeader != 2 {
		t.Errorf("Expected player 2 (teammate) to be next leader, got %d", deal.CurrentTrick.NextLeader)
	}
}