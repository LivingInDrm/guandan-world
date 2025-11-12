import React from 'react';
import type { GameState, Player, TrickInfo, PlayAction } from '../../types';
import { PlayerStatus } from '../../types';
import CardDisplay from './CardDisplay';

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
  isCurrentTurn: boolean;
}

const PlayerArea: React.FC<PlayerAreaProps> = ({ 
  player, 
  position, 
  status,
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
      </div>
    </div>
  );
};

interface PlayedCardsDisplayProps {
  play: PlayAction | null;
  position: 'bottom' | 'left' | 'top' | 'right';
}

const PlayedCardsDisplay: React.FC<PlayedCardsDisplayProps> = ({ play, position }) => {
  const getPositionClasses = () => {
    const base = 'absolute flex items-center justify-center';
    switch (position) {
      case 'bottom':
        return `${base} bottom-2 left-1/2 transform -translate-x-1/2`;
      case 'top':
        return `${base} top-2 left-1/2 transform -translate-x-1/2`;
      case 'left':
        return `${base} left-2 top-1/2 transform -translate-y-1/2`;
      case 'right':
        return `${base} right-2 top-1/2 transform -translate-y-1/2`;
      default:
        return base;
    }
  };

  const getFlexDirection = () => {
    return position === 'left' || position === 'right' ? 'flex-col' : 'flex-row';
  };

  if (!play) {
    return <div className={getPositionClasses()} />;
  }

  if (play.is_pass) {
    return (
      <div className={getPositionClasses()}>
        <div className="bg-gray-200 px-3 py-1 rounded text-sm text-gray-600 font-medium">
          不出
        </div>
      </div>
    );
  }

  if (play.cards && play.cards.length > 0) {
    return (
      <div className={getPositionClasses()}>
        <div className={`flex ${getFlexDirection()} items-center gap-0.5`}>
          {play.cards.map((card, index) => (
            <CardDisplay
              key={card.id}
              card={card}
              size="small"
              stackIndex={index}
            />
          ))}
        </div>
      </div>
    );
  }

  return <div className={getPositionClasses()} />;
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

  // Provide default values if teamLevels is undefined
  const safeTeamLevels = teamLevels || [2, 2];

  return (
    <div className="absolute top-4 left-4 bg-white border border-gray-300 rounded-lg p-3 shadow-sm">
      <div className="text-sm font-medium mb-2">等级信息</div>
      <div className="space-y-1">
        <div className="flex items-center justify-between">
          <span className="text-xs text-gray-600">队伍1 (座位0,2):</span>
          <span className="font-medium text-blue-600">{getLevelText(safeTeamLevels[0])}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-xs text-gray-600">队伍2 (座位1,3):</span>
          <span className="font-medium text-red-600">{getLevelText(safeTeamLevels[1])}</span>
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
    
    // 优先检查是否轮到该玩家
    if (trickInfo.current_turn === seat) {
      return PlayerStatus.PLAYING;
    }
    
    // 然后检查最后一次出牌记录
    const playerPlays = trickInfo.plays.filter(p => p.player_seat === seat);
    const lastPlay = playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
    
    if (lastPlay) {
      return lastPlay.is_pass ? PlayerStatus.PASSED : PlayerStatus.PLAYED;
    }
    
    // 默认等待状态
    return PlayerStatus.WAITING;
  };

  const getPlayForSeat = (seat: number): PlayAction | null => {
    if (!trickInfo) return null;
    
    // 轮到该玩家时，清空显示（等待玩家做决策）
    if (trickInfo.current_turn === seat) {
      return null;
    }
    
    // 显示该玩家最后一次出牌
    const playerPlays = trickInfo.plays.filter(p => p.player_seat === seat);
    return playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
  };

  // Extract data from nested structure: gameState.current_match
  const currentMatch = (gameState as any).current_match;
  const teamLevels = currentMatch?.team_levels || [2, 2];
  const currentLevel = currentMatch?.current_deal?.level || 2;

  return (
    <div className="relative w-full h-96 bg-green-100 border border-gray-300 rounded-lg">
      {/* Team Level Display */}
      <TeamLevelDisplay 
        teamLevels={teamLevels} 
        currentLevel={currentLevel} 
      />
      
      {/* Central Play Area */}
      <div className="absolute inset-0 flex items-center justify-center">
        <div className="relative w-64 h-40 bg-green-200 border-2 border-green-400 rounded-lg">
          <PlayedCardsDisplay 
            play={getPlayForSeat(currentPlayerSeat)}
            position="bottom"
          />
          <PlayedCardsDisplay 
            play={getPlayForSeat((currentPlayerSeat + 1) % 4)}
            position="left"
          />
          <PlayedCardsDisplay 
            play={getPlayForSeat((currentPlayerSeat + 2) % 4)}
            position="top"
          />
          <PlayedCardsDisplay 
            play={getPlayForSeat((currentPlayerSeat + 3) % 4)}
            position="right"
          />
        </div>
      </div>
      
      {/* Player Areas - positioned around the board */}
      <PlayerArea
        player={players[currentPlayerSeat]}
        position="bottom"
        status={getPlayerStatus(currentPlayerSeat)}
        isCurrentTurn={trickInfo?.current_turn === currentPlayerSeat}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 1) % 4]}
        position="left"
        status={getPlayerStatus((currentPlayerSeat + 1) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 1) % 4}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 2) % 4]}
        position="top"
        status={getPlayerStatus((currentPlayerSeat + 2) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 2) % 4}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 3) % 4]}
        position="right"
        status={getPlayerStatus((currentPlayerSeat + 3) % 4)}
        isCurrentTurn={trickInfo?.current_turn === (currentPlayerSeat + 3) % 4}
      />
    </div>
  );
};

export default GameBoard;