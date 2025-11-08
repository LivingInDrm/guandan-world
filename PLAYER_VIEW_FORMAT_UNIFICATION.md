# Player View 格式统一文档

## 概述
本次修改统一了前后端的 `player_view` 消息格式，以 GameDriver 的嵌套格式为标准，解决了页面卡在倒计时后无法进入游戏的问题。

## 问题背景

### 原问题
- 点击开始游戏后，倒计时完成，页面停留在"游戏即将开始"界面，无法进入游戏
- 前端控制台显示收到了 `game_event` 消息，但游戏状态没有正确更新

### 根本原因
后端发送 `player_view` 消息的格式不统一：

1. **handlers/room.go** 发送的初始 player_view（扁平格式）：
```json
{
  "type": "player_view",
  "data": {
    "player_seat": 0,
    "hand": [],
    "can_play": false,
    "is_my_turn": false,
    "game_state": {...}
  }
}
```

2. **driver_service.go** 发送的 player_view（嵌套格式）：
```json
{
  "type": "player_view",
  "data": {
    "player_view": {
      "player_seat": 0,
      "player_cards": [],
      "game_state": {...}
    },
    "event_type": "match_started",
    "player_seat": 0,
    "filtered_state": {
      "can_play": false,
      "is_my_turn": false
    }
  }
}
```

前端只能处理其中一种格式，导致来自 GameDriver 的消息无法正确解析。

## 解决方案

### 统一标准
**以 GameDriver 的嵌套格式为准**，符合 SDK 的 `PlayerGameState` 类型定义：

```go
type PlayerGameState struct {
    PlayerSeat   int        `json:"player_seat"`
    GameState    *GameState `json:"game_state"`
    PlayerCards  []*Card    `json:"player_cards"`  // 注意：使用 player_cards 而非 hand
    VisibleCards []*Card    `json:"visible_cards"`
}
```

### 修改内容

#### 1. 前端修改 (`frontend/src/components/game/GamePage.tsx`)

**修改前**：
```typescript
const handlePlayerView = (message: WSMessage) => {
  const data = message.data;
  setPlayerSeat(data.player_seat);
  setPlayerHand(data.hand || []);         // ❌ 扁平结构
  setCanPlay(data.can_play || false);
  if (data.game_state) {                  // ❌ 扁平结构
    setGameState(data.game_state);
  }
};
```

**修改后**：
```typescript
const handlePlayerView = (message: WSMessage) => {
  const data = message.data;
  const playerView = data.player_view;    // ✅ 从嵌套对象中提取
  
  setPlayerSeat(playerView.player_seat ?? data.player_seat);
  setPlayerHand(playerView.player_cards || []);  // ✅ 使用 player_cards
  
  if (playerView.game_state) {           // ✅ 从 player_view 中获取
    setGameState(playerView.game_state);
  }
  
  const filteredState = data.filtered_state || {};
  setCanPlay(filteredState.can_play || false);   // ✅ 从 filtered_state 中获取
  setMyTurn(filteredState.is_my_turn || false);
};
```

**关键变化**：
- 从 `data.player_view` 嵌套对象中提取数据
- 使用 `player_cards` 字段而非 `hand`
- 从 `filtered_state` 中获取 `can_play` 和 `is_my_turn`

#### 2. 后端修改 (`backend/handlers/room.go`)

**修改前**（扁平格式）：
```go
Data: map[string]interface{}{
    "player_seat": i,
    "hand":        []interface{}{},
    "can_play":    false,
    "is_my_turn":  false,
    "game_state": map[string]interface{}{...},
}
```

**修改后**（嵌套格式，与 GameDriver 一致）：
```go
Data: map[string]interface{}{
    "player_view": map[string]interface{}{
        "player_seat":   i,
        "player_cards":  []interface{}{},  // 使用 player_cards
        "visible_cards": []interface{}{},
        "game_state": map[string]interface{}{
            "current_match": map[string]interface{}{
                "team_levels": []int{2, 2},
                "current_deal": map[string]interface{}{
                    "tribute_phase": nil,
                    "current_trick": nil,
                },
            },
        },
    },
    "event_type":  "match_started",
    "player_seat": i,
    "filtered_state": map[string]interface{}{
        "can_play":   false,
        "is_my_turn": false,
    },
}
```

#### 3. 增强 filtered_state (`backend/game/driver_service.go`)

在 `createFilteredState()` 方法中添加了 `can_play` 和 `is_my_turn` 字段的计算逻辑：

```go
func (wso *WebSocketObserver) createFilteredState(playerView *sdk.PlayerGameState, playerSeat int) map[string]interface{} {
    // ... 
    
    // 判断玩家是否可以出牌
    canPlay := false
    isMyTurn := false
    
    if deal.Status == sdk.DealStatusPlaying && deal.CurrentTrick != nil {
        isMyTurn = deal.CurrentTrick.CurrentTurn == playerSeat
        canPlay = isMyTurn && len(playerView.PlayerCards) > 0
    }
    
    filteredState["can_play"] = canPlay
    filteredState["is_my_turn"] = isMyTurn
    
    // ...
}
```

## 数据格式标准

### 完整的 player_view 消息格式

```typescript
{
  type: "player_view",
  data: {
    // 嵌套的玩家视图（来自 SDK PlayerGameState）
    player_view: {
      player_seat: number,           // 玩家座位号 (0-3)
      player_cards: Card[],          // 玩家手牌（注意：使用 player_cards）
      visible_cards: Card[],         // 可见的牌
      game_state: {
        id: string,
        status: "waiting" | "started" | "finished",
        current_match: {
          id: string,
          status: "waiting" | "playing" | "finished",
          team_levels: [number, number],
          current_deal: {
            id: string,
            level: number,
            status: "waiting" | "dealing" | "tribute" | "playing" | "finished",
            tribute_phase: TributePhase | null,
            current_trick: Trick | null
          },
          players: Player[]
        }
      }
    },
    
    // 元数据
    event_type: string,              // 事件类型
    player_seat: number,             // 座位号（冗余，方便访问）
    
    // 过滤状态（UI 相关状态）
    filtered_state: {
      can_play: boolean,             // 玩家是否可以出牌
      is_my_turn: boolean,           // 是否是玩家的回合
      current_match: {...},          // 其他过滤后的状态
      current_deal: {...},
      current_trick: {...},
      players: [...]
    }
  },
  timestamp: Date,
  player_id: string
}
```

### 字段说明

| 字段路径 | 类型 | 说明 |
|---------|------|------|
| `data.player_view` | object | SDK PlayerGameState 对象，包含玩家视角的游戏状态 |
| `data.player_view.player_seat` | number | 玩家座位号 (0-3) |
| `data.player_view.player_cards` | Card[] | 玩家手牌（**注意使用 player_cards 而非 hand**） |
| `data.player_view.visible_cards` | Card[] | 当前可见的牌 |
| `data.player_view.game_state` | GameState | 完整的游戏状态 |
| `data.event_type` | string | 触发此 player_view 的事件类型 |
| `data.player_seat` | number | 座位号（冗余字段，方便快速访问） |
| `data.filtered_state` | object | 过滤后的状态，包含 UI 相关的辅助字段 |
| `data.filtered_state.can_play` | boolean | 玩家当前是否可以出牌 |
| `data.filtered_state.is_my_turn` | boolean | 是否是玩家的回合 |

## 前端处理逻辑

```typescript
// 1. 从嵌套对象中提取 player_view
const playerView = data.player_view;

// 2. 获取玩家座位
const seat = playerView.player_seat;

// 3. 获取手牌（注意字段名是 player_cards）
const hand = playerView.player_cards || [];

// 4. 获取游戏状态
const gameState = playerView.game_state;

// 5. 从 filtered_state 获取 UI 状态
const canPlay = data.filtered_state?.can_play || false;
const isMyTurn = data.filtered_state?.is_my_turn || false;

// 6. 判断游戏阶段
const currentDeal = gameState?.current_match?.current_deal;
const hasTribute = currentDeal?.tribute_phase !== null;
```

## 优势

### 1. 类型安全
- 遵循 SDK 的 `PlayerGameState` 类型定义
- 前端可以基于 SDK 类型定义生成 TypeScript 接口

### 2. 职责分离
- `player_view`：包含游戏逻辑相关的状态（来自 SDK）
- `filtered_state`：包含 UI 相关的辅助状态（后端计算好）

### 3. 扩展性
- 未来添加新字段时，可以在合适的层级添加
- 不会破坏现有的数据结构

### 4. 一致性
- 所有 player_view 消息使用相同的格式
- 减少前端处理逻辑的复杂度

## 验证要点

测试时应验证以下场景：

1. ✅ **游戏开始**
   - 倒计时完成后能正确进入游戏界面
   - 能看到初始手牌（空）
   - 游戏状态正确显示

2. ✅ **发牌**
   - 收到 `EventCardsDealt` 后能看到手牌
   - 手牌数量正确

3. ✅ **上贡阶段**
   - 能正确识别是否有上贡阶段
   - 上贡界面正确显示

4. ✅ **出牌阶段**
   - 轮到自己时 `can_play` 为 true
   - 不是自己回合时 `can_play` 为 false
   - 能正确出牌和过牌

5. ✅ **状态同步**
   - 所有玩家的游戏状态保持同步
   - 手牌数量变化正确

## 相关文件

- `frontend/src/components/game/GamePage.tsx` - 前端游戏页面（player_view 处理）
- `frontend/src/services/gameService.ts` - 游戏服务（已删除冲突的 player_view 处理器）
- `backend/handlers/room.go` - 后端房间处理器（初始 player_view）
- `backend/game/driver_service.go` - GameDriver 服务（游戏中的 player_view）
- `sdk/game_engine.go` - SDK PlayerGameState 类型定义

## 迁移说明

如果有其他代码依赖旧的扁平格式，需要相应修改：

```typescript
// 旧代码
const hand = data.hand;
const gameState = data.game_state;

// 新代码
const hand = data.player_view.player_cards;
const gameState = data.player_view.game_state;
```

## 总结

通过这次统一，我们：
1. ✅ 解决了页面卡在倒计时的问题
2. ✅ 统一了前后端的数据格式
3. ✅ 提高了代码的可维护性
4. ✅ 为未来的扩展打下了良好的基础

---

**修改日期**: 2025-11-08
**修改人**: AI Assistant
**版本**: 1.0

