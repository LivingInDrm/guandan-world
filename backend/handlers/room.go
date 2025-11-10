package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/room"
	"guandan-world/backend/websocket"
	"guandan-world/sdk"

	"github.com/gin-gonic/gin"
)

// DriverServiceInterface defines the interface for game driver service
type DriverServiceInterface interface {
	StartGameWithDriver(roomID string, players []sdk.Player) error
}

// WSManagerInterface defines the interface for WebSocket manager
type WSManagerInterface interface {
	BroadcastToRoom(roomID string, message *websocket.WSMessage)
	SendToPlayer(playerID string, message *websocket.WSMessage) error
}

// RoomHandler handles room-related HTTP requests
type RoomHandler struct {
	roomService   room.RoomService
	authService   auth.AuthService
	driverService DriverServiceInterface
	wsManager     WSManagerInterface
}

// NewRoomHandler creates a new room handler
func NewRoomHandler(
	roomService room.RoomService,
	authService auth.AuthService,
	driverService DriverServiceInterface,
	wsManager WSManagerInterface,
) *RoomHandler {
	return &RoomHandler{
		roomService:   roomService,
		authService:   authService,
		driverService: driverService,
		wsManager:     wsManager,
	}
}

// CreateRoomRequest represents a room creation request
type CreateRoomRequest struct {
	// Room creation doesn't need additional parameters for now
	// The owner is determined from the authenticated user
}

// JoinRoomRequest represents a room join request
type JoinRoomRequest struct {
	// Room ID is in the URL path, player ID comes from auth context
}

// RoomResponse represents a room response
type RoomResponse struct {
	Room *room.Room `json:"room"`
}

// RoomListResponse represents a room list response
type RoomListResponse struct {
	*room.RoomListResponse
}

// CreateRoom handles room creation
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID in context",
		})
		return
	}

	// Create room
	newRoom, err := h.roomService.CreateRoom(userIDStr)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "room_creation_failed"

		// Handle specific error cases
		if strings.Contains(err.Error(), "player is already in room") {
			statusCode = http.StatusConflict
			errorCode = "already_in_room"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, RoomResponse{
		Room: newRoom,
	})
}

// GetRooms handles room list retrieval
func (h *RoomHandler) GetRooms(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	limit := 12 // Default as specified in requirements

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Parse status filter
	var statusFilter *room.RoomStatus
	if statusStr := c.Query("status"); statusStr != "" {
		switch statusStr {
		case "waiting":
			status := room.RoomStatusWaiting
			statusFilter = &status
		case "ready":
			status := room.RoomStatusReady
			statusFilter = &status
		case "playing":
			status := room.RoomStatusPlaying
			statusFilter = &status
		}
	}

	// Get room list
	response, err := h.roomService.GetRoomList(page, limit, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "room_list_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RoomListResponse{
		RoomListResponse: response,
	})
}

// GetRoom handles single room retrieval
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Room ID is required",
		})
		return
	}

	// Get room
	roomData, err := h.roomService.GetRoom(roomID)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() != "room not found" {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "room_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RoomResponse{
		Room: roomData,
	})
}

// JoinRoom handles room joining
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Room ID is required",
		})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID in context",
		})
		return
	}

	// Join room
	updatedRoom, err := h.roomService.JoinRoom(roomID, userIDStr)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "join_room_failed"

		// Handle specific error cases
		switch err.Error() {
		case "room not found":
			statusCode = http.StatusNotFound
			errorCode = "room_not_found"
		case "room is full":
			statusCode = http.StatusConflict
			errorCode = "room_full"
		case "room is not accepting new players":
			statusCode = http.StatusConflict
			errorCode = "room_not_accepting"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RoomResponse{
		Room: updatedRoom,
	})
}

// LeaveRoom handles room leaving
func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Room ID is required",
		})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID in context",
		})
		return
	}

	// Leave room
	updatedRoom, err := h.roomService.LeaveRoom(roomID, userIDStr)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "leave_room_failed"

		// Handle specific error cases
		switch err.Error() {
		case "room not found":
			statusCode = http.StatusNotFound
			errorCode = "room_not_found"
		case "player is not in this room":
			statusCode = http.StatusConflict
			errorCode = "not_in_room"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	// If room was closed (updatedRoom is nil), return success message
	if updatedRoom == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully left room (room was closed)",
		})
		return
	}

	c.JSON(http.StatusOK, RoomResponse{
		Room: updatedRoom,
	})
}

// StartGame handles game start
func (h *RoomHandler) StartGame(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Room ID is required",
		})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID in context",
		})
		return
	}

	// Get room info before starting (need player info)
	roomData, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "room_not_found",
			Message: err.Error(),
		})
		return
	}

	// Start game (updates room status)
	err = h.roomService.StartGame(roomID, userIDStr)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "start_game_failed"

		switch {
		case err.Error() == "room not found":
			statusCode = http.StatusNotFound
			errorCode = "room_not_found"
		case err.Error() == "only room owner can start the game":
			statusCode = http.StatusForbidden
			errorCode = "not_room_owner"
		case err.Error() == "room is not ready to start game":
			statusCode = http.StatusConflict
			errorCode = "room_not_ready"
		case err.Error() == "room must have 4 players to start game":
			statusCode = http.StatusConflict
			errorCode = "insufficient_players"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	// Get updated room state
	updatedRoom, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "room_lookup_failed",
			Message: "Game started but failed to retrieve updated room state",
		})
		return
	}

	// Start game preparation sequence in background (UI + game engine)
	go h.runGamePrepareSequence(roomID, roomData.Players)

	// Return success immediately
	c.JSON(http.StatusOK, RoomResponse{
		Room: updatedRoom,
	})
}

// runGamePrepareSequence runs the complete game preparation sequence
// This includes UI preparation (countdown) and game engine startup
func (h *RoomHandler) runGamePrepareSequence(roomID string, players [4]*room.Player) {
	// Recover from any panic to prevent server crash
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in runGamePrepareSequence for room %s: %v", roomID, r)

			// Attempt to revert room status
			if revertErr := h.roomService.RevertGameStart(roomID); revertErr != nil {
				log.Printf("Failed to revert room status after panic for room %s: %v", roomID, revertErr)
			}

			// Notify players of the error
			h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
				Type: "error",
				Data: map[string]interface{}{
					"error":   "game_start_panic",
					"message": "An unexpected error occurred while starting the game",
					"room_id": roomID,
				},
				Timestamp: time.Now(),
			})

			// Send room update to reflect status change
			if roomData, err := h.roomService.GetRoom(roomID); err == nil {
				h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
					Type: "room_update",
					Data: map[string]interface{}{
						"action": "game_start_failed",
						"room":   roomData,
					},
					Timestamp: time.Now(),
				})
			}
		}
	}()

	// Step 1: Send game prepare event
	h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
		Type: "game_prepare",
		Data: map[string]interface{}{
			"room_id": roomID,
		},
		Timestamp: time.Now(),
	})

	// Step 2: Send countdown events (3 seconds)
	for i := 3; i > 0; i-- {
		time.Sleep(1 * time.Second)
		h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
			Type: "countdown",
			Data: map[string]interface{}{
				"countdown": i,
				"room_id":   roomID,
			},
			Timestamp: time.Now(),
		})
	}

	// Step 3: Send game begin event
	h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
		Type: "game_begin",
		Data: map[string]interface{}{
			"room_id":     roomID,
			"has_tribute": false, // Will be determined by game engine
		},
		Timestamp: time.Now(),
	})

	// Step 3.5: Send initial player view to each player
	// This allows frontend to render game interface immediately
	// SDK will send updated player views as the game progresses
	log.Printf("Sending initial player views for room %s", roomID)
	for i, player := range players {
		if player != nil {
			// Create player list for filtered state
			playersList := make([]map[string]interface{}, 4)
			for j, p := range players {
				if p != nil {
					playersList[j] = map[string]interface{}{
						"id":         p.ID,
						"username":   p.Username,
						"seat":       p.Seat,
						"online":     p.Online,
						"hand_count": 0, // No cards dealt yet
					}
				}
			}

			// Send initial player view
			// Use nested structure to match GameDriver format
			// Format: data.player_view contains the actual PlayerGameState
			initialView := &websocket.WSMessage{
				Type: "player_view",
				Data: map[string]interface{}{
					// Nested player_view object (matches GameDriver format)
					"player_view": map[string]interface{}{
						"player_seat":   i,
						"player_cards":  []interface{}{}, // Empty hand initially (SDK uses player_cards)
						"visible_cards": []interface{}{}, // No visible cards yet
						"game_state": map[string]interface{}{
							"id":     "",
							"status": "waiting",
							"current_match": map[string]interface{}{
								"id":          "",
								"status":      "waiting",
								"team_levels": []int{2, 2}, // Initial levels for both teams
								"current_deal": map[string]interface{}{
									"id":            "",
									"level":         2, // Starting level
									"status":        "waiting",
									"tribute_phase": nil,
									"current_trick": nil,
								},
								"players": playersList,
							},
						},
				},
				// Additional metadata (matches GameDriver format)
				"event_type":  "match_started",
				"player_seat": i,
			},
				Timestamp: time.Now(),
				PlayerID:  player.ID,
			}

			if err := h.wsManager.SendToPlayer(player.ID, initialView); err != nil {
				log.Printf("Failed to send initial player view to player %s: %v", player.ID, err)
			} else {
				log.Printf("Sent initial player view to player %s (seat %d)", player.ID, i)
			}
		}
	}

	// Step 4: Convert room players to SDK players
	sdkPlayers, err := h.convertToSDKPlayers(players)
	if err != nil {
		log.Printf("Failed to convert players for room %s: %v", roomID, err)

		// Revert room status
		if revertErr := h.roomService.RevertGameStart(roomID); revertErr != nil {
			log.Printf("Failed to revert room status for room %s: %v", roomID, revertErr)
		}

		// Notify players of the error
		h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
			Type: "error",
			Data: map[string]interface{}{
				"error":   "invalid_player_data",
				"message": "Invalid player data: " + err.Error(),
				"room_id": roomID,
			},
			Timestamp: time.Now(),
		})

		// Send room update
		if roomData, getErr := h.roomService.GetRoom(roomID); getErr == nil {
			h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
				Type: "room_update",
				Data: map[string]interface{}{
					"action": "game_start_failed",
					"room":   roomData,
				},
				Timestamp: time.Now(),
			})
		}

		return
	}

	// Step 5: Start game engine
	if err := h.driverService.StartGameWithDriver(roomID, sdkPlayers); err != nil {
		log.Printf("Failed to start game engine for room %s: %v", roomID, err)

		// Revert room status back to ready
		if revertErr := h.roomService.RevertGameStart(roomID); revertErr != nil {
			log.Printf("Failed to revert room status for room %s: %v", roomID, revertErr)
		}

		// Notify players of the error
		h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
			Type: "error",
			Data: map[string]interface{}{
				"error":   "game_start_failed",
				"message": "Failed to start game engine: " + err.Error(),
				"room_id": roomID,
			},
			Timestamp: time.Now(),
		})

		// Send room update to reflect status change back to ready
		if room, err := h.roomService.GetRoom(roomID); err == nil {
			h.wsManager.BroadcastToRoom(roomID, &websocket.WSMessage{
				Type: "room_update",
				Data: map[string]interface{}{
					"action": "game_start_failed",
					"room":   room,
				},
				Timestamp: time.Now(),
			})
		}

		return
	}

	log.Printf("Game successfully started for room %s", roomID)
}

// convertToSDKPlayers converts room players to SDK player format
func (h *RoomHandler) convertToSDKPlayers(players [4]*room.Player) ([]sdk.Player, error) {
	sdkPlayers := make([]sdk.Player, 4)
	for i := 0; i < 4; i++ {
		if players[i] == nil {
			return nil, fmt.Errorf("player at seat %d is nil", i)
		}
		sdkPlayers[i] = sdk.Player{
			ID:       players[i].ID,
			Username: players[i].Username,
			Seat:     players[i].Seat,
			Online:   players[i].Online,
			AutoPlay: false, // Initially not in auto-play mode
		}
	}
	return sdkPlayers, nil
}

// GetMyRoom handles getting the current user's room
func (h *RoomHandler) GetMyRoom(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user ID in context",
		})
		return
	}

	// Get player's room
	playerRoom, err := h.roomService.GetPlayerRoom(userIDStr)
	if err != nil {
		if err.Error() == "player is not in any room" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_in_room",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "room_lookup_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RoomResponse{
		Room: playerRoom,
	})
}

// Note: Room routes are registered directly in main.go
// This keeps the route registration centralized with other services