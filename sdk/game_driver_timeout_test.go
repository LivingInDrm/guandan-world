package sdk

import (
	"sync"
	"testing"
	"time"
)

// TestGameDriverConfig_TimeoutStrategy tests that GameDriverConfig includes TimeoutStrategy
func TestGameDriverConfig_TimeoutStrategy(t *testing.T) {
	config := DefaultGameDriverConfig()

	if config.TimeoutStrategy == nil {
		t.Error("DefaultGameDriverConfig should include a TimeoutStrategy")
	}

	// Verify the strategy is the default implementation
	if _, ok := config.TimeoutStrategy.(*DefaultTimeoutStrategy); !ok {
		t.Errorf("Expected DefaultTimeoutStrategy, got %T", config.TimeoutStrategy)
	}

	// Verify timeout durations are set
	if config.PlayDecisionTimeout != 30*time.Second {
		t.Errorf("PlayDecisionTimeout = %v, want 30s", config.PlayDecisionTimeout)
	}

	if config.TributeTimeout != 20*time.Second {
		t.Errorf("TributeTimeout = %v, want 20s", config.TributeTimeout)
	}
}

// TestGameDriverConfig_CustomTimeoutStrategy tests setting a custom timeout strategy
func TestGameDriverConfig_CustomTimeoutStrategy(t *testing.T) {
	customStrategy := NewDefaultTimeoutStrategy()
	
	config := &GameDriverConfig{
		PlayDecisionTimeout: 15 * time.Second,
		TributeTimeout:      10 * time.Second,
		TimeoutStrategy:     customStrategy,
	}

	if config.TimeoutStrategy != customStrategy {
		t.Error("Custom TimeoutStrategy not set correctly")
	}
}

// TestNewGameDriver_TimeoutFields tests that NewGameDriver initializes timeout-related fields
func TestNewGameDriver_TimeoutFields(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// Verify timeoutStats is initialized
	if driver.timeoutStats == nil {
		t.Error("timeoutStats should be initialized")
	}

	if len(driver.timeoutStats) != 0 {
		t.Errorf("timeoutStats should be empty initially, got %d entries", len(driver.timeoutStats))
	}

	// Verify config has TimeoutStrategy
	if driver.config.TimeoutStrategy == nil {
		t.Error("config.TimeoutStrategy should be set by default")
	}
}

// TestNewGameDriver_WithCustomConfig tests NewGameDriver with custom config
func TestNewGameDriver_WithCustomConfig(t *testing.T) {
	engine := NewGameEngine()
	customStrategy := NewDefaultTimeoutStrategy()
	
	config := &GameDriverConfig{
		PlayDecisionTimeout: 15 * time.Second,
		TributeTimeout:      10 * time.Second,
		TimeoutStrategy:     customStrategy,
	}

	driver := NewGameDriver(engine, config)

	if driver.config.TimeoutStrategy != customStrategy {
		t.Error("Custom TimeoutStrategy not preserved in GameDriver")
	}

	if driver.timeoutStats == nil {
		t.Error("timeoutStats should be initialized even with custom config")
	}
}

// TestNewGameDriver_NilTimeoutStrategy tests that nil TimeoutStrategy is replaced with default
func TestNewGameDriver_NilTimeoutStrategy(t *testing.T) {
	engine := NewGameEngine()
	
	config := &GameDriverConfig{
		PlayDecisionTimeout: 15 * time.Second,
		TributeTimeout:      10 * time.Second,
		TimeoutStrategy:     nil, // Explicitly nil
	}

	driver := NewGameDriver(engine, config)

	if driver.config.TimeoutStrategy == nil {
		t.Error("nil TimeoutStrategy should be replaced with default")
	}

	if _, ok := driver.config.TimeoutStrategy.(*DefaultTimeoutStrategy); !ok {
		t.Errorf("Expected DefaultTimeoutStrategy, got %T", driver.config.TimeoutStrategy)
	}
}

// TestGameDriver_GetTimeoutStats tests the GetTimeoutStats method
func TestGameDriver_GetTimeoutStats(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	stats := driver.GetTimeoutStats()
	if stats == nil {
		t.Error("GetTimeoutStats should not return nil")
	}

	if len(stats) != 0 {
		t.Errorf("Initial stats should be empty, got %d entries", len(stats))
	}

	// Manually add some stats to test retrieval
	driver.timeoutStats[0] = &PlayerTimeoutStats{
		PlayDecisionTimeouts: 2,
		TributeTimeouts:      1,
		TotalTimeouts:        3,
	}

	stats = driver.GetTimeoutStats()
	if len(stats) != 1 {
		t.Errorf("Expected 1 player stat, got %d", len(stats))
	}

	if stats[0].TotalTimeouts != 3 {
		t.Errorf("Expected TotalTimeouts=3, got %d", stats[0].TotalTimeouts)
	}
}

// TestPlayerTimeoutStats_Structure tests the PlayerTimeoutStats structure
func TestPlayerTimeoutStats_Structure(t *testing.T) {
	stats := &PlayerTimeoutStats{
		PlayDecisionTimeouts: 5,
		TributeTimeouts:      3,
		TotalTimeouts:        8,
	}

	if stats.PlayDecisionTimeouts != 5 {
		t.Errorf("PlayDecisionTimeouts = %d, want 5", stats.PlayDecisionTimeouts)
	}

	if stats.TributeTimeouts != 3 {
		t.Errorf("TributeTimeouts = %d, want 3", stats.TributeTimeouts)
	}

	if stats.TotalTimeouts != 8 {
		t.Errorf("TotalTimeouts = %d, want 8", stats.TotalTimeouts)
	}
}

// TestGetTimeoutStats_ReturnsDeepCopy tests that GetTimeoutStats returns a deep copy
func TestGetTimeoutStats_ReturnsDeepCopy(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// Add initial stats
	driver.timeoutStats[0] = &PlayerTimeoutStats{
		PlayDecisionTimeouts: 5,
		TributeTimeouts:      3,
		TotalTimeouts:        8,
	}

	// Get stats
	stats := driver.GetTimeoutStats()

	// Modify the returned stats
	stats[0].TotalTimeouts = 999
	stats[1] = &PlayerTimeoutStats{TotalTimeouts: 100} // Add new entry

	// Verify internal stats are unchanged
	internalStats := driver.GetTimeoutStats()
	if internalStats[0].TotalTimeouts != 8 {
		t.Errorf("Internal stats were modified! Expected TotalTimeouts=8, got %d", internalStats[0].TotalTimeouts)
	}

	if _, exists := internalStats[1]; exists {
		t.Error("New entry added to returned map affected internal state")
	}
}

// TestIncrementPlayDecisionTimeout tests the incrementPlayDecisionTimeout method
func TestIncrementPlayDecisionTimeout(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// Increment for seat 0
	driver.incrementPlayDecisionTimeout(0)

	stats := driver.GetTimeoutStats()
	if stats[0].PlayDecisionTimeouts != 1 {
		t.Errorf("PlayDecisionTimeouts = %d, want 1", stats[0].PlayDecisionTimeouts)
	}
	if stats[0].TotalTimeouts != 1 {
		t.Errorf("TotalTimeouts = %d, want 1", stats[0].TotalTimeouts)
	}
	if stats[0].TributeTimeouts != 0 {
		t.Errorf("TributeTimeouts = %d, want 0", stats[0].TributeTimeouts)
	}

	// Increment again
	driver.incrementPlayDecisionTimeout(0)
	stats = driver.GetTimeoutStats()
	if stats[0].PlayDecisionTimeouts != 2 {
		t.Errorf("PlayDecisionTimeouts = %d, want 2", stats[0].PlayDecisionTimeouts)
	}
	if stats[0].TotalTimeouts != 2 {
		t.Errorf("TotalTimeouts = %d, want 2", stats[0].TotalTimeouts)
	}
}

// TestIncrementTributeTimeout tests the incrementTributeTimeout method
func TestIncrementTributeTimeout(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	// Increment for seat 2
	driver.incrementTributeTimeout(2)

	stats := driver.GetTimeoutStats()
	if stats[2].TributeTimeouts != 1 {
		t.Errorf("TributeTimeouts = %d, want 1", stats[2].TributeTimeouts)
	}
	if stats[2].TotalTimeouts != 1 {
		t.Errorf("TotalTimeouts = %d, want 1", stats[2].TotalTimeouts)
	}
	if stats[2].PlayDecisionTimeouts != 0 {
		t.Errorf("PlayDecisionTimeouts = %d, want 0", stats[2].PlayDecisionTimeouts)
	}

	// Increment again
	driver.incrementTributeTimeout(2)
	stats = driver.GetTimeoutStats()
	if stats[2].TributeTimeouts != 2 {
		t.Errorf("TributeTimeouts = %d, want 2", stats[2].TributeTimeouts)
	}
	if stats[2].TotalTimeouts != 2 {
		t.Errorf("TotalTimeouts = %d, want 2", stats[2].TotalTimeouts)
	}
}

// TestTimeoutStats_ConcurrentAccess tests thread safety of timeout stats
func TestTimeoutStats_ConcurrentAccess(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	const numGoroutines = 100
	const incrementsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // *2 for both play and tribute increments

	// Concurrently increment play decision timeouts
	for i := 0; i < numGoroutines; i++ {
		go func(seat int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				driver.incrementPlayDecisionTimeout(seat % 4) // Distribute across 4 seats
			}
		}(i)
	}

	// Concurrently increment tribute timeouts
	for i := 0; i < numGoroutines; i++ {
		go func(seat int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				driver.incrementTributeTimeout(seat % 4) // Distribute across 4 seats
			}
		}(i)
	}

	wg.Wait()

	// Verify total counts
	stats := driver.GetTimeoutStats()
	totalPlay := 0
	totalTribute := 0
	totalAll := 0

	for seat := 0; seat < 4; seat++ {
		if s, ok := stats[seat]; ok {
			totalPlay += s.PlayDecisionTimeouts
			totalTribute += s.TributeTimeouts
			totalAll += s.TotalTimeouts
		}
	}

	expectedTotal := numGoroutines * incrementsPerGoroutine
	if totalPlay != expectedTotal {
		t.Errorf("PlayDecisionTimeouts total = %d, want %d", totalPlay, expectedTotal)
	}
	if totalTribute != expectedTotal {
		t.Errorf("TributeTimeouts total = %d, want %d", totalTribute, expectedTotal)
	}
	if totalAll != expectedTotal*2 {
		t.Errorf("TotalTimeouts = %d, want %d", totalAll, expectedTotal*2)
	}
}

// TestTimeoutStats_ConcurrentReadWrite tests concurrent reads and writes
func TestTimeoutStats_ConcurrentReadWrite(t *testing.T) {
	engine := NewGameEngine()
	driver := NewGameDriver(engine, nil)

	const duration = 100 * time.Millisecond
	done := make(chan bool)

	// Writer goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				driver.incrementPlayDecisionTimeout(0)
			case <-done:
				return
			}
		}
	}()

	// Reader goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(1 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = driver.GetTimeoutStats() // Just read, don't check values
				case <-done:
					return
				}
			}
		}()
	}

	// Run for a duration
	time.Sleep(duration)
	close(done)
	wg.Wait()

	// If we get here without data race, test passes
	stats := driver.GetTimeoutStats()
	if stats[0].PlayDecisionTimeouts == 0 {
		t.Error("Expected some play decision timeouts to be recorded")
	}
}
