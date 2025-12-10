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
        'relative w-full h-[30rem] rounded-lg overflow-hidden',
        'bg-gradient-to-b from-table-50 via-table-100 to-table-200',
        className
      )}
    >
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: 'radial-gradient(circle at center, transparent 60%, rgba(0,0,0,0.06) 100%)',
        }}
      />

      {topLeftSlot}

      <div className="absolute inset-0 flex items-center justify-center">
        <div
          className={cn(
            'relative w-[420px] h-[240px] rounded-[20px] p-6',
            'bg-table-300/20 backdrop-blur-sm',
            'border border-table-400/35',
            'shadow-inset-soft'
          )}
        >
          {renderCenter()}
        </div>
      </div>

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
