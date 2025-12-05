import React from 'react';
import type { Player } from '../../types';
import type { PlayAction } from '../../types/proto';
import { PlayerStatus } from '../../types';
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
  playStates?: [number, number, number, number];
  turnDeadline?: { playerSeat: number; deadlineAtMs: number } | null;
}

interface PlayerAreaProps {
  player: Player | null;
  position: 'bottom' | 'left' | 'top' | 'right';
  status: PlayerStatus;
  isCurrentTurn: boolean;
  deadlineAtMs?: number;
}

const PlayerArea: React.FC<PlayerAreaProps> = ({ 
  player, 
  position, 
  status,
  isCurrentTurn,
  deadlineAtMs
}) => {
  const getPositionClasses = () => {
    switch (position) {
      case 'bottom':
        return 'absolute bottom-4 left-1/2 transform -translate-x-1/2';
      case 'left':
        return 'absolute left-4 top-1/2 transform -translate-y-1/2';
      case 'top':
        return 'absolute top-4 left-1/2 transform -translate-x-1/2';
      case 'right':
        return 'absolute right-4 top-1/2 transform -translate-y-1/2';
      default:
        return '';
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case PlayerStatus.WAITING:
        return 'bg-gray-200 text-gray-600';
      case PlayerStatus.PLAYING:
        return 'bg-yellow-200 text-yellow-800';
      case PlayerStatus.PLAYED:
        return 'bg-green-200 text-green-800';
      case PlayerStatus.PASSED:
        return 'bg-red-200 text-red-800';
      case PlayerStatus.FINISHED:
        return 'bg-blue-200 text-blue-800';
      default:
        return 'bg-gray-200 text-gray-600';
    }
  };

  const getStatusText = () => {
    switch (status) {
      case PlayerStatus.WAITING:
        return '等待';
      case PlayerStatus.PLAYING:
        return '出牌中';
      case PlayerStatus.PLAYED:
        return '已出牌';
      case PlayerStatus.PASSED:
        return '不出';
      case PlayerStatus.FINISHED:
        return '已结束';
      default:
        return '等待';
    }
  };

  if (!player) {
    return (
      <div className={`${getPositionClasses()} w-24 h-16`}>
        <div className="bg-gray-100 border-2 border-dashed border-gray-300 rounded-lg p-2 text-center">
          <div className="text-sm text-gray-400">空座位</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`${getPositionClasses()} flex items-center gap-2`}>
      <div className={`w-32 border-2 rounded-lg p-2 text-center ${
        isCurrentTurn ? 'border-yellow-400 shadow-lg' : 'border-gray-300'
      }`}>
        <div className="text-sm font-medium truncate">{player.username}</div>
        <div className={`text-xs px-2 py-1 rounded mt-1 ${getStatusColor()}`}>
          {getStatusText()}
        </div>
      </div>
      {isCurrentTurn && deadlineAtMs && (
        <Countdown deadlineAtMs={deadlineAtMs} size="small" />
      )}
    </div>
  );
};

const GameBoard: React.FC<GameBoardProps> = ({ 
  teamLevels,
  currentLevel,
  plays,
  currentTurn,
  players, 
  currentPlayerSeat,
  playStates,
  turnDeadline
}) => {
  const getPlayerStatus = (seat: number): PlayerStatus => {
    if (playStates) {
      const state = playStates[seat];
      if (state === 0 && currentTurn === seat) {
        return PlayerStatus.PLAYING;
      }
      switch (state) {
        case 0: return PlayerStatus.WAITING;
        case 1: return PlayerStatus.PLAYED;
        case 2: return PlayerStatus.PASSED;
        case 3: return PlayerStatus.FINISHED;
      }
    }
    
    if (currentTurn === seat) {
      return PlayerStatus.PLAYING;
    }
    
    const playerPlays = plays.filter(p => p.playerSeat === seat);
    const lastPlay = playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
    
    if (lastPlay) {
      return lastPlay.isPass ? PlayerStatus.PASSED : PlayerStatus.PLAYED;
    }
    
    return PlayerStatus.WAITING;
  };

  const getPlayForSeat = (seat: number): PlayAction | null => {
    if (currentTurn === seat) {
      return null;
    }
    
    const playerPlays = plays.filter(p => p.playerSeat === seat);
    return playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
  };

  return (
    <div className="relative w-full h-96 bg-green-100 border border-gray-300 rounded-lg">
      <TeamLevelDisplay 
        teamLevels={teamLevels} 
        currentLevel={currentLevel}
        currentPlayerSeat={currentPlayerSeat}
      />
      
      <div className="absolute inset-0 flex items-center justify-center">
        <div className="relative w-[420px] h-[240px] bg-green-200 border-2 border-green-400 rounded-lg">
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
        </div>
      </div>
      
      <PlayerArea
        player={players[currentPlayerSeat]}
        position="bottom"
        status={getPlayerStatus(currentPlayerSeat)}
        isCurrentTurn={currentTurn === currentPlayerSeat}
        deadlineAtMs={turnDeadline?.playerSeat === currentPlayerSeat ? turnDeadline.deadlineAtMs : undefined}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 1) % 4]}
        position="left"
        status={getPlayerStatus((currentPlayerSeat + 1) % 4)}
        isCurrentTurn={currentTurn === (currentPlayerSeat + 1) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 1) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 2) % 4]}
        position="top"
        status={getPlayerStatus((currentPlayerSeat + 2) % 4)}
        isCurrentTurn={currentTurn === (currentPlayerSeat + 2) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 2) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 3) % 4]}
        position="right"
        status={getPlayerStatus((currentPlayerSeat + 3) % 4)}
        isCurrentTurn={currentTurn === (currentPlayerSeat + 3) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 3) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
    </div>
  );
};

export default GameBoard;