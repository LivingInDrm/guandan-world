import React, { useState } from 'react';
import type { Story } from '@ladle/react';
import PlayerHand from './PlayerHand';
import type { Card } from '../../types/proto';

function createCard(rank: number, suit: number, deckIndex: number): Card {
  return { rank, suit, deckIndex } as Card;
}

function createMockHand(): Card[] {
  const cards: Card[] = [];
  let deckIndex = 0;
  
  cards.push(createCard(16, -1, deckIndex++));
  cards.push(createCard(15, -1, deckIndex++));
  
  for (let rank = 14; rank >= 2; rank--) {
    for (let suit = 0; suit < 4; suit++) {
      if (cards.length < 27) {
        cards.push(createCard(rank, suit, deckIndex++));
      }
    }
  }
  
  return cards.slice(0, 27);
}

function createFewCards(): Card[] {
  return [
    createCard(14, 0, 0),
    createCard(14, 1, 1),
    createCard(13, 2, 2),
    createCard(10, 3, 3),
    createCard(7, 0, 4),
    createCard(3, 1, 5),
  ];
}

function createJokerHand(): Card[] {
  return [
    createCard(16, -1, 0),
    createCard(16, -1, 1),
    createCard(15, -1, 2),
    createCard(15, -1, 3),
    createCard(14, 0, 4),
    createCard(14, 1, 5),
    createCard(2, 0, 6),
    createCard(2, 1, 7),
  ];
}

const InteractiveWrapper: React.FC<{ initialCards: Card[]; initialSelected?: Card[]; currentLevel?: number }> = ({ 
  initialCards, 
  initialSelected = [],
  currentLevel = 10
}) => {
  const [selectedCards, setSelectedCards] = useState<Card[]>(initialSelected);
  
  return (
    <div className="p-4">
      <PlayerHand
        cards={initialCards}
        selectedCards={selectedCards}
        onCardSelect={setSelectedCards}
        currentLevel={currentLevel}
      />
    </div>
  );
};

export const Default: Story = () => {
  const cards = createMockHand();
  return <InteractiveWrapper initialCards={cards} />;
};
Default.meta = {
  description: '默认状态 - 27张手牌',
};

export const Empty: Story = () => {
  return <InteractiveWrapper initialCards={[]} />;
};
Empty.meta = {
  description: '空手牌状态',
};

export const FewCards: Story = () => {
  const cards = createFewCards();
  return <InteractiveWrapper initialCards={cards} />;
};
FewCards.meta = {
  description: '少量手牌 - 6张',
};

export const WithJokers: Story = () => {
  const cards = createJokerHand();
  return <InteractiveWrapper initialCards={cards} />;
};
WithJokers.meta = {
  description: '包含大小王的手牌',
};

export const WithSelectedCards: Story = () => {
  const cards = createFewCards();
  const selected = [cards[0], cards[1]];
  return <InteractiveWrapper initialCards={cards} initialSelected={selected} />;
};
WithSelectedCards.meta = {
  description: '预选中部分牌',
};
