import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import { cn } from '@/lib/utils';

interface LayoutProps {
  children?: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="relative min-h-screen overflow-hidden">
      {/* 背景层 - 东方雅韵氛围 */}
      <div className="fixed inset-0 -z-10">
        {/* 主渐变背景 */}
        <div className="absolute inset-0 bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]" />

        {/* 装饰性光晕 */}
        <div
          className="absolute top-0 left-1/4 w-[600px] h-[600px] opacity-40"
          style={{
            background: 'radial-gradient(circle, hsla(158,55%,42%,0.12) 0%, transparent 70%)',
          }}
        />
        <div
          className="absolute bottom-0 right-1/4 w-[800px] h-[800px] opacity-30"
          style={{
            background: 'radial-gradient(circle, hsla(42,95%,52%,0.1) 0%, transparent 70%)',
          }}
        />

        {/* 微妙的纹理叠加 */}
        <div
          className="absolute inset-0 opacity-[0.02]"
          style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23000000' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          }}
        />
      </div>

      {/* 顶部装饰边 */}
      <div className="fixed top-0 inset-x-0 h-1 bg-gradient-to-r from-transparent via-state-active/60 to-transparent z-50" />

      {/* Header */}
      <Header />

      {/* 主内容区 */}
      <main className={cn(
        "relative container mx-auto px-4 py-6",
        "min-h-[calc(100vh-80px)]",
      )}>
        {children || <Outlet />}
      </main>

      {/* 底部装饰 */}
      <div className="fixed bottom-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-stroke/30 to-transparent" />
    </div>
  );
};

export default Layout;
