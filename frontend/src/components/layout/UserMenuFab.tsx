import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { User, Settings, LogOut } from 'lucide-react';
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

interface UserMenuFabProps {
  className?: string;
}

const UserMenuFab: React.FC<UserMenuFabProps> = ({ className }) => {
  const { user, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();
  const [profileOpen, setProfileOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (!isAuthenticated || !user) {
    return null;
  }

  const avatarUrl = user ? getAvatarUrl(user.avatar_key) : null;

  return (
    <>
      <div className={cn("fixed top-4 right-4 mobile-landscape:top-2 mobile-landscape:right-3 z-50", className)}>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              className={cn(
                "relative flex items-center justify-center",
                "w-12 h-12 mobile-landscape:w-9 mobile-landscape:h-9",
                "rounded-xl",
                "cursor-pointer",
                "group",
                // 背景 - 翡翠绿渐变
                "bg-gradient-to-br from-[hsl(158,55%,38%)] to-[hsl(158,55%,28%)]",
                // 边框光泽
                "border border-white/20",
                // 阴影 - 翡翠绿光晕
                "shadow-[0_4px_16px_rgba(45,106,79,0.3),0_2px_4px_rgba(0,0,0,0.1)]",
                "hover:shadow-[0_8px_24px_rgba(45,106,79,0.4),0_4px_8px_rgba(0,0,0,0.15)]",
                // 过渡
                "transition-all duration-300 ease-out",
                "hover:scale-105",
                "active:scale-95",
              )}
              title="菜单"
            >
              {/* 光泽效果 */}
              <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-white/25 via-transparent to-transparent pointer-events-none" />

              {/* 头像 */}
              <Avatar
                src={avatarUrl || userProfileIcon}
                alt="菜单"
                fallback={user.nickname || user.username}
                size="sm"
                className={cn(
                  "relative z-10",
                  "w-9 h-9 mobile-landscape:w-7 mobile-landscape:h-7",
                  "ring-2 ring-white/30",
                  "group-hover:ring-white/50",
                  "transition-all duration-300",
                )}
              />
            </button>
          </DropdownMenuTrigger>

          <DropdownMenuContent align="end" sideOffset={8} className="min-w-[180px]">
            {/* 用户信息头部 */}
            <div className="px-3 py-3 mb-1.5 border-b border-stroke/20">
              <div className="flex items-center gap-3">
                {/* 小头像 */}
                <div className={cn(
                  "w-10 h-10 rounded-lg overflow-hidden",
                  "ring-2 ring-[hsl(158,55%,42%)]/20",
                  "shadow-sm",
                )}>
                  <Avatar
                    src={avatarUrl || userProfileIcon}
                    alt={user.nickname || user.username}
                    fallback={user.nickname || user.username}
                    size="sm"
                    className="w-full h-full"
                  />
                </div>

                {/* 用户名信息 */}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-display font-bold truncate text-fg-primary">
                    {user.nickname || user.username}
                  </p>
                  <p className="text-xs truncate text-fg-secondary">
                    @{user.username}
                  </p>
                </div>
              </div>
            </div>

            <DropdownMenuItem onSelect={() => setProfileOpen(true)}>
              <User className="w-4 h-4 mr-3 text-[hsl(158,55%,42%)]" />
              <span className="font-medium">个人中心</span>
            </DropdownMenuItem>

            <DropdownMenuItem onSelect={() => setSettingsOpen(true)}>
              <Settings className="w-4 h-4 mr-3 text-[hsl(158,55%,42%)]" />
              <span className="font-medium">设置</span>
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem onSelect={handleLogout} destructive>
              <LogOut className="w-4 h-4 mr-3" />
              <span className="font-medium">退出登录</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
};

export default UserMenuFab;
