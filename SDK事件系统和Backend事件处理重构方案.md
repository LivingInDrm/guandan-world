# SDK事件系统和Backend事件处理重构方案

## 目标

解决两个核心问题：
1. SDK 事件定义分散且混乱
2. Backend 事件处理逻辑职责不清、缺乏安全过滤

## 一、SDK 事件定义重构

### 1.1 创建 `sdk/events/` 目录结构

```
sdk/events/
├── types.go           # 事件类型常量定义
├── event.go           # GameEvent 结构和基础接口
├── match_events.go    # 比赛级事件构造函数
├── deal_events.go     # 牌局级事件构造函数  
├── trick_events.go    # 轮次级事件构造函数
├── tribute_events.go  # 进贡事件构造函数
└── player_events.go   # 玩家事件构造函数
```

### 1.2 迁移事件类型常量

**文件**: `sdk/events/types.go`

- 从 `sdk/game_engine.go` 迁移 `GameEventType` 和所有常量定义（line 10-38）
- 保留原位置的 type alias，确保向后兼容

### 1.3 迁移 GameEvent 结构

**文件**: `sdk/events/event.go`

- 从 `sdk/game_engine.go` 迁移 `GameEvent` 结构（line 40-47）
- 迁移 `GameEventHandler` 类型（line 50-58）
- 保留原位置的 type alias

### 1.4 创建事件构造函数

#### 比赛级事件 (`sdk/events/match_events.go`)

提取以下事件创建逻辑：
- `NewMatchStartedEvent(*Match)` ← 从 `game_engine.go:368-373`
- `NewMatchEndedEvent(*Match, *MatchResult)` ← 从 `game_engine.go:852-863`

#### 牌局级事件 (`sdk/events/deal_events.go`)

提取以下事件创建逻辑：
- `NewDealStartedEvent(*Deal, [2]int)` ← 从 `game_engine.go:396-406`
- `NewDealEndedEvent(*Deal, *DealResult)` ← 从 `game_engine.go:830-840`

#### 轮次级事件 (`sdk/events/trick_events.go`)

提取以下事件创建逻辑：
- `NewTrickStartedEvent(*Trick, map[int][]*Card)` ← 从 `game_engine.go:781-791`
- `NewTrickEndedEvent(*Trick)` ← 从 `game_engine.go:868-879`
- `NewPlayerPlayedEvent(playerSeat int, cards []*Card, *Deal)` ← 从 `game_engine.go:502-512`
- `NewPlayerPassedEvent(playerSeat int, *Deal)` ← 从 `game_engine.go:556-565`

#### 进贡事件 (`sdk/events/tribute_events.go`)

提取以下事件创建逻辑：
- `NewTributeRulesSetEvent(*DealResult, map[int]int, bool, string)` ← 从 `game_engine.go:429-443`
- `NewTributeImmunityEvent(*TributePhase, map[string]interface{})` ← 从 `game_engine.go:455-463`
- `NewTributePoolCreatedEvent([]map[string]interface{}, []int, []*Card, int)` ← 从 `game_engine.go:1082-1093`
- `NewTributeGivenEvent(giver, receiver int, *Card, string, bool)` ← 从 `game_engine.go:1105-1120`
- `NewTributeSelectedEvent(playerID int, *Card, []*Card, []*Card, int)` ← 从 `game_engine.go:1197-1215`
- `NewReturnTributeEvent(playerID, target int, *Card, *Card)` ← 从 `game_engine.go:1261-1278`
- `NewTributeCompletedEvent(*TributePhase)` ← 从 `game_engine.go:1134-1138`

#### 玩家事件 (`sdk/events/player_events.go`)

提取以下事件创建逻辑：
- `NewPlayerDisconnectEvent(playerSeat int)` ← 从 `game_engine.go:670-679`
- `NewPlayerReconnectEvent(playerSeat int)` ← 从 `game_engine.go:701-710`
- `NewPlayerTimeoutEvent(playerSeat int, actionType string)` ← 从 `game_driver.go:690-698`

### 1.5 修改 `sdk/game_engine.go`

- 移除事件类型常量和结构定义，改为 import `sdk/events`
- 将所有内联事件构造替换为调用 `events.New*Event()`
- 简化方法体，提高可读性

示例：
```go
// Before
event := &GameEvent{
    Type: EventMatchStarted,
    Data: match,
    Timestamp: time.Now(),
}
ge.emitEvent(event)

// After  
event := events.NewMatchStartedEvent(match)
ge.emitEvent(event)
```

### 1.6 修改 `sdk/game_driver.go`

- Import `sdk/events`
- 替换超时事件创建为 `events.NewPlayerTimeoutEvent(playerSeat, actionType)`

---

## 二、Backend 事件处理重构

### 2.1 创建 `backend/game/event_broadcaster.go`

实现事件广播策略和安全过滤：

```go
// EventBroadcaster 负责决定事件广播策略和安全过滤
type EventBroadcaster struct {
    wsManager WSManagerInterface
    engine    sdk.GameEngineInterface
}

// shouldBroadcast 判断事件是否应该广播
func (eb *EventBroadcaster) shouldBroadcast(eventType sdk.GameEventType) bool

// shouldSendPlayerView 判断事件是否需要发送 PlayerView
func (eb *EventBroadcaster) shouldSendPlayerView(eventType sdk.GameEventType) bool

// broadcastPublicEvent 广播公开事件（过滤敏感数据）
func (eb *EventBroadcaster) broadcastPublicEvent(roomID string, event *sdk.GameEvent)

// sendPlayerViews 发送玩家视角数据
func (eb *EventBroadcaster) sendPlayerViews(roomID string, eventType sdk.GameEventType)
```

**职责划分**：

1. **广播策略**：判断哪些事件类型需要广播
2. **安全过滤**：从事件中提取公开信息，移除敏感数据（如其他玩家手牌）
3. **单播策略**：判断哪些事件需要单独发送 PlayerView

### 2.2 修改 `backend/game/driver_service.go`

#### WebSocketObserver 重构

```go
type WebSocketObserver struct {
    roomID      string
    broadcaster *EventBroadcaster  // 使用broadcaster代替直接调用wsManager
}

func (wso *WebSocketObserver) OnGameEvent(event *sdk.GameEvent) {
    // 委托给 broadcaster 处理所有广播逻辑
    wso.broadcaster.broadcastPublicEvent(wso.roomID, event)
    
    // 根据事件类型发送 PlayerView
    if wso.broadcaster.shouldSendPlayerView(event.Type) {
        wso.broadcaster.sendPlayerViews(wso.roomID, event.Type)
    }
    
    // 日志记录（可选）
    logSignificantEvents(event)
}
```

移除 `sendPlayerViews()` 方法（line 677-720），迁移到 `EventBroadcaster`

### 2.3 安全过滤规则

定义事件数据公开规则：

| 事件类型 | 过滤规则 |
|---------|---------|
| EventDealStarted | 保留 deal_level, team_levels；移除 player_cards |
| EventPlayerPlayed | 保留 player_seat, cards；移除 deal_state.player_cards |
| EventPlayerPassed | 保留 player_seat；移除 deal_state.player_cards |
| EventTributeGiven | 保留 giver, receiver, card；移除其他玩家手牌 |
| EventReturnTribute | 保留 player, return_card, target_player；移除手牌 |

### 2.4 PlayerView 发送策略

明确哪些事件触发 PlayerView 发送：

**需要发送 PlayerView 的事件**：
- EventMatchStarted
- EventDealStarted
- EventCardsDealt
- EventTributeGiven
- EventReturnTribute
- EventTributeCompleted
- EventTrickStarted
- EventPlayerPlayed
- EventPlayerPassed
- EventTrickEnded
- EventDealEnded
- EventMatchEnded

**不需要发送 PlayerView 的事件**：
- EventTributeRulesSet（仅规则声明，无手牌变化）
- EventTributeImmunity（仅状态声明）
- EventTributePoolCreated（双下选牌前）
- EventPlayerTimeout（状态通知）
- EventPlayerDisconnect/Reconnect（连接状态）

---

## 三、实施步骤

### 阶段 1：SDK 事件系统重构

1. 创建 `sdk/events/` 目录和基础文件
2. 迁移类型定义到 `types.go` 和 `event.go`
3. 在原位置添加 type alias，保持向后兼容
4. 实现各类事件构造函数
5. 修改 `game_engine.go` 和 `game_driver.go`，替换所有事件构造调用

### 阶段 2：Backend 事件处理重构

1. 创建 `backend/game/event_broadcaster.go`
2. 实现广播策略、安全过滤、PlayerView 发送逻辑
3. 修改 `WebSocketObserver`，委托给 `EventBroadcaster`
4. 移除 `driver_service.go` 中的 `sendPlayerViews()` 方法

---

## 四、向后兼容性

### SDK 层

在 `sdk/game_engine.go` 保留 type alias：

```go
// Type aliases for backward compatibility
type GameEventType = events.GameEventType
type GameEvent = events.GameEvent
type GameEventHandler = events.GameEventHandler

// Event type constants for backward compatibility
const (
    EventMatchStarted = events.EventMatchStarted
    EventDealStarted  = events.EventDealStarted
    // ... 所有事件类型
)
```

### Backend 层

`backend/game/driver_service.go` 的 API 不变，只改变内部实现。

---

## 五、优势

### SDK 层改进

1. **集中管理**：所有事件定义和构造集中在 `sdk/events/` 包中
2. **职责清晰**：按事件领域分文件（match/deal/trick/tribute/player）
3. **易于维护**：新增/修改事件只需修改对应文件
4. **代码简化**：`game_engine.go` 从 1300+ 行缩减到 900 行左右
5. **可测试性**：事件构造函数独立可测

### Backend 层改进

1. **职责分离**：`WebSocketObserver` 只负责调度，`EventBroadcaster` 负责策略
2. **安全性**：统一的安全过滤入口，杜绝敏感数据泄露
3. **策略集中**：广播策略和 PlayerView 发送策略集中管理
4. **易于扩展**：新增事件类型只需修改策略函数
5. **可维护性**：清晰的架构边界，降低心智负担
