import React from 'react';
import { cn } from '@/lib/utils';

export interface TeamLevelDisplayProps {
  teamLevels: [number, number];
  currentLevel: number;
  currentPlayerSeat: number;
}

export const getLevelText = (level: number): string => {
  if (level <= 10) return level.toString();
  switch (level) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    default: return level.toString();
  }
};

const TeamLevelDisplay: React.FC<TeamLevelDisplayProps> = ({
  teamLevels,
  currentLevel,
  currentPlayerSeat
}) => {
  const safeTeamLevels = teamLevels || [2, 2];
  const myTeamIndex = currentPlayerSeat % 2;
  const opponentTeamIndex = 1 - myTeamIndex;
  const myLevel = safeTeamLevels[myTeamIndex];
  const opponentLevel = safeTeamLevels[opponentTeamIndex];

  const LevelCard: React.FC<{
    label: string;
    level: number;
    isMyTeam: boolean;
    isCurrentLevel: boolean;
  }> = ({ label, level, isMyTeam, isCurrentLevel }) => (
    <div
      className={cn(
        "relative overflow-hidden",
        "w-14 mobile-landscape:w-11 flex flex-col items-center",
        "rounded-lg mobile-landscape:rounded-md",
        // 深色翡翠主题背景
        "bg-gradient-to-br from-[hsl(158,25%,14%)] via-[hsl(158,28%,11%)] to-[hsl(158,30%,8%)]",
        // 边框 - 我方翡翠绿，对方红宝石
        "border",
        isMyTeam
          ? "border-[hsl(158,55%,35%)]/40"
          : "border-[hsl(0,55%,45%)]/30",
        // 阴影
        "shadow-[0_4px_16px_rgba(0,0,0,0.4)]",
        // 当前等级高亮
        isCurrentLevel && (isMyTeam
          ? "ring-2 ring-[hsl(158,55%,45%)]/60 shadow-[0_0_20px_hsla(158,55%,40%,0.25)]"
          : "ring-2 ring-[hsl(0,55%,50%)]/50 shadow-[0_0_20px_hsla(0,55%,45%,0.2)]"
        ),
      )}
    >
      {/* 装饰性光效 */}
      <div className={cn(
        "absolute top-0 left-0 w-8 mobile-landscape:w-6 h-8 mobile-landscape:h-6 rounded-br-full pointer-events-none",
        isMyTeam
          ? "bg-gradient-to-br from-[hsl(158,55%,40%)]/15 to-transparent"
          : "bg-gradient-to-br from-[hsl(0,55%,50%)]/10 to-transparent"
      )} />

      {/* 标签 */}
      <div className={cn(
        "w-full text-center py-1 mobile-landscape:py-0.5",
        "text-[11px] mobile-landscape:text-[9px] font-semibold tracking-wider uppercase",
        "border-b",
        isMyTeam
          ? "bg-[hsl(158,55%,25%)]/30 border-[hsl(158,55%,40%)]/20 text-[hsl(158,55%,70%)]"
          : "bg-[hsl(0,40%,30%)]/20 border-[hsl(0,55%,50%)]/15 text-[hsl(0,55%,75%)]"
      )}>
        {label}
      </div>

      {/* 等级数字 */}
      <div className="py-2 mobile-landscape:py-1">
        <span className={cn(
          "text-2xl mobile-landscape:text-lg font-display font-black tracking-wide",
          isMyTeam
            ? "bg-gradient-to-b from-[hsl(158,55%,65%)] via-[hsl(158,55%,55%)] to-[hsl(158,55%,45%)] bg-clip-text text-transparent"
            : "bg-gradient-to-b from-[hsl(0,55%,70%)] via-[hsl(0,55%,60%)] to-[hsl(0,55%,50%)] bg-clip-text text-transparent"
        )}>
          {getLevelText(level)}
        </span>
      </div>

      {/* 当前等级指示器 */}
      <div className="h-3 mobile-landscape:h-2 flex items-center justify-center pb-0.5">
        {isCurrentLevel && (
          <div className={cn(
            "w-0 h-0",
            "border-l-[5px] border-r-[5px] border-b-[6px]",
            "mobile-landscape:border-l-[4px] mobile-landscape:border-r-[4px] mobile-landscape:border-b-[5px]",
            "border-l-transparent border-r-transparent",
            isMyTeam
              ? "border-b-[hsl(158,55%,55%)]"
              : "border-b-[hsl(0,55%,60%)]",
            "drop-shadow-sm"
          )} />
        )}
      </div>
    </div>
  );

  return (
    <div className="absolute top-3 md:top-4 left-3 md:left-4 mobile-landscape:top-2 mobile-landscape:left-2 flex gap-2 mobile-landscape:gap-1.5">
      <LevelCard
        label="我方"
        level={myLevel}
        isMyTeam={true}
        isCurrentLevel={myLevel === currentLevel}
      />
      <LevelCard
        label="对方"
        level={opponentLevel}
        isMyTeam={false}
        isCurrentLevel={opponentLevel === currentLevel}
      />
    </div>
  );
};

export default TeamLevelDisplay;
