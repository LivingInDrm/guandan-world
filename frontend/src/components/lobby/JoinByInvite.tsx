import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Loader2, AlertCircle, Home } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { apiClient } from '../../services/api';
import { Button, Card } from '../ui';

const JoinByInvite: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { isAuthenticated, isInitialized } = useAuthStore();
  const { setCurrentRoom } = useRoomStore();

  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const roomCode = searchParams.get('code');

  useEffect(() => {
    if (!isInitialized) return;

    if (!roomCode || roomCode.length !== 4) {
      setError('无效的邀请链接');
      return;
    }

    if (!isAuthenticated) {
      navigate('/login', { 
        state: { inviteCode: roomCode },
        replace: true 
      });
      return;
    }

    const joinRoom = async () => {
      setIsJoining(true);
      setError(null);

      try {
        const response = await apiClient.joinRoomByCode(roomCode);
        if (response.success && response.data) {
          setCurrentRoom(response.data);
          navigate(`/game/${response.data.id}`, { replace: true });
        } else {
          setError(response.error || '加入房间失败');
        }
      } catch (err: unknown) {
        const error = err as { response?: { message?: string }; message?: string };
        const errorMessage = error.response?.message || error.message || '加入房间失败';
        if (errorMessage.includes('already in room')) {
          try {
            const myRoomResponse = await apiClient.getMyRoom();
            if (myRoomResponse.success && myRoomResponse.data) {
              setCurrentRoom(myRoomResponse.data);
              navigate(`/game/${myRoomResponse.data.id}`, { replace: true });
              return;
            }
          } catch {
            // ignore
          }
        }
        setError(errorMessage);
      } finally {
        setIsJoining(false);
      }
    };

    joinRoom();
  }, [isInitialized, isAuthenticated, roomCode, navigate, setCurrentRoom]);

  if (!isInitialized || isJoining) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Card variant="elevated" interactive={false} className="p-8 text-center">
          <Loader2 className="w-12 h-12 animate-spin mx-auto mb-4 text-action-primary" />
          <p className="text-fg-secondary">
            {!isInitialized ? '加载中...' : '正在加入房间...'}
          </p>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Card variant="elevated" interactive={false} className="p-8 text-center max-w-md">
          <AlertCircle className="w-12 h-12 mx-auto mb-4 text-error" />
          <h2 className="text-xl font-semibold text-fg-primary mb-2">加入失败</h2>
          <p className="text-fg-secondary mb-6">{error}</p>
          <Button intent="primary" onClick={() => navigate('/lobby')}>
            <Home className="w-4 h-4" />
            返回大厅
          </Button>
        </Card>
      </div>
    );
  }

  return null;
};

export default JoinByInvite;
