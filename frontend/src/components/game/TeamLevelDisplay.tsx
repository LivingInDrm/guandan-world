import React from 'react';

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
    <div className={`w-16 flex flex-col items-center rounded-lg overflow-hidden shadow-md ${
      isMyTeam ? 'bg-gray-600' : 'bg-red-800'
    }`}>
      <div className={`w-full text-center py-1 text-white text-sm font-medium ${
        isMyTeam ? 'bg-gray-700' : 'bg-red-900'
      }`}>
        {label}
      </div>
      <div className="py-1.5 text-white text-2xl font-bold">
        {getLevelText(level)}
      </div>
      <div className="h-4 flex items-center justify-center">
        {isCurrentLevel && (
          <div className="w-0 h-0 border-l-[6px] border-r-[6px] border-b-[8px] border-l-transparent border-r-transparent border-b-white" />
        )}
      </div>
    </div>
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
