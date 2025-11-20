# Protobuf 数据结构定义

本目录包含掼蛋游戏的 Protocol Buffers 定义，作为前后端数据结构的**单一来源**。

## 文件说明

| 文件 | 说明 | 主要类型 |
|-----|------|---------|
| `common.proto` | 公共数据结构 | `Card`, `PlayerBasicInfo`, `PlayAction`, `CompType` |
| `event.proto` | 游戏事件定义 | `GameEvent`, `EventType`, 44种事件 Payload |
| `view.proto` | 玩家视图定义 | `PlayerView`, `TributeView`, `DealStatus` |

## 使用指南

### 📱 前端开发者（TypeScript/JavaScript）

查看 **[FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)**

- WebSocket 消息格式
- TypeScript 类型映射
- 事件类型速查表
- 常见场景示例代码

### 🔧 后端开发者（Go）

查看 **[BACKEND_GUIDE.md](./BACKEND_GUIDE.md)**

- 代码生成和导入
- 创建和序列化 Protobuf 消息
- 字段填写规则
- 最佳实践

## 代码生成

**自动生成：** Proto 文件会在以下时机自动生成：
- 运行 `npm install`（根目录）
- 运行 `npm run dev` 或 `npm run build`（frontend 目录）

**手动生成：**

```bash
# 生成所有代码（Go + TypeScript）
make proto

# 仅生成 Go 代码
make proto-go

# 仅生成 TypeScript 代码
make proto-ts

# 仅生成 JavaScript 代码
make proto-js
```

**首次克隆仓库：**

```bash
# 在根目录运行（会自动生成 TypeScript 类型）
npm install
```

## 核心概念

### 三个主要数据流

1. **game_event** (`event.proto`)
   - 游戏事件通知（所有玩家收到相同内容）
   - 用于动画、提示、日志

2. **player_view** (`view.proto`)
   - 玩家视角状态（每个玩家收到不同内容）
   - 包含私有手牌信息
   - **唯一可信的状态来源**

3. **Common 类型** (`common.proto`)
   - 共享数据结构
   - 被 event 和 view 引用

### 数据传输方式

```
Proto 定义 → JSON (WebSocket) → 前端
          ↓
       Go 代码 (Backend)
```

## 向前兼容性

- ✅ 添加新字段使用新的编号
- ✅ 废弃字段用 `[deprecated = true]` 标记
- ❌ 不要删除字段或修改字段编号

## 目录结构

```
proto/
├── common.proto          # 公共数据结构定义
├── event.proto           # 游戏事件定义
├── view.proto            # 玩家视图定义
├── README.md             # 本文档（总览）
├── FRONTEND_GUIDE.md     # 前端使用指南
├── BACKEND_GUIDE.md      # 后端使用指南
└── [生成的代码]
    ├── common/           # Go: common.pb.go
    ├── event/            # Go: event.pb.go
    └── view/             # Go: view.pb.go
```

## 参考资料

- **架构文档：** `docs/EventSystemArchitecture.md`
- **API 文档：** `backend/API-Documentation.md`
- **WebSocket 协议：** `backend/API-Documentation.md#websocket-protocol`
