import React from 'react';
import { Card } from '@/components/ui';
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
    <Card
      variant="base"
      interactive={false}
      className={cn(
        'w-16 mobile-landscape:w-12 flex flex-col items-center rounded-md mobile-landscape:rounded overflow-hidden',
        'shadow-elevation-2 bg-gradient-to-b text-fg-inverse',
        isMyTeam
          ? 'from-[hsl(var(--primitive-primary-500))] to-[hsl(var(--primitive-primary-700))]'
          : 'from-[hsl(var(--primitive-danger-500))] to-[hsl(var(--primitive-danger-700))]'
      )}
    >
      <div className="w-full text-center py-1 mobile-landscape:py-0.5 text-sm mobile-landscape:text-[10px] font-semibold backdrop-brightness-95 bg-black/10">
        {label}
      </div>
      <div className="py-1.5 mobile-landscape:py-1 text-2xl mobile-landscape:text-lg font-extrabold tracking-wide drop-shadow-sm">
        {getLevelText(level)}
      </div>
      <div className="h-4 mobile-landscape:h-3 flex items-center justify-center">
        {isCurrentLevel && (
          <div className="w-0 h-0 border-l-[7px] border-r-[7px] border-b-[9px] mobile-landscape:border-l-[5px] mobile-landscape:border-r-[5px] mobile-landscape:border-b-[7px] border-l-transparent border-r-transparent border-b-fg-inverse drop-shadow" />
        )}
      </div>
    </Card>
  );

  return (
    <div className="absolute top-4 left-4 mobile-landscape:top-2 mobile-landscape:left-2 flex gap-2 mobile-landscape:gap-1">
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
