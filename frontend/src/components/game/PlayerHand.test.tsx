import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import PlayerHand from './PlayerHand';
import { Card } from '../../types';

// Mock cards for testing
const mockCards: Card[] = [
  { id: '1', suit: 0, rank: 14, is_joker: false }, // A♠
  { id: '2', suit: 1, rank: 14, is_joker: false }, // A♥
  { id: '3', suit: 0, rank: 13, is_joker: false }, // K♠
  { id: '4', suit: 1, rank: 13, is_joker: false }, // K♥
  { id: '5', suit: 2, rank: 5, is_joker: false },  // 5♣
  { id: '6', suit: 3, rank: 5, is_joker: false },  // 5♦
  { id: '7', suit: 0, rank: 15, is_joker: true },  // Small Joker
  { id: '8', suit: 0, rank: 16, is_joker: true },  // Big Joker
];

describe('PlayerHand', () => {
  const mockOnCardSelect = vi.fn();

  beforeEach(() => {
    mockOnCardSelect.mockClear();
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
        disabled={true}
      />
    );

    // Try to click on a card
    const cards = screen.getAllByText('A');
    fireEvent.click(cards[0].closest('.cursor-pointer')!);

    // Should not call onCardSelect when disabled
    expect(mockOnCardSelect).not.toHaveBeenCalled();

    // Buttons should be disabled
    expect(screen.getByText('全选')).toBeDisabled();
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
      { id: '1', suit: 0, rank: 5, is_joker: false }, // 5♠
      { id: '2', suit: 1, rank: 5, is_joker: false }, // 5♥
      { id: '3', suit: 2, rank: 5, is_joker: false }, // 5♣
      { id: '4', suit: 3, rank: 5, is_joker: false }, // 5♦
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
      { id: '1', suit: 0, rank: 15, is_joker: true }, // Small Joker
      { id: '2', suit: 0, rank: 16, is_joker: true }, // Big Joker
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