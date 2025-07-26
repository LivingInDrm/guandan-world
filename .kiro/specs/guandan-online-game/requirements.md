# Requirements Document

## Introduction

掼蛋在线对战平台是一个支持四人在线对战的掼蛋游戏系统。该系统提供完整的用户管理、房间管理、游戏流程控制功能，让用户能够在线体验完整的掼蛋游戏，包括登录注册、创建加入房间、四人对战、发牌出牌、结算升级等核心功能。

## Requirements

### Requirement 1 - 用户认证系统

**User Story:** 作为一个掼蛋爱好者，我希望能够注册账号并登录系统，这样我就能够参与在线游戏并保存我的游戏进度。

#### Acceptance Criteria

1. WHEN 用户访问系统 THEN 系统 SHALL 显示登录页面
2. WHEN 用户点击注册 THEN 系统 SHALL 提供用户名和密码注册功能
3. WHEN 用户成功注册 THEN 系统 SHALL 自动登录用户并跳转至房间大厅
4. WHEN 已注册用户输入正确的用户名和密码点击登录 THEN 系统 SHALL 登录用户并跳转至房间大厅
5. WHEN 用户输入用户名密码错误 THEN 系统 SHALL 显示错误信息


### Requirement 2 - 房间大厅管理

**User Story:** 作为一个登录用户，我希望能够查看所有可用房间并选择加入，这样我就能找到合适的游戏房间开始游戏。

#### Acceptance Criteria

1. WHEN 用户进入房间大厅 THEN 系统 SHALL 显示所有在线房间列表
2. WHEN 显示房间列表 THEN 系统 SHALL 优先显示等待中房间，按已加入人数降序排列
3. WHEN 房间列表超过12个 THEN 系统 SHALL 实现分页加载功能
4. WHEN 显示每个房间 THEN 系统 SHALL 显示房间状态、已加入人数、用户名及座次
5. WHEN 房间未满且游戏未开始 THEN 系统 SHALL 显示可点击的加入按钮
6. WHEN 房间已满或游戏进行中 THEN 系统 SHALL 将加入按钮置灰不可点击
7. WHEN 用户点击加入按钮 THEN 系统 SHALL 将用户加入对应房间
8. WHEN 用户点击创建新房间 THEN 系统 SHALL 创建新房间并将用户设为房主


### Requirement 3 - 房间内等待管理

**User Story:** 作为房间内的用户，我希望能够看到其他玩家的加入情况并等待游戏开始，这样我就能知道何时可以开始游戏。

#### Acceptance Criteria

1. WHEN 用户进入房间 THEN 系统 SHALL 自动分配座位号（0-3）
2. WHEN 房间创建者进入 THEN 系统 SHALL 设置其为房主
3. WHEN 房主查看界面 THEN 系统 SHALL 显示开始游戏按钮
4. WHEN 房间人数少于4人 THEN 系统 SHALL 将开始游戏按钮置灰
5. WHEN 房间人数达到4人 THEN 系统 SHALL 激活开始游戏按钮
6. WHEN 非房主用户退出房间 THEN 系统 SHALL 移除该用户并更新房间状态
7. WHEN 房主退出房间 THEN 系统 SHALL 随机选择新房主
8. WHEN 最后一个用户退出 THEN 系统 SHALL 自动关闭房间


### Requirement 4 - 游戏开始流程

**User Story:** 作为房主，我希望能够在人数足够时开始游戏，这样所有玩家就能进入游戏状态。

#### Acceptance Criteria

1. WHEN 房主点击开始游戏 THEN 系统 SHALL 验证房间人数为4人
2. WHEN 人数验证通过 THEN 系统 SHALL 将所有玩家同步进入准备页
3. WHEN 进入准备页 THEN 系统 SHALL 显示3秒倒计时
4. WHEN 倒计时结束 THEN 系统 SHALL 将所有玩家进入牌局页面


### Requirement 5 - 发牌和上贡阶段

**User Story:** 作为游戏玩家，我希望系统能够正确处理发牌、上贡、还贡流程，这样游戏就能按照掼蛋规则正确进行。

#### Acceptance Criteria

1. WHEN 进入新局 THEN 系统 SHALL 为每位玩家发27张牌
2. WHEN 发牌完成 THEN 系统 SHALL 判定上贡和免贡情况
3. WHEN 需要上贡 THEN 系统 SHALL 显示上贡信息和操作界面
4. WHEN 双下情况 THEN 系统 SHALL 提供贡池供胜方选择，Rank 1有3秒选择时间，若超时则自动选择最大的牌
5. WHEN 非双下情况 THEN 系统 SHALL 展示贡牌3秒后自动上贡
6. WHEN 上贡完成 THEN 系统 SHALL 进入还贡阶段，由得到贡的玩家还贡给上贡者，并将还贡结果公示3秒
7. WHEN 上贡还贡完成 THEN 系统 SHALL 确定当前Deal的先出牌者并进入出牌阶段


### Requirement 6 - 出牌游戏界面

**User Story:** 作为游戏玩家，我希望有清晰的游戏界面显示我的手牌和其他玩家状态，这样我就能做出正确的出牌决策。

#### Acceptance Criteria

1. WHEN 进入出牌阶段 THEN 系统 SHALL 在左上角显示两队当前等级和本局等级
2. WHEN 显示玩家手牌 THEN 系统 SHALL 按点数分组并按大小排序显示
3. WHEN 排列手牌 THEN 系统 SHALL 横向按点数从大到小排列，纵向按花色优先级堆叠
4. WHEN 显示出牌区 THEN 系统 SHALL 为每个玩家显示当前状态（已过牌/等待出牌/已出牌/未轮到/已结束）
5. WHEN 轮到玩家出牌 THEN 系统 SHALL 显示20秒倒计时
6. WHEN 出牌超时 THEN 系统 SHALL 自动执行Pass操作
7. WHEN 当前玩家轮次 THEN 系统 SHALL 显示出牌和不出按钮


### Requirement 7 - 出牌操作控制

**User Story:** 作为当前出牌玩家，我希望能够选择手牌进行出牌或选择不出，这样我就能按照游戏策略进行操作。

#### Acceptance Criteria

1. WHEN 轮到用户出牌 THEN 系统 SHALL 允许用户选择手牌
2. WHEN 用户选择牌后点击出牌 THEN 系统 SHALL 验证出牌合法性并执行出牌
3. WHEN 用户点击不出 THEN 系统 SHALL 执行Pass操作
4. WHEN 用户20秒内未操作 THEN 系统 SHALL 自动Pass
5. WHEN 出牌不合法 THEN 系统 SHALL 显示错误提示并要求重新选择
6. WHEN 出牌成功 THEN 系统 SHALL 更新游戏状态并轮到下一位玩家


### Requirement 8 - Deal结算系统

**User Story:** 作为游戏玩家，我希望在每局结束后能看到详细的结算信息，这样我就能了解本局的结果和升级情况。

#### Acceptance Criteria

1. WHEN Deal结束 THEN 系统 SHALL 显示出完牌顺序（Rank 1-4）
2. WHEN 显示排名 THEN 系统 SHALL 按队伍分组显示结果
3. WHEN 计算结果 THEN 系统 SHALL 判断胜方队伍和胜利类型
4. WHEN 显示升级 THEN 系统 SHALL 显示各队升级结果（+3/+2/+1）
5. WHEN 显示统计 THEN 系统 SHALL 显示Deal启动时间和持续时间
6. WHEN 任一队达到A级 THEN 系统 SHALL 显示Match结算和最终比分
7. WHEN 未达到A级 THEN 系统 SHALL 提供继续或退出选项


### Requirement 9 - Match管理系统

**User Story:** 作为游戏玩家，我希望能够进行多局连续游戏直到有队伍获胜，这样就能体验完整的掼蛋比赛。

#### Acceptance Criteria
1. WHEN Deal结束且未达到A级 THEN 系统 SHALL 允许玩家选择继续当前Match
2. WHEN 玩家选择继续 THEN 系统 SHALL 更新等级并开始下一局
3. WHEN 玩家选择退出 THEN 系统 SHALL 结束Match并返回房间大厅
4. WHEN Match结束 THEN 系统 SHALL 显示两队成员、等级、时间统计
5. WHEN 有队伍达到A级 THEN 系统 SHALL 显示最终胜负结果


### Requirement 10 - 用户断线和托管

**User Story:** 作为其他玩家，我希望当有玩家断线时游戏能够继续进行，这样就不会因为个别玩家的网络问题影响整个游戏。

#### Acceptance Criteria

1. WHEN 用户断线 THEN 系统 SHALL 自动标记该用户为托管状态
2. WHEN 托管用户为Trick Leader THEN 系统 SHALL 自动打出当前最小牌
3. WHEN 托管用户非Trick Leader THEN 系统 SHALL 自动执行Pass操作
4. WHEN 用户连续2次超时 THEN 系统 SHALL 将用户设为托管状态
5. WHEN 游戏进行中 THEN 系统 SHALL 禁止新用户加入房间


### Requirement 11 - 操作时间控制

**User Story:** 作为游戏玩家，我希望有合理的操作时间限制，这样游戏就能保持适当的节奏。

#### Acceptance Criteria

1. WHEN 轮到玩家操作 THEN 系统 SHALL 显示20秒倒计时
2. WHEN 操作时间到期 THEN 系统 SHALL 自动执行Pass操作
3. WHEN 玩家连续2次超时 THEN 系统 SHALL 将玩家设为托管状态
4. WHEN 托管状态激活 THEN 系统 SHALL 为该玩家自动执行后续操作