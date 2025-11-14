// proto_adapter_basic.go - 基础类型 Proto 适配器
//
// 职责:
// - Card 类型的双向转换（SDK ↔ Proto）
// - Player 类型的双向转换（SDK ↔ Proto）
// - 时间工具函数（毫秒时间戳 ↔ time.Time）
//
// 依赖:
// - 无外部依赖（基础层）
//
// 被依赖:
// - proto_adapter_comp.go: 使用 ToProtoCards, FromProtoCards
// - proto_adapter_action.go: 使用 ToProtoCards, FromProtoCards, timeFromMillis
// - proto_adapter_tribute.go: 使用 ToProtoCard, FromProtoCard
// - proto_adapter_result.go: 使用 ToProtoCard, FromProtoCard
package sdk

import (
	"time"

	pb "guandan-world/proto/gen/go/common"
)

// ==================== Card Adapters ====================

// ToProtoCard 转换 SDK Card 到 Proto Card
// 特殊处理：
// - Color: 王牌为 "Joker"
// - Level: 当前级别牌
func ToProtoCard(c *Card) *pb.Card {
	if c == nil {
		return nil
	}
	return &pb.Card{
		Number:    int32(c.Number),
		RawNumber: int32(c.RawNumber),
		Color:     c.Color,
		Level:     int32(c.Level),
		Name:      c.Name,
		DeckIndex: int32(c.DeckIndex),
	}
}

// ToProtoCards 批量转换 SDK Cards 到 Proto Cards
// 注意: 如果输入切片中有 nil 元素，输出切片对应位置也会是 nil
func ToProtoCards(cards []*Card) []*pb.Card {
	if cards == nil {
		return nil
	}
	result := make([]*pb.Card, len(cards))
	for i, card := range cards {
		result[i] = ToProtoCard(card)
	}
	return result
}

// FromProtoCard 转换 Proto Card 到 SDK Card
func FromProtoCard(pc *pb.Card) *Card {
	if pc == nil {
		return nil
	}
	card := &Card{
		Number:    int(pc.Number),
		RawNumber: int(pc.RawNumber),
		Color:     pc.Color,
		Level:     int(pc.Level),
		Name:      pc.Name,
		DeckIndex: int(pc.DeckIndex),
	}
	return card
}

// FromProtoCards 批量转换 Proto Cards 到 SDK Cards
// 注意: 如果输入切片中有 nil 元素，输出切片对应位置也会是 nil
func FromProtoCards(pcs []*pb.Card) []*Card {
	if pcs == nil {
		return nil
	}
	result := make([]*Card, len(pcs))
	for i, pc := range pcs {
		result[i] = FromProtoCard(pc)
	}
	return result
}

// ==================== Player Adapters ====================

// ToProtoPlayer 转换 SDK Player 到 Proto Player
func ToProtoPlayer(p *Player) *pb.Player {
	if p == nil {
		return nil
	}
	return &pb.Player{
		Id:       p.ID,
		Username: p.Username,
		Seat:     int32(p.Seat),
		Online:   p.Online,
		AutoPlay: p.AutoPlay,
	}
}

// ToProtoPlayers 批量转换 SDK Players 到 Proto Players (slice)
// 注意: 如果输入切片中有 nil 元素，输出切片对应位置也会是 nil
func ToProtoPlayers(players []*Player) []*pb.Player {
	if players == nil {
		return nil
	}
	result := make([]*pb.Player, len(players))
	for i, player := range players {
		result[i] = ToProtoPlayer(player)
	}
	return result
}

// ToProtoPlayersArray 转换固定长度数组 [4]*Player 到 Proto Players
// 注意: 输入数组固定为4个元素，对应掼蛋游戏的4个玩家座位
func ToProtoPlayersArray(players [4]*Player) []*pb.Player {
	result := make([]*pb.Player, 4)
	for i := 0; i < 4; i++ {
		result[i] = ToProtoPlayer(players[i])
	}
	return result
}

// FromProtoPlayer 转换 Proto Player 到 SDK Player
func FromProtoPlayer(pp *pb.Player) *Player {
	if pp == nil {
		return nil
	}
	return &Player{
		ID:       pp.Id,
		Username: pp.Username,
		Seat:     int(pp.Seat),
		Online:   pp.Online,
		AutoPlay: pp.AutoPlay,
	}
}

// FromProtoPlayers 批量转换 Proto Players 到 SDK Players (slice)
// 注意: 如果输入切片中有 nil 元素，输出切片对应位置也会是 nil
func FromProtoPlayers(pps []*pb.Player) []*Player {
	if pps == nil {
		return nil
	}
	result := make([]*Player, len(pps))
	for i, pp := range pps {
		result[i] = FromProtoPlayer(pp)
	}
	return result
}

// FromProtoPlayersArray 转换 Proto Players 到固定长度数组 [4]*Player
// 注意: 输入切片长度不足4时，剩余位置为 nil
func FromProtoPlayersArray(pps []*pb.Player) [4]*Player {
	var result [4]*Player
	for i := 0; i < 4 && i < len(pps); i++ {
		result[i] = FromProtoPlayer(pps[i])
	}
	return result
}

// ==================== Time Helper ====================

// timeFromMillis 从毫秒时间戳转换为 time.Time
// 特殊处理：
// - ms <= 0: 返回 time.Time{} (零值)，表示时间未设置
// - ms > 0: 转换为对应的时间
func timeFromMillis(ms int64) time.Time {
	if ms <= 0 {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}
