import React from 'react';
import { RoomStatus, type RoomInfo } from '../../types';
import { Card, Button } from '../ui';

interface RoomCardProps {
  room: RoomInfo;
  onJoinRoom: (roomId: string) => void;
  currentUserId: string;
}

const RoomCard: React.FC<RoomCardProps> = ({ room, onJoinRoom, currentUserId }) => {
  const getStatusText = (status: RoomStatus) => {
    switch (status) {
      case RoomStatus.WAITING:
        return '等待中';
      case RoomStatus.READY:
        return '准备中';
      case RoomStatus.PLAYING:
        return '游戏中';
      case RoomStatus.CLOSED:
        return '已关闭';
      default:
        return '未知';
    }
  };

  const getStatusColor = (status: RoomStatus) => {
    switch (status) {
      case RoomStatus.WAITING:
        return 'bg-primary/20 text-primary';
      case RoomStatus.READY:
        return 'bg-accent-light text-amber-700';
      case RoomStatus.PLAYING:
        return 'bg-team-us/20 text-team-us-dark';
      case RoomStatus.CLOSED:
        return 'bg-disabled-bg text-disabled-text';
      default:
        return 'bg-disabled-bg text-disabled-text';
    }
  };

  const isUserInRoom = room.players.some(player => player.id === currentUserId);
  const canJoin = room.can_join && !isUserInRoom;
  const isOwner = room.owner === currentUserId;

  const handleJoinClick = () => {
    if (canJoin) {
      onJoinRoom(room.id);
    }
  };

  const handleRoomClick = () => {
    onJoinRoom(room.id);
  };

  return (
    <Card variant="default" className="p-5 hover:shadow-card-hover transition-shadow">
      <div className="flex justify-between items-start mb-4">
        <div className="flex-1">
          <div className="flex items-center space-x-2 mb-2">
            <h3 className="font-medium text-foreground">房间 #{room.id.slice(-6)}</h3>
            {isOwner && (
              <span className="inline-flex items-center px-2 py-0.5 rounded-sm text-xs font-medium bg-accent-light text-amber-700">
                房主
              </span>
            )}
          </div>
          <div className="flex items-center space-x-3">
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(room.status)}`}>
              {getStatusText(room.status)}
            </span>
            <span className="text-sm text-table-400">
              {room.player_count}/4 人
            </span>
          </div>
        </div>
      </div>

      <div className="mb-4">
        <h4 className="text-sm font-medium text-table-400 mb-2">玩家列表</h4>
        <div className="grid grid-cols-2 gap-2">
          {Array.from({ length: 4 }).map((_, index) => {
            const player = room.players[index];
            return (
              <div
                key={index}
                className={`p-2 rounded-sm text-sm text-center ${
                  player
                    ? player.id === currentUserId
                      ? 'bg-primary/10 text-primary border border-primary/30'
                      : 'bg-table-50 text-table-400'
                    : 'bg-table-100 text-table-300 border-2 border-dashed border-table-200'
                }`}
              >
                {player ? (
                  <div className="flex items-center justify-center space-x-1">
                    <span>{player.username}</span>
                    {player.online ? (
                      <div className="w-2 h-2 bg-badge-level rounded-full"></div>
                    ) : (
                      <div className="w-2 h-2 bg-disabled-text rounded-full"></div>
                    )}
                  </div>
                ) : (
                  '等待玩家'
                )}
              </div>
            );
          })}
        </div>
      </div>

      <div className="flex justify-end">
        {isUserInRoom ? (
          <Button variant="secondary" size="sm" onClick={handleRoomClick}>
            返回房间
          </Button>
        ) : canJoin ? (
          <Button variant="primary" size="sm" onClick={handleJoinClick}>
            加入房间
          </Button>
        ) : (
          <Button variant="ghost" size="sm" disabled>
            {room.status === RoomStatus.PLAYING ? '游戏中' : '房间已满'}
          </Button>
        )}
      </div>
    </Card>
  );
};

export default RoomCard;
