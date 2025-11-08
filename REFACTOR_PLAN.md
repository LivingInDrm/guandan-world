# 游戏状态管理重构方案

## 一、问题根因总结

当前系统反复出现页面逻辑跳转不对的根本原因：

1. **状态机设计缺失** - 前后端没有统一的状态定义和转移规则
2. **消息时序不可靠** - 多层异步导致消息顺序不确定
3. **状态转移条件脆弱** - 过于严格的阶段检查，缺乏兜底机制
4. **数据结构不一致** - player_view 等消息在不同来源有不同结构
5. **容错机制缺失** - 没有状态同步、超时恢复等机制

## 二、重构原则

✅ **以现有接口为准** - 尽量使用已实现的后端接口
✅ **渐进式改进** - 分阶段实施，每阶段都能独立验证
✅ **向后兼容** - 不破坏现有功能
✅ **最小化改动** - 优先前端改进，后端只做必要调整

## 三、核心重构方案

### 阶段 1: 建立统一状态机（优先级：高）

#### 1.1 定义统一的游戏阶段

**新增共享类型定义** (前后端共享):

```typescript
// 统一的游戏阶段枚举
enum GamePhase {
  // 房间阶段
  ROOM_WAITING = 'room_waiting',           // 等待玩家加入
  ROOM_READY = 'room_ready',               // 4人已满，可以开始
  
  // 游戏准备阶段
  GAME_PREPARING = 'game_preparing',       // 倒计时准备中
  
  // 游戏进行阶段
  DEAL_PREPARING = 'deal_preparing',       // 准备发牌（短暂）
  TRIBUTE_PHASE = 'tribute_phase',         // 上贡阶段
  PLAYING_PHASE = 'playing_phase',         // 出牌阶段
  
  // 结算阶段
  DEAL_SETTLING = 'deal_settling',         // 局结算
  MATCH_SETTLING = 'match_settling',       // 比赛结算
  MATCH_FINISHED = 'match_finished'        // 比赛结束
}

// 阶段转移映射表
const PHASE_TRANSITIONS = {
  [GamePhase.ROOM_WAITING]: [GamePhase.ROOM_READY],
  [GamePhase.ROOM_READY]: [GamePhase.GAME_PREPARING, GamePhase.ROOM_WAITING],
  [GamePhase.GAME_PREPARING]: [GamePhase.DEAL_PREPARING],
  [GamePhase.DEAL_PREPARING]: [GamePhase.TRIBUTE_PHASE, GamePhase.PLAYING_PHASE],
  [GamePhase.TRIBUTE_PHASE]: [GamePhase.PLAYING_PHASE],
  [GamePhase.PLAYING_PHASE]: [GamePhase.DEAL_SETTLING],
  [GamePhase.DEAL_SETTLING]: [GamePhase.DEAL_PREPARING, GamePhase.MATCH_SETTLING],
  [GamePhase.MATCH_SETTLING]: [GamePhase.MATCH_FINISHED],
  [GamePhase.MATCH_FINISHED]: [GamePhase.ROOM_WAITING]
};
```

#### 1.2 后端改动：统一 player_view 数据结构

**目标**: 所有地方发送的 player_view 使用同一数据结构

**改动位置 1**: `handlers/room.go` 的 `runGamePrepareSequence()`

```go
// 修改初始 player_view 的结构，与 SDK 发送的保持一致
initialView := &websocket.WSMessage{
    Type: "player_view",
    Data: map[string]interface{}{
        "player_seat": i,
        "hand":        []interface{}{},
        "can_play":    false,
        "is_my_turn":  false,
        "phase":       "game_preparing",  // ✅ 新增：明确告知当前阶段
        "game_state": map[string]interface{}{
            "team_levels": []int{2, 2},
            "current_deal": map[string]interface{}{
                "status": "preparing",
                "tribute_phase": nil,  // 明确标识还没上贡
            },
        },
    },
    PlayerID: player.ID,
}
```

**改动位置 2**: `game/driver_service.go` 的 `sendPlayerViews()`

```go
// 统一 player_view 结构
wsMessage := &websocket.WSMessage{
    Type: websocket.MSG_PLAYER_VIEW,
    Data: map[string]interface{}{
        "player_seat": playerSeat,
        "hand":        playerView.PlayerCards,
        "can_play":    playerView.CanPlay,
        "is_my_turn":  playerView.IsMyTurn,
        "phase":       inferPhaseFromGameState(playerView.GameState),  // ✅ 新增
        "game_state":  playerView.GameState,
    },
    PlayerID: playerID,
}
```

**新增辅助函数**:

```go
// inferPhaseFromGameState 从游戏状态推断当前阶段
func inferPhaseFromGameState(gameState *sdk.GameState) string {
    if gameState == nil || gameState.CurrentMatch == nil {
        return "room_waiting"
    }
    
    currentDeal := gameState.CurrentMatch.CurrentDeal
    if currentDeal == nil {
        // 比赛已开始但没有当前局
        if gameState.Status == sdk.MatchStatusFinished {
            return "match_finished"
        }
        return "deal_preparing"
    }
    
    // 根据 DealStatus 映射阶段
    switch currentDeal.Status {
    case sdk.DealStatusPreparing:
        return "deal_preparing"
    case sdk.DealStatusTribute:
        return "tribute_phase"
    case sdk.DealStatusPlaying:
        return "playing_phase"
    case sdk.DealStatusEnded:
        // 判断是局结算还是比赛结算
        if gameState.Status == sdk.MatchStatusFinished {
            return "match_settling"
        }
        return "deal_settling"
    default:
        return "room_waiting"
    }
}
```

#### 1.3 前端改动：基于 phase 字段的状态管理

**新建状态管理器** `src/services/gameStateManager.ts`:

```typescript
import { GamePhase } from '../types';

export class GameStateManager {
  private currentPhase: GamePhase = GamePhase.ROOM_WAITING;
  private phaseStartTime: number = Date.now();
  private phaseTimeout: NodeJS.Timeout | null = null;
  
  // 状态转移
  transitionTo(newPhase: GamePhase, reason: string): boolean {
    // 验证转移是否合法
    const allowedTransitions = PHASE_TRANSITIONS[this.currentPhase] || [];
    if (!allowedTransitions.includes(newPhase)) {
      console.warn(`Invalid phase transition: ${this.currentPhase} -> ${newPhase}`);
      // ✅ 关键：允许某些"跳跃式"转移作为恢复机制
      if (!this.isRecoveryTransition(newPhase)) {
        return false;
      }
    }
    
    console.log(`[StateManager] ${this.currentPhase} -> ${newPhase} (${reason})`);
    
    // 清理旧阶段的超时
    this.clearPhaseTimeout();
    
    // 更新状态
    this.currentPhase = newPhase;
    this.phaseStartTime = Date.now();
    
    // 设置新阶段的超时
    this.setupPhaseTimeout(newPhase);
    
    return true;
  }
  
  // 检查是否是恢复性转移（允许跳过中间状态）
  private isRecoveryTransition(newPhase: GamePhase): boolean {
    // 从后端同步状态时，允许直接跳转到任何阶段
    return [
      GamePhase.TRIBUTE_PHASE,
      GamePhase.PLAYING_PHASE,
      GamePhase.DEAL_SETTLING
    ].includes(newPhase);
  }
  
  // 为每个阶段设置合理的超时
  private setupPhaseTimeout(phase: GamePhase): void {
    const timeouts: Record<GamePhase, number> = {
      [GamePhase.GAME_PREPARING]: 10000,  // 10秒
      [GamePhase.DEAL_PREPARING]: 5000,   // 5秒
      [GamePhase.TRIBUTE_PHASE]: 60000,   // 60秒
      [GamePhase.PLAYING_PHASE]: -1,      // 无超时（游戏中）
      [GamePhase.DEAL_SETTLING]: 30000,   // 30秒
      // ... 其他阶段
    };
    
    const timeout = timeouts[phase];
    if (timeout > 0) {
      this.phaseTimeout = setTimeout(() => {
        console.error(`Phase timeout: ${phase}`);
        // ✅ 超时时主动同步状态
        this.onPhaseTimeout(phase);
      }, timeout);
    }
  }
  
  private clearPhaseTimeout(): void {
    if (this.phaseTimeout) {
      clearTimeout(this.phaseTimeout);
      this.phaseTimeout = null;
    }
  }
  
  // 超时处理回调（由外部设置）
  onPhaseTimeout: (phase: GamePhase) => void = () => {};
  
  getCurrentPhase(): GamePhase {
    return this.currentPhase;
  }
  
  getPhaseElapsedTime(): number {
    return Date.now() - this.phaseStartTime;
  }
}
```

**修改 GamePage.tsx**:

```typescript
import { GameStateManager } from '../../services/gameStateManager';

const GamePage: React.FC = () => {
  // 使用状态管理器
  const stateManager = useRef(new GameStateManager());
  
  // 当前阶段直接从 stateManager 获取
  const [currentPhase, setCurrentPhase] = useState(stateManager.current.getCurrentPhase());
  
  // 设置超时回调
  useEffect(() => {
    stateManager.current.onPhaseTimeout = (phase) => {
      console.error(`Phase ${phase} timeout, requesting state sync`);
      // 主动同步状态
      queryGameState();
    };
  }, []);
  
  // ✅ 处理 player_view 消息：基于 phase 字段
  const handlePlayerView = (message: WSMessage) => {
    const data = message.data;
    if (!data) return;
    
    // 更新游戏数据
    setPlayerSeat(data.player_seat);
    setPlayerHand(data.hand || []);
    setCanPlay(data.can_play || false);
    setMyTurn(data.is_my_turn || false);
    setGameState(data.game_state);
    
    // ✅ 关键改动：基于后端明确告知的 phase 进行状态转移
    if (data.phase) {
      const success = stateManager.current.transitionTo(
        data.phase as GamePhase,
        'player_view message'
      );
      
      if (success) {
        setCurrentPhase(data.phase as GamePhase);
      }
    }
  };
  
  // ✅ 主动查询游戏状态（容错机制）
  const queryGameState = async () => {
    if (!roomId) return;
    
    try {
      const response = await apiClient.getGameState(roomId);
      if (response.success && response.data) {
        const { phase, player_view } = response.data;
        
        // 强制同步到后端状态
        stateManager.current.transitionTo(phase, 'state sync');
        setCurrentPhase(phase);
        
        // 更新玩家视图
        if (player_view) {
          handlePlayerView({ type: 'player_view', data: player_view, timestamp: new Date().toISOString() });
        }
      }
    } catch (error) {
      console.error('Failed to query game state:', error);
    }
  };
  
  // ✅ 处理游戏准备：不再直接转移状态，等待 player_view
  const handleGamePrepare = (message: WSMessage) => {
    // 只显示倒计时界面，不转移状态
    // 状态转移由 player_view 的 phase 字段驱动
    setCountdown(3);
  };
  
  // 其他逻辑...
};
```

### 阶段 2: 新增状态同步接口（优先级：高）

#### 2.1 后端新增接口

**新增 HTTP 接口**: `GET /api/game/state/:room_id`

```go
// handlers/game_driver.go 新增
type GameStateResponse struct {
    Phase       string                 `json:"phase"`        // 当前阶段
    PlayerView  map[string]interface{} `json:"player_view"`  // 玩家视图
    RoomInfo    map[string]interface{} `json:"room_info"`    // 房间信息
    Timestamp   time.Time              `json:"timestamp"`
}

func (h *GameDriverHandler) GetGameState(c *gin.Context) {
    roomID := c.Param("room_id")
    userID := c.GetString("user_id")  // 从认证中间件获取
    
    // 1. 获取房间信息
    room, err := h.roomService.GetRoom(roomID)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: "Room not found"})
        return
    }
    
    // 2. 确定玩家座位
    playerSeat := -1
    for i, player := range room.Players {
        if player != nil && player.ID == userID {
            playerSeat = i
            break
        }
    }
    
    if playerSeat == -1 {
        c.JSON(http.StatusForbidden, ErrorResponse{Error: "Player not in room"})
        return
    }
    
    // 3. 根据房间状态确定阶段
    var phase string
    var playerView map[string]interface{}
    
    switch room.Status {
    case room.RoomStatusWaiting:
        phase = "room_waiting"
    case room.RoomStatusReady:
        phase = "room_ready"
    case room.RoomStatusPlaying:
        // 游戏进行中，从游戏引擎获取详细状态
        gameStatus, err := h.driverService.GetGameStatus(roomID)
        if err != nil {
            // 游戏还未真正启动，处于准备阶段
            phase = "game_preparing"
        } else {
            // 从引擎获取玩家视图
            engine := h.driverService.GetEngine(roomID)
            if engine != nil {
                view := engine.GetPlayerView(playerSeat)
                phase = inferPhaseFromGameState(view.GameState)
                playerView = convertPlayerViewToMap(view)
            } else {
                phase = "game_preparing"
            }
        }
    default:
        phase = "room_waiting"
    }
    
    // 4. 返回完整状态
    c.JSON(http.StatusOK, GameStateResponse{
        Phase:      phase,
        PlayerView: playerView,
        RoomInfo: map[string]interface{}{
            "room_id": roomID,
            "status":  room.Status.String(),
            "players": room.Players,
        },
        Timestamp: time.Now(),
    })
}
```

**在 main.go 中注册路由**:

```go
// 添加到已有的游戏路由组
gameDriver.GET("/state/:room_id", gameDriverHandler.GetGameState)
```

#### 2.2 前端调用状态同步接口

**新增 API 方法** `src/services/api.ts`:

```typescript
export class ApiClient {
  // ... 已有方法
  
  // 新增：获取游戏状态
  async getGameState(roomId: string): Promise<ApiResponse<GameStateResponse>> {
    return this.request<GameStateResponse>({
      method: 'GET',
      url: `/api/game/state/${roomId}`,
    });
  }
}

// 类型定义
interface GameStateResponse {
  phase: GamePhase;
  player_view?: PlayerView;
  room_info: RoomInfo;
  timestamp: string;
}
```

**在 GamePage 中使用**:

```typescript
// 定期轮询（可选，作为兜底）
useEffect(() => {
  if (!roomId || currentPhase === GamePhase.ROOM_WAITING) return;
  
  const pollInterval = setInterval(() => {
    queryGameState();  // 每 10 秒轮询一次
  }, 10000);
  
  return () => clearInterval(pollInterval);
}, [roomId, currentPhase]);

// WebSocket 重连后立即同步
useEffect(() => {
  if (isConnected) {
    queryGameState();  // 连接成功后立即同步状态
  }
}, [isConnected]);

// 页面恢复焦点时同步（处理用户切换标签页的情况）
useEffect(() => {
  const handleVisibilityChange = () => {
    if (document.visibilityState === 'visible') {
      queryGameState();
    }
  };
  
  document.addEventListener('visibilitychange', handleVisibilityChange);
  return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
}, []);
```

### 阶段 3: 优化消息发送顺序（优先级：中）

#### 3.1 后端改动：序列化关键消息发送

**修改 `handlers/room.go` 的 `runGamePrepareSequence()`**:

```go
func (h *RoomHandler) runGamePrepareSequence(roomID string, players [4]*room.Player) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic in runGamePrepareSequence: %v", r)
            h.revertAndNotifyError(roomID, r)
        }
    }()
    
    // ✅ 改进：批量构建所有初始消息，然后按顺序发送
    messages := h.buildGameStartSequence(roomID, players)
    
    // 依次发送，确保顺序
    for _, msg := range messages {
        h.wsManager.BroadcastToRoom(roomID, msg.message)
        
        // 如果有延迟要求，等待
        if msg.delayMs > 0 {
            time.Sleep(time.Duration(msg.delayMs) * time.Millisecond)
        }
    }
    
    // ✅ 确保所有初始消息发送完毕后，再启动游戏引擎
    time.Sleep(100 * time.Millisecond)  // 给消息发送留出时间
    
    // 启动游戏引擎
    sdkPlayers := h.convertToSDKPlayers(players)
    if err := h.driverService.StartGameWithDriver(roomID, sdkPlayers); err != nil {
        log.Printf("Failed to start game engine: %v", err)
        h.revertAndNotifyError(roomID, err)
    }
}

type scheduledMessage struct {
    message *websocket.WSMessage
    delayMs int
}

func (h *RoomHandler) buildGameStartSequence(roomID string, players [4]*room.Player) []scheduledMessage {
    messages := []scheduledMessage{
        // 1. 游戏准备通知
        {
            message: &websocket.WSMessage{
                Type: "game_prepare",
                Data: map[string]interface{}{
                    "room_id": roomID,
                },
                Timestamp: time.Now(),
            },
            delayMs: 0,
        },
    }
    
    // 2. 倒计时消息
    for i := 3; i > 0; i-- {
        messages = append(messages, scheduledMessage{
            message: &websocket.WSMessage{
                Type: "countdown",
                Data: map[string]interface{}{
                    "countdown": i,
                    "room_id":   roomID,
                },
                Timestamp: time.Now().Add(time.Duration(3-i) * time.Second),
            },
            delayMs: 1000,
        })
    }
    
    // 3. 游戏开始通知
    messages = append(messages, scheduledMessage{
        message: &websocket.WSMessage{
            Type: "game_begin",
            Data: map[string]interface{}{
                "room_id": roomID,
            },
            Timestamp: time.Now().Add(3 * time.Second),
        },
        delayMs: 0,
    })
    
    // 4. 为每个玩家发送初始视图（改用广播，避免顺序问题）
    // ✅ 注意：这里发送的是一个特殊的"初始化"消息，真正的 player_view 由 SDK 发送
    messages = append(messages, scheduledMessage{
        message: &websocket.WSMessage{
            Type: "game_initialized",
            Data: map[string]interface{}{
                "room_id": roomID,
                "phase":   "game_preparing",
                "message": "Game engine is starting...",
            },
            Timestamp: time.Now().Add(3 * time.Second),
        },
        delayMs: 0,
    })
    
    return messages
}
```

#### 3.2 SDK 事件发送优化

**修改 `game/driver_service.go` 的 `WebSocketObserver`**:

```go
type WebSocketObserver struct {
    roomID          string
    wsManager       WSManagerInterface
    engine          sdk.GameEngineInterface
    eventQueue      chan *sdk.GameEvent  // ✅ 新增：事件队列
    stopChan        chan struct{}
    isStarted       bool
    mu              sync.Mutex
}

func NewWebSocketObserver(roomID string, wsManager WSManagerInterface, engine sdk.GameEngineInterface) *WebSocketObserver {
    obs := &WebSocketObserver{
        roomID:     roomID,
        wsManager:  wsManager,
        engine:     engine,
        eventQueue: make(chan *sdk.GameEvent, 100),  // 缓冲队列
        stopChan:   make(chan struct{}),
        isStarted:  false,
    }
    
    // 启动事件处理协程
    go obs.processEventQueue()
    
    return obs
}

// OnGameEvent 接收事件，放入队列
func (wso *WebSocketObserver) OnGameEvent(event *sdk.GameEvent) {
    select {
    case wso.eventQueue <- event:
        // 事件已入队
    default:
        log.Printf("Warning: Event queue full, dropping event: %s", event.Type)
    }
}

// processEventQueue 串行处理事件队列
func (wso *WebSocketObserver) processEventQueue() {
    for {
        select {
        case event := <-wso.eventQueue:
            wso.handleEvent(event)  // 串行处理
        case <-wso.stopChan:
            return
        }
    }
}

// handleEvent 处理单个事件
func (wso *WebSocketObserver) handleEvent(event *sdk.GameEvent) {
    // 转换并广播事件
    wsMessage := &websocket.WSMessage{
        Type: websocket.MSG_GAME_EVENT,
        Data: map[string]interface{}{
            "event_type":  string(event.Type),
            "event_data":  event.Data,
            "timestamp":   event.Timestamp,
            "player_seat": event.PlayerSeat,
        },
        Timestamp: event.Timestamp,
    }
    
    wso.wsManager.BroadcastToRoom(wso.roomID, wsMessage)
    
    // 对于需要发送 player_view 的事件
    if wso.shouldSendPlayerViews(event.Type) {
        // ✅ 关键：同步发送 player_view，确保顺序
        wso.sendPlayerViewsSync(event.Type)
    }
}

// sendPlayerViewsSync 同步发送 player_view（确保顺序）
func (wso *WebSocketObserver) sendPlayerViewsSync(eventType sdk.GameEventType) {
    if wso.engine == nil {
        return
    }
    
    // 串行为每个玩家发送视图
    for playerSeat := 0; playerSeat < 4; playerSeat++ {
        playerView := wso.engine.GetPlayerView(playerSeat)
        if playerView == nil {
            continue
        }
        
        gameState := playerView.GameState
        if gameState == nil || gameState.CurrentMatch == nil {
            continue
        }
        
        currentMatch := gameState.CurrentMatch
        if playerSeat >= len(currentMatch.Players) || currentMatch.Players[playerSeat] == nil {
            continue
        }
        
        playerID := currentMatch.Players[playerSeat].ID
        phase := inferPhaseFromGameState(gameState)
        
        // 构建统一的 player_view
        wsMessage := &websocket.WSMessage{
            Type: websocket.MSG_PLAYER_VIEW,
            Data: map[string]interface{}{
                "player_seat": playerSeat,
                "hand":        playerView.PlayerCards,
                "can_play":    playerView.CanPlay,
                "is_my_turn":  playerView.IsMyTurn,
                "phase":       phase,  // ✅ 关键字段
                "game_state":  gameState,
            },
            Timestamp: time.Now(),
            PlayerID:  playerID,
        }
        
        // 同步发送（等待发送完成）
        if err := wso.wsManager.SendToPlayer(playerID, wsMessage); err != nil {
            log.Printf("Failed to send player view to %s: %v", playerID, err)
        }
        
        // 短暂延迟，避免消息过于密集
        time.Sleep(10 * time.Millisecond)
    }
}
```

### 阶段 4: 增强前端容错机制（优先级：中）

#### 4.1 消息去重和排序

**新建 `src/services/messageQueue.ts`**:

```typescript
interface QueuedMessage {
  message: WSMessage;
  receivedAt: number;
  processed: boolean;
}

export class MessageQueue {
  private queue: QueuedMessage[] = [];
  private processedIds = new Set<string>();
  private maxQueueSize = 100;
  
  // 添加消息到队列
  enqueue(message: WSMessage): void {
    // 生成消息ID（用于去重）
    const messageId = this.generateMessageId(message);
    
    // 去重检查
    if (this.processedIds.has(messageId)) {
      console.log(`[MessageQueue] Duplicate message ignored: ${message.type}`);
      return;
    }
    
    // 添加到队列
    this.queue.push({
      message,
      receivedAt: Date.now(),
      processed: false
    });
    
    // 限制队列大小
    if (this.queue.length > this.maxQueueSize) {
      const removed = this.queue.shift();
      if (removed) {
        this.processedIds.delete(this.generateMessageId(removed.message));
      }
    }
  }
  
  // 处理队列中的消息
  processMessages(handler: (message: WSMessage) => void): void {
    // 按接收时间排序
    this.queue.sort((a, b) => a.receivedAt - b.receivedAt);
    
    // 处理未处理的消息
    for (const item of this.queue) {
      if (!item.processed) {
        handler(item.message);
        item.processed = true;
        this.processedIds.add(this.generateMessageId(item.message));
      }
    }
    
    // 清理已处理的旧消息
    const now = Date.now();
    this.queue = this.queue.filter(
      item => !item.processed || (now - item.receivedAt < 60000)  // 保留1分钟
    );
  }
  
  private generateMessageId(message: WSMessage): string {
    // 基于消息类型、时间戳和部分数据生成ID
    const dataHash = JSON.stringify(message.data).substring(0, 50);
    return `${message.type}-${message.timestamp}-${dataHash}`;
  }
  
  clear(): void {
    this.queue = [];
    this.processedIds.clear();
  }
}
```

**在 GamePage 中使用**:

```typescript
const GamePage: React.FC = () => {
  const messageQueue = useRef(new MessageQueue());
  
  // WebSocket 消息处理
  useEffect(() => {
    if (!roomId || !isConnected) return;
    
    const handleMessage = (message: WSMessage) => {
      // 先入队
      messageQueue.current.enqueue(message);
      
      // 然后处理队列
      messageQueue.current.processMessages((msg) => {
        handleMessageProcessing(msg);
      });
    };
    
    wsClient.on(WS_MESSAGE_TYPES.GAME_EVENT, handleMessage);
    wsClient.on(WS_MESSAGE_TYPES.PLAYER_VIEW, handleMessage);
    wsClient.on(WS_MESSAGE_TYPES.GAME_PREPARE, handleMessage);
    // ... 其他消息类型
    
    return () => {
      wsClient.off(WS_MESSAGE_TYPES.GAME_EVENT, handleMessage);
      // ... 清理其他监听器
    };
  }, [roomId, isConnected]);
};
```

#### 4.2 状态不一致检测

**新建 `src/services/stateValidator.ts`**:

```typescript
export class StateValidator {
  // 验证状态转移是否合理
  validateTransition(
    currentPhase: GamePhase,
    newPhase: GamePhase,
    gameState: any
  ): { valid: boolean; reason?: string } {
    // 1. 基本转移规则检查
    const allowedTransitions = PHASE_TRANSITIONS[currentPhase] || [];
    if (!allowedTransitions.includes(newPhase)) {
      return {
        valid: false,
        reason: `Invalid transition: ${currentPhase} -> ${newPhase}`
      };
    }
    
    // 2. 游戏状态一致性检查
    if (newPhase === GamePhase.TRIBUTE_PHASE) {
      if (!gameState?.current_deal?.tribute_phase) {
        return {
          valid: false,
          reason: 'Tribute phase indicated but no tribute data'
        };
      }
    }
    
    if (newPhase === GamePhase.PLAYING_PHASE) {
      if (!gameState?.current_deal?.current_trick) {
        console.warn('Playing phase but no current trick (might be starting)');
      }
    }
    
    return { valid: true };
  }
  
  // 检测状态异常
  detectAnomalies(
    currentPhase: GamePhase,
    phaseStartTime: number,
    gameState: any
  ): string[] {
    const anomalies: string[] = [];
    const elapsedMs = Date.now() - phaseStartTime;
    
    // 检测阶段停留时间过长
    const maxDurations: Record<GamePhase, number> = {
      [GamePhase.GAME_PREPARING]: 15000,
      [GamePhase.DEAL_PREPARING]: 10000,
      [GamePhase.TRIBUTE_PHASE]: 120000,
      // ...
    };
    
    const maxDuration = maxDurations[currentPhase];
    if (maxDuration && elapsedMs > maxDuration) {
      anomalies.push(`Phase ${currentPhase} duration exceeded: ${elapsedMs}ms`);
    }
    
    // 检测数据完整性
    if (currentPhase === GamePhase.PLAYING_PHASE) {
      if (!gameState?.current_deal) {
        anomalies.push('Playing phase but no current deal');
      }
    }
    
    return anomalies;
  }
}
```

### 阶段 5: 监控和调试工具（优先级：低）

#### 5.1 开发环境状态面板

**新建 `src/components/debug/StateDebugPanel.tsx`**:

```typescript
import React, { useEffect, useState } from 'react';
import { useGameStore } from '../../store/gameStore';

export const StateDebugPanel: React.FC = () => {
  const { currentPhase, gameState } = useGameStore();
  const [messages, setMessages] = useState<any[]>([]);
  const [phaseHistory, setPhaseHistory] = useState<string[]>([]);
  
  // 只在开发环境显示
  if (import.meta.env.PROD) {
    return null;
  }
  
  return (
    <div className="fixed bottom-0 right-0 w-96 bg-black bg-opacity-90 text-white text-xs p-4 max-h-96 overflow-y-auto">
      <h3 className="font-bold mb-2">🐛 State Debug Panel</h3>
      
      {/* 当前状态 */}
      <div className="mb-4">
        <div className="text-yellow-400">Current Phase:</div>
        <div className="font-mono">{currentPhase}</div>
      </div>
      
      {/* 阶段历史 */}
      <div className="mb-4">
        <div className="text-yellow-400">Phase History:</div>
        <div className="font-mono text-xs">
          {phaseHistory.slice(-5).map((p, i) => (
            <div key={i}>{p}</div>
          ))}
        </div>
      </div>
      
      {/* 最近消息 */}
      <div className="mb-4">
        <div className="text-yellow-400">Recent Messages:</div>
        <div className="font-mono text-xs">
          {messages.slice(-10).map((msg, i) => (
            <div key={i} className="mb-1">
              <span className="text-green-400">{msg.type}</span>
              {msg.data?.phase && (
                <span className="text-blue-400"> [phase: {msg.data.phase}]</span>
              )}
            </div>
          ))}
        </div>
      </div>
      
      {/* 游戏状态摘要 */}
      <div>
        <div className="text-yellow-400">Game State:</div>
        <div className="font-mono text-xs">
          {gameState?.current_deal && (
            <>
              <div>Deal Status: {gameState.current_deal.status}</div>
              <div>Tribute: {gameState.current_deal.tribute_phase ? 'Yes' : 'No'}</div>
              <div>Trick: {gameState.current_deal.current_trick ? 'Active' : 'None'}</div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};
```

## 四、实施计划

### 第一周：阶段 1 + 阶段 2

**目标**: 建立统一状态机 + 状态同步接口

- [ ] Day 1-2: 定义统一类型，修改后端 player_view 结构
- [ ] Day 3-4: 实现前端 GameStateManager
- [ ] Day 5-6: 新增 GET /api/game/state/:room_id 接口
- [ ] Day 7: 集成测试，验证状态同步

**验收标准**:
- 所有 player_view 消息包含 phase 字段
- 前端可通过 queryGameState() 主动同步状态
- 状态转移基于 phase 字段，不再依赖严格的顺序

### 第二周：阶段 3

**目标**: 优化消息发送顺序

- [ ] Day 1-3: 重构 runGamePrepareSequence，序列化消息发送
- [ ] Day 4-5: 优化 WebSocketObserver 事件队列
- [ ] Day 6-7: 压力测试，验证消息顺序稳定性

**验收标准**:
- 游戏启动时消息按预期顺序到达
- 倒计时结束后不再卡住
- 10 次连续测试无状态转移失败

### 第三周：阶段 4 + 阶段 5

**目标**: 容错机制 + 调试工具

- [ ] Day 1-2: 实现前端消息队列和去重
- [ ] Day 3-4: 实现状态验证器
- [ ] Day 5-6: 开发调试面板
- [ ] Day 7: 全面测试，修复遗留问题

**验收标准**:
- 重复消息被正确过滤
- 状态异常能被检测并主动恢复
- 调试面板能清晰显示状态流转

## 五、回滚计划

如果重构出现问题，可以按以下步骤回滚：

1. **阶段 1 回滚**: 移除 phase 字段，恢复原有的状态判断逻辑
2. **阶段 2 回滚**: 删除新增的状态同步接口
3. **阶段 3 回滚**: 恢复原有的并行消息发送逻辑
4. **阶段 4/5 回滚**: 移除辅助工具，不影响核心功能

## 六、测试策略

### 单元测试

- [ ] GameStateManager 状态转移逻辑
- [ ] MessageQueue 消息去重
- [ ] StateValidator 验证规则

### 集成测试

- [ ] 完整游戏流程（无 AI）
- [ ] 4 AI 玩家完整对局
- [ ] WebSocket 断线重连恢复

### 压力测试

- [ ] 10 个房间并发游戏
- [ ] 快速点击开始游戏（防止重复触发）
- [ ] 网络延迟模拟（300ms+）

### 兼容性测试

- [ ] 不同浏览器（Chrome, Firefox, Safari）
- [ ] 移动端浏览器
- [ ] 不同网络环境

## 七、预期效果

重构完成后，系统应达到以下目标：

✅ **状态转移可靠**: 任何情况下状态都能正确转移，不会卡住
✅ **容错性强**: 消息丢失、乱序、重复都能正确处理
✅ **可观测性好**: 开发者能清楚看到状态流转和消息传递
✅ **可维护性高**: 状态管理逻辑集中，易于理解和修改
✅ **用户体验佳**: 无卡顿、无跳转失败、响应及时

## 八、附录：关键接口汇总

### 现有接口

```
# 房间管理
POST   /api/rooms                    - 创建房间
GET    /api/rooms                    - 房间列表
GET    /api/rooms/:id                - 房间详情
POST   /api/rooms/:id/join           - 加入房间
POST   /api/rooms/:id/leave          - 离开房间
POST   /api/rooms/:id/start          - 开始游戏（触发倒计时和引擎启动）
GET    /api/rooms/my                 - 我的房间

# 游戏驱动（主要用于测试，也可用于状态查询）
POST   /api/game/driver/start              - 直接启动游戏引擎
POST   /api/game/driver/play-decision      - 提交出牌
POST   /api/game/driver/tribute-select     - 提交贡牌
POST   /api/game/driver/tribute-return     - 提交还贡
GET    /api/game/driver/status/:room_id    - 获取游戏状态（引擎层）
POST   /api/game/driver/stop/:room_id      - 停止游戏
```

### 新增接口

```
# 状态同步（重构新增）
GET    /api/game/state/:room_id      - 获取完整游戏状态（含阶段和玩家视图）
```

### WebSocket 消息类型

```
# 客户端 → 服务器
join_room           - 加入房间
leave_room          - 离开房间

# 服务器 → 客户端
room_update         - 房间更新（玩家加入/离开）
game_prepare        - 游戏准备（倒计时开始）
countdown           - 倒计时（3, 2, 1）
game_begin          - 游戏开始
game_initialized    - 游戏引擎已初始化（新增）
player_view         - 玩家视图（含 phase 字段）✨
game_event          - 游戏事件
game_action         - 需要玩家操作
error               - 错误消息
```

---

**文档版本**: v1.0  
**创建日期**: 2025-01-08  
**最后更新**: 2025-01-08

