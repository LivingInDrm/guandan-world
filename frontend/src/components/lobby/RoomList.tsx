import React from 'react';
import { RoomStatus, type RoomInfo } from '../../types';
import { Button } from '../ui-next';
import PlayerMiniCard from './PlayerMiniCard';
import Pagination from './Pagination';

interface RoomListProps {
  rooms: RoomInfo[];
  isLoading: boolean;
  currentPage: number;
  totalCount: number;
  limit: number;
  onPageChange: (page: number) => void;
  onJoinRoom: (roomId: string) => void;
  currentUserId: string;
}

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
      return 'bg-action-primary/20 text-action-primary';
    case RoomStatus.READY:
      return 'bg-action-secondary/20 text-action-secondary';
    case RoomStatus.PLAYING:
      return 'bg-team-us/20 text-team-us';
    case RoomStatus.CLOSED:
      return 'bg-state-disabled/20 text-fg-secondary';
    default:
      return 'bg-state-disabled/20 text-fg-secondary';
  }
};

const getStatusDot = (status: RoomStatus) => {
  switch (status) {
    case RoomStatus.WAITING:
      return 'bg-action-primary animate-pulse';
    case RoomStatus.READY:
      return 'bg-action-secondary';
    case RoomStatus.PLAYING:
      return 'bg-team-us';
    case RoomStatus.CLOSED:
      return 'bg-state-disabled';
    default:
      return 'bg-state-disabled';
  }
};

const TableHeader = () => (
  <thead className="bg-surface-elevated/50">
    <tr>
      <th className="px-4 py-3 text-center text-sm font-medium text-fg-secondary">房间</th>
      <th className="px-4 py-3 text-center text-sm font-medium text-fg-secondary">状态</th>
      <th className="px-4 py-3 text-center text-sm font-medium text-fg-secondary">玩家</th>
      <th className="px-4 py-3 text-center text-sm font-medium text-fg-secondary">人数</th>
      <th className="px-4 py-3 text-center text-sm font-medium text-fg-secondary">操作</th>
    </tr>
  </thead>
);

const RoomList: React.FC<RoomListProps> = ({
  rooms,
  isLoading,
  currentPage,
  totalCount,
  limit,
  onPageChange,
  onJoinRoom,
  currentUserId
}) => {
  const sortedRooms = [...rooms].sort((a, b) => {
    const statusPriority = {
      [RoomStatus.WAITING]: 0,
      [RoomStatus.READY]: 1,
      [RoomStatus.PLAYING]: 2,
      [RoomStatus.CLOSED]: 3
    };
    
    const statusDiff = statusPriority[a.status] - statusPriority[b.status];
    if (statusDiff !== 0) {
      return statusDiff;
    }
    
    return b.player_count - a.player_count;
  });

  if (isLoading && rooms.length === 0) {
    return (
      <div className="bg-surface-base rounded-lg shadow-elevation-1 border border-stroke overflow-x-auto">
        <table className="w-full min-w-[640px]">
          <TableHeader />
          <tbody>
            {Array.from({ length: 6 }).map((_, index) => (
              <tr key={index} className="border-t border-stroke animate-pulse">
                <td className="px-4 py-4 text-center"><div className="h-4 bg-surface-elevated rounded w-20 mx-auto"></div></td>
                <td className="px-4 py-4 text-center"><div className="h-5 bg-surface-elevated rounded w-16 mx-auto"></div></td>
                <td className="px-4 py-4">
                  <div className="flex justify-center gap-3">
                    {Array.from({ length: 4 }).map((_, i) => (
                      <div key={i} className="flex flex-col items-center gap-1">
                        <div className="h-10 w-10 bg-surface-elevated rounded-full"></div>
                        <div className="h-3 bg-surface-elevated rounded w-10"></div>
                      </div>
                    ))}
                  </div>
                </td>
                <td className="px-4 py-4 text-center"><div className="h-4 bg-surface-elevated rounded w-8 mx-auto"></div></td>
                <td className="px-4 py-4 text-center"><div className="h-8 bg-surface-elevated rounded w-16 mx-auto"></div></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  if (rooms.length === 0 && !isLoading) {
    return (
      <div className="text-center py-12">
        <div className="text-fg-secondary text-6xl mb-4">&#127918;</div>
        <h3 className="text-lg font-medium text-fg-primary mb-2">暂无房间</h3>
        <p className="text-fg-secondary mb-6">成为第一个创建房间的玩家吧！</p>
      </div>
    );
  }

  const totalPages = Math.ceil(totalCount / limit);

  const getPlayersWithSlots = (room: RoomInfo) => {
    const slots: (typeof room.players[0] | null)[] = [null, null, null, null];
    room.players.forEach(player => {
      if (player && player.seat >= 0 && player.seat < 4) {
        slots[player.seat] = player;
      }
    });
    return slots;
  };

  return (
    <div className="space-y-6">
      <div className="bg-surface-base rounded-lg shadow-elevation-1 border border-stroke overflow-x-auto">
        <table className="w-full min-w-[640px]">
          <TableHeader />
          <tbody>
            {sortedRooms.map((room) => {
              const isUserInRoom = room.players.some(player => player.id === currentUserId);
              const canJoin = room.can_join && !isUserInRoom;
              const playersWithSlots = getPlayersWithSlots(room);

              return (
                <tr key={room.id} className="border-t border-stroke hover:bg-surface-elevated/30 transition-colors">
                  <td className="px-4 py-4 text-center">
                    <span className="font-mono text-sm text-fg-primary">#{room.id.slice(-6)}</span>
                  </td>
                  <td className="px-4 py-4 text-center">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(room.status)}`}>
                      <span className={`w-1.5 h-1.5 rounded-full ${getStatusDot(room.status)}`}></span>
                      {getStatusText(room.status)}
                    </span>
                  </td>
                  <td className="px-4 py-4">
                    <div className="flex justify-center gap-3">
                      {playersWithSlots.map((player, index) => (
                        <PlayerMiniCard
                          key={index}
                          player={player}
                          isCurrentUser={player?.id === currentUserId}
                        />
                      ))}
                    </div>
                  </td>
                  <td className="px-4 py-4 text-center">
                    <span className="text-sm text-fg-secondary">{room.player_count}/4</span>
                  </td>
                  <td className="px-4 py-4 text-center">
                    {isUserInRoom ? (
                      <Button intent="neutral" size="sm" onClick={() => onJoinRoom(room.id)}>
                        返回
                      </Button>
                    ) : canJoin ? (
                      <Button intent="primary" size="sm" onClick={() => onJoinRoom(room.id)}>
                        加入
                      </Button>
                    ) : (
                      <Button intent="neutral" size="sm" disabled>
                        加入
                      </Button>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {isLoading && rooms.length > 0 && (
        <div className="text-center py-4">
          <div className="inline-flex items-center text-fg-secondary">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-action-primary mr-2"></div>
            更新中...
          </div>
        </div>
      )}

      {totalPages > 1 && (
        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={onPageChange}
        />
      )}
    </div>
  );
};

export default RoomList;
