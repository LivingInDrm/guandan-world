import React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';

interface LayoutProps {
  children?: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-100">
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
                欢迎，<span className="font-semibold">{user.username}</span>
              </span>
              <button
                onClick={handleLogout}
                className="bg-white/20 hover:bg-white/30 backdrop-blur-sm border border-white/30 px-4 py-1.5 rounded-lg text-sm transition-all"
              >
                退出登录
              </button>
            </div>
          )}
        </div>
      </header>
      <main className="container mx-auto p-4">
        {children || <Outlet />}
      </main>
    </div>
  );
};

export default Layout;