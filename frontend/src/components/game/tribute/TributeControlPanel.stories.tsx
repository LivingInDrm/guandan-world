import React, { useState } from 'react';
import type { Story } from '@ladle/react';
import TributeControlPanel from './TributeControlPanel';
import type { Card } from '../../../types';

function createCard(rank: number, suit: number, deckIndex: number): Card {
  return { rank, suit, deckIndex };
}

const mockHandCards: Card[] = [
  createCard(3, 2, 0),
  createCard(5, 0, 1),
  createCard(7, 1, 2),
  createCard(9, 3, 3),
  createCard(11, 0, 4),
  createCard(13, 1, 5),
];

const mockHandWithBomb: Card[] = [
  createCard(3, 2, 0),
  createCard(9, 0, 1),
  createCard(9, 1, 2),
  createCard(9, 2, 3),
  createCard(9, 3, 4),
  createCard(13, 1, 5),
];

const Wrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="min-h-screen bg-gray-800 p-4">{children}</div>
);

export const Default: Story = () => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);

  return (
    <Wrapper>
      <TributeControlPanel
        cards={mockHandCards}
        selectedCards={selectedCards}
        onCardSelect={setSelectedCards}
        currentLevel={2}
        canReturnTribute={true}
        turnDeadlineAtMs={Date.now() + 30000}
        onReturnTribute={(deckIndex) => console.log('Return tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
Default.meta = {
  description: '默认状态 - 可还贡，未选中牌',
};

export const WithSelectedCard: Story = () => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([mockHandCards[0]]);

  return (
    <Wrapper>
      <TributeControlPanel
        cards={mockHandCards}
        selectedCards={selectedCards}
        onCardSelect={setSelectedCards}
        currentLevel={2}
        canReturnTribute={true}
        turnDeadlineAtMs={Date.now() + 30000}
        onReturnTribute={(deckIndex) => console.log('Return tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
WithSelectedCard.meta = {
  description: '已选中一张牌 - 可点击确认还贡',
};

export const Disabled: Story = () => {
  return (
    <Wrapper>
      <TributeControlPanel
        cards={mockHandCards}
        selectedCards={[]}
        onCardSelect={() => {}}
        currentLevel={2}
        canReturnTribute={false}
        turnDeadlineAtMs={Date.now() + 30000}
        onReturnTribute={(deckIndex) => console.log('Return tribute:', deckIndex)}
        disabled={true}
      />
    </Wrapper>
  );
};
Disabled.meta = {
  description: '禁用状态 - 按钮不可点击',
};

export const Interactive: Story = () => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);

  const handleReturnTribute = (deckIndex: number) => {
    console.log('Return tribute:', deckIndex);
    alert(`还贡成功！deckIndex: ${deckIndex}`);
    setSelectedCards([]);
  };

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="text-white text-center">
          点击「提示」自动选牌，点击手牌手动选择，选中一张后点击「确认还贡」
        </div>
        <TributeControlPanel
          cards={mockHandCards}
          selectedCards={selectedCards}
          onCardSelect={setSelectedCards}
          currentLevel={2}
          canReturnTribute={true}
          turnDeadlineAtMs={Date.now() + 60000}
          onReturnTribute={handleReturnTribute}
        />
      </div>
    </Wrapper>
  );
};
Interactive.meta = {
  description: '交互演示 - 可点击提示、选牌、确认还贡',
};

export const WithBombs: Story = () => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="text-white text-center">
          手牌含4张9（炸弹），点击「提示」应优先选择非炸弹牌（梅花3或红桃K）
        </div>
        <TributeControlPanel
          cards={mockHandWithBomb}
          selectedCards={selectedCards}
          onCardSelect={setSelectedCards}
          currentLevel={2}
          canReturnTribute={true}
          turnDeadlineAtMs={Date.now() + 60000}
          onReturnTribute={(deckIndex) => console.log('Return tribute:', deckIndex)}
        />
      </div>
    </Wrapper>
  );
};
WithBombs.meta = {
  description: '手牌含炸弹 - 验证提示优先选非炸弹牌',
};
