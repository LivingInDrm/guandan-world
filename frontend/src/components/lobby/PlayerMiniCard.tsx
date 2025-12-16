import React from 'react';
import { Avatar } from '../ui';
import type { Player } from '../../types';
import { cn } from '@/lib/utils';
import { getAvatarUrl } from '../../utils/avatar';
import { UserPlus } from 'lucide-react';

interface PlayerMiniCardProps {
  player?: Player | null;
  isCurrentUser?: boolean;
}

const PlayerMiniCard: React.FC<PlayerMiniCardProps> = ({ player, isCurrentUser = false }) => {
  if (!player) {
    return (
      <div className="flex flex-col items-center w-14 mobile-landscape:w-11">
        <div className={cn(
          "h-10 w-10 mobile-landscape:h-8 mobile-landscape:w-8 rounded-full",
          "border-2 border-dashed border-stroke/40",
          "bg-gradient-to-br from-[hsl(40,8%,97%)] to-[hsl(40,8%,94%)]",
          "flex items-center justify-center",
          "transition-all duration-200",
          "hover:border-stroke/60 hover:bg-[hsl(40,8%,95%)]",
          "group",
        )}>
          <UserPlus className={cn(
            "w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3",
            "text-fg-secondary/40 group-hover:text-fg-secondary/60",
            "transition-colors duration-200",
          )} />
        </div>
        <span className="text-[10px] mobile-landscape:text-[8px] text-fg-secondary/50 mt-1.5 mobile-landscape:mt-1 truncate w-full text-center font-medium">
          空位
        </span>
      </div>
    );
  }

  const displayName = player.nickname || player.username;

  return (
    <div
      className={cn(
        "flex flex-col items-center w-14 mobile-landscape:w-11 p-1.5 mobile-landscape:p-1 rounded-xl mobile-landscape:rounded-lg",
        "transition-all duration-200",
        isCurrentUser ? [
          "bg-gradient-to-br from-state-active/10 to-state-active/20",
          "ring-2 ring-state-active/50 ring-offset-1 ring-offset-white",
          "shadow-[0_2px_8px_hsla(42,95%,52%,0.2)]",
        ] : [
          "hover:bg-[hsl(40,8%,97%)]",
        ],
      )}
      title={displayName}
    >
      <div className="relative">
        <Avatar
          size="sm"
          src={getAvatarUrl(player.avatar_key) ?? undefined}
          alt={displayName}
          fallback={displayName}
          className={cn(
            "mobile-landscape:w-8 mobile-landscape:h-8",
            "ring-2 ring-white shadow-[0_2px_6px_rgba(0,0,0,0.08)]",
            !player.online && "opacity-50 grayscale",
            "transition-all duration-200",
          )}
        />
        {/* 在线状态指示器 - 更精致的设计 */}
        <span className={cn(
          "absolute -bottom-0.5 -right-0.5",
          "w-3.5 h-3.5 mobile-landscape:w-2.5 mobile-landscape:h-2.5 rounded-full",
          "border-2 border-white",
          "transition-all duration-200",
          player.online
            ? "bg-[hsl(158,55%,45%)] shadow-[0_0_6px_hsla(158,55%,42%,0.5)]"
            : "bg-[hsl(35,6%,70%)]",
        )} />
      </div>
      <span className={cn(
        "text-[10px] mobile-landscape:text-[8px] mt-1.5 mobile-landscape:mt-1 truncate w-full text-center font-semibold",
        "transition-colors duration-200",
        player.online ? "text-fg-primary/80" : "text-fg-secondary/50",
        isCurrentUser && "text-state-active-muted",
      )}>
        {displayName}
      </span>
    </div>
  );
};

export default PlayerMiniCard;
