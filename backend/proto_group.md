# Proto 改造分组清单

## 分组原则

- **自底向上**：先定义基础类型，再定义依赖它们的复杂类型
- **控制工作量**：每组 3-8 个 message，避免单组过大
- **依赖隔离**：同组内结构相互独立或依赖关系简单，跨组依赖单向

---

## 第一组：基础类型层（Foundation）

**依赖**：无  
**工作量**：5 个文件，约 150 行  
**完成标准**：所有基础类型和枚举定义完成，可编译生成 Go 代码

### 文件清单

#### 1. `card.proto`
```
message Card
```

**对应 Go**：`sdk/card.go:Card`

---

#### 2. `player.proto`
```
message Player
```

**对应 Go**：`sdk/types.go:Player`

---

#### 3. `enums.proto`
```
enum VictoryType
enum DealStatus
enum MatchStatus
enum TrickStatus
enum TributeStatus
enum CompType
enum GameEventType
```


**对应 Go**：`sdk/match.go:VictoryType`, `sdk/types.go:*Status`, `sdk/comp.go:CompType`, `sdk/game_engine.go:GameEventType`

---

## 第二组：牌型层（CardComp）

**依赖**：第一组（Card, CompType）  
**工作量**：1 个文件，约 120 行  
**完成标准**：所有牌型组合定义完成，oneof 多态实现正确

### 文件清单

#### 4. `proto/game/card_comp.proto`
```
message CardComp
message SingleComp
message PairComp
message TripleComp
message FullHouseComp
message StraightComp
message PlateComp
message TubeComp
message BombComp
message StraightFlushComp
```

**对应 Go**：`sdk/comp.go:CardComp 接口及其实现`

---

## 第三组：游戏动作层（Actions & Phases）

**依赖**：第一组 + 第二组  
**工作量**：1 个文件，约 100 行  
**完成标准**：Trick、PlayAction、TributePhase 定义完成

### 文件清单

#### 5. `proto/game/actions.proto`
```
message PlayAction
message Trick
message TributePhase
```

**对应 Go**：`sdk/types.go:PlayAction, Trick, TributePhase`

---

## 第四组：统计结果层（Statistics & Results）

**依赖**：第一组  
**工作量**：1 个文件，约 150 行  
**完成标准**：所有统计和结果结构定义完成

### 文件清单

#### 6. `proto/game/result.proto`
```
message PlayerDealStats
message TributeInfo
message DealStatistics
message DealResult
message TeamMatchStats
message MatchStatistics
message MatchResult
message TeamUpgrades
```

**对应 Go**：`sdk/result.go:PlayerDealStats, TributeInfo, DealStatistics, DealResult, TeamMatchStats, MatchStatistics, MatchResult`

---

## 第五组：游戏实体层（Game Entities）

**依赖**：第一组 + 第三组 + 第四组  
**工作量**：2 个文件，约 120 行  
**完成标准**：Deal、Match、PlayerView 定义完成

### 文件清单

#### 7. `proto/game/deal.proto`
```
message Deal
message PlayerHand
```

**对应 Go**：`sdk/types.go:Deal`

---

#### 8. `proto/game/match.proto`
```
message Match
message PlayerView
```

**对应 Go**：`sdk/types.go:Match`, `sdk/game_engine.go:PlayerView`

---

## 第6A组：事件基础与生命周期（Event Foundation & Lifecycle）

**依赖**：所有前面的组  
**工作量**：1 个文件，约 120 行  
**完成标准**：GameEvent 基础消息及比赛/局生命周期事件定义完成

### 文件清单

#### 9. `proto/messages/game_event.proto`
```
message GameEvent
message MatchStartedEvent
message DealStartedEvent
message CardsDealtEvent
message DealEndedEvent
message MatchEndedEvent
```

**对应 Go**：`sdk/game_engine.go:GameEvent` 及对应事件的 Data 字段

---

## 第6B组：进贡事件（Tribute Events）

**依赖**：第6A组 + 所有前面的组  
**工作量**：1 个文件，约 180 行  
**完成标准**：所有进贡阶段相关事件定义完成

### 文件清单

#### 10. `proto/messages/tribute_events.proto`
```
message TributePhaseEvent
message TributeRulesSetEvent
message TributeImmunityEvent
message TributePoolCreatedEvent
message TributeStartedEvent
message TributeGivenEvent
message TributeSelectedEvent
message ReturnTributeEvent
message TributeCompletedEvent
```

**对应 Go**：`sdk/game_engine.go:GameEvent` 中进贡相关事件的 Data 字段

---

## 第6C组：出牌事件（Trick Events）

**依赖**：第6A组 + 所有前面的组  
**工作量**：1 个文件，约 80 行  
**完成标准**：所有出牌阶段相关事件定义完成

### 文件清单

#### 11. `proto/messages/trick_events.proto`
```
message TrickStartedEvent
message PlayerPlayedEvent
message PlayerPassedEvent
message TrickEndedEvent
```

**对应 Go**：`sdk/game_engine.go:GameEvent` 中出牌相关事件的 Data 字段

---

## 第6D组：连接事件（Connection Events）

**依赖**：第6A组 + 所有前面的组  
**工作量**：1 个文件，约 60 行  
**完成标准**：所有连接、断线、超时事件定义完成

### 文件清单

#### 12. `proto/messages/connection_events.proto`
```
message PlayerTimeoutEvent
message PlayerDisconnectEvent
message PlayerReconnectEvent
```

**对应 Go**：`sdk/game_engine.go:GameEvent` 中连接相关事件的 Data 字段

---

## 第6E组：WebSocket消息（WebSocket Messages）

**依赖**：第6A组（GameEvent）  
**工作量**：1 个文件，约 80 行  
**完成标准**：WebSocket 消息封装完成

### 文件清单

#### 13. `proto/messages/ws_message.proto`
```
message WSMessage
message WSRequest
message WSResponse
message ErrorInfo
```

**对应 Go**：`backend/websocket/manager.go:WSMessage`

---

## 依赖关系图

```
第一组（基础）
    ↓
第二组（牌型）
    ↓
第三组（动作）  ←── 第一组
    ↓
第四组（结果）  ←── 第一组
    ↓
第五组（实体）  ←── 第一、三、四组
    ↓
第6A组（事件基础）  ←── 所有组
    ↓
第6B组（进贡事件）  ←── 6A + 所有前组
第6C组（出牌事件）  ←── 6A + 所有前组
第6D组（连接事件）  ←── 6A + 所有前组
    ↓
第6E组（WS消息）  ←── 6A
```
