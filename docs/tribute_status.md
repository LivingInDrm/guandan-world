## 双下场景处理逻辑详解

### 前提条件
- **Winners** = `[rank1, rank2]` (同队，如 seat 0 和 seat 2)
- **TributePairs** = `[{Giver: rank3}, {Giver: rank4}]` (如 seat 1 和 seat 3)
- 初始状态: `TributeStatus = Waiting`, `PoolCards = []`

---

### 调用栈
```
GameDriver.runTributePhase()           ← 唯一循环点
  └── GameEngine.StepTribute(input)    ← 单步处理
        └── ProcessTributeStep()       ← 纯函数计算
              ├── processTributeWaiting()
              ├── processTributeSelecting()
              ├── processTributeReturning()
              └── processTributeFinished()
```

---

### 逐步执行流程

| Step | 函数调用 | 状态 | Input | 操作 | 事件 | StatusChanged | 延迟 |
|:----:|----------|:----:|:-----:|------|------|:-------------:|:----:|
| **1** | `processTributeWaiting` | Waiting | nil | rank3、rank4 各选最大牌入池 | 2× `TributeCardSubmitted` | **true** (→Selecting) | **WaitingDelay** |
| **2** | `processTributeSelecting` | Selecting | nil | PoolCards=2，需 rank1 选择 | - | false | - |
| | | | | **返回 Action: Select(rank1)** | | | |
| **3** | `processTributeSelecting` | Selecting | rank1选牌 | rank1 获牌，池剩 1 张 | 1× `TributeCardSelected` | false | - |
| **4** | `processTributeSelecting` | Selecting | nil | PoolCards=1，自动分配给 rank2 | 1× `TributeCardSelected` | false | - |
| **5** | `processTributeSelecting` | Selecting | nil | PoolCards=0 | - | **true** (→Returning) | **SelectingDelay** |
| **6** | `processTributeReturning` | Returning | nil | 需 rank1 还贡 | - | false | - |
| | | | | **返回 Action: Return(rank1→giver)** | | | |
| **7** | `processTributeReturning` | Returning | rank1还牌 | rank1 → giver | 1× `ReturnTribute` | false | - |
| **8** | `processTributeReturning` | Returning | nil | 需 rank2 还贡 | - | false | - |
| | | | | **返回 Action: Return(rank2→giver)** | | | |
| **9** | `processTributeReturning` | Returning | rank2还牌 | rank2 → giver | 1× `ReturnTribute` | false | - |
| **10** | `processTributeReturning` | Returning | nil | 全部完成 | - | **true** (→Finished) | **ReturningDelay** |
| **11** | `processTributeFinished` | Finished | nil | 验证完整性 | 1× `TributeCompleted` | false | - |
| | | | | **Completed = true** | | | **FinishedDelay** |

---

### 延迟触发点（共 4 次）

```
Step 1:  Waiting → Selecting     ══▶  TributeWaitingDelay   (默认 2s)
Step 5:  Selecting → Returning   ══▶  TributeSelectingDelay (默认 2s)
Step 10: Returning → Finished    ══▶  TributeReturningDelay (默认 2s)
Step 11: Completed = true        ══▶  TributeFinishedDelay  (默认 2s)
```

---

### Driver 循环逻辑（伪代码）

```go
for step := 0; step < 20; step++ {
    result := engine.StepTribute(pendingInput)
    pendingInput = nil
    
    // 1. 状态转换 → 延迟
    if result.StatusChanged {
        sleep(getTributeDelay(result.PrevStatus))  // Waiting/Selecting/Returning
    }
    
    // 2. 完成 → 延迟并退出
    if result.Completed {
        sleep(TributeFinishedDelay)
        break
    }
    
    // 3. 需要用户输入 → 获取输入，下轮带入
    if result.Action != nil {
        pendingInput = getTributeInput(result.Action)
    }
    
    // 4. 无 Action、无 Completed → 自动继续下一步
}
```

---

### 状态转换图

```
                    ┌──────────────────────────────────────────────────────┐
                    │                                                      │
  ┌─────────┐  Step1  ┌───────────┐  Step5  ┌───────────┐  Step10  ┌──────────┐
  │ Waiting │ ──────▶ │ Selecting │ ──────▶ │ Returning │ ───────▶ │ Finished │
  └─────────┘         └───────────┘         └───────────┘          └──────────┘
       │                    │                     │                      │
       │               Step2,3,4              Step6,7,8,9           Step11
       │              (循环处理)              (循环处理)           (验证完成)
       │                    │                     │                      │
       ▼                    ▼                     ▼                      ▼
   WaitingDelay        (无延迟)              (无延迟)             FinishedDelay
                      SelectingDelay        ReturningDelay
                      (状态转换时)          (状态转换时)
```

---

### 总结

| 维度 | 数值 |
|------|------|
| 总步数 | 11 步 |
| 用户输入次数 | 4 次 (Step 2, 6, 8 返回 Action) |
| 状态转换次数 | 3 次 (Waiting→Selecting, Selecting→Returning, Returning→Finished) |
| 延迟触发次数 | 4 次 (3 次状态转换 + 1 次完成) |
| 事件总数 | 7 个 (2×Submitted + 2×Selected + 2×Return + 1×Completed) |