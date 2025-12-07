import React, { useState } from 'react';
import type { Story } from '@ladle/react';
import PlayerControlPanel from './PlayerControlPanel';
import type { Card, PlayAction } from '../../types/proto';
import { CompType } from '../../types/generated/common';

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

function createPlayAction(cards: Card[], compType: CompType, playerSeat = 1): PlayAction {
  return {
    playerSeat,
    cards,
    compType,
    timestampMs: Date.now(),
    isPass: false,
  };
}

interface WrapperProps {
  initialCards: Card[];
  initialSelected?: Card[];
  currentLevel?: number;
  canPlay?: boolean;
  isMyTurn?: boolean;
  disabled?: boolean;
  plays?: PlayAction[];
  leader?: number;
  playerSeat?: number;
  dealLevel?: number;
}

const InteractiveWrapper: React.FC<WrapperProps> = ({
  initialCards,
  initialSelected = [],
  currentLevel = 10,
  canPlay = true,
  isMyTurn = true,
  disabled = false,
  plays = [],
  leader = 0,
  playerSeat = 0,
  dealLevel = 2,
}) => {
  const [selectedCards, setSelectedCards] = useState<Card[]>(initialSelected);
  const turnDeadlineAtMs = Date.now() + 30000;

  const handlePlayCards = (cards: Card[]) => {
    console.log('Play cards:', cards);
    setSelectedCards([]);
  };

  const handlePass = () => {
    console.log('Pass');
    setSelectedCards([]);
  };

  const handleHint = (cards: Card[]) => {
    console.log('Hint:', cards);
    setSelectedCards(cards);
  };

  return (
    <div className="inline-block border border-dashed border-gray-400">
      <PlayerControlPanel
        cards={initialCards}
        selectedCards={selectedCards}
        onCardSelect={setSelectedCards}
        currentLevel={currentLevel}
        canPlay={canPlay}
        isMyTurn={isMyTurn}
        turnDeadlineAtMs={turnDeadlineAtMs}
        onPlayCards={handlePlayCards}
        onPass={handlePass}
        onHint={handleHint}
        plays={plays}
        leader={leader}
        playerSeat={playerSeat}
        dealLevel={dealLevel}
        disabled={disabled}
      />
    </div>
  );
};

export const Default: Story = () => {
  const cards = createMockHand();
  return <InteractiveWrapper initialCards={cards} />;
};
Default.meta = {
  description: '默认状态 - 玩家回合，可出牌',
};

export const NotMyTurn: Story = () => {
  const cards = createMockHand();
  return <InteractiveWrapper initialCards={cards} isMyTurn={false} />;
};
NotMyTurn.meta = {
  description: '非玩家回合 - 控制面板禁用',
};

export const Disabled: Story = () => {
  const cards = createMockHand();
  return <InteractiveWrapper initialCards={cards} disabled={true} />;
};
Disabled.meta = {
  description: '禁用状态',
};

export const EmptyHand: Story = () => {
  return <InteractiveWrapper initialCards={[]} />;
};
EmptyHand.meta = {
  description: '空手牌状态',
};

export const WithSelectedCards: Story = () => {
  const cards = createFewCards();
  const selected = [cards[0], cards[1]];
  return <InteractiveWrapper initialCards={cards} initialSelected={selected} />;
};
WithSelectedCards.meta = {
  description: '预选中部分牌（一对A）',
};

export const WithJokers: Story = () => {
  const cards = createJokerHand();
  return <InteractiveWrapper initialCards={cards} />;
};
WithJokers.meta = {
  description: '包含大小王的手牌',
};

export const FewCards: Story = () => {
  const cards = createFewCards();
  return <InteractiveWrapper initialCards={cards} />;
};
FewCards.meta = {
  description: '少量手牌 - 6张',
};

export const WithSingle: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([createCard(10, 3, 100)], CompType.COMP_TYPE_SINGLE)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithSingle.meta = {
  description: '对手出单张 - 方片10',
};

export const WithPair: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([createCard(11, 0, 100), createCard(11, 0, 101)], CompType.COMP_TYPE_PAIR)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithPair.meta = {
  description: '对手出对子 - 黑桃JJ',
};

export const WithTriple: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([createCard(14, 0, 100), createCard(14, 1, 101), createCard(14, 2, 102)], CompType.COMP_TYPE_TRIPLE)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithTriple.meta = {
  description: '对手出三张 - AAA',
};

export const WithPlate: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([
    createCard(11, 0, 100), createCard(11, 1, 101),
    createCard(12, 0, 102), createCard(12, 1, 103),
    createCard(13, 0, 104), createCard(13, 1, 105),
  ], CompType.COMP_TYPE_PLATE)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithPlate.meta = {
  description: '对手出钢板 - JJQQKK',
};

export const WithTube: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([
    createCard(12, 0, 100), createCard(12, 1, 101), createCard(12, 2, 102),
    createCard(13, 0, 103), createCard(13, 1, 104), createCard(13, 2, 105),
  ], CompType.COMP_TYPE_TUBE)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithTube.meta = {
  description: '对手出钢管 - QQQKKK',
};

export const WithFullHouse: Story = () => {
  const cards = createMockHand();
  const plays = [createPlayAction([
    createCard(9, 0, 100), createCard(9, 1, 101), createCard(9, 2, 102),
    createCard(3, 0, 103), createCard(3, 1, 104),
  ], CompType.COMP_TYPE_FULL_HOUSE)];
  return <InteractiveWrapper initialCards={cards} plays={plays} leader={1} />;
};
WithFullHouse.meta = {
  description: '对手出葫芦 - 99933',
};
