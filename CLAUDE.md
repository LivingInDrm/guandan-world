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
go run simulate_match.go    # Verbose mode with SmartAutoPlayAlgorithm
go run simulate_match.go -q # Quiet mode
```

## High-Level Architecture

### Project Structure
```
guandan-world/
├── sdk/                  # Core game logic (pure Go, no external dependencies)
│   ├── game_engine.go    # Main game orchestration, event handling
│   ├── match.go          # Multi-round match management
│   ├── deal.go           # Individual round logic
│   ├── trick.go          # Single play sequence
│   ├── card.go           # Card representation
│   ├── comp.go           # Card combination recognition
│   ├── tribute.go        # Tribute phase logic
│   ├── validator.go      # Move validation
│   └── *_test.go         # Comprehensive test coverage
├── backend/              # API server (Gin framework)
│   ├── auth/             # JWT authentication service
│   ├── room/             # Room management service
│   ├── handlers/         # HTTP/WebSocket handlers
│   └── main.go           # Server entry point
├── frontend/             # React + TypeScript + Vite
│   └── (basic setup)     # UI to be implemented
├── ai/                   # AI player implementations
│   └── smart_ai.go       # SmartAutoPlayAlgorithm
└── simulator/            # Game simulation tools
    └── simulate_match.go # Match simulator using AI
```

### Key SDK Components

#### Core Game Engine
- **game_engine.go**: Event-driven architecture with GameEngineInterface
- **match.go**: Manages complete game session until team reaches Ace
- **deal.go**: Individual round with phases: dealing → tribute → playing → settlement
- **trick.go**: Single play sequence where each player plays/passes

#### Card System
- **card.go**: Card representation with ranks (2-A + Jokers) and suits
- **comp.go**: Recognizes 10+ card combinations (单张, 对子, 顺子, 炸弹, etc.)
- **validator.go**: Validates moves according to Guandan rules
- **wildcard support**: Red heart cards of current level act as wildcards

#### Game Rules Implementation
- **tribute.go**: Implements tribute exchange between winning/losing teams
- **first_out.go**: Determines first player each round (deal winner or tribute receiver)
- **result.go**: Calculates round results and team rankings
- **Special rules**: 
  - Joker bombs beat everything
  - Straight flushes beat regular bombs
  - Larger bombs beat smaller ones
  - Wildcards can form any combination

### Event System

The SDK emits events for all state changes:
```go
// Match events
EventMatchStarted, EventMatchEnded

// Deal events  
EventDealStarted, EventCardsDealt, EventDealEnded

// Tribute events
EventTributePhase, EventTributeGiven, EventReturnTribute

// Playing events
EventTrickStarted, EventPlayerPlayed, EventPlayerPassed, EventTrickWon

// System events
EventPlayerTimeout, EventPlayerDisconnect, EventPlayerReconnect
```

### API Architecture

#### Authentication (Implemented)
- JWT-based authentication with register/login/logout
- Token validation middleware
- User session management

#### Room Management (Implemented but not registered)
- Create/join/leave rooms
- Room listing with pagination
- Start game when ready

#### Game Protocol
- REST API for room operations
- WebSocket for real-time game updates (planned)
- Event-based state synchronization

### AI System

**SmartAutoPlayAlgorithm** features:
- Analyzes hand strength and card combinations
- Strategic decision making for tribute phase
- Intelligent card selection based on game state
- Supports both leading and following plays

### Testing Strategy

- **SDK**: 100% core rule coverage with 17 test files
- **Backend**: Basic auth tests implemented  
- **AI**: Smart algorithm tests with various scenarios
- **Integration**: simulate_match.go for end-to-end testing

### Current Implementation Status

✅ **Completed**:
- Complete Guandan rule engine with all card types
- Event-driven game architecture
- JWT authentication system
- Room management service
- Smart AI algorithm
- Match simulation system

🚧 **In Progress**:
- WebSocket implementation for real-time play
- Frontend UI development
- Database persistence layer

📋 **Planned**:
- Player statistics and rankings
- Game replay system
- Tournament support
- Mobile app support

## Important Notes

- The SDK is pure Go with no external dependencies for maximum portability
- Backend room routes need registration in main.go to activate
- CORS is currently open for development - restrict for production
- Game state is in-memory only - implement persistence for production
- Tribute phase has complex rules - see tribute.go for full implementation
- Card rankings change based on current level (2-A)