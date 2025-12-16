import React, { useMemo, type ReactNode } from 'react';
import { Crown } from 'lucide-react';
import type { Player } from '../../types';
import type { PlayerPosition } from './GameTable';
import { Avatar, Card } from '../ui';
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

const getPositionStyles = (position: PlayerPosition): { className: string; style: React.CSSProperties } => {
  const offset = 'var(--player-card-offset)';
  const playAreaOffsetY = 'var(--play-area-offset-y, 0)';
  switch (position) {
    case 'bottom':
      return {
        className: 'absolute left-1/2 transform -translate-x-1/2',
        style: { bottom: offset },
      };
    case 'left':
      return {
        className: 'absolute',
        style: {
          left: offset,
          top: '50%',
          transform: `translateY(calc(-50% + ${playAreaOffsetY}))`,
        },
      };
    case 'top':
      return {
        className: 'absolute left-1/2 transform -translate-x-1/2',
        style: { top: offset },
      };
    case 'right':
      return {
        className: 'absolute',
        style: {
          right: offset,
          top: '50%',
          transform: `translateY(calc(-50% + ${playAreaOffsetY}))`,
        },
      };
    default:
      return { className: '', style: {} };
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

  const positionStyles = getPositionStyles(position);

  // 空座位
  if (!player) {
    return (
      <div className={cn(positionStyles.className, className)} style={positionStyles.style}>
        <Card
          variant="emphasis"
          interactive={false}
          className={cn(
            "flex flex-col items-center p-2 mobile-landscape:p-1 w-fit",
            "bg-black/30 backdrop-blur-sm",
            "border border-white/15",
          )}
        >
          <div
            className={cn(
              "rounded-full",
              "bg-white/5 border-2 border-dashed border-white/20",
              "flex items-center justify-center",
              "!w-[var(--player-avatar-size)] !h-[var(--player-avatar-size)]",
            )}
          >
            <span className="text-white/30 text-2xl mobile-landscape:text-lg">+</span>
          </div>
          <span className="mt-2 mobile-landscape:mt-0.5 text-xs mobile-landscape:text-[10px] font-medium text-white/40">
            空座位
          </span>
        </Card>
      </div>
    );
  }

  return (
    <div className={cn(positionStyles.className, className)} style={positionStyles.style}>
      <div className="relative">
        {/* 房主徽章 */}
        {isOwner && (
          <div className={cn(
            "absolute -top-3 -right-3 mobile-landscape:-top-2 mobile-landscape:-right-2 z-10",
            "w-7 h-7 mobile-landscape:w-5 mobile-landscape:h-5 rounded-full",
            "bg-gradient-to-br from-state-active to-[hsl(42,95%,45%)]",
            "flex items-center justify-center",
            "shadow-[0_2px_8px_hsla(42,95%,52%,0.5)]",
            "border-2 mobile-landscape:border border-white/30",
          )}>
            <Crown className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3 text-white" />
          </div>
        )}

        {/* 玩家卡片 */}
        <Card
          variant="emphasis"
          interactive={false}
          className={cn(
            "flex flex-col items-center p-2 mobile-landscape:p-1 w-fit",
            // 玻璃态背景
            "bg-black/30 backdrop-blur-sm",
            "border border-white/15",
            // 高亮状态
            isHighlighted && [
              "ring-2 mobile-landscape:ring-1 ring-state-active/60",
              "shadow-[0_0_20px_hsla(42,95%,52%,0.3)]",
              "bg-black/40",
            ],
          )}
        >
          <Avatar
            src={avatarSrc ?? undefined}
            alt={player.username}
            fallback={player.username}
            size="2xl"
            ringState={isHighlighted ? 'active' : 'none'}
            className={cn(
              "shadow-[var(--elevation-2)]",
              isHighlighted && "shadow-[0_0_12px_hsla(42,95%,52%,0.4)]",
              // 使用 CSS 变量控制响应式尺寸
              "!w-[var(--player-avatar-size)] !h-[var(--player-avatar-size)]",
            )}
          />
          <span className={cn(
            "mt-2 mobile-landscape:mt-0.5 max-w-[80px] mobile-landscape:max-w-[60px] text-xs mobile-landscape:text-[10px] font-medium truncate text-center",
            isHighlighted ? "text-state-active" : "text-white/90",
          )}>
            {player.nickname || player.username}
          </span>
        </Card>

        {/* 状态插槽 */}
        {statusSlot && (
          <div className={cn(
            'absolute top-1/2 -translate-y-1/2',
            position === 'left' ? 'left-full ml-3 mobile-landscape:ml-2' : 'right-full mr-3 mobile-landscape:mr-2'
          )}>
            {statusSlot}
          </div>
        )}
      </div>
    </div>
  );
};

export default React.memo(PlayerCard);
