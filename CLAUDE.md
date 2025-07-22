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
npm run dev                 # Start dev server
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

### Match Simulation
```bash
go run simulate_match.go    # Verbose mode
go run simulate_match.go -q # Quiet mode
```

## High-Level Architecture

### Project Structure
- **sdk/**: Core game logic (pure Go, no external dependencies)
  - Event-driven architecture with GameEngineInterface
  - State machine pattern for Match → Deal → Trick flow
  - Complete game rule implementation including tribute system
  
- **backend/**: API server (Gin framework)
  - Currently minimal - health check only
  - Planned: REST API + WebSocket for real-time updates
  
- **frontend/**: React + TypeScript + Vite
  - Basic SPA setup
  - Communicates with backend via REST (WebSocket planned)

### Key SDK Components
- **game_engine.go**: Main game orchestration, event handling
- **match.go**: Multi-round match management (until team reaches Ace)
- **deal.go**: Individual round logic within a match
- **trick.go**: Single play sequence within a deal
- **card.go**: Card representation and operations
- **comp.go**: Card combination validation (pairs, straights, etc.)
- **tribute.go**: Special card exchange phase between deals
- **validator.go**: Move validation logic

### Game Flow
1. **Match**: Complete game session (multiple deals)
2. **Deal**: Individual round with dealing → tribute → playing → settlement
3. **Trick**: Single play sequence where each player plays cards
4. **Events**: State changes trigger events (EventMatchStarted, EventPlayerPlayed, etc.)

### Data Models
- Game state uses event-driven updates
- Clear state transitions: Waiting → Playing → Finished
- No database integration yet - game state managed in memory

### Testing Strategy
- SDK has comprehensive test coverage (17 test files)
- Use standard Go testing framework
- Test files follow *_test.go convention
- No frontend tests implemented yet

## Important Notes
- SDK is designed as pure business logic with no I/O dependencies
- Backend needs WebSocket implementation for real-time gameplay
- Database schema not yet implemented
- Authentication system planned but not implemented
- The codebase follows clean architecture principles