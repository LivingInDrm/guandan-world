package room

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"guandan-world/backend/auth"
)

// Mock auth service for testing
type mockAuthService struct {
	users map[string]*auth.User
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		users: make(map[string]*auth.User),
	}
}

func (m *mockAuthService) Register(username, password string) (*auth.User, error) {
	// Check if username already exists
	for _, user := range m.users {
		if user.Username == username {
			return nil, errors.New("username already exists")
		}
	}
	
	user := &auth.User{
		ID:       "user_" + username,
		Username: username,
		Online:   true,
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *mockAuthService) Login(username, password string) (*auth.AuthToken, error) {
	return nil, nil
}

func (m *mockAuthService) ValidateToken(token string) (*auth.User, error) {
	return nil, nil
}

func (m *mockAuthService) Logout(token string) error {
	return nil
}

func (m *mockAuthService) GetUserByID(userID string) (*auth.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func TestRoomService_CreateRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	user, err := authSvc.Register("testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test creating room
	room, err := roomSvc.CreateRoom(user.ID)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	// Verify room properties
	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}
	if room.Status != RoomStatusWaiting {
		t.Errorf("Expected room status to be waiting, got %v", room.Status)
	}
	if room.Owner != user.ID {
		t.Errorf("Expected room owner to be %s, got %s", user.ID, room.Owner)
	}
	if room.PlayerCount != 1 {
		t.Errorf("Expected player count to be 1, got %d", room.PlayerCount)
	}
	if room.Players[0] == nil {
		t.Error("First player should not be nil")
	}
	if room.Players[0].ID != user.ID {
		t.Errorf("Expected first player ID to be %s, got %s", user.ID, room.Players[0].ID)
	}
	if room.Players[0].Seat != 0 {
		t.Errorf("Expected first player seat to be 0, got %d", room.Players[0].Seat)
	}
}

func TestRoomService_CreateRoom_InvalidOwner(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Test creating room with invalid owner
	_, err := roomSvc.CreateRoom("invalid_user")
	if err == nil {
		t.Error("Expected error when creating room with invalid owner")
	}
}

func TestRoomService_CreateRoom_PlayerAlreadyInRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	user, _ := authSvc.Register("testuser", "password")

	// Create first room
	_, err := roomSvc.CreateRoom(user.ID)
	if err != nil {
		t.Fatalf("Failed to create first room: %v", err)
	}

	// Try to create second room with same user
	_, err = roomSvc.CreateRoom(user.ID)
	if err == nil {
		t.Error("Expected error when user tries to create second room")
	}
}

func TestRoomService_JoinRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	owner, _ := authSvc.Register("owner", "password")
	player, _ := authSvc.Register("player", "password")

	// Create room
	room, err := roomSvc.CreateRoom(owner.ID)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	// Join room
	updatedRoom, err := roomSvc.JoinRoom(room.ID, player.ID)
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}

	// Verify room state
	if updatedRoom.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", updatedRoom.PlayerCount)
	}
	if updatedRoom.Status != RoomStatusWaiting {
		t.Errorf("Expected room status to be waiting, got %v", updatedRoom.Status)
	}

	// Find the joined player
	var joinedPlayer *Player
	for _, p := range updatedRoom.Players {
		if p != nil && p.ID == player.ID {
			joinedPlayer = p
			break
		}
	}

	if joinedPlayer == nil {
		t.Error("Joined player not found in room")
	}
	if joinedPlayer.Seat < 0 || joinedPlayer.Seat > 3 {
		t.Errorf("Invalid seat number: %d", joinedPlayer.Seat)
	}
}

func TestRoomService_JoinRoom_RoomFull(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	users := make([]*auth.User, 5)
	for i := 0; i < 5; i++ {
		user, _ := authSvc.Register(fmt.Sprintf("user%d", i), "password")
		users[i] = user
	}

	// Create room with first user
	room, _ := roomSvc.CreateRoom(users[0].ID)

	// Fill room with 3 more users
	for i := 1; i < 4; i++ {
		_, err := roomSvc.JoinRoom(room.ID, users[i].ID)
		if err != nil {
			t.Fatalf("Failed to join room with user %d: %v", i, err)
		}
	}

	// Verify room is ready
	updatedRoom, _ := roomSvc.GetRoom(room.ID)
	if updatedRoom.Status != RoomStatusReady {
		t.Errorf("Expected room status to be ready, got %v", updatedRoom.Status)
	}
	if updatedRoom.PlayerCount != 4 {
		t.Errorf("Expected player count to be 4, got %d", updatedRoom.PlayerCount)
	}

	// Try to join with 5th user (should fail)
	_, err := roomSvc.JoinRoom(room.ID, users[4].ID)
	if err == nil {
		t.Error("Expected error when trying to join full room")
	}
}

func TestRoomService_LeaveRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	owner, _ := authSvc.Register("owner", "password")
	player, _ := authSvc.Register("player", "password")

	// Create room and join
	room, _ := roomSvc.CreateRoom(owner.ID)
	roomSvc.JoinRoom(room.ID, player.ID)

	// Leave room
	updatedRoom, err := roomSvc.LeaveRoom(room.ID, player.ID)
	if err != nil {
		t.Fatalf("Failed to leave room: %v", err)
	}

	// Verify room state
	if updatedRoom.PlayerCount != 1 {
		t.Errorf("Expected player count to be 1, got %d", updatedRoom.PlayerCount)
	}

	// Verify player is removed
	for _, p := range updatedRoom.Players {
		if p != nil && p.ID == player.ID {
			t.Error("Player should have been removed from room")
		}
	}
}

func TestRoomService_LeaveRoom_OwnerLeaves(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	owner, _ := authSvc.Register("owner", "password")
	player, _ := authSvc.Register("player", "password")

	// Create room and join
	room, _ := roomSvc.CreateRoom(owner.ID)
	roomSvc.JoinRoom(room.ID, player.ID)

	// Owner leaves room
	updatedRoom, err := roomSvc.LeaveRoom(room.ID, owner.ID)
	if err != nil {
		t.Fatalf("Failed to leave room: %v", err)
	}

	// Verify new owner
	if updatedRoom.Owner != player.ID {
		t.Errorf("Expected new owner to be %s, got %s", player.ID, updatedRoom.Owner)
	}
	if updatedRoom.PlayerCount != 1 {
		t.Errorf("Expected player count to be 1, got %d", updatedRoom.PlayerCount)
	}
}

func TestRoomService_LeaveRoom_LastPlayerLeaves(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	owner, _ := authSvc.Register("owner", "password")

	// Create room
	room, _ := roomSvc.CreateRoom(owner.ID)

	// Owner leaves room (last player)
	updatedRoom, err := roomSvc.LeaveRoom(room.ID, owner.ID)
	if err != nil {
		t.Fatalf("Failed to leave room: %v", err)
	}

	// Room should be closed (nil returned)
	if updatedRoom != nil {
		t.Error("Expected room to be closed when last player leaves")
	}

	// Verify room is deleted
	_, err = roomSvc.GetRoom(room.ID)
	if err == nil {
		t.Error("Expected error when getting deleted room")
	}
}

func TestRoomService_GetRoomList(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users and rooms
	for i := 0; i < 5; i++ {
		user, err := authSvc.Register(fmt.Sprintf("getroomlistuser%d_%d", i, time.Now().UnixNano()), "password")
		if err != nil {
			t.Fatalf("Failed to register user %d: %v", i, err)
		}
		_, err = roomSvc.CreateRoom(user.ID)
		if err != nil {
			t.Fatalf("Failed to create room for user %d: %v", i, err)
		}
	}

	// Get room list
	response, err := roomSvc.GetRoomList(1, 3, nil)
	if err != nil {
		t.Fatalf("Failed to get room list: %v", err)
	}

	// Verify response
	if response.TotalCount != 5 {
		t.Errorf("Expected total count to be 5, got %d", response.TotalCount)
	}
	if len(response.Rooms) != 3 {
		t.Errorf("Expected 3 rooms in response, got %d", len(response.Rooms))
	}
	if response.Page != 1 {
		t.Errorf("Expected page to be 1, got %d", response.Page)
	}
	if response.Limit != 3 {
		t.Errorf("Expected limit to be 3, got %d", response.Limit)
	}

	// Verify room info
	for _, roomInfo := range response.Rooms {
		if roomInfo.Status != RoomStatusWaiting {
			t.Errorf("Expected room status to be waiting, got %v", roomInfo.Status)
		}
		if roomInfo.PlayerCount != 1 {
			t.Errorf("Expected player count to be 1, got %d", roomInfo.PlayerCount)
		}
		if !roomInfo.CanJoin {
			t.Error("Expected room to be joinable")
		}
	}
}

func TestRoomService_GetRoomList_WithStatusFilter(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create rooms with different statuses
	users := make([]*auth.User, 5)
	for i := 0; i < 5; i++ {
		user, _ := authSvc.Register(fmt.Sprintf("user%d", i), "password")
		users[i] = user
	}

	// Create a waiting room
	_, _ = roomSvc.CreateRoom(users[0].ID)

	// Create a ready room (4 players)
	readyRoom, _ := roomSvc.CreateRoom(users[1].ID)
	for i := 2; i < 5; i++ {
		roomSvc.JoinRoom(readyRoom.ID, users[i].ID)
	}

	// Get only waiting rooms
	waitingStatus := RoomStatusWaiting
	response, err := roomSvc.GetRoomList(1, 10, &waitingStatus)
	if err != nil {
		t.Fatalf("Failed to get room list: %v", err)
	}

	if response.TotalCount != 1 {
		t.Errorf("Expected 1 waiting room, got %d", response.TotalCount)
	}
	if len(response.Rooms) != 1 {
		t.Errorf("Expected 1 room in response, got %d", len(response.Rooms))
	}
	if response.Rooms[0].Status != RoomStatusWaiting {
		t.Errorf("Expected room status to be waiting, got %v", response.Rooms[0].Status)
	}
}

func TestRoomService_StartGame(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	users := make([]*auth.User, 4)
	for i := 0; i < 4; i++ {
		user, _ := authSvc.Register(fmt.Sprintf("user%d", i), "password")
		users[i] = user
	}

	// Create room and fill it
	room, _ := roomSvc.CreateRoom(users[0].ID)
	for i := 1; i < 4; i++ {
		roomSvc.JoinRoom(room.ID, users[i].ID)
	}

	// Start game
	err := roomSvc.StartGame(room.ID, users[0].ID)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Verify room status
	updatedRoom, _ := roomSvc.GetRoom(room.ID)
	if updatedRoom.Status != RoomStatusPlaying {
		t.Errorf("Expected room status to be playing, got %v", updatedRoom.Status)
	}
}

func TestRoomService_StartGame_NotOwner(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test users
	owner, _ := authSvc.Register("owner", "password")
	player, _ := authSvc.Register("player", "password")

	// Create room
	room, _ := roomSvc.CreateRoom(owner.ID)
	roomSvc.JoinRoom(room.ID, player.ID)

	// Try to start game as non-owner
	err := roomSvc.StartGame(room.ID, player.ID)
	if err == nil {
		t.Error("Expected error when non-owner tries to start game")
	}
}

func TestRoomService_StartGame_NotReady(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	owner, _ := authSvc.Register("owner", "password")

	// Create room (only 1 player)
	room, _ := roomSvc.CreateRoom(owner.ID)

	// Try to start game with insufficient players
	err := roomSvc.StartGame(room.ID, owner.ID)
	if err == nil {
		t.Error("Expected error when trying to start game with insufficient players")
	}
}

func TestRoomService_GetPlayerRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	user, _ := authSvc.Register("testuser", "password")

	// Create room
	room, _ := roomSvc.CreateRoom(user.ID)

	// Get player room
	playerRoom, err := roomSvc.GetPlayerRoom(user.ID)
	if err != nil {
		t.Fatalf("Failed to get player room: %v", err)
	}

	if playerRoom.ID != room.ID {
		t.Errorf("Expected room ID to be %s, got %s", room.ID, playerRoom.ID)
	}
}

func TestRoomService_GetPlayerRoom_NotInRoom(t *testing.T) {
	authSvc := newMockAuthService()
	roomSvc := NewRoomService(authSvc)

	// Create test user
	user, _ := authSvc.Register("testuser", "password")

	// Try to get room for user not in any room
	_, err := roomSvc.GetPlayerRoom(user.ID)
	if err == nil {
		t.Error("Expected error when getting room for user not in any room")
	}
}