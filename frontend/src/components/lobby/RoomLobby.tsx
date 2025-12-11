import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Zap, Plus, Users } from 'lucide-react';
import { useRoomStore } from '../../store/roomStore';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import RoomList from './RoomList';
import CreateRoomModal from './CreateRoomModal';
import { Button, Card } from '../ui';

const RoomLobby: React.FC = () => {
  const { user } = useAuthStore();
  const location = useLocation();
  const navigate = useNavigate();
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
    setPage,
    setCurrentRoom
  } = useRoomStore();

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null);
  const [isCheckingRoom, setIsCheckingRoom] = useState(true);

  // Check if user is already in a room
  const checkUserRoom = async () => {
    if (!user) {
      setIsCheckingRoom(false);
      return;
    }

    try {
      const response = await apiClient.getMyRoom();
      if (response.success && response.data) {
        // User is already in a room, auto-redirect
        // Set current room before navigation to avoid blank page
        setCurrentRoom(response.data);
        navigate(`/game/${response.data.id}`, { replace: true });
        return;
      }
    } catch (err: any) {
      // 404 means user is not in any room, which is normal
      if (err.status !== 404) {
        console.error('Failed to check user room:', err);
      }
    } finally {
      setIsCheckingRoom(false);
    }
  };

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
        // Set current room before navigation to avoid blank page
        setCurrentRoom(response.data);
        // Navigate to the game page
        navigate(`/game/${response.data.id}`);
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
      if (response.success && response.data) {
        // Set current room before navigation to avoid blank page
        setCurrentRoom(response.data);
        // Navigate to the game page
        navigate(`/game/${roomId}`);
      } else {
        setError(response.error || '加入房间失败');
      }
    } catch (err: any) {
      console.error('Join room error:', err);
      setError(err.message || '加入房间失败');
    }
  };

  // Check user's room status on mount
  useEffect(() => {
    checkUserRoom();
  }, [user]);

  // Auto-refresh room list every 5 seconds
  useEffect(() => {
    if (user && !isCheckingRoom) {
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
  }, [user, isCheckingRoom]);

  // Refresh room list when returning from room page with shouldRefresh flag
  useEffect(() => {
    const state = location.state as { shouldRefresh?: boolean } | null;
    if (state?.shouldRefresh) {
      loadRoomList();
      // Clear the state to prevent refresh on subsequent renders
      window.history.replaceState({}, document.title);
    }
  }, [location]);

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
        <p className="text-fg-secondary">请先登录</p>
      </div>
    );
  }

  if (isCheckingRoom) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-action-primary mx-auto mb-4"></div>
          <p className="text-fg-secondary">检查房间状态...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto p-6 bg-surface-elevated/20 rounded-lg">
      {error && (
        <div className="bg-error/10 border border-error/30 rounded-sm p-4 mb-6">
          <div className="flex justify-between items-center">
            <p className="text-error">{error}</p>
            <button
              onClick={clearError}
              className="text-error/70 hover:text-error"
            >
              ✕
            </button>
          </div>
        </div>
      )}

      <div className="flex flex-col lg:flex-row gap-6 lg:gap-8">
        <div className="w-full lg:w-64 shrink-0">
          <Card variant="emphasis" className="p-4" interactive={false}>
            <div className="flex flex-col gap-3">
              <Button intent="primary" size="lg" className="w-full" onClick={() => {}}>
                <Zap className="w-5 h-5" />
                快速开始
              </Button>
              <Button intent="secondary" size="lg" className="w-full" onClick={() => setShowCreateModal(true)}>
                <Plus className="w-5 h-5" />
                创建房间
              </Button>
              <Button intent="neutral" size="lg" className="w-full" onClick={() => {}}>
                <Users className="w-5 h-5" />
                加入游戏
              </Button>
            </div>
          </Card>
        </div>

        <div className="flex-1 min-w-0">
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
        </div>
      </div>

      <CreateRoomModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onConfirm={handleCreateRoom}
      />
    </div>
  );
};

export default RoomLobby;