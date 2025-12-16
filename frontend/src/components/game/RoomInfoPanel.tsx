import React, { useState, useCallback } from 'react';
import { Copy, Check, Users } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface RoomInfoPanelProps {
  roomCode: string;
}

const RoomInfoPanel: React.FC<RoomInfoPanelProps> = ({ roomCode }) => {
  const [copied, setCopied] = useState(false);
  const [copyError, setCopyError] = useState(false);

  const inviteLink = `${window.location.origin}/join?code=${roomCode}`;

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(inviteLink);
      setCopied(true);
      setCopyError(false);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
      setCopyError(true);
      setTimeout(() => setCopyError(false), 3000);
    }
  }, [inviteLink]);

  return (
    <div className="absolute top-3 md:top-4 left-3 md:left-4 mobile-landscape:top-2 mobile-landscape:left-2 z-10">
      <div
        className={cn(
          "relative overflow-hidden",
          "rounded-xl mobile-landscape:rounded-lg",
          // 深色翡翠主题背景
          "bg-gradient-to-br from-[hsl(158,25%,14%)] via-[hsl(158,28%,11%)] to-[hsl(158,30%,8%)]",
          // 边框与阴影
          "border border-[hsl(158,55%,30%)]/25",
          "shadow-[0_4px_20px_rgba(0,0,0,0.4),0_0_30px_hsla(158,55%,30%,0.08)]",
        )}
      >
        {/* 装饰性光效 */}
        <div className="absolute top-0 left-0 w-16 mobile-landscape:w-10 h-16 mobile-landscape:h-10 bg-gradient-to-br from-[hsl(158,55%,40%)]/15 to-transparent rounded-br-full pointer-events-none" />

        {/* 顶部标签栏 */}
        <div className={cn(
          "flex items-center justify-center gap-1.5 px-4 mobile-landscape:px-2.5 py-1.5 mobile-landscape:py-1",
          "bg-gradient-to-r from-[hsl(158,55%,25%)]/40 via-[hsl(158,55%,30%)]/30 to-[hsl(158,55%,25%)]/40",
          "border-b border-[hsl(158,55%,40%)]/20",
        )}>
          <Users className="w-3.5 mobile-landscape:w-3 h-3.5 mobile-landscape:h-3 text-[hsl(158,55%,55%)]" />
          <span className="text-xs mobile-landscape:text-[10px] font-semibold tracking-wider text-[hsl(158,55%,70%)] uppercase">
            房间码
          </span>
        </div>

        {/* 房间码显示 */}
        <div className="px-5 mobile-landscape:px-3 py-2.5 mobile-landscape:py-1.5">
          <div className={cn(
            "text-3xl mobile-landscape:text-xl font-display font-black tracking-[0.25em] mobile-landscape:tracking-[0.15em] pl-[0.25em] mobile-landscape:pl-[0.15em] text-center",
            "bg-gradient-to-b from-state-active via-[hsl(42,95%,58%)] to-[hsl(42,95%,42%)]",
            "bg-clip-text text-transparent",
            "drop-shadow-[0_0_12px_hsla(42,95%,52%,0.4)]",
          )}>
            {roomCode}
          </div>
        </div>

        {/* 邀请按钮 */}
        <div className="px-3 mobile-landscape:px-2 pb-3 mobile-landscape:pb-2">
          <button
            onClick={handleCopy}
            className={cn(
              "w-full flex items-center justify-center gap-2 mobile-landscape:gap-1.5",
              "px-4 mobile-landscape:px-2.5 py-2 mobile-landscape:py-1 rounded-lg mobile-landscape:rounded-md",
              "text-sm mobile-landscape:text-xs font-semibold tracking-wide",
              "transition-all duration-200",
              copied
                ? "bg-[hsl(158,55%,35%)] text-white border border-[hsl(158,55%,50%)]/50"
                : "bg-white/10 hover:bg-white/15 text-white/90 hover:text-white border border-white/20 hover:border-white/30",
              "shadow-[0_2px_8px_rgba(0,0,0,0.2),inset_0_1px_0_rgba(255,255,255,0.05)]",
            )}
          >
            {copied ? (
              <>
                <Check className="w-4 mobile-landscape:w-3 h-4 mobile-landscape:h-3" />
                <span>已复制</span>
              </>
            ) : (
              <>
                <Copy className="w-4 mobile-landscape:w-3 h-4 mobile-landscape:h-3" />
                <span className="mobile-landscape:hidden">复制邀请链接</span>
                <span className="hidden mobile-landscape:inline">邀请</span>
              </>
            )}
          </button>
          {copyError && (
            <p className="text-xs mobile-landscape:text-[10px] text-center mt-1.5 text-red-400/80">
              复制失败
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default RoomInfoPanel;
