package sdk

import (
	"testing"
	"time"
)

func TestMatch_NewMatch(t *testing.T) {
	// Test valid match creation
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Verify match properties
	if match.Status != MatchStatusWaiting {
		t.Errorf("Expected status waiting, got %v", match.Status)
	}
	
	if match.TeamLevels[0] != 2 || match.TeamLevels[1] != 2 {
		t.Errorf("Expected both teams to start at level 2, got %v", match.TeamLevels)
	}
	
	if match.Winner != -1 {
		t.Errorf("Expected no winner initially, got %d", match.Winner)
	}
	
	if len(match.DealHistory) != 0 {
		t.Errorf("Expected empty deal history, got %d deals", len(match.DealHistory))
	}
	
	// Verify players are assigned correctly
	for _, player := range players {
		matchPlayer := match.Players[player.Seat]
		if matchPlayer == nil {
			t.Errorf("Player not assigned to seat %d", player.Seat)
			continue
		}
		
		if matchPlayer.ID != player.ID {
			t.Errorf("Expected player ID %s, got %s", player.ID, matchPlayer.ID)
		}
		
		if !matchPlayer.Online {
			t.Errorf("Expected player %s to be online", player.ID)
		}
		
		if matchPlayer.AutoPlay {
			t.Errorf("Expected player %s to not be in auto-play", player.ID)
		}
	}
}

func TestMatch_NewMatch_ErrorCases(t *testing.T) {
	// Test insufficient players
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
	}
	
	_, err := NewMatch(players)
	if err == nil {
		t.Error("Expected error for insufficient players")
	}
	
	// Test duplicate seats
	players = []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 0},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	_, err = NewMatch(players)
	if err == nil {
		t.Error("Expected error for duplicate seats")
	}
	
	// Test invalid seat numbers
	players = []Player{
		{ID: "player1", Username: "Player1", Seat: -1},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	_, err = NewMatch(players)
	if err == nil {
		t.Error("Expected error for invalid seat number")
	}
}

func TestMatch_TeamAssignment(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test team assignments
	if match.GetTeamForPlayer(0) != 0 {
		t.Errorf("Expected player 0 to be on team 0")
	}
	if match.GetTeamForPlayer(1) != 1 {
		t.Errorf("Expected player 1 to be on team 1")
	}
	if match.GetTeamForPlayer(2) != 0 {
		t.Errorf("Expected player 2 to be on team 0")
	}
	if match.GetTeamForPlayer(3) != 1 {
		t.Errorf("Expected player 3 to be on team 1")
	}
	
	// Test teammates
	teammates0 := match.GetTeammates(0)
	if len(teammates0) != 2 || teammates0[0] != 0 || teammates0[1] != 2 {
		t.Errorf("Expected teammates [0, 2], got %v", teammates0)
	}
	
	teammates1 := match.GetTeammates(1)
	if len(teammates1) != 2 || teammates1[0] != 1 || teammates1[1] != 3 {
		t.Errorf("Expected teammates [1, 3], got %v", teammates1)
	}
	
	// Test opponents
	opponents0 := match.GetOpponents(0)
	if len(opponents0) != 2 || opponents0[0] != 1 || opponents0[1] != 3 {
		t.Errorf("Expected opponents [1, 3], got %v", opponents0)
	}
}

func TestMatch_PlayerManagement(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test disconnect
	err = match.HandlePlayerDisconnect(0)
	if err != nil {
		t.Errorf("Failed to handle disconnect: %v", err)
	}
	
	if match.IsPlayerOnline(0) {
		t.Error("Expected player 0 to be offline after disconnect")
	}
	
	if !match.IsPlayerAutoPlay(0) {
		t.Error("Expected player 0 to be in auto-play after disconnect")
	}
	
	// Test reconnect
	err = match.HandlePlayerReconnect(0)
	if err != nil {
		t.Errorf("Failed to handle reconnect: %v", err)
	}
	
	if !match.IsPlayerOnline(0) {
		t.Error("Expected player 0 to be online after reconnect")
	}
	
	if match.IsPlayerAutoPlay(0) {
		t.Error("Expected player 0 to not be in auto-play after reconnect")
	}
	
	// Test auto-play setting
	err = match.SetPlayerAutoPlay(1, true)
	if err != nil {
		t.Errorf("Failed to set auto-play: %v", err)
	}
	
	if !match.IsPlayerAutoPlay(1) {
		t.Error("Expected player 1 to be in auto-play")
	}
}

func TestMatch_LevelManagement(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test initial levels
	if match.GetTeamLevel(0) != 2 {
		t.Errorf("Expected team 0 level 2, got %d", match.GetTeamLevel(0))
	}
	
	if match.GetPlayerTeamLevel(0) != 2 {
		t.Errorf("Expected player 0 team level 2, got %d", match.GetPlayerTeamLevel(0))
	}
	
	if match.GetCurrentLevel() != 2 {
		t.Errorf("Expected current level 2, got %d", match.GetCurrentLevel())
	}
	
	// Test level updates
	result := &DealResult{
		WinningTeam: 0,
		Upgrades:    [2]int{2, 0}, // Team 0 gets +2
	}
	
	match.updateTeamLevels(result)
	
	if match.GetTeamLevel(0) != 4 {
		t.Errorf("Expected team 0 level 4, got %d", match.GetTeamLevel(0))
	}
	
	if match.GetTeamLevel(1) != 2 {
		t.Errorf("Expected team 1 level 2, got %d", match.GetTeamLevel(1))
	}
	
	// Test level queries
	if !match.IsTeamAtLevel(0, 4) {
		t.Error("Expected team 0 to be at level 4")
	}
	
	if match.IsTeamAtLevel(1, 4) {
		t.Error("Expected team 1 to not be at level 4")
	}
	
	if match.GetLeadingTeam() != 0 {
		t.Errorf("Expected team 0 to be leading, got %d", match.GetLeadingTeam())
	}
	
	if match.GetLevelDifference() != 2 {
		t.Errorf("Expected level difference 2, got %d", match.GetLevelDifference())
	}
}

func TestMatch_FinishConditions(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Initially not finished
	if match.IsAnyTeamAtALevel() {
		t.Error("Expected match to not be finished initially")
	}
	
	if match.CanStartNewDeal() != true {
		t.Error("Expected to be able to start new deal initially")
	}
	
	// Simulate team reaching A level
	match.TeamLevels[0] = 14
	
	if !match.IsAnyTeamAtALevel() {
		t.Error("Expected match to be finished when team reaches A level")
	}
	
	if match.CanStartNewDeal() {
		t.Error("Expected to not be able to start new deal when finished")
	}
	
	// Test winner determination
	winner := match.getWinningTeam()
	if winner != 0 {
		t.Errorf("Expected team 0 to win, got %d", winner)
	}
	
	// Test both teams at A level
	match.TeamLevels[1] = 14
	winner = match.getWinningTeam()
	if winner != 0 {
		t.Errorf("Expected team 0 to win (higher level), got %d", winner)
	}
	
	// Test same level - should use last deal result
	match.TeamLevels[0] = 14
	match.TeamLevels[1] = 14
	
	// Create a mock deal with result
	deal, _ := NewDeal(14, nil)
	deal.Status = DealStatusFinished
	deal.Rankings = []int{1, 3, 0, 2} // Team 1 wins
	now := time.Now()
	deal.EndTime = &now
	match.DealHistory = []*Deal{deal}
	
	winner = match.getWinningTeam()
	if winner != 1 {
		t.Errorf("Expected team 1 to win (last deal winner), got %d", winner)
	}
}

func TestMatch_DealManagement(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test initial state
	if match.GetDealCount() != 0 {
		t.Errorf("Expected 0 deals initially, got %d", match.GetDealCount())
	}
	
	if match.HasActiveDeal() {
		t.Error("Expected no active deal initially")
	}
	
	if match.GetActiveDeal() != nil {
		t.Error("Expected nil active deal initially")
	}
	
	if match.GetLastDeal() != nil {
		t.Error("Expected nil last deal initially")
	}
	
	// Create and add a deal to history
	deal, _ := NewDeal(2, nil)
	deal.Status = DealStatusFinished
	now := time.Now()
	deal.EndTime = &now
	match.DealHistory = []*Deal{deal}
	
	if match.GetDealCount() != 1 {
		t.Errorf("Expected 1 deal, got %d", match.GetDealCount())
	}
	
	if match.GetLastDeal() != deal {
		t.Error("Expected last deal to match added deal")
	}
	
	if match.GetDealByIndex(0) != deal {
		t.Error("Expected deal at index 0 to match added deal")
	}
	
	if match.GetDealByIndex(1) != nil {
		t.Error("Expected nil for invalid deal index")
	}
}

func TestMatch_PlayerQueries(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test player counts
	if match.GetPlayerCount() != 4 {
		t.Errorf("Expected 4 players, got %d", match.GetPlayerCount())
	}
	
	if match.GetOnlinePlayerCount() != 4 {
		t.Errorf("Expected 4 online players, got %d", match.GetOnlinePlayerCount())
	}
	
	if match.GetAutoPlayPlayerCount() != 0 {
		t.Errorf("Expected 0 auto-play players, got %d", match.GetAutoPlayPlayerCount())
	}
	
	// Disconnect a player
	match.HandlePlayerDisconnect(0)
	
	if match.GetOnlinePlayerCount() != 3 {
		t.Errorf("Expected 3 online players after disconnect, got %d", match.GetOnlinePlayerCount())
	}
	
	if match.GetAutoPlayPlayerCount() != 1 {
		t.Errorf("Expected 1 auto-play player after disconnect, got %d", match.GetAutoPlayPlayerCount())
	}
	
	// Test player lookups
	player := match.GetPlayerByID("player1")
	if player == nil || player.ID != "player1" {
		t.Error("Failed to find player by ID")
	}
	
	player = match.GetPlayerBySeat(0)
	if player == nil || player.Seat != 0 {
		t.Error("Failed to find player by seat")
	}
	
	player = match.GetPlayerByID("nonexistent")
	if player != nil {
		t.Error("Expected nil for nonexistent player ID")
	}
	
	player = match.GetPlayerBySeat(5)
	if player != nil {
		t.Error("Expected nil for invalid seat")
	}
}

func TestMatch_Statistics(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test initial statistics
	stats := match.GetMatchStatistics()
	if stats.TotalDeals != 0 {
		t.Errorf("Expected 0 total deals, got %d", stats.TotalDeals)
	}
	
	if stats.TeamStats[0].DealsWon != 0 {
		t.Errorf("Expected 0 deals won for team 0, got %d", stats.TeamStats[0].DealsWon)
	}
	
	// Test match result for unfinished match
	result := match.GetMatchResult()
	if result != nil {
		t.Error("Expected nil match result for unfinished match")
	}
	
	// Finish the match
	match.Status = MatchStatusFinished
	match.Winner = 0
	now := time.Now()
	match.EndTime = &now
	
	result = match.GetMatchResult()
	if result == nil {
		t.Fatal("Expected match result for finished match")
	}
	
	if result.Winner != 0 {
		t.Errorf("Expected winner 0, got %d", result.Winner)
	}
}

func TestMatch_String(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Test string representation for active match
	str := match.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	// Test string representation for finished match
	match.Status = MatchStatusFinished
	match.Winner = 1
	match.TeamLevels = [2]int{12, 14}
	
	str = match.String()
	if str == "" {
		t.Error("Expected non-empty string representation for finished match")
	}
}

func TestMatch_FinishDeal(t *testing.T) {
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}
	
	// Create a deal
	deal, _ := NewDeal(2, nil)
	match.CurrentDeal = deal
	
	// Create deal result
	result := &DealResult{
		WinningTeam: 0,
		Upgrades:    [2]int{1, 0},
	}
	
	// Finish the deal
	err = match.FinishDeal(result)
	if err != nil {
		t.Errorf("Failed to finish deal: %v", err)
	}
	
	// Verify deal was added to history
	if len(match.DealHistory) != 1 {
		t.Errorf("Expected 1 deal in history, got %d", len(match.DealHistory))
	}
	
	// Verify levels were updated
	if match.TeamLevels[0] != 3 {
		t.Errorf("Expected team 0 level 3, got %d", match.TeamLevels[0])
	}
	
	// Verify current deal was cleared
	if match.CurrentDeal != nil {
		t.Error("Expected current deal to be cleared")
	}
	
	// Verify status is waiting (not finished)
	if match.Status != MatchStatusWaiting {
		t.Errorf("Expected status waiting, got %v", match.Status)
	}
	
	// Test finishing with A level
	match.CurrentDeal = deal
	result = &DealResult{
		WinningTeam: 0,
		Upgrades:    [2]int{11, 0}, // This should bring team 0 to level 14
	}
	
	err = match.FinishDeal(result)
	if err != nil {
		t.Errorf("Failed to finish deal with A level: %v", err)
	}
	
	// Verify match is finished
	if match.Status != MatchStatusFinished {
		t.Errorf("Expected status finished, got %v", match.Status)
	}
	
	if match.Winner != 0 {
		t.Errorf("Expected winner 0, got %d", match.Winner)
	}
	
	if match.EndTime == nil {
		t.Error("Expected end time to be set")
	}
}