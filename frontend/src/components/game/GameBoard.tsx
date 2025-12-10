import React from 'react';
import type { Player } from '../../types';
import type { PlayAction } from '../../types/proto';
import type { PlayerPosition } from './GameTable';
import GameTable from './GameTable';
import PlayerCard from './PlayerCard';
import TeamLevelDisplay from './TeamLevelDisplay';
import Countdown from './Countdown';
import PlayedCardsDisplay from './PlayedCardsDisplay';

interface GameBoardProps {
  teamLevels: [number, number];
  currentLevel: number;
  plays: PlayAction[];
  currentTurn: number;
  players: (Player | null)[];
  currentPlayerSeat: number;
  turnDeadline?: { playerSeat: number; deadlineAtMs: number } | null;
}

const GameBoard: React.FC<GameBoardProps> = ({ 
  teamLevels,
  currentLevel,
  plays,
  currentTurn,
  players, 
  currentPlayerSeat,
  turnDeadline
}) => {
  const getPlayForSeat = (seat: number): PlayAction | null => {
    if (currentTurn === seat) {
      return null;
    }
    const playerPlays = plays.filter(p => p.playerSeat === seat);
    return playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
  };

  const renderPlayer = (player: Player | null, position: PlayerPosition, seat: number) => {
    if (position === 'bottom') {
      return null;
    }

    const isCurrentTurn = currentTurn === seat;
    const deadlineAtMs = isCurrentTurn && 
                         turnDeadline?.playerSeat === seat 
                         ? turnDeadline.deadlineAtMs 
                         : undefined;

    return (
      <PlayerCard
        player={player}
        position={position}
        isHighlighted={isCurrentTurn}
        statusSlot={
          deadlineAtMs ? (
            <Countdown deadlineAtMs={deadlineAtMs} size="small" />
          ) : undefined
        }
      />
    );
  };

  const renderCenter = () => (
    <>
      <PlayedCardsDisplay 
        play={getPlayForSeat(currentPlayerSeat)}
        position="bottom"
        currentLevel={currentLevel}
      />
      <PlayedCardsDisplay 
        play={getPlayForSeat((currentPlayerSeat + 1) % 4)}
        position="left"
        currentLevel={currentLevel}
      />
      <PlayedCardsDisplay 
        play={getPlayForSeat((currentPlayerSeat + 2) % 4)}
        position="top"
        currentLevel={currentLevel}
      />
      <PlayedCardsDisplay 
        play={getPlayForSeat((currentPlayerSeat + 3) % 4)}
        position="right"
        currentLevel={currentLevel}
      />
    </>
  );

  return (
    <GameTable
      players={players}
      currentPlayerSeat={currentPlayerSeat}
      renderPlayer={renderPlayer}
      renderCenter={renderCenter}
      topLeftSlot={
        <TeamLevelDisplay 
          teamLevels={teamLevels} 
          currentLevel={currentLevel}
          currentPlayerSeat={currentPlayerSeat}
        />
      }
    />
  );
};

export default GameBoard;
