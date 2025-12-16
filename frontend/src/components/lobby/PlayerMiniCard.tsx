import React from 'react';
import { Avatar } from '../ui';
import type { Player } from '../../types';
import { cn } from '@/lib/utils';
import { getAvatarUrl } from '../../utils/avatar';

interface PlayerMiniCardProps {
  player?: Player | null;
  isCurrentUser?: boolean;
}

const PlayerMiniCard: React.FC<PlayerMiniCardProps> = ({ player, isCurrentUser = false }) => {
  if (!player) {
    return (
      <div className="flex flex-col items-center w-14">
        <div className={cn(
          "h-10 w-10 rounded-full",
          "border-2 border-dashed border-stroke/50",
          "bg-surface-elevated/30",
          "flex items-center justify-center",
          "transition-colors duration-fast",
        )}>
          <span className="text-fg-secondary/50 text-lg">+</span>
        </div>
        <span className="text-[10px] text-fg-secondary/60 mt-1.5 truncate w-full text-center">
          空位
        </span>
      </div>
    );
  }

  const displayName = player.nickname || player.username;

  return (
    <div
      className={cn(
        "flex flex-col items-center w-14 p-1 rounded-lg",
        "transition-all duration-fast",
        isCurrentUser && [
          "bg-state-active/10",
          "ring-2 ring-state-active/60 ring-offset-1 ring-offset-surface-base",
          "shadow-[0_0_8px_hsla(42,95%,52%,0.3)]",
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
            !player.online && "opacity-50 grayscale",
          )}
        />
        {/* 在线状态指示器 */}
        <span className={cn(
          "absolute -bottom-0.5 -right-0.5",
          "w-3 h-3 rounded-full border-2 border-surface-base",
          player.online
            ? "bg-[hsl(158,55%,42%)] shadow-[0_0_4px_hsla(158,55%,42%,0.6)]"
            : "bg-state-disabled",
        )} />
      </div>
      <span className={cn(
        "text-[10px] mt-1.5 truncate w-full text-center font-medium",
        player.online ? "text-fg-primary" : "text-fg-secondary/60",
        isCurrentUser && "text-state-active",
      )}>
        {displayName}
      </span>
    </div>
  );
};

export default PlayerMiniCard;
