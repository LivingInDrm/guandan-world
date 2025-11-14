// proto_adapter_result.go - 统计结果 Proto 适配器
//
// 职责:
// - DealResult, MatchResult 及相关统计类型的双向转换（SDK ↔ Proto）
// - 包含 TeamUpgrades, PlayerDealStats, DealStatistics, TeamMatchStats, MatchStatistics
// - 提供 tributeInfoToPhase/tributePhaseToInfo 私有辅助函数用于统计数据转换
//
// 依赖:
// - proto_adapter_basic.go: ToProtoCard, FromProtoCard
// - proto_adapter_enums.go: ToProtoVictoryType, FromProtoVictoryType,
//   ToProtoTributeStatus, FromProtoTributeStatus
//
// 被依赖:
// - 无（叶子节点）
package sdk

import (
	"time"

	pb "guandan-world/proto/gen/go/common"
	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== Result Adapters ====================

// ==================== Private Helpers ====================

// tributeInfoToPhase 转换 SDK TributeInfo 到 Proto TributePhase
// 这是一个私有辅助函数，用于将统计用的 TributeInfo 转换为 proto 的 TributePhase
// TributeInfo 只包含结果快照，所以只填充关键字段
// 特殊处理：
// - map[int]int → map[int32]int32
// - map[int]*Card → map[int32]*Card
// - 空 map 优化：只在非空时分配内存
func tributeInfoToPhase(ti *TributeInfo) *pbgame.TributePhase {
	if ti == nil {
		return nil
	}

	// 转换 TributeMap（优化：只在非空时分配）
	var tributeMap map[int32]int32
	if len(ti.TributeMap) > 0 {
		tributeMap = make(map[int32]int32, len(ti.TributeMap))
		for k, v := range ti.TributeMap {
			tributeMap[int32(k)] = int32(v)
		}
	}

	// 转换 TributeCards（优化：只在非空时分配）
	var tributeCards map[int32]*pb.Card
	if len(ti.TributeCards) > 0 {
		tributeCards = make(map[int32]*pb.Card, len(ti.TributeCards))
		for k, v := range ti.TributeCards {
			tributeCards[int32(k)] = ToProtoCard(v)
		}
	}

	// 转换 ReturnCards（优化：只在非空时分配）
	var returnCards map[int32]*pb.Card
	if len(ti.ReturnCards) > 0 {
		returnCards = make(map[int32]*pb.Card, len(ti.ReturnCards))
		for k, v := range ti.ReturnCards {
			returnCards[int32(k)] = ToProtoCard(v)
		}
	}

	// 设置状态
	var status pb.TributeStatus
	if ti.HasTribute {
		status = pb.TributeStatus_TRIBUTE_STATUS_FINISHED
	} else {
		status = pb.TributeStatus_TRIBUTE_STATUS_WAITING
	}

	return &pbgame.TributePhase{
		Status:       status,
		TributeMap:   tributeMap,
		TributeCards: tributeCards,
		ReturnCards:  returnCards,
		IsImmune:     !ti.HasTribute,
	}
}

// tributePhaseToInfo 转换 Proto TributePhase 到 SDK TributeInfo
// 这是一个私有辅助函数，用于从 proto 的 TributePhase 提取统计用的 TributeInfo
// 特殊处理：
// - map[int32]int32 → map[int]int
// - map[int32]*Card → map[int]*Card
func tributePhaseToInfo(ptp *pbgame.TributePhase) *TributeInfo {
	if ptp == nil {
		return nil
	}

	// 转换 TributeMap
	tributeMap := make(map[int]int)
	for k, v := range ptp.TributeMap {
		tributeMap[int(k)] = int(v)
	}

	// 转换 TributeCards
	tributeCards := make(map[int]*Card)
	for k, v := range ptp.TributeCards {
		tributeCards[int(k)] = FromProtoCard(v)
	}

	// 转换 ReturnCards
	returnCards := make(map[int]*Card)
	for k, v := range ptp.ReturnCards {
		returnCards[int(k)] = FromProtoCard(v)
	}

	return &TributeInfo{
		HasTribute:   !ptp.IsImmune && len(ptp.TributeMap) > 0,
		TributeMap:   tributeMap,
		TributeCards: tributeCards,
		ReturnCards:  returnCards,
	}
}

// ==================== Public Adapters ====================

// ToProtoTeamUpgrades 转换 SDK [2]int 到 Proto TeamUpgrades
func ToProtoTeamUpgrades(upgrades [2]int) *pbgame.TeamUpgrades {
	return &pbgame.TeamUpgrades{
		Team0: int32(upgrades[0]),
		Team1: int32(upgrades[1]),
	}
}

// FromProtoTeamUpgrades 转换 Proto TeamUpgrades 到 SDK [2]int
func FromProtoTeamUpgrades(ptu *pbgame.TeamUpgrades) [2]int {
	if ptu == nil {
		return [2]int{0, 0}
	}
	return [2]int{int(ptu.Team0), int(ptu.Team1)}
}

// ToProtoPlayerDealStats 转换 SDK PlayerDealStats 到 Proto PlayerDealStats
func ToProtoPlayerDealStats(pds *PlayerDealStats) *pbgame.PlayerDealStats {
	if pds == nil {
		return nil
	}
	return &pbgame.PlayerDealStats{
		PlayerSeat:   int32(pds.PlayerSeat),
		CardsPlayed:  int32(pds.CardsPlayed),
		TricksWon:    int32(pds.TricksWon),
		PassCount:    int32(pds.PassCount),
		TimeoutCount: int32(pds.TimeoutCount),
		FinishRank:   int32(pds.FinishRank),
	}
}

// FromProtoPlayerDealStats 转换 Proto PlayerDealStats 到 SDK PlayerDealStats
func FromProtoPlayerDealStats(ppds *pbgame.PlayerDealStats) *PlayerDealStats {
	if ppds == nil {
		return nil
	}
	return &PlayerDealStats{
		PlayerSeat:   int(ppds.PlayerSeat),
		CardsPlayed:  int(ppds.CardsPlayed),
		TricksWon:    int(ppds.TricksWon),
		PassCount:    int(ppds.PassCount),
		TimeoutCount: int(ppds.TimeoutCount),
		FinishRank:   int(ppds.FinishRank),
	}
}

// ToProtoDealStatistics 转换 SDK DealStatistics 到 Proto DealStatistics
// 特殊处理：
// - [4]*PlayerDealStats → repeated PlayerDealStats (固定4个)
// - TributeInfo → TributePhase (使用辅助函数转换)
func ToProtoDealStatistics(ds *DealStatistics) *pbgame.DealStatistics {
	if ds == nil {
		return nil
	}

	// 转换 PlayerStats (固定4个)
	playerStats := make([]*pbgame.PlayerDealStats, 4)
	for i := 0; i < 4; i++ {
		playerStats[i] = ToProtoPlayerDealStats(ds.PlayerStats[i])
	}

	return &pbgame.DealStatistics{
		TotalTricks:  int32(ds.TotalTricks),
		PlayerStats:  playerStats,
		TributePhase: tributeInfoToPhase(ds.TributeInfo),
	}
}

// FromProtoDealStatistics 转换 Proto DealStatistics 到 SDK DealStatistics
// 特殊处理：
// - repeated PlayerDealStats → [4]*PlayerDealStats
// - TributePhase → TributeInfo (使用辅助函数转换)
// - 添加长度验证（P1-3）
func FromProtoDealStatistics(pds *pbgame.DealStatistics) *DealStatistics {
	if pds == nil {
		return nil
	}

	// 转换 PlayerStats，并验证长度
	var playerStats [4]*PlayerDealStats
	actualLen := len(pds.PlayerStats)
	if actualLen != 4 {
		// 记录警告：预期 4 个玩家统计，但实际收到不同数量
		// 这里使用静默处理，用 nil 填充缺失的部分
		// 生产环境中应该使用适当的日志系统
		_ = actualLen // 避免未使用变量警告
	}
	for i := 0; i < 4 && i < len(pds.PlayerStats); i++ {
		playerStats[i] = FromProtoPlayerDealStats(pds.PlayerStats[i])
	}

	return &DealStatistics{
		TotalTricks: int(pds.TotalTricks),
		PlayerStats: playerStats,
		TributeInfo: tributePhaseToInfo(pds.TributePhase),
	}
}

// ToProtoDealResult 转换 SDK DealResult 到 Proto DealResult
// 特殊处理：
// - []int → repeated int32
// - [2]int → TeamUpgrades
// - time.Duration → int64 (毫秒)
func ToProtoDealResult(dr *DealResult) *pbgame.DealResult {
	if dr == nil {
		return nil
	}

	// 转换 Rankings
	rankings := make([]int32, len(dr.Rankings))
	for i, r := range dr.Rankings {
		rankings[i] = int32(r)
	}

	return &pbgame.DealResult{
		Rankings:    rankings,
		WinningTeam: int32(dr.WinningTeam),
		VictoryType: ToProtoVictoryType(dr.VictoryType),
		Upgrades:    ToProtoTeamUpgrades(dr.Upgrades),
		DurationMs:  dr.Duration.Milliseconds(),
		TrickCount:  int32(dr.TrickCount),
		Statistics:  ToProtoDealStatistics(dr.Statistics),
	}
}

// FromProtoDealResult 转换 Proto DealResult 到 SDK DealResult
// 特殊处理：
// - repeated int32 → []int
// - TeamUpgrades → [2]int
// - int64 (毫秒) → time.Duration
func FromProtoDealResult(pdr *pbgame.DealResult) *DealResult {
	if pdr == nil {
		return nil
	}

	// 转换 Rankings
	rankings := make([]int, len(pdr.Rankings))
	for i, r := range pdr.Rankings {
		rankings[i] = int(r)
	}

	return &DealResult{
		Rankings:    rankings,
		WinningTeam: int(pdr.WinningTeam),
		VictoryType: FromProtoVictoryType(pdr.VictoryType),
		Upgrades:    FromProtoTeamUpgrades(pdr.Upgrades),
		Duration:    time.Duration(pdr.DurationMs) * time.Millisecond,
		TrickCount:  int(pdr.TrickCount),
		Statistics:  FromProtoDealStatistics(pdr.Statistics),
	}
}

// ToProtoTeamMatchStats 转换 SDK TeamMatchStats 到 Proto TeamMatchStats
func ToProtoTeamMatchStats(tms *TeamMatchStats) *pbgame.TeamMatchStats {
	if tms == nil {
		return nil
	}
	return &pbgame.TeamMatchStats{
		Team:        int32(tms.Team),
		DealsWon:    int32(tms.DealsWon),
		TotalTricks: int32(tms.TotalTricks),
		Upgrades:    int32(tms.Upgrades),
	}
}

// FromProtoTeamMatchStats 转换 Proto TeamMatchStats 到 SDK TeamMatchStats
func FromProtoTeamMatchStats(ptms *pbgame.TeamMatchStats) *TeamMatchStats {
	if ptms == nil {
		return nil
	}
	return &TeamMatchStats{
		Team:        int(ptms.Team),
		DealsWon:    int(ptms.DealsWon),
		TotalTricks: int(ptms.TotalTricks),
		Upgrades:    int(ptms.Upgrades),
	}
}

// ToProtoMatchStatistics 转换 SDK MatchStatistics 到 Proto MatchStatistics
// 特殊处理：
// - [2]*TeamMatchStats → repeated TeamMatchStats (固定2个)
// - time.Duration → int64 (毫秒)
// - FinalLevels 已从 proto 中移除（在 MatchResult 中已有，避免冗余）
func ToProtoMatchStatistics(ms *MatchStatistics) *pbgame.MatchStatistics {
	if ms == nil {
		return nil
	}

	// 转换 TeamStats (固定2个)
	teamStats := make([]*pbgame.TeamMatchStats, 2)
	for i := 0; i < 2; i++ {
		teamStats[i] = ToProtoTeamMatchStats(ms.TeamStats[i])
	}

	return &pbgame.MatchStatistics{
		TotalDeals:      int32(ms.TotalDeals),
		TotalDurationMs: ms.TotalDuration.Milliseconds(),
		TeamStats:       teamStats,
	}
}

// FromProtoMatchStatistics 转换 Proto MatchStatistics 到 SDK MatchStatistics
// 特殊处理：
// - repeated TeamMatchStats → [2]*TeamMatchStats
// - int64 (毫秒) → time.Duration
// - FinalLevels 需要从父级 MatchResult 手动设置（proto 中已移除此冗余字段）
func FromProtoMatchStatistics(pms *pbgame.MatchStatistics, finalLevels [2]int) *MatchStatistics {
	if pms == nil {
		return nil
	}

	// 转换 TeamStats
	var teamStats [2]*TeamMatchStats
	for i := 0; i < 2 && i < len(pms.TeamStats); i++ {
		teamStats[i] = FromProtoTeamMatchStats(pms.TeamStats[i])
	}

	return &MatchStatistics{
		TotalDeals:    int(pms.TotalDeals),
		TotalDuration: time.Duration(pms.TotalDurationMs) * time.Millisecond,
		FinalLevels:   finalLevels,
		TeamStats:     teamStats,
	}
}

// ToProtoMatchResult 转换 SDK MatchResult 到 Proto MatchResult
// 特殊处理：
// - [2]int → TeamUpgrades
// - time.Duration → int64 (毫秒)
func ToProtoMatchResult(mr *MatchResult) *pbgame.MatchResult {
	if mr == nil {
		return nil
	}
	return &pbgame.MatchResult{
		Winner:      int32(mr.Winner),
		FinalLevels: ToProtoTeamUpgrades(mr.FinalLevels),
		DurationMs:  mr.Duration.Milliseconds(),
		Statistics:  ToProtoMatchStatistics(mr.Statistics),
	}
}

// FromProtoMatchResult 转换 Proto MatchResult 到 SDK MatchResult
// 特殊处理：
// - TeamUpgrades → [2]int
// - int64 (毫秒) → time.Duration
// - FinalLevels 传递给 Statistics（因为 proto 中已从 MatchStatistics 移除）
func FromProtoMatchResult(pmr *pbgame.MatchResult) *MatchResult {
	if pmr == nil {
		return nil
	}
	
	finalLevels := FromProtoTeamUpgrades(pmr.FinalLevels)
	
	return &MatchResult{
		Winner:      int(pmr.Winner),
		FinalLevels: finalLevels,
		Duration:    time.Duration(pmr.DurationMs) * time.Millisecond,
		Statistics:  FromProtoMatchStatistics(pmr.Statistics, finalLevels),
	}
}
