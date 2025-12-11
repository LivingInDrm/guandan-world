## 目标

完善 `ui-next/` 基础组件库，统一使用新设计 token（`--ds-*`），为业务组件迁移做准备。

---

## 核心变更

### Phase 1: 重命名现有组件

| 操作 | 旧 | 新 | 变更内容 |
|------|----|----|----------|
| **重命名** | ActionButton.tsx | Button.tsx | 增加 `danger` variant |
| **重命名** | ActionButton.stories.tsx | Button.stories.tsx | 更新导入和组件名 |
| **重命名** | Surface.tsx | Card.tsx | 仅重命名，逻辑不变 |
| **重命名** | Surface.stories.tsx | Card.stories.tsx | 更新导入和组件名 |
| **更新** | index.ts | - | 导出 Button, Card（移除 ActionButton, Surface） |

### Phase 2: 新建组件

| 组件 | 文件 | 复杂度 |
|------|------|--------|
| Dialog | Dialog.tsx + Dialog.stories.tsx | 中 |
| Input | Input.tsx + Input.stories.tsx | 低 |
| DropdownMenu | DropdownMenu.tsx + DropdownMenu.stories.tsx | 中 |
| Slider | Slider.tsx + Slider.stories.tsx | 低 |

---

## 详细实现

### 1. Button（扩展 ActionButton）

**位置**: `ui-next/Button.tsx`

**修改内容**:
```tsx
// 在 buttonVariants 的 intent 中增加 danger
intent: {
  primary: "bg-ds-action-primary text-ds-text-inverse",
  secondary: "bg-ds-action-secondary text-ds-text-inverse",
  neutral: "bg-ds-action-neutral text-ds-text-inverse",
  danger: "bg-ds-primitive-danger-500 text-ds-text-inverse", // 新增
}
```

**导出名**:
```tsx
// 修改前
export { ActionButton, actionButtonVariants }

// 修改后
export { Button, buttonVariants }
```

**Stories 更新**:
- 重命名文件：`ActionButton.stories.tsx` → `Button.stories.tsx`
- 更新导入：`import { ActionButton }` → `import { Button }`
- 新增 story：`export const DangerVariant`（展示 danger variant）

---

### 2. Card（重命名 Surface）

**位置**: `ui-next/Card.tsx`

**修改内容**:
```tsx
// 仅重命名，逻辑不变
const cardVariants = cva(/* 原 surfaceVariants 内容 */)

export interface CardProps extends React.HTMLAttributes<HTMLDivElement>,
  VariantProps<typeof cardVariants> {}

const Card = React.forwardRef<HTMLDivElement, CardProps>(/* ... */)

export { Card, cardVariants }
```

**Stories 更新**:
- 重命名文件：`Surface.stories.tsx` → `Card.stories.tsx`
- 更新导入：`import { Surface }` → `import { Card }`
- 更新所有 `<Surface>` → `<Card>`

---

### 3. Dialog

**位置**: `ui-next/Dialog.tsx`

**基于旧版改造，Token 替换**:
```tsx
// Overlay
"bg-black/60 backdrop-blur-sm" // 保持

// Content
"bg-ds-surface-elevated rounded-ds-lg shadow-ds-elevation-3"

// Title
"text-ds-text-primary"

// Description
"text-ds-text-secondary"
```

**Stories**:
- Default（基础弹窗）
- WithFooter（带 Footer 的表单弹窗）
- LongContent（长内容滚动）

---

### 4. Input

**位置**: `ui-next/Input.tsx`

**Token 替换**:
```tsx
"border-ds-border bg-ds-surface-base text-ds-text-primary"
"rounded-ds-sm shadow-ds-elevation-1"
"focus-visible:border-ds-state-active focus-visible:shadow-ds-elevation-2"
"disabled:bg-ds-state-disabled disabled:text-ds-text-secondary"
```

**Label**:
```tsx
"text-ds-text-primary"
```

**Stories**:
- Default
- WithLabel
- Disabled
- Error（可选，如需 error 状态）

---

### 5. DropdownMenu

**位置**: `ui-next/DropdownMenu.tsx`

**Token 替换**:
```tsx
// Content
"bg-ds-surface-elevated shadow-ds-elevation-2 rounded-ds-md"

// Item
"text-ds-text-primary hover:bg-ds-surface-emphasis"

// Separator
"bg-ds-border"
```

**Stories**:
- Default（基础菜单）
- WithSeparator（带分隔符）

---

### 6. Slider

**位置**: `ui-next/Slider.tsx`

**Token 替换**:
```tsx
// Track
"bg-ds-state-disabled rounded-ds-sm"

// Range
"bg-ds-state-active"

// Thumb
"border-ds-state-active bg-ds-surface-base"
```

**Stories**:
- Default
- WithValue（展示当前值）

---

## 实施步骤

### Step 1: 重命名现有组件
1. `ActionButton.tsx` → `Button.tsx`（增加 danger）
2. `ActionButton.stories.tsx` → `Button.stories.tsx`
3. `Surface.tsx` → `Card.tsx`
4. `Surface.stories.tsx` → `Card.stories.tsx`
5. 更新 `index.ts` 导出

### Step 2: 新建组件（按优先级）
1. Dialog（中优先级，通用弹窗）
2. Input（中优先级，表单）
3. DropdownMenu（低优先级，菜单）
4. Slider（低优先级，滑块）

### Step 3: 验证
- 运行 `npm run ladle`
- 检查所有 stories
- 验证 TypeScript 无错误

---

## 验证标准

### 重命名部分
- [ ] Button 有 4 个 variant（primary/secondary/neutral/danger）
- [ ] Card 保持 3 个 variant（base/elevated/emphasis）
- [ ] Ladle 中 Button/Card stories 正常显示
- [ ] `index.ts` 正确导出 Button, Card

### 新组件部分
- [ ] 每个组件使用 `--ds-*` token
- [ ] 每个组件有对应的 stories
- [ ] TypeScript 无错误
- [ ] API 与旧版兼容

---

## 文件清单

| 操作 | 文件 |
|------|------|
| **重命名** | ActionButton.tsx → Button.tsx |
| **重命名** | ActionButton.stories.tsx → Button.stories.tsx |
| **重命名** | Surface.tsx → Card.tsx |
| **重命名** | Surface.stories.tsx → Card.stories.tsx |
| **更新** | index.ts |
| **新建** | Dialog.tsx, Dialog.stories.tsx |
| **新建** | Input.tsx, Input.stories.tsx |
| **新建** | DropdownMenu.tsx, DropdownMenu.stories.tsx |
| **新建** | Slider.tsx, Slider.stories.tsx |

**总计**: 4 个重命名 + 1 个更新 + 8 个新建

---

## 注意事项

1. **重命名后无需保留旧组件**（已确认选项 A）
2. **Avatar 和 Badge 保持不变**（已实现）
3. **所有组件统一游戏风格**（圆角、阴影、过渡）
4. **保持 API 兼容**，便于业务组件迁移