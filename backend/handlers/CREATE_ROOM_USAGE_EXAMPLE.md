# POST /api/rooms API 使用示例

## API概述

创建一个新的游戏房间。用户成为房间的创建者和房主。

- **端点**: `POST /api/rooms`
- **方法**: POST
- **认证**: 必需 (JWT Bearer Token)
- **请求体**: 无
- **响应状态码**: 201 Created (成功) / 401/409 (失败)

---

## 使用场景

### 场景1: 用户想要开始一局新游戏

用户登录后，创建一个房间等待其他玩家加入。

---

## 请求示例

### cURL
```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### JavaScript (Fetch API)
```javascript
async function createRoom(authToken) {
  try {
    const response = await fetch('http://localhost:8080/api/rooms', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json'
      }
    });

    if (response.ok) {
      const data = await response.json();
      console.log('房间创建成功:', data.room);
      return data.room;
    } else {
      const error = await response.json();
      console.error('创建失败:', error);
      throw new Error(error.message);
    }
  } catch (error) {
    console.error('网络错误:', error);
    throw error;
  }
}

// 使用示例
const token = localStorage.getItem('authToken');
createRoom(token)
  .then(room => {
    console.log(`房间 ${room.id} 创建成功`);
    console.log(`房间状态: ${room.status}`);
    console.log(`玩家数量: ${room.player_count}`);
  })
  .catch(error => {
    console.error('无法创建房间:', error);
  });
```

### TypeScript (React Hook)
```typescript
import { useState } from 'react';

interface Room {
  id: string;
  status: 'waiting' | 'ready' | 'playing';
  owner: string;
  player_count: number;
  players: Player[];
  created_at: string;
  updated_at: string;
}

interface Player {
  id: string;
  username: string;
  seat: number;
  online: boolean;
}

function useCreateRoom() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createRoom = async (authToken: string): Promise<Room | null> => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:8080/api/rooms', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || '创建房间失败');
      }

      const data = await response.json();
      return data.room;
    } catch (err) {
      setError(err instanceof Error ? err.message : '未知错误');
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { createRoom, loading, error };
}

// 在组件中使用
function CreateRoomButton() {
  const { createRoom, loading, error } = useCreateRoom();
  const authToken = useAuthToken(); // 假设有个hook获取token

  const handleCreateRoom = async () => {
    const room = await createRoom(authToken);
    if (room) {
      console.log('房间创建成功:', room);
      // 跳转到房间页面
      window.location.href = `/room/${room.id}`;
    }
  };

  return (
    <div>
      <button 
        onClick={handleCreateRoom} 
        disabled={loading}
      >
        {loading ? '创建中...' : '创建房间'}
      </button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
}
```

### Python (requests)
```python
import requests

def create_room(auth_token):
    """创建游戏房间"""
    url = 'http://localhost:8080/api/rooms'
    headers = {
        'Authorization': f'Bearer {auth_token}',
        'Content-Type': 'application/json'
    }
    
    try:
        response = requests.post(url, headers=headers)
        response.raise_for_status()
        
        data = response.json()
        room = data['room']
        
        print(f"房间创建成功!")
        print(f"房间ID: {room['id']}")
        print(f"房间状态: {room['status']}")
        print(f"房主ID: {room['owner']}")
        print(f"玩家数量: {room['player_count']}")
        
        return room
        
    except requests.exceptions.HTTPError as e:
        if response.status_code == 401:
            print("错误: 未认证或Token无效")
        elif response.status_code == 409:
            error_data = response.json()
            print(f"错误: {error_data.get('message', '用户已在其他房间中')}")
        else:
            print(f"HTTP错误: {e}")
        return None
        
    except requests.exceptions.RequestException as e:
        print(f"网络错误: {e}")
        return None

# 使用示例
if __name__ == '__main__':
    token = 'your_jwt_token_here'
    room = create_room(token)
    
    if room:
        print(f"\n✓ 成功创建房间 {room['id']}")
    else:
        print("\n✗ 创建房间失败")
```

### Go (标准库)
```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Room struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	Owner       string    `json:"owner"`
	PlayerCount int       `json:"player_count"`
	Players     []Player  `json:"players"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Seat     int    `json:"seat"`
	Online   bool   `json:"online"`
}

type RoomResponse struct {
	Room *Room `json:"room"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func CreateRoom(authToken string) (*Room, error) {
	url := "http://localhost:8080/api/rooms"
	
	// 创建请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 设置headers
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")
	
	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 检查状态码
	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("创建房间失败 (%d): %s", resp.StatusCode, errResp.Message)
		}
		return nil, fmt.Errorf("创建房间失败: HTTP %d", resp.StatusCode)
	}
	
	// 解析成功响应
	var roomResp RoomResponse
	if err := json.Unmarshal(body, &roomResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	
	return roomResp.Room, nil
}

func main() {
	token := "your_jwt_token_here"
	
	room, err := CreateRoom(token)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	fmt.Printf("✓ 房间创建成功!\n")
	fmt.Printf("  房间ID: %s\n", room.ID)
	fmt.Printf("  状态: %s\n", room.Status)
	fmt.Printf("  房主: %s\n", room.Owner)
	fmt.Printf("  玩家数: %d/4\n", room.PlayerCount)
	fmt.Printf("  创建时间: %s\n", room.CreatedAt.Format(time.RFC3339))
}
```

---

## 成功响应示例

### HTTP 201 Created

```json
{
  "room": {
    "id": "room_1699234567890123456",
    "status": "waiting",
    "owner": "user_1699234567890123456",
    "player_count": 1,
    "players": [
      {
        "id": "user_1699234567890123456",
        "username": "Alice",
        "seat": 0,
        "online": true
      },
      null,
      null,
      null
    ],
    "created_at": "2024-11-06T10:30:00Z",
    "updated_at": "2024-11-06T10:30:00Z"
  }
}
```

**响应字段说明**:

- `id`: 房间唯一标识符
- `status`: 房间状态
  - `waiting`: 等待玩家加入（1-3人）
  - `ready`: 4人已满，可以开始游戏
  - `playing`: 游戏进行中
- `owner`: 房主用户ID（创建者）
- `player_count`: 当前玩家数量
- `players`: 玩家数组（固定4个位置）
  - `seat 0-3`: 四个座位，空座位为 null
  - `id`: 玩家用户ID
  - `username`: 玩家用户名
  - `seat`: 座位号（0-3）
  - `online`: 是否在线
- `created_at`: 房间创建时间
- `updated_at`: 最后更新时间

---

## 错误响应示例

### 1. 未认证 (401 Unauthorized)

```json
{
  "error": "unauthorized",
  "message": "User not authenticated"
}
```

**原因**: 
- 未提供 Authorization header
- Token格式错误
- Token无效或已过期

**解决方法**: 
- 确保用户已登录
- 使用 `POST /api/auth/login` 获取有效Token
- 正确设置 Authorization header

---

### 2. 用户已在房间中 (409 Conflict)

```json
{
  "error": "already_in_room",
  "message": "player is already in room room_1699234567890123456"
}
```

**原因**: 
- 用户已经创建了一个房间
- 用户已经加入了其他用户的房间

**解决方法**: 
- 先离开当前房间: `POST /api/rooms/{room_id}/leave`
- 或直接前往当前房间: `GET /api/rooms/my`

---

## 完整工作流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant Backend as 后端API
    participant RoomService as 房间服务

    User->>Frontend: 点击"创建房间"
    Frontend->>Frontend: 检查是否有Token
    
    alt 没有Token
        Frontend->>User: 跳转到登录页
    else 有Token
        Frontend->>Backend: POST /api/rooms<br/>Authorization: Bearer {token}
        Backend->>Backend: 验证Token
        
        alt Token无效
            Backend-->>Frontend: 401 Unauthorized
            Frontend->>User: 显示"请重新登录"
        else Token有效
            Backend->>RoomService: CreateRoom(userID)
            RoomService->>RoomService: 检查用户是否在其他房间
            
            alt 用户已在房间
                RoomService-->>Backend: Error: already_in_room
                Backend-->>Frontend: 409 Conflict
                Frontend->>User: 显示"您已在其他房间中"
            else 用户未在房间
                RoomService->>RoomService: 创建房间
                RoomService->>RoomService: 添加用户到座位0
                RoomService-->>Backend: Room对象
                Backend-->>Frontend: 201 Created + Room数据
                Frontend->>User: 跳转到房间页面
                Frontend->>Backend: 建立WebSocket连接
            end
        end
    end
```

---

## 最佳实践

### 1. 错误处理

```typescript
async function createRoomWithErrorHandling(token: string) {
  try {
    const response = await fetch('http://localhost:8080/api/rooms', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    // 处理不同的错误状态
    if (response.status === 401) {
      // Token无效，重定向到登录
      window.location.href = '/login';
      throw new Error('请重新登录');
    }

    if (response.status === 409) {
      // 用户已在房间中
      const error = await response.json();
      // 提取房间ID
      const roomId = error.message.match(/room_\d+/)?.[0];
      if (roomId) {
        // 直接跳转到现有房间
        window.location.href = `/room/${roomId}`;
      }
      throw new Error('您已在房间中');
    }

    if (!response.ok) {
      throw new Error('创建房间失败');
    }

    const data = await response.json();
    return data.room;

  } catch (error) {
    console.error('创建房间错误:', error);
    throw error;
  }
}
```

### 2. 加载状态管理

```typescript
function CreateRoomComponent() {
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async () => {
    setIsCreating(true);
    setError(null);

    try {
      const room = await createRoom(authToken);
      // 成功：跳转到房间
      navigate(`/room/${room.id}`);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <button 
      onClick={handleCreate} 
      disabled={isCreating}
    >
      {isCreating ? '创建中...' : '创建房间'}
    </button>
  );
}
```

### 3. 重试机制

```typescript
async function createRoomWithRetry(
  token: string, 
  maxRetries: number = 3
): Promise<Room> {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await createRoom(token);
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      
      // 如果是网络错误，等待后重试
      if (error.message.includes('网络')) {
        await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
        continue;
      }
      
      // 其他错误不重试
      throw error;
    }
  }
}
```

---

## 测试建议

### 单元测试示例 (Jest)

```typescript
import { createRoom } from './api';

describe('createRoom', () => {
  it('应该成功创建房间', async () => {
    const mockResponse = {
      room: {
        id: 'room_123',
        status: 'waiting',
        owner: 'user_1',
        player_count: 1,
        players: [{ id: 'user_1', username: 'Test', seat: 0, online: true }],
        created_at: '2024-11-06T10:00:00Z',
        updated_at: '2024-11-06T10:00:00Z'
      }
    };

    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: async () => mockResponse
    });

    const room = await createRoom('valid_token');
    expect(room.id).toBe('room_123');
    expect(room.status).toBe('waiting');
  });

  it('应该处理401错误', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 401,
      json: async () => ({ error: 'unauthorized' })
    });

    await expect(createRoom('invalid_token')).rejects.toThrow();
  });
});
```

---

## 相关API

创建房间后，您可能需要使用以下相关API：

- `GET /api/rooms/my` - 获取当前用户所在的房间
- `GET /api/rooms` - 获取所有房间列表
- `POST /api/rooms/:id/start` - 开始游戏（房主）
- `ws://localhost:8080/ws` - WebSocket连接（实时更新）

---

## 常见问题

**Q: 可以同时创建多个房间吗？**  
A: 不可以。一个用户同时只能在一个房间中（无论是创建的还是加入的）。

**Q: 创建房间后如何邀请其他玩家？**  
A: 将房间ID分享给其他玩家，他们可以使用 `POST /api/rooms/{room_id}/join` 加入。

**Q: 房间会自动关闭吗？**  
A: 当最后一个玩家离开时，房间会自动关闭。

**Q: 创建房间失败后如何重试？**  
A: 检查错误原因。如果是 409 错误（已在房间中），先离开当前房间再重试。

**Q: 如何获取认证Token？**  
A: 使用 `POST /api/auth/login` 或 `POST /api/auth/register` 接口。

---

## 联系与支持

如有问题或建议，请参考:
- API文档: `backend/API-Documentation.md`
- 测试用例: `backend/handlers/room_test.go`
- 技术文档: `backend/Technical-Documentation.md`




