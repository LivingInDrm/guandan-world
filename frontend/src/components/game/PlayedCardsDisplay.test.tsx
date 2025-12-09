import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import PlayedCardsDisplay from './PlayedCardsDisplay';
import type { PlayAction } from '../../types/proto';
import type { Card } from '../../types/proto';
import { CompType } from '../../types/proto';

const mockCards: Card[] = [
  { suit: 0, rank: 14, deckIndex: 0 },
  { suit: 1, rank: 13, deckIndex: 1 },
  { suit: 2, rank: 5, deckIndex: 2 },
];

describe('PlayedCardsDisplay', () => {
  it('renders empty container when play is null', () => {
    const { container } = render(
      <PlayedCardsDisplay play={null} position="bottom" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div).toBeEmptyDOMElement();
  });

  it('renders pass indicator when isPass is true', () => {
    const passPlay: PlayAction = {
      playerSeat: 0,
      isPass: true,
      cards: [],
      compType: CompType.COMP_TYPE_FOLD,
      timestampMs: Date.now(),
    };

    render(<PlayedCardsDisplay play={passPlay} position="bottom" />);
    
    expect(screen.getByText('不出')).toBeInTheDocument();
  });

  it('renders cards when play has cards', () => {
    const cardPlay: PlayAction = {
      playerSeat: 0,
      isPass: false,
      cards: mockCards,
      compType: CompType.COMP_TYPE_TRIPLE,
      timestampMs: Date.now(),
    };

    const { container } = render(
      <PlayedCardsDisplay play={cardPlay} position="bottom" />
    );
    
    const cardElements = container.querySelectorAll('[role="button"]');
    expect(cardElements.length).toBe(3);
  });

  it('renders empty container when play has empty cards and isPass is false', () => {
    const emptyPlay: PlayAction = {
      playerSeat: 0,
      isPass: false,
      cards: [],
      compType: CompType.COMP_TYPE_UNSPECIFIED,
      timestampMs: Date.now(),
    };

    const { container } = render(
      <PlayedCardsDisplay play={emptyPlay} position="bottom" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div).toBeEmptyDOMElement();
  });

  it('applies correct position classes for bottom', () => {
    const { container } = render(
      <PlayedCardsDisplay play={null} position="bottom" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div.className).toContain('bottom-4');
    expect(div.className).toContain('left-1/2');
  });

  it('applies correct position classes for top', () => {
    const { container } = render(
      <PlayedCardsDisplay play={null} position="top" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div.className).toContain('top-4');
    expect(div.className).toContain('left-1/2');
  });

  it('applies correct position classes for left', () => {
    const { container } = render(
      <PlayedCardsDisplay play={null} position="left" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div.className).toContain('left-4');
    expect(div.className).toContain('top-1/2');
  });

  it('applies correct position classes for right', () => {
    const { container } = render(
      <PlayedCardsDisplay play={null} position="right" />
    );
    
    const div = container.firstChild as HTMLElement;
    expect(div.className).toContain('right-4');
    expect(div.className).toContain('top-1/2');
  });
});
