import React from 'react';
import { Card } from '@/components/ui-next';
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
        'w-16 flex flex-col items-center rounded-ds-md overflow-hidden',
        'shadow-ds-elevation-2 bg-gradient-to-b text-ds-text-inverse',
        isMyTeam
          ? 'from-[hsl(var(--ds-primitive-primary-500))] to-[hsl(var(--ds-primitive-primary-700))]'
          : 'from-[hsl(var(--ds-primitive-danger-500))] to-[hsl(var(--ds-primitive-danger-700))]'
      )}
    >
      <div className="w-full text-center py-1 text-sm font-semibold backdrop-brightness-95 bg-black/10">
        {label}
      </div>
      <div className="py-1.5 text-2xl font-extrabold tracking-wide drop-shadow-sm">
        {getLevelText(level)}
      </div>
      <div className="h-4 flex items-center justify-center">
        {isCurrentLevel && (
          <div className="w-0 h-0 border-l-[7px] border-r-[7px] border-b-[9px] border-l-transparent border-r-transparent border-b-ds-text-inverse drop-shadow" />
        )}
      </div>
    </Card>
  );

  return (
    <div className="absolute top-4 left-4 flex gap-2">
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
