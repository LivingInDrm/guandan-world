import React from 'react';
import { Play, LogOut, Loader2, Users } from 'lucide-react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import GameTable from './GameTable';
import PlayerCard from './PlayerCard';
import RoomInfoPanel from './RoomInfoPanel';
import { Button } from '../ui';
import { cn } from '@/lib/utils';

interface WaitingBoardProps {
  players: (Player | null)[];
  currentPlayerSeat: number;
  roomId: string;
  roomCode: string;
  ownerId: string;
  currentUserId: string;
  isConnected: boolean;
  onStartGame: () => void;
  onLeaveRoom: () => void;
  isStarting: boolean;
  isLeaving: boolean;
  className?: string;
}

const WaitingBoard: React.FC<WaitingBoardProps> = ({
  players,
  currentPlayerSeat,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  roomId: _roomId,
  roomCode,
  ownerId,
  currentUserId,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  isConnected: _isConnected,
  onStartGame,
  onLeaveRoom,
  isStarting,
  isLeaving,
  className,
}) => {
  const playerCount = players.filter(p => p !== null).length;
  const isRoomOwner = currentUserId === ownerId;
  const canStartGame = isRoomOwner && playerCount === 4;
  const mySeat = players.findIndex(p => p?.id === currentUserId);
  const effectiveSeat = currentPlayerSeat ?? (mySeat >= 0 ? mySeat : 0);

  const renderPlayer = (player: Player | null, position: PlayerPosition) => {
    const isCurrentUser = player?.id === currentUserId;

    return (
      <PlayerCard
        player={player}
        position={position}
        isOwner={player?.id === ownerId}
        isHighlighted={isCurrentUser}
        statusSlot={
          player && (
            <div className="flex items-center gap-1.5">
              <div className={cn(
                "w-2 h-2 rounded-full",
                player.online
                  ? "bg-[hsl(158,55%,42%)] shadow-[0_0_4px_hsla(158,55%,42%,0.6)]"
                  : "bg-state-disabled",
              )} />
              <span className={cn(
                "text-xs",
                player.online ? "text-white/80" : "text-white/50",
              )}>
                {player.online ? '在线' : '离线'}
              </span>
            </div>
          )
        }
      />
    );
  };

  const renderCenter = () => (
    <div className="flex flex-col items-center justify-center h-full gap-5 mobile-landscape:gap-2 relative z-10">
      {/* 人数状态 */}
      <div className={cn(
        "flex items-center gap-2 px-4 py-2 mobile-landscape:px-3 mobile-landscape:py-1.5 rounded-full",
        "bg-black/20 backdrop-blur-sm",
        "border border-white/10",
      )}>
        <Users className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3 text-state-active" />
        <span className="text-white/90 text-sm mobile-landscape:text-xs font-medium">
          {playerCount}/4 位玩家
        </span>
        {playerCount < 4 && (
          <span className="text-white/50 text-xs mobile-landscape:hidden">
            (等待中...)
          </span>
        )}
      </div>

      {/* 开始游戏按钮 */}
      <Button
        onClick={onStartGame}
        disabled={!canStartGame || isStarting}
        size="lg"
        className={cn(
          "px-10 py-3 mobile-landscape:px-6 mobile-landscape:py-2 mobile-landscape:text-sm",
          canStartGame
            ? [
                "bg-gradient-to-r from-state-active to-[hsl(42,95%,45%)]",
                "hover:from-[hsl(42,95%,55%)] hover:to-[hsl(42,95%,48%)]",
                "text-primitive-neutral-900 font-bold",
                "shadow-[var(--elevation-3),var(--shadow-relief)]",
                "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-lg)]",
                "animate-pulse-glow",
              ]
            : [
                "bg-white/10",
                "text-white/60",
                "border border-white/20",
              ],
        )}
      >
        {isStarting ? (
          <>
            <Loader2 className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4 animate-spin" />
            {isRoomOwner ? '开始中...' : '等待中...'}
          </>
        ) : (
          <>
            <Play className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4" />
            {isRoomOwner ? '开始游戏' : '等待房主开始'}
          </>
        )}
      </Button>

      {/* 离开房间按钮 */}
      <Button
        onClick={onLeaveRoom}
        disabled={isLeaving}
        intent="neutral"
        size="sm"
        className={cn(
          "text-white/70 hover:text-white mobile-landscape:text-xs",
          "bg-transparent hover:bg-white/10",
          "border border-white/20 hover:border-white/30",
        )}
      >
        {isLeaving ? (
          <>
            <Loader2 className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3 animate-spin" />
            离开中...
          </>
        ) : (
          <>
            <LogOut className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3" />
            离开房间
          </>
        )}
      </Button>
    </div>
  );

  return (
    <GameTable
      players={players}
      currentPlayerSeat={effectiveSeat}
      renderPlayer={renderPlayer}
      renderCenter={renderCenter}
      topLeftSlot={<RoomInfoPanel roomCode={roomCode} />}
      className={className}
    />
  );
};

export default WaitingBoard;
