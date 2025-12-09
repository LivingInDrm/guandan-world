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
        <div className="h-8 w-8 rounded-full border-2 border-dashed border-table-300 flex items-center justify-center">
          <span className="text-table-300 text-xs">?</span>
        </div>
        <span className="text-xs text-table-300 mt-1 truncate w-full text-center">
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
        isCurrentUser && "ring-2 ring-primary ring-offset-1 rounded-lg"
      )}
      title={displayName}
    >
      <Avatar
        size="md"
        src={getAvatarUrl(player.avatar_key)}
        alt={displayName}
        fallback={displayName}
      />
      <span className={cn(
        "text-xs mt-1 truncate w-full text-center",
        player.online ? "text-table-400" : "text-table-300"
      )}>
        {displayName}
      </span>
    </div>
  );
};

export default PlayerMiniCard;
