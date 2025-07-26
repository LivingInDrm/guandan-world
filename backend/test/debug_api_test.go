//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"time"
	"guandan-world/backend/test"
)

func main() {
	fmt.Println("=== Debug API Game Test ===")
	fmt.Println("This test will help debug the connection issue")
	fmt.Println()

	// Use a test token
	authToken := "test-token-debug"
	serverURL := "localhost:8080"
	
	// Create a minimal test
	fmt.Println("Step 1: Creating tester...")
	tester := test.NewAPIGameTester(serverURL, authToken, true)
	defer tester.Close()
	
	roomID := fmt.Sprintf("debug-room-%d", time.Now().Unix())
	fmt.Printf("Step 2: Using room ID: %s\n", roomID)
	
	// Try to connect without actually starting game
	fmt.Println("Step 3: Testing WebSocket connection...")
	
	// Let's see what happens when we try to start
	fmt.Println("Step 4: Attempting to start game...")
	err := tester.StartGame(roomID)
	if err != nil {
		fmt.Printf("Error starting game: %v\n", err)
		return
	}
	
	fmt.Println("Step 5: Game start request sent, waiting for events...")
	
	// Wait a bit to see what happens
	time.Sleep(5 * time.Second)
	
	// Check event log
	eventLog := tester.GetEventLog()
	fmt.Printf("\nReceived %d events:\n", len(eventLog))
	for i, event := range eventLog {
		fmt.Printf("  %d: %s\n", i+1, event)
	}
	
	fmt.Println("\nTest complete.")
}