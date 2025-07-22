package sdk

import (
	"fmt"
	"time"
)

// DealResultCalculator handles the calculation of deal results and statistics
type DealResultCalculator struct {
	level int
}

// NewDealResultCalculator creates a new deal result calculator
func NewDealResultCalculator(level int) *DealResultCalculator {
	return &DealResultCalculator{
		level: level,
	}
}

// CalculateDealResult calculates the complete result of a finished deal
func (drc *DealResultCalculator) CalculateDealResult(deal *Deal, match *Match) (*DealResult, error) {
	if deal == nil {
		return nil, fmt.Errorf("deal cannot be nil")
	}

	if deal.Status != DealStatusFinished {
		return nil, fmt.Errorf("deal is not finished: %s", deal.Status)
	}

	if len(deal.Rankings) == 0 {
		return nil, fmt.Errorf("no rankings available")
	}

	if match == nil {
		return nil, fmt.Errorf("match cannot be nil")
	}

	// Calculate basic result information
	winningTeam := match.GetTeamForPlayer(deal.Rankings[0])
	victoryType := drc.calculateVictoryType(deal.Rankings, match)
	upgrades := drc.calculateUpgrades(victoryType, winningTeam)

	// Calculate duration
	duration := time.Duration(0)
	if deal.EndTime != nil {
		duration = deal.EndTime.Sub(deal.StartTime)
	}

	// Calculate statistics
	stats := drc.calculateDealStatistics(deal)

	result := &DealResult{
		Rankings:    deal.Rankings,
		WinningTeam: winningTeam,
		VictoryType: victoryType,
		Upgrades:    upgrades,
		Duration:    duration,
		TrickCount:  len(deal.TrickHistory),
		Statistics:  stats,
	}

	return result, nil
}

// calculateVictoryType determines the type of victory based on rankings
func (drc *DealResultCalculator) calculateVictoryType(rankings []int, match *Match) VictoryType {
	if len(rankings) < 4 {
		return VictoryTypePartnerLast
	}

	// Get teams for each ranking position
	teams := make([]int, 4)
	for i, playerSeat := range rankings {
		teams[i] = match.GetTeamForPlayer(playerSeat)
	}

	// 掼蛋升级规则：
	// 1. rank1, rank2 同队 → Double Down (+3级)
	// 2. rank1, rank3 同队 → Single Last (+2级)
	// 3. rank1, rank4 同队 → Partner Last (+1级)

	// Check for rank1, rank2 same team (+3级)
	if teams[0] == teams[1] {
		return VictoryTypeDoubleDown
	}

	// Check for rank1, rank3 same team (+2级)
	if teams[0] == teams[2] {
		return VictoryTypeSingleLast
	}

	// rank1, rank4 same team or other cases (+1级)
	return VictoryTypePartnerLast
}

// calculateUpgrades calculates level upgrades for each team based on victory type
func (drc *DealResultCalculator) calculateUpgrades(victoryType VictoryType, winningTeam int) [2]int {
	upgrades := [2]int{0, 0}

	if winningTeam < 0 || winningTeam > 1 {
		return upgrades
	}

	switch victoryType {
	case VictoryTypePartnerLast:
		// rank1, rank4 同队或其他情况：升1级
		upgrades[winningTeam] = 1
	case VictoryTypeSingleLast:
		// rank1, rank3 同队：升2级
		upgrades[winningTeam] = 2
	case VictoryTypeDoubleDown:
		// rank1, rank2 同队：升3级
		upgrades[winningTeam] = 3
	}

	return upgrades
}

// calculateDealStatistics calculates detailed statistics for the deal
func (drc *DealResultCalculator) calculateDealStatistics(deal *Deal) *DealStatistics {
	stats := &DealStatistics{
		TotalTricks: len(deal.TrickHistory),
		PlayerStats: [4]*PlayerDealStats{},
		TributeInfo: drc.calculateTributeInfo(deal.TributePhase),
	}

	// Initialize player stats
	for i := 0; i < 4; i++ {
		stats.PlayerStats[i] = &PlayerDealStats{
			PlayerSeat:   i,
			CardsPlayed:  0,
			TricksWon:    0,
			PassCount:    0,
			TimeoutCount: 0,
			FinishRank:   drc.getPlayerFinishRank(i, deal.Rankings),
		}
	}

	// Calculate statistics from trick history
	for _, trick := range deal.TrickHistory {
		// Count tricks won
		if trick.Winner >= 0 && trick.Winner < 4 {
			stats.PlayerStats[trick.Winner].TricksWon++
		}

		// Count plays and passes for each player
		for _, play := range trick.Plays {
			if play.PlayerSeat >= 0 && play.PlayerSeat < 4 {
				if play.IsPass {
					stats.PlayerStats[play.PlayerSeat].PassCount++
				} else {
					stats.PlayerStats[play.PlayerSeat].CardsPlayed += len(play.Cards)
				}
			}
		}
	}

	return stats
}

// calculateTributeInfo extracts tribute information from tribute phase
func (drc *DealResultCalculator) calculateTributeInfo(tributePhase *TributePhase) *TributeInfo {
	if tributePhase == nil {
		return &TributeInfo{
			HasTribute: false,
		}
	}

	info := &TributeInfo{
		HasTribute:   true,
		TributeMap:   make(map[int]int),
		TributeCards: make(map[int]*Card),
		ReturnCards:  make(map[int]*Card),
	}

	// Copy tribute information
	for giver, receiver := range tributePhase.TributeMap {
		info.TributeMap[giver] = receiver
	}

	for giver, card := range tributePhase.TributeCards {
		info.TributeCards[giver] = card
	}

	for receiver, card := range tributePhase.ReturnCards {
		info.ReturnCards[receiver] = card
	}

	return info
}

// getPlayerFinishRank returns the finish rank for a player (1-4, or 0 if not finished)
func (drc *DealResultCalculator) getPlayerFinishRank(playerSeat int, rankings []int) int {
	for rank, seat := range rankings {
		if seat == playerSeat {
			return rank + 1 // Convert 0-based to 1-based ranking
		}
	}
	return 0 // Player didn't finish (shouldn't happen in completed deal)
}

// DealStatistics contains detailed statistics for a completed deal
type DealStatistics struct {
	TotalTricks int                 `json:"total_tricks"`
	PlayerStats [4]*PlayerDealStats `json:"player_stats"`
	TributeInfo *TributeInfo        `json:"tribute_info"`
}

// PlayerDealStats contains statistics for a single player in a deal
type PlayerDealStats struct {
	PlayerSeat   int `json:"player_seat"`
	CardsPlayed  int `json:"cards_played"`
	TricksWon    int `json:"tricks_won"`
	PassCount    int `json:"pass_count"`
	TimeoutCount int `json:"timeout_count"`
	FinishRank   int `json:"finish_rank"` // 1-4, with 1 being first to finish
}

// TributeInfo contains information about the tribute phase
type TributeInfo struct {
	HasTribute   bool          `json:"has_tribute"`
	TributeMap   map[int]int   `json:"tribute_map"`   // giver -> receiver
	TributeCards map[int]*Card `json:"tribute_cards"` // giver -> card
	ReturnCards  map[int]*Card `json:"return_cards"`  // receiver -> card
}

// Enhanced DealResult with statistics
type DealResult struct {
	Rankings    []int           `json:"rankings"`     // Order of finishing (seat numbers)
	WinningTeam int             `json:"winning_team"` // 0 or 1
	VictoryType VictoryType     `json:"victory_type"`
	Upgrades    [2]int          `json:"upgrades"` // Level upgrades for each team
	Duration    time.Duration   `json:"duration"`
	TrickCount  int             `json:"trick_count"`
	Statistics  *DealStatistics `json:"statistics"`
}

// GetTeamRankings returns the rankings grouped by team
func (dr *DealResult) GetTeamRankings(match *Match) map[int][]int {
	teamRankings := make(map[int][]int)
	teamRankings[0] = make([]int, 0)
	teamRankings[1] = make([]int, 0)

	for rank, playerSeat := range dr.Rankings {
		team := match.GetTeamForPlayer(playerSeat)
		teamRankings[team] = append(teamRankings[team], rank+1) // Convert to 1-based ranking
	}

	return teamRankings
}

// GetWinningPlayers returns the player seats of the winning team
func (dr *DealResult) GetWinningPlayers(match *Match) []int {
	winningPlayers := make([]int, 0)

	for _, playerSeat := range dr.Rankings {
		if match.GetTeamForPlayer(playerSeat) == dr.WinningTeam {
			winningPlayers = append(winningPlayers, playerSeat)
		}
	}

	return winningPlayers
}

// GetLosingPlayers returns the player seats of the losing team
func (dr *DealResult) GetLosingPlayers(match *Match) []int {
	losingPlayers := make([]int, 0)
	losingTeam := 1 - dr.WinningTeam // Flip team (0->1, 1->0)

	for _, playerSeat := range dr.Rankings {
		if match.GetTeamForPlayer(playerSeat) == losingTeam {
			losingPlayers = append(losingPlayers, playerSeat)
		}
	}

	return losingPlayers
}

// IsDoubleDown returns true if this was a double down victory (rank1, rank2 same team)
func (dr *DealResult) IsDoubleDown() bool {
	return dr.VictoryType == VictoryTypeDoubleDown
}

// IsSingleLast returns true if this was a single last victory (rank1, rank3 same team)
func (dr *DealResult) IsSingleLast() bool {
	return dr.VictoryType == VictoryTypeSingleLast
}

// IsPartnerLast returns true if this was a partner last victory (rank1, rank4 same team)
func (dr *DealResult) IsPartnerLast() bool {
	return dr.VictoryType == VictoryTypePartnerLast
}

// GetUpgradeForTeam returns the level upgrade for a specific team
func (dr *DealResult) GetUpgradeForTeam(team int) int {
	if team < 0 || team > 1 {
		return 0
	}
	return dr.Upgrades[team]
}

// String returns a human-readable description of the deal result
func (dr *DealResult) String() string {
	victoryDesc := ""
	switch dr.VictoryType {
	case VictoryTypePartnerLast:
		victoryDesc = "Partner Last"
	case VictoryTypeSingleLast:
		victoryDesc = "Single Last"
	case VictoryTypeDoubleDown:
		victoryDesc = "Double Down"
	}

	return fmt.Sprintf("Team %d wins with %s (+%d levels)",
		dr.WinningTeam, victoryDesc, dr.Upgrades[dr.WinningTeam])
}

// MatchResult represents the result of a completed match
type MatchResult struct {
	Winner      int              `json:"winner"`       // Winning team (0 or 1)
	FinalLevels [2]int           `json:"final_levels"` // Final levels of both teams
	Duration    time.Duration    `json:"duration"`     // Total match duration
	Statistics  *MatchStatistics `json:"statistics"`   // Detailed match statistics
}

// MatchStatistics contains detailed statistics for a completed match
type MatchStatistics struct {
	TotalDeals    int                `json:"total_deals"`
	TotalDuration time.Duration      `json:"total_duration"`
	FinalLevels   [2]int             `json:"final_levels"`
	TeamStats     [2]*TeamMatchStats `json:"team_stats"`
}

// TeamMatchStats contains statistics for a team in a match
type TeamMatchStats struct {
	Team        int `json:"team"`         // Team number (0 or 1)
	DealsWon    int `json:"deals_won"`    // Number of deals won
	TotalTricks int `json:"total_tricks"` // Total tricks won across all deals
	Upgrades    int `json:"upgrades"`     // Total level upgrades gained
}

// GetWinningTeamStats returns the statistics for the winning team
func (mr *MatchResult) GetWinningTeamStats() *TeamMatchStats {
	if mr.Winner >= 0 && mr.Winner < 2 && mr.Statistics != nil {
		return mr.Statistics.TeamStats[mr.Winner]
	}
	return nil
}

// GetLosingTeamStats returns the statistics for the losing team
func (mr *MatchResult) GetLosingTeamStats() *TeamMatchStats {
	losingTeam := 1 - mr.Winner // Flip team (0->1, 1->0)
	if losingTeam >= 0 && losingTeam < 2 && mr.Statistics != nil {
		return mr.Statistics.TeamStats[losingTeam]
	}
	return nil
}

// String returns a human-readable description of the match result
func (mr *MatchResult) String() string {
	return fmt.Sprintf("Team %d wins! Final levels: Team 0: %d, Team 1: %d (Duration: %v)",
		mr.Winner, mr.FinalLevels[0], mr.FinalLevels[1], mr.Duration)
}
