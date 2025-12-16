import React, { useState, useEffect } from 'react';
import { cn } from '@/lib/utils';
import { RotateCcw, Smartphone } from 'lucide-react';

interface OrientationGuardProps {
  children: React.ReactNode;
  /** 最小宽度阈值，低于此宽度时才显示提示（默认1024px） */
  minWidthThreshold?: number;
}

/**
 * OrientationGuard - 竖屏旋转提示组件
 *
 * 当移动端/平板设备处于竖屏模式时，显示旋转提示。
 * 桌面端（宽度 >= minWidthThreshold）不显示提示。
 */
const OrientationGuard: React.FC<OrientationGuardProps> = ({
  children,
  minWidthThreshold = 1024,
}) => {
  const [isPortrait, setIsPortrait] = useState(false);
  const [isMobileDevice, setIsMobileDevice] = useState(false);

  useEffect(() => {
    const checkOrientation = () => {
      const isNarrowScreen = window.innerWidth < minWidthThreshold;
      const isPortraitOrientation = window.innerHeight > window.innerWidth;

      setIsMobileDevice(isNarrowScreen);
      setIsPortrait(isNarrowScreen && isPortraitOrientation);
    };

    // 初始检查
    checkOrientation();

    // 监听窗口变化
    window.addEventListener('resize', checkOrientation);
    window.addEventListener('orientationchange', checkOrientation);

    return () => {
      window.removeEventListener('resize', checkOrientation);
      window.removeEventListener('orientationchange', checkOrientation);
    };
  }, [minWidthThreshold]);

  // 桌面端或横屏时直接显示内容
  if (!isPortrait) {
    return <>{children}</>;
  }

  // 竖屏时显示旋转提示
  return (
    <>
      {/* 背景内容（模糊） */}
      <div className="blur-sm pointer-events-none select-none opacity-30">
        {children}
      </div>

      {/* 旋转提示遮罩 */}
      <div
        className={cn(
          "fixed inset-0 z-[9999]",
          "flex flex-col items-center justify-center",
          "bg-gradient-to-br from-[hsl(158,35%,18%)] via-[hsl(158,30%,14%)] to-[hsl(158,25%,10%)]",
          "text-white",
          "safe-area-inset",
        )}
      >
        {/* 装饰背景 */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          {/* 翡翠光晕 */}
          <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-[hsl(158,55%,42%)] opacity-10 rounded-full blur-3xl" />
          <div className="absolute bottom-1/4 right-1/4 w-48 h-48 bg-[hsl(42,95%,52%)] opacity-8 rounded-full blur-3xl" />

          {/* 装饰边框线 */}
          <div className="absolute inset-4 border border-[hsl(158,55%,42%)]/20 rounded-2xl" />
          <div className="absolute inset-8 border border-[hsl(42,95%,52%)]/10 rounded-xl" />
        </div>

        {/* 主内容 */}
        <div className="relative z-10 flex flex-col items-center px-8 py-12 max-w-sm text-center">
          {/* 旋转动画图标 */}
          <div className="relative mb-8">
            <div className={cn(
              "relative w-24 h-24 flex items-center justify-center",
              "rounded-2xl",
              "bg-gradient-to-br from-[hsl(158,55%,35%)] to-[hsl(158,55%,25%)]",
              "border border-[hsl(158,55%,42%)]/40",
              "shadow-[0_8px_32px_rgba(0,0,0,0.3),0_0_40px_hsla(158,55%,42%,0.2)]",
            )}>
              {/* 手机图标 */}
              <Smartphone className="w-10 h-10 text-[hsl(158,55%,65%)] animate-[rotate-phone_2s_ease-in-out_infinite]" />
            </div>

            {/* 旋转箭头 */}
            <div className="absolute -right-2 top-1/2 -translate-y-1/2">
              <RotateCcw className="w-8 h-8 text-[hsl(42,95%,52%)] animate-[pulse_1.5s_ease-in-out_infinite]" />
            </div>
          </div>

          {/* 标题 */}
          <h2
            className={cn(
              "text-2xl font-bold mb-3 font-display tracking-wide",
              "bg-gradient-to-r from-[hsl(42,95%,68%)] via-[hsl(42,95%,52%)] to-[hsl(42,95%,68%)]",
              "bg-clip-text text-transparent",
            )}
          >
            请旋转设备
          </h2>

          {/* 说明文字 */}
          <p className="text-[hsl(158,55%,65%)] text-sm leading-relaxed mb-6">
            为了获得最佳游戏体验，
            <br />
            请将设备横向放置
          </p>

          {/* 方向指示 */}
          <div className={cn(
            "flex items-center gap-3 px-5 py-3 rounded-full",
            "bg-black/30 border border-[hsl(158,55%,42%)]/30",
          )}>
            <div className="w-12 h-6 rounded border-2 border-[hsl(42,95%,52%)]/60 flex items-center justify-center">
              <div className="w-1 h-3 bg-[hsl(42,95%,52%)]/60 rounded-full" />
            </div>
            <span className="text-white/60 text-xs">横屏模式</span>
          </div>

          {/* 品牌标识 */}
          <div className="mt-10 flex items-center gap-2 text-white/30 text-xs">
            <span className="text-[hsl(158,55%,50%)]">♠</span>
            <span>掼蛋世界</span>
            <span className="text-[hsl(5,75%,50%)]">♥</span>
          </div>
        </div>
      </div>

      {/* 旋转动画样式 */}
      <style>{`
        @keyframes rotate-phone {
          0%, 100% {
            transform: rotate(0deg);
          }
          25% {
            transform: rotate(-15deg);
          }
          50% {
            transform: rotate(0deg);
          }
          75% {
            transform: rotate(90deg);
          }
        }
      `}</style>
    </>
  );
};

export default OrientationGuard;
