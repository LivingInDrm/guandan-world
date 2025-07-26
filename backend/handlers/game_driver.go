package handlers

import (
	"net/http"

	"guandan-world/backend/game"
	"guandan-world/sdk"

	"github.com/gin-gonic/gin"
)

// GameDriverHandler handles game operations using the GameDriver architecture
type GameDriverHandler struct {
	driverService *game.DriverService
}

// NewGameDriverHandler creates a new game driver handler
func NewGameDriverHandler(driverService *game.DriverService) *GameDriverHandler {
	return &GameDriverHandler{
		driverService: driverService,
	}
}

// StartGameWithDriverRequest represents the request to start a game with driver
type StartGameWithDriverRequest struct {
	RoomID  string       `json:"room_id" binding:"required"`
	Players []sdk.Player `json:"players" binding:"required,len=4"`
}

// StartGameWithDriver starts a new game using the GameDriver
// @Summary Start a new game with GameDriver
// @Description Starts a new game for a room using the SDK's GameDriver architecture
// @Tags game-driver
// @Accept json
// @Produce json
// @Param request body StartGameWithDriverRequest true "Start game request"
// @Success 200 {object} map[string]interface{} "Game started successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/start [post]
func (h *GameDriverHandler) StartGameWithDriver(c *gin.Context) {
	var req StartGameWithDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Start game using driver service
	err := h.driverService.StartGameWithDriver(req.RoomID, req.Players)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Game started with driver",
		"room_id": req.RoomID,
	})
}

// PlayDecisionRequest represents a play decision request
type PlayDecisionRequest struct {
	RoomID     string   `json:"room_id" binding:"required"`
	PlayerSeat int      `json:"player_seat" binding:"min=0,max=3"`
	Action     string   `json:"action" binding:"required,oneof=play pass"`
	CardIDs    []string `json:"card_ids,omitempty"`
}

// SubmitPlayDecision submits a player's play decision
// @Summary Submit play decision
// @Description Submits a player's decision to play cards or pass
// @Tags game-driver
// @Accept json
// @Produce json
// @Param request body PlayDecisionRequest true "Play decision"
// @Success 200 {object} map[string]interface{} "Decision submitted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/play-decision [post]
func (h *GameDriverHandler) SubmitPlayDecision(c *gin.Context) {
	var req PlayDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get current game level for card parsing
	gameStatus, err := h.driverService.GetGameStatus(req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get game status: " + err.Error(),
		})
		return
	}

	// Extract level from game status (default to 2 if not found)
	level := 2
	if gameStatus != nil {
		if matchDetails, ok := gameStatus["match_details"]; ok {
			if matchDetailsMap, ok := matchDetails.(map[string]interface{}); ok {
				if teamLevels, ok := matchDetailsMap["TeamLevels"]; ok {
					if teamLevelsMap, ok := teamLevels.(map[int]int); ok {
						// Use team 0's level as current level
						if lvl, exists := teamLevelsMap[0]; exists {
							level = lvl
						}
					}
				}
			}
		}
	}

	// Convert card IDs to cards using SDK function
	var cards []*sdk.Card
	if req.Action == "play" && len(req.CardIDs) > 0 {
		cards = make([]*sdk.Card, len(req.CardIDs))
		for i, cardID := range req.CardIDs {
			card, parseErr := sdk.ParseCardFromID(cardID, level)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "Invalid card ID: " + cardID,
				})
				return
			}
			cards[i] = card
		}
	}

	// Create play decision
	decision := &sdk.PlayDecision{
		Action: sdk.ActionType(req.Action),
		Cards:  cards,
	}

	// Submit to driver service
	submitErr := h.driverService.SubmitPlayDecision(req.RoomID, req.PlayerSeat, decision)
	if submitErr != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: submitErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Play decision submitted",
	})
}

// TributeSelectionRequest represents a tribute selection request
type TributeSelectionRequest struct {
	RoomID     string `json:"room_id" binding:"required"`
	PlayerSeat int    `json:"player_seat" binding:"min=0,max=3"`
	CardID     string `json:"card_id" binding:"required"`
}

// SubmitTributeSelection submits a tribute selection
// @Summary Submit tribute selection
// @Description Submits a player's tribute card selection (for double-down)
// @Tags game-driver
// @Accept json
// @Produce json
// @Param request body TributeSelectionRequest true "Tribute selection"
// @Success 200 {object} map[string]interface{} "Selection submitted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/tribute-select [post]
func (h *GameDriverHandler) SubmitTributeSelection(c *gin.Context) {
	var req TributeSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Submit to driver service
	err := h.driverService.SubmitTributeSelection(req.RoomID, req.PlayerSeat, req.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tribute selection submitted",
	})
}

// ReturnTributeRequest represents a return tribute request
type ReturnTributeRequest struct {
	RoomID     string `json:"room_id" binding:"required"`
	PlayerSeat int    `json:"player_seat" binding:"min=0,max=3"`
	CardID     string `json:"card_id" binding:"required"`
}

// SubmitReturnTribute submits a return tribute
// @Summary Submit return tribute
// @Description Submits a player's return tribute card
// @Tags game-driver
// @Accept json
// @Produce json
// @Param request body ReturnTributeRequest true "Return tribute"
// @Success 200 {object} map[string]interface{} "Return submitted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/tribute-return [post]
func (h *GameDriverHandler) SubmitReturnTribute(c *gin.Context) {
	var req ReturnTributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	// Submit to driver service
	err := h.driverService.SubmitReturnTribute(req.RoomID, req.PlayerSeat, req.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Return tribute submitted",
	})
}

// GetGameStatus gets the current game status
// @Summary Get game status
// @Description Gets the current status of a game
// @Tags game-driver
// @Accept json
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {object} map[string]interface{} "Game status"
// @Failure 404 {object} ErrorResponse "Game not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/status/{room_id} [get]
func (h *GameDriverHandler) GetGameStatus(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Room ID is required",
		})
		return
	}

	status, err := h.driverService.GetGameStatus(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// StopGame stops a game
// @Summary Stop game
// @Description Stops an active game
// @Tags game-driver
// @Accept json
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {object} map[string]interface{} "Game stopped successfully"
// @Failure 404 {object} ErrorResponse "Game not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/game/driver/stop/{room_id} [post]
func (h *GameDriverHandler) StopGame(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Room ID is required",
		})
		return
	}

	err := h.driverService.StopGame(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Game stopped",
		"room_id": roomID,
	})
}
