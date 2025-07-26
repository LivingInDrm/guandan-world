package test

import (
	"testing"
)

// TestAPIGameTesterCreation tests creating an API game tester
func TestAPIGameTesterCreation(t *testing.T) {
	tester := NewAPIGameTester("localhost:8080", "test-token", false)
	
	if tester == nil {
		t.Fatal("Failed to create APIGameTester")
	}
	
	if tester.serverURL != "localhost:8080" {
		t.Errorf("Expected serverURL to be localhost:8080, got %s", tester.serverURL)
	}
	
	if tester.authToken != "test-token" {
		t.Errorf("Expected authToken to be test-token, got %s", tester.authToken)
	}
	
	if tester.verbose != false {
		t.Error("Expected verbose to be false")
	}
	
	if len(tester.aiAlgorithms) != 0 {
		t.Errorf("Expected empty AI algorithms map, got %d", len(tester.aiAlgorithms))
	}
	
	tester.Close()
}

// TestEventObserverFunctionality tests the event observer
func TestEventObserverFunctionality(t *testing.T) {
	observer := NewTestEventObserver(true)
	
	if observer == nil {
		t.Fatal("Failed to create TestEventObserver")
	}
	
	// Test various event types
	testEvents := []struct {
		eventType string
		data      map[string]interface{}
	}{
		{
			eventType: "match_started",
			data:      map[string]interface{}{},
		},
		{
			eventType: "deal_started",
			data: map[string]interface{}{
				"event_data": map[string]interface{}{
					"deal_level": float64(2),
				},
			},
		},
		{
			eventType: "player_played",
			data: map[string]interface{}{
				"event_data": map[string]interface{}{
					"player_seat": float64(0),
					"cards":       []interface{}{},
				},
			},
		},
	}
	
	// Test that events don't panic
	for _, test := range testEvents {
		observer.OnGameEvent(test.eventType, test.data)
	}
}

// TestParseCard tests card parsing
func TestParseCard(t *testing.T) {
	tester := NewAPIGameTester("localhost:8080", "test-token", false)
	defer tester.Close()
	
	cardData := map[string]interface{}{
		"number": float64(3),
		"color":  "Spade",
	}
	
	card := tester.parseCard(cardData)
	if card == nil {
		t.Fatal("Failed to parse card")
	}
	
	if card.Number != 3 {
		t.Errorf("Expected card number 3, got %d", card.Number)
	}
	
	if card.Color != "Spade" {
		t.Errorf("Expected card color Spade, got %s", card.Color)
	}
}

// TestGameActive tests game active state management
func TestGameActive(t *testing.T) {
	tester := NewAPIGameTester("localhost:8080", "test-token", false)
	defer tester.Close()
	
	// Initially should not be active
	if tester.isGameActive() {
		t.Error("Game should not be active initially")
	}
	
	// Set active
	tester.setGameActive(true)
	if !tester.isGameActive() {
		t.Error("Game should be active after setting to true")
	}
	
	// Set inactive
	tester.setGameActive(false)
	if tester.isGameActive() {
		t.Error("Game should not be active after setting to false")
	}
}

// TestEventLog tests event logging
func TestEventLog(t *testing.T) {
	tester := NewAPIGameTester("localhost:8080", "test-token", false)
	defer tester.Close()
	
	// Log some events
	tester.log("Event 1")
	tester.log("Event 2")
	tester.log("Event 3")
	
	// Get event log
	eventLog := tester.GetEventLog()
	if len(eventLog) != 3 {
		t.Errorf("Expected 3 events in log, got %d", len(eventLog))
	}
	
	if eventLog[0] != "Event 1" {
		t.Errorf("Expected first event to be 'Event 1', got '%s'", eventLog[0])
	}
}

// TestRunAPIGameTest tests the convenience function
func TestRunAPIGameTestStructure(t *testing.T) {
	// This test just verifies the function exists and has correct signature
	// We can't actually run it without a backend
	fn := RunAPIGameTest
	if fn == nil {
		t.Error("RunAPIGameTest function not found")
	}
}

// BenchmarkAPIGameTester benchmarks tester creation
func BenchmarkAPIGameTester(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tester := NewAPIGameTester("localhost:8080", "test-token", false)
		tester.Close()
	}
}