import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import MatchResult from './MatchResult';
import type { MatchEndedPayload } from '../../types/generated/event';
import type { Player } from '../../types';

// Mock data
const mockMatchResult: MatchEndedPayload = {
  winner: 0,
  finalLevels: [14, 12],
  durationMs: 1800000,
  totalDeals: 8
};

const mockPlayers: Player[] = [
  { id: '1', username: 'Alice', seat: 0, online: true, auto_play: false },
  { id: '2', username: 'Bob', seat: 1, online: true, auto_play: false },
  { id: '3', username: 'Charlie', seat: 2, online: true, auto_play: false },
  { id: '4', username: 'David', seat: 3, online: true, auto_play: false }
];

describe('MatchResult', () => {
  let mockOnReturnToLobby: ReturnType<typeof vi.fn>;
  let mockOnPlayAgain: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnReturnToLobby = vi.fn();
    mockOnPlayAgain = vi.fn();
  });

  it('renders match result correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check main title
    expect(screen.getByText('🎉 比赛结束 🎉')).toBeInTheDocument();
    
    // Check winner announcement
    expect(screen.getByText('队伍1获得最终胜利！')).toBeInTheDocument();
    expect(screen.getByText('恭喜率先达到A级！')).toBeInTheDocument();
  });

  it('displays winning team correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check winning team section
    expect(screen.getByText('🏆 队伍1')).toBeInTheDocument();
    expect(screen.getByText('冠军队伍')).toBeInTheDocument();
    
    // Check winning team players (Team 0: Alice, Charlie)
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Charlie')).toBeInTheDocument();
    expect(screen.getByText('座位1')).toBeInTheDocument(); // Alice's seat
    expect(screen.getByText('座位3')).toBeInTheDocument(); // Charlie's seat
  });

  it('displays losing team correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check losing team section
    expect(screen.getByText('队伍2')).toBeInTheDocument();
    expect(screen.getByText('亚军队伍')).toBeInTheDocument();
    
    // Check losing team players (Team 1: Bob, David)
    expect(screen.getByText('Bob')).toBeInTheDocument();
    expect(screen.getByText('David')).toBeInTheDocument();
    expect(screen.getByText('座位2')).toBeInTheDocument(); // Bob's seat
    expect(screen.getByText('座位4')).toBeInTheDocument(); // David's seat
  });

  it('displays final levels correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check final levels display
    const finalLevelElements = screen.getAllByText('A'); // Team 0 final level
    expect(finalLevelElements.length).toBeGreaterThan(0);
    
    const team1LevelElements = screen.getAllByText('Q'); // Team 1 final level
    expect(team1LevelElements.length).toBeGreaterThan(0);
  });

  it('displays match statistics correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check statistics section
    expect(screen.getByText('比赛统计')).toBeInTheDocument();
    expect(screen.getByText('8')).toBeInTheDocument(); // Total deals
    expect(screen.getByText('30:00')).toBeInTheDocument(); // Total duration
    expect(screen.getByText('A vs Q')).toBeInTheDocument(); // Final levels comparison
    expect(screen.getByText('队伍1')).toBeInTheDocument(); // Winner in stats
    
    // Check statistic labels
    expect(screen.getByText('总局数')).toBeInTheDocument();
    expect(screen.getByText('总时长')).toBeInTheDocument();
    expect(screen.getByText('最终等级')).toBeInTheDocument();
    expect(screen.getByText('获胜队伍')).toBeInTheDocument();
  });

  it('displays congratulations message correctly', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Check congratulations section
    expect(screen.getByText('感谢参与本次掼蛋对战！')).toBeInTheDocument();
    expect(screen.getByText(/经过 8 局激烈的对战/)).toBeInTheDocument();
    expect(screen.getByText(/队伍1 成功率先达到A级/)).toBeInTheDocument();
    
    // Check motivational text
    expect(screen.getByText('🎯 精彩对局')).toBeInTheDocument();
    expect(screen.getByText('🤝 友谊第一')).toBeInTheDocument();
    expect(screen.getByText('🏆 比赛第二')).toBeInTheDocument();
  });

  it('handles different victory scenarios', () => {
    const team1WinResult = {
      ...mockMatchResult,
      winner: 1,
      final_levels: [11, 14] as [number, number] // Team 1 wins with A, Team 0 at J
    };

    render(
      <MatchResult
        matchResult={team1WinResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    expect(screen.getByText('队伍2获得最终胜利！')).toBeInTheDocument();
    expect(screen.getByText('🏆 队伍2')).toBeInTheDocument();
  });

  it('formats duration correctly for different time values', () => {
    const longMatchResult = {
      ...mockMatchResult,
      durationMs: 3665000 // 1 hour, 1 minute, 5 seconds
    };

    render(
      <MatchResult
        matchResult={longMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    expect(screen.getByText('1:01:05')).toBeInTheDocument();
  });

  it('handles different level displays correctly', () => {
    const differentLevelsResult = {
      ...mockMatchResult,
      final_levels: [13, 10] as [number, number] // K vs 10
    };

    render(
      <MatchResult
        matchResult={differentLevelsResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    const kElements = screen.getAllByText('K');
    expect(kElements.length).toBeGreaterThan(0);
    const tenElements = screen.getAllByText('10');
    expect(tenElements.length).toBeGreaterThan(0);
    expect(screen.getByText('K vs 10')).toBeInTheDocument();
  });

  it('calls onReturnToLobby when return button is clicked', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    const returnButton = screen.getByText('返回大厅');
    expect(returnButton).toBeInTheDocument();
    
    fireEvent.click(returnButton);
    expect(mockOnReturnToLobby).toHaveBeenCalledTimes(1);
  });

  it('displays correct team assignments', () => {
    render(
      <MatchResult
        matchResult={mockMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    // Team 0 should have seats 0 and 2 (Alice and Charlie)
    // Team 1 should have seats 1 and 3 (Bob and David)
    
    // Check that the correct players are displayed
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Charlie')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
    expect(screen.getByText('David')).toBeInTheDocument();
    
    // Check that the winning team section exists
    expect(screen.getByText('🏆 队伍1')).toBeInTheDocument();
    expect(screen.getByText('冠军队伍')).toBeInTheDocument();
    
    // Check that the losing team section exists
    expect(screen.getByText('队伍2')).toBeInTheDocument();
    expect(screen.getByText('亚军队伍')).toBeInTheDocument();
  });

  it('handles edge case with minimum match duration', () => {
    const quickMatchResult = {
      ...mockMatchResult,
      durationMs: 65000 // 1 minute, 5 seconds
    };

    render(
      <MatchResult
        matchResult={quickMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    expect(screen.getByText('1:05')).toBeInTheDocument();
  });

  it('displays correct congratulations message with dynamic content', () => {
    const customMatchResult = {
      ...mockMatchResult,
      totalDeals: 12
    };

    render(
      <MatchResult
        matchResult={customMatchResult}
        players={mockPlayers}
        onReturnToLobby={mockOnReturnToLobby}
        onPlayAgain={mockOnPlayAgain}
      />
    );

    expect(screen.getByText(/经过 12 局激烈的对战/)).toBeInTheDocument();
  });
});