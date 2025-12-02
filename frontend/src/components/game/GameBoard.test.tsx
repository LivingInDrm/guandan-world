import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import GameBoard from './GameBoard';
import type { Player } from '../../types';
import type { PlayAction } from '../../types/proto';
import { CompType } from '../../types/proto';

const mockPlayers: (Player | null)[] = [
  { id: '1', username: 'Player1', seat: 0, online: true, auto_play: false },
  { id: '2', username: 'Player2', seat: 1, online: true, auto_play: false },
  { id: '3', username: 'Player3', seat: 2, online: true, auto_play: false },
  { id: '4', username: 'Player4', seat: 3, online: true, auto_play: false },
];

const mockPlays: PlayAction[] = [
  {
    playerSeat: 0,
    cards: [{ suit: 0, rank: 5, deckIndex: 0 }],
    isPass: false,
    timestampMs: 0,
    compType: CompType.COMP_TYPE_SINGLE,
  }
];

describe('GameBoard', () => {
  it('renders game board with all player areas', () => {
    render(
      <GameBoard
        teamLevels={[3, 4]}
        currentLevel={5}
        plays={[]}
        currentTurn={-1}
        players={mockPlayers}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('Player1')).toBeInTheDocument();
    expect(screen.getByText('Player2')).toBeInTheDocument();
    expect(screen.getByText('Player3')).toBeInTheDocument();
    expect(screen.getByText('Player4')).toBeInTheDocument();
  });

  it('displays team levels correctly', () => {
    render(
      <GameBoard
        teamLevels={[3, 4]}
        currentLevel={5}
        plays={[]}
        currentTurn={-1}
        players={mockPlayers}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('4')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('shows correct player status when trick is active', () => {
    render(
      <GameBoard
        teamLevels={[3, 4]}
        currentLevel={5}
        plays={mockPlays}
        currentTurn={1}
        players={mockPlayers}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('已出牌')).toBeInTheDocument();
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
        teamLevels={[3, 4]}
        currentLevel={5}
        plays={[]}
        currentTurn={-1}
        players={playersWithEmpty}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('空座位')).toBeInTheDocument();
  });

  it('highlights current turn player', () => {
    const { container } = render(
      <GameBoard
        teamLevels={[3, 4]}
        currentLevel={5}
        plays={mockPlays}
        currentTurn={1}
        players={mockPlayers}
        currentPlayerSeat={0}
      />
    );

    const playerAreas = container.querySelectorAll('.border-yellow-400');
    expect(playerAreas.length).toBeGreaterThan(0);
  });

  it('converts high card levels to letters correctly', () => {
    render(
      <GameBoard
        teamLevels={[11, 14]}
        currentLevel={13}
        plays={[]}
        currentTurn={-1}
        players={mockPlayers}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('J')).toBeInTheDocument();
    expect(screen.getByText('A')).toBeInTheDocument();
    expect(screen.getByText('K')).toBeInTheDocument();
  });
});
