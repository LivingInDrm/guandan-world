import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { getAvatarUrl } from '../../utils/avatar';
import userProfileIcon from '../../assets/user_profile.png';
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

const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();
  const [profileOpen, setProfileOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const avatarUrl = user ? getAvatarUrl(user.avatar_key) : null;

  return (
    <>
      <header className="relative py-4 px-6 bg-gradient-to-b from-[hsl(var(--primitive-primary-700))] to-[hsl(var(--primitive-primary-500))] text-[hsl(var(--color-text-inverse))] shadow-[var(--elevation-2),var(--shadow-relief)] rounded-lg">
        <div className="absolute inset-x-0 bottom-0 h-0.5 bg-[hsl(var(--color-state-active))] shadow-[var(--shadow-glow-md)] rounded-b-lg" />

        <div className="container mx-auto flex items-center">
          <div className="flex-1" />
          
          <h1 className="text-2xl font-bold tracking-wide text-center">
            <span className="bg-gradient-to-r from-[hsl(var(--primitive-accent-300))] via-[hsl(var(--primitive-accent-500))] to-[hsl(var(--primitive-accent-300))] bg-clip-text text-transparent drop-shadow-sm">
              一起掼蛋 | GuanDan Together
            </span>
          </h1>
          
          <div className="flex-1 flex justify-end">
            {isAuthenticated && user && (
              <div className="flex items-center gap-4">
                <span className="text-sm">
                  欢迎，<span className="font-bold text-[hsl(var(--primitive-accent-300))]">{user.nickname || user.username}</span>
                </span>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button
                      className="rounded-full hover:ring-4 hover:ring-[hsl(var(--color-state-active)/.5)] hover:shadow-[var(--shadow-glow-md)] transition-all duration-normal cursor-pointer group"
                      title="菜单"
                    >
                      <Avatar
                        src={avatarUrl || userProfileIcon}
                        alt="菜单"
                        fallback={user.nickname || user.username}
                        size="sm"
                        ringState="normal"
                        className="group-hover:scale-110 transition-transform duration-fast"
                      />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onSelect={() => setProfileOpen(true)}>
                      个人中心
                    </DropdownMenuItem>
                    <DropdownMenuItem onSelect={() => setSettingsOpen(true)}>
                      设置
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onSelect={handleLogout}>
                      退出登录
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            )}
          </div>
        </div>
      </header>

      <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
};

export default Header;
