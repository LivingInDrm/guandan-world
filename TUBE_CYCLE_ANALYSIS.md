# 钢管（Tube）A-2 和 K-A 循环处理分析

## 总结

### 测试结果

| 场景 | 是否识别为合法 | 处理方式 |
|------|--------------|---------|
| **无变化牌** | | |
| QQ,KK,AA | ❌ 否 | 缺少特殊处理 |
| **1个变化牌** | | |
| AA + 22 + 3 + wild | ✅ 是 | 通用模式 TUBE_PATTERN_0112 |
| QQ + KK + A + wild | ❌ 否 | 缺少 K-A 循环检测 |
| **2个变化牌** | | |
| AA + KK + 2×wild | ✅ 是 | **特殊处理** (comp.go:1637-1651) |
| AA + 22 + 2×wild | ✅ 是 | 通用模式 TUBE_PATTERN_0011 |

---

## 详细分析

### 1. 无变化牌情况

#### QQ,KK,AA 钢管

**输入**: QQ, KK, AA（6张）

**排序后**（sortCardsForConsecutive按RawNumber排序）:
```
A(1), A(1), Q(12), Q(12), K(13), K(13)
```

**验证过程**:
```go
wildcardCount == 0
uniqueNumbers = {1, 12, 13}  // 3个唯一数字 ✓

cardNumbers = [1, 1, 12, 12, 13, 13]
computeRelativeDiffs(cardNumbers, 6):
  基准值 = 1
  结果 = [0, 0, 11, 11, 12, 12]

TUBE_PATTERN_TRIPLET = [0, 0, 1, 1, 2, 2]
[0,0,11,11,12,12] ≠ [0,0,1,1,2,2]  ❌ 不匹配

返回 false
```

**问题**: 缺少对 {1, 12, 13} 唯一数字组合的特殊检测

---

### 2. 一个变化牌情况

#### AA + 22 + 3 + 变化牌（✅ 成功）

**输入**: AA, 22, 3, 变化牌（6张）

**排序后**:
```
A(1), A(1), 2(2), 2(2), 3(3), wild
```

**验证过程**:
```go
wildcardCount == 1
cardNumbers = [1, 1, 2, 2, 3, -1]

firstFive = computeRelativeDiffs([1,1,2,2,3], 5)
  基准值 = 1
  firstFive = [0, 0, 1, 1, 2]

TUBE_PATTERN_0112 = [0, 0, 1, 1, 2]
匹配！✅

返回 true (comp.go:1611-1612)
```

**成功原因**: A-2-3 是自然连续（RawNumber: 1→2→3），通过通用模式识别

---

#### QQ + KK + A + 变化牌（❌ 失败）

**输入**: QQ, KK, A, 变化牌（6张）

**排序后**:
```
A(1), Q(12), Q(12), K(13), K(13), wild
```

**验证过程**:
```go
wildcardCount == 1
cardNumbers = [1, 12, 12, 13, 13, -1]

firstFive = computeRelativeDiffs([1,12,12,13,13], 5)
  基准值 = 1
  firstFive = [0, 11, 11, 12, 12]

TUBE_PATTERN_0112 = [0, 0, 1, 1, 2]  ❌ 不匹配
TUBE_PATTERN_0122 = [0, 0, 1, 2, 2]  ❌ 不匹配
TUBE_PATTERN_1122 = [0, 1, 1, 2, 2]  ❌ 不匹配

返回 false
```

**失败原因**: A的RawNumber=1破坏了Q-K-A的连续性，没有特殊检测

---

### 3. 两个变化牌情况

#### AA + KK + 2个变化牌（✅ 成功）

**输入**: AA, KK, 2个变化牌（6张）

**排序后**:
```
A(1), A(1), K(13), K(13), wild, wild
```

**验证过程**:
```go
wildcardCount == 2
cardNumbers = [1, 1, 13, 13, -1, -1]

// ✅ A-K循环特殊处理 (comp.go:1637-1651)
actualNumbers = [1, 1, 13, 13]  // 过滤掉变化牌
uniqueActual = removeDuplicates(actualNumbers)
// uniqueActual = [1, 13]

if len(uniqueActual) == 2 &&
    ((uniqueActual[0] == 1 && uniqueActual[1] == 13) ||
     (uniqueActual[0] == 13 && uniqueActual[1] == 1)) {
    // A-K循环钢管：AA + KK + 两个变化牌
    return true, sortedCards  ✅
}
```

**成功原因**: 有专门的 A-K 循环特殊处理代码！

---

#### AA + 22 + 2个变化牌（✅ 成功）

**输入**: AA, 22, 2个变化牌（6张）

**排序后**:
```
A(1), A(1), 2(2), 2(2), wild, wild
```

**验证过程**:
```go
wildcardCount == 2
cardNumbers = [1, 1, 2, 2, -1, -1]

firstFour = computeRelativeDiffs([1,1,2,2], 4)
  基准值 = 1
  firstFour = [0, 0, 1, 1]

// 匹配 TUBE_PATTERN_0011 (comp.go:1655-1664)
TUBE_PATTERN_0011 = [0, 0, 1, 1]
匹配！✅

if sortedCards[3].Number < 14 {  // 第4张牌是2，Number=2 < 14
    return true, sortedCards  ✅
}
```

**成功原因**: 通过通用模式 TUBE_PATTERN_0011 识别（不是特殊 A-2 处理）

---

## 核心问题

### RawNumber 的双重性

```go
// card.go:76-83
if number == 14 {  // Ace
    card.RawNumber = 1  // A在连续性判断中被当作1
}
```

**设计意图**: 让 A 既能在低位（A-2-3），也能在高位（Q-K-A）

**实际问题**: 
- ✅ A-2-3 可以通过，因为 RawNumber 1→2→3 自然连续
- ❌ Q-K-A 无法通过，因为排序后变成 A(1)...Q(12)-K(13)，不连续

---

## 钢板（Plate）的完整解决方案

钢板在所有情况下都正确处理了 A-2 和 K-A 循环：

### 无变化牌（comp.go:1367-1380）

```go
card1Num := triple1.Cards[0].RawNumber
card2Num := triple2.Cards[0].RawNumber

// ✅ A-2 循环检测
if card1Num == 1 && card2Num == 2 {
    return true, ...
}

// ✅ K-A 循环检测（排序后A在前）
if card1Num == 1 && card2Num == 13 {
    return true, ...
}

// ✅ K-A 循环检测（理论情况）
if card1Num == 13 && card2Num == 1 {
    return true, ...
}
```

### 有变化牌（comp.go:1432-1447）

```go
// ✅ K-A 循环检测
if tripleNum == 13 && pairNum == 1 {
    return true, ...
}
if tripleNum == 1 && pairNum == 13 {
    return true, ...
}
```

---

## 钢管缺少的处理

### 无变化牌（❌ 缺少）

需要添加类似钢板的特殊检测：

```go
if wildcardCount == 0 {
    uniqueNumbers := make(map[int]bool)
    for _, num := range cardNumbers {
        uniqueNumbers[num] = true
    }
    
    if len(uniqueNumbers) == 3 {
        // ✅ 新增：Q-K-A 特殊检测
        uniqueSlice := mapToSortedSlice(uniqueNumbers)
        if len(uniqueSlice) == 3 && 
           uniqueSlice[0] == 1 &&   // A
           uniqueSlice[1] == 12 &&  // Q
           uniqueSlice[2] == 13 {   // K
            // 验证每个数字出现2次（对子）
            return true, sortedCards
        }
        
        // 原有的通用模式匹配
        temp := computeRelativeDiffs(cardNumbers, 6)
        if matchesPattern(temp, TUBE_PATTERN_TRIPLET) {
            return true, sortedCards
        }
    }
}
```

### 一个变化牌（❌ 缺少）

需要添加 Q-K-A 循环检测：

```go
if wildcardCount == 1 {
    // 提取非变化牌的数字
    actualNumbers := filterNonWildcards(cardNumbers)
    
    // ✅ 新增：Q-K-A with 1 wildcard
    if len(actualNumbers) == 5 {
        uniqueActual := removeDuplicates(actualNumbers)
        // uniqueActual = [1, 12, 13] 表示 A, Q, K
        if len(uniqueActual) == 3 &&
           containsAll(uniqueActual, [1, 12, 13]) {
            // 重新排序：Q-Q, K-K, A-wild
            return true, reorderAsQKA(sortedCards)
        }
    }
    
    // 原有的通用模式匹配
    ...
}
```

---

## 设计原则总结

| 原则 | 说明 |
|------|------|
| **RawNumber双重性** | Number=14（牌力），RawNumber=1（连续性） |
| **排序优先** | 先按RawNumber排序，再模式识别 |
| **特殊循环** | A-2 和 K-A 需要特殊检测 |
| **双向检测** | 检查 1→13 和 13→1 两个方向 |
| **一致性** | 钢板、顺子、钢管应有相同的循环处理 |

---

## 现状对比表

| 牌型 | A-2循环(0变化牌) | K-A循环(0变化牌) | K-A循环(1变化牌) | K-A循环(2变化牌) |
|------|----------------|----------------|----------------|----------------|
| 顺子 | ❌ | ✅ 10-J-Q-K-A | ✅ | ✅ |
| 钢板 | ✅ | ✅ | ✅ | ✅ |
| 钢管 | ❌ | ❌ **缺少** | ❌ **缺少** | ✅ |

---

## 修复总结

钢管（Tube）需要补充：

1. **无变化牌**: 添加 {1, 12, 13} 特殊检测
2. **一个变化牌**: 添加 Q-K-A 循环检测和重排序
3. **保持一致性**: 与钢板和顺子的处理逻辑对齐
