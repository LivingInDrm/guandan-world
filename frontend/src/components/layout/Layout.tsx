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
      <header className="bg-blue-600 text-white p-4">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-xl font-bold">掼蛋在线对战</h1>
          
          {isAuthenticated && user && (
            <div className="flex items-center gap-4">
              <span className="text-sm">
                欢迎，<span className="font-semibold">{user.username}</span>
              </span>
              <button
                onClick={handleLogout}
                className="bg-blue-700 hover:bg-blue-800 px-4 py-1.5 rounded text-sm transition-colors"
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