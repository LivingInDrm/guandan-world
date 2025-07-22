package sdk

import (
	"testing"
	"time"
)

func TestDealResultCalculator_CalculateDealResult(t *testing.T) {
	// Create test players
	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}

	// Create test match
	match, err := NewMatch(players)
	if err != nil {
		t.Fatalf("Failed to create match: %v", err)
	}

	// Create test deal
	deal, err := NewDeal(5, nil)
	if err != nil {
		t.Fatalf("Failed to create deal: %v", err)
	}

	// Set up finished deal
	deal.Status = DealStatusFinished
	deal.Rankings = []int{0, 1, 2, 3} // Team 0 gets 1st(0) and 3rd(2), should be double down
	now := time.Now()
	deal.EndTime = &now
	deal.StartTime = now.Add(-30 * time.Minute)

	// Add some trick history for statistics
	trick1 := &Trick{
		ID:     "trick1",
		Winner: 0,
		Status: TrickStatusFinished,
		Plays: []*PlayAction{
			{PlayerSeat: 0, Cards: []*Card{{Number: 5, Color: "Hearts"}}, IsPass: false},
			{PlayerSeat: 1, IsPass: true},
			{PlayerSeat: 2, Cards: []*Card{{Number: 6, Color: "Hearts"}}, IsPass: false},
			{PlayerSeat: 3, IsPass: true},
		},
	}
	deal.TrickHistory = []*Trick{trick1}

	// Create calculator and calculate result
	calculator := NewDealResultCalculator(5)
	result, err := calculator.CalculateDealResult(deal, match)

	if err != nil {
		t.Fatalf("Failed to calculate deal result: %v", err)
	}

	// Verify basic result
	if result.WinningTeam != 0 {
		t.Errorf("Expected winning team 0, got %d", result.WinningTeam)
	}

	if result.VictoryType != VictoryTypeDoubleDown {
		t.Errorf("Expected double down victory, got %v", result.VictoryType)
	}

	if result.Upgrades[0] != 2 {
		t.Errorf("Expected 2 upgrades for winning team, got %d", result.Upgrades[0])
	}

	if result.Upgrades[1] != 0 {
		t.Errorf("Expected 0 upgrades for losing team, got %d", result.Upgrades[1])
	}

	if len(result.Rankings) != 4 {
		t.Errorf("Expected 4 rankings, got %d", len(result.Rankings))
	}

	// Verify statistics
	if result.Statistics == nil {
		t.Fatal("Expected statistics to be calculated")
	}

	if result.Statistics.TotalTricks != 1 {
		t.Errorf("Expected 1 total trick, got %d", result.Statistics.TotalTricks)
	}

	// Verify player stats
	if result.Statistics.PlayerStats[0].TricksWon != 1 {
		t.Errorf("Expected player 0 to win 1 trick, got %d", result.Statistics.PlayerStats[0].TricksWon)
	}

	if result.Statistics.PlayerStats[0].CardsPlayed != 1 {
		t.Errorf("Expected player 0 to play 1 card, got %d", result.Statistics.PlayerStats[0].CardsPlayed)
	}

	if result.Statistics.PlayerStats[1].PassCount != 1 {
		t.Errorf("Expected player 1 to pass 1 time, got %d", result.Statistics.PlayerStats[1].PassCount)
	}
}

func TestDealResultCalculator_VictoryTypes(t *testing.T) {
	// Create test players and match
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

	calculator := NewDealResultCalculator(5)

	testCases := []struct {
		name             string
		rankings         []int
		expectedType     VictoryType
		expectedUpgrades [2]int
	}{
		{
			name:             "Partner Last Victory - rank1,rank4 same team",
			rankings:         []int{0, 1, 3, 2}, // Team 0 gets 1st(0) and 4th(2) - Partner Last
			expectedType:     VictoryTypePartnerLast,
			expectedUpgrades: [2]int{1, 0},
		},
		{
			name:             "Partner Last Victory - alternating teams",
			rankings:         []int{0, 3, 1, 2}, // Team 0 gets 1st(0) and 4th(2) - Partner Last
			expectedType:     VictoryTypePartnerLast,
			expectedUpgrades: [2]int{1, 0},
		},
		{
			name:             "Single Last Victory - rank1,rank3 same team",
			rankings:         []int{0, 1, 2, 3}, // Team 0 gets 1st and 3rd
			expectedType:     VictoryTypeSingleLast,
			expectedUpgrades: [2]int{2, 0},
		},
		{
			name:             "Double Down Victory - rank1,rank2 same team",
			rankings:         []int{0, 2, 1, 3}, // Team 0 gets 1st and 2nd
			expectedType:     VictoryTypeDoubleDown,
			expectedUpgrades: [2]int{3, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			victoryType := calculator.calculateVictoryType(tc.rankings, match)
			if victoryType != tc.expectedType {
				t.Errorf("Expected victory type %v, got %v", tc.expectedType, victoryType)
			}

			winningTeam := match.GetTeamForPlayer(tc.rankings[0])
			upgrades := calculator.calculateUpgrades(victoryType, winningTeam)

			if upgrades[winningTeam] != tc.expectedUpgrades[winningTeam] {
				t.Errorf("Expected %d upgrades for winning team, got %d",
					tc.expectedUpgrades[winningTeam], upgrades[winningTeam])
			}
		})
	}
}

func TestDealResult_HelperMethods(t *testing.T) {
	// Create test players and match
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

	result := &DealResult{
		Rankings:    []int{0, 1, 2, 3}, // Team 0 gets 1st and 3rd (single last)
		WinningTeam: 0,
		VictoryType: VictoryTypeSingleLast,
		Upgrades:    [2]int{2, 0},
	}

	// Test GetTeamRankings
	teamRankings := result.GetTeamRankings(match)
	if len(teamRankings[0]) != 2 || len(teamRankings[1]) != 2 {
		t.Errorf("Expected 2 players per team in rankings")
	}

	// Team 0 should have ranks 1 and 3 (since rankings are [0, 1, 2, 3])
	expectedTeam0Ranks := []int{1, 3}
	for i, rank := range teamRankings[0] {
		if rank != expectedTeam0Ranks[i] {
			t.Errorf("Expected team 0 rank %d, got %d", expectedTeam0Ranks[i], rank)
		}
	}

	// Test GetWinningPlayers
	winningPlayers := result.GetWinningPlayers(match)
	expectedWinners := []int{0, 2}
	if len(winningPlayers) != len(expectedWinners) {
		t.Errorf("Expected %d winning players, got %d", len(expectedWinners), len(winningPlayers))
	}

	// Test GetLosingPlayers
	losingPlayers := result.GetLosingPlayers(match)
	expectedLosers := []int{1, 3}
	if len(losingPlayers) != len(expectedLosers) {
		t.Errorf("Expected %d losing players, got %d", len(expectedLosers), len(losingPlayers))
	}

	// Test victory type checks
	if !result.IsSingleLast() {
		t.Error("Expected IsSingleLast to return true")
	}

	if result.IsDoubleDown() {
		t.Error("Expected IsDoubleDown to return false")
	}

	// Test GetUpgradeForTeam
	if result.GetUpgradeForTeam(0) != 2 {
		t.Errorf("Expected 2 upgrades for team 0, got %d", result.GetUpgradeForTeam(0))
	}

	if result.GetUpgradeForTeam(1) != 0 {
		t.Errorf("Expected 0 upgrades for team 1, got %d", result.GetUpgradeForTeam(1))
	}
}

func TestDealResultCalculator_ErrorCases(t *testing.T) {
	calculator := NewDealResultCalculator(5)

	// Test nil deal
	_, err := calculator.CalculateDealResult(nil, nil)
	if err == nil {
		t.Error("Expected error for nil deal")
	}

	// Test unfinished deal
	deal, _ := NewDeal(5, nil)
	deal.Status = DealStatusPlaying

	players := []Player{
		{ID: "player1", Username: "Player1", Seat: 0},
		{ID: "player2", Username: "Player2", Seat: 1},
		{ID: "player3", Username: "Player3", Seat: 2},
		{ID: "player4", Username: "Player4", Seat: 3},
	}
	match, _ := NewMatch(players)

	_, err = calculator.CalculateDealResult(deal, match)
	if err == nil {
		t.Error("Expected error for unfinished deal")
	}

	// Test deal with no rankings
	deal.Status = DealStatusFinished
	deal.Rankings = []int{}

	_, err = calculator.CalculateDealResult(deal, match)
	if err == nil {
		t.Error("Expected error for deal with no rankings")
	}

	// Test nil match
	deal.Rankings = []int{0, 1, 2, 3}
	_, err = calculator.CalculateDealResult(deal, nil)
	if err == nil {
		t.Error("Expected error for nil match")
	}
}

func TestDealResultCalculator_TributeInfo(t *testing.T) {
	calculator := NewDealResultCalculator(5)

	// Test with no tribute phase
	tributeInfo := calculator.calculateTributeInfo(nil)
	if tributeInfo.HasTribute {
		t.Error("Expected HasTribute to be false for nil tribute phase")
	}

	// Test with tribute phase
	tributePhase := &TributePhase{
		Status:     TributeStatusFinished,
		TributeMap: map[int]int{0: 1, 2: 3},
		TributeCards: map[int]*Card{
			0: {Number: 14, Color: "Spades"},
			2: {Number: 13, Color: "Hearts"},
		},
		ReturnCards: map[int]*Card{
			1: {Number: 3, Color: "Clubs"},
			3: {Number: 4, Color: "Diamonds"},
		},
	}

	tributeInfo = calculator.calculateTributeInfo(tributePhase)
	if !tributeInfo.HasTribute {
		t.Error("Expected HasTribute to be true")
	}

	if len(tributeInfo.TributeMap) != 2 {
		t.Errorf("Expected 2 tribute mappings, got %d", len(tributeInfo.TributeMap))
	}

	if len(tributeInfo.TributeCards) != 2 {
		t.Errorf("Expected 2 tribute cards, got %d", len(tributeInfo.TributeCards))
	}

	if len(tributeInfo.ReturnCards) != 2 {
		t.Errorf("Expected 2 return cards, got %d", len(tributeInfo.ReturnCards))
	}
}

func TestDealResult_String(t *testing.T) {
	result := &DealResult{
		WinningTeam: 1,
		VictoryType: VictoryTypeDoubleDown,
		Upgrades:    [2]int{0, 2},
	}

	str := result.String()
	expected := "Team 1 wins with Double Down (+2 levels)"
	if str != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, str)
	}
}

func TestMatchResult_HelperMethods(t *testing.T) {
	stats := &MatchStatistics{
		TotalDeals:    3,
		TotalDuration: 2 * time.Hour,
		FinalLevels:   [2]int{14, 12},
		TeamStats: [2]*TeamMatchStats{
			{Team: 0, DealsWon: 2, TotalTricks: 45, Upgrades: 5},
			{Team: 1, DealsWon: 1, TotalTricks: 35, Upgrades: 3},
		},
	}

	result := &MatchResult{
		Winner:      0,
		FinalLevels: [2]int{14, 12},
		Duration:    2 * time.Hour,
		Statistics:  stats,
	}

	// Test GetWinningTeamStats
	winningStats := result.GetWinningTeamStats()
	if winningStats == nil {
		t.Fatal("Expected winning team stats")
	}

	if winningStats.Team != 0 {
		t.Errorf("Expected winning team 0, got %d", winningStats.Team)
	}

	if winningStats.DealsWon != 2 {
		t.Errorf("Expected 2 deals won, got %d", winningStats.DealsWon)
	}

	// Test GetLosingTeamStats
	losingStats := result.GetLosingTeamStats()
	if losingStats == nil {
		t.Fatal("Expected losing team stats")
	}

	if losingStats.Team != 1 {
		t.Errorf("Expected losing team 1, got %d", losingStats.Team)
	}

	if losingStats.DealsWon != 1 {
		t.Errorf("Expected 1 deal won, got %d", losingStats.DealsWon)
	}

	// Test String method
	str := result.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}
