package sdk

import (
	"errors"
	"fmt"
	"time"
)

// NewMatch creates a new match with the given players
func NewMatch(players []Player) (*Match, error) {
	if len(players) != 4 {
		return nil, errors.New("exactly 4 players are required")
	}

	// Validate player seats
	seatTaken := make(map[int]bool)
	for _, player := range players {
		if player.Seat < 0 || player.Seat > 3 {
			return nil, fmt.Errorf("invalid seat number %d for player %s", player.Seat, player.ID)
		}
		if seatTaken[player.Seat] {
			return nil, fmt.Errorf("seat %d is already taken", player.Seat)
		}
		seatTaken[player.Seat] = true
	}

	// Create match
	match := &Match{
		ID:          generateMatchID(),
		Status:      MatchStatusWaiting,
		TeamLevels:  [2]int{2, 2}, // Both teams start at level 2
		Winner:      -1,           // No winner yet
		StartTime:   time.Now(),
		DealHistory: make([]*Deal, 0),
	}

	// Assign players to seats
	for _, player := range players {
		playerCopy := player
		playerCopy.Online = true
		playerCopy.AutoPlay = false
		match.Players[player.Seat] = &playerCopy
	}

	return match, nil
}

// StartNewDeal starts a new deal in the match
func (m *Match) StartNewDeal() error {
	if m.Status == MatchStatusFinished {
		return errors.New("match is already finished")
	}

	// Determine the level for this deal
	level := m.getHighestTeamLevel()

	// Create new deal
	deal, err := NewDeal(level, m.getLastDealResult())
	if err != nil {
		return fmt.Errorf("failed to create deal: %w", err)
	}

	// Set current deal and update status
	m.CurrentDeal = deal
	m.Status = MatchStatusPlaying

	// Start the deal (deal cards and begin play)
	err = deal.StartDeal()
	if err != nil {
		return fmt.Errorf("failed to start deal: %w", err)
	}

	return nil
}

// HandlePlayerDisconnect handles a player disconnection
func (m *Match) HandlePlayerDisconnect(playerSeat int) error {
	if playerSeat < 0 || playerSeat > 3 {
		return fmt.Errorf("invalid player seat: %d", playerSeat)
	}

	if m.Players[playerSeat] == nil {
		return fmt.Errorf("no player at seat %d", playerSeat)
	}

	// Mark player as offline and enable auto-play
	m.Players[playerSeat].Online = false
	m.Players[playerSeat].AutoPlay = true

	return nil
}

// HandlePlayerReconnect handles a player reconnection
func (m *Match) HandlePlayerReconnect(playerSeat int) error {
	if playerSeat < 0 || playerSeat > 3 {
		return fmt.Errorf("invalid player seat: %d", playerSeat)
	}

	if m.Players[playerSeat] == nil {
		return fmt.Errorf("no player at seat %d", playerSeat)
	}

	// Mark player as online and disable auto-play
	m.Players[playerSeat].Online = true
	m.Players[playerSeat].AutoPlay = false

	return nil
}

// SetPlayerAutoPlay sets the auto-play status for a player
func (m *Match) SetPlayerAutoPlay(playerSeat int, enabled bool) error {
	if playerSeat < 0 || playerSeat > 3 {
		return fmt.Errorf("invalid player seat: %d", playerSeat)
	}

	if m.Players[playerSeat] == nil {
		return fmt.Errorf("no player at seat %d", playerSeat)
	}

	m.Players[playerSeat].AutoPlay = enabled
	return nil
}

// FinishDeal finishes the current deal and updates match state
func (m *Match) FinishDeal(result *DealResult) error {
	if m.CurrentDeal == nil {
		return errors.New("no active deal to finish")
	}

	// Add deal to history
	m.DealHistory = append(m.DealHistory, m.CurrentDeal)

	// Update team levels based on result
	m.updateTeamLevels(result)

	// Check if match is finished (any team reached A level)
	if m.isMatchFinished() {
		m.Status = MatchStatusFinished
		m.Winner = m.getWinningTeam()
		now := time.Now()
		m.EndTime = &now
	} else {
		m.Status = MatchStatusWaiting
	}

	// Clear current deal
	m.CurrentDeal = nil

	return nil
}

// GetTeamForPlayer returns the team number (0 or 1) for a given player seat
func (m *Match) GetTeamForPlayer(playerSeat int) int {
	// Team 0: seats 0, 2
	// Team 1: seats 1, 3
	return playerSeat % 2
}

// GetTeammates returns the seat numbers of teammates for a given player
func (m *Match) GetTeammates(playerSeat int) []int {
	team := m.GetTeamForPlayer(playerSeat)
	teammates := make([]int, 0, 2)

	for seat := 0; seat < 4; seat++ {
		if m.GetTeamForPlayer(seat) == team {
			teammates = append(teammates, seat)
		}
	}

	return teammates
}

// GetOpponents returns the seat numbers of opponents for a given player
func (m *Match) GetOpponents(playerSeat int) []int {
	team := m.GetTeamForPlayer(playerSeat)
	opponents := make([]int, 0, 2)

	for seat := 0; seat < 4; seat++ {
		if m.GetTeamForPlayer(seat) != team {
			opponents = append(opponents, seat)
		}
	}

	return opponents
}

// IsPlayerOnline checks if a player is online
func (m *Match) IsPlayerOnline(playerSeat int) bool {
	if playerSeat < 0 || playerSeat > 3 || m.Players[playerSeat] == nil {
		return false
	}
	return m.Players[playerSeat].Online
}

// IsPlayerAutoPlay checks if a player is in auto-play mode
func (m *Match) IsPlayerAutoPlay(playerSeat int) bool {
	if playerSeat < 0 || playerSeat > 3 || m.Players[playerSeat] == nil {
		return false
	}
	return m.Players[playerSeat].AutoPlay
}

// getHighestTeamLevel returns the highest level among both teams
func (m *Match) getHighestTeamLevel() int {
	if m.TeamLevels[0] > m.TeamLevels[1] {
		return m.TeamLevels[0]
	}
	return m.TeamLevels[1]
}

// getLastDealResult returns the result of the last deal, or nil if no deals yet
func (m *Match) getLastDealResult() *DealResult {
	if len(m.DealHistory) == 0 {
		return nil
	}

	lastDeal := m.DealHistory[len(m.DealHistory)-1]
	if lastDeal.EndTime == nil {
		return nil // Deal not finished
	}

	// Calculate proper deal result using the result calculator
	result, err := lastDeal.CalculateResult(m)
	if err != nil {
		// Fallback to basic result if calculation fails
		return &DealResult{
			Rankings:    lastDeal.Rankings,
			WinningTeam: m.GetTeamForPlayer(lastDeal.Rankings[0]),
			VictoryType: VictoryTypePartnerLast,
			Upgrades:    [2]int{1, 0},
		}
	}

	return result
}

// updateTeamLevels updates team levels based on deal result
func (m *Match) updateTeamLevels(result *DealResult) {
	if result == nil {
		return
	}

	// Apply upgrades to teams
	for team := 0; team < 2; team++ {
		m.TeamLevels[team] += result.Upgrades[team]

		// Cap at A level (14)
		if m.TeamLevels[team] > 14 {
			m.TeamLevels[team] = 14
		}
	}
}

// isMatchFinished checks if the match is finished (any team reached A level)
func (m *Match) isMatchFinished() bool {
	return m.TeamLevels[0] >= 14 || m.TeamLevels[1] >= 14
}

// getWinningTeam returns the winning team number, or -1 if no winner
func (m *Match) getWinningTeam() int {
	if m.TeamLevels[0] >= 14 && m.TeamLevels[1] >= 14 {
		// Both teams reached A - determine winner by who reached first or higher
		if m.TeamLevels[0] > m.TeamLevels[1] {
			return 0
		} else if m.TeamLevels[1] > m.TeamLevels[0] {
			return 1
		} else {
			// Same level - winner is determined by last deal result
			if len(m.DealHistory) > 0 {
				lastResult := m.getLastDealResult()
				if lastResult != nil {
					return lastResult.WinningTeam
				}
			}
			return 0 // Default to team 0
		}
	} else if m.TeamLevels[0] >= 14 {
		return 0
	} else if m.TeamLevels[1] >= 14 {
		return 1
	}

	return -1 // No winner yet
}

// generateMatchID generates a unique ID for the match
func generateMatchID() string {
	return fmt.Sprintf("match_%d", time.Now().UnixNano())
}

// VictoryType represents the type of victory in a deal
type VictoryType string

const (
	VictoryTypeDoubleDown  VictoryType = "double_down"  // rank1, rank2 same team (+3 levels)
	VictoryTypeSingleLast  VictoryType = "single_last"  // rank1, rank3 same team (+2 levels)
	VictoryTypePartnerLast VictoryType = "partner_last" // rank1, rank4 same team (+1 level)
)

// CanStartNewDeal checks if a new deal can be started
func (m *Match) CanStartNewDeal() bool {
	return m.Status == MatchStatusWaiting && !m.isMatchFinished()
}

// GetMatchStatistics returns comprehensive match statistics
func (m *Match) GetMatchStatistics() *MatchStatistics {
	stats := &MatchStatistics{
		TotalDeals:    len(m.DealHistory),
		TotalDuration: time.Duration(0),
		FinalLevels:   m.TeamLevels,
		TeamStats:     [2]*TeamMatchStats{},
	}

	// Initialize team stats
	for team := 0; team < 2; team++ {
		stats.TeamStats[team] = &TeamMatchStats{
			Team:        team,
			DealsWon:    0,
			TotalTricks: 0,
			Upgrades:    0,
		}
	}

	// Calculate total duration
	if m.EndTime != nil {
		stats.TotalDuration = m.EndTime.Sub(m.StartTime)
	} else {
		stats.TotalDuration = time.Since(m.StartTime)
	}

	// Calculate team statistics from deal history
	for _, deal := range m.DealHistory {
		if result, err := deal.CalculateResult(m); err == nil {
			// Count deals won
			stats.TeamStats[result.WinningTeam].DealsWon++

			// Count upgrades
			for team := 0; team < 2; team++ {
				stats.TeamStats[team].Upgrades += result.Upgrades[team]
			}

			// Count tricks won by each team
			if result.Statistics != nil {
				for _, playerStats := range result.Statistics.PlayerStats {
					if playerStats != nil {
						team := m.GetTeamForPlayer(playerStats.PlayerSeat)
						stats.TeamStats[team].TotalTricks += playerStats.TricksWon
					}
				}
			}
		}
	}

	return stats
}

// GetMatchResult creates a complete match result (for finished matches)
func (m *Match) GetMatchResult() *MatchResult {
	if m.Status != MatchStatusFinished {
		return nil
	}

	duration := time.Duration(0)
	if m.EndTime != nil {
		duration = m.EndTime.Sub(m.StartTime)
	}

	return &MatchResult{
		Winner:      m.Winner,
		FinalLevels: m.TeamLevels,
		Duration:    duration,
		Statistics:  m.GetMatchStatistics(),
	}
}

// GetCurrentLevel returns the current level for the match (highest team level)
func (m *Match) GetCurrentLevel() int {
	return m.getHighestTeamLevel()
}

// GetTeamLevel returns the current level for a specific team
func (m *Match) GetTeamLevel(team int) int {
	if team < 0 || team > 1 {
		return 0
	}
	return m.TeamLevels[team]
}

// GetPlayerTeamLevel returns the current level for a player's team
func (m *Match) GetPlayerTeamLevel(playerSeat int) int {
	team := m.GetTeamForPlayer(playerSeat)
	return m.GetTeamLevel(team)
}

// IsTeamAtLevel checks if a team has reached a specific level
func (m *Match) IsTeamAtLevel(team int, level int) bool {
	return m.GetTeamLevel(team) >= level
}

// IsAnyTeamAtALevel checks if any team has reached A level (14)
func (m *Match) IsAnyTeamAtALevel() bool {
	return m.isMatchFinished()
}

// GetLeadingTeam returns the team with the higher level, or -1 if tied
func (m *Match) GetLeadingTeam() int {
	if m.TeamLevels[0] > m.TeamLevels[1] {
		return 0
	} else if m.TeamLevels[1] > m.TeamLevels[0] {
		return 1
	}
	return -1 // Tied
}

// GetLevelDifference returns the level difference between teams (team0 - team1)
func (m *Match) GetLevelDifference() int {
	return m.TeamLevels[0] - m.TeamLevels[1]
}

// GetDealCount returns the number of completed deals
func (m *Match) GetDealCount() int {
	return len(m.DealHistory)
}

// GetLastDeal returns the last completed deal, or nil if no deals
func (m *Match) GetLastDeal() *Deal {
	if len(m.DealHistory) == 0 {
		return nil
	}
	return m.DealHistory[len(m.DealHistory)-1]
}

// GetDealByIndex returns a deal by its index in history, or nil if invalid index
func (m *Match) GetDealByIndex(index int) *Deal {
	if index < 0 || index >= len(m.DealHistory) {
		return nil
	}
	return m.DealHistory[index]
}

// HasActiveDeal checks if there's currently an active deal
func (m *Match) HasActiveDeal() bool {
	return m.CurrentDeal != nil && m.CurrentDeal.Status != DealStatusFinished
}

// GetActiveDeal returns the current active deal, or nil if none
func (m *Match) GetActiveDeal() *Deal {
	if m.HasActiveDeal() {
		return m.CurrentDeal
	}
	return nil
}

// GetPlayerCount returns the number of players in the match
func (m *Match) GetPlayerCount() int {
	count := 0
	for _, player := range m.Players {
		if player != nil {
			count++
		}
	}
	return count
}

// GetOnlinePlayerCount returns the number of online players
func (m *Match) GetOnlinePlayerCount() int {
	count := 0
	for _, player := range m.Players {
		if player != nil && player.Online {
			count++
		}
	}
	return count
}

// GetAutoPlayPlayerCount returns the number of players in auto-play mode
func (m *Match) GetAutoPlayPlayerCount() int {
	count := 0
	for _, player := range m.Players {
		if player != nil && player.AutoPlay {
			count++
		}
	}
	return count
}

// GetPlayerByID returns a player by their ID, or nil if not found
func (m *Match) GetPlayerByID(playerID string) *Player {
	for _, player := range m.Players {
		if player != nil && player.ID == playerID {
			return player
		}
	}
	return nil
}

// GetPlayerBySeat returns a player by their seat number, or nil if invalid/empty
func (m *Match) GetPlayerBySeat(seat int) *Player {
	if seat < 0 || seat > 3 {
		return nil
	}
	return m.Players[seat]
}

// String returns a human-readable description of the match
func (m *Match) String() string {
	status := string(m.Status)
	if m.Status == MatchStatusFinished {
		return fmt.Sprintf("Match %s: %s (Winner: Team %d, Levels: %d-%d)",
			m.ID, status, m.Winner, m.TeamLevels[0], m.TeamLevels[1])
	}
	return fmt.Sprintf("Match %s: %s (Levels: %d-%d, Deals: %d)",
		m.ID, status, m.TeamLevels[0], m.TeamLevels[1], len(m.DealHistory))
}

// generateDealID generates a unique ID for a deal
func generateDealID() string {
	return fmt.Sprintf("deal_%d", time.Now().UnixNano())
}
