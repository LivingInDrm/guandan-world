# 事件系统 Proto 使用指南

## 概述

事件系统已重构为基于 Protocol Buffers 的单一来源定义，所有事件类型和字段通过 `proto/event.proto` 定义，代码自动生成。

## 目录结构

```
proto/
├── event.proto          # Proto 定义文件（单一来源）
└── event/
    └── event.pb.go      # 生成的 Go 代码
```

## 代码生成

### 生成 Go 代码

```bash
make proto-go
```

### 未来支持的生成目标

```bash
make proto-js   # 生成 JavaScript 代码（需要配置）
make proto-ts   # 生成 TypeScript 代码（需要配置）
```

## 使用示例

### 1. 导入生成的包

```go
import (
    pb "guandan-world/proto/event"
)
```

### 2. 创建比赛开始事件

```go
event := &pb.GameEvent{
    MatchId:      "match-123",
    DealIndex:    -1,  // 比赛级事件不填局索引
    TrickIndex:   -1,  // 比赛级事件不填轮索引
    ActorSeat:    -1,  // 非玩家行为事件不填
    Seq:          1,
    CreatedAtMs:  time.Now().UnixMilli(),
    Payload: &pb.GameEvent_MatchStarted{
        MatchStarted: &pb.MatchStartedPayload{
            Players: []*pb.PlayerBasicInfo{
                {Id: "p1", Username: "Player1", Seat: 0, TeamNum: 0},
                {Id: "p2", Username: "Player2", Seat: 1, TeamNum: 1},
                {Id: "p3", Username: "Player3", Seat: 2, TeamNum: 0},
                {Id: "p4", Username: "Player4", Seat: 3, TeamNum: 1},
            },
            InitialLevels: []int32{2, 2},
        },
    },
}
```

### 3. 创建玩家出牌事件

```go
event := &pb.GameEvent{
    MatchId:      "match-123",
    DealIndex:    0,     // 第一局
    TrickIndex:   5,     // 第六轮
    ActorSeat:    2,     // 座位号2的玩家
    Seq:          42,
    CreatedAtMs:  time.Now().UnixMilli(),
    Payload: &pb.GameEvent_PlayerPlayed{
        PlayerPlayed: &pb.PlayerPlayedPayload{
            Cards: []*pb.Card{
                {Suit: 0, Rank: 13, DeckIndex: 12},  // 黑桃K
                {Suit: 1, Rank: 13, DeckIndex: 38},  // 红桃K
            },
        },
    },
}
```

### 4. 序列化和反序列化

```go
import "google.golang.org/protobuf/proto"

// 序列化为二进制
data, err := proto.Marshal(event)
if err != nil {
    log.Fatal(err)
}

// 反序列化
newEvent := &pb.GameEvent{}
err = proto.Unmarshal(data, newEvent)
if err != nil {
    log.Fatal(err)
}
```

### 5. JSON 序列化

```go
import "google.golang.org/protobuf/encoding/protojson"

// 转换为 JSON
jsonData, err := protojson.Marshal(event)
if err != nil {
    log.Fatal(err)
}

// 从 JSON 解析
newEvent := &pb.GameEvent{}
err = protojson.Unmarshal(jsonData, newEvent)
if err != nil {
    log.Fatal(err)
}
```

### 6. 事件类型判断

```go
switch payload := event.Payload.(type) {
case *pb.GameEvent_MatchStarted:
    fmt.Printf("比赛开始，玩家数：%d\n", len(payload.MatchStarted.Players))
    
case *pb.GameEvent_PlayerPlayed:
    fmt.Printf("玩家 %d 出了 %d 张牌\n", event.ActorSeat, len(payload.PlayerPlayed.Cards))
    
case *pb.GameEvent_DealEnded:
    fmt.Printf("牌局结束，获胜队伍：%d\n", payload.DealEnded.WinningTeam)
    
default:
    fmt.Println("未知事件类型")
}
```

## 字段填写规则

### 元数据字段

| 字段            | 必填 | 填写条件                               |
|----------------|------|---------------------------------------|
| `match_id`     | ✅   | 所有事件                              |
| `deal_index`   | ⚠️   | 局级及以下事件，其他填 -1              |
| `trick_index`  | ⚠️   | 轮级及玩家行为事件，其他填 -1          |
| `actor_seat`   | ⚠️   | 玩家行为事件，其他填 -1                |
| `seq`          | ✅   | 所有事件                              |
| `created_at_ms`| ✅   | 所有事件                              |

### 事件层级分类

**比赛级**（match_id 必填，其他索引 -1）
- `MatchStartedPayload`
- `MatchEndedPayload`

**局级**（match_id + deal_index 必填）
- `DealStartedPayload`
- `CardsDealtPayload`
- `DealEndedPayload`

**进贡阶段**（match_id + deal_index 必填）
- `TributePhaseStartedPayload`
- `TributeExemptedPayload`
- `TributeCardSubmittedPayload`
- `TributeCardSelectedPayload`
- `ReturnTributePayload`
- `TributeCompletedPayload`

**轮次级**（match_id + deal_index + trick_index 必填）
- `TrickStartedPayload`
- `TrickEndedPayload`

**玩家行为级**（match_id + deal_index + trick_index + actor_seat 必填）
- `PlayerPlayedPayload`
- `PlayerPassedPayload`
- `PlayerTimeoutPayload`
- `PlayerDisconnectPayload`
- `PlayerReconnectPayload`

## 向前兼容性

### 添加新字段

在现有 message 中添加新字段时，使用新的字段编号：

```protobuf
message DealEndedPayload {
  repeated int32 rankings = 1;
  VictoryType victory_type = 2;
  int32 winning_team = 3;
  repeated int32 level_change = 4;
  int64 duration_ms = 5;
  int32 trick_count = 6;
  // 新增字段使用新编号
  string mvp_player_id = 7;  // 未来添加
}
```

### 添加新事件类型

在 `GameEvent.payload` oneof 中添加新条目：

```protobuf
message GameEvent {
  // ... 元数据字段 ...
  
  oneof payload {
    // ... 现有事件 ...
    PlayerReconnectPayload player_reconnect = 54;
    
    // 新增事件类型
    GamePausedPayload game_paused = 60;  // 未来添加
    GameResumedPayload game_resumed = 61;
  }
}
```

### 废弃字段

不要删除字段，使用 `deprecated` 标记：

```protobuf
message Card {
  int32 suit = 1;
  int32 rank = 2;
  int32 deck_index = 3;
  string legacy_id = 4 [deprecated = true];  // 废弃字段
}
```

## 最佳实践

1. **永远不要手动编辑生成的 `.pb.go` 文件**
2. **修改 proto 定义后立即重新生成代码**：`make proto-go`
3. **使用字段名访问而非字段编号**
4. **利用 protobuf 的默认值机制**：零值字段不会被序列化
5. **在需要区分零值和未设置时，使用 `optional` 关键字**

## 依赖包

生成的 Go 代码依赖以下包：

```go
import (
    "google.golang.org/protobuf/proto"           // 核心序列化
    "google.golang.org/protobuf/encoding/protojson" // JSON 支持
)
```

确保在 `go.mod` 中包含：

```
require google.golang.org/protobuf v1.36.10
```

## 下一步

1. 在 `sdk/` 中创建适配器层，将现有 `GameEvent` struct 转换为 proto 事件
2. 逐步迁移事件创建代码使用 proto 类型
3. 为前端配置 TypeScript/JavaScript 代码生成
4. 实现事件序列化到数据库（如需持久化）
