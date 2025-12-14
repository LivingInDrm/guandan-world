import React from 'react';
import type { Card } from '../../../types';
import PlayerHand from '../PlayerHand';
import TributeControls from './TributeControls';

interface TributeControlPanelProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  currentLevel?: number;
  canReturnTribute: boolean;
  turnDeadlineAtMs: number;
  onReturnTribute: (deckIndex: number) => void;
  disabled?: boolean;
}

const TributeControlPanel: React.FC<TributeControlPanelProps> = ({
  cards,
  selectedCards,
  onCardSelect,
  currentLevel,
  canReturnTribute,
  turnDeadlineAtMs,
  onReturnTribute,
  disabled = false,
}) => {
  const handleReturnTribute = () => {
    if (selectedCards.length === 1) {
      onReturnTribute(selectedCards[0].deckIndex);
    }
  };

  return (
    <div className="flex flex-col">
      <div className="flex justify-center">
        <TributeControls
          selectedCards={selectedCards}
          canReturnTribute={canReturnTribute}
          turnDeadlineAtMs={turnDeadlineAtMs}
          onReturnTribute={handleReturnTribute}
          onHint={onCardSelect}
          handCards={cards}
          dealLevel={currentLevel ?? 2}
          disabled={disabled}
        />
      </div>
      <PlayerHand
        cards={cards}
        selectedCards={selectedCards}
        onCardSelect={onCardSelect}
        currentLevel={currentLevel}
      />
    </div>
  );
};

export default TributeControlPanel;
