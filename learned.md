# 调试经验总结

## 1. Go 并发死锁问题

### 问题
同步事件处理 + 锁重入导致死锁：
```go
func (ge *GameEngine) StartDeal() error {
    ge.mutex.Lock()           // 持有写锁
    defer ge.mutex.Unlock()
    
    ge.emitEvent(event)       // 同步调用 handler
        -> handler()          // handler 调用 GetPlayerView()
            -> GetPlayerView() tries ge.mutex.RLock()  // 死锁！
}
```

### 解决方案
异步事件处理：
```go
func (ge *GameEngine) emitEvent(event *GameEvent) {
    for _, handler := range handlers {
        h := handler
        go func() {  // 异步执行，避免锁重入
            defer recover() // panic recovery
            h(event)
        }()
    }
}
```

### 教训
- **持有锁时避免同步调用外部代码**
- **事件处理器应该异步执行**
- **添加 panic recovery 保护稳定性**

---

## 2. React useEffect 闭包陷阱

### 问题
useEffect 依赖变化导致 handler 重新创建，捕获了旧的闭包值：
```typescript
useEffect(() => {
    const handlePlayerView = () => {
        if (currentPhase === 'GAME_PREPARE') {  // 使用闭包值
            // ... 永远不会进入，因为 currentPhase 是旧值
        }
    };
    wsClient.on('player_view', handlePlayerView);
}, [currentPhase, otherDeps]);  // currentPhase 变化 → handler 重新创建 → 捕获旧值
```

### 解决方案
从 store 获取最新值，而不是使用闭包：
```typescript
const handlePlayerView = () => {
    const latestPhase = useGameStore.getState().currentPhase;  // 获取最新值
    if (latestPhase === 'GAME_PREPARE') {
        // 正确使用最新值
    }
};
```

### 教训
- **useEffect 中的 handler 会捕获创建时的闭包值**
- **需要最新值时，从 store.getState() 获取**
- **减少不必要的 useEffect 依赖项**

---

## 3. 前后端数据结构对齐

### 问题
后端发送嵌套结构，前端期望扁平结构：
```javascript
// 后端实际结构
gameState = {
    current_match: {
        team_levels: [2, 2],
        current_deal: { level: 2 }
    }
}

// 前端错误访问
gameState.team_levels        // undefined ❌
gameState.current_deal       // undefined ❌
```

### 解决方案
前端适配后端的标准结构：
```javascript
const teamLevels = gameState.current_match?.team_levels || [2, 2];
const currentDeal = gameState.current_match?.current_deal;
```

### 教训
- **类型定义要与实际数据结构一致**
- **添加运行时日志验证数据结构**
- **使用可选链 + 默认值防御性编程**

---

## 4. 状态污染问题

### 问题
多个消息源更新同一个状态，导致状态被错误覆盖：
```typescript
// player_view 正确设置 gameState
setGameState(playerView.game_state);  // ✅

// game_event 错误覆盖
setGameState(message.data);  // ❌ message.data 是事件对象，不是游戏状态
```

### 解决方案
明确状态的单一数据源：
- **gameState 只从 player_view 更新**
- **game_event 只用于事件通知，不更新状态**

### 教训
- **为每个状态定义唯一的数据源**
- **避免多个地方更新同一个状态**
- **事件通知 ≠ 状态更新**

---

## 5. 调试方法论

### 有效的调试步骤
1. **查看完整日志** - 前端 console + 后端 docker logs
2. **对比预期流程** - 应该收到什么消息 vs 实际收到什么
3. **添加调试日志** - 在关键路径打印数据结构
4. **逐层排查** - 从后端 → WebSocket → 前端 store → 组件
5. **数据结构验证** - 展开对象查看实际内容

### 关键调试技巧
```typescript
// 打印实际数据结构
console.log('数据检查:', {
    hasField: !!obj.field,
    actualKeys: Object.keys(obj),
    fullObject: obj
});
```

### 教训
- **不要假设数据结构，要验证**
- **使用渐进式调试，逐步缩小范围**
- **保留调试日志直到问题完全解决**

---

## 6. 代码防御性编程

### 最佳实践
```typescript
// 数组操作前验证
const safeArray = Array.isArray(data) ? data : [];

// 可选链 + 默认值
const level = gameState.current_match?.current_deal?.level || 2;

// 数据完整性检查
if (!currentMatch || !currentDeal) {
    return <Loading />;
}
```

### 教训
- **永远不要假设数据一定存在**
- **为所有可选字段提供默认值**
- **在渲染前验证数据完整性**

---

## 总结

这次调试涉及：
- ✅ 后端并发死锁（Go 锁重入）
- ✅ 前端状态管理（React 闭包、状态污染）
- ✅ 数据结构对齐（类型不匹配）
- ✅ 防御性编程（空指针保护）

核心教训：**系统性思考 + 逐层验证 + 防御性编程**

