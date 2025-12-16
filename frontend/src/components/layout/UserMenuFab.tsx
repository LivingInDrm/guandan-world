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
      <div className={cn("fixed top-3 right-3 z-50", className)}>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              className={cn(
                "flex items-center justify-center",
                "w-11 h-11",
                "rounded-full cursor-pointer",
                "bg-surface-base/80 backdrop-blur-md",
                "border border-stroke/50",
                "shadow-card-interactive",
                "hover:shadow-card-elevated",
                "hover:border-state-active/50",
                "active:scale-95",
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

      <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
};

export default UserMenuFab;
