package test

import (
	"fmt"
	"log"
	"time"
)

// ExampleAPIGameTest demonstrates how to use the APIGameTester
func ExampleAPIGameTest() {
	// Configuration
	serverURL := "localhost:8080"
	authToken := "test-auth-token" // You need to obtain this from authentication
	verbose := true

	// Run the test
	if err := RunAPIGameTest(serverURL, authToken, verbose); err != nil {
		log.Fatalf("API game test failed: %v", err)
	}
}

// ExampleManualAPIGameTest shows manual step-by-step usage
func ExampleManualAPIGameTest() error {
	// Create tester
	tester := NewAPIGameTester("localhost:8080", "test-auth-token", true)
	defer tester.Close()

	// Create unique room ID
	roomID := fmt.Sprintf("test-room-%d", time.Now().Unix())

	// Start game
	log.Printf("Starting game in room %s...", roomID)
	if err := tester.StartGame(roomID); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	// Run until complete
	log.Println("Running game until completion...")
	if err := tester.RunUntilComplete(); err != nil {
		return fmt.Errorf("game failed: %w", err)
	}

	// Get event log
	eventLog := tester.GetEventLog()
	log.Printf("Game completed with %d events", len(eventLog))

	// Print some events
	log.Println("First 10 events:")
	for i, event := range eventLog {
		if i >= 10 {
			break
		}
		log.Printf("  %d: %s", i+1, event)
	}

	return nil
}

// ExampleIntegrationTest shows how to integrate with actual backend
func ExampleIntegrationTest() error {
	// This example assumes the backend is running with:
	// - Auth service enabled
	// - Room service enabled
	// - Game driver service enabled
	// - WebSocket manager running

	// Step 1: Authenticate and get token
	// In a real test, you would call the auth API to get a token
	// authToken := authenticateTestUser()

	// Step 2: Create a test room
	// roomID := createTestRoom(authToken)

	// Step 3: Run the game test
	// return RunAPIGameTest("localhost:8080", authToken, true)

	// For now, just return nil
	log.Println("Integration test example - implement with actual backend")
	return nil
}

// TestMultipleGames shows how to run multiple concurrent games
func TestMultipleGames() error {
	numGames := 3
	results := make(chan error, numGames)

	// Start multiple games concurrently
	for i := 0; i < numGames; i++ {
		go func(gameNum int) {
			serverURL := "localhost:8080"
			authToken := fmt.Sprintf("test-token-%d", gameNum)
			verbose := false // Less verbose for multiple games

			log.Printf("Starting game %d...", gameNum)
			err := RunAPIGameTest(serverURL, authToken, verbose)
			if err != nil {
				log.Printf("Game %d failed: %v", gameNum, err)
			} else {
				log.Printf("Game %d completed successfully", gameNum)
			}
			results <- err
		}(i)
	}

	// Wait for all games to complete
	var firstError error
	for i := 0; i < numGames; i++ {
		if err := <-results; err != nil && firstError == nil {
			firstError = err
		}
	}

	if firstError != nil {
		return fmt.Errorf("at least one game failed: %w", firstError)
	}

	log.Printf("All %d games completed successfully", numGames)
	return nil
}

// main function for testing
func main() {
	// Example 1: Simple API game test
	fmt.Println("=== Example 1: Simple API Game Test ===")
	ExampleAPIGameTest()

	// Example 2: Manual step-by-step test
	fmt.Println("\n=== Example 2: Manual API Game Test ===")
	if err := ExampleManualAPIGameTest(); err != nil {
		log.Printf("Manual test failed: %v", err)
	}

	// Example 3: Multiple concurrent games
	fmt.Println("\n=== Example 3: Multiple Concurrent Games ===")
	if err := TestMultipleGames(); err != nil {
		log.Printf("Multiple games test failed: %v", err)
	}
}