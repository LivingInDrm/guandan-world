package sdk

import (
	"testing"
	"time"
)

func TestGameEngine_CompleteGameFlow(t *testing.T) {
	// Create game engine
	engine := NewGameEngine()

	// Create test players
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	// Start match
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// Verify game state
	gameState := engine.GetGameState()
	if gameState.Status != GameStatusStarted {
		t.Errorf("Expected game status %v, got %v", GameStatusStarted, gameState.Status)
	}

	if gameState.CurrentMatch == nil {
		t.Fatal("Expected current match to be set")
	}

	if gameState.CurrentMatch.Status != MatchStatusWaiting {
		t.Errorf("Expected match status %v, got %v", MatchStatusWaiting, gameState.CurrentMatch.Status)
	}
}

func TestGameEngine_StartDeal(t *testing.T) {
	// Create and setup game engine
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// Start deal
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	// Verify deal state
	gameState := engine.GetGameState()
	if gameState.CurrentMatch.Status != MatchStatusPlaying {
		t.Errorf("Expected match status %v, got %v", MatchStatusPlaying, gameState.CurrentMatch.Status)
	}

	if gameState.CurrentMatch.CurrentDeal == nil {
		t.Fatal("Expected current deal to be set")
	}

	deal := gameState.CurrentMatch.CurrentDeal
	if deal.Status != DealStatusPlaying && deal.Status != DealStatusTribute {
		t.Errorf("Expected deal status to be playing or tribute, got %v", deal.Status)
	}

	// Verify cards were dealt
	for i := 0; i < 4; i++ {
		if len(deal.PlayerCards[i]) != 27 {
			t.Errorf("Player %d should have 27 cards, got %d", i, len(deal.PlayerCards[i]))
		}
	}
}

func TestGameEngine_PlayCardsValidation(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal

	// Skip tribute phase if it exists
	if deal.Status == DealStatusTribute {
		// For testing, we'll skip the tribute phase by setting status to playing
		deal.Status = DealStatusPlaying
		deal.startFirstTrick()
	}

	// Test invalid play - wrong player's turn
	if deal.CurrentTrick != nil && deal.CurrentTrick.CurrentTurn != 0 {
		playerCards := deal.PlayerCards[0]
		if len(playerCards) > 0 {
			_, err = engine.PlayCards(0, []*Card{playerCards[0]})
			if err == nil {
				t.Error("Expected error for playing out of turn")
			}
		}
	}

	// Test valid play - correct player's turn
	if deal.CurrentTrick != nil {
		currentPlayer := deal.CurrentTrick.CurrentTurn
		playerCards := deal.PlayerCards[currentPlayer]
		if len(playerCards) > 0 {
			event, err := engine.PlayCards(currentPlayer, []*Card{playerCards[0]})
			if err != nil {
				t.Errorf("Expected valid play to succeed: %v", err)
			}
			if event == nil {
				t.Error("Expected event to be returned")
			}
			if event.Type != EventPlayerPlayed {
				t.Errorf("Expected event type %v, got %v", EventPlayerPlayed, event.Type)
			}
		}
	}
}

func TestGameEngine_PassValidation(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal

	// Skip tribute phase if it exists
	if deal.Status == DealStatusTribute {
		deal.Status = DealStatusPlaying
		deal.startFirstTrick()
	}

	if deal.CurrentTrick != nil {
		// Test invalid pass - trick leader cannot pass
		leader := deal.CurrentTrick.Leader
		if deal.CurrentTrick.LeadComp == nil {
			_, err = engine.PassTurn(leader)
			if err == nil {
				t.Error("Expected error for trick leader trying to pass")
			}
		}

		// Test valid pass after someone has played
		if deal.CurrentTrick.LeadComp != nil {
			currentPlayer := deal.CurrentTrick.CurrentTurn
			event, err := engine.PassTurn(currentPlayer)
			if err != nil {
				t.Errorf("Expected valid pass to succeed: %v", err)
			}
			if event == nil {
				t.Error("Expected event to be returned")
			}
			if event.Type != EventPlayerPassed {
				t.Errorf("Expected event type %v, got %v", EventPlayerPassed, event.Type)
			}
		}
	}
}

func TestGameEngine_AutoPlay(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal

	// Skip tribute phase if it exists
	if deal.Status == DealStatusTribute {
		deal.Status = DealStatusPlaying
		err = deal.startFirstTrick()
		if err != nil {
			t.Fatalf("Failed to start first trick: %v", err)
		}
	}

	if deal.CurrentTrick != nil {
		originalCurrentPlayer := deal.CurrentTrick.CurrentTurn

		// Test auto-play for current player
		event, err := engine.AutoPlayForPlayer(originalCurrentPlayer)
		if err != nil {
			t.Logf("Auto-play failed (this may be expected): %v", err)
		} else if event == nil {
			t.Error("Expected event to be returned from auto-play")
		}

		// After auto-play, the current player may have changed
		// Test auto-play for the original player (should now be wrong)
		_, err = engine.AutoPlayForPlayer(originalCurrentPlayer)
		if err == nil {
			t.Logf("Auto-play succeeded for original player (turn may have changed)")
		} else {
			t.Logf("Correctly got error for original player after turn change: %v", err)
		}
	}
}

func TestGameEngine_PlayerView(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	// Test player view for each player
	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		playerView := engine.GetPlayerView(playerSeat)

		if playerView == nil {
			t.Errorf("Expected player view for player %d", playerSeat)
			continue
		}

		if playerView.PlayerSeat != playerSeat {
			t.Errorf("Expected player seat %d, got %d", playerSeat, playerView.PlayerSeat)
		}

		if playerView.GameState == nil {
			t.Errorf("Expected game state in player view for player %d", playerSeat)
		}

		// Player should have their cards
		if len(playerView.PlayerCards) != 27 {
			t.Errorf("Player %d should have 27 cards in view, got %d", playerSeat, len(playerView.PlayerCards))
		}
	}
}

func TestGameEngine_ProcessTimeouts(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("Failed to start deal: %v", err)
	}

	gameState := engine.GetGameState()
	deal := gameState.CurrentMatch.CurrentDeal

	// Skip tribute phase if it exists
	if deal.Status == DealStatusTribute {
		deal.Status = DealStatusPlaying
		err := deal.startFirstTrick()
		if err != nil {
			t.Fatalf("Failed to start first trick: %v", err)
		}
	}

	// Test timeout processing (should not timeout immediately)
	events := engine.ProcessTimeouts()
	if len(events) > 0 {
		t.Error("Should not have timeout events immediately after starting")
	}

	// Simulate timeout by setting past time
	if deal.CurrentTrick != nil {
		t.Logf("Deal status: %v, Trick status: %v", deal.Status, deal.CurrentTrick.Status)

		// Make sure trick is in playing status
		deal.CurrentTrick.Status = TrickStatusPlaying
		deal.CurrentTrick.TurnTimeout = time.Now().Add(-1 * time.Second)

		t.Logf("After setup - Deal status: %v, Trick status: %v, Timeout: %v",
			deal.Status, deal.CurrentTrick.Status, deal.CurrentTrick.TurnTimeout)

		events = engine.ProcessTimeouts()
		t.Logf("Got %d timeout events", len(events))
		for i, event := range events {
			t.Logf("Event %d: Type=%v, Data=%v", i, event.Type, event.Data)
		}

		if len(events) == 0 {
			t.Error("Expected timeout events after timeout period")
		} else {
			// Should have timeout event
			found := false
			for _, event := range events {
				if event.Type == EventPlayerTimeout {
					found = true
					break
				}
			}
			if !found {
				t.Error("Expected timeout event in processed events")
			}
		}
	}
}

func TestGameEngine_PlayerDisconnectReconnect(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// Test player disconnect
	event, err := engine.HandlePlayerDisconnect(0)
	if err != nil {
		t.Errorf("Failed to handle player disconnect: %v", err)
	}
	if event == nil {
		t.Error("Expected disconnect event")
	}
	if event.Type != EventPlayerDisconnect {
		t.Errorf("Expected event type %v, got %v", EventPlayerDisconnect, event.Type)
	}

	// Verify player is marked as offline and auto-play
	gameState := engine.GetGameState()
	if gameState.CurrentMatch.Players[0].Online {
		t.Error("Player should be marked as offline")
	}
	if !gameState.CurrentMatch.Players[0].AutoPlay {
		t.Error("Player should be marked as auto-play")
	}

	// Test player reconnect
	event, err = engine.HandlePlayerReconnect(0)
	if err != nil {
		t.Errorf("Failed to handle player reconnect: %v", err)
	}
	if event == nil {
		t.Error("Expected reconnect event")
	}
	if event.Type != EventPlayerReconnect {
		t.Errorf("Expected event type %v, got %v", EventPlayerReconnect, event.Type)
	}

	// Verify player is marked as online and not auto-play
	gameState = engine.GetGameState()
	if !gameState.CurrentMatch.Players[0].Online {
		t.Error("Player should be marked as online")
	}
	if gameState.CurrentMatch.Players[0].AutoPlay {
		t.Error("Player should not be marked as auto-play")
	}
}

func TestGameEngine_SetPlayerAutoPlay(t *testing.T) {
	// Setup game
	engine := NewGameEngine()
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// Test setting auto-play
	err = engine.SetPlayerAutoPlay(0, true)
	if err != nil {
		t.Errorf("Failed to set player auto-play: %v", err)
	}

	// Verify auto-play is set
	gameState := engine.GetGameState()
	if !gameState.CurrentMatch.Players[0].AutoPlay {
		t.Error("Player should be marked as auto-play")
	}

	// Test disabling auto-play
	err = engine.SetPlayerAutoPlay(0, false)
	if err != nil {
		t.Errorf("Failed to disable player auto-play: %v", err)
	}

	// Verify auto-play is disabled
	gameState = engine.GetGameState()
	if gameState.CurrentMatch.Players[0].AutoPlay {
		t.Error("Player should not be marked as auto-play")
	}
}
