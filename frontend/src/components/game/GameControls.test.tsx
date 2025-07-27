import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import GameControls from './GameControls';
import { Card } from '../../types';

// Mock cards for testing
const mockCards: Card[] = [
  { id: '1', suit: 0, rank: 14, is_joker: false }, // A♠
  { id: '2', suit: 1, rank: 14, is_joker: false }, // A♥
  { id: '3', suit: 0, rank: 13, is_joker: false }, // K♠
];

describe('GameControls', () => {
  const mockOnPlayCards = vi.fn();
  const mockOnPass = vi.fn();

  beforeEach(() => {
    mockOnPlayCards.mockClear();
    mockOnPass.mockClear();
  });

  it('renders correctly when it is player turn', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('轮到您出牌')).toBeInTheDocument();
    expect(screen.getByText('出牌')).toBeInTheDocument();
    expect(screen.getByText('不出')).toBeInTheDocument();
  });

  it('renders correctly when it is not player turn', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={false}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('等待其他玩家')).toBeInTheDocument();
    expect(screen.getByText('等待其他玩家操作...')).toBeInTheDocument();
  });

  it('shows countdown timer when it is player turn', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={15}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('15秒')).toBeInTheDocument();
  });

  it('disables play button when no cards selected', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    expect(playButton).toBeDisabled();
  });

  it('enables play button when valid cards are selected', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const playButton = screen.getByText('出牌 (1张)').closest('button');
    expect(playButton).not.toBeDisabled();
  });

  it('shows validation message for valid single card', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('单牌 - 可以出牌')).toBeInTheDocument();
  });

  it('shows validation message for valid pair', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0], mockCards[1]]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('对子 - 可以出牌')).toBeInTheDocument();
  });

  it('shows validation error for invalid pair', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0], mockCards[2]]} // A and K
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('两张牌必须是对子')).toBeInTheDocument();
  });

  it('calls onPlayCards when play button is clicked with valid cards', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const playButton = screen.getByText('出牌 (1张)').closest('button');
    fireEvent.click(playButton!);

    expect(mockOnPlayCards).toHaveBeenCalledWith([mockCards[0]]);
  });

  it('calls onPass when pass button is clicked', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const passButton = screen.getByText('不出').closest('button');
    fireEvent.click(passButton!);

    expect(mockOnPass).toHaveBeenCalled();
  });

  it('disables all buttons when disabled prop is true', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        disabled={true}
      />
    );

    const playButton = screen.getByText('出牌 (1张)').closest('button');
    const passButton = screen.getByText('不出').closest('button');

    expect(playButton).toBeDisabled();
    expect(passButton).toBeDisabled();
  });

  it('disables buttons when not player turn', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={false}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const playButton = screen.getByText('出牌 (1张)').closest('button');
    const passButton = screen.getByText('不出').closest('button');

    expect(playButton).toBeDisabled();
    expect(passButton).toBeDisabled();
  });

  it('disables buttons when cannot play', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={false}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    const playButton = screen.getByText('出牌 (1张)').closest('button');
    const passButton = screen.getByText('不出').closest('button');

    expect(playButton).toBeDisabled();
    expect(passButton).toBeDisabled();
  });

  it('shows help text when no cards selected', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('请选择要出的牌，或点击"不出"跳过')).toBeInTheDocument();
  });

  it('countdown timer changes color based on time left', async () => {
    const { rerender } = render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={15}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    // Should show blue color for > 10 seconds
    expect(screen.getByText('15秒').closest('div')).toHaveClass('text-blue-600');

    // Change to 8 seconds - should show orange
    rerender(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={8}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('8秒').closest('div')).toHaveClass('text-orange-600');

    // Change to 3 seconds - should show red
    rerender(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={3}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('3秒').closest('div')).toHaveClass('text-red-600');
  });

  it('validates three of a kind correctly', () => {
    const threeCards = [
      { id: '1', suit: 0, rank: 5, is_joker: false },
      { id: '2', suit: 1, rank: 5, is_joker: false },
      { id: '3', suit: 2, rank: 5, is_joker: false },
    ];

    render(
      <GameControls
        selectedCards={threeCards}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('三张 - 可以出牌')).toBeInTheDocument();
  });

  it('shows error for invalid three cards', () => {
    const invalidThreeCards = [
      { id: '1', suit: 0, rank: 5, is_joker: false },
      { id: '2', suit: 1, rank: 6, is_joker: false },
      { id: '3', suit: 2, rank: 7, is_joker: false },
    ];

    render(
      <GameControls
        selectedCards={invalidThreeCards}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('三张牌必须是同点数')).toBeInTheDocument();
  });

  it('handles complex card combinations', () => {
    const fiveCards = [
      { id: '1', suit: 0, rank: 5, is_joker: false },
      { id: '2', suit: 1, rank: 6, is_joker: false },
      { id: '3', suit: 2, rank: 7, is_joker: false },
      { id: '4', suit: 3, rank: 8, is_joker: false },
      { id: '5', suit: 0, rank: 9, is_joker: false },
    ];

    render(
      <GameControls
        selectedCards={fiveCards}
        canPlay={true}
        isMyTurn={true}
        turnTimeoutSeconds={20}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
      />
    );

    expect(screen.getByText('5张牌型 - 可以出牌')).toBeInTheDocument();
  });
});