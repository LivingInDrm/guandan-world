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
  isOwner?: boolean;
  isTeammate?: boolean;
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
  isTeammate = false,
  className,
}) => {
  const avatarSrc = useMemo(() => {
    return player ? getAvatarByUsername(player.username) : null;
  }, [player]);

  if (!player) {
    return (
      <div className={cn(getPositionClasses(position), className)}>
        <div className="flex flex-col items-center p-1.5 rounded-xl bg-card/70 backdrop-blur-sm shadow-card w-fit">
          <div className="w-14 h-14 rounded-xl bg-muted shadow-card ring-2 ring-border/50" />
          <span className="mt-1 text-xs font-medium text-muted-foreground">空座位</span>
        </div>
      </div>
    );
  }

  const avatarRingClass = cn(
    'rounded-xl shadow-card ring-4 transition-all duration-300',
    isHighlighted ? 'ring-yellow-400 ring-glow-active' : (
      isTeammate ? 'ring-emerald-400' : 'ring-slate-400'
    )
  );

  return (
    <div className={cn(getPositionClasses(position), className)}>
      <div className="relative">
        {isOwner && (
          <div className="absolute -top-2 -right-2 z-10 bg-gradient-to-r from-yellow-400 to-yellow-500 text-white text-xs px-1.5 py-0.5 rounded-md font-bold shadow-md">
            房主
          </div>
        )}
        {isTeammate && !isOwner && (
          <div className="absolute -top-1 -right-1 z-10">
            <div className="w-5 h-5 rounded-full bg-gradient-to-br from-emerald-400 to-emerald-500 flex items-center justify-center shadow-md ring-2 ring-white/50">
              <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
        )}
        <div className={cn(
          'flex flex-col items-center p-1.5 rounded-xl bg-card/80 backdrop-blur-sm shadow-card w-fit',
          isHighlighted && 'bg-card/90'
        )}>
          <Avatar
            src={avatarSrc}
            alt={player.username}
            fallback={player.username}
            size="2xl"
            className={avatarRingClass}
          />
          <span className="mt-1.5 max-w-[80px] text-xs font-semibold text-foreground truncate text-center">
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
