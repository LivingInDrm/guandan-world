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
      <header className="bg-gradient-to-r from-emerald-600 to-emerald-700 text-white p-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold tracking-wide">
            <span className="bg-gradient-to-r from-yellow-200 via-amber-100 to-yellow-200 bg-clip-text text-transparent drop-shadow-sm">
              一起掼蛋吧
            </span>
          </h1>
          
          {isAuthenticated && user && (
            <div className="flex items-center gap-4">
              <span className="text-sm">
                欢迎，<span className="font-semibold">{user.nickname || user.username}</span>
              </span>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    className="p-1.5 hover:bg-white/20 rounded-lg transition-all"
                    title="菜单"
                  >
                    <img
                      src={avatarUrl || userProfileIcon}
                      alt="菜单"
                      className="w-6 h-6 rounded-full object-cover"
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
      </header>

      <ProfileDialog open={profileOpen} onOpenChange={setProfileOpen} />
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />
    </>
  );
};

export default Header;
