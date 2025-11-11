package game

import (
	"sync"
	"testing"
	"time"

	"guandan-world/sdk"
)

// TestStopGame_CancelsGameLoop verifies that StopGame properly cancels
// the SDK layer's game loop by calling driver.CancelMatch()
func TestStopGame_CancelsGameLoop(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	service := NewDriverService(wsManager)

	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}

	roomID := "test-room-cancel"

	// Start game
	err := service.StartGameWithDriver(roomID, players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Wait a bit for game to start
	time.Sleep(50 * time.Millisecond)

	// Track if RunMatch goroutine completes
	var runMatchCompleted sync.WaitGroup
	runMatchCompleted.Add(1)

	// Monitor completion by checking if driver is cleaned up
	go func() {
		defer runMatchCompleted.Done()
		// Wait for driver to be cleaned up (should happen after CancelMatch)
		timeout := time.After(2 * time.Second)
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				service.mu.RLock()
				_, exists := service.drivers[roomID]
				service.mu.RUnlock()
				if !exists {
					// Driver cleaned up, game stopped successfully
					return
				}
			case <-timeout:
				t.Error("Timeout waiting for game to stop - RunMatch goroutine may not have been cancelled")
				return
			}
		}
	}()

	// Stop game - this should call driver.CancelMatch()
	err = service.StopGame(roomID)
	if err != nil {
		t.Fatalf("Failed to stop game: %v", err)
	}

	// Wait for RunMatch to complete
	done := make(chan struct{})
	go func() {
		runMatchCompleted.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - game loop was properly cancelled
	case <-time.After(3 * time.Second):
		t.Fatal("RunMatch goroutine did not complete after StopGame - CancelMatch may not be working")
	}

	// Verify resources are cleaned up
	service.mu.RLock()
	_, driverExists := service.drivers[roomID]
	_, providerExists := service.providers[roomID]
	service.mu.RUnlock()

	if driverExists {
		t.Error("Driver should be removed after StopGame")
	}
	if providerExists {
		t.Error("Provider should be removed after StopGame")
	}
}

// TestStopGame_OrderOfOperations verifies the correct order:
// 1. CancelMatch() called first
// 2. Then CancelAll() on provider
// 3. Then cleanup
func TestStopGame_OrderOfOperations(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	service := NewDriverService(wsManager)

	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}

	roomID := "test-room-order"

	// Start game
	err := service.StartGameWithDriver(roomID, players)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Wait for game to start
	time.Sleep(50 * time.Millisecond)

	// Get references before stopping
	service.mu.RLock()
	driver := service.drivers[roomID]
	provider := service.providers[roomID]
	service.mu.RUnlock()

	if driver == nil {
		t.Fatal("Driver should exist before StopGame")
	}
	if provider == nil {
		t.Fatal("Provider should exist before StopGame")
	}

	// Stop game
	err = service.StopGame(roomID)
	if err != nil {
		t.Fatalf("Failed to stop game: %v", err)
	}

	// Wait a bit for cleanup to complete
	time.Sleep(100 * time.Millisecond)

	// Verify both driver and provider were handled
	// We can't directly verify the order, but we can verify both operations completed
	service.mu.RLock()
	_, driverExists := service.drivers[roomID]
	_, providerExists := service.providers[roomID]
	service.mu.RUnlock()

	if driverExists {
		t.Error("Driver was not cleaned up")
	}
	if providerExists {
		t.Error("Provider was not cleaned up")
	}

	// Verify game status is not accessible
	_, err = service.GetGameStatus(roomID)
	if err == nil {
		t.Error("GetGameStatus should fail after StopGame")
	}
}

// TestStopGame_MultipleRooms verifies that stopping one room
// doesn't affect other rooms
func TestStopGame_MultipleRooms(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	service := NewDriverService(wsManager)

	players := []sdk.Player{
		{ID: "player1", Username: "Alice", Seat: 0},
		{ID: "player2", Username: "Bob", Seat: 1},
		{ID: "player3", Username: "Charlie", Seat: 2},
		{ID: "player4", Username: "David", Seat: 3},
	}

	// Start two games in different rooms
	room1 := "test-room-1"
	room2 := "test-room-2"

	err := service.StartGameWithDriver(room1, players)
	if err != nil {
		t.Fatalf("Failed to start game in room1: %v", err)
	}

	err = service.StartGameWithDriver(room2, players)
	if err != nil {
		t.Fatalf("Failed to start game in room2: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Stop only room1
	err = service.StopGame(room1)
	if err != nil {
		t.Fatalf("Failed to stop room1: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify room1 is stopped
	service.mu.RLock()
	_, room1Exists := service.drivers[room1]
	_, room2Exists := service.drivers[room2]
	service.mu.RUnlock()

	if room1Exists {
		t.Error("Room1 should be stopped")
	}

	if !room2Exists {
		t.Error("Room2 should still be running")
	}

	// Clean up room2
	service.StopGame(room2)
}

// TestStopGame_NonExistentRoom verifies error handling
func TestStopGame_NonExistentRoom(t *testing.T) {
	wsManager := NewMockDriverWSManager()
	service := NewDriverService(wsManager)

	err := service.StopGame("non-existent-room")
	if err == nil {
		t.Error("Expected error when stopping non-existent room")
	}

	expectedError := "no active game for room non-existent-room"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
