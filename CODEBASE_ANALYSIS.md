# Guandan World Codebase Analysis

## Project Overview

**Guandan World (掼蛋世界)** is a full-stack online multiplayer card game platform built with Go and React. It implements the traditional Chinese card game "Guandan" (掼蛋) with real-time multiplayer functionality, featuring complete game mechanics, user authentication, and room management.

## Architecture Overview

The project follows a modular microservices architecture with clear separation of concerns:

```
Frontend (React/TypeScript) ← WebSocket/HTTP → Backend (Go/Gin) ← → SDK (Game Logic)
                                                       ↓
                                               Database (PostgreSQL) + Cache (Redis)
```

## Directory Structure

```
guandan-world/
├── frontend/           # React TypeScript frontend application
├── backend/           # Go backend API server with WebSocket support
├── sdk/              # Core game logic and rules engine
├── ai/               # AI algorithms for automated gameplay
├── simulator/        # Game simulation and testing tools
├── cmd/              # Command-line utilities
├── docs/             # Game rules and design documentation
├── test-data/        # Test datasets and scenarios
└── Integration files # Docker configs, scripts, and documentation
```

## Core Modules

### 1. Frontend (`/frontend/`)

**Technology Stack:**
- React 19.1.0 with TypeScript
- Vite for build tooling
- ESLint for code quality

**Key Features:**
- Health check integration with backend
- Modern React development setup
- Production-ready build pipeline

**Dependencies:**
```json
{
  "react": "^19.1.0",
  "react-dom": "^19.1.0",
  "@vitejs/plugin-react": "^4.6.0",
  "typescript": "~5.8.3",
  "vite": "^7.0.4"
}
```

### 2. Backend (`/backend/`)

**Technology Stack:**
- Go 1.23+ with Gin web framework
- JWT authentication
- WebSocket support via Gorilla WebSocket
- Comprehensive testing with testify

**Core Services:**
- **Authentication Service**: JWT-based user auth with 24-hour token expiration
- **Room Service**: Game room management and player coordination
- **WebSocket Manager**: Real-time communication hub
- **Game Driver Service**: Game state management and flow control

**API Endpoints:**
- `GET /healthz` - Health check
- `POST /auth/*` - Authentication routes
- `GET /ws` - WebSocket connection with token validation
- `POST /api/rooms/*` - Room management (create, join, leave, start)
- `POST /api/game/driver/*` - Game control (start, play, tribute, status)

**Key Dependencies:**
```go
require (
    github.com/gin-gonic/gin v1.10.1
    github.com/golang-jwt/jwt/v5 v5.2.3
    github.com/gorilla/websocket v1.5.3
    golang.org/x/crypto v0.40.0
)
```

### 3. SDK (`/sdk/`) - Game Logic Engine

**Core Components:**

#### 3.1 Card System (`card.go`)
- **Card Representation**: Numbers 1-16 (1=Ace, 11=Jack, 12=Queen, 13=King, 14=Ace, 15=Black Joker, 16=Red Joker)
- **Special Rules**: Level cards, wildcard logic (red heart level cards)
- **Features**: Card comparison, cloning, JSON serialization, string representation

#### 3.2 Card Combination System (`comp.go`)
**Supported Card Types:**
- **Basic Types**: Single, Pair, Triple, Full House
- **Sequence Types**: Straight (5+ cards), Plate (6+ consecutive triples), Tube (6+ consecutive pairs)
- **Bomb Types**: 4-8 card bombs, Joker bombs, Straight flush
- **Smart Recognition**: Automatic card type detection with wildcard handling

**Priority System (High to Low):**
1. Joker Bomb (highest priority)
2. Straight Flush (beats all except Joker Bomb)
3. 6+ Card Bombs (beat Straight Flush)
4. 5-Card Bombs (beat regular types, less than Straight Flush)
5. 4-Card Bombs (beat regular types)
6. Regular Types (same type comparison by value)

#### 3.3 Game Flow Management
- **Match**: Complete game until team reaches Ace level
- **Deal**: Single round with dealing, tribute, playing phases
- **Trick**: One round of card plays from all players
- **Tribute System**: Complex tribute exchange mechanics

**Data Structures:**
```go
type Match struct {
    ID          string
    Status      MatchStatus
    Players     [4]*Player
    CurrentDeal *Deal
    TeamLevels  [2]int    // Team levels
    Winner      int       // Winning team
}

type Deal struct {
    Level        int
    Status       DealStatus
    CurrentTrick *Trick
    TributePhase *TributePhase
    PlayerCards  [4][]*Card
    Rankings     []int
}
```

### 4. AI Module (`/ai/`)

**Components:**
- **Auto Play Algorithm**: Automated gameplay decisions
- **Smart Algorithm**: Advanced AI strategies
- Comprehensive test coverage

### 5. Simulator (`/simulator/`)

**Features:**
- Match simulation capabilities
- Performance testing and validation
- Input provider for automated testing

## Infrastructure & DevOps

### Docker Configuration

**Multi-service setup** via `docker-compose.yml`:
- **Frontend**: React app served on port 3000
- **Backend**: Go API server on port 8080
- **PostgreSQL**: Database on port 5432
- **Redis**: Caching layer on port 6379

**Network**: All services connected via `guandan-network` bridge

### Testing Strategy

**Comprehensive Test Coverage:**
- **Unit Tests**: Card logic, game rules, API endpoints
- **Integration Tests**: Full game flow testing
- **WebSocket Tests**: Real-time communication testing
- **Performance Tests**: Load testing and stress testing

**Test Files Pattern:**
- `*_test.go` for unit tests
- `*_integration_test.go` for integration tests
- Test helpers and utilities in dedicated files

## Game Features

### ✅ Implemented Features
- Real-time multiplayer battles
- Complete Guandan rule implementation
- User authentication system
- Game room management
- Real-time chat functionality
- Game replay functionality
- Complete card system
- Smart card type recognition
- Wildcard handling
- Bomb priority logic

### Game Rules Implementation

**Tribute System**: Complex tribute exchange between teams based on previous deal results
**Level Progression**: Teams advance through levels (2-A) based on finishing order
**Special Cards**: Level cards and wildcards with dynamic behavior
**Team Play**: 4-player game with 2 teams (seats 0,2 vs 1,3)

## Development Workflow

### Local Development
```bash
# Backend
cd backend && go run main.go

# Frontend  
cd frontend && npm install && npm run dev

# Database (optional)
docker-compose up postgres redis
```

### Testing
```bash
# Backend tests
cd backend && go test ./...

# SDK tests
cd sdk && go test ./...

# Frontend tests
cd frontend && npm test
```

### Production Deployment
```bash
# Build and deploy
docker-compose build
docker-compose up -d
```

## Technical Highlights

1. **Modular Architecture**: Clear separation between game logic (SDK), backend services, and frontend
2. **Real-time Communication**: WebSocket-based real-time gameplay
3. **Comprehensive Testing**: Extensive test coverage across all modules
4. **Docker Integration**: Full containerization for consistent deployment
5. **Advanced Game Logic**: Complex card combination recognition and game rule implementation
6. **Scalable Design**: Microservices architecture ready for scaling

## Security Considerations

- JWT-based authentication with configurable expiration
- CORS configuration for cross-origin requests
- Input validation for all API endpoints
- Secure WebSocket connections with token validation

## Performance Features

- Efficient card comparison algorithms
- Optimized game state management
- Redis caching for improved performance
- Asynchronous WebSocket communication

## Documentation Quality

The project includes extensive documentation:
- API documentation for backend services
- Game rules documentation in `/docs/`
- Technical implementation guides
- Integration test reports
- Performance analysis reports

This codebase represents a production-ready, well-architected online card game platform with comprehensive features and robust testing coverage.