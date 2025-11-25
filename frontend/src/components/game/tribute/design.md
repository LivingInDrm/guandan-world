# 上贡组件重构实施方案

## 目标
重构上贡阶段组件，实现统一的6阶段交互流程，支持抗贡、单下、末游、双下所有场景。

## 核心设计原则
- **状态驱动**：完全依赖后端 `player_view` 和事件推送
- **配置优先**：用数据配置代替大量组件
- **前端缓冲**：用定时器控制展示节奏，同时监听后端状态跳跃
- **复用优先**：复用现有 `PlayerHand`、`CardDisplay` 组件

---

## 文件结构

```
frontend/src/components/game/tribute/
├── TributeFlow.tsx              (主容器，250行)
├── TributePhaseContent.tsx      (内容路由，200行)
├── SelectionPanel.tsx           (选贡面板，150行)
├── ReturnPanel.tsx              (还贡面板，200行)
├── ProgressIndicator.tsx        (进度条，80行)
├── phaseConfigs.tsx             (阶段配置，100行)
└── types.ts                     (类型定义，50行)

复用现有组件：
- PlayerHand (已存在)
- CardDisplay (已存在)
```

---

## 阶段 1：创建类型定义 (types.ts)

### 文件：`frontend/src/components/game/tribute/types.ts`

**内容**：
- 定义 `UIPhase` 枚举：`'START' | 'IMMUNITY_CHECK' | 'SUBMITTING' | 'SELECTING' | 'RETURNING' | 'FINISHED'`
- 定义 `PhaseConfig` 接口：`{ title, icon?, content, duration? }`
- 定义 `ReturnTask` 接口：`{ receiver, giver, done }`
- 导出所有共享类型

**验证点**：
- TypeScript 编译通过
- 类型覆盖所有阶段和场景

---

## 阶段 2：实现阶段配置 (phaseConfigs.tsx)

### 文件：`frontend/src/components/game/tribute/phaseConfigs.tsx`

**内容**：
- 导出 `getPhaseConfig(phase, tributePhase)` 函数
- 返回每个阶段的配置对象：
  - `title`: 标题文本
  - `icon`: 图标（可选）
  - `renderContent`: 渲染函数
  - `duration`: 自动切换时长（毫秒）

**配置示例**：
```typescript
{
  START: {
    title: '上贡阶段开始',
    icon: '🎴',
    duration: 3000,
    renderContent: (props) => <StartContent {...props} />
  },
  IMMUNITY_CHECK: {
    title: '抗贡检测',
    icon: '🛡️',
    duration: (isImmune) => isImmune ? 3000 : 2000,
    renderContent: (props) => <ImmunityCheckContent {...props} />
  },
  // ... 其他阶段
}
```

**验证点**：
- 所有6个阶段都有配置
- 配置可根据运行时数据动态调整

---

## 阶段 3：实现主容器 (TributeFlow.tsx)

### 文件：`frontend/src/components/game/tribute/TributeFlow.tsx`

**职责**：
1. 管理 `uiPhase` 状态
2. 处理阶段自动切换（定时器）
3. 监听后端状态变化（状态跳跃）
4. 路由到内容组件

**核心逻辑**：

### 状态管理
```typescript
const [uiPhase, setUIPhase] = useState<UIPhase>('START');
```

### 自动切换逻辑（阶段0-1）
```typescript
useEffect(() => {
  if (uiPhase === 'START') {
    const timer = setTimeout(() => setUIPhase('IMMUNITY_CHECK'), 3000);
    return () => clearTimeout(timer);
  }
}, [uiPhase]);

useEffect(() => {
  if (uiPhase === 'IMMUNITY_CHECK') {
    const delay = tributePhase.is_immune ? 3000 : 2000;
    const timer = setTimeout(() => {
      setUIPhase(tributePhase.is_immune ? 'FINISHED' : 'SUBMITTING');
    }, delay);
    return () => clearTimeout(timer);
  }
}, [uiPhase, tributePhase.is_immune]);
```

### 状态跳跃监听（关键）
```typescript
// 监听后端已到 RETURNING，立即跳过动画
useEffect(() => {
  if (tributePhase.status === 'RETURNING' && 
      uiPhase !== 'RETURNING' && 
      uiPhase !== 'FINISHED') {
    setUIPhase('RETURNING');
  }
}, [tributePhase.status, uiPhase]);

// 监听后端已到 FINISHED
useEffect(() => {
  if (tributePhase.status === 'FINISHED' && uiPhase !== 'FINISHED') {
    setUIPhase('FINISHED');
  }
}, [tributePhase.status, uiPhase]);
```

**Props 接口**：
```typescript
interface TributeFlowProps {
  tributePhase: TributePhase;
  players: Player[];
  currentPlayerSeat: number;
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
  onReturnTribute: (deckIndex: number) => void;
}
```

**验证点**：
- 阶段按预期自动切换
- 后端快速推进时能正确跳过中间阶段
- 所有 props 正确传递给子组件

---

## 阶段 4：实现内容路由 (TributePhaseContent.tsx)

### 文件：`frontend/src/components/game/tribute/TributePhaseContent.tsx`

**职责**：
- 接收 `phase` 和所有 props
- 根据阶段渲染对应内容
- 使用配置驱动，减少条件分支

**核心结构**：
```typescript
const TributePhaseContent = ({ phase, tributePhase, players, ...props }) => {
  const config = getPhaseConfig(phase, tributePhase);
  
  return (
    <div className="tribute-phase-content">
      {/* 头部 */}
      <div className="phase-header">
        {config.icon && <span className="icon">{config.icon}</span>}
        <h2>{config.title}</h2>
      </div>
      
      {/* 主体内容 */}
      <div className="phase-body">
        {renderPhaseContent(phase, tributePhase, players, props)}
      </div>
    </div>
  );
};
```

### 内容渲染函数
```typescript
const renderPhaseContent = (phase, tributePhase, players, props) => {
  switch (phase) {
    case 'START':
      return <StartContent tributePhase={tributePhase} players={players} />;
    
    case 'IMMUNITY_CHECK':
      return tributePhase.is_immune 
        ? <ImmuneSuccessContent tributePhase={tributePhase} />
        : <ImmuneFailContent />;
    
    case 'SUBMITTING':
      return <SubmittingContent poolCards={tributePhase.pool_cards} />;
    
    case 'SELECTING':
      return <SelectionPanel {...props} tributePhase={tributePhase} />;
    
    case 'RETURNING':
      return <ReturnPanel {...props} tributePhase={tributePhase} />;
    
    case 'FINISHED':
      return <ResultContent tributePhase={tributePhase} players={players} />;
  }
};
```

### 内联简单展示组件
- `StartContent`: 显示上贡规则和参与角色
- `ImmuneSuccessContent`: 显示抗贡成功信息
- `ImmuneFailContent`: 显示未触发抗贡
- `SubmittingContent`: 显示贡牌池和提示
- `ResultContent`: 显示最终结果

**验证点**：
- 每个阶段正确渲染
- 数据正确传递给展示组件
- 样式统一美观

---

## 阶段 5：实现选贡面板 (SelectionPanel.tsx)

### 文件：`frontend/src/components/game/tribute/SelectionPanel.tsx`

**职责**：处理阶段3的所有场景（单下自动、双下选择、等待）

**Props 接口**：
```typescript
interface SelectionPanelProps {
  tributePhase: TributePhase;
  currentPlayerSeat: number;
  players: Player[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
}
```

**核心逻辑**：
```typescript
const SelectionPanel = ({ tributePhase, currentPlayerSeat, onSelectTribute }) => {
  const { pool_cards, selecting_player } = tributePhase;
  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  
  const isDoubleDown = pool_cards.length > 1;
  const isMyTurn = selecting_player === currentPlayerSeat;
  
  // 场景1: 单下自动分配
  if (!isDoubleDown) {
    return (
      <div className="auto-select">
        <p>rank1 自动获得贡牌</p>
        {/* 显示动画 */}
      </div>
    );
  }
  
  // 场景2: 双下，我选择
  if (isMyTurn) {
    return (
      <div className="manual-select">
        <p>请选择一张贡牌</p>
        <div className="pool-cards">
          {pool_cards.map(card => (
            <div 
              key={card.id}
              className={selectedCard === card.deckIndex ? 'selected' : ''}
              onClick={() => setSelectedCard(card.deckIndex)}
            >
              <CardDisplay card={card} />
            </div>
          ))}
        </div>
        <button 
          disabled={selectedCard === null}
          onClick={() => {
            onSelectTribute(selectedCard);
            setSelectedCard(null);
          }}
        >
          确认选择
        </button>
      </div>
    );
  }
  
  // 场景3: 双下，等待
  return (
    <div className="waiting-select">
      <p>等待 {getPlayerName(selecting_player)} 选择贡牌...</p>
      <div className="pool-cards disabled">
        {pool_cards.map(card => (
          <CardDisplay key={card.id} card={card} />
        ))}
      </div>
    </div>
  );
};
```

**验证点**：
- 单下场景显示自动分配提示
- 双下场景我方可选择，他方显示等待
- 选择后正确调用 `onSelectTribute`
- 复用 `CardDisplay` 组件显示贡池卡牌

---

## 阶段 6：实现还贡面板 (ReturnPanel.tsx)

### 文件：`frontend/src/components/game/tribute/ReturnPanel.tsx`

**职责**：统一布局，上半部分显示进度，下半部分根据角色显示操作区

**Props 接口**：
```typescript
interface ReturnPanelProps {
  tributePhase: TributePhase;
  currentPlayerSeat: number;
  players: Player[];
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onReturnTribute: (deckIndex: number) => void;
}
```

**核心结构**：
```typescript
const ReturnPanel = ({ tributePhase, currentPlayerSeat, playerHand, selectedCards, onCardSelect, onReturnTribute, players }) => {
  // 解析还贡任务
  const tasks = tributePhase.tribute_pairs
    .filter(p => p.receiver !== -1)
    .map(p => ({
      receiver: p.receiver,
      giver: p.giver,
      done: !!p.return_card
    }));
  
  const myTask = tasks.find(t => t.receiver === currentPlayerSeat);
  const needReturn = myTask && !myTask.done;
  
  return (
    <div className="return-panel">
      {/* 上半部分：进度 */}
      <div className="progress-section">
        <h3>还贡进度</h3>
        {tasks.map((task, i) => (
          <div key={i} className={task.done ? 'task-done' : 'task-pending'}>
            <span className="status-icon">{task.done ? '✓' : '⏳'}</span>
            <span className="task-desc">
              {getPlayerName(task.receiver)} → {getPlayerName(task.giver)}
              {task.done ? ' 已还贡' : ' 还贡中...'}
            </span>
          </div>
        ))}
      </div>
      
      {/* 下半部分：操作区 */}
      <div className="action-section">
        {needReturn ? (
          <div className="return-selector">
            <p>请选择一张牌还给 {getPlayerName(myTask.giver)}</p>
            
            {/* 复用 PlayerHand 组件 */}
            <PlayerHand 
              cards={playerHand}
              selectedCards={selectedCards}
              onCardSelect={onCardSelect}
              disabled={false}
            />
            
            <button
              disabled={selectedCards.length !== 1}
              onClick={() => {
                onReturnTribute(selectedCards[0].deckIndex);
              }}
            >
              确认还贡
            </button>
          </div>
        ) : (
          <div className="waiting-return">
            <p>等待其他玩家还贡...</p>
          </div>
        )}
      </div>
    </div>
  );
};
```

**验证点**：
- 进度正确显示所有还贡任务
- 需要还贡的玩家看到手牌选择界面
- 不需要还贡的玩家看到等待提示
- 复用 `PlayerHand` 组件正常工作
- 选择后正确调用 `onReturnTribute`

---

## 阶段 7：实现进度指示器 (ProgressIndicator.tsx)

### 文件：`frontend/src/components/game/tribute/ProgressIndicator.tsx`

**职责**：显示当前所处阶段

**Props 接口**：
```typescript
interface ProgressIndicatorProps {
  currentPhase: UIPhase;
  isImmune?: boolean;
}
```

**核心实现**：
```typescript
const ProgressIndicator = ({ currentPhase, isImmune }) => {
  const phases = [
    { key: 'START', label: '开始' },
    { key: 'IMMUNITY_CHECK', label: '检测' },
    { key: 'SUBMITTING', label: '贡牌' },
    { key: 'SELECTING', label: '选贡' },
    { key: 'RETURNING', label: '还贡' },
    { key: 'FINISHED', label: '完成' }
  ];
  
  // 如果抗贡，标记跳过的阶段
  const currentIndex = phases.findIndex(p => p.key === currentPhase);
  
  return (
    <div className="progress-indicator">
      {phases.map((phase, index) => {
        const isDone = index < currentIndex;
        const isCurrent = index === currentIndex;
        const isSkipped = isImmune && ['SUBMITTING', 'SELECTING', 'RETURNING'].includes(phase.key);
        
        return (
          <div 
            key={phase.key}
            className={`phase-step ${
              isDone ? 'done' : 
              isCurrent ? 'current' : 
              isSkipped ? 'skipped' : 'pending'
            }`}
          >
            {phase.label}
          </div>
        );
      })}
    </div>
  );
};
```

**验证点**：
- 正确高亮当前阶段
- 抗贡场景正确标记跳过的阶段
- 样式清晰美观

---

## 阶段 8：集成到 GamePage

### 修改文件：`frontend/src/components/game/GamePage.tsx`

**修改点**：

### 导入新组件
```typescript
import TributeFlow from './tribute/TributeFlow';
```

### 删除旧导入
```typescript
// 删除：import TributePhase from './TributePhase';
```

### 修改渲染逻辑（约 line 635-645）
```typescript
case GamePageState.TRIBUTE_PHASE:
  return tributeInfo && room ? (
    <TributeFlow
      tributePhase={tributeInfo}
      players={room.players}
      currentPlayerSeat={playerSeat || 0}
      playerHand={playerHand}
      selectedCards={selectedCards}
      onCardSelect={setSelectedCards}
      onSelectTribute={handleSelectTribute}
      onReturnTribute={handleReturnTribute}
    />
  ) : null;
```

**验证点**：
- 上贡阶段正确渲染新组件
- 所有回调函数正确连接
- 旧组件完全移除

---

## 阶段 9：样式实现

### 创建文件：`frontend/src/components/game/tribute/TributeFlow.css`

**样式范围**：
1. **TributeFlow 容器**：整体布局
2. **ProgressIndicator**：进度条样式
3. **TributePhaseContent**：阶段内容通用样式
4. **SelectionPanel**：选贡面板样式
5. **ReturnPanel**：还贡面板样式（上下布局）
6. **响应式适配**：移动端友好

**关键样式设计**：
- 使用 Tailwind CSS 或自定义 CSS
- 保持与现有 GameBoard、PlayerHand 风格一致
- 添加过渡动画（淡入淡出、卡牌飞行）

**验证点**：
- 布局美观，层次清晰
- 动画流畅自然
- 移动端正常显示

---

## 阶段 10：清理旧代码

### 删除文件
- `frontend/src/components/game/TributePhase.tsx`
- `frontend/src/components/game/TributePhase.test.tsx`

### 检查引用
- 全局搜索 `TributePhase`，确保无遗留引用
- 更新 `frontend/src/components/game/index.ts` 导出

**验证点**：
- 旧文件完全移除
- 无编译错误
- 无运行时警告

---

## 验收标准 (DoD)

### 功能验收
- [ ] 所有场景（抗贡、单下、末游、双下）流程正常
- [ ] 6个阶段按顺序展示
- [ ] 用户交互点（选贡、还贡）正常工作
- [ ] 复用组件（PlayerHand、CardDisplay）集成正常

### 视觉验收
- [ ] 进度指示器清晰显示当前阶段
- [ ] 动画流畅（贡牌飞行、自动切换）
- [ ] 等待提示明确（倒计时、玩家名称）
- [ ] 移动端布局正常

### 性能验收
- [ ] 阶段切换无卡顿
- [ ] 定时器正确清理（无内存泄漏）
- [ ] 后端状态跳跃时前端立即响应

### 代码质量
- [ ] TypeScript 无编译错误
- [ ] 无 ESLint 警告
- [ ] 组件职责单一清晰
- [ ] 代码注释完整

---

## 实施顺序建议

1. **先建骨架**：types.ts → TributeFlow.tsx（空壳）
2. **再填内容**：TributePhaseContent.tsx → 各阶段展示组件
3. **后加交互**：SelectionPanel.tsx → ReturnPanel.tsx
4. **最后优化**：ProgressIndicator.tsx → 样式 → 动画

---

## 风险与应对

| 风险 | 应对措施 |
|------|---------|
| **后端状态推进太快** | 用状态跳跃监听 + 立即切换机制 |
| **定时器冲突** | 每次切换阶段时清理上一个定时器 |
| **PlayerHand 复用问题** | 确保 props 接口兼容，必要时包装 |
| **样式冲突** | 使用 CSS Modules 或独立命名空间 |

---

## 估算

- **开发工作量**：约 1000 行代码
- **文件数量**：7 个新文件
- **复杂度**：中等（主要是状态管理和条件渲染）
- **依赖变更**：无新依赖，复用现有组件

---

## 后续优化方向

1. **动画增强**：使用 Framer Motion 或 React Spring
2. **音效反馈**：添加贡牌、还贡音效
3. **可访问性**：添加 ARIA 标签和键盘导航
4. **测试覆盖**：单元测试 + 集成测试
5. **国际化**：多语言支持


已生成完整的实施方案。

## 方案概要

**目标**：重构上贡组件，实现统一的6阶段流程（开始→检测→贡牌→选贡→还贡→完成）

**核心组件**（仅7个文件）：
1. `TributeFlow` - 主容器（状态管理+阶段路由）
2. `TributePhaseContent` - 内容路由（配置驱动）
3. `SelectionPanel` - 选贡面板（合并3种场景）
4. `ReturnPanel` - 还贡面板（统一布局）
5. `ProgressIndicator` - 进度条
6. `phaseConfigs` - 阶段配置
7. `types` - 类型定义

**关键设计**：
- 前端定时器控制展示节奏
- 同时监听后端状态，支持跳跃
- 复用现有 `PlayerHand`、`CardDisplay`
- 配置驱动减少重复代码

**实施步骤**：10个阶段，从类型定义到清理旧代码
