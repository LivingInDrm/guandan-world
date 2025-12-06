import React, { useMemo } from 'react';
import type { Player } from '../../types';
import type { PlayAction } from '../../types/proto';
import TeamLevelDisplay from './TeamLevelDisplay';
import Countdown from './Countdown';
import PlayedCardsDisplay from './PlayedCardsDisplay';

import avatar001 from '../../assets/avatar_set/avatar_001.jpg';
import avatar002 from '../../assets/avatar_set/avatar_002.jpg';
import avatar003 from '../../assets/avatar_set/avatar_003.jpg';
import avatar004 from '../../assets/avatar_set/avatar_004.jpg';
import avatar005 from '../../assets/avatar_set/avatar_005.jpg';
import avatar006 from '../../assets/avatar_set/avatar_006.jpg';
import avatar007 from '../../assets/avatar_set/avatar_007.jpg';
import avatar008 from '../../assets/avatar_set/avatar_008.jpg';
import avatar009 from '../../assets/avatar_set/avatar_009.jpg';
import avatar010 from '../../assets/avatar_set/avatar_010.jpg';
import avatar011 from '../../assets/avatar_set/avatar_011.jpg';
import avatar012 from '../../assets/avatar_set/avatar_012.jpg';
import avatar013 from '../../assets/avatar_set/avatar_013.jpg';
import avatar014 from '../../assets/avatar_set/avatar_014.jpg';
import avatar015 from '../../assets/avatar_set/avatar_015.jpg';
import avatar016 from '../../assets/avatar_set/avatar_016.jpg';
import avatar017 from '../../assets/avatar_set/avatar_017.jpg';
import avatar018 from '../../assets/avatar_set/avatar_018.jpg';
import avatar019 from '../../assets/avatar_set/avatar_019.jpg';
import avatar020 from '../../assets/avatar_set/avatar_020.jpg';

const avatars = [
  avatar001, avatar002, avatar003, avatar004, avatar005,
  avatar006, avatar007, avatar008, avatar009, avatar010,
  avatar011, avatar012, avatar013, avatar014, avatar015,
  avatar016, avatar017, avatar018, avatar019, avatar020,
];

const getAvatarByUsername = (username: string): string => {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = ((hash << 5) - hash) + username.charCodeAt(i);
    hash = hash & hash;
  }
  return avatars[Math.abs(hash) % avatars.length];
};

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
  isCurrentTurn: boolean;
  deadlineAtMs?: number;
}

const PlayerArea: React.FC<PlayerAreaProps> = ({ 
  player, 
  position, 
  isCurrentTurn,
  deadlineAtMs
}) => {
  const avatar = useMemo(() => {
    return player ? getAvatarByUsername(player.username) : null;
  }, [player?.username]);

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

  if (!player) {
    return (
      <div className={`${getPositionClasses()}`}>
        <div className="flex flex-col items-center px-3 py-2 rounded-xl bg-white/70 backdrop-blur-sm shadow-sm w-fit">
          <div className="w-14 h-14 rounded-lg bg-gray-200 shadow-sm ring-1 ring-white/70" />
          <span className="mt-1 text-xs font-medium text-slate-400">空座位</span>
        </div>
      </div>
    );
  }

  return (
    <div className={`${getPositionClasses()}`}>
      <div className="relative">
        <div className={`flex flex-col items-center px-3 py-2 rounded-xl bg-white/70 backdrop-blur-sm shadow-sm w-fit ${
          isCurrentTurn ? 'ring-2 ring-yellow-400' : ''
        }`}>
          <img
            src={avatar!}
            alt={player.username}
            className="w-14 h-14 rounded-lg object-cover shadow-sm ring-1 ring-white/70"
          />
          <span className="mt-1 max-w-[80px] text-xs font-medium text-slate-800 truncate text-center">
            {player.username}
          </span>
        </div>
        {isCurrentTurn && deadlineAtMs && (
          <div className={`absolute top-1/2 -translate-y-1/2 ${position === 'right' ? 'right-full mr-2' : 'left-full ml-2'}`}>
            <Countdown deadlineAtMs={deadlineAtMs} size="small" />
          </div>
        )}
      </div>
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
  turnDeadline
}) => {
  const getPlayForSeat = (seat: number): PlayAction | null => {
    if (currentTurn === seat) {
      return null;
    }
    
    const playerPlays = plays.filter(p => p.playerSeat === seat);
    return playerPlays.length > 0 ? playerPlays[playerPlays.length - 1] : null;
  };

  return (
    <div 
      className="relative w-full h-[30rem] rounded-lg overflow-hidden"
      style={{
        background: 'linear-gradient(180deg, #EAF4EF 0%, #DDEEE5 40%, #D2E8DD 100%)',
      }}
    >
      <div 
        className="absolute inset-0 pointer-events-none"
        style={{
          background: 'radial-gradient(circle at center, transparent 60%, rgba(0,0,0,0.06) 100%)',
        }}
      />
      <TeamLevelDisplay 
        teamLevels={teamLevels} 
        currentLevel={currentLevel}
        currentPlayerSeat={currentPlayerSeat}
      />
      
      <div className="absolute inset-0 flex items-center justify-center">
        <div 
          className="relative w-[420px] h-[240px] rounded-[20px]"
          style={{
            background: 'rgba(140, 170, 150, 0.18)',
            backdropFilter: 'blur(10px)',
            border: '1.5px solid rgba(80, 110, 90, 0.35)',
            boxShadow: 'inset 0 0 12px rgba(0,0,0,0.08)',
            padding: '24px',
          }}
        >
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
        player={players[(currentPlayerSeat + 1) % 4]}
        position="left"
        isCurrentTurn={currentTurn === (currentPlayerSeat + 1) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 1) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 2) % 4]}
        position="top"
        isCurrentTurn={currentTurn === (currentPlayerSeat + 2) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 2) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
      
      <PlayerArea
        player={players[(currentPlayerSeat + 3) % 4]}
        position="right"
        isCurrentTurn={currentTurn === (currentPlayerSeat + 3) % 4}
        deadlineAtMs={turnDeadline?.playerSeat === (currentPlayerSeat + 3) % 4 ? turnDeadline.deadlineAtMs : undefined}
      />
    </div>
  );
};

export default GameBoard;