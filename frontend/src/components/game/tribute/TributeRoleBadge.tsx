import React from 'react';
import { cn } from '@/lib/utils';
import { ArrowUp, ArrowDown, Check, Loader2 } from 'lucide-react';

export interface TributeRoleBadgeProps {
  role: 'giver' | 'receiver' | null;
  isSubmitted?: boolean;
  isReceived?: boolean;
  isCurrentSelector?: boolean;
}

const TributeRoleBadge: React.FC<TributeRoleBadgeProps> = ({
  role,
  isSubmitted = false,
  isReceived = false,
  isCurrentSelector = false,
}) => {
  if (!role) {
    return null;
  }

  // 上贡方 (Giver) - 红宝石色调
  if (role === 'giver') {
    const isCompleted = isSubmitted;

    return (
      <div
        className={cn(
          "flex items-center gap-1 mobile-landscape:gap-0.5",
          "px-2 mobile-landscape:px-1.5 py-1 mobile-landscape:py-0.5",
          "rounded-md mobile-landscape:rounded",
          "text-xs mobile-landscape:text-[10px] font-semibold tracking-wide whitespace-nowrap",
          "border",
          "transition-all duration-300",
          isCompleted
            ? [
                "bg-[hsl(158,55%,25%)]/30",
                "border-[hsl(158,55%,40%)]/40",
                "text-[hsl(158,55%,65%)]",
              ]
            : [
                "bg-[hsl(0,55%,30%)]/25",
                "border-[hsl(0,55%,50%)]/35",
                "text-[hsl(0,55%,70%)]",
                "shadow-[0_0_8px_hsla(0,55%,50%,0.15)]",
              ]
        )}
      >
        {isCompleted ? (
          <Check className="w-3 mobile-landscape:w-2.5 h-3 mobile-landscape:h-2.5" />
        ) : (
          <ArrowUp className="w-3 mobile-landscape:w-2.5 h-3 mobile-landscape:h-2.5" />
        )}
        <span>{isCompleted ? '已上贡' : '待上贡'}</span>
      </div>
    );
  }

  // 收贡方 (Receiver) - 翡翠/金色色调
  if (role === 'receiver') {
    const isCompleted = isReceived;
    const isSelecting = isCurrentSelector && !isCompleted;

    const getText = () => {
      if (isCompleted) return '已收贡';
      if (isSelecting) return '选牌中';
      return '待收贡';
    };

    return (
      <div
        className={cn(
          "flex items-center gap-1 mobile-landscape:gap-0.5",
          "px-2 mobile-landscape:px-1.5 py-1 mobile-landscape:py-0.5",
          "rounded-md mobile-landscape:rounded",
          "text-xs mobile-landscape:text-[10px] font-semibold tracking-wide whitespace-nowrap",
          "border",
          "transition-all duration-300",
          isCompleted
            ? [
                "bg-[hsl(158,55%,25%)]/30",
                "border-[hsl(158,55%,40%)]/40",
                "text-[hsl(158,55%,65%)]",
              ]
            : isSelecting
              ? [
                  "bg-state-active/20",
                  "border-state-active/50",
                  "text-state-active",
                  "shadow-[0_0_12px_hsla(42,95%,52%,0.25)]",
                  "animate-pulse",
                ]
              : [
                  "bg-[hsl(158,55%,30%)]/25",
                  "border-[hsl(158,55%,45%)]/35",
                  "text-[hsl(158,55%,65%)]",
                  "shadow-[0_0_8px_hsla(158,55%,45%,0.15)]",
                ]
        )}
      >
        {isCompleted ? (
          <Check className="w-3 mobile-landscape:w-2.5 h-3 mobile-landscape:h-2.5" />
        ) : isSelecting ? (
          <Loader2 className="w-3 mobile-landscape:w-2.5 h-3 mobile-landscape:h-2.5 animate-spin" />
        ) : (
          <ArrowDown className="w-3 mobile-landscape:w-2.5 h-3 mobile-landscape:h-2.5" />
        )}
        <span>{getText()}</span>
      </div>
    );
  }

  return null;
};

export default React.memo(TributeRoleBadge);
