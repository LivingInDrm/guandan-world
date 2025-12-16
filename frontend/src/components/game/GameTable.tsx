import React, { type ReactNode } from 'react';
import type { Player } from '../../types';
import { cn } from '@/lib/utils';

export type PlayerPosition = 'bottom' | 'left' | 'top' | 'right';

export interface GameTableProps {
  players: (Player | null)[];
  currentPlayerSeat: number;
  renderPlayer: (player: Player | null, position: PlayerPosition, seat: number) => ReactNode;
  renderCenter: () => ReactNode;
  topLeftSlot?: ReactNode;
  className?: string;
}

const POSITIONS: PlayerPosition[] = ['bottom', 'left', 'top', 'right'];

const GameTable: React.FC<GameTableProps> = ({
  players,
  currentPlayerSeat,
  renderPlayer,
  renderCenter,
  topLeftSlot,
  className,
}) => {
  return (
    <div
      className={cn(
        'relative w-full rounded-2xl overflow-hidden',
        'h-[var(--table-height)]',
        'bg-gradient-table',
        'border-4 border-stroke-game mobile-landscape:border-2',
        'shadow-card-elevated',
        'safe-area-inset',
        className
      )}
    >
      {/* 牌桌纹理 */}
      <div
        className="absolute inset-0 pointer-events-none opacity-[0.03]"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
        }}
      />

      {/* 中心区域渐变光晕 */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: 'radial-gradient(ellipse 60% 50% at center, hsla(158,55%,45%,0.1) 0%, transparent 70%)',
        }}
      />

      {/* 边缘暗角 */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: 'radial-gradient(circle at center, transparent 50%, hsla(0,0%,0%,0.15) 100%)',
        }}
      />

      {/* 顶部高光 */}
      <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent" />

      {topLeftSlot}

      {/* 中央出牌区 - 向上偏移 */}
      <div
        className="absolute inset-0 flex items-center justify-center"
        style={{ transform: 'translateY(var(--play-area-offset-y, 0))' }}
      >
        <div
          className={cn(
            'relative rounded-xl',
            'w-[var(--table-center-width)] h-[var(--table-center-height)]',
            'p-[var(--table-padding)]',
            'bg-gradient-play-area',
            'shadow-[inset_0_2px_8px_hsla(0,0%,0%,0.25),inset_0_-1px_0_hsla(158,55%,50%,0.1)]',
            'border border-stroke-play-area',
          )}
        >
          {/* 出牌区内部装饰 */}
          <div className="absolute inset-2 mobile-landscape:inset-1 rounded-lg border border-dashed border-white/5" />
          {renderCenter()}
        </div>
      </div>

      {/* 角落装饰 */}
      <div className="absolute top-4 left-4 mobile-landscape:top-2 mobile-landscape:left-2 text-2xl mobile-landscape:text-base opacity-10 text-white">♠</div>
      <div className="absolute top-4 right-4 mobile-landscape:top-2 mobile-landscape:right-2 text-2xl mobile-landscape:text-base opacity-10 text-suit-red">♥</div>
      <div className="absolute bottom-4 left-4 mobile-landscape:bottom-2 mobile-landscape:left-2 text-2xl mobile-landscape:text-base opacity-10 text-white">♣</div>
      <div className="absolute bottom-4 right-4 mobile-landscape:bottom-2 mobile-landscape:right-2 text-2xl mobile-landscape:text-base opacity-10 text-suit-red">♦</div>

      {/* 玩家位置 */}
      {POSITIONS.map((position, relativeSeat) => {
        const seat = (currentPlayerSeat + relativeSeat) % 4;
        return (
          <React.Fragment key={seat}>
            {renderPlayer(players[seat], position, seat)}
          </React.Fragment>
        );
      })}
    </div>
  );
};

export default React.memo(GameTable);
