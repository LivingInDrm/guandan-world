#Tribute，在DealPrepare阶段后，如果不是Match的首个Deal，则会进入该阶段Tribute，它基于上一Deal的结算结果，由败方队伍给胜方队伍贡献1-2张手牌，具体方式如下：

1. 免贡判定（Tribute Immunity Check）

根据上一Deal的获胜情形，判断是否满足免贡条件。若满足免贡条件，则跳过上贡与还贡
三种情形的免贡规则：
	•	Double Down：若 Rank3 和 Rank4 合计持有 两张 Big Joker，则触发免贡。
	•	Single Last：若 Rank4 单独持有 两张 Big Joker，则触发免贡。
	•	Partner Last：若 Rank3 单独持有 两张 Big Joker，则触发免贡。


2. 上贡（Tribute Phase）
若未满足免贡条件，则进入正式上贡流程。

各情形的上贡流程：
	•	Double Down：Rank3 和 Rank4 各上交 1 张贡牌，放入贡牌池；Rank1 优先从贡牌池中挑选其一；Rank2 获得剩下的一张贡牌。
	•	Single Last：Rank4 上交 1 张贡牌，直接交给 Rank1。
	•	Partner Last：Rank3 上交 1 张贡牌，直接交给 Rank1。

上交贡牌时，自动选取玩家前手牌中 除红桃 Trump 外最大的一张牌

3. 还贡（Return Tribute Phase）
若触发了上贡流程，上贡结束后，由胜方玩家从自己手牌中任选一张牌，向贡牌来源归还一张手牌。
