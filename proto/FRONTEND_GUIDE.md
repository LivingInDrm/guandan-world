# Protobuf 前端使用指南

## 核心概念

游戏状态通过 WebSocket 推送三种消息：

| 消息类型 | 用途 | 数据来源 |
|---------|------|---------|
| `game_event` | 游戏事件通知（动画、提示） | `proto/event.proto` |
| `player_view` | 玩家视角状态（**唯一可信来源**） | `proto/view.proto` |
| `game_action` | 请求玩家操作 | - |

**重要：** 渲染界面时，始终以 `player_view` 为准，`game_event` 仅用于动画。

---

## TypeScript 类型映射

```typescript
// Card - 卡牌
interface Card {
  suit: number;        // 0=♠, 1=♥, 2=♣, 3=♦, -1=Joker
  rank: number;        // 2-14=普通牌, 15=小王, 16=大王
  deck_index: number;  // 0-107
}

// PlayerView - 玩家视角（核心）
interface PlayerView {
  match_id: string;
  deal_index: number;
  player_seat: number;           // 0-3
  player_cards: Card[];          // 手牌
  team_levels: [number, number]; // 两队等级
  deal_level: number;            // 当前局等级
  deal_status: DealStatus;       // 牌局状态
  
  // 仅 deal_status = PLAYING 时有效
  current_turn?: number;         // 当前出牌玩家
  plays?: PlayAction[];          // 本轮出牌记录
  leader?: number;               // 本轮领先玩家
}

// PlayAction - 出牌动作
interface PlayAction {
  player_seat: number;
  cards: Card[];
  comp_type: CompType;  // SINGLE, PAIR, BOMB...
  is_pass: boolean;
}
```

---

## WebSocket 消息格式

### 1. game_event - 游戏事件

```json
{
  "type": "game_event",
  "data": {
    "event_type": "player_played",
    "event_data": {
      "cards": [{"suit": 0, "rank": 3, "deck_index": 1}]
    },
    "player_seat": 0
  }
}
```

### 2. player_view - 状态更新

```json
{
  "type": "player_view",
  "data": {
    "player_seat": 0,
    "player_cards": [...],
    "deal_status": "DEAL_STATUS_PLAYING",
    "current_turn": 1,
    "plays": []
  }
}
```

### 3. game_action - 操作请求

```json
{
  "type": "game_action",
  "data": {
    "action_type": "play_decision_required",  // 或 tribute_selection_required, return_tribute_required
    "player_seat": 0,
    "timeout": 30
  }
}
```

---

## 事件类型速查

| EventType | 触发时机 | 前端动作 |
|-----------|---------|---------|
| `MATCH_STARTED` | 比赛开始 | 初始化界面 |
| `DEAL_STARTED` | 新局开始 | 更新等级显示 |
| `CARDS_DEALT` | 发牌完成 | 等待 player_view |
| `TRIBUTE_STARTED` | 进贡开始 | 显示进贡界面 |
| `TRIBUTE_COMPLETED` | 进贡结束 | 准备出牌 |
| `TRICK_STARTED` | 新轮开始 | 清空牌桌 |
| `PLAYER_PLAYED` | 玩家出牌 | 播放动画 |
| `PLAYER_PASSED` | 玩家过牌 | 显示"过" |
| `TRICK_ENDED` | 轮次结束 | 显示赢家 |
| `DEAL_ENDED` | 局结束 | 显示结算 |

**完整列表见 `proto/event.proto` 中的 `EventType` 枚举。**

---

## DealStatus 状态机

```
WAITING → DEALING → TRIBUTE → PLAYING → FINISHED
                      ↓
                 (可能跳过)
```

| 状态 | 说明 | 可用字段 |
|-----|------|---------|
| `WAITING` | 等待开始 | `team_levels`, `deal_level` |
| `TRIBUTE` | 进贡阶段 | 上述 + `player_cards` |
| `PLAYING` | 出牌阶段 | 所有字段 |

---

## 关键场景示例

### 场景 1：出牌流程

```typescript
// 1. 收到操作请求
wsClient.on('game_action', (msg) => {
  if (msg.data.action_type === 'play_decision_required') {
    enableCardSelection();
  }
});

// 2. 提交出牌
await apiClient.playCards(roomId, seat, cardIds);

// 3. 收到事件（动画）
wsClient.on('game_event', (msg) => {
  if (msg.data.event_type === 'player_played') {
    showAnimation(msg.data);
  }
});

// 4. 收到状态更新（更新UI）
wsClient.on('player_view', (msg) => {
  setPlayerCards(msg.data.player_cards);  // 手牌已减少
  setCurrentTurn(msg.data.current_turn);  // 下一个玩家
});
```

### 场景 2：双下进贡流程

```
1. DEAL_ENDED → rankings: [0, 2, 1, 3]
2. TRIBUTE_STARTED → tribute_type: DOUBLE_DOWN, givers: [1,3]
3. TRIBUTE_CARD_SUBMITTED × 2 → 败方提交贡牌
4. game_action → tribute_selection_required (rank1 选牌)
5. TRIBUTE_CARD_SELECTED × 2 → 胜方选牌完成
6. game_action → return_tribute_required
7. TRIBUTE_CARD_RETURNED × 2 → 还贡完成
8. TRIBUTE_COMPLETED
9. player_view → deal_status: PLAYING
```

---

## 实用代码片段

### 卡牌显示

```typescript
const SUIT = ['♠', '♥', '♣', '♦'];
const RANK = ['', '', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K', 'A', '小王', '大王'];

function cardStr(c: Card): string {
  return c.rank >= 15 ? RANK[c.rank] : SUIT[c.suit] + RANK[c.rank];
}
```

### 玩家队伍

```typescript
function getTeam(seat: number): number { return seat % 2; }
function getTeammate(seat: number): number { return (seat + 2) % 4; }
```

### 状态判断

```typescript
function isMyTurn(view: PlayerView, mySeat: number): boolean {
  return view.deal_status === 'DEAL_STATUS_PLAYING' && view.current_turn === mySeat;
}
```

---

## 常见问题

**Q: game_event 和 player_view 的区别？**  
A: `game_event` 用于动画，所有玩家看到相同内容；`player_view` 用于渲染界面，每个玩家看到不同内容（包含手牌）。**永远以 player_view 为准。**

**Q: 如何处理断线重连？**  
A: 重连后发送 `join_room`，等待服务端推送最新 `player_view`，根据它恢复界面。

**Q: 超时后会怎样？**  
A: 服务端自动执行默认操作（随机出牌/过牌），发送 `PLAYER_TIMEOUT` 事件。

---

## 参考

- **Proto 总览：** [README.md](./README.md)
- **后端使用（Go）：** [BACKEND_GUIDE.md](./BACKEND_GUIDE.md)
- **Proto 定义：** `common.proto`, `event.proto`, `view.proto`
- **HTTP API：** `backend/API-Documentation.md`
- **架构文档：** `docs/EventSystemArchitecture.md`
