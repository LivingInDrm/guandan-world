import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import DealResult from './DealResult';
import type { Player } from '../../types';
import type { DealEndedPayload } from '../../types/generated/event';
import { VictoryType } from '../../types/proto';

// Mock data
const mockPlayers: Player[] = [
  { id: '1', username: 'Player1', seat: 0, online: true, auto_play: false },
  { id: '2', username: 'Player2', seat: 1, online: true, auto_play: false },
  { id: '3', username: 'Player3', seat: 2, online: true, auto_play: false },
  { id: '4', username: 'Player4', seat: 3, online: true, auto_play: false }
];

const mockDealResult: DealEndedPayload = {
  dealLevel: 5,
  rankings: [0, 2, 1, 3],
  victoryType: VictoryType.VICTORY_TYPE_DOUBLE_DOWN,
  winningTeam: 0,
  levelChange: [3, 0],
  durationMs: 300000,
  trickCount: 15,
  playerStats: [
    {
      playerSeat: 0,
      cardsPlayed: 25,
      tricksWon: 8,
      passCount: 3,
      finishRank: 1
    },
    {
      playerSeat: 1,
      cardsPlayed: 22,
      tricksWon: 2,
      passCount: 5,
      finishRank: 3
    },
    {
      playerSeat: 2,
      cardsPlayed: 24,
      tricksWon: 5,
      passCount: 2,
      finishRank: 2
    },
    {
      playerSeat: 3,
      cardsPlayed: 20,
      tricksWon: 0,
      passCount: 8,
      finishRank: 4
    }
  ],
  nextDealDeadlineMs: 0
};

const mockTeamLevels: [number, number] = [5, 3]; // Team 0 at level 5, Team 1 at level 3

describe('DealResult', () => {
  let mockOnContinue: ReturnType<typeof vi.fn>;
  let mockOnExit: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnContinue = vi.fn();
    mockOnExit = vi.fn();
  });



  it('renders deal result correctly', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    // Check title
    expect(screen.getByText('局结算')).toBeInTheDocument();
    
    // Check victory information
    expect(screen.getByText(/双下.*队伍1获胜.*\+3级/)).toBeInTheDocument();
    
    // Check team sections
    expect(screen.getByText('队伍1 (胜方)')).toBeInTheDocument();
    expect(screen.getByText('队伍2 (负方)')).toBeInTheDocument();
  });

  it('displays team rankings correctly', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    // Check winning team (Team 0: Player1, Player3) - players appear in both team section and stats table
    const player1Elements = screen.getAllByText('Player1');
    expect(player1Elements.length).toBeGreaterThan(0);
    const player3Elements = screen.getAllByText('Player3');
    expect(player3Elements.length).toBeGreaterThan(0);
    const rank1Elements = screen.getAllByText('第1名');
    expect(rank1Elements.length).toBeGreaterThan(0);
    const rank2Elements = screen.getAllByText('第2名');
    expect(rank2Elements.length).toBeGreaterThan(0);

    // Check losing team (Team 1: Player2, Player4)
    const player2Elements = screen.getAllByText('Player2');
    expect(player2Elements.length).toBeGreaterThan(0);
    const player4Elements = screen.getAllByText('Player4');
    expect(player4Elements.length).toBeGreaterThan(0);
    const rank3Elements = screen.getAllByText('第3名');
    expect(rank3Elements.length).toBeGreaterThan(0);
    const rank4Elements = screen.getAllByText('第4名');
    expect(rank4Elements.length).toBeGreaterThan(0);
  });

  it('displays team levels and upgrades correctly', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    // Check team levels - use more specific queries
    const teamLevelElements = screen.getAllByText('5');
    expect(teamLevelElements.length).toBeGreaterThan(0); // Team 0 level appears
    
    const team1LevelElements = screen.getAllByText('3');
    expect(team1LevelElements.length).toBeGreaterThan(0); // Team 1 level appears
    
    // Check upgrades
    expect(screen.getByText('+3级')).toBeInTheDocument(); // Winning team upgrade
    expect(screen.getByText('+0级')).toBeInTheDocument(); // Losing team upgrade
  });

  it('displays deal statistics correctly', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    // Check statistics section
    expect(screen.getByText('本局统计')).toBeInTheDocument();
    expect(screen.getByText('5:00')).toBeInTheDocument(); // Duration formatted
    expect(screen.getByText('15')).toBeInTheDocument(); // Trick count
    expect(screen.getByText('双下')).toBeInTheDocument(); // Victory type
  });

  it('displays player statistics table correctly', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    // Check player statistics table
    expect(screen.getByText('玩家统计')).toBeInTheDocument();
    
    // Check table headers
    expect(screen.getByText('玩家')).toBeInTheDocument();
    expect(screen.getByText('排名')).toBeInTheDocument();
    expect(screen.getByText('出牌次数')).toBeInTheDocument();
    expect(screen.getByText('获胜轮次')).toBeInTheDocument();
    expect(screen.getByText('过牌次数')).toBeInTheDocument();

    // Check some player stats - use more specific queries
    expect(screen.getByText('25')).toBeInTheDocument(); // Player1 cards played
    const tricksWonElements = screen.getAllByText('8');
    expect(tricksWonElements.length).toBeGreaterThan(0); // Player1 tricks won and pass count
  });

  it('handles different victory types correctly', () => {
    const singleLastResult = {
      ...mockDealResult,
      victoryType: VictoryType.VICTORY_TYPE_SINGLE_LAST,
      levelChange: [2, 0]
    };

    render(
      <DealResult
        dealResult={singleLastResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText(/单贡.*\+2级/)).toBeInTheDocument();
  });

  it('handles high levels correctly (J, Q, K, A)', () => {
    const highLevelTeams: [number, number] = [11, 14]; // J and A levels

    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={highLevelTeams}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('J')).toBeInTheDocument(); // Team 0 level
    expect(screen.getByText('A')).toBeInTheDocument(); // Team 1 level
  });

  it('shows continue button when match is not finished', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('继续游戏')).toBeInTheDocument();
    expect(screen.getByText('退出房间')).toBeInTheDocument();
  });

  it('shows only exit button when match is finished', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={true}
        currentPlayerSeat={0}
      />
    );

    expect(screen.queryByText('继续游戏')).not.toBeInTheDocument();
    expect(screen.getByText('返回大厅')).toBeInTheDocument();
    expect(screen.getByText('比赛结束')).toBeInTheDocument();
  });

  it('calls onContinue when continue button is clicked', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    fireEvent.click(screen.getByText('继续游戏'));
    expect(mockOnContinue).toHaveBeenCalledTimes(1);
  });

  it('calls onExit when exit button is clicked', () => {
    render(
      <DealResult
        dealResult={mockDealResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    fireEvent.click(screen.getByText('退出房间'));
    expect(mockOnExit).toHaveBeenCalledTimes(1);
  });

  it('formats duration correctly for different time values', () => {
    const shortDurationResult = {
      ...mockDealResult,
      durationMs: 65000
    };

    render(
      <DealResult
        dealResult={shortDurationResult}
        players={mockPlayers}
        teamLevels={mockTeamLevels}
        onContinue={mockOnContinue}
        onExit={mockOnExit}
        isMatchFinished={false}
        currentPlayerSeat={0}
      />
    );

    expect(screen.getByText('1:05')).toBeInTheDocument();
  });
});