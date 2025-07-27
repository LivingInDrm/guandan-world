import React from 'react';
import { RoomStatus, type RoomInfo } from '../../types';

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
        return 'bg-green-100 text-green-800';
      case RoomStatus.READY:
        return 'bg-yellow-100 text-yellow-800';
      case RoomStatus.PLAYING:
        return 'bg-blue-100 text-blue-800';
      case RoomStatus.CLOSED:
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
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

  return (
    <div className="bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow">
      <div className="p-6">
        {/* Header */}
        <div className="flex justify-between items-start mb-4">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <h3 className="font-medium text-gray-900">房间 #{room.id.slice(-6)}</h3>
              {isOwner && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                  房主
                </span>
              )}
            </div>
            <div className="flex items-center space-x-3">
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(room.status)}`}>
                {getStatusText(room.status)}
              </span>
              <span className="text-sm text-gray-600">
                {room.player_count}/4 人
              </span>
            </div>
          </div>
        </div>

        {/* Players */}
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-700 mb-2">玩家列表</h4>
          <div className="grid grid-cols-2 gap-2">
            {Array.from({ length: 4 }).map((_, index) => {
              const player = room.players[index];
              return (
                <div
                  key={index}
                  className={`p-2 rounded text-sm text-center ${
                    player
                      ? player.id === currentUserId
                        ? 'bg-blue-50 text-blue-700 border border-blue-200'
                        : 'bg-gray-50 text-gray-700'
                      : 'bg-gray-100 text-gray-400 border-2 border-dashed border-gray-300'
                  }`}
                >
                  {player ? (
                    <div className="flex items-center justify-center space-x-1">
                      <span>{player.username}</span>
                      {player.online ? (
                        <div className="w-2 h-2 bg-green-400 rounded-full"></div>
                      ) : (
                        <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
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

        {/* Action button */}
        <div className="flex justify-end">
          {isUserInRoom ? (
            <button
              disabled
              className="px-4 py-2 bg-gray-100 text-gray-500 rounded-lg text-sm font-medium cursor-not-allowed"
            >
              已在房间内
            </button>
          ) : canJoin ? (
            <button
              onClick={handleJoinClick}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
            >
              加入房间
            </button>
          ) : (
            <button
              disabled
              className="px-4 py-2 bg-gray-100 text-gray-500 rounded-lg text-sm font-medium cursor-not-allowed"
            >
              {room.status === RoomStatus.PLAYING ? '游戏中' : '房间已满'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default RoomCard;