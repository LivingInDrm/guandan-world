package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"guandan-world/backend/auth"
	"guandan-world/backend/room"
)

// Message type constants
const (
	// Room management messages
	MSG_JOIN_ROOM  = "join_room"
	MSG_LEAVE_ROOM = "leave_room"
	MSG_START_GAME = "start_game"
	
	// Game operation messages
	MSG_PLAY_CARDS = "play_cards"
	MSG_PASS       = "pass"
	MSG_TRIBUTE_SELECT = "tribute_select"
	MSG_TRIBUTE_RETURN = "tribute_return"
	
	// Status and notification messages
	MSG_ROOM_UPDATE = "room_update"
	MSG_GAME_EVENT  = "game_event"
	MSG_PLAYER_VIEW = "player_view"
	MSG_GAME_ACTION = "game_action"
	MSG_ERROR       = "error"
	MSG_PING        = "ping"
	MSG_PONG        = "pong"
)

// JoinRoomData represents the data for joining a room
type JoinRoomData struct {
	RoomID string `json:"room_id"`
}

// LeaveRoomData represents the data for leaving a room
type LeaveRoomData struct {
	RoomID string `json:"room_id"`
}

// StartGameData represents the data for starting a game
type StartGameData struct {
	RoomID string `json:"room_id"`
}

// PlayCardsData represents the data for playing cards
type PlayCardsData struct {
	Cards []string `json:"cards"` // Card IDs
}

// TributeSelectData represents the data for tribute selection
type TributeSelectData struct {
	CardID string `json:"card_id"`
}

// TributeReturnData represents the data for tribute return
type TributeReturnData struct {
	CardID string `json:"card_id"`
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	PlayerID  string      `json:"player_id,omitempty"`
}

// WSConnection represents a WebSocket connection with associated metadata
type WSConnection struct {
	conn     *websocket.Conn
	playerID string
	roomID   string
	send     chan []byte
	manager  *WSManager
	lastPing time.Time
	mu       sync.RWMutex
}

// WSManager manages WebSocket connections and message routing
type WSManager struct {
	// Connection management
	connections map[string]*WSConnection // playerID -> connection
	rooms       map[string]map[string]*WSConnection // roomID -> playerID -> connection
	
	// Services
	authService auth.AuthService
	roomService room.RoomService
	
	// Channels for connection lifecycle
	register   chan *WSConnection
	unregister chan *WSConnection
	broadcast  chan *BroadcastMessage
	
	// Configuration
	upgrader websocket.Upgrader
	
	// Synchronization
	mu sync.RWMutex
	
	// Heartbeat configuration
	pingInterval time.Duration
	pongTimeout  time.Duration
	
	// Message handlers
	messageHandlers map[string]MessageHandler
}

// BroadcastMessage represents a message to be broadcast to a room
type BroadcastMessage struct {
	RoomID  string
	Message *WSMessage
	Exclude string // PlayerID to exclude from broadcast
}

// MessageHandler defines the interface for handling WebSocket messages
type MessageHandler func(conn *WSConnection, message *WSMessage) error

// NewWSManager creates a new WebSocket manager
func NewWSManager(authService auth.AuthService, roomService room.RoomService) *WSManager {
	manager := &WSManager{
		connections:     make(map[string]*WSConnection),
		rooms:          make(map[string]map[string]*WSConnection),
		authService:    authService,
		roomService:    roomService,
		register:       make(chan *WSConnection, 256),
		unregister:     make(chan *WSConnection, 256),
		broadcast:      make(chan *BroadcastMessage, 256),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement proper origin checking for production
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		pingInterval:    30 * time.Second,
		pongTimeout:     60 * time.Second,
		messageHandlers: make(map[string]MessageHandler),
	}
	
	// Register default message handlers
	manager.registerDefaultHandlers()
	
	return manager
}

// registerDefaultHandlers registers the default message handlers
func (m *WSManager) registerDefaultHandlers() {
	m.messageHandlers[MSG_PING] = m.handlePing
	m.messageHandlers[MSG_JOIN_ROOM] = m.handleJoinRoom
	m.messageHandlers[MSG_LEAVE_ROOM] = m.handleLeaveRoom
	m.messageHandlers[MSG_START_GAME] = m.handleStartGame
	m.messageHandlers[MSG_PLAY_CARDS] = m.handlePlayCards
	m.messageHandlers[MSG_PASS] = m.handlePass
	m.messageHandlers[MSG_TRIBUTE_SELECT] = m.handleTributeSelect
	m.messageHandlers[MSG_TRIBUTE_RETURN] = m.handleTributeReturn
}

// RegisterHandler registers a custom message handler
func (m *WSManager) RegisterHandler(messageType string, handler MessageHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messageHandlers[messageType] = handler
}

// Run starts the WebSocket manager's main loop
func (m *WSManager) Run() {
	// Start heartbeat ticker
	heartbeatTicker := time.NewTicker(m.pingInterval)
	defer heartbeatTicker.Stop()
	
	for {
		select {
		case conn := <-m.register:
			m.handleRegister(conn)
			
		case conn := <-m.unregister:
			m.handleUnregister(conn)
			
		case broadcastMsg := <-m.broadcast:
			m.handleBroadcast(broadcastMsg)
			
		case <-heartbeatTicker.C:
			m.handleHeartbeat()
		}
	}
}

// HandleWebSocket handles WebSocket upgrade and connection
func (m *WSManager) HandleWebSocket(w http.ResponseWriter, r *http.Request, playerID string) error {
	// Validate player authentication
	user, err := m.authService.GetUserByID(playerID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return fmt.Errorf("invalid player: %w", err)
	}
	
	// Upgrade connection
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}
	
	// Create WSConnection
	wsConn := &WSConnection{
		conn:     conn,
		playerID: user.ID,
		send:     make(chan []byte, 256),
		manager:  m,
		lastPing: time.Now(),
	}
	
	// Register connection
	m.register <- wsConn
	
	// Start goroutines for reading and writing
	go wsConn.readPump()
	go wsConn.writePump()
	
	return nil
}

// handleRegister handles new connection registration
func (m *WSManager) handleRegister(conn *WSConnection) {
	if conn == nil {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Close existing connection if any
	if existingConn, exists := m.connections[conn.playerID]; exists {
		// Close send channel safely
		select {
		case <-existingConn.send:
			// Channel already closed
		default:
			close(existingConn.send)
		}
		
		if existingConn.conn != nil {
			existingConn.conn.Close()
		}
		
		// Remove from room if in one
		if existingConn.roomID != "" {
			if roomConns, exists := m.rooms[existingConn.roomID]; exists {
				delete(roomConns, conn.playerID)
				if len(roomConns) == 0 {
					delete(m.rooms, existingConn.roomID)
				}
			}
		}
	}
	
	// Register new connection
	m.connections[conn.playerID] = conn
	
	// Check if player is in a room and add to room connections
	if playerRoom, err := m.roomService.GetPlayerRoom(conn.playerID); err == nil {
		conn.roomID = playerRoom.ID
		if _, exists := m.rooms[playerRoom.ID]; !exists {
			m.rooms[playerRoom.ID] = make(map[string]*WSConnection)
		}
		m.rooms[playerRoom.ID][conn.playerID] = conn
	}
	
	log.Printf("Player %s connected via WebSocket", conn.playerID)
}

// handleUnregister handles connection disconnection
func (m *WSManager) handleUnregister(conn *WSConnection) {
	if conn == nil {
		return
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Remove from connections
	if existingConn, exists := m.connections[conn.playerID]; exists && existingConn == conn {
		delete(m.connections, conn.playerID)
		// Close send channel safely
		select {
		case <-conn.send:
			// Channel already closed
		default:
			close(conn.send)
		}
	}
	
	// Remove from room connections
	if conn.roomID != "" {
		if roomConns, exists := m.rooms[conn.roomID]; exists {
			delete(roomConns, conn.playerID)
			if len(roomConns) == 0 {
				delete(m.rooms, conn.roomID)
			}
		}
		
		// Notify room service about disconnection
		// This will trigger auto-play for the disconnected player
		m.notifyPlayerDisconnected(conn.playerID, conn.roomID)
	}
	
	log.Printf("Player %s disconnected from WebSocket", conn.playerID)
}

// handleBroadcast handles broadcasting messages to room members
func (m *WSManager) handleBroadcast(broadcastMsg *BroadcastMessage) {
	if broadcastMsg == nil {
		return
	}
	
	m.mu.RLock()
	roomConns, exists := m.rooms[broadcastMsg.RoomID]
	if !exists {
		m.mu.RUnlock()
		return
	}
	
	// Create a copy of connections to avoid holding lock during send
	connections := make([]*WSConnection, 0, len(roomConns))
	for playerID, conn := range roomConns {
		if playerID != broadcastMsg.Exclude {
			connections = append(connections, conn)
		}
	}
	m.mu.RUnlock()
	
	// Serialize message
	messageData, err := json.Marshal(broadcastMsg.Message)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}
	
	// Send to all connections
	for _, conn := range connections {
		if conn != nil && conn.send != nil {
			select {
			case conn.send <- messageData:
			default:
				// Connection is blocked, close it
				select {
				case m.unregister <- conn:
				default:
					// Unregister channel might be full, ignore
				}
			}
		}
	}
}

// handleHeartbeat handles periodic heartbeat checks
func (m *WSManager) handleHeartbeat() {
	m.mu.RLock()
	connections := make([]*WSConnection, 0, len(m.connections))
	for _, conn := range m.connections {
		connections = append(connections, conn)
	}
	m.mu.RUnlock()
	
	now := time.Now()
	for _, conn := range connections {
		conn.mu.RLock()
		lastPing := conn.lastPing
		conn.mu.RUnlock()
		
		if now.Sub(lastPing) > m.pongTimeout {
			// Connection is stale, close it
			log.Printf("Closing stale connection for player %s", conn.playerID)
			m.unregister <- conn
		} else {
			// Send ping
			conn.sendPing()
		}
	}
}

// BroadcastToRoom broadcasts a message to all connections in a room
func (m *WSManager) BroadcastToRoom(roomID string, message *WSMessage) {
	if message == nil || roomID == "" {
		return
	}
	
	broadcastMsg := &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
	
	select {
	case m.broadcast <- broadcastMsg:
	default:
		log.Printf("Broadcast channel full, dropping message for room %s", roomID)
	}
}

// BroadcastToRoomExcept broadcasts a message to all connections in a room except one player
func (m *WSManager) BroadcastToRoomExcept(roomID string, message *WSMessage, excludePlayerID string) {
	if message == nil || roomID == "" {
		return
	}
	
	broadcastMsg := &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
		Exclude: excludePlayerID,
	}
	
	select {
	case m.broadcast <- broadcastMsg:
	default:
		log.Printf("Broadcast channel full, dropping message for room %s", roomID)
	}
}

// SendToPlayer sends a message to a specific player
func (m *WSManager) SendToPlayer(playerID string, message *WSMessage) error {
	m.mu.RLock()
	conn, exists := m.connections[playerID]
	m.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("player %s not connected", playerID)
	}
	
	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	select {
	case conn.send <- messageData:
		return nil
	default:
		return fmt.Errorf("connection send channel full for player %s", playerID)
	}
}

// GetRoomConnections returns the number of active connections in a room
func (m *WSManager) GetRoomConnections(roomID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if roomConns, exists := m.rooms[roomID]; exists {
		return len(roomConns)
	}
	return 0
}

// IsPlayerConnected checks if a player is currently connected
func (m *WSManager) IsPlayerConnected(playerID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, exists := m.connections[playerID]
	return exists
}

// notifyPlayerDisconnected notifies about player disconnection
func (m *WSManager) notifyPlayerDisconnected(playerID, roomID string) {
	// This would typically trigger auto-play or other disconnection handling
	// For now, we'll just log it
	log.Printf("Player %s disconnected from room %s", playerID, roomID)
	
	// TODO: Integrate with game service to handle player disconnection
	// This should trigger the SDK's HandlePlayerDisconnect method
}

// handlePing handles ping messages for heartbeat
func (m *WSManager) handlePing(conn *WSConnection, message *WSMessage) error {
	// Update last ping time
	conn.mu.Lock()
	conn.lastPing = time.Now()
	conn.mu.Unlock()
	
	// Send pong response
	pongMsg := &WSMessage{
		Type:      MSG_PONG,
		Data:      nil,
		Timestamp: time.Now(),
	}
	
	messageData, err := json.Marshal(pongMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal pong message: %w", err)
	}
	
	select {
	case conn.send <- messageData:
		return nil
	default:
		return fmt.Errorf("failed to send pong message")
	}
}

// handleJoinRoom handles room joining requests
func (m *WSManager) handleJoinRoom(conn *WSConnection, message *WSMessage) error {
	// Parse join room data
	var joinData JoinRoomData
	if err := parseMessageData(message.Data, &joinData); err != nil {
		return fmt.Errorf("invalid join room data: %w", err)
	}
	
	// Validate room ID
	if joinData.RoomID == "" {
		return fmt.Errorf("room ID is required")
	}
	
	// Join room through room service
	room, err := m.roomService.JoinRoom(joinData.RoomID, conn.playerID)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}
	
	// Update connection room ID
	m.mu.Lock()
	oldRoomID := conn.roomID
	conn.roomID = joinData.RoomID
	
	// Remove from old room if any
	if oldRoomID != "" && oldRoomID != joinData.RoomID {
		if oldRoomConns, exists := m.rooms[oldRoomID]; exists {
			delete(oldRoomConns, conn.playerID)
			if len(oldRoomConns) == 0 {
				delete(m.rooms, oldRoomID)
			}
		}
	}
	
	// Add to new room
	if _, exists := m.rooms[joinData.RoomID]; !exists {
		m.rooms[joinData.RoomID] = make(map[string]*WSConnection)
	}
	m.rooms[joinData.RoomID][conn.playerID] = conn
	m.mu.Unlock()
	
	// Broadcast room update to all room members
	roomUpdateMsg := &WSMessage{
		Type: MSG_ROOM_UPDATE,
		Data: map[string]interface{}{
			"action": "player_joined",
			"room":   room,
			"player_id": conn.playerID,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoom(joinData.RoomID, roomUpdateMsg)
	
	return nil
}

// handleLeaveRoom handles room leaving requests
func (m *WSManager) handleLeaveRoom(conn *WSConnection, message *WSMessage) error {
	// Parse leave room data
	var leaveData LeaveRoomData
	if err := parseMessageData(message.Data, &leaveData); err != nil {
		return fmt.Errorf("invalid leave room data: %w", err)
	}
	
	// Validate room ID
	if leaveData.RoomID == "" {
		return fmt.Errorf("room ID is required")
	}
	
	// Verify player is in the specified room
	if conn.roomID != leaveData.RoomID {
		return fmt.Errorf("player is not in the specified room")
	}
	
	// Leave room through room service
	room, err := m.roomService.LeaveRoom(leaveData.RoomID, conn.playerID)
	if err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}
	
	// Update connection room ID
	m.mu.Lock()
	conn.roomID = ""
	
	// Remove from room connections
	if roomConns, exists := m.rooms[leaveData.RoomID]; exists {
		delete(roomConns, conn.playerID)
		if len(roomConns) == 0 {
			delete(m.rooms, leaveData.RoomID)
		}
	}
	m.mu.Unlock()
	
	// Broadcast room update if room still exists
	if room != nil {
		roomUpdateMsg := &WSMessage{
			Type: MSG_ROOM_UPDATE,
			Data: map[string]interface{}{
				"action": "player_left",
				"room":   room,
				"player_id": conn.playerID,
			},
			Timestamp: time.Now(),
		}
		
		m.BroadcastToRoom(leaveData.RoomID, roomUpdateMsg)
	}
	
	return nil
}

// handleStartGame handles game start requests
func (m *WSManager) handleStartGame(conn *WSConnection, message *WSMessage) error {
	// Parse start game data
	var startData StartGameData
	if err := parseMessageData(message.Data, &startData); err != nil {
		return fmt.Errorf("invalid start game data: %w", err)
	}
	
	// Validate room ID
	if startData.RoomID == "" {
		return fmt.Errorf("room ID is required")
	}
	
	// Verify player is in the specified room
	if conn.roomID != startData.RoomID {
		return fmt.Errorf("player is not in the specified room")
	}
	
	// Start game through room service
	err := m.roomService.StartGame(startData.RoomID, conn.playerID)
	if err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	
	// Get updated room info
	room, err := m.roomService.GetRoom(startData.RoomID)
	if err != nil {
		return fmt.Errorf("failed to get room info: %w", err)
	}
	
	// Broadcast game start to all room members
	gameStartMsg := &WSMessage{
		Type: MSG_ROOM_UPDATE,
		Data: map[string]interface{}{
			"action": "game_started",
			"room":   room,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoom(startData.RoomID, gameStartMsg)
	
	// TODO: Initialize game engine and start the actual game
	// This will be implemented in task 9 (游戏服务)
	
	return nil
}

// handlePlayCards handles card playing requests
func (m *WSManager) handlePlayCards(conn *WSConnection, message *WSMessage) error {
	// Parse play cards data
	var playData PlayCardsData
	if err := parseMessageData(message.Data, &playData); err != nil {
		return fmt.Errorf("invalid play cards data: %w", err)
	}
	
	// Validate player is in a room
	if conn.roomID == "" {
		return fmt.Errorf("player is not in a room")
	}
	
	// Validate cards
	if len(playData.Cards) == 0 {
		return fmt.Errorf("no cards specified")
	}
	
	// TODO: Forward to game service for processing
	// This will be implemented in task 9 (游戏服务)
	
	// For now, just acknowledge the request
	ackMsg := &WSMessage{
		Type: "play_cards_ack",
		Data: map[string]interface{}{
			"player_id": conn.playerID,
			"cards":     playData.Cards,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoomExcept(conn.roomID, ackMsg, conn.playerID)
	
	return nil
}

// handlePass handles pass requests
func (m *WSManager) handlePass(conn *WSConnection, message *WSMessage) error {
	// Validate player is in a room
	if conn.roomID == "" {
		return fmt.Errorf("player is not in a room")
	}
	
	// TODO: Forward to game service for processing
	// This will be implemented in task 9 (游戏服务)
	
	// For now, just acknowledge the request
	passMsg := &WSMessage{
		Type: "pass_ack",
		Data: map[string]interface{}{
			"player_id": conn.playerID,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoomExcept(conn.roomID, passMsg, conn.playerID)
	
	return nil
}

// handleTributeSelect handles tribute selection requests
func (m *WSManager) handleTributeSelect(conn *WSConnection, message *WSMessage) error {
	// Parse tribute select data
	var selectData TributeSelectData
	if err := parseMessageData(message.Data, &selectData); err != nil {
		return fmt.Errorf("invalid tribute select data: %w", err)
	}
	
	// Validate player is in a room
	if conn.roomID == "" {
		return fmt.Errorf("player is not in a room")
	}
	
	// Validate card ID
	if selectData.CardID == "" {
		return fmt.Errorf("card ID is required")
	}
	
	// TODO: Forward to game service for processing
	// This will be implemented in task 9 (游戏服务)
	
	// For now, just acknowledge the request
	selectMsg := &WSMessage{
		Type: "tribute_select_ack",
		Data: map[string]interface{}{
			"player_id": conn.playerID,
			"card_id":   selectData.CardID,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoomExcept(conn.roomID, selectMsg, conn.playerID)
	
	return nil
}

// handleTributeReturn handles tribute return requests
func (m *WSManager) handleTributeReturn(conn *WSConnection, message *WSMessage) error {
	// Parse tribute return data
	var returnData TributeReturnData
	if err := parseMessageData(message.Data, &returnData); err != nil {
		return fmt.Errorf("invalid tribute return data: %w", err)
	}
	
	// Validate player is in a room
	if conn.roomID == "" {
		return fmt.Errorf("player is not in a room")
	}
	
	// Validate card ID
	if returnData.CardID == "" {
		return fmt.Errorf("card ID is required")
	}
	
	// TODO: Forward to game service for processing
	// This will be implemented in task 9 (游戏服务)
	
	// For now, just acknowledge the request
	returnMsg := &WSMessage{
		Type: "tribute_return_ack",
		Data: map[string]interface{}{
			"player_id": conn.playerID,
			"card_id":   returnData.CardID,
		},
		Timestamp: time.Now(),
	}
	
	m.BroadcastToRoomExcept(conn.roomID, returnMsg, conn.playerID)
	
	return nil
}

// readPump handles reading messages from the WebSocket connection
func (c *WSConnection) readPump() {
	defer func() {
		select {
		case c.manager.unregister <- c:
		default:
			// Channel might be closed, ignore
		}
		if c.conn != nil {
			c.conn.Close()
		}
	}()
	
	// Set read deadline and pong handler
	c.conn.SetReadDeadline(time.Now().Add(c.manager.pongTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastPing = time.Now()
		c.mu.Unlock()
		c.conn.SetReadDeadline(time.Now().Add(c.manager.pongTimeout))
		return nil
	})
	
	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for player %s: %v", c.playerID, err)
			}
			break
		}
		
		// Parse message
		var message WSMessage
		if err := json.Unmarshal(messageData, &message); err != nil {
			log.Printf("Failed to unmarshal message from player %s: %v", c.playerID, err)
			continue
		}
		
		// Set player ID and timestamp
		message.PlayerID = c.playerID
		message.Timestamp = time.Now()
		
		// Handle message
		c.handleMessage(&message)
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *WSConnection) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to write message to player %s: %v", c.playerID, err)
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sendPing sends a ping message to the connection
func (c *WSConnection) sendPing() {
	select {
	case c.send <- []byte(`{"type":"ping","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`):
	default:
		// Channel full, connection will be closed by heartbeat
	}
}

// handleMessage routes messages to appropriate handlers
func (c *WSConnection) handleMessage(message *WSMessage) {
	c.manager.mu.RLock()
	handler, exists := c.manager.messageHandlers[message.Type]
	c.manager.mu.RUnlock()
	
	if !exists {
		log.Printf("Unknown message type '%s' from player %s", message.Type, c.playerID)
		return
	}
	
	if err := handler(c, message); err != nil {
		log.Printf("Error handling message type '%s' from player %s: %v", message.Type, c.playerID, err)
		
		// Send error response
		errorMsg := &WSMessage{
			Type: "error",
			Data: map[string]interface{}{
				"message": err.Error(),
				"original_type": message.Type,
			},
			Timestamp: time.Now(),
		}
		
		if messageData, marshalErr := json.Marshal(errorMsg); marshalErr == nil {
			select {
			case c.send <- messageData:
			default:
				// Channel full, ignore
			}
		}
	}
}

// parseMessageData parses message data into the specified struct
func parseMessageData(data interface{}, target interface{}) error {
	// Convert to JSON and back to parse into target struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	
	return nil
}