import React, { useMemo, type ReactNode } from 'react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import { Avatar } from '../ui';
import { getAvatarByUsername } from '../../utils/avatar';
import { cn } from '@/lib/utils';

export type { PlayerPosition };

export interface PlayerCardProps {
  player: Player | null;
  position: PlayerPosition;
  statusSlot?: ReactNode;
  isHighlighted?: boolean;
  className?: string;
}

const getPositionClasses = (position: PlayerPosition): string => {
  switch (position) {
    case 'bottom':
      return 'absolute bottom-4 left-1/2 transform -translate-x-1/2';
    case 'left':
      return 'absolute left-4 top-1/2 transform -translate-y-1/2';
    case 'top':
      return 'absolute top-4 left-1/2 transform -translate-x-1/2';
    case 'right':
      return 'absolute right-4 top-1/2 transform -translate-y-1/2';
    default:
      return '';
  }
};

const PlayerCard: React.FC<PlayerCardProps> = ({
  player,
  position,
  statusSlot,
  isHighlighted = false,
  className,
}) => {
  const avatarSrc = useMemo(() => {
    return player ? getAvatarByUsername(player.username) : null;
  }, [player]);

  if (!player) {
    return (
      <div className={cn(getPositionClasses(position), className)}>
        <div className="flex flex-col items-center p-1 rounded-xl bg-card/70 backdrop-blur-sm shadow-card w-fit">
          <div className="w-14 h-14 rounded-lg bg-muted shadow-card ring-1 ring-border" />
          <span className="mt-1 text-xs font-medium text-muted-foreground">空座位</span>
        </div>
      </div>
    );
  }

  return (
    <div className={cn(getPositionClasses(position), className)}>
      <div className="relative">
        <div className={cn(
          'flex flex-col items-center p-1 rounded-xl bg-card/70 backdrop-blur-sm shadow-card w-fit',
          isHighlighted && 'ring-2 ring-accent'
        )}>
          <Avatar
            src={avatarSrc}
            alt={player.username}
            fallback={player.username}
            size="2xl"
            className="rounded-lg shadow-card ring-1 ring-card/70"
          />
          <span className="mt-1 max-w-[80px] text-xs font-medium text-foreground truncate text-center">
            {player.username}
          </span>
        </div>
        {statusSlot && (
          <div className={cn(
            'absolute top-1/2 -translate-y-1/2',
            position === 'left' ? 'left-full ml-2' : 'right-full mr-2'
          )}>
            {statusSlot}
          </div>
        )}
      </div>
    </div>
  );
};

export default React.memo(PlayerCard);
