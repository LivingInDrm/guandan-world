import React from 'react';
import { Outlet } from 'react-router-dom';

interface LayoutProps {
  children?: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-blue-600 text-white p-4">
        <h1 className="text-xl font-bold">掼蛋在线对战</h1>
      </header>
      <main className="container mx-auto p-4">
        {children || <Outlet />}
      </main>
    </div>
  );
};

export default Layout;