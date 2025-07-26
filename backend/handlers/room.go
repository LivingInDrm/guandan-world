package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"guandan-world/backend/auth"
	"guandan-world/backend/room"

	"github.com/gin-gonic/gin"
)

// RoomHandler handles room-related HTTP requests
type RoomHandler struct {
	roomService room.RoomService
	authService auth.AuthService
}

// NewRoomHandler creates a new room handler
func NewRoomHandler(roomService room.RoomService, authService auth.AuthService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		authService: authService,
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

	// Start game
	err := h.roomService.StartGame(roomID, userIDStr)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "start_game_failed"

		// Handle specific error cases
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

	// Get updated room to return current state
	updatedRoom, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "room_lookup_failed",
			Message: "Game started but failed to retrieve updated room state",
		})
		return
	}

	c.JSON(http.StatusOK, RoomResponse{
		Room: updatedRoom,
	})
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

// RegisterRoutes registers all room routes
func (h *RoomHandler) RegisterRoutes(router *gin.Engine, authHandler *AuthHandler) {
	rooms := router.Group("/api/rooms")
	{
		// Public routes (with authentication)
		rooms.Use(authHandler.JWTMiddleware())
		
		rooms.GET("", h.GetRooms)                    // GET /api/rooms - list rooms with pagination
		rooms.POST("", h.CreateRoom)                 // POST /api/rooms - create room
		rooms.GET("/my", h.GetMyRoom)                // GET /api/rooms/my - get current user's room
		rooms.GET("/:id", h.GetRoom)                 // GET /api/rooms/:id - get specific room
		rooms.POST("/:id/join", h.JoinRoom)          // POST /api/rooms/:id/join - join room
		rooms.POST("/:id/leave", h.LeaveRoom)        // POST /api/rooms/:id/leave - leave room
		rooms.POST("/:id/start", h.StartGame)        // POST /api/rooms/:id/start - start game
	}
}