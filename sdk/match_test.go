package sdk

import (
	"testing"
	"time"
)

func TestNewMatch(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("NewMatch failed: %v", err)
	}
	
	if match == nil {
		t.Fatal("NewMatch should return a non-nil match")
	}
	
	if match.ID == "" {
		t.Error("Match should have a non-empty ID")
	}
	
	if match.Status != MatchStatusWaiting {
		t.Errorf("New match should have status %v, got %v", MatchStatusWaiting, match.Status)
	}
	
	if match.TeamLevels[0] != 2 || match.TeamLevels[1] != 2 {
		t.Error("Both teams should start at level 2")
	}
	
	if match.Winner != -1 {
		t.Error("New match should have no winner")
	}
	
	// Check players are assigned correctly
	for i, player := range players {
		if match.Players[i] == nil {
			t.Errorf("Player at seat %d should not be nil", i)
			continue
		}
		if match.Players[i].ID != player.ID {
			t.Errorf("Player at seat %d should have ID %s, got %s", i, player.ID, match.Players[i].ID)
		}
		if !match.Players[i].Online {
			t.Errorf("Player at seat %d should be online", i)
		}
		if match.Players[i].AutoPlay {
			t.Errorf("Player at seat %d should not be in auto-play", i)
		}
	}
}

func TestNewMatchValidation(t *testing.T) {
	// Test with wrong number of players
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
	}
	
	_, err := NewMatch(players)
	if err == nil {
		t.Error("NewMatch should fail with less than 4 players")
	}
	
	// Test with invalid seat numbers
	players = []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 5}, // Invalid seat
	}
	
	_, err = NewMatch(players)
	if err == nil {
		t.Error("NewMatch should fail with invalid seat number")
	}
	
	// Test with duplicate seats
	players = []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 2}, // Duplicate seat
	}
	
	_, err = NewMatch(players)
	if err == nil {
		t.Error("NewMatch should fail with duplicate seat numbers")
	}
}

func TestMatchTeamMethods(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	
	// Test GetTeamForPlayer
	if match.GetTeamForPlayer(0) != 0 {
		t.Error("Player 0 should be on team 0")
	}
	if match.GetTeamForPlayer(1) != 1 {
		t.Error("Player 1 should be on team 1")
	}
	if match.GetTeamForPlayer(2) != 0 {
		t.Error("Player 2 should be on team 0")
	}
	if match.GetTeamForPlayer(3) != 1 {
		t.Error("Player 3 should be on team 1")
	}
	
	// Test GetTeammates
	teammates := match.GetTeammates(0)
	if len(teammates) != 2 {
		t.Error("Player 0 should have 2 teammates (including self)")
	}
	if teammates[0] != 0 || teammates[1] != 2 {
		t.Error("Player 0's teammates should be 0 and 2")
	}
	
	// Test GetOpponents
	opponents := match.GetOpponents(0)
	if len(opponents) != 2 {
		t.Error("Player 0 should have 2 opponents")
	}
	if opponents[0] != 1 || opponents[1] != 3 {
		t.Error("Player 0's opponents should be 1 and 3")
	}
}

func TestMatchPlayerManagement(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	
	// Test HandlePlayerDisconnect
	err := match.HandlePlayerDisconnect(0)
	if err != nil {
		t.Errorf("HandlePlayerDisconnect failed: %v", err)
	}
	
	if match.IsPlayerOnline(0) {
		t.Error("Player 0 should be offline after disconnect")
	}
	
	if !match.IsPlayerAutoPlay(0) {
		t.Error("Player 0 should be in auto-play after disconnect")
	}
	
	// Test HandlePlayerReconnect
	err = match.HandlePlayerReconnect(0)
	if err != nil {
		t.Errorf("HandlePlayerReconnect failed: %v", err)
	}
	
	if !match.IsPlayerOnline(0) {
		t.Error("Player 0 should be online after reconnect")
	}
	
	if match.IsPlayerAutoPlay(0) {
		t.Error("Player 0 should not be in auto-play after reconnect")
	}
	
	// Test SetPlayerAutoPlay
	err = match.SetPlayerAutoPlay(1, true)
	if err != nil {
		t.Errorf("SetPlayerAutoPlay failed: %v", err)
	}
	
	if !match.IsPlayerAutoPlay(1) {
		t.Error("Player 1 should be in auto-play")
	}
}

func TestMatchPlayerManagementValidation(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	
	// Test with invalid seat numbers
	err := match.HandlePlayerDisconnect(-1)
	if err == nil {
		t.Error("HandlePlayerDisconnect should fail with invalid seat")
	}
	
	err = match.HandlePlayerDisconnect(4)
	if err == nil {
		t.Error("HandlePlayerDisconnect should fail with invalid seat")
	}
	
	err = match.HandlePlayerReconnect(-1)
	if err == nil {
		t.Error("HandlePlayerReconnect should fail with invalid seat")
	}
	
	err = match.SetPlayerAutoPlay(5, true)
	if err == nil {
		t.Error("SetPlayerAutoPlay should fail with invalid seat")
	}
}

func TestMatchStartNewDeal(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	
	// Test starting new deal
	err := match.StartNewDeal()
	if err != nil {
		t.Errorf("StartNewDeal failed: %v", err)
	}
	
	if match.CurrentDeal == nil {
		t.Error("Match should have a current deal after StartNewDeal")
	}
	
	if match.Status != MatchStatusPlaying {
		t.Errorf("Match status should be %v after StartNewDeal, got %v", MatchStatusPlaying, match.Status)
	}
	
	if match.CurrentDeal.Level != 2 {
		t.Error("First deal should be at level 2")
	}
}

func TestMatchFinishDeal(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	match.StartNewDeal()
	
	// Create a mock deal result
	result := &DealResult{
		Rankings:    []int{0, 2, 1, 3}, // Team 0 wins
		WinningTeam: 0,
		VictoryType: VictoryTypeNormal,
		Upgrades:    [2]int{1, 0}, // Team 0 gets +1 level
		Duration:    5 * time.Minute,
		TrickCount:  10,
	}
	
	// Test finishing deal
	err := match.FinishDeal(result)
	if err != nil {
		t.Errorf("FinishDeal failed: %v", err)
	}
	
	if len(match.DealHistory) != 1 {
		t.Error("Deal should be added to history")
	}
	
	if match.TeamLevels[0] != 3 {
		t.Error("Team 0 should be at level 3 after winning")
	}
	
	if match.TeamLevels[1] != 2 {
		t.Error("Team 1 should still be at level 2")
	}
	
	if match.Status != MatchStatusWaiting {
		t.Error("Match should be waiting for next deal")
	}
	
	if match.CurrentDeal != nil {
		t.Error("Current deal should be cleared after finishing")
	}
}

func TestMatchFinishWithWinner(t *testing.T) {
	players := []Player{
		{ID: "1", Username: "Player1", Seat: 0},
		{ID: "2", Username: "Player2", Seat: 1},
		{ID: "3", Username: "Player3", Seat: 2},
		{ID: "4", Username: "Player4", Seat: 3},
	}
	
	match, _ := NewMatch(players)
	match.StartNewDeal()
	
	// Set team 0 to level 13 (one away from A)
	match.TeamLevels[0] = 13
	
	// Create a winning deal result
	result := &DealResult{
		Rankings:    []int{0, 2, 1, 3}, // Team 0 wins
		WinningTeam: 0,
		VictoryType: VictoryTypeNormal,
		Upgrades:    [2]int{1, 0}, // Team 0 gets +1 level (reaches A)
	}
	
	err := match.FinishDeal(result)
	if err != nil {
		t.Errorf("FinishDeal failed: %v", err)
	}
	
	if match.Status != MatchStatusFinished {
		t.Error("Match should be finished when team reaches A")
	}
	
	if match.Winner != 0 {
		t.Error("Team 0 should be the winner")
	}
	
	if match.EndTime == nil {
		t.Error("Match should have an end time")
	}
}