# Protobuf 后端使用指南（Go）

本文档面向 Go 后端开发者。

> **前端开发者请查看：** [FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)  
> **Proto 总览：** [README.md](./README.md)

## 快速开始

### 1. 代码生成

```bash
make proto-go
```

### 2. 导入包

```go
import (
    pbcommon "guandan-world/proto/common"
    pbevent "guandan-world/proto/event"
    pbview "guandan-world/proto/view"
)
```

### 3. 使用示例

```go
// 创建玩家出牌事件
event := &pbevent.GameEvent{
    MatchId:     "match-123",
    DealIndex:   ptrint32(0),
    TrickIndex:  ptrint32(5),
    ActorSeat:   ptrint32(2),
    Seq:         42,
    CreatedAtMs: time.Now().UnixMilli(),
    Type:        pbevent.EventType_EVENT_TYPE_PLAYER_PLAYED,
    Payload: &pbevent.GameEvent_PlayerPlayed{
        PlayerPlayed: &pbevent.PlayerPlayedPayload{
            Cards: []*pbcommon.Card{
                {Suit: 0, Rank: 13, DeckIndex: 12},
            },
        },
    },
}

// 序列化为 JSON（用于 WebSocket 传输）
jsonData, _ := protojson.Marshal(event)

// 创建玩家视图
view := &pbview.PlayerView{
    MatchId:     "match-123",
    DealIndex:   0,
    PlayerSeat:  0,
    PlayerCards: []*pbcommon.Card{...},
    TeamLevels:  []int32{2, 2},
    DealLevel:   2,
    DealStatus:  pbview.DealStatus_DEAL_STATUS_PLAYING,
    CurrentTurn: ptrint32(1),
    Plays:       []*pbcommon.PlayAction{...},
    Leader:      ptrint32(0),
}
```

---

## 字段填写规则

### GameEvent 元数据字段

| 字段 | 必填 | 填写规则 |
|-----|------|---------|
| `match_id` | ✅ | 所有事件 |
| `deal_index` | ⚠️ | 局级及以下事件，否则不设置（nil） |
| `trick_index` | ⚠️ | 轮级及玩家行为事件，否则不设置 |
| `actor_seat` | ⚠️ | 玩家行为事件，否则不设置 |
| `seq` | ✅ | 全局递增序列号 |
| `created_at_ms` | ✅ | 事件创建时间（Unix 毫秒） |
| `type` | ✅ | 对应的 `EventType` 枚举值 |

### 事件层级分类

```
比赛级 (match_id)
  ├── MATCH_STARTED
  └── MATCH_ENDED

局级 (match_id + deal_index)
  ├── DEAL_STARTED
  ├── CARDS_DEALT
  ├── DEAL_ENDED
  └── 进贡阶段 (6个事件)

轮次级 (match_id + deal_index + trick_index)
  ├── TRICK_STARTED
  └── TRICK_ENDED

玩家行为级 (match_id + deal_index + trick_index + actor_seat)
  ├── PLAYER_PLAYED
  ├── PLAYER_PASSED
  ├── PLAYER_TIMEOUT
  ├── PLAYER_DISCONNECT
  └── PLAYER_RECONNECT
```

---

## 序列化

### JSON（用于 WebSocket）

```go
import "google.golang.org/protobuf/encoding/protojson"

// 转换为 JSON
jsonBytes, err := protojson.Marshal(event)

// 从 JSON 解析
event := &pbevent.GameEvent{}
err = protojson.Unmarshal(jsonBytes, event)
```

### 二进制（用于持久化）

```go
import "google.golang.org/protobuf/proto"

// 序列化
data, err := proto.Marshal(event)

// 反序列化
event := &pbevent.GameEvent{}
err = proto.Unmarshal(data, event)
```

---

## 类型判断

### 事件类型判断

```go
switch payload := event.Payload.(type) {
case *pbevent.GameEvent_MatchStarted:
    fmt.Printf("比赛开始: %d 名玩家\n", len(payload.MatchStarted.Players))
    
case *pbevent.GameEvent_PlayerPlayed:
    fmt.Printf("玩家 %d 出牌: %d 张\n", 
        event.GetActorSeat(), len(payload.PlayerPlayed.Cards))
    
case *pbevent.GameEvent_DealEnded:
    fmt.Printf("局结束，获胜队伍: %d\n", payload.DealEnded.WinningTeam)
}
```

### 状态判断

```go
// 检查牌局状态
if view.DealStatus == pbview.DealStatus_DEAL_STATUS_PLAYING {
    currentTurn := view.GetCurrentTurn()  // 使用 Get 方法安全访问 optional 字段
    if currentTurn >= 0 {
        fmt.Printf("当前轮到玩家 %d\n", currentTurn)
    }
}
```

---

## 最佳实践

### 1. 使用 Get 方法访问 optional 字段

```go
// ✅ 推荐：使用 Get 方法
currentTurn := view.GetCurrentTurn()  // 如果未设置返回 0

// ❌ 不推荐：直接访问可能为 nil
if view.CurrentTurn != nil {
    turn := *view.CurrentTurn
}
```

### 2. 利用 Type 字段快速判断

```go
// ✅ 推荐：先检查 Type 字段
if event.Type == pbevent.EventType_EVENT_TYPE_PLAYER_PLAYED {
    cards := event.GetPlayerPlayed().Cards
}

// ❌ 不推荐：直接 switch payload（效率低）
switch event.Payload.(type) {
case *pbevent.GameEvent_PlayerPlayed:
    // ...
}
```

### 3. 不要手动编辑生成的代码

```bash
# 修改 proto 文件后，立即重新生成
make proto-go
```

### 4. 使用默认值机制

```go
// 零值字段不会被序列化（节省空间）
event := &pbevent.GameEvent{
    MatchId: "match-123",
    Seq:     1,
    // deal_index, trick_index, actor_seat 不设置
}
```

---

## 向前兼容性

### 添加新字段

```protobuf
message DealEndedPayload {
  // 现有字段
  repeated int32 rankings = 1;
  VictoryType victory_type = 2;
  
  // 新增字段使用新编号
  string mvp_player_id = 7;  // ✅ 兼容旧版本
}
```

### 废弃字段

```protobuf
message Card {
  int32 suit = 1;
  int32 rank = 2;
  string legacy_id = 3 [deprecated = true];  // ✅ 标记废弃，不删除
}
```

---

## 工具函数

```go
// 创建 optional int32
func ptrint32(v int32) *int32 { return &v }

// 卡牌转字符串（调试用）
func CardString(c *pbcommon.Card) string {
    suits := []string{"♠", "♥", "♣", "♦"}
    ranks := []string{"", "", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "小王", "大王"}
    if c.Rank == 15 || c.Rank == 16 {
        return ranks[c.Rank]
    }
    return suits[c.Suit] + ranks[c.Rank]
}
```

---

## 依赖

```go
require (
    google.golang.org/protobuf v1.36.10
)
```

---

## 参考

- **Proto 总览：** [README.md](./README.md)
- **前端使用：** [FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)
- **Proto 定义：** `common.proto`, `event.proto`, `view.proto`
- **事件系统架构：** `docs/EventSystemArchitecture.md`
- **API 文档：** `backend/API-Documentation.md`
