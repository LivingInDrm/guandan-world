import type { Story } from '@ladle/react';
import PlayedCardsDisplay from './PlayedCardsDisplay';
import type { PlayAction } from '../../types/proto';
import type { Card } from '../../types/proto';

const mockCards: Card[] = [
  { suit: 0, rank: 14, deckIndex: 0 },
  { suit: 1, rank: 13, deckIndex: 1 },
  { suit: 2, rank: 5, deckIndex: 2 },
];

const manyCards: Card[] = [
  { suit: 0, rank: 14, deckIndex: 0 },
  { suit: 1, rank: 14, deckIndex: 1 },
  { suit: 2, rank: 14, deckIndex: 2 },
  { suit: 3, rank: 14, deckIndex: 3 },
  { suit: 0, rank: 13, deckIndex: 4 },
  { suit: 1, rank: 13, deckIndex: 5 },
  { suit: 2, rank: 13, deckIndex: 6 },
  { suit: 3, rank: 13, deckIndex: 7 },
];

const jokerCards: Card[] = [
  { suit: -1, rank: 16, deckIndex: 0 },
  { suit: -1, rank: 15, deckIndex: 1 },
  { suit: 0, rank: 14, deckIndex: 2 },
];

const Container: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div 
    className="relative bg-green-200 border-2 border-green-400 rounded-lg"
    style={{ width: 500, height: 300 }}
  >
    {children}
  </div>
);

export const Empty: Story = () => (
  <Container>
    <PlayedCardsDisplay play={null} position="bottom" currentLevel={13} />
  </Container>
);

export const Pass: Story = () => {
  const passPlay: PlayAction = {
    playerSeat: 0,
    isPass: true,
    cards: [],
  };

  return (
    <Container>
      <PlayedCardsDisplay play={passPlay} position="bottom" currentLevel={13} />
    </Container>
  );
};

export const SingleCard: Story = () => {
  const play: PlayAction = {
    playerSeat: 0,
    isPass: false,
    cards: [mockCards[0]],
  };

  return (
    <Container>
      <PlayedCardsDisplay play={play} position="bottom" currentLevel={13} />
    </Container>
  );
};

export const MultipleCards: Story = () => {
  const play: PlayAction = {
    playerSeat: 0,
    isPass: false,
    cards: mockCards,
  };

  return (
    <Container>
      <PlayedCardsDisplay play={play} position="bottom" currentLevel={13} />
    </Container>
  );
};

export const ManyCards: Story = () => {
  const play: PlayAction = {
    playerSeat: 0,
    isPass: false,
    cards: manyCards,
  };

  return (
    <Container>
      <PlayedCardsDisplay play={play} position="bottom" currentLevel={13} />
    </Container>
  );
};

export const WithJokers: Story = () => {
  const play: PlayAction = {
    playerSeat: 0,
    isPass: false,
    cards: jokerCards,
  };

  return (
    <Container>
      <PlayedCardsDisplay play={play} position="bottom" currentLevel={13} />
    </Container>
  );
};

export const AllPositions: Story = () => {
  const play: PlayAction = {
    playerSeat: 0,
    isPass: false,
    cards: mockCards,
  };

  const passPlay: PlayAction = {
    playerSeat: 1,
    isPass: true,
    cards: [],
  };

  return (
    <div 
      className="relative bg-green-200 border-2 border-green-400 rounded-lg"
      style={{ width: 500, height: 300 }}
    >
      <PlayedCardsDisplay play={play} position="bottom" currentLevel={13} />
      <PlayedCardsDisplay play={passPlay} position="left" currentLevel={13} />
      <PlayedCardsDisplay play={play} position="top" currentLevel={13} />
      <PlayedCardsDisplay play={null} position="right" currentLevel={13} />
    </div>
  );
};

AllPositions.meta = {
  description: 'bottom: cards, left: pass, top: cards, right: empty',
};
