# Guandan World API Documentation

## Overview

Guandan World (掼蛋世界) is a Chinese card game platform that provides RESTful APIs and WebSocket connections for real-time game play. This document describes all available endpoints for frontend development.

## Base Information

- **Base URL**: `http://localhost:8080`
- **API Prefix**: `/api`
- **WebSocket URL**: `ws://localhost:8080/ws`
- **Authentication**: JWT Bearer Token
- **Content-Type**: `application/json`

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## API Endpoints

### 1. Authentication APIs

#### 1.1 Register
- **Endpoint**: `POST /api/auth/register`
- **Description**: Register a new user account
- **Authentication**: Not required
- **Request Body**:
```json
{
  "username": "string",  // Required, unique username
  "password": "string"   // Required, password
}
```
- **Success Response** (201):
```json
{
  "user": {
    "id": "user_1234567890",
    "username": "testuser",
    "status": "online",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-01T00:00:00Z"
  },
  "token": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": "user_1234567890",
    "expires_at": "2024-01-02T00:00:00Z"
  }
}
```
- **Error Responses**:
  - 400: Invalid request format
  - 409: Username already exists

#### 1.2 Login
- **Endpoint**: `POST /api/auth/login`
- **Description**: Login with username and password
- **Authentication**: Not required
- **Request Body**:
```json
{
  "username": "string",  // Required
  "password": "string"   // Required
}
```
- **Success Response** (200): Same as register response
- **Error Responses**:
  - 400: Invalid request format
  - 401: Authentication failed (wrong username/password)

#### 1.3 Logout
- **Endpoint**: `POST /api/auth/logout`
- **Description**: Logout current user
- **Authentication**: Required
- **Success Response** (200):
```json
{
  "message": "Successfully logged out"
}
```
- **Error Responses**:
  - 400: Missing or invalid token
  - 401: Unauthorized

#### 1.4 Get Current User
- **Endpoint**: `GET /api/auth/me`
- **Description**: Get current authenticated user information
- **Authentication**: Required
- **Success Response** (200):
```json
{
  "user": {
    "id": "user_1234567890",
    "username": "testuser",
    "status": "online",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-01T00:00:00Z"
  }
}
```
- **Error Responses**:
  - 401: Unauthorized

### 2. Room Management APIs

#### 2.1 Create Room
- **Endpoint**: `POST /api/rooms/create`
- **Description**: Create a new game room (creator becomes owner)
- **Authentication**: Required
- **Request Body**: None
- **Success Response** (201):
```json
{
  "room": {
    "id": "room_1234567890",
    "owner_id": "user_1234567890",
    "status": "waiting",  // waiting | ready | playing
    "created_at": "2024-01-01T00:00:00Z",
    "players": [
      {
        "id": "user_1234567890",
        "username": "testuser",
        "status": "online",
        "is_owner": true,
        "seat": 0
      }
    ]
  }
}
```
- **Error Responses**:
  - 401: Unauthorized
  - 409: Player is already in a room

#### 2.2 Join Room
- **Endpoint**: `POST /api/rooms/join`
- **Description**: Join an existing room
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890"  // Required
}
```
- **Success Response** (200): Returns updated room object
- **Error Responses**:
  - 400: Invalid request
  - 404: Room not found
  - 409: Room is full or not accepting new players

#### 2.3 Leave Room
- **Endpoint**: `POST /api/rooms/leave`
- **Description**: Leave current room
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890"  // Required
}
```
- **Success Response** (200): Returns updated room object or success message if room was closed
- **Error Responses**:
  - 400: Invalid request
  - 404: Room not found
  - 409: Player is not in this room

#### 2.4 Get Room List
- **Endpoint**: `GET /api/rooms`
- **Description**: Get paginated list of rooms
- **Authentication**: Required
- **Query Parameters**:
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 12, max: 50)
  - `status`: Filter by status (waiting | ready | playing)
- **Success Response** (200):
```json
{
  "rooms": [
    {
      "id": "room_1234567890",
      "owner_id": "user_1234567890",
      "status": "waiting",
      "created_at": "2024-01-01T00:00:00Z",
      "players": [...]
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 12,
    "total": 100,
    "total_pages": 9
  }
}
```

#### 2.5 Get Room Details
- **Endpoint**: `GET /api/rooms/:id`
- **Description**: Get details of a specific room
- **Authentication**: Required
- **Success Response** (200): Returns room object
- **Error Responses**:
  - 404: Room not found

#### 2.6 Start Game
- **Endpoint**: `POST /api/rooms/:id/start`
- **Description**: Start game in a room (owner only)
- **Authentication**: Required
- **Success Response** (200): Returns updated room object with status "playing"
- **Error Responses**:
  - 403: Not room owner
  - 404: Room not found
  - 409: Room not ready (need 4 players)

### 3. Game Control APIs (Driver)

These APIs are used during active gameplay to control the game flow.

#### 3.1 Start Game with Driver
- **Endpoint**: `POST /api/game/driver/start`
- **Description**: Initialize game engine for a room
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890",
  "players": [
    {
      "id": "user_1",
      "username": "Player1",
      "seat": 0
    },
    {
      "id": "user_2",
      "username": "Player2", 
      "seat": 1
    },
    {
      "id": "user_3",
      "username": "Player3",
      "seat": 2
    },
    {
      "id": "user_4",
      "username": "Player4",
      "seat": 3
    }
  ]
}
```
- **Success Response** (200):
```json
{
  "success": true,
  "message": "Game started with driver",
  "room_id": "room_1234567890"
}
```

#### 3.2 Submit Play Decision
- **Endpoint**: `POST /api/game/driver/play-decision`
- **Description**: Submit a player's decision to play cards or pass
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890",
  "player_seat": 0,  // 0-3
  "action": "play",  // "play" or "pass"
  "card_ids": ["3H", "3D", "3C"]  // Required if action is "play"
}
```
- **Success Response** (200):
```json
{
  "success": true,
  "message": "Play decision submitted"
}
```

#### 3.3 Submit Tribute Selection
- **Endpoint**: `POST /api/game/driver/tribute-select`
- **Description**: Submit tribute card selection (for losers in double-down)
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890",
  "player_seat": 0,
  "card_id": "AH"  // Card ID to give as tribute
}
```

#### 3.4 Submit Return Tribute
- **Endpoint**: `POST /api/game/driver/tribute-return`
- **Description**: Submit return tribute card (winners returning card to losers)
- **Authentication**: Required
- **Request Body**:
```json
{
  "room_id": "room_1234567890",
  "player_seat": 0,
  "card_id": "3H"  // Card ID to return
}
```

#### 3.5 Get Game Status
- **Endpoint**: `GET /api/game/driver/status/:room_id`
- **Description**: Get current game status
- **Authentication**: Required
- **Success Response** (200):
```json
{
  "room_id": "room_1234567890",
  "game_status": "playing",
  "deal_status": "tribute_phase",
  "turn_info": {
    "current_player": 0,
    "trick_leader": 0
  },
  "match_details": {
    "TeamLevels": {
      "0": 2,  // Team 0 at level 2
      "1": 2   // Team 1 at level 2
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 3.6 Stop Game
- **Endpoint**: `POST /api/game/driver/stop/:room_id`
- **Description**: Stop an active game
- **Authentication**: Required
- **Success Response** (200):
```json
{
  "success": true,
  "message": "Game stopped",
  "room_id": "room_1234567890"
}
```

### 4. Health Check

#### 4.1 Health Check
- **Endpoint**: `GET /healthz`
- **Description**: Check if server is running
- **Authentication**: Not required
- **Success Response** (200):
```json
{
  "status": "pong"
}
```

## WebSocket Protocol

### Connection

Connect to WebSocket using the authentication token:
```
ws://localhost:8080/ws?token=<jwt_token>
```

### Message Format

All WebSocket messages follow this format:
```json
{
  "type": "message_type",
  "data": {
    // Message-specific data
  },
  "timestamp": "2024-01-01T00:00:00Z",
  "player_id": "user_1234567890"  // Optional, set by server
}
```

### Message Types

#### Incoming Messages (Client → Server)

1. **Join Room**
```json
{
  "type": "join_room",
  "data": {
    "room_id": "room_1234567890"
  }
}
```

2. **Leave Room**
```json
{
  "type": "leave_room",
  "data": {
    "room_id": "room_1234567890"
  }
}
```

3. **Start Game**
```json
{
  "type": "start_game",
  "data": {
    "room_id": "room_1234567890"
  }
}
```

4. **Play Cards** (Not implemented yet)
```json
{
  "type": "play_cards",
  "data": {
    "cards": ["3H", "3D", "3C"]
  }
}
```

5. **Pass** (Not implemented yet)
```json
{
  "type": "pass",
  "data": {}
}
```

6. **Ping** (Heartbeat)
```json
{
  "type": "ping",
  "data": {}
}
```

#### Outgoing Messages (Server → Client)

1. **Room Update**
```json
{
  "type": "room_update",
  "data": {
    "action": "player_joined",  // player_joined | player_left | game_started
    "room": { /* room object */ },
    "player_id": "user_1234567890"
  }
}
```

2. **Game Event**
```json
{
  "type": "game_event",
  "data": {
    "event_type": "match_started",  // Various game events
    "event_data": { /* event-specific data */ },
    "player_seat": 0,
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

3. **Game Action Request**
```json
{
  "type": "game_action",
  "data": {
    "action_type": "play_decision_required",  // play_decision_required | tribute_selection_required | return_tribute_required
    "player_seat": 0,
    "room_id": "room_1234567890",
    "timeout": 30,  // seconds
    // Additional action-specific data
  }
}
```

4. **Error**
```json
{
  "type": "error",
  "data": {
    "error": "error message",
    "room_id": "room_1234567890"
  }
}
```

5. **Pong** (Heartbeat response)
```json
{
  "type": "pong",
  "data": {}
}
```

### Game Event Types

The following event types are emitted during gameplay:

- `match_started`: Match begins
- `match_ended`: Match completes
- `deal_started`: New deal (round) starts
- `deal_ended`: Deal completes
- `cards_dealt`: Cards distributed to players
- `tribute_phase`: Tribute phase begins
- `tribute_given`: Tribute card given
- `tribute_returned`: Return tribute given
- `tribute_completed`: Tribute phase ends
- `trick_started`: New trick begins
- `player_played`: Player played cards
- `player_passed`: Player passed
- `trick_won`: Trick winner determined
- `player_finished`: Player finished all cards
- `player_timeout`: Player failed to act in time
- `player_disconnect`: Player disconnected
- `player_reconnect`: Player reconnected

## Card Format

Cards are represented by a string ID combining number and suit:
- Numbers: 2-10, J, Q, K, A
- Suits: H (Hearts ♥), D (Diamonds ♦), C (Clubs ♣), S (Spades ♠)
- Jokers: RJ (Red Joker), BJ (Black Joker)

Examples: "3H", "AD", "RJ"

## Error Response Format

All error responses follow this format:
```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

Common error codes:
- `invalid_request`: Request format is invalid
- `unauthorized`: Authentication required or failed
- `forbidden`: Operation not allowed
- `not_found`: Resource not found
- `conflict`: Operation conflicts with current state
- `internal_error`: Server error

## Development Notes

1. **Authentication**: Store the JWT token after login and include it in all protected API calls
2. **WebSocket**: Maintain a persistent WebSocket connection for real-time game updates
3. **Room Management**: Players can only be in one room at a time
4. **Game Flow**: Use REST APIs for game control actions, receive updates via WebSocket
5. **Heartbeat**: Send periodic ping messages to keep WebSocket connection alive
6. **Error Handling**: Always check for error responses and handle appropriately

## Current Implementation Status

✅ **Implemented**:
- All authentication endpoints
- All room management endpoints
- Game driver control endpoints
- WebSocket connection and basic message handling
- Game event broadcasting

🚧 **Partially Implemented**:
- WebSocket game control messages (play_cards, pass, etc.)
- Player-specific WebSocket routing

📋 **Not Implemented**:
- Direct game control via WebSocket
- Game replay
- Statistics and rankings