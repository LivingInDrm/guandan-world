package room

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"guandan-world/backend/auth"
)

// RoomStatus represents the status of a room
type RoomStatus int

const (
	RoomStatusWaiting RoomStatus = iota // 等待玩家加入
	RoomStatusReady                     // 4人已满，可以开始
	RoomStatusPlaying                   // 游戏进行中
	RoomStatusClosed                    // 房间已关闭
)

// String returns the string representation of RoomStatus
func (rs RoomStatus) String() string {
	switch rs {
	case RoomStatusWaiting:
		return "waiting"
	case RoomStatusReady:
		return "ready"
	case RoomStatusPlaying:
		return "playing"
	case RoomStatusClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// Player represents a player in a room
type Player struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Nickname  *string `json:"nickname,omitempty"`
	AvatarKey *string `json:"avatar_key,omitempty"`
	Seat      int     `json:"seat"`
	Online    bool    `json:"online"`
}

// Room represents a game room
type Room struct {
	ID          string     `json:"id"`
	RoomCode    string     `json:"room_code"`
	Status      RoomStatus `json:"status"`
	Players     [4]*Player `json:"players"`      // Fixed 4 seats (0-3)
	Owner       string     `json:"owner"`        // Owner user ID
	PlayerCount int        `json:"player_count"` // Number of players currently in room
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RoomInfo provides room information for listing
type RoomInfo struct {
	ID          string     `json:"id"`
	RoomCode    string     `json:"room_code"`
	Status      RoomStatus `json:"status"`
	PlayerCount int        `json:"player_count"`
	Players     []*Player  `json:"players"`
	Owner       string     `json:"owner"`
	CanJoin     bool       `json:"can_join"`
	CreatedAt   time.Time  `json:"created_at"`
}

// RoomListResponse represents the response for room list queries
type RoomListResponse struct {
	Rooms      []*RoomInfo `json:"rooms"`
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

// RoomService interface defines room management operations
type RoomService interface {
	CreateRoom(ownerID string) (*Room, error)
	JoinRoom(roomID, playerID string) (*Room, error)
	JoinRoomByCode(roomCode, playerID string) (*Room, error)
	LeaveRoom(roomID, playerID string) (*Room, error)
	GetRoom(roomID string) (*Room, error)
	GetRoomList(page, limit int, statusFilter *RoomStatus) (*RoomListResponse, error)
	StartGame(roomID, playerID string) error
	CloseRoom(roomID string) error
	GetPlayerRoom(playerID string) (*Room, error)
	RevertGameStart(roomID string) error
	QuickJoin(playerID string) (*Room, error)
}

// roomService implements RoomService interface
type roomService struct {
	rooms       map[string]*Room  // roomID -> Room
	roomCodes   map[string]string // roomCode -> roomID
	playerRooms map[string]string // playerID -> roomID
	authService auth.AuthService
	mu          sync.RWMutex
}

// NewRoomService creates a new room service
func NewRoomService(authService auth.AuthService) RoomService {
	rand.Seed(time.Now().UnixNano())
	return &roomService{
		rooms:       make(map[string]*Room),
		roomCodes:   make(map[string]string),
		playerRooms: make(map[string]string),
		authService: authService,
	}
}

// generateUniqueRoomCode generates a unique 4-digit room code
func (s *roomService) generateUniqueRoomCode() (string, error) {
	for i := 0; i < 100; i++ {
		code := fmt.Sprintf("%04d", rand.Intn(10000))
		if _, exists := s.roomCodes[code]; !exists {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique room code: too many active rooms")
}

// CreateRoom creates a new room with the specified owner
func (s *roomService) CreateRoom(ownerID string) (*Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate owner exists
	owner, err := s.authService.GetUserByID(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner: %w", err)
	}

	// Check if player is already in a room
	if existingRoomID, exists := s.playerRooms[ownerID]; exists {
		return nil, fmt.Errorf("player is already in room %s", existingRoomID)
	}

	// Generate room ID
	roomID := fmt.Sprintf("room_%d", time.Now().UnixNano())

	// Generate unique room code
	roomCode, err := s.generateUniqueRoomCode()
	if err != nil {
		return nil, err
	}

	// Create room
	room := &Room{
		ID:          roomID,
		RoomCode:    roomCode,
		Status:      RoomStatusWaiting,
		Players:     [4]*Player{},
		Owner:       ownerID,
		PlayerCount: 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add owner as first player (seat 0)
	room.Players[0] = &Player{
		ID:        owner.ID,
		Username:  owner.Username,
		Nickname:  owner.Nickname,
		AvatarKey: owner.AvatarKey,
		Seat:      0,
		Online:    owner.Online,
	}

	// Store room, room code, and player mapping
	s.rooms[roomID] = room
	s.roomCodes[roomCode] = roomID
	s.playerRooms[ownerID] = roomID

	return room, nil
}

// JoinRoom adds a player to an existing room
func (s *roomService) JoinRoom(roomID, playerID string) (*Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate player exists
	player, err := s.authService.GetUserByID(playerID)
	if err != nil {
		return nil, fmt.Errorf("invalid player: %w", err)
	}

	// Check if player is already in a room
	if existingRoomID, exists := s.playerRooms[playerID]; exists {
		if existingRoomID == roomID {
			// Player is already in this room, return current room
			return s.rooms[roomID], nil
		}
		return nil, fmt.Errorf("player is already in room %s", existingRoomID)
	}

	// Get room
	room, exists := s.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	// Check room status
	if room.Status != RoomStatusWaiting {
		return nil, errors.New("room is not accepting new players")
	}

	// Check if room is full
	if room.PlayerCount >= 4 {
		return nil, errors.New("room is full")
	}

	// Find available seat
	var availableSeat int = -1
	for i := 0; i < 4; i++ {
		if room.Players[i] == nil {
			availableSeat = i
			break
		}
	}

	if availableSeat == -1 {
		return nil, errors.New("no available seats")
	}

	// Add player to room
	room.Players[availableSeat] = &Player{
		ID:        player.ID,
		Username:  player.Username,
		Nickname:  player.Nickname,
		AvatarKey: player.AvatarKey,
		Seat:      availableSeat,
		Online:    player.Online,
	}

	room.PlayerCount++
	room.UpdatedAt = time.Now()

	// Update room status if full
	if room.PlayerCount == 4 {
		room.Status = RoomStatusReady
	}

	// Store player mapping
	s.playerRooms[playerID] = roomID

	return room, nil
}

// LeaveRoom removes a player from a room
func (s *roomService) LeaveRoom(roomID, playerID string) (*Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get room
	room, exists := s.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	// Check if player is in this room
	currentRoomID, inRoom := s.playerRooms[playerID]
	if !inRoom || currentRoomID != roomID {
		return nil, errors.New("player is not in this room")
	}

	// Find and remove player
	var playerSeat int = -1
	for i := 0; i < 4; i++ {
		if room.Players[i] != nil && room.Players[i].ID == playerID {
			playerSeat = i
			break
		}
	}

	if playerSeat == -1 {
		return nil, errors.New("player not found in room")
	}

	// Remove player
	room.Players[playerSeat] = nil
	room.PlayerCount--
	room.UpdatedAt = time.Now()

	// Remove player mapping
	delete(s.playerRooms, playerID)

	// Handle owner leaving
	if room.Owner == playerID {
		// Find new owner from remaining players
		newOwner := ""
		for i := 0; i < 4; i++ {
			if room.Players[i] != nil {
				newOwner = room.Players[i].ID
				break
			}
		}

		if newOwner != "" {
			room.Owner = newOwner
		} else {
			// No players left, close room
			room.Status = RoomStatusClosed
			delete(s.roomCodes, room.RoomCode)
			delete(s.rooms, roomID)
			return nil, nil // Room closed
		}
	}

	// Check if game should end due to insufficient players
	if room.Status == RoomStatusPlaying && room.PlayerCount < 4 {
		// Game cannot continue with less than 4 players, close room
		room.Status = RoomStatusClosed
		// Remove all remaining player mappings
		for i := 0; i < 4; i++ {
			if room.Players[i] != nil {
				delete(s.playerRooms, room.Players[i].ID)
			}
		}
		delete(s.roomCodes, room.RoomCode)
		delete(s.rooms, roomID)
		return nil, nil // Room closed due to game ending
	}

	// Update room status
	if room.PlayerCount < 4 && room.Status == RoomStatusReady {
		room.Status = RoomStatusWaiting
	}

	return room, nil
}

// GetRoom retrieves a room by ID
func (s *roomService) GetRoom(roomID string) (*Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	return room, nil
}

// GetRoomList retrieves a paginated list of rooms with optional status filtering
func (s *roomService) GetRoomList(page, limit int, statusFilter *RoomStatus) (*RoomListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 12 // Default to 12 as specified in requirements
	}

	// Collect and filter rooms
	var allRooms []*Room
	for _, room := range s.rooms {
		// Skip closed rooms
		if room.Status == RoomStatusClosed {
			continue
		}

		// Apply status filter if provided
		if statusFilter != nil && room.Status != *statusFilter {
			continue
		}

		allRooms = append(allRooms, room)
	}

	// Sort rooms: waiting rooms first, then by player count descending
	sort.Slice(allRooms, func(i, j int) bool {
		roomA, roomB := allRooms[i], allRooms[j]

		// Waiting rooms first
		if roomA.Status == RoomStatusWaiting && roomB.Status != RoomStatusWaiting {
			return true
		}
		if roomA.Status != RoomStatusWaiting && roomB.Status == RoomStatusWaiting {
			return false
		}

		// Then by player count descending
		if roomA.PlayerCount != roomB.PlayerCount {
			return roomA.PlayerCount > roomB.PlayerCount
		}

		// Finally by creation time (newer first)
		return roomA.CreatedAt.After(roomB.CreatedAt)
	})

	totalCount := len(allRooms)

	// Apply pagination
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalCount {
		// Page beyond available data
		return &RoomListResponse{
			Rooms:      []*RoomInfo{},
			TotalCount: totalCount,
			Page:       page,
			Limit:      limit,
		}, nil
	}

	if endIndex > totalCount {
		endIndex = totalCount
	}

	pagedRooms := allRooms[startIndex:endIndex]

	// Convert to RoomInfo
	roomInfos := make([]*RoomInfo, len(pagedRooms))
	for i, room := range pagedRooms {
		// Collect non-nil players
		var players []*Player
		for _, player := range room.Players {
			if player != nil {
				players = append(players, player)
			}
		}

		// Determine if room can be joined
		canJoin := room.Status == RoomStatusWaiting && room.PlayerCount < 4

		roomInfos[i] = &RoomInfo{
			ID:          room.ID,
			RoomCode:    room.RoomCode,
			Status:      room.Status,
			PlayerCount: room.PlayerCount,
			Players:     players,
			Owner:       room.Owner,
			CanJoin:     canJoin,
			CreatedAt:   room.CreatedAt,
		}
	}

	return &RoomListResponse{
		Rooms:      roomInfos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}, nil
}

// StartGame starts a game in the specified room
func (s *roomService) StartGame(roomID, playerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get room
	room, exists := s.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	// Check if player is the owner
	if room.Owner != playerID {
		return errors.New("only room owner can start the game")
	}

	// Check room status
	if room.Status != RoomStatusReady {
		return errors.New("room is not ready to start game")
	}

	// Check if room is full
	if room.PlayerCount != 4 {
		return errors.New("room must have 4 players to start game")
	}

	// Update room status
	room.Status = RoomStatusPlaying
	room.UpdatedAt = time.Now()

	return nil
}

// RevertGameStart reverts room status from playing back to ready
// Used when game engine fails to start after status was changed
func (s *roomService) RevertGameStart(roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	// Only revert if currently in playing status
	if room.Status == RoomStatusPlaying {
		room.Status = RoomStatusReady
		room.UpdatedAt = time.Now()
	}

	return nil
}

// CloseRoom closes a room
func (s *roomService) CloseRoom(roomID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	// Remove all player mappings
	for _, player := range room.Players {
		if player != nil {
			delete(s.playerRooms, player.ID)
		}
	}

	// Remove room code and room
	delete(s.roomCodes, room.RoomCode)
	delete(s.rooms, roomID)

	return nil
}

// GetPlayerRoom finds the room that a player is currently in
func (s *roomService) GetPlayerRoom(playerID string) (*Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roomID, exists := s.playerRooms[playerID]
	if !exists {
		return nil, errors.New("player is not in any room")
	}

	room, exists := s.rooms[roomID]
	if !exists {
		// Clean up stale mapping
		delete(s.playerRooms, playerID)
		return nil, errors.New("room not found")
	}

	return room, nil
}

// JoinRoomByCode allows a player to join a room using a room code
func (s *roomService) JoinRoomByCode(roomCode, playerID string) (*Room, error) {
	s.mu.RLock()
	roomID, exists := s.roomCodes[roomCode]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("room code not found")
	}

	return s.JoinRoom(roomID, playerID)
}

// QuickJoin finds the best available room (most players, not full) and joins it
// If no suitable room exists, creates a new room
func (s *roomService) QuickJoin(playerID string) (*Room, error) {
	s.mu.RLock()
	if existingRoomID, exists := s.playerRooms[playerID]; exists {
		s.mu.RUnlock()
		return nil, fmt.Errorf("player is already in room %s", existingRoomID)
	}

	var candidates []*Room
	for _, room := range s.rooms {
		if room.Status == RoomStatusWaiting && room.PlayerCount < 4 {
			candidates = append(candidates, room)
		}
	}
	s.mu.RUnlock()

	if len(candidates) == 0 {
		return s.CreateRoom(playerID)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].PlayerCount > candidates[j].PlayerCount
	})

	rand.Shuffle(len(candidates), func(i, j int) {
		if candidates[i].PlayerCount == candidates[j].PlayerCount {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		}
	})

	for _, room := range candidates {
		result, err := s.JoinRoom(room.ID, playerID)
		if err == nil {
			return result, nil
		}
	}

	return s.CreateRoom(playerID)
}