# 前端访问数据与Proto定义不一致问题扫描报告

## 问题总结

基于 proto 定义扫描前端代码，发现 **15 处** 字段访问错误。

---

## 核心问题说明

### 数据结构理解错误

**后端发送的数据：**
```json
{
  "type": "game_event",
  "data": {
    "event_type": "EVENT_TYPE_TRIBUTE_STARTED",  // 字符串
    "event_data": {  // ← 这是整个 GameEvent proto
      "matchId": "...",
      "dealIndex": 0,
      "type": 20,
      "tributeStarted": {  // ← 这才是 TributeStartedPayload
        "tributeType": 1,
        "givers": [2, 3],
        "receivers": [0, 1]
      }
    }
  }
}
```

**前端当前错误理解：**
```typescript
// GamePage.tsx:147-149
const { event_type, event_data, player_seat, ...restData } = message.data || {};
const payload = event_data || restData;  // payload 是整个 GameEvent

// ❌ 错误：直接访问 payload.xxx，以为 payload 就是具体的 XxxPayload
payload.tribute_info  // 实际应该是 payload.tributeStarted.xxx
```

**正确的访问方式：**
```typescript
// payload 是 GameEvent，需要先访问 oneof 字段
payload.tributeStarted.tributeType  // ✅
payload.dealStarted.dealLevel      // ✅
payload.trickEnded.trickWinner     // ✅
```

---

## 详细问题清单

### 文件：`frontend/src/components/game/GamePage.tsx`

---

#### ❌ 问题 1：Line 159 - 事件类型名称错误

**当前代码：**
```typescript
case 'tribute_started':
```

**问题：** 后端发送的 `event_type` 是枚举的 `.String()` 值：`"EVENT_TYPE_TRIBUTE_STARTED"`

**正确写法：**
```typescript
case 'EVENT_TYPE_TRIBUTE_STARTED':
```

**影响：** 所有 case 分支的事件类型名称都需要修正。

---

#### ❌ 问题 2：Line 162 - TributeStartedPayload 访问错误

**当前代码：**
```typescript
case 'tribute_started':
  setTributeInfo(payload.tribute_info);  // ❌ 
```

**Proto 定义：**
```protobuf
message TributeStartedPayload {
  TributeType tribute_type = 1;
  repeated int32 givers = 2;
  repeated int32 receivers = 3;
}
```

**问题分析：**
1. `payload` 是 `GameEvent`，不是 `TributeStartedPayload`
2. `TributeStartedPayload` 没有 `tribute_info` 字段
3. 进贡信息应该从 `player_view` 消息的 `TributeView` 获取

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_TRIBUTE_STARTED':
  setCurrentPhase(GamePageState.TRIBUTE_PHASE);
  // tributeInfo 从后续的 player_view 消息中获取，不从事件中获取
  // 如果需要，可以访问 payload.tributeStarted.tributeType 等
```

**正确写法（当前 snake_case）：**
```typescript
case 'EVENT_TYPE_TRIBUTE_STARTED':
  setCurrentPhase(GamePageState.TRIBUTE_PHASE);
  // 如需访问：payload.tribute_started.tribute_type
```

---

#### ❌ 问题 3：Line 204 - MatchEndedPayload 访问错误

**当前代码：**
```typescript
case 'match_ended':
  setMatchResult(payload.match_result || payload);  // ❌
```

**Proto 定义：**
```protobuf
message MatchEndedPayload {
  int32 winner = 1;
  repeated int32 final_levels = 2;
  int64 duration_ms = 3;
  int32 total_deals = 4;
}
```

**问题：** `payload` 是 `GameEvent`，`payload.match_result` 不存在

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_MATCH_ENDED':
  setMatchResult(payload.matchEnded);  // matchEnded 就是 MatchEndedPayload
```

**正确写法（当前 snake_case）：**
```typescript
case 'EVENT_TYPE_MATCH_ENDED':
  setMatchResult(payload.match_ended);
```

---

#### ❌ 问题 4：Line 207-215 - TrickStartedPayload 访问错误

**当前代码：**
```typescript
case 'trick_started':
  console.log('🎲 Trick Started:', payload.trick);  // ❌
```

**Proto 定义：**
```protobuf
message TrickStartedPayload {
  int32 leader = 1;
  bool is_first_trick = 2;
  repeated int32 remaining_players = 3;
}
```

**问题：** `TrickStartedPayload` 没有 `trick` 字段

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_TRICK_STARTED':
  console.log('🎲 Trick Started:', {
    leader: payload.trickStarted.leader,
    isFirstTrick: payload.trickStarted.isFirstTrick,
    remainingPlayers: payload.trickStarted.remainingPlayers
  });
```

---

#### ❌ 问题 5：Line 223-224 - TrickEndedPayload 访问错误

**当前代码：**
```typescript
case 'trick_ended':
  console.log('✅ Trick Ended:', {
    winner: payload.winner,           // ❌ 不存在
    next_leader: payload.next_leader  // ❌ 不存在
  });
```

**Proto 定义：**
```protobuf
message TrickEndedPayload {
  int32 trick_winner = 1;
}
```

**问题：** 
1. `payload` 是 `GameEvent`，需要访问 `payload.trick_ended`
2. 字段名是 `trick_winner`，不是 `winner`
3. 没有 `next_leader` 字段（`trick_winner` 就是下一轮的 leader）

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_TRICK_ENDED':
  console.log('✅ Trick Ended:', {
    winner: payload.trickEnded.trickWinner
  });
```

**正确写法（当前 snake_case）：**
```typescript
case 'EVENT_TYPE_TRICK_ENDED':
  console.log('✅ Trick Ended:', {
    winner: payload.trick_ended.trick_winner
  });
```

---

#### ❌ 问题 6：Line 228-232 - DealStartedPayload 访问错误

**当前代码：**
```typescript
case 'deal_started':
  console.log('🎴 Deal Started:', {
    level: payload.deal_level,        // ❌
    team0_level: payload.team0_level, // ❌
    team1_level: payload.team1_level  // ❌
  });
```

**Proto 定义：**
```protobuf
message DealStartedPayload {
  int32 deal_level = 1;
  repeated int32 team_levels = 2;  // [team0_level, team1_level]
}
```

**问题：**
1. `payload` 是 `GameEvent`，需要访问 `payload.deal_started`
2. `team_levels` 是数组，不是独立的 `team0_level` 和 `team1_level`

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_DEAL_STARTED':
  console.log('🎴 Deal Started:', {
    level: payload.dealStarted.dealLevel,
    team0_level: payload.dealStarted.teamLevels[0],
    team1_level: payload.dealStarted.teamLevels[1]
  });
```

**正确写法（当前 snake_case）：**
```typescript
case 'EVENT_TYPE_DEAL_STARTED':
  console.log('🎴 Deal Started:', {
    level: payload.deal_started.deal_level,
    team0_level: payload.deal_started.team_levels[0],
    team1_level: payload.deal_started.team_levels[1]
  });
```

---

#### ❌ 问题 7：Line 234-239 - 未定义的事件类型

**当前代码：**
```typescript
case 'tribute_rules_set':      // ❌ Proto 中不存在
case 'tribute_immunity':       // ❌ 应该是 EVENT_TYPE_TRIBUTE_EXEMPTED
```

**Proto 中的事件类型：**
- `EVENT_TYPE_TRIBUTE_EXEMPTED`（免贡）

**问题：** `tribute_rules_set` 在 proto 中不存在

**正确写法：**
```typescript
case 'EVENT_TYPE_TRIBUTE_EXEMPTED':  // 免贡事件
  console.log('🛡️ Tribute Exempted:', payload.tributeExempted);
```

---

#### ❌ 问题 8：Line 240-246 - TributeCardSubmittedPayload 访问错误

**当前代码：**
```typescript
case 'tribute_given':  // ❌ 应该是 EVENT_TYPE_TRIBUTE_CARD_SUBMITTED
  console.log('⬆️ Tribute Given:', {
    giver: payload.giver,      // ❌ 不存在
    receiver: payload.receiver, // ❌ 不存在
    card: payload.card          // ❌ 应该是 submitted_card
  });
```

**Proto 定义：**
```protobuf
message TributeCardSubmittedPayload {
  guandan.common.Card submitted_card = 1;
}
```

**问题：**
1. 事件类型名错误
2. `giver` 和 `receiver` 信息在 `GameEvent.actor_seat`，不在 Payload 中
3. 字段名是 `submitted_card`，不是 `card`

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_TRIBUTE_CARD_SUBMITTED':
  console.log('⬆️ Tribute Submitted:', {
    actor: player_seat,  // 从 message.data.player_seat 获取
    card: payload.tributeCardSubmitted.submittedCard
  });
```

**正确写法（当前 snake_case）：**
```typescript
case 'EVENT_TYPE_TRIBUTE_CARD_SUBMITTED':
  console.log('⬆️ Tribute Submitted:', {
    actor: player_seat,
    card: payload.tribute_card_submitted.submitted_card
  });
```

---

#### ❌ 问题 9：Line 247-252 - 未定义的事件类型

**当前代码：**
```typescript
case 'tribute_selected':   // ❌ 应该是 EVENT_TYPE_TRIBUTE_CARD_SELECTED
case 'return_tribute':     // ❌ 应该是 EVENT_TYPE_TRIBUTE_CARD_RETURNED
```

**Proto 中的事件类型：**
- `EVENT_TYPE_TRIBUTE_CARD_SELECTED`（选贡牌）
- `EVENT_TYPE_TRIBUTE_CARD_RETURNED`（还贡）

**正确写法（改为 camelCase 后）：**
```typescript
case 'EVENT_TYPE_TRIBUTE_CARD_SELECTED':
  console.log('✅ Tribute Selected:', {
    card: payload.tributeCardSelected.selectedCard,
    isAuto: payload.tributeCardSelected.isAuto
  });
  break;

case 'EVENT_TYPE_TRIBUTE_CARD_RETURNED':
  console.log('⬇️ Return Tribute:', {
    card: payload.tributeCardReturned.returnedCard,
    targetPlayer: payload.tributeCardReturned.targetPlayer,
    isAuto: payload.tributeCardReturned.isAuto
  });
  break;
```

---

#### ❌ 问题 10：Line 174-196 - DealEndedPayload 访问（部分正确）

**当前代码：**
```typescript
const dealEndedPayload = payload as DealEndedPayload;  // ❌ payload 不是 DealEndedPayload

const dealResult: DealResultType = {
  rankings: dealEndedPayload.rankings || [],           // ✅ 字段名正确
  winning_team: dealEndedPayload.winningTeam || 0,     // ✅ 
  victory_type: dealEndedPayload.victoryType,          // ✅
  upgrades: (dealEndedPayload.levelChange || [0, 0]),  // ✅
  duration: dealEndedPayload.durationMs || 0,          // ✅
  trick_count: dealEndedPayload.trickCount || 0,       // ✅
  // ...
};
```

**问题：** `payload` 是 `GameEvent`，需要先访问 `payload.deal_ended` 或 `payload.dealEnded`

**正确写法（改为 camelCase 后）：**
```typescript
const dealEndedPayload = payload.dealEnded as DealEndedPayload;

const dealResult: DealResultType = {
  rankings: dealEndedPayload.rankings || [],
  winning_team: dealEndedPayload.winningTeam || 0,
  victory_type: dealEndedPayload.victoryType,
  upgrades: (dealEndedPayload.levelChange || [0, 0]) as [number, number],
  duration: dealEndedPayload.durationMs || 0,
  trick_count: dealEndedPayload.trickCount || 0,
  // ...
};
```

**正确写法（当前 snake_case）：**
```typescript
const dealEndedPayload = payload.deal_ended as DealEndedPayload;
```

---

## 事件类型名称映射表

| 前端当前使用（错误） | Proto 定义（正确） | 说明 |
|-------------------|------------------|-----|
| `tribute_started` | `EVENT_TYPE_TRIBUTE_STARTED` | 进贡开始 |
| `tribute_completed` | `EVENT_TYPE_TRIBUTE_COMPLETED` | 进贡完成 |
| `tribute_given` | `EVENT_TYPE_TRIBUTE_CARD_SUBMITTED` | 贡牌提交 |
| `tribute_selected` | `EVENT_TYPE_TRIBUTE_CARD_SELECTED` | 选贡牌 |
| `return_tribute` | `EVENT_TYPE_TRIBUTE_CARD_RETURNED` | 还贡 |
| `tribute_immunity` | `EVENT_TYPE_TRIBUTE_EXEMPTED` | 免贡 |
| `tribute_rules_set` | ❌ 不存在 | 应删除 |
| `deal_completed` | `EVENT_TYPE_DEAL_ENDED` | 局结束 |
| `deal_started` | `EVENT_TYPE_DEAL_STARTED` | 局开始 |
| `match_completed` | `EVENT_TYPE_MATCH_ENDED` | 比赛结束 |
| `trick_started` | `EVENT_TYPE_TRICK_STARTED` | 轮次开始 |
| `trick_ended` / `trick_completed` | `EVENT_TYPE_TRICK_ENDED` | 轮次结束 |
| `player_played` | `EVENT_TYPE_PLAYER_PLAYED` | 玩家出牌 |
| `player_passed` | `EVENT_TYPE_PLAYER_PASSED` | 玩家过牌 |

---

## GameEvent oneof payload 访问模式

### 当前 snake_case（后端配置 UseProtoNames: true）

```typescript
// payload 是整个 GameEvent
switch (event_type) {
  case 'EVENT_TYPE_TRIBUTE_STARTED':
    payload.tribute_started.tribute_type
    payload.tribute_started.givers
    payload.tribute_started.receivers
    
  case 'EVENT_TYPE_DEAL_STARTED':
    payload.deal_started.deal_level
    payload.deal_started.team_levels[0]
    
  case 'EVENT_TYPE_DEAL_ENDED':
    payload.deal_ended.rankings
    payload.deal_ended.winning_team
    payload.deal_ended.victory_type
    
  case 'EVENT_TYPE_TRICK_ENDED':
    payload.trick_ended.trick_winner
    
  case 'EVENT_TYPE_TRIBUTE_CARD_SUBMITTED':
    payload.tribute_card_submitted.submitted_card
}
```

### 改为 camelCase（后端配置 UseProtoNames: false）

```typescript
// payload 是整个 GameEvent
switch (event_type) {
  case 'EVENT_TYPE_TRIBUTE_STARTED':
    payload.tributeStarted.tributeType
    payload.tributeStarted.givers
    payload.tributeStarted.receivers
    
  case 'EVENT_TYPE_DEAL_STARTED':
    payload.dealStarted.dealLevel
    payload.dealStarted.teamLevels[0]
    
  case 'EVENT_TYPE_DEAL_ENDED':
    payload.dealEnded.rankings
    payload.dealEnded.winningTeam
    payload.dealEnded.victoryType
    
  case 'EVENT_TYPE_TRICK_ENDED':
    payload.trickEnded.trickWinner
    
  case 'EVENT_TYPE_TRIBUTE_CARD_SUBMITTED':
    payload.tributeCardSubmitted.submittedCard
}
```

---

## 修复优先级

### 高优先级（导致功能失效）

1. ✅ **Issue 2 (Line 162)**: `payload.tribute_info` 访问导致进贡阶段无法正确设置信息
2. ✅ **Issue 10 (Line 174)**: `payload as DealEndedPayload` 导致局结算数据全部 undefined

### 中优先级（日志输出错误，不影响功能）

3. **Issue 4-9 (Line 207-252)**: 所有 console.log 中的字段访问错误

### 低优先级（可选）

4. **事件类型名称**: 从简化名改为 Proto 标准名（建议保持一致性）

---

## 建议的修复方案

### 方案 1：最小改动（推荐）

1. 保持前端使用简化的事件类型名（如 `tribute_started`）
2. 修正所有 payload 访问，统一加上 oneof 字段：
   ```typescript
   payload.tribute_started.xxx  // 而不是 payload.xxx
   ```
3. 后端保持 `UseProtoNames: true`

### 方案 2：完全统一（推荐长期）

1. 后端改为 `UseProtoNames: false`
2. 前端改用 Proto 标准事件类型名（`EVENT_TYPE_XXX`）
3. 统一使用 camelCase 字段名
4. 修正所有 payload 访问模式

---

## 额外发现

### ✅ 正确的代码

**`frontend/src/utils/converters.ts`**
- 所有 proto 字段访问都正确使用 camelCase
- `trick_id: proto.trickIndex?.toString()` ✅
- 转换器工作正常

**`frontend/src/components/game/GamePage.tsx:214`**
```typescript
case 'player_played':
  console.log('🃏 Player Played:', {
    player: player_seat,
    cards: payload.cards  // ✅ 这个需要改为 payload.player_played.cards
  });
```

---

## 总结

- **错误数量**: 15 处
- **主要问题**: 混淆了 `GameEvent` 和具体的 `XxxPayload`
- **根本原因**: 数据结构理解错误
- **影响范围**: 游戏事件处理逻辑、日志输出
- **修复难度**: 中等（需要系统性修改所有事件处理分支）
