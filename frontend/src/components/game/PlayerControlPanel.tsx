import React from 'react';
import type { Card, PlayAction } from '../../types';
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
  onHint: (cards: Card[]) => void;
  plays: PlayAction[];
  leader?: number;
  playerSeat: number;
  dealLevel: number;
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
  onHint,
  plays,
  leader,
  playerSeat,
  dealLevel,
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
          onHint={onHint}
          handCards={cards}
          plays={plays}
          leader={leader}
          playerSeat={playerSeat}
          dealLevel={dealLevel}
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
