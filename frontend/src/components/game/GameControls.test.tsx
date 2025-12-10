import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import GameControls from './GameControls';
import type { Card } from '../../types/proto';

const mockCards: Card[] = [
  { suit: 0, rank: 14, deckIndex: 0 },
  { suit: 1, rank: 14, deckIndex: 1 },
  { suit: 0, rank: 13, deckIndex: 2 },
];

describe('GameControls', () => {
  let mockOnPlayCards: ReturnType<typeof vi.fn>;
  let mockOnPass: ReturnType<typeof vi.fn>;
  let mockOnHint: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnPlayCards = vi.fn();
    mockOnPass = vi.fn();
    mockOnHint = vi.fn();
  });

  const futureDeadline = () => Date.now() + 20000;

  const defaultProps = {
    handCards: mockCards,
    plays: [],
    playerSeat: 0,
    dealLevel: 2,
  };

  it('renders all four control buttons', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    expect(screen.getByText('不出')).toBeInTheDocument();
    expect(screen.getByText('提示')).toBeInTheDocument();
    expect(screen.getByText('出牌')).toBeInTheDocument();
  });

  it('disables play button when no cards selected', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
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
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    expect(playButton).not.toBeDisabled();
  });

  it('calls onPlayCards when play button is clicked with valid cards', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    fireEvent.click(playButton!);

    expect(mockOnPlayCards).toHaveBeenCalledWith([mockCards[0]]);
  });

  it('calls onPass when pass button is clicked', () => {
    const nonEmptyPlays = [
      { playerSeat: 1, cards: mockCards, isPass: false, compType: 1 }
    ];
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
        plays={nonEmptyPlays}
      />
    );

    const passButton = screen.getByText('不出').closest('button');
    fireEvent.click(passButton!);

    expect(mockOnPass).toHaveBeenCalled();
  });

  it('calls onHint when hint button is clicked', () => {
    render(
      <GameControls
        selectedCards={[]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const hintButton = screen.getByText('提示').closest('button');
    fireEvent.click(hintButton!);

    expect(mockOnHint).toHaveBeenCalled();
  });

  it('disables all buttons when disabled prop is true', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        disabled={true}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    const passButton = screen.getByText('不出').closest('button');
    const hintButton = screen.getByText('提示').closest('button');

    expect(playButton).toBeDisabled();
    expect(passButton).toBeDisabled();
    expect(hintButton).toBeDisabled();
  });

  it('disables buttons when not player turn', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0]]}
        canPlay={true}
        isMyTurn={false}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
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
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    const passButton = screen.getByText('不出').closest('button');

    expect(playButton).toBeDisabled();
    expect(passButton).toBeDisabled();
  });

  it('disables play button for invalid card combination', () => {
    render(
      <GameControls
        selectedCards={[mockCards[0], mockCards[2]]}
        canPlay={true}
        isMyTurn={true}
        turnDeadlineAtMs={futureDeadline()}
        onPlayCards={mockOnPlayCards}
        onPass={mockOnPass}
        onHint={mockOnHint}
        {...defaultProps}
      />
    );

    const playButton = screen.getByText('出牌').closest('button');
    expect(playButton).toBeDisabled();
  });
});
