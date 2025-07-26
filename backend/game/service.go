package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"guandan-world/backend/websocket"
	"guandan-world/sdk"
)

// WSManagerInterface defines the interface for WebSocket manager
type WSManagerInterface interface {
	BroadcastToRoom(roomID string, message *websocket.WSMessage)
	SendToPlayer(playerID string, message *websocket.WSMessage) error
}

// GameService provides pure coordination between SDK and WebSocket layer
// It contains NO game logic - all game rules and state management are handled by the SDK
type GameService struct {
	// Game engines by room ID
	engines map[string]sdk.GameEngineInterface
	
	// WebSocket manager for real-time communication
	wsManager WSManagerInterface
	
	// Synchronization
	mu sync.RWMutex
	
	// Timeout handling
	timeoutTicker *time.Ticker
	stopTimeout   chan struct{}
}

// NewGameService creates a new game coordination service
func NewGameService(wsManager WSManagerInterface) *GameService {
	service := &GameService{
		engines:     make(map[string]sdk.GameEngineInterface),
		wsManager:   wsManager,
		stopTimeout: make(chan struct{}),
	}
	
	// Start timeout processing
	service.startTimeoutProcessor()
	
	return service
}

// StartGame initializes a new game for the specified room
// This method creates a new SDK GameEngine and registers event handlers
func (gs *GameService) StartGame(roomID string, players []sdk.Player) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	// Validate input
	if roomID == "" {
		return fmt.Errorf("room ID is required")
	}
	
	if len(players) != 4 {
		return fmt.Errorf("exactly 4 players are required, got %d", len(players))
	}
	
	// Check if game already exists for this room
	if _, exists := gs.engines[roomID]; exists {
		return fmt.Errorf("game already exists for room %s", roomID)
	}
	
	// Create new game engine
	engine := sdk.NewGameEngine()
	
	// Register all event handlers for SDK events -> WebSocket message conversion
	gs.registerEventHandlers(engine, roomID)
	
	// Start the match in SDK
	if err := engine.StartMatch(players); err != nil {
		return fmt.Errorf("failed to start match in SDK: %w", err)
	}
	
	// Start the first deal
	if err := engine.StartDeal(); err != nil {
		return fmt.Errorf("failed to start deal in SDK: %w", err)
	}
	
	// Store engine
	gs.engines[roomID] = engine
	
	log.Printf("Game started for room %s with players: %v", roomID, players)
	return nil
}

// PlayCards delegates player card play to SDK without any validation or processing
func (gs *GameService) PlayCards(roomID string, playerSeat int, cardIDs []string) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Convert card IDs to Card objects (this is just data conversion, not game logic)
	cards, err := gs.convertCardIDs(cardIDs)
	if err != nil {
		return fmt.Errorf("failed to convert card IDs: %w", err)
	}
	
	// Delegate directly to SDK - no validation or processing here
	_, err = engine.PlayCards(playerSeat, cards)
	if err != nil {
		return fmt.Errorf("SDK rejected play: %w", err)
	}
	
	return nil
}

// PassTurn delegates player pass to SDK without any validation or processing
func (gs *GameService) PassTurn(roomID string, playerSeat int) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK - no validation or processing here
	_, err := engine.PassTurn(playerSeat)
	if err != nil {
		return fmt.Errorf("SDK rejected pass: %w", err)
	}
	
	return nil
}

// SubmitTributeSelection delegates tribute selection to SDK
func (gs *GameService) SubmitTributeSelection(roomID string, playerSeat int, cardID string) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	err := engine.SubmitTributeSelection(playerSeat, cardID)
	if err != nil {
		return fmt.Errorf("SDK rejected tribute selection: %w", err)
	}
	
	return nil
}

// SubmitReturnTribute delegates tribute return to SDK
func (gs *GameService) SubmitReturnTribute(roomID string, playerSeat int, cardID string) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	err := engine.SubmitReturnTribute(playerSeat, cardID)
	if err != nil {
		return fmt.Errorf("SDK rejected tribute return: %w", err)
	}
	
	return nil
}

// HandlePlayerDisconnect delegates player disconnection to SDK
func (gs *GameService) HandlePlayerDisconnect(roomID string, playerSeat int) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	_, err := engine.HandlePlayerDisconnect(playerSeat)
	if err != nil {
		return fmt.Errorf("SDK failed to handle disconnect: %w", err)
	}
	
	return nil
}

// HandlePlayerReconnect delegates player reconnection to SDK
func (gs *GameService) HandlePlayerReconnect(roomID string, playerSeat int) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	_, err := engine.HandlePlayerReconnect(playerSeat)
	if err != nil {
		return fmt.Errorf("SDK failed to handle reconnect: %w", err)
	}
	
	return nil
}

// GetPlayerView gets player-specific game state from SDK
func (gs *GameService) GetPlayerView(roomID string, playerSeat int) (*sdk.PlayerGameState, error) {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	return engine.GetPlayerView(playerSeat), nil
}

// GetGameState gets complete game state from SDK
func (gs *GameService) GetGameState(roomID string) (*sdk.GameState, error) {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Delegate directly to SDK
	return engine.GetGameState(), nil
}

// EndGame removes the game engine for a room
func (gs *GameService) EndGame(roomID string) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	if _, exists := gs.engines[roomID]; !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	delete(gs.engines, roomID)
	log.Printf("Game ended for room %s", roomID)
	return nil
}

// Stop stops the game service and cleanup resources
func (gs *GameService) Stop() {
	if gs.timeoutTicker != nil {
		gs.timeoutTicker.Stop()
	}
	
	close(gs.stopTimeout)
	
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	// Clear all engines
	gs.engines = make(map[string]sdk.GameEngineInterface)
}

// SyncGameState forces a complete state synchronization for a room
// This can be used when a player reconnects or needs a full state refresh
func (gs *GameService) SyncGameState(roomID string) error {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("no active game for room %s", roomID)
	}
	
	// Get complete game state
	gameState := engine.GetGameState()
	if gameState == nil {
		return fmt.Errorf("failed to get game state for room %s", roomID)
	}
	
	// Broadcast complete game state to all players
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_EVENT,
		Data: map[string]interface{}{
			"event_type": "game_state_sync",
			"game_state": gameState,
			"timestamp":  time.Now(),
		},
		Timestamp: time.Now(),
	}
	
	gs.wsManager.BroadcastToRoom(roomID, wsMessage)
	
	// Send individual player views
	gs.sendPlayerViews(roomID, "game_state_sync")
	
	log.Printf("Synchronized game state for room %s", roomID)
	return nil
}

// BroadcastGameEvent broadcasts a custom game event to all players in a room
// This provides a way to send additional events beyond the standard SDK events
func (gs *GameService) BroadcastGameEvent(roomID string, eventType string, eventData interface{}) error {
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_EVENT,
		Data: map[string]interface{}{
			"event_type": eventType,
			"event_data": eventData,
			"timestamp":  time.Now(),
		},
		Timestamp: time.Now(),
	}
	
	gs.wsManager.BroadcastToRoom(roomID, wsMessage)
	
	log.Printf("Broadcast custom event %s to room %s", eventType, roomID)
	return nil
}

// GetRealTimeGameStatus returns real-time status information for a room
// This provides quick access to current game status without full state
func (gs *GameService) GetRealTimeGameStatus(roomID string) (map[string]interface{}, error) {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no active game for room %s", roomID)
	}
	
	gameState := engine.GetGameState()
	if gameState == nil {
		return nil, fmt.Errorf("failed to get game state")
	}
	
	status := map[string]interface{}{
		"game_status": gameState.Status,
		"room_id": roomID,
		"timestamp": time.Now(),
	}
	
	if gameState.CurrentMatch != nil {
		match := gameState.CurrentMatch
		status["match_status"] = match.Status
		status["team_levels"] = match.TeamLevels
		
		if match.CurrentDeal != nil {
			deal := match.CurrentDeal
			status["deal_status"] = deal.Status
			status["deal_level"] = deal.Level
			
			if deal.CurrentTrick != nil {
				trick := deal.CurrentTrick
				status["current_turn"] = trick.CurrentTurn
				status["trick_leader"] = trick.Leader
				status["trick_status"] = trick.Status
			}
			
			if deal.TributePhase != nil {
				tribute := deal.TributePhase
				status["tribute_status"] = tribute.Status
				status["selecting_player"] = tribute.SelectingPlayer
			}
		}
	}
	
	return status, nil
}

// registerEventHandlers registers all SDK event handlers for WebSocket message conversion
func (gs *GameService) registerEventHandlers(engine sdk.GameEngineInterface, roomID string) {
	// Match events
	engine.RegisterEventHandler(sdk.EventMatchStarted, gs.createEventHandler(roomID, "match_started"))
	engine.RegisterEventHandler(sdk.EventMatchEnded, gs.createEventHandler(roomID, "match_ended"))
	
	// Deal events
	engine.RegisterEventHandler(sdk.EventDealStarted, gs.createEventHandler(roomID, "deal_started"))
	engine.RegisterEventHandler(sdk.EventCardsDealt, gs.createEventHandler(roomID, "cards_dealt"))
	engine.RegisterEventHandler(sdk.EventDealEnded, gs.createEventHandler(roomID, "deal_ended"))
	
	// Tribute events
	engine.RegisterEventHandler(sdk.EventTributePhase, gs.createEventHandler(roomID, "tribute_phase"))
	engine.RegisterEventHandler(sdk.EventTributeRulesSet, gs.createEventHandler(roomID, "tribute_rules_set"))
	engine.RegisterEventHandler(sdk.EventTributeImmunity, gs.createEventHandler(roomID, "tribute_immunity"))
	engine.RegisterEventHandler(sdk.EventTributePoolCreated, gs.createEventHandler(roomID, "tribute_pool_created"))
	engine.RegisterEventHandler(sdk.EventTributeStarted, gs.createEventHandler(roomID, "tribute_started"))
	engine.RegisterEventHandler(sdk.EventTributeGiven, gs.createEventHandler(roomID, "tribute_given"))
	engine.RegisterEventHandler(sdk.EventTributeSelected, gs.createEventHandler(roomID, "tribute_selected"))
	engine.RegisterEventHandler(sdk.EventReturnTribute, gs.createEventHandler(roomID, "return_tribute"))
	engine.RegisterEventHandler(sdk.EventTributeCompleted, gs.createEventHandler(roomID, "tribute_completed"))
	
	// Trick events
	engine.RegisterEventHandler(sdk.EventTrickStarted, gs.createEventHandler(roomID, "trick_started"))
	engine.RegisterEventHandler(sdk.EventPlayerPlayed, gs.createEventHandler(roomID, "player_played"))
	engine.RegisterEventHandler(sdk.EventPlayerPassed, gs.createEventHandler(roomID, "player_passed"))
	engine.RegisterEventHandler(sdk.EventTrickEnded, gs.createEventHandler(roomID, "trick_ended"))
	
	// Player events
	engine.RegisterEventHandler(sdk.EventPlayerTimeout, gs.createEventHandler(roomID, "player_timeout"))
	engine.RegisterEventHandler(sdk.EventPlayerDisconnect, gs.createEventHandler(roomID, "player_disconnect"))
	engine.RegisterEventHandler(sdk.EventPlayerReconnect, gs.createEventHandler(roomID, "player_reconnect"))
}

// createEventHandler creates a generic event handler that converts SDK events to WebSocket messages
func (gs *GameService) createEventHandler(roomID string, messageType string) sdk.GameEventHandler {
	return func(event *sdk.GameEvent) {
		// Convert SDK event to WebSocket message
		wsMessage := &websocket.WSMessage{
			Type:      websocket.MSG_GAME_EVENT,
			Data: map[string]interface{}{
				"event_type": messageType,
				"event_data": event.Data,
				"timestamp":  event.Timestamp,
				"player_seat": event.PlayerSeat,
			},
			Timestamp: event.Timestamp,
		}
		
		// Broadcast to all players in the room
		gs.wsManager.BroadcastToRoom(roomID, wsMessage)
		
		// Also send player-specific views for state-changing events
		// Use goroutine to avoid deadlock with SDK locks
		if gs.shouldSendPlayerViews(event.Type) {
			go gs.sendPlayerViews(roomID, event.Type)
		}
		
		log.Printf("Converted SDK event %s to WebSocket message for room %s", event.Type, roomID)
	}
}

// shouldSendPlayerViews determines if player-specific views should be sent for an event type
func (gs *GameService) shouldSendPlayerViews(eventType sdk.GameEventType) bool {
	switch eventType {
	case sdk.EventDealStarted,    // Cards are dealt when deal starts
		 sdk.EventCardsDealt,     // Explicit cards dealt event (if implemented)
		 sdk.EventTributeCompleted,
		 sdk.EventPlayerPlayed,
		 sdk.EventPlayerPassed,
		 sdk.EventTrickEnded,
		 sdk.EventDealEnded:
		return true
	default:
		return false
	}
}

// sendPlayerViews sends player-specific game state to each player
// This implements player view state filtering based on SDK GetPlayerView
func (gs *GameService) sendPlayerViews(roomID string, eventType sdk.GameEventType) {
	gs.mu.RLock()
	engine, exists := gs.engines[roomID]
	gs.mu.RUnlock()
	
	if !exists {
		return
	}
	
	// Send player view to each player (seats 0-3)
	for playerSeat := 0; playerSeat < 4; playerSeat++ {
		// Use SDK's GetPlayerView for proper state filtering
		playerView := engine.GetPlayerView(playerSeat)
		if playerView == nil {
			continue
		}
		
		// Get player ID from the game state
		if playerView.GameState != nil && 
		   playerView.GameState.CurrentMatch != nil && 
		   playerSeat < len(playerView.GameState.CurrentMatch.Players) &&
		   playerView.GameState.CurrentMatch.Players[playerSeat] != nil {
			
			playerID := playerView.GameState.CurrentMatch.Players[playerSeat].ID
			
			// Create filtered player view message
			wsMessage := &websocket.WSMessage{
				Type: websocket.MSG_PLAYER_VIEW,
				Data: map[string]interface{}{
					"player_view": playerView,
					"event_type":  eventType,
					"player_seat": playerSeat,
					"filtered_state": gs.createFilteredState(playerView, playerSeat),
				},
				Timestamp: time.Now(),
				PlayerID:  playerID,
			}
			
			// Send to specific player
			if err := gs.wsManager.SendToPlayer(playerID, wsMessage); err != nil {
				log.Printf("Failed to send player view to player %s: %v", playerID, err)
			}
		}
	}
}

// createFilteredState creates a filtered state object for a specific player
// This ensures each player only sees information they should have access to
func (gs *GameService) createFilteredState(playerView *sdk.PlayerGameState, playerSeat int) map[string]interface{} {
	filteredState := map[string]interface{}{
		"player_seat": playerSeat,
		"player_cards": playerView.PlayerCards, // Only this player's cards
		"visible_cards": playerView.VisibleCards, // Cards visible to all
	}
	
	// Add game state information that's visible to all players
	if playerView.GameState != nil {
		filteredState["game_status"] = playerView.GameState.Status
		
		if playerView.GameState.CurrentMatch != nil {
			match := playerView.GameState.CurrentMatch
			
			// Team levels are visible to all
			filteredState["team_levels"] = match.TeamLevels
			
			// Player names and seats are visible to all (but not their cards)
			players := make([]map[string]interface{}, 4)
			for i, player := range match.Players {
				if player != nil {
					players[i] = map[string]interface{}{
						"id": player.ID,
						"username": player.Username,
						"seat": player.Seat,
						"online": player.Online,
						"auto_play": player.AutoPlay,
					}
				}
			}
			filteredState["players"] = players
			
			// Current deal information (filtered)
			if match.CurrentDeal != nil {
				deal := match.CurrentDeal
				dealInfo := map[string]interface{}{
					"level": deal.Level,
					"status": deal.Status,
				}
				
				// Current trick information (visible to all)
				if deal.CurrentTrick != nil {
					trick := deal.CurrentTrick
					dealInfo["current_trick"] = map[string]interface{}{
						"leader": trick.Leader,
						"current_turn": trick.CurrentTurn,
						"status": trick.Status,
						"plays": trick.Plays, // All plays are visible
					}
				}
				
				// Tribute phase information (visible to all)
				if deal.TributePhase != nil {
					tribute := deal.TributePhase
					dealInfo["tribute_phase"] = map[string]interface{}{
						"status": tribute.Status,
						"tribute_map": tribute.TributeMap,
						"selecting_player": tribute.SelectingPlayer,
						"is_immune": tribute.IsImmune,
						// Pool cards are visible during selection
						"pool_cards": tribute.PoolCards,
					}
				}
				
				filteredState["current_deal"] = dealInfo
			}
		}
	}
	
	return filteredState
}

// startTimeoutProcessor starts the timeout processing goroutine
func (gs *GameService) startTimeoutProcessor() {
	gs.timeoutTicker = time.NewTicker(5 * time.Second) // Check timeouts every 5 seconds
	
	go func() {
		for {
			select {
			case <-gs.timeoutTicker.C:
				gs.processTimeouts()
			case <-gs.stopTimeout:
				return
			}
		}
	}()
}

// processTimeouts processes timeouts for all active games
// This implements the ProcessTimeouts requirement with real-time event handling
func (gs *GameService) processTimeouts() {
	gs.mu.RLock()
	engines := make(map[string]sdk.GameEngineInterface)
	for roomID, engine := range gs.engines {
		engines[roomID] = engine
	}
	gs.mu.RUnlock()
	
	// Process timeouts for each engine
	for roomID, engine := range engines {
		timeoutEvents := engine.ProcessTimeouts()
		
		// Handle each timeout event with real-time broadcasting
		for _, event := range timeoutEvents {
			gs.handleTimeoutEvent(roomID, event)
		}
		
		if len(timeoutEvents) > 0 {
			log.Printf("Processed %d timeout events for room %s", len(timeoutEvents), roomID)
		}
	}
}

// handleTimeoutEvent handles individual timeout events with real-time broadcasting
func (gs *GameService) handleTimeoutEvent(roomID string, event *sdk.GameEvent) {
	// Convert timeout event to WebSocket message
	wsMessage := &websocket.WSMessage{
		Type: websocket.MSG_GAME_EVENT,
		Data: map[string]interface{}{
			"event_type": "timeout_processed",
			"event_data": event.Data,
			"timestamp":  event.Timestamp,
			"player_seat": event.PlayerSeat,
			"original_event_type": string(event.Type),
		},
		Timestamp: event.Timestamp,
	}
	
	// Broadcast timeout event to all players in the room
	gs.wsManager.BroadcastToRoom(roomID, wsMessage)
	
	// Send updated player views for timeout events that change game state
	if gs.shouldSendPlayerViewsForTimeout(event.Type) {
		gs.sendPlayerViews(roomID, event.Type)
	}
	
	log.Printf("Handled timeout event %s for room %s, player seat %d", 
		event.Type, roomID, event.PlayerSeat)
}

// shouldSendPlayerViewsForTimeout determines if player views should be sent for timeout events
func (gs *GameService) shouldSendPlayerViewsForTimeout(eventType sdk.GameEventType) bool {
	switch eventType {
	case sdk.EventPlayerTimeout,
		 sdk.EventPlayerPlayed,  // Auto-play after timeout
		 sdk.EventPlayerPassed,  // Auto-pass after timeout
		 sdk.EventTrickEnded:    // Trick might end due to timeout
		return true
	default:
		return false
	}
}

// convertCardIDs converts string card IDs to Card objects
// This is pure data conversion, not game logic
func (gs *GameService) convertCardIDs(cardIDs []string) ([]*sdk.Card, error) {
	if len(cardIDs) == 0 {
		return nil, fmt.Errorf("no card IDs provided")
	}
	
	cards := make([]*sdk.Card, len(cardIDs))
	for i, cardID := range cardIDs {
		card, err := gs.parseCardFromID(cardID)
		if err != nil {
			return nil, fmt.Errorf("invalid card ID %s: %w", cardID, err)
		}
		cards[i] = card
	}
	
	return cards, nil
}

// parseCardFromID parses a card ID string into a Card object
// Card ID format: "Color_Number" (e.g., "Heart_5", "Joker_15")
func (gs *GameService) parseCardFromID(cardID string) (*sdk.Card, error) {
	if cardID == "" {
		return nil, fmt.Errorf("empty card ID")
	}
	
	// Split by underscore
	parts := make([]string, 0, 2)
	lastIndex := 0
	for i, char := range cardID {
		if char == '_' {
			if i > lastIndex {
				parts = append(parts, cardID[lastIndex:i])
			}
			lastIndex = i + 1
		}
	}
	if lastIndex < len(cardID) {
		parts = append(parts, cardID[lastIndex:])
	}
	
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid card ID format, expected 'Color_Number'")
	}
	
	color := parts[0]
	numberStr := parts[1]
	
	// Parse number
	var number int
	switch numberStr {
	case "1":
		number = 1
	case "2":
		number = 2
	case "3":
		number = 3
	case "4":
		number = 4
	case "5":
		number = 5
	case "6":
		number = 6
	case "7":
		number = 7
	case "8":
		number = 8
	case "9":
		number = 9
	case "10":
		number = 10
	case "11":
		number = 11
	case "12":
		number = 12
	case "13":
		number = 13
	case "14":
		number = 14
	case "15":
		number = 15
	case "16":
		number = 16
	default:
		return nil, fmt.Errorf("invalid card number: %s", numberStr)
	}
	
	// Create card with default level (will be set by game engine)
	card, err := sdk.NewCard(number, color, 2) // Default level 2
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	
	return card, nil
}