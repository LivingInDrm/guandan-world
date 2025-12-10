import React from 'react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import GameTable from './GameTable';
import PlayerCard from './PlayerCard';
import { Button } from '../ui';

interface WaitingBoardProps {
  players: (Player | null)[];
  currentPlayerSeat: number;
  roomId: string;
  ownerId: string;
  currentUserId: string;
  isConnected: boolean;
  onStartGame: () => void;
  onLeaveRoom: () => void;
  isStarting: boolean;
  isLeaving: boolean;
}

const WaitingBoard: React.FC<WaitingBoardProps> = ({
  players,
  currentPlayerSeat,
  roomId,
  ownerId,
  currentUserId,
  isConnected,
  onStartGame,
  onLeaveRoom,
  isStarting,
  isLeaving,
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
            <div className="flex items-center gap-1">
              <div className={`w-2 h-2 rounded-full ${player.online ? 'bg-green-500' : 'bg-muted-foreground'}`} />
              <span className="text-xs text-muted-foreground">
                {player.online ? '在线' : '离线'}
              </span>
            </div>
          )
        }
      />
    );
  };

  const renderCenter = () => (
    <div className="flex flex-col items-center justify-center h-full gap-4">
      {isRoomOwner ? (
        canStartGame ? (
          <Button
            onClick={onStartGame}
            disabled={isStarting}
            size="lg"
            className="px-8"
          >
            {isStarting ? '开始中...' : '开始游戏'}
          </Button>
        ) : (
          <div className="text-center">
            <p className="text-lg font-medium text-foreground mb-1">
              等待玩家加入
            </p>
            <p className="text-sm text-muted-foreground">
              ({playerCount}/4)
            </p>
          </div>
        )
      ) : (
        <p className="text-lg text-muted-foreground">
          等待房主开始游戏...
        </p>
      )}

      <Button
        onClick={onLeaveRoom}
        disabled={isLeaving}
        variant="secondary"
        size="sm"
      >
        {isLeaving ? '离开中...' : '离开房间'}
      </Button>
    </div>
  );

  const renderTopLeft = () => (
    <div className="absolute top-4 left-4 bg-card/80 backdrop-blur-sm rounded-lg p-3 shadow-sm z-10">
      <div className="space-y-1.5 text-sm">
        <div className="text-muted-foreground">
          房间ID: <span className="font-mono text-foreground">{roomId}</span>
        </div>
        <div className="text-muted-foreground">
          玩家数量: <span className="text-foreground">{playerCount}/4</span>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-xs text-muted-foreground">
            {isConnected ? '已连接' : '连接断开'}
          </span>
        </div>
      </div>
    </div>
  );

  return (
    <div className="max-w-6xl mx-auto p-6">
      <GameTable
        players={players}
        currentPlayerSeat={effectiveSeat}
        renderPlayer={renderPlayer}
        renderCenter={renderCenter}
        topLeftSlot={renderTopLeft()}
      />
    </div>
  );
};

export default WaitingBoard;
