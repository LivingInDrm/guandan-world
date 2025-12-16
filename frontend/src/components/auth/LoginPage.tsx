import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Home } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { apiClient } from '../../services/api';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import { Button } from '../ui';
import { cn } from '@/lib/utils';

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
      <div className="min-h-screen flex items-center justify-center">
        <div className={cn(
          "text-center p-8 rounded-xl",
          "bg-surface-base/80 backdrop-blur-sm",
          "shadow-card-interactive",
          "border border-stroke/30",
        )}>
          <div className={cn(
            "w-16 h-16 mx-auto mb-4",
            "border-4 border-state-active/30 border-t-state-active",
            "rounded-full animate-spin",
          )} />
          <p className="text-fg-primary font-display text-lg">正在加入房间...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      {/* 背景装饰 */}
      <div className="fixed inset-0 -z-10">
        <div className="absolute inset-0 bg-gradient-to-br from-surface-muted via-surface-elevated to-surface-base" />
        <div
          className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[800px] h-[800px]"
          style={{
            background: 'radial-gradient(circle, hsla(158,55%,42%,0.08) 0%, transparent 60%)',
          }}
        />
        <div
          className="absolute bottom-0 right-0 w-[600px] h-[600px]"
          style={{
            background: 'radial-gradient(circle, hsla(42,95%,52%,0.06) 0%, transparent 60%)',
          }}
        />
      </div>

      <div className="w-full max-w-md space-y-6">
        {/* 标题 */}
        <div className="text-center mb-8">
          <h1 className="font-display text-4xl font-bold tracking-wider mb-2">
            <span className={cn(
              "bg-gradient-primary",
              "bg-clip-text text-transparent",
            )}>
              一起掼蛋
            </span>
          </h1>
          <p className="text-fg-secondary text-sm tracking-wide">
            GuanDan Together
          </p>
        </div>

        {/* 错误提示 */}
        {joinError && (
          <div className={cn(
            "p-4 rounded-lg text-center space-y-3",
            "bg-action-danger/10 border border-action-danger/30",
          )}>
            <p className="text-sm text-action-danger">{joinError}</p>
            <Button intent="secondary" size="sm" onClick={handleGoToLobby}>
              <Home className="w-4 h-4" />
              返回大厅
            </Button>
          </div>
        )}

        {/* 邀请码提示 */}
        {inviteCode && !joinError && (
          <div className={cn(
            "p-4 rounded-lg text-center",
            "bg-action-primary/10 border border-action-primary/30",
          )}>
            <p className="text-sm text-fg-primary">
              请登录后加入房间{' '}
              <span className="font-mono font-bold text-state-active">{inviteCode}</span>
            </p>
          </div>
        )}

        {/* 登录/注册表单 */}
        {isLogin ? <LoginForm /> : <RegisterForm />}

        {/* 切换按钮 */}
        <div className="text-center pt-2">
          <button
            type="button"
            onClick={() => setIsLogin(!isLogin)}
            className={cn(
              "text-sm font-medium",
              "text-action-primary hover:text-action-primary/80",
              "underline-offset-4 hover:underline",
              "transition-colors duration-fast",
            )}
          >
            {isLogin ? '没有账号？立即注册' : '已有账号？立即登录'}
          </button>
        </div>

        {/* 底部装饰 */}
        <div className="flex items-center justify-center gap-2 pt-6 opacity-40">
          <span className="text-2xl">♠</span>
          <span className="text-2xl text-suit-red">♥</span>
          <span className="text-2xl">♣</span>
          <span className="text-2xl text-suit-red">♦</span>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
