# P0 贡牌/还贡数据一致性问题修复报告

## 问题描述
在掼蛋游戏（使用两副牌，共108张）中，当玩家手牌存在同点数同花色的重复牌时，贡牌/还贡移交逻辑可能删除错误的副本，导致：
- 原始贡牌仍留在贡者手牌
- 同时被加入接收者手牌
- 形成重复引用与手牌污染

## 根本原因
原 `cardsEqual` 实现仅比较 `Number` 和 `Color`，未使用唯一标识 `DeckIndex`：
```go
// 修复前 - 错误的实现
func cardsEqual(card1, card2 *Card) bool {
    return card1.Number == card2.Number && card1.Color == card2.Color
}
```

## 修复方案
将所有卡牌比较逻辑改为基于 `DeckIndex` 的精确匹配：
```go
// 修复后 - 正确的实现
func cardsEqual(card1, card2 *Card) bool {
    return card1.DeckIndex == card2.DeckIndex
}
```

## 修改文件
1. `sdk/tribute.go:337` - `TributeManager.cardsEqual`
2. `sdk/tribute.go:429` - `TributePhase.cardsEqual`
3. `sdk/deal.go:361` - `Deal.cardsEqual`

## 验证结果
### ✅ 重复牌场景测试
- 玩家手牌有两张黑桃A（DeckIndex=12 和 DeckIndex=64）
- 指定 DeckIndex=12 作为贡牌
- **结果**: 正确删除 DeckIndex=12，保留 DeckIndex=64
- **接收者**: 正确收到 DeckIndex=12

### ✅ 现有功能回归测试
- 正常贡牌流程：通过
- 抗贡检测：通过
- 自动选择最大牌（排除红桃Trump）：通过

## 影响范围
- 贡牌/还贡移交：`tribute.go:257, 267`
- 选贡逻辑：`tribute.go:357`
- 出牌验证和移除：`deal.go:315, 348`

## 测试代码
新增 `sdk/tribute_duplicate_cards_test.go`，包含：
- `TestTributeWithDuplicateCards` - 重复牌场景验证
- `TestRemoveCardFromHandPrecision` - 精确删除验证
- `TestCardsEqualByDeckIndex` - DeckIndex 比较验证

## 结论
✅ 问题已修复，所有验证通过
✅ 现有功能未受影响
✅ 数据一致性得到保障
