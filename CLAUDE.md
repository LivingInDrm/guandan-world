# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Guandan World (掼蛋世界) is a Chinese card game platform implementing the game of Guandan. The project uses a microservices architecture with Go backend and React frontend.

## Common Development Commands

### Quick Start
```bash
# Start all services with Docker Compose
docker-compose up --build

# Services will be available at:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```

### Backend Development
```bash
cd backend
go mod download
go run main.go              # Run server
go test ./...               # Run tests
go test ./... -v            # Verbose test output
go build -v ./...           # Build backend
```

### Frontend Development
```bash
cd frontend
npm install                 # Install dependencies
npm run dev                 # Start dev server (http://localhost:5173)
npm run build               # Production build
npm run lint                # Run linter
npm run preview             # Preview production build
```

### SDK Testing
```bash
cd sdk
go test ./...               # Run all SDK tests
go test ./... -v            # Verbose output
go test -run TestCardCreation  # Run specific test
go test ./... -cover        # With coverage
```

### AI Testing
```bash
cd ai
go test ./...               # Run AI tests
go test -run TestSmartAI -v # Run specific AI tests
```

### Match Simulation
```bash
# From project root
go run cmd/simulate/main.go    # Verbose mode with SmartAutoPlayAlgorithm
go run cmd/simulate/main.go -q # Quiet mode
```

### Proto Generation
```bash
make proto        # Generate all proto files (Go + frontend)
make proto-go     # Generate Go proto files only
make proto-ts     # Generate TypeScript proto files
make clean-proto  # Clean generated proto files
```

## High-Level Architecture

### Project Structure
```
guandan-world/
├── sdk/                  # Core game logic (pure Go, no external dependencies)
│   ├── game_engine.go    # Main game orchestration, event handling
│   ├── game_driver.go    # Game loop driver with timeout management
│   ├── match.go          # Multi-round match management
│   ├── deal.go           # Individual round logic
│   ├── comp.go           # Card combination recognition (base types)
│   ├── comp_*.go         # Card combo implementations (straight, fullhouse, plate, tube)
│   ├── tribute.go        # Tribute phase logic
│   ├── event.go          # Event types and observers (uses proto types)
│   ├── event_converter.go # SDK to proto type conversion
│   └── types.go          # Core type definitions
├── backend/              # API server (Gin framework)
│   ├── auth/             # JWT authentication service
│   ├── room/             # Room management service
│   ├── game/             # Game driver service
│   ├── handlers/         # HTTP/WebSocket handlers
│   ├── websocket/        # WebSocket manager
│   └── main.go           # Server entry point
├── frontend/             # React + TypeScript + Vite
│   └── src/
│       ├── components/   # React components
│       ├── services/     # API services
│       └── store/        # Zustand state management
├── ai/                   # AI player implementations
│   ├── smart_algorithm.go      # SmartAutoPlayAlgorithm
│   └── auto_play_algorithm.go  # Algorithm interface
├── simulator/            # Game simulation tools
│   ├── match_simulator_v2.go   # Main simulator
│   └── simulating_input_provider.go
├── proto/                # Protocol Buffer definitions
│   ├── common.proto      # Common types
│   ├── event.proto       # Game events
│   └── view.proto        # Player views
└── cmd/simulate/         # Simulation entry point
```

### Key SDK Components

#### Core Game Engine (`game_engine.go`)
- Event-driven architecture with `GameEngineInterface`
- Thread-safe with RWMutex for concurrent access
- Emits events through observer pattern (`RegisterObserver`, `On`)
- All game state changes are atomic and emit appropriate events

#### Game Driver (`game_driver.go`)
- Orchestrates game loop with timeout management
- Uses `InputProvider` interface for player decisions
- Handles auto-play for disconnected players
- `RunMatch()` runs complete match, `RunDeal()` runs single deal

#### Card System
- **card.go**: Card representation with ranks (2-A + Jokers) and suits
- **comp.go**: Base CardComp interface and CompType enum
- **comp_straight.go**: Straight combinations (顺子)
- **comp_fullhouse.go**: Full house combinations (三带二)
- **comp_plate.go**: Plate combinations (钢板/连对)
- **comp_tube.go**: Tube combinations (管子/三连)
- **Wildcard support**: Red heart cards of current level act as wildcards

#### Game Flow
1. **Match**: Manages complete game session until team reaches Ace
2. **Deal**: Individual round with phases: dealing → tribute → playing → settlement
3. **Trick**: Single play sequence where each player plays/passes

### Event System

The SDK uses Protocol Buffers for event types (defined in `proto/event.proto`):
```go
// Key event types
EventMatchStarted, EventMatchEnded
EventDealStarted, EventCardsDealt, EventDealEnded
EventTributeStarted, EventTributeExempted, EventTributeCompleted
EventTrickStarted, EventTrickEnded
EventPlayerPlayed, EventPlayerPassed
EventPlayerTimeout, EventPlayerDisconnect, EventPlayerReconnect
```

### API Architecture

#### Authentication
- JWT-based with `/api/auth/register`, `/api/auth/login`, `/api/auth/logout`
- Token validation middleware for protected routes

#### Room Management
- CRUD operations at `/api/rooms/*`
- WebSocket at `/ws` for real-time game updates

#### Game Driver API
- `/api/game/driver/start` - Start game with driver
- `/api/game/driver/play-decision` - Submit play decision
- `/api/game/driver/tribute-*` - Tribute phase actions

### AI System

**SmartAutoPlayAlgorithm** (`ai/smart_algorithm.go`):
- Implements `AutoPlayAlgorithm` interface
- Analyzes hand strength and card combinations
- Uses scoring system for card selection
- Supports both leading and following plays

### Go Workspace

The project uses Go workspace (`go.work`) to manage multiple modules:
```
guandan-world/
├── go.work           # Workspace definition
├── sdk/go.mod        # SDK module
├── backend/go.mod    # Backend module (depends on SDK)
├── ai/go.mod         # AI module (depends on SDK)
└── simulator/go.mod  # Simulator module (depends on SDK, AI)
```

## Important Notes

- The SDK is pure Go with minimal external dependencies for maximum portability
- Proto files must be regenerated after changes: `make proto`
- Frontend automatically runs `make proto-ts` before dev/build
- Game state is event-sourced - all changes emit events
- Card rankings change based on current level (2-A)
- Tribute phase has complex rules - see `tribute.go` for implementation
