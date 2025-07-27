import React from 'react';
import { GameState, PlayerStatus, Player, TrickInfo } from '../../types';

interface GameBoardProps {
  gameState: GameState;
  players: (Player | null)[];
  currentPlayerSeat: number;
  trickInfo: TrickInfo | null;
}

interface PlayerAreaProps {
  player: Player | null;
  position: 'bottom' | 'left' | 'top' | 'right';
  status: PlayerStatus;
  lastPlay: any; // Cards played in current trick
  isCurrentTurn: boolean;
}

const PlayerArea: React.FC<PlayerAreaProps> = ({ 
  player, 
  position, 
  status, 
  lastPlay, 
  isCurrentTurn 
}) => {
  const getPositionClasses = () => {
    switch (position) {
      case 'bottom':
        return 'absolute bottom-4 left-1/2 transform -translate-x-1/2';
      case 'left':
        return 'absolute left-4 top-1/2 transform -translate-y-1/2 -rotate-90';
      case 'top':
        return 'absolute top-4 left-1/2 transform -translate-x-1/2 rotate-180';
      case 'right':
        return 'absolute right-4 top-1/2 transform -translate-y-1/2 rotate-90';
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
    <div className={`${getPositionClasses()} w-32 h-20`}>
      <div className={`border-2 rounded-lg p-2 text-center ${
        isCurrentTurn ? 'border-yellow-400 shadow-lg' : 'border-gray-300'
      }`}>
        <div className="text-sm font-medium truncate">{player.username}</div>
        <div className={`text-xs px-2 py-1 rounded mt-1 ${getStatusColor()}`}>
          {getStatusText()}
        </div>
        {lastPlay && lastPlay.cards && lastPlay.cards.length > 0 && (
          <div className="text-xs text-gray-500 mt-1">
            出牌: {lastPlay.cards.length}张
          </div>
        )}
      </div>
    </div>
  );
};

interface TeamLevelDisplayProps {
  teamLevels: [number, number];
  currentLevel: number;
}

const TeamLevelDisplay: React.FC<TeamLevelDisplayProps> = ({ teamLevels, currentLevel }) => {
  const getLevelText = (level: number) => {
    if (level <= 10) return level.toString();
    switch (level) {
      case 11: return 'J';
      case 12: return 'Q';
      case 13: return 'K';
      case 14: return 'A';
      default: return level.toString();
    }
  };

  return (
    <div className="absolute top-4 left-4 bg-white border border-gray-300 rounded-lg p-3 shadow-sm">
      <div className="text-sm font-medium mb-2">等级信息</div>
      <div className="space-y-1">
        <div className="flex items-center justify-between">
          <span className="text-xs text-gray-600">队伍1 (座位0,2):</span>
          <span className="font-medium text-blue-600">{getLevelText(teamLevels[0])}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-xs text-gray-600">队伍2 (座位1,3):</span>
          <span className="font-medium text-red-600">{getLevelText(teamLevels[1])}</span>
        </div>
        <div className="border-t pt-1 mt-2">
          <div className="flex items-center justify-between">
            <span className="text-xs text-gray-600">本局等级:</span>
            <span className="font-medium text-green-600">{getLevelText(currentLevel)}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

const GameBoard: React.FC<GameBoardProps> = ({ 
  gameState, 
  players, 
  currentPlayerSeat, 
  trickInfo 
}) => {
  const getPlayerStatus = (seat: number): PlayerStatus => {
    if (!trickInfo) return PlayerStatus.WAITING;
    
    const play = trickInfo.plays.find(p => p.player_seat === seat);
    if (play) {
      return play.is_pass ? PlayerStatus.PASSED : PlayerStatus.PLAYED;
    }
    
    if (trickInfo.current_turn === seat) {
      return PlayerStatus.PLAYING;
    }
    
    return PlayerStatus.WAITING;
  };

  const getLastPlay = (seat: number) => {
    if (!trickInfo) return null;
    return trickInfo.plays.find(p => p.player_seat === seat) || null;
  };

  return (
    <div className="relative w-full h-96 bg-green-100 border border-gray-300 rounded-lg">
      {/* Team Level Display */}
      <TeamLevelDisplay 
        teamLevels={gameState.team_levels} 
        currentLevel={gameState.current_deal.level} 
      />
      
      {/* Central Play Area */}
      <div className="absolute inset-0 flex items-center justify-center">
        <div className="w-48 h-32 bg-green-200 border-2 border-green-400 rounded-lg flex items-center justify-center">
          <div className="text-center">
            <div className="text-sm text-gray-600 mb-1">出牌区</div>
            {trickInfo && trickInfo.plays.length > 0 && (
              <div className="text-xs text-gray-500">
                当前轮次: {trickInfo.plays.length}/4
              </div>
            )}
          </div>
        </div>
      </div>
      
      {/* Player Areas - positioned around the board */}
      <PlayerArea
        player={players[currentPlayerSeat]}
        position="bottom"
        status={getPlayerStatus(currentPlayerSeat)}
        lastPlay={getLastPlay(currentPlayerSeat)}
        isCurrentTurn={trickInfo?.current_turn === currentPlayerSeat}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 1) % 4]}
        position="left"
        status={getPlayerStatus((currentPlayerSeat + 1) % 4)}
        lastPlay={getLastPlay((currentPlayerSeat + 1) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 1) % 4}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 2) % 4]}
        position="top"
        status={getPlayerStatus((currentPlayerSeat + 2) % 4)}
        lastPlay={getLastPlay((currentPlayerSeat + 2) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 2) % 4}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 3) % 4]}
        position="right"
        status={getPlayerStatus((currentPlayerSeat + 3) % 4)}
        lastPlay={getLastPlay((currentPlayerSeat + 3) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 3) % 4}
      />
    </div>
  );
};

export default GameBoard;