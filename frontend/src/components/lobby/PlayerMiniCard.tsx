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
        <div className="h-10 w-10 rounded-full border-2 border-dashed border-stroke flex items-center justify-center">
          <span className="text-fg-secondary text-xs">?</span>
        </div>
        <span className="text-xs text-fg-secondary mt-1 truncate w-full text-center">
          空位
        </span>
      </div>
    );
  }

  const displayName = player.nickname || player.username;

  return (
    <div
      className={cn(
        "flex flex-col items-center w-14",
        isCurrentUser && "ring-2 ring-state-active ring-offset-1 rounded-lg"
      )}
      title={displayName}
    >
      <Avatar
        size="sm"
        src={getAvatarUrl(player.avatar_key) ?? undefined}
        alt={displayName}
        fallback={displayName}
      />
      <span className={cn(
        "text-xs mt-1 truncate w-full text-center",
        player.online ? "text-fg-primary" : "text-fg-secondary"
      )}>
        {displayName}
      </span>
    </div>
  );
};

export default PlayerMiniCard;
