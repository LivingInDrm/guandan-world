import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import GameBoard from './GameBoard';
import { GameState, Player, TrickInfo, TrickStatus, PlayerStatus } from '../../types';

// Mock data for testing
const mockPlayers: (Player | null)[] = [
  { id: '1', username: 'Player1', seat: 0, online: true, auto_play: false },
  { id: '2', username: 'Player2', seat: 1, online: true, auto_play: false },
  { id: '3', username: 'Player3', seat: 2, online: true, auto_play: false },
  { id: '4', username: 'Player4', seat: 3, online: true, auto_play: false },
];

const mockGameState: GameState = {
  match_id: 'match1',
  current_deal: {
    id: 'deal1',
    level: 5,
    status: 'playing',
    current_trick: null,
    player_cards: [[], [], [], []],
    rankings: [],
    start_time: '2024-01-01T00:00:00Z',
  },
  team_levels: [3, 4],
  current_player: 0,
  status: 'playing',
};

const mockTrickInfo: TrickInfo = {
  id: 'trick1',
  leader: 0,
  current_turn: 1,
  plays: [
    {
      player_seat: 0,
      cards: [{ id: 'card1', suit: 0, rank: 5, is_joker: false }],
      is_pass: false,
      timestamp: '2024-01-01T00:00:00Z',
    },
  ],
  winner: -1,
  lead_comp: null,
  status: TrickStatus.PLAYING,
  start_time: '2024-01-01T00:00:00Z',
  turn_timeout: '2024-01-01T00:00:20Z',
  next_leader: 0,
};

describe('GameBoard', () => {
  it('renders game board with all player areas', () => {
    render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={null}
      />
    );

    // Check if all players are rendered
    expect(screen.getByText('Player1')).toBeInTheDocument();
    expect(screen.getByText('Player2')).toBeInTheDocument();
    expect(screen.getByText('Player3')).toBeInTheDocument();
    expect(screen.getByText('Player4')).toBeInTheDocument();

    // Check if play area is rendered
    expect(screen.getByText('出牌区')).toBeInTheDocument();
  });

  it('displays team levels correctly', () => {
    render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={null}
      />
    );

    // Check team level display
    expect(screen.getByText('等级信息')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument(); // Team 0 level
    expect(screen.getByText('4')).toBeInTheDocument(); // Team 1 level
    expect(screen.getByText('5')).toBeInTheDocument(); // Current deal level
  });

  it('shows correct player status when trick is active', () => {
    render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={mockTrickInfo}
      />
    );

    // Player 0 has played
    expect(screen.getByText('已出牌')).toBeInTheDocument();
    
    // Player 1 is current turn
    expect(screen.getByText('出牌中')).toBeInTheDocument();
  });

  it('handles empty player slots', () => {
    const playersWithEmpty: (Player | null)[] = [
      mockPlayers[0],
      null,
      mockPlayers[2],
      mockPlayers[3],
    ];

    render(
      <GameBoard
        gameState={mockGameState}
        players={playersWithEmpty}
        currentPlayerSeat={0}
        trickInfo={null}
      />
    );

    expect(screen.getByText('空座位')).toBeInTheDocument();
  });

  it('highlights current turn player', () => {
    const { container } = render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={mockTrickInfo}
      />
    );

    // Check if current turn player has highlight border
    const playerAreas = container.querySelectorAll('.border-yellow-400');
    expect(playerAreas.length).toBeGreaterThan(0);
  });

  it('displays trick information when available', () => {
    render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={mockTrickInfo}
      />
    );

    // Check if trick progress is shown
    expect(screen.getByText('当前轮次: 1/4')).toBeInTheDocument();
  });

  it('converts high card levels to letters correctly', () => {
    const gameStateWithHighLevel: GameState = {
      ...mockGameState,
      team_levels: [11, 14], // J and A
      current_deal: {
        ...mockGameState.current_deal,
        level: 13, // K
      },
    };

    render(
      <GameBoard
        gameState={gameStateWithHighLevel}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={null}
      />
    );

    expect(screen.getByText('J')).toBeInTheDocument(); // Team 0 level
    expect(screen.getByText('A')).toBeInTheDocument(); // Team 1 level
    expect(screen.getByText('K')).toBeInTheDocument(); // Current deal level
  });

  it('shows last play information', () => {
    render(
      <GameBoard
        gameState={mockGameState}
        players={mockPlayers}
        currentPlayerSeat={0}
        trickInfo={mockTrickInfo}
      />
    );

    // Check if last play info is shown for player who played
    expect(screen.getByText('出牌: 1张')).toBeInTheDocument();
  });
});