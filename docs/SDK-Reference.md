# 掼蛋世界 SDK 参考文档

## 目录
- [概述](#概述)
- [架构设计](#架构设计)
- [核心模块](#核心模块)
- [API 接口](#api-接口)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)

## 概述

掼蛋世界 SDK 是一个完整的掼蛋游戏引擎，提供了从基础卡牌系统到完整游戏循环的所有功能。SDK 采用模块化设计，支持灵活的扩展和定制。

### 主要特性
- 完整的掼蛋游戏规则实现
- 模块化架构，易于扩展
- 支持多种输入源（AI、人工、网络）
- 事件驱动的游戏状态管理
- 详细的统计数据收集
- 完整的上贡系统（包含双下、单落、对落等情况）

### 模块概览

| 模块文件 | 主要功能 | 核心类型 | 说明 |
|---------|---------|---------|------|
| **types.go** | 核心类型定义 | `Player`, `Match`, `Deal`, `Trick` | 定义所有核心数据结构和枚举类型 |
| **card.go** | 卡牌系统 | `Card` | 实现掼蛋卡牌的完整体系，包含变化牌逻辑 |
| **dealer.go** | 发牌器 | `Dealer` | 负责洗牌、发牌等操作 |
| **comp.go** | 牌型组合 | `CardComp`, `CompType` | 识别和比较所有掼蛋牌型 |
| **validator.go** | 规则验证 | `PlayValidator` | 验证出牌和游戏规则的合法性 |
| **trick.go** | 轮次管理 | `Trick`, `PlayAction` | 管理单个出牌轮次的生命周期 |
| **tribute.go** | 上贡系统 | `TributeManager`, `TributePhase` | 实现完整的上贡规则和免贡逻辑 |
| **deal.go** | 牌局管理 | `Deal` | 管理单个牌局的完整流程 |
| **match.go** | 比赛管理 | `Match` | 管理多局比赛直到升A |
| **result.go** | 结果处理 | `DealResult`, `DealResultCalculator` | 计算游戏结果和统计数据 |
| **game_engine.go** | 游戏引擎 | `GameEngine`, `GameEvent` | 核心控制器，协调所有组件 |
| **game_driver.go** | 游戏驱动器 | `GameDriver`, `PlayerInputProvider` | 高级封装，提供完整游戏循环 |

## 架构设计

### 核心架构图
```
┌─────────────────────────────────────────────────────────────┐
│                      Game Driver                            │
│  ┌───────────────┐  ┌─────────────────┐  ┌──────────────┐  │
│  │ Input Provider│  │   Game Engine   │  │   Observers  │  │
│  └───────────────┘  └─────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                      Core Modules                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │  Card   │  │ Tribute │  │ Trick   │  │     Match       │ │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │ Dealer  │  │ Comp    │  │  Deal   │  │   Validator     │ │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 设计原则
- **单一职责**：每个模块专注于特定功能
- **松耦合**：模块间通过接口交互
- **事件驱动**：通过事件系统实现状态通知
- **可扩展性**：支持插件化的输入提供者和观察者

## 核心模块

### 1. 类型系统 (types.go)

定义了SDK中所有核心数据结构和枚举类型。

#### 主要类型

##### Player 玩家信息
```go
type Player struct {
    ID       string `json:"id"`        // 玩家唯一标识
    Username string `json:"username"`  // 玩家用户名
    Seat     int    `json:"seat"`      // 座位号(0-3)
    Online   bool   `json:"online"`    // 是否在线
    AutoPlay bool   `json:"auto_play"` // 是否自动出牌
}
```

##### Match 比赛实例
```go
type Match struct {
    ID          string      `json:"id"`           // 比赛唯一标识
    Status      MatchStatus `json:"status"`       // 比赛状态
    Players     [4]*Player  `json:"players"`      // 4个玩家
    CurrentDeal *Deal       `json:"current_deal"` // 当前牌局
    DealHistory []*Deal     `json:"deal_history"` // 历史牌局
    TeamLevels  [2]int      `json:"team_levels"`  // 两队等级
    Winner      int         `json:"winner"`       // 获胜队伍
    StartTime   time.Time   `json:"start_time"`   // 开始时间
    EndTime     *time.Time  `json:"end_time"`     // 结束时间
}
```

##### Deal 牌局实例
```go
type Deal struct {
    ID           string        `json:"id"`            // 牌局唯一标识
    Level        int           `json:"level"`         // 当前级别
    Status       DealStatus    `json:"status"`        // 牌局状态
    CurrentTrick *Trick        `json:"current_trick"` // 当前轮次
    TrickHistory []*Trick      `json:"trick_history"` // 轮次历史
    TributePhase *TributePhase `json:"tribute_phase"` // 上贡阶段
    PlayerCards  [4][]*Card    `json:"player_cards"`  // 玩家手牌
    Rankings     []int         `json:"rankings"`      // 玩家排名
}
```

##### Trick 出牌轮次
```go
type Trick struct {
    ID          string        `json:"id"`           // 轮次唯一标识
    Leader      int           `json:"leader"`       // 轮次领导者
    CurrentTurn int           `json:"current_turn"` // 当前出牌者
    Plays       []*PlayAction `json:"plays"`        // 出牌记录
    Winner      int           `json:"winner"`       // 轮次获胜者
    LeadComp    CardComp      `json:"lead_comp"`    // 领先牌型
    Status      TrickStatus   `json:"status"`       // 轮次状态
}
```

#### 状态枚举

##### MatchStatus 比赛状态
- `MatchStatusWaiting` - 等待开始
- `MatchStatusPlaying` - 进行中
- `MatchStatusFinished` - 已结束

##### DealStatus 牌局状态
- `DealStatusWaiting` - 等待开始
- `DealStatusDealing` - 发牌中
- `DealStatusTribute` - 上贡阶段
- `DealStatusPlaying` - 游戏中
- `DealStatusFinished` - 已结束

##### TrickStatus 轮次状态
- `TrickStatusWaiting` - 等待开始
- `TrickStatusPlaying` - 进行中
- `TrickStatusFinished` - 已结束

### 2. 卡牌系统 (card.go)

实现掼蛋游戏的完整卡牌体系。

#### Card 卡牌结构
```go
type Card struct {
    Number    int    // 牌面数值 (2-16)
    RawNumber int    // 原始数值 (用于顺子判断)
    Color     string // 花色 (Spade/Heart/Diamond/Club/Joker)
    Level     int    // 当前级别 (用于变化牌判断)
    Name      string // 牌面名称
}
```

#### 核心方法

##### NewCard 创建卡牌
```go
func NewCard(number int, color string, level int) (*Card, error)
```
- 创建新的卡牌实例
- 自动处理 Ace 转换(1→14)
- 验证数值和花色的合法性

##### IsWildcard 判断变化牌
```go
func (c *Card) IsWildcard() bool
```
- 判断是否为红桃变化牌
- 变化牌：红桃且数值等于当前级别

##### GreaterThan 比较大小
```go
func (c *Card) GreaterThan(other *Card) bool
```
- 按掼蛋规则比较牌的大小
- 考虑级别牌和王牌的特殊规则

##### ToShortString 简化显示
```go
func (c *Card) ToShortString() string
```
- 返回简化的字符串表示
- 格式：`9H`(红桃9), `QS`(黑桃Q), `BJ`(大王)

### 3. 发牌器 (dealer.go)

负责洗牌和发牌操作。

#### Dealer 发牌器
```go
type Dealer struct {
    level int     // 当前级别
    deck  []*Card // 牌堆
}
```

#### 核心方法

##### NewDealer 创建发牌器
```go
func NewDealer(level int) (*Dealer, error)
```

##### CreateFullDeck 创建完整牌堆
```go
func (d *Dealer) CreateFullDeck() []*Card
```
- 创建108张牌的完整牌堆
- 每种花色2-A各2张，共104张
- 大小王各2张，共4张

##### ShuffleDeck 洗牌
```go
func (d *Dealer) ShuffleDeck()
```
- 使用Fisher-Yates算法洗牌
- 保证随机性

##### DealCards 发牌
```go
func (d *Dealer) DealCards() ([4][]*Card, error)
```
- 给4个玩家各发27张牌
- 自动对每个玩家的手牌排序

### 4. 牌型组合系统 (comp.go)

实现掼蛋游戏的所有牌型识别和比较。

#### CardComp 牌型接口
```go
type CardComp interface {
    GreaterThan(other CardComp) bool  // 比较大小
    IsBomb() bool                     // 是否为炸弹
    GetCards() []*Card                // 获取卡牌
    String() string                   // 字符串表示
    IsValid() bool                    // 是否合法
    GetType() CompType                // 获取牌型
}
```

#### CompType 牌型类型
```go
type CompType int

const (
    TypeFold         CompType = iota // 弃牌
    TypeIllegal                      // 非法牌型
    TypeSingle                       // 单张
    TypePair                         // 对子
    TypeTriple                       // 三张
    TypeFullHouse                    // 葫芦
    TypeStraight                     // 顺子
    TypePlate                        // 钢板(连对)
    TypeTube                         // 钢管(连三)
    TypeJokerBomb                    // 王炸
    TypeNaiveBomb                    // 同数炸弹
    TypeStraightFlush                // 同花顺
)
```

#### 核心功能
- 自动识别所有掼蛋牌型
- 支持变化牌的特殊规则
- 实现完整的牌型比较逻辑
- 处理各种炸弹类型的优先级

### 5. 验证器 (validator.go)

提供游戏规则验证功能。

#### PlayValidator 出牌验证器
```go
type PlayValidator struct {
    level int // 当前级别
}
```

#### 核心方法

##### ValidatePlay 验证出牌
```go
func (pv *PlayValidator) ValidatePlay(playerSeat int, cards []*Card, 
    playerCards []*Card, currentTrick *Trick) error
```
- 验证出牌合法性
- 检查卡牌归属
- 验证牌型有效性
- 验证是否能压过当前牌型

##### ValidatePass 验证过牌
```go
func (pv *PlayValidator) ValidatePass(playerSeat int, currentTrick *Trick) error
```
- 验证过牌的合法性
- 首家不能过牌

### 6. 轮次管理 (trick.go)

管理单个出牌轮次的完整生命周期。

#### Trick 轮次实例
核心属性已在类型系统中描述。

#### 核心方法

##### NewTrick 创建轮次
```go
func NewTrick(leader int) (*Trick, error)
```

##### PlayCards 处理出牌
```go
func (t *Trick) PlayCards(playerSeat int, cards []*Card, comp CardComp) error
```
- 处理玩家出牌
- 更新轮次状态
- 判断轮次是否结束

##### PassTurn 处理过牌
```go
func (t *Trick) PassTurn(playerSeat int) error
```
- 处理玩家过牌
- 移动到下一位玩家

##### ProcessTimeout 处理超时
```go
func (t *Trick) ProcessTimeout() error
```
- 自动过牌处理超时

### 7. 上贡系统 (tribute.go)

实现完整的掼蛋上贡规则。

#### TributeManager 上贡管理器
```go
type TributeManager struct {
    level int // 当前级别
}
```

#### TributePhase 上贡阶段
```go
type TributePhase struct {
    Status           TributeStatus `json:"status"`            // 上贡状态
    TributeMap       map[int]int   `json:"tribute_map"`       // 上贡映射
    TributeCards     map[int]*Card `json:"tribute_cards"`     // 上贡卡牌
    ReturnCards      map[int]*Card `json:"return_cards"`      // 还贡卡牌
    PoolCards        []*Card       `json:"pool_cards"`        // 贡牌池(双下)
    SelectingPlayer  int           `json:"selecting_player"`  // 选牌玩家
    SelectTimeout    time.Time     `json:"select_timeout"`    // 选择超时
    IsImmune         bool          `json:"is_immune"`         // 是否免贡
    SelectionResults map[int]int   `json:"selection_results"` // 选择结果
}
```

#### 上贡规则

##### 胜利类型与上贡规则
1. **双下 (Double Down)**：rank1、rank2同队
   - rank3、rank4各上交1张最大非红桃主牌到贡牌池
   - rank1优先从池中选择1张，rank2获得剩余1张
   - 双方分别还贡1张

2. **单落 (Single Last)**：rank1、rank3同队
   - rank4直接上贡1张最大非红桃主牌给rank1
   - rank1还贡1张

3. **对落 (Partner Last)**：rank1、rank4同队
   - rank3直接上贡1张最大非红桃主牌给rank1
   - rank1还贡1张

##### 免贡条件
败方队伍合计持有2张及以上大王时触发免贡。

#### 核心方法

##### NewTributePhase 创建上贡阶段
```go
func NewTributePhase(lastResult *DealResult) (*TributePhase, error)
```

##### CheckTributeImmunity 检查免贡
```go
func (tm *TributeManager) CheckTributeImmunity(lastResult *DealResult, 
    playerHands [4][]*Card) bool
```

##### ProcessTribute 处理上贡
```go
func (tm *TributeManager) ProcessTribute(tributePhase *TributePhase, 
    playerHands [4][]*Card) error
```

##### ApplyTributeToHands 应用上贡效果
```go
func (tm *TributeManager) ApplyTributeToHands(tributePhase *TributePhase, 
    playerHands *[4][]*Card) error
```

### 8. 牌局管理 (deal.go)

管理单个牌局的完整生命周期。

#### Deal 牌局实例
核心属性已在类型系统中描述。

#### 核心方法

##### NewDeal 创建牌局
```go
func NewDeal(level int, lastResult *DealResult) (*Deal, error)
```
- 创建新牌局
- 根据上局结果初始化上贡阶段

##### StartDeal 开始牌局
```go
func (d *Deal) StartDeal() error
```
- 发牌给所有玩家
- 处理上贡阶段
- 开始首轮出牌

##### PlayCards 处理出牌
```go
func (d *Deal) PlayCards(playerSeat int, cards []*Card) error
```
- 验证并处理玩家出牌
- 更新游戏状态
- 检查牌局是否结束

##### PassTurn 处理过牌
```go
func (d *Deal) PassTurn(playerSeat int) error
```

### 9. 比赛管理 (match.go)

管理完整比赛的生命周期。

#### Match 比赛实例
核心属性已在类型系统中描述。

#### 核心方法

##### NewMatch 创建比赛
```go
func NewMatch(players []Player) (*Match, error)
```
- 验证4个玩家
- 初始化队伍等级

##### StartNewDeal 开始新牌局
```go
func (m *Match) StartNewDeal() error
```
- 基于当前等级创建新牌局
- 自动处理上贡逻辑

##### FinishDeal 结束牌局
```go
func (m *Match) FinishDeal(result *DealResult) error
```
- 更新队伍等级
- 检查比赛是否结束

##### GetTeamForPlayer 获取玩家队伍
```go
func (m *Match) GetTeamForPlayer(playerSeat int) int
```
- 返回玩家所属队伍(0或1)
- 队伍分配：0,2为一队，1,3为一队

### 10. 结果处理 (result.go)

计算和统计游戏结果。

#### DealResult 牌局结果
```go
type DealResult struct {
    Rankings    []int         `json:"rankings"`     // 玩家排名
    WinningTeam int           `json:"winning_team"` // 获胜队伍
    VictoryType VictoryType   `json:"victory_type"` // 胜利类型
    Upgrades    [2]int        `json:"upgrades"`     // 升级数
    Duration    time.Duration `json:"duration"`     // 用时
    TrickCount  int           `json:"trick_count"`  // 轮次数
    Statistics  *DealStatistics `json:"statistics"` // 详细统计
}
```

#### VictoryType 胜利类型
```go
type VictoryType string

const (
    VictoryTypePartnerLast VictoryType = "Partner Last" // 对落 (+1级)
    VictoryTypeSingleLast  VictoryType = "Single Last"  // 单落 (+2级)
    VictoryTypeDoubleDown  VictoryType = "Double Down"  // 双下 (+3级)
)
```

#### DealResultCalculator 结果计算器
```go
type DealResultCalculator struct {
    level int
}
```

##### CalculateDealResult 计算牌局结果
```go
func (drc *DealResultCalculator) CalculateDealResult(deal *Deal, 
    match *Match) (*DealResult, error)
```
- 计算胜利类型和升级数
- 生成详细统计数据

### 11. 游戏引擎 (game_engine.go)

游戏的核心控制器，协调所有组件。

#### GameEngine 游戏引擎
```go
type GameEngine struct {
    id            string                               // 引擎ID
    status        GameStatus                           // 游戏状态
    currentMatch  *Match                               // 当前比赛
    eventHandlers map[GameEventType][]GameEventHandler // 事件处理器
    mutex         sync.RWMutex                         // 并发控制
    createdAt     time.Time                            // 创建时间
    updatedAt     time.Time                            // 更新时间
}
```

#### GameEngineInterface 引擎接口
```go
type GameEngineInterface interface {
    // 生命周期管理
    StartMatch(players []Player) error
    StartDeal() error
    
    // 游戏操作
    PlayCards(playerSeat int, cards []*Card) (*GameEvent, error)
    PassTurn(playerSeat int) (*GameEvent, error)
    
    // 上贡操作
    ProcessTributePhase() (*TributeAction, error)
    SubmitTributeSelection(playerID int, cardID string) error
    SubmitReturnTribute(playerID int, cardID string) error
    
    // 状态查询
    GetGameState() (*GameState, error)
    GetPlayerGameState(playerSeat int) (*PlayerGameState, error)
    GetTurnInfo() (*TurnInfo, error)
    
    // 事件处理
    RegisterEventHandler(eventType GameEventType, handler GameEventHandler)
    UnregisterEventHandler(eventType GameEventType, handler GameEventHandler)
}
```

#### 游戏事件

##### GameEvent 游戏事件
```go
type GameEvent struct {
    Type       GameEventType `json:"type"`        // 事件类型
    Data       interface{}   `json:"data"`        // 事件数据
    Timestamp  time.Time     `json:"timestamp"`   // 时间戳
    PlayerSeat int           `json:"player_seat"` // 相关玩家
}
```

##### 事件类型
- `EventMatchStarted` - 比赛开始
- `EventDealStarted` - 牌局开始
- `EventCardsDealt` - 发牌完成
- `EventTributePhase` - 上贡阶段
- `EventTrickStarted` - 轮次开始
- `EventPlayerPlayed` - 玩家出牌
- `EventPlayerPassed` - 玩家过牌
- `EventTrickEnded` - 轮次结束
- `EventDealEnded` - 牌局结束
- `EventMatchEnded` - 比赛结束

### 12. 游戏驱动器 (game_driver.go)

高级封装，提供完整的游戏循环管理。

#### GameDriver 游戏驱动器
```go
type GameDriver struct {
    engine        GameEngineInterface // 游戏引擎
    inputProvider PlayerInputProvider // 输入提供者
    observers     []EventObserver     // 事件观察者
    config        *GameDriverConfig   // 驱动器配置
}
```

#### PlayerInputProvider 输入提供者接口
```go
type PlayerInputProvider interface {
    RequestPlayDecision(ctx context.Context, playerSeat int, 
        hand []*Card, trickInfo *TrickInfo) (*PlayDecision, error)
    RequestTributeSelection(ctx context.Context, playerSeat int, 
        options []*Card) (*Card, error)
    RequestReturnTribute(ctx context.Context, playerSeat int, 
        hand []*Card) (*Card, error)
}
```

#### EventObserver 事件观察者接口
```go
type EventObserver interface {
    OnGameEvent(event *GameEvent)
}
```

## API 接口

### 核心接口方法

#### 游戏生命周期

```go
// 创建游戏引擎
engine := NewGameEngine(gameID)

// 开始比赛
err := engine.StartMatch(players)

// 开始新牌局
err := engine.StartDeal()
```

#### 出牌操作

```go
// 玩家出牌
event, err := engine.PlayCards(playerSeat, cards)

// 玩家过牌
event, err := engine.PassTurn(playerSeat)
```

#### 上贡操作

```go
// 处理上贡阶段
action, err := engine.ProcessTributePhase()

// 提交选牌(双下)
err := engine.SubmitTributeSelection(playerID, cardID)

// 提交还贡
err := engine.SubmitReturnTribute(playerID, cardID)
```

#### 状态查询

```go
// 获取游戏状态
state, err := engine.GetGameState()

// 获取玩家状态
playerState, err := engine.GetPlayerGameState(playerSeat)

// 获取轮次信息
turnInfo, err := engine.GetTurnInfo()
```

#### 事件处理

```go
// 注册事件处理器
engine.RegisterEventHandler(EventPlayerPlayed, func(event *GameEvent) {
    // 处理玩家出牌事件
})
```

## 使用示例

### 基础游戏循环

```go
package main

import (
    "fmt"
    "github.com/guandan-world/sdk"
)

func main() {
    // 创建游戏引擎
    engine := sdk.NewGameEngine("game1")
    
    // 创建玩家
    players := []sdk.Player{
        {ID: "p1", Username: "Alice", Seat: 0},
        {ID: "p2", Username: "Bob", Seat: 1},
        {ID: "p3", Username: "Charlie", Seat: 2},
        {ID: "p4", Username: "David", Seat: 3},
    }
    
    // 注册事件处理器
    engine.RegisterEventHandler(sdk.EventPlayerPlayed, func(event *sdk.GameEvent) {
        fmt.Printf("Player %d played cards\n", event.PlayerSeat)
    })
    
    // 开始比赛
    err := engine.StartMatch(players)
    if err != nil {
        panic(err)
    }
    
    // 开始首局
    err = engine.StartDeal()
    if err != nil {
        panic(err)
    }
    
    // 游戏循环
    for {
        // 获取当前轮次信息
        turnInfo, err := engine.GetTurnInfo()
        if err != nil {
            break
        }
        
        // 获取当前玩家状态
        playerState, err := engine.GetPlayerGameState(turnInfo.CurrentPlayer)
        if err != nil {
            break
        }
        
        // 模拟玩家出牌(这里需要实际的决策逻辑)
        if len(playerState.PlayerCards) > 0 {
            // 出单张
            cards := []*sdk.Card{playerState.PlayerCards[0]}
            _, err := engine.PlayCards(turnInfo.CurrentPlayer, cards)
            if err != nil {
                // 尝试过牌
                engine.PassTurn(turnInfo.CurrentPlayer)
            }
        }
    }
}
```

### 自定义输入提供者

```go
type MyInputProvider struct{}

func (p *MyInputProvider) RequestPlayDecision(ctx context.Context, 
    playerSeat int, hand []*Card, trickInfo *TrickInfo) (*PlayDecision, error) {
    
    // 实现自定义决策逻辑
    if trickInfo.IsLeader {
        // 首家出单张
        return &PlayDecision{
            Action: ActionPlay,
            Cards:  []*Card{hand[0]},
        }, nil
    } else {
        // 跟牌或过牌
        return &PlayDecision{
            Action: ActionPass,
        }, nil
    }
}

func (p *MyInputProvider) RequestTributeSelection(ctx context.Context, 
    playerSeat int, options []*Card) (*Card, error) {
    // 选择第一张牌
    return options[0], nil
}

func (p *MyInputProvider) RequestReturnTribute(ctx context.Context, 
    playerSeat int, hand []*Card) (*Card, error) {
    // 返回最小的牌
    return hand[len(hand)-1], nil
}
```

### 使用游戏驱动器

```go
// 创建游戏驱动器
engine := sdk.NewGameEngine("game1")
driver := sdk.NewGameDriver(engine, nil)

// 设置输入提供者
driver.SetInputProvider(&MyInputProvider{})

// 添加观察者
driver.AddObserver(&MyObserver{})

// 运行完整比赛
result, err := driver.RunMatch(players)
if err != nil {
    panic(err)
}

fmt.Printf("Match completed, winner: Team %d\n", result.Winner)
```

## 最佳实践

### 1. 错误处理
- 始终检查方法返回的错误
- 对无效状态转换进行适当处理
- 使用日志记录错误详情

### 2. 并发安全
- SDK内部使用读写锁保护共享状态
- 外部调用者应避免并发修改同一游戏实例
- 事件处理器应设计为无状态或线程安全

### 3. 内存管理
- 及时清理已结束的游戏实例
- 避免在事件处理器中持有大量引用
- 合理使用事件处理器注册/注销

### 4. 性能优化
- 批量处理多个操作时使用事务
- 对频繁查询的状态进行缓存
- 避免在热路径中进行复杂计算

### 5. 扩展性
- 通过接口而非具体类型进行交互
- 使用组合而非继承扩展功能
- 保持模块间的低耦合

### 6. 测试
- 为每个模块编写单元测试
- 使用模拟对象测试复杂交互
- 编写集成测试验证完整流程

---

*本文档描述了掼蛋世界 SDK v1.0 的完整功能。如有疑问或需要技术支持，请参考代码注释或联系开发团队。* 