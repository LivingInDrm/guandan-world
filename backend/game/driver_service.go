package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"guandan-world/backend/websocket"
	eventpb "guandan-world/proto/event"
	"guandan-world/sdk"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// protoJSONMarshaler is the configuration for serializing proto messages to JSON
// Uses snake_case field names and emits all fields including zero values
// EmitUnpopulated must be true to ensure player_seat=0 is serialized (seat 0 is valid)
var protoJSONMarshaler = protojson.MarshalOptions{
	UseProtoNames:   true, // Use snake_case field names (match_id, deal_status)
	EmitUnpopulated: true, // Emit all fields, including zero values (required for player_seat=0)
}

// marshalProtoToRawJSON serializes a proto message to json.RawMessage
// using the configured protoJSONMarshaler (snake_case fields, no zero values)
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
	provider := NewRoomInputProvider(roomID, ds.wsManager)
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
	// Defensive check: context must not be nil
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}
	
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
	wsData := map[string]interface{}{
		"action_type": "play_decision_required",
		"player_seat": playerSeat,
		"hand":        hand,
		"trick_info":  trickInfo,
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		wsData["timeout"] = secs
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      wsData,
		Timestamp: time.Now(),
	}

	// Get player ID and send message
	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send play request: %w", err)
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

	// Send request to player
	wsData := map[string]interface{}{
		"action_type": "tribute_selection_required",
		"player_seat": playerSeat,
		"options":     options,
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		wsData["timeout"] = secs
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      wsData,
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send tribute selection request: %w", err)
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

	// Send request to player
	wsData := map[string]interface{}{
		"action_type": "return_tribute_required",
		"player_seat": playerSeat,
		"hand":        hand,
	}
	
	// Include timeout if context has a deadline
	if secs, ok := getRemainingTimeout(ctx); ok {
		wsData["timeout"] = secs
	}
	
	wsMessage := &websocket.WSMessage{
		Type:      websocket.MSG_GAME_ACTION,
		Data:      wsData,
		Timestamp: time.Now(),
	}

	if err := rip.sendToPlayer(playerSeat, wsMessage); err != nil {
		return nil, fmt.Errorf("failed to send return tribute request: %w", err)
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
}

// sendToPlayer sends a message to a specific player
// SECURITY WARNING: Current implementation broadcasts to entire room, exposing private game state.
// This is a TEMPORARY development-only approach and MUST NOT ship to production.
// TODO: Implement targeted messaging using WSManagerInterface.SendToPlayer with actual player ID.
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
	engine    sdk.GameEngineInterface // Reference to engine for accessing player views
}

// NewWebSocketObserver creates a new WebSocket observer
func NewWebSocketObserver(roomID string, wsManager WSManagerInterface, engine sdk.GameEngineInterface) *WebSocketObserver {
	return &WebSocketObserver{
		roomID:    roomID,
		wsManager: wsManager,
		engine:    engine,
	}
}

// OnGameEvent implements sdk.EventObserver
func (wso *WebSocketObserver) OnGameEvent(event *sdk.GameEvent) {
	// Serialize the entire event to JSON using protojson
	eventJSON := marshalProtoToRawJSON(event, "GameEvent")
	if eventJSON == nil {
		return
	}

	// Convert SDK event to WebSocket message
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_EVENT,
		Data: map[string]interface{}{
			"event_type":  event.Type.String(),
			"event_data":  eventJSON,
			"timestamp":   time.UnixMilli(event.CreatedAtMs),
			"player_seat": event.GetActorSeat(),
		},
		Timestamp: time.UnixMilli(event.CreatedAtMs),
	}

	// Broadcast to all players in the room
	wso.wsManager.BroadcastToRoom(wso.roomID, wsMessage)

	// Send player-specific views only for key events
	// These are events where hand cards change or game phase transitions occur
	switch event.Type {
	case eventpb.EventType_EVENT_TYPE_MATCH_STARTED, // Match begins
		eventpb.EventType_EVENT_TYPE_DEAL_STARTED,                // New deal starts
		eventpb.EventType_EVENT_TYPE_CARDS_DEALT,                 // Cards dealt to players
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_SUBMITTED,      // Tribute given (hand changes)
		eventpb.EventType_EVENT_TYPE_TRIBUTE_CARD_RETURNED,       // Return tribute (hand changes)
		eventpb.EventType_EVENT_TYPE_TRIBUTE_COMPLETED,           // Tribute phase completed
		eventpb.EventType_EVENT_TYPE_TRICK_STARTED,               // New trick starts
		eventpb.EventType_EVENT_TYPE_PLAYER_PLAYED,               // Player played cards (hand changes)
		eventpb.EventType_EVENT_TYPE_PLAYER_PASSED,               // Player passed (turn changes)
		eventpb.EventType_EVENT_TYPE_TRICK_ENDED,                 // Trick ends
		eventpb.EventType_EVENT_TYPE_DEAL_ENDED,                  // Deal ends
		eventpb.EventType_EVENT_TYPE_MATCH_ENDED:                 // Match ends
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
