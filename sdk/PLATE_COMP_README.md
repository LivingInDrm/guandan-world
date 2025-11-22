# 钢板（Plate）枚举验证方案

## 概述

`sdk/plate_comp.go` 实现了基于枚举的钢板验证逻辑，参考 `sdk/tube_comp.go` 的设计模式，提供更清晰、更易维护的验证方案。

## 设计特点

### 1. 枚举方案
预定义了所有13种可能的钢板类型：
- 12种普通连续：A-2, 2-3, 3-4, ..., Q-K
- 1种循环：K-A

### 2. 统一验证流程
```
输入6张牌 
  → 排序并统计 
  → 枚举检查所有可能的组合 
  → 选择最大的组合 
  → 规范化并返回
```

### 3. 比较键值映射

| 钢板类型 | pair表示 | 比较键值 | 说明 |
|---------|---------|---------|------|
| A-2     | [1,2]   | 1       | 最小 |
| 2-3     | [2,3]   | 2       |      |
| ...     | ...     | ...     |      |
| Q-K     | [12,13] | 12      |      |
| K-A     | [13,1]  | 13      | 最大 |

## 核心函数

### plateSatisfyNew
```go
func plateSatisfyNew(cards []*Card) (bool, []*Card, int)
```

**功能**：验证6张牌是否构成有效钢板

**参数**：
- `cards`: 待验证的6张牌

**返回值**：
- `bool`: 是否为合法钢板
- `[]*Card`: 规范化后的牌组（按pair顺序排列）
- `int`: 比较键值（1-13），无效时为0

**验证步骤**：
1. 检查输入有效性（牌数=6）
2. 连续性排序
3. 检查王牌（不允许）
4. 统计万能牌数量
5. 统计非万能牌的数字分布
6. 健壮性检查（任何数字不超过3张）
7. 枚举所有可能的钢板组合
8. 选择最大的组合
9. 规范化卡牌顺序
10. 返回结果

### 辅助函数

#### getAllConsecutivePairs
返回所有13种钢板枚举

#### canFormPlateWithWildcards
检查给定的牌+万能牌能否凑成指定钢板
- 精确匹配：需要的万能牌数 == 可用的万能牌数

#### findAllValidPlatePairs
找出所有可以凑成的钢板组合

#### getPlatePairComparisonKey
获取钢板的比较键值（直接返回pair[0]）

#### selectBestPlatePair
从多个可能的组合中选择最大的

#### normalizePlateWithPair
根据最佳组合规范化卡牌顺序
- 前3张：第一个数字的牌
- 后3张：第二个数字的牌
- 万能牌填补缺失位置

## 与旧实现对比

| 方面 | plateSatisfy（旧） | plateSatisfyNew（新） |
|------|-------------------|---------------------|
| 逻辑结构 | 多层嵌套if-else | 清晰的枚举+选择流程 |
| 循环处理 | 散落在多个分支 | 统一在枚举表中 |
| 万能牌 | 通过葫芦间接判断 | 直接计算缺失牌数 |
| 可维护性 | 低 | 高 |
| 代码行数 | ~120行 | ~80行 |

## 使用方法

### 当前状态
新函数 `plateSatisfyNew` 已实现，但**尚未集成**到 `comp.go` 中。

### 后续集成步骤

1. 在 `comp.go` 中修改 `NewPlate` 函数：
```go
func NewPlate(cards []*Card) *Plate {
    valid := false
    var sortedCards []*Card
    var normalizedCards []*Card
    var comparisonKey int  // 新增

    if len(cards) == 6 {
        // 使用新的枚举验证逻辑
        var ok bool
        ok, sortedCards, comparisonKey = plateSatisfyNew(cards)  // 替换
        valid = ok
        
        if valid {
            normalizedCards = sortedCards  // 新逻辑已规范化
        }
    } else {
        sortedCards = sortCards(cards)
    }

    return &Plate{
        BaseComp: BaseComp{
            Cards:           cards,
            NormalizedCards: normalizedCards,
            Valid:           valid,
            Type:            TypePlate,
        },
        ComparisonKey: comparisonKey,  // 可选：添加到结构体
    }
}
```

2. 可选：在 `Plate` 结构体中添加 `ComparisonKey` 字段（参考 `Tube`）：
```go
type Plate struct {
    BaseComp
    ComparisonKey int
}
```

3. 更新 `getPlateComparisonKey` 函数使用预计算的键值（参考 `Tube.GreaterThan`）

4. 标记旧函数为 Deprecated：
```go
// plateSatisfy 旧版钢板验证逻辑
// Deprecated: 使用 plateSatisfyNew 替代
func plateSatisfy(cards []*Card) (bool, []*Card) {
    // ...
}
```

## 测试验证

运行演示程序：
```bash
go run demo_plate.go
```

输出示例：
```
测试基本钢板识别:
   ✓ A-2钢板（最小）: 有效
   ✓ 5-6钢板: 有效
   ✓ K-A钢板（最大）: 有效
   ✓ Q-K钢板: 有效

测试循环情况:
   A-2钢板: RawNumbers: 1 1 1 2 2 2 → 有效
   K-A钢板: RawNumbers: 13 13 13 1 1 1 → 有效

非法情况验证:
   ✓ 3-5钢板（不连续）: 无效
   ✓ 5-6钢板（牌数不够）: 无效
```

## 文件结构

```
sdk/
├── comp.go              # 主牌型文件（未修改）
├── comp_util.go         # 工具函数
├── tube_comp.go         # 钢管枚举验证（参考）
├── plate_comp.go        # 钢板枚举验证（新增）✨
└── plate_comp_test.go   # 单元测试（新增）
```

## 优势

1. **可读性强**：枚举方案一目了然
2. **易于扩展**：新增规则只需修改枚举表
3. **统一风格**：与 `tube_comp.go` 保持一致
4. **循环自然**：A-2 和 K-A 作为普通枚举项
5. **精确匹配**：万能牌数量精确验证，避免多余或不足

## 注意事项

- 新实现独立于现有代码，不影响任何功能
- 后续可以由用户自行决定何时集成到 `comp.go`
- 建议先运行完整测试套件验证一致性
- 可以标记旧函数为 Deprecated 保持向后兼容性
