# GET /api/rooms API 使用示例

## 概述

本文档提供 `GET /api/rooms` API 的实际使用示例，包括各种场景下的请求格式和响应示例。

## 目录

1. [基础用法](#基础用法)
2. [分页查询](#分页查询)
3. [状态过滤](#状态过滤)
4. [组合查询](#组合查询)
5. [错误处理](#错误处理)
6. [前端集成示例](#前端集成示例)

---

## 基础用法

### 示例 1: 获取房间列表（默认参数）

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543210",
      "status": "waiting",
      "player_count": 3,
      "players": [
        {
          "id": "user_1699876500123",
          "username": "player1",
          "seat": 0,
          "online": true
        },
        {
          "id": "user_1699876500456",
          "username": "player2",
          "seat": 1,
          "online": true
        },
        {
          "id": "user_1699876500789",
          "username": "player3",
          "seat": 2,
          "online": true
        }
      ],
      "owner": "user_1699876500123",
      "can_join": true,
      "created_at": "2025-11-06T10:30:00Z"
    },
    {
      "id": "room_1699876543211",
      "status": "ready",
      "player_count": 4,
      "players": [
        {
          "id": "user_1699876501000",
          "username": "alice",
          "seat": 0,
          "online": true
        },
        {
          "id": "user_1699876501111",
          "username": "bob",
          "seat": 1,
          "online": true
        },
        {
          "id": "user_1699876501222",
          "username": "charlie",
          "seat": 2,
          "online": true
        },
        {
          "id": "user_1699876501333",
          "username": "david",
          "seat": 3,
          "online": true
        }
      ],
      "owner": "user_1699876501000",
      "can_join": false,
      "created_at": "2025-11-06T10:32:00Z"
    }
  ],
  "total_count": 2,
  "page": 1,
  "limit": 12
}
```

**说明**:
- 默认返回第1页，每页12个房间
- 房间按照排序规则排列（waiting优先，玩家数多的优先）
- `can_join` 字段表示是否可以加入该房间

---

## 分页查询

### 示例 2: 获取第2页，每页5个房间

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?page=2&limit=5" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543216",
      "status": "waiting",
      "player_count": 1,
      "players": [
        {
          "id": "user_1699876502000",
          "username": "user6",
          "seat": 0,
          "online": true
        }
      ],
      "owner": "user_1699876502000",
      "can_join": true,
      "created_at": "2025-11-06T10:35:00Z"
    }
    // ... 更多房间 ...
  ],
  "total_count": 15,
  "page": 2,
  "limit": 5
}
```

**说明**:
- `page=2` 表示第2页
- `limit=5` 表示每页5个房间
- `total_count=15` 表示总共15个房间
- 第2页将返回第6-10个房间

### 示例 3: 请求超出范围的页码

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?page=100&limit=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [],
  "total_count": 15,
  "page": 100,
  "limit": 10
}
```

**说明**:
- 页码超出范围时返回空数组
- `total_count` 仍然显示总数
- 响应状态码仍为200（不是错误）

---

## 状态过滤

### 示例 4: 只获取等待中的房间

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?status=waiting" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543210",
      "status": "waiting",
      "player_count": 3,
      "players": [
        // ... 玩家信息 ...
      ],
      "owner": "user_1699876500123",
      "can_join": true,
      "created_at": "2025-11-06T10:30:00Z"
    },
    {
      "id": "room_1699876543215",
      "status": "waiting",
      "player_count": 2,
      "players": [
        // ... 玩家信息 ...
      ],
      "owner": "user_1699876501500",
      "can_join": true,
      "created_at": "2025-11-06T10:33:00Z"
    }
  ],
  "total_count": 2,
  "page": 1,
  "limit": 12
}
```

**说明**:
- `status=waiting` 只返回等待中的房间
- 所有返回的房间都是waiting状态
- `can_join` 都是true（waiting状态且未满）

### 示例 5: 只获取准备开始的房间

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?status=ready" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543211",
      "status": "ready",
      "player_count": 4,
      "players": [
        // ... 4个玩家信息 ...
      ],
      "owner": "user_1699876501000",
      "can_join": false,
      "created_at": "2025-11-06T10:32:00Z"
    }
  ],
  "total_count": 1,
  "page": 1,
  "limit": 12
}
```

**说明**:
- `status=ready` 返回4人已满的房间
- `player_count` 都是4
- `can_join` 都是false（已满）

### 示例 6: 只获取游戏中的房间

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?status=playing" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543220",
      "status": "playing",
      "player_count": 4,
      "players": [
        // ... 4个玩家信息 ...
      ],
      "owner": "user_1699876503000",
      "can_join": false,
      "created_at": "2025-11-06T10:40:00Z"
    }
  ],
  "total_count": 1,
  "page": 1,
  "limit": 12
}
```

**说明**:
- `status=playing` 返回游戏进行中的房间
- `can_join` 都是false（游戏中无法加入）

---

## 组合查询

### 示例 7: 分页 + 状态过滤

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms?status=waiting&page=1&limit=3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543210",
      "status": "waiting",
      "player_count": 3,
      "players": [/* ... */],
      "owner": "user_1699876500123",
      "can_join": true,
      "created_at": "2025-11-06T10:30:00Z"
    },
    {
      "id": "room_1699876543215",
      "status": "waiting",
      "player_count": 2,
      "players": [/* ... */],
      "owner": "user_1699876501500",
      "can_join": true,
      "created_at": "2025-11-06T10:33:00Z"
    },
    {
      "id": "room_1699876543216",
      "status": "waiting",
      "player_count": 1,
      "players": [/* ... */],
      "owner": "user_1699876502000",
      "can_join": true,
      "created_at": "2025-11-06T10:35:00Z"
    }
  ],
  "total_count": 8,
  "page": 1,
  "limit": 3
}
```

**说明**:
- 先过滤waiting状态（8个）
- 然后分页，返回前3个
- `total_count=8` 是过滤后的总数

---

## 错误处理

### 示例 8: 未提供认证Token

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms"
```

**响应** (401 Unauthorized):
```json
{
  "error": "missing_token",
  "message": "Authorization header is required"
}
```

### 示例 9: 使用无效的Token

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms" \
  -H "Authorization: Bearer invalid-token-12345"
```

**响应** (401 Unauthorized):
```json
{
  "error": "invalid_token",
  "message": "failed to parse token: token is malformed: token contains an invalid number of segments"
}
```

### 示例 10: 使用过期的Token

**请求**:
```bash
curl -X GET "http://localhost:8080/api/rooms" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.expired..."
```

**响应** (401 Unauthorized):
```json
{
  "error": "invalid_token",
  "message": "token expired"
}
```

---

## 前端集成示例

### React/TypeScript 示例

```typescript
// types.ts
export interface Player {
  id: string;
  username: string;
  seat: number;
  online: boolean;
}

export interface RoomInfo {
  id: string;
  status: 'waiting' | 'ready' | 'playing';
  player_count: number;
  players: Player[];
  owner: string;
  can_join: boolean;
  created_at: string;
}

export interface RoomListResponse {
  rooms: RoomInfo[];
  total_count: number;
  page: number;
  limit: number;
}

// api.ts
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

// 获取房间列表
export async function getRoomList(
  page: number = 1,
  limit: number = 12,
  status?: 'waiting' | 'ready' | 'playing'
): Promise<RoomListResponse> {
  const params: any = { page, limit };
  if (status) {
    params.status = status;
  }

  const response = await axios.get(`${API_BASE_URL}/rooms`, {
    params,
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });

  return response.data;
}

// RoomList.tsx - 房间列表组件
import React, { useEffect, useState } from 'react';
import { getRoomList } from './api';
import type { RoomInfo, RoomListResponse } from './types';

interface RoomListProps {
  statusFilter?: 'waiting' | 'ready' | 'playing';
}

export const RoomList: React.FC<RoomListProps> = ({ statusFilter }) => {
  const [roomData, setRoomData] = useState<RoomListResponse | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const limit = 12;

  useEffect(() => {
    loadRooms();
  }, [currentPage, statusFilter]);

  const loadRooms = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getRoomList(currentPage, limit, statusFilter);
      setRoomData(data);
    } catch (err: any) {
      setError(err.response?.data?.message || '加载房间列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handlePrevPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    if (roomData && currentPage * limit < roomData.total_count) {
      setCurrentPage(currentPage + 1);
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'waiting': return '等待中';
      case 'ready': return '准备就绪';
      case 'playing': return '游戏中';
      default: return status;
    }
  };

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  if (error) {
    return (
      <div className="error">
        <p>错误: {error}</p>
        <button onClick={loadRooms}>重试</button>
      </div>
    );
  }

  if (!roomData || roomData.rooms.length === 0) {
    return <div className="empty">暂无房间</div>;
  }

  return (
    <div className="room-list">
      <div className="room-list-header">
        <h2>房间列表</h2>
        <p>共 {roomData.total_count} 个房间</p>
      </div>

      <div className="rooms">
        {roomData.rooms.map((room) => (
          <div key={room.id} className="room-card">
            <div className="room-header">
              <h3>房间 {room.id}</h3>
              <span className={`status status-${room.status}`}>
                {getStatusText(room.status)}
              </span>
            </div>

            <div className="room-info">
              <p>玩家数: {room.player_count}/4</p>
              <p>房主: {room.players.find(p => p.id === room.owner)?.username}</p>
            </div>

            <div className="players">
              {room.players.map((player) => (
                <div key={player.id} className="player">
                  <span>{player.username}</span>
                  {player.id === room.owner && <span className="owner-badge">房主</span>}
                </div>
              ))}
            </div>

            {room.can_join && (
              <button className="join-button">加入房间</button>
            )}
          </div>
        ))}
      </div>

      <div className="pagination">
        <button 
          onClick={handlePrevPage} 
          disabled={currentPage === 1}
        >
          上一页
        </button>
        <span>第 {currentPage} 页</span>
        <button 
          onClick={handleNextPage}
          disabled={currentPage * limit >= roomData.total_count}
        >
          下一页
        </button>
      </div>
    </div>
  );
};

// App.tsx - 使用示例
import React, { useState } from 'react';
import { RoomList } from './RoomList';

export const App: React.FC = () => {
  const [statusFilter, setStatusFilter] = useState<'waiting' | 'ready' | 'playing' | undefined>(undefined);

  return (
    <div className="app">
      <div className="filters">
        <button onClick={() => setStatusFilter(undefined)}>全部</button>
        <button onClick={() => setStatusFilter('waiting')}>等待中</button>
        <button onClick={() => setStatusFilter('ready')}>准备就绪</button>
        <button onClick={() => setStatusFilter('playing')}>游戏中</button>
      </div>

      <RoomList statusFilter={statusFilter} />
    </div>
  );
};
```

### JavaScript/Fetch 示例

```javascript
// 简单的JavaScript示例
async function getRooms(page = 1, limit = 12, status = null) {
  const params = new URLSearchParams({
    page: page.toString(),
    limit: limit.toString(),
  });
  
  if (status) {
    params.append('status', status);
  }

  const token = localStorage.getItem('token');
  
  const response = await fetch(`http://localhost:8080/api/rooms?${params}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }

  return await response.json();
}

// 使用示例
getRooms(1, 12, 'waiting')
  .then(data => {
    console.log('房间列表:', data.rooms);
    console.log('总数:', data.total_count);
  })
  .catch(error => {
    console.error('获取房间列表失败:', error);
  });
```

---

## 常见使用场景

### 场景 1: 大厅页面加载房间列表
```bash
GET /api/rooms?limit=12
```

### 场景 2: 只显示可加入的房间
```bash
GET /api/rooms?status=waiting
```

### 场景 3: 翻页浏览所有房间
```bash
GET /api/rooms?page=1&limit=10
GET /api/rooms?page=2&limit=10
```

### 场景 4: 查看游戏中的房间
```bash
GET /api/rooms?status=playing
```

### 场景 5: 实时刷新房间列表
```javascript
// 每5秒刷新一次
setInterval(async () => {
  const data = await getRooms();
  updateRoomList(data);
}, 5000);
```

---

## 参数说明总结

| 参数 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | number | 否 | 1 | 页码，从1开始 |
| limit | number | 否 | 12 | 每页数量，范围1-50 |
| status | string | 否 | - | 状态过滤：waiting/ready/playing |

## 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| rooms | array | 房间信息数组 |
| rooms[].id | string | 房间ID |
| rooms[].status | string | 房间状态 |
| rooms[].player_count | number | 当前玩家数 |
| rooms[].players | array | 玩家信息数组 |
| rooms[].owner | string | 房主用户ID |
| rooms[].can_join | boolean | 是否可加入 |
| rooms[].created_at | string | 创建时间 (ISO 8601) |
| total_count | number | 总房间数（过滤后） |
| page | number | 当前页码 |
| limit | number | 每页数量 |

---

## 最佳实践

1. **分页加载**: 使用合理的limit值（建议12或24）避免一次加载过多数据
2. **状态过滤**: 在不同场景下使用适当的状态过滤
3. **错误处理**: 始终处理401错误，提示用户重新登录
4. **实时更新**: 使用定时刷新或WebSocket保持数据最新
5. **性能优化**: 避免频繁请求，使用防抖/节流
6. **用户体验**: 显示加载状态和友好的错误提示

---

## 相关API

- **POST /api/rooms** - 创建房间
- **GET /api/rooms/:id** - 获取单个房间详情
- **GET /api/rooms/my** - 获取当前用户的房间
- **POST /api/rooms/:id/join** - 加入房间
- **POST /api/rooms/:id/leave** - 离开房间
- **POST /api/rooms/:id/start** - 开始游戏



