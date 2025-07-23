package sdk

import (
	"fmt"
	"strings"
)

// SimulatedPlayer 模拟玩家
type SimulatedPlayer struct {
	Player
	AutoPlayAlgorithm AutoPlayAlgorithm
}

// AutoPlayAlgorithm 自动出牌算法接口
type AutoPlayAlgorithm interface {
	// SelectCardsToPlay 选择要出的牌
	// 参数:
	//   hand: 手牌
	//   currentTrick: 当前轮次信息
	//   isLeader: 是否为首出
	// 返回值:
	//   []*Card: 要出的牌，如果返回nil表示过牌
	SelectCardsToPlay(hand []*Card, currentTrick *Trick, isLeader bool) []*Card

	// SelectTributeCard 选择要进贡的牌
	// 参数:
	//   hand: 手牌
	//   excludeHeartTrump: 是否排除红桃主牌
	// 返回值:
	//   *Card: 选择的贡牌
	SelectTributeCard(hand []*Card, excludeHeartTrump bool) *Card

	// SelectReturnCard 选择要还贡的牌
	// 参数:
	//   hand: 手牌
	//   avoidBreakingBomb: 是否避免破坏炸弹
	// 返回值:
	//   *Card: 选择的还贡牌
	SelectReturnCard(hand []*Card, avoidBreakingBomb bool) *Card
}

// SimpleAutoPlayAlgorithm 简单的自动出牌算法实现
type SimpleAutoPlayAlgorithm struct {
	level int // 当前级别，用于判断主牌
}

// NewSimpleAutoPlayAlgorithm 创建简单自动出牌算法
func NewSimpleAutoPlayAlgorithm(level int) *SimpleAutoPlayAlgorithm {
	return &SimpleAutoPlayAlgorithm{
		level: level,
	}
}

// SelectCardsToPlay 实现自动出牌逻辑
func (s *SimpleAutoPlayAlgorithm) SelectCardsToPlay(hand []*Card, currentTrick *Trick, isLeader bool) []*Card {
	if isLeader {
		// 首出：出张数尽可能多的合法非炸弹牌
		return s.selectLeaderPlay(hand)
	} else {
		// 跟牌：如果能压过则出牌，否则过牌
		return s.selectFollowPlay(hand, currentTrick)
	}
}

// selectLeaderPlay 首出选牌逻辑
func (s *SimpleAutoPlayAlgorithm) selectLeaderPlay(hand []*Card) []*Card {
	// 如果没有手牌，返回nil
	if len(hand) == 0 {
		return nil
	}

	// 尝试各种可能的出牌组合，从多到少

	// 尝试出三张
	if triples := s.findTriples(hand); len(triples) > 0 {
		return triples[0]
	}

	// 尝试出对子
	if pairs := s.findPairs(hand); len(pairs) > 0 {
		return pairs[0]
	}

	// 默认出单张（选最小的）- 这个总是有效的
	smallest := hand[0]
	for _, card := range hand {
		if card.LessThan(smallest) {
			smallest = card
		}
	}
	return []*Card{smallest}
}

// selectFollowPlay 跟牌选牌逻辑
func (s *SimpleAutoPlayAlgorithm) selectFollowPlay(hand []*Card, currentTrick *Trick) []*Card {
	if currentTrick == nil || currentTrick.LeadComp == nil {
		return nil
	}

	leadType := currentTrick.LeadComp.GetType()
	leadCards := currentTrick.LeadComp.GetCards()

	// 根据领出的牌型寻找能压过的牌
	switch leadType {
	case TypeSingle:
		return s.findBeatingSingle(hand, leadCards[0])
	case TypePair:
		return s.findBeatingPair(hand, currentTrick.LeadComp)
	case TypeTriple:
		return s.findBeatingTriple(hand, currentTrick.LeadComp)
	case TypeStraight:
		return s.findBeatingStraight(hand, currentTrick.LeadComp)
	case TypeFullHouse:
		return s.findBeatingFullHouse(hand, currentTrick.LeadComp)
	case TypePlate:
		return s.findBeatingPlate(hand, currentTrick.LeadComp)
	case TypeTube:
		return s.findBeatingTube(hand, currentTrick.LeadComp)
	default:
		// 对于炸弹类型，暂时不跟
		return nil
	}
}

// 以下是各种牌型的查找方法

func (s *SimpleAutoPlayAlgorithm) findPairs(hand []*Card) [][]*Card {
	pairs := make([][]*Card, 0)
	cardCount := make(map[int]int)
	cardsByNumber := make(map[int][]*Card)

	// 统计每个数字的牌数
	for _, card := range hand {
		cardCount[card.Number]++
		cardsByNumber[card.Number] = append(cardsByNumber[card.Number], card)
	}

	// 找出所有对子
	for number, count := range cardCount {
		if count >= 2 {
			cards := cardsByNumber[number]
			pairs = append(pairs, cards[:2])
		}
	}

	// 按牌面值从小到大排序
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i][0].GreaterThan(pairs[j][0]) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	return pairs
}

func (s *SimpleAutoPlayAlgorithm) findTriples(hand []*Card) [][]*Card {
	triples := make([][]*Card, 0)
	cardCount := make(map[int]int)
	cardsByNumber := make(map[int][]*Card)

	// 统计每个数字的牌数
	for _, card := range hand {
		if card.Color != "Joker" { // 三张不能包含王
			cardCount[card.Number]++
			cardsByNumber[card.Number] = append(cardsByNumber[card.Number], card)
		}
	}

	// 找出所有三张
	for number, count := range cardCount {
		if count >= 3 {
			cards := cardsByNumber[number]
			triples = append(triples, cards[:3])
		}
	}

	// 按牌面值从小到大排序
	for i := 0; i < len(triples); i++ {
		for j := i + 1; j < len(triples); j++ {
			if triples[i][0].GreaterThan(triples[j][0]) {
				triples[i], triples[j] = triples[j], triples[i]
			}
		}
	}

	return triples
}

func (s *SimpleAutoPlayAlgorithm) findStraights(hand []*Card) [][]*Card {
	// 简化实现：暂时返回空
	return make([][]*Card, 0)
}

func (s *SimpleAutoPlayAlgorithm) findFullHouses(hand []*Card) [][]*Card {
	// 简化实现：暂时返回空
	return make([][]*Card, 0)
}

func (s *SimpleAutoPlayAlgorithm) findPlates(hand []*Card) [][]*Card {
	// 简化实现：暂时返回空
	return make([][]*Card, 0)
}

func (s *SimpleAutoPlayAlgorithm) findTubes(hand []*Card) [][]*Card {
	// 简化实现：暂时返回空
	return make([][]*Card, 0)
}

// 查找能打过的牌

func (s *SimpleAutoPlayAlgorithm) findBeatingSingle(hand []*Card, leadCard *Card) []*Card {
	for _, card := range hand {
		if card.GreaterThan(leadCard) {
			return []*Card{card}
		}
	}

	// 检查是否有炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}

	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingPair(hand []*Card, leadComp CardComp) []*Card {
	pairs := s.findPairs(hand)
	for _, pair := range pairs {
		comp := NewPair(pair)
		if comp.IsValid() && comp.GreaterThan(leadComp) {
			return pair
		}
	}

	// 检查是否有炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}

	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingTriple(hand []*Card, leadComp CardComp) []*Card {
	triples := s.findTriples(hand)
	for _, triple := range triples {
		comp := NewTriple(triple)
		if comp.IsValid() && comp.GreaterThan(leadComp) {
			return triple
		}
	}

	// 检查是否有炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}

	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingStraight(hand []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时只考虑炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}
	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingFullHouse(hand []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时只考虑炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}
	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingPlate(hand []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时只考虑炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}
	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBeatingTube(hand []*Card, leadComp CardComp) []*Card {
	// 简化实现：暂时只考虑炸弹
	if bombs := s.findBombs(hand); len(bombs) > 0 {
		return bombs[0]
	}
	return nil
}

func (s *SimpleAutoPlayAlgorithm) findBombs(hand []*Card) [][]*Card {
	bombs := make([][]*Card, 0)

	// 检查王炸
	jokerCount := 0
	jokers := make([]*Card, 0)
	for _, card := range hand {
		if card.Color == "Joker" {
			jokerCount++
			jokers = append(jokers, card)
		}
	}
	if jokerCount == 4 {
		bombs = append(bombs, jokers)
	}

	// 检查普通炸弹（4张或以上相同数字）
	cardCount := make(map[int]int)
	cardsByNumber := make(map[int][]*Card)

	for _, card := range hand {
		if card.Color != "Joker" {
			cardCount[card.Number]++
			cardsByNumber[card.Number] = append(cardsByNumber[card.Number], card)
		}
	}

	for number, count := range cardCount {
		if count >= 4 {
			bombs = append(bombs, cardsByNumber[number])
		}
	}

	return bombs
}

// SelectTributeCard 选择贡牌
func (s *SimpleAutoPlayAlgorithm) SelectTributeCard(hand []*Card, excludeHeartTrump bool) *Card {
	if len(hand) == 0 {
		return nil
	}

	// 选择最大的牌（排除红桃主牌）
	var maxCard *Card
	for _, card := range hand {
		if excludeHeartTrump && card.IsWildcard() {
			continue
		}
		if maxCard == nil || card.GreaterThan(maxCard) {
			maxCard = card
		}
	}

	// 如果没有找到合适的牌（全是红桃主牌），则选第一张
	if maxCard == nil && len(hand) > 0 {
		maxCard = hand[0]
	}

	return maxCard
}

// SelectReturnCard 选择还贡的牌
func (s *SimpleAutoPlayAlgorithm) SelectReturnCard(hand []*Card, avoidBreakingBomb bool) *Card {
	if len(hand) == 0 {
		return nil
	}

	// 如果需要避免破坏炸弹，先统计牌数
	cardCount := make(map[int]int)
	if avoidBreakingBomb {
		for _, card := range hand {
			if card.Color != "Joker" {
				cardCount[card.Number]++
			}
		}
	}

	// 选择最小的牌（避免破坏炸弹）
	var minCard *Card
	for _, card := range hand {
		// 如果需要避免破坏炸弹，检查这张牌是否会破坏炸弹
		if avoidBreakingBomb && card.Color != "Joker" {
			if cardCount[card.Number] >= 4 {
				continue // 跳过会破坏炸弹的牌
			}
		}

		if minCard == nil || card.LessThan(minCard) {
			minCard = card
		}
	}

	// 如果没有找到合适的牌，则选最小的
	if minCard == nil && len(hand) > 0 {
		minCard = hand[0]
		for _, card := range hand {
			if card.LessThan(minCard) {
				minCard = card
			}
		}
	}

	return minCard
}

// MatchSimulator 比赛模拟器
type MatchSimulator struct {
	engine         *GameEngine
	players        []SimulatedPlayer
	eventLog       []string
	verbose        bool
	dealCount      int
	trickCount     int
	currentDealNum int
}

// NewMatchSimulator 创建新的比赛模拟器
func NewMatchSimulator(verbose bool) *MatchSimulator {
	return &MatchSimulator{
		engine:   NewGameEngine(),
		players:  make([]SimulatedPlayer, 4),
		eventLog: make([]string, 0),
		verbose:  verbose,
	}
}

// SimulateMatch 模拟完整的比赛
func (ms *MatchSimulator) SimulateMatch() error {
	// 创建4个模拟玩家
	for i := 0; i < 4; i++ {
		ms.players[i] = SimulatedPlayer{
			Player: Player{
				ID:       fmt.Sprintf("player_%d", i),
				Username: fmt.Sprintf("Player %d", i+1),
				Seat:     i,
				Online:   true,
				AutoPlay: true,
			},
			AutoPlayAlgorithm: NewSimpleAutoPlayAlgorithm(2), // 假设从2开始打
		}
	}

	// 注册事件处理器
	ms.registerEventHandlers()

	// 将Player类型转换为[]Player
	players := make([]Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = ms.players[i].Player
	}

	// 开始比赛
	if err := ms.engine.StartMatch(players); err != nil {
		return fmt.Errorf("failed to start match: %w", err)
	}

	ms.log("Match started with 4 players")

	// 主游戏循环（添加安全计数器防止无限循环）
	maxDeals := 10 // 进一步降低最大局数限制，防止测试超时
	for !ms.engine.IsGameFinished() && ms.currentDealNum < maxDeals {
		ms.currentDealNum++

		// 记录当前队伍情况和等级
		ms.logTeamStatus()

		// 开始新的一局
		if err := ms.engine.StartDeal(); err != nil {
			return fmt.Errorf("failed to start deal: %w", err)
		}

		ms.log(fmt.Sprintf("\n=== Deal %d started ===", ms.currentDealNum))

		// 处理这一局直到结束
		if err := ms.processDeal(); err != nil {
			return fmt.Errorf("failed to process deal: %w", err)
		}

		ms.dealCount++
	}

	if ms.currentDealNum >= maxDeals {
		ms.log(fmt.Sprintf("\n=== Match terminated after %d deals (safety limit) ===", ms.dealCount))
	} else {
		ms.log(fmt.Sprintf("\n=== Match finished after %d deals ===", ms.dealCount))
	}

	// 打印最终结果
	ms.printMatchSummary()

	return nil
}

// processDeal 处理一局游戏
func (ms *MatchSimulator) processDeal() error {
	gameState := ms.engine.GetGameState()
	if gameState.CurrentMatch == nil || gameState.CurrentMatch.CurrentDeal == nil {
		return fmt.Errorf("no active deal")
	}

	deal := gameState.CurrentMatch.CurrentDeal

	// 处理贡牌阶段
	if deal.Status == DealStatusTribute {
		if err := ms.processTributePhase(); err != nil {
			return fmt.Errorf("failed to process tribute: %w", err)
		}

		// 立即重新获取状态，因为processTributePhase可能已经更新了deal状态
		gameState = ms.engine.GetGameState()
		if gameState.CurrentMatch == nil || gameState.CurrentMatch.CurrentDeal == nil {
			return nil // Deal已结束
		}
		deal = gameState.CurrentMatch.CurrentDeal

		// 如果贡牌阶段已完成，立即打印贡牌详情
		if deal.TributePhase != nil && deal.TributePhase.Status == TributeStatusFinished && !deal.TributePhase.IsImmune {
			ms.logTributeDetails(deal.TributePhase)
			ms.logPlayerHands("After Tribute", deal)
		}
	} else if deal.TributePhase != nil && deal.TributePhase.IsImmune {
		// 处理免贡的情况 - 贡牌阶段被跳过
		ms.log("Tribute phase skipped due to immunity")
	}

	// 等待deal状态变为playing（贡牌阶段结束后）- 简化逻辑
	// 重新获取最新状态
	gameState = ms.engine.GetGameState()
	if gameState.CurrentMatch == nil || gameState.CurrentMatch.CurrentDeal == nil {
		return nil // Deal已结束
	}
	deal = gameState.CurrentMatch.CurrentDeal

	if deal.Status != DealStatusPlaying {
		ms.log(fmt.Sprintf("Deal status is %v instead of playing, skipping this deal", deal.Status))
		return nil // 跳过这个deal而不是报错
	}

	// 游戏主循环（添加安全计数器）
	maxTricks := 200 // 增加每局最大轮数限制，确保玩家能出完所有牌
	trickCounter := 0
	for deal.Status == DealStatusPlaying && trickCounter < maxTricks {
		trickCounter++

		// 检查是否有活跃的trick
		if deal.CurrentTrick == nil {
			ms.log("Warning: No current trick in playing state")
			break
		}

		// 获取当前轮到谁出牌
		currentPlayer := deal.CurrentTrick.CurrentTurn

		// 获取玩家视图
		playerView := ms.engine.GetPlayerView(currentPlayer)
		playerHand := playerView.PlayerCards

		// 判断是否为首出
		isLeader := deal.CurrentTrick.LeadComp == nil

		// 如果是trick的第一次出牌（首出），输出所有玩家手牌
		if isLeader && len(deal.CurrentTrick.Plays) == 0 {
			ms.logPlayerHands(fmt.Sprintf("New Trick Started (Leader: Player %d)", currentPlayer), deal)
		}

		// 使用自动算法选择出牌
		algorithm := ms.players[currentPlayer].AutoPlayAlgorithm
		selectedCards := algorithm.SelectCardsToPlay(playerHand, deal.CurrentTrick, isLeader)

		// 执行出牌或过牌
		if selectedCards != nil && len(selectedCards) > 0 {
			_, err := ms.engine.PlayCards(currentPlayer, selectedCards)
			if err != nil {
				// 如果出牌失败，尝试过牌（非首出时）或强制出最小单张（首出时）
				ms.log(fmt.Sprintf("Player %d failed to play cards: %v", currentPlayer, err))
				if !isLeader {
					_, err = ms.engine.PassTurn(currentPlayer)
					if err != nil {
						return fmt.Errorf("player %d failed to pass: %w", currentPlayer, err)
					}
				} else {
					// 首出失败，强制出最小的单张
					smallest := playerHand[0]
					for _, card := range playerHand {
						if card.LessThan(smallest) {
							smallest = card
						}
					}
					_, err = ms.engine.PlayCards(currentPlayer, []*Card{smallest})
					if err != nil {
						return fmt.Errorf("player %d failed to play emergency single card: %w", currentPlayer, err)
					}
				}
			}
		} else if !isLeader {
			// 过牌
			_, err := ms.engine.PassTurn(currentPlayer)
			if err != nil {
				return fmt.Errorf("player %d failed to pass: %w", currentPlayer, err)
			}
		} else {
			// 首出时没有选中牌，强制出最小单张
			if len(playerHand) > 0 {
				smallest := playerHand[0]
				for _, card := range playerHand {
					if card.LessThan(smallest) {
						smallest = card
					}
				}
				_, err := ms.engine.PlayCards(currentPlayer, []*Card{smallest})
				if err != nil {
					return fmt.Errorf("player %d failed to play forced single card: %w", currentPlayer, err)
				}
			}
		}

		// 检查游戏状态更新
		gameState = ms.engine.GetGameState()
		if gameState.CurrentMatch.CurrentDeal == nil {
			break // Deal已结束
		}
		deal = gameState.CurrentMatch.CurrentDeal

		ms.trickCount++
	}

	if trickCounter >= maxTricks {
		ms.log(fmt.Sprintf("WARNING: Deal terminated after %d tricks (safety limit). Deal status: %v", trickCounter, deal.Status))
	}

	return nil
}

// processTributePhase 处理贡牌阶段
func (ms *MatchSimulator) processTributePhase() error {
	// 添加安全计数器防止无限循环
	maxTributeActions := 10 // 增加最大循环次数，确保完成整个流程
	actionsProcessed := 0

	for actionsProcessed < maxTributeActions {
		actionsProcessed++

		// 调用新的贡牌接口
		action, err := ms.engine.ProcessTributePhase()
		if err != nil {
			ms.log(fmt.Sprintf("ProcessTributePhase error: %v", err))
			// 如果贡牌阶段出错，直接进入游戏阶段
			return nil
		}

		// 如果没有待处理的动作，检查贡牌阶段是否真正完成
		if action == nil {
			// 获取当前状态以确认贡牌阶段是否完成
			gameState := ms.engine.GetGameState()
			if gameState.CurrentMatch != nil && gameState.CurrentMatch.CurrentDeal != nil {
				deal := gameState.CurrentMatch.CurrentDeal
				if deal.Status == DealStatusPlaying {
					ms.log("Tribute phase completed and game phase started")
					break
				} else if deal.TributePhase != nil && deal.TributePhase.Status == TributeStatusFinished {
					// 贡牌阶段标记为完成但 Deal 状态还没更新，再调用一次 ProcessTributePhase
					ms.log("Tribute phase finished, triggering state transition")
					continue
				}
			}
			ms.log("No tribute action available")
			break
		}

		// 根据动作类型处理
		switch action.Type {
		case TributeActionSelect:
			// 双下选牌：选择最大的牌
			if len(action.Options) > 0 {
				selectedCard := action.Options[0]
				for _, card := range action.Options {
					if card.GreaterThan(selectedCard) {
						selectedCard = card
					}
				}

				ms.log(fmt.Sprintf("Player %d selecting tribute card: %s", action.PlayerID, selectedCard))
				if err := ms.engine.SubmitTributeSelection(action.PlayerID, selectedCard.GetID()); err != nil {
					return fmt.Errorf("failed to submit tribute selection: %w", err)
				}
			}

		case TributeActionReturn:
			// 还贡：选择最小的牌（避免破坏炸弹）
			algorithm := ms.players[action.PlayerID].AutoPlayAlgorithm
			returnCard := algorithm.SelectReturnCard(action.Options, true)

			if returnCard != nil {
				ms.log(fmt.Sprintf("Player %d returning card: %s", action.PlayerID, returnCard))
				if err := ms.engine.SubmitReturnTribute(action.PlayerID, returnCard.GetID()); err != nil {
					return fmt.Errorf("failed to submit return tribute: %w", err)
				}
			}
		}
	}

	if actionsProcessed >= maxTributeActions {
		ms.log("Tribute phase terminated due to safety limit")
	}

	return nil
}

// registerEventHandlers 注册事件处理器
func (ms *MatchSimulator) registerEventHandlers() {
	// 注册各种事件的处理器
	ms.engine.RegisterEventHandler(EventMatchStarted, ms.handleMatchStarted)
	ms.engine.RegisterEventHandler(EventDealStarted, ms.handleDealStarted)
	ms.engine.RegisterEventHandler(EventCardsDealt, ms.handleCardsDealt)
	ms.engine.RegisterEventHandler(EventTributePhase, ms.handleTributePhase)
	ms.engine.RegisterEventHandler(EventTributeImmunity, ms.handleTributeImmunity)
	ms.engine.RegisterEventHandler(EventTributeStarted, ms.handleTributeStarted)
	ms.engine.RegisterEventHandler(EventTributeGiven, ms.handleTributeGiven)
	ms.engine.RegisterEventHandler(EventTributeSelected, ms.handleTributeSelected)
	ms.engine.RegisterEventHandler(EventReturnTribute, ms.handleReturnTribute)
	ms.engine.RegisterEventHandler(EventTributeCompleted, ms.handleTributeCompleted)
	ms.engine.RegisterEventHandler(EventTrickStarted, ms.handleTrickStarted)
	ms.engine.RegisterEventHandler(EventPlayerPlayed, ms.handlePlayerPlayed)
	ms.engine.RegisterEventHandler(EventPlayerPassed, ms.handlePlayerPassed)
	ms.engine.RegisterEventHandler(EventTrickEnded, ms.handleTrickEnded)
	ms.engine.RegisterEventHandler(EventDealEnded, ms.handleDealEnded)
	ms.engine.RegisterEventHandler(EventMatchEnded, ms.handleMatchEnded)
}

// 事件处理方法

func (ms *MatchSimulator) handleMatchStarted(event *GameEvent) {
	ms.log("Event: Match Started")
}

func (ms *MatchSimulator) handleDealStarted(event *GameEvent) {
	ms.log("Event: Deal Started")
	// 输出发牌后每个玩家的初始手牌
	if deal, ok := event.Data.(*Deal); ok {
		ms.logPlayerHands("Deal Started", deal)

	}
}

func (ms *MatchSimulator) handleCardsDealt(event *GameEvent) {
	ms.log("Event: Cards Dealt")
}

func (ms *MatchSimulator) handleTributePhase(event *GameEvent) {
	ms.log("Event: Tribute Phase")
}

func (ms *MatchSimulator) handleTributeImmunity(event *GameEvent) {
	ms.log("Event: Tribute Immunity triggered - No tribute required this deal")
}

func (ms *MatchSimulator) handleTributeStarted(event *GameEvent) {
	ms.log("Event: Tribute Started - Tribute phase begins")
}

func (ms *MatchSimulator) handleTributeGiven(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if giver, ok := data["giver"].(int); ok {
			if receiver, ok := data["receiver"].(int); ok {
				if card, ok := data["card"].(*Card); ok {
					ms.log(fmt.Sprintf("Event: Tribute Given - Player %d gives %s to Player %d",
						giver, card.ToShortString(), receiver))
				}
			}
		}
	}
}

func (ms *MatchSimulator) handleTributeSelected(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if player, ok := data["player"].(int); ok {
			if cardID, ok := data["cardID"].(string); ok {
				ms.log(fmt.Sprintf("Event: Tribute Selected - Player %d selected card %s (Double-down selection)",
					player, cardID))
			}
		}
	}
}

func (ms *MatchSimulator) handleReturnTribute(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if returner, ok := data["player"].(int); ok {
			if cardID, ok := data["cardID"].(string); ok {
				ms.log(fmt.Sprintf("Event: Return Tribute - Player %d returns card %s",
					returner, cardID))
			}
		}
	}
}

func (ms *MatchSimulator) handleTributeCompleted(event *GameEvent) {
	ms.log("Event: Tribute Completed")
	// 贡牌详情已经在processDeal中打印，这里不再重复打印
}

func (ms *MatchSimulator) handleTrickStarted(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if leader, ok := data["leader"].(int); ok {
			ms.log(fmt.Sprintf("Event: New Trick Started, Leader: Player %d", leader))
			// TODO: 需要以异步方式或在主循环中添加手牌输出，避免死锁
		}
	}
}

func (ms *MatchSimulator) handlePlayerPlayed(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		cards := data["cards"].([]*Card)

		// 将出牌转换为简化格式
		var cardStrs []string
		for _, card := range cards {
			cardStrs = append(cardStrs, card.ToShortString())
		}

		ms.log(fmt.Sprintf("Event: Player %d played %d cards: [%s]",
			playerSeat, len(cards), strings.Join(cardStrs, ",")))
	}
}

func (ms *MatchSimulator) handlePlayerPassed(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		playerSeat := data["player_seat"].(int)
		ms.log(fmt.Sprintf("Event: Player %d passed", playerSeat))
	}
}

func (ms *MatchSimulator) handleTrickEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if winner, ok := data["winner"].(int); ok {
			ms.log(fmt.Sprintf("Event: Trick Ended, Winner: Player %d", winner))
			// TODO: 在不持有锁的情况下输出手牌信息
			// 暂时禁用手牌输出以避免死锁问题
		}
	}
}

func (ms *MatchSimulator) handleDealEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if result, ok := data["result"].(*DealResult); ok {
			ms.log(fmt.Sprintf("Event: Deal Ended, Rankings: %v, Victory Type: %v",
				result.Rankings, result.VictoryType))
		}
	}
}

func (ms *MatchSimulator) handleMatchEnded(event *GameEvent) {
	if data, ok := event.Data.(map[string]interface{}); ok {
		if winner, ok := data["winner"].(int); ok {
			ms.log(fmt.Sprintf("Event: Match Ended, Winner: Team %d", winner))
		}
	}
}

// log 记录日志
func (ms *MatchSimulator) log(message string) {
	ms.eventLog = append(ms.eventLog, message)
	if ms.verbose {
		fmt.Println(message)
	}
}

// formatPlayerHands 格式化所有玩家的手牌为简化输出
func (ms *MatchSimulator) formatPlayerHands(deal *Deal) string {
	if deal == nil {
		return "No active deal"
	}

	var result []string

	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		playerCards := deal.PlayerCards[playerSeat]
		var cardStrs []string

		// 将每张牌转换为简化格式
		for _, card := range playerCards {
			cardStrs = append(cardStrs, card.ToShortString())
		}

		result = append(result, fmt.Sprintf("Player %d (%d cards): [%s]",
			playerSeat, len(playerCards), strings.Join(cardStrs, ",")))
	}

	return strings.Join(result, "\n")
}

// logTeamStatus 输出当前队伍情况和等级
func (ms *MatchSimulator) logTeamStatus() {
	gameState := ms.engine.GetGameState()
	if gameState.CurrentMatch == nil {
		return
	}

	match := gameState.CurrentMatch
	ms.log(fmt.Sprintf("=== Team Status Before Deal %d ===", ms.currentDealNum))
	ms.log(fmt.Sprintf("Team 0 (Players 0,2): Level %d", match.TeamLevels[0]))
	ms.log(fmt.Sprintf("Team 1 (Players 1,3): Level %d", match.TeamLevels[1]))

	// 显示玩家名称
	ms.log("Players:")
	for i := 0; i < 4; i++ {
		teamNum := i % 2
		ms.log(fmt.Sprintf("  Player %d (%s) - Team %d", i, match.Players[i].Username, teamNum))
	}
}

// logTributeDetails 输出贡牌阶段的详细信息
func (ms *MatchSimulator) logTributeDetails(tributePhase *TributePhase) {
	if tributePhase == nil {
		return
	}

	ms.log("=== Tribute Details ===")

	// 输出贡牌映射关系
	if len(tributePhase.TributeMap) > 0 {
		ms.log("Tribute Map (Giver -> Receiver):")
		for giver, receiver := range tributePhase.TributeMap {
			ms.log(fmt.Sprintf("  Player %d -> Player %d", giver, receiver))
		}
	}

	// 输出具体的贡牌
	if len(tributePhase.TributeCards) > 0 {
		for giver, card := range tributePhase.TributeCards {
			receiver := tributePhase.TributeMap[giver]
			ms.log(fmt.Sprintf("Tribute Cards: Player %d gave %s to Player %d", giver, card.ToShortString(), receiver))
		}
	}

	// 输出还贡牌
	if len(tributePhase.ReturnCards) > 0 {
		for returner, card := range tributePhase.ReturnCards {
			ms.log(fmt.Sprintf("Return Cards: Player %d returned %s", returner, card.ToShortString()))
		}
	}

	// 如果有抗贡（免贡）
	if tributePhase.Status == TributeStatusFinished && len(tributePhase.TributeCards) == 0 && len(tributePhase.TributeMap) == 0 {
		ms.log("Tribute was skipped (Immunity)")
	}
}

// logPlayerHands 输出所有玩家的手牌
func (ms *MatchSimulator) logPlayerHands(context string, deal *Deal) {
	if ms.verbose {
		handInfo := ms.formatPlayerHands(deal)
		ms.log(fmt.Sprintf("%s - Player Hands:", context))
		ms.log(handInfo)
	}
}

// printMatchSummary 打印比赛总结
func (ms *MatchSimulator) printMatchSummary() {
	gameState := ms.engine.GetGameState()
	match := gameState.CurrentMatch

	fmt.Println("\n========== Match Summary ==========")
	fmt.Printf("Total Deals: %d\n", ms.dealCount)
	fmt.Printf("Total Tricks: %d\n", ms.trickCount)

	if match != nil {
		fmt.Printf("Winner: Team %d\n", match.Winner)
		fmt.Printf("Final Levels: Team 0: Level %d, Team 1: Level %d\n",
			match.TeamLevels[0],
			match.TeamLevels[1])

		if match.EndTime != nil {
			duration := match.EndTime.Sub(match.StartTime)
			fmt.Printf("Duration: %v\n", duration)
		}
	}

	fmt.Println("===================================")
}

// RunMatchSimulation 运行比赛模拟的便捷函数
func RunMatchSimulation(verbose bool) error {
	simulator := NewMatchSimulator(verbose)
	return simulator.SimulateMatch()
}

// RunVerboseDemo 运行详细模式演示（用于调试和学习）
func RunVerboseDemo() error {
	fmt.Println("🎮 掼蛋比赛模拟器 - 详细模式演示")
	fmt.Println("=====================================")
	fmt.Println("🚀 开始模拟比赛（详细模式）...")

	simulator := NewMatchSimulator(true) // 启用详细模式
	err := simulator.SimulateMatch()

	if err != nil {
		fmt.Printf("❌ 模拟失败: %v\n", err)
		return err
	}

	fmt.Println("\n✅ 模拟完成!")
	return nil
}
