import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { getAvatarUrl } from '../../utils/avatar';
import userProfileIcon from '../../assets/user_profile.png';
import { cn } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/DropdownMenu';
import { Avatar } from '../ui/Avatar';
import ProfileDialog from '../profile/ProfileDialog';
import SettingsDialog from '../settings/SettingsDialog';

interface HeaderProps {
  collapsible?: boolean;
}

const Header: React.FC<HeaderProps> = ({ collapsible = false }) => {
  const { user, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();
  const [profileOpen, setProfileOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [isHovered, setIsHovered] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const avatarUrl = user ? getAvatarUrl(user.avatar_key) : null;

  // 可收起模式：默认收起，hover 展开
  if (collapsible) {
    return (
      <>
        <div
          className="fixed top-0 left-0 right-0 z-50"
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
        >
          {/* 触发区域 - 始终可见 */}
          <div className={cn(
            "h-2 bg-gradient-to-b from-black/20 to-transparent",
            "cursor-pointer",
          )} />

          {/* Header 主体 - 可收起 */}
          <header
            className={cn(
              "relative mx-4 mobile-landscape:mx-2",
              "py-3 px-6 mobile-landscape:py-2 mobile-landscape:px-3",
              "bg-gradient-header-bar",
              "rounded-b-xl mobile-landscape:rounded-b-lg",
              "border border-t-0 border-stroke-header",
              "shadow-card-elevated",
              "text-fg-inverse",
              "transition-all duration-300 ease-out",
              isHovered
                ? "translate-y-0 opacity-100"
                : "-translate-y-full opacity-0 pointer-events-none",
            )}
          >
            {/* 底部金色装饰线 */}
            <div className={cn(
              "absolute inset-x-4 mobile-landscape:inset-x-2 -bottom-px h-[2px] mobile-landscape:h-px rounded-full",
              "bg-gradient-to-r from-transparent via-state-active to-transparent",
              "shadow-[0_0_12px_hsla(42,95%,52%,0.6)]",
            )} />

            <div className="container mx-auto flex items-center">
              <div className="flex-1" />
              <div className="text-center">
                <h1 className="text-xl mobile-landscape:text-lg font-display font-bold tracking-wider">
                  <span className={cn(
                    "bg-gradient-gold",
                    "bg-clip-text text-transparent",
                    "drop-shadow-[0_2px_4px_rgba(0,0,0,0.3)]",
                  )}>
                    一起掼蛋
                  </span>
                </h1>
              </div>
              <div className="flex-1 flex justify-end">
                {isAuthenticated && user && (
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <button
                        className={cn(
                          "rounded-full cursor-pointer",
                          "ring-2 ring-white/20",
                          "hover:ring-state-active/60",
                          "transition-all duration-normal",
                        )}
                      >
                        <Avatar
                          src={avatarUrl || userProfileIcon}
                          alt="菜单"
                          fallback={user.nickname || user.username}
                          size="sm"
                        />
                      </button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="min-w-[140px]">
                      <DropdownMenuItem onSelect={() => setProfileOpen(true)}>
                        个人中心
                      </DropdownMenuItem>
                      <DropdownMenuItem onSelect={() => setSettingsOpen(true)}>
                        设置
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onSelect={handleLogout} destructive>
                        退出登录
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                )}
              </div>
            </div>
          </header>
        </div>

        <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
        <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
      </>
    );
  }

  // 常规模式
  return (
    <>
      <header
        className={cn(
          "relative mx-4 mt-4 mb-2 mobile-landscape:mx-2 mobile-landscape:mt-2 mobile-landscape:mb-1",
          "py-4 px-6 mobile-landscape:py-2 mobile-landscape:px-3",
          "bg-gradient-header-bar",
          "rounded-xl mobile-landscape:rounded-lg",
          "border border-stroke-header",
          "shadow-card-elevated",
          "text-fg-inverse",
        )}
        style={{ minHeight: 'var(--header-height)' }}
      >
        {/* 顶部高光线 */}
        <div className="absolute inset-x-0 top-0 h-px rounded-t-xl bg-gradient-to-r from-transparent via-white/30 to-transparent" />

        {/* 底部金色装饰线 */}
        <div className={cn(
          "absolute inset-x-4 mobile-landscape:inset-x-2 -bottom-px h-[2px] mobile-landscape:h-px rounded-full",
          "bg-gradient-to-r from-transparent via-state-active to-transparent",
          "shadow-[0_0_12px_hsla(42,95%,52%,0.6)]",
        )} />

        <div className="container mx-auto flex items-center">
          {/* 左侧空白区域 - 保持标题居中 */}
          <div className="flex-1" />

          {/* 中央标题 */}
          <div className="text-center">
            <h1 className="text-2xl mobile-landscape:text-lg font-display font-bold tracking-wider">
              <span className={cn(
                "bg-gradient-gold",
                "bg-clip-text text-transparent",
                "drop-shadow-[0_2px_4px_rgba(0,0,0,0.3)]",
              )}>
                一起掼蛋
              </span>
              <span className="mx-3 text-white/30 mobile-landscape:hidden">|</span>
              <span className="text-white/90 text-lg font-normal tracking-wide mobile-landscape:hidden">
                GuanDan Together
              </span>
            </h1>
          </div>

          {/* 右侧用户区域 */}
          <div className="flex-1 flex justify-end">
            {isAuthenticated && user && (
              <div className="flex items-center gap-4 mobile-landscape:gap-2">
                <span className="text-sm mobile-landscape:text-xs text-white/80 mobile-landscape:hidden">
                  欢迎，
                  <span className={cn(
                    "font-bold ml-1",
                    "bg-gradient-gold-text",
                    "bg-clip-text text-transparent",
                  )}>
                    {user.nickname || user.username}
                  </span>
                </span>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button
                      className={cn(
                        "rounded-full cursor-pointer",
                        "ring-2 ring-white/20 ring-offset-2 ring-offset-transparent",
                        "hover:ring-state-active/60",
                        "hover:shadow-[var(--shadow-glow-md)]",
                        "transition-all duration-normal",
                        "group",
                      )}
                      title="菜单"
                    >
                      <Avatar
                        src={avatarUrl || userProfileIcon}
                        alt="菜单"
                        fallback={user.nickname || user.username}
                        size="sm"
                        className="group-hover:scale-105 transition-transform duration-fast"
                      />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="min-w-[140px]">
                    <DropdownMenuItem onSelect={() => setProfileOpen(true)}>
                      个人中心
                    </DropdownMenuItem>
                    <DropdownMenuItem onSelect={() => setSettingsOpen(true)}>
                      设置
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onSelect={handleLogout} destructive>
                      退出登录
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            )}
          </div>
        </div>

        {/* 装饰性角标 */}
        <div className="absolute top-2 left-3 mobile-landscape:top-1 mobile-landscape:left-2 w-6 h-6 mobile-landscape:w-4 mobile-landscape:h-4 opacity-20">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" />
          </svg>
        </div>
        <div className="absolute top-2 right-3 mobile-landscape:top-1 mobile-landscape:right-2 w-6 h-6 mobile-landscape:w-4 mobile-landscape:h-4 opacity-20">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" />
          </svg>
        </div>
      </header>

      <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
};

export default Header;
