import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Home } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { apiClient } from '../../services/api';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import { Button } from '../ui';

interface LocationState {
  from?: { pathname: string };
  inviteCode?: string;
}

const LoginPage: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [isJoiningRoom, setIsJoiningRoom] = useState(false);
  const [joinError, setJoinError] = useState<string | null>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated } = useAuthStore();
  const { setCurrentRoom } = useRoomStore();

  const state = location.state as LocationState | null;
  const inviteCode = state?.inviteCode;

  useEffect(() => {
    if (!isAuthenticated) return;

    const handlePostAuth = async () => {
      if (inviteCode) {
        setIsJoiningRoom(true);
        setJoinError(null);

        try {
          const response = await apiClient.joinRoomByCode(inviteCode);
          if (response.success && response.data) {
            setCurrentRoom(response.data);
            navigate(`/game/${response.data.id}`, { replace: true });
            return;
          }
          setJoinError(response.error || '加入房间失败');
        } catch (err: unknown) {
          const error = err as { message?: string };
          if (error.message?.includes('already in room')) {
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
          setJoinError(error.message || '加入房间失败');
        } finally {
          setIsJoiningRoom(false);
        }
        return;
      }

      const targetPath = state?.from?.pathname;
      if (targetPath && targetPath !== '/login') {
        navigate(targetPath, { replace: true });
      } else {
        navigate('/lobby', { replace: true });
      }
    };

    handlePostAuth();
  }, [isAuthenticated, inviteCode, state?.from, navigate, setCurrentRoom]);

  const handleGoToLobby = () => {
    navigate('/lobby', { replace: true });
  };

  if (isJoiningRoom) {
    return (
      <div className="min-h-screen bg-surface-base flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-action-primary mx-auto mb-4"></div>
          <p className="text-fg-secondary">正在加入房间...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-surface-base flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-6">
        {joinError && (
          <div className="bg-error/10 border border-error/30 rounded-sm p-4 text-center space-y-3">
            <p className="text-sm text-error">{joinError}</p>
            <Button intent="secondary" size="sm" onClick={handleGoToLobby}>
              <Home className="w-4 h-4" />
              返回大厅
            </Button>
          </div>
        )}
        {inviteCode && !joinError && (
          <div className="bg-action-primary/10 border border-action-primary/30 rounded-sm p-3 text-center">
            <p className="text-sm text-fg-primary">
              请登录后加入房间 <span className="font-mono font-bold">{inviteCode}</span>
            </p>
          </div>
        )}
        {isLogin ? <LoginForm /> : <RegisterForm />}
        
        <div className="text-center">
          <button
            type="button"
            onClick={() => setIsLogin(!isLogin)}
            className="text-action-primary hover:text-action-primary/80 text-sm font-medium transition-colors"
          >
            {isLogin ? '没有账号？立即注册' : '已有账号？立即登录'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
