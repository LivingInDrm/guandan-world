# 掼蛋在线对战游戏设计文档

## 概述

掼蛋在线对战平台是一个基于WebSocket的实时多人在线游戏系统。系统采用前后端分离架构，后端使用Go语言实现游戏逻辑和WebSocket服务，前端使用React+TypeScript实现用户界面。系统支持用户认证、房间管理、实时游戏对战、断线托管、超时处理等核心功能，提供完整的掼蛋游戏体验，包括发牌、上贡还贡、出牌验证、Deal结算、Match管理等完整流程。

## 架构

### 整体架构

```
┌─────────────────┐    WebSocket/HTTP    ┌─────────────────┐
│   React前端     │ ◄─────────────────► │   Go后端服务    │
│                 │                      │                 │
│ - 用户界面      │                      │ - WebSocket服务 │
│ - 状态管理      │                      │ - 房间管理      │
│ - WebSocket客户端│                      │ - 用户认证      │
└─────────────────┘                      │ - 游戏协调      │
                                         └─────────────────┘
                                                   │
                                                   ▼
                                         ┌─────────────────┐
                                         │   掼蛋游戏SDK   │
                                         │                 │
                                         │ - 游戏引擎      │
                                         │ - 牌规则逻辑    │
                                         │ - Trick/Deal/Match    │
                                         │ - 状态管理      │
                                         └─────────────────┘
                                                   │
                                                   ▼
                                         ┌─────────────────┐
                                         │   内存存储      │
                                         │                 │
                                         │ - 用户会话      │
                                         │ - 房间状态      │
                                         └─────────────────┘
```

### 技术栈

**后端：**
- Go 1.21+
- Gin Web框架
- Gorilla WebSocket
- 现有的掼蛋牌规则SDK

**前端：**
- React 19
- TypeScript
- Vite构建工具
- WebSocket客户端

**部署：**
- Docker容器化
- Docker Compose本地开发

## 组件和接口

### 后端组件

#### 1. 用户认证服务 (AuthService)

```go
type AuthService struct {
    users map[string]*User
    sessions map[string]*Session
}

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Session struct {
    UserID    string    `json:"user_id"`
    Token     string    `json:"token"`
    CreatedAt time.Time `json:"created_at"`
}

```

**接口：**
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `POST /api/logout` - 用户登出

#### 2. 房间管理服务 (RoomService)

```go
type RoomService struct {
    rooms map[string]*Room
    mutex sync.RWMutex
}

type Room struct {
    ID          string             `json:"id"`
    Name        string             `json:"name"`
    Status      RoomStatus         `json:"status"`
    Players     map[int]*Player    `json:"players"` // seat -> player
    Owner       string             `json:"owner"`
    CreatedAt   time.Time          `json:"created_at"`
    GameEngine  *sdk.GameEngine    `json:"game_engine,omitempty"` // 引用SDK游戏引擎
}

type RoomStatus string
const (
    RoomStatusWaiting RoomStatus = "waiting"
    RoomStatusPlaying RoomStatus = "playing"
)

type Player struct {
    UserID   string `json:"user_id"` // 和user_id一致,通过 NewMatch(playerIDs []string) 显式传入 SDK
    Username string `json:"username"`
    Seat     int    `json:"seat"`
    Online   bool   `json:"online"`
}

```

**接口：**
- `GET /api/rooms` - 获取房间列表，支持分页参数（page, limit），默认每页12个房间
- `POST /api/rooms` - 创建房间
- `POST /api/rooms/:id/join` - 加入房间
- `POST /api/rooms/:id/leave` - 离开房间

**房间列表排序策略：**
- 优先显示等待中的房间（status = "waiting"）
- 按已加入人数降序排列
- 游戏进行中的房间排在后面

#### 3. WebSocket连接管理 (WSManager)

```go
type WSManager struct {
    clients    map[string]*WSClient
    rooms      map[string]map[string]*WSClient
    register   chan *WSClient
    unregister chan *WSClient
    broadcast  chan *WSMessage
}

type WSClient struct {
    ID     string
    UserID string
    RoomID string
    Conn   *websocket.Conn
    Send   chan *WSMessage
}

type WSMessage struct {
    Type    string      `json:"type"`
    RoomID  string      `json:"room_id,omitempty"`
    UserID  string      `json:"user_id,omitempty"`
    Data    interface{} `json:"data"`
}
```

**连接管理策略：**
- 自动心跳检测和断线重连
- 游戏进行中禁止新用户加入房间
- 断线用户自动标记为托管状态

#### 4. 游戏协调服务 (GameService)

游戏协调服务作为SDK和外部系统的桥梁，完全基于SDK的事件系统工作：

```go
type GameService struct {
    games     map[string]*sdk.GameEngine  // 使用SDK的GameEngine
    wsManager *WSManager
    mutex     sync.RWMutex
}

// GameService完全依赖SDK的事件系统
func (gs *GameService) CreateGame(roomID string, playerIDs []string) error {
    // 创建SDK游戏引擎
    game := sdk.NewGameEngine()
    
    // 注册事件处理器，将SDK事件转换为WebSocket消息
    game.RegisterEventHandler(sdk.EventPlayerPlayed, gs.handlePlayerPlayed)
    game.RegisterEventHandler(sdk.EventDealEnded, gs.handleDealEnded)
    game.RegisterEventHandler(sdk.EventMatchEnded, gs.handleMatchEnded)
    game.RegisterEventHandler(sdk.EventPlayerTimeout, gs.handlePlayerTimeout)
    // ... 注册所有事件处理器
    
    // 初始化游戏
    err := game.StartMatch(playerIDs)
    if err != nil {
        return err
    }
    
    gs.games[roomID] = game
    return nil
}

// 处理玩家操作 - 完全委托给SDK
func (gs *GameService) HandlePlayerAction(roomID string, action *PlayerAction) error {
    game, exists := gs.games[roomID]
    if !exists {
        return errors.New("game not found")
    }
    
    var events []*sdk.GameEvent
    var err error
    
    // 所有操作都通过SDK接口
    switch action.Type {
    case "start_game":
        events, err = game.StartMatch()
    case "play_cards":
        event, playErr := game.PlayCards(action.PlayerSeat, action.Cards)
        events = []*sdk.GameEvent{event}
        err = playErr
    case "pass":
        event, passErr := game.PassTurn(action.PlayerSeat)
        events = []*sdk.GameEvent{event}
        err = passErr
    case "select_tribute":
        event, tributeErr := game.SelectTribute(action.PlayerSeat, action.Card)
        events = []*sdk.GameEvent{event}
        err = tributeErr
    default:
        return errors.New("unknown action type")
    }
    
    if err != nil {
        return err
    }
    
    // 处理SDK返回的事件
    for _, event := range events {
        gs.handleGameEvent(roomID, event)
    }
    
    return nil
}

// 事件处理器 - 将SDK事件转换为WebSocket消息
func (gs *GameService) handleGameEvent(roomID string, event *sdk.GameEvent) {
    message := &WSMessage{
        Type:   string(event.Type),
        RoomID: roomID,
        Data:   event.Data,
    }
    gs.wsManager.BroadcastToRoom(roomID, message)
}

// 定时处理超时事件
func (gs *GameService) ProcessTimeouts() {
    for roomID, game := range gs.games {
        timeoutEvents := game.ProcessTimeouts()
        for _, event := range timeoutEvents {
            gs.handleGameEvent(roomID, event)
        }
    }
}

type PlayerAction struct {
    Type       string      `json:"type"`
    PlayerSeat int         `json:"player_seat"`
    Cards      []*sdk.Card `json:"cards,omitempty"`
    Card       *sdk.Card   `json:"card,omitempty"`
}
```

#### 5. 掼蛋游戏引擎 (完全独立的SDK)

基于现有的`sdk/card.go`和`sdk/comp.go`，扩展实现完全独立的游戏引擎：

```go

// sdk/game_engine.go - 游戏引擎主入口
type GameEngine struct {
    ID           string           `json:"id"`
    CurrentMatch *Match           `json:"current_match"`
    // 内部状态管理
    eventHandlers map[GameEventType][]GameEventHandler
}


// sdk/match.go - 比赛管理
type Match struct {
    ID           string      `json:"id"`
    Status       MatchStatus `json:"status"`
    Players      [4]*Player   `json:"players"`
    CurrentDeal  *Deal       `json:"current_deal"`
    DealHistory  []*Deal     `json:"deal_history"`
    TeamLevels   [2]int      `json:"team_levels"`    // 两队当前等级
    Winner       int         `json:"winner"`         // 获胜队伍 (-1表示未结束)
    StartTime    time.Time   `json:"start_time"`
    EndTime      *time.Time  `json:"end_time,omitempty"`
}

// sdk/deal.go - 局管理
type Deal struct {
    ID              string        `json:"id"`
    Level           int           `json:"level"`           // 本局等级
    Status          DealStatus    `json:"status"`
    CurrentTrick    *Trick        `json:"current_trick"`
    TrickHistory    []*Trick      `json:"trick_history"`
    TributePhase    *TributePhase `json:"tribute_phase,omitempty"`
    PlayerCards     [4][]*Card    `json:"player_cards"`    // 每个玩家当前的手牌
    Rankings        []int         `json:"rankings"`        // 按照出完牌顺序的座次
    StartTime       time.Time     `json:"start_time"`
    EndTime         *time.Time    `json:"end_time,omitempty"`
}

// sdk/trick.go - 轮次管理
type Trick struct {
    ID           string           `json:"id"`
    Leader       int              `json:"leader"`          // 当前Trick在领牌的座位号
    CurrentTurn  int              `json:"current_turn"`    // 当前轮到的玩家
    Plays        []*PlayAction    `json:"plays"`           // 当前Trick下每个玩家的出牌/过牌行为
    Winner       int              `json:"winner"`          // 当前Trick获胜玩家座位号
    LeadComp     CardComp         `json:"lead_comp"`       // 领先牌组，当前在桌面的牌组
    Status       TrickStatus      `json:"status"`
    StartTime    time.Time        `json:"start_time"`
    TurnTimeout  time.Time        `json:"turn_timeout"`    // 当前玩家操作超时时间
}
```

**SDK完全独立的核心方法：**

```go
// 游戏引擎主要接口
type GameEngineInterface interface {
    // 游戏生命周期
    StartMatch(players []Player) error
    StartDeal() error
    
    // 游戏操作
    PlayCards(playerSeat int, cards []*Card) (*GameEvent, error)
    PassTurn(playerSeat int) (*GameEvent, error)
    SelectTribute(playerSeat int, card *Card) (*GameEvent, error)
    
    // 状态查询
    GetGameState() *GameState
    GetPlayerView(playerSeat int) *PlayerGameState  // 玩家视角的状态
    IsGameFinished() bool
    
    // 事件处理
    RegisterEventHandler(eventType GameEventType, handler GameEventHandler)
    ProcessTimeouts() []*GameEvent  // 处理超时事件
    
    // 玩家管理
    HandlePlayerDisconnect(playerSeat int) (*GameEvent, error)
    HandlePlayerReconnect(playerSeat int) (*GameEvent, error)
    SetPlayerAutoPlay(playerSeat int, enabled bool) error
}

// 游戏事件系统
type GameEvent struct {
    Type      GameEventType `json:"type"`
    Data      interface{}   `json:"data"`
    Timestamp time.Time     `json:"timestamp"`
    PlayerSeat int          `json:"player_seat,omitempty"`
}

type GameEventType string
const (
    EventMatchStarted    GameEventType = "match_started"
    EventDealStarted     GameEventType = "deal_started"
    EventCardsDealt      GameEventType = "cards_dealt"
    EventTributePhase    GameEventType = "tribute_phase"
    EventTrickStarted    GameEventType = "trick_started"
    EventPlayerPlayed    GameEventType = "player_played"
    EventPlayerPassed    GameEventType = "player_passed"
    EventTrickEnded      GameEventType = "trick_ended"
    EventDealEnded       GameEventType = "deal_ended"
    EventMatchEnded      GameEventType = "match_ended"
    EventPlayerTimeout   GameEventType = "player_timeout"
    EventPlayerDisconnect GameEventType = "player_disconnect"
    EventPlayerReconnect GameEventType = "player_reconnect"
)
```

### 前端组件

#### 1. 路由结构

```
/login          - 登录页面
/register       - 注册页面
/lobby          - 房间大厅
/room/:id       - 房间页面
/game/:id       - 游戏页面
```

#### 2. 主要React组件

```typescript
// 应用主组件
interface App {
  user: User | null;
  currentRoom: Room | null;
  currentGame: Game | null;
}

// 用户认证组件
interface AuthComponents {
  LoginForm: React.FC;
  RegisterForm: React.FC;
}

// 房间相关组件
interface RoomComponents {
  RoomLobby: React.FC;
  RoomList: React.FC<{ rooms: Room[] }>;
  RoomCard: React.FC<{ room: Room }>;
  RoomWaiting: React.FC<{ room: Room }>;
}

// 游戏相关组件
interface GameComponents {
  GameBoard: React.FC<{ game: Game }>;
  PlayerHand: React.FC<{ cards: Card[], onCardSelect: Function }>;
  PlayArea: React.FC<{ players: Player[] }>;
  GameControls: React.FC<{ onPlay: Function, onPass: Function }>;
  TributePhase: React.FC<{ tribute: TributePhase }>;
  DealResult: React.FC<{ result: DealResult }>;
}
```

#### 3. 状态管理

```typescript
interface AppState {
  auth: {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
  };
  rooms: {
    list: Room[];
    current: Room | null;
    loading: boolean;
  };
  game: {
    current: Game | null;
    selectedCards: Card[];
    gameState: GameState;
  };
  websocket: {
    connected: boolean;
    reconnecting: boolean;
  };
}
```

### WebSocket消息协议

#### 客户端到服务端消息

```typescript
// 加入房间
{
  type: "join_room",
  data: { room_id: string }
}

// 开始游戏
{
  type: "start_game",
  data: { room_id: string }
}

// 出牌
{
  type: "play_cards",
  data: { 
    room_id: string,
    cards: Card[]
  }
}

// 不出
{
  type: "pass",
  data: { room_id: string }
}

// 选择贡牌
{
  type: "select_tribute",
  data: {
    room_id: string,
    card: Card
  }
}
```

#### 服务端到客户端消息

```typescript
// 房间状态更新
{
  type: "room_updated",
  data: Room
}

// 游戏状态更新
{
  type: "game_updated",
  data: Game
}

// 玩家出牌
{
  type: "player_played",
  data: {
    player_id: string,
    cards: Card[]
  }
}

// 轮到玩家操作
{
  type: "player_turn",
  data: {
    player_id: string,
    timeout: number
  }
}

// Deal结束
{
  type: "deal_finished",
  data: DealResult
}
```

## 数据模型

### 架构分层

数据模型按照职责分为两层：
- **SDK层**：包含所有游戏逻辑相关的数据结构，完全独立
- **服务器层**：只包含网络通信和用户管理相关的数据结构

### SDK层数据结构

SDK层包含完整的游戏数据模型，在前面的"掼蛋游戏引擎"部分已经详细定义：

- `GameEngine` - 游戏引擎主结构
- `Match` - 比赛管理
- `Deal` - 局管理  
- `Trick` - 轮次管理
- `Player` - 游戏玩家
- `TributePhase` - 上贡阶段
- `GameEvent` - 游戏事件系统
- `DealResult` / `MatchResult` - 结算结果

所有这些结构都在SDK中定义，确保游戏逻辑的完整性和独立性。

### 服务器层数据结构

服务器层只保留与外部通信和协调相关的最小数据结构：

#### 房间和用户管理

#### 用户认证

### 数据流向

```
前端 ←→ 服务器层 ←→ SDK层
     (WebSocket)   (方法调用/事件)
```

- **前端 → 服务器层**：用户操作（WebSocket消息）
- **服务器层 → SDK层**：调用SDK方法（PlayCards, PassTurn等）
- **SDK层 → 服务器层**：触发事件（GameEvent）
- **服务器层 → 前端**：广播状态更新（WebSocket消息）

## 错误处理

### 错误类型定义

```go
type GameError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

const (
    ErrInvalidCards     = "INVALID_CARDS"
    ErrNotYourTurn      = "NOT_YOUR_TURN"
    ErrRoomFull         = "ROOM_FULL"
    ErrGameInProgress   = "GAME_IN_PROGRESS"
    ErrInsufficientPlayers = "INSUFFICIENT_PLAYERS"
    ErrConnectionLost   = "CONNECTION_LOST"
    ErrTimeout          = "TIMEOUT"
)
```

### 错误处理策略

1. **网络错误**：自动重连机制，最多重试3次
2. **游戏规则错误**：返回具体错误信息，允许重新操作
3. **超时错误**：自动执行默认操作（Pass或托管）
4. **断线处理**：标记为托管状态，游戏继续进行

## 测试策略

### 单元测试

1. **牌规则测试**：利用现有的测试用例
2. **游戏逻辑测试**：
   - 发牌逻辑
   - 上贡还贡逻辑
   - 出牌验证逻辑
   - 结算逻辑

### 集成测试

1. **WebSocket通信测试**：
   - 连接建立和断开
   - 消息发送和接收
   - 房间状态同步

2. **端到端测试**：
   - 完整游戏流程
   - 多用户并发测试
   - 断线重连测试

### 性能测试

1. **并发测试**：支持多个房间同时进行游戏
2. **内存测试**：长时间运行的内存泄漏检测
3. **WebSocket连接测试**：大量并发连接的稳定性

## 部署和扩展

### 容器化部署

```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend
```

### 扩展性考虑

1. **水平扩展**：
   - 使用Redis作为共享状态存储
   - 负载均衡器分发WebSocket连接

2. **数据持久化**：
   - 游戏历史记录存储
   - 用户统计数据

3. **监控和日志**：
   - 游戏状态监控
   - 性能指标收集
   - 错误日志记录

## 安全考虑

### 认证和授权

1. **JWT Token**：用于用户身份验证
2. **WebSocket认证**：连接时验证token有效性
3. **操作权限**：验证用户是否有权执行特定操作

### 防作弊机制

1. **服务端验证**：所有游戏操作在服务端验证
2. **状态同步**：客户端状态仅用于显示，不影响游戏逻辑
3. **超时机制**：防止恶意延迟游戏进程

### 数据安全

1. **输入验证**：所有用户输入进行严格验证
2. **SQL注入防护**：使用参数化查询
3. **XSS防护**：前端输出转义