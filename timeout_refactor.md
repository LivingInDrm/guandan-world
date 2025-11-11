# 超时处理架构重构：统一到SDK GameDriver层

## 目标

将超时管理职责统一到SDK的GameDriver层，实现：
- Deal/Tribute/Trick层与时间概念解耦
- Backend Driver只负责I/O通信
- 统一超时事件通知和统计

---

## 架构调整

### 职责划分

```
Backend Driver 层
├─ 配置超时参数传递给SDK
├─ WebSocket消息收发
├─ 监听context取消信号（游戏停止/断线）
└─ 转发SDK事件到前端

SDK GameDriver 层
├─ 超时检测和触发
├─ 执行默认决策策略
├─ 发出EventPlayerTimeout事件
└─ 统计超时次数

Deal/Tribute/Trick 层
└─ 纯游戏逻辑，无时间概念
```

---

## 任务清单

### 任务1：SDK新增超时策略接口 ✅

**文件**: `sdk/timeout_strategy.go`（新建）

**内容**:
- 定义`TimeoutStrategy`接口
  - `GetDefaultPlayDecision(hand, trickInfo)` - 出牌超时默认策略
  - `GetDefaultTributeCard(options)` - 贡牌超时默认策略  
  - `GetDefaultReturnCard(hand)` - 还贡超时默认策略
- 实现`DefaultTimeoutStrategy`
  - Leader超时 → 出最小单牌
  - 非Leader超时 → PASS
  - 贡牌超时 → 选最大牌
  - 还贡超时 → 选最小牌

**依赖**: 无

**实施记录**:
- ✅ 已创建 `sdk/timeout_strategy.go`
- ✅ 已创建 `sdk/timeout_strategy_test.go` (19个测试用例，全部通过)
- ✅ 已通过code review并修复关键问题：
  - 修复了nil指针风险：改用安全的nil检查和跳过逻辑
  - 改进了接口文档：明确说明nil返回值情况
  - 添加了防御性编程：所有方法都能安全处理nil元素
- ✅ 编译验证通过
- ✅ 单元测试全部通过

---

### 任务2：GameDriver增强超时配置 ✅

**文件**: `sdk/game_driver.go`, `sdk/types.go`

**修改点**:

1. **扩展GameDriverConfig**
   - 新增字段: `TimeoutStrategy TimeoutStrategyInterface` (策略接口)
   - 保留现有: `PlayDecisionTimeout`, `TributeTimeout`

2. **新增GameDriver字段**
   - `timeoutStats map[int]*PlayerTimeoutStats` - 超时统计
   - `gameCancelCtx context.Context` - 游戏取消信号（传递给Backend）

**依赖**: 任务1

**实施记录**:
- ✅ 已在 `sdk/types.go` 中新增 `PlayerTimeoutStats` 类型定义
  - `PlayDecisionTimeouts` - 出牌决策超时次数
  - `TributeTimeouts` - 贡牌超时次数  
  - `TotalTimeouts` - 总超时次数
- ✅ 已修改 `GameDriverConfig` 添加 `TimeoutStrategy` 字段
- ✅ 已修改 `DefaultGameDriverConfig()` 初始化默认超时策略
- ✅ 已修改 `GameDriver` 添加超时管理字段：
  - `timeoutStats` - 超时统计map
  - `timeoutMu` - 读写锁保护并发访问
  - `gameCancelCtx` - 游戏取消上下文
  - `cancelFunc` - 取消函数
- ✅ 已修改 `NewGameDriver()` 初始化超时相关字段
  - 初始化 `timeoutStats` 为空map
  - 如果config未设置超时策略，使用默认策略
- ✅ 已新增线程安全的方法：
  - `GetTimeoutStats()` - 返回深拷贝，避免外部修改
  - `incrementPlayDecisionTimeout()` - 内部方法，线程安全增加出牌超时计数
  - `incrementTributeTimeout()` - 内部方法，线程安全增加贡牌超时计数
- ✅ 已创建 `game_driver_timeout_test.go` (13个测试用例，全部通过)
  - 包含并发安全测试和深拷贝验证
- ✅ 通过code review并修复关键问题：
  - 添加 `sync.RWMutex` 保护 `timeoutStats` 并发访问
  - 修改 `GetTimeoutStats()` 返回深拷贝防止外部修改
  - 添加内部helper方法用于线程安全更新统计
- ✅ 编译验证通过
- ✅ 单元测试全部通过（包含race detector验证）

---

### 任务3：GameDriver集成超时检测 ✅

**文件**: `sdk/game_driver.go`

**修改点**:

1. **RunMatch方法**
   - 创建 `gameCancelCtx, cancelFunc = context.WithCancel(context.Background())`
   - 保存到GameDriver字段
   - defer调用 `cancelFunc()` 确保游戏结束时取消

2. **runTrick方法改造**
   - 请求决策时：
     - 记录 `requestStartTime = time.Now()`
     - 启动超时检测goroutine
     - 调用 `inputProvider.RequestPlayDecision(gd.gameCancelCtx, ...)`
   - 超时检测goroutine:
     - 每1秒检查 `time.Since(requestStartTime) > config.PlayDecisionTimeout`
     - 检测到超时：
       - 调用 `handleTimeout(playerSeat, "play_decision")`
       - 使用 `TimeoutStrategy.GetDefaultPlayDecision()`
       - 继续游戏流程
   - 清理：defer关闭超时检测goroutine

3. **runTributePhase方法改造**
   - 贡牌选择时同样逻辑
   - 使用 `TimeoutStrategy.GetDefaultTributeCard()`
   - 还贡时使用 `TimeoutStrategy.GetDefaultReturnCard()`

4. **新增handleTimeout方法**
   - 发出 `EventPlayerTimeout` 事件
   - 记录超时统计到 `timeoutStats`
   - 日志记录

**依赖**: 任务2

**实施记录**:
- ✅ 已修改 `RunMatch()` 方法
  - 在开始时创建 `gameCancelCtx` 和 `cancelFunc`
  - 使用 defer 确保游戏结束时调用 `cancelFunc()`
  - 所有子操作基于 `gameCancelCtx`，实现级联取消
  - 添加 `cancelMu` 互斥锁保护 `cancelFunc` 的并发访问
- ✅ 已新增 `handleTimeout()` 方法
  - 根据 actionType 调用相应的统计增加方法
  - 创建并发送 `EventPlayerTimeout` 事件
  - 包含超时类型和玩家座位信息
  - 事件payload使用 `"action"` key 与GameEngine保持一致
- ✅ 已改造 `runTrick()` 方法
  - 使用 `context.WithTimeout(gd.gameCancelCtx, ...)` 创建带超时的上下文
  - **在 `cancel()` 之前捕获 `ctx.Err()`** 避免错误分类
  - 使用 `errors.Is()` 判断错误类型
  - 检测 `context.DeadlineExceeded` 判断超时
  - 超时时调用 `TimeoutStrategy.GetDefaultPlayDecision()` 生成默认决策
  - 检测 `context.Canceled` 判断游戏取消
  - 继续游戏流程，不中断
- ✅ 已改造 `runTributePhase()` 方法
  - 为 `TributeActionSelect` 添加超时处理
    - 超时时调用 `TimeoutStrategy.GetDefaultTributeCard()`
    - **添加fallback：如果策略返回nil，选择第一个非nil选项**
  - 为 `TributeActionReturn` 添加超时处理
    - 超时时调用 `TimeoutStrategy.GetDefaultReturnCard()`
    - **添加fallback：如果策略返回nil，选择第一个非nil选项**
  - 使用 `gameCancelCtx` 作为基础上下文
  - **在 `cancel()` 之前捕获 `ctx.Err()`** 避免错误分类
- ✅ 已新增 `CancelMatch()` 方法
  - 线程安全地取消当前正在运行的比赛
  - 使用 `cancelMu` 保护并发访问
- ✅ 已创建 `game_driver_task3_test.go` (4个测试用例，全部通过)
  - `TestHandleTimeout` - 验证超时处理和事件发送，包括事件payload格式
  - `TestRunMatch_InitializesGameCancelContext` - 验证上下文初始化
  - `TestRunMatch_CleansUpContext` - 验证上下文清理
  - `TestTimeoutStrategy_Integration` - 集成测试验证完整超时流程
- ✅ 通过code review并修复关键问题：
  - **修复错误分类问题**：在cancel()之前捕获ctx.Err()，使用errors.Is()判断
  - **修复nil处理**：为TimeoutStrategy返回nil时添加fallback逻辑
  - **修复并发安全**：添加cancelMu保护cancelFunc访问
  - **修复事件payload**：使用"action"而非"action_type"，移除重复的player_seat
  - **添加context重用保护**：RunMatch开始时检查并取消之前的context
- ✅ 编译验证通过
- ✅ 所有新测试通过
- ✅ 向后兼容：所有任务1、2的测试继续通过
- ✅ Race detector验证通过：无数据竞争

**设计改进**:
- 采用上下文级联取消机制：游戏取消会自动传播到所有待处理的请求
- 超时检测基于context而非独立goroutine，更简洁高效
- 明确区分超时(DeadlineExceeded)、取消(Canceled)和其他错误三种情况
- 超时时使用默认策略继续游戏，而取消时返回错误中止游戏
- 错误分类逻辑健壮：在cancel()之前捕获ctx.Err()，避免误判
- Nil安全：为TimeoutStrategy返回nil的情况提供fallback
- 线程安全：所有共享状态都有适当的并发保护

---

### 任务4：清理Deal/Tribute/Trick时间字段 ✅

**文件**: `sdk/types.go`, `sdk/deal.go`, `sdk/tribute.go`, `sdk/trick.go`

**修改点**:

1. **sdk/types.go**
   - 移除 `Trick.TurnTimeout time.Time`
   - 移除 `TributePhase.SelectTimeout time.Time`
   - 移除 `Deal.TurnTimeoutSecs int`
   - 删除 `TimeoutAction` 类型定义

2. **sdk/deal.go**
   - 删除 `GetTimeoutActions()` 方法
   - 移除所有设置 `TurnTimeout` 的代码
   - 移除 `NewDeal()` 中的 `TurnTimeoutSecs: 20` 初始化

3. **sdk/tribute.go**
   - 删除 `handleTimeout()` 方法
   - 移除所有设置 `SelectTimeout` 的代码

4. **sdk/trick.go**
   - 删除 `ProcessTimeout()` 方法

**依赖**: 任务3完成后才能删除

**实施记录**:
- ✅ **sdk/types.go 修改完成**
  - 移除 `Deal.TurnTimeoutSecs` 字段
  - 移除 `Trick.TurnTimeout` 字段
  - 移除 `TributePhase.SelectTimeout` 字段
  - 删除 `TimeoutAction` 类型定义
  - 为 `TributeAction.Timeout` 添加文档注释：`// DRIVER-MANAGED: Set by GameDriver, not by Tribute layer`
  
- ✅ **sdk/deal.go 修改完成**
  - 删除 `GetTimeoutActions()` 方法（214-270行）
  - 移除 `NewDeal()` 中的 `TurnTimeoutSecs: 20` 初始化
  - 移除 `PlayCards()` 中设置 `TurnTimeout` 的代码（157行）
  - 移除 `PassTurn()` 中设置 `TurnTimeout` 的代码（202行）
  - 移除 `startFirstTrick()` 中设置 `TurnTimeout` 的代码（369行）
  
- ✅ **sdk/tribute.go 修改完成**
  - 删除 `handleTimeout()` 方法（473-491行）
  - 删除 `HandleTimeoutAction()` 方法（657-707行，code review发现）
  - 移除 `NewTributePhase()` 中设置 `SelectTimeout` 的代码（54行）
  - 移除 `selectTribute()` 中设置 `SelectTimeout` 的代码（439行）
  - 移除 `ProcessTributePhaseAction()` 中设置 `Timeout` 字段（539, 555行）
  - 移除 `GetTributeStatusInfo()` 中设置 `Timeout` 字段（671, 684行）
  - 移除未使用的 `time` 包导入
  - 修复 IsImmune 硬编码问题：使用 `phase.IsImmune` 而非固定 `false`
  
- ✅ **sdk/trick.go 修改完成**
  - 删除 `ProcessTimeout()` 方法（185-196行）
  - 移除 `NewTrick()` 中的 `TurnTimeout: time.Time{}` 初始化
  
- ✅ **相关文件修复**
  - **sdk/game_engine.go**:
    - 移除 `GetPlayerGameState()` 中的超时计算逻辑（618-625行）
    - 将 `ProcessTimeouts()` 改为返回空数组并添加废弃注释（准备任务5删除）
    - 移除 `checkPendingActions()` 中设置 `TurnTimeout` 的代码（783行）
  - **sdk/deal_test.go**: 删除 `TestDealGetTimeoutActions` 测试
  - **sdk/trick_test.go**: 删除 `TestTrickProcessTimeout` 测试，移除 `time` 包导入
  - **sdk/game_flow_test.go**: 更新超时测试注释，移除 `time` 包导入
  
- ✅ **Code Review 修复**
  - **Critical问题**:
    - 移除 TributeManager.HandleTimeoutAction 方法（Tribute层不应有超时处理）
    - 移除所有 Tribute 层对 Timeout 字段的设置（4处）
    - 为 TributeAction.Timeout 添加文档说明其为 Driver 管理
  - **Minor问题**:
    - 修复 linter warning S1009：移除冗余的 nil 检查
    - 修复 IsImmune 硬编码：使用实际 phase.IsImmune 值
    - 移除未使用的 time 导入
    
- ✅ **验证结果**
  - 编译通过：`ok guandan-world/sdk 0.603s`
  - 所有 SDK linter warnings 清除
  - Deal/Tribute/Trick 层完全解耦时间概念
  - 所有超时逻辑已迁移到 GameDriver 层

**设计改进**:
- 层级职责更清晰：Deal/Tribute/Trick 专注纯游戏逻辑，完全无时间概念
- 代码质量提升：移除冗余代码，修复潜在bug
- 文档完善：关键字段添加使用说明，防止误用
- Tribute 层完全解耦超时：不设置、不处理、不管理超时

---

### 任务5：清理GameEngine的ProcessTimeouts ✅

**文件**: `sdk/game_engine.go`

**修改点**:

1. **删除ProcessTimeouts方法**（`game_engine.go:651-723`）
   - 该方法调用 `deal.GetTimeoutActions()` 已被移除
   - 功能已迁移到GameDriver层

2. **保留EventPlayerTimeout事件类型**
   - 由GameDriver的handleTimeout触发

**依赖**: 任务4

**实施记录**:
- ✅ **sdk/game_engine.go 修改完成**
  - 删除 `ProcessTimeouts()` 方法实现（643-649行）
  - 从 `GameEngineInterface` 接口删除 `ProcessTimeouts()` 定义（255-262行）
  - 删除接口中的 `// ProcessTimeouts 处理超时情况` 注释（254行）
  
- ✅ **sdk/game_engine_test.go 修改完成**
  - 删除 `TestGameEngineProcessTimeouts` 测试（230-242行）
  - 添加注释说明测试已移至 GameDriver
  
- ✅ **sdk/game_flow_test.go 修改完成**
  - 删除 `TestGameEngine_ProcessTimeouts` 测试函数（298-332行）
  - 添加注释指向 GameDriver 测试
  
- ✅ **验证结果**
  - 编译通过
  - 所有 GameEngine 相关测试通过：`ok guandan-world/sdk 0.412s`
  - ProcessTimeouts 完全从代码库移除（仅剩文档注释）
  - EventPlayerTimeout 事件类型保留，由 GameDriver 使用

**设计改进**:
- GameEngine 接口更简洁：移除已废弃的超时处理方法
- 职责更明确：GameEngine 不再负责超时检测和处理
- 超时逻辑完全统一到 GameDriver 层
- 测试覆盖率保持：超时测试已在 game_driver_task3_test.go

---

### 任务6：Backend简化RoomInputProvider ✅

**文件**: `backend/game/driver_service.go`

**修改点**:

1. **RequestPlayDecision简化**
   
   保留:
   - `ctx context.Context` 参数（SDK传入的gameCancelCtx）
   - 创建 `decisionChan`
   - 发送 WebSocket 消息
   - `select` 等待 `decisionChan` 或 `ctx.Done()`
   - `ctx.Done()` 时返回 `ctx.Err()`
   
   移除:
   - `if ctx == nil { ctx, cancel = context.WithTimeout(...) }` 逻辑
   - `validateAndFixPlayDecision()` 方法
   - `generateDefaultPlayDecision()` 方法
   - 所有默认决策相关逻辑

2. **RequestTributeSelection简化**
   
   保留:
   - `ctx` 参数和 `ctx.Done()` 监听
   - 基础收发逻辑
   
   移除:
   - 自己创建WithTimeout context的代码
   - 超时自动选最大牌逻辑

3. **RequestReturnTribute简化**
   
   保留:
   - `ctx` 参数和 `ctx.Done()` 监听
   - 基础收发逻辑
   
   移除:
   - 自己创建WithTimeout context的代码
   - 超时自动选最小牌逻辑

4. **删除辅助方法**
   - `validateAndFixPlayDecision()`
   - `generateDefaultPlayDecision()`
   - `findSmallestCard()`

**依赖**: 任务3完成确保SDK能处理超时

**实施记录**:
- ✅ **RequestPlayDecision 简化完成** (`backend/game/driver_service.go:299-342`)
  - 移除nil context检查和WithTimeout创建
  - 移除nil decision检查和默认策略生成
  - 移除validateAndFixPlayDecision调用
  - ctx.Done()时直接返回ctx.Err()
  - 添加nil context防御检查
  - 使用双值接收处理channel关闭
  - 从ctx.Deadline()动态计算超时时间
  
- ✅ **RequestTributeSelection 简化完成** (`backend/game/driver_service.go:345-404`)
  - 移除WithTimeout创建逻辑
  - 移除超时自动选最大牌逻辑
  - 添加nil context防御检查
  - 使用双值接收处理channel关闭
  - 从ctx.Deadline()动态计算超时时间
  
- ✅ **RequestReturnTribute 简化完成** (`backend/game/driver_service.go:407-472`)
  - 移除WithTimeout创建逻辑
  - 移除超时自动选最小牌逻辑
  - 添加nil context防御检查
  - 使用双值接收处理channel关闭
  - 从ctx.Deadline()动态计算超时时间
  
- ✅ **辅助方法删除完成**
  - 删除 validateAndFixPlayDecision() 方法（原581-652行）
  - 删除 generateDefaultPlayDecision() 方法（原514-526行）
  - 删除 findSmallestCard() 方法（原500-512行）
  
- ✅ **SDK GameDriver增强** (`sdk/game_driver.go`)
  - 在runTrick中添加nil decision检查和fallback
  - 使用TimeoutStrategy生成默认决策
  
- ✅ **新增测试** (`backend/game/driver_service_task6_test.go`)
  - TestRequestPlayDecision_NilContext - 验证nil context检查
  - TestRequestTributeSelection_NilContext
  - TestRequestReturnTribute_NilContext
  - TestSubmitTributeSelection_NilCard - 验证nil card检查
  - TestSubmitReturnTribute_NilCard
  - TestRequestPlayDecision_CanceledChannel - 验证channel关闭处理
  
- ✅ **新增测试** (`backend/game/driver_service_timeout_test.go`)
  - TestRequestPlayDecision_TimeoutInMessage - 验证超时从context获取
  - TestRequestPlayDecision_NoDeadline - 验证无deadline情况
  
- ✅ **Code Review修复** (详见TASK6_CODE_REVIEW_FIXES.md)
  - **Critical修复1**: 使用双值接收检测channel关闭，返回明确错误
  - **Critical修复2**: 从ctx.Deadline()动态获取超时，保持前后端一致
  - **Critical修复3**: 添加安全警告注释标记sendToPlayer广播问题
  - **Major修复4**: 添加nil decision/card验证
  - **Major修复5**: 添加nil context防御检查，防止无限阻塞
  - **Minor修复7**: SubmitTributeSelection/SubmitReturnTribute添加nil检查
  
- ✅ **MockDriverWSManager并发修复** (`backend/game/driver_service_test.go`)
  - 添加sync.RWMutex保护broadcasts/messages map
  - 修复concurrent map writes问题
  
- ✅ **验证结果**
  - 编译通过：`go build ./backend/game/...`
  - 所有测试通过：`ok guandan-world/backend/game 0.684s`
  - Race detector通过：无数据竞争
  - 新增8个测试全部通过

**设计改进**:
- RoomInputProvider职责更清晰：纯I/O通信，不处理业务逻辑
- 超时决策完全由SDK GameDriver管理
- 错误处理更健壮：区分channel关闭、超时、取消三种情况
- 超时配置统一：从context动态获取，避免前后端不一致
- 防御性编程：添加nil检查，避免静默失败

---

### 任务7：DriverService配置超时策略 ✅

**文件**: `backend/game/driver_service.go`

**修改点**:

1. **StartGameWithDriver方法**
   ```
   创建GameDriver时：
   - 配置 TimeoutStrategy = DefaultTimeoutStrategy
   - 配置 PlayDecisionTimeout = 10秒（测试）/ 30秒（生产）
   - 配置 TributeTimeout = 10秒（测试）/ 20秒（生产）
   ```

2. **监听超时事件**
   - WebSocketObserver已监听所有事件
   - EventPlayerTimeout会自动广播到前端
   - 无需额外处理

**依赖**: 任务6

**实施记录**:
- ✅ **StartGameWithDriver配置完成** (`backend/game/driver_service.go:64-84`)
  - 配置 `config.TimeoutStrategy = sdk.NewDefaultTimeoutStrategy()`
  - 根据环境变量APP_ENV配置超时时间：
    - test/dev环境: 10秒（快速迭代）
    - prod环境: 使用DefaultGameDriverConfig默认值（30秒/20秒）
  - 添加详细注释说明生产环境建议配置
  
- ✅ **环境检测函数** (`backend/game/driver_service.go:14-21`)
  - 新增 `getEnvironment()` 函数
  - 读取APP_ENV环境变量
  - 默认为"prod"确保生产安全
  
- ✅ **超时计算helper** (`backend/game/driver_service.go:23-35`)
  - 新增 `getRemainingTimeout()` 函数
  - 从context.Deadline()计算剩余秒数
  - 负数clamp为0，防止前端显示错误
  - 统一三个Request方法使用
  
- ✅ **超时事件自动广播验证**
  - WebSocketObserver.OnGameEvent已实现通用事件处理
  - EventPlayerTimeout自动转换为WebSocket消息
  - 通过BroadcastToRoom广播到房间
  - 无需额外代码
  
- ✅ **新增测试** (`backend/game/driver_service_task7_test.go`)
  - TestWebSocketObserver_PlayerTimeoutEvent - 验证超时事件广播
  - TestDriverService_TimeoutStrategyConfigured - 验证策略配置
  
- ✅ **新增测试** (`backend/game/driver_service_review_fixes_test.go`)
  - TestGetEnvironment (4个子测试) - 验证环境检测
  - TestGetRemainingTimeout (3个子测试) - 验证超时计算
  - TestGetRemainingTimeout_PastDeadline - 验证过期deadline处理
  - TestRequestPlayDecision_NegativeTimeout - 验证无负数超时
  - TestDriverService_ProductionTimeouts - 验证生产环境配置
  - TestDriverService_TestEnvironmentTimeouts - 验证测试环境配置
  
- ✅ **Code Review修复** (详见TASK7_CODE_REVIEW_FIXES.md)
  - **Critical修复1**: 通过APP_ENV环境变量区分环境，生产使用30s/20s
  - **Major修复2**: getRemainingTimeout将负数clamp为0
  - **Minor修复3**: 提取重复代码为helper函数
  
- ✅ **验证结果**
  - 编译通过：`go build ./backend/game/...`
  - 所有测试通过：`ok guandan-world/backend/game 0.851s`
  - Race detector通过：无数据竞争
  - 新增11个测试全部通过（任务7: 2个，code review: 9个）

**事件数据格式**:
```json
{
  "type": "game_event",
  "data": {
    "event_type": "player_timeout",
    "event_data": {
      "action": "play_decision" | "tribute_select" | "return_tribute"
    },
    "timestamp": "2025-11-11T...",
    "player_seat": 0-3
  }
}
```

**环境配置**:
```bash
# 开发/测试环境
export APP_ENV=test  # 或 dev，使用10秒超时

# 生产环境
export APP_ENV=prod  # 或不设置（默认prod），使用30秒/20秒超时
```

**设计改进**:
- 超时策略配置清晰明确
- 环境感知配置：测试快速迭代，生产充足思考时间
- 超时时间动态计算：避免前后端配置不一致
- 代码复用：提取helper函数减少重复
- 防御性编程：负数超时clamp为0

---

### 任务8：更新WebSocketObserver ✅

**文件**: `backend/game/driver_service.go`

**修改点**:

1. **WebSocketObserver.OnGameEvent**
   - 确认 `EventPlayerTimeout` 在switch中被处理
   - 事件数据包含：
     - `player_seat`
     - `action_type` (play_decision/tribute_select/return_tribute)
     - `timestamp`

**依赖**: 任务3

**实施记录**:
- ✅ 已在 `driver_service.go` 的 `OnGameEvent` 方法中添加 `EventPlayerTimeout` 到日志switch
  - 超时事件现在会被记录到日志中，便于监控和调试
  - 事件通过现有的通用逻辑正确广播到WebSocket客户端
- ✅ 验证事件数据格式符合要求：
  - `player_seat`：通过 `event.PlayerSeat` 字段传递
  - `action`：在 `event.Data` 中（注：SDK使用"action"而非"action_type"，与任务3保持一致）
  - `timestamp`：通过 `event.Timestamp` 字段传递
- ✅ WebSocket消息结构：
  ```json
  {
    "type": "game_event",
    "data": {
      "event_type": "player_timeout",
      "event_data": {"action": "play_decision|tribute_select|return_tribute"},
      "player_seat": <seat>,
      "timestamp": <time>
    }
  }
  ```
- ✅ 已创建 `driver_service_task8_test.go`（3个测试用例，全部通过）
  - `TestWebSocketObserver_PlayerTimeoutEvent_AllActionTypes`：验证所有三种超时类型的事件处理
  - `TestWebSocketObserver_PlayerTimeoutEvent_LoggingEnabled`：验证超时事件记录到日志
  - 测试覆盖play_decision、tribute_select、return_tribute三种超时场景
- ✅ 编译验证通过
- ✅ 所有测试通过（包括任务7的现有测试）
- ✅ 向后兼容：所有现有backend/game测试继续通过

**设计说明**:
- 超时事件通过通用事件广播机制处理，无需特殊逻辑
- 明确在日志switch中添加EventPlayerTimeout，确保重要的超时事件被记录
- 前端可以接收并显示"玩家X超时，自动出牌"等提示信息

---

## 架构流程

### 超时处理流程（重构后）

```
1. GameDriver.RunMatch()
   ├─ 创建 gameCancelCtx (WithCancel)
   │
2. GameDriver.runTrick()
   ├─ 记录 requestStartTime
   ├─ 启动超时检测goroutine
   │  └─ 每1秒检查: time.Since(requestStartTime) > PlayDecisionTimeout
   │
   ├─ 调用 inputProvider.RequestPlayDecision(gameCancelCtx, ...)
   │
   ├─ 超时触发
   │  ├─ handleTimeout(playerSeat, "play_decision")
   │  │  ├─ 发出 EventPlayerTimeout
   │  │  └─ 记录 timeoutStats
   │  ├─ decision = TimeoutStrategy.GetDefaultPlayDecision()
   │  └─ 执行 engine.PlayCards(decision)
   │
3. Backend RoomInputProvider.RequestPlayDecision(ctx)
   ├─ 发送 WebSocket 消息
   ├─ select {
   │    case decision := <-decisionChan:
   │        return decision, nil
   │    case <-ctx.Done():
   │        return nil, ctx.Err()  // 游戏取消/断线
   │  }
```

---

## 关键设计点

### 1. 超时检测机制

**GameDriver主动检测**:
- 请求决策前记录 `requestStartTime`
- 启动独立goroutine每1秒检查
- 检测到超时立即执行默认策略
- 不依赖Backend的context timeout

### 2. Context用途变更

**之前**: `context.WithTimeout` - 超时控制  
**现在**: `context.WithCancel` - 取消信号

**取消场景**:
- 游戏结束（RunMatch返回）
- 用户主动停止游戏
- 玩家断线（可选，通过provider.CancelAll）

### 3. 时间字段迁移

**从Deal层完全移除**:
- Deal/Trick/Tribute不再存储任何时间字段
- GameDriver在请求决策时临时记录开始时间
- 超时检测在GameDriver的goroutine中进行

### 4. 事件通知

**EventPlayerTimeout开始生效**:
- SDK定义已存在但未触发
- 重构后GameDriver会真正发出此事件
- 前端可显示"玩家X超时，自动出牌"提示

---

## 文件改动清单

### SDK新增
- `sdk/timeout_strategy.go` - 超时策略接口和默认实现

### SDK修改
- `sdk/game_driver.go` - 核心超时管理逻辑
- `sdk/types.go` - 移除时间字段
- `sdk/deal.go` - 移除GetTimeoutActions和超时设置
- `sdk/tribute.go` - 移除handleTimeout
- `sdk/trick.go` - 移除ProcessTimeout
- `sdk/game_engine.go` - 删除ProcessTimeouts方法

### Backend修改
- `backend/game/driver_service.go` - 简化RoomInputProvider，配置超时策略

---

## 依赖关系

```
任务1 (超时策略接口)
  ↓
任务2 (GameDriver配置扩展)
  ↓
任务3 (GameDriver集成超时检测) ← 核心任务
  ↓
任务4 (清理Deal/Tribute时间字段)
任务5 (清理GameEngine)
  ↓
任务6 (Backend简化) ← 关键简化
  ↓
任务7 (DriverService配置)
任务8 (WebSocketObserver)
```

---

## 关键约束

1. **向后兼容事件**: 保持所有现有GameEvent类型不变
2. **前端事件变化**: WebSocket消息格式不变，EventPlayerTimeout开始真正触发（之前未生效）
3. **性能要求**: 超时检测goroutine开销可控（1秒间隔）
4. **并发安全**: GameDriver的超时统计需要mutex保护
5. **Context语义**: Backend的context仅用于取消，不含timeout