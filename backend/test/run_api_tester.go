//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"guandan-world/backend/test"
)

func main() {
	// Parse command line flags
	var (
		serverURL = flag.String("server", "localhost:8080", "Server URL")
		authToken = flag.String("token", "", "Authentication token")
		roomID    = flag.String("room", "", "Room ID (auto-generated if empty)")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
		help      = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help || *authToken == "" {
		fmt.Println("Usage: go run run_api_test.go -token <auth_token> [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  go run run_api_test.go -token abc123 -verbose")
		os.Exit(0)
	}

	// Create tester
	tester := test.NewAPIGameTester(*serverURL, *authToken, *verbose)
	defer tester.Close()

	// Use provided room ID or generate one
	if *roomID == "" {
		*roomID = fmt.Sprintf("test-room-%d", time.Now().Unix())
	}

	log.Printf("Starting API game test...")
	log.Printf("Server: %s", *serverURL)
	log.Printf("Room ID: %s", *roomID)
	log.Printf("Verbose: %v", *verbose)

	// Start game
	if err := tester.StartGame(*roomID); err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}

	log.Println("Game started successfully, waiting for completion...")

	// Run until complete
	if err := tester.RunUntilComplete(); err != nil {
		log.Fatalf("Game failed: %v", err)
	}

	// Print summary
	eventLog := tester.GetEventLog()
	log.Printf("\n=== Game Completed Successfully ===")
	log.Printf("Total events: %d", len(eventLog))

	if !*verbose && len(eventLog) > 0 {
		log.Println("\nLast 5 events:")
		start := len(eventLog) - 5
		if start < 0 {
			start = 0
		}
		for i := start; i < len(eventLog); i++ {
			log.Printf("  - %s", eventLog[i])
		}
	}
}
