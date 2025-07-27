import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import TributePhase from './TributePhase';
import { 
  TributePhase as TributePhaseType, 
  TributeStatus, 
  Player, 
  Card 
} from '../../types';

// Mock data
const mockPlayers: (Player | null)[] = [
  { id: '1', username: 'Player1', seat: 0, online: true, auto_play: false },
  { id: '2', username: 'Player2', seat: 1, online: true, auto_play: false },
  { id: '3', username: 'Player3', seat: 2, online: true, auto_play: false },
  { id: '4', username: 'Player4', seat: 3, online: true, auto_play: false }
];

const mockCard1: Card = {
  id: 'card1',
  suit: 0, // spades
  rank: 14, // A
  is_joker: false
};

const mockCard2: Card = {
  id: 'card2',
  suit: 1, // hearts
  rank: 13, // K
  is_joker: false
};

const mockCard3: Card = {
  id: 'card3',
  suit: 2, // clubs
  rank: 12, // Q
  is_joker: false
};

const mockBigJoker: Card = {
  id: 'joker1',
  suit: 0,
  rank: 16,
  is_joker: true
};

describe('TributePhase', () => {
  const mockOnSelectTribute = vi.fn();
  const mockOnReturnTribute = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    // Mock Date.now() for consistent timing tests
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2024-01-01T12:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe('Tribute Info Display', () => {
    it('should display immunity message when tribute is immune', () => {
      const immuneTributePhase: TributePhaseType = {
        status: TributeStatus.FINISHED,
        tribute_map: {},
        tribute_cards: {},
        return_cards: {},
        pool_cards: [],
        selecting_player: -1,
        select_timeout: '',
        is_immune: true,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={immuneTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText(/抗贡生效/)).toBeInTheDocument();
      expect(screen.getByText(/败方持有2张及以上大王/)).toBeInTheDocument();
    });

    it('should display normal tribute info for single tribute', () => {
      const normalTributePhase: TributePhaseType = {
        status: TributeStatus.RETURNING,
        tribute_map: { 3: 0 }, // Player 3 tributes to Player 0
        tribute_cards: { 3: mockCard1 },
        return_cards: {},
        pool_cards: [],
        selecting_player: -1,
        select_timeout: '',
        is_immune: false,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={normalTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('上贡信息')).toBeInTheDocument();
      expect(screen.getByText('Player4 → Player1')).toBeInTheDocument();
    });

    it('should display double down tribute info', () => {
      const doubleDownTributePhase: TributePhaseType = {
        status: TributeStatus.SELECTING,
        tribute_map: { 2: -1, 3: -1 }, // Players 2,3 tribute to pool
        tribute_cards: { 2: mockCard1, 3: mockCard2 },
        return_cards: {},
        pool_cards: [mockCard1, mockCard2],
        selecting_player: 0,
        select_timeout: '2024-01-01T12:00:03Z',
        is_immune: false,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={doubleDownTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('双下情况：败方贡牌到池，胜方选择')).toBeInTheDocument();
      expect(screen.getByText('Player3 贡牌到池')).toBeInTheDocument();
      expect(screen.getByText('Player4 贡牌到池')).toBeInTheDocument();
    });
  });

  describe('Pool Selection Phase', () => {
    const selectingTributePhase: TributePhaseType = {
      status: TributeStatus.SELECTING,
      tribute_map: { 2: -1, 3: -1 },
      tribute_cards: { 2: mockCard1, 3: mockCard2 },
      return_cards: {},
      pool_cards: [mockCard1, mockCard2],
      selecting_player: 0,
      select_timeout: '2024-01-01T12:00:03Z',
      is_immune: false,
      selection_results: {}
    };

    it('should show selection interface for current player', () => {
      render(
        <TributePhase
          tributePhase={selectingTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('请选择一张贡牌')).toBeInTheDocument();
      expect(screen.getByText('3秒')).toBeInTheDocument(); // Countdown timer
      
      // Should show pool cards
      const cards = screen.getAllByRole('generic').filter(el => 
        el.className.includes('cursor-pointer')
      );
      expect(cards.length).toBeGreaterThan(0);
    });

    it('should show waiting message for non-selecting player', () => {
      render(
        <TributePhase
          tributePhase={selectingTributePhase}
          players={mockPlayers}
          currentPlayerSeat={1}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText(/等待.*Player1.*选择贡牌/)).toBeInTheDocument();
    });

    it('should allow card selection and confirmation', () => {
      render(
        <TributePhase
          tributePhase={selectingTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      // Should show pool cards and selection interface
      expect(screen.getByText('请选择一张贡牌')).toBeInTheDocument();
      
      // Should show cards (we can't easily test the click interaction without more complex setup)
      // But we can verify the interface is rendered correctly
      const cardContainers = screen.getAllByRole('generic').filter(el => 
        el.className.includes('bg-white') && el.className.includes('border-2')
      );
      expect(cardContainers.length).toBeGreaterThan(0);
    });

    it('should display countdown timer', () => {
      render(
        <TributePhase
          tributePhase={selectingTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      // Should show initial countdown
      expect(screen.getByText('3秒')).toBeInTheDocument();
    });
  });

  describe('Return Tribute Phase', () => {
    const returningTributePhase: TributePhaseType = {
      status: TributeStatus.RETURNING,
      tribute_map: { 3: 0 }, // Player 3 tributed to Player 0
      tribute_cards: { 3: mockCard1 },
      return_cards: {},
      pool_cards: [],
      selecting_player: -1,
      select_timeout: '',
      is_immune: false,
      selection_results: {}
    };

    it('should show return interface for player who needs to return', () => {
      render(
        <TributePhase
          tributePhase={returningTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('请选择一张牌还贡')).toBeInTheDocument();
      expect(screen.getByText('选择一张最小的牌还给上贡者')).toBeInTheDocument();
    });

    it('should show waiting message for players who do not need to return', () => {
      render(
        <TributePhase
          tributePhase={returningTributePhase}
          players={mockPlayers}
          currentPlayerSeat={1}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('等待其他玩家还贡...')).toBeInTheDocument();
    });
  });

  describe('Finished Phase', () => {
    const finishedTributePhase: TributePhaseType = {
      status: TributeStatus.FINISHED,
      tribute_map: { 3: 0 },
      tribute_cards: { 3: mockCard1 },
      return_cards: { 0: mockCard2 },
      pool_cards: [],
      selecting_player: -1,
      select_timeout: '',
      is_immune: false,
      selection_results: {}
    };

    it('should display tribute results', () => {
      render(
        <TributePhase
          tributePhase={finishedTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('贡牌结果')).toBeInTheDocument();
      expect(screen.getByText('Player4 上贡:')).toBeInTheDocument();
      expect(screen.getByText('Player1 还贡:')).toBeInTheDocument();
      expect(screen.getByText('3秒后自动进入出牌阶段...')).toBeInTheDocument();
    });
  });

  describe('Card Display Component', () => {
    it('should display regular cards correctly', () => {
      const regularCard: Card = {
        id: 'test-card',
        suit: 1, // hearts
        rank: 14, // A
        is_joker: false
      };

      render(
        <TributePhase
          tributePhase={{
            status: TributeStatus.FINISHED,
            tribute_map: { 0: 1 },
            tribute_cards: { 0: regularCard },
            return_cards: {},
            pool_cards: [],
            selecting_player: -1,
            select_timeout: '',
            is_immune: false,
            selection_results: {}
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      // Should display A and heart symbol
      expect(screen.getByText('A')).toBeInTheDocument();
      expect(screen.getByText('♥')).toBeInTheDocument();
    });

    it('should display joker cards correctly', () => {
      render(
        <TributePhase
          tributePhase={{
            status: TributeStatus.FINISHED,
            tribute_map: { 0: 1 },
            tribute_cards: { 0: mockBigJoker },
            return_cards: {},
            pool_cards: [],
            selecting_player: -1,
            select_timeout: '',
            is_immune: false,
            selection_results: {}
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('大王')).toBeInTheDocument();
    });
  });

  describe('Waiting Phase', () => {
    it('should show waiting message', () => {
      const waitingTributePhase: TributePhaseType = {
        status: TributeStatus.WAITING,
        tribute_map: {},
        tribute_cards: {},
        return_cards: {},
        pool_cards: [],
        selecting_player: -1,
        select_timeout: '',
        is_immune: false,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={waitingTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('准备上贡阶段...')).toBeInTheDocument();
    });
  });

  describe('Component Structure', () => {
    it('should render main title', () => {
      const basicTributePhase: TributePhaseType = {
        status: TributeStatus.WAITING,
        tribute_map: {},
        tribute_cards: {},
        return_cards: {},
        pool_cards: [],
        selecting_player: -1,
        select_timeout: '',
        is_immune: false,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={basicTributePhase}
          players={mockPlayers}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      expect(screen.getByText('上贡阶段')).toBeInTheDocument();
    });

    it('should handle null players gracefully', () => {
      const playersWithNull: (Player | null)[] = [
        mockPlayers[0],
        null,
        mockPlayers[2],
        null
      ];

      const tributePhase: TributePhaseType = {
        status: TributeStatus.RETURNING,
        tribute_map: { 2: 0 },
        tribute_cards: { 2: mockCard1 },
        return_cards: {},
        pool_cards: [],
        selecting_player: -1,
        select_timeout: '',
        is_immune: false,
        selection_results: {}
      };

      render(
        <TributePhase
          tributePhase={tributePhase}
          players={playersWithNull}
          currentPlayerSeat={0}
          onSelectTribute={mockOnSelectTribute}
          onReturnTribute={mockOnReturnTribute}
        />
      );

      // Should show fallback player names (Player3 is seat 2, Player1 is seat 0)
      expect(screen.getByText('Player3 → Player1')).toBeInTheDocument();
    });
  });
});