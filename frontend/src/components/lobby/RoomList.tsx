import React from 'react';
import { RoomStatus, type RoomInfo } from '../../types';
import RoomCard from './RoomCard';
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
  // Sort rooms: waiting rooms first, then by player count descending
  const sortedRooms = [...rooms].sort((a, b) => {
    // Priority order: waiting > ready > playing > closed
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
    
    // Same status, sort by player count descending
    return b.player_count - a.player_count;
  });

  if (isLoading && rooms.length === 0) {
    return (
      <div className="space-y-4">
        {/* Loading skeleton */}
        {Array.from({ length: 6 }).map((_, index) => (
          <div key={index} className="bg-white rounded-lg shadow-sm border p-6 animate-pulse">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="h-4 bg-gray-200 rounded w-24 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-32 mb-4"></div>
                <div className="flex space-x-4">
                  <div className="h-3 bg-gray-200 rounded w-16"></div>
                  <div className="h-3 bg-gray-200 rounded w-20"></div>
                </div>
              </div>
              <div className="h-8 bg-gray-200 rounded w-20"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (rooms.length === 0 && !isLoading) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-400 text-6xl mb-4">🎮</div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">暂无房间</h3>
        <p className="text-gray-600 mb-6">成为第一个创建房间的玩家吧！</p>
      </div>
    );
  }

  const totalPages = Math.ceil(totalCount / limit);

  return (
    <div className="space-y-6">
      {/* Room grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {sortedRooms.map((room) => (
          <RoomCard
            key={room.id}
            room={room}
            onJoinRoom={onJoinRoom}
            currentUserId={currentUserId}
          />
        ))}
      </div>

      {/* Loading overlay for refresh */}
      {isLoading && rooms.length > 0 && (
        <div className="text-center py-4">
          <div className="inline-flex items-center text-gray-600">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
            更新中...
          </div>
        </div>
      )}

      {/* Pagination */}
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