package sdk

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// GameEventType defines the types of game events
type GameEventType string

const (
	EventMatchStarted     GameEventType = "match_started"
	EventDealStarted      GameEventType = "deal_started"
	EventCardsDealt       GameEventType = "cards_dealt"
	EventTributePhase     GameEventType = "tribute_phase"
	EventTrickStarted     GameEventType = "trick_started"
	EventPlayerPlayed     GameEventType = "player_played"
	EventPlayerPassed     GameEventType = "player_passed"
	EventTrickEnded       GameEventType = "trick_ended"
	EventDealEnded        GameEventType = "deal_ended"
	EventMatchEnded       GameEventType = "match_ended"
	EventPlayerTimeout    GameEventType = "player_timeout"
	EventPlayerDisconnect GameEventType = "player_disconnect"
	EventPlayerReconnect  GameEventType = "player_reconnect"
)

// GameEvent represents a game event with associated data
type GameEvent struct {
	Type       GameEventType `json:"type"`
	Data       interface{}   `json:"data"`
	Timestamp  time.Time     `json:"timestamp"`
	PlayerSeat int           `json:"player_seat,omitempty"`
}

// GameEventHandler is a function type for handling game events
type GameEventHandler func(*GameEvent)

// GameState represents the complete state of the game
type GameState struct {
	ID           string      `json:"id"`
	Status       GameStatus  `json:"status"`
	CurrentMatch *Match      `json:"current_match,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// PlayerGameState represents the game state from a specific player's perspective
type PlayerGameState struct {
	PlayerSeat   int         `json:"player_seat"`
	GameState    *GameState  `json:"game_state"`
	PlayerCards  []*Card     `json:"player_cards"`
	VisibleCards []*Card     `json:"visible_cards"`
}

// GameStatus represents the current status of the game
type GameStatus string

const (
	GameStatusWaiting   GameStatus = "waiting"
	GameStatusStarted   GameStatus = "started"
	GameStatusFinished  GameStatus = "finished"
)

// GameEngine is the main game engine that manages the complete game lifecycle
type GameEngine struct {
	id            string
	status        GameStatus
	currentMatch  *Match
	eventHandlers map[GameEventType][]GameEventHandler
	mutex         sync.RWMutex
	createdAt     time.Time
	updatedAt     time.Time
}

// GameEngineInterface defines the public interface for the game engine
type GameEngineInterface interface {
	// Game lifecycle
	StartMatch(players []Player) error
	StartDeal() error
	
	// Game operations
	PlayCards(playerSeat int, cards []*Card) (*GameEvent, error)
	PassTurn(playerSeat int) (*GameEvent, error)
	SelectTribute(playerSeat int, card *Card) (*GameEvent, error)
	
	// State queries
	GetGameState() *GameState
	GetPlayerView(playerSeat int) *PlayerGameState
	IsGameFinished() bool
	
	// Event handling
	RegisterEventHandler(eventType GameEventType, handler GameEventHandler)
	ProcessTimeouts() []*GameEvent
	
	// Player management
	HandlePlayerDisconnect(playerSeat int) (*GameEvent, error)
	HandlePlayerReconnect(playerSeat int) (*GameEvent, error)
	SetPlayerAutoPlay(playerSeat int, enabled bool) error
}

// NewGameEngine creates a new game engine instance
func NewGameEngine() *GameEngine {
	now := time.Now()
	return &GameEngine{
		id:            generateID(),
		status:        GameStatusWaiting,
		eventHandlers: make(map[GameEventType][]GameEventHandler),
		createdAt:     now,
		updatedAt:     now,
	}
}

// StartMatch initializes a new match with the given players
func (ge *GameEngine) StartMatch(players []Player) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if len(players) != 4 {
		return errors.New("exactly 4 players are required")
	}
	
	if ge.status != GameStatusWaiting {
		return errors.New("game is not in waiting status")
	}
	
	// Create new match
	match, err := NewMatch(players)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}
	
	ge.currentMatch = match
	ge.status = GameStatusStarted
	ge.updatedAt = time.Now()
	
	// Emit match started event
	event := &GameEvent{
		Type:      EventMatchStarted,
		Data:      match,
		Timestamp: time.Now(),
	}
	ge.emitEvent(event)
	
	return nil
}

// StartDeal starts a new deal in the current match
func (ge *GameEngine) StartDeal() error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil {
		return errors.New("no active match")
	}
	
	err := ge.currentMatch.StartNewDeal()
	if err != nil {
		return fmt.Errorf("failed to start deal: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Emit deal started event
	event := &GameEvent{
		Type:      EventDealStarted,
		Data:      ge.currentMatch.CurrentDeal,
		Timestamp: time.Now(),
	}
	ge.emitEvent(event)
	
	return nil
}

// PlayCards handles a player's card play action
func (ge *GameEngine) PlayCards(playerSeat int, cards []*Card) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}
	
	err := ge.currentMatch.CurrentDeal.PlayCards(playerSeat, cards)
	if err != nil {
		return nil, fmt.Errorf("failed to play cards: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Create and emit player played event
	event := &GameEvent{
		Type:       EventPlayerPlayed,
		Data:       map[string]interface{}{
			"player_seat": playerSeat,
			"cards":       cards,
			"deal_state":  ge.currentMatch.CurrentDeal,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)
	
	return event, nil
}

// PassTurn handles a player's pass action
func (ge *GameEngine) PassTurn(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}
	
	err := ge.currentMatch.CurrentDeal.PassTurn(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to pass turn: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Create and emit player passed event
	event := &GameEvent{
		Type:       EventPlayerPassed,
		Data:       map[string]interface{}{
			"player_seat": playerSeat,
			"deal_state":  ge.currentMatch.CurrentDeal,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)
	
	return event, nil
}

// SelectTribute handles tribute selection during tribute phase
func (ge *GameEngine) SelectTribute(playerSeat int, card *Card) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil || ge.currentMatch.CurrentDeal == nil {
		return nil, errors.New("no active deal")
	}
	
	err := ge.currentMatch.CurrentDeal.SelectTribute(playerSeat, card)
	if err != nil {
		return nil, fmt.Errorf("failed to select tribute: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Create tribute selection event
	event := &GameEvent{
		Type:       EventTributePhase,
		Data:       map[string]interface{}{
			"player_seat": playerSeat,
			"card":        card,
			"tribute_phase": ge.currentMatch.CurrentDeal.TributePhase,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)
	
	return event, nil
}

// GetGameState returns the current complete game state
func (ge *GameEngine) GetGameState() *GameState {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()
	
	return &GameState{
		ID:           ge.id,
		Status:       ge.status,
		CurrentMatch: ge.currentMatch,
		CreatedAt:    ge.createdAt,
		UpdatedAt:    ge.updatedAt,
	}
}

// GetPlayerView returns the game state from a specific player's perspective
func (ge *GameEngine) GetPlayerView(playerSeat int) *PlayerGameState {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()
	
	gameState := ge.GetGameState()
	playerView := &PlayerGameState{
		PlayerSeat: playerSeat,
		GameState:  gameState,
	}
	
	// Add player-specific information if there's an active deal
	if ge.currentMatch != nil && ge.currentMatch.CurrentDeal != nil {
		if playerSeat >= 0 && playerSeat < 4 {
			playerView.PlayerCards = ge.currentMatch.CurrentDeal.PlayerCards[playerSeat]
		}
		
		// Add visible cards (cards played in current trick)
		if ge.currentMatch.CurrentDeal.CurrentTrick != nil {
			playerView.VisibleCards = ge.getVisibleCardsForPlayer(playerSeat)
		}
	}
	
	return playerView
}

// IsGameFinished checks if the game is finished
func (ge *GameEngine) IsGameFinished() bool {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()
	
	return ge.status == GameStatusFinished
}

// RegisterEventHandler registers an event handler for a specific event type
func (ge *GameEngine) RegisterEventHandler(eventType GameEventType, handler GameEventHandler) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.eventHandlers[eventType] == nil {
		ge.eventHandlers[eventType] = make([]GameEventHandler, 0)
	}
	ge.eventHandlers[eventType] = append(ge.eventHandlers[eventType], handler)
}

// ProcessTimeouts processes any pending timeouts and returns resulting events
func (ge *GameEngine) ProcessTimeouts() []*GameEvent {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	events := make([]*GameEvent, 0)
	
	if ge.currentMatch != nil && ge.currentMatch.CurrentDeal != nil {
		timeoutEvents := ge.currentMatch.CurrentDeal.ProcessTimeouts()
		events = append(events, timeoutEvents...)
		
		// Emit all timeout events
		for _, event := range timeoutEvents {
			ge.emitEvent(event)
		}
	}
	
	return events
}

// HandlePlayerDisconnect handles a player disconnection
func (ge *GameEngine) HandlePlayerDisconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil {
		return nil, errors.New("no active match")
	}
	
	err := ge.currentMatch.HandlePlayerDisconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle disconnect: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Create disconnect event
	event := &GameEvent{
		Type:       EventPlayerDisconnect,
		Data:       map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   true,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)
	
	return event, nil
}

// HandlePlayerReconnect handles a player reconnection
func (ge *GameEngine) HandlePlayerReconnect(playerSeat int) (*GameEvent, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil {
		return nil, errors.New("no active match")
	}
	
	err := ge.currentMatch.HandlePlayerReconnect(playerSeat)
	if err != nil {
		return nil, fmt.Errorf("failed to handle reconnect: %w", err)
	}
	
	ge.updatedAt = time.Now()
	
	// Create reconnect event
	event := &GameEvent{
		Type:       EventPlayerReconnect,
		Data:       map[string]interface{}{
			"player_seat": playerSeat,
			"auto_play":   false,
		},
		Timestamp:  time.Now(),
		PlayerSeat: playerSeat,
	}
	ge.emitEvent(event)
	
	return event, nil
}

// SetPlayerAutoPlay sets the auto-play status for a player
func (ge *GameEngine) SetPlayerAutoPlay(playerSeat int, enabled bool) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if ge.currentMatch == nil {
		return errors.New("no active match")
	}
	
	return ge.currentMatch.SetPlayerAutoPlay(playerSeat, enabled)
}

// emitEvent emits an event to all registered handlers
func (ge *GameEngine) emitEvent(event *GameEvent) {
	handlers, exists := ge.eventHandlers[event.Type]
	if !exists {
		return
	}
	
	// Call all handlers for this event type
	for _, handler := range handlers {
		go handler(event) // Run handlers in goroutines to avoid blocking
	}
}

// getVisibleCardsForPlayer returns the cards visible to a specific player
func (ge *GameEngine) getVisibleCardsForPlayer(playerSeat int) []*Card {
	visibleCards := make([]*Card, 0)
	
	if ge.currentMatch.CurrentDeal.CurrentTrick != nil {
		for _, play := range ge.currentMatch.CurrentDeal.CurrentTrick.Plays {
			if play.Cards != nil {
				visibleCards = append(visibleCards, play.Cards...)
			}
		}
	}
	
	return visibleCards
}

// generateID generates a unique ID for the game engine
func generateID() string {
	return fmt.Sprintf("game_%d", time.Now().UnixNano())
}