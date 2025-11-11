package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"guandan-world/ai"
	"guandan-world/backend/handlers"
	wsmanager "guandan-world/backend/websocket"
	"guandan-world/sdk"

	"github.com/gorilla/websocket"
)

// AIPlayerClient 代表一个独立的AI玩家客户端
type AIPlayerClient struct {
	// 配置
	serverURL   string
	roomID      string
	username    string
	password    string
	verbose     bool
	playDelay   time.Duration

	// 认证信息
	authToken string
	userID    string

	// HTTP客户端
	httpClient *HTTPClient

	// WebSocket连接
	wsConn *websocket.Conn

	// AI算法
	aiAlgorithm ai.AutoPlayAlgorithm

	// 游戏状态
	playerSeat  int
	hand        []*sdk.Card
	currentRank int
	gameActive  bool
	mu          sync.RWMutex

	// 控制
	stopChan chan struct{}
	doneChan chan struct{}
}

// NewAIPlayerClient 创建新的AI玩家客户端
func NewAIPlayerClient(serverURL, roomID, username, password string, verbose bool) *AIPlayerClient {
	return NewAIPlayerClientWithDelay(serverURL, roomID, username, password, verbose, 5*time.Second)
}

// NewAIPlayerClientWithDelay 创建新的AI玩家客户端并指定出牌延迟
func NewAIPlayerClientWithDelay(serverURL, roomID, username, password string, verbose bool, playDelay time.Duration) *AIPlayerClient {
	return &AIPlayerClient{
		serverURL:   serverURL,
		roomID:      roomID,
		username:    username,
		password:    password,
		verbose:     verbose,
		playDelay:   playDelay,
		httpClient:  NewHTTPClient(serverURL, &http.Client{Timeout: 10 * time.Second}),
		aiAlgorithm: ai.NewSmartAutoPlayAlgorithm(2), // 从2级开始
		currentRank: 2,
		playerSeat:  -1, // 未分配座位
		hand:        make([]*sdk.Card, 0),
		stopChan:    make(chan struct{}),
		doneChan:    make(chan struct{}),
	}
}

// Start 启动AI玩家客户端
func (c *AIPlayerClient) Start() error {
	c.log("Starting AI player client...")

	// 1. 注册或登录
	if err := c.registerOrLogin(); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	c.log("Authentication successful")

	// 2. 加入房间
	if err := c.joinRoom(); err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}
	c.log(fmt.Sprintf("Joined room %s, seat: %d", c.roomID, c.playerSeat))

	// 3. 连接WebSocket
	if err := c.connectWebSocket(); err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}
	c.log("WebSocket connected")

	// 4. 发送join_room消息
	if err := c.sendJoinRoomMessage(); err != nil {
		return fmt.Errorf("failed to send join room message: %w", err)
	}
	c.log("Sent join_room message")

	// 5. 启动消息处理循环
	go c.handleMessages()

	c.setGameActive(true)
	c.log("AI player client started, waiting for game...")

	return nil
}

// Stop 停止客户端
func (c *AIPlayerClient) Stop() {
	close(c.stopChan)
	if c.wsConn != nil {
		c.wsConn.Close()
	}
	<-c.doneChan
}

// Wait 等待游戏结束
func (c *AIPlayerClient) Wait() {
	<-c.doneChan
}

// registerOrLogin 注册或登录账号
func (c *AIPlayerClient) registerOrLogin() error {
	// 先尝试注册
	registerReq := handlers.RegisterRequest{
		Username: c.username,
		Password: c.password,
	}

	var authResp handlers.AuthResponse
	err := c.httpClient.CallAPI("POST", "/api/auth/register", registerReq, &authResp)

	if err != nil {
		// 如果注册失败（可能是用户已存在），尝试登录
		c.log("Registration failed, trying login...")

		loginReq := handlers.LoginRequest{
			Username: c.username,
			Password: c.password,
		}

		err = c.httpClient.CallAPI("POST", "/api/auth/login", loginReq, &authResp)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
	}

	c.authToken = authResp.Token.Token
	c.userID = authResp.User.ID
	c.httpClient.SetAuthToken(c.authToken)

	return nil
}

// joinRoom 加入房间
func (c *AIPlayerClient) joinRoom() error {
	url := fmt.Sprintf("/api/rooms/%s/join", c.roomID)
	var roomResp handlers.RoomResponse

	if err := c.httpClient.CallAPI("POST", url, nil, &roomResp); err != nil {
		return err
	}

	// 找到自己的座位号
	for _, player := range roomResp.Room.Players {
		if player.ID == c.userID {
			c.playerSeat = player.Seat
			break
		}
	}

	if c.playerSeat == -1 {
		return fmt.Errorf("failed to get seat assignment")
	}

	return nil
}

// connectWebSocket 连接WebSocket
func (c *AIPlayerClient) connectWebSocket() error {
	wsURL := fmt.Sprintf("ws://%s/ws?token=%s", c.serverURL, c.authToken)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	c.wsConn = conn
	return nil
}

// sendJoinRoomMessage 发送加入房间的WebSocket消息
func (c *AIPlayerClient) sendJoinRoomMessage() error {
	msg := wsmanager.WSMessage{
		Type: wsmanager.MSG_JOIN_ROOM,
		Data: map[string]interface{}{
			"room_id": c.roomID,
		},
		Timestamp: time.Now(),
	}

	return c.wsConn.WriteJSON(msg)
}

// handleMessages 处理WebSocket消息
func (c *AIPlayerClient) handleMessages() {
	defer func() {
		c.wsConn.Close()
		c.setGameActive(false)
		close(c.doneChan)
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
		}

		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			c.log(fmt.Sprintf("WebSocket read error: %v", err))
			return
		}

		var wsMsg wsmanager.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			c.log(fmt.Sprintf("Failed to parse WebSocket message: %v", err))
			continue
		}

		// 处理不同类型的消息
		switch wsMsg.Type {
		case wsmanager.MSG_GAME_EVENT:
			c.handleGameEvent(&wsMsg)
		case wsmanager.MSG_GAME_ACTION:
			c.handleGameAction(&wsMsg)
		case wsmanager.MSG_ERROR:
			c.log(fmt.Sprintf("Error message: %v", wsMsg.Data))
		}
	}
}

// handleGameEvent 处理游戏事件
func (c *AIPlayerClient) handleGameEvent(msg *wsmanager.WSMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	eventType, ok := data["event_type"].(string)
	if !ok {
		return
	}

	c.logVerbose(fmt.Sprintf("Game event: %s", eventType))

	// 处理游戏结束事件
	if eventType == "match_ended" {
		c.log("Match ended")
		c.setGameActive(false)
		// 不立即关闭，等待外部调用Stop
	}

	// 处理发牌事件 - 更新手牌
	if eventType == "cards_dealt" {
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			c.updateHandFromEvent(eventData)
		}
	}

	// 处理出牌事件 - 从手牌中移除
	if eventType == "player_played" {
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			playerSeat := -1
			if ps, ok := eventData["player_seat"].(float64); ok {
				playerSeat = int(ps)
			}
			
			// 记录所有玩家的出牌
			if cardsData, ok := eventData["cards"].([]interface{}); ok {
				c.logVerbose(fmt.Sprintf("Player %d played %d cards", playerSeat, len(cardsData)))
			}
			
			// 只移除自己的手牌
			if playerSeat == c.playerSeat {
				c.removeCardsFromHand(eventData)
			}
		}
	}
}

// handleGameAction 处理游戏操作请求
func (c *AIPlayerClient) handleGameAction(msg *wsmanager.WSMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	actionType, ok := data["action_type"].(string)
	if !ok {
		return
	}

	playerSeat, ok := data["player_seat"].(float64)
	if !ok || int(playerSeat) != c.playerSeat {
		// 不是针对我的操作请求
		return
	}

	c.log(fmt.Sprintf("Action request: %s", actionType))

	switch actionType {
	case "play_decision_required":
		go c.handlePlayDecisionRequest(data)
	case "tribute_selection_required":
		go c.handleTributeSelectionRequest(data)
	case "return_tribute_required":
		go c.handleReturnTributeRequest(data)
	}
}

// updateHandFromEvent 从事件数据中更新手牌
func (c *AIPlayerClient) updateHandFromEvent(eventData map[string]interface{}) {
	// 查找对应座位的手牌
	if hands, ok := eventData["hands"].(map[string]interface{}); ok {
		seatKey := fmt.Sprintf("%d", c.playerSeat)
		if handData, ok := hands[seatKey].([]interface{}); ok {
			c.mu.Lock()
			c.hand = ParseCards(handData, c.currentRank)
			c.mu.Unlock()
			c.logVerbose(fmt.Sprintf("Updated hand: %d cards", len(c.hand)))
		}
	}
}

// removeCardsFromHand 从手牌中移除已出的牌
func (c *AIPlayerClient) removeCardsFromHand(eventData map[string]interface{}) {
	if cards, ok := eventData["cards"].([]interface{}); ok {
		playedCards := ParseCards(cards, c.currentRank)
		cardIDs := make(map[string]bool)
		for _, card := range playedCards {
			cardIDs[card.GetID()] = true
		}

		c.mu.Lock()
		newHand := make([]*sdk.Card, 0, len(c.hand))
		for _, card := range c.hand {
			if !cardIDs[card.GetID()] {
				newHand = append(newHand, card)
			}
		}
		c.hand = newHand
		c.mu.Unlock()

		c.logVerbose(fmt.Sprintf("Removed cards, remaining: %d", len(c.hand)))
	}
}

// Helper methods

func (c *AIPlayerClient) log(message string) {
	log.Printf("[%s] %s", c.username, message)
}

func (c *AIPlayerClient) logVerbose(message string) {
	if c.verbose {
		log.Printf("[%s] %s", c.username, message)
	}
}

func (c *AIPlayerClient) setGameActive(active bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gameActive = active
}

func (c *AIPlayerClient) isGameActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.gameActive
}

// handlePlayDecisionRequest 处理出牌决策请求
func (c *AIPlayerClient) handlePlayDecisionRequest(data map[string]interface{}) {
	// 解析手牌
	handData, ok := data["hand"].([]interface{})
	if !ok {
		c.log("Failed to parse hand data")
		return
	}

	hand := ParseCards(handData, c.currentRank)

	// 更新手牌状态
	c.mu.Lock()
	c.hand = hand
	c.mu.Unlock()

	// 解析trick信息
	trickData, ok := data["trick_info"].(map[string]interface{})
	if !ok {
		c.log("Failed to parse trick info")
		return
	}

	isLeader, ok := trickData["is_leader"].(bool)
	if !ok {
		c.log("Failed to parse is_leader")
		return
	}

	trickInfo := &sdk.TrickInfo{
		IsLeader: isLeader,
	}

	// 如果有leadComp，解析它
	if leadCompData, ok := trickData["lead_comp"]; ok && leadCompData != nil {
		c.logVerbose(fmt.Sprintf("Raw lead_comp data: %+v", leadCompData))
		if leadCompMap, ok := leadCompData.(map[string]interface{}); ok {
			c.logVerbose(fmt.Sprintf("lead_comp map: %+v", leadCompMap))
			if cardsData, ok := leadCompMap["cards"].([]interface{}); ok {
				c.logVerbose(fmt.Sprintf("lead_comp cards count: %d", len(cardsData)))
				leadCards := ParseCards(cardsData, c.currentRank)
				c.logVerbose(fmt.Sprintf("Parsed lead cards count: %d", len(leadCards)))
				if len(leadCards) > 0 {
					trickInfo.LeadComp = sdk.FromCardList(leadCards, nil)
					c.logVerbose(fmt.Sprintf("LeadComp set: type=%v, valid=%v", trickInfo.LeadComp.GetType(), trickInfo.LeadComp.IsValid()))
				}
			} else {
				c.logVerbose("Failed to parse cards from lead_comp")
			}
		} else {
			c.logVerbose(fmt.Sprintf("lead_comp is not a map, type: %T", leadCompData))
		}
	} else {
		c.logVerbose("No lead_comp in trick_info")
	}

	// 使用AI算法获取决策
	selectedCards := c.aiAlgorithm.SelectCardsToPlay(hand, trickInfo)

	c.logVerbose(fmt.Sprintf("AI selected %d cards, isLeader=%v, hasLeadComp=%v", 
		len(selectedCards), isLeader, trickInfo.LeadComp != nil))
	if len(selectedCards) > 0 {
		cardIDs := make([]string, len(selectedCards))
		for i, card := range selectedCards {
			cardIDs[i] = card.GetID()
		}
		c.logVerbose(fmt.Sprintf("Selected cards: %v", cardIDs))
	}
	if trickInfo.LeadComp != nil {
		c.logVerbose(fmt.Sprintf("LeadComp: type=%v, cards=%v", 
			trickInfo.LeadComp.GetType(), trickInfo.LeadComp.GetCards()))
	}

	// 安全检查：如果是首出但算法返回空结果，强制出最小的牌
	if isLeader && len(selectedCards) == 0 {
		c.log("WARNING: AI algorithm returned no cards for trick leader, forcing smallest card")
		if len(hand) > 0 {
			// 找到最小的牌
			smallest := hand[0]
			for _, card := range hand[1:] {
				if !card.GreaterThan(smallest) {
					smallest = card
				}
			}
			selectedCards = []*sdk.Card{smallest}
		} else {
			c.log("ERROR: No cards in hand!")
			return
		}
	}

	// 等待配置的延迟时间
	time.Sleep(c.playDelay)

	// 构建决策
	var action string
	var cardIDs []string

	if len(selectedCards) == 0 {
		action = "pass"
		c.logVerbose("Decision: PASS")
	} else {
		action = "play"
		cardIDs = make([]string, len(selectedCards))
		for i, card := range selectedCards {
			cardIDs[i] = card.GetID()
		}
		c.logVerbose(fmt.Sprintf("Decision: PLAY %d cards", len(cardIDs)))
	}

	// 提交决策
	req := handlers.PlayDecisionRequest{
		RoomID:     c.roomID,
		PlayerSeat: c.playerSeat,
		Action:     action,
		CardIDs:    cardIDs,
	}

	if err := c.httpClient.CallAPI("POST", "/api/game/driver/play-decision", req, nil); err != nil {
		c.log(fmt.Sprintf("Failed to submit play decision: %v", err))
	} else {
		c.logVerbose("Play decision submitted successfully")
	}
}

// handleTributeSelectionRequest 处理贡牌选择请求
func (c *AIPlayerClient) handleTributeSelectionRequest(data map[string]interface{}) {
	// 解析选项
	optionsData, ok := data["options"].([]interface{})
	if !ok {
		c.log("Failed to parse tribute options")
		return
	}

	options := ParseCards(optionsData, c.currentRank)

	// 选择最大的牌
	selectedCard := SelectMaxCard(options)

	if selectedCard == nil {
		c.log("No valid tribute card to select")
		return
	}

	c.logVerbose(fmt.Sprintf("Tribute selection: %s", selectedCard.GetID()))

	// 提交选择
	req := handlers.TributeSelectionRequest{
		RoomID:     c.roomID,
		PlayerSeat: c.playerSeat,
		CardID:     selectedCard.GetID(),
	}

	if err := c.httpClient.CallAPI("POST", "/api/game/driver/tribute-select", req, nil); err != nil {
		c.log(fmt.Sprintf("Failed to submit tribute selection: %v", err))
	} else {
		c.logVerbose("Tribute selection submitted successfully")
	}
}

// handleReturnTributeRequest 处理还贡请求
func (c *AIPlayerClient) handleReturnTributeRequest(data map[string]interface{}) {
	// 解析手牌
	handData, ok := data["hand"].([]interface{})
	if !ok {
		c.log("Failed to parse hand data for return tribute")
		return
	}

	hand := ParseCards(handData, c.currentRank)

	// 更新手牌状态
	c.mu.Lock()
	c.hand = hand
	c.mu.Unlock()

	// 使用AI算法选择还贡牌
	returnCard := c.aiAlgorithm.SelectReturnTributeCard(hand, nil)

	if returnCard == nil {
		// 如果AI没有选择，选最小的牌
		returnCard = SelectMinCard(hand)
	}

	if returnCard == nil {
		c.log("No valid return tribute card")
		return
	}

	c.logVerbose(fmt.Sprintf("Return tribute: %s", returnCard.GetID()))

	// 提交还贡
	req := handlers.ReturnTributeRequest{
		RoomID:     c.roomID,
		PlayerSeat: c.playerSeat,
		CardID:     returnCard.GetID(),
	}

	if err := c.httpClient.CallAPI("POST", "/api/game/driver/tribute-return", req, nil); err != nil {
		c.log(fmt.Sprintf("Failed to submit return tribute: %v", err))
	} else {
		c.logVerbose("Return tribute submitted successfully")
	}
}
