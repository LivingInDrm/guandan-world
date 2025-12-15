package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"guandan-world/backend/websocket"
	actionpb "guandan-world/proto/action"
	eventpb "guandan-world/proto/event"
	viewpb "guandan-world/proto/view"
	"guandan-world/sdk"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const pendingEventsThreshold = 50

// protoJSONMarshaler is the configuration for serializing proto messages to JSON
// Uses camelCase field names and emits all fields including zero values
// EmitUnpopulated must be true to ensure playerSeat=0 is serialized (seat 0 is valid)
// UseEnumNumbers serializes enums as integers instead of strings for type safety and efficiency
var protoJSONMarshaler = protojson.MarshalOptions{
	UseProtoNames:   false, // Use camelCase field names (matchId, dealStatus)
	EmitUnpopulated: true,  // Emit all fields, including zero values (required for playerSeat=0)
	UseEnumNumbers:  true,  // Serialize enums as integers (VictoryType: 1 instead of "VICTORY_TYPE_DOUBLE_DOWN")
}

// marshalProtoToRawJSON serializes a proto message to json.RawMessage
// using the configured protoJSONMarshaler (camelCase fields, all fields including zero values)
// Returns nil and logs error if marshaling fails
func marshalProtoToRawJSON(msg proto.Message, logPrefix string) json.RawMessage {
	jsonBytes, err := protoJSONMarshaler.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s to JSON: %v", logPrefix, err)
		return nil
	}
	return json.RawMessage(jsonBytes)
}

// getEnvironment returns the current environment (test, dev, prod)
// Defaults to "prod" if APP_ENV is not set
func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "prod"
	}
	return env
}

// getRemainingTimeout calculates remaining seconds until context deadline
// Returns 0 if deadline has passed or no deadline exists
func getRemainingTimeout(ctx context.Context) (int, bool) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, false
	}
	
	secs := int(time.Until(deadline).Seconds())
	if secs < 0 {
		secs = 0
	}
	return secs, true
}

// WSManagerInterface defines the interface for WebSocket management
type WSManagerInterface interface {
	BroadcastToRoom(roomID string, message *websocket.WSMessage)
	SendToPlayer(playerID string, message *websocket.WSMessage) error
}

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

	// Create game driver with timeout configuration
	config := sdk.DefaultGameDriverConfig()
	
	// Configure timeout strategy (default strategy for automated decisions)
	config.TimeoutStrategy = sdk.NewDefaultTimeoutStrategy()
	
	// Configure timeout durations based on environment
	// DefaultGameDriverConfig provides production defaults: 30s for play decisions, 20s for tribute
	// Override for test/development environments to speed up iteration
	env := getEnvironment()
	if env == "test" || env == "dev" {
		config.PlayDecisionTimeout = 10 * time.Second
		config.TributeTimeout = 10 * time.Second
	}
	// else: use production defaults from DefaultGameDriverConfig (30s/20s)
	
	driver := sdk.NewGameDriver(engine, config)

	// Create and set input provider for this room
	provider := NewRoomInputProvider(roomID, ds.wsManager, players)
	driver.SetInputProvider(provider)

	// Add WebSocket observer for real-time events with engine reference
	observer := NewWebSocketObserver(roomID, ds.wsManager, engine)
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
		observer.Stop() // Stop event processing and handle remaining events
		ds.mu.Lock()
		delete(ds.drivers, roomID)
		delete(ds.providers, roomID)
		ds.mu.Unlock()
	}()

	return nil
}

// SubmitPlayDecision submits a player's play decision to the driver
func (ds *DriverService) SubmitPlayDecision(
	roomID string,
	playerSeat int,
	action string,
	deckIndexes []int,
) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Convert DeckIndexes to Cards by looking up in player's hand
	var cards []*sdk.Card
	if action == "play" && len(deckIndexes) > 0 {
		for _, deckIndex := range deckIndexes {
			card, err := ds.findCardByDeckIndex(provider, playerSeat, deckIndex)
			if err != nil {
				return fmt.Errorf("failed to find card with deck_index %d: %w", deckIndex, err)
			}
			cards = append(cards, card)
		}
	}

	// Create SDK PlayDecision object
	decision := &sdk.PlayDecision{
		Action: sdk.ActionType(action),
		Cards:  cards,
	}

	// Submit decision to the input provider
	return provider.SubmitPlayDecision(playerSeat, decision)
}

// SubmitTributeSelection submits a tribute selection to the driver
func (ds *DriverService) SubmitTributeSelection(roomID string, playerSeat int, deckIndex int) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Find the card by DeckIndex
	card, err := ds.findCardByDeckIndex(provider, playerSeat, deckIndex)
	if err != nil {
		return err
	}

	// Submit selection to the input provider
	return provider.SubmitTributeSelection(playerSeat, card)
}

// SubmitReturnTribute submits a return tribute to the driver
func (ds *DriverService) SubmitReturnTribute(roomID string, playerSeat int, deckIndex int) error {
	ds.mu.RLock()
	provider, exists := ds.providers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// Find the card by DeckIndex
	card, err := ds.findCardByDeckIndex(provider, playerSeat, deckIndex)
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

// HandlePlayerReconnect handles player reconnection and syncs game state
func (ds *DriverService) HandlePlayerReconnect(roomID, playerID string) error {
	ds.mu.RLock()
	driver, exists := ds.drivers[roomID]
	ds.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	engine := driver.GetEngine()
	if engine == nil {
		return fmt.Errorf("no game engine for room %s", roomID)
	}

	gameState := engine.GetGameState()
	if gameState == nil || gameState.CurrentMatch == nil {
		return fmt.Errorf("no active match for room %s", roomID)
	}

	playerSeat := -1
	for seat, player := range gameState.CurrentMatch.Players {
		if player != nil && player.ID == playerID {
			playerSeat = seat
			break
		}
	}

	if playerSeat < 0 {
		return fmt.Errorf("player %s not found in room %s", playerID, roomID)
	}

	_, err := engine.HandlePlayerReconnect(playerSeat)
	if err != nil {
		log.Printf("Failed to handle player reconnect in SDK: %v", err)
	}

	playerView := engine.GetPlayerView(playerSeat)
	if playerView != nil {
		playerViewJSON := marshalProtoToRawJSON(playerView, "PlayerView")
		if playerViewJSON != nil {
			wsMessage := &websocket.WSMessage{
				Type: websocket.MSG_PLAYER_VIEW,
				Data: map[string]interface{}{
					"player_view": playerViewJSON,
					"event_type":  "reconnect",
					"player_seat": playerSeat,
				},
				Timestamp: time.Now(),
				PlayerID:  playerID,
			}

			if err := ds.wsManager.SendToPlayer(playerID, wsMessage); err != nil {
				log.Printf("Failed to send player view on reconnect: %v", err)
			}
		}

		if playerView.DealStatus == viewpb.DealStatus_DEAL_STATUS_TRIBUTE {
			tributeView := engine.GetTributeView(playerSeat)
			if tributeView != nil {
				tributeViewJSON := marshalProtoToRawJSON(tributeView, "TributeView")
				if tributeViewJSON != nil {
					tributeMsg := &websocket.WSMessage{
						Type: "tribute_view",
						Data: map[string]interface{}{
							"tribute_view": tributeViewJSON,
							"event_type":   "reconnect",
							"player_seat":  playerSeat,
						},
						Timestamp: time.Now(),
						PlayerID:  playerID,
					}

					if err := ds.wsManager.SendToPlayer(playerID, tributeMsg); err != nil {
						log.Printf("Failed to send tribute view on reconnect: %v", err)
					}
				}
			}
		}
	}

	log.Printf("Player %s reconnected to room %s (seat %d)", playerID, roomID, playerSeat)
	return nil
}

// StopGame stops the game for a room
func (ds *DriverService) StopGame(roomID string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	driver, exists := ds.drivers[roomID]
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}

	// 1. First, cancel SDK layer's game loop
	// This triggers context cancellation which propagates to all pending operations
	driver.CancelMatch()

	// 2. Then cancel all pending input requests
	// This closes all input channels to unblock waiting goroutines
	if provider, ok := ds.providers[roomID]; ok {
		provider.CancelAll()
	}

	// 3. Clean up resources
	delete(ds.drivers, roomID)
	delete(ds.providers, roomID)

	// 4. Notify clients
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

// findCardByDeckIndex finds a card by its DeckIndex from the provider's context
func (ds *DriverService) findCardByDeckIndex(provider *RoomInputProvider, playerSeat int, deckIndex int) (*sdk.Card, error) {
	// Get the last options provided to this player
	options := provider.GetLastOptions(playerSeat)
	if options == nil {
		return nil, fmt.Errorf("no card options available for player %d", playerSeat)
	}

	// Find the card
	for _, card := range options {
		if card.DeckIndex == deckIndex {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card with deck_index %d not found in available options", deckIndex)
}

// RoomInputProvider implements sdk.PlayerInputProvider for a specific room
type RoomInputProvider struct {
	roomID         string
	wsManager      WSManagerInterface
	seatToPlayerID map[int]string // Maps seat number to player ID

	// Channels for receiving player decisions
	playDecisions     map[int]chan *sdk.PlayDecision
	tributeSelections map[int]chan *sdk.Card
	returnTributes    map[int]chan *sdk.Card

	// Store last options for card lookup
	lastOptions map[int][]*sdk.Card

	mu sync.RWMutex
}

// NewRoomInputProvider creates a new input provider for a room
func NewRoomInputProvider(roomID string, wsManager WSManagerInterface, players []sdk.Player) *RoomInputProvider {
	// Build seat to player ID mapping
	seatToPlayerID := make(map[int]string)
	for _, p := range players {
		seatToPlayerID[p.Seat] = p.ID
	}

	return &RoomInputProvider{
		roomID:            roomID,
		wsManager:         wsManager,
		seatToPlayerID:    seatToPlayerID,
		playDecisions:     make(map[int]chan *sdk.PlayDecision),
		tributeSelections: make(map[int]chan *sdk.Card),
		returnTributes:    make(map[int]chan *sdk.Card),
		lastOptions:       make(map[int][]*sdk.Card),
	}
}

// RequestPlayDecision implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestPlayDecision(ctx context.Context, playerSeat int, hand []*sdk.Card, trickInfo *sdk.TrickInfo) (*sdk.PlayDecision, error) {
	// Defensive check: context must not be nil
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}
	
	// Create channel for this request and store hand for card lookup
	rip.mu.Lock()
	rip.lastOptions[playerSeat] = hand  // ✅ 保存SDK的真实Card对象
	decisionChan := make(chan *sdk.PlayDecision, 1)
	rip.playDecisions[playerSeat] = decisionChan
	rip.mu.Unlock()

	defer func() {
		rip.mu.Lock()
		delete(rip.playDecisions, playerSeat)
		rip.mu.Unlock()
	}()

	// Build proto GameAction message
	gameAction := &actionpb.GameAction{
		ActionType: actionpb.GameActionType_GAME_ACTION_TYPE_PLAY_DECISION,
		PlayerSeat: int32(playerSeat),
		Hand:       sdk.ConvertCardsToProto(hand),
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		gameAction.Timeout = int32(secs)
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      map[string]interface{}{
			"game_action": marshalProtoToRawJSON(gameAction, "GameAction"),
		},
		Timestamp: time.Now(),
	}

	// Get player ID and send message
	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send play request: %w", err)
	}

	// Broadcast deadline to all players in the room
	if deadline, ok := ctx.Deadline(); ok {
		turnDeadline := &actionpb.TurnDeadline{
			PlayerSeat:   int32(playerSeat),
			ActionType:   actionpb.GameActionType_GAME_ACTION_TYPE_PLAY_DECISION,
			DeadlineAtMs: deadline.UnixMilli(),
		}
		rip.wsManager.BroadcastToRoom(rip.roomID, &websocket.WSMessage{
			Type: websocket.MSG_TURN_DEADLINE,
			Data: map[string]interface{}{
				"turn_deadline": marshalProtoToRawJSON(turnDeadline, "TurnDeadline"),
			},
			Timestamp: time.Now(),
		})
	}

	// Wait for decision or context cancellation
	select {
	case decision, ok := <-decisionChan:
		if !ok {
			return nil, fmt.Errorf("play decision request canceled for player %d", playerSeat)
		}
		if decision == nil {
			return nil, fmt.Errorf("received nil decision for player %d", playerSeat)
		}
		return decision, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RequestTributeSelection implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestTributeSelection(ctx context.Context, playerSeat int, options []*sdk.Card) (*sdk.Card, error) {
	// Defensive check: context must not be nil
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}
	
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

	// Build proto GameAction message
	gameAction := &actionpb.GameAction{
		ActionType: actionpb.GameActionType_GAME_ACTION_TYPE_TRIBUTE_SELECTION,
		PlayerSeat: int32(playerSeat),
		Options:    sdk.ConvertCardsToProto(options),
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		gameAction.Timeout = int32(secs)
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      map[string]interface{}{
			"game_action": marshalProtoToRawJSON(gameAction, "GameAction"),
		},
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send tribute selection request: %w", err)
	}

	// Broadcast deadline to all players in the room
	if deadline, ok := ctx.Deadline(); ok {
		turnDeadline := &actionpb.TurnDeadline{
			PlayerSeat:   int32(playerSeat),
			ActionType:   actionpb.GameActionType_GAME_ACTION_TYPE_TRIBUTE_SELECTION,
			DeadlineAtMs: deadline.UnixMilli(),
		}
		rip.wsManager.BroadcastToRoom(rip.roomID, &websocket.WSMessage{
			Type: websocket.MSG_TURN_DEADLINE,
			Data: map[string]interface{}{
				"turn_deadline": marshalProtoToRawJSON(turnDeadline, "TurnDeadline"),
			},
			Timestamp: time.Now(),
		})
	}

	// Wait for selection or context cancellation
	select {
	case card, ok := <-selectionChan:
		if !ok {
			return nil, fmt.Errorf("tribute selection request canceled for player %d", playerSeat)
		}
		if card == nil {
			return nil, fmt.Errorf("received nil tribute selection for player %d", playerSeat)
		}
		return card, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RequestReturnTribute implements sdk.PlayerInputProvider
func (rip *RoomInputProvider) RequestReturnTribute(ctx context.Context, playerSeat int, hand []*sdk.Card) (*sdk.Card, error) {
	// Defensive check: context must not be nil
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}
	
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

	// Build proto GameAction message
	// Note: for return tribute, hand cards are sent as options (cards available for return)
	gameAction := &actionpb.GameAction{
		ActionType: actionpb.GameActionType_GAME_ACTION_TYPE_RETURN_TRIBUTE,
		PlayerSeat: int32(playerSeat),
		Options:    sdk.ConvertCardsToProto(hand),
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		gameAction.Timeout = int32(secs)
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      map[string]interface{}{
			"game_action": marshalProtoToRawJSON(gameAction, "GameAction"),
		},
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send return tribute request: %w", err)
	}

	// Broadcast deadline to all players in the room
	if deadline, ok := ctx.Deadline(); ok {
		turnDeadline := &actionpb.TurnDeadline{
			PlayerSeat:   int32(playerSeat),
			ActionType:   actionpb.GameActionType_GAME_ACTION_TYPE_RETURN_TRIBUTE,
			DeadlineAtMs: deadline.UnixMilli(),
		}
		rip.wsManager.BroadcastToRoom(rip.roomID, &websocket.WSMessage{
			Type: websocket.MSG_TURN_DEADLINE,
			Data: map[string]interface{}{
				"turn_deadline": marshalProtoToRawJSON(turnDeadline, "TurnDeadline"),
			},
			Timestamp: time.Now(),
		})
	}

	// Wait for return or context cancellation
	select {
	case card, ok := <-returnChan:
		if !ok {
			return nil, fmt.Errorf("return tribute request canceled for player %d", playerSeat)
		}
		if card == nil {
			return nil, fmt.Errorf("received nil return tribute for player %d", playerSeat)
		}
		return card, nil
	case <-ctx.Done():
		return nil, ctx.Err()
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
	// Validate input to prevent nil pointer exceptions
	if card == nil {
		return fmt.Errorf("card cannot be nil for player %d", playerSeat)
	}
	
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
	// Validate input to prevent nil pointer exceptions
	if card == nil {
		return fmt.Errorf("card cannot be nil for player %d", playerSeat)
	}
	
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
	rip.seatToPlayerID = make(map[int]string)
}

// sendToPlayer sends a message to a specific player
func (rip *RoomInputProvider) sendToPlayer(playerSeat int, message *websocket.WSMessage) error {
	// Look up player ID from seat number
	rip.mu.RLock()
	playerID, exists := rip.seatToPlayerID[playerSeat]
	rip.mu.RUnlock()

	if !exists {
		return fmt.Errorf("player seat %d not found in room %s", playerSeat, rip.roomID)
	}

	// Send message to specific player
	return rip.wsManager.SendToPlayer(playerID, message)
}

// WebSocketObserver implements sdk.EventObserver for WebSocket broadcasting
// Uses an ordered event queue to ensure events are processed in seq order
type WebSocketObserver struct {
	roomID    string
	wsManager WSManagerInterface
	engine    sdk.GameEngineInterface // Reference to engine for accessing player views

	// Ordered event processing
	eventQueue    chan *sdk.GameEvent        // Event input channel
	stopChan      chan struct{}              // Stop signal
	lastSeq       int64                      // Last processed seq
	pendingEvents map[int64]*sdk.GameEvent   // Out-of-order events cache
	mu            sync.Mutex                 // Protects pendingEvents and lastSeq
}

// NewWebSocketObserver creates a new WebSocket observer
func NewWebSocketObserver(roomID string, wsManager WSManagerInterface, engine sdk.GameEngineInterface) *WebSocketObserver {
	wso := &WebSocketObserver{
		roomID:        roomID,
		wsManager:     wsManager,
		engine:        engine,
		eventQueue:    make(chan *sdk.GameEvent, 256),
		stopChan:      make(chan struct{}),
		lastSeq:       0,
		pendingEvents: make(map[int64]*sdk.GameEvent),
	}
	go wso.processEventLoop()
	return wso
}

// OnGameEvent implements sdk.EventObserver
// Events are queued and processed in seq order by processEventLoop
func (wso *WebSocketObserver) OnGameEvent(event *sdk.GameEvent) {
	select {
	case wso.eventQueue <- event:
		// Successfully queued
	case <-wso.stopChan:
		// Observer stopped
	default:
		log.Printf("[WebSocketObserver] Event queue full for room %s, seq=%d", wso.roomID, event.Seq)
	}
}

// Stop stops the event processing loop and processes remaining events
func (wso *WebSocketObserver) Stop() {
	close(wso.stopChan)

	// Process remaining events in queue
	for {
		select {
		case event := <-wso.eventQueue:
			wso.processEventInOrder(event)
		default:
			return
		}
	}
}

// processEventLoop processes events from the queue in order
func (wso *WebSocketObserver) processEventLoop() {
	for {
		select {
		case event := <-wso.eventQueue:
			wso.processEventInOrder(event)
		case <-wso.stopChan:
			return
		}
	}
}

// processEventInOrder ensures events are processed in seq order
func (wso *WebSocketObserver) processEventInOrder(event *sdk.GameEvent) {
	wso.mu.Lock()
	defer wso.mu.Unlock()

	seq := event.Seq

	// MATCH_STARTED 特殊处理：只有当 seq <= lastSeq 时（新比赛 seq 重新计数）才需要立即重置
	if event.Type == eventpb.EventType_EVENT_TYPE_MATCH_STARTED && seq <= wso.lastSeq {
		if len(wso.pendingEvents) > 0 {
			log.Printf("[WebSocketObserver] MATCH_STARTED forcing reset, discarding %d pending events for room %s", 
				len(wso.pendingEvents), wso.roomID)
		}
		wso.lastSeq = seq - 1
		wso.pendingEvents = make(map[int64]*sdk.GameEvent)
	}

	// First event: auto-initialize lastSeq (only MATCH_STARTED can trigger initialization)
	if wso.lastSeq == 0 && len(wso.pendingEvents) == 0 && event.Type == eventpb.EventType_EVENT_TYPE_MATCH_STARTED {
		wso.lastSeq = seq - 1
	}

	// Already processed, ignore
	if seq <= wso.lastSeq {
		return
	}

	// Exactly the next seq
	if seq == wso.lastSeq+1 {
		wso.handleEvent(event)
		wso.lastSeq = seq

		// Process consecutive cached events
		for {
			nextSeq := wso.lastSeq + 1
			if nextEvent, exists := wso.pendingEvents[nextSeq]; exists {
				wso.handleEvent(nextEvent)
				delete(wso.pendingEvents, nextSeq)
				wso.lastSeq = nextSeq
			} else {
				break
			}
		}
		return
	}

	// Out of order, cache it
	wso.pendingEvents[seq] = event

	// Prevent cache from growing too large
	if len(wso.pendingEvents) > pendingEventsThreshold {
		log.Printf("[WebSocketObserver] ERROR: Too many pending events for room %s, lastSeq=%d, pendingCount=%d, forcing process",
			wso.roomID, wso.lastSeq, len(wso.pendingEvents))
		wso.forceProcessPending()
	}
}

// forceProcessPending processes all pending events in order, skipping any gaps
func (wso *WebSocketObserver) forceProcessPending() {
	// Sort all pending seqs
	seqs := make([]int64, 0, len(wso.pendingEvents))
	for seq := range wso.pendingEvents {
		seqs = append(seqs, seq)
	}
	sort.Slice(seqs, func(i, j int) bool { return seqs[i] < seqs[j] })

	// Calculate missing seqs for logging
	var missingSeqs []int64
	if len(seqs) > 0 {
		for expected := wso.lastSeq + 1; expected < seqs[0]; expected++ {
			missingSeqs = append(missingSeqs, expected)
		}
		for i := 0; i < len(seqs)-1; i++ {
			for gap := seqs[i] + 1; gap < seqs[i+1]; gap++ {
				missingSeqs = append(missingSeqs, gap)
			}
		}
	}

	// Log warning if there are missing seqs
	if len(missingSeqs) > 0 {
		log.Printf("[WebSocketObserver] ERROR: Skipping missing seqs for room %s: %v (lastSeq=%d, processing=%v)",
			wso.roomID, missingSeqs, wso.lastSeq, seqs)
	}

	// Process all in order
	for _, seq := range seqs {
		wso.handleEvent(wso.pendingEvents[seq])
	}

	if len(seqs) > 0 {
		wso.lastSeq = seqs[len(seqs)-1]
	}
	wso.pendingEvents = make(map[int64]*sdk.GameEvent)
}

// handleEvent processes a single event (original OnGameEvent logic)
func (wso *WebSocketObserver) handleEvent(event *sdk.GameEvent) {
	// Serialize the entire event to JSON using protojson
	eventJSON := marshalProtoToRawJSON(event, "GameEvent")
	if eventJSON == nil {
		return
	}

	// Send GameEvent directly without wrapping
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_EVENT,
		Data:      eventJSON, // Direct GameEvent JSON (flattened structure)
		Timestamp: time.UnixMilli(event.CreatedAtMs),
	}

	// Broadcast to all players in the room
	wso.wsManager.BroadcastToRoom(wso.roomID, wsMessage)

	// Send player-specific views only for key events
	// These are events where hand cards change or game phase transitions occur
	switch event.Type {
	case eventpb.EventType_EVENT_TYPE_MATCH_STARTED, // Match begins
		eventpb.EventType_EVENT_TYPE_DEAL_STARTED,           // New deal starts
		eventpb.EventType_EVENT_TYPE_CARDS_DEALT,            // Cards dealt to players
		eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED,        // Tribute phase starts
		eventpb.EventType_EVENT_TYPE_TRIBUTE_EXEMPTED,       // Tribute exempted
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SUBMITTED, // Tribute given (hand changes)
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SELECTED,  // Tribute card selected from pool
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_RETURNED,  // Return tribute (hand changes)
		eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED,      // Tribute phase completed
		eventpb.EventType_EVENT_TYPE_TRICK_STARTED,          // New trick starts
		eventpb.EventType_EVENT_TYPE_PLAYER_PLAYED,          // Player played cards (hand changes)
		eventpb.EventType_EVENT_TYPE_PLAYER_PASSED,          // Player passed (turn changes)
		eventpb.EventType_EVENT_TYPE_TRICK_ENDED,            // Trick ends
		eventpb.EventType_EVENT_TYPE_DEAL_ENDED,             // Deal ends
		eventpb.EventType_EVENT_TYPE_MATCH_ENDED:            // Match ends
		wso.sendPlayerViews(event.Type)
	}

	// Log significant events
	switch event.Type {
	case eventpb.EventType_EVENT_TYPE_MATCH_STARTED, eventpb.EventType_EVENT_TYPE_MATCH_ENDED,
		eventpb.EventType_EVENT_TYPE_DEAL_STARTED, eventpb.EventType_EVENT_TYPE_DEAL_ENDED,
		eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED,
		eventpb.EventType_EVENT_TYPE_PLAYER_TIMEOUT:
		log.Printf("Game event %s for room %s", event.Type.String(), wso.roomID)
	}
}

// sendPlayerViews sends player-specific game state to each player
func (wso *WebSocketObserver) sendPlayerViews(eventType eventpb.EventType) {
	if wso.engine == nil {
		return
	}

	// Get game state once to extract player IDs
	gameState := wso.engine.GetGameState()
	if gameState == nil || gameState.CurrentMatch == nil {
		return
	}

	// Send player view to each player (seats 0-3)
	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		// Get player view for this seat
		playerView := wso.engine.GetPlayerView(playerSeat)
		if playerView == nil {
			continue
		}

		// Get player ID from the game state
		if playerSeat < len(gameState.CurrentMatch.Players) &&
			gameState.CurrentMatch.Players[playerSeat] != nil {

			playerID := gameState.CurrentMatch.Players[playerSeat].ID

			// Serialize PlayerView to JSON using protojson
			playerViewJSON := marshalProtoToRawJSON(playerView, "PlayerView")
			if playerViewJSON == nil {
				continue
			}

			// Create filtered player view message
			wsMessage := &websocket.WSMessage{
				Type: websocket.MSG_PLAYER_VIEW,
				Data: map[string]interface{}{
					"player_view": playerViewJSON,
					"event_type":  eventType.String(),
					"player_seat": playerSeat,
				},
				Timestamp: time.Now(),
				PlayerID:  playerID,
			}

			// Send to specific player
			if err := wso.wsManager.SendToPlayer(playerID, wsMessage); err != nil {
				log.Printf("Failed to send player view to player %s: %v", playerID, err)
			}

			// If in tribute phase, also send TributeView
			if eventType >= eventpb.EventType_EVENT_TYPE_TRIBUTE_STARTED &&
				eventType <= eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED {
				
				tributeView := wso.engine.GetTributeView(playerSeat)
				if tributeView != nil {
					// Serialize TributeView to JSON using protojson
					tributeViewJSON := marshalProtoToRawJSON(tributeView, "TributeView")
					if tributeViewJSON != nil {
						tributeMsg := &websocket.WSMessage{
							Type: "tribute_view",
							Data: map[string]interface{}{
								"tribute_view": tributeViewJSON,
								"event_type":   eventType.String(),
								"player_seat":  playerSeat,
							},
							Timestamp: time.Now(),
							PlayerID:  playerID,
						}

						if err := wso.wsManager.SendToPlayer(playerID, tributeMsg); err != nil {
							log.Printf("Failed to send tribute view to player %s: %v", playerID, err)
						}
					}
				}
			}
		}
	}
}
