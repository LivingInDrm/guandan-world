# 开发规范

本文档定义了 guandan-world 代码库的开发规则和最佳实践。

---

## 日志规范

### 日志包

使用项目统一的日志包：

```go
import "guandan-world/pkg/log"
```

**不要使用** 标准库 `log` 包，避免命名冲突。

### 日志级别

| 级别 | 用途 | 示例 |
|-----|------|-----|
| `Debug` | 开发调试信息 | 状态转换、变量值 |
| `Info` | 关键业务事件 | 游戏开始、比赛结束 |
| `Warn` | 可恢复异常 | 超时重试、降级处理 |
| `Error` | 业务错误 | 非法操作、处理失败 |

### 使用方式

```go
log.Info("game started", "match_id", matchID, "player_count", 4)
log.Error("invalid move", "player_id", playerID, "error", err)
```

### 参数规范

- 使用 key-value 形式，key 使用 snake_case，参数必须成对出现；常用key：`match_id`, `room_id`, `player_id`

### 初始化

日志系统已经在 `backend/main.go` 中初始化，其他模块无需关心，直接使用即可。

### 日志输出

- 格式：`MM-DDThh:mm:ss.sss LEVEL [package.Type] [Method] message key=value`
- 文件：`./logs/yyyy-mm-dd.log`
- 示例：`12-01T14:32:05.123 INFO  [sdk.GameEngine] [StartMatch] game started match_id=abc player_count=4`

---

## 异常检测规范

### 分层检测原则（Trust Boundary）

根据函数所在层次决定检测职责：

| 层次 | 定义 | 检测职责 |
|------|------|---------|
| **Public API** | 接口定义的方法、对外暴露 | 检测所有参数 + 状态前提 |
| **Internal** | 仅被同包 Public API 调用 | 信任调用者，仅检测自身新增逻辑 |
| **Pure Function** | 无副作用的纯函数 | 通过类型约束，不做 nil 检测 |

### 检测函数命名规范

| 前缀 | 层次 | 失败行为 | 使用场景 |
|------|------|---------|---------|
| `require` | Public API | `log.Warn` + 返回 `error` | 外部调用，失败 = 业务错误 |
| `must` | Internal | `log.Error` + `panic` | 调用者已验证，失败 = 编程错误 |

```go
// Public API - log.Warn + 返回 error
func (ge *GameEngine) requireActiveDeal() (*Deal, error) {
    if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
        log.Warn("require failed: no active deal")
        return nil, errors.New("no active deal")
    }
    return ge.currentMatch.CurrentDeal, nil
}

// Internal - panic + log
func (ge *GameEngine) mustActiveDeal() *Deal {
    if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
        log.Error("must failed: no active deal")
        panic("must: no active deal")
    }
    return ge.currentMatch.CurrentDeal
}
```

### 使用示例

```go
// Public API：使用 require
func (ge *GameEngine) PlayCards(...) (*GameEvent, error) {
    deal, err := ge.requireActiveDeal()
    if err != nil { return nil, err }
    // ...
}

// Internal：使用 must（调用者已验证）
func (ge *GameEngine) checkPreActionStateTransitions() []*GameEvent {
    deal := ge.mustActiveDeal()
    // ...
}
```
