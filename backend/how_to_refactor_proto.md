# Proto 改造标准流程

---

## 统一规范

### 1. 目录结构

```
下述路径在代码库根目录下，注意严格遵循，不要放到backend或者其他目录下。

proto/
├── common/              # 基础类型（Card, Player, 枚举）
├── game/                # 游戏类型（CardComp, Deal, Match, Result）
└── messages/            # 事件消息（GameEvent, WSMessage）

proto/gen/go/            # 生成的 Go 代码（对应上述结构）
├── common/
├── game/
└── messages/
```

---

### 2. Package 规范

```protobuf
// 文件头固定格式
syntax = "proto3";
package guandan.{目录名};                           // common / game / messages
option go_package = "guandan-world/proto/gen/go/{目录名}";

// Import 路径：从 proto/ 开始的相对路径
import "common/card.proto";
import "game/actions.proto";
```

---

### 3. 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| Message | `UpperCamelCase` | `Card`, `DealResult` |
| 字段 | `snake_case` | `deck_index`, `player_seat` |
| 时间字段 | `xxx_ms` | `start_time_ms`, `duration_ms` |
| 布尔字段 | `is_/has_/can_` | `is_pass`, `has_tribute` |
| Enum 类型 | `UpperCamelCase` | `VictoryType`, `DealStatus` |
| Enum 值 | `UPPER_SNAKE_CASE` + 类型前缀 | `VICTORY_TYPE_DOUBLE_DOWN` |
| Enum 零值 | `{TYPE}_UNSPECIFIED = 0` | `VICTORY_TYPE_UNSPECIFIED = 0` |
| Oneof | `snake_case` | `oneof payload { ... }` |

---

### 4. 类型映射表

| Go 类型 | Proto 类型 | 说明 |
|---------|-----------|------|
| `int` | `int32` | 座位号、牌数等小整数 |
| `int64` | `int64` | ID、时间戳 |
| `string` | `string` | |
| `bool` | `bool` | |
| `[]T` | `repeated T` | |
| `map[K]V` | `map<K, V>` | |
| `*T` | `T` | proto3 message 默认可选 |
| `time.Time` | `int64 xxx_ms` | 毫秒时间戳 |
| `time.Duration` | `int64 xxx_ms` | 毫秒数 |
| `[2]int` (Team) | 独立 message | `TeamUpgrades { int32 team0; int32 team1; }` |
| `[4]T` (Players) | `repeated T` | 注释标注"固定 4 个" |
| 枚举 `type X string` | `enum X` | 必须有 `XXX_UNSPECIFIED = 0` |
| 接口 `interface{}` | `oneof` | 列举所有可能的类型 |

---

### 5. 字段编号规则

| 范围 | 用途 | 说明 |
|------|------|------|
| 1-15 | 核心字段 | 高频访问，单字节编码 |
| 16-99 | 常用字段 | |
| 100+ | 预留扩展 | |
| oneof | 从 10 开始 | 为公共字段预留 1-9 |

**废弃字段**：使用 `reserved`，永不重用编号
```protobuf
reserved 4, 5;
reserved "old_field_name";
```

---

### 6. 注释规范

```protobuf
// Message 总注释：说明用途
message Card {
  // 字段注释：含义 + 值域 + 特殊说明
  int32 number = 1;     // 牌的数字值 (1-16): 11=J, 12=Q, 13=K, 14=A, 15/16=Joker
  string color = 2;     // 花色: "Spade", "Club", "Heart", "Diamond", "" for Joker
  int32 deck_index = 3; // 唯一索引 (0-107)
}
```

**原则**：每个字段必须有注释，说明值域、特殊情况

---

### 7. 适配器规范

**函数命名与访问权限**：
```go
// SDK 基础适配器：公开（大写开头），供其他包复用
// Go → Proto
func ToProtoCard(c *Card) *pb.Card
func ToProtoCards(cards []*Card) []*pb.Card        // 注意: nil元素保留为nil
func ToProtoPlayersArray(players [4]*Player) []*pb.Player

// Proto → Go
func FromProtoCard(pc *pb.Card) *Card
func FromProtoCards(pcs []*pb.Card) []*Card        // 注意: nil元素保留为nil
func FromProtoPlayersArray(pps []*pb.Player) [4]*Player

// Backend 业务适配器：私有（小写开头），内部使用
func toGameEventProto(e *GameEvent) *pbmsg.GameEvent
func fromGameEventProto(pe *pbmsg.GameEvent) *GameEvent
```

**Import 别名**：
```go
import (
    pb "guandan-world/proto/gen/go/common"
    pbgame "guandan-world/proto/gen/go/game"
    pbmsg "guandan-world/proto/gen/go/messages"
)
```

**文件位置**：
```
sdk/proto_adapter.go              # SDK ↔ Proto 基础适配器（公开）
backend/game/event_adapter.go     # GameEvent 适配器（私有）
backend/websocket/serializer.go   # 序列化工具（私有）
```

---

### 8. 编译命令

```bash
make proto-deps   # 安装编译依赖（protoc-gen-go）
make proto        # 编译所有 proto
make proto-check  # 仅检查语法
make proto-clean  # 清理生成文件
```

---

### 9. 快速检查清单

每个 proto 文件完成后，逐项检查：

```markdown
### 文件：`xxx.proto`

**基础**
- [ ] 文件位置：proto/{common|game|messages}/
- [ ] package：guandan.{目录名}
- [ ] go_package：guandan-world/proto/gen/go/{目录名}
- [ ] import 路径从 proto/ 开始

**命名**
- [ ] Message：UpperCamelCase
- [ ] 字段：snake_case
- [ ] 时间字段：xxx_ms
- [ ] 枚举值：UPPER_SNAKE_CASE + 类型前缀
- [ ] 枚举有 UNSPECIFIED = 0

**字段**
- [ ] 核心字段在 1-15
- [ ] oneof 从 10 开始
- [ ] 字段编号无冲突
- [ ] 每个字段有注释

**验证**
- [ ] 编译通过：make proto
- [ ] 生成文件在：proto/gen/go/{目录名}/
- [ ] 字段数匹配：Go struct = Proto message
```

---

### 10. 关键原则

1. **一致性优先**：所有 proto 遵循相同规范
2. **类型安全**：枚举用 enum，接口用 oneof
3. **向后兼容**：字段编号预留，废弃用 reserved
4. **可读性**：每个字段必须注释
5. **简洁性**：避免过度嵌套，独立定义复杂类型

---

## 改造流程（5 步骤）

对每个 proto 的改造分成 **5 个步骤**，形成标准化流程。这 5 步是一个**从抽象到具体、从设计到实现**的完整闭环：

```
问题定义 → 方案设计 → 代码实现 → 正确性验证 → 使用准备
   ↓           ↓           ↓             ↓            ↓
 分析        设计        实现          验证         适配器
(What)      (How)      (Code)       (Check)      (Use)
```

每一步都在降低下一步的风险，避免"边写边改"的混乱状态。

---

### 步骤 1：读取并分析 Go 结构（Analyze）

**解决问题**：**"我需要转换什么？"**

- 明确源数据的完整结构（避免遗漏字段）
- 识别特殊类型（知道哪里有坑）
- 理解数据的使用上下文（避免过度设计或设计不足）

**防止问题**：盲目开始导致字段遗漏、类型理解错误

**目标**：完整理解源数据结构

**操作**：
1. 定位 Go struct 定义文件和行号
2. 记录所有字段：名称、类型、JSON tag、注释
3. 识别嵌套结构和依赖关系
4. 查找该结构的使用场景（被哪些函数创建/读取）

**输出检查清单**：
```markdown
- [ ] Go struct 位置已确认：`文件:行号`
- [ ] 所有字段（共 N 个）已列出
- [ ] 嵌套结构已识别（需要 import 哪些 proto）
- [ ] 特殊类型已标记（time.Time / interface{} / 固定数组 / 指针）
```

---

### 步骤 2：设计 Proto Message（Design）

**解决问题**：**"怎么转换最合理？"**

- 确定类型映射策略（`time.Time` 用毫秒？`[4]T` 用 repeated 还是独立字段？）
- 规划字段编号（避免后续冲突，预留扩展空间）
- 设计多态表达（`interface{}` 怎么用 `oneof`？）

**防止问题**：转换策略不一致、字段编号混乱、后续难以扩展

**目标**：完成字段映射和类型转换设计

**操作**：
1. **字段映射**：逐字段设计 proto 字段（参考上方"类型映射表"）
2. **字段编号规划**：参考上方"字段编号规则"
3. **命名转换**：Go `CamelCase` → Proto `snake_case`

**输出检查清单**：
```markdown
- [ ] 所有字段已映射（Go 字段数 = Proto 字段数）
- [ ] 时间类型统一为 int64 毫秒（xxx_ms）
- [ ] 固定数组处理方式已确定（独立 message / repeated + 注释长度）
- [ ] 枚举有 UNSPECIFIED = 0 值
- [ ] 接口类型使用 oneof（列出所有变体）
- [ ] 字段编号无冲突、有预留
- [ ] 命名符合 snake_case 规范
```

---

### 步骤 3：编写 Proto 文件（Implement）

**解决问题**：**"把设计写成 proto 代码"**

- 将步骤 2 的设计落地为 `.proto` 文件
- 添加注释保证可维护性
- 遵循 protobuf 规范和最佳实践

**防止问题**：proto 语法错误、可读性差、缺少文档

**目标**：实现完整的 .proto 定义

**操作**：
1. **文件头**：按照上方"Package 规范"编写
2. **Message 定义**：按照上方"命名规范"和"注释规范"编写
3. **嵌套结构独立定义**：每个复杂类型都是独立 message
4. **对于 oneof**：列举所有可能的变体

**输出检查清单**：
```markdown
- [ ] 文件头完整（syntax, package, go_package, imports）
- [ ] Message 总注释已添加
- [ ] 每个字段有注释（特别是枚举值、值域范围）
- [ ] 嵌套结构独立定义（不使用内联 message）
- [ ] 无循环依赖（import 关系是 DAG）
- [ ] 字段编号连续、无重复
- [ ] Proto 文件可读性好（分组、对齐、空行）
```

---

### 步骤 4：编译和验证（Verify）

**解决问题**：**"proto 定义是否正确可用？"**

- 编译检查语法和依赖正确性
- 确认生成的 Go 代码符合预期
- 验证字段完整性（数量、类型）

**防止问题**：错误累积到后续步骤、生成代码不可用、字段缺失未发现

**目标**：确保 proto 定义正确且可用

**操作**：
1. **编译生成 Go 代码**：使用上方"编译命令"
2. **检查编译错误**：语法错误、类型不匹配、Import 路径错误、重复字段编号
3. **检查生成的 Go 代码**：
   ```bash
   cat proto/gen/go/xxx/yyy.pb.go | grep "type.*struct"
   ```
4. **字段数量核对**：
   ```bash
   # Go 原始字段数
   grep -E "^\s+\w+\s+" sdk/xxx.go | wc -l
   # Proto 字段数
   grep -E "^\s+\w+.*=\s+\d+;" proto/xxx/yyy.proto | wc -l
   ```

**输出检查清单**：
```markdown
- [ ] protoc 编译无错误
- [ ] 生成的 .pb.go 文件存在
- [ ] 生成的 Go struct 字段数正确
- [ ] 导入路径可用（能 import 使用）
- [ ] Enum 生成了常量定义
- [ ] Oneof 生成了接口和实现
```

---

### 步骤 5：编写适配器骨架（Adapter）

**解决问题**：**"运行时怎么转换？"**

- 提前设计转换函数接口
- 识别复杂转换逻辑（需要特殊处理的地方）
- 为后续实际编写适配器代码做准备

**防止问题**：转换逻辑混乱、遗漏边界情况处理、适配器代码难以维护

**目标**：预先设计类型转换函数

**操作**：
1. **定义转换函数签名**：按照上方"适配器规范"
2. **识别复杂转换**：
   - 时间转换：`time.Time` ↔ `int64` (毫秒)
   - 枚举转换：`string` ↔ `enum`
   - 接口转换：`CardComp` ↔ `oneof`
   - 固定数组：`[4]T` ↔ `repeated T`

3. **标注转换逻辑**：
   ```go
   // ToProtoCard 转换 SDK Card 到 Proto Card
   // 特殊处理：
   // - Color: 王牌为空字符串
   // - Level: 当前级别牌
   func ToProtoCard(c *Card) *pb.Card {
       if c == nil {
           return nil
       }
       return &pb.Card{
           Number:    int32(c.Number),
           RawNumber: int32(c.RawNumber),
           Color:     c.Color,
           Level:     int32(c.Level),
           Name:      c.Name,
           DeckIndex: int32(c.DeckIndex),
       }
   }
   ```

**输出检查清单**：
```markdown
- [ ] 转换函数签名已定义（Go→Proto, Proto→Go）
- [ ] 复杂转换逻辑已标注（时间/枚举/接口/数组）
- [ ] Nil 处理已考虑
- [ ] 批量转换函数已定义（slice/map）
- [ ] 适配器代码位置已确定（文件路径）
```

---
