package sdk

import (
	"testing"
	"time"

	pbgame "guandan-world/proto/gen/go/game"
)

func TestTeamUpgradesAdapter(t *testing.T) {
	upgrades := [2]int{3, 0}

	proto := ToProtoTeamUpgrades(upgrades)
	if proto.Team0 != 3 || proto.Team1 != 0 {
		t.Errorf("ToProtoTeamUpgrades failed: got (%d, %d), want (3, 0)", proto.Team0, proto.Team1)
	}

	back := FromProtoTeamUpgrades(proto)
	if back != upgrades {
		t.Errorf("FromProtoTeamUpgrades failed: got %v, want %v", back, upgrades)
	}
}

func TestPlayerDealStatsAdapter(t *testing.T) {
	stats := &PlayerDealStats{
		PlayerSeat:   0,
		CardsPlayed:  27,
		TricksWon:    5,
		PassCount:    3,
		TimeoutCount: 0,
		FinishRank:   1,
	}

	proto := ToProtoPlayerDealStats(stats)
	if proto.PlayerSeat != 0 {
		t.Errorf("PlayerSeat: got %d, want 0", proto.PlayerSeat)
	}
	if proto.CardsPlayed != 27 {
		t.Errorf("CardsPlayed: got %d, want 27", proto.CardsPlayed)
	}
	if proto.TricksWon != 5 {
		t.Errorf("TricksWon: got %d, want 5", proto.TricksWon)
	}
	if proto.FinishRank != 1 {
		t.Errorf("FinishRank: got %d, want 1", proto.FinishRank)
	}

	back := FromProtoPlayerDealStats(proto)
	if back.PlayerSeat != stats.PlayerSeat || back.CardsPlayed != stats.CardsPlayed {
		t.Errorf("FromProtoPlayerDealStats failed")
	}
}

func TestPlayerDealStatsAdapterNil(t *testing.T) {
	var stats *PlayerDealStats
	proto := ToProtoPlayerDealStats(stats)
	if proto != nil {
		t.Errorf("ToProtoPlayerDealStats(nil) should return nil")
	}

	back := FromProtoPlayerDealStats(nil)
	if back != nil {
		t.Errorf("FromProtoPlayerDealStats(nil) should return nil")
	}
}

func TestDealStatisticsAdapter(t *testing.T) {
	stats := &DealStatistics{
		TotalTricks: 10,
		PlayerStats: [4]*PlayerDealStats{
			{PlayerSeat: 0, CardsPlayed: 27, TricksWon: 5, FinishRank: 1},
			{PlayerSeat: 1, CardsPlayed: 27, TricksWon: 3, FinishRank: 3},
			{PlayerSeat: 2, CardsPlayed: 27, TricksWon: 2, FinishRank: 2},
			{PlayerSeat: 3, CardsPlayed: 27, TricksWon: 0, FinishRank: 4},
		},
		TributeInfo: &TributeInfo{
			HasTribute: false,
		},
	}

	proto := ToProtoDealStatistics(stats)
	if proto.TotalTricks != 10 {
		t.Errorf("TotalTricks: got %d, want 10", proto.TotalTricks)
	}
	if len(proto.PlayerStats) != 4 {
		t.Errorf("PlayerStats length: got %d, want 4", len(proto.PlayerStats))
	}
	if proto.PlayerStats[0].TricksWon != 5 {
		t.Errorf("PlayerStats[0].TricksWon: got %d, want 5", proto.PlayerStats[0].TricksWon)
	}

	back := FromProtoDealStatistics(proto)
	if back.TotalTricks != 10 {
		t.Errorf("Back TotalTricks: got %d, want 10", back.TotalTricks)
	}
	if len(back.PlayerStats) != 4 {
		t.Errorf("Back PlayerStats length: got %d, want 4", len(back.PlayerStats))
	}
	if back.PlayerStats[0].TricksWon != 5 {
		t.Errorf("Back PlayerStats[0].TricksWon: got %d, want 5", back.PlayerStats[0].TricksWon)
	}
}

func TestDealResultAdapter(t *testing.T) {
	result := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		WinningTeam: 0,
		VictoryType: VictoryTypeDoubleDown,
		Upgrades:    [2]int{3, 0},
		Duration:    5 * time.Minute,
		TrickCount:  12,
		Statistics: &DealStatistics{
			TotalTricks: 12,
			PlayerStats: [4]*PlayerDealStats{
				{PlayerSeat: 0, FinishRank: 1},
				{PlayerSeat: 1, FinishRank: 3},
				{PlayerSeat: 2, FinishRank: 2},
				{PlayerSeat: 3, FinishRank: 4},
			},
			TributeInfo: &TributeInfo{HasTribute: false},
		},
	}

	proto := ToProtoDealResult(result)
	if len(proto.Rankings) != 4 {
		t.Errorf("Rankings length: got %d, want 4", len(proto.Rankings))
	}
	if proto.Rankings[0] != 0 {
		t.Errorf("Rankings[0]: got %d, want 0", proto.Rankings[0])
	}
	if proto.WinningTeam != 0 {
		t.Errorf("WinningTeam: got %d, want 0", proto.WinningTeam)
	}
	if proto.TrickCount != 12 {
		t.Errorf("TrickCount: got %d, want 12", proto.TrickCount)
	}
	if proto.DurationMs != 5*60*1000 {
		t.Errorf("DurationMs: got %d, want %d", proto.DurationMs, 5*60*1000)
	}

	back := FromProtoDealResult(proto)
	if len(back.Rankings) != 4 {
		t.Errorf("Back Rankings length: got %d, want 4", len(back.Rankings))
	}
	if back.Rankings[0] != 0 {
		t.Errorf("Back Rankings[0]: got %d, want 0", back.Rankings[0])
	}
	if back.WinningTeam != 0 {
		t.Errorf("Back WinningTeam: got %d, want 0", back.WinningTeam)
	}
	if back.Duration != 5*time.Minute {
		t.Errorf("Back Duration: got %v, want %v", back.Duration, 5*time.Minute)
	}
}

func TestTeamMatchStatsAdapter(t *testing.T) {
	stats := &TeamMatchStats{
		Team:        0,
		DealsWon:    3,
		TotalTricks: 45,
		Upgrades:    8,
	}

	proto := ToProtoTeamMatchStats(stats)
	if proto.Team != 0 {
		t.Errorf("Team: got %d, want 0", proto.Team)
	}
	if proto.DealsWon != 3 {
		t.Errorf("DealsWon: got %d, want 3", proto.DealsWon)
	}
	if proto.TotalTricks != 45 {
		t.Errorf("TotalTricks: got %d, want 45", proto.TotalTricks)
	}
	if proto.Upgrades != 8 {
		t.Errorf("Upgrades: got %d, want 8", proto.Upgrades)
	}

	back := FromProtoTeamMatchStats(proto)
	if back.Team != stats.Team || back.DealsWon != stats.DealsWon {
		t.Errorf("FromProtoTeamMatchStats failed")
	}
}

func TestMatchStatisticsAdapter(t *testing.T) {
	stats := &MatchStatistics{
		TotalDeals:    5,
		TotalDuration: 30 * time.Minute,
		FinalLevels:   [2]int{14, 7},
		TeamStats: [2]*TeamMatchStats{
			{Team: 0, DealsWon: 3, TotalTricks: 45, Upgrades: 8},
			{Team: 1, DealsWon: 2, TotalTricks: 30, Upgrades: 5},
		},
	}

	proto := ToProtoMatchStatistics(stats)
	if proto.TotalDeals != 5 {
		t.Errorf("TotalDeals: got %d, want 5", proto.TotalDeals)
	}
	if proto.TotalDurationMs != 30*60*1000 {
		t.Errorf("TotalDurationMs: got %d, want %d", proto.TotalDurationMs, 30*60*1000)
	}
	if len(proto.TeamStats) != 2 {
		t.Errorf("TeamStats length: got %d, want 2", len(proto.TeamStats))
	}
	if proto.TeamStats[0].DealsWon != 3 {
		t.Errorf("TeamStats[0].DealsWon: got %d, want 3", proto.TeamStats[0].DealsWon)
	}

	back := FromProtoMatchStatistics(proto, [2]int{14, 7})
	if back.TotalDeals != 5 {
		t.Errorf("Back TotalDeals: got %d, want 5", back.TotalDeals)
	}
	if back.TotalDuration != 30*time.Minute {
		t.Errorf("Back TotalDuration: got %v, want %v", back.TotalDuration, 30*time.Minute)
	}
	if back.FinalLevels != [2]int{14, 7} {
		t.Errorf("Back FinalLevels: got %v, want [14 7]", back.FinalLevels)
	}
	if len(back.TeamStats) != 2 {
		t.Errorf("Back TeamStats length: got %d, want 2", len(back.TeamStats))
	}
}

func TestMatchResultAdapter(t *testing.T) {
	result := &MatchResult{
		Winner:      0,
		FinalLevels: [2]int{14, 7},
		Duration:    45 * time.Minute,
		Statistics: &MatchStatistics{
			TotalDeals:    5,
			TotalDuration: 45 * time.Minute,
			FinalLevels:   [2]int{14, 7},
			TeamStats: [2]*TeamMatchStats{
				{Team: 0, DealsWon: 3},
				{Team: 1, DealsWon: 2},
			},
		},
	}

	proto := ToProtoMatchResult(result)
	if proto.Winner != 0 {
		t.Errorf("Winner: got %d, want 0", proto.Winner)
	}
	if proto.DurationMs != 45*60*1000 {
		t.Errorf("DurationMs: got %d, want %d", proto.DurationMs, 45*60*1000)
	}
	if proto.FinalLevels.Team0 != 14 || proto.FinalLevels.Team1 != 7 {
		t.Errorf("FinalLevels: got (%d, %d), want (14, 7)", proto.FinalLevels.Team0, proto.FinalLevels.Team1)
	}

	back := FromProtoMatchResult(proto)
	if back.Winner != 0 {
		t.Errorf("Back Winner: got %d, want 0", back.Winner)
	}
	if back.Duration != 45*time.Minute {
		t.Errorf("Back Duration: got %v, want %v", back.Duration, 45*time.Minute)
	}
	if back.FinalLevels[0] != 14 || back.FinalLevels[1] != 7 {
		t.Errorf("Back FinalLevels: got %v, want [14, 7]", back.FinalLevels)
	}
}

func TestResultNilHandling(t *testing.T) {
	if proto := ToProtoTeamUpgrades([2]int{}); proto == nil {
		t.Error("ToProtoTeamUpgrades should never return nil")
	}

	if back := FromProtoTeamUpgrades(nil); back != [2]int{0, 0} {
		t.Errorf("FromProtoTeamUpgrades(nil) should return [0, 0], got %v", back)
	}

	if proto := ToProtoDealResult(nil); proto != nil {
		t.Error("ToProtoDealResult(nil) should return nil")
	}

	if back := FromProtoDealResult(nil); back != nil {
		t.Error("FromProtoDealResult(nil) should return nil")
	}

	if proto := ToProtoMatchResult(nil); proto != nil {
		t.Error("ToProtoMatchResult(nil) should return nil")
	}

	if back := FromProtoMatchResult(nil); back != nil {
		t.Error("FromProtoMatchResult(nil) should return nil")
	}
}

func TestDurationConversion(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		millis   int64
	}{
		{5 * time.Minute, 5 * 60 * 1000},
		{30 * time.Second, 30 * 1000},
		{1 * time.Hour, 60 * 60 * 1000},
		{0, 0},
	}

	for _, tc := range testCases {
		result := &DealResult{
			Duration: tc.duration,
		}
		proto := ToProtoDealResult(result)
		if proto.DurationMs != tc.millis {
			t.Errorf("Duration %v: got %d ms, want %d ms", tc.duration, proto.DurationMs, tc.millis)
		}

		back := FromProtoDealResult(proto)
		if back.Duration != tc.duration {
			t.Errorf("Duration back conversion: got %v, want %v", back.Duration, tc.duration)
		}
	}
}

func TestRoundTripDealResult(t *testing.T) {
	original := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		WinningTeam: 0,
		VictoryType: VictoryTypeSingleLast,
		Upgrades:    [2]int{2, 0},
		Duration:    10 * time.Minute,
		TrickCount:  15,
		Statistics: &DealStatistics{
			TotalTricks: 15,
			PlayerStats: [4]*PlayerDealStats{
				{PlayerSeat: 0, CardsPlayed: 27, TricksWon: 8, PassCount: 2, TimeoutCount: 0, FinishRank: 1},
				{PlayerSeat: 1, CardsPlayed: 27, TricksWon: 3, PassCount: 5, TimeoutCount: 0, FinishRank: 3},
				{PlayerSeat: 2, CardsPlayed: 27, TricksWon: 4, PassCount: 4, TimeoutCount: 0, FinishRank: 2},
				{PlayerSeat: 3, CardsPlayed: 27, TricksWon: 0, PassCount: 8, TimeoutCount: 1, FinishRank: 4},
			},
			TributeInfo: &TributeInfo{
				HasTribute: true,
				TributeMap: map[int]int{0: 2},
			},
		},
	}

	proto := ToProtoDealResult(original)
	back := FromProtoDealResult(proto)

	if len(back.Rankings) != len(original.Rankings) {
		t.Fatalf("Rankings length mismatch: got %d, want %d", len(back.Rankings), len(original.Rankings))
	}
	for i := range original.Rankings {
		if back.Rankings[i] != original.Rankings[i] {
			t.Errorf("Rankings[%d]: got %d, want %d", i, back.Rankings[i], original.Rankings[i])
		}
	}

	if back.WinningTeam != original.WinningTeam {
		t.Errorf("WinningTeam: got %d, want %d", back.WinningTeam, original.WinningTeam)
	}

	if back.VictoryType != original.VictoryType {
		t.Errorf("VictoryType: got %s, want %s", back.VictoryType, original.VictoryType)
	}

	if back.Upgrades != original.Upgrades {
		t.Errorf("Upgrades: got %v, want %v", back.Upgrades, original.Upgrades)
	}

	if back.Duration != original.Duration {
		t.Errorf("Duration: got %v, want %v", back.Duration, original.Duration)
	}

	if back.TrickCount != original.TrickCount {
		t.Errorf("TrickCount: got %d, want %d", back.TrickCount, original.TrickCount)
	}

	if back.Statistics.TotalTricks != original.Statistics.TotalTricks {
		t.Errorf("Statistics.TotalTricks: got %d, want %d", back.Statistics.TotalTricks, original.Statistics.TotalTricks)
	}

	for i := 0; i < 4; i++ {
		if back.Statistics.PlayerStats[i].TricksWon != original.Statistics.PlayerStats[i].TricksWon {
			t.Errorf("PlayerStats[%d].TricksWon: got %d, want %d", i,
				back.Statistics.PlayerStats[i].TricksWon, original.Statistics.PlayerStats[i].TricksWon)
		}
	}
}

func TestRoundTripMatchResult(t *testing.T) {
	original := &MatchResult{
		Winner:      1,
		FinalLevels: [2]int{9, 14},
		Duration:    60 * time.Minute,
		Statistics: &MatchStatistics{
			TotalDeals:    8,
			TotalDuration: 60 * time.Minute,
			FinalLevels:   [2]int{9, 14},
			TeamStats: [2]*TeamMatchStats{
				{Team: 0, DealsWon: 3, TotalTricks: 40, Upgrades: 7},
				{Team: 1, DealsWon: 5, TotalTricks: 55, Upgrades: 12},
			},
		},
	}

	proto := ToProtoMatchResult(original)
	back := FromProtoMatchResult(proto)

	if back.Winner != original.Winner {
		t.Errorf("Winner: got %d, want %d", back.Winner, original.Winner)
	}

	if back.FinalLevels != original.FinalLevels {
		t.Errorf("FinalLevels: got %v, want %v", back.FinalLevels, original.FinalLevels)
	}

	if back.Duration != original.Duration {
		t.Errorf("Duration: got %v, want %v", back.Duration, original.Duration)
	}

	if back.Statistics.TotalDeals != original.Statistics.TotalDeals {
		t.Errorf("Statistics.TotalDeals: got %d, want %d", back.Statistics.TotalDeals, original.Statistics.TotalDeals)
	}

	for i := 0; i < 2; i++ {
		if back.Statistics.TeamStats[i].DealsWon != original.Statistics.TeamStats[i].DealsWon {
			t.Errorf("TeamStats[%d].DealsWon: got %d, want %d", i,
				back.Statistics.TeamStats[i].DealsWon, original.Statistics.TeamStats[i].DealsWon)
		}
	}
}

func BenchmarkDealResultAdapter(b *testing.B) {
	result := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		WinningTeam: 0,
		VictoryType: VictoryTypeDoubleDown,
		Upgrades:    [2]int{3, 0},
		Duration:    5 * time.Minute,
		TrickCount:  12,
		Statistics: &DealStatistics{
			TotalTricks: 12,
			PlayerStats: [4]*PlayerDealStats{
				{PlayerSeat: 0, CardsPlayed: 27, TricksWon: 8, FinishRank: 1},
				{PlayerSeat: 1, CardsPlayed: 27, TricksWon: 3, FinishRank: 3},
				{PlayerSeat: 2, CardsPlayed: 27, TricksWon: 1, FinishRank: 2},
				{PlayerSeat: 3, CardsPlayed: 27, TricksWon: 0, FinishRank: 4},
			},
			TributeInfo: &TributeInfo{HasTribute: false},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proto := ToProtoDealResult(result)
		_ = FromProtoDealResult(proto)
	}
}

func BenchmarkMatchResultAdapter(b *testing.B) {
	result := &MatchResult{
		Winner:      0,
		FinalLevels: [2]int{14, 7},
		Duration:    45 * time.Minute,
		Statistics: &MatchStatistics{
			TotalDeals:    5,
			TotalDuration: 45 * time.Minute,
			FinalLevels:   [2]int{14, 7},
			TeamStats: [2]*TeamMatchStats{
				{Team: 0, DealsWon: 3, TotalTricks: 45, Upgrades: 8},
				{Team: 1, DealsWon: 2, TotalTricks: 30, Upgrades: 5},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proto := ToProtoMatchResult(result)
		_ = FromProtoMatchResult(proto)
	}
}

func ExampleToProtoDealResult() {
	result := &DealResult{
		Rankings:    []int{0, 2, 1, 3},
		WinningTeam: 0,
		VictoryType: VictoryTypeDoubleDown,
		Upgrades:    [2]int{3, 0},
		Duration:    5 * time.Minute,
		TrickCount:  12,
	}

	proto := ToProtoDealResult(result)
	_ = proto
}

func ExampleFromProtoMatchResult() {
	proto := &pbgame.MatchResult{
		Winner: 0,
		FinalLevels: &pbgame.TeamUpgrades{
			Team0: 14,
			Team1: 7,
		},
		DurationMs: 45 * 60 * 1000,
	}

	result := FromProtoMatchResult(proto)
	_ = result
}
