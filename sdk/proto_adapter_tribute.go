// proto_adapter_tribute.go - 进贡阶段 Proto 适配器
//
// 职责:
// - TributePhase 类型的双向转换（SDK ↔ Proto）
// - 包含复杂的 map 类型转换和空值优化
//
// 依赖:
// - proto_adapter_basic.go: ToProtoCard, FromProtoCard, ToProtoCards, FromProtoCards
// - proto_adapter_enums.go: ToProtoTributeStatus, FromProtoTributeStatus
//
// 被依赖:
// - 无（叶子节点）
//
// 注意:
// - tributeInfoToPhase 和 tributePhaseToInfo 辅助函数已移至 proto_adapter_result.go
package sdk

import (
	pb "guandan-world/proto/gen/go/common"
	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== TributePhase Adapters ====================

// ToProtoTributePhase 转换 SDK TributePhase 到 Proto TributePhase
// 特殊处理：
// - TributeMap: map[int]int → map[int32]int32
// - TributeCards: map[int]*Card → map[int32]*Card（需要遍历）
// - ReturnCards: map[int]*Card → map[int32]*Card（需要遍历）
// - SelectionResults: map[int]int → map[int32]int32
// - 空 map 优化：只在非空时分配内存
func ToProtoTributePhase(tp *TributePhase) *pbgame.TributePhase {
	if tp == nil {
		return nil
	}

	// 转换 TributeMap（优化：只在非空时分配）
	var tributeMap map[int32]int32
	if len(tp.TributeMap) > 0 {
		tributeMap = make(map[int32]int32, len(tp.TributeMap))
		for k, v := range tp.TributeMap {
			tributeMap[int32(k)] = int32(v)
		}
	}

	// 转换 TributeCards（过滤 nil 值，防止序列化后变成空 Card；优化：只在非空时分配）
	var tributeCards map[int32]*pb.Card
	if len(tp.TributeCards) > 0 {
		tributeCards = make(map[int32]*pb.Card, len(tp.TributeCards))
		for k, v := range tp.TributeCards {
			if v != nil {
				tributeCards[int32(k)] = ToProtoCard(v)
			}
		}
	}

	// 转换 ReturnCards（过滤 nil 值，防止序列化后变成空 Card；优化：只在非空时分配）
	var returnCards map[int32]*pb.Card
	if len(tp.ReturnCards) > 0 {
		returnCards = make(map[int32]*pb.Card, len(tp.ReturnCards))
		for k, v := range tp.ReturnCards {
			if v != nil {
				returnCards[int32(k)] = ToProtoCard(v)
			}
		}
	}

	// 转换 SelectionResults（优化：只在非空时分配）
	var selectionResults map[int32]int32
	if len(tp.SelectionResults) > 0 {
		selectionResults = make(map[int32]int32, len(tp.SelectionResults))
		for k, v := range tp.SelectionResults {
			selectionResults[int32(k)] = int32(v)
		}
	}

	return &pbgame.TributePhase{
		Status:           ToProtoTributeStatus(tp.Status),
		TributeMap:       tributeMap,
		TributeCards:     tributeCards,
		ReturnCards:      returnCards,
		PoolCards:        ToProtoCards(tp.PoolCards),
		SelectingPlayer:  int32(tp.SelectingPlayer),
		IsImmune:         tp.IsImmune,
		SelectionResults: selectionResults,
	}
}

// FromProtoTributePhase 转换 Proto TributePhase 到 SDK TributePhase
// 特殊处理：
// - map[int32]int32 → map[int]int
// - map[int32]*Card → map[int]*Card（需要遍历）
func FromProtoTributePhase(ptp *pbgame.TributePhase) *TributePhase {
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

	// 转换 SelectionResults
	selectionResults := make(map[int]int)
	for k, v := range ptp.SelectionResults {
		selectionResults[int(k)] = int(v)
	}

	return &TributePhase{
		Status:           FromProtoTributeStatus(ptp.Status),
		TributeMap:       tributeMap,
		TributeCards:     tributeCards,
		ReturnCards:      returnCards,
		PoolCards:        FromProtoCards(ptp.PoolCards),
		SelectingPlayer:  int(ptp.SelectingPlayer),
		IsImmune:         ptp.IsImmune,
		SelectionResults: selectionResults,
	}
}
