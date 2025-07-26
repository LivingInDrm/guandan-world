# Game Driver API Documentation

## Overview

The Game Driver API provides a complete game management solution using the SDK's GameDriver architecture. This API encapsulates the full game flow including player input handling, event observation, and real-time WebSocket communication.

## Architecture

### Key Components

1. **GameDriver** (SDK): Orchestrates the game loop and manages game flow
2. **DriverService** (Backend): Wraps the GameDriver for HTTP/WebSocket API
3. **PlayerInputProvider**: Handles player decisions via WebSocket
4. **EventObserver**: Broadcasts game events to connected clients

### Design Benefits

- **Complete Game Flow**: The GameDriver manages the entire match lifecycle automatically
- **Real-time Communication**: Players receive action requests and submit decisions via WebSocket
- **Event-driven Updates**: All game state changes are broadcast as events
- **Timeout Handling**: Built-in timeout management for player decisions

## API Endpoints

### Start Game with Driver

```http
POST /api/game/driver/start
Authorization: Bearer <token>
Content-Type: application/json

{
  "room_id": "room-123",
  "players": [
    {"id": "player1", "username": "Alice", "seat": 0},
    {"id": "player2", "username": "Bob", "seat": 1},
    {"id": "player3", "username": "Charlie", "seat": 2},
    {"id": "player4", "username": "David", "seat": 3}
  ]
}

Response:
{
  "success": true,
  "message": "Game started with driver",
  "room_id": "room-123"
}
```

### Submit Play Decision

```http
POST /api/game/driver/play-decision
Authorization: Bearer <token>
Content-Type: application/json

{
  "room_id": "room-123",
  "player_seat": 0,
  "action": "play",  // or "pass"
  "card_ids": ["Heart_5", "Heart_6", "Heart_7"]  // only for "play" action
}

Response:
{
  "success": true,
  "message": "Play decision submitted"
}
```

### Submit Tribute Selection

```http
POST /api/game/driver/tribute-select
Authorization: Bearer <token>
Content-Type: application/json

{
  "room_id": "room-123",
  "player_seat": 0,
  "card_id": "Spade_13"
}

Response:
{
  "success": true,
  "message": "Tribute selection submitted"
}
```

### Submit Return Tribute

```http
POST /api/game/driver/tribute-return
Authorization: Bearer <token>
Content-Type: application/json

{
  "room_id": "room-123",
  "player_seat": 1,
  "card_id": "Diamond_3"
}

Response:
{
  "success": true,
  "message": "Return tribute submitted"
}
```

### Get Game Status

```http
GET /api/game/driver/status/{room_id}
Authorization: Bearer <token>

Response:
{
  "room_id": "room-123",
  "game_status": "started",
  "deal_status": "playing",
  "turn_info": {
    "current_player": 0,
    "is_leader": true,
    "is_new_trick": false,
    "has_active_trick": true,
    "lead_comp": null
  },
  "match_details": {
    "team_levels": [2, 2],
    "players": [
      {"seat": 0, "username": "Alice", "team_num": 0},
      {"seat": 1, "username": "Bob", "team_num": 1},
      {"seat": 2, "username": "Charlie", "team_num": 0},
      {"seat": 3, "username": "David", "team_num": 1}
    ]
  },
  "timestamp": "2024-01-26T10:30:00Z"
}
```

### Stop Game

```http
POST /api/game/driver/stop/{room_id}
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Game stopped",
  "room_id": "room-123"
}
```

## WebSocket Protocol

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=<auth_token>');
```

### Message Types

#### Game Events (Server → Client)

```json
{
  "type": "game_event",
  "data": {
    "event_type": "deal_started",
    "event_data": {
      "deal_level": 2,
      "team0_level": 2,
      "team1_level": 2
    },
    "timestamp": "2024-01-26T10:30:00Z",
    "player_seat": -1
  },
  "timestamp": "2024-01-26T10:30:00Z"
}
```

#### Action Requests (Server → Client)

```json
{
  "type": "game_action",
  "data": {
    "action_type": "play_decision_required",
    "player_seat": 0,
    "hand": [
      {"number": 3, "color": "Heart", "id": "Heart_3"},
      {"number": 4, "color": "Heart", "id": "Heart_4"}
    ],
    "trick_info": {
      "is_leader": true,
      "lead_comp": null
    },
    "timeout": 30,
    "room_id": "room-123"
  },
  "timestamp": "2024-01-26T10:30:00Z"
}
```

### Game Flow

1. **Start Game**: Call `/api/game/driver/start` to initialize the game
2. **Game Events**: Receive real-time events via WebSocket
3. **Action Requests**: When it's a player's turn, they receive an action request
4. **Submit Decision**: Player submits their decision via the appropriate API endpoint
5. **Timeout Handling**: If no decision is received within the timeout, a default action is taken

### Event Types

- `match_started`: Match has begun
- `deal_started`: New deal started
- `tribute_rules_set`: Tribute requirements determined
- `tribute_pool_created`: Double-down tribute pool created
- `tribute_given`: Tribute card given
- `tribute_selected`: Card selected from tribute pool
- `return_tribute`: Return tribute completed
- `tribute_completed`: Tribute phase finished
- `trick_started`: New trick begun
- `player_played`: Player played cards
- `player_passed`: Player passed
- `trick_ended`: Trick completed
- `deal_ended`: Deal finished
- `match_ended`: Match completed

### Action Types

- `play_decision_required`: Player needs to play or pass
- `tribute_selection_required`: Player needs to select from tribute pool (double-down)
- `return_tribute_required`: Player needs to return a tribute card

## Error Handling

All endpoints return appropriate HTTP status codes:

- `200 OK`: Success
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource not found (e.g., game doesn't exist)
- `500 Internal Server Error`: Server error

Error response format:
```json
{
  "error": "Detailed error message"
}
```

## Implementation Example

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws?token=' + authToken);

// Listen for game events
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'game_event') {
    // Handle game state updates
    updateGameUI(message.data);
  } else if (message.type === 'game_action') {
    // Handle action requests
    if (message.data.action_type === 'play_decision_required') {
      showPlayDecisionUI(message.data);
    }
  }
};

// Submit play decision
async function submitPlayDecision(roomId, playerSeat, action, cardIds) {
  const response = await fetch('/api/game/driver/play-decision', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + authToken,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      room_id: roomId,
      player_seat: playerSeat,
      action: action,
      card_ids: cardIds
    })
  });
  
  return response.json();
}
```

## Configuration

The GameDriver uses configurable timeouts:

- `PlayDecisionTimeout`: 30 seconds (default)
- `TributeTimeout`: 20 seconds (default)

These can be adjusted in the backend configuration for different game modes or testing purposes.