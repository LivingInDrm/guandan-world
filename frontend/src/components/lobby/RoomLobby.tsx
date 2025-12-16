import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Zap, Plus, Users, X, AlertCircle } from 'lucide-react';
import { useRoomStore } from '../../store/roomStore';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import RoomList from './RoomList';
import CreateRoomModal from './CreateRoomModal';
import JoinRoomModal from './JoinRoomModal';
import { Button, Card, CardContent } from '../ui';
import { cn } from '@/lib/utils';
import OrientationGuard from '../ui/OrientationGuard';
import Header from '../layout/Header';

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
  const [showJoinModal, setShowJoinModal] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null);
  const [isCheckingRoom, setIsCheckingRoom] = useState(false);
  const [isQuickStarting, setIsQuickStarting] = useState(false);

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

  // Handle room join by code
  const handleJoinByCode = async (roomCode: string) => {
    if (!user) return;

    const response = await apiClient.joinRoomByCode(roomCode);
    if (response.success && response.data) {
      setCurrentRoom(response.data);
      navigate(`/game/${response.data.id}`);
    } else {
      throw new Error(response.error || '加入房间失败');
    }
  };

  // Handle quick start
  const handleQuickStart = async () => {
    if (!user || isQuickStarting) return;

    setIsQuickStarting(true);
    try {
      const response = await apiClient.quickJoin();
      if (response.success && response.data) {
        setCurrentRoom(response.data);
        navigate(`/game/${response.data.id}`);
      } else {
        setError(response.error || '快速开始失败');
      }
    } catch (err: any) {
      console.error('Quick start error:', err);
      setError(err.message || '快速开始失败');
    } finally {
      setIsQuickStarting(false);
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
      <OrientationGuard>
        <div className="fixed inset-0 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
          <div className="flex items-center justify-center h-full">
            <p className="text-fg-secondary">请先登录</p>
          </div>
        </div>
      </OrientationGuard>
    );
  }

  return (
    <OrientationGuard>
      <div className="fixed inset-0 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
        {/* 背景装饰层 */}
        <div className="absolute inset-0 -z-10 pointer-events-none">
          <div
            className="absolute top-0 left-1/4 w-[600px] h-[600px] opacity-40"
            style={{
              background: 'radial-gradient(circle, hsla(158,55%,42%,0.12) 0%, transparent 70%)',
            }}
          />
          <div
            className="absolute bottom-0 right-1/4 w-[800px] h-[800px] opacity-30"
            style={{
              background: 'radial-gradient(circle, hsla(42,95%,52%,0.1) 0%, transparent 70%)',
            }}
          />
          <div
            className="absolute inset-0 opacity-[0.02]"
            style={{
              backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23000000' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
            }}
          />
        </div>

        {/* 顶部装饰边 */}
        <div className="absolute top-0 inset-x-0 h-1 bg-gradient-to-r from-transparent via-state-active/60 to-transparent z-50" />

        {/* Header */}
        <Header />

        {/* 主内容区 */}
        <main className={cn(
          "absolute inset-0 pt-14 mobile-landscape:pt-10",
          "container mx-auto px-4 py-6",
          "mobile-landscape:px-1 mobile-landscape:py-1 mobile-landscape:max-w-none",
          "overflow-auto",
        )}>
          <div className={cn(
            "max-w-6xl mx-auto mobile-landscape:max-w-none",
            "p-6 mobile-landscape:p-2 rounded-2xl mobile-landscape:rounded-lg",
            "bg-surface-base/60 backdrop-blur-sm",
            "border border-stroke/30 mobile-landscape:border-0",
            "shadow-[var(--elevation-2)] mobile-landscape:shadow-none",
            "mobile-landscape:bg-transparent",
          )}>
            {/* 错误提示 */}
            {error && (
              <div className={cn(
                "mb-6 p-4 rounded-xl",
                "bg-action-danger/10 border border-action-danger/30",
                "animate-in fade-in-0 slide-in-from-top-2 duration-normal",
              )}>
                <div className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-2 text-action-danger">
                    <AlertCircle className="w-5 h-5 shrink-0" />
                    <p className="text-sm">{error}</p>
                  </div>
                  <button
                    onClick={clearError}
                    className={cn(
                      "p-1 rounded-md",
                      "text-action-danger/70 hover:text-action-danger",
                      "hover:bg-action-danger/10",
                      "transition-colors duration-fast",
                    )}
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}

            <div className="flex flex-col mobile-landscape:flex-row lg:flex-row gap-4 mobile-landscape:gap-3 lg:gap-8">
              {/* 左侧操作面板 */}
              <div className="w-full mobile-landscape:w-48 lg:w-72 mobile-landscape:shrink-0 shrink-0">
                <Card
                  variant="elevated"
                  interactive={false}
                  className={cn(
                    "overflow-hidden",
                    "bg-gradient-to-b from-surface-base to-surface-elevated/80",
                    "border border-stroke/50",
                    "shadow-card-interactive",
                  )}
                >
                  {/* 装饰顶部 */}
                  <div className="h-1 bg-gradient-to-r from-transparent via-state-active to-transparent" />

                  <CardContent className="p-5 mobile-landscape:p-3">
                    {/* 标题 */}
                    <h2 className={cn(
                      "text-lg mobile-landscape:text-base font-display font-bold text-center mb-5 mobile-landscape:mb-3",
                      "bg-gradient-primary",
                      "bg-clip-text text-transparent",
                    )}>
                      开始游戏
                    </h2>

                    <div className="flex flex-col gap-3 mobile-landscape:gap-2">
                      {/* 快速开始按钮 */}
                      <Button
                        intent="primary"
                        size="lg"
                        className={cn(
                          "w-full",
                          "mobile-landscape:h-9 mobile-landscape:text-sm",
                          "bg-gradient-active",
                          "hover:brightness-110",
                          "shadow-card-interactive",
                          "hover:shadow-card-elevated",
                          "text-primitive-neutral-900 font-bold",
                        )}
                        onClick={handleQuickStart}
                        disabled={isQuickStarting}
                      >
                        <Zap className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4" />
                        快速开始
                      </Button>

                      <div className="flex items-center gap-3 my-1 mobile-landscape:my-0.5">
                        <div className="flex-1 h-px bg-gradient-to-r from-transparent via-stroke to-transparent" />
                        <span className="text-xs text-fg-secondary">或</span>
                        <div className="flex-1 h-px bg-gradient-to-r from-transparent via-stroke to-transparent" />
                      </div>

                      {/* 创建房间按钮 */}
                      <Button
                        intent="secondary"
                        size="lg"
                        className={cn(
                          "w-full",
                          "mobile-landscape:h-9 mobile-landscape:text-sm",
                          "bg-gradient-primary",
                          "hover:brightness-110",
                          "border-none text-white",
                          "shadow-card-interactive",
                          "hover:shadow-card-elevated",
                        )}
                        onClick={() => setShowCreateModal(true)}
                      >
                        <Plus className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4" />
                        创建房间
                      </Button>

                      {/* 加入游戏按钮 */}
                      <Button
                        intent="tertiary"
                        size="lg"
                        className={cn(
                          "w-full",
                          "mobile-landscape:h-9 mobile-landscape:text-sm",
                          "border-stroke/50 hover:border-state-active/50",
                          "hover:bg-state-active/5",
                        )}
                        onClick={() => setShowJoinModal(true)}
                      >
                        <Users className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4" />
                        输入房间码
                      </Button>
                    </div>

                    {/* 底部装饰 */}
                    <div className="flex items-center justify-center gap-2 mt-5 pt-4 border-t border-stroke/30 mobile-landscape:hidden">
                      <span className="text-lg opacity-30">♠</span>
                      <span className="text-lg text-suit-red opacity-30">♥</span>
                      <span className="text-lg opacity-30">♣</span>
                      <span className="text-lg text-suit-red opacity-30">♦</span>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* 右侧房间列表 */}
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

            <JoinRoomModal
              open={showJoinModal}
              onClose={() => setShowJoinModal(false)}
              onJoin={handleJoinByCode}
            />
          </div>
        </main>

        {/* 底部装饰 */}
        <div className="absolute bottom-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stroke/30 to-transparent" />
      </div>
    </OrientationGuard>
  );
};

export default RoomLobby;