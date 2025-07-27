import React, { useEffect, useState } from 'react';
import { useRoomStore } from '../../store/roomStore';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import RoomList from './RoomList';
import CreateRoomModal from './CreateRoomModal';
import type { RoomInfo } from '../../types';

const RoomLobby: React.FC = () => {
  const { user } = useAuthStore();
  const {
    roomList,
    totalCount,
    currentPage,
    limit,
    isLoading,
    error,
    setRoomList,
    setLoading,
    setError,
    clearError,
    setPage
  } = useRoomStore();

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null);

  // Load room list
  const loadRoomList = async (page: number = currentPage) => {
    if (!user) return;

    setLoading(true);
    clearError();

    try {
      const response = await apiClient.getRoomList(page, limit);
      if (response.success && response.data) {
        setRoomList(response.data);
      } else {
        setError(response.error || '获取房间列表失败');
      }
    } catch (err: any) {
      console.error('Load room list error:', err);
      setError(err.message || '获取房间列表失败');
    } finally {
      setLoading(false);
    }
  };

  // Handle page change
  const handlePageChange = (page: number) => {
    setPage(page);
    loadRoomList(page);
  };

  // Handle room creation
  const handleCreateRoom = async () => {
    if (!user) return;

    try {
      const response = await apiClient.createRoom();
      if (response.success && response.data) {
        setShowCreateModal(false);
        // Refresh room list to show the new room
        await loadRoomList(1);
        setPage(1);
      } else {
        setError(response.error || '创建房间失败');
      }
    } catch (err: any) {
      console.error('Create room error:', err);
      setError(err.message || '创建房间失败');
    }
  };

  // Handle room join
  const handleJoinRoom = async (roomId: string) => {
    if (!user) return;

    try {
      const response = await apiClient.joinRoom(roomId);
      if (response.success) {
        // Room join successful, the user will be redirected by the backend
        // or we can navigate to the room waiting page
        console.log('Successfully joined room:', roomId);
        // TODO: Navigate to room waiting page
      } else {
        setError(response.error || '加入房间失败');
      }
    } catch (err: any) {
      console.error('Join room error:', err);
      setError(err.message || '加入房间失败');
    }
  };

  // Auto-refresh room list every 5 seconds
  useEffect(() => {
    if (user) {
      loadRoomList();
      
      const interval = setInterval(() => {
        loadRoomList();
      }, 5000);
      
      setRefreshInterval(interval);

      return () => {
        if (interval) {
          clearInterval(interval);
        }
      };
    }
  }, [user]);

  // Cleanup interval on unmount
  useEffect(() => {
    return () => {
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
    };
  }, [refreshInterval]);

  if (!user) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-gray-600">请先登录</p>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">房间大厅</h1>
          <p className="text-gray-600 mt-1">
            欢迎，{user.username}！选择一个房间开始游戏
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          创建房间
        </button>
      </div>

      {/* Error display */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <div className="flex justify-between items-center">
            <p className="text-red-600">{error}</p>
            <button
              onClick={clearError}
              className="text-red-400 hover:text-red-600"
            >
              ✕
            </button>
          </div>
        </div>
      )}

      {/* Room statistics */}
      <div className="bg-gray-50 rounded-lg p-4 mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600">{totalCount}</div>
              <div className="text-sm text-gray-600">总房间数</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">
                {roomList.filter(room => room.status === 0).length}
              </div>
              <div className="text-sm text-gray-600">等待中</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-yellow-600">
                {roomList.filter(room => room.status === 2).length}
              </div>
              <div className="text-sm text-gray-600">游戏中</div>
            </div>
          </div>
          <button
            onClick={() => loadRoomList()}
            disabled={isLoading}
            className="text-blue-600 hover:text-blue-800 disabled:text-gray-400"
          >
            {isLoading ? '刷新中...' : '手动刷新'}
          </button>
        </div>
      </div>

      {/* Room List */}
      <RoomList
        rooms={roomList}
        isLoading={isLoading}
        currentPage={currentPage}
        totalCount={totalCount}
        limit={limit}
        onPageChange={handlePageChange}
        onJoinRoom={handleJoinRoom}
        currentUserId={user.id}
      />

      {/* Create Room Modal */}
      {showCreateModal && (
        <CreateRoomModal
          onClose={() => setShowCreateModal(false)}
          onConfirm={handleCreateRoom}
        />
      )}
    </div>
  );
};

export default RoomLobby;