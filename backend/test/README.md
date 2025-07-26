# Backend API Testing

This directory contains tools for testing the Guandan backend API, particularly the Game Driver API implementation.

## APIGameTester

The `APIGameTester` is a simplified client that simulates frontend requests to test the backend API. It:

- Connects via WebSocket for real-time communication
- Handles game action requests (play decisions, tribute selection, etc.)
- Uses AI algorithms to make intelligent game decisions
- Provides event logging and observation

## Usage

### Quick Start

```bash
# First, ensure the backend is running:
cd ../
go run main.go

# In another terminal, run the test:
cd test
go run run_api_test.go -token <auth_token> -verbose
```

### Obtaining an Auth Token

You need to authenticate first to get a token:

```bash
# Register a test user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'

# Login to get token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

### Running Tests

1. **Simple test with verbose output:**
```bash
go run run_api_test.go -token <your_token> -verbose
```

2. **Test with specific room ID:**
```bash
go run run_api_test.go -token <your_token> -room test-room-123
```

3. **Programmatic usage:**
```go
import "guandan-world/backend/test"

// Simple usage
err := test.RunAPIGameTest("localhost:8080", authToken, true)

// Advanced usage
tester := test.NewAPIGameTester("localhost:8080", authToken, true)
defer tester.Close()

err := tester.StartGame("my-room-id")
if err != nil {
    // Handle error
}

err = tester.RunUntilComplete()
```

## Architecture

```
APIGameTester
├── HTTP Client (for API calls)
├── WebSocket Connection (for real-time updates)
├── AI Algorithms (for game decisions)
└── Event Observer (for logging)
```

The tester simulates all 4 players through a single WebSocket connection, making it easier to test and debug compared to managing multiple connections.

## Key Features

1. **Synchronous Action Handling**: Processes action requests one at a time to avoid race conditions
2. **AI Integration**: Uses SmartAutoPlayAlgorithm for intelligent gameplay
3. **Event Logging**: Captures all game events for debugging
4. **Simplified Design**: Single connection manages all players

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Ensure you have a valid auth token
   - Check that the backend auth service is running

2. **WebSocket Connection Failed**
   - Verify the backend is running on the correct port
   - Check that WebSocket upgrade is working

3. **Game Doesn't Start**
   - Ensure all required services are registered in backend/main.go
   - Check that the room service and driver service are initialized

### Debug Tips

- Use `-verbose` flag to see detailed event logs
- Check backend logs for server-side errors
- Ensure all 4 players are properly initialized
- Verify WebSocket messages are being sent/received

## Example Output

```
Starting API game test...
Server: localhost:8080
Room ID: test-room-1234567890
Verbose: true
[APIGameTester] WebSocket connected
[APIGameTester] Game started successfully
=== Match Started ===
=== Deal Started ===
Deal Level: 2
Player 0 played 1 cards
Player 1 played 1 cards
...
=== Match Ended ===
Winner: Team 0

=== Game Completed Successfully ===
Total events: 245
```