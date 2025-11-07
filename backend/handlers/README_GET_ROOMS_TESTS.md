# GET /api/rooms API 测试文档

## 文档导航

本目录包含 `GET /api/rooms` API 的完整测试文档：

1. **[GET_ROOMS_TEST_CASES.md](./GET_ROOMS_TEST_CASES.md)** - 详细的测试用例文档
2. **[GET_ROOMS_TEST_SUMMARY.md](./GET_ROOMS_TEST_SUMMARY.md)** - 测试执行总结
3. **[GET_ROOMS_USAGE_EXAMPLE.md](./GET_ROOMS_USAGE_EXAMPLE.md)** - API使用示例
4. **[get_rooms_api_test.go](./get_rooms_api_test.go)** - 测试实现代码

---

## 快速开始

### 运行所有测试

```bash
cd /Users/xiaochunliu/program/guandan-world
go test -v ./backend/handlers -run TestGetRoomsAPI
```

### 运行特定测试

```bash
# 测试默认参数
go test -v ./backend/handlers -run TestGetRoomsAPI_DefaultParameters

# 测试分页功能
go test -v ./backend/handlers -run TestGetRoomsAPI_Pagination

# 测试状态过滤
go test -v ./backend/handlers -run TestGetRoomsAPI_FilterByStatus
```

### 生成覆盖率报告

```bash
go test -v ./backend/handlers -run TestGetRoomsAPI -coverprofile=coverage_get_rooms.out
go tool cover -html=coverage_get_rooms.out
```

---

## API 概览

### 端点信息

- **URL**: `/api/rooms`
- **方法**: `GET`
- **认证**: 需要 JWT Token (Bearer)
- **功能**: 获取游戏房间列表

### 请求参数

| 参数 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `page` | number | 否 | 1 | 页码，从1开始 |
| `limit` | number | 否 | 12 | 每页数量，范围1-50 |
| `status` | string | 否 | - | 状态过滤：`waiting`/`ready`/`playing` |

### 响应格式

**成功响应** (200 OK):
```json
{
  "rooms": [
    {
      "id": "room_1699876543210",
      "status": "waiting",
      "player_count": 3,
      "players": [...],
      "owner": "user_1699876500123",
      "can_join": true,
      "created_at": "2025-11-06T10:30:00Z"
    }
  ],
  "total_count": 10,
  "page": 1,
  "limit": 12
}
```

**错误响应** (401 Unauthorized):
```json
{
  "error": "missing_token",
  "message": "Authorization header is required"
}
```

---

## 测试覆盖

### 测试分类 (27个测试用例)

| 分类 | 测试数量 | 状态 |
|------|---------|------|
| 正常流程测试 | 3 | ✅ |
| 认证相关测试 | 3 | ✅ |
| 分页功能测试 | 3 | ✅ |
| 状态过滤测试 | 5 | ✅ |
| 数据完整性测试 | 3 | ✅ |
| 边界条件测试 | 4 | ✅ |
| 边缘情况测试 | 3 | ✅ |
| 集成测试 | 3 | ✅ |
| **总计** | **27** | **✅ 100%** |

### 代码覆盖率

- **Handler层**: ~95%
- **Service层**: ~100%
- **整体覆盖率**: > 95%

---

## 测试用例详情

### 1. 正常流程测试 (3个)

✅ **TestGetRoomsAPI_DefaultParameters**
- 验证使用默认参数获取房间列表
- 检查响应格式和字段完整性

✅ **TestGetRoomsAPI_EmptyList**
- 验证没有房间时返回空列表
- 确保响应格式正确

✅ **TestGetRoomsAPI_RoomSorting**
- 验证房间排序规则
- waiting状态优先，玩家数多的优先

### 2. 认证相关测试 (3个)

✅ **TestGetRoomsAPI_NoAuthToken**
- 验证未提供Token时返回401

✅ **TestGetRoomsAPI_InvalidToken**
- 验证无效Token时返回401

✅ **TestGetRoomsAPI_ExpiredToken**
- 验证过期Token时返回401

### 3. 分页功能测试 (3个)

✅ **TestGetRoomsAPI_Pagination**
- 验证分页参数正确工作
- 测试多页数据的正确性

✅ **TestGetRoomsAPI_PageOutOfRange**
- 验证页码超出范围时返回空数组

✅ **TestGetRoomsAPI_InvalidPaginationParams**
- 验证无效参数使用默认值

### 4. 状态过滤测试 (5个)

✅ **TestGetRoomsAPI_FilterByStatusWaiting**
- 验证waiting状态过滤

✅ **TestGetRoomsAPI_FilterByStatusReady**
- 验证ready状态过滤

✅ **TestGetRoomsAPI_FilterByStatusPlaying**
- 验证playing状态过滤

✅ **TestGetRoomsAPI_InvalidStatusFilter**
- 验证无效状态值被忽略

✅ **TestGetRoomsAPI_CombinedPaginationAndFilter**
- 验证分页和过滤组合使用

### 5. 数据完整性测试 (3个)

✅ **TestGetRoomsAPI_RoomInfoIntegrity**
- 验证房间信息字段完整性

✅ **TestGetRoomsAPI_CanJoinAccuracy**
- 验证can_join字段准确性

✅ **TestGetRoomsAPI_PlayerInfoAccuracy**
- 验证玩家信息正确性

### 6. 边界条件测试 (4个)

✅ **TestGetRoomsAPI_LargeNumberOfRooms**
- 测试100个房间的性能

✅ **TestGetRoomsAPI_ConcurrentAccess**
- 测试20个并发请求

✅ **TestGetRoomsAPI_DynamicRoomUpdates**
- 测试房间状态动态更新

✅ **TestGetRoomsAPI_RoomClosureHandling**
- 测试房间关闭后的列表更新

### 7. 边缘情况测试 (3个)

✅ **TestGetRoomsAPI_MinimumLimit**
- 测试limit=1的情况

✅ **TestGetRoomsAPI_MaximumLimit**
- 测试limit=50的情况

✅ **TestGetRoomsAPI_NegativeOrZeroPage**
- 测试负数和零页码

### 8. 集成测试 (3个)

✅ **TestGetRoomsAPI_CompleteRoomLifecycle**
- 测试完整的房间生命周期

✅ **TestGetRoomsAPI_MultiRoomMultiStatus**
- 测试多房间多状态综合场景

✅ **TestGetRoomsAPI_IntegrationWithCreateRoom**
- 测试与创建房间API的集成

---

## 关键特性

### 排序规则
1. waiting状态房间优先
2. 相同状态下，玩家数多的优先
3. 相同玩家数下，创建时间新的优先

### 状态说明
- **waiting**: 1-3人，可加入
- **ready**: 4人已满，可以开始游戏
- **playing**: 游戏进行中，不可加入
- **closed**: 房间已关闭，不在列表中

### 分页规则
- 默认page=1, limit=12
- page必须 > 0
- limit范围: 1-50
- 页码超出范围返回空数组

---

## 使用示例

### cURL 示例

```bash
# 获取第1页，默认每页12个
curl -X GET "http://localhost:8080/api/rooms" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取第2页，每页10个
curl -X GET "http://localhost:8080/api/rooms?page=2&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 只获取等待中的房间
curl -X GET "http://localhost:8080/api/rooms?status=waiting" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### TypeScript/React 示例

```typescript
// API调用
async function getRoomList(
  page: number = 1,
  limit: number = 12,
  status?: 'waiting' | 'ready' | 'playing'
): Promise<RoomListResponse> {
  const params: any = { page, limit };
  if (status) params.status = status;

  const response = await axios.get('/api/rooms', {
    params,
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });

  return response.data;
}

// 组件使用
const RoomList: React.FC = () => {
  const [rooms, setRooms] = useState<RoomInfo[]>([]);
  
  useEffect(() => {
    getRoomList(1, 12, 'waiting').then(data => {
      setRooms(data.rooms);
    });
  }, []);
  
  return (
    <div>
      {rooms.map(room => (
        <RoomCard key={room.id} room={room} />
      ))}
    </div>
  );
};
```

---

## 性能指标

### 响应时间
- 正常情况: < 100ms
- 100个房间: < 1s
- 并发20个请求: < 5s

### 测试执行时间
- 单个测试: 0.1-12秒
- 全部27个测试: ~40秒
- 性能测试占用时间最长

---

## 错误处理

### 401 Unauthorized
```json
{
  "error": "missing_token",
  "message": "Authorization header is required"
}
```

### 500 Internal Server Error
```json
{
  "error": "room_list_failed",
  "message": "error message"
}
```

---

## 相关资源

### API文档
- [API-Documentation.md](../API-Documentation.md) - 完整的API文档
- [Technical-Documentation.md](../Technical-Documentation.md) - 技术文档

### 相关API
- `POST /api/rooms` - 创建房间
  - [CREATE_ROOM_TEST_CASES.md](./CREATE_ROOM_TEST_CASES.md)
- `GET /api/rooms/:id` - 获取单个房间
- `GET /api/rooms/my` - 获取当前用户的房间
- `POST /api/rooms/:id/join` - 加入房间
- `POST /api/rooms/:id/leave` - 离开房间
- `POST /api/rooms/:id/start` - 开始游戏

### 代码位置
- **Handler**: `backend/handlers/room.go:94-141`
- **Service**: `backend/room/service.go:302-405`
- **Tests**: `backend/handlers/get_rooms_api_test.go`

---

## 维护指南

### 添加新测试
1. 在 `get_rooms_api_test.go` 中添加测试函数
2. 使用命名规范: `TestGetRoomsAPI_TestName`
3. 更新 `GET_ROOMS_TEST_CASES.md` 文档
4. 更新测试总结

### 修改API
1. 修改handler或service代码
2. 运行所有测试确保不破坏现有功能
3. 添加新的测试用例覆盖新功能
4. 更新相关文档

### 性能优化
1. 关注 `TestGetRoomsAPI_LargeNumberOfRooms` 的执行时间
2. 监控并发测试的稳定性
3. 优化排序和过滤算法

---

## 常见问题

### Q: 为什么有些测试执行时间较长？
A: 大数据测试（100个房间）和并发测试需要创建大量测试数据，执行时间相对较长是正常的。

### Q: 如何只运行快速测试？
A: 使用特定的测试名称过滤，避免运行大数据和并发测试。

### Q: 测试之间会互相影响吗？
A: 不会。每个测试使用独立的用户名前缀，测试数据相互隔离。

### Q: 为什么没有测试limit>50的情况？
A: API实现限制最大limit为50，超过会使用默认值12。这在 `TestGetRoomsAPI_InvalidPaginationParams` 中已测试。

---

## 贡献指南

如果你想为测试添加新的用例：

1. **确保测试独立**: 不依赖其他测试的执行
2. **使用唯一命名**: 为测试用户使用唯一的前缀
3. **清晰的断言**: 每个断言都有清晰的含义
4. **适当的注释**: 复杂逻辑添加注释
5. **更新文档**: 同步更新测试文档

---

## 总结

GET /api/rooms API 的测试套件提供了全面的测试覆盖，包括：
- ✅ 27个测试用例全部通过
- ✅ 代码覆盖率 > 95%
- ✅ 包含性能和并发测试
- ✅ 完善的文档和示例

测试确保了API的稳定性、可靠性和性能，为生产环境部署提供了信心。



