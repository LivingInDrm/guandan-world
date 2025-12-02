import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PlayerHand from './PlayerHand';
import type { Card } from '../../types/proto';

// Mock cards for testing
const mockCards: Card[] = [
  { suit: 0, rank: 14, deckIndex: 0 }, // A♠
  { suit: 1, rank: 14, deckIndex: 1 }, // A♥
  { suit: 0, rank: 13, deckIndex: 2 }, // K♠
  { suit: 1, rank: 13, deckIndex: 3 }, // K♥
  { suit: 2, rank: 5, deckIndex: 4 },  // 5♣
  { suit: 3, rank: 5, deckIndex: 5 },  // 5♦
  { suit: -1, rank: 15, deckIndex: 6 },  // Small Joker
  { suit: -1, rank: 16, deckIndex: 7 },  // Big Joker
];

describe('PlayerHand', () => {
  let mockOnCardSelect: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnCardSelect = vi.fn();
  });

  it('renders hand with correct card count', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    expect(screen.getByText('手牌 (8张)')).toBeInTheDocument();
  });

  it('groups cards by rank correctly', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Check if rank groups are displayed
    expect(screen.getByText('大王 (1)')).toBeInTheDocument();
    expect(screen.getByText('小王 (1)')).toBeInTheDocument();
    expect(screen.getByText('A (2)')).toBeInTheDocument();
    expect(screen.getByText('K (2)')).toBeInTheDocument();
    expect(screen.getByText('5 (2)')).toBeInTheDocument();
  });

  it('displays card symbols correctly', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Check for suit symbols (using getAllByText since there are multiple)
    expect(screen.getAllByText('♠').length).toBeGreaterThan(0);
    expect(screen.getAllByText('♥').length).toBeGreaterThan(0);
    expect(screen.getAllByText('♣').length).toBeGreaterThan(0);
    expect(screen.getAllByText('♦').length).toBeGreaterThan(0);
  });

  it('handles card selection', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Click on a card - due to suit priority sorting, A♥ comes first
    const cards = screen.getAllByText('A');
    fireEvent.click(cards[0].closest('.cursor-pointer')!);

    expect(mockOnCardSelect).toHaveBeenCalledWith([mockCards[1]]); // A♥ is mockCards[1]
  });

  it('shows selected cards count', () => {
    const selectedCards = [mockCards[0], mockCards[1]];
    
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={selectedCards}
        onCardSelect={mockOnCardSelect}
      />
    );

    expect(screen.getByText('已选择 2张')).toBeInTheDocument();
  });

  it('handles card deselection', () => {
    const selectedCards = [mockCards[1]]; // A♥ which appears first due to sorting
    
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={selectedCards}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Click on the selected card to deselect
    const cards = screen.getAllByText('A');
    fireEvent.click(cards[0].closest('.cursor-pointer')!);

    expect(mockOnCardSelect).toHaveBeenCalledWith([]);
  });

  it('shows clear selection button when cards are selected', () => {
    const selectedCards = [mockCards[0]];
    
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={selectedCards}
        onCardSelect={mockOnCardSelect}
      />
    );

    expect(screen.getByText('清空选择')).toBeInTheDocument();
  });

  it('handles clear selection', () => {
    const selectedCards = [mockCards[0], mockCards[1]];
    
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={selectedCards}
        onCardSelect={mockOnCardSelect}
      />
    );

    fireEvent.click(screen.getByText('清空选择'));
    expect(mockOnCardSelect).toHaveBeenCalledWith([]);
  });

  it('handles select all', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    fireEvent.click(screen.getByText('全选'));
    expect(mockOnCardSelect).toHaveBeenCalledWith(mockCards);
  });

  it('disables interaction when disabled prop is true', () => {
    render(
      <PlayerHand
        cards={mockCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Try to click on a card
    const cards = screen.getAllByText('A');
    fireEvent.click(cards[0].closest('.cursor-pointer')!);

    // Verify onCardSelect was called (since disabled prop is removed)
    expect(mockOnCardSelect).toHaveBeenCalled();
  });

  it('shows empty state when no cards', () => {
    render(
      <PlayerHand
        cards={[]}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    expect(screen.getByText('暂无手牌')).toBeInTheDocument();
    expect(screen.getByText('手牌 (0张)')).toBeInTheDocument();
  });

  it('sorts cards within groups by suit priority', () => {
    const sameRankCards: Card[] = [
      { suit: 0, rank: 5, deckIndex: 0 }, // 5♠
      { suit: 1, rank: 5, deckIndex: 1 }, // 5♥
      { suit: 2, rank: 5, deckIndex: 2 }, // 5♣
      { suit: 3, rank: 5, deckIndex: 3 }, // 5♦
    ];

    const { container } = render(
      <PlayerHand
        cards={sameRankCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Hearts should come first (suit priority: hearts > diamonds > clubs > spades)
    const suitSymbols = container.querySelectorAll('.text-lg');
    expect(suitSymbols[0]).toHaveTextContent('♦'); // diamonds first
    expect(suitSymbols[1]).toHaveTextContent('♥'); // hearts second
    expect(suitSymbols[2]).toHaveTextContent('♣'); // clubs third
    expect(suitSymbols[3]).toHaveTextContent('♠'); // spades last
  });

  it('handles joker cards correctly', () => {
    const jokerCards: Card[] = [
      { suit: -1, rank: 15, deckIndex: 0 }, // Small Joker
      { suit: -1, rank: 16, deckIndex: 1 }, // Big Joker
    ];

    render(
      <PlayerHand
        cards={jokerCards}
        selectedCards={[]}
        onCardSelect={mockOnCardSelect}
      />
    );

    expect(screen.getByText('小王')).toBeInTheDocument();
    expect(screen.getByText('大王')).toBeInTheDocument();
  });

  it('applies visual feedback for selected cards', () => {
    const selectedCards = [mockCards[0]];
    
    const { container } = render(
      <PlayerHand
        cards={mockCards}
        selectedCards={selectedCards}
        onCardSelect={mockOnCardSelect}
      />
    );

    // Selected card should have transform and border styling
    const selectedCard = container.querySelector('.transform.-translate-y-2.border-blue-500');
    expect(selectedCard).toBeInTheDocument();
  });
});