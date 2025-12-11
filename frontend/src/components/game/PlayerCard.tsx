import React, { useMemo, type ReactNode } from 'react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import { Avatar, Card, Badge } from '../ui-next';
import { getAvatarByUsername } from '../../utils/avatar';
import { cn } from '@/lib/utils';

export type { PlayerPosition };

export interface PlayerCardProps {
  player: Player | null;
  position: PlayerPosition;
  statusSlot?: ReactNode;
  isHighlighted?: boolean;
  isOwner?: boolean;
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
  isOwner = false,
  className,
}) => {
  const avatarSrc = useMemo(() => {
    return player ? getAvatarByUsername(player.username) : null;
  }, [player]);

  if (!player) {
    return (
      <div className={cn(getPositionClasses(position), className)}>
        <Card variant="base" interactive={false} className="flex flex-col items-center p-1 w-fit">
          <div className="w-24 h-24 rounded-ds-lg bg-ds-surface-elevated shadow-ds-elevation-1 ring-1 ring-ds-border" />
          <span className="mt-1 text-xs font-medium text-ds-text-secondary">空座位</span>
        </Card>
      </div>
    );
  }

  return (
    <div className={cn(getPositionClasses(position), className)}>
      <div className="relative">
        {isOwner && (
          <Badge variant="owner" size="sm" className="absolute -top-2 -right-2 z-10">
            房主
          </Badge>
        )}
        <Card variant="emphasis" interactive={false} className="flex flex-col items-center p-1 w-fit">
          <Avatar
            src={avatarSrc ?? undefined}
            alt={player.username}
            fallback={player.username}
            size="2xl"
            ringState={isHighlighted ? 'active' : 'none'}
          />
          <span className="mt-1 max-w-[80px] text-xs font-medium text-ds-text-primary truncate text-center">
            {player.username}
          </span>
        </Card>
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
