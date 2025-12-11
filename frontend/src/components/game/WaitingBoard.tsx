import React from 'react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import GameTable from './GameTable';
import PlayerCard from './PlayerCard';
import { Button, Card } from '../ui-next';

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
              <div className={`w-2 h-2 rounded-full ${player.online ? 'bg-[hsl(var(--ds-primitive-success-500))]' : 'bg-[hsl(var(--ds-primitive-neutral-500))]'}`} />
              <span className="text-xs text-[hsl(var(--ds-color-text-secondary))]">
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
      <Button
        onClick={onStartGame}
        disabled={!canStartGame || isStarting}
        size="lg"
        className="px-8"
      >
        {isStarting 
          ? (isRoomOwner ? '开始中...' : '等待中...') 
          : (isRoomOwner ? '开始游戏' : '等待开始')
        }
      </Button>

      <Button
        onClick={onLeaveRoom}
        disabled={isLeaving}
        intent="secondary"
        size="sm"
      >
        {isLeaving ? '离开中...' : '离开房间'}
      </Button>
    </div>
  );

  const renderTopLeft = () => (
    <Card variant="elevated" interactive={false} className="absolute top-4 left-4 z-10 p-3">
      <div className="text-sm text-[hsl(var(--ds-color-text-secondary))]">
        房间ID: <span className="font-mono text-[hsl(var(--ds-color-text-primary))]">{roomId}</span>
      </div>
    </Card>
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
