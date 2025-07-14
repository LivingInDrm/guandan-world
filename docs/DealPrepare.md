
1. 洗牌与发牌
	- 将Deck中的牌随机打散
    - 随机标记其中一张为Start Card
    - 将打散后的Deck，按照顺时针，一张张分发至四位玩家手中（最后每人 27 张）
    - 记录Starting Card的牌面和归属玩家

2. 确定本局Level，并标记Trump和Wild Card
	- 读取上一Deal获胜的队伍，用其当前Level，作为当前Deal的Level
	- 若当前Deal为Match的首个Deal，Level默认为双方初始值 2。
    - 根据Level 标记 当前Deal的Trump和Wild Card

3. 主牌（Trump）识别逻辑：当前Deal中所有点数和Level一致的8张牌，都是Trump

4. 万能牌（Wildcard）的识别逻辑：Suit为红桃的Trump，同时为WildCard
