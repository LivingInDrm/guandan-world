package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"guandan-world/ai"
	"guandan-world/backend/handlers"
	wsmanager "guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// APIGameTester 简化的游戏API测试客户端
type APIGameTester struct {
	// 配置
	serverURL  string
	authToken  string
	roomID     string
	verbose    bool
	
	// HTTP客户端
	httpClient *http.Client
	
	// WebSocket连接
	wsConn     *websocket.Conn
	wsURL      string
	
	// AI算法
	aiAlgorithms map[int]ai.AutoPlayAlgorithm
	
	// 状态
	gameActive bool
	mu         sync.RWMutex
	
	// 日志
	eventLog   []string
	observer   *TestEventObserver
}

// NewAPIGameTester 创建API游戏测试器
func NewAPIGameTester(serverURL, authToken string, verbose bool) *APIGameTester {
	return &APIGameTester{
		serverURL:    serverURL,
		authToken:    authToken,
		verbose:      verbose,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		aiAlgorithms: make(map[int]ai.AutoPlayAlgorithm),
		eventLog:     make([]string, 0),
		observer:     NewTestEventObserver(verbose),
	}
}

// StartGame 开始游戏测试
func (t *APIGameTester) StartGame(roomID string) error {
	t.roomID = roomID
	
	// 1. 连接WebSocket
	if err := t.connectWebSocket(); err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}
	
	// 2. 启动WebSocket消息处理
	go t.handleWebSocketMessages()
	
	// 2.5. 让WebSocket连接加入房间
	// 注意：在实际应用中，客户端应该通过WebSocket发送join_room消息
	// 但由于我们使用单个连接代表4个玩家，这里跳过此步骤
	
	// 给WebSocket一些时间来建立连接
	time.Sleep(100 * time.Millisecond)
	
	// 3. 调用开始游戏API
	players := []sdk.Player{
		{ID: "test_player_0", Username: "TestPlayer1", Seat: 0},
		{ID: "test_player_1", Username: "TestPlayer2", Seat: 1},
		{ID: "test_player_2", Username: "TestPlayer3", Seat: 2},
		{ID: "test_player_3", Username: "TestPlayer4", Seat: 3},
	}
	
	// 初始化AI算法
	for i := 0; i < 4; i++ {
		t.aiAlgorithms[i] = ai.NewSmartAutoPlayAlgorithm(2) // 从2级开始
	}
	
	req := handlers.StartGameWithDriverRequest{
		RoomID:  roomID,
		Players: players,
	}
	
	if err := t.callAPI("POST", "/api/game/driver/start", req, nil); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	
	t.gameActive = true
	t.log("Game started successfully")
	
	return nil
}

// RunUntilComplete 运行直到游戏结束
func (t *APIGameTester) RunUntilComplete() error {
	// 等待游戏结束
	for t.isGameActive() {
		time.Sleep(100 * time.Millisecond)
	}
	
	t.log("Game completed")
	return nil
}

// connectWebSocket 连接WebSocket
func (t *APIGameTester) connectWebSocket() error {
	// 构建WebSocket URL
	wsURL := fmt.Sprintf("ws://%s/ws?token=%s", t.serverURL, t.authToken)
	
	// 连接WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}
	
	t.wsConn = conn
	t.log("WebSocket connected")
	
	return nil
}

// handleWebSocketMessages 处理WebSocket消息
func (t *APIGameTester) handleWebSocketMessages() {
	defer func() {
		t.wsConn.Close()
		t.setGameActive(false)
	}()
	
	for {
		_, message, err := t.wsConn.ReadMessage()
		if err != nil {
			t.log(fmt.Sprintf("WebSocket read error: %v", err))
			return
		}
		
		var wsMsg wsmanager.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			t.log(fmt.Sprintf("Failed to parse WebSocket message: %v", err))
			continue
		}
		
		// 处理不同类型的消息
		switch wsMsg.Type {
		case wsmanager.MSG_GAME_EVENT:
			t.handleGameEvent(&wsMsg)
		case wsmanager.MSG_GAME_ACTION:
			t.handleGameAction(&wsMsg)
		case wsmanager.MSG_ERROR:
			t.log(fmt.Sprintf("Error message: %v", wsMsg.Data))
		}
	}
}

// handleGameEvent 处理游戏事件
func (t *APIGameTester) handleGameEvent(msg *wsmanager.WSMessage) {
	data := msg.Data.(map[string]interface{})
	eventType := data["event_type"].(string)
	
	// 转发给观察者
	t.observer.OnGameEvent(eventType, data)
	
	// 检查游戏是否结束
	if eventType == "match_ended" {
		t.setGameActive(false)
	}
}

// handleGameAction 处理游戏操作请求
func (t *APIGameTester) handleGameAction(msg *wsmanager.WSMessage) {
	data := msg.Data.(map[string]interface{})
	actionType := data["action_type"].(string)
	playerSeat := int(data["player_seat"].(float64))
	
	t.log(fmt.Sprintf("Action request: %s for player %d", actionType, playerSeat))
	
	switch actionType {
	case "play_decision_required":
		t.handlePlayDecisionRequest(playerSeat, data)
	case "tribute_selection_required":
		t.handleTributeSelectionRequest(playerSeat, data)
	case "return_tribute_required":
		t.handleReturnTributeRequest(playerSeat, data)
	}
}

// handlePlayDecisionRequest 处理出牌决策请求
func (t *APIGameTester) handlePlayDecisionRequest(playerSeat int, data map[string]interface{}) {
	// 解析手牌
	handData := data["hand"].([]interface{})
	hand := make([]*sdk.Card, 0, len(handData))
	for _, cardData := range handData {
		cardMap := cardData.(map[string]interface{})
		card := t.parseCard(cardMap)
		if card != nil {
			hand = append(hand, card)
		}
	}
	
	// 解析trick信息
	trickData := data["trick_info"].(map[string]interface{})
	trickInfo := &sdk.TrickInfo{
		IsLeader: trickData["is_leader"].(bool),
	}
	
	// 如果有leadComp，解析它
	if leadCompData, ok := trickData["lead_comp"]; ok && leadCompData != nil {
		// 这里简化处理，实际需要根据具体格式解析
		trickInfo.LeadComp = nil
	}
	
	// 使用AI算法获取决策
	ai := t.aiAlgorithms[playerSeat]
	selectedCards := ai.SelectCardsToPlay(hand, trickInfo)
	
	// 构建决策
	var action string
	var cardIDs []string
	
	if selectedCards == nil || len(selectedCards) == 0 {
		action = "pass"
	} else {
		action = "play"
		cardIDs = make([]string, len(selectedCards))
		for i, card := range selectedCards {
			cardIDs[i] = card.GetID()
		}
	}
	
	// 提交决策
	req := handlers.PlayDecisionRequest{
		RoomID:     t.roomID,
		PlayerSeat: playerSeat,
		Action:     action,
		CardIDs:    cardIDs,
	}
	
	go func() {
		if err := t.callAPI("POST", "/api/game/driver/play-decision", req, nil); err != nil {
			t.log(fmt.Sprintf("Failed to submit play decision: %v", err))
		}
	}()
}

// handleTributeSelectionRequest 处理贡牌选择请求
func (t *APIGameTester) handleTributeSelectionRequest(playerSeat int, data map[string]interface{}) {
	// 解析选项
	optionsData := data["options"].([]interface{})
	options := make([]*sdk.Card, 0, len(optionsData))
	for _, cardData := range optionsData {
		cardMap := cardData.(map[string]interface{})
		card := t.parseCard(cardMap)
		if card != nil {
			options = append(options, card)
		}
	}
	
	// 选择最大的牌
	var selectedCard *sdk.Card
	if len(options) > 0 {
		selectedCard = options[0]
		for _, card := range options[1:] {
			if card.GreaterThan(selectedCard) {
				selectedCard = card
			}
		}
	}
	
	if selectedCard == nil {
		return
	}
	
	// 提交选择
	req := handlers.TributeSelectionRequest{
		RoomID:     t.roomID,
		PlayerSeat: playerSeat,
		CardID:     selectedCard.GetID(),
	}
	
	go func() {
		if err := t.callAPI("POST", "/api/game/driver/tribute-select", req, nil); err != nil {
			t.log(fmt.Sprintf("Failed to submit tribute selection: %v", err))
		}
	}()
}

// handleReturnTributeRequest 处理还贡请求
func (t *APIGameTester) handleReturnTributeRequest(playerSeat int, data map[string]interface{}) {
	// 解析手牌
	handData := data["hand"].([]interface{})
	hand := make([]*sdk.Card, 0, len(handData))
	for _, cardData := range handData {
		cardMap := cardData.(map[string]interface{})
		card := t.parseCard(cardMap)
		if card != nil {
			hand = append(hand, card)
		}
	}
	
	// 使用AI算法选择还贡牌
	ai := t.aiAlgorithms[playerSeat]
	returnCard := ai.SelectReturnTributeCard(hand, nil)
	
	if returnCard == nil && len(hand) > 0 {
		// 如果AI没有选择，选最小的牌
		returnCard = hand[0]
		for _, card := range hand[1:] {
			if card.LessThan(returnCard) {
				returnCard = card
			}
		}
	}
	
	if returnCard == nil {
		return
	}
	
	// 提交还贡
	req := handlers.ReturnTributeRequest{
		RoomID:     t.roomID,
		PlayerSeat: playerSeat,
		CardID:     returnCard.GetID(),
	}
	
	go func() {
		if err := t.callAPI("POST", "/api/game/driver/tribute-return", req, nil); err != nil {
			t.log(fmt.Sprintf("Failed to submit return tribute: %v", err))
		}
	}()
}

// parseCard 解析卡牌数据
func (t *APIGameTester) parseCard(cardMap map[string]interface{}) *sdk.Card {
	number := int(cardMap["number"].(float64))
	color := cardMap["color"].(string)
	
	// 假设当前等级为2
	card, err := sdk.NewCard(number, color, 2)
	if err != nil {
		return nil
	}
	
	return card
}

// callAPI 调用HTTP API
func (t *APIGameTester) callAPI(method, path string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error
	
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
	}
	
	url := fmt.Sprintf("http://%s%s", t.serverURL, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+t.authToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		var errResp handlers.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("HTTP %d: failed to parse error response", resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Error)
	}
	
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}
	
	return nil
}

// Helper methods

func (t *APIGameTester) log(message string) {
	t.eventLog = append(t.eventLog, message)
	if t.verbose {
		log.Printf("[APIGameTester] %s", message)
	}
}

func (t *APIGameTester) isGameActive() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.gameActive
}

func (t *APIGameTester) setGameActive(active bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.gameActive = active
}

// GetEventLog 获取事件日志
func (t *APIGameTester) GetEventLog() []string {
	return t.eventLog
}

// Close 关闭测试器
func (t *APIGameTester) Close() {
	if t.wsConn != nil {
		t.wsConn.Close()
	}
}

// TestEventObserver 测试事件观察者（简化版）
type TestEventObserver struct {
	verbose bool
}

func NewTestEventObserver(verbose bool) *TestEventObserver {
	return &TestEventObserver{verbose: verbose}
}

func (o *TestEventObserver) OnGameEvent(eventType string, data map[string]interface{}) {
	if !o.verbose {
		return
	}
	
	// 参考 match_simulator_observer.go 的日志格式
	switch eventType {
	case "match_started":
		log.Println("=== Match Started ===")
		
	case "deal_started":
		log.Println("=== Deal Started ===")
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			if dealLevel, ok := eventData["deal_level"].(float64); ok {
				log.Printf("Deal Level: %d", int(dealLevel))
			}
		}
		
	case "tribute_rules_set":
		log.Println("=== Tribute Rules Set ===")
		
	case "tribute_immunity":
		log.Println("=== Tribute Immunity Check ===")
		
	case "trick_started":
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			if leader, ok := eventData["leader"].(float64); ok {
				log.Printf("New Trick Started, Leader: Player %d", int(leader))
			}
		}
		
	case "player_played":
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			if playerSeat, ok := eventData["player_seat"].(float64); ok {
				if cards, ok := eventData["cards"].([]interface{}); ok {
					log.Printf("Player %d played %d cards", int(playerSeat), len(cards))
				}
			}
		}
		
	case "player_passed":
		if playerSeat, ok := data["player_seat"].(float64); ok {
			log.Printf("Player %d passed", int(playerSeat))
		}
		
	case "deal_ended":
		log.Println("=== Deal Ended ===")
		
	case "match_ended":
		log.Println("=== Match Ended ===")
		if eventData, ok := data["event_data"].(map[string]interface{}); ok {
			if winner, ok := eventData["winner"].(float64); ok {
				log.Printf("Winner: Team %d", int(winner))
			}
		}
	}
}

// RunAPIGameTest 运行API游戏测试的便捷函数
func RunAPIGameTest(serverURL, authToken string, verbose bool) error {
	tester := NewAPIGameTester(serverURL, authToken, verbose)
	defer tester.Close()
	
	// 生成唯一的房间ID
	roomID := fmt.Sprintf("test-room-%d", time.Now().Unix())
	
	// 开始游戏
	if err := tester.StartGame(roomID); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	
	// 运行直到完成
	if err := tester.RunUntilComplete(); err != nil {
		return fmt.Errorf("game failed: %w", err)
	}
	
	log.Printf("Test completed successfully. Total events: %d", len(tester.GetEventLog()))
	
	return nil
}