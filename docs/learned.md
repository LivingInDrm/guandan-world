# 调试经验总结

## 1. Go 并发死锁问题

### 问题
同步事件处理 + 锁重入导致死锁：
```go
func (ge *GameEngine) StartDeal() error {
    ge.mutex.Lock()           // 持有写锁
    defer ge.mutex.Unlock()
    
    ge.emitEvent(event)       // 同步调用 handler
        -> handler()          // handler 调用 GetPlayerView()
            -> GetPlayerView() tries ge.mutex.RLock()  // 死锁！
}
```

### 解决方案
异步事件处理：
```go
func (ge *GameEngine) emitEvent(event *GameEvent) {
    for _, handler := range handlers {
        h := handler
        go func() {  // 异步执行，避免锁重入
            defer recover() // panic recovery
            h(event)
        }()
    }
}
```

### 教训
- **持有锁时避免同步调用外部代码**
- **事件处理器应该异步执行**
- **添加 panic recovery 保护稳定性**

---

## 2. React useEffect 闭包陷阱

### 问题
useEffect 依赖变化导致 handler 重新创建，捕获了旧的闭包值：
```typescript
useEffect(() => {
    const handlePlayerView = () => {
        if (currentPhase === 'GAME_PREPARE') {  // 使用闭包值
            // ... 永远不会进入，因为 currentPhase 是旧值
        }
    };
    wsClient.on('player_view', handlePlayerView);
}, [currentPhase, otherDeps]);  // currentPhase 变化 → handler 重新创建 → 捕获旧值
```

### 解决方案
从 store 获取最新值，而不是使用闭包：
```typescript
const handlePlayerView = () => {
    const latestPhase = useGameStore.getState().currentPhase;  // 获取最新值
    if (latestPhase === 'GAME_PREPARE') {
        // 正确使用最新值
    }
};
```

### 教训
- **useEffect 中的 handler 会捕获创建时的闭包值**
- **需要最新值时，从 store.getState() 获取**
- **减少不必要的 useEffect 依赖项**

---

## 3. 前后端数据结构对齐

### 问题
后端发送嵌套结构，前端期望扁平结构：
```javascript
// 后端实际结构
gameState = {
    current_match: {
        team_levels: [2, 2],
        current_deal: { level: 2 }
    }
}

// 前端错误访问
gameState.team_levels        // undefined ❌
gameState.current_deal       // undefined ❌
```

### 解决方案
前端适配后端的标准结构：
```javascript
const teamLevels = gameState.current_match?.team_levels || [2, 2];
const currentDeal = gameState.current_match?.current_deal;
```

### 教训
- **类型定义要与实际数据结构一致**
- **添加运行时日志验证数据结构**
- **使用可选链 + 默认值防御性编程**

---

## 4. 状态污染问题

### 问题
多个消息源更新同一个状态，导致状态被错误覆盖：
```typescript
// player_view 正确设置 gameState
setGameState(playerView.game_state);  // ✅

// game_event 错误覆盖
setGameState(message.data);  // ❌ message.data 是事件对象，不是游戏状态
```

### 解决方案
明确状态的单一数据源：
- **gameState 只从 player_view 更新**
- **game_event 只用于事件通知，不更新状态**

### 教训
- **为每个状态定义唯一的数据源**
- **避免多个地方更新同一个状态**
- **事件通知 ≠ 状态更新**

---

## 5. 调试方法论

### 有效的调试步骤
1. **查看完整日志** - 前端 console + 后端 docker logs
2. **对比预期流程** - 应该收到什么消息 vs 实际收到什么
3. **添加调试日志** - 在关键路径打印数据结构
4. **逐层排查** - 从后端 → WebSocket → 前端 store → 组件
5. **数据结构验证** - 展开对象查看实际内容

### 关键调试技巧
```typescript
// 打印实际数据结构
console.log('数据检查:', {
    hasField: !!obj.field,
    actualKeys: Object.keys(obj),
    fullObject: obj
});
```

### 教训
- **不要假设数据结构，要验证**
- **使用渐进式调试，逐步缩小范围**
- **保留调试日志直到问题完全解决**

---

## 6. 代码防御性编程

### 最佳实践
```typescript
// 数组操作前验证
const safeArray = Array.isArray(data) ? data : [];

// 可选链 + 默认值
const level = gameState.current_match?.current_deal?.level || 2;

// 数据完整性检查
if (!currentMatch || !currentDeal) {
    return <Loading />;
}
```

### 教训
- **永远不要假设数据一定存在**
- **为所有可选字段提供默认值**
- **在渲染前验证数据完整性**

---

## 7. AI 算法兜底逻辑缺失

### 问题
AI 算法在首出时返回 nil，导致玩家无法出牌：
```
Error: invalid pass: cannot pass as trick leader - must play cards
```

**根本原因**：
1. **算法层**：`selectBestGroup` 只返回评分为正的牌组，导致所有牌组评分都不为正时返回 nil
2. **客户端层**：没有对首出时返回空结果进行安全检查，直接发送 "pass" 动作

### 原始代码问题
```go
// 算法层 - smart_algorithm.go
func (algo *SmartAutoPlayAlgorithm) selectBestFirstPlay(hand []*sdk.Card) []*sdk.Card {
    // ... 识别和评分 ...
    bestGroup := algo.selectBestGroup(groups)
    
    // 问题：如果 bestGroup 为 nil，虽然有单牌兜底，但某些极端情况下仍可能失败
    if bestGroup == nil {
        return []*sdk.Card{algo.findSmallestCard(hand)}
    }
    return bestGroup.Cards
}

// 客户端层 - ai_player_client.go
selectedCards := c.aiAlgorithm.SelectCardsToPlay(hand, trickInfo)

// 问题：没有检查 isLeader 时是否返回了空结果
if len(selectedCards) == 0 {
    action = "pass"  // ❌ 如果是首出，这会导致错误！
}
```

### 解决方案

**1. 算法层增强兜底逻辑**
```go
func (algo *SmartAutoPlayAlgorithm) selectBestFirstPlay(hand []*sdk.Card) []*sdk.Card {
    // 确保手牌不为空
    if len(hand) == 0 {
        return nil
    }
    
    // ... 识别和评分 ...
    
    // 如果找到评分为正的牌组
    if bestGroup != nil && bestGroup.Score > 0 {
        return bestGroup.Cards
    }
    
    // 如果所有牌组评分都不为正，选择评分最高的（即使是负分）
    if len(groups) > 0 {
        best := groups[0]
        for _, group := range groups[1:] {
            if group.Score > best.Score {
                best = group
            }
        }
        // 如果评分不太差，就出这个牌组
        if best.Score > -5.0 {
            return best.Cards
        }
    }
    
    // 最后的兜底：返回最小单牌
    smallestCard := algo.findSmallestCard(hand)
    if smallestCard != nil {
        return []*sdk.Card{smallestCard}
    }
    
    return nil  // 理论上不应该到这里
}
```

**2. 客户端增加安全检查**
```go
selectedCards := c.aiAlgorithm.SelectCardsToPlay(hand, trickInfo)

// 安全检查：如果是首出但算法返回空结果，强制出最小的牌
if isLeader && len(selectedCards) == 0 {
    c.log("WARNING: AI algorithm returned no cards for trick leader, forcing smallest card")
    if len(hand) > 0 {
        smallest := hand[0]
        for _, card := range hand[1:] {
            if !card.GreaterThan(smallest) {
                smallest = card
            }
        }
        selectedCards = []*sdk.Card{smallest}
    } else {
        c.log("ERROR: No cards in hand!")
        return
    }
}

if len(selectedCards) == 0 {
    action = "pass"  // ✅ 现在安全了
}
```

### 教训
1. **关键路径必须有多层兜底保护**
   - 算法层保证有合理返回值
   - 客户端层检查异常情况
   - 永远不要假设上游一定正确

2. **强制性约束要在代码中体现**
   - "首出必须出牌" 这个规则应该在代码逻辑中明确体现
   - 不能仅依赖算法"应该"返回正确结果

3. **防御性编程的层次**
   - 算法层：多重兜底逻辑（评分牌组 → 最高评分 → 最小单牌）
   - 调用层：验证结果是否符合业务规则
   - 系统层：记录异常日志便于排查

4. **评分系统的合理性**
   - 评分为负不代表不能出
   - 应该区分"相对优劣"和"能否执行"
   - 首出时即使最差的选择也要执行

---

## 8. 前后端数据格式不匹配导致显示异常

### 问题
前端手牌显示 "NaN (0)"，虽然标题显示有 27 张牌，但实际牌面无法显示。

### 根本原因
**前后端 Card 数据结构不匹配**：

**后端 SDK Card 结构** (`sdk/card.go`)：
```go
type Card struct {
    Number    int    // 牌的数字值 (1-16)
    RawNumber int    // 原始数字值
    Color     string // 花色 (Spade/Heart/Club/Diamond/Joker)
    Level     int    // 当前级别
    Name      string // 牌的名称
}
// ❌ 没有 JSON tag，序列化后字段名是 Number, Color 等
```

**前端 TypeScript 期望的结构** (`frontend/src/types/index.ts`)：
```typescript
export interface Card {
  id: string;
  suit: number;      // 0=spades, 1=hearts, 2=clubs, 3=diamonds
  rank: number;      // 2-14 (11=J, 12=Q, 13=K, 14=A), 15/16=jokers
  is_joker: boolean;
}
```

### 问题表现

1. **后端发送**：`{ Number: 2, Color: "Spade", ... }`
2. **前端接收**：尝试访问 `card.rank`，结果是 `undefined`
3. **分组失败**：
```typescript
const groupedCards = safeCards.reduce((groups, card) => {
    const rank = card.rank;  // undefined!
    // groups[undefined] = []
}
```
4. **显示 NaN**：`getRankText(undefined)` → `"NaN"`

### 解决方案

**在后端添加自定义 JSON 序列化**，无需改动前端代码：

```go
// 添加花色转数字的方法
func (c *Card) GetSuitNumber() int {
	switch c.Color {
	case "Spade":   return 0
	case "Heart":   return 1
	case "Club":    return 2
	case "Diamond": return 3
	case "Joker":   return -1
	default:        return -1
	}
}

// 自定义 JSON 序列化
func (c *Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":       c.GetID(),           // "Spade_5"
		"suit":     c.GetSuitNumber(),   // 0
		"rank":     c.Number,            // 5
		"is_joker": c.Color == "Joker",  // false
	})
}
```

### 方案优势

1. **改动集中**：只修改一个文件（`sdk/card.go`）
2. **前端零改动**：完全兼容现有前端代码
3. **符合标准**：使用国际扑克牌标准命名（suit/rank）
4. **不破坏内部逻辑**：Go 代码内部仍使用原字段名
5. **类型安全**：通过测试验证序列化正确性

### 测试验证

添加了完整的测试用例：
- `TestCardGetSuitNumber`：验证花色转数字
- `TestCardMarshalJSON`：验证单个 Card 序列化
- `TestCardSliceMarshalJSON`：验证 Card 数组序列化

所有测试通过，确保：
- ✅ 序列化输出格式正确
- ✅ 所有牌型（普通牌、大小王）都正确
- ✅ 不破坏现有功能

### 教训

1. **API 契约要明确定义**
   - 前后端要约定统一的数据格式
   - 最好有 API 文档或 OpenAPI 规范
   - 字段命名要使用行业标准

2. **Go JSON 序列化最佳实践**
   - 默认使用字段名序列化（首字母大写会导出）
   - 使用 JSON tag 控制字段名：`json:"field_name"`
   - 需要复杂转换时使用 `MarshalJSON()` 方法

3. **前端防御性编程很重要**
   - `PlayerHand` 组件已经做了防御（`Array.isArray` 检查）
   - 但字段不存在时仍会出问题
   - 可以添加数据验证层（如 zod）

4. **问题排查思路**
   - 症状：显示 "NaN"
   - 定位：前端解析到 undefined
   - 原因：后端字段名不匹配
   - 方案：后端自定义序列化

5. **改动后端 vs 前端的权衡**
   - ✅ 改后端：改动集中、符合标准、影响小
   - ❌ 改前端：改动分散、需要修改多个组件、测试成本高

---

## 总结

这次调试涉及：
- ✅ 后端并发死锁（Go 锁重入）
- ✅ 前端状态管理（React 闭包、状态污染）
- ✅ 数据结构对齐（类型不匹配）
- ✅ 防御性编程（空指针保护）
- ✅ AI 算法兜底逻辑（多层保护）
- ✅ 前后端数据格式不匹配（JSON 序列化）

核心教训：**系统性思考 + 逐层验证 + 防御性编程 + 多层兜底保护 + API 契约明确**

