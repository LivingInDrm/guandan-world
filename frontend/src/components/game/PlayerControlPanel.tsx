import React from 'react';
import type { Card } from '../../types';
import PlayerHand from './PlayerHand';
import GameControls from './GameControls';

interface PlayerControlPanelProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  currentLevel?: number;
  canPlay: boolean;
  isMyTurn: boolean;
  turnDeadlineAtMs: number;
  onPlayCards: (cards: Card[]) => void;
  onPass: () => void;
  disabled?: boolean;
}

const PlayerControlPanel: React.FC<PlayerControlPanelProps> = ({
  cards,
  selectedCards,
  onCardSelect,
  currentLevel,
  canPlay,
  isMyTurn,
  turnDeadlineAtMs,
  onPlayCards,
  onPass,
  disabled = false,
}) => {
  return (
    <div className="flex flex-col">
      <div className="flex justify-center">
        <GameControls
          selectedCards={selectedCards}
          canPlay={canPlay}
          isMyTurn={isMyTurn}
          turnDeadlineAtMs={turnDeadlineAtMs}
          onPlayCards={onPlayCards}
          onPass={onPass}
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

export default PlayerControlPanel;
