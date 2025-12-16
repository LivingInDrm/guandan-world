import React from 'react';
import { RoomStatus, type RoomInfo } from '../../types';
import { Button } from '../ui';
import PlayerMiniCard from './PlayerMiniCard';
import Pagination from './Pagination';
import { cn } from '@/lib/utils';

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

const getStatusStyles = (status: RoomStatus) => {
  switch (status) {
    case RoomStatus.WAITING:
      return {
        badge: 'bg-gradient-to-r from-state-active/10 to-state-active/20 text-state-active-muted border border-state-active/25',
        dot: 'bg-state-active shadow-[0_0_6px_hsla(42,95%,52%,0.6)]',
        animate: true,
      };
    case RoomStatus.READY:
      return {
        badge: 'bg-gradient-to-r from-action-primary/10 to-action-primary/15 text-action-primary-dark border border-action-primary/25',
        dot: 'bg-action-primary shadow-[0_0_4px_hsla(158,55%,42%,0.5)]',
        animate: false,
      };
    case RoomStatus.PLAYING:
      return {
        badge: 'bg-gradient-to-r from-[hsl(200,80%,50%)]/10 to-[hsl(200,80%,50%)]/15 text-[hsl(200,80%,40%)] border border-[hsl(200,80%,50%)]/25',
        dot: 'bg-[hsl(200,80%,50%)] shadow-[0_0_4px_hsla(200,80%,50%,0.5)]',
        animate: false,
      };
    case RoomStatus.CLOSED:
      return {
        badge: 'bg-stroke/10 text-fg-secondary/70 border border-stroke/20',
        dot: 'bg-state-disabled',
        animate: false,
      };
    default:
      return {
        badge: 'bg-stroke/10 text-fg-secondary/70 border border-stroke/20',
        dot: 'bg-state-disabled',
        animate: false,
      };
  }
};

const TableHeader = () => (
  <thead>
    <tr className={cn(
      "bg-gradient-to-r from-[hsl(158,55%,32%)] via-[hsl(158,55%,38%)] to-[hsl(158,55%,32%)]",
      "text-white",
    )}>
      <th className="px-4 py-2.5 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center text-sm mobile-landscape:text-xs font-display font-bold tracking-wide">房间</th>
      <th className="px-4 py-2.5 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center text-sm mobile-landscape:text-xs font-display font-bold tracking-wide">状态</th>
      <th className="px-4 py-2.5 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center text-sm mobile-landscape:text-xs font-display font-bold tracking-wide">玩家</th>
      <th className="px-4 py-2.5 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center text-sm mobile-landscape:text-xs font-display font-bold tracking-wide">人数</th>
      <th className="px-4 py-2.5 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center text-sm mobile-landscape:text-xs font-display font-bold tracking-wide">操作</th>
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
      <div className={cn(
        "bg-surface-base rounded-xl overflow-hidden",
        "shadow-[var(--elevation-2),var(--shadow-relief)]",
        "border border-stroke/50",
      )}>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px]">
            <TableHeader />
            <tbody>
              {Array.from({ length: 6 }).map((_, index) => (
                <tr key={index} className="border-t border-stroke/50 animate-pulse">
                  <td className="px-4 py-4 text-center">
                    <div className="h-4 bg-surface-elevated rounded-md w-20 mx-auto" />
                  </td>
                  <td className="px-4 py-4 text-center">
                    <div className="h-6 bg-surface-elevated rounded-full w-16 mx-auto" />
                  </td>
                  <td className="px-4 py-4">
                    <div className="flex justify-center gap-3">
                      {Array.from({ length: 4 }).map((_, i) => (
                        <div key={i} className="flex flex-col items-center gap-1">
                          <div className="h-10 w-10 bg-surface-elevated rounded-full" />
                          <div className="h-3 bg-surface-elevated rounded w-10" />
                        </div>
                      ))}
                    </div>
                  </td>
                  <td className="px-4 py-4 text-center">
                    <div className="h-4 bg-surface-elevated rounded w-8 mx-auto" />
                  </td>
                  <td className="px-4 py-4 text-center">
                    <div className="h-8 bg-surface-elevated rounded-lg w-16 mx-auto" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    );
  }

  if (rooms.length === 0 && !isLoading) {
    return (
      <div className={cn(
        "text-center py-16 px-8 rounded-xl",
        "bg-surface-base/80 backdrop-blur-sm",
        "border border-stroke/30",
        "shadow-[var(--elevation-1)]",
      )}>
        {/* 装饰图案 */}
        <div className="flex items-center justify-center gap-3 mb-6 opacity-40">
          <span className="text-4xl">♠</span>
          <span className="text-4xl text-suit-red">♥</span>
          <span className="text-4xl">♣</span>
          <span className="text-4xl text-suit-red">♦</span>
        </div>
        <h3 className={cn(
          "text-xl font-display font-bold mb-2",
          "bg-gradient-primary",
          "bg-clip-text text-transparent",
        )}>
          暂无房间
        </h3>
        <p className="text-fg-secondary">成为第一个创建房间的玩家吧！</p>
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
    <div className="relative">
      {/* 左右翻页按钮 */}
      <Pagination
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={onPageChange}
        variant="arrows"
      />

      <div className={cn(
        "bg-surface-base rounded-xl mobile-landscape:rounded-lg overflow-hidden",
        "shadow-[var(--elevation-2),var(--shadow-relief)]",
        "border border-stroke/50",
      )}>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] mobile-landscape:min-w-[560px]">
            <TableHeader />
            <tbody>
              {sortedRooms.slice(0, 5).map((room, rowIndex) => {
                const isUserInRoom = room.players.some(player => player.id === currentUserId);
                const canJoin = room.can_join && !isUserInRoom;
                const playersWithSlots = getPlayersWithSlots(room);
                const statusStyles = getStatusStyles(room.status);

                return (
                  <tr
                    key={room.id}
                    className={cn(
                      "border-t border-stroke/30",
                      "hover:bg-surface-elevated/50",
                      "transition-colors duration-fast",
                      rowIndex % 2 === 0 ? "bg-surface-base" : "bg-surface-elevated/20",
                    )}
                  >
                    <td className="px-4 py-2 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center">
                      <span className={cn(
                        "font-mono text-sm mobile-landscape:text-xs font-medium",
                        "text-fg-primary",
                      )}>
                        #{room.id.slice(-6)}
                      </span>
                    </td>
                    <td className="px-4 py-2 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center">
                      <span className={cn(
                        "inline-flex items-center gap-1.5 mobile-landscape:gap-1 px-3 py-1 mobile-landscape:px-2 rounded-full text-xs mobile-landscape:text-[10px] font-medium",
                        statusStyles.badge,
                      )}>
                        <span className={cn("w-1.5 h-1.5 mobile-landscape:w-1 mobile-landscape:h-1 rounded-full", statusStyles.dot)} />
                        {getStatusText(room.status)}
                      </span>
                    </td>
                    <td className="px-4 py-2 mobile-landscape:px-2 mobile-landscape:py-1.5">
                      <div className="flex justify-center gap-3 mobile-landscape:gap-1.5">
                        {playersWithSlots.map((player, index) => (
                          <PlayerMiniCard
                            key={index}
                            player={player}
                            isCurrentUser={player?.id === currentUserId}
                          />
                        ))}
                      </div>
                    </td>
                    <td className="px-4 py-2 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center">
                      <span className={cn(
                        "inline-flex items-center justify-center",
                        "min-w-[2.5rem] mobile-landscape:min-w-[2rem] px-2 mobile-landscape:px-1.5 py-0.5 rounded-md",
                        "bg-surface-elevated/80 text-sm mobile-landscape:text-xs font-medium text-fg-secondary",
                        "border border-stroke/30",
                      )}>
                        {room.player_count}/4
                      </span>
                    </td>
                    <td className="px-4 py-2 mobile-landscape:px-2 mobile-landscape:py-1.5 text-center">
                      {isUserInRoom ? (
                        <Button
                          intent="neutral"
                          size="sm"
                          onClick={() => onJoinRoom(room.id)}
                          className="border-state-active/50 text-state-active hover:bg-state-active/10 mobile-landscape:text-xs mobile-landscape:px-2 mobile-landscape:py-1"
                        >
                          返回
                        </Button>
                      ) : canJoin ? (
                        <Button
                          intent="primary"
                          size="sm"
                          onClick={() => onJoinRoom(room.id)}
                          className={cn(
                            "bg-gradient-primary",
                            "hover:brightness-110",
                            "shadow-elevation-1",
                            "hover:shadow-jade",
                            "mobile-landscape:text-xs mobile-landscape:px-2 mobile-landscape:py-1",
                          )}
                        >
                          加入
                        </Button>
                      ) : (
                        <Button intent="neutral" size="sm" disabled className="mobile-landscape:text-xs mobile-landscape:px-2 mobile-landscape:py-1">
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
      </div>

      {/* 加载指示器 */}
      {isLoading && rooms.length > 0 && (
        <div className="absolute inset-0 flex items-center justify-center bg-surface-base/50 backdrop-blur-sm rounded-xl">
          <div className="inline-flex items-center gap-2 text-fg-secondary">
            <div className={cn(
              "w-4 h-4 rounded-full border-2",
              "border-state-active/30 border-t-state-active",
              "animate-spin",
            )} />
            <span className="text-sm">更新中...</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default RoomList;
