package game

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// DriverService provides complete game management using SDK's GameDriver
// This service encapsulates the full game flow including input handling and event observation
type DriverService struct {
	// Game drivers by room ID
	drivers map[string]*sdk.GameDriver

	// Input providers for each room
	providers map[string]*RoomInputProvider

	// WebSocket manager for real-time communication
	wsManager WSManagerInterface

	// Synchronization
	mu sync.RWMutex
}

// NewDriverService creates a new game driver service
func NewDriverService(wsManager WSManagerInterface) *DriverService {
	return &DriverService{
		drivers:   make(map[string]*sdk.GameDriver),
		providers: make(map[string]*RoomInputProvider),
		wsManager: wsManager,
	}
}

// StartGameWithDriver starts a new game using the GameDriver architecture
func (ds *DriverService) StartGameWithDriver(roomID string, players []sdk.Player) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Validate input
	if roomID == "" {
		return fmt.Errorf("room ID is required")
	}

	if len(players) != 4 {
		return fmt.Errorf("exactly 4 players are required, got %d", len(players))
	}

	// Check if game already exists
	if _, exists := ds.drivers[roomID]; exists {
		return fmt.Errorf("game already exists for room %s", roomID)
	}

	// Create game engine
	engine := sdk.NewGameEngine()

	// Create game driver with shorter timeouts for testing
	config := sdk.DefaultGameDriverConfig()
	config.PlayDecisionTimeout = 1 * time.Second // Short timeout for tests
	config.TributeTimeout = 1 * time.Second
	driver := sdk.NewGameDriver(engine, config)

	// Create and set input provider for this room
	provider := NewRoomInputProvider(roomID, ds.wsManager)
	driver.SetInputProvider(provider)

	// Add WebSocket observer for real-time events
	observer := NewWebSocketObserver(roomID, ds.wsManager)
	driver.AddObserver(observer)

	// Store driver and provider
	ds.drivers[roomID] = driver
	ds.providers[roomID] = provider

	// Start the match in a goroutine
	go func() {
		log.Printf("Starting match for room %s with GameDriver", roomID)

		result, err := driver.RunMatch(players)
		if err != nil {
			log.Printf("Match error for room %s: %v", roomID, err)
			// Send error event to clients
			ds.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
				Type: websocket.MSG_ERROR,
				Data: map[string]interface{}{
					"error":   err.Error(),
					"room_id": roomID,
				},
				Timestamp: time.Now(),
			})
		} else {
			log.Printf("Match completed for room %s, winner: team %d", roomID, result.Winner)
			// Match completed event is already sent by the observer
		}

		// Clean up after match
		ds.mu.Lock()
		delete(ds.drivers, roomID)
		delete(ds.providers, roomID)
		ds.mu.Unlock()
	}()

	return nil
}

// SubmitPlayDecision submits a player's play decision to the driver
func (ds *DriverService) SubmitPlayDecision(roomID string, playerSeat int, decision *sdk.PlayDecision) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Submit decision to the input provider
	return provider.SubmitPlayDecision(playerSeat, decision)
}

// SubmitTributeSelection submits a tribute selection to the driver
func (ds *DriverService) SubmitTributeSelection(roomID string, playerSeat int, cardID string) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Find the card by ID
	card, err := ds.findCardByID(provider, playerSeat, cardID)
	if err != nil {
		return err
	}

	// Submit selection to the input provider
	return provider.SubmitTributeSelection(playerSeat, card)
}

// SubmitReturnTribute submits a return tribute to the driver
func (ds *DriverService) SubmitReturnTribute(roomID string, playerSeat int, cardID string) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Find the card by ID
	card, err := ds.findCardByID(provider, playerSeat, cardID)
	if err != nil {
		return err
	}

	// Submit return to the input provider
	return provider.SubmitReturnTribute(playerSeat, card)
}

// GetGameStatus gets the current game status for a room
func (ds *DriverService) GetGameStatus(roomID string) (map[string]interface{}, error) {
	ds.mu.RLock()
	driver, exists := ds.drivers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no active game for room %s", roomID)
	}

	// Get engine from driver
	engine := driver.GetEngine()
	if engine == nil {
		return nil, fmt.Errorf("no game engine for room %s", roomID)
	}

	// Get current state
	gameState := engine.GetGameState()
	dealStatus := engine.GetCurrentDealStatus()
	turnInfo := engine.GetCurrentTurnInfo()
	matchDetails := engine.GetMatchDetails()

	status := map[string]interface{}{
		"room_id":     roomID,
		"game_status": gameState.Status,
		"deal_status": dealStatus,
		"timestamp":   time.Now(),
	}

	if turnInfo != nil {
		status["turn_info"] = turnInfo
	}

	if matchDetails != nil {
		status["match_details"] = matchDetails
	}

	return status, nil
}

// StopGame stops the game for a room
func (ds *DriverService) StopGame(roomID string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	_, exists := ds.drivers[roomID]
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Cancel any pending input requests
	if provider, ok := ds.providers[roomID]; ok {
		provider.CancelAll()
	}

	// Clean up
	delete(ds.drivers, roomID)
	delete(ds.providers, roomID)

	// Notify clients
	ds.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
		Type: websocket.MSG_GAME_EVENT,
		Data: map[string]interface{}{
			"event_type": "game_stopped",
			"room_id":    roomID,
		},
		Timestamp: time.Now(),
	})

	log.Printf("Game stopped for room %s", roomID)
	return nil
}

// findCardByID finds a card by its ID from the provider's context
func (ds *DriverService) findCardByID(provider *RoomInputProvider, playerSeat int, cardID string) (*sdk.Card, error) {
	// Get the last options provided to this player
	options := provider.GetLastOptions(playerSeat)
	if options == nil {
		return nil, fmt.Errorf("no card options available for player %d", playerSeat)
	}

	// Find the card
	for _, card := range options {
		if card.GetID() == cardID {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card %s not found in available options", cardID)
}

// RoomInputProvider implements sdk.PlayerInputProvider for a specific room
type RoomInputProvider struct {
	roomID    string
	wsManager WSManagerInterface

	// Channels for receiving player decisions
	playDecisions     map[int]chan *sdk.PlayDecision
	tributeSelections map[int]chan *sdk.Card
	returnTributes    map[int]chan *sdk.Card

	// Store last options for card lookup
	lastOptions map[int][]*sdk.Card

	mu sync.RWMutex
}

// NewRoomInputProvider creates a new input provider for a room
func NewRoomInputProvider(roomID string, wsManager WSManagerInterface) *RoomInputProvider {
	return &RoomInputProvider{
		roomID:            roomID,
		wsManager:         wsManager,
		playDecisions:     make(map[int]chan *sdk.PlayDecision),
		tributeSelections: make(map[int]chan *sdk.Card),
		returnTributes:    make(map[int]chan *sdk.Card),
		lastOptions:       make(map[int][]*sdk.Card),
	}
}

// RequestPlayDecision implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestPlayDecision(ctx context.Context, playerSeat int, hand []*sdk.Card, trickInfo *sdk.TrickInfo) (*sdk.PlayDecision, error) {
	// Create channel for this request
	rip.mu.Lock()
	decisionChan := make(chan *sdk.PlayDecision, 1)
	rip.playDecisions[playerSeat] = decisionChan
	rip.mu.Unlock()

	defer func() {
		rip.mu.Lock()
		delete(rip.playDecisions, playerSeat)
		rip.mu.Unlock()
	}()

	// Send request to player via WebSocket
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_ACTION,
		Data: map[string]interface{}{
			"action_type": "play_decision_required",
			"player_seat": playerSeat,
			"hand":        hand,
			"trick_info":  trickInfo,
			"timeout":     30, // seconds
		},
		Timestamp: time.Now(),
	}

	// Get player ID and send message
	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send play request: %w", err)
	}

	// If no context provided, create one with default timeout
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Wait for decision or timeout
	select {
	case decision := <-decisionChan:
		// 添加 nil 检查防止空指针异常
		if decision == nil {
			return nil, fmt.Errorf("received nil decision from player %d", playerSeat)
		}
		return decision, nil
	case <-ctx.Done():
		// Timeout - return a default decision (pass if not leader, play smallest card if leader)
		if trickInfo.IsLeader && len(hand) > 0 {
			// Play smallest single card
			smallestCard := hand[0]
			for _, card := range hand {
				if card.LessThan(smallestCard) {
					smallestCard = card
				}
			}
			return &sdk.PlayDecision{
				Action: sdk.ActionPlay,
				Cards:  []*sdk.Card{smallestCard},
			}, nil
		}
		return &sdk.PlayDecision{
			Action: sdk.ActionPass,
		}, nil
	}
}

// RequestTributeSelection implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestTributeSelection(ctx context.Context, playerSeat int, options []*sdk.Card) (*sdk.Card, error) {
	// Store options for lookup
	rip.mu.Lock()
	rip.lastOptions[playerSeat] = options
	selectionChan := make(chan *sdk.Card, 1)
	rip.tributeSelections[playerSeat] = selectionChan
	rip.mu.Unlock()

	defer func() {
		rip.mu.Lock()
		delete(rip.tributeSelections, playerSeat)
		rip.mu.Unlock()
	}()

	// Send request to player
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_ACTION,
		Data: map[string]interface{}{
			"action_type": "tribute_selection_required",
			"player_seat": playerSeat,
			"options":     options,
			"timeout":     20, // seconds
		},
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send tribute selection request: %w", err)
	}

	// If no context provided, create one with default timeout
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
	}

	// Wait for selection or timeout
	select {
	case card := <-selectionChan:
		return card, nil
	case <-ctx.Done():
		// Timeout - select the largest card
		if len(options) > 0 {
			largestCard := options[0]
			for _, card := range options {
				if card.GreaterThan(largestCard) {
					largestCard = card
				}
			}
			return largestCard, nil
		}
		return nil, fmt.Errorf("no options available")
	}
}

// RequestReturnTribute implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestReturnTribute(ctx context.Context, playerSeat int, hand []*sdk.Card) (*sdk.Card, error) {
	// Store hand as options for lookup
	rip.mu.Lock()
	rip.lastOptions[playerSeat] = hand
	returnChan := make(chan *sdk.Card, 1)
	rip.returnTributes[playerSeat] = returnChan
	rip.mu.Unlock()

	defer func() {
		rip.mu.Lock()
		delete(rip.returnTributes, playerSeat)
		rip.mu.Unlock()
	}()

	// Send request to player
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_ACTION,
		Data: map[string]interface{}{
			"action_type": "return_tribute_required",
			"player_seat": playerSeat,
			"hand":        hand,
			"timeout":     20, // seconds
		},
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send return tribute request: %w", err)
	}

	// If no context provided, create one with default timeout
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
	}

	// Wait for return or timeout
	select {
	case card := <-returnChan:
		return card, nil
	case <-ctx.Done():
		// Timeout - return the smallest non-trump card
		if len(hand) > 0 {
			smallestCard := hand[0]
			for _, card := range hand {
				if card.LessThan(smallestCard) {
					smallestCard = card
				}
			}
			return smallestCard, nil
		}
		return nil, fmt.Errorf("no cards available")
	}
}

// SubmitPlayDecision submits a play decision from a player
func (rip *RoomInputProvider) SubmitPlayDecision(playerSeat int, decision *sdk.PlayDecision) error {
	// 添加输入验证防止空指针异常
	if decision == nil {
		return fmt.Errorf("decision cannot be nil for player %d", playerSeat)
	}

	rip.mu.RLock()
	decisionChan, exists := rip.playDecisions[playerSeat]
	rip.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no pending play decision for player %d", playerSeat)
	}

	select {
	case decisionChan <- decision:
		return nil
	default:
		return fmt.Errorf("decision channel is full for player %d", playerSeat)
	}
}

// SubmitTributeSelection submits a tribute selection from a player
func (rip *RoomInputProvider) SubmitTributeSelection(playerSeat int, card *sdk.Card) error {
	rip.mu.RLock()
	selectionChan, exists := rip.tributeSelections[playerSeat]
	rip.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no pending tribute selection for player %d", playerSeat)
	}

	select {
	case selectionChan <- card:
		return nil
	default:
		return fmt.Errorf("selection channel is full for player %d", playerSeat)
	}
}

// SubmitReturnTribute submits a return tribute from a player
func (rip *RoomInputProvider) SubmitReturnTribute(playerSeat int, card *sdk.Card) error {
	rip.mu.RLock()
	returnChan, exists := rip.returnTributes[playerSeat]
	rip.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no pending return tribute for player %d", playerSeat)
	}

	select {
	case returnChan <- card:
		return nil
	default:
		return fmt.Errorf("return channel is full for player %d", playerSeat)
	}
}

// GetLastOptions returns the last options provided to a player
func (rip *RoomInputProvider) GetLastOptions(playerSeat int) []*sdk.Card {
	rip.mu.RLock()
	defer rip.mu.RUnlock()
	return rip.lastOptions[playerSeat]
}

// CancelAll cancels all pending input requests
func (rip *RoomInputProvider) CancelAll() {
	rip.mu.Lock()
	defer rip.mu.Unlock()

	// Close all channels to unblock waiting goroutines
	for _, ch := range rip.playDecisions {
		close(ch)
	}
	for _, ch := range rip.tributeSelections {
		close(ch)
	}
	for _, ch := range rip.returnTributes {
		close(ch)
	}

	// Clear maps
	rip.playDecisions = make(map[int]chan *sdk.PlayDecision)
	rip.tributeSelections = make(map[int]chan *sdk.Card)
	rip.returnTributes = make(map[int]chan *sdk.Card)
	rip.lastOptions = make(map[int][]*sdk.Card)
}

// sendToPlayer sends a message to a specific player
func (rip *RoomInputProvider) sendToPlayer(playerSeat int, message *websocket.WSMessage) error {
	// For now, broadcast to room with player seat info
	// In a real implementation, this would send to the specific player
	message.Data.(map[string]interface{})["room_id"] = rip.roomID
	rip.wsManager.BroadcastToRoom(rip.roomID, message)
	return nil
}

// WebSocketObserver implements sdk.EventObserver for WebSocket broadcasting
type WebSocketObserver struct {
	roomID    string
	wsManager WSManagerInterface
}

// NewWebSocketObserver creates a new WebSocket observer
func NewWebSocketObserver(roomID string, wsManager WSManagerInterface) *WebSocketObserver {
	return &WebSocketObserver{
		roomID:    roomID,
		wsManager: wsManager,
	}
}

// OnGameEvent implements sdk.EventObserver
func (wso *WebSocketObserver) OnGameEvent(event *sdk.GameEvent) {
	// Convert SDK event to WebSocket message
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

	// Broadcast to all players in the room
	wso.wsManager.BroadcastToRoom(wso.roomID, wsMessage)

	// Log significant events
	switch event.Type {
	case sdk.EventMatchStarted, sdk.EventMatchEnded,
		sdk.EventDealStarted, sdk.EventDealEnded,
		sdk.EventTributeCompleted:
		log.Printf("Game event %s for room %s", event.Type, wso.roomID)
	}
}
